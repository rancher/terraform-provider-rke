package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//Schemas

func rkeClusterPrivateRegistriesFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"url": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Registry URL",
		},
		"is_default": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Set as default registry",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Registry password",
		},
		"user": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Registry user",
		},
	}
	return s
}
