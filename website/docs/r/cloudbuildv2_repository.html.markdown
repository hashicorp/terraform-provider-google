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
subcategory: "Cloud Build v2"
description: |-
  Beta only: The Cloudbuildv2 Repository resource
---

# google_cloudbuildv2_repository

Beta only: The Cloudbuildv2 Repository resource

## Example Usage - Repository in GitHub Connection
Creates a Repository resource inside a Connection to github.com
```hcl
resource "google_secret_manager_secret" "github-token-secret" {
  provider = google-beta
  secret_id = "github-token-secret"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "github-token-secret-version" {
  provider = google-beta
  secret = google_secret_manager_secret.github-token-secret.id
  secret_data = file("my-github-token.txt")
}

data "google_iam_policy" "p4sa-secretAccessor" {
  provider = google-beta
  binding {
    role = "roles/secretmanager.secretAccessor"
    // Here, 123456789 is the Google Cloud project number for my-project-name.
    members = ["serviceAccount:service-123456789@gcp-sa-cloudbuild.iam.gserviceaccount.com"]
  }
}

resource "google_secret_manager_secret_iam_policy" "policy" {
  provider = google-beta
  secret_id = google_secret_manager_secret.github-token-secret.secret_id
  policy_data = data.google_iam_policy.p4sa-secretAccessor.policy_data
}

resource "google_cloudbuildv2_connection" "my-connection" {
  provider = google-beta
  location = "us-west1"
  name = "my-connection"

  github_config {
    app_installation_id = 123123
    authorizer_credential {
      oauth_token_secret_version = google_secret_manager_secret_version.github-token-secret-version.id
    }
  }
}

resource "google_cloudbuildv2_repository" "my-repository" {
  provider = google-beta
  location = "us-west1"
  name = "my-repo"
  parent_connection = google_cloudbuildv2_connection.my-connection.name
  remote_uri = "https://github.com/myuser/myrepo.git"
}

```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  Name of the repository.
  
* `parent_connection` -
  (Required)
  The connection for the resource
  
* `remote_uri` -
  (Required)
  Required. Git Clone HTTPS URI.
  


- - -

* `annotations` -
  (Optional)
  Allows clients to store small amounts of arbitrary data.
  
* `location` -
  (Optional)
  The location for the resource
  
* `project` -
  (Optional)
  The project for the resource
  


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/connections/{{parent_connection}}/repositories/{{name}}`

* `create_time` -
  Output only. Server assigned timestamp for when the connection was created.
  
* `etag` -
  This checksum is computed by the server based on the value of other fields, and may be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.
  
* `update_time` -
  Output only. Server assigned timestamp for when the connection was updated.
  
## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Repository can be imported using any of these accepted formats:

```
$ terraform import google_cloudbuildv2_repository.default projects/{{project}}/locations/{{location}}/connections/{{parent_connection}}/repositories/{{name}}
$ terraform import google_cloudbuildv2_repository.default {{project}}/{{location}}/{{parent_connection}}/{{name}}
$ terraform import google_cloudbuildv2_repository.default {{location}}/{{parent_connection}}/{{name}}
```



