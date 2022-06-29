package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	rkeClusterCloudProviderAzureName                    = "azure"
	rkeClusterCloudProviderAzureLoadBalancerSkuBasic    = "basic"
	rkeClusterCloudProviderAzureLoadBalancerSkuStandard = "standard"
)

var (
	rkeClusterCloudProviderAzureLoadBalancerSkuList = []string{
		rkeClusterCloudProviderAzureLoadBalancerSkuBasic,
		rkeClusterCloudProviderAzureLoadBalancerSkuStandard,
	}
)

//Schemas

func rkeClusterCloudProviderAzureFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"aad_client_id": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "The ClientID for an AAD application with RBAC access to talk to Azure RM APIs",
		},
		"aad_client_secret": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "The ClientSecret for an AAD application with RBAC access to talk to Azure RM APIs",
		},
		"subscription_id": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "The ID of the Azure Subscription that the cluster is deployed in",
		},
		"tenant_id": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "The AAD Tenant ID for the Subscription that the cluster is deployed in",
		},
		"aad_client_cert_password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "The password of the client certificate for an AAD application with RBAC access to talk to Azure RM APIs",
		},
		"aad_client_cert_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The path of a client certificate for an AAD application with RBAC access to talk to Azure RM APIs",
		},
		"cloud": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The cloud environment identifier. Takes values from https://github.com/Azure/go-autorest/blob/ec5f4903f77ed9927ac95b19ab8e44ada64c1356/autorest/azure/environments.go#L13",
		},
		"cloud_provider_backoff": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable exponential backoff to manage resource request retries",
		},
		"cloud_provider_backoff_duration": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Backoff duration",
		},
		"cloud_provider_backoff_exponent": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Backoff exponent",
		},
		"cloud_provider_backoff_jitter": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Backoff jitter",
		},
		"cloud_provider_backoff_retries": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Backoff retry limit",
		},
		"cloud_provider_rate_limit": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable rate limiting",
		},
		"cloud_provider_rate_limit_bucket": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"cloud_provider_rate_limit_qps": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Rate limit QPS",
		},
		"load_balancer_sku": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      rkeClusterCloudProviderAzureLoadBalancerSkuBasic,
			Description:  "Load balancer type (basic | standard). Must be standard for auto-scaling",
			ValidateFunc: validation.StringInSlice(rkeClusterCloudProviderAzureLoadBalancerSkuList, true),
		},
		"location": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The location of the resource group that the cluster is deployed in",
		},
		"maximum_load_balancer_rule_count": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      0,
			ValidateFunc: validation.IntBetween(0, 148),
			Description:  "Maximum allowed LoadBalancer Rule Count is the limit enforced by Azure Load balancer",
		},
		"primary_availability_set_name": {
			Type:     schema.TypeString,
			Optional: true,
			Description: "The name of the availability set that should be used as the load balancer backend" +
				"If this is set, the Azure cloudprovider will only add nodes from that availability set to the load" +
				"balancer backend pool. If this is not set, and multiple agent pools (availability sets) are used, then" +
				"the cloudprovider will try to add all nodes to a single backend pool which is forbidden." +
				"In other words, if you use multiple agent pools (availability sets), you MUST set this field.",
		},
		"primary_scale_set_name": {
			Type:     schema.TypeString,
			Optional: true,
			Description: "The name of the scale set that should be used as the load balancer backend." +
				"If this is set, the Azure cloudprovider will only add nodes from that scale set to the load" +
				"balancer backend pool. If this is not set, and multiple agent pools (scale sets) are used, then" +
				"the cloudprovider will try to add all nodes to a single backend pool which is forbidden." +
				"In other words, if you use multiple agent pools (scale sets), you MUST set this field.",
		},
		"resource_group": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the resource group that the cluster is deployed in",
		},
		"route_table_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "(Optional in 1.6) The name of the route table attached to the subnet that the cluster is deployed in",
		},
		"security_group_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the security group attached to the cluster's subnet",
		},
		"subnet_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the Subnet that the cluster is deployed in",
		},
		"use_instance_metadata": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Use instance metadata service where possible",
		},
		"use_managed_identity_extension": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Use managed service identity for the virtual machine to access Azure ARM APIs",
		},
		"vm_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The type of azure nodes. If not set, it will be default to standard.",
		},
		"vnet_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the VNet that the cluster is deployed in",
		},
		"vnet_resource_group": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The name of the resource group that the Vnet is deployed in",
		},
	}
	return s
}
