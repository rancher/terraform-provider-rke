package rke

import (
	"sort"

	"github.com/rancher/rke/pki"
)

const (
	rkeClusterCertificatesKubeAdminCertName = pki.KubeAdminCertName
)

// Flatteners

func flattenRKEClusterCertificates(in map[string]pki.CertificatePKI) (string, string, string, []interface{}) {
	var caCrt, clientCrt, clientKey string
	outLen := len(in)
	if in == nil || outLen == 0 {
		return caCrt, clientCrt, clientKey, []interface{}{}
	}
	sortedKeys := make([]string, 0, outLen)
	for k := range in {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)
	out := make([]interface{}, outLen)
	for i, k := range sortedKeys {
		v := in[k]

		if k == pki.CACertName {
			caCrt = v.CertificatePEM
		}

		if k == pki.KubeAdminCertName {
			clientCrt = v.CertificatePEM
			clientKey = v.KeyPEM
		}

		obj := map[string]interface{}{
			"id":              k,
			"certificate":     v.CertificatePEM,
			"key":             v.KeyPEM,
			"config":          v.Config,
			"name":            v.Name,
			"common_name":     v.CommonName,
			"ou_name":         v.OUName,
			"env_name":        v.EnvName,
			"path":            v.Path,
			"key_env_name":    v.KeyEnvName,
			"key_path":        v.KeyPath,
			"config_env_name": v.ConfigEnvName,
			"config_path":     v.ConfigPath,
		}
		out[i] = obj
	}
	return caCrt, clientCrt, clientKey, out
}
