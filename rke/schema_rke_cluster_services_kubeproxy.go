package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func rkeClusterServicesKubeproxyFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"extra_args": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Extra arguments that are added to the kubeproxy services",
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
			Description: "Extra arguments that can be specified multiple times which are added to the kubeproxy services",
			Elem: &schema.Resource{
				Schema: rkeClusterServicesKubeProxyExtraArgsArrayFields(),
			},
		},
		"win_extra_args_array": {
			Type:        schema.TypeSet,
			Optional:    true,
			MaxItems:    1,
			Description: "Extra arguments for Windows systems that can be specified multiple times which are added to the kubeproxy services",
			Elem: &schema.Resource{
				Schema: rkeClusterServicesKubeProxyExtraArgsArrayFields(),
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
			Description: "Extra env added to the worker nodes",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"image": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Docker image of the kubeproxy service",
		},
	}
	return s
}

func rkeClusterServicesKubeProxyExtraArgsArrayFields() map[string]*schema.Schema {
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
