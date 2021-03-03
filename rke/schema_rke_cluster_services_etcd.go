package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func rkeClusterServicesEtcdBackupConfigS3Fields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"access_key": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"bucket_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"custom_ca": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"endpoint": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"folder": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"region": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"secret_key": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
	}
	return s
}

func rkeClusterServicesEtcdBackupConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"interval_hours": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  12,
		},
		"retention": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  6,
		},
		"s3_backup_config": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: rkeClusterServicesEtcdBackupConfigS3Fields(),
			},
		},
		"safe_timestamp": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"timeout": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  300,
		},
	}
	return s
}

func rkeClusterServicesEtcdFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"backup_config": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: rkeClusterServicesEtcdBackupConfigFields(),
			},
		},
		"ca_cert": {
			Type:      schema.TypeString,
			Optional:  true,
			Computed:  true,
			Sensitive: true,
		},
		"cert": {
			Type:      schema.TypeString,
			Optional:  true,
			Computed:  true,
			Sensitive: true,
		},
		"creation": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"external_urls": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"extra_args": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"extra_binds": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"extra_env": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"gid": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
		"image": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"key": {
			Type:      schema.TypeString,
			Optional:  true,
			Computed:  true,
			Sensitive: true,
		},
		"path": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"retention": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"snapshot": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"uid": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  0,
		},
	}
	return s
}
