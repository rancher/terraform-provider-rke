package rke

import (
	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

// Flatteners

func flattenRKEClusterIngress(in rancher.IngressConfig) []interface{} {
	obj := make(map[string]interface{})

	if len(in.DNSPolicy) > 0 {
		obj["dns_policy"] = in.DNSPolicy
	}

	if len(in.ExtraArgs) > 0 {
		obj["extra_args"] = toMapInterface(in.ExtraArgs)
	}

	if len(in.NodeSelector) > 0 {
		obj["node_selector"] = toMapInterface(in.NodeSelector)
	}

	if len(in.Options) > 0 {
		obj["options"] = toMapInterface(in.Options)
	}

	if len(in.Provider) > 0 {
		obj["provider"] = in.Provider
	}

	return []interface{}{obj}
}

// Expanders

func expandRKEClusterIngress(p []interface{}) rancher.IngressConfig {
	obj := rancher.IngressConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["dns_policy"].(string); ok && len(v) > 0 {
		obj.DNSPolicy = v
	}

	if v, ok := in["extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ExtraArgs = toMapString(v)
	}

	if v, ok := in["node_selector"].(map[string]interface{}); ok && len(v) > 0 {
		obj.NodeSelector = toMapString(v)
	}

	if v, ok := in["options"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Options = toMapString(v)
	}

	if v, ok := in["provider"].(string); ok && len(v) > 0 {
		obj.Provider = v
	}

	return obj
}
