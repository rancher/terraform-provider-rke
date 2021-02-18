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
	rkeClusterNetworkPluginAciName     = "aci"
)

var (
	rkeClusterNetworkPluginDefault = rkeClusterNetworkPluginCanalName
	rkeClusterNetworkPluginList    = []string{
		rkeClusterNetworkPluginCalicoName,
		rkeClusterNetworkPluginCanalName,
		rkeClusterNetworkPluginFlannelName,
		rkeClusterNetworkPluginNonelName,
		rkeClusterNetworkPluginWeaveName,
		rkeClusterNetworkPluginAciName,
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

func rkeClusterNetworkAciFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"system_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"apic_hosts": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"token": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"apic_user_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"apic_user_key": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"apic_user_crt": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"encap_type": {
			Type:     schema.TypeString,
			Required: true,
		},
		"mcast_range_start": {
			Type:     schema.TypeString,
			Required: true,
		},
		"mcast_range_end": {
			Type:     schema.TypeString,
			Required: true,
		},
		"aep": {
			Type:     schema.TypeString,
			Required: true,
		},
		"vrf_name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"vrf_tenant": {
			Type:     schema.TypeString,
			Required: true,
		},
		"l3out": {
			Type:     schema.TypeString,
			Required: true,
		},
		"node_subnet": {
			Type:     schema.TypeString,
			Required: true,
		},
		"l3out_external_networks": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"extern_dynamic": {
			Type:     schema.TypeString,
			Required: true,
		},
		"extern_static": {
			Type:     schema.TypeString,
			Required: true,
		},
		"node_svc_subnet": {
			Type:     schema.TypeString,
			Required: true,
		},
		"kube_api_vlan": {
			Type:     schema.TypeString,
			Required: true,
		},
		"service_vlan": {
			Type:     schema.TypeString,
			Required: true,
		},
		"infra_vlan": {
			Type:     schema.TypeString,
			Required: true,
		},
		"snat_port_range_start": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"snat_port_range_end": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"snat_ports_per_node": {
			Type:     schema.TypeString,
			Optional: true,
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
		"aci_network_provider": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Aci network provider config",
			Elem: &schema.Resource{
				Schema: rkeClusterNetworkAciFields(),
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
