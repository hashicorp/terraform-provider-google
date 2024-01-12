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
subcategory: "OrgPolicy"
description: |-
  An organization policy gives you programmatic control over your organization's cloud resources.  Using Organization Policies, you will be able to configure constraints across your entire resource hierarchy.
---

# google_org_policy_policy

An organization policy gives you programmatic control over your organization's cloud resources.  Using Organization Policies, you will be able to configure constraints across your entire resource hierarchy.

For more information, see:
* [Understanding Org Policy concepts](https://cloud.google.com/resource-manager/docs/organization-policy/overview)
* [The resource hierarchy](https://cloud.google.com/resource-manager/docs/cloud-platform-resource-hierarchy)
* [All valid constraints](https://cloud.google.com/resource-manager/docs/organization-policy/org-policy-constraints)
## Example Usage - enforce_policy
A test of an enforce orgpolicy policy for a project
```hcl
resource "google_org_policy_policy" "primary" {
  name   = "projects/${google_project.basic.name}/policies/iam.disableServiceAccountKeyUpload"
  parent = "projects/${google_project.basic.name}"

  spec {
    rules {
      enforce = "FALSE"
    }
  }
}

resource "google_project" "basic" {
  project_id = "id"
  name       = "id"
  org_id     = "123456789"
}


```
## Example Usage - folder_policy
A test of an orgpolicy policy for a folder
```hcl
resource "google_org_policy_policy" "primary" {
  name   = "${google_folder.basic.name}/policies/gcp.resourceLocations"
  parent = google_folder.basic.name

  spec {
    inherit_from_parent = true

    rules {
      deny_all = "TRUE"
    }
  }
}

resource "google_folder" "basic" {
  parent       = "organizations/123456789"
  display_name = "folder"
}


```
## Example Usage - organization_policy
A test of an orgpolicy policy for an organization
```hcl
resource "google_org_policy_policy" "primary" {
  name   = "organizations/123456789/policies/gcp.detailedAuditLoggingMode"
  parent = "organizations/123456789"

  spec {
    reset = true
  }
}


```
## Example Usage - project_policy
A test of an orgpolicy policy for a project
```hcl
resource "google_org_policy_policy" "primary" {
  name   = "projects/${google_project.basic.name}/policies/gcp.resourceLocations"
  parent = "projects/${google_project.basic.name}"

  spec {
    rules {
      condition {
        description = "A sample condition for the policy"
        expression  = "resource.matchLabels('labelKeys/123', 'labelValues/345')"
        location    = "sample-location.log"
        title       = "sample-condition"
      }

      values {
        allowed_values = ["projects/allowed-project"]
        denied_values  = ["projects/denied-project"]
      }
    }

    rules {
      allow_all = "TRUE"
    }
  }
}

resource "google_project" "basic" {
  project_id = "id"
  name       = "id"
  org_id     = "123456789"
}


```
## Example Usage - dry_run_spec
```hcl
resource "google_org_policy_custom_constraint" "constraint" {
  name         = "custom.disableGkeAutoUpgrade%{random_suffix}"
  parent       = "organizations/123456789"
  display_name = "Disable GKE auto upgrade"
  description  = "Only allow GKE NodePool resource to be created or updated if AutoUpgrade is not enabled where this custom constraint is enforced."

  action_type    = "ALLOW"
  condition      = "resource.management.autoUpgrade == false"
  method_types   = ["CREATE"]
  resource_types = ["container.googleapis.com/NodePool"]
}

resource "google_org_policy_policy" "primary" {
  name   = "organizations/123456789/policies/${google_org_policy_custom_constraint.constraint.name}"
  parent = "organizations/123456789"

  spec {
    rules {
      enforce = "FALSE"
    }
  }
  dry_run_spec {
    inherit_from_parent = false
    reset               = false
    rules {
      enforce = "FALSE"
    }
  }
}

```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  Immutable. The resource name of the Policy. Must be one of the following forms, where constraint_name is the name of the constraint which this Policy configures: * `projects/{project_number}/policies/{constraint_name}` * `folders/{folder_id}/policies/{constraint_name}` * `organizations/{organization_id}/policies/{constraint_name}` For example, "projects/123/policies/compute.disableSerialPortAccess". Note: `projects/{project_id}/policies/{constraint_name}` is also an acceptable name for API requests, but responses will return the name using the equivalent project number.
  
* `parent` -
  (Required)
  The parent of the resource.
  


- - -

* `dry_run_spec` -
  (Optional)
  Dry-run policy. Audit-only policy, can be used to monitor how the policy would have impacted the existing and future resources if it's enforced.
  
* `spec` -
  (Optional)
  Basic information about the Organization Policy.
  


The `dry_run_spec` block supports:
    
* `etag` -
  An opaque tag indicating the current version of the policy, used for concurrency control. This field is ignored if used in a `CreatePolicy` request. When the policy` is returned from either a `GetPolicy` or a `ListPolicies` request, this `etag` indicates the version of the current policy to use when executing a read-modify-write loop. When the policy is returned from a `GetEffectivePolicy` request, the `etag` will be unset.
    
* `inherit_from_parent` -
  (Optional)
  Determines the inheritance behavior for this policy. If `inherit_from_parent` is true, policy rules set higher up in the hierarchy (up to the closest root) are inherited and present in the effective policy. If it is false, then no rules are inherited, and this policy becomes the new root for evaluation. This field can be set only for policies which configure list constraints.
    
* `reset` -
  (Optional)
  Ignores policies set above this resource and restores the `constraint_default` enforcement behavior of the specific constraint at this resource. This field can be set in policies for either list or boolean constraints. If set, `rules` must be empty and `inherit_from_parent` must be set to false.
    
* `rules` -
  (Optional)
  In policies for boolean constraints, the following requirements apply: - There must be one and only one policy rule where condition is unset. - Boolean policy rules with conditions must set `enforced` to the opposite of the policy rule without a condition. - During policy evaluation, policy rules with conditions that are true for a target resource take precedence.
    
* `update_time` -
  Output only. The time stamp this was previously updated. This represents the last time a call to `CreatePolicy` or `UpdatePolicy` was made for that policy.
    
The `rules` block supports:
    
* `allow_all` -
  (Optional)
  Setting this to `"TRUE"` means that all values are allowed. This field can be set only in policies for list constraints.
    
* `condition` -
  (Optional)
  A condition which determines whether this rule is used in the evaluation of the policy. When set, the `expression` field in the `Expr' must include from 1 to 10 subexpressions, joined by the "||" or "&&" operators. Each subexpression must be of the form "resource.matchTag('/tag_key_short_name, 'tag_value_short_name')". or "resource.matchTagId('tagKeys/key_id', 'tagValues/value_id')". where key_name and value_name are the resource names for Label Keys and Values. These names are available from the Tag Manager Service. An example expression is: "resource.matchTag('123456789/environment, 'prod')". or "resource.matchTagId('tagKeys/123', 'tagValues/456')".
    
* `deny_all` -
  (Optional)
  Setting this to `"TRUE"` means that all values are denied. This field can be set only in policies for list constraints.
    
* `enforce` -
  (Optional)
  If `"TRUE"`, then the policy is enforced. If `"FALSE"`, then any configuration is acceptable. This field can be set only in policies for boolean constraints.
    
* `values` -
  (Optional)
  List of values to be used for this policy rule. This field can be set only in policies for list constraints.
    
The `condition` block supports:
    
* `description` -
  (Optional)
  Optional. Description of the expression. This is a longer text which describes the expression, e.g. when hovered over it in a UI.
    
* `expression` -
  (Optional)
  Textual representation of an expression in Common Expression Language syntax.
    
* `location` -
  (Optional)
  Optional. String indicating the location of the expression for error reporting, e.g. a file name and a position in the file.
    
* `title` -
  (Optional)
  Optional. Title for the expression, i.e. a short string describing its purpose. This can be used e.g. in UIs which allow to enter the expression.
    
The `values` block supports:
    
* `allowed_values` -
  (Optional)
  List of values allowed at this resource.
    
* `denied_values` -
  (Optional)
  List of values denied at this resource.
    
The `spec` block supports:
    
* `etag` -
  An opaque tag indicating the current version of the `Policy`, used for concurrency control. This field is ignored if used in a `CreatePolicy` request. When the `Policy` is returned from either a `GetPolicy` or a `ListPolicies` request, this `etag` indicates the version of the current `Policy` to use when executing a read-modify-write loop. When the `Policy` is returned from a `GetEffectivePolicy` request, the `etag` will be unset.
    
* `inherit_from_parent` -
  (Optional)
  Determines the inheritance behavior for this `Policy`. If `inherit_from_parent` is true, PolicyRules set higher up in the hierarchy (up to the closest root) are inherited and present in the effective policy. If it is false, then no rules are inherited, and this Policy becomes the new root for evaluation. This field can be set only for Policies which configure list constraints.
    
* `reset` -
  (Optional)
  Ignores policies set above this resource and restores the `constraint_default` enforcement behavior of the specific `Constraint` at this resource. This field can be set in policies for either list or boolean constraints. If set, `rules` must be empty and `inherit_from_parent` must be set to false.
    
* `rules` -
  (Optional)
  Up to 10 PolicyRules are allowed. In Policies for boolean constraints, the following requirements apply: - There must be one and only one PolicyRule where condition is unset. - BooleanPolicyRules with conditions must set `enforced` to the opposite of the PolicyRule without a condition. - During policy evaluation, PolicyRules with conditions that are true for a target resource take precedence.
    
* `update_time` -
  Output only. The time stamp this was previously updated. This represents the last time a call to `CreatePolicy` or `UpdatePolicy` was made for that `Policy`.
    
The `rules` block supports:
    
* `allow_all` -
  (Optional)
  Setting this to `"TRUE"` means that all values are allowed. This field can be set only in Policies for list constraints.
    
* `condition` -
  (Optional)
  A condition which determines whether this rule is used in the evaluation of the policy. When set, the `expression` field in the `Expr' must include from 1 to 10 subexpressions, joined by the "||" or "&&" operators. Each subexpression must be of the form "resource.matchTag('/tag_key_short_name, 'tag_value_short_name')". or "resource.matchTagId('tagKeys/key_id', 'tagValues/value_id')". where key_name and value_name are the resource names for Label Keys and Values. These names are available from the Tag Manager Service. An example expression is: "resource.matchTag('123456789/environment, 'prod')". or "resource.matchTagId('tagKeys/123', 'tagValues/456')".
    
* `deny_all` -
  (Optional)
  Setting this to `"TRUE"` means that all values are denied. This field can be set only in Policies for list constraints.
    
* `enforce` -
  (Optional)
  If `"TRUE"`, then the `Policy` is enforced. If `"FALSE"`, then any configuration is acceptable. This field can be set only in Policies for boolean constraints.
    
* `values` -
  (Optional)
  List of values to be used for this PolicyRule. This field can be set only in Policies for list constraints.
    
The `condition` block supports:
    
* `description` -
  (Optional)
  Optional. Description of the expression. This is a longer text which describes the expression, e.g. when hovered over it in a UI.
    
* `expression` -
  (Optional)
  Textual representation of an expression in Common Expression Language syntax.
    
* `location` -
  (Optional)
  Optional. String indicating the location of the expression for error reporting, e.g. a file name and a position in the file.
    
* `title` -
  (Optional)
  Optional. Title for the expression, i.e. a short string describing its purpose. This can be used e.g. in UIs which allow to enter the expression.
    
The `values` block supports:
    
* `allowed_values` -
  (Optional)
  List of values allowed at this resource.
    
* `denied_values` -
  (Optional)
  List of values denied at this resource.
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `{{parent}}/policies/{{name}}`

* `etag` -
  Optional. An opaque tag indicating the current state of the policy, used for concurrency control. This 'etag' is computed by the server based on the value of other fields, and may be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.
  
## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Policy can be imported using any of these accepted formats:
* `{{parent}}/policies/{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Policy using one of the formats above. For example:


```tf
import {
  id = "{{parent}}/policies/{{name}}"
  to = google_org_policy_policy.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Policy can be imported using one of the formats above. For example:

```
$ terraform import google_org_policy_policy.default {{parent}}/policies/{{name}}
```



