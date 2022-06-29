package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//Schemas

func rkeClusterRestoreFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"restore": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Restore RKE cluster",
		},
		"snapshot_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Snapshot name",
		},
	}
	return s
}
