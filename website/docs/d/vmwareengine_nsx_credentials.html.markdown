---
subcategory: "Cloud VMware Engine"
description: |-
  Get NSX Credentials of a Private Cloud.
---

# google\_vmwareengine\_nsx_credentials

Use this data source to get NSX credentials for a Private Cloud.

To get more information about private cloud NSX credentials, see:
* [API documentation](https://cloud.google.com/vmware-engine/docs/reference/rest/v1/projects.locations.privateClouds/showNsxCredentials)

## Example Usage

```hcl
data "google_vmwareengine_nsx_credentials" "ds" {
	parent =  "projects/my-project/locations/us-west1-a/privateClouds/my-cloud"
}
```

## Argument Reference

The following arguments are supported:

* `parent` - (Required) The resource name of the private cloud which contains the NSX.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `username` - The username of the NSX Credential.
* `password` - The password of the NSX Credential.