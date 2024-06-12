---
subcategory: "Cloud Composer"
description: |-
  User workloads ConfigMap used by Airflow tasks that run with Kubernetes Executor or KubernetesPodOperator.
---

# google\_composer\_user\_workloads\_config\_map

Provides access to Kubernetes ConfigMap configuration for a given project, region and Composer Environment.

> **Warning:** This data source is in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

## Example Usage

```hcl
resource "google_composer_environment" "example" {
    name = "example-environment"
    config{
        software_config {
            image_version = "composer-3-airflow-2"
        }
    }
}
resource "google_composer_user_workloads_config_map" "example" {
    environment = google_composer_environment.example.name
    name = "example-config-map"
    data = {
        db_host: "dbhost:5432",
        api_host: "apihost:443",
    }
}
data "google_composer_user_workloads_config_map" "example" {
    environment = google_composer_environment.example.name
    name = resource.google_composer_user_workloads_config_map.example.name
}
output "debug" {
    value = data.google_composer_user_workloads_config_map.example
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the ConfigMap.

* `environment` - (Required) Environment where the ConfigMap is stored.

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

* `region` - (Optional) The location or Compute Engine region of the environment.

## Attributes Reference

The following attributes are exported:

* `id` - An identifier for the resource in format `projects/{{project}}/locations/{{region}}/environments/{{environment}}/userWorkloadsConfigMaps/{{name}}`

* `data` - The "data" field of Kubernetes ConfigMap, organized in key-value pairs.
      For details see: https://kubernetes.io/docs/concepts/configuration/configmap/
