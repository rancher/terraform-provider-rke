package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

var (
	testRKEClusterDNSNodelocalConf      *rancher.Nodelocal
	testRKEClusterDNSNodelocalInterface []interface{}
	testRKEClusterDNSConf               *rancher.DNSConfig
	testRKEClusterDNSInterface          []interface{}
)

func init() {
	testRKEClusterDNSNodelocalConf = &rancher.Nodelocal{
		NodeSelector: map[string]string{
			"sel1": "value1",
			"sel2": "value2",
		},
		IPAddress: "ip_address",
	}
	testRKEClusterDNSNodelocalInterface = []interface{}{
		map[string]interface{}{
			"node_selector": map[string]interface{}{
				"sel1": "value1",
				"sel2": "value2",
			},
			"ip_address": "ip_address",
		},
	}
	testRKEClusterDNSConf = &rancher.DNSConfig{
		Nodelocal: testRKEClusterDNSNodelocalConf,
		NodeSelector: map[string]string{
			"sel1": "value1",
			"sel2": "value2",
		},
		Provider:            "kube-dns",
		ReverseCIDRs:        []string{"rev1", "rev2"},
		UpstreamNameservers: []string{"up1", "up2"},
	}
	testRKEClusterDNSInterface = []interface{}{
		map[string]interface{}{
			"nodelocal": testRKEClusterDNSNodelocalInterface,
			"node_selector": map[string]interface{}{
				"sel1": "value1",
				"sel2": "value2",
			},
			"provider":             "kube-dns",
			"reverse_cidrs":        []interface{}{"rev1", "rev2"},
			"upstream_nameservers": []interface{}{"up1", "up2"},
		},
	}
}

func TestFlattenRKEClusterDNSNodelocal(t *testing.T) {

	cases := []struct {
		Input          *rancher.Nodelocal
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterDNSNodelocalConf,
			testRKEClusterDNSNodelocalInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterDNSNodelocal(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterDNS(t *testing.T) {

	cases := []struct {
		Input          *rancher.DNSConfig
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterDNSConf,
			testRKEClusterDNSInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterDNS(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterDNSNodelocal(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rancher.Nodelocal
	}{
		{
			testRKEClusterDNSNodelocalInterface,
			testRKEClusterDNSNodelocalConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterDNSNodelocal(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterDNS(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rancher.DNSConfig
	}{
		{
			testRKEClusterDNSInterface,
			testRKEClusterDNSConf,
		},
	}

	for _, tc := range cases {
		output := expandRKEClusterDNS(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
