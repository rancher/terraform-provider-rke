## 1.0.1 (June 30, 2020)

FEATURES:



ENHANCEMENTS:

* Updated go to 1.13
* Updated RKE to v1.1.3
* Adapt provider release to met terraform registry requirements

BUG FIXES:

* Updated `cloud_provider` and `bastion_host` arguments on `rancher2_cluster` resource to fix false diff 
* Updated `ignore_docker_version` argument on `rancher2_cluster` resource to fix provider crash 

## 1.0.0 (May 7, 2020)

FEATURES:



ENHANCEMENTS:

* Set arguments as computed to avoid false diff:
  * `upgrade_strategy` on `rke_cluster` resource 
  * `ssh_key_path` on `bastion_host` argument 
  * `audit_log` on `services.kube_api` argument
* Add `external` as allowed value on `cloud_provider` argument
* Add `nodelocal` argument to `dns` argument on `rke_cluster` resource
* Updated RKE to v1.1.1

BUG FIXES:



## 1.0.0-rc5 (April 3, 2020)

FEATURES:

* **Deprecated Argument:** `nodes_conf` - Use `cluster_yaml` instead
* **Deprecated Argument:** `internal_kube_config_yaml` - Use `kube_config_yaml` instead
* **New Argument:** `cluster_yaml` - RKE cluster config yaml
* **New Import:** `rke_cluster` - RKE cluster import is supported

ENHANCEMENTS:

* Updated `hashicorp/terraform-plugin-sdk` go modules and vendor files to v1.8.0
* Updated go modules and vendor files to support [RKE v1.1.0](https://github.com/rancher/rke/releases/tag/v1.1.0)
* Added `upgrade_strategy` argument to `rke_cluster` resource
* Updated `kubernetes_version` argument to get default and available k8s versions from rke metadata

BUG FIXES:

* Fixed computed fields to avoid inconsistent plan
* Disabled `debug` config option until next rke release (Breaking logs)
* Fixed segmentation fault with deprecated fields
* Fixed `vsphere_cloud_provider` argument to avoid false diff
* Fixed segmentation fault with deprecated fields

## 1.0.0-rc4 (March 13, 2020)

FEATURES:



ENHANCEMENTS:

* Updated `log_file` to sync logs


BUG FIXES:

* Fixed computed fields to avoid inconsistent plan

## 1.0.0-rc3 (February 28, 2020)

FEATURES:



ENHANCEMENTS:

* Added `cert_dir`, `custom_certs` and `update_only` arguments to rke cluster configuration
* Refactored `rke_cluster` resource: 
  * Added uuid as tfstate id
  * Added `CustomizeDiff` to control changes
  * Saving state on any execution
* Added `kube_api.audit_log.Configuration.policy` argument to `services` argument
* Added `dind` support
* Added acceptance tests
* Added `debug` and `log_file` provider configuration


BUG FIXES:

* Fixed k8s version upgrade on `rke_cluster` resource


## 1.0.0-rc2 (February 6, 2020)

FEATURES:



ENHANCEMENTS:

* Updated go modules and vendor files to support RKE v1.0.4
* Added `mtu` argument to network configuration
* Save `rke_cluster` resource data in tfstate even if `clusterUp` fails, to be able to retry

BUG FIXES:

* Fix `nodes.port` argument definition on `rke_cluster` resource
* Fix false diffs setting: 
  * set `extra_*` and `image` arguments as computed on all `services` nested arguments
  * set `bastion_host` and `cloud_provider` arguments as non computed
  * set `hostname_override` and `internal_address` arguments as computed on `nodes` argument

## 1.0.0-rc1 (January 21, 2020)

FEATURES:

* **Deprecated Datasource:** `rke_node_parameter` - Use `rke_cluster` resource + dynamic instead
* **New Argument:** `services`
* **New Argument:** `services.etcd`  
* **Deprecated Argument:** `services_etcd` - Use `services.etcd` instead
* **New Argument:** `services.kube_api` 
* **Deprecated Argument:** `services_kube_api` - Use `services.kube_api` instead
* **New Argument:** `services.kube_controller` 
* **Deprecated Argument:** `services_kube_controller` - Use `services.kube_controller` instead
* **New Argument:** `services.kubelet` 
* **Deprecated Argument:** `services_kubelet` - Use `services.kubelet` instead
* **New Argument:** `services.kubeproxy` 
* **Deprecated Argument:** `services_kubeproxy` - Use `services.kubeproxy` instead
* **New Argument:** `services.scheduler` 
* **Deprecated Argument:** `services_scheduler` - Use `services.scheduler` instead
* **New Argument:** `cloud_provider.aws_cloud_provider`
* **Deprecated Argument:** `cloud_provider.aws_cloud_config` - Use `cloud_provider.aws_cloud_provider` instead
* **New Argument:** `cloud_provider.azure_cloud_provider`
* **Deprecated Argument:** `cloud_provider.azure_cloud_config` - Use `cloud_provider.azure_cloud_provider` instead
* **New Argument:** `cloud_provider.custom_cloud_provider`
**Deprecated Argument:** `cloud_provider.custom_cloud_config` - Use `cloud_provider.custom_cloud_provider` instead
* **New Argument:** `cloud_provider.openstack_cloud_provider`
* **Deprecated Argument:** `cloud_provider.openstack_cloud_config` - Use `cloud_provider.openstack_cloud_provider` instead
* **New Argument:** `cloud_provider.vsphere_cloud_provider`
* **Deprecated Argument:** `cloud_provider.vsphere_cloud_config` - Use `cloud_provider.vsphere_cloud_provider` instead


ENHANCEMENTS:

* Compatible with RKE v1.0.0

BUG FIXES:



**Important** There are some breaking changes from previous provider version. Some provider arguments has been deprecated, please take a look to [Documentation](website/docs)