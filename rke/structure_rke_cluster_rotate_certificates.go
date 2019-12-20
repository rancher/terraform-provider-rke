package rke

import (
	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

// Flatteners

func flattenRKEClusterRotateCertificates(in *rancher.RotateCertificates) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	obj["ca_certificates"] = in.CACertificates

	if len(in.Services) > 0 {
		obj["services"] = toArrayInterface(in.Services)
	}

	return []interface{}{obj}
}

// Expanders

func expandRKEClusterRotateCertificates(p []interface{}) *rancher.RotateCertificates {
	obj := &rancher.RotateCertificates{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["ca_certificates"].(bool); ok {
		obj.CACertificates = v
	}

	if v, ok := in["services"].([]interface{}); ok && len(v) > 0 {
		obj.Services = toArrayString(v)
	}

	return obj
}
