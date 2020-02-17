package rke

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/rke/cluster"
	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

// Flatteners

func flattenRKECluster(d *schema.ResourceData, in *cluster.Cluster) error {

	if in == nil {
		return nil
	}

	if in.AddonJobTimeout > 0 {
		d.Set("addon_job_timeout", int(in.AddonJobTimeout))
	}

	if len(in.Addons) > 0 {
		d.Set("addons", in.Addons)
	}

	if len(in.AddonsInclude) > 0 {
		d.Set("addons_include", in.AddonsInclude)
	}

	err := d.Set("authentication", flattenRKEClusterAuthentication(in.Authentication))
	if err != nil {
		return err
	}

	err = d.Set("authorization", flattenRKEClusterAuthorization(in.Authorization))
	if err != nil {
		return err
	}

	err = d.Set("bastion_host", flattenRKEClusterBastionHost(in.BastionHost))
	if err != nil {
		return err
	}

	v, ok := d.Get("cloud_provider").([]interface{})
	if !ok {
		v = []interface{}{}
	}
	err = d.Set("cloud_provider", flattenRKEClusterCloudProvider(in.CloudProvider, v))
	if err != nil {
		return err
	}

	if len(in.ClusterName) > 0 {
		d.Set("cluster_name", in.ClusterName)
	}

	if in.DNS != nil {
		err := d.Set("dns", flattenRKEClusterDNS(in.DNS))
		if err != nil {
			return err
		}
	}

	d.Set("ignore_docker_version", in.IgnoreDockerVersion)

	err = d.Set("ingress", flattenRKEClusterIngress(in.Ingress))
	if err != nil {
		return err
	}

	if len(in.Version) > 0 {
		d.Set("kubernetes_version", in.Version)
	}

	err = d.Set("monitoring", flattenRKEClusterMonitoring(in.Monitoring))
	if err != nil {
		return err
	}

	err = d.Set("network", flattenRKEClusterNetwork(in.Network))
	if err != nil {
		return err
	}

	if in.Nodes != nil {
		err := d.Set("nodes", flattenRKEClusterNodes(in.Nodes))
		if err != nil {
			return err
		}
	}

	if len(in.PrefixPath) > 0 {
		d.Set("prefix_path", in.PrefixPath)
	}

	if in.PrivateRegistries != nil {
		err := d.Set("private_registries", flattenRKEClusterPrivateRegistries(in.PrivateRegistries))
		if err != nil {
			return err
		}
	}

	err = d.Set("restore", flattenRKEClusterRestore(in.Restore))
	if err != nil {
		return err
	}

	if in.RotateCertificates != nil {
		err := d.Set("rotate_certificates", flattenRKEClusterRotateCertificates(in.RotateCertificates))
		if err != nil {
			return err
		}
	}

	v, ok = d.Get("services").([]interface{})
	if !ok {
		v = []interface{}{}
	}
	services, err := flattenRKEClusterServices(in.Services, v)
	if err != nil {
		return err
	}
	err = d.Set("services", services)
	if err != nil {
		return err
	}

	d.Set("ssh_agent_auth", in.SSHAgentAuth)

	if len(in.SSHCertPath) > 0 {
		d.Set("ssh_cert_path", in.SSHCertPath)
	}

	if len(in.SSHKeyPath) > 0 {
		d.Set("ssh_key_path", in.SSHKeyPath)
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

	return nil
}

// Expanders

func expandRKECluster(in *schema.ResourceData) (string, error) {
	if in == nil {
		return "", nil
	}

	obj := &rancher.RancherKubernetesEngineConfig{}

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

	if v, ok := in.Get("ignore_docker_version").(bool); ok {
		obj.IgnoreDockerVersion = v
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

	if v, ok := in.Get("ssh_agent_auth").(bool); ok {
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

	policyJSON := ""
	if v, ok := in.Get("services").([]interface{}); ok && len(v) > 0 {
		services, err := expandRKEClusterServices(v)
		if err != nil {
			return "", err
		}
		obj.Services = services
		if obj.Services.KubeAPI.AuditLog != nil && obj.Services.KubeAPI.AuditLog.Configuration != nil {
			policyJSON, err = interfaceToJSON(obj.Services.KubeAPI.AuditLog.Configuration.Policy)
			if err != nil {
				return "", err
			}
		}
	}

	objYml, err := interfaceToYaml(obj)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal yaml RKE cluster: %v", err)
	}

	objFixed, err := patchRKEClusterYaml(objYml, policyJSON)
	if err != nil {
		return "", fmt.Errorf("Failed to patch RKE cluster yaml: %v", err)
	}

	return objFixed, nil
}

// patchRKEClusterYaml is needed due to auditv1.Policy{} doesn't provide yaml tags
func patchRKEClusterYaml(str, policyJSON string) (string, error) {
	if len(policyJSON) == 0 {
		return str, nil
	}

	fixedPolicy := make(map[string]interface{})
	err := jsonToInterface(policyJSON, &fixedPolicy)
	if err != nil {
		return "", fmt.Errorf("ummarshalling policy json: %s", err)
	}

	out := make(map[string]interface{})
	err = ghodssyamlToInterface(str, &out)
	if err != nil {
		return "", fmt.Errorf("ummarshalling RKE cluster yaml: %s", err)
	}

	if services, ok := out["services"].(map[string]interface{}); ok {
		if kubeapi, ok := services["kube-api"].(map[string]interface{}); ok {
			if auditlog, ok := kubeapi["audit_log"].(map[string]interface{}); ok {
				if _, ok := auditlog["configuration"].(map[string]interface{}); ok {
					out["services"].(map[string]interface{})["kube-api"].(map[string]interface{})["audit_log"].(map[string]interface{})["configuration"].(map[string]interface{})["policy"] = fixedPolicy
				}
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

	// setting up the flags
	obj := cluster.GetExternalFlags(false, updateOnly, disablePortCheck, "", clusterFilePath)
	// Custom certificates and certificate dir flags
	if v, ok := in.Get("cert_dir").(string); ok && len(v) > 0 {
		obj.CertificateDir = v
	}
	obj.CustomCerts = in.Get("custom_certs").(bool)

	return obj
}
