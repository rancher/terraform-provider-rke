# Droplet(ノード)の定義
resource "digitalocean_droplet" "rke-node" {
  image    = "rancheros"
  name     = "rke-nodes-${count.index + 1}"
  region   = var.region
  size     = var.droplet_size
  ssh_keys = [digitalocean_ssh_key.key.id]
  count    = 4

  private_networking = true

  provisioner "remote-exec" {
    connection {
      type        = "ssh"
      user        = "rancher"
      host        = self.ipv4_address
      private_key = tls_private_key.node-key.private_key_pem
    }

    inline = [
      "while ! docker ps > /dev/null ; do sleep 2 ; done",
    ]
  }
}

