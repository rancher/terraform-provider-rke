package rke

import (
	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

// Flatteners

func flattenRKEClusterServicesScheduler(in rancher.SchedulerService) []interface{} {
	obj := make(map[string]interface{})

	if len(in.ExtraArgs) > 0 {
		obj["extra_args"] = toMapInterface(in.ExtraArgs)
	}

	if len(in.ExtraBinds) > 0 {
		obj["extra_binds"] = toArrayInterface(in.ExtraBinds)
	}

	if len(in.ExtraEnv) > 0 {
		obj["extra_env"] = toArrayInterface(in.ExtraEnv)
	}

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	return []interface{}{obj}
}

// Expanders

func expandRKEClusterServicesScheduler(p []interface{}) rancher.SchedulerService {
	obj := rancher.SchedulerService{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ExtraArgs = toMapString(v)
	}

	if v, ok := in["extra_binds"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraBinds = toArrayString(v)
	}

	if v, ok := in["extra_env"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraEnv = toArrayString(v)
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	return obj
}
