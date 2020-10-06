package rke

import (
	rancher "github.com/rancher/rke/types"
)

// Flatteners

func flattenRKEClusterCloudProvider(in rancher.CloudProvider, p []interface{}) []interface{} {
	if len(in.Name) == 0 {
		return nil
	}

	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	obj["name"] = in.Name

	if in.AWSCloudProvider != nil {
		obj["aws_cloud_provider"] = flattenRKEClusterCloudProviderAws(in.AWSCloudProvider)
	}

	if in.AzureCloudProvider != nil {
		v, ok := obj["azure_cloud_provider"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		obj["azure_cloud_provider"] = flattenRKEClusterCloudProviderAzure(in.AzureCloudProvider, v)
	}

	if len(in.CustomCloudProvider) > 0 {
		obj["custom_cloud_provider"] = in.CustomCloudProvider
	}

	if in.OpenstackCloudProvider != nil {
		v, ok := obj["openstack_cloud_provider"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		obj["openstack_cloud_provider"] = flattenRKEClusterCloudProviderOpenstack(in.OpenstackCloudProvider, v)
	}

	if in.VsphereCloudProvider != nil {
		v, ok := obj["vsphere_cloud_provider"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		obj["vsphere_cloud_provider"] = flattenRKEClusterCloudProviderVsphere(in.VsphereCloudProvider, v)
	}

	return []interface{}{obj}
}

// Expanders

func expandRKEClusterCloudProvider(p []interface{}) rancher.CloudProvider {
	obj := rancher.CloudProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["aws_cloud_provider"].([]interface{}); ok && len(v) > 0 {
		obj.AWSCloudProvider = expandRKEClusterCloudProviderAws(v)
	}

	if v, ok := in["azure_cloud_provider"].([]interface{}); ok && len(v) > 0 {
		obj.AzureCloudProvider = expandRKEClusterCloudProviderAzure(v)
	}

	if v, ok := in["custom_cloud_provider"].(string); ok && len(v) > 0 {
		obj.CustomCloudProvider = v
	}

	if v, ok := in["name"].(string); ok && len(v) > 0 {
		obj.Name = v
	}

	if v, ok := in["openstack_cloud_provider"].([]interface{}); ok && len(v) > 0 {
		obj.OpenstackCloudProvider = expandRKEClusterCloudProviderOpenstack(v)
	}

	if v, ok := in["vsphere_cloud_provider"].([]interface{}); ok && len(v) > 0 {
		obj.VsphereCloudProvider = expandRKEClusterCloudProviderVsphere(v)
	}

	return obj
}
