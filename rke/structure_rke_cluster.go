package rke

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rancher/rke/cluster"
	rancher "github.com/rancher/rke/types"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiserverconfigv1 "k8s.io/apiserver/pkg/apis/config/v1"
)

// Flatteners

func flattenRKEClusterFlag(d *schema.ResourceData, in *cluster.ExternalFlags) {
	if in == nil {
		return
	}

	d.Set("update_only", in.UpdateOnly)
	d.Set("disable_port_check", in.DisablePortCheck)
	d.Set("dind", in.DinD)
	d.Set("custom_certs", in.CustomCerts)
	if len(in.CertificateDir) > 0 {
		d.Set("cert_dir", in.CertificateDir)
	}
}

func flattenRKECluster(d *schema.ResourceData, in *cluster.Cluster) error {
	if in == nil {
		return nil
	}

	var err error
	if v, ok := d.Get("addon_job_timeout").(int); ok && v > 0 && in.AddonJobTimeout > 0 {
		d.Set("addon_job_timeout", in.AddonJobTimeout)
	}

	if v, ok := d.Get("addons").(string); ok && len(v) > 0 && len(in.Addons) > 0 {
		d.Set("addons", in.Addons)
	}

	if v, ok := d.Get("addons_include").([]interface{}); ok && len(v) > 0 && len(in.AddonsInclude) > 0 {
		err = d.Set("addons_include", toArrayInterface(in.AddonsInclude))
		if err != nil {
			return err
		}
	}

	if v, ok := d.Get("authentication").([]interface{}); ok && len(v) > 0 {
		err = d.Set("authentication", flattenRKEClusterAuthentication(in.Authentication))
		if err != nil {
			return err
		}
	}

	if v, ok := d.Get("authorization").([]interface{}); ok && len(v) > 0 {
		err = d.Set("authorization", flattenRKEClusterAuthorization(in.Authorization))
		if err != nil {
			return err
		}
	}

	if v, ok := d.Get("bastion_host").([]interface{}); ok && len(v) > 0 {
		err = d.Set("bastion_host", flattenRKEClusterBastionHost(in.BastionHost))
		if err != nil {
			return err
		}
	}

	if v, ok := d.Get("cloud_provider").([]interface{}); ok && len(v) > 0 {
		err = d.Set("cloud_provider", flattenRKEClusterCloudProvider(in.CloudProvider, v))
		if err != nil {
			return err
		}
	}

	if v, ok := d.Get("cluster_name").(string); ok && len(v) > 0 && len(in.ClusterName) > 0 {
		d.Set("cluster_name", in.ClusterName)
	}

	if v, ok := d.Get("dns").([]interface{}); ok && len(v) > 0 && in.DNS != nil {
		err := d.Set("dns", flattenRKEClusterDNS(in.DNS))
		if err != nil {
			return err
		}
	}

	d.Set("dind", in.DinD)

	if _, ok := d.Get("enable_cri_dockerd").(bool); ok && in.EnableCRIDockerd != nil {
		d.Set("enable_cri_dockerd", *in.EnableCRIDockerd)
	}

	if _, ok := d.Get("ignore_docker_version").(bool); ok && in.IgnoreDockerVersion != nil {
		d.Set("ignore_docker_version", *in.IgnoreDockerVersion)
	}

	if v, ok := d.Get("ingress").([]interface{}); ok && len(v) > 0 {
		err = d.Set("ingress", flattenRKEClusterIngress(in.Ingress))
		if err != nil {
			return err
		}
	}

	if len(in.Version) > 0 {
		d.Set("kubernetes_version", in.Version)
	}

	if v, ok := d.Get("monitoring").([]interface{}); ok && len(v) > 0 {
		err = d.Set("monitoring", flattenRKEClusterMonitoring(in.Monitoring))
		if err != nil {
			return err
		}
	}

	if v, ok := d.Get("network").([]interface{}); ok && len(v) > 0 {
		err = d.Set("network", flattenRKEClusterNetwork(in.Network))
		if err != nil {
			return err
		}
	}

	if v, ok := d.Get("nodes").([]interface{}); ok && len(v) > 0 && in.Nodes != nil && !in.DinD {
		nodes := flattenRKEClusterNodes(in.Nodes, v)
		err := d.Set("nodes", nodes)
		if err != nil {
			return err
		}
	}

	if v, ok := d.Get("prefix_path").(string); ok && len(v) > 0 && len(in.PrefixPath) > 0 {
		d.Set("prefix_path", in.PrefixPath)
	}

	if v, ok := d.Get("private_registries").([]interface{}); ok && len(v) > 0 && in.PrivateRegistries != nil {
		err := d.Set("private_registries", flattenRKEClusterPrivateRegistries(in.PrivateRegistries))
		if err != nil {
			return err
		}
	}

	if v, ok := d.Get("restore").([]interface{}); ok && len(v) > 0 {
		err = d.Set("restore", flattenRKEClusterRestore(in.Restore))
		if err != nil {
			return err
		}
	}

	if v, ok := d.Get("rotate_certificates").([]interface{}); in.RotateCertificates != nil && ok && len(v) > 0 {
		err := d.Set("rotate_certificates", flattenRKEClusterRotateCertificates(in.RotateCertificates))
		if err != nil {
			return err
		}
	}

	if v, ok := d.Get("services").([]interface{}); ok && len(v) > 0 {
		services, err := flattenRKEClusterServices(in.Services, v)
		if err != nil {
			return err
		}
		err = d.Set("services", services)
		if err != nil {
			return err
		}
	}

	if _, ok := d.Get("ssh_agent_auth").(bool); ok {
		d.Set("ssh_agent_auth", in.SSHAgentAuth)
	}

	if v, ok := d.Get("ssh_cert_path").(string); ok && len(v) > 0 && len(in.SSHCertPath) > 0 {
		d.Set("ssh_cert_path", in.SSHCertPath)
	}

	if v, ok := d.Get("ssh_key_path").(string); ok && len(v) > 0 && len(in.SSHKeyPath) > 0 {
		d.Set("ssh_key_path", in.SSHKeyPath)
	}

	if v, ok := d.Get("system_images").([]interface{}); ok && len(v) > 0 {
		d.Set("system_images", in.SystemImages)
	}

	if v, ok := d.Get("upgrade_strategy").([]interface{}); ok && len(v) > 0 && in.UpgradeStrategy != nil {
		err = d.Set("upgrade_strategy", flattenRKEClusterNodeUpgradeStrategy(in.UpgradeStrategy))
		if err != nil {
			return err
		}
	}

	// computed values
	d.Set("api_server_url", "") // nolint
	if in.ControlPlaneHosts != nil && len(in.ControlPlaneHosts) > 0 {
		apiServerURL := fmt.Sprintf("https://" + in.ControlPlaneHosts[0].Address + ":6443")
		d.Set("api_server_url", apiServerURL)
	}

	caCrt, clientCrt, clientKey, certificates := flattenRKEClusterCertificates(in.Certificates)
	d.Set("ca_crt", caCrt)          // nolint
	d.Set("client_cert", clientCrt) // nolint
	d.Set("client_key", clientKey)  // nolint
	d.Set("certificates", certificates)
	d.Set("kube_admin_user", rkeClusterCertificatesKubeAdminCertName)
	d.Set("cluster_domain", in.ClusterDomain)        // nolint
	d.Set("cluster_cidr", in.ClusterCIDR)            // nolint
	d.Set("cluster_dns_server", in.ClusterDNSServer) // nolint

	err = d.Set("etcd_hosts", flattenRKEClusterNodesComputed(in.EtcdHosts))
	if err != nil {
		return err
	}

	err = d.Set("control_plane_hosts", flattenRKEClusterNodesComputed(in.ControlPlaneHosts))
	if err != nil {
		return err
	}

	err = d.Set("worker_hosts", flattenRKEClusterNodesComputed(in.WorkerHosts))
	if err != nil {
		return err
	}

	err = d.Set("inactive_hosts", flattenRKEClusterNodesComputed(in.InactiveHosts))
	if err != nil {
		return err
	}

	err = d.Set("running_system_images", flattenRKEClusterSystemImages(in.SystemImages))
	if err != nil {
		return err
	}

	rkeClusterYaml, _, err := expandRKECluster(d)
	err = d.Set("rke_cluster_yaml", rkeClusterYaml)
	if err != nil {
		return err
	}

	return nil
}

// Expanders

func expandRKECluster(in *schema.ResourceData) (string, *rancher.RancherKubernetesEngineConfig, error) {
	if in == nil {
		return "", nil, nil
	}

	obj := &rancher.RancherKubernetesEngineConfig{}

	if v, ok := in.Get("cluster_yaml").(string); ok && len(v) > 0 {
		var err error
		obj, err = cluster.ParseConfig(v)
		if err != nil {
			return "", nil, err
		}
	}

	if v, ok := in.Get("addon_job_timeout").(int); ok && v > 0 {
		obj.AddonJobTimeout = v
	}

	if v, ok := in.Get("addons").(string); ok && len(v) > 0 {
		obj.Addons = v
	}

	if v, ok := in.Get("addons_include").([]interface{}); ok && len(v) > 0 {
		obj.AddonsInclude = toArrayString(v)
	}

	if v, ok := in.Get("authentication").([]interface{}); ok && len(v) > 0 {
		obj.Authentication = expandRKEClusterAuthentication(v)
	}

	if v, ok := in.Get("authorization").([]interface{}); ok && len(v) > 0 {
		obj.Authorization = expandRKEClusterAuthorization(v)
	}

	if v, ok := in.Get("bastion_host").([]interface{}); ok && len(v) > 0 {
		obj.BastionHost = expandRKEClusterBastionHost(v)
	}

	if v, ok := in.Get("cloud_provider").([]interface{}); ok && len(v) > 0 {
		obj.CloudProvider = expandRKEClusterCloudProvider(v)
	}

	if v, ok := in.Get("cluster_name").(string); ok && len(v) > 0 {
		obj.ClusterName = v
	}

	if v, ok := in.Get("dns").([]interface{}); ok && len(v) > 0 {
		obj.DNS = expandRKEClusterDNS(v)
	}

	if v, ok := in.Get("enable_cri_dockerd").(bool); ok && obj.EnableCRIDockerd == nil {
		obj.EnableCRIDockerd = &v
	}

	if v, ok := in.Get("ignore_docker_version").(bool); ok && obj.IgnoreDockerVersion == nil {
		obj.IgnoreDockerVersion = &v
	}

	if v, ok := in.Get("ingress").([]interface{}); ok && len(v) > 0 {
		obj.Ingress = expandRKEClusterIngress(v)
	}

	if v, ok := in.Get("kubernetes_version").(string); ok && len(v) > 0 {
		obj.Version = v
	}

	if v, ok := in.Get("monitoring").([]interface{}); ok && len(v) > 0 {
		obj.Monitoring = expandRKEClusterMonitoring(v)
	}

	if v, ok := in.Get("network").([]interface{}); ok && len(v) > 0 {
		obj.Network = expandRKEClusterNetwork(v)
	}

	if v, ok := in.Get("nodes").([]interface{}); ok && len(v) > 0 {
		obj.Nodes = expandRKEClusterNodes(v)
	}

	if v, ok := in.Get("prefix_path").(string); ok && len(v) > 0 {
		obj.PrefixPath = v
	}

	if v, ok := in.Get("private_registries").([]interface{}); ok && len(v) > 0 {
		obj.PrivateRegistries = expandRKEClusterPrivateRegistries(v)
	}

	if v, ok := in.Get("restore").([]interface{}); ok && len(v) > 0 {
		obj.Restore = expandRKEClusterRestore(v)
	}

	if v, ok := in.Get("rotate_certificates").([]interface{}); ok && len(v) > 0 {
		obj.RotateCertificates = expandRKEClusterRotateCertificates(v)
	}

	if v, ok := in.Get("ssh_agent_auth").(bool); ok && v {
		obj.SSHAgentAuth = v
	}

	if v, ok := in.Get("ssh_cert_path").(string); ok && len(v) > 0 {
		obj.SSHCertPath = v
	}

	if v, ok := in.Get("ssh_key_path").(string); ok && len(v) > 0 {
		obj.SSHKeyPath = v
	}

	if v, ok := in.Get("system_images").([]interface{}); ok && len(v) > 0 {
		obj.SystemImages = expandRKEClusterSystemImages(v)
	}

	if v, ok := in.Get("upgrade_strategy").([]interface{}); ok {
		obj.UpgradeStrategy = expandRKEClusterNodeUpgradeStrategy(v)
	}

	if v, ok := in.Get("services").([]interface{}); ok && len(v) > 0 {
		services, err := expandRKEClusterServices(v)
		if err != nil {
			return "", nil, err
		}
		obj.Services = services
	}

	if v, ok := in.Get("dind").(bool); ok && v {
		if obj.Services.Kubeproxy.ExtraArgs == nil {
			obj.Services.Kubeproxy.ExtraArgs = make(map[string]string)
		}
		obj.Services.Kubeproxy.ExtraArgs["conntrack-max-per-core"] = "0"
	}

	if k8sVersionRequiresCri(obj.Version) && obj.EnableCRIDockerd != nil && !*obj.EnableCRIDockerd {
		return "", nil, fmt.Errorf("kubernetes version %s requires enable_cri_dockerd to be set to true", obj.Version)
	}

	objYml, err := patchRKEClusterYaml(obj)
	if err != nil {
		return "", nil, fmt.Errorf("Failed to patch RKE cluster yaml: %v", err)
	}

	return objYml, obj, nil
}

// patchRKEClusterYaml is needed due to auditv1.Policy{} doesn't provide yaml tags
func patchRKEClusterYaml(in *rancher.RancherKubernetesEngineConfig) (string, error) {
	outFixed := make(map[string]interface{})
	if in.Services.KubeAPI.AuditLog != nil && in.Services.KubeAPI.AuditLog.Configuration != nil {
		inJSON, err := interfaceToJSON(in.Services.KubeAPI.AuditLog.Configuration.Policy)
		if err != nil {
			return "", err
		}
		if len(inJSON) > 0 {
			outFixed["audit_log"], err = jsonToMapInterface(inJSON)
			if err != nil {
				return "", fmt.Errorf("unmarshalling auditlog json: %s", err)
			}
		}
	}
	if in.Services.KubeAPI.EventRateLimit != nil && in.Services.KubeAPI.EventRateLimit.Configuration != nil {
		if len(in.Services.KubeAPI.EventRateLimit.Configuration.TypeMeta.Kind) == 0 {
			in.Services.KubeAPI.EventRateLimit.Configuration.TypeMeta.Kind = clusterServicesKubeAPIEventRateLimitConfigKindDefault

		}
		if len(in.Services.KubeAPI.EventRateLimit.Configuration.TypeMeta.APIVersion) == 0 {
			in.Services.KubeAPI.EventRateLimit.Configuration.TypeMeta.APIVersion = clusterServicesKubeAPIEventRateLimitConfigAPIDefault
		}
		inJSON, err := interfaceToJSON(in.Services.KubeAPI.EventRateLimit.Configuration)
		if err != nil {
			return "", err
		}
		if len(inJSON) > 0 {
			outFixed["event_rate_limit"], err = jsonToMapInterface(inJSON)
			if err != nil {
				return "", fmt.Errorf("unmarshalling event_rate_limit json: %s", err)
			}
		}
	}
	if in.Services.KubeAPI.SecretsEncryptionConfig != nil && in.Services.KubeAPI.SecretsEncryptionConfig.CustomConfig != nil {
		customConfigV1Str, err := interfaceToGhodssyaml(in.Services.KubeAPI.SecretsEncryptionConfig.CustomConfig)
		if err != nil {
			return "", fmt.Errorf("marshalling custom_config yaml: %v", err)
		}
		customConfigV1 := &apiserverconfigv1.EncryptionConfiguration{
			TypeMeta: metav1.TypeMeta{
				Kind:       clusterServicesKubeAPISecretsEncryptionConfigKindDefault,
				APIVersion: clusterServicesKubeAPISecretsEncryptionConfigAPIDefault,
			},
		}
		err = ghodssyamlToInterface(customConfigV1Str, customConfigV1)
		if err != nil {
			return "", fmt.Errorf("unmarshalling custom_config yaml: %v", err)
		}
		inJSON, err := interfaceToJSON(customConfigV1)
		if err != nil {
			return "", err
		}
		if len(inJSON) > 0 {
			outFixed["secrets_encryption_config"], err = jsonToMapInterface(inJSON)
			if err != nil {
				return "", fmt.Errorf("unmarshalling eventrate json: %s", err)
			}
		}
	}

	outYml, err := interfaceToYaml(in)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal yaml RKE cluster: %v", err)
	}

	if len(outFixed) == 0 {
		return outYml, nil
	}

	out := make(map[string]interface{})
	err = ghodssyamlToInterface(outYml, &out)
	if err != nil {
		return "", fmt.Errorf("unmarshalling RKE cluster yaml: %s", err)
	}

	if services, ok := out["services"].(map[string]interface{}); ok {
		if kubeapi, ok := services["kube-api"].(map[string]interface{}); ok {
			if auditlog, ok := kubeapi["audit_log"].(map[string]interface{}); ok && outFixed["audit_log"] != nil {
				if _, ok := auditlog["configuration"].(map[string]interface{}); ok {
					out["services"].(map[string]interface{})["kube-api"].(map[string]interface{})["audit_log"].(map[string]interface{})["configuration"].(map[string]interface{})["policy"] = outFixed["audit_log"]
				}
			}
			if _, ok := kubeapi["event_rate_limit"].(map[string]interface{}); ok && outFixed["event_rate_limit"] != nil {
				out["services"].(map[string]interface{})["kube-api"].(map[string]interface{})["event_rate_limit"].(map[string]interface{})["configuration"] = outFixed["event_rate_limit"]
			}
			if _, ok := kubeapi["secrets_encryption_config"].(map[string]interface{}); ok && outFixed["secrets_encryption_config"] != nil {
				out["services"].(map[string]interface{})["kube-api"].(map[string]interface{})["secrets_encryption_config"].(map[string]interface{})["custom_config"] = outFixed["secrets_encryption_config"]
			}
		}
	}

	outYaml, err := interfaceToGhodssyaml(out)
	if err != nil {
		return "", fmt.Errorf("marshalling RKE cluster patched yaml: %s", err)
	}

	return outYaml, nil
}

func expandRKEClusterFlag(in *schema.ResourceData, clusterFilePath string) cluster.ExternalFlags {
	if in == nil {
		return cluster.ExternalFlags{}
	}

	updateOnly := in.Get("update_only").(bool)
	disablePortCheck := in.Get("disable_port_check").(bool)
	dind := in.Get("dind").(bool)

	if dind {
		updateOnly = false
	}

	// setting up the flags
	obj := cluster.GetExternalFlags(false, updateOnly, disablePortCheck, false, "", clusterFilePath)
	obj.DinD = dind
	if !dind {
		// Custom certificates and certificate dir flags
		if v, ok := in.Get("cert_dir").(string); ok && len(v) > 0 {
			obj.CertificateDir = v
		}
		obj.CustomCerts = in.Get("custom_certs").(bool)
	}

	return obj
}

func k8sVersionRequiresCri(kubernetesVersion string) bool {
	version, err := getClusterVersion(kubernetesVersion)
	if err != nil {
		// This debug / error is not supposed to happen, the kubernetesVersion should be validated by the provider.
		log.Debugf("Unable to get the semantic version for kubernetesVersion, value: %s", kubernetesVersion)
		return false
	}
	return parsedRangeAtLeast124(version)
}
