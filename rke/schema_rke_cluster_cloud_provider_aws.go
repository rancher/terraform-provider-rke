package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	rkeClusterCloudProviderAwsName = "aws"
)

//Schemas

func rkeClusterCloudProviderAwsGlobalFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"disable_security_group_ingress": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Disables the automatic ingress creation",
		},
		"disable_strict_zone_check": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Setting this to true will disable the check and provide a warning that the check was skipped",
		},
		"elb_security_group": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Use these ELB security groups instead create new",
		},
		"kubernetes_cluster_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The cluster id we'll use to identify our cluster resources",
		},
		"kubernetes_cluster_tag": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Legacy cluster id we'll use to identify our cluster resources",
		},
		"role_arn": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IAM role to assume when interaction with AWS APIs",
		},
		"route_table_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Enables using a specific RouteTable",
		},
		"subnet_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Enables using a specific subnet to use for ELB's",
		},
		"vpc": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The AWS VPC flag enables the possibility to run the master components on a different aws account, on a different cloud provider or on-premises. If the flag is set also the KubernetesClusterTag must be provided",
		},
		"zone": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The AWS zone",
		},
	}
	return s
}

func rkeClusterCloudProviderAwsServiceOverrideFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"service": {
			Type:     schema.TypeString,
			Required: true,
		},
		"key": {
			Type:       schema.TypeString,
			Optional:   true,
			Deprecated: "Use service instead",
		},
		"region": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"signing_method": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"signing_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"signing_region": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"url": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
	return s
}

func rkeClusterCloudProviderAwsFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"global": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: rkeClusterCloudProviderAwsGlobalFields(),
			},
		},
		"service_override": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: rkeClusterCloudProviderAwsServiceOverrideFields(),
			},
		},
	}
	return s
}
