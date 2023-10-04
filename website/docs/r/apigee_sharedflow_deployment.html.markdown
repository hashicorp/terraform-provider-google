---
subcategory: "Apigee"
description: |-
  Deploys a revision of a sharedflow.
---

# google\_apigee\_sharedflow\_deployment

Deploys a revision of a sharedflow.


To get more information about SharedflowDeployment, see:

* [API documentation](https://cloud.google.com/apigee/docs/reference/apis/apigee/rest/v1/organizations.environments.sharedflows.revisions.deployments)
* How-to Guides
    * [sharedflows.revisions.deployments](https://cloud.google.com/apigee/docs/reference/apis/apigee/rest/v1/organizations.environments.sharedflows.revisions.deployments)

## Argument Reference

The following arguments are supported:


* `org_id` -
  (Required)
  The Apigee Organization associated with the Sharedflow

* `environment` -
  (Required)
  The resource ID of the environment.

* `sharedflow_id` -
  (Required)
  Id of the Sharedflow to be deployed.

* `revision` -
  (Required)
  Revision of the Sharedflow to be deployed.


- - -


* `service_account` -
  (Optional)
  The service account represents the identity of the deployed proxy, and determines what permissions it has. The format must be {ACCOUNT_ID}@{PROJECT}.iam.gserviceaccount.com.


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `organizations/{{org_id}}/environments/{{environment}}/sharedflows/{{sharedflow_id}}/revisions/{{revision}}/deployments`


## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import


SharedflowDeployment can be imported using any of these accepted formats:

* `organizations/{{org_id}}/environments/{{environment}}/sharedflows/{{sharedflow_id}}/revisions/{{revision}}/deployments/{{name}}`
* `{{org_id}}/{{environment}}/{{sharedflow_id}}/{{revision}}/{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import SharedflowDeployment using one of the formats above. For example:

```tf
import {
  id = "organizations/{{org_id}}/environments/{{environment}}/sharedflows/{{sharedflow_id}}/revisions/{{revision}}/deployments/{{name}}"
  to = google_apigee_flowhook.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), SharedflowDeployment can be imported using one of the formats above. For example:

```
$ terraform import google_apigee_sharedflow_deployment.default organizations/{{org_id}}/environments/{{environment}}/sharedflows/{{sharedflow_id}}/revisions/{{revision}}/deployments/{{name}}
$ terraform import google_apigee_sharedflow_deployment.default {{org_id}}/{{environment}}/{{sharedflow_id}}/{{revision}}/{{name}}
```
