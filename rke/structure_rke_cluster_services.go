package rke

import (
	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

// Flatteners

func flattenRKEClusterServices(in rancher.RKEConfigServices, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	v, ok := obj["etcd"].([]interface{})
	if !ok {
		v = []interface{}{}
	}
	obj["etcd"] = flattenRKEClusterServicesEtcd(in.Etcd, v)
	obj["kube_api"] = flattenRKEClusterServicesKubeAPI(in.KubeAPI)
	obj["kube_controller"] = flattenRKEClusterServicesKubeController(in.KubeController)
	obj["kubelet"] = flattenRKEClusterServicesKubelet(in.Kubelet)
	obj["kubeproxy"] = flattenRKEClusterServicesKubeproxy(in.Kubeproxy)
	obj["scheduler"] = flattenRKEClusterServicesScheduler(in.Scheduler)

	return []interface{}{obj}
}

// Expanders

func expandRKEClusterServices(p []interface{}) (rancher.RKEConfigServices, error) {
	obj := rancher.RKEConfigServices{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["etcd"].([]interface{}); ok && len(v) > 0 {
		etcd, err := expandRKEClusterServicesEtcd(v)
		if err != nil {
			return obj, err
		}
		obj.Etcd = etcd
	}

	if v, ok := in["kube_api"].([]interface{}); ok && len(v) > 0 {
		obj.KubeAPI = expandRKEClusterServicesKubeAPI(v)
	}

	if v, ok := in["kube_controller"].([]interface{}); ok && len(v) > 0 {
		obj.KubeController = expandRKEClusterServicesKubeController(v)
	}

	if v, ok := in["kubelet"].([]interface{}); ok && len(v) > 0 {
		obj.Kubelet = expandRKEClusterServicesKubelet(v)
	}

	if v, ok := in["kubeproxy"].([]interface{}); ok && len(v) > 0 {
		obj.Kubeproxy = expandRKEClusterServicesKubeproxy(v)
	}

	if v, ok := in["scheduler"].([]interface{}); ok && len(v) > 0 {
		obj.Scheduler = expandRKEClusterServicesScheduler(v)
	}

	return obj, nil
}
