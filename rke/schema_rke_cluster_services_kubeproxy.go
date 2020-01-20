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
