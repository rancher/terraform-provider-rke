package rke

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-getter/helper/url"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/rancher/rke/cluster"
	"github.com/rancher/rke/hosts"
	"github.com/rancher/rke/k8s"
	"github.com/rancher/rke/log"
	"github.com/rancher/rke/pki"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/util/cert"
)

const rkeKubeConfigFileName = "kube_config_cluster.yml"

func resourceRKECluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceRKEClusterUp,
		Read:   resourceRKEClusterRead,
		Update: resourceRKEClusterUp,
		Delete: resourceRKEClusterDelete,
		Schema: map[string]*schema.Schema{
			"disable_port_check": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"nodes_conf": {
				Type:          schema.TypeList,
				MinItems:      1,
				Optional:      true,
				Description:   "Kubernetes nodes(YAML or JSON)",
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"nodes"},
			},
			"nodes": {
				Type:        schema.TypeList,
				MinItems:    1,
				Optional:    true,
				Description: "Kubernetes nodes",
				Elem: &schema.Resource{
					Schema: nodeSchema(),
				},
				ConflictsWith: []string{"nodes_conf"},
			},
			"services_etcd": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Docker image of the etcd service",
						},
						"extra_args": {
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Description: "Extra arguments that are added to the etcd services",
						},
						"extra_binds": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Computed:    true,
							Description: "Extra binds added to the nodes",
						},
						"extra_env": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Computed:    true,
							Description: "Extra env is to provide extra env variable to the docker container running kubernetes service",
						},
						"external_urls": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Computed:    true,
							Description: "List of etcd urls",
						},
						"ca_cert": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "External CA certificate",
						},
						"cert": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "External Client certificate",
						},
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "External Client key",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "External etcd prefix",
						},
						"snapshot": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Etcd Recurring snapshot Service",
						},
						"retention": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Etcd snapshot Retention period",
						},
						"creation": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Etcd snapshot Creation period",
						},
					},
				},
			},
			"services_kube_api": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Docker image of the kube-api service",
						},
						"extra_args": {
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Description: "Extra arguments that are added to the kube-api services",
						},
						"extra_binds": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Computed:    true,
							Description: "Extra binds added to the nodes",
						},
						"extra_env": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Computed:    true,
							Description: "Extra env is to provide extra env variable to the docker container running kubernetes service",
						},
						"service_cluster_ip_range": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Virtual IP range that will be used by Kubernetes services",
						},
						"service_node_port_range": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Port range for services defined with NodePort type",
						},
						"pod_security_policy": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Enabled/Disable PodSecurityPolicy",
						},
					},
				},
			},
			"services_kube_controller": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Docker image of the kube-controller service",
						},
						"extra_args": {
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Description: "Extra arguments that are added to the kube-controller services",
						},
						"extra_binds": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Computed:    true,
							Description: "Extra binds added to the nodes",
						},
						"extra_env": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Computed:    true,
							Description: "Extra env is to provide extra env variable to the docker container running kubernetes service",
						},
						"cluster_cidr": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "CIDR Range for Pods in cluster",
						},
						"service_cluster_ip_range": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Virtual IP range that will be used by Kubernetes services",
						},
					},
				},
			},
			"services_scheduler": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Docker image of the scheduler service",
						},
						"extra_args": {
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Description: "Extra arguments that are added to the scheduler services",
						},
						"extra_binds": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Optional:    true,
							Description: "Extra binds added to the nodes",
						},
						"extra_env": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Optional:    true,
							Description: "Extra env is to provide extra env variable to the docker container running kubernetes service",
						},
					},
				},
			},
			"services_kubelet": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Docker image of the kubelet service",
						},
						"extra_args": {
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Description: "Extra arguments that are added to the kubelet services",
						},
						"extra_binds": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Computed:    true,
							Description: "Extra binds added to the nodes",
						},
						"extra_env": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Computed:    true,
							Description: "Extra env is to provide extra env variable to the docker container running kubernetes service",
						},
						"cluster_domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Domain of the cluster (default: "cluster.local")`,
						},
						"infra_container_image": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The image whose network/ipc namespaces containers in each pod will use",
						},
						"cluster_dns_server": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Cluster DNS service ip",
						},
						"fail_swap_on": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Fail if swap is enabled",
						},
					},
				},
			},
			"services_kubeproxy": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"image": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Docker image of the kubeproxy service",
						},
						"extra_args": {
							Type:        schema.TypeMap,
							Optional:    true,
							Computed:    true,
							Description: "Extra arguments that are added to the kubeproxy services",
						},
						"extra_binds": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Computed:    true,
							Description: "Extra binds added to the nodes",
						},
						"extra_env": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Computed:    true,
							Description: "Extra env is to provide extra env variable to the docker container running kubernetes service",
						},
					},
				},
			},
			"network": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Network configuration used in the kubernetes cluster",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plugin": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							Description:  "Network Plugin That will be used in kubernetes cluster",
							ValidateFunc: validation.StringInSlice([]string{"flannel", "calico", "canal", "weave"}, false),
						},
						"options": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Plugin options to configure network properties",
						},
					},
				},
			},
			"authentication": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Authentication configuration used in the cluster",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"strategy": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							Description:  "Authentication strategy that will be used in kubernetes cluster",
							ValidateFunc: validation.StringInSlice([]string{"x509"}, false),
						},
						"options": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Authentication options",
						},
						"sans": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "List of additional hostnames and IPs to include in the api server PKI cert",
						},
					},
				},
			},
			"addons": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "YAML manifest for user provided addons to be deployed on the cluster",
			},
			"addons_include": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "List of urls or paths for addons",
			},
			"addon_job_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
				Description:  "Timeout in seconds for status check on addon deployment jobs",
			},
			"system_images": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "List of images used internally for proxy, cert download ,kubedns and more",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"etcd": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"alpine": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"nginx_proxy": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cert_downloader": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kubernetes_services_sidecar": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kube_dns": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"dnsmasq": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kube_dns_sidecar": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kube_dns_autoscaler": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"kubernetes": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"flannel": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"flannel_cni": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"calico_node": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"calico_cni": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"calico_controllers": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"calico_ctl": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"canal_node": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"canal_cni": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"canal_flannel": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"weave_node": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"weave_cni": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"pod_infra_container": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ingress": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ingress_backend": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"metrics_server": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"ssh_key_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "SSH Private Key Path",
			},
			"ssh_agent_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "SSH Agent Auth enable",
			},
			"bastion_host": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Bastion/Jump Host configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Address of Bastion Host",
						},
						"port": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
							Description:  "SSH Port of Bastion Host",
						},
						"user": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SSH User to Bastion Host",
						},
						"ssh_agent_auth": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "SSH Agent Auth enable",
						},
						"ssh_key": {
							Type:        schema.TypeString,
							Sensitive:   true,
							Optional:    true,
							Description: "SSH Private Key",
						},
						"ssh_key_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "SSH Private Key",
						},
					},
				},
			},
			"monitoring": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Monitoring Config",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Monitoring server provider",
						},
						"options": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Metrics server options",
						},
					},
				},
			},
			"authorization": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Authorization mode configuration used in the cluster",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							Description:  "Authorization mode used by kubernetes",
							ValidateFunc: validation.StringInSlice([]string{"rbac", "none"}, false),
						},
						"options": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Authorization mode options",
						},
					},
				},
			},
			"ignore_docker_version": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable/Disable strict docker version checking",
			},
			"kubernetes_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(func() []string {
					keys := make([]string, 0, len(v3.K8sVersionToRKESystemImages))
					for k := range v3.K8sVersionToRKESystemImages {
						keys = append(keys, k)
					}
					return keys
				}(), false),
				Description: "Kubernetes version to use (if kubernetes image is specified, image version takes precedence)",
			},
			"private_registries": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of private registries and their credentials",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "URL for the registry",
						},
						"user": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User name for registry access",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "Password for registry access",
						},
						"is_default": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Default registry",
						},
					},
				},
			},
			"ingress": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Ingress controller used in the cluster",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"provider": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Ingress controller type used by kubernetes",
						},
						"options": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Ingress controller options",
						},
						"node_selector": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Ingress controller used in the cluster",
						},
						"extra_args": {
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "Ingress controller extra arguments",
						},
					},
				},
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Cluster Name used in the kube config",
			},
			"cloud_provider": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Cloud Provider options",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the Cloud Provider",
						},
						"aws_cloud_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "AWS cloud config file",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{},
							},
						},
						"azure_cloud_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Azure cloud config file",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cloud": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The cloud environment identifier. Takes values from https://github.com/Azure/go-autorest/blob/ec5f4903f77ed9927ac95b19ab8e44ada64c1356/autorest/azure/environments.go#L13",
									},
									"tenant_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The AAD Tenant ID for the Subscription that the cluster is deployed in",
									},
									"subscription_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The ID of the Azure Subscription that the cluster is deployed in",
									},
									"resource_group": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name of the resource group that the cluster is deployed in",
									},
									"location": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The location of the resource group that the cluster is deployed in",
									},
									"vnet_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name of the VNet that the cluster is deployed in",
									},
									"vnet_resource_group": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name of the resource group that the Vnet is deployed in",
									},
									"route_table_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "(Optional in 1.6) The name of the route table attached to the subnet that the cluster is deployed in",
									},
									"primary_availability_set_name": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "The name of the availability set that should be used as the load balancer backend" +
											"If this is set, the Azure cloudprovider will only add nodes from that availability set to the load" +
											"balancer backend pool. If this is not set, and multiple agent pools (availability sets) are used, then" +
											"the cloudprovider will try to add all nodes to a single backend pool which is forbidden." +
											"In other words, if you use multiple agent pools (availability sets), you MUST set this field.",
									},
									"vm_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The type of azure nodes. Candidate valudes are: vmss and standard. If not set, it will be default to standard.",
									},
									"primary_scale_set_name": {
										Type:     schema.TypeString,
										Optional: true,
										Description: "The name of the scale set that should be used as the load balancer backend." +
											"If this is set, the Azure cloudprovider will only add nodes from that scale set to the load" +
											"balancer backend pool. If this is not set, and multiple agent pools (scale sets) are used, then" +
											"the cloudprovider will try to add all nodes to a single backend pool which is forbidden." +
											"In other words, if you use multiple agent pools (scale sets), you MUST set this field.",
									},
									"aad_client_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The ClientID for an AAD application with RBAC access to talk to Azure RM APIs",
									},
									"aad_client_secret": {
										Type:        schema.TypeString,
										Optional:    true,
										Sensitive:   true,
										Description: "The ClientSecret for an AAD application with RBAC access to talk to Azure RM APIs",
									},
									"aad_client_cert_path": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The path of a client certificate for an AAD application with RBAC access to talk to Azure RM APIs",
									},
									"aad_client_cert_password": {
										Type:        schema.TypeString,
										Optional:    true,
										Sensitive:   true,
										Description: "The password of the client certificate for an AAD application with RBAC access to talk to Azure RM APIs",
									},
									"cloud_provider_backoff": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Enable exponential backoff to manage resource request retries",
									},
									"cloud_provider_backoff_retries": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Backoff retry limit",
									},
									"cloud_provider_backoff_exponent": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Backoff exponent",
									},
									"cloud_provider_backoff_duration": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Backoff duration",
									},
									"cloud_provider_backoff_jitter": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Backoff jitter",
									},
									"cloud_provider_rate_limit": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Enable rate limiting",
									},
									"cloud_provider_rate_limit_qps": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Rate limit QPS",
									},
									"cloud_provider_rate_limit_bucket": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Rate limit Bucket Size",
									},
									"use_instance_metadata": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Use instance metadata service where possible",
									},
									"use_managed_identity_extension": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Use managed service identity for the virtual machine to access Azure ARM APIs",
									},
									"maximum_load_balancer_rule_count": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Maximum allowed LoadBalancer Rule Count is the limit enforced by Azure Load balancer",
									},
								},
							},
						},
						"vsphere_cloud_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Vsphere cloud config file",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"global": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"user": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"password": {
													Type:      schema.TypeString,
													Optional:  true,
													Sensitive: true,
												},
												"server": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"port": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"insecure_flag": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"datacenter": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"datacenters": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"datastore": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"working_dir": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"soap_roundtrip_count": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"vm_uuid": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"vm_name": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"virtual_center": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"server": {
													Type:     schema.TypeString,
													Required: true,
												},
												"user": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"password": {
													Type:      schema.TypeString,
													Optional:  true,
													Sensitive: true,
												},
												"port": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"datacenters": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"soap_roundtrip_count": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
									"network": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"public_network": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"disk": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"scsi_controller_type": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"workspace": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"server": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"datacenter": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"folder": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"default_datastore": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"resourcepool_path": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"openstack_cloud_config": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "OpenStack cloud config file",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"global": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"auth_url": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"username": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"user_id": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"password": {
													Type:      schema.TypeString,
													Optional:  true,
													Sensitive: true,
												},
												"tenant_id": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"tenant_name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"trust_id": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"domain_id": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"domain_name": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"region": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"ca_file": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"load_balancer": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"lb_version": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"use_octavia": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"subnet_id": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"floating_network_id": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"lb_method": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"lb_provider": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"create_monitor": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"monitor_delay": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"monitor_timeout": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"monitor_max_retries": {
													Type:     schema.TypeInt,
													Optional: true,
												},
												"manage_security_groups": {
													Type:     schema.TypeBool,
													Optional: true,
												},
											},
										},
									},
									"block_storage": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bs_version": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"trust_device_path": {
													Type:     schema.TypeBool,
													Optional: true,
												},
												"ignore_volume_az": {
													Type:     schema.TypeBool,
													Optional: true,
												},
											},
										},
									},
									"route": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"router_id": {
													Type:     schema.TypeString,
													Optional: true,
												},
											},
										},
									},
									"metadata": {
										Type:     schema.TypeList,
										MaxItems: 1,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"search_order": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"request_timeout": {
													Type:     schema.TypeInt,
													Optional: true,
												},
											},
										},
									},
								},
							},
						},
						"custom_cloud_config": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CustomCloudProvider is a multiline string that represent a custom cloud config file",
						},
					},
				},
			},
			"prefix_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Kubernetes directory path",
			},
			"ca_crt": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"client_cert": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"client_key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"kube_config_yaml": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rke_cluster_yaml": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key": {
							Type:      schema.TypeString,
							Sensitive: true,
							Computed:  true,
						},
						"config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"common_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ou_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"env_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_env_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"key_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config_env_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"kube_admin_user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_server_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_cidr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_dns_server": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"etcd_hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"worker_hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"control_plane_hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"inactive_hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceRKEClusterUp(d *schema.ResourceData, meta interface{}) error {
	if err := clusterUp(d); err != nil {
		return wrapErrWithRKEOutputs(err)
	}
	return wrapErrWithRKEOutputs(resourceRKEClusterRead(d, meta))
}

func resourceRKEClusterRead(d *schema.ResourceData, meta interface{}) error {
	currentCluster, err := readClusterState(d)
	if err != nil {
		if _, ok := err.(*nodeUnreachableError); ok {
			d.SetId("")
			return nil
		}
		return wrapErrWithRKEOutputs(err)
	}
	return wrapErrWithRKEOutputs(clusterToState(currentCluster, d))
}

func resourceRKEClusterDelete(d *schema.ResourceData, meta interface{}) error {
	if err := clusterRemove(d); err != nil {
		if _, ok := err.(*nodeUnreachableError); !ok {
			return wrapErrWithRKEOutputs(err)
		}
	}
	d.SetId("")
	return nil
}

func clusterUp(d *schema.ResourceData) error {
	rkeConfig, parseErr := parseResourceRKEConfig(d)
	if parseErr != nil {
		return parseErr
	}

	disablePortCheck := d.Get("disable_port_check").(bool)

	// create tmp dir for configDir
	tempDir, tempDirErr := createTempDir()
	if tempDirErr != nil {
		return tempDirErr
	}
	defer os.RemoveAll(tempDir) // nolint

	// deploy
	clusterFilePath := filepath.Join(tempDir, "cluster.yml")
	apiURL, caCrt, clientCert, clientKey, clusterUpErr := realClusterUp(context.Background(),
		rkeConfig, nil, nil, nil,
		clusterFilePath, "", false, disablePortCheck)
	if clusterUpErr != nil {
		return clusterUpErr
	}

	// set keys to resourceData
	return setRKEClusterKeys(d, apiURL, caCrt, clientCert, clientKey, tempDir, rkeConfig)
}

func clusterRemove(d *schema.ResourceData) error {
	rkeConfig, parseErr := parseResourceRKEConfig(d)
	if parseErr != nil {
		return parseErr
	}

	// create tmp dir for configDir
	tempDir, tempDirErr := createTempDir()
	if tempDirErr != nil {
		return tempDirErr
	}
	defer os.RemoveAll(tempDir) // nolint
	if err := writeKubeConfigFile(tempDir, d); err != nil {
		return err
	}
	clusterFilePath := filepath.Join(tempDir, "cluster.yml")
	if err := writeClusterConfig(rkeConfig, clusterFilePath); err != nil {
		return err
	}

	return realClusterRemove(context.Background(),
		rkeConfig, nil, nil, clusterFilePath, "")
}

func realClusterUp( // nolint: gocyclo
	ctx context.Context,
	rkeConfig *v3.RancherKubernetesEngineConfig,
	dockerDialerFactory, localConnDialerFactory hosts.DialerFactory,
	k8sWrapTransport k8s.WrapTransport,
	clusterFilePath, configDir string, updateOnly, disablePortCheck bool) (string, string, string, string, error) {

	log.Infof(ctx, "Building Kubernetes cluster")
	var APIURL, caCrt, clientCert, clientKey string
	kubeCluster, err := cluster.ParseCluster(ctx, rkeConfig, clusterFilePath, configDir,
		dockerDialerFactory, localConnDialerFactory, k8sWrapTransport)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, err
	}

	err = kubeCluster.TunnelHosts(ctx, false)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, newNodeUnreachableError(err)
	}

	currentCluster, err := kubeCluster.GetClusterState(ctx)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, err
	}
	if !disablePortCheck {
		if err = kubeCluster.CheckClusterPorts(ctx, currentCluster); err != nil {
			return APIURL, caCrt, clientCert, clientKey, err
		}
	}

	err = cluster.SetUpAuthentication(ctx, kubeCluster, currentCluster)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, err
	}

	err = cluster.ReconcileCluster(ctx, kubeCluster, currentCluster, updateOnly)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, err
	}

	err = kubeCluster.SetUpHosts(ctx)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, err
	}

	if err = kubeCluster.PrePullK8sImages(ctx); err != nil {
		return APIURL, caCrt, clientCert, clientKey, err
	}

	err = kubeCluster.DeployControlPlane(ctx)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, err
	}

	// Apply Authz configuration after deploying controlplane
	err = cluster.ApplyAuthzResources(ctx, kubeCluster.RancherKubernetesEngineConfig, clusterFilePath, configDir, k8sWrapTransport)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, err
	}

	err = kubeCluster.SaveClusterState(ctx, rkeConfig)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, err
	}

	err = kubeCluster.DeployWorkerPlane(ctx)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, err
	}

	if err = kubeCluster.CleanDeadLogs(ctx); err != nil {
		return APIURL, caCrt, clientCert, clientKey, err
	}

	err = kubeCluster.SyncLabelsAndTaints(ctx, currentCluster)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, err
	}

	err = cluster.ConfigureCluster(ctx, kubeCluster.RancherKubernetesEngineConfig, kubeCluster.Certificates, clusterFilePath, configDir, k8sWrapTransport, false)
	if err != nil {
		return APIURL, caCrt, clientCert, clientKey, err
	}
	if len(kubeCluster.ControlPlaneHosts) > 0 {
		APIURL = fmt.Sprintf("https://" + kubeCluster.ControlPlaneHosts[0].Address + ":6443")
		clientCert = string(cert.EncodeCertPEM(kubeCluster.Certificates[pki.KubeAdminCertName].Certificate))
		clientKey = string(cert.EncodePrivateKeyPEM(kubeCluster.Certificates[pki.KubeAdminCertName].Key))
	}
	caCrt = string(cert.EncodeCertPEM(kubeCluster.Certificates[pki.CACertName].Certificate))

	if err := checkAllIncluded(kubeCluster); err != nil {
		return APIURL, caCrt, clientCert, clientKey, err
	}

	log.Infof(ctx, "Finished building Kubernetes cluster successfully")
	return APIURL, caCrt, clientCert, clientKey, nil
}

func checkAllIncluded(cluster *cluster.Cluster) error {
	if len(cluster.InactiveHosts) == 0 {
		return nil
	}

	var names []string
	for _, host := range cluster.InactiveHosts {
		names = append(names, host.Address)
	}

	return fmt.Errorf("Provisioning incomplete, host(s) [%s] skipped because they could not be contacted", strings.Join(names, ","))
}

func realClusterRemove(
	ctx context.Context,
	rkeConfig *v3.RancherKubernetesEngineConfig,
	dialerFactory hosts.DialerFactory,
	k8sWrapTransport k8s.WrapTransport,
	clusterFilePath, configDir string) error {

	log.Infof(ctx, "Tearing down Kubernetes cluster")
	kubeCluster, err := cluster.ParseCluster(ctx, rkeConfig, clusterFilePath, configDir, dialerFactory, nil, k8sWrapTransport)
	if err != nil {
		return err
	}

	err = kubeCluster.TunnelHosts(ctx, false)
	if err != nil {
		return newNodeUnreachableError(err)
	}

	log.Infof(ctx, "Starting Cluster removal")
	err = kubeCluster.ClusterRemove(ctx)
	if err != nil {
		return err
	}

	log.Infof(ctx, "Cluster removed successfully")
	return nil
}

func setRKEClusterKeys(d *schema.ResourceData, apiURL, caCrt, clientCert, clientKey string, configDir string, rkeConfig *v3.RancherKubernetesEngineConfig) error {

	parsedURL, err := url.Parse(apiURL)
	if err != nil {
		return err
	}
	d.Set("ca_crt", caCrt)           // nolint
	d.Set("client_cert", clientCert) // nolint
	d.Set("client_key", clientKey)   // nolint

	kubeConfig, err := readKubeConfig(configDir)
	if err != nil {
		return err
	}
	if kubeConfig != "" {
		d.Set("kube_config_yaml", kubeConfig) // nolint
	}

	yamlRkeConfig, err := yaml.Marshal(*rkeConfig)
	if err != nil {
		return err
	}
	d.Set("rke_cluster_yaml", string(yamlRkeConfig)) // nolint

	d.SetId(parsedURL.Hostname())
	return nil
}

func readClusterState(d *schema.ResourceData) (*cluster.Cluster, error) {
	apiURL := fmt.Sprintf("https://%s:6443", d.Id())
	caCrt := d.Get("ca_crt").(string)
	clientCert := d.Get("client_cert").(string)
	clientKey := d.Get("client_key").(string)

	requiredValues := []string{apiURL, caCrt, clientCert, clientKey}
	for _, v := range requiredValues {
		if v == "" {
			d.SetId("")
			return nil, nil
		}
	}

	rkeConfig, parseErr := parseResourceRKEConfig(d)
	if parseErr != nil {
		return nil, parseErr
	}

	// create tmp dir for cluster.yml and kube_config_cluster.yml
	tempDir, tempDirErr := createTempDir()
	if tempDirErr != nil {
		return nil, tempDirErr
	}
	defer os.RemoveAll(tempDir) // nolint
	if err := writeKubeConfigFile(tempDir, d); err != nil {
		return nil, err
	}

	clusterFilePath := filepath.Join(tempDir, "cluster.yml")
	if err := writeClusterConfig(rkeConfig, clusterFilePath); err != nil {
		return nil, err
	}

	ctx := context.Background()
	kubeCluster, err := cluster.ParseCluster(ctx, rkeConfig, clusterFilePath,
		"", nil, nil, nil)
	if err != nil {
		return nil, err
	}

	err = kubeCluster.TunnelHosts(ctx, false)
	if err != nil {
		return nil, newNodeUnreachableError(err)
	}

	return kubeCluster.GetClusterState(ctx)
}

func readKubeConfig(dir string) (string, error) {
	kubeConfigPath := filepath.Join(dir, rkeKubeConfigFileName)
	if _, err := os.Stat(kubeConfigPath); err == nil {
		var data []byte
		if data, err = ioutil.ReadFile(kubeConfigPath); err != nil {
			return "", err
		}
		return string(data), nil
	}
	return "", nil
}

func writeKubeConfigFile(dir string, d *schema.ResourceData) error {
	if rawKubeConfig, ok := d.GetOk("kube_config_yaml"); ok {
		strConf := rawKubeConfig.(string)
		if strConf != "" {
			kubeConfigPath := filepath.Join(dir, rkeKubeConfigFileName)
			if err := ioutil.WriteFile(kubeConfigPath, []byte(strConf), 0640); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeClusterConfig(cluster *v3.RancherKubernetesEngineConfig, configFile string) error {
	yamlConfig, err := yaml.Marshal(*cluster)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(configFile, []byte(string(yamlConfig)), 0640)
}

func createTempDir() (string, error) {
	// create tmp dir for configDir
	var workDir, tempDir string
	var err error
	if workDir, err = os.Getwd(); err != nil {
		return "", err
	}
	if tempDir, err = ioutil.TempDir(workDir, "terraform-provider-rke-"); err != nil {
		return "", err
	}
	return tempDir, nil
}

type nodeUnreachableError struct {
	actual error
}

func newNodeUnreachableError(actual error) *nodeUnreachableError {
	return &nodeUnreachableError{actual: actual}
}

func (n *nodeUnreachableError) Error() string {
	return n.actual.Error()
}
