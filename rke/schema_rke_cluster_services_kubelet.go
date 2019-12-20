package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func rkeClusterServicesKubeletFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_dns_server": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Cluster DNS service ip",
		},
		"cluster_domain": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "cluster.local",
			Description: "Domain of the cluster",
		},
		"extra_args": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Extra arguments that are added to the kubelet services",
		},
		"extra_binds": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Extra binds added to the worker nodes",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"extra_env": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Extra env added to the nodes",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"fail_swap_on": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Fail if swap is enabled",
		},
		"generate_serving_certificate": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"image": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Docker image of the kubelet service",
		},
		"infra_container_image": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The image whose network/ipc namespaces containers in each pod will use",
		},
	}
	return s
}
