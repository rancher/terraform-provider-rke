# terraform-provider-rke


[![Go Report Card](https://goreportcard.com/badge/github.com/yamamoto-febc/terraform-provider-rke)](https://goreportcard.com/report/github.com/yamamoto-febc/terraform-provider-rke)
[![Build Status](https://travis-ci.org/yamamoto-febc/terraform-provider-rke.svg?branch=master)](https://travis-ci.org/yamamoto-febc/terraform-provider-rke)

Terraform RKE providers can easily deploy Kubernetes clusters with [Rancher Kubernetes Engine](https://github.com/rancher/rke).  

#### Compatible Versions

- Terraform: v0.11+(includes v0.12)
- RKE: v0.2.5

**Current master is focusing on RKE v0.2.x. If you use RKE v0.1.x, please see [rke/v0.1.x](https://github.com/yamamoto-febc/terraform-provider-rke/tree/rke/v0.1.x) branch.**

## Installation

- Install [Terraform](https://www.terraform.io/downloads.html) v0.11+(includes v0.12) 
- Add the [`terraform-provider-rke` plugin binary](https://github.com/yamamoto-febc/terraform-provider-rke/releases/latest) In `~/.terraform.d/plugins/` 

## Usage

*Note: Following examples are compatible with Terraform v0.12.*

### Target Node Requirements

It is same as the [requirements of RKE](https://github.com/rancher/rke/blob/master/README.md#requirements).  

### Examples

#### Minimal example

```hcl
resource rke_cluster "cluster" {
  nodes {
      address = "1.2.3.4"
      user    = "rancher"
      role    = ["controlplane", "worker", "etcd"]
  }
}

###############################################################################
# If you need kubeconfig.yml for using kubectl, please uncomment follows.
###############################################################################
//resource "local_file" "kube_cluster_yaml" {
//  filename = format("%s/%s" , path.root, "kube_config_cluster.yml")
//  content = rke_cluster.cluster.kube_config_yaml
//}

###############################################################################
# If you need cluster.yml for using rke, please uncomment follows.
###############################################################################
//resource "local_file" "rke_cluster_yaml" {
//  filename = format("%s/%s", path.root, "cluster.yml")
//  content = rke_cluster.cluster.rke_cluster_yaml
//}

###############################################################################
# You can also use an output.
###############################################################################
// output "rke_cluster_yaml" {
//  sensitive = true
//  value = rke_cluster.cluster.rke_cluster_yaml
//}
```

* default k8s version: `v1.14.3-rancher1-1`
* default network plugin: `canal`

#### Dynamic multiple nodes example

```hcl
variable "nodes" {
  type = list(object({
    address = string,
    user    = string,
  }))
  default = [
    {
      address = "192.2.0.1"
      user    = "ubuntu"
    },
    {
      address = "192.2.0.2"
      user    = "ubuntu"
    },
  ]
}

resource rke_cluster "cluster" {
  dynamic nodes {
    for_each = var.nodes
    content {
      address = nodes.value.address
      user    = nodes.value.user
      role    = ["controlplane", "worker", "etcd"]
      ssh_key = file("~/.ssh/id_rsa")
    }
  }
}

###############################################################################
# If you need kubeconfig.yml for using kubectl, please uncomment follows.
###############################################################################
//resource "local_file" "kube_cluster_yaml" {
//  filename = format("%s/%s" , path.root, "kube_config_cluster.yml")
//  content = rke_cluster.cluster.kube_config_yaml
//}
```

#### With cloud provider

- [AWS(EC2)](examples/aws_ec2)
- [DigitalOcean](examples/digitalocean)
- [OpenStack](https://github.com/mcapuccini/terraform-openstack-rke)

#### With Kubernetes provider 

You can use RKE provider and [Kubernetes provider](https://www.terraform.io/docs/providers/kubernetes/index.html) together.

```hcl
resource rke_cluster "cluster" {
  nodes {
    address = "1.2.3.4"
    user    = "ubuntu"
    role    = ["controlplane", "worker", "etcd"]
    ssh_key = file("~/.ssh/id_rsa")
  }
}

provider "kubernetes" {
  host     = rke_cluster.cluster.api_server_url
  username = rke_cluster.cluster.kube_admin_user

  client_certificate     = rke_cluster.cluster.client_cert
  client_key             = rke_cluster.cluster.client_key
  cluster_ca_certificate = rke_cluster.cluster.ca_crt
  # load_config_file = false
}

resource "kubernetes_namespace" "example" {
  metadata {
    name = "terraform-example-namespace"
  }
}
```

#### Deploying Rancher 2.0 using terraform-provider-rke

You can view examples to deploying Rancher 2.0
 
 - [Use own SSL certificates](examples/rancher_server_minimal/example.tf)
 - [Use SSL-passthrough](examples/rancher_server_ssl_passthrough/example.tf)

#### Full example

You can view full example of tffile, [here](examples/full/example.tf).

## License

 `terraform-provider-rke` Copyright (C) 2018-2019 Kazumichi Yamamoto.

  This project is published under [Apache 2.0 License](LICENSE.txt).
  
## Author

  * Kazumichi Yamamoto ([@yamamoto-febc](https://github.com/yamamoto-febc))
