package rke

import (
	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

// Flatteners

func flattenRKEClusterAuthentication(in rancher.AuthnConfig) []interface{} {
	obj := make(map[string]interface{})

	if len(in.SANs) > 0 {
		obj["sans"] = toArrayInterface(in.SANs)
	}

	if len(in.Strategy) > 0 {
		obj["strategy"] = in.Strategy
	}

	return []interface{}{obj}
}

// Expanders

func expandRKEClusterAuthentication(p []interface{}) rancher.AuthnConfig {
	obj := rancher.AuthnConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["sans"].([]interface{}); ok && len(v) > 0 {
		obj.SANs = toArrayString(v)
	}

	if v, ok := in["strategy"].(string); ok && len(v) > 0 {
		obj.Strategy = v
	}

	return obj
}
