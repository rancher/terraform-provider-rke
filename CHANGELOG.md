## 1.0.0-rc2 (Unreleased)

FEATURES:



ENHANCEMENTS:

* Updated go modules and vendor files to support RKE v1.0.2

BUG FIXES:



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