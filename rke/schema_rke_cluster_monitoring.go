package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func rkeClusterMonitoringFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"node_selector": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Node selector key pair",
		},
		"options": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Monitoring options",
		},
		"provider": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Monitoring provider",
		},
	}
	return s
}
