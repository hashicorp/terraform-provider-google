---
subcategory: "Compute Engine"
layout: "google"
page_title: "Google: google_compute_security_policy"
sidebar_current: "docs-google-compute-security-policy"
description: |-
  Creates a Security Policy resource for Google Compute Engine.
---

# google\_compute\_security\_policy

A Security Policy defines an IP blacklist or whitelist that protects load balanced Google Cloud services by denying or permitting traffic from specified IP ranges. For more information
see the [official documentation](https://cloud.google.com/armor/docs/configure-security-policies)
and the [API](https://cloud.google.com/compute/docs/reference/rest/beta/securityPolicies).

Security Policy is used by [`google_compute_backend_service`](https://www.terraform.io/docs/providers/google/r/compute_backend_service.html#security_policy).

## Example Usage

```hcl
resource "google_compute_security_policy" "policy" {
  name = "my-policy"

  rule {
    action   = "deny(403)"
    priority = "1000"
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["9.9.9.0/24"]
      }
    }
    description = "Deny access to IPs in 9.9.9.0/24"
  }

  rule {
    action   = "allow"
    priority = "2147483647"
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["*"]
      }
    }
    description = "default rule"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the security policy.

- - -

* `description` - (Optional) An optional description of this security policy. Max size is 2048.

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `rule` - (Optional) The set of rules that belong to this policy. There must always be a default
    rule (rule with priority 2147483647 and match "\*"). If no rules are provided when creating a
    security policy, a default rule with action "allow" will be added. Structure is [documented below](#nested_rule).

* `adaptive_protection_config` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) Configuration for [Google Cloud Armor Adaptive Protection](https://cloud.google.com/armor/docs/adaptive-protection-overview?hl=en). Structure is [documented below](#nested_adaptive_protection_config).

<a name="nested_rule"></a>The `rule` block supports:

* `action` - (Required) Action to take when `match` matches the request. Valid values:
  * "allow" : allow access to target
  * "deny(status)" : deny access to target, returns the  HTTP response code specified (valid values are 403, 404 and 502)

* `priority` - (Required) An unique positive integer indicating the priority of evaluation for a rule.
    Rules are evaluated from highest priority (lowest numerically) to lowest priority (highest numerically) in order.

* `match` - (Required) A match condition that incoming traffic is evaluated against.
    If it evaluates to true, the corresponding `action` is enforced. Structure is [documented below](#nested_match).

* `description` - (Optional) An optional description of this rule. Max size is 64.

* `preview` - (Optional) When set to true, the `action` specified above is not enforced.
    Stackdriver logs for requests that trigger a preview action are annotated as such.

<a name="nested_match"></a>The `match` block supports:

* `config` - (Optional) The configuration options available when specifying `versioned_expr`.
    This field must be specified if `versioned_expr` is specified and cannot be specified if `versioned_expr` is not specified.
    Structure is [documented below](#nested_config).

* `versioned_expr` - (Optional) Predefined rule expression. If this field is specified, `config` must also be specified.
    Available options:
    * SRC_IPS_V1: Must specify the corresponding `src_ip_ranges` field in `config`.

* `expr` - (Optional) User defined CEVAL expression. A CEVAL expression is used to specify match criteria
    such as origin.ip, source.region_code and contents in the request header.
    Structure is [documented below](#nested_expr).

<a name="nested_config"></a>The `config` block supports:

* `src_ip_ranges` - (Required) Set of IP addresses or ranges (IPV4 or IPV6) in CIDR notation
    to match against inbound traffic. There is a limit of 10 IP ranges per rule. A value of '\*' matches all IPs
    (can be used to override the default behavior).

<a name="nested_expr"></a>The `expr` block supports:

* `expression` - (Required) Textual representation of an expression in Common Expression Language syntax.
    The application context of the containing message determines which well-known feature set of CEL is supported.

<a name="nested_adaptive_protection_config"></a>The `adaptive_protection_config` block supports:

* `layer_7_ddos_defense_config` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) Configuration for [Google Cloud Armor Adaptive Protection Layer 7 DDoS Defense](https://cloud.google.com/armor/docs/adaptive-protection-overview?hl=en). Structure is [documented below](#nested_layer_7_ddos_defense_config).

<a name="nested_layer_7_ddos_defense_config"></a>The `layer_7_ddos_defense_config` block supports:

* `enable` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) If set to true, enables CAAP for L7 DDoS detection.

* `rule_visibility` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) Rule visibility can be one of the following: STANDARD - opaque rules. (default) PREMIUM - transparent rules.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `id` - an identifier for the resource with format `projects/{{project}}/global/securityPolicies/{{name}}`

* `fingerprint` - Fingerprint of this resource.

* `self_link` - The URI of the created resource.

## Import

Security policies can be imported using any of the following formats

```
$ terraform import google_compute_security_policy.policy projects/{{project}}/global/securityPolicies/{{name}}
$ terraform import google_compute_security_policy.policy {{project}}/{{name}}
$ terraform import google_compute_security_policy.policy {{name}}
```
