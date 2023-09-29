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
subcategory: "AssuredWorkloads"
description: |-
  The AssuredWorkloads Workload resource
---

# google_assured_workloads_workload

The AssuredWorkloads Workload resource

## Example Usage - basic_workload
A basic test of a assuredworkloads api
```hcl
resource "google_assured_workloads_workload" "primary" {
  billing_account   = "billingAccounts/000000-0000000-0000000-000000"
  compliance_regime = "FEDRAMP_MODERATE"
  display_name      = "Workload Example"
  location          = "us-west1"
  organization      = "123456789"

  kms_settings {
    next_rotation_time = "9999-10-02T15:01:23Z"
    rotation_period    = "10368000s"
  }

  provisioned_resources_parent = "folders/519620126891"

  resource_settings {
    resource_type = "CONSUMER_PROJECT"
  }

  resource_settings {
    resource_type = "ENCRYPTION_KEYS_PROJECT"
  }

  resource_settings {
    resource_id   = "ring"
    resource_type = "KEYRING"
  }

  labels = {
    label-one = "value-one"
  }
}


```

## Argument Reference

The following arguments are supported:

* `billing_account` -
  (Required)
  Required. Input only. The billing account used for the resources which are direct children of workload. This billing account is initially associated with the resources created as part of Workload creation. After the initial creation of these resources, the customer can change the assigned billing account. The resource name has the form `billingAccounts/{billing_account_id}`. For example, 'billingAccounts/012345-567890-ABCDEF`.
  
* `compliance_regime` -
  (Required)
  Required. Immutable. Compliance Regime associated with this workload. Possible values: COMPLIANCE_REGIME_UNSPECIFIED, IL4, CJIS, FEDRAMP_HIGH, FEDRAMP_MODERATE, US_REGIONAL_ACCESS, HIPAA, EU_REGIONS_AND_SUPPORT, CA_REGIONS_AND_SUPPORT, ITAR, AU_REGIONS_AND_US_SUPPORT, ASSURED_WORKLOADS_FOR_PARTNERS
  
* `display_name` -
  (Required)
  Required. The user-assigned display name of the Workload. When present it must be between 4 to 30 characters. Allowed characters are: lowercase and uppercase letters, numbers, hyphen, and spaces. Example: My Workload
  
* `location` -
  (Required)
  The location for the resource
  
* `organization` -
  (Required)
  The organization for the resource
  


- - -

* `kms_settings` -
  (Optional)
  Input only. Settings used to create a CMEK crypto key. When set a project with a KMS CMEK key is provisioned. This field is mandatory for a subset of Compliance Regimes.
  
* `labels` -
  (Optional)
  Optional. Labels applied to the workload.

**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
Please refer to the field `effective_labels` for all of the labels present on the resource.
  
* `provisioned_resources_parent` -
  (Optional)
  Input only. The parent resource for the resources managed by this Assured Workload. May be either an organization or a folder. Must be the same or a child of the Workload parent. If not specified all resources are created under the Workload parent. Formats: folders/{folder_id}, organizations/{organization_id}
  
* `resource_settings` -
  (Optional)
  Input only. Resource properties that are used to customize workload resources. These properties (such as custom project id) will be used to create workload resources if possible. This field is optional.
  


The `kms_settings` block supports:
    
* `next_rotation_time` -
  (Required)
  Required. Input only. Immutable. The time at which the Key Management Service will automatically create a new version of the crypto key and mark it as the primary.
    
* `rotation_period` -
  (Required)
  Required. Input only. Immutable. will be advanced by this period when the Key Management Service automatically rotates a key. Must be at least 24 hours and at most 876,000 hours.
    
The `resource_settings` block supports:
    
* `resource_id` -
  (Optional)
  Resource identifier. For a project this represents project_number. If the project is already taken, the workload creation will fail.
    
* `resource_type` -
  (Optional)
  Indicates the type of resource. This field should be specified to correspond the id to the right project type (CONSUMER_PROJECT or ENCRYPTION_KEYS_PROJECT) Possible values: RESOURCE_TYPE_UNSPECIFIED, CONSUMER_PROJECT, ENCRYPTION_KEYS_PROJECT, KEYRING, CONSUMER_FOLDER
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `organizations/{{organization}}/locations/{{location}}/workloads/{{name}}`

* `create_time` -
  Output only. Immutable. The Workload creation timestamp.
  
* `effective_labels` -
  All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.
  
* `name` -
  Output only. The resource name of the workload.
  
* `resources` -
  Output only. The resources associated with this workload. These resources will be created when creating the workload. If any of the projects already exist, the workload creation will fail. Always read only.
  
* `terraform_labels` -
  The combination of labels configured directly on the resource and default labels configured on the provider.
  
## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Workload can be imported using any of these accepted formats:

```
$ terraform import google_assured_workloads_workload.default organizations/{{organization}}/locations/{{location}}/workloads/{{name}}
$ terraform import google_assured_workloads_workload.default {{organization}}/{{location}}/{{name}}
```



