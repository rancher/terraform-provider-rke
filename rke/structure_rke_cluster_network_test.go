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
	testRKEClusterNetworkAciConf          *rancher.AciNetworkProvider
	testRKEClusterNetworkAciInterface     []interface{}
	testRKEClusterNetworkConfCalico       rancher.NetworkConfig
	testRKEClusterNetworkInterfaceCalico  []interface{}
	testRKEClusterNetworkConfCanal        rancher.NetworkConfig
	testRKEClusterNetworkInterfaceCanal   []interface{}
	testRKEClusterNetworkConfFlannel      rancher.NetworkConfig
	testRKEClusterNetworkInterfaceFlannel []interface{}
	testRKEClusterNetworkConfWeave        rancher.NetworkConfig
	testRKEClusterNetworkInterfaceWeave   []interface{}
	testRKEClusterNetworkConfAci          rancher.NetworkConfig
	testRKEClusterNetworkInterfaceAci     []interface{}
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
	testRKEClusterNetworkAciConf = &rancher.AciNetworkProvider{
		SystemIdentifier:      "demo",
		ApicHosts:             []string{"10.80.20.180"},
		Token:                 "a819d237-25f1-41da-91a4-985e87cef864",
		ApicUserName:          "admin",
		ApicUserKey:           "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1J",
		ApicUserCrt:           "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JS",
		EncapType:             "vxlan",
		McastRangeStart:       "225.39.1.1",
		McastRangeEnd:         "225.22.255.255",
		AEP:                   "hypf-aep",
		VRFName:               "k8s12_vrf",
		VRFTenant:             "common",
		L3Out:                 "k8s12",
		L3OutExternalNetworks: []string{"10.80.20.180"},
		NodeSubnet:            "192.168.12.1/24",
		DynamicExternalSubnet: "10.3.0.1/24",
		StaticExternalSubnet:  "10.4.0.1/24",
		ServiceGraphSubnet:    "10.5.0.1/24",
		KubeAPIVlan:           "23",
		ServiceVlan:           "24",
		InfraVlan:             "3901",
		SnatPortRangeStart:    "5000",
		SnatPortRangeEnd:      "65000",
		SnatPortsPerNode:      "3000",
	}
	testRKEClusterNetworkAciInterface = []interface{}{
		map[string]interface{}{
			"system_id":               "demo",
			"apic_hosts":              []interface{}{"10.80.20.180"},
			"token":                   "a819d237-25f1-41da-91a4-985e87cef864",
			"apic_user_name":          "admin",
			"apic_user_key":           "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1J",
			"apic_user_crt":           "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JS",
			"encap_type":              "vxlan",
			"mcast_range_start":       "225.39.1.1",
			"mcast_range_end":         "225.22.255.255",
			"aep":                     "hypf-aep",
			"vrf_name":                "k8s12_vrf",
			"vrf_tenant":              "common",
			"l3out":                   "k8s12",
			"l3out_external_networks": []interface{}{"10.80.20.180"},
			"node_subnet":             "192.168.12.1/24",
			"extern_dynamic":          "10.3.0.1/24",
			"extern_static":           "10.4.0.1/24",
			"node_svc_subnet":         "10.5.0.1/24",
			"kube_api_vlan":           "23",
			"service_vlan":            "24",
			"infra_vlan":              "3901",
			"snat_port_range_start":   "5000",
			"snat_port_range_end":     "65000",
			"snat_ports_per_node":     "3000",
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
	testRKEClusterNetworkConfAci = rancher.NetworkConfig{
		AciNetworkProvider: testRKEClusterNetworkAciConf,
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
			"option3": "value3",
		},
		Plugin: rkeClusterNetworkPluginAciName,
	}
	testRKEClusterNetworkInterfaceAci = []interface{}{
		map[string]interface{}{
			"aci_network_provider": testRKEClusterNetworkAciInterface,
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
				"option3": "value3",
			},
			"plugin": rkeClusterNetworkPluginAciName,
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

func TestFlattenRKEClusterNetworkAci(t *testing.T) {

	cases := []struct {
		Input          *rancher.AciNetworkProvider
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterNetworkAciConf,
			testRKEClusterNetworkAciInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterNetworkAci(tc.Input)
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
		{
			testRKEClusterNetworkConfAci,
			testRKEClusterNetworkInterfaceAci,
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

func TestExpandRKEClusterNetworkAci(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rancher.AciNetworkProvider
	}{
		{
			testRKEClusterNetworkAciInterface,
			testRKEClusterNetworkAciConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterNetworkAci(tc.Input)
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
		{
			testRKEClusterNetworkInterfaceAci,
			testRKEClusterNetworkConfAci,
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
