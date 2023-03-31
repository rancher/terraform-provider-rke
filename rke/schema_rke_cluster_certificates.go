package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

//Schemas

func rkeClusterCertificatesFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"certificate": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"key": {
			Type:      schema.TypeString,
			Sensitive: true,
			Computed:  true,
		},
		"config": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"common_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ou_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"env_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"path": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"key_env_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"key_path": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"config_env_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"config_path": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
	return s
}
