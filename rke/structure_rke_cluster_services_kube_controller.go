package rke

import (
	"encoding/json"
	rancher "github.com/rancher/rke/types"
)

// Flatteners

func flattenRKEClusterServicesKubeController(in rancher.KubeControllerService) ([]interface{}, error) {
	obj := make(map[string]interface{})

	if len(in.ClusterCIDR) > 0 {
		obj["cluster_cidr"] = in.ClusterCIDR
	}

	if len(in.ExtraArgs) > 0 {
		obj["extra_args"] = toMapInterface(in.ExtraArgs)
	}

	if len(in.WindowsExtraArgs) > 0 {
		obj["windows_extra_args"] = toMapInterface(in.WindowsExtraArgs)
	}

	if len(in.ExtraArgsArray) > 0 {
		j, err := json.Marshal(in.ExtraArgsArray)
		if err != nil {
			return nil, err
		}
		obj["extra_args_array"] = string(j)
	}

	if len(in.WindowsExtraArgsArray) > 0 {
		j, err := json.Marshal(in.WindowsExtraArgsArray)
		if err != nil {
			return nil, err
		}
		obj["windows_extra_args_array"] = string(j)
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

	if len(in.ServiceClusterIPRange) > 0 {
		obj["service_cluster_ip_range"] = in.ServiceClusterIPRange
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandRKEClusterServicesKubeController(p []interface{}) (rancher.KubeControllerService, error) {
	obj := rancher.KubeControllerService{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["cluster_cidr"].(string); ok && len(v) > 0 {
		obj.ClusterCIDR = v
	}

	if v, ok := in["extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ExtraArgs = toMapString(v)
	}

	if v, ok := in["windows_extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.WindowsExtraArgs = toMapString(v)
	}

	if v, ok := in["extra_args_array"].(string); ok && len(v) != 0 {
		array, err := jsonToMapStringSlice(v)
		if err != nil {
			return rancher.KubeControllerService{}, err
		}
		obj.ExtraArgsArray = array
	}

	if v, ok := in["windows_extra_args_array"].(string); ok && len(v) != 0 {
		array, err := jsonToMapStringSlice(v)
		if err != nil {
			return rancher.KubeControllerService{}, err
		}
		obj.WindowsExtraArgsArray = array
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

	if v, ok := in["service_cluster_ip_range"].(string); ok && len(v) > 0 {
		obj.ServiceClusterIPRange = v
	}

	return obj, nil
}
