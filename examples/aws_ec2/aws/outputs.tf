output "ssh_username" {
  value = "ubuntu"
}

output "addresses" {
  value = aws_instance.rke-node[*].public_dns
}

output "internal_ips" {
  value = aws_instance.rke-node[*].private_ip
}
