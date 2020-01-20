package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

var (
	testRKEClusterRestoreConf      rancher.RestoreConfig
	testRKEClusterRestoreInterface []interface{}
)

func init() {
	testRKEClusterRestoreConf = rancher.RestoreConfig{
		Restore:      true,
		SnapshotName: "snapshot_name",
	}
	testRKEClusterRestoreInterface = []interface{}{
		map[string]interface{}{
			"restore":       true,
			"snapshot_name": "snapshot_name",
		},
	}
}

func TestFlattenRKEClusterRestore(t *testing.T) {

	cases := []struct {
		Input          rancher.RestoreConfig
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterRestoreConf,
			testRKEClusterRestoreInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterRestore(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterRestore(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.RestoreConfig
	}{
		{
			testRKEClusterRestoreInterface,
			testRKEClusterRestoreConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterRestore(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
