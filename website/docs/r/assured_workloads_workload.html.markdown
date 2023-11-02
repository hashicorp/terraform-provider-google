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
  compliance_regime = "FEDRAMP_MODERATE"
  display_name      = "{{display}}"
  location          = "us-west1"
  organization      = "123456789"
  billing_account   = "billingAccounts/000000-0000000-0000000-000000"

  kms_settings {
    next_rotation_time = "9999-10-02T15:01:23Z"
    rotation_period    = "10368000s"
  }

  provisioned_resources_parent = "folders/519620126891"

  resource_settings {
    display_name  = "folder-display-name"
    resource_type = "CONSUMER_FOLDER"
  }

  resource_settings {
    resource_type = "ENCRYPTION_KEYS_PROJECT"
  }

  resource_settings {
    resource_id   = "ring"
    resource_type = "KEYRING"
  }

  violation_notifications_enabled = true

  labels = {
    label-one = "value-one"
  }
}


```
## Example Usage - sovereign_controls_workload
A Sovereign Controls test of the assuredworkloads api
```hcl
resource "google_assured_workloads_workload" "primary" {
  compliance_regime         = "EU_REGIONS_AND_SUPPORT"
  display_name              = "display"
  location                  = "europe-west9"
  organization              = "123456789"
  billing_account           = "billingAccounts/000000-0000000-0000000-000000"
  enable_sovereign_controls = true

  kms_settings {
    next_rotation_time = "9999-10-02T15:01:23Z"
    rotation_period    = "10368000s"
  }

  resource_settings {
    resource_type = "CONSUMER_FOLDER"
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
  provider                  = google-beta
}

```

## Argument Reference

The following arguments are supported:

* `compliance_regime` -
  (Required)
  Required. Immutable. Compliance Regime associated with this workload. Possible values: COMPLIANCE_REGIME_UNSPECIFIED, IL4, CJIS, FEDRAMP_HIGH, FEDRAMP_MODERATE, US_REGIONAL_ACCESS, HIPAA, HITRUST, EU_REGIONS_AND_SUPPORT, CA_REGIONS_AND_SUPPORT, ITAR, AU_REGIONS_AND_US_SUPPORT, ASSURED_WORKLOADS_FOR_PARTNERS, ISR_REGIONS, ISR_REGIONS_AND_SUPPORT, CA_PROTECTED_B, IL5, IL2, JP_REGIONS_AND_SUPPORT
  
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

* `billing_account` -
  (Optional)
  Optional. Input only. The billing account used for the resources which are direct children of workload. This billing account is initially associated with the resources created as part of Workload creation. After the initial creation of these resources, the customer can change the assigned billing account. The resource name has the form `billingAccounts/{billing_account_id}`. For example, `billingAccounts/012345-567890-ABCDEF`.
  
* `enable_sovereign_controls` -
  (Optional)
  Optional. Indicates the sovereignty status of the given workload. Currently meant to be used by Europe/Canada customers.
  
* `kms_settings` -
  (Optional)
  **DEPRECATED** Input only. Settings used to create a CMEK crypto key. When set, a project with a KMS CMEK key is provisioned. This field is deprecated as of Feb 28, 2022. In order to create a Keyring, callers should specify, ENCRYPTION_KEYS_PROJECT or KEYRING in ResourceSettings.resource_type field.
  
* `labels` -
  (Optional)
  Optional. Labels applied to the workload.

**Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
Please refer to the field `effective_labels` for all of the labels present on the resource.
  
* `partner` -
  (Optional)
  Optional. Partner regime associated with this workload. Possible values: PARTNER_UNSPECIFIED, LOCAL_CONTROLS_BY_S3NS, SOVEREIGN_CONTROLS_BY_T_SYSTEMS, SOVEREIGN_CONTROLS_BY_SIA_MINSAIT, SOVEREIGN_CONTROLS_BY_PSN
  
* `partner_permissions` -
  (Optional)
  Optional. Permissions granted to the AW Partner SA account for the customer workload
  
* `provisioned_resources_parent` -
  (Optional)
  Input only. The parent resource for the resources managed by this Assured Workload. May be either empty or a folder resource which is a child of the Workload parent. If not specified all resources are created under the parent organization. Format: folders/{folder_id}
  
* `resource_settings` -
  (Optional)
  Input only. Resource properties that are used to customize workload resources. These properties (such as custom project id) will be used to create workload resources if possible. This field is optional.
  
* `violation_notifications_enabled` -
  (Optional)
  Optional. Indicates whether the e-mail notification for a violation is enabled for a workload. This value will be by default True, and if not present will be considered as true. This should only be updated via updateWorkload call. Any Changes to this field during the createWorkload call will not be honored. This will always be true while creating the workload.
  


The `kms_settings` block supports:
    
* `next_rotation_time` -
  (Required)
  Required. Input only. Immutable. The time at which the Key Management Service will automatically create a new version of the crypto key and mark it as the primary.
    
* `rotation_period` -
  (Required)
  Required. Input only. Immutable. will be advanced by this period when the Key Management Service automatically rotates a key. Must be at least 24 hours and at most 876,000 hours.
    
The `partner_permissions` block supports:
    
* `assured_workloads_monitoring` -
  (Optional)
  Optional. Allow partner to view violation alerts.
    
* `data_logs_viewer` -
  (Optional)
  Allow the partner to view inspectability logs and monitoring violations.
    
* `service_access_approver` -
  (Optional)
  Optional. Allow partner to view access approval logs.
    
The `resource_settings` block supports:
    
* `display_name` -
  (Optional)
  User-assigned resource display name. If not empty it will be used to create a resource with the specified name.
    
* `resource_id` -
  (Optional)
  Resource identifier. For a project this represents projectId. If the project is already taken, the workload creation will fail. For KeyRing, this represents the keyring_id. For a folder, don't set this value as folder_id is assigned by Google.
    
* `resource_type` -
  (Optional)
  Indicates the type of resource. This field should be specified to correspond the id to the right project type (CONSUMER_PROJECT or ENCRYPTION_KEYS_PROJECT) Possible values: RESOURCE_TYPE_UNSPECIFIED, CONSUMER_PROJECT, ENCRYPTION_KEYS_PROJECT, KEYRING, CONSUMER_FOLDER
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `organizations/{{organization}}/locations/{{location}}/workloads/{{name}}`

* `compliance_status` -
  Output only. Count of active Violations in the Workload.
  
* `compliant_but_disallowed_services` -
  Output only. Urls for services which are compliant for this Assured Workload, but which are currently disallowed by the ResourceUsageRestriction org policy. Invoke workloads.restrictAllowedResources endpoint to allow your project developers to use these services in their environment.
  
* `create_time` -
  Output only. Immutable. The Workload creation timestamp.
  
* `effective_labels` -
  All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.
  
* `ekm_provisioning_response` -
  Optional. Represents the Ekm Provisioning State of the given workload.
  
* `kaj_enrollment_state` -
  Output only. Represents the KAJ enrollment state of the given workload. Possible values: KAJ_ENROLLMENT_STATE_UNSPECIFIED, KAJ_ENROLLMENT_STATE_PENDING, KAJ_ENROLLMENT_STATE_COMPLETE
  
* `name` -
  Output only. The resource name of the workload.
  
* `resources` -
  Output only. The resources associated with this workload. These resources will be created when creating the workload. If any of the projects already exist, the workload creation will fail. Always read only.
  
* `saa_enrollment_response` -
  Output only. Represents the SAA enrollment response of the given workload. SAA enrollment response is queried during workloads.get call. In failure cases, user friendly error message is shown in SAA details page.
  
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
* `organizations/{{organization}}/locations/{{location}}/workloads/{{name}}`
* `{{organization}}/{{location}}/{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Workload using one of the formats above. For example:


```tf
import {
  id = "organizations/{{organization}}/locations/{{location}}/workloads/{{name}}"
  to = google_assured_workloads_workload.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Workload can be imported using one of the formats above. For example:

```
$ terraform import google_assured_workloads_workload.default organizations/{{organization}}/locations/{{location}}/workloads/{{name}}
$ terraform import google_assured_workloads_workload.default {{organization}}/{{location}}/{{name}}
```



