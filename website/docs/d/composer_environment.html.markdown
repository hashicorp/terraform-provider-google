---
subcategory: "Cloud Composer"
description: |-
  Provides Cloud Composer environment configuration data.
---

# google_composer_environment

Provides access to Cloud Composer environment configuration in a region for a given project.

## Example Usage

```hcl
resource "google_composer_environment" "composer_env" {
    name = "composer-environment"
}

data "google_composer_environment" "composer_env" {
    name = google_composer_environment.test.name

    depends_on = [google_composer_environment.composer_env]
}

output "debug" {
    value = data.google_composer_environment.composer_env.config
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the environment.

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

* `region` - (Optional) The location or Compute Engine region of the environment.

## Attributes Reference

The following attributes are exported:

* `id` - An identifier for the resource in format `projects/{{project}}/locations/{{region}}/environments/{{name}}`

* `config` - Configuration parameters for the environment.
    Full structure is provided by [composer environment resource documentation](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/composer_environment#config).

    * `config.0.gke_cluster` -
    The Kubernetes Engine cluster used to run the environment.

    * `config.0.dag_gcs_prefix` -
    The Cloud Storage prefix of the DAGs for the environment.

    * `config.0.airflow_uri` -
    The URI of the Apache Airflow Web UI hosted within the
    environment.
