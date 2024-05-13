---
subcategory: "Cloud Quotas"
---

# google_cloud_quotas_quota_infos

Provides information about all quotas for a given project, folder or organization.

## Example Usage

```hcl
data "google_cloud_quotas_quota_infos" "my_quota_infos" {
    parent      = "projects/my-project"	
    service 	= "compute.googleapis.com"
}
```

## Argument Reference

The following arguments are supported:

* `parent` - (Required) Parent value of QuotaInfo resources. Listing across different resource containers (such as 'projects/-') is not allowed. Allowed parents are "projects/[project-id / number]" or "folders/[folder-id / number]" or "organizations/[org-id / number].

* `service` - (Required) The name of the service in which the quotas are defined.


## Attributes Reference

The following attributes are exported:

* `quota_infos` - (Output) The list of QuotaInfo.

<a name="nested_quota_infos"></a> The `quota_infos` block supports:

* `name` - (Output) Resource name of this QuotaInfo, for example: `projects/123/locations/global/services/compute.googleapis.com/quotaInfos/CpusPerProjectPerRegion`.
* `metric` - (Output) The metric of the quota. It specifies the resources consumption the quota is defined for, for example: `compute.googleapis.com/cpus`.
* `is_precise` - (Output) Whether this is a precise quota. A precise quota is tracked with absolute precision. In contrast, an imprecise quota is not tracked with precision.
* `refresh_interval` - (Output) The reset time interval for the quota. Refresh interval applies to rate quota only. Example: "minute" for per minute, "day" for per day, or "10 seconds" for every 10 seconds.
* `container_type` - (Output) The container type of the QuotaInfo.
* `dimensions` - (Output) The dimensions the quota is defined on.
* `metric_display_name` - (Output) The display name of the quota metric.
* `quota_display_name` - (Output) The display name of the quota.
* `metric_unit` - (Output) The unit in which the metric value is reported, e.g., `MByte`.
* `quota_increase_eligibility` - (Output) Whether it is eligible to request a higher quota value for this quota.
* `is_fixed` - (Output) Whether the quota value is fixed or adjustable.
* `dimensions_infos` - (Output) The collection of dimensions info ordered by their dimensions from more specific ones to less specific ones.
* `is_concurrent` - (Output) Whether the quota is a concurrent quota. Concurrent quotas are enforced on the total number of concurrent operations in flight at any given time.
* `service_request_quota_uri` - (Output) URI to the page where users can request more quota for the cloud service, for example: `https://console.cloud.google.com/iam-admin/quotas`.

<a name="nested_quota_increase_eligibility"></a> The `quota_increase_eligibility` block supports:

* `is_eligible` - Whether a higher quota value can be requested for the quota.
* `ineligibility_reason` - The enumeration of reasons when it is ineligible to request increase adjustment.

<a name="nested_dimensions_infos"></a> The `dimensions_infos` block supports:
* `dimensions` - The map of dimensions for this dimensions info. The key of a map entry is "region", "zone" or the name of a service specific dimension, and the value of a map entry is the value of the dimension. If a dimension does not appear in the map of dimensions, the dimensions info applies to all the dimension values except for those that have another DimenisonInfo instance configured for the specific value. Example: {"provider" : "Foo Inc"} where "provider" is a service specific dimension of a quota.

* `details` - The quota details for a map of dimensions.
* `applicable_locations` - The applicable regions or zones of this dimensions info. The field will be set to `['global']` for quotas that are not per region or per zone. Otherwise, it will be set to the list of locations this dimension info is applicable to.

<a name="nested_details"></a> The `details` block supports:
* `value` - The value currently in effect and being enforced.
