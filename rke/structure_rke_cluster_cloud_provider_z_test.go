package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

var (
	testRKEClusterCloudProviderConfAzure          rancher.CloudProvider
	testRKEClusterCloudProviderInterfaceAzure     []interface{}
	testRKEClusterCloudProviderConfOpenstack      rancher.CloudProvider
	testRKEClusterCloudProviderInterfaceOpenstack []interface{}
	testRKEClusterCloudProviderConfVsphere        rancher.CloudProvider
	testRKEClusterCloudProviderInterfaceVsphere   []interface{}
	testRKEClusterCloudProviderConf               rancher.CloudProvider
	testRKEClusterCloudProviderInterface          []interface{}
)

func init() {
	testRKEClusterCloudProviderConfAzure = rancher.CloudProvider{
		AzureCloudProvider: testRKEClusterCloudProviderAzureConf,
		Name:               "azure-test",
	}
	testRKEClusterCloudProviderInterfaceAzure = []interface{}{
		map[string]interface{}{
			"azure_cloud_provider": testRKEClusterCloudProviderAzureInterface,
			"name":                 "azure-test",
		},
	}
	testRKEClusterCloudProviderConfOpenstack = rancher.CloudProvider{
		Name:                   "openstack-test",
		OpenstackCloudProvider: testRKEClusterCloudProviderOpenstackConf,
	}
	testRKEClusterCloudProviderInterfaceOpenstack = []interface{}{
		map[string]interface{}{
			"name":                     "openstack-test",
			"openstack_cloud_provider": testRKEClusterCloudProviderOpenstackInterface,
		},
	}
	testRKEClusterCloudProviderConfVsphere = rancher.CloudProvider{
		Name:                 "vsphere-test",
		VsphereCloudProvider: testRKEClusterCloudProviderVsphereConf,
	}
	testRKEClusterCloudProviderInterfaceVsphere = []interface{}{
		map[string]interface{}{
			"name":                   "vsphere-test",
			"vsphere_cloud_provider": testRKEClusterCloudProviderVsphereInterface,
		},
	}
	testRKEClusterCloudProviderConf = rancher.CloudProvider{
		CustomCloudProvider: "XXXXXXXXXXXX",
		Name:                "test",
	}
	testRKEClusterCloudProviderInterface = []interface{}{
		map[string]interface{}{
			"custom_cloud_provider": "XXXXXXXXXXXX",
			"name":                  "test",
		},
	}
}

func TestFlattenRKEClusterCloudProvider(t *testing.T) {

	cases := []struct {
		Input          rancher.CloudProvider
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterCloudProviderConfAzure,
			testRKEClusterCloudProviderInterfaceAzure,
		},
		{
			testRKEClusterCloudProviderConfOpenstack,
			testRKEClusterCloudProviderInterfaceOpenstack,
		},
		{
			testRKEClusterCloudProviderConfVsphere,
			testRKEClusterCloudProviderInterfaceVsphere,
		},
		{
			testRKEClusterCloudProviderConf,
			testRKEClusterCloudProviderInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterCloudProvider(tc.Input, tc.ExpectedOutput)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterCloudProvider(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.CloudProvider
	}{
		{
			testRKEClusterCloudProviderInterfaceAzure,
			testRKEClusterCloudProviderConfAzure,
		},
		{
			testRKEClusterCloudProviderInterfaceOpenstack,
			testRKEClusterCloudProviderConfOpenstack,
		},
		{
			testRKEClusterCloudProviderInterfaceVsphere,
			testRKEClusterCloudProviderConfVsphere,
		},
		{
			testRKEClusterCloudProviderInterface,
			testRKEClusterCloudProviderConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterCloudProvider(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
