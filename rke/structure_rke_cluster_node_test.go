package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

var (
	testRKEClusterNodesTaintsConf      []rancher.RKETaint
	testRKEClusterNodesTaintsInterface []interface{}
	testRKEClusterNodesConf            []rancher.RKEConfigNode
	testRKEClusterNodesInterface       []interface{}
)

func init() {
	testRKEClusterNodesTaintsConf = []rancher.RKETaint{
		{
			Key:    "key",
			Value:  "value",
			Effect: "recipient",
		},
	}
	testRKEClusterNodesTaintsInterface = []interface{}{
		map[string]interface{}{
			"key":    "key",
			"value":  "value",
			"effect": "recipient",
		},
	}
	testRKEClusterNodesConf = []rancher.RKEConfigNode{
		{
			Address:          "url.terraform.test",
			DockerSocket:     "docker.sock",
			HostnameOverride: "terra-test",
			InternalAddress:  "192.168.1.1",
			Labels: map[string]string{
				"label_one": "one",
				"label_two": "two",
			},
			NodeName:     "test1",
			Port:         "22",
			Role:         []string{"worker"},
			SSHAgentAuth: true,
			SSHCert:      "XXXXXXXX",
			SSHCertPath:  "/home/user/.ssh",
			SSHKey:       "XXXXXXXX",
			SSHKeyPath:   "/home/user/.ssh",
			User:         "test",
			Taints:       testRKEClusterNodesTaintsConf,
		},
	}
	testRKEClusterNodesInterface = []interface{}{
		map[string]interface{}{
			"address":           "url.terraform.test",
			"docker_socket":     "docker.sock",
			"hostname_override": "terra-test",
			"internal_address":  "192.168.1.1",
			"labels": map[string]interface{}{
				"label_one": "one",
				"label_two": "two",
			},
			"node_name":      "test1",
			"port":           "22",
			"role":           []interface{}{"worker"},
			"ssh_agent_auth": true,
			"ssh_cert":       "XXXXXXXX",
			"ssh_cert_path":  "/home/user/.ssh",
			"ssh_key":        "XXXXXXXX",
			"ssh_key_path":   "/home/user/.ssh",
			"user":           "test",
			"taints":         testRKEClusterNodesTaintsInterface,
		},
	}
}

func TestFlattenRKEClusterNodes(t *testing.T) {

	cases := []struct {
		Input          []rancher.RKEConfigNode
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterNodesConf,
			testRKEClusterNodesInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterNodes(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterNodes(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []rancher.RKEConfigNode
	}{
		{
			testRKEClusterNodesInterface,
			testRKEClusterNodesConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterNodes(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
