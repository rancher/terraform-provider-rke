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
		"win_extra_args": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Extra arguments for Windows systems that are added to the scheduler services",
		},
		"extra_args_array": {
			Type:        schema.TypeSet,
			Optional:    true,
			MaxItems:    1,
			Description: "Extra arguments that can be added multiple times which are added to the kubelet services",
			Elem: &schema.Resource{
				Schema: rkeClusterServicesKubeletExtraArgsArrayFields(),
			},
		},
		"win_extra_args_array": {
			Type:        schema.TypeSet,
			Optional:    true,
			MaxItems:    1,
			Description: "Extra arguments for Windows systems that can be added multiple times which are added to the kubelet services",
			Elem: &schema.Resource{
				Schema: rkeClusterServicesKubeletExtraArgsArrayFields(),
			},
		},
		"extra_binds": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Extra binds added to the worker nodes",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"extra_env": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
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
			Computed:    true,
			Description: "Docker image of the kubelet service",
		},
		"infra_container_image": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The image whose network/ipc namespaces containers in each pod will use",
		},
	}
	return s
}

func rkeClusterServicesKubeletExtraArgsArrayFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"extra_arg": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"argument": {
						Required: true,
						Type:     schema.TypeString,
					},
					"values": {
						Type:     schema.TypeList,
						Required: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
	}
}
