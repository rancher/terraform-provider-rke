resource rke_cluster "cluster" {
  nodes = [
    {
      address = "1.2.3.4"
      user    = "ubuntu"
      role    = ["controlplane", "worker", "etcd"]
      ssh_key = "${file("~/.ssh/id_rsa")}"
    },
  ]
}

provider "kubernetes" {
  host     = "${rke_cluster.cluster.api_server_url}"
  username = "${rke_cluster.cluster.kube_admin_user}"

  client_certificate     = "${rke_cluster.cluster.client_cert}"
  client_key             = "${rke_cluster.cluster.client_key}"
  cluster_ca_certificate = "${rke_cluster.cluster.ca_crt}"
}

resource "kubernetes_namespace" "example" {
  metadata {
    name = "terraform-example-namespace"
  }
}
