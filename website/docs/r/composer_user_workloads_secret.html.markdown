---
subcategory: "Cloud Composer"
description: |-
  User workloads Secret used by Airflow tasks that run with Kubernetes Executor or KubernetesPodOperator.
---

# google_composer_user_workloads_secret

~> **Warning:** These resources are in beta, and should be used with the terraform-provider-google-beta provider.
See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html) for more details on beta resources.

User workloads Secret used by Airflow tasks that run with Kubernetes Executor or KubernetesPodOperator. 
Intended for Composer 3 Environments.

## Example Usage

```hcl
resource "google_composer_environment" "example" {
  name              = "example-environment"
  project           = "example-project"
  region            = "us-central1"
  config {
    software_config {
      image_version = "example-image-version"
    }
  }
}

resource "google_composer_user_workloads_secret" "example" {
  name = "example-secret"
  project = "example-project"
  region = "us-central1"
  environment = google_composer_environment.example.name
  data = {
    email: base64encode("example-email"),
    password: base64encode("example-password"),
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  Name of the Kubernetes Secret.

* `region` -
  (Optional)
  The location or Compute Engine region for the environment.

* `project` -
  (Optional)
  The ID of the project in which the resource belongs.
  If it is not provided, the provider project is used.

* `environment` -
  Environment where the Kubernetes Secret will be stored and used.

* `data` -
  (Optional) 
  The "data" field of Kubernetes Secret, organized in key-value pairs,
  which can contain sensitive values such as a password, a token, or a key. 
  Content of this field will not be displayed in CLI output, 
  but it will be stored in terraform state file. To protect sensitive data, 
  follow the best practices outlined in the HashiCorp documentation: 
  https://developer.hashicorp.com/terraform/language/state/sensitive-data.
  The values for all keys have to be base64-encoded strings. 
  For details see: https://kubernetes.io/docs/concepts/configuration/secret/



## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{region}}/environments/{{environment}}/userWorkloadsSecrets/{{name}}`

## Import

Secret can be imported using any of these accepted formats:

* `projects/{{project}}/locations/{{region}}/environments/{{environment}}/userWorkloadsSecrets/{{name}}`
* `{{project}}/{{region}}/{{environment}}/{{name}}`
* `{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import User Workloads Secret using one of the formats above. For example:

```tf
import {
  id = "projects/{{project}}/locations/{{region}}/environments/{{environment}}/userWorkloadsSecrets/{{name}}"
  to = google_composer_user_workloads_secret.example
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Environment can be imported using one of the formats above. For example:

```
$ terraform import google_composer_user_workloads_secret.example projects/{{project}}/locations/{{region}}/environments/{{environment}}/userWorkloadsSecrets/{{name}}
$ terraform import google_composer_user_workloads_secret.example {{project}}/{{region}}/{{environment}}/{{name}}
$ terraform import google_composer_user_workloads_secret.example {{name}}
```
