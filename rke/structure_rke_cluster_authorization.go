package rke

import (
	rancher "github.com/rancher/rke/types"
)

// Flatteners

func flattenRKEClusterAuthorization(in rancher.AuthzConfig) []interface{} {
	obj := make(map[string]interface{})

	if len(in.Mode) > 0 {
		obj["mode"] = in.Mode
	}

	if len(in.Options) > 0 {
		obj["options"] = toMapInterface(in.Options)
	}

	return []interface{}{obj}
}

// Expanders

func expandRKEClusterAuthorization(p []interface{}) rancher.AuthzConfig {
	obj := rancher.AuthzConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["mode"].(string); ok && len(v) > 0 {
		obj.Mode = v
	}

	if v, ok := in["options"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Options = toMapString(v)
	}

	return obj
}
