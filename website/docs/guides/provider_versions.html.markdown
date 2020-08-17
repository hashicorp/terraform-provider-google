---
layout: "google"
page_title: "Google Provider Versions"
sidebar_current: "docs-google-provider-versions"
description: |-
  How to use the different Google Provider versions
---

# Google Provider Versions

Starting with the `1.19.0` release, there are two versions of the Google
provider:

* `google`

* `google-beta`

This documentation (https://www.terraform.io/docs/providers/google/) is shared
between both providers, and all generally available (GA) products and features
are available in both versions of the provider.

The `google-beta` provider is distinct from the `google` provider in that it
supports GCP products and features that are in beta, while `google` does not.
Fields and resources that are only present in `google-beta` will be marked as
such in the shared provider documentation.

`1.X` series releases of the `google` provider supported beta features; from
`2.0.0` onwards, beta features are only supported in `google-beta`.

Beta GCP features have no deprecation policy and no SLA, but are otherwise considered to be feature-complete
with only minor outstanding issues after their Alpha period. Beta is when GCP
features are publicly announced, and is when they generally become publicly
available. For more information see [the official documentation on GCP launch stages](https://cloud.google.com/terms/launch-stages).

The `google-beta` provider sends all requests to the beta endpoint for GCP if
one exists for that product, regardless of whether the request contains any beta
features.

-> In short, using `google-beta` over `google` is similar to using `gcloud beta`
over `gcloud`. Features that are exclusively available in `google-beta` are GCP
features that are not yet GA, and they will be made available in `google` after
their GA launch.

## Using the `google-beta` provider

To use the `google-beta` provider, simply set the `provider` field on each
resource where you want to use `google-beta`.

```hcl
resource "google_compute_instance" "beta-instance" {
  provider = google-beta
  # ...
}
```

To customize the behavior of the beta provider, you can define a `google-beta`
provider block, which accepts the same arguments as the `google` provider block.

```hcl
provider "google-beta" {
  credentials = "${file("account.json")}"
  project     = "my-project-id"
  region      = "us-central1"
}
```

~> If the `provider` field is omitted, Terraform will implicitly use the `google`
 provider by default even if you have only defined a `google-beta` provider block.

## Using both provider versions together

It is safe to use both provider versions in the same configuration.

In each resource, state which provider that resource should be used with.
We recommend that you set `provider = google` even though it is the default,
for clarity.

```hcl
resource "google_compute_instance" "ga-instance" {
  provider = google

  # ...
}

resource "google_compute_instance" "beta-instance" {
  provider = google-beta

  # ...
}
```

You can define parallel provider blocks - they will not interfere with each other.

```hcl
provider "google" {
  credentials = "${file("account.json")}"
  project     = "my-project-id"
  region      = "us-central1"
}

provider "google-beta" {
  credentials = "${file("account.json")}"
  project     = "my-project-id"
  region      = "us-central1"
}
```

## Importing resources with `google-beta`
By default, Terraform will always import resources using the `google` provider.
To import resources with `google-beta`, you need to explicitly specify a provider
with the `-provider` flag, similarly to if you were using a provider alias.


```bash
terraform import -provider=google-beta google_compute_instance.beta-instance my-instance
```

## Converting resources between versions

Resources can safely be converted from one version to the other without needing to rebuild infrastructure.

To go from GA to beta, change the `provider` field from `"google"` to `"google-beta"`.

To go from beta to GA, do the reverse. If you were previously using beta fields that you no longer wish to use:

1. (Optional) Explicitly set the fields back to their default values in your Terraform config file, and run `terraform apply`.
1. Change the `provider` field to `"google"`.
1. Remove any beta fields from your Terraform config.
1. Run  `terraform plan` or `terraform refresh`+`terraform show` to see that the beta fields are no longer in state.
