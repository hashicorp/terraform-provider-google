---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_ssl_policy"
sidebar_current: "docs-google-datasource-compute-ssl-policy"
description: |-
  Gets an SSL Policy within GCE, for use with Target HTTPS and Target SSL Proxies.
---

# google\_compute\_ssl\_policy

Gets an SSL Policy within GCE from its name, for use with Target HTTPS and Target SSL Proxies.
    For more information see [the official documentation](https://cloud.google.com/compute/docs/load-balancing/ssl-policies).

## Example Usage

```tf
data "google_compute_ssl_policy" "my-ssl-policy" {
  name = "production-ssl-policy"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SSL Policy.

- - -

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `enabled_features` - The set of enabled encryption ciphers as a result of the policy config

* `description` - Description of this SSL Policy.

* `min_tls_version` - The minimum supported TLS version of this policy.

* `profile` - The Google-curated or custom profile used by this policy.

* `custom_features` - If the `profile` is `CUSTOM`, these are the custom encryption
    ciphers supported by the profile. If the `profile` is *not* `CUSTOM`, this
    attribute will be empty.

* `fingerprint` - Fingerprint of this resource.

* `self_link` - The URI of the created resource.
