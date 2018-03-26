---
layout: "google"
page_title: "Google: google_compute_ssl_policy"
sidebar_current: "docs-google-compute-ssl-policy"
description: |-
  Manages an SSL Policy within GCE, for use with Target HTTPS and Target SSL Proxies.
---

# google\_compute\_ssl\_policy

Manages an SSL Policy within GCE, for use with Target HTTPS and Target SSL Proxies. For more information see
[the official documentation](https://cloud.google.com/compute/docs/load-balancing/ssl-policies)
and
[API](https://cloud.google.com/compute/docs/reference/rest/beta/sslPolicies).

## Example Usage

```hcl
resource "google_compute_ssl_policy" "prod-ssl-policy" {
  name    = "production-ssl-policy"
  profile = "MODERN"
}

resource "google_compute_ssl_policy" "nonprod-ssl-policy" {
  name            = "nonprod-ssl-policy"
  profile         = "MODERN"
  min_tls_version = "TLS_1_2"
}

resource "google_compute_ssl_policy" "custom-ssl-policy" {
  name            = "custom-ssl-policy"
  min_tls_version = "TLS_1_2"
  profile         = "CUSTOM"
  custom_features = ["TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384", "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the resource, required by GCE.
    Changing this forces a new resource to be created.

- - -

* `description` - (Optional) Description of this subnetwork. Changing this forces a new resource to be created.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `min_tls_version` - (Optional) The minimum TLS version to support. Must be one of `TLS_1_0`, `TLS_1_1`, or `TLS_1_2`. 
    Default is `TLS_1_0`.

* `profile` - (Optional) The Google-curated SSL profile to use. Must be one of `COMPATIBLE`, `MODERN`, 
    `RESTRICTED`, or `CUSTOM`. See the 
    [official documentation](https://cloud.google.com/compute/docs/load-balancing/ssl-policies#profilefeaturesupport) 
    for information on what cipher suites each profile provides. If `CUSTOM` is used, the `custom_features` attribute 
    **must be set**. Default is `COMPATIBLE`.

* `custom_features` - (Required with `CUSTOM` profile) The specific encryption ciphers to use. See the 
    [official documentation](https://cloud.google.com/compute/docs/load-balancing/ssl-policies#profilefeaturesupport) 
    for which ciphers are available to use. **Note**: this argument *must* be present when using the `CUSTOM` profile. 
    This argument *must not* be present when using any other profile.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `fingerprint` - Fingerprint of this resource.

* `self_link` - The URI of the created resource.

## Import

SSL Policies can be imported using the GCP canonical `name` of the Policy. For example, an SSL Policy named `production-ssl-policy` 
    would be imported by running:

```bash
$ terraform import google_compute_ssl_policy.my-policy production-ssl-policy
```
