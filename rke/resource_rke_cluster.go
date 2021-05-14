package rke

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/rke/cluster"
	"github.com/rancher/rke/cmd"
	"github.com/rancher/rke/dind"
	"github.com/rancher/rke/hosts"
	"github.com/rancher/rke/pki"
	v3 "github.com/rancher/rke/types"
	log "github.com/sirupsen/logrus"
)

const rkeClusterDINDWaitTime = 3

func resourceRKECluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceRKEClusterCreate,
		Read:   resourceRKEClusterRead,
		Update: resourceRKEClusterUpdate,
		Delete: resourceRKEClusterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRKEClusterImport,
		},
		Schema: rkeClusterFields(),
		CustomizeDiff: func(d *schema.ResourceDiff, i interface{}) error {
			if changedKeys := getChangedKeys(d); len(changedKeys) > 0 {
				log.Infof("[rke_provider] rke cluster changed arguments: %v", changedKeys)
				if log.IsLevelEnabled(log.DebugLevel) {
					for k := range changedKeys {
						old, new := d.GetChange(k)
						log.Debugf("[rke_provider] %s values old: %v new: %v", k, old, new)
					}
				}
				computedFields := []string{
					"rke_state",
					"kube_config_yaml",
					"rke_cluster_yaml",
				}

				if changedKeys["rotate_certificates"] || changedKeys["cluster_yaml"] {
					for _, key := range []string{"ca_crt", "client_cert", "client_key", "certificates", "kube_admin_user"} {
						computedFields = append(computedFields, key)
					}
				}

				if changedKeys["dns"] || changedKeys["services"] || changedKeys["cluster_yaml"] {
					for _, key := range []string{"cluster_domain", "cluster_cidr", "cluster_dns_server"} {
						computedFields = append(computedFields, key)
					}
				}

				if changedKeys["nodes"] || changedKeys["cluster_yaml"] {
					for _, key := range []string{"api_server_url", "etcd_hosts", "control_plane_hosts", "inactive_hosts", "worker_hosts"} {
						computedFields = append(computedFields, key)
					}
				}

				if changedKeys["kubernetes_version"] || changedKeys["system_images"] || changedKeys["cluster_yaml"] {
					computedFields = append(computedFields, "running_system_images")
				}

				for _, key := range computedFields {
					if err := d.SetNewComputed(key); err != nil {
						return err
					}
				}
			}
			return nil
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

func resourceRKEClusterCreate(d *schema.ResourceData, meta interface{}) error {
	log.Info("Creating RKE cluster...")
	if delay, ok := d.Get("delay_on_creation").(int); ok && delay > 0 {
		time.Sleep(time.Duration(delay) * time.Second)
	}
	if err := clusterUp(d); err != nil {
		return meta.(*Config).saveRKEOutput(err)
	}
	return resourceRKEClusterRead(d, meta)
}

func resourceRKEClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Info("Updating RKE cluster...")

	restored, err := clusterRestore(d)
	if err != nil {
		return meta.(*Config).saveRKEOutput(err)
	}
	if !restored {
		if err := clusterUp(d); err != nil {
			return meta.(*Config).saveRKEOutput(err)
		}
	}
	return resourceRKEClusterRead(d, meta)
}

func resourceRKEClusterRead(d *schema.ResourceData, meta interface{}) error {
	log.Infof("Reading RKE cluster %s ...", d.Id())
	currentCluster, err := readClusterState(d)
	if err != nil {
		return meta.(*Config).saveRKEOutput(err)
	}

	return meta.(*Config).saveRKEOutput(flattenRKECluster(d, currentCluster))
}

func resourceRKEClusterDelete(d *schema.ResourceData, meta interface{}) error {
	log.Info("Deleting RKE cluster...")
	err := clusterDelete(d)
	if err != nil {
		return meta.(*Config).saveRKEOutput(err)
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

	// setting up the flags, dialers and context
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
	// set init cluster state to resourceData
	flattenRKEClusterFlag(d, &flags)
	err = setRKEClusterState(d, tempDir)
	if err != nil {
		return fmt.Errorf("Failed setting initial cluster state err:%v", err)
	}

	_, _, _, _, _, clusterUpErr := cmd.ClusterUp(context.Background(), dialers, flags, map[string]interface{}{})

	// set cluster state to resourceData
	err = setRKEClusterState(d, tempDir)
	if clusterUpErr != nil {
		return fmt.Errorf("Failed running cluster err:%v", clusterUpErr)
	}
	if err != nil {
		return fmt.Errorf("Failed setting cluster state err:%v", err)
	}

	return nil
}

func clusterRestore(d *schema.ResourceData) (bool, error) {
	rkeConfig, _, clusterFilePath, tempDir, err := getRKEClusterConfig(d)
	defer removeTempDir(tempDir)
	if err != nil {
		return false, err
	}

	if !rkeConfig.Restore.Restore {
		return false, nil
	}
	if len(rkeConfig.Restore.SnapshotName) == 0 {
		return false, fmt.Errorf("Failed restoring cluster: snapshop_name must be provided")
	}

	// setting up the flags, dialers and context
	flags := expandRKEClusterFlag(d, clusterFilePath)
	dialers := hosts.DialersOptions{}

	// set restore to false to force diff on next apply
	rkeConfig.Restore.Restore = false
	_, _, _, _, _, clusterRestoreErr := cmd.RestoreEtcdSnapshot(context.Background(), rkeConfig, dialers, flags, map[string]interface{}{}, rkeConfig.Restore.SnapshotName)

	// set cluster state to resourceData
	flattenRKEClusterFlag(d, &flags)
	err = setRKEClusterState(d, tempDir)
	if clusterRestoreErr != nil {
		return false, fmt.Errorf("Failed restoring cluster err:%v", clusterRestoreErr)
	}
	if err != nil {
		return false, fmt.Errorf("Failed setting cluster state err:%v", err)
	}

	return true, nil
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
	flags := cluster.GetExternalFlags(false, false, false, false, "", clusterFilePath)

	// Omiting ClusterRemove  errors
	cmd.ClusterRemove(context.Background(), rkeConfig, hosts.DialersOptions{}, flags)

	return nil
}

func getRKEClusterConfig(d *schema.ResourceData) (*v3.RancherKubernetesEngineConfig, string, string, string, error) {
	rkeClusterYaml, _, err := expandRKECluster(d)
	if err != nil {
		return nil, "", "", "", err
	}

	rkeConfig, err := cluster.ParseConfig(rkeClusterYaml)
	if err != nil {
		return nil, "", "", "", fmt.Errorf("Failed to parse cluster config: %v\n%s", err, rkeClusterYaml)
	}

	if rkeConfig.Services.KubeAPI.EventRateLimit != nil && rkeConfig.Services.KubeAPI.EventRateLimit.Configuration != nil {
		if len(rkeConfig.Services.KubeAPI.EventRateLimit.Configuration.TypeMeta.Kind) == 0 {
			rkeConfig.Services.KubeAPI.EventRateLimit.Configuration.TypeMeta.Kind = clusterServicesKubeAPIEventRateLimitConfigKindDefault
		}
		if len(rkeConfig.Services.KubeAPI.EventRateLimit.Configuration.TypeMeta.APIVersion) == 0 {
			rkeConfig.Services.KubeAPI.EventRateLimit.Configuration.TypeMeta.APIVersion = clusterServicesKubeAPIEventRateLimitConfigAPIDefault
		}
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

	readedCluster.DinD = false
	if v, ok := d.Get("dind").(bool); ok && v {
		readedCluster.DinD = v
	}

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
	if strConf, ok := d.Get("kube_config_yaml").(string); ok && len(strConf) > 0 {
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
	tempDir, err := ioutil.TempDir(workDir, "rke-provider-tmp-")
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

func getChangedKeys(d *schema.ResourceDiff) map[string]bool {
	targetKeys := []string{
		"addon_job_timeout",
		"addons",
		"addons_include",
		"authentication",
		"authorization",
		"bastion_host",
		"cloud_provider",
		"cluster_name",
		"cluster_yaml",
		"dns",
		"ingress",
		"kubernetes_version",
		"monitoring",
		"network",
		"nodes",
		"prefix_path",
		"private_registries",
		"restore",
		"rotate_certificates",
		"services",
		"ssh_agent_auth",
		"ssh_cert_path",
		"ssh_key_path",
		"system_images",
	}
	changedKeys := map[string]bool{}
	for _, key := range targetKeys {
		if v := d.HasChange(key); v {
			changedKeys[key] = v
		}
		if key == "services" && changedKeys[key] {
			old, new := d.GetChange("services")
			oldstr, _ := expandRKEClusterServices(old.([]interface{}))
			newstr, _ := expandRKEClusterServices(new.([]interface{}))
			if reflect.DeepEqual(oldstr, newstr) {
				delete(changedKeys, key)
			}
		}
		if key == "nodes" && changedKeys[key] {
			old, new := d.GetChange("nodes")
			oldInput := old.([]interface{})
			oldInputLen := len(oldInput)
			newInput := new.([]interface{})
			newInputLen := len(newInput)
			// Indexing old and new input by Address
			oldInputIndexAddress := map[string]int{}
			for i := range oldInput {
				if row, ok := oldInput[i].(map[string]interface{}); ok {
					if v, ok := row["address"].(string); ok {
						oldInputIndexAddress[v] = i
					}
				}
			}
			// Sorting new input
			sortedNewInput := make([]interface{}, len(newInput))
			newRows := []interface{}{}
			lastIndex := 0
			for i := range newInput {
				if row, ok := newInput[i].(map[string]interface{}); ok {
					if address, ok := row["address"].(string); ok {
						if v, ok := oldInputIndexAddress[address]; ok {
							if v > i && oldInputLen > newInputLen {
								v = v - (v - i)
							}
							sortedNewInput[v] = row
							lastIndex++
							continue
						}
					}
					newRows = append(newRows, row)
				}
			}
			for i := range newRows {
				sortedNewInput[lastIndex+i] = newRows[i]
			}
			if reflect.DeepEqual(oldInput, sortedNewInput) {
				delete(changedKeys, key)
			}
		}
	}
	return changedKeys
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
