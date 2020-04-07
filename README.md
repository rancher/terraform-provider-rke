Terraform Provider for RKE
==================================

[![Go Report Card](https://goreportcard.com/badge/github.com/rancher/terraform-provider-rke)](https://goreportcard.com/report/github.com/rancher/terraform-provider-rke)

Terraform RKE providers can easily deploy Kubernetes clusters with [Rancher Kubernetes Engine](https://github.com/rancher/rke).  

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) >= 0.11.x
- [Go](https://golang.org/doc/install) 1.12 to build the provider plugin
- [Docker](https://docs.docker.com/install/) 17.03.x to run acceptance tests

Installing The Provider
-----------------------

Download the binary for your platform from the [releases](https://github.com/rancher/terraform-provider-rke/releases) page

```
wget https://github.com/rancher/terraform-provider-rke/releases/download/x.x.x/terraform-provider-rke_<platform>
```

- rename the `terraform-provider-rke_<OS_ARCH>` to `terraform-provider-rke_v0.14.1_x4`
  - the `v0.14.1` is based on the downloaded version
- put it to `project_directory`.
  - e.g. `Users/$(whoami)/works/src/github.com/quickstart/azure/.terraform/plugins/darwin_amd64`
- run `terraform init`
you should get the following response
```console
> terraform init
Initializing modules...

Initializing the backend...

Initializing provider plugins...

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

If using a custom directory use the `-plug-in=/my/custom/path` when running `terraform init` to specify the path.

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

If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it. Documentation about the provider specific configuration options can be found on the [provider's website](https://www.terraform.io/docs/providers/rke/index.html).

Developing the Provider
-----------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.12+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

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

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, a running rancher system, a rancher API key and a working k8s cluster imported are needed.

To run acceptance tests, you can simply run `make testacc`. `scripts/gotestacc.sh` will be run, deploying all needed requirements, running acceptance tests and cleanup.

```sh
$ make testacc
```

Currently acceptance tests only works on Linux platforms, on Mac OSX are not yet supported.

Provider examples
-----------------

You can view some tf file examples, [here](examples).

On Openstack you can use [terraform-openstack-rke](https://github.com/remche/terraform-openstack-rke) module.
