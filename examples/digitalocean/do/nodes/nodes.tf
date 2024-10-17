# Droplet(ノード)の定義
resource "digitalocean_droplet" "rke-node" {
  image    = "rancheros"
  name     = "${var.name}${count.index + 1}.${var.cluster_domain}"
  region   = var.region
  size     = var.droplet_size
  ssh_keys = [var.do_key_id]
  count    = var.node_count

  private_networking = true

  provisioner "remote-exec" {
    connection {
      type        = "ssh"
      user        = var.ssh_username
      host        = self.ipv4_address
      private_key = var.private_key
    }

    inline = [
      "while ! docker ps > /dev/null ; do sleep 2 ; done",
    ]
  }
}

