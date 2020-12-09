package rke

import (
	"github.com/rancher/rke/hosts"
	rancher "github.com/rancher/rke/types"
)

// Flatteners

func flattenRKEClusterNodeDrainInput(in *rancher.NodeDrainInput) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	obj["delete_local_data"] = in.DeleteLocalData
	obj["force"] = in.Force
	obj["grace_period"] = in.GracePeriod
	if in.IgnoreDaemonSets != nil {
		obj["ignore_daemon_sets"] = *in.IgnoreDaemonSets
	}
	obj["timeout"] = in.Timeout

	return []interface{}{obj}
}

func flattenRKEClusterNodeUpgradeStrategy(in *rancher.NodeUpgradeStrategy) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if in.Drain != nil {
		obj["drain"] = *in.Drain
	}

	if in.DrainInput != nil {
		obj["drain_input"] = flattenRKEClusterNodeDrainInput(in.DrainInput)
	}

	if len(in.MaxUnavailableControlplane) > 0 {
		obj["max_unavailable_controlplane"] = in.MaxUnavailableControlplane
	}

	if len(in.MaxUnavailableWorker) > 0 {
		obj["max_unavailable_worker"] = in.MaxUnavailableWorker
	}

	return []interface{}{obj}
}

func flattenRKEClusterNodes(input []rancher.RKEConfigNode, p []interface{}) []interface{} {
	if input == nil || len(input) == 0 {
		return []interface{}{}
	}

	// Sorting input array by data interface
	pIndexAddress := map[string]int{}
	for i := range p {
		if row, ok := p[i].(map[string]interface{}); ok {
			if v, ok := row["address"].(string); ok {
				pIndexAddress[v] = i
			}
		}
	}
	pLen := len(p)
	inputLen := len(input)
	sortedInput := make([]rancher.RKEConfigNode, inputLen)
	newNodes := []rancher.RKEConfigNode{}
	lastIndex := 0
	for i := range sortedInput {
		if v, ok := pIndexAddress[input[i].Address]; ok {
			if v > i && pLen > inputLen {
				v = v - (v - i)
			}
			sortedInput[v] = input[i]
			lastIndex++
			continue
		}
		newNodes = append(newNodes, input[i])
	}

	for i := range newNodes {
		sortedInput[lastIndex+i] = newNodes[i]
	}

	out := make([]interface{}, len(sortedInput))
	for i, in := range sortedInput {
		var obj map[string]interface{}
		if v, ok := pIndexAddress[in.Address]; ok {
			if row, ok := p[v].(map[string]interface{}); ok {
				obj = row
			}
		}
		if obj == nil {
			obj = make(map[string]interface{})
		}

		obj["address"] = in.Address
		obj["role"] = toArrayInterface(in.Role)
		obj["user"] = in.User

		if len(in.DockerSocket) > 0 {
			obj["docker_socket"] = in.DockerSocket
		}

		if v, ok := obj["hostname_override"].(string); ok && len(v) > 0 && len(in.HostnameOverride) > 0 {
			obj["hostname_override"] = in.HostnameOverride
		}
		if v, ok := obj["internal_address"].(string); ok && len(v) > 0 && len(in.InternalAddress) > 0 {
			obj["internal_address"] = in.InternalAddress
		}

		if len(in.Labels) > 0 {
			obj["labels"] = toMapInterface(in.Labels)
		}

		if len(in.NodeName) > 0 {
			obj["node_name"] = in.NodeName
		}

		if v, ok := obj["port"].(string); ok && len(v) > 0 && len(in.Port) > 0 {
			obj["port"] = in.Port
		}

		obj["ssh_agent_auth"] = in.SSHAgentAuth

		if v, ok := obj["ssh_cert"].(string); ok && len(v) > 0 && len(in.SSHCert) > 0 {
			obj["ssh_cert"] = in.SSHCert
		}

		if v, ok := obj["ssh_cert_path"].(string); ok && len(v) > 0 && len(in.SSHCertPath) > 0 {
			obj["ssh_cert_path"] = in.SSHCertPath
		}

		if v, ok := obj["ssh_key"].(string); ok && len(v) > 0 && len(in.SSHKey) > 0 {
			obj["ssh_key"] = in.SSHKey
		}
		if v, ok := obj["ssh_key_path"].(string); ok && len(v) > 0 && len(in.SSHKeyPath) > 0 {
			obj["ssh_key_path"] = in.SSHKeyPath
		}

		if in.Taints != nil {
			obj["taints"] = flattenRKEClusterTaints(in.Taints)
		}

		out[i] = obj
	}

	return out
}

func flattenRKEClusterNodesComputed(p []*hosts.Host) []interface{} {
	out := []interface{}{}

	for _, in := range p {
		obj := make(map[string]interface{})

		if len(in.Address) > 0 {
			obj["address"] = in.Address
		}

		if len(in.NodeName) > 0 {
			obj["node_name"] = in.NodeName
		}

		out = append(out, obj)
	}

	return out
}

// Expanders

func expandRKEClusterNodeDrainInput(p []interface{}) *rancher.NodeDrainInput {
	obj := &rancher.NodeDrainInput{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["delete_local_data"].(bool); ok {
		obj.DeleteLocalData = v
	}

	if v, ok := in["force"].(bool); ok {
		obj.Force = v
	}

	if v, ok := in["grace_period"].(int); ok {
		obj.GracePeriod = v
	}

	if v, ok := in["ignore_daemon_sets"].(bool); ok {
		obj.IgnoreDaemonSets = &v
	}

	if v, ok := in["timeout"].(int); ok {
		obj.Timeout = v
	}

	return obj
}

func expandRKEClusterNodeUpgradeStrategy(p []interface{}) *rancher.NodeUpgradeStrategy {
	obj := &rancher.NodeUpgradeStrategy{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["drain"].(bool); ok {
		obj.Drain = &v
	}

	if v, ok := in["drain_input"].([]interface{}); ok {
		obj.DrainInput = expandRKEClusterNodeDrainInput(v)
	}

	if v, ok := in["max_unavailable_controlplane"].(string); ok && len(v) > 0 {
		obj.MaxUnavailableControlplane = v
	}

	if v, ok := in["max_unavailable_worker"].(string); ok && len(v) > 0 {
		obj.MaxUnavailableWorker = v
	}

	return obj
}

func expandRKEClusterNodes(p []interface{}) []rancher.RKEConfigNode {
	out := []rancher.RKEConfigNode{}
	if len(p) == 0 || p[0] == nil {
		return out
	}

	for i := range p {
		in := p[i].(map[string]interface{})
		obj := rancher.RKEConfigNode{}

		if v, ok := in["address"].(string); ok && len(v) > 0 {
			obj.Address = v
		}

		if v, ok := in["role"].([]interface{}); ok && len(v) > 0 {
			obj.Role = toArrayString(v)
		}

		if v, ok := in["docker_socket"].(string); ok && len(v) > 0 {
			obj.DockerSocket = v
		}

		if v, ok := in["hostname_override"].(string); ok && len(v) > 0 {
			obj.HostnameOverride = v
		}

		if v, ok := in["internal_address"].(string); ok && len(v) > 0 {
			obj.InternalAddress = v
		}

		if v, ok := in["labels"].(map[string]interface{}); ok && len(v) > 0 {
			obj.Labels = toMapString(v)
		}

		if v, ok := in["node_name"].(string); ok && len(v) > 0 {
			obj.NodeName = v
		}

		if v, ok := in["port"].(string); ok && len(v) > 0 {
			obj.Port = v
		}

		if v, ok := in["ssh_agent_auth"].(bool); ok {
			obj.SSHAgentAuth = v
		}

		if v, ok := in["ssh_cert"].(string); ok && len(v) > 0 {
			obj.SSHCert = v
		}

		if v, ok := in["ssh_cert_path"].(string); ok && len(v) > 0 {
			obj.SSHCertPath = v
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

		if v, ok := in["taints"].([]interface{}); ok && len(v) > 0 {
			obj.Taints = expandRKEClusterTaints(v)
		}

		out = append(out, obj)
	}

	return out
}
