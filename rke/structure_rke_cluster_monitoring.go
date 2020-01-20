package rke

import (
	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

// Flatteners

func flattenRKEClusterMonitoring(in rancher.MonitoringConfig) []interface{} {
	obj := make(map[string]interface{})

	if len(in.Options) > 0 {
		obj["options"] = toMapInterface(in.Options)
	}

	if len(in.Provider) > 0 {
		obj["provider"] = in.Provider
	}

	return []interface{}{obj}
}

// Expanders

func expandRKEClusterMonitoring(p []interface{}) rancher.MonitoringConfig {
	obj := rancher.MonitoringConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["options"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Options = toMapString(v)
	}

	if v, ok := in["provider"].(string); ok && len(v) > 0 {
		obj.Provider = v
	}

	return obj
}
