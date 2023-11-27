package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/rke/types"
)

var (
	testRKEClusterServicesKubeproxyConf      rancher.KubeproxyService
	testRKEClusterServicesKubeproxyInterface []interface{}
)

func init() {
	testRKEClusterServicesKubeproxyConf = rancher.KubeproxyService{}
	testRKEClusterServicesKubeproxyConf.ExtraArgs = map[string]string{
		"arg_one": "one",
		"arg_two": "two",
	}
	testRKEClusterServicesKubeproxyConf.WindowsExtraArgs = map[string]string{
		"arg_one": "one",
		"arg_two": "two",
	}
	testRKEClusterServicesKubeproxyConf.ExtraArgsArray = map[string][]string{
		"arg1": {"v1", "v2"},
		"arg2": {"v1", "v2"},
	}
	testRKEClusterServicesKubeproxyConf.WindowsExtraArgsArray = map[string][]string{
		"arg1": {"v1", "v2"},
		"arg2": {"v1", "v2"},
	}
	testRKEClusterServicesKubeproxyConf.ExtraBinds = []string{"bind_one", "bind_two"}
	testRKEClusterServicesKubeproxyConf.ExtraEnv = []string{"env_one", "env_two"}
	testRKEClusterServicesKubeproxyConf.Image = "image"
	testRKEClusterServicesKubeproxyInterface = []interface{}{
		map[string]interface{}{
			"extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"windows_extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"extra_args_array":         "{\"arg1\":[\"v1\",\"v2\"],\"arg2\":[\"v1\",\"v2\"]}",
			"windows_extra_args_array": "{\"arg1\":[\"v1\",\"v2\"],\"arg2\":[\"v1\",\"v2\"]}",
			"extra_binds":              []interface{}{"bind_one", "bind_two"},
			"extra_env":                []interface{}{"env_one", "env_two"},
			"image":                    "image",
		},
	}
}

func TestFlattenRKEClusterServicesKubeproxy(t *testing.T) {

	cases := []struct {
		Input          rancher.KubeproxyService
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterServicesKubeproxyConf,
			testRKEClusterServicesKubeproxyInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenRKEClusterServicesKubeproxy(tc.Input)
		if err != nil {
			t.Fatalf("Unexpected error from flattener: %v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterServicesKubeproxy(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.KubeproxyService
	}{
		{
			testRKEClusterServicesKubeproxyInterface,
			testRKEClusterServicesKubeproxyConf,
		},
	}

	for _, tc := range cases {
		output, err := expandRKEClusterServicesKubeproxy(tc.Input)
		if err != nil {
			t.Fatalf("Unexpected error from expander: %v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
