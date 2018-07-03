package rke

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"github.com/rancher/rke/cluster"
	"github.com/rancher/rke/hosts"
	"github.com/rancher/rke/pki"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/stretchr/testify/assert"
)

var (
	dummyCertificate   *x509.Certificate
	dummyPrivateKey    *rsa.PrivateKey
	dummyPrivateKeyPEM string
)

const dummyCertPEM = `-----BEGIN CERTIFICATE-----
MIIDujCCAqKgAwIBAgIIE31FZVaPXTUwDQYJKoZIhvcNAQEFBQAwSTELMAkGA1UE
BhMCVVMxEzARBgNVBAoTCkdvb2dsZSBJbmMxJTAjBgNVBAMTHEdvb2dsZSBJbnRl
cm5ldCBBdXRob3JpdHkgRzIwHhcNMTQwMTI5MTMyNzQzWhcNMTQwNTI5MDAwMDAw
WjBpMQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwN
TW91bnRhaW4gVmlldzETMBEGA1UECgwKR29vZ2xlIEluYzEYMBYGA1UEAwwPbWFp
bC5nb29nbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEfRrObuSW5T7q
5CnSEqefEmtH4CCv6+5EckuriNr1CjfVvqzwfAhopXkLrq45EQm8vkmf7W96XJhC
7ZM0dYi1/qOCAU8wggFLMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAa
BgNVHREEEzARgg9tYWlsLmdvb2dsZS5jb20wCwYDVR0PBAQDAgeAMGgGCCsGAQUF
BwEBBFwwWjArBggrBgEFBQcwAoYfaHR0cDovL3BraS5nb29nbGUuY29tL0dJQUcy
LmNydDArBggrBgEFBQcwAYYfaHR0cDovL2NsaWVudHMxLmdvb2dsZS5jb20vb2Nz
cDAdBgNVHQ4EFgQUiJxtimAuTfwb+aUtBn5UYKreKvMwDAYDVR0TAQH/BAIwADAf
BgNVHSMEGDAWgBRK3QYWG7z2aLV29YG2u2IaulqBLzAXBgNVHSAEEDAOMAwGCisG
AQQB1nkCBQEwMAYDVR0fBCkwJzAloCOgIYYfaHR0cDovL3BraS5nb29nbGUuY29t
L0dJQUcyLmNybDANBgkqhkiG9w0BAQUFAAOCAQEAH6RYHxHdcGpMpFE3oxDoFnP+
gtuBCHan2yE2GRbJ2Cw8Lw0MmuKqHlf9RSeYfd3BXeKkj1qO6TVKwCh+0HdZk283
TZZyzmEOyclm3UGFYe82P/iDFt+CeQ3NpmBg+GoaVCuWAARJN/KfglbLyyYygcQq
0SgeDh8dRKUiaW3HQSoYvTvdTuqzwK4CXsr3b5/dAOY8uMuG/IAR3FgwTbZ1dtoW
RvOTa8hYiU6A475WuZKyEHcwnGYe57u2I2KbMgcKjPniocj4QzgYsVAVKW3IwaOh
yE+vPxsiUkvQHdO2fojCkY8jg70jxM+gu59tPDNbw3Uh/2Ij310FgTHsnGQMyA==
-----END CERTIFICATE-----
`

func init() {
	block, _ := pem.Decode([]byte(dummyCertPEM))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}
	dummyCertificate = cert

	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic("failed to generate rsa private key: " + err.Error())
	}
	dummyPrivateKey = key
	dummyPrivateKeyPEM = privateKeyToPEM(key)
}

type dummyResourceData struct {
	values map[string]interface{}
}

func (d *dummyResourceData) GetOk(key string) (interface{}, bool) {
	v, ok := d.values[key]
	return v, ok
}

type dummyStateBuilder struct {
	values map[string]interface{}
}

func (d *dummyStateBuilder) Set(k string, v interface{}) error {
	d.values[k] = v
	return nil
}
func (d *dummyStateBuilder) SetId(id string) { // nolint
	d.values["Id"] = id
}

func TestParseResourceRKEConfigNode(t *testing.T) {

	testcases := []struct {
		caseName     string
		resourceData map[string]interface{}
		expectNodes  []v3.RKEConfigNode
		expectErrStr string
	}{
		{
			caseName: "minimum fields",
			resourceData: map[string]interface{}{
				"nodes": []interface{}{
					map[string]interface{}{
						"address": "192.2.0.1",
						"role":    []interface{}{"etcd"},
					},
				},
			},
			expectNodes: []v3.RKEConfigNode{
				{
					Address: "192.2.0.1",
					Role:    []string{"etcd"},
				},
			},
		},
		{
			caseName: "with both role and roles",
			resourceData: map[string]interface{}{
				"nodes": []interface{}{
					map[string]interface{}{
						"address": "192.2.0.1",
						"role":    []interface{}{"etcd"},
						"roles":   "etcd",
					},
				},
			},
			expectErrStr: "cannot specify both role and roles for a node",
		},
		{
			caseName: "without both role and roles",
			resourceData: map[string]interface{}{
				"nodes": []interface{}{
					map[string]interface{}{
						"address": "192.2.0.1",
					},
				},
			},
			expectErrStr: "either role or roles is required",
		},
		{
			caseName: "invalid role",
			resourceData: map[string]interface{}{
				"nodes": []interface{}{
					map[string]interface{}{
						"address": "192.2.0.1",
						"role":    []interface{}{"xxx"},
					},
				},
			},
			expectErrStr: `"role" must be one of [controlplane/etcd/worker]`,
		},
		{
			caseName: "invalid roles",
			resourceData: map[string]interface{}{
				"nodes": []interface{}{
					map[string]interface{}{
						"address": "192.2.0.1",
						"roles":   "xxx",
					},
				},
			},
			expectErrStr: `"role" must be one of [controlplane/etcd/worker]`,
		},
		{
			caseName: "use roles attr",
			resourceData: map[string]interface{}{
				"nodes": []interface{}{
					map[string]interface{}{
						"address": "192.2.0.1",
						"roles":   "controlplane,worker,etcd",
					},
				},
			},
			expectNodes: []v3.RKEConfigNode{
				{
					Address: "192.2.0.1",
					Role:    []string{"controlplane", "worker", "etcd"},
				},
			},
		},
		{
			caseName: "all fields",
			resourceData: map[string]interface{}{
				"nodes": []interface{}{
					map[string]interface{}{
						"node_name":         "node_name",
						"address":           "192.2.0.1",
						"port":              22,
						"internal_address":  "192.2.0.2",
						"role":              []interface{}{"controlplane", "worker", "etcd"},
						"hostname_override": "hostname_override",
						"user":              "rancher",
						"docker_socket":     "/var/run/docker.sock",
						"ssh_agent_auth":    true,
						"ssh_key":           "ssh_key",
						"ssh_key_path":      "ssh_key_path",
						"labels": map[string]interface{}{
							"foo": "foo",
							"bar": "bar",
						},
					},
				},
			},
			expectNodes: []v3.RKEConfigNode{
				{
					NodeName:         "node_name",
					Address:          "192.2.0.1",
					Port:             "22",
					InternalAddress:  "192.2.0.2",
					Role:             []string{"controlplane", "worker", "etcd"},
					HostnameOverride: "hostname_override",
					User:             "rancher",
					DockerSocket:     "/var/run/docker.sock",
					SSHAgentAuth:     true,
					SSHKey:           "ssh_key",
					SSHKeyPath:       "ssh_key_path",
					Labels: map[string]string{
						"foo": "foo",
						"bar": "bar",
					},
				},
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.caseName, func(t *testing.T) {
			d := &dummyResourceData{values: testcase.resourceData}
			nodes, err := parseResourceRKEConfigNodes(d)
			if testcase.expectErrStr == "" {
				assert.NoError(t, err)
				assert.EqualValues(t, testcase.expectNodes, nodes)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, testcase.expectErrStr)
			}
		})
	}
}

const (
	nodeJSON = `{
	  "address": "192.2.0.1",
      "role": ["controlplane", "worker", "etcd"]
	}`
	nodeYAML = `---
address: 192.2.0.2
role:
- controlplane
- worker
- etcd
`
)

func TestParseResourceRKEConfigNodesConf(t *testing.T) {

	testcases := []struct {
		caseName     string
		resourceData map[string]interface{}
		expectNodes  []v3.RKEConfigNode
		expectErrStr string
	}{
		{
			caseName: "JSON",
			resourceData: map[string]interface{}{
				"nodes_conf": []interface{}{nodeJSON},
			},
			expectNodes: []v3.RKEConfigNode{
				{
					Address: "192.2.0.1",
					Role:    []string{"controlplane", "worker", "etcd"},
				},
			},
		},
		{
			caseName: "YAML",
			resourceData: map[string]interface{}{
				"nodes_conf": []interface{}{nodeYAML},
			},
			expectNodes: []v3.RKEConfigNode{
				{
					Address: "192.2.0.2",
					Role:    []string{"controlplane", "worker", "etcd"},
				},
			},
		},
		{
			caseName: "both JSON and YAML",
			resourceData: map[string]interface{}{
				"nodes_conf": []interface{}{nodeJSON, nodeYAML},
			},
			expectNodes: []v3.RKEConfigNode{
				{
					Address: "192.2.0.1",
					Role:    []string{"controlplane", "worker", "etcd"},
				},
				{
					Address: "192.2.0.2",
					Role:    []string{"controlplane", "worker", "etcd"},
				},
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.caseName, func(t *testing.T) {
			d := &dummyResourceData{values: testcase.resourceData}
			nodes, err := parseResourceRKEConfigNodesConf(d)
			if testcase.expectErrStr == "" {
				assert.NoError(t, err)
				assert.EqualValues(t, testcase.expectNodes, nodes)
			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, testcase.expectErrStr)
			}
		})
	}
}

func TestParseResourceETCDService(t *testing.T) {
	testcases := []struct {
		caseName      string
		resourceData  map[string]interface{}
		expectService *v3.ETCDService
	}{
		{
			caseName: "all fields",
			resourceData: map[string]interface{}{
				"services_etcd": []interface{}{
					map[string]interface{}{
						"image": "image",
						"extra_args": map[string]interface{}{
							"foo": "foo",
							"bar": "bar",
						},
						"extra_binds":   []interface{}{"/etc1", "/etc2"},
						"extra_env":     []interface{}{"env1=val1", "env2=val2"},
						"external_urls": []interface{}{"https://etcd1.example.com", "https://etcd2.example.com"},
						"ca_cert":       "ca_cert",
						"cert":          "cert",
						"key":           "key",
						"path":          "path",
						"snapshot":      true,
						"retention":     "retention",
						"creation":      "creation",
					},
				},
			},
			expectService: &v3.ETCDService{
				BaseService: v3.BaseService{
					Image: "image",
					ExtraArgs: map[string]string{
						"foo": "foo",
						"bar": "bar",
					},
					ExtraBinds: []string{"/etc1", "/etc2"},
					ExtraEnv:   []string{"env1=val1", "env2=val2"},
				},
				ExternalURLs: []string{"https://etcd1.example.com", "https://etcd2.example.com"},
				CACert:       "ca_cert",
				Cert:         "cert",
				Key:          "key",
				Path:         "path",
				Snapshot:     true,
				Retention:    "retention",
				Creation:     "creation",
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.caseName, func(t *testing.T) {
			d := &dummyResourceData{values: testcase.resourceData}
			service, err := parseResourceETCDService(d)
			assert.NoError(t, err)
			assert.EqualValues(t, testcase.expectService, service)
		})
	}
}

func TestParseResourceKubeAPIService(t *testing.T) {
	testcases := []struct {
		caseName      string
		resourceData  map[string]interface{}
		expectService *v3.KubeAPIService
	}{
		{
			caseName: "all fields",
			resourceData: map[string]interface{}{
				"services_kube_api": []interface{}{
					map[string]interface{}{
						"image": "image",
						"extra_args": map[string]interface{}{
							"foo": "foo",
							"bar": "bar",
						},
						"extra_binds":              []interface{}{"/etc1", "/etc2"},
						"extra_env":                []interface{}{"env1=val1", "env2=val2"},
						"service_cluster_ip_range": "10.240.0.0/16",
						"service_node_port_range":  "30000-31000",
						"pod_security_policy":      true,
					},
				},
			},
			expectService: &v3.KubeAPIService{
				BaseService: v3.BaseService{
					Image: "image",
					ExtraArgs: map[string]string{
						"foo": "foo",
						"bar": "bar",
					},
					ExtraBinds: []string{"/etc1", "/etc2"},
					ExtraEnv:   []string{"env1=val1", "env2=val2"},
				},
				ServiceClusterIPRange: "10.240.0.0/16",
				ServiceNodePortRange:  "30000-31000",
				PodSecurityPolicy:     true,
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.caseName, func(t *testing.T) {
			d := &dummyResourceData{values: testcase.resourceData}
			service, err := parseResourceKubeAPIService(d)
			assert.NoError(t, err)
			assert.EqualValues(t, testcase.expectService, service)
		})
	}
}

func TestParseResourceKubeControllerService(t *testing.T) {
	testcases := []struct {
		caseName      string
		resourceData  map[string]interface{}
		expectService *v3.KubeControllerService
	}{
		{
			caseName: "all fields",
			resourceData: map[string]interface{}{
				"services_kube_controller": []interface{}{
					map[string]interface{}{
						"image": "image",
						"extra_args": map[string]interface{}{
							"foo": "foo",
							"bar": "bar",
						},
						"extra_binds":              []interface{}{"/etc1", "/etc2"},
						"extra_env":                []interface{}{"env1=val1", "env2=val2"},
						"cluster_cidr":             "10.240.0.0/16",
						"service_cluster_ip_range": "10.240.0.0/16",
					},
				},
			},
			expectService: &v3.KubeControllerService{
				BaseService: v3.BaseService{
					Image: "image",
					ExtraArgs: map[string]string{
						"foo": "foo",
						"bar": "bar",
					},
					ExtraBinds: []string{"/etc1", "/etc2"},
					ExtraEnv:   []string{"env1=val1", "env2=val2"},
				},
				ClusterCIDR:           "10.240.0.0/16",
				ServiceClusterIPRange: "10.240.0.0/16",
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.caseName, func(t *testing.T) {
			d := &dummyResourceData{values: testcase.resourceData}
			service, err := parseResourceKubeControllerService(d)
			assert.NoError(t, err)
			assert.EqualValues(t, testcase.expectService, service)
		})
	}
}

func TestParseResourceSchedulerService(t *testing.T) {
	testcases := []struct {
		caseName      string
		resourceData  map[string]interface{}
		expectService *v3.SchedulerService
	}{
		{
			caseName: "all fields",
			resourceData: map[string]interface{}{
				"services_scheduler": []interface{}{
					map[string]interface{}{
						"image": "image",
						"extra_args": map[string]interface{}{
							"foo": "foo",
							"bar": "bar",
						},
						"extra_binds": []interface{}{"/etc1", "/etc2"},
						"extra_env":   []interface{}{"env1=val1", "env2=val2"},
					},
				},
			},
			expectService: &v3.SchedulerService{
				BaseService: v3.BaseService{
					Image: "image",
					ExtraArgs: map[string]string{
						"foo": "foo",
						"bar": "bar",
					},
					ExtraBinds: []string{"/etc1", "/etc2"},
					ExtraEnv:   []string{"env1=val1", "env2=val2"},
				},
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.caseName, func(t *testing.T) {
			d := &dummyResourceData{values: testcase.resourceData}
			service, err := parseResourceSchedulerService(d)
			assert.NoError(t, err)
			assert.EqualValues(t, testcase.expectService, service)
		})
	}
}

func TestParseResourceKubeletService(t *testing.T) {
	testcases := []struct {
		caseName      string
		resourceData  map[string]interface{}
		expectService *v3.KubeletService
	}{
		{
			caseName: "all fields",
			resourceData: map[string]interface{}{
				"services_kubelet": []interface{}{
					map[string]interface{}{
						"image": "image",
						"extra_args": map[string]interface{}{
							"foo": "foo",
							"bar": "bar",
						},
						"extra_binds":           []interface{}{"/etc1", "/etc2"},
						"extra_env":             []interface{}{"env1=val1", "env2=val2"},
						"cluster_domain":        "example.com",
						"infra_container_image": "alpine:latest",
						"cluster_dns_server":    "192.2.0.1",
						"fail_swap_on":          true,
					},
				},
			},
			expectService: &v3.KubeletService{
				BaseService: v3.BaseService{
					Image: "image",
					ExtraArgs: map[string]string{
						"foo": "foo",
						"bar": "bar",
					},
					ExtraBinds: []string{"/etc1", "/etc2"},
					ExtraEnv:   []string{"env1=val1", "env2=val2"},
				},
				ClusterDomain:       "example.com",
				InfraContainerImage: "alpine:latest",
				ClusterDNSServer:    "192.2.0.1",
				FailSwapOn:          true,
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.caseName, func(t *testing.T) {
			d := &dummyResourceData{values: testcase.resourceData}
			service, err := parseResourceKubeletService(d)
			assert.NoError(t, err)
			assert.EqualValues(t, testcase.expectService, service)
		})
	}
}

func TestParseResourceKubeproxyService(t *testing.T) {
	testcases := []struct {
		caseName      string
		resourceData  map[string]interface{}
		expectService *v3.KubeproxyService
	}{
		{
			caseName: "all fields",
			resourceData: map[string]interface{}{
				"services_kubeproxy": []interface{}{
					map[string]interface{}{
						"image": "image",
						"extra_args": map[string]interface{}{
							"foo": "foo",
							"bar": "bar",
						},
						"extra_binds": []interface{}{"/etc1", "/etc2"},
						"extra_env":   []interface{}{"env1=val1", "env2=val2"},
					},
				},
			},
			expectService: &v3.KubeproxyService{
				BaseService: v3.BaseService{
					Image: "image",
					ExtraArgs: map[string]string{
						"foo": "foo",
						"bar": "bar",
					},
					ExtraBinds: []string{"/etc1", "/etc2"},
					ExtraEnv:   []string{"env1=val1", "env2=val2"},
				},
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.caseName, func(t *testing.T) {
			d := &dummyResourceData{values: testcase.resourceData}
			service, err := parseResourceKubeproxyService(d)
			assert.NoError(t, err)
			assert.EqualValues(t, testcase.expectService, service)
		})
	}
}

func TestParseResourceNetwork(t *testing.T) {
	testcases := []struct {
		caseName      string
		resourceData  map[string]interface{}
		expectNetwork *v3.NetworkConfig
	}{
		{
			caseName: "all fields",
			resourceData: map[string]interface{}{
				"network": []interface{}{
					map[string]interface{}{
						"plugin": "calico",
						"options": map[string]interface{}{
							"foo": "foo",
							"bar": "bar",
						},
					},
				},
			},
			expectNetwork: &v3.NetworkConfig{
				Plugin: "calico",
				Options: map[string]string{
					"foo": "foo",
					"bar": "bar",
				},
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.caseName, func(t *testing.T) {
			d := &dummyResourceData{values: testcase.resourceData}
			network, err := parseResourceNetwork(d)
			assert.NoError(t, err)
			assert.EqualValues(t, testcase.expectNetwork, network)
		})
	}
}

func TestParseResourceAuthentication(t *testing.T) {
	testcases := []struct {
		caseName     string
		resourceData map[string]interface{}
		expectConfig *v3.AuthnConfig
	}{
		{
			caseName: "all fields",
			resourceData: map[string]interface{}{
				"authentication": []interface{}{
					map[string]interface{}{
						"strategy": "x509",
						"options": map[string]interface{}{
							"foo": "foo",
							"bar": "bar",
						},
						"sans": []interface{}{
							"192.2.0.1",
							"test.example.com",
						},
					},
				},
			},
			expectConfig: &v3.AuthnConfig{
				Strategy: "x509",
				Options: map[string]string{
					"foo": "foo",
					"bar": "bar",
				},
				SANs: []string{
					"192.2.0.1",
					"test.example.com",
				},
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.caseName, func(t *testing.T) {
			d := &dummyResourceData{values: testcase.resourceData}
			config, err := parseResourceAuthentication(d)
			assert.NoError(t, err)
			assert.EqualValues(t, testcase.expectConfig, config)
		})
	}
}

func TestParseResourceAddons(t *testing.T) {
	d := &dummyResourceData{values: map[string]interface{}{"addons": "addons: yaml"}}
	addon, err := parseResourceAddons(d)
	assert.NoError(t, err)
	assert.EqualValues(t, "addons: yaml", addon)
}

func TestParseResourceAddonsInclude(t *testing.T) {
	expect := []string{
		"https://example.com/addon1.yaml",
		"https://example.com/addon2.yaml",
	}
	d := &dummyResourceData{
		values: map[string]interface{}{
			"addons_include": []interface{}{
				"https://example.com/addon1.yaml",
				"https://example.com/addon2.yaml",
			},
		},
	}
	includes, err := parseResourceAddonsInclude(d)
	assert.NoError(t, err)
	assert.EqualValues(t, expect, includes)
}

func TestParseResourceAddonJobTimeout(t *testing.T) {
	d := &dummyResourceData{values: map[string]interface{}{"addon_job_timeout": 10}}
	v, err := parseResourceAddonJobTimeout(d)
	assert.NoError(t, err)
	assert.EqualValues(t, 10, v)
}

func TestParseResourceSystemImages(t *testing.T) {
	testcases := []struct {
		caseName     string
		resourceData map[string]interface{}
		expectConfig *v3.RKESystemImages
	}{
		{
			caseName: "all fields",
			resourceData: map[string]interface{}{
				"system_images": []interface{}{
					map[string]interface{}{
						"etcd":                        "etcd",
						"alpine":                      "alpine",
						"nginx_proxy":                 "nginx_proxy",
						"cert_downloader":             "cert_downloader",
						"kubernetes_services_sidecar": "kubernetes_services_sidecar",
						"kube_dns":                    "kube_dns",
						"dnsmasq":                     "dnsmasq",
						"kube_dns_sidecar":            "kube_dns_sidecar",
						"kube_dns_autoscaler":         "kube_dns_autoscaler",
						"kubernetes":                  "kubernetes",
						"flannel":                     "flannel",
						"flannel_cni":                 "flannel_cni",
						"calico_node":                 "calico_node",
						"calico_cni":                  "calico_cni",
						"calico_controllers":          "calico_controllers",
						"calico_ctl":                  "calico_ctl",
						"canal_node":                  "canal_node",
						"canal_cni":                   "canal_cni",
						"canal_flannel":               "canal_flannel",
						"weave_node":                  "weave_node",
						"weave_cni":                   "weave_cni",
						"pod_infra_container":         "pod_infra_container",
						"ingress":                     "ingress",
						"ingress_backend":             "ingress_backend",
						"dashboard":                   "dashboard",
						"heapster":                    "heapster",
						"grafana":                     "grafana",
						"influxdb":                    "influxdb",
						"tiller":                      "tiller",
					},
				},
			},
			expectConfig: &v3.RKESystemImages{
				Etcd:                      "etcd",
				Alpine:                    "alpine",
				NginxProxy:                "nginx_proxy",
				CertDownloader:            "cert_downloader",
				KubernetesServicesSidecar: "kubernetes_services_sidecar",
				KubeDNS:                   "kube_dns",
				DNSmasq:                   "dnsmasq",
				KubeDNSSidecar:            "kube_dns_sidecar",
				KubeDNSAutoscaler:         "kube_dns_autoscaler",
				Kubernetes:                "kubernetes",
				Flannel:                   "flannel",
				FlannelCNI:                "flannel_cni",
				CalicoNode:                "calico_node",
				CalicoCNI:                 "calico_cni",
				CalicoControllers:         "calico_controllers",
				CalicoCtl:                 "calico_ctl",
				CanalNode:                 "canal_node",
				CanalCNI:                  "canal_cni",
				CanalFlannel:              "canal_flannel",
				WeaveNode:                 "weave_node",
				WeaveCNI:                  "weave_cni",
				PodInfraContainer:         "pod_infra_container",
				Ingress:                   "ingress",
				IngressBackend:            "ingress_backend",
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.caseName, func(t *testing.T) {
			d := &dummyResourceData{values: testcase.resourceData}
			config, err := parseResourceSystemImages(d)
			assert.NoError(t, err)
			assert.EqualValues(t, testcase.expectConfig, config)
		})
	}
}

func TestParseResourceSSHKeyPath(t *testing.T) {
	d := &dummyResourceData{values: map[string]interface{}{"ssh_key_path": "ssh_key_path"}}
	keyPath, err := parseResourceSSHKeyPath(d)
	assert.NoError(t, err)
	assert.EqualValues(t, "ssh_key_path", keyPath)
}

func TestParseResourceSSHAgentAuth(t *testing.T) {
	d := &dummyResourceData{values: map[string]interface{}{"ssh_agent_auth": true}}
	auth, err := parseResourceSSHAgentAuth(d)
	assert.NoError(t, err)
	assert.EqualValues(t, true, auth)
}

func TestParseResourceBastionHost(t *testing.T) {

	testcases := []struct {
		caseName     string
		resourceData map[string]interface{}
		expectConfig *v3.BastionHost
	}{
		{
			caseName: "all fields",
			resourceData: map[string]interface{}{
				"bastion_host": []interface{}{
					map[string]interface{}{
						"address":        "192.2.0.1",
						"port":           22,
						"user":           "rancher",
						"ssh_agent_auth": true,
						"ssh_key":        "ssh_key",
						"ssh_key_path":   "ssh_key_path",
					},
				},
			},
			expectConfig: &v3.BastionHost{
				Address:      "192.2.0.1",
				Port:         "22",
				User:         "rancher",
				SSHAgentAuth: true,
				SSHKey:       "ssh_key",
				SSHKeyPath:   "ssh_key_path",
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.caseName, func(t *testing.T) {
			d := &dummyResourceData{values: testcase.resourceData}
			host, err := parseResourceBastionHost(d)
			assert.NoError(t, err)
			assert.EqualValues(t, testcase.expectConfig, host)
		})
	}
}

func TestParseResourceAuthorization(t *testing.T) {
	testcases := []struct {
		caseName     string
		resourceData map[string]interface{}
		expectConfig *v3.AuthzConfig
	}{
		{
			caseName: "all fields",
			resourceData: map[string]interface{}{
				"authorization": []interface{}{
					map[string]interface{}{
						"mode": "rbac",
						"options": map[string]interface{}{
							"foo": "foo",
							"bar": "bar",
						},
					},
				},
			},
			expectConfig: &v3.AuthzConfig{
				Mode: "rbac",
				Options: map[string]string{
					"foo": "foo",
					"bar": "bar",
				},
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.caseName, func(t *testing.T) {
			d := &dummyResourceData{values: testcase.resourceData}
			config, err := parseResourceAuthorization(d)
			assert.NoError(t, err)
			assert.EqualValues(t, testcase.expectConfig, config)
		})
	}
}

func TestParseResourceIgnoreDockerVersion(t *testing.T) {
	d := &dummyResourceData{values: map[string]interface{}{"ignore_docker_version": true}}
	ignore, err := parseResourceIgnoreDockerVersion(d)
	assert.NoError(t, err)
	assert.EqualValues(t, true, ignore)
}

func TestParseResourceKubernetesVersion(t *testing.T) {
	d := &dummyResourceData{
		values: map[string]interface{}{
			"kubernetes_version": "1.8.9",
		},
	}
	version, err := parseResourceVersion(d)
	assert.NoError(t, err)
	assert.EqualValues(t, "1.8.9", version)
}

func TestParseResourcePrivateRegistries(t *testing.T) {
	testcases := []struct {
		caseName     string
		resourceData map[string]interface{}
		expectConfig []v3.PrivateRegistry
	}{
		{
			caseName: "all fields",
			resourceData: map[string]interface{}{
				"private_registries": []interface{}{
					map[string]interface{}{
						"url":      "https://example.com",
						"user":     "rancher",
						"password": "p@ssword",
					},
				},
			},
			expectConfig: []v3.PrivateRegistry{
				{
					URL:      "https://example.com",
					User:     "rancher",
					Password: "p@ssword",
				},
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.caseName, func(t *testing.T) {
			d := &dummyResourceData{values: testcase.resourceData}
			config, err := parseResourcePrivateRegistries(d)
			assert.NoError(t, err)
			assert.EqualValues(t, testcase.expectConfig, config)
		})
	}
}

func TestParseResourceIngress(t *testing.T) {
	testcases := []struct {
		caseName     string
		resourceData map[string]interface{}
		expectConfig *v3.IngressConfig
	}{
		{
			caseName: "all fields",
			resourceData: map[string]interface{}{
				"ingress": []interface{}{
					map[string]interface{}{
						"provider": "nginx",
						"options": map[string]interface{}{
							"foo": "foo",
							"bar": "bar",
						},
						"node_selector": map[string]interface{}{
							"role": "worker",
						},
						"extra_args": map[string]interface{}{
							"foo": "foo",
							"bar": "bar",
						},
					},
				},
			},
			expectConfig: &v3.IngressConfig{
				Provider: "nginx",
				Options: map[string]string{
					"foo": "foo",
					"bar": "bar",
				},
				NodeSelector: map[string]string{
					"role": "worker",
				},
				ExtraArgs: map[string]string{
					"foo": "foo",
					"bar": "bar",
				},
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.caseName, func(t *testing.T) {
			d := &dummyResourceData{values: testcase.resourceData}
			config, err := parseResourceIngress(d)
			assert.NoError(t, err)
			assert.EqualValues(t, testcase.expectConfig, config)
		})
	}
}

func TestParseResourceClusterName(t *testing.T) {
	d := &dummyResourceData{
		values: map[string]interface{}{
			"cluster_name": "rke",
		},
	}
	name, err := parseResourceClusterName(d)
	assert.NoError(t, err)
	assert.EqualValues(t, "rke", name)
}

func TestParseResourceCloudProvider(t *testing.T) {
	testcases := []struct {
		caseName     string
		resourceData map[string]interface{}
		expectConfig *v3.CloudProvider
	}{
		{
			caseName: "all fields",
			resourceData: map[string]interface{}{
				"cloud_provider": []interface{}{
					map[string]interface{}{
						"name": "sakuracloud",
						// TODO not implements on RKE v0.1.8-rc2
						//"aws_cloud_provider": []interface{}{
						//	map[string]interface{}{},
						//},
						"azure_cloud_config": []interface{}{
							map[string]interface{}{
								"cloud":                            "cloud",
								"tenant_id":                        "tenant_id",
								"subscription_id":                  "subscription_id",
								"resource_group":                   "resource_group",
								"location":                         "location",
								"vnet_name":                        "vnet_name",
								"vnet_resource_group":              "vnet_resource_group",
								"route_table_name":                 "route_table_name",
								"primary_availability_set_name":    "primary_availability_set_name",
								"vm_type":                          "vm_type",
								"primary_scale_set_name":           "primary_scale_set_name",
								"aad_client_id":                    "aad_client_id",
								"aad_client_secret":                "aad_client_secret",
								"aad_client_cert_path":             "aad_client_cert_path",
								"aad_client_cert_password":         "aad_client_cert_password",
								"cloud_provider_backoff":           true,
								"cloud_provider_backoff_retries":   1,
								"cloud_provider_backoff_exponent":  2,
								"cloud_provider_backoff_duration":  3,
								"cloud_provider_backoff_jitter":    4,
								"cloud_provider_rate_limit":        true,
								"cloud_provider_rate_limit_qps":    1,
								"cloud_provider_rate_limit_bucket": 2,
								"use_instance_metadata":            true,
								"use_managed_identity_extension":   true,
								"maximum_load_balancer_rule_count": 1,
							},
						},
						"vsphere_cloud_config": []interface{}{
							map[string]interface{}{
								"global": []interface{}{
									map[string]interface{}{
										"user":                 "user",
										"password":             "password",
										"server":               "server",
										"port":                 "port",
										"insecure_flag":        true,
										"datacenter":           "datacenter",
										"datacenters":          "datacenters",
										"datastore":            "datastore",
										"working_dir":          "working_dir",
										"soap_roundtrip_count": 1,
										"vm_uuid":              "vm_uuid",
										"vm_name":              "vm_name",
									},
								},
								"virtual_center": []interface{}{
									map[string]interface{}{
										"server":               "192.2.0.1",
										"user":                 "user",
										"password":             "password",
										"port":                 "port",
										"datacenters":          "datacenters",
										"soap_roundtrip_count": 1,
									},
								},
								"network": []interface{}{
									map[string]interface{}{
										"public_network": "public_network",
									},
								},
								"disk": []interface{}{
									map[string]interface{}{
										"scsi_controller_type": "scsi_controller_type",
									},
								},
								"workspace": []interface{}{
									map[string]interface{}{
										"server":            "server",
										"datacenter":        "datacenter",
										"folder":            "folder",
										"default_datastore": "default_datastore",
										"resourcepool_path": "resourcepool_path",
									},
								},
							},
						},
						"openstack_cloud_config": []interface{}{
							map[string]interface{}{
								"global": []interface{}{
									map[string]interface{}{
										"auth_url":    "auth_url",
										"username":    "username",
										"user_id":     "user_id",
										"password":    "password",
										"tenant_id":   "tenant_id",
										"tenant_name": "tenant_name",
										"trust_id":    "trust_id",
										"domain_id":   "domain_id",
										"domain_name": "domain_name",
										"region":      "region",
										"ca_file":     "ca_file",
									},
								},
								"load_balancer": []interface{}{
									map[string]interface{}{
										"lb_version":             "lb_version",
										"use_octavia":            true,
										"subnet_id":              "subnet_id",
										"floating_network_id":    "floating_network_id",
										"lb_method":              "lb_method",
										"lb_provider":            "lb_provider",
										"create_monitor":         true,
										"monitor_delay":          1,
										"monitor_timeout":        2,
										"monitor_max_retries":    3,
										"manage_security_groups": true,
									},
								},
								"block_storage": []interface{}{
									map[string]interface{}{
										"bs_version":        "bs_version",
										"trust_device_path": true,
										"ignore_volume_az":  true,
									},
								},
								"route": []interface{}{
									map[string]interface{}{
										"router_id": "router_id",
									},
								},
								"metadata": []interface{}{
									map[string]interface{}{
										"search_order":    "search_order",
										"request_timeout": 1,
									},
								},
							},
						},
						"custom_cloud_config": "custom_cloud_config",
					},
				},
			},
			expectConfig: &v3.CloudProvider{
				Name: "sakuracloud",
				AzureCloudProvider: &v3.AzureCloudProvider{
					Cloud:                      "cloud",
					TenantID:                   "tenant_id",
					SubscriptionID:             "subscription_id",
					ResourceGroup:              "resource_group",
					Location:                   "location",
					VnetName:                   "vnet_name",
					VnetResourceGroup:          "vnet_resource_group",
					RouteTableName:             "route_table_name",
					PrimaryAvailabilitySetName: "primary_availability_set_name",
					VMType:                       "vm_type",
					PrimaryScaleSetName:          "primary_scale_set_name",
					AADClientID:                  "aad_client_id",
					AADClientSecret:              "aad_client_secret",
					AADClientCertPath:            "aad_client_cert_path",
					AADClientCertPassword:        "aad_client_cert_password",
					CloudProviderBackoff:         true,
					CloudProviderBackoffRetries:  1,
					CloudProviderBackoffExponent: 2,
					CloudProviderBackoffDuration: 3,
					CloudProviderBackoffJitter:   4,
					CloudProviderRateLimit:       true,
					CloudProviderRateLimitQPS:    1,
					CloudProviderRateLimitBucket: 2,
					UseInstanceMetadata:          true,
					UseManagedIdentityExtension:  true,
					MaximumLoadBalancerRuleCount: 1,
				},
				VsphereCloudProvider: &v3.VsphereCloudProvider{
					Global: v3.GlobalVsphereOpts{
						User:              "user",
						Password:          "password",
						VCenterIP:         "server",
						VCenterPort:       "port",
						InsecureFlag:      true,
						Datacenter:        "datacenter",
						Datacenters:       "datacenters",
						DefaultDatastore:  "datastore",
						WorkingDir:        "working_dir",
						RoundTripperCount: 1,
						VMUUID:            "vm_uuid",
						VMName:            "vm_name",
					},
					VirtualCenter: map[string]v3.VirtualCenterConfig{
						"192.2.0.1": {
							User:              "user",
							Password:          "password",
							VCenterPort:       "port",
							Datacenters:       "datacenters",
							RoundTripperCount: 1,
						},
					},
					Network: v3.NetworkVshpereOpts{
						PublicNetwork: "public_network",
					},
					Disk: v3.DiskVsphereOpts{
						SCSIControllerType: "scsi_controller_type",
					},
					Workspace: v3.WorkspaceVsphereOpts{
						VCenterIP:        "server",
						Datacenter:       "datacenter",
						Folder:           "folder",
						DefaultDatastore: "default_datastore",
						ResourcePoolPath: "resourcepool_path",
					},
				},
				OpenstackCloudProvider: &v3.OpenstackCloudProvider{
					Global: v3.GlobalOpenstackOpts{
						AuthURL:    "auth_url",
						Username:   "username",
						UserID:     "user_id",
						Password:   "password",
						TenantID:   "tenant_id",
						TenantName: "tenant_name",
						TrustID:    "trust_id",
						DomainID:   "domain_id",
						DomainName: "domain_name",
						Region:     "region",
						CAFile:     "ca_file",
					},
					LoadBalancer: v3.LoadBalancerOpenstackOpts{
						LBVersion:            "lb_version",
						UseOctavia:           true,
						SubnetID:             "subnet_id",
						FloatingNetworkID:    "floating_network_id",
						LBMethod:             "lb_method",
						LBProvider:           "lb_provider",
						CreateMonitor:        true,
						MonitorDelay:         1,
						MonitorTimeout:       2,
						MonitorMaxRetries:    3,
						ManageSecurityGroups: true,
					},
					BlockStorage: v3.BlockStorageOpenstackOpts{
						BSVersion:       "bs_version",
						TrustDevicePath: true,
						IgnoreVolumeAZ:  true,
					},
					Route: v3.RouteOpenstackOpts{
						RouterID: "router_id",
					},
					Metadata: v3.MetadataOpenstackOpts{
						SearchOrder:    "search_order",
						RequestTimeout: 1,
					},
				},
				CustomCloudProvider: "custom_cloud_config",
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.caseName, func(t *testing.T) {
			d := &dummyResourceData{values: testcase.resourceData}
			config, err := parseResourceCloudProvider(d)
			assert.NoError(t, err)
			assert.EqualValues(t, testcase.expectConfig, config)
		})
	}
}

func TestParseResourcePrefixPath(t *testing.T) {
	d := &dummyResourceData{values: map[string]interface{}{"prefix_path": "prefix_path"}}
	prefixPath, err := parseResourcePrefixPath(d)
	assert.NoError(t, err)
	assert.EqualValues(t, "prefix_path", prefixPath)
}

func TestClusterToState(t *testing.T) {

	testcases := []struct {
		caseName string
		cluster  *cluster.Cluster
		state    map[string]interface{}
	}{
		{
			caseName: "all fields",
			cluster: &cluster.Cluster{
				RancherKubernetesEngineConfig: v3.RancherKubernetesEngineConfig{
					Services: v3.RKEConfigServices{
						Etcd: v3.ETCDService{
							BaseService: v3.BaseService{
								Image: "etcd:latest",
								ExtraArgs: map[string]string{
									"foo": "bar",
									"bar": "foo",
								},
								ExtraBinds: []string{"/bind1", "/bind2"},
								ExtraEnv:   []string{"env1=val1", "env2=val2"},
							},
							ExternalURLs: []string{
								"https://ext1.example.com",
								"https://ext2.example.com",
							},
							CACert:    "ca_cert",
							Cert:      "cert",
							Key:       "key",
							Path:      "path",
							Snapshot:  true,
							Retention: "retention",
							Creation:  "creation",
						},
						KubeAPI: v3.KubeAPIService{
							BaseService: v3.BaseService{
								Image: "kube_api:latest",
								ExtraArgs: map[string]string{
									"foo": "bar",
									"bar": "foo",
								},
								ExtraBinds: []string{"/bind1", "/bind2"},
								ExtraEnv:   []string{"env1=val1", "env2=val2"},
							},
							ServiceClusterIPRange: "10.240.0.0/16",
							ServiceNodePortRange:  "30000-31000",
							PodSecurityPolicy:     true,
						},
						KubeController: v3.KubeControllerService{
							BaseService: v3.BaseService{
								Image: "kube_controller:latest",
								ExtraArgs: map[string]string{
									"foo": "bar",
									"bar": "foo",
								},
								ExtraBinds: []string{"/bind1", "/bind2"},
								ExtraEnv:   []string{"env1=val1", "env2=val2"},
							},
							ClusterCIDR:           "10.200.0.0/8",
							ServiceClusterIPRange: "10.240.0.0/16",
						},
						Scheduler: v3.SchedulerService{
							BaseService: v3.BaseService{
								Image: "scheduler:latest",
								ExtraArgs: map[string]string{
									"foo": "bar",
									"bar": "foo",
								},
								ExtraBinds: []string{"/bind1", "/bind2"},
								ExtraEnv:   []string{"env1=val1", "env2=val2"},
							},
						},
						Kubelet: v3.KubeletService{
							BaseService: v3.BaseService{
								Image: "kubelet:latest",
								ExtraArgs: map[string]string{
									"foo": "bar",
									"bar": "foo",
								},
								ExtraBinds: []string{"/bind1", "/bind2"},
								ExtraEnv:   []string{"env1=val1", "env2=val2"},
							},
							ClusterDomain:       "example.com",
							InfraContainerImage: "alpine:latest",
							ClusterDNSServer:    "192.2.0.1",
							FailSwapOn:          true,
						},
						Kubeproxy: v3.KubeproxyService{
							BaseService: v3.BaseService{
								Image: "kubeproxy:latest",
								ExtraArgs: map[string]string{
									"foo": "bar",
									"bar": "foo",
								},
								ExtraBinds: []string{"/bind1", "/bind2"},
								ExtraEnv:   []string{"env1=val1", "env2=val2"},
							},
						},
					},
					Network: v3.NetworkConfig{
						Plugin: "calico",
						Options: map[string]string{
							"foo": "bar",
							"bar": "foo",
						},
					},
					Authentication: v3.AuthnConfig{
						Strategy: "x509",
						Options: map[string]string{
							"foo": "bar",
							"bar": "foo",
						},
						SANs: []string{"sans1", "sans2"},
					},
					Addons: "addons: yaml",
					AddonsInclude: []string{
						"https://example.com/addon1.yaml",
						"https://example.com/addon2.yaml",
					},
					AddonJobTimeout: 10,
					SystemImages: v3.RKESystemImages{
						Etcd:                      "etcd",
						Alpine:                    "alpine",
						NginxProxy:                "nginx_proxy",
						CertDownloader:            "cert_downloader",
						KubernetesServicesSidecar: "kubernetes_services_sidecar",
						KubeDNS:                   "kube_dns",
						DNSmasq:                   "dnsmasq",
						KubeDNSSidecar:            "kube_dns_sidecar",
						KubeDNSAutoscaler:         "kube_dns_autoscaler",
						Kubernetes:                "kubernetes",
						Flannel:                   "flannel",
						FlannelCNI:                "flannel_cni",
						CalicoNode:                "calico_node",
						CalicoCNI:                 "calico_cni",
						CalicoControllers:         "calico_controllers",
						CalicoCtl:                 "calico_ctl",
						CanalNode:                 "canal_node",
						CanalCNI:                  "canal_cni",
						CanalFlannel:              "canal_flannel",
						WeaveNode:                 "weave_node",
						WeaveCNI:                  "weave_cni",
						PodInfraContainer:         "pod_infra_container",
						Ingress:                   "ingress",
						IngressBackend:            "ingress_backend",
					},
					SSHKeyPath:   "ssh_key_path",
					SSHAgentAuth: true,
					BastionHost: v3.BastionHost{
						Address:      "192.2.0.1",
						Port:         "22",
						User:         "rancher",
						SSHAgentAuth: true,
						SSHKey:       "ssh_key",
						SSHKeyPath:   "ssh_key_path",
					},
					Authorization: v3.AuthzConfig{
						Mode: "rbac",
						Options: map[string]string{
							"foo": "bar",
							"bar": "foo",
						},
					},
					IgnoreDockerVersion: true,
					Version:             "1.8.9",
					PrivateRegistries: []v3.PrivateRegistry{
						{
							URL:      "https://registry1.example.com",
							User:     "user1",
							Password: "password1",
						},
						{
							URL:      "https://registry2.example.com",
							User:     "user2",
							Password: "password2",
						},
					},
					Ingress: v3.IngressConfig{
						Provider: "nginx",
						Options: map[string]string{
							"foo": "bar",
							"bar": "foo",
						},
						NodeSelector: map[string]string{
							"role": "worker",
						},
						ExtraArgs: map[string]string{
							"foo": "foo",
							"bar": "bar",
						},
					},
					ClusterName: "example",
					CloudProvider: v3.CloudProvider{
						Name: "sakuracloud",
						AzureCloudProvider: &v3.AzureCloudProvider{
							Cloud:                      "cloud",
							TenantID:                   "tenant_id",
							SubscriptionID:             "subscription_id",
							ResourceGroup:              "resource_group",
							Location:                   "location",
							VnetName:                   "vnet_name",
							VnetResourceGroup:          "vnet_resource_group",
							RouteTableName:             "route_table_name",
							PrimaryAvailabilitySetName: "primary_availability_set_name",
							VMType:                       "vm_type",
							PrimaryScaleSetName:          "primary_scale_set_name",
							AADClientID:                  "aad_client_id",
							AADClientSecret:              "aad_client_secret",
							AADClientCertPath:            "aad_client_cert_path",
							AADClientCertPassword:        "aad_client_cert_password",
							CloudProviderBackoff:         true,
							CloudProviderBackoffRetries:  1,
							CloudProviderBackoffExponent: 2,
							CloudProviderBackoffDuration: 3,
							CloudProviderBackoffJitter:   4,
							CloudProviderRateLimit:       true,
							CloudProviderRateLimitQPS:    1,
							CloudProviderRateLimitBucket: 2,
							UseInstanceMetadata:          true,
							UseManagedIdentityExtension:  true,
							MaximumLoadBalancerRuleCount: 1,
						},
						VsphereCloudProvider: &v3.VsphereCloudProvider{
							Global: v3.GlobalVsphereOpts{
								User:              "user",
								Password:          "password",
								VCenterIP:         "server",
								VCenterPort:       "port",
								InsecureFlag:      true,
								Datacenter:        "datacenter",
								Datacenters:       "datacenters",
								DefaultDatastore:  "datastore",
								WorkingDir:        "working_dir",
								RoundTripperCount: 1,
								VMUUID:            "vm_uuid",
								VMName:            "vm_name",
							},
							VirtualCenter: map[string]v3.VirtualCenterConfig{
								"192.2.0.1": {
									User:              "user",
									Password:          "password",
									VCenterPort:       "port",
									Datacenters:       "datacenters",
									RoundTripperCount: 1,
								},
							},
							Network: v3.NetworkVshpereOpts{
								PublicNetwork: "public_network",
							},
							Disk: v3.DiskVsphereOpts{
								SCSIControllerType: "scsi_controller_type",
							},
							Workspace: v3.WorkspaceVsphereOpts{
								VCenterIP:        "server",
								Datacenter:       "datacenter",
								Folder:           "folder",
								DefaultDatastore: "default_datastore",
								ResourcePoolPath: "resourcepool_path",
							},
						},
						OpenstackCloudProvider: &v3.OpenstackCloudProvider{
							Global: v3.GlobalOpenstackOpts{
								AuthURL:    "auth_url",
								Username:   "username",
								UserID:     "user_id",
								Password:   "password",
								TenantID:   "tenant_id",
								TenantName: "tenant_name",
								TrustID:    "trust_id",
								DomainID:   "domain_id",
								DomainName: "domain_name",
								Region:     "region",
								CAFile:     "ca_file",
							},
							LoadBalancer: v3.LoadBalancerOpenstackOpts{
								LBVersion:            "lb_version",
								UseOctavia:           true,
								SubnetID:             "subnet_id",
								FloatingNetworkID:    "floating_network_id",
								LBMethod:             "lb_method",
								LBProvider:           "lb_provider",
								CreateMonitor:        true,
								MonitorDelay:         1,
								MonitorTimeout:       2,
								MonitorMaxRetries:    3,
								ManageSecurityGroups: true,
							},
							BlockStorage: v3.BlockStorageOpenstackOpts{
								BSVersion:       "bs_version",
								TrustDevicePath: true,
								IgnoreVolumeAZ:  true,
							},
							Route: v3.RouteOpenstackOpts{
								RouterID: "router_id",
							},
							Metadata: v3.MetadataOpenstackOpts{
								SearchOrder:    "search_order",
								RequestTimeout: 1,
							},
						},
						CustomCloudProvider: "custom_cloud_config",
					},
					PrefixPath: "prefix_path",
				},
				EtcdHosts: []*hosts.Host{
					{
						RKEConfigNode: v3.RKEConfigNode{
							NodeName: "etcd1",
							Address:  "192.2.0.1",
						},
					},
					{
						RKEConfigNode: v3.RKEConfigNode{
							NodeName: "etcd2",
							Address:  "192.2.0.2",
						},
					},
				},
				WorkerHosts: []*hosts.Host{
					{
						RKEConfigNode: v3.RKEConfigNode{
							NodeName: "host",
							Address:  "192.2.0.1",
						},
					},
				},
				ControlPlaneHosts: []*hosts.Host{
					{
						RKEConfigNode: v3.RKEConfigNode{
							NodeName: "host",
							Address:  "192.2.0.1",
						},
					},
				},
				InactiveHosts: []*hosts.Host{
					{
						RKEConfigNode: v3.RKEConfigNode{
							NodeName: "host",
							Address:  "192.2.0.1",
						},
					},
				},
				Certificates: map[string]pki.CertificatePKI{
					"example": {
						Certificate:   dummyCertificate,
						Key:           dummyPrivateKey,
						Config:        "config",
						Name:          "name",
						CommonName:    "common_name",
						OUName:        "ou_name",
						EnvName:       "env_name",
						Path:          "path",
						KeyEnvName:    "key_env_name",
						KeyPath:       "key_path",
						ConfigEnvName: "config_env_name",
						ConfigPath:    "config_path",
					},
				},
				ClusterDomain:    "example.com",
				ClusterCIDR:      "10.200.0.0/8",
				ClusterDNSServer: "192.2.0.1",
			},
			state: map[string]interface{}{
				"services_etcd": []interface{}{
					map[string]interface{}{
						"image": "etcd:latest",
						"extra_args": map[string]string{
							"foo": "bar",
							"bar": "foo",
						},
						"extra_binds": []string{"/bind1", "/bind2"},
						"extra_env":   []string{"env1=val1", "env2=val2"},
						"external_urls": []string{
							"https://ext1.example.com",
							"https://ext2.example.com",
						},
						"ca_cert":   "ca_cert",
						"cert":      "cert",
						"key":       "key",
						"path":      "path",
						"snapshot":  true,
						"retention": "retention",
						"creation":  "creation",
					},
				},
				"services_kube_api": []interface{}{
					map[string]interface{}{
						"image": "kube_api:latest",
						"extra_args": map[string]string{
							"foo": "bar",
							"bar": "foo",
						},
						"extra_binds":              []string{"/bind1", "/bind2"},
						"extra_env":                []string{"env1=val1", "env2=val2"},
						"service_cluster_ip_range": "10.240.0.0/16",
						"service_node_port_range":  "30000-31000",
						"pod_security_policy":      true,
					},
				},
				"services_kube_controller": []interface{}{
					map[string]interface{}{
						"image": "kube_controller:latest",
						"extra_args": map[string]string{
							"foo": "bar",
							"bar": "foo",
						},
						"extra_binds":              []string{"/bind1", "/bind2"},
						"extra_env":                []string{"env1=val1", "env2=val2"},
						"cluster_cidr":             "10.200.0.0/8",
						"service_cluster_ip_range": "10.240.0.0/16",
					},
				},
				"services_scheduler": []interface{}{
					map[string]interface{}{
						"image": "scheduler:latest",
						"extra_args": map[string]string{
							"foo": "bar",
							"bar": "foo",
						},
						"extra_binds": []string{"/bind1", "/bind2"},
						"extra_env":   []string{"env1=val1", "env2=val2"},
					},
				},
				"services_kubelet": []interface{}{
					map[string]interface{}{
						"image": "kubelet:latest",
						"extra_args": map[string]string{
							"foo": "bar",
							"bar": "foo",
						},
						"extra_binds":           []string{"/bind1", "/bind2"},
						"extra_env":             []string{"env1=val1", "env2=val2"},
						"cluster_domain":        "example.com",
						"infra_container_image": "alpine:latest",
						"cluster_dns_server":    "192.2.0.1",
						"fail_swap_on":          true,
					},
				},
				"services_kubeproxy": []interface{}{
					map[string]interface{}{
						"image": "kubeproxy:latest",
						"extra_args": map[string]string{
							"foo": "bar",
							"bar": "foo",
						},
						"extra_binds": []string{"/bind1", "/bind2"},
						"extra_env":   []string{"env1=val1", "env2=val2"},
					},
				},
				"network": []interface{}{
					map[string]interface{}{
						"plugin": "calico",
						"options": map[string]string{
							"foo": "bar",
							"bar": "foo",
						},
					},
				},
				"authentication": []interface{}{
					map[string]interface{}{
						"strategy": "x509",
						"options": map[string]string{
							"foo": "bar",
							"bar": "foo",
						},
						"sans": []string{"sans1", "sans2"},
					},
				},
				"addons": "addons: yaml",
				"addons_include": []string{
					"https://example.com/addon1.yaml",
					"https://example.com/addon2.yaml",
				},
				"addon_job_timeout": 10,
				"system_images": []interface{}{
					map[string]interface{}{
						"etcd":                        "etcd",
						"alpine":                      "alpine",
						"nginx_proxy":                 "nginx_proxy",
						"cert_downloader":             "cert_downloader",
						"kubernetes_services_sidecar": "kubernetes_services_sidecar",
						"kube_dns":                    "kube_dns",
						"dnsmasq":                     "dnsmasq",
						"kube_dns_sidecar":            "kube_dns_sidecar",
						"kube_dns_autoscaler":         "kube_dns_autoscaler",
						"kubernetes":                  "kubernetes",
						"flannel":                     "flannel",
						"flannel_cni":                 "flannel_cni",
						"calico_node":                 "calico_node",
						"calico_cni":                  "calico_cni",
						"calico_controllers":          "calico_controllers",
						"calico_ctl":                  "calico_ctl",
						"canal_node":                  "canal_node",
						"canal_cni":                   "canal_cni",
						"canal_flannel":               "canal_flannel",
						"weave_node":                  "weave_node",
						"weave_cni":                   "weave_cni",
						"pod_infra_container":         "pod_infra_container",
						"ingress":                     "ingress",
						"ingress_backend":             "ingress_backend",
					},
				},
				"ssh_key_path":   "ssh_key_path",
				"ssh_agent_auth": true,
				"bastion_host": []interface{}{
					map[string]interface{}{
						"address":        "192.2.0.1",
						"port":           22,
						"user":           "rancher",
						"ssh_agent_auth": true,
						"ssh_key":        "ssh_key",
						"ssh_key_path":   "ssh_key_path",
					},
				},
				"authorization": []interface{}{
					map[string]interface{}{
						"mode": "rbac",
						"options": map[string]string{
							"foo": "bar",
							"bar": "foo",
						},
					},
				},
				"ignore_docker_version": true,
				"kubernetes_version":    "1.8.9",
				"private_registries": []interface{}{
					map[string]interface{}{
						"url":      "https://registry1.example.com",
						"user":     "user1",
						"password": "password1",
					},
					map[string]interface{}{
						"url":      "https://registry2.example.com",
						"user":     "user2",
						"password": "password2",
					},
				},
				"ingress": []interface{}{
					map[string]interface{}{
						"provider": "nginx",
						"options": map[string]string{
							"foo": "bar",
							"bar": "foo",
						},
						"node_selector": map[string]string{
							"role": "worker",
						},
						"extra_args": map[string]string{
							"foo": "foo",
							"bar": "bar",
						},
					},
				},
				"cluster_name":    "example",
				"kube_admin_user": "kube-admin",
				"api_server_url":  "https://192.2.0.1:6443",
				"cloud_provider": []interface{}{
					map[string]interface{}{
						"name": "sakuracloud",
						"azure_cloud_config": []interface{}{
							map[string]interface{}{
								"cloud":                            "cloud",
								"tenant_id":                        "tenant_id",
								"subscription_id":                  "subscription_id",
								"resource_group":                   "resource_group",
								"location":                         "location",
								"vnet_name":                        "vnet_name",
								"vnet_resource_group":              "vnet_resource_group",
								"route_table_name":                 "route_table_name",
								"primary_availability_set_name":    "primary_availability_set_name",
								"vm_type":                          "vm_type",
								"primary_scale_set_name":           "primary_scale_set_name",
								"aad_client_id":                    "aad_client_id",
								"aad_client_secret":                "aad_client_secret",
								"aad_client_cert_path":             "aad_client_cert_path",
								"aad_client_cert_password":         "aad_client_cert_password",
								"cloud_provider_backoff":           true,
								"cloud_provider_backoff_retries":   1,
								"cloud_provider_backoff_exponent":  2,
								"cloud_provider_backoff_duration":  3,
								"cloud_provider_backoff_jitter":    4,
								"cloud_provider_rate_limit":        true,
								"cloud_provider_rate_limit_qps":    1,
								"cloud_provider_rate_limit_bucket": 2,
								"use_instance_metadata":            true,
								"use_managed_identity_extension":   true,
								"maximum_load_balancer_rule_count": 1,
							},
						},
						"vsphere_cloud_config": []interface{}{
							map[string]interface{}{
								"global": []interface{}{
									map[string]interface{}{
										"user":                 "user",
										"password":             "password",
										"server":               "server",
										"port":                 "port",
										"insecure_flag":        true,
										"datacenter":           "datacenter",
										"datacenters":          "datacenters",
										"datastore":            "datastore",
										"working_dir":          "working_dir",
										"soap_roundtrip_count": 1,
										"vm_uuid":              "vm_uuid",
										"vm_name":              "vm_name",
									},
								},
								"virtual_center": []interface{}{
									map[string]interface{}{
										"server":               "192.2.0.1",
										"user":                 "user",
										"password":             "password",
										"port":                 "port",
										"datacenters":          "datacenters",
										"soap_roundtrip_count": 1,
									},
								},
								"network": []interface{}{
									map[string]interface{}{
										"public_network": "public_network",
									},
								},
								"disk": []interface{}{
									map[string]interface{}{
										"scsi_controller_type": "scsi_controller_type",
									},
								},
								"workspace": []interface{}{
									map[string]interface{}{
										"server":            "server",
										"datacenter":        "datacenter",
										"folder":            "folder",
										"default_datastore": "default_datastore",
										"resourcepool_path": "resourcepool_path",
									},
								},
							},
						},
						"openstack_cloud_config": []interface{}{
							map[string]interface{}{
								"global": []interface{}{
									map[string]interface{}{
										"auth_url":    "auth_url",
										"username":    "username",
										"user_id":     "user_id",
										"password":    "password",
										"tenant_id":   "tenant_id",
										"tenant_name": "tenant_name",
										"trust_id":    "trust_id",
										"domain_id":   "domain_id",
										"domain_name": "domain_name",
										"region":      "region",
										"ca_file":     "ca_file",
									},
								},
								"load_balancer": []interface{}{
									map[string]interface{}{
										"lb_version":             "lb_version",
										"use_octavia":            true,
										"subnet_id":              "subnet_id",
										"floating_network_id":    "floating_network_id",
										"lb_method":              "lb_method",
										"lb_provider":            "lb_provider",
										"create_monitor":         true,
										"monitor_delay":          1,
										"monitor_timeout":        2,
										"monitor_max_retries":    3,
										"manage_security_groups": true,
									},
								},
								"block_storage": []interface{}{
									map[string]interface{}{
										"bs_version":        "bs_version",
										"trust_device_path": true,
										"ignore_volume_az":  true,
									},
								},
								"route": []interface{}{
									map[string]interface{}{
										"router_id": "router_id",
									},
								},
								"metadata": []interface{}{
									map[string]interface{}{
										"search_order":    "search_order",
										"request_timeout": 1,
									},
								},
							},
						},
						"custom_cloud_config": "custom_cloud_config",
					},
				},
				"prefix_path": "prefix_path",
				"certificates": []interface{}{
					map[string]interface{}{
						"id":              "example",
						"certificate":     dummyCertPEM,
						"key":             dummyPrivateKeyPEM,
						"config":          "config",
						"name":            "name",
						"common_name":     "common_name",
						"ou_name":         "ou_name",
						"env_name":        "env_name",
						"path":            "path",
						"key_env_name":    "key_env_name",
						"key_path":        "key_path",
						"config_env_name": "config_env_name",
						"config_path":     "config_path",
					},
				},
				"cluster_domain":     "example.com",
				"cluster_cidr":       "10.200.0.0/8",
				"cluster_dns_server": "192.2.0.1",
				"etcd_hosts": []map[string]interface{}{
					{
						"node_name": "etcd1",
						"address":   "192.2.0.1",
					},
					{
						"node_name": "etcd2",
						"address":   "192.2.0.2",
					},
				},
				"worker_hosts": []map[string]interface{}{
					{
						"node_name": "host",
						"address":   "192.2.0.1",
					},
				},
				"control_plane_hosts": []map[string]interface{}{
					{
						"node_name": "host",
						"address":   "192.2.0.1",
					},
				},
				"inactive_hosts": []map[string]interface{}{
					{
						"node_name": "host",
						"address":   "192.2.0.1",
					},
				},
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.caseName, func(t *testing.T) {
			d := &dummyStateBuilder{values: map[string]interface{}{}}
			err := clusterToState(testcase.cluster, d)
			assert.NoError(t, err)
			assert.EqualValues(t, testcase.state, d.values)
		})
	}

}
