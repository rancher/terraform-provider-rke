package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

var (
	testRKEClusterCloudProviderOpenstackBlockStorageConf      rancher.BlockStorageOpenstackOpts
	testRKEClusterCloudProviderOpenstackBlockStorageInterface []interface{}
	testRKEClusterCloudProviderOpenstackGlobalConf            rancher.GlobalOpenstackOpts
	testRKEClusterCloudProviderOpenstackGlobalInterface       []interface{}
	testRKEClusterCloudProviderOpenstackLoadBalancerConf      rancher.LoadBalancerOpenstackOpts
	testRKEClusterCloudProviderOpenstackLoadBalancerInterface []interface{}
	testRKEClusterCloudProviderOpenstackMetadataConf          rancher.MetadataOpenstackOpts
	testRKEClusterCloudProviderOpenstackMetadataInterface     []interface{}
	testRKEClusterCloudProviderOpenstackRouteConf             rancher.RouteOpenstackOpts
	testRKEClusterCloudProviderOpenstackRouteInterface        []interface{}
	testRKEClusterCloudProviderOpenstackConf                  *rancher.OpenstackCloudProvider
	testRKEClusterCloudProviderOpenstackInterface             []interface{}
)

func init() {
	testRKEClusterCloudProviderOpenstackBlockStorageConf = rancher.BlockStorageOpenstackOpts{
		BSVersion:       "test",
		IgnoreVolumeAZ:  true,
		TrustDevicePath: true,
	}
	testRKEClusterCloudProviderOpenstackBlockStorageInterface = []interface{}{
		map[string]interface{}{
			"bs_version":        "test",
			"ignore_volume_az":  true,
			"trust_device_path": true,
		},
	}
	testRKEClusterCloudProviderOpenstackGlobalConf = rancher.GlobalOpenstackOpts{
		AuthURL:    "auth.terraform.test",
		Password:   "XXXXXXXX",
		TenantID:   "YYYYYYYY",
		Username:   "user",
		CAFile:     "ca_file",
		DomainID:   "domain_id",
		DomainName: "domain_name",
		Region:     "region",
		TenantName: "tenant",
		TrustID:    "VVVVVVVV",
	}
	testRKEClusterCloudProviderOpenstackGlobalInterface = []interface{}{
		map[string]interface{}{
			"auth_url":    "auth.terraform.test",
			"password":    "XXXXXXXX",
			"tenant_id":   "YYYYYYYY",
			"username":    "user",
			"ca_file":     "ca_file",
			"domain_id":   "domain_id",
			"domain_name": "domain_name",
			"region":      "region",
			"tenant_name": "tenant",
			"trust_id":    "VVVVVVVV",
		},
	}
	testRKEClusterCloudProviderOpenstackLoadBalancerConf = rancher.LoadBalancerOpenstackOpts{
		CreateMonitor:        true,
		FloatingNetworkID:    "test",
		LBMethod:             "method",
		LBProvider:           "provider",
		LBVersion:            "version",
		ManageSecurityGroups: true,
		MonitorDelay:         "30s",
		MonitorMaxRetries:    5,
		MonitorTimeout:       "10s",
		SubnetID:             "subnet",
		UseOctavia:           true,
	}
	testRKEClusterCloudProviderOpenstackLoadBalancerInterface = []interface{}{
		map[string]interface{}{
			"create_monitor":         true,
			"floating_network_id":    "test",
			"lb_method":              "method",
			"lb_provider":            "provider",
			"lb_version":             "version",
			"manage_security_groups": true,
			"monitor_delay":          "30s",
			"monitor_max_retries":    5,
			"monitor_timeout":        "10s",
			"subnet_id":              "subnet",
			"use_octavia":            true,
		},
	}
	testRKEClusterCloudProviderOpenstackMetadataConf = rancher.MetadataOpenstackOpts{
		RequestTimeout: 30,
		SearchOrder:    "order",
	}
	testRKEClusterCloudProviderOpenstackMetadataInterface = []interface{}{
		map[string]interface{}{
			"request_timeout": 30,
			"search_order":    "order",
		},
	}
	testRKEClusterCloudProviderOpenstackRouteConf = rancher.RouteOpenstackOpts{
		RouterID: "test",
	}
	testRKEClusterCloudProviderOpenstackRouteInterface = []interface{}{
		map[string]interface{}{
			"router_id": "test",
		},
	}
	testRKEClusterCloudProviderOpenstackConf = &rancher.OpenstackCloudProvider{
		BlockStorage: testRKEClusterCloudProviderOpenstackBlockStorageConf,
		Global:       testRKEClusterCloudProviderOpenstackGlobalConf,
		LoadBalancer: testRKEClusterCloudProviderOpenstackLoadBalancerConf,
		Metadata:     testRKEClusterCloudProviderOpenstackMetadataConf,
		Route:        testRKEClusterCloudProviderOpenstackRouteConf,
	}
	testRKEClusterCloudProviderOpenstackInterface = []interface{}{
		map[string]interface{}{
			"block_storage": testRKEClusterCloudProviderOpenstackBlockStorageInterface,
			"global":        testRKEClusterCloudProviderOpenstackGlobalInterface,
			"load_balancer": testRKEClusterCloudProviderOpenstackLoadBalancerInterface,
			"metadata":      testRKEClusterCloudProviderOpenstackMetadataInterface,
			"route":         testRKEClusterCloudProviderOpenstackRouteInterface,
		},
	}
}

func TestFlattenRKEClusterCloudProviderOpenstackBlockStorage(t *testing.T) {

	cases := []struct {
		Input          rancher.BlockStorageOpenstackOpts
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterCloudProviderOpenstackBlockStorageConf,
			testRKEClusterCloudProviderOpenstackBlockStorageInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterCloudProviderOpenstackBlockStorage(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterCloudProviderOpenstackGlobal(t *testing.T) {

	cases := []struct {
		Input          rancher.GlobalOpenstackOpts
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterCloudProviderOpenstackGlobalConf,
			testRKEClusterCloudProviderOpenstackGlobalInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterCloudProviderOpenstackGlobal(tc.Input, testRKEClusterCloudProviderOpenstackGlobalInterface)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterCloudProviderOpenstackLoadBalancer(t *testing.T) {

	cases := []struct {
		Input          rancher.LoadBalancerOpenstackOpts
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterCloudProviderOpenstackLoadBalancerConf,
			testRKEClusterCloudProviderOpenstackLoadBalancerInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterCloudProviderOpenstackLoadBalancer(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterCloudProviderOpenstackMetadata(t *testing.T) {

	cases := []struct {
		Input          rancher.MetadataOpenstackOpts
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterCloudProviderOpenstackMetadataConf,
			testRKEClusterCloudProviderOpenstackMetadataInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterCloudProviderOpenstackMetadata(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterCloudProviderOpenstackRoute(t *testing.T) {

	cases := []struct {
		Input          rancher.RouteOpenstackOpts
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterCloudProviderOpenstackRouteConf,
			testRKEClusterCloudProviderOpenstackRouteInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterCloudProviderOpenstackRoute(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterCloudProviderOpenstack(t *testing.T) {

	cases := []struct {
		Input          *rancher.OpenstackCloudProvider
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterCloudProviderOpenstackConf,
			testRKEClusterCloudProviderOpenstackInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterCloudProviderOpenstack(tc.Input, testRKEClusterCloudProviderOpenstackInterface)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterCloudProviderOpenstackBlockStorage(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.BlockStorageOpenstackOpts
	}{
		{
			testRKEClusterCloudProviderOpenstackBlockStorageInterface,
			testRKEClusterCloudProviderOpenstackBlockStorageConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterCloudProviderOpenstackBlockStorage(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterCloudProviderOpenstackGlobal(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.GlobalOpenstackOpts
	}{
		{
			testRKEClusterCloudProviderOpenstackGlobalInterface,
			testRKEClusterCloudProviderOpenstackGlobalConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterCloudProviderOpenstackGlobal(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterCloudProviderOpenstackLoadBalancer(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.LoadBalancerOpenstackOpts
	}{
		{
			testRKEClusterCloudProviderOpenstackLoadBalancerInterface,
			testRKEClusterCloudProviderOpenstackLoadBalancerConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterCloudProviderOpenstackLoadBalancer(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterCloudProviderOpenstackMetadata(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.MetadataOpenstackOpts
	}{
		{
			testRKEClusterCloudProviderOpenstackMetadataInterface,
			testRKEClusterCloudProviderOpenstackMetadataConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterCloudProviderOpenstackMetadata(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterCloudProviderOpenstackRoute(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.RouteOpenstackOpts
	}{
		{
			testRKEClusterCloudProviderOpenstackRouteInterface,
			testRKEClusterCloudProviderOpenstackRouteConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterCloudProviderOpenstackRoute(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterCloudProviderOpenstack(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rancher.OpenstackCloudProvider
	}{
		{
			testRKEClusterCloudProviderOpenstackInterface,
			testRKEClusterCloudProviderOpenstackConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterCloudProviderOpenstack(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
