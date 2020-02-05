package rke

import (
	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
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
