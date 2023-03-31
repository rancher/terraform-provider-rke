package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	rkeClusterTaintEffectNoExecute        = "NoExecute"
	rkeClusterTaintEffectNoSchedule       = "NoSchedule"
	rkeClusterTaintEffectPreferNoSchedule = "PreferNoSchedule"
)

var (
	rkeClusterTaintEffectTypes = []string{
		rkeClusterTaintEffectNoExecute,
		rkeClusterTaintEffectNoSchedule,
		rkeClusterTaintEffectPreferNoSchedule,
	}
)

//Schemas

func rkeClusterTaintFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"value": {
			Type:     schema.TypeString,
			Required: true,
		},
		"effect": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      rkeClusterTaintEffectNoSchedule,
			ValidateFunc: validation.StringInSlice(rkeClusterTaintEffectTypes, true),
		},
	}

	return s
}
