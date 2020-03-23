package rke

import (
	"sort"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/rancher/kontainer-driver-metadata/rke"
)

//Schemas

func rkeClusterFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"cert_dir": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specify a certificate dir path",
		},
		"custom_certs": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Use custom certificates from a cert dir",
		},
		"dind": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "RKE k8s cluster dind (experimental)",
		},
		"dind_storage_driver": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "RKE k8s cluster dind storage driver (experimental)",
		},
		"dind_dns_server": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "RKE k8s cluster dind storage driver (experimental)",
		},
		"delay_on_creation": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "RKE k8s cluster delay on creation",
			ValidateFunc: validation.IntAtLeast(0),
		},
		"disable_port_check": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable/Disable RKE k8s cluster port checking",
		},
		"addon_job_timeout": {
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
			Description:  "RKE k8s cluster addon deployment timeout in seconds for status check",
		},
		"addons": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "RKE k8s cluster user addons YAML manifest to be deployed",
		},
		"addons_include": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "RKE k8s cluster user addons YAML manifest urls or paths to be deployed",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"authentication": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "RKE k8s cluster authentication configuration",
			Elem: &schema.Resource{
				Schema: rkeClusterAuthenticationFields(),
			},
		},
		"authorization": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "RKE k8s cluster authorization mode configuration",
			Elem: &schema.Resource{
				Schema: rkeClusterAuthorizationFields(),
			},
		},
		"bastion_host": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "RKE k8s cluster bastion Host configuration",
			Elem: &schema.Resource{
				Schema: rkeClusterBastionHostFields(),
			},
		},
		"cloud_provider": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "RKE k8s cluster cloud provider configuration",
			Elem: &schema.Resource{
				Schema: rkeClusterCloudProviderFields(),
			},
		},
		"cluster_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "RKE k8s cluster name used in the kube config",
		},
		"dns": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "RKE k8s cluster DNS Config",
			Elem: &schema.Resource{
				Schema: rkeClusterDNSFields(),
			},
		},
		"ignore_docker_version": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Enable/Disable RKE k8s cluster strict docker version checking",
		},
		"ingress": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "RKE k8s cluster ingress controller configuration",
			Elem: &schema.Resource{
				Schema: rkeClusterIngressFields(),
			},
		},
		"kubernetes_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "K8s version to deploy (if kubernetes image is specified, image version takes precedence)",
			ValidateFunc: validation.StringInSlice(func() []string {
				rkeData := rke.DriverData
				versions := make([]*version.Version, 0, len(rkeData.K8sVersionRKESystemImages))
				for k := range rkeData.K8sVersionRKESystemImages {
					v, _ := version.NewVersion(k)
					versions = append(versions, v)
				}
				sort.Sort(sort.Reverse(version.Collection(versions)))
				keys := make([]string, len(versions))
				for i := range versions {
					keys[i] = "v" + versions[i].String()
				}
				return keys
			}(), false),
		},
		"monitoring": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "RKE k8s cluster monitoring Config",
			Elem: &schema.Resource{
				Schema: rkeClusterMonitoringFields(),
			},
		},
		"network": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "RKE k8s cluster network configuration",
			Elem: &schema.Resource{
				Schema: rkeClusterNetworkFields(),
			},
		},
		"nodes_conf": {
			Type:        schema.TypeList,
			MinItems:    1,
			Optional:    true,
			Sensitive:   true,
			Description: "RKE k8s cluster nodes (YAML | JSON)",
			Elem: &schema.Schema{
				Type:      schema.TypeString,
				Sensitive: true,
			},
			ConflictsWith: []string{"nodes"},
		},
		"nodes": {
			Type:        schema.TypeList,
			MinItems:    1,
			Optional:    true,
			Description: "RKE k8s cluster nodes",
			Elem: &schema.Resource{
				Schema: rkeClusterNodeFields(),
			},
			ConflictsWith: []string{"nodes_conf"},
		},
		"prefix_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "RKE k8s directory path",
		},
		"private_registries": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "RKE k8s cluster private docker registries",
			Elem: &schema.Resource{
				Schema: rkeClusterPrivateRegistriesFields(),
			},
		},
		"restore": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "RKE k8s cluster restore configuration",
			Elem: &schema.Resource{
				Schema: rkeClusterRestoreFields(),
			},
		},
		"rotate_certificates": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "RKE k8s cluster rotate certificates configuration",
			Elem: &schema.Resource{
				Schema: rkeClusterRotateCertificatesFields(),
			},
		},
		"services": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "RKE k8s cluster services",
			Elem: &schema.Resource{
				Schema: rkeClusterServicesFields(),
			},
		},
		"services_etcd": {
			Type:       schema.TypeList,
			MaxItems:   1,
			Optional:   true,
			Deprecated: "Use services.etcd instead",
			Elem: &schema.Resource{
				Schema: rkeClusterServicesEtcdFields(),
			},
		},
		"services_kube_api": {
			Type:       schema.TypeList,
			MaxItems:   1,
			Optional:   true,
			Deprecated: "Use services.kube_api instead",
			Elem: &schema.Resource{
				Schema: rkeClusterServicesKubeAPIFields(),
			},
		},
		"services_kube_controller": {
			Type:       schema.TypeList,
			MaxItems:   1,
			Optional:   true,
			Deprecated: "Use services.kube_controller instead",
			Elem: &schema.Resource{
				Schema: rkeClusterServicesKubeControllerFields(),
			},
		},
		"services_kubelet": {
			Type:       schema.TypeList,
			MaxItems:   1,
			Optional:   true,
			Deprecated: "Use services.kubelet instead",
			Elem: &schema.Resource{
				Schema: rkeClusterServicesKubeletFields(),
			},
		},
		"services_kubeproxy": {
			Type:       schema.TypeList,
			MaxItems:   1,
			Optional:   true,
			Deprecated: "Use services.kubeproxy instead",
			Elem: &schema.Resource{
				Schema: rkeClusterServicesKubeproxyFields(),
			},
		},
		"services_scheduler": {
			Type:       schema.TypeList,
			MaxItems:   1,
			Optional:   true,
			Deprecated: "Use services.scheduler instead",
			Elem: &schema.Resource{
				Schema: rkeClusterServicesSchedulerFields(),
			},
		},
		"ssh_agent_auth": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "SSH Agent Auth enable",
		},
		"ssh_cert_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "SSH Certificate Path",
		},
		"ssh_key_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "SSH Private Key Path",
		},
		"system_images": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "RKE k8s cluster system images list",
			Elem: &schema.Resource{
				Schema: rkeClusterSystemImagesFields(),
			},
		},
		"update_only": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Skip idempotent deployment of control and etcd plane",
		},
		// Computed fields
		"ca_crt": {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "RKE k8s cluster CA certificate",
		},
		"client_cert": {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "RKE k8s cluster client certificate",
		},
		"client_key": {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "RKE k8s cluster client key",
		},
		"rke_state": {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "RKE k8s cluster state",
		},
		"kube_config_yaml": {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "RKE k8s cluster kube config yaml",
		},
		"internal_kube_config_yaml": {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "RKE k8s cluster internal kube config yaml",
		},
		"rke_cluster_yaml": {
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
			Description: "RKE k8s cluster config yaml",
		},
		"certificates": {
			Type:        schema.TypeList,
			Computed:    true,
			Sensitive:   true,
			Description: "RKE k8s cluster certificates",
			Elem: &schema.Resource{
				Schema: rkeClusterCertificatesFields(),
			},
		},
		"kube_admin_user": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "RKE k8s cluster admin user",
		},
		"api_server_url": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "RKE k8s cluster api server url",
		},
		"cluster_domain": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "RKE k8s cluster domain",
		},
		"cluster_cidr": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "RKE k8s cluster cidr",
		},
		"cluster_dns_server": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "RKE k8s cluster dns server",
		},
		"control_plane_hosts": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "RKE k8s cluster control plane nodes",
			Elem: &schema.Resource{
				Schema: rkeClusterNodeComputedFields(),
			},
		},
		"etcd_hosts": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "RKE k8s cluster etcd nodes",
			Elem: &schema.Resource{
				Schema: rkeClusterNodeComputedFields(),
			},
		},
		"inactive_hosts": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "RKE k8s cluster inactive nodes",
			Elem: &schema.Resource{
				Schema: rkeClusterNodeComputedFields(),
			},
		},
		"worker_hosts": {
			Type:        schema.TypeList,
			Computed:    true,
			Description: "RKE k8s cluster worker nodes",
			Elem: &schema.Resource{
				Schema: rkeClusterNodeComputedFields(),
			},
		},
		"running_system_images": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Computed:    true,
			Description: "RKE k8s cluster running system images list",
			Elem: &schema.Resource{
				Schema: rkeClusterSystemImagesFields(),
			},
		},
	}
}
