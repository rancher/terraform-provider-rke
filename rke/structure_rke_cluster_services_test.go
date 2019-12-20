package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

var (
	testRKEClusterServicesConf      rancher.RKEConfigServices
	testRKEClusterServicesInterface []interface{}
)

func init() {
	testRKEClusterServicesConf = rancher.RKEConfigServices{
		Etcd:           testRKEClusterServicesETCDConf,
		KubeAPI:        testRKEClusterServicesKubeAPIConf,
		KubeController: testRKEClusterServicesKubeControllerConf,
		Kubelet:        testRKEClusterServicesKubeletConf,
		Kubeproxy:      testRKEClusterServicesKubeproxyConf,
		Scheduler:      testRKEClusterServicesSchedulerConf,
	}
	testRKEClusterServicesInterface = []interface{}{
		map[string]interface{}{
			"etcd":            testRKEClusterServicesETCDInterface,
			"kube_api":        testRKEClusterServicesKubeAPIInterface,
			"kube_controller": testRKEClusterServicesKubeControllerInterface,
			"kubelet":         testRKEClusterServicesKubeletInterface,
			"kubeproxy":       testRKEClusterServicesKubeproxyInterface,
			"scheduler":       testRKEClusterServicesSchedulerInterface,
		},
	}
}

func TestFlattenRKEClusterServices(t *testing.T) {

	cases := []struct {
		Input          rancher.RKEConfigServices
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterServicesConf,
			testRKEClusterServicesInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterServices(tc.Input, testRKEClusterServicesInterface)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterServices(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.RKEConfigServices
	}{
		{
			testRKEClusterServicesInterface,
			testRKEClusterServicesConf,
		},
	}

	for _, tc := range cases {
		output, err := expandRKEClusterServices(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
