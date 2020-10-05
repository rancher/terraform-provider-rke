package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/rke/types"
)

var (
	testRKEClusterServicesKubeControllerConf      rancher.KubeControllerService
	testRKEClusterServicesKubeControllerInterface []interface{}
)

func init() {
	testRKEClusterServicesKubeControllerConf = rancher.KubeControllerService{
		ClusterCIDR:           "10.42.0.0/16",
		ServiceClusterIPRange: "10.43.0.0/16",
	}
	testRKEClusterServicesKubeControllerConf.ExtraArgs = map[string]string{
		"arg_one": "one",
		"arg_two": "two",
	}
	testRKEClusterServicesKubeControllerConf.ExtraBinds = []string{"bind_one", "bind_two"}
	testRKEClusterServicesKubeControllerConf.ExtraEnv = []string{"env_one", "env_two"}
	testRKEClusterServicesKubeControllerConf.Image = "image"
	testRKEClusterServicesKubeControllerInterface = []interface{}{
		map[string]interface{}{
			"cluster_cidr": "10.42.0.0/16",
			"extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"extra_binds":              []interface{}{"bind_one", "bind_two"},
			"extra_env":                []interface{}{"env_one", "env_two"},
			"image":                    "image",
			"service_cluster_ip_range": "10.43.0.0/16",
		},
	}
}

func TestFlattenRKEClusterServicesKubeController(t *testing.T) {

	cases := []struct {
		Input          rancher.KubeControllerService
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterServicesKubeControllerConf,
			testRKEClusterServicesKubeControllerInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterServicesKubeController(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterServicesKubeController(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.KubeControllerService
	}{
		{
			testRKEClusterServicesKubeControllerInterface,
			testRKEClusterServicesKubeControllerConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterServicesKubeController(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
