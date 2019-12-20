package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	rkeClusterDNSProviderKube = "kube-dns"
	rkeClusterDNSProviderCore = "coredns"
	rkeClusterDNSProviderNone = "none"
)

var (
	rkeClusterDNSProviderList = []string{
		rkeClusterDNSProviderKube,
		rkeClusterDNSProviderCore,
		rkeClusterDNSProviderNone,
	}
)

//Schemas

func rkeClusterDNSFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"node_selector": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "NodeSelector key pair",
		},
		"provider": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      rkeClusterDNSProviderCore,
			Description:  "DNS provider",
			ValidateFunc: validation.StringInSlice(rkeClusterDNSProviderList, true),
		},
		"reverse_cidrs": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "ReverseCIDRs",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"upstream_nameservers": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Upstream nameservers",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}

	return s
}
