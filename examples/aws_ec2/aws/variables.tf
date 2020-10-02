variable "region" {
  default = "us-east-1"
}

variable "instance_type" {
  default = "t3.large"
}

variable "cluster_id" {
  default = "rke"
}

variable "docker_install_url" {
  default = "https://releases.rancher.com/install-docker/19.03.sh"
}
