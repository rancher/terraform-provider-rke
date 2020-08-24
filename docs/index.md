---
page_title: "RKE Provider"
---

# RKE Provider

The RKE provider is used to interact with Rancher Kubernetes Engine kubernetes clusters.

## Example Usage

```hcl
# Configure the RKE provider
provider "rke" {
  debug = true
  log_file = "<RKE_LOG_FILE>"
}
```

## Argument Reference

The following arguments are supported:

* `debug` - (Optional) Enable RKE debug logs. It can also be sourced from the `RKE_DEBUG` environment variable. Default `false` (bool)
* `log_file` - (Optional) Save RKE logs to a file. It can also be sourced from the `RKE_LOG_FILE` environment variable (string)
