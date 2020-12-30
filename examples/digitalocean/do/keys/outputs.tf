output "do_key_id" {
  value = digitalocean_ssh_key.key.id
}

output "private_key" {
  value = tls_private_key.node-key.private_key_pem
}
