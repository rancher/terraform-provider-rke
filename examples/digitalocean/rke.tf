### Example works for RKE v1.13.5-rancher1-2

variable "do_token" {
  default = ""
}

module "nodes" {
  source = "./do"
  do_token = var.do_token
  # region       = "nyc1"
  # droplet_size = "t2.micro"
}

resource "rke_cluster" "cluster" {
  nodes {
    internal_address = module.nodes.internal_addresses[0]
    address          = module.nodes.addresses[0]
    user             = module.nodes.ssh_username
    ssh_key          = module.nodes.private_key
    role             = ["controlplane", "etcd"]
  }
  nodes {
    internal_address = module.nodes.internal_addresses[1]
    address          = module.nodes.addresses[1]
    user             = module.nodes.ssh_username
    ssh_key          = module.nodes.private_key
    role             = ["worker"]
  }
  nodes {
    internal_address = module.nodes.internal_addresses[2]
    address          = module.nodes.addresses[2]
    user             = module.nodes.ssh_username
    ssh_key          = module.nodes.private_key
    role             = ["worker"]
  }
  nodes {
    internal_address = module.nodes.internal_addresses[3]
    address          = module.nodes.addresses[3]
    user             = module.nodes.ssh_username
    ssh_key          = module.nodes.private_key
    role             = ["worker"]
  }

  services {
    kube_api {
      extra_args = {
        kubelet-preferred-address-types = "InternalIP,ExternalIP,Hostname"
        feature-gates = "VolumeSnapshotDataSource=true,KubeletPluginsWatcher=true,CSINodeInfo=true,CSIDriverRegistry=true"
      }
    }
    kubelet {
      extra_args = {
        cloud-provider = "external"
        feature-gates = "VolumeSnapshotDataSource=true,KubeletPluginsWatcher=true,CSINodeInfo=true,CSIDriverRegistry=true"
      }
    }
  }

  ingress {
    provider = "none"
  }
  
  addon_job_timeout = 60
  addons = "${data.template_file.addons.rendered}"

  addons_include = [
    # Cloud Controller Manager for DigitalOcean
    "${path.module}/files/ccm-digitalocean-v0.1.15.yaml",
    # CSI driver for DO, only v1.0.1 version works with rke v1.13.5-rancher1-2
    "${path.module}/files/csi-digitalocean-v1.0.1.yaml",
    # Ingess-nginx for generic cloud (with LoadBalancer type service)
    "${path.module}/files/ingress-mandatory.yaml",
    "${path.module}/files/ingress-cloud-generic.yaml"

  ]
}

resource "local_file" "kube_cluster_yaml" {
  filename = "./kube_config_cluster.yml"
  content  = rke_cluster.cluster.kube_config_yaml
}

