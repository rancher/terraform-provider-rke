variable "region" {
  default = "us-east-1"
}

variable "instance_type" {
  default = "t2.micro"
}

variable "cluster_id" {
  default = "rke"
}

variable "docker_install_url" {
  default = "https://releases.rancher.com/install-docker/18.09.sh"
}
