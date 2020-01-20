package rke

import (
	"github.com/rancher/rke/pki"
)

const (
	rkeClusterCertificatesKubeAdminCertName = pki.KubeAdminCertName
)

// Flatteners

func flattenRKEClusterCertificates(in map[string]pki.CertificatePKI) []interface{} {
	out := []interface{}{}

	for k, v := range in {
		/*certPEM := ""
		if v.Certificate != nil {
			certPEM = certificateToPEM(v.Certificate)
		}
		privateKeyPEM := ""
		if v.Key != nil {
			privateKeyPEM = privateKeyToPEM(v.Key)
		}*/

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

		out = append(out, obj)
	}

	return out
}
