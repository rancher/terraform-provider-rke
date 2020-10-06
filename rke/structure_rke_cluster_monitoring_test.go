package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/rke/types"
)

var (
	testRKEClusterMonitoringConf      rancher.MonitoringConfig
	testRKEClusterMonitoringInterface []interface{}
)

func init() {
	testRKEClusterMonitoringConf = rancher.MonitoringConfig{
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Provider: "test",
	}
	testRKEClusterMonitoringInterface = []interface{}{
		map[string]interface{}{
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
			"provider": "test",
		},
	}
}

func TestFlattenRKEClusterMonitoring(t *testing.T) {

	cases := []struct {
		Input          rancher.MonitoringConfig
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterMonitoringConf,
			testRKEClusterMonitoringInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterMonitoring(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterMonitoring(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.MonitoringConfig
	}{
		{
			testRKEClusterMonitoringInterface,
			testRKEClusterMonitoringConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterMonitoring(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
