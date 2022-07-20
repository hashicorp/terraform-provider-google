---
page_title: "Google Provider Versions"
description: |-
  How to use the different Google Provider versions
---

# Google Provider Versions

Starting with the `1.19.0` release, there are two versions of the Google Cloud Platform
provider:

* `google`

* `google-beta`

This documentation (https://registry.terraform.io/providers/hashicorp/google/latest/docs) is shared
between both providers, and all generally available (GA) products and features
are available in both versions of the provider.

You may see beta features referenced as Preview since Google simplified the [product launch stages](https://cloud.google.com/blog/products/gcp/google-cloud-gets-simplified-product-launch-stages) in late 2020.

The `google-beta` provider is distinct from the `google` provider in that it
supports GCP products and features that are in Preview, while `google` does not.
Fields and resources that are only present in `google-beta` are clearly marked in the provider documentation.

Pre-GA products and features might have limited support, and changes to pre-GA products and features might not be compatible with other pre-GA versions. For more information, see the [launch stage descriptions](https://cloud.google.com/products#product-launch-stages).

The `google-beta` provider sends all requests to the beta endpoint for GCP if
one exists for that product, regardless of whether the request contains any beta
features.

-> In short, using `google-beta` over `google` is similar to using `gcloud beta`
over `gcloud`. Features that are exclusively available in `google-beta` are GCP
features that are not yet GA, and they will be made available in `google` after
they go to GA.

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
  project     = "my-project-id"
  region      = "us-central1"
}

provider "google-beta" {
  project     = "my-project-id"
  region      = "us-central1"
}
```

## Importing resources with `google-beta`
By default, Terraform will always import resources using the `google` provider.
To import resources with `google-beta`, you need to explicitly specify a provider
with the `-provider` flag, similarly to if you were using a provider alias.


```bash
terraform import google_compute_instance.beta-instance my-instance
```

## Converting resources between versions

Resources can safely be converted from one version to the other without needing to rebuild infrastructure.

To go from GA to beta, change the `provider` field from `"google"` to `"google-beta"`.

To go from beta to GA, do the reverse. If you were previously using beta fields that you no longer wish to use:

1. (Optional) Explicitly set the fields back to their default values in your Terraform config file, and run `terraform apply`.
1. Change the `provider` field to `"google"`.
1. Remove any beta fields from your Terraform config.
1. Run `terraform plan` or `terraform refresh`+`terraform show` to see that the beta fields are no longer in state.
