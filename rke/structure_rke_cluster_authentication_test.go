package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

var (
	testRKEClusterAuthenticationConf      rancher.AuthnConfig
	testRKEClusterAuthenticationInterface []interface{}
)

func init() {
	testRKEClusterAuthenticationConf = rancher.AuthnConfig{
		SANs:     []string{"sans1", "sans2"},
		Strategy: "strategy",
	}
	testRKEClusterAuthenticationInterface = []interface{}{
		map[string]interface{}{
			"sans":     []interface{}{"sans1", "sans2"},
			"strategy": "strategy",
		},
	}
}

func TestFlattenRKEClusterAuthentication(t *testing.T) {

	cases := []struct {
		Input          rancher.AuthnConfig
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterAuthenticationConf,
			testRKEClusterAuthenticationInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterAuthentication(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterAuthentication(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.AuthnConfig
	}{
		{
			testRKEClusterAuthenticationInterface,
			testRKEClusterAuthenticationConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterAuthentication(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
