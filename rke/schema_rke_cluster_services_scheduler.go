package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//Schemas

func rkeClusterServicesSchedulerFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"extra_args": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Extra arguments that are added to the scheduler services",
		},
		"extra_args_array": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Extra arguments array that is added to the scheduler services",
		},
		"extra_binds": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Extra binds added to the controlplane nodes",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"extra_env": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Extra env added to the controlplane nodes",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"image": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Docker image of the scheduler service",
		},
	}
	return s
}
