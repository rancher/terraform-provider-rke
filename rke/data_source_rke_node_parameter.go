package rke

import (
	"encoding/json"

	"github.com/hashicorp/terraform/helper/schema"
	"gopkg.in/yaml.v2"
)

func dataSourceRKENodeParameter() *schema.Resource {
	return &schema.Resource{
		Read:   resourceRKENodeParameterRead,
		Schema: nodeDataSourceSchema(),
	}
}

func resourceRKENodeParameterRead(d *schema.ResourceData, _ interface{}) error {

	nodeValues := map[string]interface{}{
		"node_name":         d.Get("node_name"),
		"address":           d.Get("address"),
		"port":              d.Get("port"),
		"internal_address":  d.Get("internal_address"),
		"role":              d.Get("role"),
		"roles":             d.Get("roles"),
		"hostname_override": d.Get("hostname_override"),
		"user":              d.Get("user"),
		"docker_socket":     d.Get("docker_socket"),
		"ssh_agent_auth":    d.Get("ssh_agent_auth"),
		"ssh_key":           d.Get("ssh_key"),
		"ssh_key_path":      d.Get("ssh_key_path"),
		"labels":            d.Get("labels"),
	}

	node, err := parseResourceRKEConfigNode(nodeValues)
	if err != nil {
		return err
	}

	// to YAML
	strYAML, err := yaml.Marshal(&node)
	if err != nil {
		return err
	}
	d.Set("yaml", string(strYAML))

	// to JSON
	strJSON, err := json.Marshal(&node)
	if err != nil {
		return err
	}
	d.Set("json", string(strJSON))

	d.SetId(d.Get("address").(string))
	return nil
}
