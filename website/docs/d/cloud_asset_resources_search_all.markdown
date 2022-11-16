---
subcategory: "Cloud Asset Inventory"
page_title: "Google: google_cloud_asset_resources_search_all"
description: |-
  Retrieve all the resources within a given accessible CRM scope (project/folder/organization).
---

# google\_cloud\_asset\_resources\_search\_all

Retrieve all the resources within a given accessible CRM scope (project/folder/organization). See the
[REST API](https://cloud.google.com/asset-inventory/docs/reference/rest/v1p1beta1/resources/searchAll)
for more details.

~> **Warning:** This resource is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

## Example Usage - searching for all projects in an org

```hcl
data google_cloud_asset_resources_search_all projects {
  provider = google-beta
  scope = "organizations/0123456789"
  asset_types = [
    "cloudresourcemanager.googleapis.com/Project"
  ]
}
```

## Example Usage - searching for all projects with CloudBuild API enabled

```hcl
data google_cloud_asset_resources_search_all cloud_build_projects {
  provider = google-beta
  scope = "organizations/0123456789"
  asset_types = [
    "serviceusage.googleapis.com/Service"
  ]
  query = "displayName:cloudbuild.googleapis.com AND state:ENABLED"
}
```

## Example Usage - searching for all Service Accounts in a project

```hcl
data google_cloud_asset_resources_search_all project_service_accounts {
  provider = google-beta
  scope = "projects/my-project-id"
  asset_types = [
    "iam.googleapis.com/ServiceAccount"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `scope` - (Required) A scope can be a project, a folder, or an organization. The allowed value must be: organization number (such as "organizations/123"), folder number (such as "folders/1234"), project number (such as "projects/12345") or project id (such as "projects/abc")
* `asset_types` - (Optional) A list of asset types that this request searches for. If empty, it will search all the [supported asset types](https://cloud.google.com/asset-inventory/docs/supported-asset-types). 
* `query` - (Optional) The query statement. See [how to construct a query](https://cloud.google.com/asset-inventory/docs/searching-resources#how_to_construct_a_query) for more information. If not specified or empty, it will search all the resources within the specified `scope` and `asset_types`.


## Attributes Reference

The following attributes are exported:

* `results` - A list of search results based on provided inputs. Structure is [defined below](#nested_results).

<a name="nested_results"></a>The `results` block supports:

* `name` - The full resource name. See [Resource Names](https://cloud.google.com/apis/design/resource_names#full_resource_name) for more information.
* `asset_type` - The type of this resource. 
* `project` - The project that this resource belongs to, in the form of `projects/{project_number}`.
* `display_name` - The display name of this resource.
* `description` - One or more paragraphs of text description of this resource. Maximum length could be up to 1M bytes.
* `additional_attributes` - Additional searchable attributes of this resource. Informational only. The exact set of attributes is subject to change. For example: project id, DNS name etc.
* `location` - Location can be `global`, regional like `us-east1`, or zonal like `us-west1-b`.
* `labels` - Labels associated with this resource.
* `network_tags` - Network tags associated with this resource.
