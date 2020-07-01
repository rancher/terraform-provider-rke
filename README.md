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
- [Go](https://golang.org/doc/install) 1.13 to build the provider plugin
- [Docker](https://docs.docker.com/install/) 17.03.x to run acceptance tests

Installing The Provider
-----------------------

**Note:** If you are testing [terraform 0.13-beta](https://www.hashicorp.com/blog/announcing-the-terraform-0-13-beta/), manual installation is not required. Provider will be downloaded by `terraform init` from [terraform rke registry](https://registry.terraform.io/providers/rancher/rke).

By the moment, the provider should be installed manually. Downloadable packages are available at [rke provider releases](https://github.com/rancher/terraform-provider-rke/releases)

How to install:
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
| /my/custom/path | Custom plugin directory. Use the `-plug-in=/my/custom/path` when running `terraform init` |

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

Using the provider
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
