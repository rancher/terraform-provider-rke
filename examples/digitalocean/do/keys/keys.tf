resource "tls_private_key" "node-key" {
  algorithm = "RSA"
}

resource "digitalocean_ssh_key" "key" {
  name       = "rke-node-key"
  public_key = tls_private_key.node-key.public_key_openssh
}

