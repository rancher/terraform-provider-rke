package rke

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func nodeSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"node_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the host provisioned via docker machine",
		},
		"address": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "IP or FQDN that is fully resolvable and used for SSH communication",
		},
		"port": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
			Description:  "Port used for SSH communication",
		},
		"internal_address": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Internal address that will be used for components communication",
		},
		"role": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			// cannot use ConflictsWith in this context. see https://github.com/terraform-providers/terraform-provider-google/pull/1062
			// ConflictsWith: []string{"roles"},
			Description: "Node role in kubernetes cluster [controlplane/worker/etcd])",
		},
		"roles": {
			Type:     schema.TypeString,
			Optional: true,
			// cannot use ConflictsWith in this context. see https://github.com/terraform-providers/terraform-provider-google/pull/1062
			// ConflictsWith: []string{"role"},
			Deprecated:  "roles is a workaround when a role can not be specified in list",
			Description: "Node role in kubernetes cluster [controlplane/worker/etcd], specified by a comma-separated string",
		},
		"hostname_override": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "HostnameOverride",
		},
		"user": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "SSH user that will be used by RKE",
		},
		"docker_socket": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Docker socket on the node that will be used in tunneling",
		},
		"ssh_agent_auth": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "SSH Agent Auth enable",
		},
		"ssh_key": {
			Type:        schema.TypeString,
			Sensitive:   true,
			Optional:    true,
			Description: "SSH Private Key",
		},
		"ssh_key_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "SSH Private Key path",
		},
		"ssh_cert": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "SSH Certificate",
		},
		"ssh_cert_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "SSH Certificate path",
		},
		"labels": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Node Labels",
		},
	}
}

func nodeDataSourceSchema() map[string]*schema.Schema {
	nodeSchema := nodeSchema()

	nodeSchema["yaml"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "RKE Node YAML",
	}
	nodeSchema["json"] = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "RKE Node JSON",
	}

	return nodeSchema
}
