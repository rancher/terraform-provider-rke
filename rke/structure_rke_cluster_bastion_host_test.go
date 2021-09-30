package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/rke/types"
)

var (
	testRKEClusterBastionHostConf      rancher.BastionHost
	testRKEClusterBastionHostInterface []interface{}
)

func init() {
	testRKEClusterBastionHostConf = rancher.BastionHost{
		Address:      "bastion.terraform.test",
		Port:         "22",
		SSHAgentAuth: true,
		SSHCert:      "XXXXXXXX",
		SSHCertPath:  "/home/user/.ssh",
		SSHKey:       "XXXXXXXX",
		SSHKeyPath:   "/home/user/.ssh",
		User:         "test",
	}
	testRKEClusterBastionHostInterface = []interface{}{
		map[string]interface{}{
			"address":        "bastion.terraform.test",
			"port":           "22",
			"ssh_agent_auth": true,
			"ssh_cert":       "XXXXXXXX",
			"ssh_cert_path":  "/home/user/.ssh",
			"ssh_key":        "XXXXXXXX",
			"ssh_key_path":   "/home/user/.ssh",
			"user":           "test",
		},
	}
}

func TestFlattenRKEClusterBastionHost(t *testing.T) {

	cases := []struct {
		Input          rancher.BastionHost
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterBastionHostConf,
			testRKEClusterBastionHostInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterBastionHost(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterBastionHost(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.BastionHost
	}{
		{
			testRKEClusterBastionHostInterface,
			testRKEClusterBastionHostConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterBastionHost(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
