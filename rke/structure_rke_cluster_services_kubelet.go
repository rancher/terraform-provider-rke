package rke

import (
	rancher "github.com/rancher/rke/types"
)

// Flatteners

func flattenRKEClusterServicesKubelet(in rancher.KubeletService) []interface{} {
	obj := make(map[string]interface{})

	if len(in.ClusterDNSServer) > 0 {
		obj["cluster_dns_server"] = in.ClusterDNSServer
	}

	if len(in.ClusterDomain) > 0 {
		obj["cluster_domain"] = in.ClusterDomain
	}

	if len(in.ExtraArgs) > 0 {
		obj["extra_args"] = toMapInterface(in.ExtraArgs)
	}

	if len(in.ExtraBinds) > 0 {
		obj["extra_binds"] = toArrayInterface(in.ExtraBinds)
	}

	if len(in.ExtraEnv) > 0 {
		obj["extra_env"] = toArrayInterface(in.ExtraEnv)
	}

	obj["fail_swap_on"] = in.FailSwapOn
	obj["generate_serving_certificate"] = in.GenerateServingCertificate

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	if len(in.InfraContainerImage) > 0 {
		obj["infra_container_image"] = in.InfraContainerImage
	}

	return []interface{}{obj}
}

// Expanders

func expandRKEClusterServicesKubelet(p []interface{}) rancher.KubeletService {
	obj := rancher.KubeletService{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["cluster_dns_server"].(string); ok && len(v) > 0 {
		obj.ClusterDNSServer = v
	}

	if v, ok := in["cluster_domain"].(string); ok && len(v) > 0 {
		obj.ClusterDomain = v
	}

	if v, ok := in["extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ExtraArgs = toMapString(v)
	}

	if v, ok := in["extra_binds"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraBinds = toArrayString(v)
	}

	if v, ok := in["extra_env"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraEnv = toArrayString(v)
	}

	if v, ok := in["fail_swap_on"].(bool); ok {
		obj.FailSwapOn = v
	}

	if v, ok := in["generate_serving_certificate"].(bool); ok {
		obj.GenerateServingCertificate = v
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["infra_container_image"].(string); ok && len(v) > 0 {
		obj.InfraContainerImage = v
	}

	return obj
}
