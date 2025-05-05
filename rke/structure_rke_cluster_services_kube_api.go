package rke

import (
	"encoding/json"
	"fmt"
	rancher "github.com/rancher/rke/types"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	apiserverconfigv1 "k8s.io/apiserver/pkg/apis/apiserver/v1"
	auditv1 "k8s.io/apiserver/pkg/apis/audit/v1"
	eventratelimitapi "k8s.io/kubernetes/plugin/pkg/admission/eventratelimit/apis/eventratelimit"
)

// Flatteners

func flattenRKEClusterServicesKubeAPIAuditLogConfig(in *rancher.AuditLogConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	obj["format"] = in.Format
	obj["max_age"] = in.MaxAge
	obj["max_backup"] = in.MaxBackup
	obj["max_size"] = in.MaxSize
	obj["path"] = in.Path

	if in.Policy != nil {
		// convert Policy to a map to maintain json order
		policyMap, err := interfaceToMap(in.Policy)
		if err != nil {
			return []interface{}{}, fmt.Errorf("interface to map err: %v", err)
		}
		policyStr, err := interfaceToJSON(policyMap)
		if err != nil {
			return []interface{}{}, fmt.Errorf("interface to json err: %v", err)
		}
		obj["policy"] = policyStr
	}

	return []interface{}{obj}, nil
}

func flattenRKEClusterServicesKubeAPIAuditLog(in *rancher.AuditLog) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	obj["enabled"] = in.Enabled

	config, err := flattenRKEClusterServicesKubeAPIAuditLogConfig(in.Configuration)
	if err != nil {
		return []interface{}{}, fmt.Errorf("Flattening RKEClusterServicesKubeAPIAuditLogConfig err: %v", err)
	}
	obj["configuration"] = config

	return []interface{}{obj}, nil
}

func flattenRKEClusterServicesKubeAPIEventRateLimit(in *rancher.EventRateLimit) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	obj["enabled"] = in.Enabled

	if in.Configuration != nil {
		if len(in.Configuration.TypeMeta.Kind) == 0 {
			in.Configuration.TypeMeta.Kind = clusterServicesKubeAPIEventRateLimitConfigKindDefault

		}
		if len(in.Configuration.TypeMeta.APIVersion) == 0 {
			in.Configuration.TypeMeta.APIVersion = clusterServicesKubeAPIEventRateLimitConfigAPIDefault
		}
		configMap, err := interfaceToMap(in.Configuration)
		if err != nil {
			return []interface{}{}, fmt.Errorf("marshalling configuration map: %v", err)
		}
		configStr, err := interfaceToGhodssyaml(configMap)
		if err != nil {
			return []interface{}{}, fmt.Errorf("marshalling configuration yaml: %v", err)
		}
		obj["configuration"] = configStr
	}

	return []interface{}{obj}, nil
}

func flattenRKEClusterServicesKubeAPISecretsEncryptionConfig(in *rancher.SecretsEncryptionConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	obj["enabled"] = in.Enabled

	if in.CustomConfig != nil {
		configStr, err := interfaceToGhodssyaml(in.CustomConfig)
		if err != nil {
			return []interface{}{}, fmt.Errorf("marshalling custom_config yaml: %v", err)
		}
		obj["custom_config"] = configStr
	}

	return []interface{}{obj}, nil
}

func flattenRKEClusterServicesKubeAPI(in rancher.KubeAPIService) ([]interface{}, error) {
	obj := make(map[string]interface{})

	obj["always_pull_images"] = in.AlwaysPullImages

	if in.AuditLog != nil {
		auditLog, err := flattenRKEClusterServicesKubeAPIAuditLog(in.AuditLog)
		if err != nil {
			return []interface{}{}, err
		}
		obj["audit_log"] = auditLog
	}

	if in.EventRateLimit != nil {
		eventRate, err := flattenRKEClusterServicesKubeAPIEventRateLimit(in.EventRateLimit)
		if err != nil {
			return []interface{}{}, err
		}
		obj["event_rate_limit"] = eventRate
	}

	if len(in.ExtraArgs) > 0 {
		obj["extra_args"] = toMapInterface(in.ExtraArgs)
	}

	if len(in.WindowsExtraArgs) > 0 {
		obj["windows_extra_args"] = toMapInterface(in.WindowsExtraArgs)
	}

	if len(in.ExtraArgsArray) > 0 {
		j, err := json.Marshal(in.ExtraArgsArray)
		if err != nil {
			return nil, err
		}
		obj["extra_args_array"] = string(j)
	}

	if len(in.WindowsExtraArgsArray) > 0 {
		j, err := json.Marshal(in.WindowsExtraArgsArray)
		if err != nil {
			return nil, err
		}
		obj["windows_extra_args_array"] = string(j)
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

	if len(in.PodSecurityConfiguration) > 0 {
		obj["pod_security_configuration"] = in.PodSecurityConfiguration
	}

	if in.SecretsEncryptionConfig != nil {
		secretEnc, err := flattenRKEClusterServicesKubeAPISecretsEncryptionConfig(in.SecretsEncryptionConfig)
		if err != nil {
			return []interface{}{}, err
		}
		obj["secrets_encryption_config"] = secretEnc
	}

	if len(in.ServiceClusterIPRange) > 0 {
		obj["service_cluster_ip_range"] = in.ServiceClusterIPRange
	}

	if len(in.ServiceNodePortRange) > 0 {
		obj["service_node_port_range"] = in.ServiceNodePortRange
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandRKEClusterServicesKubeAPIAuditLogConfig(p []interface{}) (*rancher.AuditLogConfig, error) {
	obj := &rancher.AuditLogConfig{}
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil, nil
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

	if v, ok := in["policy"].(string); ok && len(v) > 0 {
		policyBytes := []byte(v)
		scheme := runtime.NewScheme()
		err := auditv1.AddToScheme(scheme)
		if err != nil {
			return nil, fmt.Errorf("error adding to scheme: %v", err)
		}
		codecs := serializer.NewCodecFactory(scheme)
		p := auditv1.Policy{}
		err = runtime.DecodeInto(codecs.UniversalDecoder(), policyBytes, &p)
		if err != nil || p.Kind != "Policy" {
			return nil, fmt.Errorf("error decoding audit policy %s\n: %v", string(policyBytes), err)
		}

		obj.Policy = &p
	}

	return obj, nil
}

func expandRKEClusterServicesKubeAPIAuditLog(p []interface{}) (*rancher.AuditLog, error) {
	obj := &rancher.AuditLog{}
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["enabled"].(bool); ok {
		obj.Enabled = v
	}

	if v, ok := in["configuration"].([]interface{}); ok && len(v) > 0 {
		config, err := expandRKEClusterServicesKubeAPIAuditLogConfig(v)
		if err != nil {
			return nil, err
		}
		obj.Configuration = config
	}

	return obj, nil
}

func expandRKEClusterServicesKubeAPIEventRateLimit(p []interface{}) (*rancher.EventRateLimit, error) {
	obj := &rancher.EventRateLimit{}
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["enabled"].(bool); ok {
		obj.Enabled = v
	}

	if v, ok := in["configuration"].(string); ok && len(v) > 0 {
		configMap, err := ghodssyamlToMapInterface(v)
		if err != nil {
			return obj, fmt.Errorf("unmarshalling configuration yaml: %v", err)
		}
		configStr, err := mapInterfaceToJSON(configMap)
		if err != nil {
			return obj, fmt.Errorf("marshalling custom_config json: %v", err)
		}
		newConfig := &eventratelimitapi.Configuration{}
		err = jsonToInterface(configStr, newConfig)
		if err != nil {
			return obj, fmt.Errorf("unmarshalling EncryptionConfiguration json:\n%s", v)
		}
		obj.Configuration = newConfig
	}

	return obj, nil
}

func expandRKEClusterServicesKubeAPISecretsEncryptionConfig(p []interface{}) (*rancher.SecretsEncryptionConfig, error) {
	obj := &rancher.SecretsEncryptionConfig{}
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["enabled"].(bool); ok {
		obj.Enabled = v
	}

	if v, ok := in["custom_config"].(string); ok && len(v) > 0 {
		configMap, err := ghodssyamlToMapInterface(v)
		if err != nil {
			return obj, fmt.Errorf("unmarshalling custom_config yaml: %v", err)
		}
		configStr, err := mapInterfaceToJSON(configMap)
		if err != nil {
			return obj, fmt.Errorf("marshalling custom_config json: %v", err)
		}
		newConfigV1 := &apiserverconfigv1.EncryptionConfiguration{}
		err = jsonToInterface(configStr, newConfigV1)
		if err != nil {
			return obj, fmt.Errorf("unmarshalling EncryptionConfiguration json: %v", err)
		}
		obj.CustomConfig = newConfigV1
	}

	return obj, nil
}

func expandRKEClusterServicesKubeAPI(p []interface{}) (rancher.KubeAPIService, error) {
	obj := rancher.KubeAPIService{}
	if p == nil || len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["always_pull_images"].(bool); ok {
		obj.AlwaysPullImages = v
	}

	if v, ok := in["audit_log"].([]interface{}); ok && len(v) > 0 {
		auditLog, err := expandRKEClusterServicesKubeAPIAuditLog(v)
		if err != nil {
			return obj, err
		}
		obj.AuditLog = auditLog
	}

	if v, ok := in["event_rate_limit"].([]interface{}); ok && len(v) > 0 {
		eventRate, err := expandRKEClusterServicesKubeAPIEventRateLimit(v)
		if err != nil {
			return obj, err
		}
		obj.EventRateLimit = eventRate
	}

	if v, ok := in["extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ExtraArgs = toMapString(v)
	}

	if v, ok := in["windows_extra_args"].(map[string]interface{}); ok && len(v) > 0 {
		obj.WindowsExtraArgs = toMapString(v)
	}

	if v, ok := in["extra_args_array"].(string); ok && len(v) != 0 {
		array, err := jsonToMapStringSlice(v)
		if err != nil {
			return rancher.KubeAPIService{}, err
		}
		obj.ExtraArgsArray = array
	}

	if v, ok := in["windows_extra_args_array"].(string); ok && len(v) != 0 {
		array, err := jsonToMapStringSlice(v)
		if err != nil {
			return rancher.KubeAPIService{}, err
		}
		obj.WindowsExtraArgsArray = array
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

	if v, ok := in["pod_security_configuration"].(string); ok && len(v) > 0 {
		obj.PodSecurityConfiguration = v
	}

	if v, ok := in["secrets_encryption_config"].([]interface{}); ok && len(v) > 0 {
		secretEnc, err := expandRKEClusterServicesKubeAPISecretsEncryptionConfig(v)
		if err != nil {
			return obj, err
		}
		obj.SecretsEncryptionConfig = secretEnc
	}

	if v, ok := in["service_cluster_ip_range"].(string); ok && len(v) > 0 {
		obj.ServiceClusterIPRange = v
	}

	if v, ok := in["service_node_port_range"].(string); ok && len(v) > 0 {
		obj.ServiceNodePortRange = v
	}

	return obj, nil
}
