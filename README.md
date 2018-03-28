# terraform-provider-rke


[![Go Report Card](https://goreportcard.com/badge/github.com/yamamoto-febc/terraform-provider-rke)](https://goreportcard.com/report/github.com/yamamoto-febc/terraform-provider-rke)
[![Build Status](https://travis-ci.org/yamamoto-febc/terraform-provider-rke.svg?branch=master)](https://travis-ci.org/yamamoto-febc/terraform-provider-rke)

Terraform RKE providers can easily deploy Kubernetes clusters with [Rancher Kubernetes Engine](https://github.com/rancher/rke).  

#### Compatible Versions

- Terraform: v0.11+
- RKE: v1.5 (Kubernetes v1.8)

## Installation

- Install [Terraform](https://www.terraform.io/downloads.html) v0.11+ 
- Add the [`terraform-provider-rke` plugin binary](https://github.com/yamamoto-febc/terraform-provider-rke/releases/latest) In `~/.terraform.d/plugins/` 

## Usage

### Target Node Requirements

It is same as the [requirements of RKE](https://github.com/rancher/rke/blob/master/README.md#requirements)

> #### Requirements of RKE
>
> - Docker versions 1.12.6, 1.13.1, or 17.03 should be installed for Kubernetes 1.8.
> - OpenSSH 7.0+ must be installed on each node for stream local forwarding to work.
> - The SSH user used for node access must be a member of the `docker` group:
> - Ports 6443, 2379, and 2380 should be opened between cluster nodes.
> - Swap disabled on worker nodes.

### Examples

#### Minimal example

```hcl
resource rke_cluster "cluster" {
  nodes = [
    {
      address = "1.2.3.4"
      user    = "rancher"
      role    = ["controlplane", "worker", "etcd"]
      ssh_key = "${file("~/.ssh/id_rsa")}"
    },
  ]
}

###############################################################################
# If you need kubeconfig.yml for using kubectl, please uncomment follows.
###############################################################################
//resource "local_file" "kube_cluster_yaml" {
//  filename = "${path.root}/kube_config_cluster.yml"
//  content = "${rke_cluster.cluster.kube_config_yaml}"
//}
```

#### Deploying Rancher 2.0 using terraform-provider-rke

You can view deploying Rancher 2.0 example of tffile, [here](examples/rancher_server/example.tf).


#### Full example

You can view full example of tffile, [here](examples/full/example.tf).

## License

 `terraform-provider-rke` Copyright (C) 2018 Kazumichi Yamamoto.

  This project is published under [Apache 2.0 License](LICENSE.txt).
  
## Author

  * Kazumichi Yamamoto ([@yamamoto-febc](https://github.com/yamamoto-febc))