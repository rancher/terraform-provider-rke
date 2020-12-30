# How To Deploy Kubernetes Clusters on DigitalOcean using Terraform and Terraform RKE Provider

This repository is an example for building a Kubernetes cluster using Terraform and Terraform RKE provider on DigitalOcean.

## What it does

By default, the following will be deployed to your DigitalOcean account:
-   One 4-node RKE cluster with nodes provisioned with 2 vCPUs and 4 GB RAM, each installed with RancherOS, configured as follows:
    +   1 node assigned the `controlplane` and `etcd` roles
    +   3 `worker` nodes
-   Additionally, the following will be deployed to the RKE cluster:
    +   [DigitalOcean Container Storage Interface](https://github.com/digitalocean/csi-digitalocean) (`csi-digitalocean`) v2.1.1
    +   [Kubernetes Cloud Controller Manager for DigitalOcean](https://github.com/digitalocean/digitalocean-cloud-controller-manager) (`digitalocean-cloud-controller-manager`) v0.1.30
    +   [NGINX Ingress Controller](https://github.com/kubernetes/ingress-nginx) (`ingress-nginx`) v0.42.0

## How to use

### Requirements

-   [terraform](https://terraform.io) v0.13+
-   [terraform-provider-rke](https://github.com/rancher/terraform-provider-rke)
-   Valid DigitalOcean API token
-   [optional] `kubectl` command

### Deploy Kubernetes Cluster on DigitalOcean

1.  Clone this repo
    ```console
    $ git clone https://github.com/rancher/terraform-provider-rke
    $ cd terraform-provider-rke/examples/digialocean
    ```
2.  Set API token to environment variables
    ```console
    $ export DIGITALOCEAN_TOKEN="<your-token>"
    ```
3.  Deploy
    ```console
    $ terraform init && terraform apply -var do_token=$DIGITALOCEAN_TOKEN
    ```

### Verify and start using the cluster

When "terraform apply" is completed, the kubeconfig file should be created in the current directory.

1.  Set KUBECONFIG environment variable for kubectl 
    ```console
    $ export KUBECONFIG=${PWD}/kube_config_cluster.yml
    ```
2.  Then, the `kubectl` command can be used, like so:
    -   Check component statuses
        ```console
        $ kubectl get cs

        NAME                 STATUS    MESSAGE              ERROR
        controller-manager   Healthy   ok                   
        scheduler            Healthy   ok                   
        etcd-0               Healthy   {"health": "true"}  
        ```
    -   Check node statuses
        ```console
        $ kubectl get nodes

        NAME                                             STATUS    ROLES               AGE       VERSION
        ip-xx-xx-xx-xx.ap-northeast-1.compute.internal   Ready     controlplane,etcd   1m        v1.10.1
        ip-xx-xx-xx-xx.ap-northeast-1.compute.internal   Ready     worker              1m        v1.10.1
        ip-xx-xx-xx-xx.ap-northeast-1.compute.internal   Ready     worker              1m        v1.10.1
        ip-xx-xx-xx-xx.ap-northeast-1.compute.internal   Ready     worker              1m        v1.10.1
        ```
### Tear down the cluster

You can use Terraform to destroy the cluster and remove it from your DigitalOcean account using the following command:

**NOTE:** Terraform cannot remove the DigitalOcean load balancer associated with the cluster, since it is created by the NGINX Ingress Controller running on the Kubernetes cluster itself, and not by Terraform.  The same will be true if you created any persistent volumes on the cluster.  You will need to remove any such resources either by using `kubectl` prior to destroying the cluster with Terraform, or manually in the DigitalOcean console after destroying the cluster.

```console
$ terraform destroy -var do_token=$DIGITALOCEAN_TOKEN
```