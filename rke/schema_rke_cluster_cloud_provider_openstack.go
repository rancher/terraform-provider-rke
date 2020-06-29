package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	rkeClusterCloudProviderOpenstackLBMonitorDelay      = "60s"
	rkeClusterCloudProviderOpenstackLBMonitorMaxRetries = 5
	rkeClusterCloudProviderOpenstackLBMonitorTimeout    = "30s"
	rkeClusterCloudProviderOpenstackName                = "openstack"
)

//Schemas

func rkeClusterCloudProviderOpenstackBlockStorageFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"bs_version": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"ignore_volume_az": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"trust_device_path": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}

	return s
}

func rkeClusterCloudProviderOpenstackGlobalFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"auth_url": {
			Type:     schema.TypeString,
			Required: true,
		},
		"password": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"ca_file": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"domain_id": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"domain_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"region": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"tenant_id": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"tenant_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"trust_id": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"username": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"user_id": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
	}
	return s
}

func rkeClusterCloudProviderOpenstackLoadBalancerFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"create_monitor": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"floating_network_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"lb_method": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"lb_provider": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"lb_version": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"manage_security_groups": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		"monitor_delay": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  rkeClusterCloudProviderOpenstackLBMonitorDelay,
		},
		"monitor_max_retries": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  rkeClusterCloudProviderOpenstackLBMonitorMaxRetries,
		},
		"monitor_timeout": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  rkeClusterCloudProviderOpenstackLBMonitorTimeout,
		},
		"subnet_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"use_octavia": {
			Type:     schema.TypeBool,
			Optional: true,
		},
	}
	return s
}

func rkeClusterCloudProviderOpenstackMetadataFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"request_timeout": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"search_order": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	return s
}

func rkeClusterCloudProviderOpenstackRouteFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"router_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	return s
}

func rkeClusterCloudProviderOpenstackFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"global": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: rkeClusterCloudProviderOpenstackGlobalFields(),
			},
		},
		"block_storage": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: rkeClusterCloudProviderOpenstackBlockStorageFields(),
			},
		},
		"load_balancer": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: rkeClusterCloudProviderOpenstackLoadBalancerFields(),
			},
		},
		"metadata": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: rkeClusterCloudProviderOpenstackMetadataFields(),
			},
		},
		"route": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: rkeClusterCloudProviderOpenstackRouteFields(),
			},
		},
	}
	return s
}
