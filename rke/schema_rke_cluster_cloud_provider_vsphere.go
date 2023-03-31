package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	rkeClusterCloudProviderVsphereName = "vsphere"
)

//Schemas

func rkeClusterCloudProviderVsphereDiskFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"scsi_controller_type": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	return s
}

func rkeClusterCloudProviderVsphereGlobalFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"datacenter": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"datacenters": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"datastore": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"insecure_flag": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"user": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"port": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"soap_roundtrip_count": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"working_dir": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"vm_uuid": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"vm_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	return s
}

func rkeClusterCloudProviderVsphereNetworkFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"public_network": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	return s
}

func rkeClusterCloudProviderVsphereVirtualCenterFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"datacenters": {
			Type:     schema.TypeString,
			Required: true,
		},
		"name": { // called server on original
			Type:     schema.TypeString,
			Required: true,
		},
		"password": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"user": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"port": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"soap_roundtrip_count": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
	return s
}

func rkeClusterCloudProviderVsphereWorkspaceFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"datacenter": {
			Type:     schema.TypeString,
			Required: true,
		},
		"server": {
			Type:     schema.TypeString,
			Required: true,
		},
		"default_datastore": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"folder": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"resourcepool_path": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	return s
}

func rkeClusterCloudProviderVsphereFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"virtual_center": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: rkeClusterCloudProviderVsphereVirtualCenterFields(),
			},
		},
		"workspace": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: rkeClusterCloudProviderVsphereWorkspaceFields(),
			},
		},
		"disk": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: rkeClusterCloudProviderVsphereDiskFields(),
			},
		},
		"global": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: rkeClusterCloudProviderVsphereGlobalFields(),
			},
		},
		"network": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: rkeClusterCloudProviderVsphereNetworkFields(),
			},
		},
	}
	return s
}
