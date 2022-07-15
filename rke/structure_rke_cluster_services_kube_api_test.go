package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/rke/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	auditv1 "k8s.io/apiserver/pkg/apis/audit/v1"
	apiserverconfigv1 "k8s.io/apiserver/pkg/apis/config/v1"
)

var (
	testRKEClusterServicesKubeAPIAuditLogConfigConf               *rancher.AuditLogConfig
	testRKEClusterServicesKubeAPIAuditLogConfigInterface          []interface{}
	testRKEClusterServicesKubeAPIAuditLogConf                     *rancher.AuditLog
	testRKEClusterServicesKubeAPIAuditLogInterface                []interface{}
	testRKEClusterServicesKubeAPIEventRateLimitConf               *rancher.EventRateLimit
	testRKEClusterServicesKubeAPIEventRateLimitInterface          []interface{}
	testRKEClusterServicesKubeAPISecretsEncryptionConfigConf      *rancher.SecretsEncryptionConfig
	testRKEClusterServicesKubeAPISecretsEncryptionConfigInterface []interface{}
	testRKEClusterServicesKubeAPIConf                             rancher.KubeAPIService
	testRKEClusterServicesKubeAPIInterface                        []interface{}
)

func init() {
	testRKEClusterServicesKubeAPIAuditLogConfigConf = &rancher.AuditLogConfig{
		Format:    "format",
		MaxAge:    5,
		MaxBackup: 10,
		MaxSize:   100,
		Path:      "path",
		Policy: &auditv1.Policy{
			Rules: []auditv1.PolicyRule{
				{
					Level: "RequestResponse",
					Resources: []auditv1.GroupResources{
						{
							Group:     "*",
							Resources: []string{"pods"},
						},
					},
				},
			},
		},
	}
	testRKEClusterServicesKubeAPIAuditLogConfigConf.Policy.TypeMeta = metav1.TypeMeta{
		Kind:       clusterServicesKubeAPIAuditLogConfigPolicyKindDefault,
		APIVersion: clusterServicesKubeAPIAuditLogConfigPolicyAPIDefault,
	}
	testRKEClusterServicesKubeAPIAuditLogConfigInterface = []interface{}{
		map[string]interface{}{
			"format":     "format",
			"max_age":    5,
			"max_backup": 10,
			"max_size":   100,
			"path":       "path",
			"policy":     `{"apiVersion":"` + clusterServicesKubeAPIAuditLogConfigPolicyAPIDefault + `","kind":"` + clusterServicesKubeAPIAuditLogConfigPolicyKindDefault + `","metadata":{"creationTimestamp":null},"rules":[{"level":"RequestResponse","resources":[{"group":"*","resources":["pods"]}]}]}`,
		},
	}
	testRKEClusterServicesKubeAPIAuditLogConf = &rancher.AuditLog{
		Enabled:       true,
		Configuration: testRKEClusterServicesKubeAPIAuditLogConfigConf,
	}
	testRKEClusterServicesKubeAPIAuditLogInterface = []interface{}{
		map[string]interface{}{
			"enabled":       true,
			"configuration": testRKEClusterServicesKubeAPIAuditLogConfigInterface,
		},
	}
	testRKEClusterServicesKubeAPIEventRateLimitConf = &rancher.EventRateLimit{
		Enabled: true,
		Configuration: &rancher.Configuration{
			Limits: []rancher.Limit{
				{
					Type:  "Server",
					Burst: 30000,
					QPS:   6000,
				},
			},
		},
	}
	testRKEClusterServicesKubeAPIEventRateLimitConf.Configuration.TypeMeta = metav1.TypeMeta{
		Kind:       clusterServicesKubeAPIEventRateLimitConfigKindDefault,
		APIVersion: clusterServicesKubeAPIEventRateLimitConfigAPIDefault,
	}
	testRKEClusterServicesKubeAPIEventRateLimitInterface = []interface{}{
		map[string]interface{}{
			"enabled":       true,
			"configuration": "apiVersion: " + clusterServicesKubeAPIEventRateLimitConfigAPIDefault + "\nkind: " + clusterServicesKubeAPIEventRateLimitConfigKindDefault + "\nlimits:\n- type: Server\n  burst: 30000\n  qps: 6000\n",
		},
	}
	testRKEClusterServicesKubeAPISecretsEncryptionConfigConf = &rancher.SecretsEncryptionConfig{
		Enabled: true,
		CustomConfig: &apiserverconfigv1.EncryptionConfiguration{
			Resources: []apiserverconfigv1.ResourceConfiguration{
				{
					Resources: []string{"secrets"},
					Providers: []apiserverconfigv1.ProviderConfiguration{
						{
							AESCBC: &apiserverconfigv1.AESConfiguration{
								Keys: []apiserverconfigv1.Key{
									{
										Name:   "k-fw5hn",
										Secret: "RTczRjFDODMwQzAyMDVBREU4NDJBMUZFNDhCNzM5N0I=",
									},
								},
							},
							Identity: &apiserverconfigv1.IdentityConfiguration{},
						},
					},
				},
			},
		},
	}
	testRKEClusterServicesKubeAPISecretsEncryptionConfigConf.CustomConfig.TypeMeta = metav1.TypeMeta{
		Kind:       clusterServicesKubeAPISecretsEncryptionConfigKindDefault,
		APIVersion: clusterServicesKubeAPISecretsEncryptionConfigAPIDefault,
	}
	testRKEClusterServicesKubeAPISecretsEncryptionConfigInterface = []interface{}{
		map[string]interface{}{
			"enabled":       true,
			"custom_config": "apiVersion: " + clusterServicesKubeAPISecretsEncryptionConfigAPIDefault + "\nkind: " + clusterServicesKubeAPISecretsEncryptionConfigKindDefault + "\nresources:\n- resources:\n  - secrets\n  providers:\n  - aescbc:\n      keys:\n      - name: k-fw5hn\n        secret: RTczRjFDODMwQzAyMDVBREU4NDJBMUZFNDhCNzM5N0I=\n    identity: {}\n",
		},
	}
	testRKEClusterServicesKubeAPIConf = rancher.KubeAPIService{
		AlwaysPullImages:        true,
		AuditLog:                testRKEClusterServicesKubeAPIAuditLogConf,
		EventRateLimit:          testRKEClusterServicesKubeAPIEventRateLimitConf,
		PodSecurityPolicy:       true,
		SecretsEncryptionConfig: testRKEClusterServicesKubeAPISecretsEncryptionConfigConf,
		ServiceClusterIPRange:   "10.43.0.0/16",
		ServiceNodePortRange:    "30000-32000",
	}
	testRKEClusterServicesKubeAPIConf.ExtraArgs = map[string]string{
		"arg_one": "one",
		"arg_two": "two",
	}
	testRKEClusterServicesKubeAPIConf.WindowsExtraArgs = map[string]string{
		"arg_one": "one",
		"arg_two": "two",
	}
	testRKEClusterServicesKubeAPIConf.ExtraArgsArray = map[string][]string{
		"arg1": {"v1"},
		"arg2": {"v2"},
	}
	testRKEClusterServicesKubeAPIConf.WindowsExtraArgsArray = map[string][]string{
		"arg1": {"v1"},
		"arg2": {"v2"},
	}
	testRKEClusterServicesKubeAPIConf.ExtraBinds = []string{"bind_one", "bind_two"}
	testRKEClusterServicesKubeAPIConf.ExtraEnv = []string{"env_one", "env_two"}
	testRKEClusterServicesKubeAPIConf.Image = "image"
	testRKEClusterServicesKubeAPIInterface = []interface{}{
		map[string]interface{}{
			"always_pull_images": true,
			"audit_log":          testRKEClusterServicesKubeAPIAuditLogInterface,
			"event_rate_limit":   testRKEClusterServicesKubeAPIEventRateLimitInterface,
			"extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"win_extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"extra_args_array": []interface{}{
				map[string]interface{}{
					"extra_arg": []interface{}{
						map[string]interface{}{
							"argument": "arg1",
							"values":   []interface{}{"v1"},
						},
						map[string]interface{}{
							"argument": "arg2",
							"values":   []interface{}{"v2"},
						},
					},
				},
			},
			"win_extra_args_array": []interface{}{
				map[string]interface{}{
					"extra_arg": []interface{}{
						map[string]interface{}{
							"argument": "arg1",
							"values":   []interface{}{"v1"},
						},
						map[string]interface{}{
							"argument": "arg2",
							"values":   []interface{}{"v2"},
						},
					},
				},
			},
			"extra_binds":               []interface{}{"bind_one", "bind_two"},
			"extra_env":                 []interface{}{"env_one", "env_two"},
			"image":                     "image",
			"pod_security_policy":       true,
			"secrets_encryption_config": testRKEClusterServicesKubeAPISecretsEncryptionConfigInterface,
			"service_cluster_ip_range":  "10.43.0.0/16",
			"service_node_port_range":   "30000-32000",
		},
	}
}

func TestFlattenRKEClusterServicesKubeAPIAuditLogConfig(t *testing.T) {

	cases := []struct {
		Input          *rancher.AuditLogConfig
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterServicesKubeAPIAuditLogConfigConf,
			testRKEClusterServicesKubeAPIAuditLogConfigInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenRKEClusterServicesKubeAPIAuditLogConfig(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterServicesKubeAPIAuditLog(t *testing.T) {

	cases := []struct {
		Input          *rancher.AuditLog
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterServicesKubeAPIAuditLogConf,
			testRKEClusterServicesKubeAPIAuditLogInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenRKEClusterServicesKubeAPIAuditLog(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterServicesKubeAPIEventRateLimit(t *testing.T) {

	cases := []struct {
		Input          *rancher.EventRateLimit
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterServicesKubeAPIEventRateLimitConf,
			testRKEClusterServicesKubeAPIEventRateLimitInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenRKEClusterServicesKubeAPIEventRateLimit(tc.Input)
		if err != nil {
			t.Fatalf("Error on flattenRKEClusterServicesKubeAPIEventRateLimit: %#v", err)
		}
		outputObject := &rancher.Configuration{}
		expectedObject := &rancher.Configuration{}
		outputStr, _ := mapInterfaceToJSON(output[0].(map[string]interface{}))
		expectedStr, _ := mapInterfaceToJSON(tc.ExpectedOutput[0].(map[string]interface{}))
		jsonToInterface(outputStr, outputObject)
		jsonToInterface(expectedStr, expectedObject)
		if !reflect.DeepEqual(outputObject, expectedObject) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterServicesKubeAPISecretsEncryptionConfig(t *testing.T) {

	cases := []struct {
		Input          *rancher.SecretsEncryptionConfig
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterServicesKubeAPISecretsEncryptionConfigConf,
			testRKEClusterServicesKubeAPISecretsEncryptionConfigInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenRKEClusterServicesKubeAPISecretsEncryptionConfig(tc.Input)
		if err != nil {
			t.Fatalf("Error on flattenRKEClusterServicesKubeAPISecretsEncryptionConfig: %#v", err)
		}
		outputObject := &apiserverconfigv1.EncryptionConfiguration{}
		expectedObject := &apiserverconfigv1.EncryptionConfiguration{}
		outputStr, _ := mapInterfaceToJSON(output[0].(map[string]interface{}))
		expectedStr, _ := mapInterfaceToJSON(tc.ExpectedOutput[0].(map[string]interface{}))
		jsonToInterface(outputStr, outputObject)
		jsonToInterface(expectedStr, expectedObject)
		if !reflect.DeepEqual(outputObject, expectedObject) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenRKEClusterServicesKubeAPI(t *testing.T) {

	cases := []struct {
		Input          rancher.KubeAPIService
		ExpectedOutput []interface{}
	}{
		{
			testRKEClusterServicesKubeAPIConf,
			testRKEClusterServicesKubeAPIInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenRKEClusterServicesKubeAPI(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		outputObject := &rancher.KubeAPIService{}
		expectedObject := &rancher.KubeAPIService{}
		outputStr, _ := mapInterfaceToJSON(output[0].(map[string]interface{}))
		expectedStr, _ := mapInterfaceToJSON(tc.ExpectedOutput[0].(map[string]interface{}))
		jsonToInterface(outputStr, outputObject)
		jsonToInterface(expectedStr, expectedObject)
		if !reflect.DeepEqual(outputObject, expectedObject) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterServicesKubeAPIAuditLogConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rancher.AuditLogConfig
	}{
		{
			testRKEClusterServicesKubeAPIAuditLogConfigInterface,
			testRKEClusterServicesKubeAPIAuditLogConfigConf,
		},
	}

	for _, tc := range cases {
		output, err := expandRKEClusterServicesKubeAPIAuditLogConfig(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput.Policy, output.Policy)
		}
	}
}

func TestExpandRKEClusterServicesKubeAPIAuditLog(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rancher.AuditLog
	}{
		{
			testRKEClusterServicesKubeAPIAuditLogInterface,
			testRKEClusterServicesKubeAPIAuditLogConf,
		},
	}

	for _, tc := range cases {
		output, err := expandRKEClusterServicesKubeAPIAuditLog(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterServicesKubeAPIEventRateLimit(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rancher.EventRateLimit
	}{
		{
			testRKEClusterServicesKubeAPIEventRateLimitInterface,
			testRKEClusterServicesKubeAPIEventRateLimitConf,
		},
	}

	for _, tc := range cases {
		output, err := expandRKEClusterServicesKubeAPIEventRateLimit(tc.Input)
		if err != nil {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRKEClusterServicesKubeAPISecretsEncryptionConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rancher.SecretsEncryptionConfig
	}{
		{
			testRKEClusterServicesKubeAPISecretsEncryptionConfigInterface,
			testRKEClusterServicesKubeAPISecretsEncryptionConfigConf,
		},
	}

	for _, tc := range cases {
		output, err := expandRKEClusterServicesKubeAPISecretsEncryptionConfig(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", tc.Input)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput.CustomConfig, output.CustomConfig)
		}
	}
}

func TestExpandRKEClusterServicesKubeAPI(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rancher.KubeAPIService
	}{
		{
			testRKEClusterServicesKubeAPIInterface,
			testRKEClusterServicesKubeAPIConf,
		},
	}

	for _, tc := range cases {
		output, err := expandRKEClusterServicesKubeAPI(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
