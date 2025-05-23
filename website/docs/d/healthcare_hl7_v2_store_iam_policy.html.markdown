---
# ----------------------------------------------------------------------------
#
#     ***     AUTO GENERATED CODE    ***    Type: Handwritten     ***
#
# ----------------------------------------------------------------------------
#
#     This code is generated by Magic Modules using the following:
#
#     Source file: https://github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/third_party/terraform/website/docs/d/healthcare_hl7_v2_store_iam_policy.html.markdown
#
#     DO NOT EDIT this file directly. Any changes made to this file will be
#     overwritten during the next generation cycle.
#
# ----------------------------------------------------------------------------
subcategory: "Cloud Healthcare"
description: |-
  A datasource to retrieve the IAM policy state for a Google Cloud Healthcare HL7v2 store.
---


# `google_healthcare_hl7_v2_store_iam_policy`
Retrieves the current IAM policy data for a Google Cloud Healthcare HL7v2 store.

## example

```hcl
data "google_healthcare_hl7_v2_store_iam_policy" "foo" {
  hl7_v2_store_id = google_healthcare_hl7_v2_store.hl7_v2_store.id
}
```

## Argument Reference

The following arguments are supported:

* `hl7_v2_store_id` - (Required) The HL7v2 store ID, in the form
    `{project_id}/{location_name}/{dataset_name}/{hl7_v2_store_name}` or
    `{location_name}/{dataset_name}/{hl7_v2_store_name}`. In the second form, the provider's
    project setting will be used as a fallback.

## Attributes Reference

The attributes are exported:

* `etag` - (Computed) The etag of the IAM policy.

* `policy_data` - (Computed) The policy data
