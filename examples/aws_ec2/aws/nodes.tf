data "aws_availability_zones" "az" {
}


data "aws_security_group" "selected" {
  name = var.security_group_name
}

# resource "aws_subnet" "subnet" {
#   vpc_id     = data.aws_security_group.selected.vpc_id
#   cidr_block = "10.0.1.0/24"
# }

# resource "aws_default_subnet" "default" {
#   availability_zone = data.aws_availability_zones.az.names[count.index]
#   tags              = local.cluster_id_tag
#   count             = length(data.aws_availability_zones.az.names)
# }

# resource "aws_security_group" "allow-all" {
#   name        = "rke-default-security-group"
#   description = "rke"

#   ingress {
#     from_port   = 0
#     to_port     = 0
#     protocol    = "-1"
#     cidr_blocks = ["0.0.0.0/0"]
#   }

#   egress {
#     from_port   = 0
#     to_port     = 0
#     protocol    = "-1"
#     cidr_blocks = ["0.0.0.0/0"]
#   }

#   tags = local.cluster_id_tag
# }

resource "aws_instance" "rke-node" {
  count = var.instance_count

  ami                    = data.aws_ami.ubuntu.id
  instance_type          = var.instance_type
  key_name               = "mbh"
  iam_instance_profile   = var.instance_profile_name
  vpc_security_group_ids = [data.aws_security_group.selected.id]
  tags                   = {  
    "kubernetes.io/cluster/${var.cluster_id}" = "owned"
    Name = "rke_0${count.index+1}"
  }
  provisioner "remote-exec" {
    connection {
      host        = coalesce(self.public_ip, self.private_ip)
      type        = "ssh"
      user        = "ubuntu"
      private_key = file("${path.module}/../mbh.pem")
    }

    inline = [
      "curl ${var.docker_install_url} | sh",
      "sudo usermod -a -G docker ubuntu",
    ]
  }
}

