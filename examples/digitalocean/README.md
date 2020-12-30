# How To Deploy Kubernetes Clusters on DigitalOcean using Terraform and Terraform RKE Provider

This repository is an example for building a Kubernetes cluster using Terraform and Terraform RKE provider on DigitalOcean.

## What it does

By default, the following will be deployed to your DigitalOcean account:
-   One 5-node RKE cluster, with each node installed with RancherOS, configured as follows:
    +   3 control plane nodes, each assigned the `controlplane`, `etcd`, and `worker` roles
        -   By default, these nodes are provisioned with 2 vCPUs and 4 GB RAM each
    +   2 `worker` nodes
        -   By default, these nodes are provisioned with 1 vCPU and 2 GB RAM each
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

### Review/revise the example cluster configuration in `rke.tf`

Most of the cluster configuration is defined in `./rke.tf`.  The file has been meticulously commented to make it easier to understand what you may want to adjust.  For instance, you can:

*   Change the `cluster_domain` (the default, if not changed, will be `cluster.local`, which is also the RKE default as of this writing)
*   Change the DigitalOcean `region` that the cluster will be deployed in
*   Adjust the size of the DigitalOcean droplets that will be provisioned as control plane and worker nodes in the cluster (`droplet_size`)
*   Increase or reduce the number of control plane and worker nodes to be deployed in the cluster (`node_count`)

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

When `terraform apply` is completed, the kubeconfig file should be created in the current directory.

1.  Set KUBECONFIG environment variable for `kubectl`
    ```console
    $ export KUBECONFIG=${PWD}/kube_config_cluster.yml
    ```
2.  Then, the `kubectl` command can be used, like so:
    -   Check component statuses
        ```console
        $ kubectl get cs
        NAME                 STATUS    MESSAGE             ERROR
        controller-manager   Healthy   ok                  
        scheduler            Healthy   ok                  
        etcd-0               Healthy   {"health":"true"}   
        etcd-1               Healthy   {"health":"true"}   
        etcd-2               Healthy   {"health":"true"} 
        ```
    -   Check node statuses
        ```console
        $ kubectl get nodes
        NAME                          STATUS   ROLES                      AGE    VERSION
        controlplane1.cluster.local   Ready    controlplane,etcd,worker   3m4s   v1.18.12
        controlplane2.cluster.local   Ready    controlplane,etcd,worker   3m2s   v1.18.12
        controlplane3.cluster.local   Ready    controlplane,etcd,worker   3m2s   v1.18.12
        worker1.cluster.local         Ready    worker                     3m4s   v1.18.12
        worker2.cluster.local         Ready    worker                     3m3s   v1.18.12
        ```
    -   Find the External IP assigned to the NGINX Ingress Controller `LoadBalancer` resource
        ```console
        $ kubectl get svc -n ingress-nginx
        NAME                                 TYPE           CLUSTER-IP     EXTERNAL-IP      PORT(S)                      AGE
        ingress-nginx-controller             LoadBalancer   10.43.91.160   104.248.108.30   80:31081/TCP,443:30350/TCP   2m16s
        ingress-nginx-controller-admission   ClusterIP      10.43.168.27   <none>           443/TCP                      2m16s

### Tear down the cluster

You can use Terraform to destroy the cluster and remove it from your DigitalOcean account using the command below.

**NOTE:** Terraform cannot remove the DigitalOcean load balancer associated with the cluster, since it is created by the NGINX Ingress Controller running on the Kubernetes cluster itself, and not by Terraform.  The same will be true if you created any persistent volumes on the cluster.  You will need to remove any such resources either by using `kubectl` prior to destroying the cluster with Terraform, or manually via the DigitalOcean console or the `doctl` CLI after destroying the cluster.

```console
$ terraform destroy -var do_token=$DIGITALOCEAN_TOKEN
```