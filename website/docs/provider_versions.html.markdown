---
layout: "google"
page_title: "Google Provider Versions"
sidebar_current: "docs-google-provider-versions"
description: |-
  How to use the different Google Provider versions
---

# Google Provider Versions

Starting with version `1.19.0`, there are two versions of the Google provider:

* terraform-provider-google
* terraform-provider-google-beta

All GA products and features are available in both versions of the provider.

From version `2.0.0` onwards, beta features are only available in the beta version of the provider (`google-beta`).
Beta GCP Features have no deprecation policy and no SLA, but are otherwise considered to be feature-complete
with only minor outstanding issues after their Alpha period. Beta is when GCP
features are publicly announced, and is when they generally become publicly
available. For more information see [the official documentation on GCP launch stages](https://cloud.google.com/terms/launch-stages).

The beta provider sends all requests to the beta endpoint for GCP if one exists for that product, regardless of whether the request contains any beta features.

## Using multiple provider versions together

To have resources at different API versions, set up provider blocks for each version:

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

In each resource, state which provider that resource should be used with:

```hcl
resource "google_compute_instance" "ga-instance" {
  provider = "google"

  # ...
}

resource "google_compute_instance" "beta-instance" {
  provider = "google-beta"

  # ...
}
```

If the `provider` field is omitted, Terraform will choose one of the versions available to it. To be in control of which version Terraform chooses, be sure to set the `provider` field.

## Converting resources between versions

Resources can safely be converted from one version to the other without needing to rebuild infrastructure.

To go from GA to beta, change the `provider` field from `"google"` to `"google-beta"`.

To go from beta to GA, do the reverse. If you were previously using beta fields that you no longer wish to use:

1. (Optional) Explicitly set the fields back to their default values in your Terraform config file, and run `terraform apply`.
1. Change the `provider` field to `"google"`.
1. Remove any beta fields from your Terraform config.
1. Run  `terraform plan` or `terraform refresh`+`terraform show` to see that the beta fields are no longer in state.
