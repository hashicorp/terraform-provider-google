---
subcategory: "Cloud VMware Engine"
description: |-
  Get information about a external access rule.
---

# google_vmwareengine_external_access_rule

Use this data source to get details about a external access rule resource.

To get more information about external address, see:
* [API documentation](https://cloud.google.com/vmware-engine/docs/reference/rest/v1/projects.locations.networkPolicies.externalAccessRules)

## Example Usage

```hcl
data "google_vmwareengine_external_access_rule" "my_external_access_rule" {
  name     = "my-external-access-rule"
  parent   = "project/my-project/locations/us-west1-a/networkPolicies/my-network-policy"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the resource.
* `parent` - (Required) The resource name of the network policy that this cluster belongs.

## Attributes Reference

See [google_vmwareengine_external_access_rule](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/vmwareengine_external_access_rule#attributes-reference) resource for details of all the available attributes.