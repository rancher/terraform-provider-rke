package rke

import (
	"reflect"
	"testing"

	rancher "github.com/rancher/types/apis/management.cattle.io/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	auditv1 "k8s.io/apiserver/pkg/apis/audit/v1"
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
		Kind:       "Policy",
		APIVersion: "audit.k8s.io/v1",
	}
	testRKEClusterServicesKubeAPIAuditLogConfigInterface = []interface{}{
		map[string]interface{}{
			"format":     "format",
			"max_age":    5,
			"max_backup": 10,
			"max_size":   100,
			"path":       "path",
			"policy":     `{"apiVersion":"audit.k8s.io/v1","kind":"Policy","metadata":{"creationTimestamp":null},"rules":[{"level":"RequestResponse","resources":[{"group":"*","resources":["pods"]}]}]}`,
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
	}
	testRKEClusterServicesKubeAPIEventRateLimitInterface = []interface{}{
		map[string]interface{}{
			"enabled": true,
		},
	}
	testRKEClusterServicesKubeAPISecretsEncryptionConfigConf = &rancher.SecretsEncryptionConfig{
		Enabled: true,
	}
	testRKEClusterServicesKubeAPISecretsEncryptionConfigInterface = []interface{}{
		map[string]interface{}{
			"enabled": true,
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
		output := flattenRKEClusterServicesKubeAPIEventRateLimit(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
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
		output := flattenRKEClusterServicesKubeAPISecretsEncryptionConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
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
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
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
		output := expandRKEClusterServicesKubeAPIEventRateLimit(tc.Input)
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
		output := expandRKEClusterServicesKubeAPISecretsEncryptionConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
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
