Terraform Provider for RKE
==================================

[![Go Report Card](https://goreportcard.com/badge/github.com/rancher/terraform-provider-rke)](https://goreportcard.com/report/github.com/rancher/terraform-provider-rke) [![Build Status](https://drone-publish.rancher.io/api/badges/rancher/terraform-provider-rke/status.svg)](https://drone-publish.rancher.io/rancher/terraform-provider-rke)

Terraform RKE providers can easily deploy Kubernetes clusters with [Rancher Kubernetes Engine](https://github.com/rancher/rke).  

- Website: https://registry.terraform.io/providers/rancher/rke
- Docs: https://registry.terraform.io/providers/rancher/rke/latest/docs
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) >= 0.12.x
- [Go](https://golang.org/doc/install) 1.14 to build the provider plugin
- [Docker](https://docs.docker.com/install/) 17.03.x to run acceptance tests

Installing The Provider
-----------------------

For terraform 0.13 or above users, manual installation is not required anymore. Provider will be downloaded by `terraform init` from [terraform rke registry](https://registry.terraform.io/providers/rancher/rke). 

For upgrade terraform to 0.13, please follow [upgrade_to_0.13 guide](https://registry.terraform.io/providers/rancher/rke/latest/docs/guides/upgrade_to_0.13)

**Note:** If you are using terraform 0.12 or lower, the provider should be installed manually. Downloadable packages are available at [rke provider releases](https://github.com/rancher/terraform-provider-rke/releases)

How to install manually:
* Get binary for your platform: `curl -L https://github.com/rancher/terraform-provider-rke/releases/download/vX.Y.Z/terraform-provider-rke_X.Y.Z_<OS>_<ARCH>.zip | unzip -`
* Grant execution permission: `chmod 755 terraform-provider-rke_vX.Y.Z`
* Place provider binary on your terraform plugin directory: `mv terraform-provider-rke_vX.Y.Z <TERRAFORM_PLUGIN_DIRECTORY>`
  * `terraform init` will search the following terraform plugin directories (More info at [terra-farm.github.io](https://terra-farm.github.io/main/installation.html)):

| Directory | Purpose |
|-|-|
| . | In case the provider is only used in a single Terraform project. |
| Location of the terraform binary (/usr/local/bin, for example.) | For airgapped installations; see terraform bundle. |
| terraform.d/plugins/<OS>_<ARCH> | For checking custom providers into a configuration’s VCS repository. Not usually desirable, but sometimes necessary in Terraform Enterprise. |
| .terraform/plugins/<OS>_<ARCH> | Automatically downloaded providers. |
| ~/.terraform.d/plugins | The user plugins directory. |
| ~/.terraform.d/plugins/<OS>_<ARCH> | The user plugins directory, with explicit OS and architecture. |
| /my/custom/path | Custom plugin directory. Use the `-plug-dir=/my/custom/path` when running `terraform init` |

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-rke`

```sh
$ mkdir -p $GOPATH/src/github.com/rancher
$ cd $GOPATH/src/github.com/rancher

$ go get github.com/rancher/terraform-provider-rke
$ go install github.com/rancher/terraform-provider-rke
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/rancher/terraform-provider-rke
$ make build
```

**Current master is focusing on RKE v1.0.x** There are some breaking changes from previous provider version.

**If you use RKE v0.2.x or v0.1.x, please set proper branch.**

Using the Provider
------------------

If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it. Documentation about the provider specific configuration options can be found at [rke provider docs](https://registry.terraform.io/providers/rancher/rke/latest/docs).

Developing the Provider
-----------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.13+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in `$GOPATH/bin` .

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-rke
...
```

To just compile provider binary on repo path and test on terraform:

```sh
$ make bin
$ terraform init
$ terraform plan
$ terraform apply
```

Testing the Provider
--------------------

In order to test structures on the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, you can simply run `make testacc`. `scripts/gotestacc.sh` will be run, deploying all needed requirements, running acceptance tests and cleanup.

```sh
$ make testacc
```

Due to network limitations on docker for Mac OSX and windows, on these platforms acceptance tests should run on a dockerized enviroment.

```sh
$ make dapper-testacc
```

Provider examples
-----------------

You can view some tf file examples, [here](examples).

On Openstack you can use [terraform-openstack-rke](https://github.com/remche/terraform-openstack-rke) module.

Branching the Provider
---------------------------

The provider is branched into two release lines `master` and `release/v1.3` that align with the minor versions of RKE. `master` is aligned with RKE v1.4 and `release/v1.3` is aligned with RKE v1.3.

Release process

* When there's an RKE/KDM release and the RKE release has new kubernetes versions, create an issue that will list the RKE version + provide links
* Backport CVEs or fixes as needed to `release/v1.3`
* Release the provider that includes the new version/s
* It depends on what RKE version just got released (v1.3 / v1.4) as to what line of the provider gets released. Could be one or the other, or both.

Releasing the Provider
---------------------------

* Create a draft of the [release](https://github.com/rancher/terraform-provider-rke/releases) and select create new tag for the version you are releasing
* Create release notes by clicking `Generate release notes`
* Copy the release notes to the CHANGELOG and update to the following format

```
# <tag version> (Month Day, Year)
FEATURES:
ENHANCEMENTS:
BUG FIXES:
```

* Create a PR to update CHANGELOG
* Copy the updated notes back to the draft release (DO NOT release with just the generated notes. Those are just a template to help you)
* Make sure the branch is up-to-date with the remote, in this example, the branch is master and the release tag is v1.24.0

```
git remote add upstream-release git@github.com:rancher/terraform-provider-rke.git
git checkout upstream-release/master
git push upstream-release v1.24.0
```
* Check that the new version has been published on the [Hashicorp registry](https://registry.terraform.io/providers/rancher/rke/latest). If it has not been published after 2 hours, create an [EIO issue](https://github.com/rancherlabs/eio) to trigger the registry sync.
