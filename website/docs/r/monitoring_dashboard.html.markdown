---
subcategory: "Cloud (Stackdriver) Monitoring"
description: |-
  A Google Stackdriver dashboard.
---

# google\_monitoring\_dashboard

A Google Stackdriver dashboard. Dashboards define the content and layout of pages in the Stackdriver web application.

To get more information about Dashboards, see:

* [API documentation](https://cloud.google.com/monitoring/api/ref_v3/rest/v1/projects.dashboards)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/monitoring/dashboards)

## Example Usage - Monitoring Dashboard Basic


```hcl
resource "google_monitoring_dashboard" "dashboard" {
  dashboard_json = <<EOF
{
  "displayName": "Demo Dashboard",
  "gridLayout": {
    "widgets": [
      {
        "blank": {}
      }
    ]
  }
}

EOF
}
```

## Example Usage - Monitoring Dashboard GridLayout


```hcl
resource "google_monitoring_dashboard" "dashboard" {
  dashboard_json = <<EOF
{
  "displayName": "Grid Layout Example",
  "gridLayout": {
    "columns": "2",
    "widgets": [
      {
        "title": "Widget 1",
        "xyChart": {
          "dataSets": [{
            "timeSeriesQuery": {
              "timeSeriesFilter": {
                "filter": "metric.type=\"agent.googleapis.com/nginx/connections/accepted_count\"",
                "aggregation": {
                  "perSeriesAligner": "ALIGN_RATE"
                }
              },
              "unitOverride": "1"
            },
            "plotType": "LINE"
          }],
          "timeshiftDuration": "0s",
          "yAxis": {
            "label": "y1Axis",
            "scale": "LINEAR"
          }
        }
      },
      {
        "text": {
          "content": "Widget 2",
          "format": "MARKDOWN"
        }
      },
      {
        "title": "Widget 3",
        "xyChart": {
          "dataSets": [{
            "timeSeriesQuery": {
              "timeSeriesFilter": {
                "filter": "metric.type=\"agent.googleapis.com/nginx/connections/accepted_count\"",
                "aggregation": {
                  "perSeriesAligner": "ALIGN_RATE"
                }
              },
              "unitOverride": "1"
            },
            "plotType": "STACKED_BAR"
          }],
          "timeshiftDuration": "0s",
          "yAxis": {
            "label": "y1Axis",
            "scale": "LINEAR"
          }
        }
      }
    ]
  }
}

EOF
}
```

## Argument Reference

The following arguments are supported:


* `dashboard_json` -
  (Required)
  The JSON representation of a dashboard, following the format at https://cloud.google.com/monitoring/api/ref_v3/rest/v1/projects.dashboards.
  The representation of an existing dashboard can be found by using the [API Explorer](https://cloud.google.com/monitoring/api/ref_v3/rest/v1/projects.dashboards/get)

  ~> **Warning:** Because this is represented as a JSON string, Terraform doesn't have underlying information to know
    which fields in the string have defaults. To prevent permanent diffs from default values, Terraform will attempt to
    suppress diffs where the value is returned in the JSON string but doesn't exist in the configuration. Consequently,
    legitmate remove-only diffs will also be suppressed. For Terraform to detect the diff, key removals must also be
    accompanied by a non-removal change (trivial or not).

- - -


* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{project_id_or_number}/dashboards/{dashboard_id}`

## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options: configuration options:

- `create` - Default is 4 minutes.
- `update` - Default is 4 minutes.
- `delete` - Default is 4 minutes.

## Import

Dashboard can be imported using any of these accepted formats:

```
$ terraform import google_monitoring_dashboard.default projects/{{project}}/dashboards/{{dashboard_id}}
$ terraform import google_monitoring_dashboard.default {{dashboard_id}}
```

## User Project Overrides

This resource supports [User Project Overrides](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#user_project_override).
