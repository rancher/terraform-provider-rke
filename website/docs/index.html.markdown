---
layout: "rke"
page_title: "Provider: RKE"
sidebar_current: "docs-rke-index"
description: |-
  The RKE provider is used to interact with Rancher Kubernetes Engine kubernetes clusters.
---

# RKE Provider

The RKE provider is used to interact with Rancher Kubernetes Engine kubernetes clusters.

## Example Usage

```hcl
# Configure the RKE provider
provider "rke" {
  log = true
}
```

## Argument Reference

The following arguments are supported:

* `log` - (Optional) Enable RKE logs. It can also be sourced from the `RKE_LOG` environment variable. Default `false` (bool)
