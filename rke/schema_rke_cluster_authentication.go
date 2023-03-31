package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	rkeClusterAuthenticationStrategyX509 = "x509"
)

var (
	rkeClusterAuthenticationStrategyList = []string{rkeClusterAuthenticationStrategyX509}
)

//Schemas

func rkeClusterAuthenticationWebhookFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"config_file": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Multiline string that represent a custom webhook config file",
		},
		"cache_timeout": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Controls how long to cache authentication decisions",
		},
	}

	return s
}
func rkeClusterAuthenticationFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"sans": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "List of additional hostnames and IPs to include in the api server PKI cert",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"strategy": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      rkeClusterAuthenticationStrategyX509,
			Description:  "Authentication strategy that will be used in RKE k8s cluster",
			ValidateFunc: validation.StringInSlice(rkeClusterAuthenticationStrategyList, false),
		},

		"webhook": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "Webhook configuration options",
			Elem: &schema.Resource{
				Schema: rkeClusterAuthenticationWebhookFields(),
			},
		},
	}
	return s
}
