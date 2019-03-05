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
# resource "local_file" "kube_cluster_yaml" {
#   filename = "${path.root}/kube_config_cluster.yml"
#   content  = rke_cluster.cluster.kube_config_yaml
# }
