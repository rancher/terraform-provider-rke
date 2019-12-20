package rke

import (
	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

// Flatteners

func flattenRKEClusterCloudProviderVsphereDisk(in rancher.DiskVsphereOpts) []interface{} {
	obj := make(map[string]interface{})

	if len(in.SCSIControllerType) > 0 {
		obj["scsi_controller_type"] = in.SCSIControllerType
	}

	return []interface{}{obj}
}

func flattenRKEClusterCloudProviderVsphereGlobal(in rancher.GlobalVsphereOpts, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if len(in.Datacenters) > 0 {
		obj["datacenters"] = in.Datacenters
	}

	if len(in.DefaultDatastore) > 0 {
		obj["datastore"] = in.DefaultDatastore
	}

	obj["insecure_flag"] = in.InsecureFlag

	if len(in.Password) > 0 {
		obj["password"] = in.Password
	}

	if len(in.VCenterPort) > 0 {
		obj["port"] = in.VCenterPort
	}

	if len(in.User) > 0 {
		obj["user"] = in.User
	}

	if in.RoundTripperCount > 0 {
		obj["soap_roundtrip_count"] = in.RoundTripperCount
	}

	if len(in.VMName) > 0 {
		obj["vm_name"] = in.VMName
	}

	if len(in.VMUUID) > 0 {
		obj["vm_uuid"] = in.VMUUID
	}

	if len(in.WorkingDir) > 0 {
		obj["working_dir"] = in.WorkingDir
	}

	return []interface{}{obj}
}

func flattenRKEClusterCloudProviderVsphereNetwork(in rancher.NetworkVshpereOpts) []interface{} {
	obj := make(map[string]interface{})

	if len(in.PublicNetwork) > 0 {
		obj["public_network"] = in.PublicNetwork
	}

	return []interface{}{obj}
}

func flattenRKEClusterCloudProviderVsphereVirtualCenter(in map[string]rancher.VirtualCenterConfig, p []interface{}) []interface{} {
	if len(in) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(in))
	lenP := len(p)
	i := 0
	for key := range in {
		var obj map[string]interface{}
		if lenP <= i {
			obj = make(map[string]interface{})
		} else {
			obj = p[i].(map[string]interface{})
		}

		obj["name"] = key
		if len(in[key].Datacenters) > 0 {
			obj["datacenters"] = in[key].Datacenters
		}

		if len(in[key].Password) > 0 {
			obj["password"] = in[key].Password
		}

		if len(in[key].VCenterPort) > 0 {
			obj["port"] = in[key].VCenterPort
		}

		if len(in[key].User) > 0 {
			obj["user"] = in[key].User
		}

		if in[key].RoundTripperCount > 0 {
			obj["soap_roundtrip_count"] = in[key].RoundTripperCount
		}
		out[i] = obj
		i++
	}

	return out
}

func flattenRKEClusterCloudProviderVsphereWorkspace(in rancher.WorkspaceVsphereOpts) []interface{} {
	obj := make(map[string]interface{})

	if len(in.Datacenter) > 0 {
		obj["datacenter"] = in.Datacenter
	}

	if len(in.Folder) > 0 {
		obj["folder"] = in.Folder
	}

	if len(in.VCenterIP) > 0 {
		obj["server"] = in.VCenterIP
	}

	if len(in.DefaultDatastore) > 0 {
		obj["default_datastore"] = in.DefaultDatastore
	}

	if len(in.ResourcePoolPath) > 0 {
		obj["resourcepool_path"] = in.ResourcePoolPath
	}

	return []interface{}{obj}
}

func flattenRKEClusterCloudProviderVsphere(in *rancher.VsphereCloudProvider, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}
	}

	obj["disk"] = flattenRKEClusterCloudProviderVsphereDisk(in.Disk)

	v, ok := obj["global"].([]interface{})
	if !ok {
		v = []interface{}{}
	}
	obj["global"] = flattenRKEClusterCloudProviderVsphereGlobal(in.Global, v)

	obj["network"] = flattenRKEClusterCloudProviderVsphereNetwork(in.Network)

	v, ok = obj["virtual_center"].([]interface{})
	if !ok {
		v = []interface{}{}
	}
	obj["virtual_center"] = flattenRKEClusterCloudProviderVsphereVirtualCenter(in.VirtualCenter, v)

	obj["workspace"] = flattenRKEClusterCloudProviderVsphereWorkspace(in.Workspace)

	return []interface{}{obj}
}

// Expanders

func expandRKEClusterCloudProviderVsphereDisk(p []interface{}) rancher.DiskVsphereOpts {
	obj := rancher.DiskVsphereOpts{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["scsi_controller_type"].(string); ok && len(v) > 0 {
		obj.SCSIControllerType = v
	}

	return obj
}

func expandRKEClusterCloudProviderVsphereGlobal(p []interface{}) rancher.GlobalVsphereOpts {
	obj := rancher.GlobalVsphereOpts{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["datacenters"].(string); ok && len(v) > 0 {
		obj.Datacenters = v
	}

	if v, ok := in["datastore"].(string); ok && len(v) > 0 {
		obj.DefaultDatastore = v
	}

	if v, ok := in["insecure_flag"].(bool); ok {
		obj.InsecureFlag = v
	}

	if v, ok := in["password"].(string); ok && len(v) > 0 {
		obj.Password = v
	}

	if v, ok := in["port"].(string); ok && len(v) > 0 {
		obj.VCenterPort = v
	}

	if v, ok := in["user"].(string); ok && len(v) > 0 {
		obj.User = v
	}

	if v, ok := in["soap_roundtrip_count"].(int); ok && v > 0 {
		obj.RoundTripperCount = v
	}

	if v, ok := in["vm_name"].(string); ok && len(v) > 0 {
		obj.VMName = v
	}

	if v, ok := in["vm_uuid"].(string); ok && len(v) > 0 {
		obj.VMUUID = v
	}

	if v, ok := in["working_dir"].(string); ok && len(v) > 0 {
		obj.WorkingDir = v
	}

	return obj
}

func expandRKEClusterCloudProviderVsphereNetwork(p []interface{}) rancher.NetworkVshpereOpts {
	obj := rancher.NetworkVshpereOpts{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["public_network"].(string); ok && len(v) > 0 {
		obj.PublicNetwork = v
	}

	return obj
}

func expandRKEClusterCloudProviderVsphereVirtualCenter(p []interface{}) map[string]rancher.VirtualCenterConfig {
	if len(p) == 0 || p[0] == nil {
		return map[string]rancher.VirtualCenterConfig{}
	}

	obj := make(map[string]rancher.VirtualCenterConfig)

	for i := range p {
		in := p[i].(map[string]interface{})
		aux := rancher.VirtualCenterConfig{}
		key := in["name"].(string)

		if v, ok := in["datacenters"].(string); ok && len(v) > 0 {
			aux.Datacenters = v
		}

		if v, ok := in["password"].(string); ok && len(v) > 0 {
			aux.Password = v
		}

		if v, ok := in["port"].(string); ok && len(v) > 0 {
			aux.VCenterPort = v
		}

		if v, ok := in["user"].(string); ok && len(v) > 0 {
			aux.User = v
		}

		if v, ok := in["soap_roundtrip_count"].(int); ok && v > 0 {
			aux.RoundTripperCount = v
		}

		obj[key] = aux
	}

	return obj
}

func expandRKEClusterCloudProviderVsphereWorkspace(p []interface{}) rancher.WorkspaceVsphereOpts {
	obj := rancher.WorkspaceVsphereOpts{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["datacenter"].(string); ok && len(v) > 0 {
		obj.Datacenter = v
	}

	if v, ok := in["folder"].(string); ok && len(v) > 0 {
		obj.Folder = v
	}

	if v, ok := in["server"].(string); ok && len(v) > 0 {
		obj.VCenterIP = v
	}

	if v, ok := in["default_datastore"].(string); ok && len(v) > 0 {
		obj.DefaultDatastore = v
	}

	if v, ok := in["resourcepool_path"].(string); ok && len(v) > 0 {
		obj.ResourcePoolPath = v
	}

	return obj
}

func expandRKEClusterCloudProviderVsphere(p []interface{}) *rancher.VsphereCloudProvider {
	obj := &rancher.VsphereCloudProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["disk"].([]interface{}); ok && len(v) > 0 {
		obj.Disk = expandRKEClusterCloudProviderVsphereDisk(v)
	}

	if v, ok := in["global"].([]interface{}); ok && len(v) > 0 {
		obj.Global = expandRKEClusterCloudProviderVsphereGlobal(v)
	}

	if v, ok := in["network"].([]interface{}); ok && len(v) > 0 {
		obj.Network = expandRKEClusterCloudProviderVsphereNetwork(v)
	}

	if v, ok := in["virtual_center"].([]interface{}); ok && len(v) > 0 {
		obj.VirtualCenter = expandRKEClusterCloudProviderVsphereVirtualCenter(v)
	}

	if v, ok := in["workspace"].([]interface{}); ok && len(v) > 0 {
		obj.Workspace = expandRKEClusterCloudProviderVsphereWorkspace(v)
	}

	return obj
}
