package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/rke/types"
)

var (
	testRKEClusterServicesKubeletConf      rancher.KubeletService
	testRKEClusterServicesKubeletInterface []interface{}
)

func init() {
	testRKEClusterServicesKubeletConf = rancher.KubeletService{
		ClusterDNSServer:           "dns.hostname.test",
		ClusterDomain:              "terraform.test",
		FailSwapOn:                 true,
		GenerateServingCertificate: true,
		InfraContainerImage:        "infra_image",
	}
	testRKEClusterServicesKubeletConf.ExtraArgs = map[string]string{
		"arg_one": "one",
		"arg_two": "two",
	}
	testRKEClusterServicesKubeletConf.WindowsExtraArgs = map[string]string{
		"arg_one": "one",
		"arg_two": "two",
	}
	testRKEClusterServicesKubeletConf.ExtraArgsArray = map[string][]string{
		"arg1": {"v1", "v2"},
		"arg2": {"v1", "v2"},
	}
	testRKEClusterServicesKubeletConf.WindowsExtraArgsArray = map[string][]string{
		"arg1": {"v1", "v2"},
		"arg2": {"v1", "v2"},
	}
	testRKEClusterServicesKubeletConf.ExtraBinds = []string{"bind_one", "bind_two"}
	testRKEClusterServicesKubeletConf.ExtraEnv = []string{"env_one", "env_two"}
	testRKEClusterServicesKubeletConf.Image = "image"
	testRKEClusterServicesKubeletInterface = []interface{}{
		map[string]interface{}{
			"cluster_dns_server": "dns.hostname.test",
			"cluster_domain":     "terraform.test",
			"extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"windows_extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"extra_args_array":             "{\"arg1\":[\"v1\",\"v2\"],\"arg2\":[\"v1\",\"v2\"]}",
			"windows_extra_args_array":     "{\"arg1\":[\"v1\",\"v2\"],\"arg2\":[\"v1\",\"v2\"]}",
			"extra_binds":                  []interface{}{"bind_one", "bind_two"},
			"extra_env":                    []interface{}{"env_one", "env_two"},
			"fail_swap_on":                 true,
			"generate_serving_certificate": true,
			"image":                        "image",
			"infra_container_image":        "infra_image",
		},
	}
}

func TestFlattenRKEClusterServicesKubelet(t *testing.T) {

	cases := []struct {
		Input          rancher.KubeletService
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterServicesKubeletConf,
			testRKEClusterServicesKubeletInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenRKEClusterServicesKubelet(tc.Input)
		if err != nil {
			t.Fatalf("Unexpected error from flattener: %v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterServicesKubelet(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.KubeletService
	}{
		{
			testRKEClusterServicesKubeletInterface,
			testRKEClusterServicesKubeletConf,
		},
	}

	for _, tc := range cases {
		output, err := expandRKEClusterServicesKubelet(tc.Input)
		if err != nil {
			t.Fatalf("Unexpected error from expander: %v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
