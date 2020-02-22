terraform {
  backend "s3" {
    bucket = "rancher-mhassine-rampup"
    key    = "rke-provider-example.state"
    region = "eu-central-1"
  }
}