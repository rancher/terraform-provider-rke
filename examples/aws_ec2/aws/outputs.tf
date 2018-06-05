output "private_key" {
  value = "${tls_private_key.node-key.private_key_pem}"
}

output "ssh_username" {
  value = "ubuntu"
}

output "addresses" {
  value = ["${aws_instance.rke-node.*.public_dns}"]
}
