provider "aws" {
  region = var.region
}
terraform {
  required_providers {
    rke = {
      source = "rancher/rke"
      version = "1.1.2"
    }
  }
}

