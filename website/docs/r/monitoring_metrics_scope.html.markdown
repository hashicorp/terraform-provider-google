---
# ----------------------------------------------------------------------------
#
#     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
#
# ----------------------------------------------------------------------------
#
#     This file is managed by Magic Modules (https:#github.com/GoogleCloudPlatform/magic-modules)
#     and is based on the DCL (https:#github.com/GoogleCloudPlatform/declarative-resource-client-library).
#     Changes will need to be made to the DCL or Magic Modules instead of here.
#
#     We are not currently able to accept contributions to this file. If changes
#     are required, please file an issue at https:#github.com/hashicorp/terraform-provider-google/issues/new/choose
#
# ----------------------------------------------------------------------------
subcategory: "Monitoring"
layout: "google"
page_title: "Google: google_monitoring_metrics_scope"
sidebar_current: "docs-google-monitoring-metrics-scope"
description: |-

---

# google\_monitoring\_metrics\_scope


## Example Usage - basic_metrics_scope
A basic example of a monitoring metrics scope
```hcl
resource "google_monitoring_metrics_scope" "primary" {
  name = "my-project-name"
}


```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  Immutable. The resource name of the Monitoring Metrics Scope. On input, the resource name can be specified with the scoping project ID or number. On output, the resource name is specified with the scoping project number. Example: `locations/global/metricsScopes/{SCOPING_PROJECT_ID_OR_NUMBER}`
  


- - -



## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `locations/global/metricsScopes/{{name}}`

* `create_time` -
  Output only. The time when this `Metrics Scope` was created.
  
* `monitored_projects` -
  Output only. The list of projects monitored by this `Metrics Scope`.
  
* `update_time` -
  Output only. The time when this `Metrics Scope` record was last updated.
  
## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 10 minutes.
- `delete` - Default is 10 minutes.

## Import

MetricsScope can be imported using any of these accepted formats:

```
$ terraform import google_monitoring_metrics_scope.default locations/global/metricsScopes/{{name}}
$ terraform import google_monitoring_metrics_scope.default {{name}}
```



