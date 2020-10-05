package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/rke/types"
)

var (
	testRKEClusterCloudProviderAwsGlobalConf               rancher.GlobalAwsOpts
	testRKEClusterCloudProviderAwsGlobalInterface          []interface{}
	testRKEClusterCloudProviderAwsServiceOverrideConf      map[string]rancher.ServiceOverride
	testRKEClusterCloudProviderAwsServiceOverrideInterface []interface{}
	testRKEClusterCloudProviderAwsConf                     *rancher.AWSCloudProvider
	testRKEClusterCloudProviderAwsInterface                []interface{}
)

func init() {
	testRKEClusterCloudProviderAwsGlobalConf = rancher.GlobalAwsOpts{
		DisableSecurityGroupIngress: true,
		DisableStrictZoneCheck:      true,
		ElbSecurityGroup:            "elb_group",
		KubernetesClusterID:         "k8s_id",
		KubernetesClusterTag:        "k8s_tag",
		RoleARN:                     "role_arn",
		RouteTableID:                "route_table_id",
		SubnetID:                    "subnet_id",
		VPC:                         "vpc",
		Zone:                        "zone",
	}
	testRKEClusterCloudProviderAwsGlobalInterface = []interface{}{
		map[string]interface{}{
			"disable_security_group_ingress": true,
			"disable_strict_zone_check":      true,
			"elb_security_group":             "elb_group",
			"kubernetes_cluster_id":          "k8s_id",
			"kubernetes_cluster_tag":         "k8s_tag",
			"role_arn":                       "role_arn",
			"route_table_id":                 "route_table_id",
			"subnet_id":                      "subnet_id",
			"vpc":                            "vpc",
			"zone":                           "zone",
		},
	}
	testRKEClusterCloudProviderAwsServiceOverrideConf = map[string]rancher.ServiceOverride{
		"service": {
			Region:        "region",
			Service:       "service",
			SigningMethod: "signing_method",
			SigningName:   "signing_name",
			SigningRegion: "signing_region",
			URL:           "url",
		},
	}
	testRKEClusterCloudProviderAwsServiceOverrideInterface = []interface{}{
		map[string]interface{}{
			"region":         "region",
			"service":        "service",
			"signing_method": "signing_method",
			"signing_name":   "signing_name",
			"signing_region": "signing_region",
			"url":            "url",
		},
	}
	testRKEClusterCloudProviderAwsConf = &rancher.AWSCloudProvider{
		Global:          testRKEClusterCloudProviderAwsGlobalConf,
		ServiceOverride: testRKEClusterCloudProviderAwsServiceOverrideConf,
	}
	testRKEClusterCloudProviderAwsInterface = []interface{}{
		map[string]interface{}{
			"global":           testRKEClusterCloudProviderAwsGlobalInterface,
			"service_override": testRKEClusterCloudProviderAwsServiceOverrideInterface,
		},
	}
}

func TestFlattenRKEClusterCloudProviderAwsGlobal(t *testing.T) {

	cases := []struct {
		Input          rancher.GlobalAwsOpts
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterCloudProviderAwsGlobalConf,
			testRKEClusterCloudProviderAwsGlobalInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterCloudProviderAwsGlobal(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterCloudProviderAwsServiceOverride(t *testing.T) {

	cases := []struct {
		Input          map[string]rancher.ServiceOverride
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterCloudProviderAwsServiceOverrideConf,
			testRKEClusterCloudProviderAwsServiceOverrideInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterCloudProviderAwsServiceOverride(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterCloudProviderAws(t *testing.T) {

	cases := []struct {
		Input          *rancher.AWSCloudProvider
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterCloudProviderAwsConf,
			testRKEClusterCloudProviderAwsInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterCloudProviderAws(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterCloudProviderAwsGlobal(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.GlobalAwsOpts
	}{
		{
			testRKEClusterCloudProviderAwsGlobalInterface,
			testRKEClusterCloudProviderAwsGlobalConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterCloudProviderAwsGlobal(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterCloudProviderAwsServiceOverride(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput map[string]rancher.ServiceOverride
	}{
		{
			testRKEClusterCloudProviderAwsServiceOverrideInterface,
			testRKEClusterCloudProviderAwsServiceOverrideConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterCloudProviderAwsServiceOverride(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterCloudProviderAws(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rancher.AWSCloudProvider
	}{
		{
			testRKEClusterCloudProviderAwsInterface,
			testRKEClusterCloudProviderAwsConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterCloudProviderAws(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
