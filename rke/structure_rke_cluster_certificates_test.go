package rke

import (
	"reflect"
	"testing"

	"github.com/rancher/rke/pki"
)

var (
	testRKEClusterCertificatesConf      map[string]pki.CertificatePKI
	testRKEClusterCertificatesInterface []interface{}
)

func init() {
	testRKEClusterCertificatesConf = map[string]pki.CertificatePKI{
		"test": {
			CertificatePEM: "certificate",
			KeyPEM:         "key",
			Config:         "config",
			Name:           "name",
			CommonName:     "common_name",
			OUName:         "ou_name",
			EnvName:        "env_name",
			Path:           "path",
			KeyEnvName:     "key_env_name",
			KeyPath:        "key_path",
			ConfigEnvName:  "config_env_name",
			ConfigPath:     "config_path",
		},
	}
	testRKEClusterCertificatesInterface = []interface{}{
		map[string]interface{}{
			"id":              "test",
			"certificate":     "certificate",
			"key":             "key",
			"config":          "config",
			"name":            "name",
			"common_name":     "common_name",
			"ou_name":         "ou_name",
			"env_name":        "env_name",
			"path":            "path",
			"key_env_name":    "key_env_name",
			"key_path":        "key_path",
			"config_env_name": "config_env_name",
			"config_path":     "config_path",
		},
	}
}

func TestFlattenRKEClusterCertificates(t *testing.T) {

	cases := []struct {
		Input          map[string]pki.CertificatePKI
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterCertificatesConf,
			testRKEClusterCertificatesInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRKEClusterCertificates(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
