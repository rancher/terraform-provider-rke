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
		"arg1": {"v1"},
		"arg2": {"v2"},
	}
	testRKEClusterServicesKubeproxyConf.WindowsExtraArgsArray = map[string][]string{
		"arg1": {"v1"},
		"arg2": {"v2"},
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
			"win_extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"extra_args_array": []interface{}{
				map[string]interface{}{
					"extra_arg": []interface{}{
						map[string]interface{}{
							"argument": "arg1",
							"values":   []interface{}{"v1"},
						},
						map[string]interface{}{
							"argument": "arg2",
							"values":   []interface{}{"v2"},
						},
					},
				},
			},
			"win_extra_args_array": []interface{}{
				map[string]interface{}{
					"extra_arg": []interface{}{
						map[string]interface{}{
							"argument": "arg1",
							"values":   []interface{}{"v1"},
						},
						map[string]interface{}{
							"argument": "arg2",
							"values":   []interface{}{"v2"},
						},
					},
				},
			},
			"extra_binds": []interface{}{"bind_one", "bind_two"},
			"extra_env":   []interface{}{"env_one", "env_two"},
			"image":       "image",
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
		output := flattenRKEClusterServicesKubeproxy(tc.Input)
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
		output := expandRKEClusterServicesKubeproxy(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
