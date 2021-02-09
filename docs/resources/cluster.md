---
page_title: "rke_cluster Resource"
---

# rke\_cluster

Provides RKE cluster resource. This can be used to create RKE clusters and retrieve their information.

RKE clusters can be defined in the provider:
- Using cluster_yaml: The full RKE cluster is defined in an RKE cluster.yml file.
- Using the TF provider arguments to define the entire cluster.
- Using a combination of both the cluster_yaml and TF provider arguments. The TF arguments will override the cluster_yaml options if collisions occur.

## Example Usage

Creating RKE cluster

```hcl
# Configure RKE provider
provider "rke" {
  log_file = "rke_debug.log"
}
# Create a new RKE cluster using config yaml
resource "rke_cluster" "foo" {
  cluster_yaml = file("cluster.yaml")
}
# Create a new RKE cluster using arguments
resource "rke_cluster" "foo2" {
  nodes {
    address = "1.2.3.4"
    user    = "ubuntu"
    role    = ["controlplane", "worker", "etcd"]
    ssh_key = file("~/.ssh/id_rsa")
  }
  upgrade_strategy {
	  drain = true
	  max_unavailable_worker = "20%"
  }
}
# Create a new RKE cluster using both. In case of conflict, arguments override cluster_yaml arguments
resource "rke_cluster" "foo2" {
  cluster_yaml = file("cluster.yaml")
  ssh_agent_auth = true
  ignore_docker_version = true
  kubernetes_version = "<K8s_VERSION>"
  upgrade_strategy {
	  drain = true
	  max_unavailable_worker = "20%"
  }
}
```

Restore RKE cluster. RKE cluster must be already managed by terraform and etcd snapshot must exist

```hcl
resource "rke_cluster" "cluster" {
  cluster_name = "foo"
  nodes {
    address = "1.2.3.4"
    user    = "ubuntu"
    role    = ["controlplane", "worker", "etcd"]
    ssh_key = file("~/.ssh/id_rsa")
  }
  restore {
    restore = true
    snapshot_name = "test.db"
  }
  upgrade_strategy {
    drain = true
    max_unavailable_worker = "20%"
  }
}
```

**Note** Once the RKE cluster is restored, `rke_cluster.restore.restore` will be set to `false` to force tf diff on next apply until user set `rke_cluster.restore.restore = false` on tf file

## Argument Reference

The following arguments are supported:

* `delay_on_creation` - (Optional) RKE k8s cluster delay on creation (int)
* `disable_port_check` - (Optional) Enable/Disable RKE k8s cluster port checking. Default `false` (bool)
* `addon_job_timeout` - (Optional) RKE k8s cluster addon deployment timeout in seconds for status check (int)
* `addons` - (Optional) RKE k8s cluster user addons YAML manifest to be deployed (string)
* `addons_include` - (Optional) RKE k8s cluster user addons YAML manifest urls or paths to be deployed (list)
* `authentication` - (Optional) RKE k8s cluster authentication configuration (list maxitems:1)
* `authorization` - (Optional) RKE k8s cluster authorization mode configuration (list maxitems:1)
* `bastion_host` - (Optional) RKE k8s cluster bastion Host configuration (list maxitems:1)
* `cert_dir` - (Optional) Specify a certificate dir path (string)
* `cloud_provider` - (Optional) RKE k8s cluster cloud provider configuration [rke-cloud-providers](https://rancher.com/docs/rke/latest/en/config-options/cloud-providers/) (list maxitems:1)
* `cluster_name` - (Optional) RKE k8s cluster name used in the kube config (string)
* `cluster_yaml` - (Optional) RKE k8s cluster config yaml encoded. Provider arguments take precedence over this one (string)
* `custom_certs` - (Optional) Use custom certificates from a cert dir (string)
* `dind` - (Optional/Experimental) Deploy RKE cluster on a dind environment. Default: `false` (bool)
* `dind_storage_driver` - (Optional/Experimental) DinD RKE cluster storage driver (string)
* `dind_dns_server` - (Optional/Experimental) DinD RKE cluster dns (string)
* `dns` - (Optional) RKE k8s cluster DNS Config (list maxitems:1)
* `ignore_docker_version` - (Optional) Enable/Disable RKE k8s cluster strict docker version checking. Default `false` (bool)
* `ingress` - (Optional) RKE k8s cluster ingress controller configuration (list maxitems:1)
* `kubernetes_version` - (Optional) K8s version to deploy. If kubernetes image is specified, image version takes precedence. Default: `rke default` (string)
* `monitoring` - (Optional) RKE k8s cluster monitoring Config (list maxitems:1)
* `network` - (Optional) RKE k8s cluster network configuration (list maxitems:1)
* `nodes` - (Optional) RKE k8s cluster nodes (list)
* `prefix_path` - (Optional/Computed) RKE k8s directory path (string)
* `private_registries` - (Optional/Computed) RKE k8s cluster private docker registries (list)
* `restore` - (Optional/Computed) RKE k8s cluster restore configuration (list maxitems:1)
* `rotate_certificates` - (Optional) RKE k8s cluster rotate certificates configuration (list maxitems:1)
* `services` - (Optional) RKE k8s cluster services (list maxitems:1)
* `services_etcd` - (DEPRECATED) Use services.etcd instead (list maxitems:1)
* `services_kube_api` - (DEPRECATED) Use services.kube_api instead (list maxitems:1)
* `services_kube_controller` - (DEPRECATED) Use services.kube_controller instead (list maxitems:1)
* `services_kubelet` - (DEPRECATED) Use services.kubelet instead (list maxitems:1)
* `services_kubeproxy` - (DEPRECATED) Use services.kubeproxy instead (list maxitems:1)
* `services_scheduler` - (DEPRECATED) Use services.scheduler instead (list maxitems:1)
* `ssh_agent_auth` - (Optional/Computed) SSH Agent Auth enable (bool)
* `ssh_cert_path` - (Optional) SSH Certificate Path (string)
* `ssh_key_path` - (Optional) SSH Private Key Path (string)
* `system_images` - (Optional) RKE k8s cluster system images list (list maxitems:1)
* `update_only` - (Optional) Skip idempotent deployment of control and etcd plane. Default `false` (bool)
* `upgrade_strategy` - (Optional) RKE k8s cluster upgrade strategy (list maxitems:1)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `ca_crt` - (Computed/Sensitive) RKE k8s cluster CA certificate (string)
* `client_cert` - (Computed/Sensitive) RKE k8s cluster client certificate (string)
* `client_key` - (Computed/Sensitive) RKE k8s cluster client key (string)
* `rke_state` - (Computed/Sensitive) RKE k8s cluster state (string)
* `kube_config_yaml` - (Computed/Sensitive) RKE k8s cluster kube config yaml (string)
* `internal_kube_config_yaml` - (Computed/Sensitive) RKE k8s cluster internal kube config yaml (string)
* `rke_cluster_yaml` - (Computed/Sensitive) RKE k8s cluster config yaml (string)
* `certificates` - (Computed/Sensitive) RKE k8s cluster certificates (string)
* `kube_admin_user` - (Computed) RKE k8s cluster admin user (string)
* `api_server_url` - (Computed) RKE k8s cluster api server url (string)
* `cluster_domain` - (Computed) RKE k8s cluster domain (string)
* `cluster_cidr` - (Computed) RKE k8s cluster cidr (string)
* `cluster_dns_server` - (Computed) RKE k8s cluster dns server (string)
* `control_plane_hosts` - (Computed) RKE k8s cluster control plane nodes (list)
* `etcd_hosts` - (Computed) RKE k8s cluster etcd nodes (list)
* `inactive_hosts` - (Computed) RKE k8s cluster inactive nodes (list)
* `worker_hosts` - (Computed) RKE k8s cluster worker nodes (list)
* `running_system_images` - (Computed) RKE k8s cluster running system images list (list)

## Nested blocks

### `authentication`

#### Arguments

* `sans` - (Optional/Computed) List of additional hostnames and IPs to include in the api server PKI cert (list)
* `strategy` - (Optional) Authentication strategy that will be used in RKE k8s cluster. Default: `x509` (string)
* `webhook` - (Optional/Computed) Webhook configuration options (list maxitem: 1)

#### `webhook`

##### Arguments

* `config_file` - (Optional) Multiline string that represent a custom webhook config file (string)
* `cache_timeout` - (Optional) Controls how long to cache authentication decisions (string)

### `authorization`

#### Arguments

* `mode` - (Optional) RKE mode for authorization. `rbac` and `none` modes are available. Default `rbac` (string)
* `options` - (Optional/Computed) RKE options for authorization (map)

### `bastion_host`

#### Arguments

* `address` - (Required) Address of Bastion Host (string)
* `user` - (Required) SSH User to Bastion Host (string)
* `port` - (Optional) SSH Port of Bastion Host. Default `22` (string)
* `ssh_agent_auth` - (Optional/Computed) SSH Agent Auth enable (bool)
* `ssh_cert` - (Optional/Sensitive) SSH Certificate Key (string)
* `ssh_cert_path` - (Optional) SSH Certificate Key Path (string)
* `ssh_key` - (Optional/Sensitive) SSH Private Key (string)
* `ssh_key_path` - (Optional) SSH Private Key Path (string)

### `cloud_provider`

#### Arguments

* `name` - (Required) Cloud Provider name. `aws`, `azure`, `custom`, `external`, `openstack`, `vsphere` are supported (string)
* `aws_cloud_config` - (DEPRECATED) Use aws_cloud_provider instead
* `aws_cloud_provider` - (Optional) AWS Cloud Provider config [rke-aws-cloud-provider](https://rancher.com/docs/rke/latest/en/config-options/cloud-providers/aws/) (list maxitems:1)
* `azure_cloud_config` - (DEPRECATED) Use azure_cloud_provider instead
* `azure_cloud_provider` - (Optional) Azure Cloud Provider config [rke-azure-cloud-provider](https://rancher.com/docs/rke/latest/en/config-options/cloud-providers/azure/) (list maxitems:1)
* `custom_cloud_config` - (DEPRECATED) Use custom_cloud_provider instead
* `custom_cloud_provider` - (Optional) Custom Cloud Provider config (string)
* `openstack_cloud_config` - (DEPRECATED) Use openstack_cloud_provider instead
* `openstack_cloud_provider` - (Optional/Computed) Openstack Cloud Provider config [rke-openstack-cloud-provider](https://rancher.com/docs/rke/latest/en/config-options/cloud-providers/openstack/) (list maxitems:1)
* `vsphere_cloud_config` - (DEPRECATED) Use vsphere_cloud_provider instead
* `vsphere_cloud_provider` - (Optional/Computed) Vsphere Cloud Provider config [rke-vsphere-cloud-provider](https://rancher.com/docs/rke/latest/en/config-options/cloud-providers/vsphere/) Extra argument `name` is required on `virtual_center` configuration. (list maxitems:1)

#### `aws_cloud_provider`

##### Arguments

* `global` - (Optional) (list maxitems:1)
* `service_override` - (Optional) (list)

##### `global`

###### Arguments

* `disable_security_group_ingress` - (Optional) Disables the automatic ingress creation. Default `false` (bool)
* `disable_strict_zone_check` - (Optional) Setting this to true will disable the check and provide a warning that the check was skipped. Default `false` (bool)
* `elb_security_group` - (Optional) Use these ELB security groups instead create new (string)
* `kubernetes_cluster_id` - (Optional) The cluster id we'll use to identify our cluster resources (string)
* `kubernetes_cluster_tag` - (Optional) Legacy cluster id we'll use to identify our cluster resources (string)
* `role_arn` - (Optional) IAM role to assume when interaction with AWS APIs (string)
* `route_table_id` - (Optional) Enables using a specific RouteTable (string)
* `subnet_id` - (Optional) Enables using a specific subnet to use for ELB's (string)
* `vpc` - (Optional) The AWS VPC flag enables the possibility to run the master components on a different aws account, on a different cloud provider or on-premises. If the flag is set also the KubernetesClusterTag must be provided (string)
* `zone` - (Optional) The AWS zone (string)

##### `service_override`

###### Arguments

* `key` - (DEPRECATED) Use service instead
* `service` - (Required) (string)
* `region` - (Optional) (string)
* `signing_method` - (Optional/Computed) (string)
* `signing_name` - (Optional) (string)
* `signing_region` - (Optional) (string)
* `url` - (Optional) (string)

#### `azure_cloud_provider`

##### Arguments

* `aad_client_id` - (Required/Sensitive) (string)
* `aad_client_secret` - (Required/Sensitive) (string)
* `subscription_id` - (Required/Sensitive) (string)
* `tenant_id` - (Required/Sensitive) (string)
* `aad_client_cert_password` - (Optional/Computed/Sensitive) (string)
* `aad_client_cert_path` - (Optional/Computed) (string)
* `cloud` - (Optional/Computed) (string)
* `cloud_provider_backoff` - (Optional/Computed) (bool)
* `cloud_provider_backoff_duration` - (Optional/Computed) (int)
* `cloud_provider_backoff_exponent` - (Optional/Computed) (int)
* `cloud_provider_backoff_jitter` - (Optional/Computed) (int)
* `cloud_provider_backoff_retries` - (Optional/Computed) (int)
* `cloud_provider_rate_limit` - (Optional/Computed) (bool)
* `cloud_provider_rate_limit_bucket` - (Optional/Computed) (int)
* `cloud_provider_rate_limit_qps` - (Optional/Computed) (int)
* `location` - (Optional/Computed) (string)
* `maximum_load_balancer_rule_count` - (Optional/Computed) (int)
* `primary_availability_set_name` - (Optional/Computed) (string)
* `primary_scale_set_name` - (Optional/Computed) (string)
* `resource_group` - (Optional/Computed) (string)
* `route_table_name` - (Optional/Computed) (string)
* `security_group_name` - (Optional/Computed) (string)
* `subnet_name` - (Optional/Computed) (string)
* `use_instance_metadata` - (Optional/Computed) (bool)
* `use_managed_identity_extension` - (Optional/Computed) (bool)
* `vm_type` - (Optional/Computed) (string)
* `vnet_name` - (Optional/Computed) (string)
* `vnet_resource_group` - (Optional/Computed) (string)

#### `openstack_cloud_provider`

##### Arguments

* `global` - (Required) (list maxitems:1)
* `block_storage` - (Optional/Computed) (list maxitems:1)
* `load_balancer` - (Optional/Computed) (list maxitems:1)
* `metadata` - (Optional/Computed) (list maxitems:1)
* `route` - (Optional/Computed) (list maxitems:1)

##### `global`

###### Arguments

* `auth_url` - (Required) (string)
* `password` - (Required/Sensitive) (string)
* `ca_file` - (Optional) (string)
* `domain_id` - (Optional/Sensitive) Required if `domain_name` not provided. (string)
* `domain_name` - (Optional) Required if `domain_id` not provided. (string)
* `region` - (Optional) (string)
* `tenant_id` - (Optional/Sensitive) Required if `tenant_name` not provided. (string)
* `tenant_name` - (Optional) Required if `tenant_id` not provided. (string)
* `trust_id` - (Optional/Sensitive) (string)
* `username` - (Optional) Required if `user_id` not provided. (string)
* `user_id` - (Optional/Sensitive) Required if `username` not provided. (string)

##### `block_storage`

###### Arguments

* `bs_version` - (Optional) (string)
* `ignore_volume_az` - (Optional) (string)
* `trust_device_path` - (Optional) (string)

##### `load_balancer`

###### Arguments

* `create_monitor` - (Optional) (bool)
* `floating_network_id` - (Optional) (string)
* `lb_method` - (Optional) (string)
* `lb_provider` - (Optional) (string)
* `lb_version` - (Optional) (string)
* `manage_security_groups` - (Optional) (bool)
* `monitor_delay` - (Optional) (string)
* `monitor_max_retries` - (Optional) (int)
* `monitor_timeout` - (Optional) (string)
* `subnet_id` - (Optional) (string)
* `use_octavia` - (Optional) (bool)

##### `metadata`

###### Arguments

* `request_timeout` - (Optional) (int)
* `search_order` - (Optional) (string)

##### `route`

###### Arguments

* `router_id` - (Optional) (string)

#### `vsphere_cloud_provider`

##### Arguments

* `virtual_center` - (Required) (List)
* `workspace` - (Required) (list maxitems:1)
* `disk` - (Optional/Computed) (list maxitems:1)
* `global` - (Optional/Computed) (list maxitems:1)
* `network` - (Optional/Computed) (list maxitems:1)

##### `virtual_center`

###### Arguments

* `datacenters` - (Required) (string)
* `name` - (Required) Name of virtualcenter config for Vsphere Cloud Provider config (string)
* `password` - (Required/Sensitive) (string)
* `user` - (Required/Sensitive) (string)
* `port` - (Optional) (string)
* `soap_roundtrip_count` - (Optional) (int)

##### `workspace`

###### Arguments

* `datacenter` - (Required) (string)
* `server` - (Required) (string)
* `default_datastore` - (Optional) (string)
* `folder` - (Optional) (string)
* `resourcepool_path` - (Optional) (string)

##### `disk`

###### Arguments

* `scsi_controller_type` - (Optional) (string)

##### `global`

###### Arguments

* `datacenters` - (Optional/Computed) (string)
* `datastore` - (Optional) (string)
* `insecure_flag` - (Optional) (bool)
* `password` - (Optional/Sensitive) (string)
* `user` - (Optional/Sensitive) (string)
* `port` - (Optional) (string)
* `soap_roundtrip_count` - (Optional) (int)
* `working_dir` - (Optional) (string)
* `vm_name` - (Optional) (string)
* `vm_uuid` - (Optional) (string)

##### `network`

###### Arguments

* `public_network` - (Optional) (string)

### `dns`

#### Arguments

* `nodelocal` - (Optional) Nodelocal dns config  (list Maxitem: 1)
* `node_selector` - (Optional) Node selector key pair (map)
* `provider` - (Optional) DNS provider. `kube-dns`, `coredns` (default), and `none` are supported (string)
* `reverse_cidrs` - (Optional) Reverse CIDRs  (list)
* `upstream_nameservers` - (Optional) Upstream nameservers  (list)

#### `nodelocal`

##### Arguments

* `ip_address` - (required) Nodelocal dns ip address (string)
* `node_selector` - (Optional) Node selector key pair (map)

### `ingress`

#### Arguments

* `dns_policy` - (Optional) Ingress controller DNS policy. `ClusterFirstWithHostNet`, `ClusterFirst`, `Default`, and `None` are supported. [K8S dns Policy](https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/#pod-s-dns-policy) (string)
* `extra_args` - (Optional) Extra arguments for the ingress controller (map)
* `node_selector` - (Optional) Node selector key pair (map)
* `options` - (Optional) Ingress controller options (map)
* `provider` - (Optional) Ingress controller provider. `nginx` (default), and `none` are supported (string)

### `monitoring`

#### Arguments

* `node_selector` - (Optional) Node selector key pair (map)
* `options` - (Optional) Monitoring options (map)
* `provider` - (Optional/Computed) Monitoring provider (string)

### `network`

#### Arguments

* `calico_network_provider` - (Optional) Calico network provider config (list maxitems:1)
* `canal_network_provider` - (Optional) Canal network provider config (list maxitems:1)
* `flannel_network_provider` - (Optional) Flannel network provider config (list maxitems:1)
* `weave_network_provider` - (Optional) Weave network provider config (list maxitems:1)
* `mtu` - (Optional) Network provider MTU. Default `0` (int)
* `options` - (Optional/Computed) Network provider options (map)
* `plugin` - (Optional) Network provider plugin. `calico`, `canal` (default), `flannel`, `none` and `weave` are supported. (string)

#### `calico_network_provider`

##### Arguments

* `cloud_provider` - (Optional/Computed) Calico cloud provider (string)

#### `canal_network_provider`

##### Arguments

* `iface` - (Optional/Computed) Canal network interface (string)

#### `flannel_network_provider`

##### Arguments

* `iface` - (Optional/Computed) Flannel network interface (string)

#### `weave_network_provider`

##### Arguments

* `password` - (Optional/Computed) Weave password (string)

### `nodes`

#### Arguments

* `address` - (Required) Address ip for node (string)
* `role` - (Required) Node roles in k8s cluster. `controlplane`, `etcd` and `worker` are supported. (list)
* `user` - (Required/Sensitive) SSH user that will be used by RKE (string)
* `docker_socket` - (Optional) Docker socket on the node that will be used in tunneling (string)
* `hostname_override` - (Optional) Hostname override for node (string)
* `internal_address` - (Optional) Internal address that will be used for components communication (string)
* `labels` - (Optional) Node labels (map)
* `node_name` - (Optional) Name of the host provisioned via docker machine (string)
* `port` - (Optional) Port used for SSH communication (string)
* `ssh_agent_auth` - (Optional/Computed) SSH Agent Auth enable (bool)
* `ssh_cert` - (Optional/Sensitive) SSH Certificate (string)
* `ssh_cert_path` - (Optional) SSH Certificate path (string)
* `ssh_key` - (Optional/Sensitive) SSH Private Key (string)
* `ssh_key_path` - (Optional) SSH Private Key path (string)
* `taints` - (Optional) Node taints (list)

#### `taint`

##### Arguments

* `key` - (Required) Taint key (string)
* `value` - (Required) Taint value (string)
* `effect` - (Optional) Taint effect. `NoExecute`, `NoSchedule` (default) and `PreferNoSchedule` are supported (string)

### `private_registries`

#### Arguments

* `url` - (Required) Registry URL (string)
* `is_default` - (Optional) Set as default registry. Default `false` (bool)
* `password` - (Optional/Sensitive) Registry password (string)
* `user` - (Optional/Sensitive) Registry user (string)

### `restore`

#### Arguments

* `restore` - (Optional) Restore cluster. Default `false` (bool)
* `snapshot_name` - (Optional) Snapshot name (string)

### `rotate_certificates`

#### Arguments

* `ca_certificates` - (Optional) Rotate CA Certificates. Default `false` (bool)
* `services` - (Optional) Services to rotate their certs. `etcd`, `kubelet`, `kube-apiserver`, `kube-proxy`, `kube-scheduler` and `kube-controller-manager` are supported (list)

### `services`

#### Arguments

* `etcd` - (Optional/Computed) Etcd options for RKE services (list maxitems:1)
* `kube_api` - (Optional/Computed) Kube API options for RKE services (list maxitems:1)
* `kube_controller` - (Optional/Computed) Kube Controller options for RKE services (list maxitems:1)
* `kubelet` - (Optional/Computed) Kubelet options for RKE services (list maxitems:1)
* `kubeproxy` - (Optional/Computed) Kubeproxy options for RKE services (list maxitems:1)
* `scheduler` - (Optional/Computed) Scheduler options for RKE services (list maxitems:1)

#### `etcd`

##### Arguments

* `backup_config` - (Optional/Computed) Backup options for etcd service. Just for Rancher v2.2.x (list maxitems:1)
* `ca_cert` - (Optional/Computed/Sensitive) TLS CA certificate for etcd service (string)
* `cert` - (Optional/Computed/Sensitive) TLS certificate for etcd service (string)
* `creation` - (Optional/Computed) Creation option for etcd service (string)
* `external_urls` - (Optional/Computed) External urls for etcd service (list)
* `extra_args` - (Optional/Computed) Extra arguments for etcd service (map)
* `extra_binds` - (Optional/Computed) Extra binds for etcd service (list)
* `extra_env` - (Optional/Computed) Extra environment for etcd service (list)
* `gid` - (Optional) Etcd service GID. Default: `0`. For Rancher v2.3.x or above (int)
* `image` - (Optional/Computed) Docker image for etcd service (string)
* `key` - (Optional/Computed/Sensitive) TLS key for etcd service (string)
* `path` - (Optional/Computed) Path for etcd service (string)
* `retention` - (Optional/Computed) Retention option for etcd service (string)
* `snapshot` - (Optional) Snapshot option for etcd service. Default `true` (bool)
* `uid` - (Optional) Etcd service UID. Default: `0`. For Rancher v2.3.x or above (int)

##### `backup_config`

###### Arguments

* `enabled` - (Optional) Enable etcd backup. Default `true` (bool)
* `interval_hours` - (Optional) Interval hours for etcd backup. Default `12` (int)
* `retention` - (Optional) Retention for etcd backup. Default `6` (int)
* `s3_backup_config` - (Optional) S3 config options for etcd backup (list maxitems:1)
* `safe_timestamp` - (Optional) Safe timestamp for etcd backup. Default: `false` (bool)

###### `s3_backup_config`

###### Arguments

* `access_key` - (Optional/Sensitive) Access key for S3 service (string)
* `bucket_name` - (Optional) Bucket name for S3 service (string)
* `custom_ca` - (Optional/Sensitive) Base64 encoded custom CA for S3 service. Use filebase64(<FILE>) for encoding file. Available from Rancher v2.2.5 (string)
* `endpoint` - (Optional) Endpoint for S3 service (string)
* `folder` - (Optional) Folder for S3 service. Available from Rancher v2.2.7 (string)
* `region` - (Optional) Region for S3 service (string)
* `secret_key` - (Optional/Sensitive) Secret key for S3 service (string)

#### `kube_api`

##### Arguments

* `always_pull_images` - (Optional/Computed) Enable [AlwaysPullImages](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/#alwayspullimages) Admission controller plugin. [Rancher docs](https://rancher.com/docs/rke/latest/en/config-options/services/#kubernetes-api-server-options) (bool)
* `audit_log` - (Optional/Computed) K8s audit log configuration. (list maxitem: 1)
* `event_rate_limit` - (Optional) K8s event rate limit configuration. (list maxitem: 1)
* `extra_args` - (Optional/Computed) Extra arguments for kube API service (map)
* `extra_binds` - (Optional/Computed) Extra binds for kube API service (list)
* `extra_env` - (Optional/Computed) Extra environment for kube API service (list)
* `image` - (Optional/Computed) Docker image for kube API service (string)
* `pod_security_policy` - (Optional/Computed) Pod Security Policy option for kube API service (bool)
* `secrets_encryption_config` - (Optional) [Encrypt k8s secret data configration](https://rancher.com/docs/rke/latest/en/config-options/secrets-encryption/). (list maxitem: 1)
* `service_cluster_ip_range` - (Optional/Computed) Service Cluster IP Range option for kube API service (string)
* `service_node_port_range` - (Optional/Computed) Service Node Port Range option for kube API service (string)

##### `audit_log`

###### Arguments

* `configuration` - (Optional/Computed) Audit log configuration. (list maxtiem: 1)
* `enabled` - (Optional/Computed) Enable audit log (bool)

###### `configuration`

###### Arguments

* `format` - (Optional/Computed) Audit log format (string)
* `max_age` - (Optional/Computed) Audit log max age (int)
* `max_backup` - (Optional/Computed) Audit log max backup. Default: `10` (int)
* `max_size` - (Optional/Computed) Audit log max size. Default: `100` (int)
* `path` - (Optional/Computed) Audit log path. Default: `/var/log/kube-audit/audit-log.json` (string)
* `policy` - (Optional/Computed) Audit policy json encoded definition. `"apiVersion"` and `"kind":"Policy","rules"` fields are required in the json. Ex. `jsonencode({"apiVersion":"audit.k8s.io/v1","kind":"Policy","rules":[{"level":"RequestResponse","resources":[{"group":"","resources":["pods"]}]}]})` [More info](https://rancher.com/docs/rke/latest/en/config-options/audit-log/) (string)

##### `event_rate_limit`

###### Arguments

* `enabled` - (Optional/Computed) Enable event rate limit (bool)
* `configuration` - (Optional) Event rate limit yaml encoded configuration. `"apiVersion"` and `"kind":"Configuration"` fields are required in the yaml. Ex. `apiVersion: eventratelimit.admission.k8s.io/v1alpha1\nkind: Configuration\nlimits:\n- type: Server\n  burst: 30000\n  qps: 6000\n` [More info](https://rancher.com/docs/rke/latest/en/config-options/rate-limiting/) (string)

##### `secrets_encryption_config`

###### Arguments

* `enabled` - (Optional/Computed) Enable secrets encryption (bool)
* `custom_config` - (Optional) Secrets encryption yaml encoded custom configuration. `"apiVersion"` and `"kind":"EncryptionConfiguration"` fields are required in the yaml. Ex. `apiVersion: apiserver.config.k8s.io/v1\nkind: EncryptionConfiguration\nresources:\n- resources:\n  - secrets\n  providers:\n  - aescbc:\n      keys:\n      - name: k-fw5hn\n        secret: RTczRjFDODMwQzAyMDVBREU4NDJBMUZFNDhCNzM5N0I=\n    identity: {}\n` [More info](https://rancher.com/docs/rke/latest/en/config-options/secrets-encryption/) (string)

#### `kube_controller`

##### Arguments

* `cluster_cidr` - (Optional/Computed) Cluster CIDR option for kube controller service (string)
* `extra_args` - (Optional/Computed) Extra arguments for kube controller service (map)
* `extra_binds` - (Optional/Computed) Extra binds for kube controller service (list)
* `extra_env` - (Optional/Computed) Extra environment for kube controller service (list)
* `image` - (Optional/Computed) Docker image for kube controller service (string)
* `service_cluster_ip_range` - (Optional/Computed) Service Cluster ip Range option for kube controller service (string)

#### `kubelet`

##### Arguments

* `cluster_dns_server` - (Optional/Computed) Cluster DNS Server option for kubelet service (string)
* `cluster_domain` - (Optional) Cluster Domain option for kubelet service. Default `cluster.local` (string)
* `extra_args` - (Optional/Computed) Extra arguments for kubelet service (map)
* `extra_binds` - (Optional/Computed) Extra binds for kubelet service (list)
* `extra_env` - (Optional/Computed) Extra environment for kubelet service (list)
* `fail_swap_on` - (Optional/Computed) Enable or disable failing when swap on is not supported (bool)
* `generate_serving_certificate` [Generate a certificate signed by the kube-ca](https://rancher.com/docs/rke/latest/en/config-options/services/#kubelet-serving-certificate-requirements). Default `false` (bool)
* `image` - (Optional/Computed) Docker image for kubelet service (string)
* `infra_container_image` - (Optional/Computed) Infra container image for kubelet service (string)

#### `kubeproxy`

##### Arguments

* `extra_args` - (Optional/Computed) Extra arguments for kubeproxy service (map)
* `extra_binds` - (Optional/Computed) Extra binds for kubeproxy service (list)
* `extra_env` - (Optional/Computed) Extra environment for kubeproxy service (list)
* `image` - (Optional/Computed) Docker image for kubeproxy service (string)

#### `scheduler`

##### Arguments

* `extra_args` - (Optional/Computed) Extra arguments for scheduler service (map)
* `extra_binds` - (Optional/Computed) Extra binds for scheduler service (list)
* `extra_env` - (Optional/Computed) Extra environment for scheduler service (list)
* `image` - (Optional/Computed) Docker image for scheduler service (string)

### `system_images`

#### Arguments

* `etcd` - (Optional) Docker image for etcd (string)
* `alpine` - (Optional) Docker image for alpine (string)
* `nginx_proxy` - (Optional) Docker image for nginx_proxy (string)
* `cert_downloader` - (Optional) Docker image for cert_downloader (string)
* `kubernetes_services_sidecar` - (Optional) Docker image for kubernetes_services_sidecar (string)
* `kube_dns` - (Optional) Docker image for kube_dns (string)
* `dnsmasq` - (Optional) Docker image for dnsmasq (string)
* `kube_dns_sidecar` - (Optional) Docker image for kube_dns_sidecar (string)
* `kube_dns_autoscaler` - (Optional) Docker image for kube_dns_autoscaler (string)
* `coredns` - (Optional) Docker image for coredns (string)
* `coredns_autoscaler` - (Optional) Docker image for coredns_autoscaler (string)
* `nodelocal` - (Optional) Docker image for nodelocal (string)
* `kubernetes` - (Optional) Docker image for kubernetes (string)
* `flannel` - (Optional) Docker image for flannel (string)
* `flannel_cni` - (Optional) Docker image for flannel_cni (string)
* `calico_node` - (Optional) Docker image for calico_node (string)
* `calico_cni` - (Optional) Docker image for calico_cni (string)
* `calico_controllers` - (Optional) Docker image for calico_controllers (string)
* `calico_ctl` - (Optional) Docker image for calico_ctl (string)
* `calico_flex_vol` - (Optional) Docker image for calico_flex_vol (string)
* `canal_node` - (Optional) Docker image for canal_node (string)
* `canal_cni` - (Optional) Docker image for canal_cni (string)
* `canal_flannel` - (Optional) Docker image for canal_flannel (string)
* `canal_flex_vol` - (Optional) Docker image for canal_flex_vol (string)
* `weave_node` - (Optional) Docker image for weave_node (string)
* `weave_cni` - (Optional) Docker image for weave_cni (string)
* `pod_infra_container` - (Optional) Docker image for pod_infra_container (string)
* `ingress` - (Optional) Docker image for ingress (string)
* `ingress_backend` - (Optional) Docker image for ingress_backend (string)
* `metrics_server` - (Optional) Docker image for metrics_server (string)
* `windows_pod_infra_container` - (Optional) Docker image for windows_pod_infra_container (string)

### `upgrade_strategy`

#### Arguments

* `drain` - (Optional/Computed) RKE drain nodes (bool)
* `drain_input` - (Optional/Computed) RKE drain node input (list Maxitems: 1)
* `max_unavailable_controlplane` - (Optional/Computed) RKE max unavailable controlplane nodes (string)
* `max_unavailable_worker` - (Optional/Computed) RKE max unavailable worker nodes (string)

#### `drain_input`

##### Arguments

* `delete_local_data` - (Optional/Computed) Delete RKE node local data (bool)
* `force` - (Optional/Computed) Force RKE node drain (bool)
* `grace_period` - (Optional/Computed) RKE node drain grace period (int)
* `ignore_daemon_sets` - (Optional/Computed) Ignore RKE daemon sets (bool)
* `timeout` - (Optional/Computed) RKE node drain timeout (int)

## Timeouts

`rke_cluster` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `30 minutes`) Used for creating clusters.
- `update` - (Default `30 minutes`) Used for cluster modifications.
- `delete` - (Default `30 minutes`) Used for deleting clusters.

## Import

rke_cluster can be imported using the RKE cluster config and state files as ID in the format `<cluster_config_file>:<rke_state_file>`

```
$ terraform import rke_cluster.foo  &lt;cluster_config_file&gt;:&lt;rke_state_file&gt;
```

As experimental feature, dind rke_cluster can be also imported adding `dind` as 3rd import parameter `<cluster_config_file>:<rke_state_file>:dind`


```
$ terraform import rke_cluster.foo  &lt;cluster_config_file&gt;:&lt;rke_state_file&gt;:dind
```
