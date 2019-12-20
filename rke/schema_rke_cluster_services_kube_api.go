package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func rkeClusterServicesKubeAPIAuditLogConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"format": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "json",
		},
		"max_age": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  30,
		},
		"max_backup": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  10,
		},
		"max_size": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  100,
		},
		"path": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "/var/log/kube-audit/audit-log.json",
		},
	}
	return s
}

func rkeClusterServicesKubeAPIAuditLogFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"configuration": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: rkeClusterServicesKubeAPIAuditLogConfigFields(),
			},
		},
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
	return s
}

func rkeClusterServicesKubeAPIEventRateLimitFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
	return s
}

func rkeClusterServicesKubeAPISecretsEncryptionConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
	return s
}

func rkeClusterServicesKubeAPIFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"always_pull_images": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable/Disable AlwaysPullImages admissions plugin",
		},
		"audit_log": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: rkeClusterServicesKubeAPIAuditLogFields(),
			},
		},
		"event_rate_limit": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: rkeClusterServicesKubeAPIEventRateLimitFields(),
			},
		},
		"extra_args": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Extra arguments that are added to the kube-api services",
		},
		"extra_binds": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Extra binds added to the controlplane nodes",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"extra_env": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Extra env added to the controlplane nodes",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"image": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"pod_security_policy": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enabled/Disable PodSecurityPolicy",
		},
		"secrets_encryption_config": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: rkeClusterServicesKubeAPISecretsEncryptionConfigFields(),
			},
		},
		"service_cluster_ip_range": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Virtual IP range that will be used by Kubernetes services",
		},
		"service_node_port_range": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Port range for services defined with NodePort type",
		},
	}
	return s
}
