---
subcategory: "Compute Engine"
page_title: "Google: google_compute_security_policy"
description: |-
  Creates a Security Policy resource for Google Compute Engine.
---

# google\_compute\_security\_policy

A Security Policy defines an IP blacklist or whitelist that protects load balanced Google Cloud services by denying or permitting traffic from specified IP ranges. For more information
see the [official documentation](https://cloud.google.com/armor/docs/configure-security-policies)
and the [API](https://cloud.google.com/compute/docs/reference/rest/beta/securityPolicies).

Security Policy is used by [`google_compute_backend_service`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_backend_service#security_policy).

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

## Example Usage - With reCAPTCHA configuration options

```hcl
resource "google_recaptcha_enterprise_key" "primary" {
  display_name = "display-name"

  labels = {
    label-one = "value-one"
   }

  project = "my-project-name"

  web_settings {
    integration_type  = "INVISIBLE"
    allow_all_domains = true
    allowed_domains   = ["localhost"]
  }
}

resource "google_compute_security_policy" "policy" {
  name        = "my-policy"
  description = "basic security policy"
  type        = "CLOUD_ARMOR"

  recaptcha_options_config {
    redirect_site_key = google_recaptcha_enterprise_key.primary.name
  }
}
```

## Example Usage - With header actions

```hcl
resource "google_compute_security_policy" "policy" {
	name = "my-policy"

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

  rule {
    action   = "allow"
    priority = "1000"
    match {
      expr {
        expression = "request.path.matches(\"/login.html\") && token.recaptcha_session.score < 0.2"
      }
    }

    header_action {
      request_headers_to_adds {
        header_name  = "reCAPTCHA-Warning"
        header_value = "high"
      }

      request_headers_to_adds {
        header_name  = "X-Resource"
        header_value = "test"
      }
    }
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

* `advanced_options_config` - (Optional) [Advanced Configuration Options](https://cloud.google.com/armor/docs/security-policy-overview#json-parsing).
    Structure is [documented below](#nested_advanced_options_config).

* `adaptive_protection_config` - (Optional) Configuration for [Google Cloud Armor Adaptive Protection](https://cloud.google.com/armor/docs/adaptive-protection-overview?hl=en). Structure is [documented below](#nested_adaptive_protection_config).

* `recaptcha_options_config` - (Optional) [reCAPTCHA Configuration Options](https://cloud.google.com/armor/docs/configure-security-policies?hl=en#use_a_manual_challenge_to_distinguish_between_human_or_automated_clients). Structure is [documented below](#nested_recaptcha_options_config).

* `type` - The type indicates the intended use of the security policy. This field can be set only at resource creation time.
  * CLOUD_ARMOR - Cloud Armor backend security policies can be configured to filter incoming HTTP requests targeting backend services.
    They filter requests before they hit the origin servers.
  * CLOUD_ARMOR_EDGE - Cloud Armor edge security policies can be configured to filter incoming HTTP requests targeting backend services
    (including Cloud CDN-enabled) as well as backend buckets (Cloud Storage).
    They filter requests before the request is served from Google's cache.
  * CLOUD_ARMOR_INTERNAL_SERVICE - Cloud Armor internal service policies can be configured to filter HTTP requests targeting services 
    managed by Traffic Director in a service mesh. They filter requests before the request is served from the application.

<a name="nested_advanced_options_config"></a>The `advanced_options_config` block supports:

* `json_parsing` - Whether or not to JSON parse the payload body. Defaults to `DISABLED`.
  * DISABLED - Don't parse JSON payloads in POST bodies.
  * STANDARD - Parse JSON payloads in POST bodies.

* `json_custom_config` - Custom configuration to apply the JSON parsing. Only applicable when
    `json_parsing` is set to `STANDARD`. Structure is [documented below](#nested_json_custom_config).

* `log_level` - Log level to use. Defaults to `NORMAL`.
  * NORMAL - Normal log level.
  * VERBOSE - Verbose log level.

<a name="nested_json_custom_config"></a>The `json_custom_config` block supports:

* `content_types` - A list of custom Content-Type header values to apply the JSON parsing. The
    format of the Content-Type header values is defined in
    [RFC 1341](https://www.ietf.org/rfc/rfc1341.txt). When configuring a custom Content-Type header
    value, only the type/subtype needs to be specified, and the parameters should be excluded.

<a name="nested_rule"></a>The `rule` block supports:

* `action` - (Required) Action to take when `match` matches the request. Valid values:
    * allow: allow access to target.
    * deny(): deny access to target, returns the HTTP response code specified (valid values are 403, 404, and 502).
    * rate_based_ban: limit client traffic to the configured threshold and ban the client if the traffic exceeds the threshold. Configure parameters for this action in RateLimitOptions. Requires rateLimitOptions to be set.
    * redirect: redirect to a different target. This can either be an internal reCAPTCHA redirect, or an external URL-based redirect via a 302 response. Parameters for this action can be configured via redirectOptions.
    * throttle: limit client traffic to the configured threshold. Configure parameters for this action in rateLimitOptions. Requires rateLimitOptions to be set for this.

* `priority` - (Required) An unique positive integer indicating the priority of evaluation for a rule.
    Rules are evaluated from highest priority (lowest numerically) to lowest priority (highest numerically) in order.

* `match` - (Required) A match condition that incoming traffic is evaluated against.
    If it evaluates to true, the corresponding `action` is enforced. Structure is [documented below](#nested_match).

* `preconfigured_waf_config` - (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html)) Preconfigured WAF configuration to be applied for the rule. If the rule does not evaluate preconfigured WAF rules, i.e., if evaluatePreconfiguredWaf() is not used, this field will have no effect. Structure is [documented below](#nested_preconfigured_waf_config).

* `description` - (Optional) An optional description of this rule. Max size is 64.

* `preview` - (Optional) When set to true, the `action` specified above is not enforced.
    Stackdriver logs for requests that trigger a preview action are annotated as such.

* `rate_limit_options` - (Optional)
    Must be specified if the `action` is "rate_based_ban" or "throttle". Cannot be specified for other actions. Structure is [documented below](#nested_rate_limit_options).

* `redirect_options` - (Optional)
    Can be specified if the `action` is "redirect". Cannot be specified for other actions. Structure is [documented below](#nested_redirect_options).

* `header_action` - (Optional)
    Additional actions that are performed on headers. Structure is [documented below](#nested_header_action).

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

<a name="nested_preconfigured_waf_config"></a>The `preconfigured_waf_config` block supports:

* `exclusion` - (Optional) An exclusion to apply during preconfigured WAF evaluation. Structure is [documented below](#nested_exclusion).

<a name="nested_exclusion"></a>The `exclusion` block supports:

* `request_header` - (Optional) Request header whose value will be excluded from inspection during preconfigured WAF evaluation. Structure is [documented below](#nested_field_params).

* `request_cookie` - (Optional) Request cookie whose value will be excluded from inspection during preconfigured WAF evaluation. Structure is [documented below](#nested_field_params).

* `request_uri` - (Optional) Request query parameter whose value will be excluded from inspection during preconfigured WAF evaluation. Note that the parameter can be in the query string or in the POST body. Structure is [documented below](#nested_field_params).

* `request_query_param` - (Optional) Request URI from the request line to be excluded from inspection during preconfigured WAF evaluation. When specifying this field, the query or fragment part should be excluded. Structure is [documented below](#nested_field_params).

* `target_rule_set` - (Required) Target WAF rule set to apply the preconfigured WAF exclusion.

* `target_rule_ids` - (Optional) A list of target rule IDs under the WAF rule set to apply the preconfigured WAF exclusion. If omitted, it refers to all the rule IDs under the WAF rule set.

<a name="nested_field_params"></a>The `request_header`, `request_cookie`, `request_uri` and `request_query_param` blocks support:

* `operator` - (Required) You can specify an exact match or a partial match by using a field operator and a field value.

  * EQUALS: The operator matches if the field value equals the specified value.
  * STARTS_WITH: The operator matches if the field value starts with the specified value.
  * ENDS_WITH: The operator matches if the field value ends with the specified value.
  * CONTAINS: The operator matches if the field value contains the specified value.
  * EQUALS_ANY: The operator matches if the field value is any value.

* `value` - (Optional) A request field matching the specified value will be excluded from inspection during preconfigured WAF evaluation.
    The field value must be given if the field `operator` is not "EQUALS_ANY", and cannot be given if the field `operator` is "EQUALS_ANY".

<a name="nested_rate_limit_options"></a>The `rate_limit_options` block supports:

* `conform_action` - (Required) Action to take for requests that are under the configured rate limit threshold. Valid option is "allow" only.

* `exceed_action` - (Required) When a request is denied, returns the HTTP response code specified.
    Valid options are "deny()" where valid values for status are 403, 404, 429, and 502.

* `rate_limit_threshold` - (Required) Threshold at which to begin ratelimiting. Structure is [documented below](#nested_threshold).

* `ban_duration_sec` - (Optional) Can only be specified if the `action` for the rule is "rate_based_ban".
    If specified, determines the time (in seconds) the traffic will continue to be banned by the rate limit after the rate falls below the threshold.

* `ban_threshold` - (Optional) Can only be specified if the `action` for the rule is "rate_based_ban".
    If specified, the key will be banned for the configured 'ban_duration_sec' when the number of requests that exceed the 'rate_limit_threshold' also
    exceed this 'ban_threshold'. Structure is [documented below](#nested_threshold).

* `enforce_on_key` - (Optional) Determines the key to enforce the rate_limit_threshold on. If not specified, defaults to "ALL".

    * ALL: A single rate limit threshold is applied to all the requests matching this rule.
    * IP: The source IP address of the request is the key. Each IP has this limit enforced separately.
    * HTTP_HEADER: The value of the HTTP header whose name is configured under "enforceOnKeyName". The key value is truncated to the first 128 bytes of the header value. If no such header is present in the request, the key type defaults to ALL.
    * XFF_IP: The first IP address (i.e. the originating client IP address) specified in the list of IPs under X-Forwarded-For HTTP header. If no such header is present or the value is not a valid IP, the key type defaults to ALL.
    * HTTP_COOKIE: The value of the HTTP cookie whose name is configured under "enforceOnKeyName". The key value is truncated to the first 128 bytes of the cookie value. If no such cookie is present in the request, the key type defaults to ALL.

* `enforce_on_key_name` - (Optional) Rate limit key name applicable only for the following key types: HTTP_HEADER -- Name of the HTTP header whose value is taken as the key value. HTTP_COOKIE -- Name of the HTTP cookie whose value is taken as the key value.

* `exceed_redirect_options` - (Optional) Parameters defining the redirect action that is used as the exceed action. Cannot be specified if the exceed action is not redirect. Structure is [documented below](#nested_exceed_redirect_options).

<a name="nested_threshold"></a>The `{ban/rate_limit}_threshold` block supports:

* `count` - (Required) Number of HTTP(S) requests for calculating the threshold.

* `interval_sec` - (Required) Interval over which the threshold is computed.

* <a  name="nested_exceed_redirect_options"></a>The `exceed_redirect_options` block supports:

* `type` - (Required) Type of the redirect action.

* `target` - (Optional) Target for the redirect action. This is required if the type is EXTERNAL_302 and cannot be specified for GOOGLE_RECAPTCHA.

<a name="nested_redirect_options"></a>The `redirect_options` block supports:

* `type` - (Required) Type of redirect action.

    * EXTERNAL_302: Redirect to an external address, configured in 'target'.
    * GOOGLE_RECAPTCHA: Redirect to Google reCAPTCHA.

* `target` - (Optional) External redirection target when "EXTERNAL_302" is set in 'type'.

<a name="nested_header_action"></a> The `header_action` block supports:

* `request_headers_to_adds` - (Required) The list of request headers to add or overwrite if they're already present. Structure is [documented below](#nested_request_headers_to_adds).

<a name="nested_request_headers_to_adds"><a> The `request_headers_to_adds` block supports:

* `header_name` - (Required) The name of the header to set.

* `header_value` - (Optional) The value to set the named header to.

<a name="nested_adaptive_protection_config"></a>The `adaptive_protection_config` block supports:

* `layer_7_ddos_defense_config` - (Optional) Configuration for [Google Cloud Armor Adaptive Protection Layer 7 DDoS Defense](https://cloud.google.com/armor/docs/adaptive-protection-overview?hl=en). Structure is [documented below](#nested_layer_7_ddos_defense_config).

<a name="nested_layer_7_ddos_defense_config"></a>The `layer_7_ddos_defense_config` block supports:

* `enable` - (Optional) If set to true, enables CAAP for L7 DDoS detection.

* `rule_visibility` - (Optional) Rule visibility can be one of the following: STANDARD - opaque rules. (default) PREMIUM - transparent rules.

<a name="nested_recaptcha_options_config"></a>The `recaptcha_options_config` block supports:

* `redirect_site_key` - (Required) A field to supply a reCAPTCHA site key to be used for all the rules using the redirect action with the type of GOOGLE_RECAPTCHA under the security policy. The specified site key needs to be created from the reCAPTCHA API. The user is responsible for the validity of the specified site key. If not specified, a Google-managed site key is used.

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
