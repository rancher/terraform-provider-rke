package rke

import (
	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
	v1 "k8s.io/api/core/v1"
)

// Flatteners

func flattenRKEClusterTaints(p []rancher.RKETaint) []interface{} {
	if len(p) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(p))
	for i, in := range p {
		obj := make(map[string]interface{})

		if len(in.Key) > 0 {
			obj["key"] = in.Key
		}

		if len(in.Value) > 0 {
			obj["value"] = in.Value
		}

		if len(in.Effect) > 0 {
			obj["effect"] = string(in.Effect)
		}

		out[i] = obj
	}

	return out
}

// Expanders

func expandRKEClusterTaints(p []interface{}) []rancher.RKETaint {
	if len(p) == 0 || p[0] == nil {
		return []rancher.RKETaint{}
	}

	obj := make([]rancher.RKETaint, len(p))

	for i := range p {
		in := p[i].(map[string]interface{})

		if v, ok := in["key"].(string); ok && len(v) > 0 {
			obj[i].Key = v
		}

		if v, ok := in["value"].(string); ok && len(v) > 0 {
			obj[i].Value = v
		}

		if v, ok := in["effect"].(string); ok && len(v) > 0 {
			obj[i].Effect = v1.TaintEffect(v)
		}
	}

	return obj
}
