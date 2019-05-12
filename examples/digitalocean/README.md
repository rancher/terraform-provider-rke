# How To Deploy Kubernetes Clusters on DigitalOcean using Terraform and Terraform RKE Provider

This repository is an examples for building a Kubernetes cluster using Terraform and Terraform RKE provider on DigitalOcean.

## How to use

### Requirements

- [terraform](https://terraform.io) v0.11+
- [terraform-provider-rke](https://github.com/yamamoto-febc/terraform-provider-rke)
- Valid DigitalOcean API token
- [optional] `kubectl` command

### Deploy Kubernetes Cluster on AWS

```console
#clone this repo
$ git clone https://github.com/yamamoto-febc/terraform-provider-rke
$ cd terraform-provider-rke/examples/digialocean

#set API token to environment variables
$ export DIGITALOCEAN_TOKEN="<your-token>"

#deploy
$ terraform init && terraform apply -var do_token=$DIGITALOCEAN_TOKEN

###########################################################################
#When "terraform apply" is completed, 
#kubeconfig file should be created in the current directory 
###########################################################################

#set KUBECONFIG environment variable for kubectl 
$ export KUBECONFIG=${PWD}/kube_config_cluster.yml 

###########################################################################
#Then, kubectl command can be used
###########################################################################

#component statuses
$ kubectl get cs

NAME                 STATUS    MESSAGE              ERROR
controller-manager   Healthy   ok                   
scheduler            Healthy   ok                   
etcd-0               Healthy   {"health": "true"}  

#nodes
$ kubectl get nodes

NAME                                             STATUS    ROLES               AGE       VERSION
ip-xx-xx-xx-xx.ap-northeast-1.compute.internal   Ready     controlplane,etcd   1m        v1.10.1
ip-xx-xx-xx-xx.ap-northeast-1.compute.internal   Ready     worker              1m        v1.10.1
ip-xx-xx-xx-xx.ap-northeast-1.compute.internal   Ready     worker              1m        v1.10.1
ip-xx-xx-xx-xx.ap-northeast-1.compute.internal   Ready     worker              1m        v1.10.1
```

