terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
    }
    tls = {
      source = "hashicorp/tls"
    }
  }
  required_version = ">= 0.13"
}
