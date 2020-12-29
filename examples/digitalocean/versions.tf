terraform {
  required_providers {
    local = {
      source = "hashicorp/local"
    }
    rke = {
      source = "rancher/rke"
    }
    template = {
      source = "hashicorp/template"
    }
  }
  required_version = ">= 0.13"
}
