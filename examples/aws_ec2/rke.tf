module "nodes" {
  source = "./aws"

  # region        = "us-east-1"
  # instance_type = "t2.micro"
  # cluster_id    = "rke"
}

resource rke_cluster "cluster" {
  cloud_provider {
    name = "aws"
  }

  nodes = [
    {
      address = "${module.nodes.addresses[0]}"
      user    = "${module.nodes.ssh_username}"
      ssh_key = "${module.nodes.private_key}"
      role    = ["controlplane", "etcd"]
    },
    {
      address = "${module.nodes.addresses[1]}"
      user    = "${module.nodes.ssh_username}"
      ssh_key = "${module.nodes.private_key}"
      role    = ["worker"]
    },
    {
      address = "${module.nodes.addresses[2]}"
      user    = "${module.nodes.ssh_username}"
      ssh_key = "${module.nodes.private_key}"
      role    = ["worker"]
    },
    {
      address = "${module.nodes.addresses[3]}"
      user    = "${module.nodes.ssh_username}"
      ssh_key = "${module.nodes.private_key}"
      role    = ["worker"]
    },
  ]
}

resource "local_file" "kube_cluster_yaml" {
  filename = "./kube_config_cluster.yml"
  content  = "${rke_cluster.cluster.kube_config_yaml}"
}