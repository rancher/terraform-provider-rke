package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"k8s.io/api/core/v1"
)

const (
	rkeClusterIngressNginx = "nginx"
	rkeClusterIngressNone  = "none"
)

var (
	rkeClusterIngressDNSPolicyClusterFirst            = string(v1.DNSClusterFirst)
	rkeClusterIngressDNSPolicyClusterFirstWithHostNet = string(v1.DNSClusterFirstWithHostNet)
	rkeClusterIngressDNSPolicyDefault                 = string(v1.DNSDefault)
	rkeClusterIngressDNSPolicyNone                    = string(v1.DNSNone)
	rkeClusterIngressProviderList                     = []string{rkeClusterIngressNginx, rkeClusterIngressNone}
	rkeClusterIngressDNSPolicyList                    = []string{
		rkeClusterIngressDNSPolicyClusterFirst,
		rkeClusterIngressDNSPolicyClusterFirstWithHostNet,
		rkeClusterIngressDNSPolicyDefault,
		rkeClusterIngressDNSPolicyNone,
	}
)

//Schemas

func rkeClusterIngressFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"dns_policy": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(rkeClusterIngressDNSPolicyList, true),
			Description:  "Ingress controller dns policy",
		},
		"extra_args": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Extra arguments for the ingress controller",
		},
		"node_selector": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Node selector key pair",
		},
		"options": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Ingress controller options",
		},
		"provider": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      rkeClusterIngressNginx,
			ValidateFunc: validation.StringInSlice(rkeClusterIngressProviderList, true),
			Description:  "Ingress controller provider",
		},
		"default_backend": {
		Type: schema.TypeBool,
		Optional: true,
		Default: true,
		Description:"Ingress Default Backend"
		}
	}
	return s
}
