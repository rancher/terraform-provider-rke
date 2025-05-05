package rke

import (
	rancher "github.com/rancher/rke/types"
)

// Flatteners

func flattenRKEClusterServices(in rancher.RKEConfigServices, p []interface{}) ([]interface{}, error) {
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
	obj["etcd"], _ = flattenRKEClusterServicesEtcd(in.Etcd, v)

	kubeAPI, err := flattenRKEClusterServicesKubeAPI(in.KubeAPI)
	if err != nil {
		return []interface{}{obj}, err
	}
	obj["kube_api"] = kubeAPI

	kubeController, err := flattenRKEClusterServicesKubeController(in.KubeController)
	if err != nil {
		return []interface{}{obj}, err
	}
	obj["kube_controller"] = kubeController

	kubelet, err := flattenRKEClusterServicesKubelet(in.Kubelet)
	if err != nil {
		return []interface{}{obj}, err
	}
	obj["kubelet"] = kubelet

	kubeproxy, err := flattenRKEClusterServicesKubeproxy(in.Kubeproxy)
	if err != nil {
		return []interface{}{obj}, err
	}
	obj["kubeproxy"] = kubeproxy

	scheduler, err := flattenRKEClusterServicesScheduler(in.Scheduler)
	if err != nil {
		return []interface{}{obj}, err
	}
	obj["scheduler"] = scheduler

	return []interface{}{obj}, nil
}

// Expanders

func expandRKEClusterServices(p []interface{}) (rancher.RKEConfigServices, error) {
	obj := rancher.RKEConfigServices{}
	if p == nil || len(p) == 0 || p[0] == nil {
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
		kubeAPI, err := expandRKEClusterServicesKubeAPI(v)
		if err != nil {
			return obj, err
		}
		obj.KubeAPI = kubeAPI
	}

	if v, ok := in["kube_controller"].([]interface{}); ok && len(v) > 0 {
		kubeController, err := expandRKEClusterServicesKubeController(v)
		if err != nil {
			return obj, err
		}
		obj.KubeController = kubeController
	}

	if v, ok := in["kubelet"].([]interface{}); ok && len(v) > 0 {
		kubelet, err := expandRKEClusterServicesKubelet(v)
		if err != nil {
			return obj, err
		}
		obj.Kubelet = kubelet
	}

	if v, ok := in["kubeproxy"].([]interface{}); ok && len(v) > 0 {
		kubeproxy, err := expandRKEClusterServicesKubeproxy(v)
		if err != nil {
			return obj, err
		}
		obj.Kubeproxy = kubeproxy
	}

	if v, ok := in["scheduler"].([]interface{}); ok && len(v) > 0 {
		scheduler, err := expandRKEClusterServicesScheduler(v)
		if err != nil {
			return obj, err
		}
		obj.Scheduler = scheduler
	}

	return obj, nil
}
