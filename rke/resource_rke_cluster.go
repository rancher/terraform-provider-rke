package rke

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/rke/cluster"
	"github.com/rancher/rke/cmd"
	"github.com/rancher/rke/dind"
	"github.com/rancher/rke/hosts"
	"github.com/rancher/rke/pki"
	v3 "github.com/rancher/types/apis/management.cattle.io/v3"
	//log "github.com/sirupsen/logrus"
)

const rkeClusterDINDWaitTime = 3

func resourceRKECluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceRKEClusterCreate,
		Read:   resourceRKEClusterRead,
		Update: resourceRKEClusterUpdate,
		Delete: resourceRKEClusterDelete,
		Schema: rkeClusterFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

func resourceRKEClusterCreate(d *schema.ResourceData, meta interface{}) error {
	if delay, ok := d.Get("delay_on_creation").(int); ok && delay > 0 {
		time.Sleep(time.Duration(delay) * time.Second)
	}
	return resourceRKEClusterUpdate(d, meta)
}

func resourceRKEClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := clusterUp(d); err != nil {
		return wrapErrWithRKEOutputs(err)
	}
	return wrapErrWithRKEOutputs(resourceRKEClusterRead(d, meta))
}

func resourceRKEClusterRead(d *schema.ResourceData, meta interface{}) error {
	currentCluster, err := readClusterState(d)
	if err != nil {
		return wrapErrWithRKEOutputs(err)
	}
	return wrapErrWithRKEOutputs(flattenRKECluster(d, currentCluster))
}

func resourceRKEClusterDelete(d *schema.ResourceData, meta interface{}) error {
	err := clusterDelete(d)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func clusterUp(d *schema.ResourceData) error {
	rkeConfig, _, clusterFilePath, tempDir, err := getRKEClusterConfig(d)
	defer removeTempDir(tempDir)
	if err != nil {
		return err
	}

	// setting up the flags and dialers
	flags := expandRKEClusterFlag(d, clusterFilePath)
	dialers := hosts.DialersOptions{}

	// setting dind if needed
	if d.Get("dind").(bool) {
		dindStorageDriver := d.Get("dind_storage_driver").(string)
		dindDNS := d.Get("dind_dns_server").(string)
		if err = prepareDINDEnv(context.Background(), rkeConfig, dindStorageDriver, dindDNS); err != nil {
			return fmt.Errorf("Failed preparing DIND environment err:%v", err)
		}
		dialers = hosts.GetDialerOptions(hosts.DindConnFactory, hosts.DindHealthcheckConnFactory, nil)
	}

	if err := cmd.ClusterInit(context.Background(), rkeConfig, dialers, flags); err != nil {
		return fmt.Errorf("Failed initializing cluster err:%v", err)
	}
	_, _, _, _, _, clusterUpErr := cmd.ClusterUp(context.Background(), dialers, flags, map[string]interface{}{})

	// set cluster state to resourceData
	err = setRKEClusterState(d, tempDir)
	if clusterUpErr != nil {
		return fmt.Errorf("Failed running cluster err:%v", clusterUpErr)
	}

	return err
}

func prepareDINDEnv(ctx context.Context, rkeConfig *v3.RancherKubernetesEngineConfig, dindStorageDriver, dindDNS string) error {
	os.Setenv("DOCKER_API_VERSION", hosts.DockerAPIVersion)
	for i := range rkeConfig.Nodes {
		address, err := dind.StartUpDindContainer(ctx, rkeConfig.Nodes[i].Address, dind.DINDNetwork, dindStorageDriver, dindDNS)
		if err != nil {
			return err
		}
		if rkeConfig.Nodes[i].HostnameOverride == "" {
			rkeConfig.Nodes[i].HostnameOverride = rkeConfig.Nodes[i].Address
		}
		rkeConfig.Nodes[i].Address = address
	}
	time.Sleep(rkeClusterDINDWaitTime * time.Second)
	return nil
}

func clusterDelete(d *schema.ResourceData) error {
	rkeConfig, _, clusterFilePath, tempDir, err := getRKEClusterConfig(d)
	defer removeTempDir(tempDir)
	if err != nil {
		return err
	}

	if d.Get("dind").(bool) {
		os.Setenv("DOCKER_API_VERSION", hosts.DockerAPIVersion)
		for _, node := range rkeConfig.Nodes {
			if err = dind.RmoveDindContainer(context.Background(), node.Address); err != nil {
				return nil
			}
		}
		return nil
	}

	// setting up the flags
	flags := cluster.GetExternalFlags(false, false, false, "", clusterFilePath)

	// Omiting ClusterRemove  errors
	cmd.ClusterRemove(context.Background(), rkeConfig, hosts.DialersOptions{}, flags)

	return nil
}

func getRKEClusterConfig(d *schema.ResourceData) (*v3.RancherKubernetesEngineConfig, string, string, string, error) {
	rkeClusterYaml, err := expandRKECluster(d)
	if err != nil {
		return nil, "", "", "", err
	}

	rkeConfig, err := cluster.ParseConfig(rkeClusterYaml)
	if err != nil {
		return nil, "", "", "", fmt.Errorf("Failed to parse cluster config: %v\n%s", err, rkeClusterYaml)
	}

	d.Set("rke_cluster_yaml", rkeClusterYaml)

	clusterFilePath, tempDir, err := writeRKEConfigFiles(d)
	if err != nil {
		return nil, "", "", "", err
	}
	return rkeConfig, rkeClusterYaml, clusterFilePath, tempDir, err

}

func setRKEClusterState(d *schema.ResourceData, configDir string) error {
	rkeState, err := readRKEStateFile(configDir)
	if err != nil {
		return err
	}
	if rkeState != "" {
		d.Set("rke_state", rkeState) // nolint
	}

	kubeConfig, err := readKubeConfig(configDir)
	if err != nil {
		return err
	}
	if kubeConfig != "" {
		d.Set("kube_config_yaml", kubeConfig)          // nolint
		d.Set("internal_kube_config_yaml", kubeConfig) // nolint
	}

	if len(d.Id()) == 0 {
		d.SetId(getNewUUID())
	}
	return nil
}

func readClusterState(d *schema.ResourceData) (*cluster.Cluster, error) {
	_, _, clusterFilePath, tempDir, err := getRKEClusterConfig(d)
	defer removeTempDir(tempDir)
	if err != nil {
		return nil, err
	}

	// setting up the flags
	flags := expandRKEClusterFlag(d, clusterFilePath)
	_, readedCluster, err := getClusterState(context.Background(), hosts.DialersOptions{}, flags)
	if err != nil {
		switch err.(type) {
		case *stateNotFoundError:
			d.SetId("")
			return nil, nil
		}
	}

	readedCluster.DinD = d.Get("dind").(bool)

	return readedCluster, err
}

func getClusterState(ctx context.Context, dialersOptions hosts.DialersOptions, flags cluster.ExternalFlags) (*cluster.FullState, *cluster.Cluster, error) {
	fullState, err := cluster.ReadStateFile(ctx, cluster.GetStateFilePath(flags.ClusterFilePath, flags.ConfigDir))
	if err != nil {
		return nil, nil, newStateNotFoundError(err)
	}

	kubeCluster, err := cluster.InitClusterObject(ctx, fullState.DesiredState.RancherKubernetesEngineConfig.DeepCopy(), flags, "")
	if err != nil {
		return nil, nil, err
	}

	err = kubeCluster.SetupDialers(ctx, hosts.DialersOptions{})
	if err != nil {
		return nil, nil, err
	}

	if fullState.CurrentState.RancherKubernetesEngineConfig == nil && fullState.DesiredState.RancherKubernetesEngineConfig != nil {
		fullState.CurrentState = fullState.DesiredState
	}

	clusterState, err := kubeCluster.GetClusterState(ctx, fullState)
	if err != nil {
		return nil, nil, err
	}

	return fullState, clusterState, nil
}

func readKubeConfig(dir string) (string, error) {
	configPath := filepath.Join(dir, pki.ClusterConfig)
	localKubeConfigPath := pki.GetLocalKubeConfig(configPath, "")
	if _, err := os.Stat(localKubeConfigPath); err == nil {
		var data []byte
		if data, err = ioutil.ReadFile(localKubeConfigPath); err != nil {
			return "", err
		}
		return string(data), nil
	}
	return "", nil
}

func readRKEStateFile(dir string) (string, error) {
	configPath := filepath.Join(dir, pki.ClusterConfig)
	stateFilePath := cluster.GetStateFilePath(configPath, "")
	if _, err := os.Stat(stateFilePath); err == nil {
		var data []byte
		if data, err = ioutil.ReadFile(stateFilePath); err != nil {
			return "", err
		}
		return string(data), nil
	}
	return "", nil
}

func writeRKEConfigFiles(d *schema.ResourceData) (string, string, error) {
	tempDir, err := createTempDir()
	if err != nil {
		return "", "", err
	}
	clusterFilePath := filepath.Join(tempDir, pki.ClusterConfig)
	if err = writeRKEConfig(clusterFilePath, d); err != nil {
		return "", tempDir, err
	}

	if err = writeKubeConfig(clusterFilePath, d); err != nil {
		return "", tempDir, err
	}
	if err = writeRKEState(clusterFilePath, d); err != nil {
		return "", tempDir, err
	}

	return clusterFilePath, tempDir, err
}

func writeRKEState(dir string, d *schema.ResourceData) error {
	if strState, ok := d.Get("rke_state").(string); ok && len(strState) > 0 {
		stateFilePath := cluster.GetStateFilePath(dir, "")
		return ioutil.WriteFile(stateFilePath, []byte(strState), 0640)
	}
	return nil
}

func writeKubeConfig(dir string, d *schema.ResourceData) error {
	if strConf, ok := d.Get("internal_kube_config_yaml").(string); ok && len(strConf) > 0 {
		localKubeConfigPath := pki.GetLocalKubeConfig(dir, "")
		return ioutil.WriteFile(localKubeConfigPath, []byte(strConf), 0640)
	}
	return nil
}

func writeRKEConfig(configFile string, d *schema.ResourceData) error {
	if strConf, ok := d.Get("rke_cluster_yaml").(string); ok && len(strConf) > 0 {
		return ioutil.WriteFile(configFile, []byte(strConf), 0640)
	}
	return nil

}

func createTempDir() (string, error) {
	// create tmp dir for configDir
	workDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	tempDir, err := ioutil.TempDir(workDir, "terraform-provider-rke-tmp-")
	if err != nil {
		return "", err
	}
	return tempDir, nil
}

func removeTempDir(tempDir string) {
	if len(tempDir) > 0 {
		os.RemoveAll(tempDir)
	}
}

type stateNotFoundError struct {
	actual error
}

func newStateNotFoundError(actual error) *stateNotFoundError {
	return &stateNotFoundError{actual: actual}
}

func (n *stateNotFoundError) Error() string {
	return n.actual.Error()
}
