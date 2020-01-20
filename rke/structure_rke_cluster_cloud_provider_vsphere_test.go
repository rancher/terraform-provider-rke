package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

var (
	testRKEClusterCloudProviderVsphereDiskConf               rancher.DiskVsphereOpts
	testRKEClusterCloudProviderVsphereDiskInterface          []interface{}
	testRKEClusterCloudProviderVsphereGlobalConf             rancher.GlobalVsphereOpts
	testRKEClusterCloudProviderVsphereGlobalInterface        []interface{}
	testRKEClusterCloudProviderVsphereNetworkConf            rancher.NetworkVshpereOpts
	testRKEClusterCloudProviderVsphereNetworkInterface       []interface{}
	testRKEClusterCloudProviderVsphereVirtualCenterConf      map[string]rancher.VirtualCenterConfig
	testRKEClusterCloudProviderVsphereVirtualCenterInterface []interface{}
	testRKEClusterCloudProviderVsphereWorkspaceConf          rancher.WorkspaceVsphereOpts
	testRKEClusterCloudProviderVsphereWorkspaceInterface     []interface{}
	testRKEClusterCloudProviderVsphereConf                   *rancher.VsphereCloudProvider
	testRKEClusterCloudProviderVsphereInterface              []interface{}
)

func init() {
	testRKEClusterCloudProviderVsphereDiskConf = rancher.DiskVsphereOpts{
		SCSIControllerType: "test",
	}
	testRKEClusterCloudProviderVsphereDiskInterface = []interface{}{
		map[string]interface{}{
			"scsi_controller_type": "test",
		},
	}
	testRKEClusterCloudProviderVsphereGlobalConf = rancher.GlobalVsphereOpts{
		Datacenters:       "auth.terraform.test",
		InsecureFlag:      true,
		Password:          "YYYYYYYY",
		VCenterPort:       "123",
		User:              "user",
		RoundTripperCount: 10,
	}
	testRKEClusterCloudProviderVsphereGlobalInterface = []interface{}{
		map[string]interface{}{
			"datacenters":          "auth.terraform.test",
			"insecure_flag":        true,
			"password":             "YYYYYYYY",
			"port":                 "123",
			"user":                 "user",
			"soap_roundtrip_count": 10,
		},
	}
	testRKEClusterCloudProviderVsphereNetworkConf = rancher.NetworkVshpereOpts{
		PublicNetwork: "test",
	}
	testRKEClusterCloudProviderVsphereNetworkInterface = []interface{}{
		map[string]interface{}{
			"public_network": "test",
		},
	}
	testRKEClusterCloudProviderVsphereVirtualCenterConf = map[string]rancher.VirtualCenterConfig{
		"test": {
			Datacenters:       "auth.terraform.test",
			Password:          "YYYYYYYY",
			VCenterPort:       "123",
			User:              "user",
			RoundTripperCount: 10,
		},
	}
	testRKEClusterCloudProviderVsphereVirtualCenterInterface = []interface{}{
		map[string]interface{}{
			"name":                 "test",
			"datacenters":          "auth.terraform.test",
			"password":             "YYYYYYYY",
			"port":                 "123",
			"user":                 "user",
			"soap_roundtrip_count": 10,
		},
	}
	testRKEClusterCloudProviderVsphereWorkspaceConf = rancher.WorkspaceVsphereOpts{
		Datacenter:       "test",
		Folder:           "folder",
		VCenterIP:        "vcenter",
		DefaultDatastore: "datastore",
		ResourcePoolPath: "resourcepool",
	}
	testRKEClusterCloudProviderVsphereWorkspaceInterface = []interface{}{
		map[string]interface{}{
			"datacenter":        "test",
			"folder":            "folder",
			"server":            "vcenter",
			"default_datastore": "datastore",
			"resourcepool_path": "resourcepool",
		},
	}
	testRKEClusterCloudProviderVsphereConf = &rancher.VsphereCloudProvider{
		Disk:          testRKEClusterCloudProviderVsphereDiskConf,
		Global:        testRKEClusterCloudProviderVsphereGlobalConf,
		Network:       testRKEClusterCloudProviderVsphereNetworkConf,
		VirtualCenter: testRKEClusterCloudProviderVsphereVirtualCenterConf,
		Workspace:     testRKEClusterCloudProviderVsphereWorkspaceConf,
	}
	testRKEClusterCloudProviderVsphereInterface = []interface{}{
		map[string]interface{}{
			"disk":           testRKEClusterCloudProviderVsphereDiskInterface,
			"global":         testRKEClusterCloudProviderVsphereGlobalInterface,
			"network":        testRKEClusterCloudProviderVsphereNetworkInterface,
			"virtual_center": testRKEClusterCloudProviderVsphereVirtualCenterInterface,
			"workspace":      testRKEClusterCloudProviderVsphereWorkspaceInterface,
		},
	}
}

func TestFlattenRKEClusterCloudProviderVsphereDisk(t *testing.T) {

	cases := []struct {
		Input          rancher.DiskVsphereOpts
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterCloudProviderVsphereDiskConf,
			testRKEClusterCloudProviderVsphereDiskInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterCloudProviderVsphereDisk(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterCloudProviderVsphereGlobal(t *testing.T) {

	cases := []struct {
		Input          rancher.GlobalVsphereOpts
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterCloudProviderVsphereGlobalConf,
			testRKEClusterCloudProviderVsphereGlobalInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterCloudProviderVsphereGlobal(tc.Input, testRKEClusterCloudProviderVsphereGlobalInterface)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterCloudProviderVsphereNetwork(t *testing.T) {

	cases := []struct {
		Input          rancher.NetworkVshpereOpts
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterCloudProviderVsphereNetworkConf,
			testRKEClusterCloudProviderVsphereNetworkInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterCloudProviderVsphereNetwork(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterCloudProviderVsphereVirtualCenter(t *testing.T) {

	cases := []struct {
		Input          map[string]rancher.VirtualCenterConfig
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterCloudProviderVsphereVirtualCenterConf,
			testRKEClusterCloudProviderVsphereVirtualCenterInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterCloudProviderVsphereVirtualCenter(tc.Input, testRKEClusterCloudProviderVsphereVirtualCenterInterface)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterCloudProviderVsphereWorkspace(t *testing.T) {

	cases := []struct {
		Input          rancher.WorkspaceVsphereOpts
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterCloudProviderVsphereWorkspaceConf,
			testRKEClusterCloudProviderVsphereWorkspaceInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterCloudProviderVsphereWorkspace(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterCloudProviderVsphere(t *testing.T) {

	cases := []struct {
		Input          *rancher.VsphereCloudProvider
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterCloudProviderVsphereConf,
			testRKEClusterCloudProviderVsphereInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterCloudProviderVsphere(tc.Input, testRKEClusterCloudProviderVsphereInterface)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterCloudProviderVsphereDisk(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.DiskVsphereOpts
	}{
		{
			testRKEClusterCloudProviderVsphereDiskInterface,
			testRKEClusterCloudProviderVsphereDiskConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterCloudProviderVsphereDisk(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterCloudProviderVsphereGlobal(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.GlobalVsphereOpts
	}{
		{
			testRKEClusterCloudProviderVsphereGlobalInterface,
			testRKEClusterCloudProviderVsphereGlobalConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterCloudProviderVsphereGlobal(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterCloudProviderVsphereNetwork(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.NetworkVshpereOpts
	}{
		{
			testRKEClusterCloudProviderVsphereNetworkInterface,
			testRKEClusterCloudProviderVsphereNetworkConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterCloudProviderVsphereNetwork(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterCloudProviderVsphereVirtualCenter(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput map[string]rancher.VirtualCenterConfig
	}{
		{
			testRKEClusterCloudProviderVsphereVirtualCenterInterface,
			testRKEClusterCloudProviderVsphereVirtualCenterConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterCloudProviderVsphereVirtualCenter(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterCloudProviderVsphereWorkspace(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.WorkspaceVsphereOpts
	}{
		{
			testRKEClusterCloudProviderVsphereWorkspaceInterface,
			testRKEClusterCloudProviderVsphereWorkspaceConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterCloudProviderVsphereWorkspace(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterCloudProviderVsphere(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rancher.VsphereCloudProvider
	}{
		{
			testRKEClusterCloudProviderVsphereInterface,
			testRKEClusterCloudProviderVsphereConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterCloudProviderVsphere(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
