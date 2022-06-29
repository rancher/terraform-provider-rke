package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//Schemas

func rkeClusterSystemImagesFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"etcd": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"alpine": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"nginx_proxy": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"cert_downloader": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"kubernetes_services_sidecar": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"kube_dns": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"dnsmasq": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"kube_dns_sidecar": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"kube_dns_autoscaler": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"coredns": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"coredns_autoscaler": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"nodelocal": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"kubernetes": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"flannel": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"flannel_cni": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"calico_node": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"calico_cni": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"calico_controllers": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"calico_ctl": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"calico_flex_vol": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"canal_node": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"canal_cni": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"canal_flannel": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"canal_flex_vol": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"weave_node": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"weave_cni": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"pod_infra_container": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"ingress": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"ingress_backend": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"metrics_server": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"windows_pod_infra_container": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"aci_cni_deploy_container": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"aci_host_container": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"aci_opflex_container": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"aci_mcast_container": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"aci_ovs_container": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"aci_controller_container": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	return s
}
