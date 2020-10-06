package rke

import (
	rancher "github.com/rancher/rke/types"
)

// Flatteners

func flattenRKEClusterRestore(in rancher.RestoreConfig) []interface{} {
	obj := make(map[string]interface{})

	obj["restore"] = in.Restore

	if len(in.SnapshotName) > 0 {
		obj["snapshot_name"] = in.SnapshotName
	}

	return []interface{}{obj}
}

// Expanders

func expandRKEClusterRestore(p []interface{}) rancher.RestoreConfig {
	obj := rancher.RestoreConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["restore"].(bool); ok {
		obj.Restore = v
	}

	if v, ok := in["snapshot_name"].(string); ok && len(v) > 0 {
		obj.SnapshotName = v
	}

	return obj
}
