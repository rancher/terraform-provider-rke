### Example works for RKE v1.19.4-rancher1-1

# Your DigitalOcean API token; typically specified at runtime
variable "do_token" {
  default = ""
}

# The DigitalOcean region where the RKE cluster should be deployed
variable "region" {
  default = "nyc1"
}

# The domain name of your RKE cluster
variable "cluster_domain" {
  default = "cluster.local"
}

# The size of the DigitalOcean droplets for the controlplane and worker nodes in the RKE cluster
variable "droplet_size" {
  type    = map
  default = {
    "controlplane" = "s-2vcpu-4gb"
    "worker"       = "s-1vcpu-2gb"
  }
}

# The number of DigitalOcean droplets that should be provisioned as controlplane and worker nodes
variable "node_count" {
  type    = map
  default = {
    "controlplane" = "3"
    "worker"       = "2"
  }
}

# Define two local variables using the Terraform range function
# This will enable us to leverage dynamic blocks when defining the controlplane and worker nodes
locals {
  nodes_controlplane_count = range(var.node_count["controlplane"])
  nodes_worker_count = range(var.node_count["worker"])
}

# Create a DigitalOcean SSH key to deploy and interface with the RKE cluster
module "keys" {
  source   = "./do/keys"
  do_token = var.do_token
}

# Create the DigitalOcean droplets to provision as RKE controlplane nodes
module "nodes_controlplane" {
  name           = "controlplane"
  source         = "./do/nodes"
  do_token       = var.do_token
  region         = var.region
  droplet_size   = var.droplet_size["controlplane"]
  do_key_id      = module.keys.do_key_id
  private_key    = module.keys.private_key
  cluster_domain = var.cluster_domain
  node_count     = var.node_count["controlplane"]
}

# Create the DigitalOcean droplets to provision as RKE worker nodes
module "nodes_worker" {
  name           = "worker"
  source         = "./do/nodes"
  do_token       = var.do_token
  region         = var.region
  droplet_size   = var.droplet_size["worker"]
  do_key_id      = module.keys.do_key_id
  private_key    = module.keys.private_key
  cluster_domain = var.cluster_domain
  node_count     = var.node_count["worker"]
}

# Once droplets are provisioned, deploy RKE
resource "rke_cluster" "cluster" {
  # Uncomment to define a specific RKE version to deploy to the cluster
#  kubernetes_version = "v1.19.4-rancher1-1"

  # Define the RKE controlplane nodes dynamically
  dynamic "nodes" {
    for_each = local.nodes_controlplane_count
    iterator = node_id
    content {
      node_name         = module.nodes_controlplane.names[node_id.value]
      internal_address  = module.nodes_controlplane.internal_addresses[node_id.value]
      address           = module.nodes_controlplane.addresses[node_id.value]
      user              = module.nodes_controlplane.ssh_username
      ssh_key           = module.keys.private_key
      hostname_override = module.nodes_controlplane.names[node_id.value]
      role              = ["controlplane", "etcd", "worker"]
    }
  }

  # Define the RKE worker nodes dynamically
  dynamic "nodes" {
    for_each = local.nodes_worker_count
    iterator = node_id
    content {
      node_name         = module.nodes_worker.names[node_id.value]
      internal_address  = module.nodes_worker.internal_addresses[node_id.value]
      address           = module.nodes_worker.addresses[node_id.value]
      user              = module.nodes_worker.ssh_username
      ssh_key           = module.keys.private_key
      hostname_override = module.nodes_worker.names[node_id.value]
      role              = ["worker"]
    }
  }

  # Configure the RKE services
  services {
    kube_api {
      # The extra_args below are required to support the DigitalOcean CSI and CCM resources that
      # will be applied as addons to the cluster
      extra_args = {
        allow-privileged = "true"
        kubelet-preferred-address-types = "InternalIP,ExternalIP,Hostname"
        feature-gates = "VolumeSnapshotDataSource=true,CSINodeInfo=true,CSIDriverRegistry=true"
      }
    }
    kubelet {
      cluster_domain = var.cluster_domain
      # The extra_args below are required to support the DigitalOcean CSI and CCM resources that
      # will be applied as addons to the cluster
      extra_args = {
        cloud-provider = "external"
        feature-gates = "VolumeSnapshotDataSource=true,CSINodeInfo=true,CSIDriverRegistry=true"
      }
    }
  }

  # The NGINX Ingress Controller will be provisioned as an addon (below), using the DigitalOcean-specific
  # Kubernetes resources, so we should bypass deploying the native RKE ingress controller
  ingress {
    provider = "none"
  }
  
  addon_job_timeout = 60
  addons = data.template_file.addons.rendered

  # Provision Kubernetes resources as addons to support RKE deployment in DigitalOcean
  addons_include = [
    # Cloud Controller Manager for DigitalOcean
    # Source: https://github.com/digitalocean/digitalocean-cloud-controller-manager/blob/v0.1.30/releases/v0.1.30.yml
    "${path.module}/files/ccm-digitalocean-v0.1.30.yaml",
    # CSI driver for DigitalOcean
    # Source: https://github.com/digitalocean/csi-digitalocean/tree/v2.1.1/deploy/kubernetes/releases/csi-digitalocean-v2.1.1
    "${path.module}/files/csi-digitalocean-v2.1.1/crds.yaml",
    "${path.module}/files/csi-digitalocean-v2.1.1/driver.yaml", # modified for RKE
    "${path.module}/files/csi-digitalocean-v2.1.1/snapshot-controller.yaml",
    # Ingess-nginx for DigitalOcean v0.42.0
    # Source: https://github.com/kubernetes/ingress-nginx/blob/controller-v0.42.0/deploy/static/provider/do/deploy.yaml
    "${path.module}/files/ingress-nginx-deploy-v0.42.0.yaml"
  ]
}

# Dump the kubeconfig
resource "local_file" "kube_cluster_yaml" {
  filename = "./kube_config_cluster.yml"
  content  = rke_cluster.cluster.kube_config_yaml
}
