package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

var (
	testRKEClusterRotateCertificatesConf      *rancher.RotateCertificates
	testRKEClusterRotateCertificatesInterface []interface{}
)

func init() {
	testRKEClusterRotateCertificatesConf = &rancher.RotateCertificates{
		CACertificates: true,
		Services:       []string{"serv1", "serv2"},
	}
	testRKEClusterRotateCertificatesInterface = []interface{}{
		map[string]interface{}{
			"ca_certificates": true,
			"services":        []interface{}{"serv1", "serv2"},
		},
	}
}

func TestFlattenRKEClusterRotateCertificates(t *testing.T) {

	cases := []struct {
		Input          *rancher.RotateCertificates
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterRotateCertificatesConf,
			testRKEClusterRotateCertificatesInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterRotateCertificates(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterRotateCertificates(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rancher.RotateCertificates
	}{
		{
			testRKEClusterRotateCertificatesInterface,
			testRKEClusterRotateCertificatesConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterRotateCertificates(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
