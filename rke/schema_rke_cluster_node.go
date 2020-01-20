package rke

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var (
	rkeClusterNodesRoles = []string{"controlplane", "etcd", "worker"}
)

//Schemas

func rkeClusterNodeFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"address": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "IP or FQDN that is fully resolvable and used for SSH communication",
		},
		"role": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "Node roles in k8s cluster [controlplane/worker/etcd])",
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(rkeClusterNodesRoles, true),
			},
		},
		"roles": {
			Type:         schema.TypeString,
			Optional:     true,
			Deprecated:   "Use role instead",
			Description:  "Node role in kubernetes cluster [controlplane/worker/etcd], specified by a comma-separated string",
			ValidateFunc: validation.StringInSlice(rkeClusterNodesRoles, true),
		},
		"user": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "SSH user that will be used by RKE",
		},
		"docker_socket": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Docker socket on the node that will be used in tunneling",
		},
		"hostname_override": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Hostname override",
		},
		"internal_address": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Internal address that will be used for components communication",
		},
		"labels": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Node Labels",
		},
		"node_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the host provisioned via docker machine",
		},
		"port": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "22",
			ValidateFunc: validation.IntBetween(1, 65535),
			Description:  "Port used for SSH communication",
		},
		"ssh_agent_auth": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "SSH Agent Auth enable",
		},
		"ssh_cert": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "SSH Certificate",
		},
		"ssh_cert_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "SSH Certificate path",
		},
		"ssh_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Sensitive:   true,
			Description: "SSH Private Key",
		},
		"ssh_key_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "SSH Private Key path",
		},
		"taints": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Node taints",
			Elem: &schema.Resource{
				Schema: rkeClusterTaintFields(),
			},
		},
	}
	return s
}

func rkeClusterNodeComputedFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"node_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"address": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
	return s
}
