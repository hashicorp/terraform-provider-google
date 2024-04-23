---
subcategory: "Dataproc"
description: |-
  A datasource to retrieve the IAM policy state for a Dataproc cluster.
---


# `google_dataproc_cluster_iam_policy`
Retrieves the current IAM policy data for a Dataproc cluster.

## example

```hcl
data "google_dataproc_cluster_iam_policy" "policy" {
  cluster     = google_dataproc_cluster.cluster.name
  region      = "us-central1"
}
```

## Argument Reference

The following arguments are supported:

* `cluster` - (Required) The name or relative resource id of the cluster to manage IAM policies for.

* `project` - (Optional) The project in which the cluster belongs. If it
    is not provided, Terraform will use the provider default.

* `region` - (Optional) The region in which the cluster belongs. If it
    is not provided, Terraform will use the provider default.

## Attributes Reference

The attributes are exported:

* `etag` - (Computed) The etag of the IAM policy.

* `policy_data` - (Computed) The policy data
