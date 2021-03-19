---
layout: "google"
page_title: "Provider: Google Cloud Platform"
sidebar_current: "docs-google-provider-x"
description: |-
   The Google provider is used to configure your Google Cloud Platform infrastructure
---

# Google Cloud Platform Provider

The Google provider is used to configure your [Google Cloud Platform](https://cloud.google.com/) infrastructure.
See the [Getting Started](/docs/providers/google/guides/getting_started.html) page for an introduction to using the provider.

A typical provider configuration will look something like:

```hcl
provider "google" {
  project     = "my-project-id"
  region      = "us-central1"
}
```

See the [provider reference](/docs/providers/google/guides/provider_reference.html)
for more details on authentication or otherwise configuring the provider.

Take advantage of [Modules](https://www.terraform.io/docs/modules/index.html)
to simplify your config by browsing the [Module Registry for GCP modules](https://registry.terraform.io/browse?provider=google).

The Google provider is jointly maintained by:

* The [Terraform Team](https://cloud.google.com/docs/terraform) at Google
* The Terraform team at [HashiCorp](https://www.hashicorp.com/)

If you have configuration questions, or general questions about using the provider, try checking out:

* The [Google Cloud Platform Community Slack](https://gcp-slack.appspot.com/) #terraform channel
* [Terraform's community resources](https://www.terraform.io/docs/extend/community/index.html)
* [HashiCorp support](https://support.hashicorp.com) for Terraform Enterprise customers

## Releases

Interested in the provider's latest features, or want to make sure you're up to date?
Check out the [`google` provider changelog](https://github.com/hashicorp/terraform-provider-google/blob/master/CHANGELOG.md)
and the [`google-beta` provider changelog](https://github.com/hashicorp/terraform-provider-google-beta/blob/master/CHANGELOG.md))
for release notes and additional information.

Per [Terraform Provider Versioning](https://www.hashicorp.com/blog/hashicorp-terraform-provider-versioning),
the Google provider follows [semantic versioning](https://semver.org/).

In practice, patch / bugfix-only releases of the provider are infrequent. Most
provider releases are either minor or major releases.

### Minor Releases

The Google provider currently aims to publish a minor release every week,
although the timing of individual releases may differ if required by the
provider team.

### Major Releases

The Google provider publishes major releases roughly yearly. An upgrade guide
will be published to help ease you through the transition between the prior
releases series and the new major release.

During major releases, all current deprecation warnings will be resolved,
removing the field in question unless the deprecation warning message specifies
another resolution.

Before a major release, deprecation warnings don't require immediate action. The
provider team aims to surface deprecation warnings early in a major release
lifecycle to give users plenty of time to safely update their configs.

## Features and Bug Requests

The Google provider's bugs and feature requests can be found in the [GitHub repo issues](https://github.com/hashicorp/terraform-provider-google/issues).
Please avoid "me too" or "+1" comments. Instead, use a thumbs up [reaction](https://blog.github.com/2016-03-10-add-reactions-to-pull-requests-issues-and-comments/)
on enhancement requests. Provider maintainers will often prioritize work based on the
number of thumbs on an issue.

Community input is appreciated on outstanding issues! We love to hear what use
cases you have for new features, and want to provide the best possible
experience for you using the Google provider.

If you have a bug or feature request without an existing issue

* and an existing resource or field is working in an unexpected way, [file a bug](https://github.com/hashicorp/terraform-provider-google/issues/new?template=bug.md).

* and you'd like the provider to support a new resource or field, [file an enhancement/feature request](https://github.com/hashicorp/terraform-provider-google/issues/new?template=enhancement.md).

The provider maintainers will often use the assignee field on an issue to mark
who is working on it.

* An issue assigned to an individual maintainer indicates that maintainer is working
on the issue
* An issue assigned to `hashibot` indicates a member of the community has taken on
the issue!

## Contributing

If you'd like to help extend the Google provider, we gladly accept community
contributions! Our full contribution guide is available at [CONTRIBUTING.md](https://github.com/hashicorp/terraform-provider-google/blob/master/.github/CONTRIBUTING.md)

Pull requests can be made against either provider repo where a maintainer will
apply them to both `google` and `google-beta`, or against [Magic Modules](https://github.com/GoogleCloudPlatform/magic-modules)
directly.
