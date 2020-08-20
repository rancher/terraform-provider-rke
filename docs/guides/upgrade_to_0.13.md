---
page_title: "Upgrade to terraform 0.13"
---

# Upgrade to terraform 0.13 

RKE provider is already published at [rke terraform registry](https://registry.terraform.io/providers/rancher/rke) as verified provider. That means that the provider is compatible with terraform 0.13 and is automatically installed by it. As RKE provider had to be manually installated on terraform prior to 0.13, [in-house step](https://www.terraform.io/upgrade-guides/0-13.html#in-house-providers) is required on the terraform upgrade process.

## Steps

These are the steps to properly update RKE provider to tf 0.13:

1. Before start, be sure that `terraform apply` doesn't show any diff
2. Update `version.tf` file adding provider definition with required version

```
terraform {
  required_providers {
...
    rke = {
      source  = "rancher/rke"
      version = "1.1.0"
    }
  }
...
}
```

3. Execute `terraform 0.13upgrade`

```
$ terraform 0.13upgrade

This command will update the configuration files in the given directory to use
the new provider source features from Terraform v0.13. It will also highlight
any providers for which the source cannot be detected, and advise how to
proceed.

We recommend using this command in a clean version control work tree, so that
you can easily see the proposed changes as a diff against the latest commit.
If you have uncommited changes already present, we recommend aborting this
command and dealing with them before running this command again.

Would you like to upgrade the module in the current directory?
  Only 'yes' will be accepted to confirm.

  Enter a value: yes

-----------------------------------------------------------------------------

Upgrade complete!

Use your version control system to review the proposed changes, make any
necessary adjustments, and then commit.
```

4. Replace previous in house provider definition on tfstate. tfstate will be updated, backup before proceed is recommended. [More info](https://www.terraform.io/upgrade-guides/0-13.html#in-house-providers)

```
$ terraform  state replace-provider 'registry.terraform.io/-/rke' 'registry.terraform.io/rancher/rke'
Terraform will perform the following actions:

  ~ Updating provider:
    - registry.terraform.io/-/rke
    + registry.terraform.io/rancher/rke

Changing 1 resources:

  rke_cluster.cluster

Do you want to make these changes?
Only 'yes' will be accepted to continue.

Enter a value: yes

Successfully replaced provider for 1 resources.
```

5. Init the provider

```
$ terraform init

Initializing the backend...

Initializing provider plugins...
- Finding rancher/rke versions matching "1.1.0"...
- Installing rancher/rke v1.1.0...
- Installed rancher/rke v1.1.0 (signed by a HashiCorp partner, key ID 2EEB0F9AD44A135C)

Partner and community providers are signed by their developers.
If you'd like to know more about provider signing, you can read about it here:
https://www.terraform.io/docs/plugins/signing.html

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

More info at [terraform upgrade guide to 0.13](https://www.terraform.io/upgrade-guides/0-13.html)
