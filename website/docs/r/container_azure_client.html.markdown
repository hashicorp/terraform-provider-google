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
subcategory: "ContainerAzure"
description: |-
  AzureClient resources hold client authentication information needed by the Anthos Multi-Cloud API to manage Azure resources on your Azure subscription.When an AzureCluster is created, an AzureClient resource needs to be provided and all operations on Azure resources associated to that cluster will authenticate to Azure services using the given client.AzureClient resources are immutable and cannot be modified upon creation.Each AzureClient resource is bound to a single Azure Active Directory Application and tenant.
---

# google_container_azure_client

AzureClient resources hold client authentication information needed by the Anthos Multi-Cloud API to manage Azure resources on your Azure subscription.When an AzureCluster is created, an AzureClient resource needs to be provided and all operations on Azure resources associated to that cluster will authenticate to Azure services using the given client.AzureClient resources are immutable and cannot be modified upon creation.Each AzureClient resource is bound to a single Azure Active Directory Application and tenant.

For more information, see:
* [Multicloud overview](https://cloud.google.com/anthos/clusters/docs/multi-cloud)
## Example Usage - basic_azure_client
A basic example of a containerazure azure client
```hcl
resource "google_container_azure_client" "primary" {
  application_id = "12345678-1234-1234-1234-123456789111"
  location       = "us-west1"
  name           = "client-name"
  tenant_id      = "12345678-1234-1234-1234-123456789111"
  project        = "my-project-name"
}

```

## Argument Reference

The following arguments are supported:

* `application_id` -
  (Required)
  The Azure Active Directory Application ID.
  
* `location` -
  (Required)
  The location for the resource
  
* `name` -
  (Required)
  The name of this resource.
  
* `tenant_id` -
  (Required)
  The Azure Active Directory Tenant ID.
  


- - -

* `project` -
  (Optional)
  The project for the resource
  


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/azureClients/{{name}}`

* `certificate` -
  Output only. The PEM encoded x509 certificate.
  
* `create_time` -
  Output only. The time at which this resource was created.
  
* `uid` -
  Output only. A globally unique identifier for the client.
  
## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Client can be imported using any of these accepted formats:
* `projects/{{project}}/locations/{{location}}/azureClients/{{name}}`
* `{{project}}/{{location}}/{{name}}`
* `{{location}}/{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Client using one of the formats above. For example:


```tf
import {
  id = "projects/{{project}}/locations/{{location}}/azureClients/{{name}}"
  to = google_container_azure_client.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Client can be imported using one of the formats above. For example:

```
$ terraform import google_container_azure_client.default projects/{{project}}/locations/{{location}}/azureClients/{{name}}
$ terraform import google_container_azure_client.default {{project}}/{{location}}/{{name}}
$ terraform import google_container_azure_client.default {{location}}/{{name}}
```



