package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

var (
	testRKEClusterIngressConf      rancher.IngressConfig
	testRKEClusterIngressInterface []interface{}
)

func init() {
	testRKEClusterIngressConf = rancher.IngressConfig{
		DNSPolicy: "test",
		ExtraArgs: map[string]string{
			"arg_one": "one",
			"arg_two": "two",
		},
		NodeSelector: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Provider: "test",
	}
	testRKEClusterIngressInterface = []interface{}{
		map[string]interface{}{
			"dns_policy": "test",
			"extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"node_selector": map[string]interface{}{
				"node_one": "one",
				"node_two": "two",
			},
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
			"provider": "test",
		},
	}
}

func TestFlattenRKEClusterIngress(t *testing.T) {

	cases := []struct {
		Input          rancher.IngressConfig
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterIngressConf,
			testRKEClusterIngressInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterIngress(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterIngress(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.IngressConfig
	}{
		{
			testRKEClusterIngressInterface,
			testRKEClusterIngressConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterIngress(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
