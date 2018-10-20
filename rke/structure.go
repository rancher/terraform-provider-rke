package rke

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/validation"
	"github.com/rancher/rke/cluster"
	"github.com/rancher/rke/hosts"
	"github.com/rancher/rke/pki"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"gopkg.in/yaml.v2"
)

type resourceData interface {
	GetOk(key string) (interface{}, bool)
}

type stateBuilder interface {
	Set(string, interface{}) error
	SetId(string)
}

var rkeConfigBuilders = []func(rkeConfig *v3.RancherKubernetesEngineConfig, d resourceData) error{
	setNodesFromResource,
	setServicesFromResource,
	setNetworkFromResource,
	setAuthenticationFromResource,
	setAddonsFromResource,
	setSystemImagesFromResource,
	setSSHSettingsFromResource,
	setBastionHostFromResource,
	setMonitoringFromResource,
	setAuthorizationFromResource,
	setMiscConfigFromResource,
	setCloudProviderFromResource,
}

func parseResourceRKEConfig(d resourceData) (*v3.RancherKubernetesEngineConfig, error) {
	var err error
	rkeConfig := &v3.RancherKubernetesEngineConfig{}

	for _, builder := range rkeConfigBuilders {
		if err = builder(rkeConfig, d); err != nil {
			return nil, err
		}
	}
	return rkeConfig, nil
}

func setNodesFromResource(rkeConfig *v3.RancherKubernetesEngineConfig, d resourceData) error {

	rawNodes, _ := d.GetOk("nodes")
	rawNodesConf, _ := d.GetOk("nodes_conf")

	nodesOK := rawNodes != nil && len(rawNodes.([]interface{})) > 0
	nodesConfOK := rawNodesConf != nil && len(rawNodesConf.([]interface{})) > 0

	var nodes []v3.RKEConfigNode
	var err error

	switch {
	case nodesOK && nodesConfOK:
		return fmt.Errorf("cannot specify both %q and %q", "nodes", "nodes_conf")
	case nodesOK:
		if nodes, err = parseResourceRKEConfigNodes(d); err != nil {
			return err
		}
	case nodesConfOK:
		if nodes, err = parseResourceRKEConfigNodesConf(d); err != nil {
			return err
		}
	default:
		return fmt.Errorf("either %q or %q is required", "nodes", "nodes_conf")
	}

	if len(nodes) > 0 {
		rkeConfig.Nodes = nodes
	}
	return nil
}

func setServicesFromResource(rkeConfig *v3.RancherKubernetesEngineConfig, d resourceData) error {
	var err error
	var etcd *v3.ETCDService
	if etcd, err = parseResourceETCDService(d); err != nil {
		return err
	}
	if etcd != nil {
		rkeConfig.Services.Etcd = *etcd
	}

	var kubeAPI *v3.KubeAPIService
	if kubeAPI, err = parseResourceKubeAPIService(d); err != nil {
		return err
	}
	if kubeAPI != nil {
		rkeConfig.Services.KubeAPI = *kubeAPI
	}

	var kubeController *v3.KubeControllerService
	if kubeController, err = parseResourceKubeControllerService(d); err != nil {
		return err
	}
	if kubeController != nil {
		rkeConfig.Services.KubeController = *kubeController
	}

	var scheduler *v3.SchedulerService
	if scheduler, err = parseResourceSchedulerService(d); err != nil {
		return err
	}
	if scheduler != nil {
		rkeConfig.Services.Scheduler = *scheduler
	}

	var kubelet *v3.KubeletService
	if kubelet, err = parseResourceKubeletService(d); err != nil {
		return err
	}
	if kubelet != nil {
		rkeConfig.Services.Kubelet = *kubelet
	}

	var kubeproxy *v3.KubeproxyService
	if kubeproxy, err = parseResourceKubeproxyService(d); err != nil {
		return err
	}
	if kubeproxy != nil {
		rkeConfig.Services.Kubeproxy = *kubeproxy
	}
	return nil
}

func setNetworkFromResource(rkeConfig *v3.RancherKubernetesEngineConfig, d resourceData) error {
	var err error
	var network *v3.NetworkConfig
	if network, err = parseResourceNetwork(d); err != nil {
		return err
	}
	if network != nil {
		rkeConfig.Network = *network
	}
	return nil
}

func setAuthenticationFromResource(rkeConfig *v3.RancherKubernetesEngineConfig, d resourceData) error {
	var err error
	var authn *v3.AuthnConfig
	if authn, err = parseResourceAuthentication(d); err != nil {
		return err
	}
	if authn != nil {
		rkeConfig.Authentication = *authn
	}
	return nil
}

func setAddonsFromResource(rkeConfig *v3.RancherKubernetesEngineConfig, d resourceData) error {
	var err error
	var addons string
	if addons, err = parseResourceAddons(d); err != nil {
		return err
	}
	rkeConfig.Addons = addons

	var addonsInclude []string
	if addonsInclude, err = parseResourceAddonsInclude(d); err != nil {
		return err
	}
	rkeConfig.AddonsInclude = addonsInclude

	var addonJobTimeout int
	if addonJobTimeout, err = parseResourceAddonJobTimeout(d); err != nil {
		return err
	}
	rkeConfig.AddonJobTimeout = addonJobTimeout

	return nil
}

func setSystemImagesFromResource(rkeConfig *v3.RancherKubernetesEngineConfig, d resourceData) error {
	var err error
	var systemImages *v3.RKESystemImages
	if systemImages, err = parseResourceSystemImages(d); err != nil {
		return err
	}
	if systemImages != nil {
		rkeConfig.SystemImages = *systemImages
	}
	return nil
}

func setSSHSettingsFromResource(rkeConfig *v3.RancherKubernetesEngineConfig, d resourceData) error {
	var err error
	var sshKeyPath string
	if sshKeyPath, err = parseResourceSSHKeyPath(d); err != nil {
		return err
	}
	rkeConfig.SSHKeyPath = sshKeyPath

	var sshAgentAuth bool
	if sshAgentAuth, err = parseResourceSSHAgentAuth(d); err != nil {
		return err
	}
	rkeConfig.SSHAgentAuth = sshAgentAuth

	return nil
}

func setBastionHostFromResource(rkeConfig *v3.RancherKubernetesEngineConfig, d resourceData) error {
	var err error
	var host *v3.BastionHost
	if host, err = parseResourceBastionHost(d); err != nil {
		return err
	}
	if host != nil {
		rkeConfig.BastionHost = *host
	}
	return nil
}

func setMonitoringFromResource(rkeConfig *v3.RancherKubernetesEngineConfig, d resourceData) error {
	var err error
	var monitoring *v3.MonitoringConfig
	if monitoring, err = parseResourceMonitoring(d); err != nil {
		return err
	}
	if monitoring != nil {
		rkeConfig.Monitoring = *monitoring
	}
	return nil
}

func setAuthorizationFromResource(rkeConfig *v3.RancherKubernetesEngineConfig, d resourceData) error {
	var err error
	var authz *v3.AuthzConfig
	if authz, err = parseResourceAuthorization(d); err != nil {
		return err
	}
	if authz != nil {
		rkeConfig.Authorization = *authz
	}
	return nil
}

func setMiscConfigFromResource(rkeConfig *v3.RancherKubernetesEngineConfig, d resourceData) error {
	var err error
	var ignoreDockerVersion bool
	if ignoreDockerVersion, err = parseResourceIgnoreDockerVersion(d); err != nil {
		return err
	}
	rkeConfig.IgnoreDockerVersion = ignoreDockerVersion

	var kubernetesVersion string
	if kubernetesVersion, err = parseResourceVersion(d); err != nil {
		return err
	}
	rkeConfig.Version = kubernetesVersion

	var registries []v3.PrivateRegistry
	registries, err = parseResourcePrivateRegistries(d)
	if err != nil {
		return err
	}
	if len(registries) > 0 {
		rkeConfig.PrivateRegistries = registries
	}

	var ingress *v3.IngressConfig
	if ingress, err = parseResourceIngress(d); err != nil {
		return err
	}
	if ingress != nil {
		rkeConfig.Ingress = *ingress
	}

	var clusterName string
	if clusterName, err = parseResourceClusterName(d); err != nil {
		return err
	}
	rkeConfig.ClusterName = clusterName

	var prefixPath string
	if prefixPath, err = parseResourcePrefixPath(d); err != nil {
		return err
	}
	rkeConfig.PrefixPath = prefixPath

	return nil
}

func setCloudProviderFromResource(rkeConfig *v3.RancherKubernetesEngineConfig, d resourceData) error {
	var err error
	var cloudProvider *v3.CloudProvider
	if cloudProvider, err = parseResourceCloudProvider(d); err != nil {
		return err
	}
	if cloudProvider != nil {
		rkeConfig.CloudProvider = *cloudProvider
	}
	return nil
}

func parseResourceRKEConfigNodes(d resourceData) ([]v3.RKEConfigNode, error) {
	nodes := []v3.RKEConfigNode{}
	if rawNodes, ok := d.GetOk("nodes"); ok {

		nodeList := rawNodes.([]interface{})
		for _, rawNode := range nodeList {
			nodeValues := rawNode.(map[string]interface{})

			node, err := parseResourceRKEConfigNode(nodeValues)
			if err != nil {
				return nil, err
			}

			nodes = append(nodes, node)
		}
	}
	return nodes, nil
}

func parseResourceRKEConfigNode(nodeValues map[string]interface{}) (v3.RKEConfigNode, error) {
	node := v3.RKEConfigNode{}

	applyMapToObj(&mapObjMapping{
		source: nodeValues,
		stringMapping: map[string]*string{
			"node_name":         &node.NodeName,
			"address":           &node.Address,
			"internal_address":  &node.InternalAddress,
			"hostname_override": &node.HostnameOverride,
			"user":              &node.User,
			"docker_socket":     &node.DockerSocket,
			"ssh_key":           &node.SSHKey,
			"ssh_key_path":      &node.SSHKeyPath,
		},
		boolMapping: map[string]*bool{
			"ssh_agent_auth": &node.SSHAgentAuth,
		},
		listStrMapping: map[string]*[]string{
			"role": &node.Role,
		},
		mapStrMapping: map[string]*map[string]string{
			"labels": &node.Labels,
		},
	})

	if v, ok := nodeValues["port"]; ok {
		p := v.(int)
		if p > 0 {
			node.Port = fmt.Sprintf("%d", p)
		}
	}

	// validate role and roles
	roleValidateFunc := validation.StringInSlice([]string{"controlplane", "etcd", "worker"}, false)
	rawRole, hasRole := nodeValues["role"]
	rawRoles, hasRoles := nodeValues["roles"]

	if !hasRole && !hasRoles {
		return node, fmt.Errorf("either role or roles is required")
	}
	if hasRole && hasRoles {
		if len(rawRole.([]interface{})) > 0 && len(rawRoles.(string)) > 0 {
			return node, fmt.Errorf("cannot specify both role and roles for a node")
		}
	}

	if hasRole && len(rawRole.([]interface{})) > 0 {
		roles := []string{}
		for _, e := range rawRole.([]interface{}) {
			strRole := e.(string)
			if _, errs := roleValidateFunc(strRole, "role"); len(errs) > 0 {
				return node, errs[0]
			}
			roles = append(roles, strRole)
		}
		node.Role = roles
	}

	if hasRoles && len(rawRoles.(string)) > 0 {
		roles := []string{}
		for _, e := range strings.Split(rawRoles.(string), ",") {
			strRole := strings.TrimSpace(e)
			if _, errs := roleValidateFunc(strRole, "role"); len(errs) > 0 {
				return node, errs[0]
			}
			roles = append(roles, strRole)
		}
		node.Role = roles
	}

	if v, ok := nodeValues["ssh_agent_auth"]; ok {
		node.SSHAgentAuth = v.(bool)
	}

	if v, ok := nodeValues["labels"]; ok {
		nodeLabels := map[string]string{}
		labels := v.(map[string]interface{})
		for k, v := range labels {
			if v, ok := v.(string); ok {
				nodeLabels[k] = v
			}
		}
		node.Labels = nodeLabels
	}

	return node, nil
}

func parseResourceRKEConfigNodesConf(d resourceData) ([]v3.RKEConfigNode, error) {
	nodes := []v3.RKEConfigNode{}
	if rawNodes, ok := d.GetOk("nodes_conf"); ok {

		nodeList := rawNodes.([]interface{})
		for _, rawNode := range nodeList {

			var node v3.RKEConfigNode
			var err error

			nodeConf := []byte(rawNode.(string))

			if err = json.Unmarshal(nodeConf, &node); err != nil {
				err = yaml.Unmarshal(nodeConf, &node)
			}

			if err != nil {
				return nil, err
			}

			nodes = append(nodes, node)
		}
	}
	return nodes, nil
}

func parseResourceETCDService(d resourceData) (*v3.ETCDService, error) {
	if rawList, ok := d.GetOk("services_etcd"); ok {
		if rawServices, ok := rawList.([]interface{}); ok && len(rawServices) > 0 {
			rawService := rawServices[0]
			etcd := &v3.ETCDService{}
			rawMap := rawService.(map[string]interface{})

			applyMapToObj(&mapObjMapping{
				source: rawMap,
				stringMapping: map[string]*string{
					//"image":     &etcd.Image,
					"ca_cert":   &etcd.CACert,
					"cert":      &etcd.Cert,
					"key":       &etcd.Key,
					"path":      &etcd.Path,
					"retention": &etcd.Retention,
					"creation":  &etcd.Creation,
				},
				boolMapping: map[string]*bool{
					"snapshot": &etcd.Snapshot,
				},
				mapStrMapping: map[string]*map[string]string{
					"extra_args": &etcd.ExtraArgs,
				},
				listStrMapping: map[string]*[]string{
					"extra_binds":   &etcd.ExtraBinds,
					"extra_env":     &etcd.ExtraEnv,
					"external_urls": &etcd.ExternalURLs,
				},
			})

			return etcd, nil
		}
	}
	return nil, nil
}

func parseResourceKubeAPIService(d resourceData) (*v3.KubeAPIService, error) {
	if rawList, ok := d.GetOk("services_kube_api"); ok {

		if rawServices, ok := rawList.([]interface{}); ok && len(rawServices) > 0 {
			rawService := rawServices[0]
			kubeAPI := &v3.KubeAPIService{}
			rawMap := rawService.(map[string]interface{})

			applyMapToObj(&mapObjMapping{
				source: rawMap,
				stringMapping: map[string]*string{
					//"image":                    &kubeAPI.Image,
					"service_cluster_ip_range": &kubeAPI.ServiceClusterIPRange,
					"service_node_port_range":  &kubeAPI.ServiceNodePortRange,
				},
				boolMapping: map[string]*bool{
					"pod_security_policy": &kubeAPI.PodSecurityPolicy,
				},
				mapStrMapping: map[string]*map[string]string{
					"extra_args": &kubeAPI.ExtraArgs,
				},
				listStrMapping: map[string]*[]string{
					"extra_binds": &kubeAPI.ExtraBinds,
					"extra_env":   &kubeAPI.ExtraEnv,
				},
			})

			return kubeAPI, nil
		}
	}
	return nil, nil
}

func parseResourceKubeControllerService(d resourceData) (*v3.KubeControllerService, error) {
	if rawList, ok := d.GetOk("services_kube_controller"); ok {
		if rawServices, ok := rawList.([]interface{}); ok && len(rawServices) > 0 {
			rawService := rawServices[0]
			kubeController := &v3.KubeControllerService{}
			rawMap := rawService.(map[string]interface{})

			applyMapToObj(&mapObjMapping{
				source: rawMap,
				stringMapping: map[string]*string{
					//"image":                    &kubeController.Image,
					"cluster_cidr":             &kubeController.ClusterCIDR,
					"service_cluster_ip_range": &kubeController.ServiceClusterIPRange,
				},
				mapStrMapping: map[string]*map[string]string{
					"extra_args": &kubeController.ExtraArgs,
				},
				listStrMapping: map[string]*[]string{
					"extra_binds": &kubeController.ExtraBinds,
					"extra_env":   &kubeController.ExtraEnv,
				},
			})

			return kubeController, nil
		}
	}
	return nil, nil
}

func parseResourceSchedulerService(d resourceData) (*v3.SchedulerService, error) {
	if rawList, ok := d.GetOk("services_scheduler"); ok {
		if rawServices, ok := rawList.([]interface{}); ok && len(rawServices) > 0 {
			rawService := rawServices[0]
			scheduler := &v3.SchedulerService{}
			rawMap := rawService.(map[string]interface{})

			applyMapToObj(&mapObjMapping{
				source:        rawMap,
				stringMapping: map[string]*string{
					//"image": &scheduler.Image,
				},
				mapStrMapping: map[string]*map[string]string{
					"extra_args": &scheduler.ExtraArgs,
				},
				listStrMapping: map[string]*[]string{
					"extra_binds": &scheduler.ExtraBinds,
					"extra_env":   &scheduler.ExtraEnv,
				},
			})

			return scheduler, nil
		}
	}
	return nil, nil
}

func parseResourceKubeletService(d resourceData) (*v3.KubeletService, error) {
	if rawList, ok := d.GetOk("services_kubelet"); ok {
		if rawServices, ok := rawList.([]interface{}); ok && len(rawServices) > 0 {
			rawService := rawServices[0]
			kubelet := &v3.KubeletService{}
			rawMap := rawService.(map[string]interface{})

			applyMapToObj(&mapObjMapping{
				source: rawMap,
				stringMapping: map[string]*string{
					//"image":                 &kubelet.Image,
					"cluster_domain": &kubelet.ClusterDomain,
					//"infra_container_image": &kubelet.InfraContainerImage,
					"cluster_dns_server": &kubelet.ClusterDNSServer,
				},
				boolMapping: map[string]*bool{
					"fail_swap_on": &kubelet.FailSwapOn,
				},
				mapStrMapping: map[string]*map[string]string{
					"extra_args": &kubelet.ExtraArgs,
				},
				listStrMapping: map[string]*[]string{
					"extra_binds": &kubelet.ExtraBinds,
					"extra_env":   &kubelet.ExtraEnv,
				},
			})

			return kubelet, nil
		}
	}
	return nil, nil
}

func parseResourceKubeproxyService(d resourceData) (*v3.KubeproxyService, error) {
	if rawList, ok := d.GetOk("services_kubeproxy"); ok {
		if rawServices, ok := rawList.([]interface{}); ok && len(rawServices) > 0 {
			rawService := rawServices[0]
			kubeproxy := &v3.KubeproxyService{}
			rawMap := rawService.(map[string]interface{})

			applyMapToObj(&mapObjMapping{
				source:        rawMap,
				stringMapping: map[string]*string{
					//"image": &kubeproxy.Image,
				},
				mapStrMapping: map[string]*map[string]string{
					"extra_args": &kubeproxy.ExtraArgs,
				},
				listStrMapping: map[string]*[]string{
					"extra_binds": &kubeproxy.ExtraBinds,
					"extra_env":   &kubeproxy.ExtraEnv,
				},
			})

			return kubeproxy, nil
		}
	}
	return nil, nil
}

func parseResourceNetwork(d resourceData) (*v3.NetworkConfig, error) {
	if rawList, ok := d.GetOk("network"); ok {
		if rawNetworks, ok := rawList.([]interface{}); ok && len(rawNetworks) > 0 {
			rawNetwork := rawNetworks[0]
			network := &v3.NetworkConfig{}
			rawMap := rawNetwork.(map[string]interface{})

			applyMapToObj(&mapObjMapping{
				source: rawMap,
				stringMapping: map[string]*string{
					"plugin": &network.Plugin,
				},
				mapStrMapping: map[string]*map[string]string{
					"options": &network.Options,
				},
			})

			return network, nil
		}
	}
	return nil, nil
}

func parseResourceAuthentication(d resourceData) (*v3.AuthnConfig, error) {
	if rawList, ok := d.GetOk("authentication"); ok {
		if rawAuthns, ok := rawList.([]interface{}); ok && len(rawAuthns) > 0 {
			rawAuthn := rawAuthns[0]
			config := &v3.AuthnConfig{}
			rawMap := rawAuthn.(map[string]interface{})

			applyMapToObj(&mapObjMapping{
				source: rawMap,
				stringMapping: map[string]*string{
					"strategy": &config.Strategy,
				},
				listStrMapping: map[string]*[]string{
					"sans": &config.SANs,
				},
				mapStrMapping: map[string]*map[string]string{
					"options": &config.Options,
				},
			})

			return config, nil
		}
	}
	return nil, nil
}

func parseResourceAddons(d resourceData) (string, error) {
	if v, ok := d.GetOk("addons"); ok {
		return v.(string), nil
	}
	return "", nil
}

func parseResourceAddonsInclude(d resourceData) ([]string, error) {
	if v, ok := d.GetOk("addons_include"); ok {
		values := []string{}
		for _, e := range v.([]interface{}) {
			values = append(values, e.(string))
		}
		return values, nil
	}
	return []string{}, nil
}

func parseResourceAddonJobTimeout(d resourceData) (int, error) {
	if v, ok := d.GetOk("addon_job_timeout"); ok {
		return v.(int), nil
	}
	return 0, nil
}

func parseResourceSystemImages(d resourceData) (*v3.RKESystemImages, error) {
	if rawList, ok := d.GetOk("system_images"); ok {
		if rawImages, ok := rawList.([]interface{}); ok && len(rawImages) > 0 {
			rawImage := rawImages[0]
			config := &v3.RKESystemImages{}
			rawMap := rawImage.(map[string]interface{})

			applyMapToObj(&mapObjMapping{
				source: rawMap,
				stringMapping: map[string]*string{
					"etcd":                        &config.Etcd,
					"alpine":                      &config.Alpine,
					"nginx_proxy":                 &config.NginxProxy,
					"cert_downloader":             &config.CertDownloader,
					"kubernetes_services_sidecar": &config.KubernetesServicesSidecar,
					"kube_dns":                    &config.KubeDNS,
					"dnsmasq":                     &config.DNSmasq,
					"kube_dns_sidecar":            &config.KubeDNSSidecar,
					"kube_dns_autoscaler":         &config.KubeDNSAutoscaler,
					"kubernetes":                  &config.Kubernetes,
					"flannel":                     &config.Flannel,
					"flannel_cni":                 &config.FlannelCNI,
					"calico_node":                 &config.CalicoNode,
					"calico_cni":                  &config.CalicoCNI,
					"calico_controllers":          &config.CalicoControllers,
					"calico_ctl":                  &config.CalicoCtl,
					"canal_node":                  &config.CanalNode,
					"canal_cni":                   &config.CanalCNI,
					"canal_flannel":               &config.CanalFlannel,
					"weave_node":                  &config.WeaveNode,
					"weave_cni":                   &config.WeaveCNI,
					"pod_infra_container":         &config.PodInfraContainer,
					"ingress":                     &config.Ingress,
					"ingress_backend":             &config.IngressBackend,
					"metrics_server":              &config.MetricsServer,
				},
			})

			return config, nil
		}
	}
	return nil, nil
}

func parseResourceSSHKeyPath(d resourceData) (string, error) {
	if v, ok := d.GetOk("ssh_key_path"); ok {
		return v.(string), nil
	}
	return "", nil
}

func parseResourceSSHAgentAuth(d resourceData) (bool, error) {
	if v, ok := d.GetOk("ssh_agent_auth"); ok {
		return v.(bool), nil
	}
	return false, nil
}

func parseResourceBastionHost(d resourceData) (*v3.BastionHost, error) {
	if rawList, ok := d.GetOk("bastion_host"); ok {
		if rawHosts, ok := rawList.([]interface{}); ok && len(rawHosts) > 0 {
			rawHost := rawHosts[0]
			config := &v3.BastionHost{}

			rawMap := rawHost.(map[string]interface{})

			applyMapToObj(&mapObjMapping{
				source: rawMap,
				stringMapping: map[string]*string{
					"address":      &config.Address,
					"user":         &config.User,
					"ssh_key":      &config.SSHKey,
					"ssh_key_path": &config.SSHKeyPath,
				},
				boolMapping: map[string]*bool{
					"ssh_agent_auth": &config.SSHAgentAuth,
				},
			})

			if v, ok := rawMap["port"]; ok {
				p := v.(int)
				if p > 0 {
					config.Port = fmt.Sprintf("%d", p)
				}
			}

			return config, nil
		}
	}
	return nil, nil
}

func parseResourceMonitoring(d resourceData) (*v3.MonitoringConfig, error) {
	if rawList, ok := d.GetOk("monitoring"); ok {
		if rawMonitorings, ok := rawList.([]interface{}); ok && len(rawMonitorings) > 0 {
			rawMonitoring := rawMonitorings[0]
			config := &v3.MonitoringConfig{}

			rawMap := rawMonitoring.(map[string]interface{})
			if v, ok := rawMap["provider"]; ok {
				config.Provider = v.(string)
			}

			if v, ok := rawMap["options"]; ok {
				options := map[string]string{}
				values := v.(map[string]interface{})
				for k, v := range values {
					if v, ok := v.(string); ok {
						options[k] = v
					}
				}
				config.Options = options
			}
			return config, nil
		}
	}
	return nil, nil
}

func parseResourceAuthorization(d resourceData) (*v3.AuthzConfig, error) {
	if rawList, ok := d.GetOk("authorization"); ok {
		if rawAuthzs, ok := rawList.([]interface{}); ok && len(rawAuthzs) > 0 {
			rawAuthz := rawAuthzs[0]
			config := &v3.AuthzConfig{}

			rawMap := rawAuthz.(map[string]interface{})
			if v, ok := rawMap["mode"]; ok {
				config.Mode = v.(string)
			}

			if v, ok := rawMap["options"]; ok {
				options := map[string]string{}
				values := v.(map[string]interface{})
				for k, v := range values {
					if v, ok := v.(string); ok {
						options[k] = v
					}
				}
				config.Options = options
			}
			return config, nil
		}
	}
	return nil, nil
}

func parseResourceIgnoreDockerVersion(d resourceData) (bool, error) {
	if v, ok := d.GetOk("ignore_docker_version"); ok {
		return v.(bool), nil
	}
	return false, nil
}

func parseResourceVersion(d resourceData) (string, error) {
	if v, ok := d.GetOk("kubernetes_version"); ok {
		return v.(string), nil
	}
	return "", nil
}

func parseResourcePrivateRegistries(d resourceData) ([]v3.PrivateRegistry, error) {
	if rawList, ok := d.GetOk("private_registries"); ok {
		if rawRegistries, ok := rawList.([]interface{}); ok && len(rawRegistries) > 0 {
			res := []v3.PrivateRegistry{}
			for _, rawRegistry := range rawRegistries {
				config := v3.PrivateRegistry{}
				rawMap := rawRegistry.(map[string]interface{})
				if v, ok := rawMap["url"]; ok {
					config.URL = v.(string)
				}
				if v, ok := rawMap["user"]; ok {
					config.User = v.(string)
				}
				if v, ok := rawMap["password"]; ok {
					config.Password = v.(string)
				}
				if v, ok := rawMap["is_default"]; ok {
					config.IsDefault = v.(bool)
				}
				res = append(res, config)
			}
			return res, nil
		}
	}
	return nil, nil
}

func parseResourceIngress(d resourceData) (*v3.IngressConfig, error) {
	if rawList, ok := d.GetOk("ingress"); ok {
		if rawIngresses, ok := rawList.([]interface{}); ok && len(rawIngresses) > 0 {
			rawIngress := rawIngresses[0]
			config := &v3.IngressConfig{}

			rawMap := rawIngress.(map[string]interface{})
			if v, ok := rawMap["provider"]; ok {
				config.Provider = v.(string)
			}

			if v, ok := rawMap["options"]; ok {
				options := map[string]string{}
				values := v.(map[string]interface{})
				for k, v := range values {
					if v, ok := v.(string); ok {
						options[k] = v
					}
				}
				config.Options = options
			}

			if v, ok := rawMap["node_selector"]; ok {
				options := map[string]string{}
				values := v.(map[string]interface{})
				for k, v := range values {
					if v, ok := v.(string); ok {
						options[k] = v
					}
				}
				config.NodeSelector = options
			}
			if v, ok := rawMap["extra_args"]; ok {
				extraArgs := map[string]string{}
				args := v.(map[string]interface{})
				for k, v := range args {
					if v, ok := v.(string); ok {
						extraArgs[k] = v
					}
				}
				config.ExtraArgs = extraArgs
			}

			return config, nil
		}
	}
	return nil, nil
}

func parseResourceClusterName(d resourceData) (string, error) {
	if v, ok := d.GetOk("cluster_name"); ok {
		return v.(string), nil
	}
	return "", nil
}

func parseResourceCloudProvider(d resourceData) (*v3.CloudProvider, error) {
	if rawList, ok := d.GetOk("cloud_provider"); ok {
		if rawProviders, ok := rawList.([]interface{}); ok && len(rawProviders) > 0 {
			rawProvider := rawProviders[0]
			config := &v3.CloudProvider{}

			rawProviderMap := rawProvider.(map[string]interface{})
			if v, ok := rawProviderMap["name"]; ok {
				config.Name = v.(string)
			}

			/*
				if v, ok := rawProviderMap["cloud_config"]; ok {
					cc := map[string]string{}
					values := v.(map[string]interface{})
					for k, v := range values {
						if v, ok := v.(string); ok {
							cc[k] = v
						}
					}
					config.CloudConfig = cc
				}
			*/

			if rawList, ok := rawProviderMap["azure_cloud_config"]; ok {
				if rawCloudConfigs, ok := rawList.([]interface{}); ok && len(rawCloudConfigs) > 0 {
					rawConfig := rawCloudConfigs[0]
					c := &v3.AzureCloudProvider{}

					rawMap := rawConfig.(map[string]interface{})

					applyMapToObj(&mapObjMapping{
						source: rawMap,
						stringMapping: map[string]*string{
							"cloud":                         &c.Cloud,
							"tenant_id":                     &c.TenantID,
							"subscription_id":               &c.SubscriptionID,
							"resource_group":                &c.ResourceGroup,
							"location":                      &c.Location,
							"vnet_name":                     &c.VnetName,
							"vnet_resource_group":           &c.VnetResourceGroup,
							"route_table_name":              &c.RouteTableName,
							"primary_availability_set_name": &c.PrimaryAvailabilitySetName,
							"vm_type":                       &c.VMType,
							"primary_scale_set_name":        &c.PrimaryScaleSetName,
							"aad_client_id":                 &c.AADClientID,
							"aad_client_secret":             &c.AADClientSecret,
							"aad_client_cert_path":          &c.AADClientCertPath,
							"aad_client_cert_password":      &c.AADClientCertPassword,
						},
						intMapping: map[string]*int{
							"cloud_provider_backoff_retries":   &c.CloudProviderBackoffRetries,
							"cloud_provider_backoff_exponent":  &c.CloudProviderBackoffExponent,
							"cloud_provider_backoff_duration":  &c.CloudProviderBackoffDuration,
							"cloud_provider_backoff_jitter":    &c.CloudProviderBackoffJitter,
							"cloud_provider_rate_limit_qps":    &c.CloudProviderRateLimitQPS,
							"cloud_provider_rate_limit_bucket": &c.CloudProviderRateLimitBucket,
							"maximum_load_balancer_rule_count": &c.MaximumLoadBalancerRuleCount,
						},
						boolMapping: map[string]*bool{
							"cloud_provider_backoff":         &c.CloudProviderBackoff,
							"cloud_provider_rate_limit":      &c.CloudProviderRateLimit,
							"use_instance_metadata":          &c.UseInstanceMetadata,
							"use_managed_identity_extension": &c.UseManagedIdentityExtension,
						},
					})

					config.AzureCloudProvider = c
				}
			}

			if rawList, ok := rawProviderMap["vsphere_cloud_config"]; ok {
				if rawCloudConfigs, ok := rawList.([]interface{}); ok && len(rawCloudConfigs) > 0 {
					rawConfig := rawCloudConfigs[0]
					c := &v3.VsphereCloudProvider{}

					rawVSphereConfigMap := rawConfig.(map[string]interface{})

					// global
					if rawList, ok := rawVSphereConfigMap["global"]; ok {
						if rawGlobals, ok := rawList.([]interface{}); ok && len(rawGlobals) > 0 {
							rawGlobal := rawGlobals[0]
							rawMap := rawGlobal.(map[string]interface{})
							global := v3.GlobalVsphereOpts{}

							applyMapToObj(&mapObjMapping{
								source: rawMap,
								stringMapping: map[string]*string{
									"user":        &global.User,
									"password":    &global.Password,
									"server":      &global.VCenterIP,
									"port":        &global.VCenterPort,
									"datacenter":  &global.Datacenter,
									"datacenters": &global.Datacenters,
									"datastore":   &global.DefaultDatastore,
									"working_dir": &global.WorkingDir,
									"vm_uuid":     &global.VMUUID,
									"vm_name":     &global.VMName,
								},
								intMapping: map[string]*int{
									"soap_roundtrip_count": &global.RoundTripperCount,
								},
								boolMapping: map[string]*bool{
									"insecure_flag": &global.InsecureFlag,
								},
							})

							c.Global = global
						}
					}

					// virtual_center
					if rawList, ok := rawVSphereConfigMap["virtual_center"]; ok {
						if rawVCs, ok := rawList.([]interface{}); ok {
							vcs := map[string]v3.VirtualCenterConfig{}
							for _, rawVC := range rawVCs {
								var server string
								rawMap := rawVC.(map[string]interface{})
								vc := v3.VirtualCenterConfig{}

								applyMapToObj(&mapObjMapping{
									source: rawMap,
									stringMapping: map[string]*string{
										"server":      &server,
										"user":        &vc.User,
										"password":    &vc.Password,
										"port":        &vc.VCenterPort,
										"datacenters": &vc.Datacenters,
									},
									intMapping: map[string]*int{
										"soap_roundtrip_count": &vc.RoundTripperCount,
									},
								})

								vcs[server] = vc
							}
							c.VirtualCenter = vcs
						}
					}

					// network
					if rawList, ok := rawVSphereConfigMap["network"]; ok {
						if rawNetworks, ok := rawList.([]interface{}); ok && len(rawNetworks) > 0 {
							rawNetwork := rawNetworks[0]
							rawMap := rawNetwork.(map[string]interface{})
							network := v3.NetworkVshpereOpts{}

							applyMapToObj(&mapObjMapping{
								source: rawMap,
								stringMapping: map[string]*string{
									"public_network": &network.PublicNetwork,
								},
							})

							c.Network = network
						}
					}

					// disk
					if rawList, ok := rawVSphereConfigMap["disk"]; ok {
						if rawDisks, ok := rawList.([]interface{}); ok && len(rawDisks) > 0 {
							rawDisk := rawDisks[0]
							rawMap := rawDisk.(map[string]interface{})
							disk := v3.DiskVsphereOpts{}

							applyMapToObj(&mapObjMapping{
								source: rawMap,
								stringMapping: map[string]*string{
									"scsi_controller_type": &disk.SCSIControllerType,
								},
							})

							c.Disk = disk
						}
					}

					// workspace
					if rawList, ok := rawVSphereConfigMap["workspace"]; ok {
						if rawWSs, ok := rawList.([]interface{}); ok && len(rawWSs) > 0 {
							rawWS := rawWSs[0]
							rawMap := rawWS.(map[string]interface{})
							ws := v3.WorkspaceVsphereOpts{}

							applyMapToObj(&mapObjMapping{
								source: rawMap,
								stringMapping: map[string]*string{
									"server":            &ws.VCenterIP,
									"datacenter":        &ws.Datacenter,
									"folder":            &ws.Folder,
									"default_datastore": &ws.DefaultDatastore,
									"resourcepool_path": &ws.ResourcePoolPath,
								},
							})

							c.Workspace = ws
						}
					}

					config.VsphereCloudProvider = c
				}
			}

			if rawList, ok := rawProviderMap["openstack_cloud_config"]; ok {
				if rawCloudConfigs, ok := rawList.([]interface{}); ok && len(rawCloudConfigs) > 0 {
					rawConfig := rawCloudConfigs[0]
					c := &v3.OpenstackCloudProvider{}

					rawOpenStackConfigMap := rawConfig.(map[string]interface{})

					// global
					if rawList, ok := rawOpenStackConfigMap["global"]; ok {
						if rawGlobals, ok := rawList.([]interface{}); ok && len(rawGlobals) > 0 {
							rawGlobal := rawGlobals[0]
							rawMap := rawGlobal.(map[string]interface{})
							global := v3.GlobalOpenstackOpts{}

							applyMapToObj(&mapObjMapping{
								source: rawMap,
								stringMapping: map[string]*string{
									"auth_url":    &global.AuthURL,
									"username":    &global.Username,
									"user_id":     &global.UserID,
									"password":    &global.Password,
									"tenant_id":   &global.TenantID,
									"tenant_name": &global.TenantName,
									"trust_id":    &global.TrustID,
									"domain_id":   &global.DomainID,
									"domain_name": &global.DomainName,
									"region":      &global.Region,
									"ca_file":     &global.CAFile,
								},
							})

							c.Global = global
						}
					}

					// load_balancer
					if rawList, ok := rawOpenStackConfigMap["load_balancer"]; ok {
						if rawLBs, ok := rawList.([]interface{}); ok && len(rawLBs) > 0 {
							rawLB := rawLBs[0]
							rawMap := rawLB.(map[string]interface{})
							lb := v3.LoadBalancerOpenstackOpts{}

							applyMapToObj(&mapObjMapping{
								source: rawMap,
								stringMapping: map[string]*string{
									"lb_version":          &lb.LBVersion,
									"subnet_id":           &lb.SubnetID,
									"floating_network_id": &lb.FloatingNetworkID,
									"lb_method":           &lb.LBMethod,
									"lb_provider":         &lb.LBProvider,
								},
								boolMapping: map[string]*bool{
									"use_octavia":            &lb.UseOctavia,
									"create_monitor":         &lb.CreateMonitor,
									"manage_security_groups": &lb.ManageSecurityGroups,
								},
								intMapping: map[string]*int{
									"monitor_delay":       &lb.MonitorDelay,
									"monitor_timeout":     &lb.MonitorTimeout,
									"monitor_max_retries": &lb.MonitorMaxRetries,
								},
							})

							c.LoadBalancer = lb
						}
					}

					// block_storage
					if rawList, ok := rawOpenStackConfigMap["block_storage"]; ok {
						if rawBSs, ok := rawList.([]interface{}); ok && len(rawBSs) > 0 {
							rawBS := rawBSs[0]
							rawMap := rawBS.(map[string]interface{})
							bs := v3.BlockStorageOpenstackOpts{}

							applyMapToObj(&mapObjMapping{
								source: rawMap,
								stringMapping: map[string]*string{
									"bs_version": &bs.BSVersion,
								},
								boolMapping: map[string]*bool{
									"trust_device_path": &bs.TrustDevicePath,
									"ignore_volume_az":  &bs.IgnoreVolumeAZ,
								},
							})

							c.BlockStorage = bs
						}
					}

					// route
					if rawList, ok := rawOpenStackConfigMap["route"]; ok {
						if rawRouters, ok := rawList.([]interface{}); ok && len(rawRouters) > 0 {
							rawRouter := rawRouters[0]
							rawMap := rawRouter.(map[string]interface{})
							router := v3.RouteOpenstackOpts{}

							applyMapToObj(&mapObjMapping{
								source: rawMap,
								stringMapping: map[string]*string{
									"router_id": &router.RouterID,
								},
							})

							c.Route = router
						}
					}

					// metadata
					if rawList, ok := rawOpenStackConfigMap["metadata"]; ok {
						if rawMetadataList, ok := rawList.([]interface{}); ok && len(rawMetadataList) > 0 {
							rawMetadata := rawMetadataList[0]
							rawMap := rawMetadata.(map[string]interface{})
							meta := v3.MetadataOpenstackOpts{}

							applyMapToObj(&mapObjMapping{
								source: rawMap,
								stringMapping: map[string]*string{
									"search_order": &meta.SearchOrder,
								},
								intMapping: map[string]*int{
									"request_timeout": &meta.RequestTimeout,
								},
							})

							c.Metadata = meta
						}
					}

					config.OpenstackCloudProvider = c
				}
			}

			if v, ok := rawProviderMap["custom_cloud_config"]; ok {
				config.CustomCloudProvider = v.(string)
			}

			return config, nil
		}
	}
	return nil, nil
}

func parseResourcePrefixPath(d resourceData) (string, error) {
	if v, ok := d.GetOk("prefix_path"); ok {
		return v.(string), nil
	}
	return "", nil
}

func clusterToState(cluster *cluster.Cluster, d stateBuilder) error {

	if cluster == nil {
		d.SetId("")
		return nil
	}

	// services
	d.Set("services_etcd", []interface{}{ // nolint
		map[string]interface{}{
			//"image":         cluster.Services.Etcd.Image,
			"extra_args":    cluster.Services.Etcd.ExtraArgs,
			"extra_binds":   cluster.Services.Etcd.ExtraBinds,
			"extra_env":     cluster.Services.Etcd.ExtraEnv,
			"external_urls": cluster.Services.Etcd.ExternalURLs,
			"ca_cert":       cluster.Services.Etcd.CACert,
			"cert":          cluster.Services.Etcd.Cert,
			"key":           cluster.Services.Etcd.Key,
			"path":          cluster.Services.Etcd.Path,
			"snapshot":      cluster.Services.Etcd.Snapshot,
			"retention":     cluster.Services.Etcd.Retention,
			"creation":      cluster.Services.Etcd.Creation,
		},
	})

	d.Set("services_kube_api", []interface{}{ // nolint
		map[string]interface{}{
			//"image":                    cluster.Services.KubeAPI.Image,
			"extra_args":               cluster.Services.KubeAPI.ExtraArgs,
			"extra_binds":              cluster.Services.KubeAPI.ExtraBinds,
			"extra_env":                cluster.Services.KubeAPI.ExtraEnv,
			"service_cluster_ip_range": cluster.Services.KubeAPI.ServiceClusterIPRange,
			"service_node_port_range":  cluster.Services.KubeAPI.ServiceNodePortRange,
			"pod_security_policy":      cluster.Services.KubeAPI.PodSecurityPolicy,
		},
	})

	d.Set("services_kube_controller", []interface{}{ // nolint
		map[string]interface{}{
			//"image":                    cluster.Services.KubeController.Image,
			"extra_args":               cluster.Services.KubeController.ExtraArgs,
			"extra_binds":              cluster.Services.KubeController.ExtraBinds,
			"extra_env":                cluster.Services.KubeController.ExtraEnv,
			"cluster_cidr":             cluster.Services.KubeController.ClusterCIDR,
			"service_cluster_ip_range": cluster.Services.KubeController.ServiceClusterIPRange,
		},
	})

	d.Set("services_scheduler", []interface{}{ // nolint
		map[string]interface{}{
			//"image":       cluster.Services.Scheduler.Image,
			"extra_args":  cluster.Services.Scheduler.ExtraArgs,
			"extra_binds": cluster.Services.Scheduler.ExtraBinds,
			"extra_env":   cluster.Services.Scheduler.ExtraEnv,
		},
	})

	d.Set("services_kubelet", []interface{}{ // nolint
		map[string]interface{}{
			//"image":                 cluster.Services.Kubelet.Image,
			"extra_args":     cluster.Services.Kubelet.ExtraArgs,
			"extra_binds":    cluster.Services.Kubelet.ExtraBinds,
			"extra_env":      cluster.Services.Kubelet.ExtraEnv,
			"cluster_domain": cluster.Services.Kubelet.ClusterDomain,
			//"infra_container_image": cluster.Services.Kubelet.InfraContainerImage,
			"cluster_dns_server": cluster.Services.Kubelet.ClusterDNSServer,
			"fail_swap_on":       cluster.Services.Kubelet.FailSwapOn,
		},
	})

	d.Set("services_kubeproxy", []interface{}{ // nolint
		map[string]interface{}{
			//"image":       cluster.Services.Kubeproxy.Image,
			"extra_args":  cluster.Services.Kubeproxy.ExtraArgs,
			"extra_binds": cluster.Services.Kubeproxy.ExtraBinds,
			"extra_env":   cluster.Services.Kubeproxy.ExtraEnv,
		},
	})

	d.Set("network", []interface{}{ // nolint
		map[string]interface{}{
			"plugin":  cluster.Network.Plugin,
			"options": cluster.Network.Options,
		},
	})

	d.Set("authentication", []interface{}{ // nolint
		map[string]interface{}{
			"strategy": cluster.Authentication.Strategy,
			"options":  cluster.Authentication.Options,
			"sans":     cluster.Authentication.SANs,
		},
	})

	d.Set("addons", cluster.Addons)                     // nolint
	d.Set("addons_include", cluster.AddonsInclude)      // nolint
	d.Set("addon_job_timeout", cluster.AddonJobTimeout) // nolint

	d.Set("ssh_key_path", cluster.SSHKeyPath)     // nolint
	d.Set("ssh_agent_auth", cluster.SSHAgentAuth) // nolint

	bastionHost := map[string]interface{}{}
	bastionHost["address"] = cluster.BastionHost.Address
	if cluster.BastionHost.Port != "" {
		if port, err := strconv.Atoi(cluster.BastionHost.Port); err == nil {
			if port > 0 {
				bastionHost["port"] = port
			}
		} else {
			return err
		}
	}
	bastionHost["user"] = cluster.BastionHost.User
	bastionHost["ssh_agent_auth"] = cluster.BastionHost.SSHAgentAuth
	bastionHost["ssh_key"] = cluster.BastionHost.SSHKey
	bastionHost["ssh_key_path"] = cluster.BastionHost.SSHKeyPath
	d.Set("bastion_host", []interface{}{bastionHost}) // nolint

	d.Set("monitoring", []interface{}{ // nolint
		map[string]interface{}{
			"provider": cluster.Monitoring.Provider,
			"options":  cluster.Monitoring.Options,
		},
	})

	d.Set("authorization", []interface{}{ // nolint
		map[string]interface{}{
			"mode":    cluster.Authorization.Mode,
			"options": cluster.Authorization.Options,
		},
	})

	d.Set("ignore_docker_version", cluster.IgnoreDockerVersion) // nolint
	d.Set("kubernetes_version", cluster.Version)                // nolint

	registries := []interface{}{}
	for _, registry := range cluster.PrivateRegistries {
		r := map[string]interface{}{}
		r["url"] = registry.URL
		r["user"] = registry.User
		r["password"] = registry.Password
		r["is_default"] = registry.IsDefault
		registries = append(registries, r)
	}
	d.Set("private_registries", registries) // nolint

	d.Set("ingress", []interface{}{ // nolint
		map[string]interface{}{
			"provider":      cluster.Ingress.Provider,
			"options":       cluster.Ingress.Options,
			"node_selector": cluster.Ingress.NodeSelector,
			"extra_args":    cluster.Ingress.ExtraArgs,
		},
	})

	d.Set("cluster_name", cluster.ClusterName)      // nolint
	d.Set("kube_admin_user", pki.KubeAdminCertName) // nolint

	var apiServerURL = ""
	if len(cluster.ControlPlaneHosts) > 0 {
		apiServerURL = fmt.Sprintf("https://" + cluster.ControlPlaneHosts[0].Address + ":6443")
	}
	d.Set("api_server_url", apiServerURL) // nolint

	cloudProvider := map[string]interface{}{
		"name":                cluster.CloudProvider.Name,
		"custom_cloud_config": cluster.CloudProvider.CustomCloudProvider,
	}
	if cp := cluster.CloudProvider.AzureCloudProvider; cp != nil {
		acp := map[string]interface{}{}
		acp["cloud"] = cp.Cloud
		acp["tenant_id"] = cp.TenantID
		acp["subscription_id"] = cp.SubscriptionID
		acp["resource_group"] = cp.ResourceGroup
		acp["location"] = cp.Location
		acp["vnet_name"] = cp.VnetName
		acp["vnet_resource_group"] = cp.VnetResourceGroup
		acp["route_table_name"] = cp.RouteTableName
		acp["primary_availability_set_name"] = cp.PrimaryAvailabilitySetName
		acp["vm_type"] = cp.VMType
		acp["primary_scale_set_name"] = cp.PrimaryScaleSetName
		acp["aad_client_id"] = cp.AADClientID
		acp["aad_client_secret"] = cp.AADClientSecret
		acp["aad_client_cert_path"] = cp.AADClientCertPath
		acp["aad_client_cert_password"] = cp.AADClientCertPassword
		acp["cloud_provider_backoff"] = cp.CloudProviderBackoff
		acp["cloud_provider_backoff_retries"] = cp.CloudProviderBackoffRetries
		acp["cloud_provider_backoff_exponent"] = cp.CloudProviderBackoffExponent
		acp["cloud_provider_backoff_duration"] = cp.CloudProviderBackoffDuration
		acp["cloud_provider_backoff_jitter"] = cp.CloudProviderBackoffJitter
		acp["cloud_provider_rate_limit"] = cp.CloudProviderRateLimit
		acp["cloud_provider_rate_limit_qps"] = cp.CloudProviderRateLimitQPS
		acp["cloud_provider_rate_limit_bucket"] = cp.CloudProviderRateLimitBucket
		acp["use_instance_metadata"] = cp.UseInstanceMetadata
		acp["use_managed_identity_extension"] = cp.UseManagedIdentityExtension
		acp["maximum_load_balancer_rule_count"] = cp.MaximumLoadBalancerRuleCount

		cloudProvider["azure_cloud_config"] = []interface{}{acp}
	}

	if cp := cluster.CloudProvider.VsphereCloudProvider; cp != nil {
		vcp := map[string]interface{}{}

		global := map[string]interface{}{}
		global["user"] = cp.Global.User
		global["password"] = cp.Global.Password
		global["server"] = cp.Global.VCenterIP
		global["port"] = cp.Global.VCenterPort
		global["insecure_flag"] = cp.Global.InsecureFlag
		global["datacenter"] = cp.Global.Datacenter
		global["datacenters"] = cp.Global.Datacenters
		global["datastore"] = cp.Global.DefaultDatastore
		global["working_dir"] = cp.Global.WorkingDir
		global["soap_roundtrip_count"] = cp.Global.RoundTripperCount
		global["vm_uuid"] = cp.Global.VMUUID
		global["vm_name"] = cp.Global.VMName
		vcp["global"] = []interface{}{global}

		var vcs []interface{}
		for k, v := range cp.VirtualCenter {
			vc := map[string]interface{}{}
			vc["server"] = k
			vc["user"] = v.User
			vc["password"] = v.Password
			vc["port"] = v.VCenterPort
			vc["datacenters"] = v.Datacenters
			vc["soap_roundtrip_count"] = v.RoundTripperCount
			vcs = append(vcs, vc)
		}
		vcp["virtual_center"] = vcs

		nw := map[string]interface{}{}
		nw["public_network"] = cp.Network.PublicNetwork
		vcp["network"] = []interface{}{nw}

		disk := map[string]interface{}{}
		disk["scsi_controller_type"] = cp.Disk.SCSIControllerType
		vcp["disk"] = []interface{}{disk}

		ws := map[string]interface{}{}
		ws["server"] = cp.Workspace.VCenterIP
		ws["datacenter"] = cp.Workspace.Datacenter
		ws["folder"] = cp.Workspace.Folder
		ws["default_datastore"] = cp.Workspace.DefaultDatastore
		ws["resourcepool_path"] = cp.Workspace.ResourcePoolPath
		vcp["workspace"] = []interface{}{ws}

		cloudProvider["vsphere_cloud_config"] = []interface{}{vcp}
	}

	if cp := cluster.CloudProvider.OpenstackCloudProvider; cp != nil {
		ocp := map[string]interface{}{}

		global := map[string]interface{}{}
		global["auth_url"] = cp.Global.AuthURL
		global["username"] = cp.Global.Username
		global["user_id"] = cp.Global.UserID
		global["password"] = cp.Global.Password
		global["tenant_id"] = cp.Global.TenantID
		global["tenant_name"] = cp.Global.TenantName
		global["trust_id"] = cp.Global.TrustID
		global["domain_id"] = cp.Global.DomainID
		global["domain_name"] = cp.Global.DomainName
		global["region"] = cp.Global.Region
		global["ca_file"] = cp.Global.CAFile
		ocp["global"] = []interface{}{global}

		lb := map[string]interface{}{}
		lb["lb_version"] = cp.LoadBalancer.LBVersion
		lb["use_octavia"] = cp.LoadBalancer.UseOctavia
		lb["subnet_id"] = cp.LoadBalancer.SubnetID
		lb["floating_network_id"] = cp.LoadBalancer.FloatingNetworkID
		lb["lb_method"] = cp.LoadBalancer.LBMethod
		lb["lb_provider"] = cp.LoadBalancer.LBProvider
		lb["create_monitor"] = cp.LoadBalancer.CreateMonitor
		lb["monitor_delay"] = cp.LoadBalancer.MonitorDelay
		lb["monitor_timeout"] = cp.LoadBalancer.MonitorTimeout
		lb["monitor_max_retries"] = cp.LoadBalancer.MonitorMaxRetries
		lb["manage_security_groups"] = cp.LoadBalancer.ManageSecurityGroups
		ocp["load_balancer"] = []interface{}{lb}

		bs := map[string]interface{}{}
		bs["bs_version"] = cp.BlockStorage.BSVersion
		bs["trust_device_path"] = cp.BlockStorage.TrustDevicePath
		bs["ignore_volume_az"] = cp.BlockStorage.IgnoreVolumeAZ
		ocp["block_storage"] = []interface{}{bs}

		ro := map[string]interface{}{}
		ro["router_id"] = cp.Route.RouterID
		ocp["route"] = []interface{}{ro}

		meta := map[string]interface{}{}
		meta["search_order"] = cp.Metadata.SearchOrder
		meta["request_timeout"] = cp.Metadata.RequestTimeout
		ocp["metadata"] = []interface{}{meta}

		cloudProvider["openstack_cloud_config"] = []interface{}{ocp}
	}

	d.Set("cloud_provider", []interface{}{cloudProvider}) // nolint

	d.Set("prefix_path", cluster.PrefixPath) // nolint

	// computed values
	certs := []interface{}{}
	for k, v := range cluster.Certificates {

		certPEM := ""
		if v.Certificate != nil {
			certPEM = certificateToPEM(v.Certificate)
		}
		privateKeyPEM := ""
		if v.Key != nil {
			privateKeyPEM = privateKeyToPEM(v.Key)
		}

		cert := map[string]interface{}{
			"id":              k,
			"certificate":     certPEM,
			"key":             privateKeyPEM,
			"config":          v.Config,
			"name":            v.Name,
			"common_name":     v.CommonName,
			"ou_name":         v.OUName,
			"env_name":        v.EnvName,
			"path":            v.Path,
			"key_env_name":    v.KeyEnvName,
			"key_path":        v.KeyPath,
			"config_env_name": v.ConfigEnvName,
			"config_path":     v.ConfigPath,
		}
		certs = append(certs, cert)
	}
	d.Set("certificates", certs) // nolint

	d.Set("cluster_domain", cluster.ClusterDomain)        // nolint
	d.Set("cluster_cidr", cluster.ClusterCIDR)            // nolint
	d.Set("cluster_dns_server", cluster.ClusterDNSServer) // nolint

	hostsMap := map[string][]*hosts.Host{
		"etcd_hosts":          cluster.EtcdHosts,
		"worker_hosts":        cluster.WorkerHosts,
		"control_plane_hosts": cluster.ControlPlaneHosts,
		"inactive_hosts":      cluster.InactiveHosts,
	}

	for key, hosts := range hostsMap {
		values := []map[string]interface{}{}
		for _, host := range hosts {
			values = append(values, map[string]interface{}{
				"node_name": host.NodeName,
				"address":   host.Address,
			})
		}
		d.Set(key, values) // nolint
	}

	return nil
}

func privateKeyToPEM(key *rsa.PrivateKey) string {
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)
	return string(pemdata)
}

func certificateToPEM(cert *x509.Certificate) string {
	pemdata := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: cert.Raw,
		},
	)
	return string(pemdata)
}
