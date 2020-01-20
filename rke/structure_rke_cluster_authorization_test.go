package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

var (
	testRKEClusterAuthorizationConf      rancher.AuthzConfig
	testRKEClusterAuthorizationInterface []interface{}
)

func init() {
	testRKEClusterAuthorizationConf = rancher.AuthzConfig{
		Mode: "rbac",
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testRKEClusterAuthorizationInterface = []interface{}{
		map[string]interface{}{
			"mode": "rbac",
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
		},
	}
}

func TestFlattenRKEClusterAuthorization(t *testing.T) {

	cases := []struct {
		Input          rancher.AuthzConfig
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterAuthorizationConf,
			testRKEClusterAuthorizationInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterAuthorization(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterAuthorization(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.AuthzConfig
	}{
		{
			testRKEClusterAuthorizationInterface,
			testRKEClusterAuthorizationConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterAuthorization(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
