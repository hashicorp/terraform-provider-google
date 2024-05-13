---
subcategory: "Filestore"
description: |-
  Get information about a Google Cloud Filestore instance.
---

# google_filestore_instance

Get info about a Google Cloud Filestore instance.

## Example Usage

```tf
data "google_filestore_instance" "my_instance" {
  name = "my-filestore-instance"
}

output "instance_ip_addresses" {
  value = data.google_filestore_instance.my_instance.networks.ip_addresses
}

output "instance_connect_mode" {
  value = data.google_filestore_instance.my_instance.networks.connect_mode
}

output "instance_file_share_name" {
  value = data.google_filestore_instance.my_instance.file_shares.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of a Filestore instance.

- - -

* `project` - (Optional) The project in which the resource belongs. If it
    is not provided, the provider project is used.

* `location` - (Optional) The name of the location of the instance. This 
    can be a region for ENTERPRISE tier instances. If it is not provided, 
    the provider region or zone is used.

## Attributes Reference

See [google_filestore_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/filestore_instance) resource for details of the available attributes.
