package rke

import (
	rancher "github.com/rancher/rke/types"
)

// Flatteners

func flattenRKEClusterCloudProviderAzure(in *rancher.AzureCloudProvider, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}
	}

	if len(in.AADClientID) > 0 {
		obj["aad_client_id"] = in.AADClientID
	}

	if len(in.AADClientSecret) > 0 {
		obj["aad_client_secret"] = in.AADClientSecret
	}

	if len(in.SubscriptionID) > 0 {
		obj["subscription_id"] = in.SubscriptionID
	}

	if len(in.TenantID) > 0 {
		obj["tenant_id"] = in.TenantID
	}

	if len(in.AADClientCertPassword) > 0 {
		obj["aad_client_cert_password"] = in.AADClientCertPassword
	}

	if len(in.AADClientCertPath) > 0 {
		obj["aad_client_cert_path"] = in.AADClientCertPath
	}

	if len(in.Cloud) > 0 {
		obj["cloud"] = in.Cloud
	}

	obj["cloud_provider_backoff"] = in.CloudProviderBackoff

	if in.CloudProviderBackoffDuration > 0 {
		obj["cloud_provider_backoff_duration"] = in.CloudProviderBackoffDuration
	}

	if in.CloudProviderBackoffExponent > 0 {
		obj["cloud_provider_backoff_exponent"] = in.CloudProviderBackoffExponent
	}

	if in.CloudProviderBackoffJitter > 0 {
		obj["cloud_provider_backoff_jitter"] = in.CloudProviderBackoffJitter
	}

	if in.CloudProviderBackoffRetries > 0 {
		obj["cloud_provider_backoff_retries"] = in.CloudProviderBackoffRetries
	}

	obj["cloud_provider_rate_limit"] = in.CloudProviderRateLimit

	if in.CloudProviderRateLimitBucket > 0 {
		obj["cloud_provider_rate_limit_bucket"] = in.CloudProviderRateLimitBucket
	}

	if in.CloudProviderRateLimitQPS > 0 {
		obj["cloud_provider_rate_limit_qps"] = in.CloudProviderRateLimitQPS
	}

	if len(in.Location) > 0 {
		obj["location"] = in.Location
	}

	if in.MaximumLoadBalancerRuleCount > 0 {
		obj["maximum_load_balancer_rule_count"] = in.MaximumLoadBalancerRuleCount
	}

	if len(in.PrimaryAvailabilitySetName) > 0 {
		obj["primary_availability_set_name"] = in.PrimaryAvailabilitySetName
	}

	if len(in.PrimaryScaleSetName) > 0 {
		obj["primary_scale_set_name"] = in.PrimaryScaleSetName
	}

	if len(in.ResourceGroup) > 0 {
		obj["resource_group"] = in.ResourceGroup
	}

	if len(in.RouteTableName) > 0 {
		obj["route_table_name"] = in.RouteTableName
	}

	if len(in.SecurityGroupName) > 0 {
		obj["security_group_name"] = in.SecurityGroupName
	}

	if len(in.SubnetName) > 0 {
		obj["subnet_name"] = in.SubnetName
	}

	obj["use_instance_metadata"] = in.UseInstanceMetadata
	obj["use_managed_identity_extension"] = in.UseManagedIdentityExtension

	if len(in.VMType) > 0 {
		obj["vm_type"] = in.VMType
	}

	if len(in.VnetName) > 0 {
		obj["vnet_name"] = in.VnetName
	}

	if len(in.VnetResourceGroup) > 0 {
		obj["vnet_resource_group"] = in.VnetResourceGroup
	}

	return []interface{}{obj}
}

// Expanders

func expandRKEClusterCloudProviderAzure(p []interface{}) *rancher.AzureCloudProvider {
	obj := &rancher.AzureCloudProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["aad_client_id"].(string); ok && len(v) > 0 {
		obj.AADClientID = v
	}

	if v, ok := in["aad_client_secret"].(string); ok && len(v) > 0 {
		obj.AADClientSecret = v
	}

	if v, ok := in["subscription_id"].(string); ok && len(v) > 0 {
		obj.SubscriptionID = v
	}

	if v, ok := in["tenant_id"].(string); ok && len(v) > 0 {
		obj.TenantID = v
	}

	if v, ok := in["aad_client_cert_password"].(string); ok && len(v) > 0 {
		obj.AADClientCertPassword = v
	}

	if v, ok := in["aad_client_cert_path"].(string); ok && len(v) > 0 {
		obj.AADClientCertPath = v
	}

	if v, ok := in["cloud"].(string); ok && len(v) > 0 {
		obj.Cloud = v
	}

	if v, ok := in["cloud_provider_backoff"].(bool); ok {
		obj.CloudProviderBackoff = v
	}

	if v, ok := in["cloud_provider_backoff_duration"].(int); ok && v > 0 {
		obj.CloudProviderBackoffDuration = v
	}

	if v, ok := in["cloud_provider_backoff_exponent"].(int); ok && v > 0 {
		obj.CloudProviderBackoffExponent = v
	}

	if v, ok := in["cloud_provider_backoff_jitter"].(int); ok && v > 0 {
		obj.CloudProviderBackoffJitter = v
	}

	if v, ok := in["cloud_provider_backoff_retries"].(int); ok && v > 0 {
		obj.CloudProviderBackoffRetries = v
	}

	if v, ok := in["cloud_provider_rate_limit"].(bool); ok {
		obj.CloudProviderRateLimit = v
	}

	if v, ok := in["cloud_provider_rate_limit_bucket"].(int); ok && v > 0 {
		obj.CloudProviderRateLimitBucket = v
	}

	if v, ok := in["cloud_provider_rate_limit_qps"].(int); ok && v > 0 {
		obj.CloudProviderRateLimitQPS = v
	}

	if v, ok := in["location"].(string); ok && len(v) > 0 {
		obj.Location = v
	}

	if v, ok := in["maximum_load_balancer_rule_count"].(int); ok && v > 0 {
		obj.MaximumLoadBalancerRuleCount = v
	}

	if v, ok := in["primary_availability_set_name"].(string); ok && len(v) > 0 {
		obj.PrimaryAvailabilitySetName = v
	}

	if v, ok := in["primary_scale_set_name"].(string); ok && len(v) > 0 {
		obj.PrimaryScaleSetName = v
	}

	if v, ok := in["resource_group"].(string); ok && len(v) > 0 {
		obj.ResourceGroup = v
	}

	if v, ok := in["route_table_name"].(string); ok && len(v) > 0 {
		obj.RouteTableName = v
	}

	if v, ok := in["security_group_name"].(string); ok && len(v) > 0 {
		obj.SecurityGroupName = v
	}

	if v, ok := in["subnet_name"].(string); ok && len(v) > 0 {
		obj.SubnetName = v
	}

	if v, ok := in["use_instance_metadata"].(bool); ok {
		obj.UseInstanceMetadata = v
	}

	if v, ok := in["use_managed_identity_extension"].(bool); ok {
		obj.UseManagedIdentityExtension = v
	}

	if v, ok := in["vm_type"].(string); ok && len(v) > 0 {
		obj.VMType = v
	}

	if v, ok := in["vnet_name"].(string); ok && len(v) > 0 {
		obj.VnetName = v
	}

	if v, ok := in["vnet_resource_group"].(string); ok && len(v) > 0 {
		obj.VnetResourceGroup = v
	}

	return obj
}
