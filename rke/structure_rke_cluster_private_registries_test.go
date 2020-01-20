package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

var (
	testRKEClusterPrivateRegistriesConf      []rancher.PrivateRegistry
	testRKEClusterPrivateRegistriesInterface []interface{}
)

func init() {
	testRKEClusterPrivateRegistriesConf = []rancher.PrivateRegistry{
		{
			IsDefault: true,
			Password:  "XXXXXXXX",
			URL:       "url.terraform.test",
			User:      "user",
		},
	}
	testRKEClusterPrivateRegistriesInterface = []interface{}{
		map[string]interface{}{
			"is_default": true,
			"password":   "XXXXXXXX",
			"url":        "url.terraform.test",
			"user":       "user",
		},
	}
}

func TestFlattenPrivateRegistries(t *testing.T) {

	cases := []struct {
		Input          []rancher.PrivateRegistry
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterPrivateRegistriesConf,
			testRKEClusterPrivateRegistriesInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterPrivateRegistries(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPrivateRegistries(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []rancher.PrivateRegistry
	}{
		{
			testRKEClusterPrivateRegistriesInterface,
			testRKEClusterPrivateRegistriesConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterPrivateRegistries(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
