output "private_key" {
  value = tls_private_key.node-key.private_key_pem
}

output "ssh_username" {
  value = "root"
}

output "addresses" {
  value = digitalocean_droplet.rke-node[*].ipv4_address
}

output "internal_addresses" {
  value = digitalocean_droplet.rke-node.*.ipv4_address_private
}