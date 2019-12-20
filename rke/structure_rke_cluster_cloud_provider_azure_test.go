package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

var (
	testRKEClusterCloudProviderAzureConf      *rancher.AzureCloudProvider
	testRKEClusterCloudProviderAzureInterface []interface{}
)

func init() {
	testRKEClusterCloudProviderAzureConf = &rancher.AzureCloudProvider{
		AADClientID:                  "XXXXXXXX",
		AADClientSecret:              "XXXXXXXXXXXX",
		SubscriptionID:               "YYYYYYYY",
		TenantID:                     "ZZZZZZZZ",
		AADClientCertPassword:        "password",
		AADClientCertPath:            "/home/user/.ssh",
		Cloud:                        "cloud",
		CloudProviderBackoff:         true,
		CloudProviderBackoffDuration: 30,
		CloudProviderBackoffExponent: 20,
		CloudProviderBackoffJitter:   10,
		CloudProviderBackoffRetries:  5,
		CloudProviderRateLimit:       true,
		CloudProviderRateLimitBucket: 15,
		CloudProviderRateLimitQPS:    100,
		Location:                     "location",
		MaximumLoadBalancerRuleCount: 150,
		PrimaryAvailabilitySetName:   "primary",
		PrimaryScaleSetName:          "primary_scale",
		ResourceGroup:                "resource_group",
		RouteTableName:               "route_table_name",
		SecurityGroupName:            "security_group_name",
		SubnetName:                   "subnet_name",
		UseInstanceMetadata:          true,
		UseManagedIdentityExtension:  true,
		VMType:                       "vm_type",
		VnetName:                     "vnet_name",
		VnetResourceGroup:            "vnet_resource_group",
	}
	testRKEClusterCloudProviderAzureInterface = []interface{}{
		map[string]interface{}{
			"aad_client_id":                    "XXXXXXXX",
			"aad_client_secret":                "XXXXXXXXXXXX",
			"subscription_id":                  "YYYYYYYY",
			"tenant_id":                        "ZZZZZZZZ",
			"aad_client_cert_password":         "password",
			"aad_client_cert_path":             "/home/user/.ssh",
			"cloud":                            "cloud",
			"cloud_provider_backoff":           true,
			"cloud_provider_backoff_duration":  30,
			"cloud_provider_backoff_exponent":  20,
			"cloud_provider_backoff_jitter":    10,
			"cloud_provider_backoff_retries":   5,
			"cloud_provider_rate_limit":        true,
			"cloud_provider_rate_limit_bucket": 15,
			"cloud_provider_rate_limit_qps":    100,
			"location":                         "location",
			"maximum_load_balancer_rule_count": 150,
			"primary_availability_set_name":    "primary",
			"primary_scale_set_name":           "primary_scale",
			"resource_group":                   "resource_group",
			"route_table_name":                 "route_table_name",
			"security_group_name":              "security_group_name",
			"subnet_name":                      "subnet_name",
			"use_instance_metadata":            true,
			"use_managed_identity_extension":   true,
			"vm_type":                          "vm_type",
			"vnet_name":                        "vnet_name",
			"vnet_resource_group":              "vnet_resource_group",
		},
	}
}

func TestFlattenRKEClusterCloudProviderAzure(t *testing.T) {

	cases := []struct {
		Input          *rancher.AzureCloudProvider
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterCloudProviderAzureConf,
			testRKEClusterCloudProviderAzureInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterCloudProviderAzure(tc.Input, testRKEClusterCloudProviderAzureInterface)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterCloudProviderAzure(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rancher.AzureCloudProvider
	}{
		{
			testRKEClusterCloudProviderAzureInterface,
			testRKEClusterCloudProviderAzureConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterCloudProviderAzure(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
