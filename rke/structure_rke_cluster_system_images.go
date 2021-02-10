package rke

import (
	rancher "github.com/rancher/rke/types"
)

// Flatteners

func flattenRKEClusterSystemImages(in rancher.RKESystemImages) []interface{} {
	obj := make(map[string]interface{})

	obj["etcd"] = in.Etcd
	obj["alpine"] = in.Alpine
	obj["nginx_proxy"] = in.NginxProxy
	obj["cert_downloader"] = in.CertDownloader
	obj["kubernetes_services_sidecar"] = in.KubernetesServicesSidecar
	obj["kube_dns"] = in.KubeDNS
	obj["dnsmasq"] = in.DNSmasq
	obj["kube_dns_sidecar"] = in.KubeDNSSidecar
	obj["kube_dns_autoscaler"] = in.KubeDNSAutoscaler
	obj["coredns"] = in.CoreDNS
	obj["coredns_autoscaler"] = in.CoreDNSAutoscaler
	obj["kubernetes"] = in.Kubernetes
	obj["flannel"] = in.Flannel
	obj["flannel_cni"] = in.FlannelCNI
	obj["calico_node"] = in.CalicoNode
	obj["calico_cni"] = in.CalicoCNI
	obj["calico_controllers"] = in.CalicoControllers
	obj["calico_ctl"] = in.CalicoCtl
	obj["calico_flex_vol"] = in.CalicoFlexVol
	obj["canal_node"] = in.CanalNode
	obj["canal_cni"] = in.CanalCNI
	obj["canal_flannel"] = in.CanalFlannel
	obj["canal_flex_vol"] = in.CanalFlexVol
	obj["weave_node"] = in.WeaveNode
	obj["weave_cni"] = in.WeaveCNI
	obj["pod_infra_container"] = in.PodInfraContainer
	obj["ingress"] = in.Ingress
	obj["ingress_backend"] = in.IngressBackend
	obj["metrics_server"] = in.MetricsServer
	obj["windows_pod_infra_container"] = in.WindowsPodInfraContainer
	obj["nodelocal"] = in.Nodelocal
	obj["aci_cni_deploy_container"] = in.AciCniDeployContainer
	obj["aci_host_container"] = in.AciHostContainer
	obj["aci_opflex_container"] = in.AciOpflexContainer
	obj["aci_mcast_container"] = in.AciMcastContainer
	obj["aci_ovs_container"] = in.AciOpenvSwitchContainer
	obj["aci_controller_container"] = in.AciControllerContainer

	return []interface{}{obj}
}

// Expanders

func expandRKEClusterSystemImages(p []interface{}) rancher.RKESystemImages {
	obj := rancher.RKESystemImages{}
	if p == nil || len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["etcd"].(string); ok && len(v) > 0 {
		obj.Etcd = v
	}

	if v, ok := in["alpine"].(string); ok && len(v) > 0 {
		obj.Alpine = v
	}

	if v, ok := in["nginx_proxy"].(string); ok && len(v) > 0 {
		obj.NginxProxy = v
	}

	if v, ok := in["cert_downloader"].(string); ok && len(v) > 0 {
		obj.CertDownloader = v
	}

	if v, ok := in["kubernetes_services_sidecar"].(string); ok && len(v) > 0 {
		obj.KubernetesServicesSidecar = v
	}

	if v, ok := in["kube_dns"].(string); ok && len(v) > 0 {
		obj.KubeDNS = v
	}

	if v, ok := in["dnsmasq"].(string); ok && len(v) > 0 {
		obj.DNSmasq = v
	}

	if v, ok := in["kube_dns_sidecar"].(string); ok && len(v) > 0 {
		obj.KubeDNSSidecar = v
	}

	if v, ok := in["kube_dns_autoscaler"].(string); ok && len(v) > 0 {
		obj.KubeDNSAutoscaler = v
	}

	if v, ok := in["coredns"].(string); ok && len(v) > 0 {
		obj.CoreDNS = v
	}

	if v, ok := in["coredns_autoscaler"].(string); ok && len(v) > 0 {
		obj.CoreDNSAutoscaler = v
	}

	if v, ok := in["kubernetes"].(string); ok && len(v) > 0 {
		obj.Kubernetes = v
	}

	if v, ok := in["flannel"].(string); ok && len(v) > 0 {
		obj.Flannel = v
	}

	if v, ok := in["flannel_cni"].(string); ok && len(v) > 0 {
		obj.FlannelCNI = v
	}

	if v, ok := in["calico_node"].(string); ok && len(v) > 0 {
		obj.CalicoNode = v
	}

	if v, ok := in["calico_cni"].(string); ok && len(v) > 0 {
		obj.CalicoCNI = v
	}

	if v, ok := in["calico_controllers"].(string); ok && len(v) > 0 {
		obj.CalicoControllers = v
	}

	if v, ok := in["calico_ctl"].(string); ok && len(v) > 0 {
		obj.CalicoCtl = v
	}

	if v, ok := in["calico_flex_vol"].(string); ok && len(v) > 0 {
		obj.CalicoFlexVol = v
	}

	if v, ok := in["canal_node"].(string); ok && len(v) > 0 {
		obj.CanalNode = v
	}

	if v, ok := in["canal_cni"].(string); ok && len(v) > 0 {
		obj.CanalCNI = v
	}

	if v, ok := in["canal_flannel"].(string); ok && len(v) > 0 {
		obj.CanalFlannel = v
	}

	if v, ok := in["canal_flex_vol"].(string); ok && len(v) > 0 {
		obj.CanalFlexVol = v
	}

	if v, ok := in["weave_node"].(string); ok && len(v) > 0 {
		obj.WeaveNode = v
	}

	if v, ok := in["weave_cni"].(string); ok && len(v) > 0 {
		obj.WeaveCNI = v
	}

	if v, ok := in["pod_infra_container"].(string); ok && len(v) > 0 {
		obj.PodInfraContainer = v
	}

	if v, ok := in["ingress"].(string); ok && len(v) > 0 {
		obj.Ingress = v
	}

	if v, ok := in["ingress_backend"].(string); ok && len(v) > 0 {
		obj.IngressBackend = v
	}

	if v, ok := in["metrics_server"].(string); ok && len(v) > 0 {
		obj.MetricsServer = v
	}

	if v, ok := in["windows_pod_infra_container"].(string); ok && len(v) > 0 {
		obj.WindowsPodInfraContainer = v
	}

	if v, ok := in["nodelocal"].(string); ok && len(v) > 0 {
		obj.Nodelocal = v
	}

	if v, ok := in["aci_cni_deploy_container"].(string); ok && len(v) > 0 {
		obj.AciCniDeployContainer = v
	}

	if v, ok := in["aci_host_container"].(string); ok && len(v) > 0 {
		obj.AciHostContainer = v
	}

	if v, ok := in["aci_opflex_container"].(string); ok && len(v) > 0 {
		obj.AciOpflexContainer = v
	}

	if v, ok := in["aci_mcast_container"].(string); ok && len(v) > 0 {
		obj.AciMcastContainer = v
	}

	if v, ok := in["aci_ovs_container"].(string); ok && len(v) > 0 {
		obj.AciOpenvSwitchContainer = v
	}

	if v, ok := in["aci_controller_container"].(string); ok && len(v) > 0 {
		obj.AciControllerContainer = v
	}

	return obj
}
