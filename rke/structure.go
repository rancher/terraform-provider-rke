package rke

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strconv"

	"github.com/rancher/rke/cluster"
	"github.com/rancher/rke/hosts"
	"github.com/rancher/types/apis/management.cattle.io/v3"
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
	nodes, err := parseResourceRKEConfigNode(d)
	if err != nil {
		return err
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

func parseResourceRKEConfigNode(d resourceData) ([]v3.RKEConfigNode, error) {
	nodes := []v3.RKEConfigNode{}
	if rawNodes, ok := d.GetOk("nodes"); ok {
		nodeList := rawNodes.([]interface{})
		for _, rawNode := range nodeList {
			nodeValues := rawNode.(map[string]interface{})
			node := v3.RKEConfigNode{}

			valueMapping := map[string]*string{
				"node_name":         &node.NodeName,
				"address":           &node.Address,
				"internal_address":  &node.InternalAddress,
				"hostname_override": &node.HostnameOverride,
				"user":              &node.User,
				"docker_socket":     &node.DockerSocket,
				"ssh_key":           &node.SSHKey,
				"ssh_key_path":      &node.SSHKeyPath,
			}

			for key, dest := range valueMapping {
				if v, ok := nodeValues[key]; ok {
					*dest = v.(string)
				}
			}

			if v, ok := nodeValues["port"]; ok {
				p := v.(int)
				if p > 0 {
					node.Port = fmt.Sprintf("%d", p)
				}
			}

			if rawRoles, ok := nodeValues["role"]; ok {
				roles := []string{}
				for _, e := range rawRoles.([]interface{}) {
					roles = append(roles, e.(string))
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
			if v, ok := rawMap["image"]; ok {
				etcd.Image = v.(string)
			}

			if v, ok := rawMap["extra_args"]; ok {
				extraArgs := map[string]string{}
				args := v.(map[string]interface{})
				for k, v := range args {
					if v, ok := v.(string); ok {
						extraArgs[k] = v
					}
				}
				etcd.ExtraArgs = extraArgs
			}

			if v, ok := rawMap["extra_binds"]; ok {
				extraBinds := []string{}
				for _, e := range v.([]interface{}) {
					extraBinds = append(extraBinds, e.(string))
				}
				etcd.ExtraBinds = extraBinds
			}

			if v, ok := rawMap["external_urls"]; ok {
				externalURLs := []string{}
				for _, e := range v.([]interface{}) {
					externalURLs = append(externalURLs, e.(string))
				}
				etcd.ExternalURLs = externalURLs
			}

			valueMapping := map[string]*string{
				"ca_cert": &etcd.CACert,
				"cert":    &etcd.Cert,
				"key":     &etcd.Key,
				"path":    &etcd.Path,
			}

			for key, dest := range valueMapping {
				if v, ok := rawMap[key]; ok {
					*dest = v.(string)
				}
			}

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
			if v, ok := rawMap["image"]; ok {
				kubeAPI.Image = v.(string)
			}

			if v, ok := rawMap["extra_args"]; ok {
				extraArgs := map[string]string{}
				args := v.(map[string]interface{})
				for k, v := range args {
					if v, ok := v.(string); ok {
						extraArgs[k] = v
					}
				}
				kubeAPI.ExtraArgs = extraArgs
			}

			if v, ok := rawMap["extra_binds"]; ok {
				extraBinds := []string{}
				for _, e := range v.([]interface{}) {
					extraBinds = append(extraBinds, e.(string))
				}
				kubeAPI.ExtraBinds = extraBinds
			}

			if v, ok := rawMap["service_cluster_ip_range"]; ok {
				kubeAPI.ServiceClusterIPRange = v.(string)
			}

			if v, ok := rawMap["pod_security_policy"]; ok {
				kubeAPI.PodSecurityPolicy = v.(bool)
			}
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
			if v, ok := rawMap["image"]; ok {
				kubeController.Image = v.(string)
			}

			if v, ok := rawMap["extra_args"]; ok {
				extraArgs := map[string]string{}
				args := v.(map[string]interface{})
				for k, v := range args {
					if v, ok := v.(string); ok {
						extraArgs[k] = v
					}
				}
				kubeController.ExtraArgs = extraArgs
			}

			if v, ok := rawMap["extra_binds"]; ok {
				extraBinds := []string{}
				for _, e := range v.([]interface{}) {
					extraBinds = append(extraBinds, e.(string))
				}
				kubeController.ExtraBinds = extraBinds
			}
			if v, ok := rawMap["cluster_cidr"]; ok {
				kubeController.ClusterCIDR = v.(string)
			}
			if v, ok := rawMap["service_cluster_ip_range"]; ok {
				kubeController.ServiceClusterIPRange = v.(string)
			}
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
			if v, ok := rawMap["image"]; ok {
				scheduler.Image = v.(string)
			}

			if v, ok := rawMap["extra_args"]; ok {
				extraArgs := map[string]string{}
				args := v.(map[string]interface{})
				for k, v := range args {
					if v, ok := v.(string); ok {
						extraArgs[k] = v
					}
				}
				scheduler.ExtraArgs = extraArgs
			}

			if v, ok := rawMap["extra_binds"]; ok {
				extraBinds := []string{}
				for _, e := range v.([]interface{}) {
					extraBinds = append(extraBinds, e.(string))
				}
				scheduler.ExtraBinds = extraBinds
			}
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
			if v, ok := rawMap["image"]; ok {
				kubelet.Image = v.(string)
			}

			if v, ok := rawMap["extra_args"]; ok {
				extraArgs := map[string]string{}
				args := v.(map[string]interface{})
				for k, v := range args {
					if v, ok := v.(string); ok {
						extraArgs[k] = v
					}
				}
				kubelet.ExtraArgs = extraArgs
			}

			if v, ok := rawMap["extra_binds"]; ok {
				extraBinds := []string{}
				for _, e := range v.([]interface{}) {
					extraBinds = append(extraBinds, e.(string))
				}
				kubelet.ExtraBinds = extraBinds
			}
			if v, ok := rawMap["cluster_domain"]; ok {
				kubelet.ClusterDomain = v.(string)
			}
			if v, ok := rawMap["infra_container_image"]; ok {
				kubelet.InfraContainerImage = v.(string)
			}
			if v, ok := rawMap["cluster_dns_server"]; ok {
				kubelet.ClusterDNSServer = v.(string)
			}
			if v, ok := rawMap["fail_swap_on"]; ok {
				kubelet.FailSwapOn = v.(bool)
			}
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
			if v, ok := rawMap["image"]; ok {
				kubeproxy.Image = v.(string)
			}

			if v, ok := rawMap["extra_args"]; ok {
				extraArgs := map[string]string{}
				args := v.(map[string]interface{})
				for k, v := range args {
					if v, ok := v.(string); ok {
						extraArgs[k] = v
					}
				}
				kubeproxy.ExtraArgs = extraArgs
			}

			if v, ok := rawMap["extra_binds"]; ok {
				extraBinds := []string{}
				for _, e := range v.([]interface{}) {
					extraBinds = append(extraBinds, e.(string))
				}
				kubeproxy.ExtraBinds = extraBinds
			}
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
			if v, ok := rawMap["plugin"]; ok {
				network.Plugin = v.(string)
			}

			if v, ok := rawMap["options"]; ok {
				options := map[string]string{}
				values := v.(map[string]interface{})
				for k, v := range values {
					if v, ok := v.(string); ok {
						options[k] = v
					}
				}
				network.Options = options
			}
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
			if v, ok := rawMap["strategy"]; ok {
				config.Strategy = v.(string)
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

			if v, ok := rawMap["sans"]; ok {
				sans := []string{}
				for _, e := range v.([]interface{}) {
					sans = append(sans, e.(string))
				}
				config.SANs = sans
			}
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

func parseResourceSystemImages(d resourceData) (*v3.RKESystemImages, error) {
	if rawList, ok := d.GetOk("system_images"); ok {
		if rawImages, ok := rawList.([]interface{}); ok && len(rawImages) > 0 {
			rawImage := rawImages[0]
			config := &v3.RKESystemImages{}
			rawMap := rawImage.(map[string]interface{})

			valueMapping := map[string]*string{
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
				"dashboard":                   &config.Dashboard,
				"heapster":                    &config.Heapster,
				"grafana":                     &config.Grafana,
				"influxdb":                    &config.Influxdb,
				"tiller":                      &config.Tiller,
			}

			for key, dest := range valueMapping {
				if v, ok := rawMap[key]; ok {
					*dest = v.(string)
				}
			}

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

			rawMap := rawProvider.(map[string]interface{})
			if v, ok := rawMap["name"]; ok {
				config.Name = v.(string)
			}

			if v, ok := rawMap["cloud_config"]; ok {
				cc := map[string]string{}
				values := v.(map[string]interface{})
				for k, v := range values {
					if v, ok := v.(string); ok {
						cc[k] = v
					}
				}
				config.CloudConfig = cc
			}

			return config, nil
		}
	}
	return nil, nil
}

func clusterToState(cluster *cluster.Cluster, d stateBuilder) error {

	if cluster == nil {
		d.SetId("")
		return nil
	}

	// node
	nodes := []interface{}{}
	for _, node := range cluster.Nodes {
		n := map[string]interface{}{}
		n["node_name"] = node.NodeName
		n["address"] = node.Address
		if port, err := strconv.Atoi(node.Port); err == nil {
			if port > 0 {
				n["port"] = port
			}
		} else {
			return err
		}
		n["internal_address"] = node.InternalAddress
		n["role"] = node.Role
		n["hostname_override"] = node.HostnameOverride
		n["user"] = node.User
		n["docker_socket"] = node.DockerSocket
		n["ssh_agent_auth"] = node.SSHAgentAuth
		n["ssh_key"] = node.SSHKey
		n["ssh_key_path"] = node.SSHKeyPath
		n["labels"] = node.Labels
		nodes = append(nodes, n)
	}
	d.Set("nodes", nodes) // nolint

	// services
	d.Set("services_etcd", []interface{}{ // nolint
		map[string]interface{}{
			"image":         cluster.Services.Etcd.Image,
			"extra_args":    cluster.Services.Etcd.ExtraArgs,
			"extra_binds":   cluster.Services.Etcd.ExtraBinds,
			"external_urls": cluster.Services.Etcd.ExternalURLs,
			"ca_cert":       cluster.Services.Etcd.CACert,
			"cert":          cluster.Services.Etcd.Cert,
			"key":           cluster.Services.Etcd.Key,
			"path":          cluster.Services.Etcd.Path,
		},
	})

	d.Set("services_kube_api", []interface{}{ // nolint
		map[string]interface{}{
			"image":                    cluster.Services.KubeAPI.Image,
			"extra_args":               cluster.Services.KubeAPI.ExtraArgs,
			"extra_binds":              cluster.Services.KubeAPI.ExtraBinds,
			"service_cluster_ip_range": cluster.Services.KubeAPI.ServiceClusterIPRange,
			"pod_security_policy":      cluster.Services.KubeAPI.PodSecurityPolicy,
		},
	})

	d.Set("services_kube_controller", []interface{}{ // nolint
		map[string]interface{}{
			"image":                    cluster.Services.KubeController.Image,
			"extra_args":               cluster.Services.KubeController.ExtraArgs,
			"extra_binds":              cluster.Services.KubeController.ExtraBinds,
			"cluster_cidr":             cluster.Services.KubeController.ClusterCIDR,
			"service_cluster_ip_range": cluster.Services.KubeController.ServiceClusterIPRange,
		},
	})

	d.Set("services_scheduler", []interface{}{ // nolint
		map[string]interface{}{
			"image":       cluster.Services.Scheduler.Image,
			"extra_args":  cluster.Services.Scheduler.ExtraArgs,
			"extra_binds": cluster.Services.Scheduler.ExtraBinds,
		},
	})

	d.Set("services_kubelet", []interface{}{ // nolint
		map[string]interface{}{
			"image":                 cluster.Services.Kubelet.Image,
			"extra_args":            cluster.Services.Kubelet.ExtraArgs,
			"extra_binds":           cluster.Services.Kubelet.ExtraBinds,
			"cluster_domain":        cluster.Services.Kubelet.ClusterDomain,
			"infra_container_image": cluster.Services.Kubelet.InfraContainerImage,
			"cluster_dns_server":    cluster.Services.Kubelet.ClusterDNSServer,
			"fail_swap_on":          cluster.Services.Kubelet.FailSwapOn,
		},
	})

	d.Set("services_kubeproxy", []interface{}{ // nolint
		map[string]interface{}{
			"image":       cluster.Services.Kubeproxy.Image,
			"extra_args":  cluster.Services.Kubeproxy.ExtraArgs,
			"extra_binds": cluster.Services.Kubeproxy.ExtraBinds,
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

	d.Set("addons", cluster.Addons)                // nolint
	d.Set("addons_include", cluster.AddonsInclude) // nolint

	d.Set("system_images", []interface{}{ // nolint
		map[string]interface{}{
			"etcd":                        cluster.SystemImages.Etcd,
			"alpine":                      cluster.SystemImages.Alpine,
			"nginx_proxy":                 cluster.SystemImages.NginxProxy,
			"cert_downloader":             cluster.SystemImages.CertDownloader,
			"kubernetes_services_sidecar": cluster.SystemImages.KubernetesServicesSidecar,
			"kube_dns":                    cluster.SystemImages.KubeDNS,
			"dnsmasq":                     cluster.SystemImages.DNSmasq,
			"kube_dns_sidecar":            cluster.SystemImages.KubeDNSSidecar,
			"kube_dns_autoscaler":         cluster.SystemImages.KubeDNSAutoscaler,
			"kubernetes":                  cluster.SystemImages.Kubernetes,
			"flannel":                     cluster.SystemImages.Flannel,
			"flannel_cni":                 cluster.SystemImages.FlannelCNI,
			"calico_node":                 cluster.SystemImages.CalicoNode,
			"calico_cni":                  cluster.SystemImages.CalicoCNI,
			"calico_controllers":          cluster.SystemImages.CalicoControllers,
			"calico_ctl":                  cluster.SystemImages.CalicoCtl,
			"canal_node":                  cluster.SystemImages.CanalNode,
			"canal_cni":                   cluster.SystemImages.CanalCNI,
			"canal_flannel":               cluster.SystemImages.CanalFlannel,
			"weave_node":                  cluster.SystemImages.WeaveNode,
			"weave_cni":                   cluster.SystemImages.WeaveCNI,
			"pod_infra_container":         cluster.SystemImages.PodInfraContainer,
			"ingress":                     cluster.SystemImages.Ingress,
			"ingress_backend":             cluster.SystemImages.IngressBackend,
			"dashboard":                   cluster.SystemImages.Dashboard,
			"heapster":                    cluster.SystemImages.Heapster,
			"grafana":                     cluster.SystemImages.Grafana,
			"influxdb":                    cluster.SystemImages.Influxdb,
			"tiller":                      cluster.SystemImages.Tiller,
		},
	})

	d.Set("ssh_key_path", cluster.SSHKeyPath)     // nolint
	d.Set("ssh_agent_auth", cluster.SSHAgentAuth) // nolint

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
		registries = append(registries, r)
	}
	d.Set("private_registries", registries) // nolint

	d.Set("ingress", []interface{}{ // nolint
		map[string]interface{}{
			"provider":      cluster.Ingress.Provider,
			"options":       cluster.Ingress.Options,
			"node_selector": cluster.Ingress.NodeSelector,
		},
	})

	d.Set("cluster_name", cluster.ClusterName) // nolint

	d.Set("cloud_provider", []interface{}{ // nolint
		map[string]interface{}{
			"name":         cluster.CloudProvider.Name,
			"cloud_config": cluster.CloudProvider.CloudConfig,
		},
	})

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
