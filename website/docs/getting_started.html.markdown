---
layout: "google"
page_title: "Getting Started with the Google provider"
sidebar_current: "docs-google-provider-getting-started"
description: |-
  Getting started with the Google Cloud Platform provider
---

# Getting Started with the Google Provider

## Before you begin

* You'll want to have created an account and project in the [Google Cloud Console](https://console.cloud.google.com/)
and set up billing on that project. Any examples in this guide will be part of
the [GCP "always free" tier](https://cloud.google.com/free/).
* [Install Terraform](https://www.terraform.io/intro/getting-started/install.html)
and read the Terraform getting started guide that follows. This guide will
assume basic proficiency with Terraform - it is an introduction to the Google
provider.

## Configuring the Provider

To configure the provider, first create a Terraform config file. Inside, you'll
want to include the following configuration:

```hcl
provider "google" {
  project = "{{YOUR GCP PROJECT}}"
  region  = "us-central1"
  zone    = "us-central1-c"
}
```

* The `project` field should include your personal project id. The `project`
indicates the default GCP project all of your resources will be created in.
Most Terraform resources will have a `project` field.
* The`region` and `zone` are [locations](https://cloud.google.com/compute/docs/regions-zones/global-regional-zonal-resources)
for your resources to be created in.
    * The `region` will be used to choose the default location for regional
    resources. Regional resources are spread across several zones.
    * The `zone` will be used to choose the default location for zonal resources.
    Zonal resources exist in a single zone. All zones are a part of a region.

Not all resources require a location. Some GCP resources are global and are
automatically spread across all of GCP.

-> Want to try out another location? Check out the [list of available regions and zones](https://cloud.google.com/compute/docs/regions-zones/#available).
Instances created in zones outside the US are not part of the always free tier
and could incur charges.

## Creating a VM instance
A [Google Compute Engine VM instance](https://cloud.google.com/compute/docs/instances/) is
named `google_compute_instance` in Terraform. The `google` part of the name
identifies the provider for Terraform, `compute` indicates the GCP product
family, and `instance` is the resource name.

Google provider resources will generally, although not always, be named after
the name used in the REST API. For example, a VM instance is called [`instance` in the API](https://cloud.google.com/compute/docs/reference/rest/v1/instances).
Most resource field names will also correspond 1:1 with their REST API names.

If you look at the [`google_compute_instance documentation`](/docs/providers/google/r/compute_instance.html),
you'll see that `project` and `zone` (VM instances are a zonal resource) are
listed as optional. You can omit them from a Terraform config and the provider
defaults will be used. If any of these fields is added to a config, it will be
used instead of the default.

Try adding the following to your state file:

```hcl
resource "google_compute_instance" "vm_instance" {
  name         = "terraform-instance"
  machine_type = "f1-micro"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  network_interface {
    # A default network is created for all GCP projects
    network       = "default"
    access_config = {
    }
  }
}
```

~> Note: Don't use `terraform apply` quite yet! You still need to add GCP
credentials below. If you want to try out provisioning your VM instance before
continuing, follow the instructions in the "Adding credentials" section below.

## Connecting to a VPC network

Like this VM instance, nearly every GCP resource will have a `name` field. They
are used as a short way to identify resources, and their name in the Cloud
Console will be the one defined in the `name` field.

When linking resources in a Terraform config though, you'll primarily want to
use a different field, the `self_link` of a resource. Like `name`, nearly every
resource has a `self_link`. They look like:

```
{{API base url}}/projects/{{your project}}/{{location type}}/{{location}}/{{resouce type}}/{{name}}
```

For a real example, let's look at what the instance `self_link` will be assuming
your project is named `foo`:

```
https://www.googleapis.com/compute/v1/projects/foo/zones/us-central1-c/instances/terraform-instance
```

A resource's `self_link` is a unique reference to that resource. When
connecting two resources in Terraform, you can use Terraform interpolation to
avoid typing all that out! Let's use a `google_compute_network` to demonstrate.

Add this block to your config:

```hcl
resource "google_compute_network" "vpc_network" {
  name                    = "terraform-network"
  auto_create_subnetworks = "true"
}
```

This will create [VPC network resource](/docs/providers/google/r/compute_network.html)
with a subnetwork (due to auto_create_subnetworks) in each region. Then, change
the network of the `google_compute_instance`.

```diff
network_interface {
-  # A default network is created for all GCP projects
-  network = "default"
+  network = "${google_compute_network.vpc_network.self_link}"
  access_config = {
```

This means that when we create the VM instance, it will use
`"terraform-network"` instead of the default VPC network for the project. If you
run `terraform plan`, you should see that `"terraform-instance"` depends on
`"terraform-network"`.

You're ready to run `terraform apply`! Except for one last thing...

~> You haven't added GCP API credentials yet.

## Adding credentials
In order to make requests against the GCP API, you need to authenticate to prove
that it's you making the request. The preferred method of provisioning resources
with Terraform is to use a [GCP service account](https://cloud.google.com/docs/authentication/getting-started),
a "robot account" that can be granted a limited set of IAM permissions.

From [the service account key page in the Cloud Console](https://pantheon.corp.google.com/apis/credentials/serviceaccountkey)
choose an existing account, or create a new one. Then, download the JSON key
file. Name it something you can remember, and store it somewhere secure on your
machine.

You supply the key to Terraform using the environment variable
`GOOGLE_APPLICATION_CREDENTIALS`, setting the value to the location of the file.

```bash
export GOOGLE_APPLICATION_CREDENTIALS={{path}}
```

-> Remember to add this line to a startup file such as `bash_profile` or
`bashrc` to store your credentials across sessions!

## Provisioning your resources
By now, your config should look something like:

```hcl
provider "google" {
  project = "{{YOUR GCP PROJECT}}"
  region  = "us-central1"
  zone    = "us-central1-c"
}

resource "google_compute_instance" "vm_instance" {
  name         = "terraform-instance"
  machine_type = "f1-micro"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  network_interface {
    # A default network is created for all GCP projects
    network       = "${google_compute_network.vpc_network.self_link}"
    access_config = {
    }
  }
}

resource "google_compute_network" "vpc_network" {
  name                    = "terraform-network"
  auto_create_subnetworks = "true"
}
```

With a Terraform config and with your credentials configured, it's time to
provision your resources:

```hcl
terraform apply
```

Congratulations! You've gotten started using the Google provider and provisioned
a virtual machine on Google Cloud Platform. Along the way, you've learned
important concepts such as:

* How a `project` contains resources
    * and how to use a default `project` in your provider
* What a resource being global, regional, or zonal means on GCP
    * and how to specify a default `region` and `zone`
* How GCP uses `name` and `self_link` to identify resources
* How to add GCP service account credentials to Terraform

Run `terraform destroy` to tear down your resources.

Afterwards, check out the [provider reference](/docs/providers/google/provider_reference.html) for more details on configuring
the provider block (including how you can eliminate it entirely!).
