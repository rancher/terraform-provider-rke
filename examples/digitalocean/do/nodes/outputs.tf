output "ssh_username" {
  value = var.ssh_username
}

output "names" {
  value = digitalocean_droplet.rke-node[*].name
}

output "addresses" {
  value = digitalocean_droplet.rke-node[*].ipv4_address
}

output "internal_addresses" {
  value = digitalocean_droplet.rke-node.*.ipv4_address_private
}
