module "nodes" {
  source = "./aws"
  region        = "eu-central-1"
  # instance_type = "t2.micro"
  # cluster_id    = "rke"
  security_group_name = "rke-default-security-group"
}

resource "rke_cluster" "cluster" {
  cloud_provider {
    name = "aws"
  }

  nodes {
    address = module.nodes.addresses[0]
    internal_address = module.nodes.internal_ips[0]
    user    = module.nodes.ssh_username
    ssh_key = file("mbh.pem")
    role    = ["controlplane", "etcd"]
  }
  nodes {
    address = module.nodes.addresses[1]
    internal_address = module.nodes.internal_ips[1]
    user    = module.nodes.ssh_username
    ssh_key = file("mbh.pem")
    role    = ["worker"]
  }
  nodes {
    address = module.nodes.addresses[2]
    internal_address = module.nodes.internal_ips[2]
    user    = module.nodes.ssh_username
    ssh_key = file("mbh.pem")
    role    = ["worker"]
  }
  nodes {
    address = module.nodes.addresses[3]
    internal_address = module.nodes.internal_ips[3]
    user    = module.nodes.ssh_username
    ssh_key = file("mbh.pem")
    role    = ["worker"]
  }
}

resource "local_file" "kube_cluster_yaml" {
  filename = "./kube_config_cluster.yml"
  content  = rke_cluster.cluster.kube_config_yaml
}

