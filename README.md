Terraform Provider for Google Cloud Platform
==================

- Website: https://www.terraform.io
- Documentation: https://www.terraform.io/docs/providers/google/index.html
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)
<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Maintainers
-----------

This provider plugin is maintained by:

* The [Google Cloud Graphite Team](https://cloudplatform.googleblog.com/2017/03/partnering-on-open-source-Google-and-HashiCorp-engineers-on-managing-GCP-infrastructure.html) at Google
* The Terraform team at [HashiCorp](https://www.hashicorp.com/)

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.12+


Using the provider
----------------------

See the [Google Provider documentation](https://www.terraform.io/docs/providers/google/index.html) to get started using the
Google provider.

We also maintain the 'google-beta' provider for preview features and features at a beta [launch stage](https://cloud.google.com/products#product-launch-stages). See [Provider Versions](https://www.terraform.io/docs/providers/google/provider_versions.html)
for more details on how to use `google-beta`.

Upgrading the provider
----------------------

The Google provider doesn't upgrade automatically once you've started using it. After a new release you can run

```bash
terraform init -upgrade
```

to upgrade to the latest stable version of the Google provider. See the [Terraform website](https://www.terraform.io/docs/configuration/providers.html#provider-versions)
for more information on provider upgrades, and how to set version constraints on your provider.

Building the provider
---------------------

Clone repository to: `$GOPATH/src/github.com/hashicorp/terraform-provider-google`

```sh
$ mkdir -p $GOPATH/src/github.com/hashicorp; cd $GOPATH/src/github.com/hashicorp
$ git clone git@github.com:hashicorp/terraform-provider-google
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/hashicorp/terraform-provider-google
$ make build
```

Developing the provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org)
installed on your machine (version 1.16.0+ is *required*). You can use [goenv](https://github.com/syndbg/goenv)
to manage your Go version. You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH),
as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`.
This will build the provider and put the provider binary in the `$GOPATH/bin`
directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-google
...
```

For guidance on common development practices such as testing changes or
vendoring libraries, see the [contribution guidelines](https://github.com/hashicorp/terraform-provider-google/blob/master/.github/CONTRIBUTING.md).
If you have other development questions we don't cover, please file an issue!
