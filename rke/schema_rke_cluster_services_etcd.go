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

func rkeClusterServicesEtcdExtraArgsArrayFields() map[string]*schema.Schema {
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
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Extra arguments that are added to the etcd services",
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
			Description: "Extra arguments that can be specified multiple times which are added to the etcd services",
			Elem: &schema.Resource{
				Schema: rkeClusterServicesEtcdExtraArgsArrayFields(),
			},
		},
		"win_extra_args_array": {
			Type:        schema.TypeSet,
			Optional:    true,
			MaxItems:    1,
			Description: "Extra arguments for Windows systems that can be specified multiple times which are added to the etcd services",
			Elem: &schema.Resource{
				Schema: rkeClusterServicesEtcdExtraArgsArrayFields(),
			},
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
