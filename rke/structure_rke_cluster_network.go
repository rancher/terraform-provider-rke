package rke

import (
	rancher "github.com/rancher/rke/types"
)

// Flatteners

func flattenRKEClusterNetworkCalico(in *rancher.CalicoNetworkProvider) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.CloudProvider) > 0 {
		obj["cloud_provider"] = in.CloudProvider
	}

	return []interface{}{obj}
}

func flattenRKEClusterNetworkCanal(in *rancher.CanalNetworkProvider) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.FlannelNetworkProvider.Iface) > 0 {
		obj["iface"] = in.FlannelNetworkProvider.Iface
	}

	return []interface{}{obj}
}

func flattenRKEClusterNetworkFlannel(in *rancher.FlannelNetworkProvider) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.Iface) > 0 {
		obj["iface"] = in.Iface
	}

	return []interface{}{obj}
}

func flattenRKEClusterNetworkWeave(in *rancher.WeaveNetworkProvider) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.Password) > 0 {
		obj["password"] = in.Password
	}

	return []interface{}{obj}
}

func flattenRKEClusterNetworkAci(in *rancher.AciNetworkProvider) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}
	if len(in.SystemIdentifier) > 0 {
		obj["system_id"] = in.SystemIdentifier
	}
	if in.ApicHosts != nil {
		obj["apic_hosts"] = toArrayInterface(in.ApicHosts)
	}
	if len(in.Token) > 0 {
		obj["token"] = in.Token
	}
	if len(in.ApicUserName) > 0 {
		obj["apic_user_name"] = in.ApicUserName
	}
	if len(in.ApicUserKey) > 0 {
		obj["apic_user_key"] = in.ApicUserKey
	}
	if len(in.ApicUserCrt) > 0 {
		obj["apic_user_crt"] = in.ApicUserCrt
	}
	if len(in.EncapType) > 0 {
		obj["encap_type"] = in.EncapType
	}
	if len(in.McastRangeStart) > 0 {
		obj["mcast_range_start"] = in.McastRangeStart
	}
	if len(in.McastRangeEnd) > 0 {
		obj["mcast_range_end"] = in.McastRangeEnd
	}
	if len(in.AEP) > 0 {
		obj["aep"] = in.AEP
	}
	if len(in.VRFName) > 0 {
		obj["vrf_name"] = in.VRFName
	}
	if len(in.VRFTenant) > 0 {
		obj["vrf_tenant"] = in.VRFTenant
	}
	if len(in.L3Out) > 0 {
		obj["l3out"] = in.L3Out
	}
	if len(in.NodeSubnet) > 0 {
		obj["node_subnet"] = in.NodeSubnet
	}
	if in.L3OutExternalNetworks != nil {
		obj["l3out_external_networks"] = toArrayInterface(in.L3OutExternalNetworks)
	}
	if len(in.DynamicExternalSubnet) > 0 {
		obj["extern_dynamic"] = in.DynamicExternalSubnet
	}
	if len(in.StaticExternalSubnet) > 0 {
		obj["extern_static"] = in.StaticExternalSubnet
	}
	if len(in.ServiceGraphSubnet) > 0 {
		obj["node_svc_subnet"] = in.ServiceGraphSubnet
	}
	if len(in.KubeAPIVlan) > 0 {
		obj["kube_api_vlan"] = in.KubeAPIVlan
	}
	if len(in.ServiceVlan) > 0 {
		obj["service_vlan"] = in.ServiceVlan
	}
	if len(in.InfraVlan) > 0 {
		obj["infra_vlan"] = in.InfraVlan
	}
	if len(in.SnatPortRangeStart) > 0 {
		obj["snat_port_range_start"] = in.SnatPortRangeStart
	}
	if len(in.SnatPortRangeEnd) > 0 {
		obj["snat_port_range_end"] = in.SnatPortRangeEnd
	}
	if len(in.SnatPortsPerNode) > 0 {
		obj["snat_ports_per_node"] = in.SnatPortsPerNode
	}

	return []interface{}{obj}
}

func flattenRKEClusterNetwork(in rancher.NetworkConfig) []interface{} {
	obj := make(map[string]interface{})

	if in.CalicoNetworkProvider != nil {
		obj["calico_network_provider"] = flattenRKEClusterNetworkCalico(in.CalicoNetworkProvider)
	}

	if in.CanalNetworkProvider != nil {
		obj["canal_network_provider"] = flattenRKEClusterNetworkCanal(in.CanalNetworkProvider)
	}

	if in.FlannelNetworkProvider != nil {
		obj["flannel_network_provider"] = flattenRKEClusterNetworkFlannel(in.FlannelNetworkProvider)
	}

	if in.WeaveNetworkProvider != nil {
		obj["weave_network_provider"] = flattenRKEClusterNetworkWeave(in.WeaveNetworkProvider)
	}

	if in.AciNetworkProvider != nil {
		obj["aci_network_provider"] = flattenRKEClusterNetworkAci(in.AciNetworkProvider)
	}

	if in.MTU > 0 {
		obj["mtu"] = in.MTU
	}

	if len(in.Options) > 0 {
		obj["options"] = toMapInterface(in.Options)
	}

	if len(in.Plugin) > 0 {
		obj["plugin"] = in.Plugin
	}

	return []interface{}{obj}
}

// Expanders

func expandRKEClusterNetworkCalico(p []interface{}) *rancher.CalicoNetworkProvider {
	obj := &rancher.CalicoNetworkProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["cloud_provider"].(string); ok && len(v) > 0 {
		obj.CloudProvider = v
	}

	return obj
}

func expandRKEClusterNetworkCanal(p []interface{}) *rancher.CanalNetworkProvider {
	obj := &rancher.CanalNetworkProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["iface"].(string); ok && len(v) > 0 {
		obj.FlannelNetworkProvider.Iface = v
	}

	return obj
}

func expandRKEClusterNetworkFlannel(p []interface{}) *rancher.FlannelNetworkProvider {
	obj := &rancher.FlannelNetworkProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["iface"].(string); ok && len(v) > 0 {
		obj.Iface = v
	}

	return obj
}

func expandRKEClusterNetworkWeave(p []interface{}) *rancher.WeaveNetworkProvider {
	obj := &rancher.WeaveNetworkProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["password"].(string); ok && len(v) > 0 {
		obj.Password = v
	}

	return obj
}

func expandRKEClusterNetworkAci(p []interface{}) *rancher.AciNetworkProvider {
	obj := &rancher.AciNetworkProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["system_id"].(string); ok && len(v) > 0 {
		obj.SystemIdentifier = v
	}
	if v, ok := in["apic_hosts"].([]interface{}); ok && len(v) > 0 {
		obj.ApicHosts = toArrayString(v)
	}
	if v, ok := in["token"].(string); ok && len(v) > 0 {
		obj.Token = v
	}
	if v, ok := in["apic_user_name"].(string); ok && len(v) > 0 {
		obj.ApicUserName = v
	}
	if v, ok := in["apic_user_key"].(string); ok && len(v) > 0 {
		obj.ApicUserKey = v
	}
	if v, ok := in["apic_user_crt"].(string); ok && len(v) > 0 {
		obj.ApicUserCrt = v
	}
	if v, ok := in["encap_type"].(string); ok && len(v) > 0 {
		obj.EncapType = v
	}
	if v, ok := in["mcast_range_start"].(string); ok && len(v) > 0 {
		obj.McastRangeStart = v
	}
	if v, ok := in["mcast_range_end"].(string); ok && len(v) > 0 {
		obj.McastRangeEnd = v
	}
	if v, ok := in["aep"].(string); ok && len(v) > 0 {
		obj.AEP = v
	}
	if v, ok := in["vrf_name"].(string); ok && len(v) > 0 {
		obj.VRFName = v
	}
	if v, ok := in["vrf_tenant"].(string); ok && len(v) > 0 {
		obj.VRFTenant = v
	}
	if v, ok := in["l3out"].(string); ok && len(v) > 0 {
		obj.L3Out = v
	}
	if v, ok := in["node_subnet"].(string); ok && len(v) > 0 {
		obj.NodeSubnet = v
	}
	if v, ok := in["l3out_external_networks"].([]interface{}); ok && len(v) > 0 {
		obj.L3OutExternalNetworks = toArrayString(v)
	}
	if v, ok := in["extern_dynamic"].(string); ok && len(v) > 0 {
		obj.DynamicExternalSubnet = v
	}
	if v, ok := in["extern_static"].(string); ok && len(v) > 0 {
		obj.StaticExternalSubnet = v
	}
	if v, ok := in["node_svc_subnet"].(string); ok && len(v) > 0 {
		obj.ServiceGraphSubnet = v
	}
	if v, ok := in["kube_api_vlan"].(string); ok && len(v) > 0 {
		obj.KubeAPIVlan = v
	}
	if v, ok := in["service_vlan"].(string); ok && len(v) > 0 {
		obj.ServiceVlan = v
	}
	if v, ok := in["infra_vlan"].(string); ok && len(v) > 0 {
		obj.InfraVlan = v
	}
	if v, ok := in["snat_port_range_start"].(string); ok && len(v) > 0 {
		obj.SnatPortRangeStart = v
	}
	if v, ok := in["snat_port_range_end"].(string); ok && len(v) > 0 {
		obj.SnatPortRangeEnd = v
	}
	if v, ok := in["snat_ports_per_node"].(string); ok && len(v) > 0 {
		obj.SnatPortsPerNode = v
	}

	return obj
}

func expandRKEClusterNetwork(p []interface{}) rancher.NetworkConfig {
	obj := rancher.NetworkConfig{}
	if len(p) == 0 || p[0] == nil {
		obj.Plugin = rkeClusterNetworkPluginDefault
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["calico_network_provider"].([]interface{}); ok && len(v) > 0 {
		obj.CalicoNetworkProvider = expandRKEClusterNetworkCalico(v)
	}

	if v, ok := in["canal_network_provider"].([]interface{}); ok && len(v) > 0 {
		obj.CanalNetworkProvider = expandRKEClusterNetworkCanal(v)
	}

	if v, ok := in["flannel_network_provider"].([]interface{}); ok && len(v) > 0 {
		obj.FlannelNetworkProvider = expandRKEClusterNetworkFlannel(v)
	}

	if v, ok := in["weave_network_provider"].([]interface{}); ok && len(v) > 0 {
		obj.WeaveNetworkProvider = expandRKEClusterNetworkWeave(v)
	}

	if v, ok := in["aci_network_provider"].([]interface{}); ok && len(v) > 0 {
		obj.AciNetworkProvider = expandRKEClusterNetworkAci(v)
	}

	if v, ok := in["mtu"].(int); ok && v > 0 {
		obj.MTU = v
	}

	if v, ok := in["options"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Options = toMapString(v)
	}

	if v, ok := in["plugin"].(string); ok && len(v) > 0 {
		obj.Plugin = v
	}

	return obj
}
