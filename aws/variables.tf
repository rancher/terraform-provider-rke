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
  default = "https://releases.rancher.com/install-docker/19.03.sh"
}

variable "instance_count" {
  default = 5
}

variable "security_group_name" {
  description = "Name of existing security group to use for rke ec2 instances"
  default     = "allow-all"
}

variable "instance_profile_name" {
  description = "Name of the IAM instance profile that EC2 instances will use"
  default     = "rke-aws"
}

