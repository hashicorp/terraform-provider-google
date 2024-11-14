---
subcategory: "Cloud Composer"
description: |-
  User workloads Secret used by Airflow tasks that run with Kubernetes Executor or KubernetesPodOperator.
---

# google_composer_user_workloads_secret

Provides access to Kubernetes Secret configuration for a given project, region and Composer Environment.

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

resource "google_composer_user_workloads_secret" "example" {
    environment = google_composer_environment.example.name
    name = "example-secret"
    data = {
        username: base64encode("username"),
        password: base64encode("password"),
    }
}

data "google_composer_user_workloads_secret" "example" {
    environment = google_composer_environment.example.name
    name = resource.google_composer_user_workloads_secret.example.name
}

output "debug" {
    value = data.google_composer_user_workloads_secret.example
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Secret.

* `environment` - (Required) Environment where the Secret is stored.

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

* `region` - (Optional) The location or Compute Engine region of the environment.

## Attributes Reference

See [google_composer_user_workloads_secret](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/composer_user_workloads_secret) resource for details of the available attributes.
