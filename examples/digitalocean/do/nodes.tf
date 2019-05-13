# Droplet(ノード)の定義
resource "digitalocean_droplet" "rke-node" {
  image    = "ubuntu-18-04-x64"
  name     = "rke-nodes-${count.index + 1}"
  region   = var.region
  size     = var.droplet_size
  ssh_keys = [digitalocean_ssh_key.key.id]
  count    = 4

  private_networking = true

  provisioner "remote-exec" {
    connection {
      type        = "ssh"
      user        = "root"
      host        = self.ipv4_address
      private_key = tls_private_key.node-key.private_key_pem
    }

    inline = [
      "curl releases.rancher.com/install-docker/18.09.4.sh | bash",
    ]
  }
}

