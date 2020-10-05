package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/rke/types"
)

var (
	testRKEClusterTaintsConf      []rancher.RKETaint
	testRKEClusterTaintsInterface []interface{}
)

func init() {
	testRKEClusterTaintsConf = []rancher.RKETaint{
		{
			Key:    "key",
			Value:  "value",
			Effect: "recipient",
		},
	}
	testRKEClusterTaintsInterface = []interface{}{
		map[string]interface{}{
			"key":    "key",
			"value":  "value",
			"effect": "recipient",
		},
	}
}

func TestFlattenRKEClusterTaints(t *testing.T) {

	cases := []struct {
		Input          []rancher.RKETaint
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterTaintsConf,
			testRKEClusterTaintsInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterTaints(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterTaints(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []rancher.RKETaint
	}{
		{
			testRKEClusterTaintsInterface,
			testRKEClusterTaintsConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterTaints(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
