package rke

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-getter/helper/url"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rancher/rke/cluster"
	"github.com/rancher/rke/cmd"
	"github.com/rancher/rke/hosts"
	"github.com/rancher/rke/log"
	"github.com/rancher/rke/pki"
	v3 "github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func resourceRKECluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceRKEClusterCreate,
		Read:   resourceRKEClusterRead,
		Update: resourceRKEClusterUpdate,
		Delete: resourceRKEClusterDelete,
		CustomizeDiff: func(d *schema.ResourceDiff, i interface{}) error {
			if isRKEConfigChanged(d) {
				computedFields := []string{
					"kube_config_yaml",
					"rke_cluster_yaml",
				}
				for _, key := range computedFields {
					d.SetNewComputed(key)
				}
			}
			return nil
		},
		Schema: ClusterSchema(),
	}
}

func resourceRKEClusterCreate(d *schema.ResourceData, meta interface{}) error {
	if delay, ok := d.GetOk("delay_on_creation"); ok && delay.(int) > 0 {
		time.Sleep(time.Duration(delay.(int)) * time.Second)
	}

	if err := clusterUp(d); err != nil {
		return wrapErrWithRKEOutputs(err)
	}
	return wrapErrWithRKEOutputs(resourceRKEClusterRead(d, meta))
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
		if _, ok := err.(*nodeUnreachableError); ok {
			d.SetId("")
			return nil
		}
		return wrapErrWithRKEOutputs(err)
	}
	return wrapErrWithRKEOutputs(clusterToState(currentCluster, d))
}

func resourceRKEClusterDelete(d *schema.ResourceData, meta interface{}) error {
	if err := clusterRemove(d); err != nil {
		if _, ok := err.(*nodeUnreachableError); !ok {
			return wrapErrWithRKEOutputs(err)
		}
	}
	d.SetId("")
	return nil
}

func clusterUp(d *schema.ResourceData) error {
	rkeConfig, parseErr := parseResourceRKEConfig(d)
	if parseErr != nil {
		return parseErr
	}
	disablePortCheck := d.Get("disable_port_check").(bool)

	clusterFilePath, tempDir, err := prepareTempRKEConfigFiles(rkeConfig, d)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir) // nolint

	// setting up the flags
	flags := cluster.GetExternalFlags(false, false, disablePortCheck, "", clusterFilePath)
	if err := cmd.ClusterInit(context.Background(), rkeConfig, hosts.DialersOptions{}, flags); err != nil {
		return err
	}

	apiURL, caCrt, clientCert, clientKey, _, clusterUpErr := cmd.ClusterUp(context.Background(), hosts.DialersOptions{}, flags)
	if clusterUpErr != nil {
		return clusterUpErr
	}

	// set keys to resourceData
	return setRKEClusterKeys(d, apiURL, caCrt, clientCert, clientKey, tempDir, rkeConfig)
}

func clusterRemove(d *schema.ResourceData) error {
	rkeConfig, parseErr := parseResourceRKEConfig(d)
	if parseErr != nil {
		return parseErr
	}

	clusterFilePath, tempDir, err := prepareTempRKEConfigFiles(rkeConfig, d)
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir) // nolint

	// setting up the flags
	flags := cluster.GetExternalFlags(false, false, false, "", clusterFilePath)

	return realClusterRemove(context.Background(), rkeConfig, hosts.DialersOptions{}, flags)
}

func realClusterRemove(
	ctx context.Context,
	rkeConfig *v3.RancherKubernetesEngineConfig,
	dialersOptions hosts.DialersOptions,
	flags cluster.ExternalFlags) error {

	log.Infof(ctx, "Tearing down Kubernetes cluster")
	kubeCluster, err := cluster.InitClusterObject(ctx, rkeConfig, flags)
	if err != nil {
		return err
	}
	if err := kubeCluster.SetupDialers(ctx, dialersOptions); err != nil {
		return err
	}

	err = kubeCluster.TunnelHosts(ctx, flags)
	if err != nil {
		return newNodeUnreachableError(err)
	}

	logrus.Debugf("Starting Cluster removal")
	err = kubeCluster.ClusterRemove(ctx)
	if err != nil {
		return err
	}

	log.Infof(ctx, "Cluster removed successfully")
	return nil
}

func setRKEClusterKeys(d *schema.ResourceData, apiURL, caCrt, clientCert, clientKey string, configDir string, rkeConfig *v3.RancherKubernetesEngineConfig) error {

	parsedURL, err := url.Parse(apiURL)
	if err != nil {
		return err
	}
	d.Set("ca_crt", caCrt)           // nolint
	d.Set("client_cert", clientCert) // nolint
	d.Set("client_key", clientKey)   // nolint

	rkeState, err := readRKEStateFile(configDir)
	if err != nil {
		return err
	}
	if rkeState != "" {
		d.Set("rke_state", rkeState)
	}

	kubeConfig, err := readKubeConfig(configDir)
	if err != nil {
		return err
	}
	if kubeConfig != "" {
		d.Set("kube_config_yaml", kubeConfig)          // nolint
		d.Set("internal_kube_config_yaml", kubeConfig) // nolint
	}

	yamlRkeConfig, err := yaml.Marshal(*rkeConfig)
	if err != nil {
		return err
	}
	d.Set("rke_cluster_yaml", string(yamlRkeConfig)) // nolint

	d.SetId(parsedURL.Hostname())
	return nil
}

func readClusterState(d *schema.ResourceData) (*cluster.Cluster, error) {
	apiURL := fmt.Sprintf("https://%s:6443", d.Id())
	caCrt := d.Get("ca_crt").(string)
	clientCert := d.Get("client_cert").(string)
	clientKey := d.Get("client_key").(string)

	requiredValues := []string{apiURL, caCrt, clientCert, clientKey}
	for _, v := range requiredValues {
		if v == "" {
			d.SetId("")
			return nil, nil
		}
	}

	rkeConfig, parseErr := parseResourceRKEConfig(d)
	if parseErr != nil {
		return nil, parseErr
	}

	yamlRkeConfig, err := yaml.Marshal(*rkeConfig)
	if err != nil {
		return nil, err
	}
	d.Set("rke_cluster_yaml", string(yamlRkeConfig)) // nolint

	clusterFilePath, tempDir, err := prepareTempRKEConfigFiles(rkeConfig, d)
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir) // nolint

	// setting up the flags
	flags := cluster.GetExternalFlags(false, false, d.Get("disable_port_check").(bool), "", clusterFilePath)
	fullState, readedCluster, err := realClusterRead(context.Background(), hosts.DialersOptions{}, flags)
	if err != nil {
		switch err.(type) {
		case *stateNotFoundError, *nodeUnreachableError:
			d.SetId("")
			return nil, nil
		}
	}

	kubeConfig, err := readKubeConfig(tempDir)
	if err != nil {
		return nil, err
	}
	if kubeConfig != "" {
		d.Set("kube_config_yaml", kubeConfig)          // nolint
		d.Set("internal_kube_config_yaml", kubeConfig) // nolint
	}

	strRKEState, err := json.MarshalIndent(fullState, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("Failed to Marshal state object: %v", err)
	}
	d.Set("rke_state", strRKEState) // nolint

	return readedCluster, err
}

func realClusterRead(ctx context.Context, dialersOptions hosts.DialersOptions, flags cluster.ExternalFlags) (*cluster.FullState, *cluster.Cluster, error) {

	fullState, err := cluster.ReadStateFile(ctx, cluster.GetStateFilePath(flags.ClusterFilePath, flags.ConfigDir))
	if err != nil {
		return nil, nil, newStateNotFoundError(err)
	}

	kubeCluster, err := cluster.InitClusterObject(ctx, fullState.DesiredState.RancherKubernetesEngineConfig.DeepCopy(), flags)
	if err != nil {
		return nil, nil, err
	}
	err = kubeCluster.SetupDialers(ctx, hosts.DialersOptions{})
	if err != nil {
		return nil, nil, err
	}

	err = kubeCluster.TunnelHosts(ctx, flags)
	if err != nil {
		return nil, nil, newNodeUnreachableError(err)
	}

	clusterState, err := kubeCluster.GetClusterState(ctx, fullState)
	if err != nil {
		return nil, nil, err
	}
	return fullState, clusterState, nil
}

func prepareTempRKEConfigFiles(rkeConfig *v3.RancherKubernetesEngineConfig, d resourceData) (string, string, error) {
	tempDir, tempDirErr := createTempDir()
	if tempDirErr != nil {
		return "", "", tempDirErr
	}
	if err := writeKubeConfigFile(tempDir, d); err != nil {
		return "", "", err
	}
	if err := writeRKEStateFile(tempDir, d); err != nil {
		return "", "", err
	}

	clusterFilePath := filepath.Join(tempDir, pki.ClusterConfig)
	if err := writeClusterConfig(rkeConfig, clusterFilePath); err != nil {
		return "", "", err
	}

	return clusterFilePath, tempDir, nil
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

func writeRKEStateFile(dir string, d resourceData) error {
	if rawRKEState, ok := d.GetOk("rke_state"); ok {
		strState := rawRKEState.(string)
		if strState != "" {
			configPath := filepath.Join(dir, pki.ClusterConfig)
			stateFilePath := cluster.GetStateFilePath(configPath, "")
			if err := ioutil.WriteFile(stateFilePath, []byte(strState), 0640); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeKubeConfigFile(dir string, d resourceData) error {
	if rawKubeConfig, ok := d.GetOk("internal_kube_config_yaml"); ok {
		strConf := rawKubeConfig.(string)
		if strConf != "" {
			configPath := filepath.Join(dir, pki.ClusterConfig)
			localKubeConfigPath := pki.GetLocalKubeConfig(configPath, "")
			if err := ioutil.WriteFile(localKubeConfigPath, []byte(strConf), 0640); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeClusterConfig(cluster *v3.RancherKubernetesEngineConfig, configFile string) error {
	yamlConfig, err := yaml.Marshal(*cluster)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFile, []byte(string(yamlConfig)), 0640)

}

func createTempDir() (string, error) {
	// create tmp dir for configDir
	var workDir, tempDir string
	var err error
	if workDir, err = os.Getwd(); err != nil {
		return "", err
	}
	if tempDir, err = ioutil.TempDir(workDir, "terraform-provider-rke-"); err != nil {
		return "", err
	}
	return tempDir, nil
}

type resourceDiffer interface {
	HasChange(string) bool
}

func isRKEConfigChanged(d resourceDiffer) bool {
	targetKeys := []string{
		"nodes_conf",
		"nodes",
		"services_etcd",
		"services_kube_api",
		"services_kube_controller",
		"services_scheduler",
		"services_kubelet",
		"services_kubeproxy",
		"network",
		"authentication",
		"addons",
		"addons_include",
		"addon_job_timeout",
		"system_images",
		"ssh_key_path",
		"ssh_agent_auth",
		"bastion_host",
		"monitoring",
		"authorization",
		"ignore_docker_version",
		"kubernetes_version",
		"private_registries",
		"ingress",
		"cluster_name",
		"cloud_provider",
		"prefix_path",
	}
	for _, key := range targetKeys {
		if d.HasChange(key) {
			return true
		}
	}
	return false
}

type nodeUnreachableError struct {
	actual error
}

func newNodeUnreachableError(actual error) *nodeUnreachableError {
	return &nodeUnreachableError{actual: actual}
}

func (n *nodeUnreachableError) Error() string {
	return n.actual.Error()
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
