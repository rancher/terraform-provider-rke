package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const authorizationModeRBAC = "rbac"

var (
	authorizationModeList = []string{authorizationModeRBAC, "none"}
)

//Schemas

func rkeClusterAuthorizationFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"mode": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      authorizationModeRBAC,
			ValidateFunc: validation.StringInSlice(authorizationModeList, true),
		},
		"options": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Authorization mode options",
		},
	}
	return s
}
