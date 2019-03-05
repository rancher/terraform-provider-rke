resource "rke_cluster" "cluster" {
  nodes {
    address = "1.2.3.4"
    user    = "ubuntu"
    role    = ["controlplane", "worker", "etcd"]
    ssh_key = file("~/.ssh/id_rsa")
  }
}

###############################################################################
# If you need kubeconfig.yml for using kubectl, please uncomment follows.
###############################################################################
// resource "local_file" "kube_cluster_yaml" {
//   filename = "${path.root}/kube_config_cluster.yml"
//   content  = rke_cluster.cluster.kube_config_yaml
// }

