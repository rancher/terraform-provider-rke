package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/rke/types"
)

var (
	testRKEClusterNetworkCalicoConf       *rancher.CalicoNetworkProvider
	testRKEClusterNetworkCalicoInterface  []interface{}
	testRKEClusterNetworkCanalConf        *rancher.CanalNetworkProvider
	testRKEClusterNetworkCanalInterface   []interface{}
	testRKEClusterNetworkFlannelConf      *rancher.FlannelNetworkProvider
	testRKEClusterNetworkFlannelInterface []interface{}
	testRKEClusterNetworkWeaveConf        *rancher.WeaveNetworkProvider
	testRKEClusterNetworkWeaveInterface   []interface{}
	testRKEClusterNetworkConfCalico       rancher.NetworkConfig
	testRKEClusterNetworkInterfaceCalico  []interface{}
	testRKEClusterNetworkConfCanal        rancher.NetworkConfig
	testRKEClusterNetworkInterfaceCanal   []interface{}
	testRKEClusterNetworkConfFlannel      rancher.NetworkConfig
	testRKEClusterNetworkInterfaceFlannel []interface{}
	testRKEClusterNetworkConfWeave        rancher.NetworkConfig
	testRKEClusterNetworkInterfaceWeave   []interface{}
)

func init() {
	testRKEClusterNetworkCalicoConf = &rancher.CalicoNetworkProvider{
		CloudProvider: "aws",
	}
	testRKEClusterNetworkCalicoInterface = []interface{}{
		map[string]interface{}{
			"cloud_provider": "aws",
		},
	}
	testRKEClusterNetworkCanalConf = &rancher.CanalNetworkProvider{
		FlannelNetworkProvider: rancher.FlannelNetworkProvider{
			Iface: "eth0",
		},
	}
	testRKEClusterNetworkCanalInterface = []interface{}{
		map[string]interface{}{
			"iface": "eth0",
		},
	}
	testRKEClusterNetworkFlannelConf = &rancher.FlannelNetworkProvider{
		Iface: "eth0",
	}
	testRKEClusterNetworkFlannelInterface = []interface{}{
		map[string]interface{}{
			"iface": "eth0",
		},
	}
	testRKEClusterNetworkWeaveConf = &rancher.WeaveNetworkProvider{
		Password: "password",
	}
	testRKEClusterNetworkWeaveInterface = []interface{}{
		map[string]interface{}{
			"password": "password",
		},
	}
	testRKEClusterNetworkConfCalico = rancher.NetworkConfig{
		CalicoNetworkProvider: testRKEClusterNetworkCalicoConf,
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Plugin: rkeClusterNetworkPluginCalicoName,
	}
	testRKEClusterNetworkInterfaceCalico = []interface{}{
		map[string]interface{}{
			"calico_network_provider": testRKEClusterNetworkCalicoInterface,
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
			"plugin": rkeClusterNetworkPluginCalicoName,
		},
	}
	testRKEClusterNetworkConfCanal = rancher.NetworkConfig{
		CanalNetworkProvider: testRKEClusterNetworkCanalConf,
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Plugin: rkeClusterNetworkPluginCanalName,
	}
	testRKEClusterNetworkInterfaceCanal = []interface{}{
		map[string]interface{}{
			"canal_network_provider": testRKEClusterNetworkCanalInterface,
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
			"plugin": rkeClusterNetworkPluginCanalName,
		},
	}
	testRKEClusterNetworkConfFlannel = rancher.NetworkConfig{
		FlannelNetworkProvider: testRKEClusterNetworkFlannelConf,
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Plugin: rkeClusterNetworkPluginFlannelName,
	}
	testRKEClusterNetworkInterfaceFlannel = []interface{}{
		map[string]interface{}{
			"flannel_network_provider": testRKEClusterNetworkFlannelInterface,
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
			"plugin": rkeClusterNetworkPluginFlannelName,
		},
	}
	testRKEClusterNetworkConfWeave = rancher.NetworkConfig{
		WeaveNetworkProvider: testRKEClusterNetworkWeaveConf,
		MTU:                  1500,
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Plugin: rkeClusterNetworkPluginWeaveName,
	}
	testRKEClusterNetworkInterfaceWeave = []interface{}{
		map[string]interface{}{
			"weave_network_provider": testRKEClusterNetworkWeaveInterface,
			"mtu":                    1500,
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
			"plugin": rkeClusterNetworkPluginWeaveName,
		},
	}
}

func TestFlattenRKEClusterNetworkCalico(t *testing.T) {

	cases := []struct {
		Input          *rancher.CalicoNetworkProvider
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterNetworkCalicoConf,
			testRKEClusterNetworkCalicoInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterNetworkCalico(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterNetworkCanal(t *testing.T) {

	cases := []struct {
		Input          *rancher.CanalNetworkProvider
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterNetworkCanalConf,
			testRKEClusterNetworkCanalInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterNetworkCanal(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterNetworkFlannel(t *testing.T) {

	cases := []struct {
		Input          *rancher.FlannelNetworkProvider
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterNetworkFlannelConf,
			testRKEClusterNetworkFlannelInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterNetworkFlannel(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterNetworkWeave(t *testing.T) {

	cases := []struct {
		Input          *rancher.WeaveNetworkProvider
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterNetworkWeaveConf,
			testRKEClusterNetworkWeaveInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterNetworkWeave(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterNetwork(t *testing.T) {

	cases := []struct {
		Input          rancher.NetworkConfig
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterNetworkConfCalico,
			testRKEClusterNetworkInterfaceCalico,
		},
		{
			testRKEClusterNetworkConfCanal,
			testRKEClusterNetworkInterfaceCanal,
		},
		{
			testRKEClusterNetworkConfFlannel,
			testRKEClusterNetworkInterfaceFlannel,
		},
		{
			testRKEClusterNetworkConfWeave,
			testRKEClusterNetworkInterfaceWeave,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterNetwork(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterNetworkCalico(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rancher.CalicoNetworkProvider
	}{
		{
			testRKEClusterNetworkCalicoInterface,
			testRKEClusterNetworkCalicoConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterNetworkCalico(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterNetworkCanal(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rancher.CanalNetworkProvider
	}{
		{
			testRKEClusterNetworkCanalInterface,
			testRKEClusterNetworkCanalConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterNetworkCanal(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterNetworkFlannel(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rancher.FlannelNetworkProvider
	}{
		{
			testRKEClusterNetworkFlannelInterface,
			testRKEClusterNetworkFlannelConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterNetworkFlannel(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterNetworkWeave(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rancher.WeaveNetworkProvider
	}{
		{
			testRKEClusterNetworkWeaveInterface,
			testRKEClusterNetworkWeaveConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterNetworkWeave(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterNetwork(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.NetworkConfig
	}{
		{
			testRKEClusterNetworkInterfaceCalico,
			testRKEClusterNetworkConfCalico,
		},
		{
			testRKEClusterNetworkInterfaceCanal,
			testRKEClusterNetworkConfCanal,
		},
		{
			testRKEClusterNetworkInterfaceFlannel,
			testRKEClusterNetworkConfFlannel,
		},
		{
			testRKEClusterNetworkInterfaceWeave,
			testRKEClusterNetworkConfWeave,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterNetwork(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
