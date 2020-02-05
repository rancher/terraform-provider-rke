package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	rkeClusterNetworkPluginCalicoName  = "calico"
	rkeClusterNetworkPluginCanalName   = "canal"
	rkeClusterNetworkPluginFlannelName = "flannel"
	rkeClusterNetworkPluginNonelName   = "none"
	rkeClusterNetworkPluginWeaveName   = "weave"
)

var (
	rkeClusterNetworkPluginDefault = rkeClusterNetworkPluginCanalName
	rkeClusterNetworkPluginList    = []string{
		rkeClusterNetworkPluginCalicoName,
		rkeClusterNetworkPluginCanalName,
		rkeClusterNetworkPluginFlannelName,
		rkeClusterNetworkPluginNonelName,
		rkeClusterNetworkPluginWeaveName,
	}
)

//Schemas

func rkeClusterNetworkCalicoFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cloud_provider": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

func rkeClusterNetworkCanalFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"iface": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

func rkeClusterNetworkFlannelFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"iface": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

func rkeClusterNetworkWeaveFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"password": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
	return s
}

func rkeClusterNetworkFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"calico_network_provider": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Calico network provider config",
			Elem: &schema.Resource{
				Schema: rkeClusterNetworkCalicoFields(),
			},
		},
		"canal_network_provider": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Canal network provider config",
			Elem: &schema.Resource{
				Schema: rkeClusterNetworkCanalFields(),
			},
		},
		"flannel_network_provider": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Flannel network provider config",
			Elem: &schema.Resource{
				Schema: rkeClusterNetworkFlannelFields(),
			},
		},
		"weave_network_provider": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Weave network provider config",
			Elem: &schema.Resource{
				Schema: rkeClusterNetworkWeaveFields(),
			},
		},
		"mtu": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      0,
			Description:  "Network provider MTU",
			ValidateFunc: validation.IntBetween(0, 9000),
		},
		"options": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Network provider options",
		},
		"plugin": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      rkeClusterNetworkPluginDefault,
			Description:  "Network provider plugin",
			ValidateFunc: validation.StringInSlice(rkeClusterNetworkPluginList, true),
		},
	}
	return s
}
