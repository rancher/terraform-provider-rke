# How To Deploy Kubernetes Clusters on DigitalOcean using Terraform and Terraform RKE Provider

This repository is an examples for building a Kubernetes cluster using Terraform and Terraform RKE provider on DigitalOcean.

## How to use

### Requirements

- [terraform](https://terraform.io) v0.13+
- [terraform-provider-rke](https://github.com/rancher/terraform-provider-rke)
- Valid DigitalOcean API token
- [optional] `kubectl` command

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
