package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/rke/types"
)

var (
	testRKEClusterServicesSchedulerConf      rancher.SchedulerService
	testRKEClusterServicesSchedulerInterface []interface{}
)

func init() {
	testRKEClusterServicesSchedulerConf = rancher.SchedulerService{}
	testRKEClusterServicesSchedulerConf.ExtraArgs = map[string]string{
		"arg_one": "one",
		"arg_two": "two",
	}
	testRKEClusterServicesSchedulerConf.WindowsExtraArgs = map[string]string{
		"arg_one": "one",
		"arg_two": "two",
	}
	testRKEClusterServicesSchedulerConf.ExtraArgsArray = map[string][]string{
		"arg1": {"v1"},
		"arg2": {"v2"},
	}
	testRKEClusterServicesSchedulerConf.WindowsExtraArgsArray = map[string][]string{
		"arg1": {"v1"},
		"arg2": {"v2"},
	}
	testRKEClusterServicesSchedulerConf.ExtraBinds = []string{"bind_one", "bind_two"}
	testRKEClusterServicesSchedulerConf.ExtraEnv = []string{"env_one", "env_two"}
	testRKEClusterServicesSchedulerConf.Image = "image"
	testRKEClusterServicesSchedulerInterface = []interface{}{
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

func TestFlattenRKEClusterServicesScheduler(t *testing.T) {

	cases := []struct {
		Input          rancher.SchedulerService
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterServicesSchedulerConf,
			testRKEClusterServicesSchedulerInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterServicesScheduler(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterServicesScheduler(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.SchedulerService
	}{
		{
			testRKEClusterServicesSchedulerInterface,
			testRKEClusterServicesSchedulerConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterServicesScheduler(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
