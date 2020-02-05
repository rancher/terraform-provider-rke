package rke

import (
	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

// Flatteners

func flattenRKEClusterBastionHost(in rancher.BastionHost) []interface{} {
	if len(in.Address) == 0 || len(in.User) == 0 {
		return nil
	}

	obj := make(map[string]interface{})

	obj["address"] = in.Address
	obj["user"] = in.User

	if len(in.Port) > 0 {
		obj["port"] = in.Port
	}

	obj["ssh_agent_auth"] = in.SSHAgentAuth

	if len(in.SSHKey) > 0 {
		obj["ssh_key"] = in.SSHKey
	}

	if len(in.SSHKeyPath) > 0 {
		obj["ssh_key_path"] = in.SSHKeyPath
	}

	return []interface{}{obj}
}

// Expanders

func expandRKEClusterBastionHost(p []interface{}) rancher.BastionHost {
	obj := rancher.BastionHost{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["address"].(string); ok && len(v) > 0 {
		obj.Address = v
	}

	if v, ok := in["port"].(string); ok && len(v) > 0 {
		obj.Port = v
	}

	if v, ok := in["ssh_agent_auth"].(bool); ok {
		obj.SSHAgentAuth = v
	}

	if v, ok := in["ssh_key"].(string); ok && len(v) > 0 {
		obj.SSHKey = v
	}

	if v, ok := in["ssh_key_path"].(string); ok && len(v) > 0 {
		obj.SSHKeyPath = v
	}

	if v, ok := in["user"].(string); ok && len(v) > 0 {
		obj.User = v
	}

	return obj
}
