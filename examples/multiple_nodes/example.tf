variable "node_addrs" {
  type    = "list"
  default = ["192.2.0.1", "192.2.0.2"]
}

data rke_node_parameter "nodes" {
  count   = "${length(var.node_addrs)}"

  address = "${var.node_addrs[count.index]}"
  user    = "ubuntu"
  role    = ["controlplane", "worker", "etcd"]
  ssh_key = "${file("~/.ssh/id_rsa")}"
}

resource rke_cluster "cluster" {
  nodes_conf = ["${data.rke_node_parameter.nodes.*.json}"]
}

###############################################################################
# If you need kubeconfig.yml for using kubectl, please uncomment follows.
###############################################################################
//resource "local_file" "kube_cluster_yaml" {
//  filename = "${path.root}/kube_config_cluster.yml"
//  content = "${rke_cluster.cluster.kube_config_yaml}"
//}

