package rke

import (
	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
)

// Flatteners

func flattenRKEClusterServicesKubeAPIAuditLogConfig(in *rancher.AuditLogConfig) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	obj["format"] = in.Format
	obj["max_age"] = in.MaxAge
	obj["max_backup"] = in.MaxBackup
	obj["max_size"] = in.MaxSize
	obj["path"] = in.Path

	return []interface{}{obj}
}

func flattenRKEClusterServicesKubeAPIAuditLog(in *rancher.AuditLog) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	obj["enabled"] = in.Enabled
	obj["configuration"] = flattenRKEClusterServicesKubeAPIAuditLogConfig(in.Configuration)

	return []interface{}{obj}
}

func flattenRKEClusterServicesKubeAPIEventRateLimit(in *rancher.EventRateLimit) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	obj["enabled"] = in.Enabled

	return []interface{}{obj}
}

func flattenRKEClusterServicesKubeAPISecretsEncryptionConfig(in *rancher.SecretsEncryptionConfig) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	obj["enabled"] = in.Enabled

	return []interface{}{obj}
}

func flattenRKEClusterServicesKubeAPI(in rancher.KubeAPIService) []interface{} {
	obj := make(map[string]interface{})

	obj["always_pull_images"] = in.AlwaysPullImages

	if in.AuditLog != nil {
		obj["audit_log"] = flattenRKEClusterServicesKubeAPIAuditLog(in.AuditLog)
	}

	if in.EventRateLimit != nil {
		obj["event_rate_limit"] = flattenRKEClusterServicesKubeAPIEventRateLimit(in.EventRateLimit)
	}

	if len(in.ExtraArgs) > 0 {
		obj["extra_args"] = toMapInterface(in.ExtraArgs)
	}

	if len(in.ExtraBinds) > 0 {
		obj["extra_binds"] = toArrayInterface(in.ExtraBinds)
	}

	if len(in.ExtraEnv) > 0 {
		obj["extra_env"] = toArrayInterface(in.ExtraEnv)
	}

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	obj["pod_security_policy"] = in.PodSecurityPolicy

	if in.SecretsEncryptionConfig != nil {
		obj["secrets_encryption_config"] = flattenRKEClusterServicesKubeAPISecretsEncryptionConfig(in.SecretsEncryptionConfig)
	}

	if len(in.ServiceClusterIPRange) > 0 {
		obj["service_cluster_ip_range"] = in.ServiceClusterIPRange
	}

	if len(in.ServiceNodePortRange) > 0 {
		obj["service_node_port_range"] = in.ServiceNodePortRange
	}

	return []interface{}{obj}
}

// Expanders

func expandRKEClusterServicesKubeAPIAuditLogConfig(p []interface{}) *rancher.AuditLogConfig {
	obj := &rancher.AuditLogConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["format"].(string); ok && len(v) > 0 {
		obj.Format = v
	}

	if v, ok := in["max_age"].(int); ok && v > 0 {
		obj.MaxAge = v
	}

	if v, ok := in["max_backup"].(int); ok && v > 0 {
		obj.MaxBackup = v
	}

	if v, ok := in["max_size"].(int); ok && v > 0 {
		obj.MaxSize = v
	}

	if v, ok := in["path"].(string); ok && len(v) > 0 {
		obj.Path = v
	}

	return obj
}

func expandRKEClusterServicesKubeAPIAuditLog(p []interface{}) *rancher.AuditLog {
	obj := &rancher.AuditLog{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["enabled"].(bool); ok {
		obj.Enabled = v
	}

	if v, ok := in["configuration"].([]interface{}); ok && len(v) > 0 {
		obj.Configuration = expandRKEClusterServicesKubeAPIAuditLogConfig(v)
	}

	return obj
}

func expandRKEClusterServicesKubeAPIEventRateLimit(p []interface{}) *rancher.EventRateLimit {
	obj := &rancher.EventRateLimit{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["enabled"].(bool); ok {
		obj.Enabled = v
	}

	return obj
}

func expandRKEClusterServicesKubeAPISecretsEncryptionConfig(p []interface{}) *rancher.SecretsEncryptionConfig {
	obj := &rancher.SecretsEncryptionConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["enabled"].(bool); ok {
		obj.Enabled = v
	}

	return obj
}

func expandRKEClusterServicesKubeAPI(p []interface{}) rancher.KubeAPIService {
	obj := rancher.KubeAPIService{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["always_pull_images"].(bool); ok {
		obj.AlwaysPullImages = v
	}

	if v, ok := in["audit_log"].([]interface{}); ok && len(v) > 0 {
		obj.AuditLog = expandRKEClusterServicesKubeAPIAuditLog(v)
	}

	if v, ok := in["event_rate_limit"].([]interface{}); ok && len(v) > 0 {
		obj.EventRateLimit = expandRKEClusterServicesKubeAPIEventRateLimit(v)
	}

	if v, ok := in["extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ExtraArgs = toMapString(v)
	}

	if v, ok := in["extra_binds"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraBinds = toArrayString(v)
	}

	if v, ok := in["extra_env"].([]interface{}); ok && len(v) > 0 {
		obj.ExtraEnv = toArrayString(v)
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["pod_security_policy"].(bool); ok {
		obj.PodSecurityPolicy = v
	}

	if v, ok := in["secrets_encryption_config"].([]interface{}); ok && len(v) > 0 {
		obj.SecretsEncryptionConfig = expandRKEClusterServicesKubeAPISecretsEncryptionConfig(v)
	}

	if v, ok := in["service_cluster_ip_range"].(string); ok && len(v) > 0 {
		obj.ServiceClusterIPRange = v
	}

	if v, ok := in["service_node_port_range"].(string); ok && len(v) > 0 {
		obj.ServiceNodePortRange = v
	}

	return obj
}
