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
  Beta only: The Cloudbuildv2 Connection resource
---

# google_cloudbuildv2_connection

Beta only: The Cloudbuildv2 Connection resource

## Example Usage - ghe
```hcl
resource "google_secret_manager_secret" "private-key-secret" {
  provider = google-beta
  secret_id = "ghe-pk-secret"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "private-key-secret-version" {
  provider = google-beta
  secret = google_secret_manager_secret.private-key-secret.id
  secret_data = file("private-key.pem")
}

resource "google_secret_manager_secret" "webhook-secret-secret" {
  provider = google-beta
  secret_id = "github-token-secret"

  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "webhook-secret-secret-version" {
  provider = google-beta
  secret = google_secret_manager_secret.webhook-secret-secret.id
  secret_data = "<webhook-secret-data>"
}

data "google_iam_policy" "p4sa-secretAccessor" {
  provider = google-beta
  binding {
    role = "roles/secretmanager.secretAccessor"
    // Here, 123456789 is the Google Cloud project number for the project that contains the connection.
    members = ["serviceAccount:service-123456789@gcp-sa-cloudbuild.iam.gserviceaccount.com"]
  }
}

resource "google_secret_manager_secret_iam_policy" "policy-pk" {
  provider = google-beta
  secret_id = google_secret_manager_secret.private-key-secret.secret_id
  policy_data = data.google_iam_policy.p4sa-secretAccessor.policy_data
}

resource "google_secret_manager_secret_iam_policy" "policy-whs" {
  provider = google-beta
  secret_id = google_secret_manager_secret.webhook-secret-secret.secret_id
  policy_data = data.google_iam_policy.p4sa-secretAccessor.policy_data
}

resource "google_cloudbuildv2_connection" "my-connection" {
  provider = google-beta
  location = "us-central1"
  name = "my-terraform-ghe-connection"

  github_enterprise_config {
    host_uri = "https://ghe.com"
    private_key_secret_version = google_secret_manager_secret_version.private-key-secret-version.id
    webhook_secret_secret_version = google_secret_manager_secret_version.webhook-secret-secret-version.id
    app_id = 200
    app_slug = "gcb-app"
    app_installation_id = 300
  }

  depends_on = [
    google_secret_manager_secret_iam_policy.policy-pk,
    google_secret_manager_secret_iam_policy.policy-whs
  ]
}

```
## Example Usage - GitHub Connection
Creates a Connection to github.com
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

```

## Argument Reference

The following arguments are supported:

* `location` -
  (Required)
  The location for the resource
  
* `name` -
  (Required)
  Immutable. The resource name of the connection, in the format `projects/{project}/locations/{location}/connections/{connection_id}`.
  


- - -

* `annotations` -
  (Optional)
  Allows clients to store small amounts of arbitrary data.
  
* `disabled` -
  (Optional)
  If disabled is set to true, functionality is disabled for this connection. Repository based API methods and webhooks processing for repositories in this connection will be disabled.
  
* `github_config` -
  (Optional)
  Configuration for connections to github.com.
  
* `github_enterprise_config` -
  (Optional)
  Configuration for connections to an instance of GitHub Enterprise.
  
* `project` -
  (Optional)
  The project for the resource
  


The `github_config` block supports:
    
* `app_installation_id` -
  (Optional)
  GitHub App installation id.
    
* `authorizer_credential` -
  (Optional)
  OAuth credential of the account that authorized the Cloud Build GitHub App. It is recommended to use a robot account instead of a human user account. The OAuth token must be tied to the Cloud Build GitHub App.
    
The `authorizer_credential` block supports:
    
* `oauth_token_secret_version` -
  (Optional)
  A SecretManager resource containing the OAuth token that authorizes the Cloud Build connection. Format: `projects/*/secrets/*/versions/*`.
    
* `username` -
  Output only. The username associated to this token.
    
The `github_enterprise_config` block supports:
    
* `host_uri` -
  (Required)
  Required. The URI of the GitHub Enterprise host this connection is for.
    
* `app_id` -
  (Optional)
  Id of the GitHub App created from the manifest.
    
* `app_installation_id` -
  (Optional)
  ID of the installation of the GitHub App.
    
* `app_slug` -
  (Optional)
  The URL-friendly name of the GitHub App.
    
* `private_key_secret_version` -
  (Optional)
  SecretManager resource containing the private key of the GitHub App, formatted as `projects/*/secrets/*/versions/*`.
    
* `service_directory_config` -
  (Optional)
  Configuration for using Service Directory to privately connect to a GitHub Enterprise server. This should only be set if the GitHub Enterprise server is hosted on-premises and not reachable by public internet. If this field is left empty, calls to the GitHub Enterprise server will be made over the public internet.
    
* `ssl_ca` -
  (Optional)
  SSL certificate to use for requests to GitHub Enterprise.
    
* `webhook_secret_secret_version` -
  (Optional)
  SecretManager resource containing the webhook secret of the GitHub App, formatted as `projects/*/secrets/*/versions/*`.
    
The `service_directory_config` block supports:
    
* `service` -
  (Required)
  Required. The Service Directory service name. Format: projects/{project}/locations/{location}/namespaces/{namespace}/services/{service}.
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/connections/{{name}}`

* `create_time` -
  Output only. Server assigned timestamp for when the connection was created.
  
* `etag` -
  This checksum is computed by the server based on the value of other fields, and may be sent on update and delete requests to ensure the client has an up-to-date value before proceeding.
  
* `installation_state` -
  Output only. Installation state of the Connection.
  
* `reconciling` -
  Output only. Set to true when the connection is being set up or updated in the background.
  
* `update_time` -
  Output only. Server assigned timestamp for when the connection was updated.
  
## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Connection can be imported using any of these accepted formats:

```
$ terraform import google_cloudbuildv2_connection.default projects/{{project}}/locations/{{location}}/connections/{{name}}
$ terraform import google_cloudbuildv2_connection.default {{project}}/{{location}}/{{name}}
$ terraform import google_cloudbuildv2_connection.default {{location}}/{{name}}
```



