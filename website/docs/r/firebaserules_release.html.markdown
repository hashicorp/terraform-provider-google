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
subcategory: "Firebaserules"
description: |-
  
---

# google_firebaserules_release



For more information, see:
* [Get started with Firebase Security Rules](https://firebase.google.com/docs/rules/get-started)
## Example Usage - firestore_release
Creates a Firebase Rules Release to Cloud Firestore
```hcl
resource "google_firebaserules_release" "primary" {
  name         = "cloud.firestore"
  ruleset_name = "projects/my-project-name/rulesets/${google_firebaserules_ruleset.firestore.name}"
  project      = "my-project-name"

  lifecycle {
    replace_triggered_by = [
      google_firebaserules_ruleset.firestore
    ]
  }
}

resource "google_firebaserules_ruleset" "firestore" {
  source {
    files {
      content = "service cloud.firestore {match /databases/{database}/documents { match /{document=**} { allow read, write: if false; } } }"
      name    = "firestore.rules"
    }
  }

  project = "my-project-name"
}

```
## Example Usage - storage_release
Creates a Firebase Rules Release for a Storage bucket
```hcl
resource "google_firebaserules_release" "primary" {
  provider     = google-beta
  name         = "firebase.storage/${google_storage_bucket.bucket.name}"
  ruleset_name = "projects/my-project-name/rulesets/${google_firebaserules_ruleset.storage.name}"
  project      = "my-project-name"

  lifecycle {
    replace_triggered_by = [
      google_firebaserules_ruleset.storage
    ]
  }
}

# Provision a non-default Cloud Storage bucket.
resource "google_storage_bucket" "bucket" {
  provider = google-beta
  project  = "my-project-name"
  name     = "bucket"
  location = "us-west1"
}

# Make the Storage bucket accessible for Firebase SDKs, authentication, and Firebase Security Rules.
resource "google_firebase_storage_bucket" "bucket" {
  provider  = google-beta
  project   = "my-project-name"
  bucket_id = google_storage_bucket.bucket.name
}

# Create a ruleset of Firebase Security Rules from a local file.
resource "google_firebaserules_ruleset" "storage" {
  provider = google-beta
  project  = "my-project-name"
  source {
    files {
      name    = "storage.rules"
      content = "service firebase.storage {match /b/{bucket}/o {match /{allPaths=**} {allow read, write: if request.auth != null;}}}"
    }
  }

  depends_on = [
    google_firebase_storage_bucket.bucket
  ]
}

```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  Format: `projects/{project_id}/releases/{release_id}`\Firestore Rules Releases will **always** have the name 'cloud.firestore'
  
* `ruleset_name` -
  (Required)
  Name of the `Ruleset` referred to by this `Release`. The `Ruleset` must exist for the `Release` to be created.
  


- - -

* `project` -
  (Optional)
  The project for the resource
  


## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/releases/{{name}}`

* `create_time` -
  Output only. Time the release was created.
  
* `disabled` -
  Disable the release to keep it from being served. The response code of NOT_FOUND will be given for executables generated from this Release.
  
* `update_time` -
  Output only. Time the release was updated.
  
## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Release can be imported using any of these accepted formats:
* `projects/{{project}}/releases/{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Release using one of the formats above. For example:


```tf
import {
  id = "projects/{{project}}/releases/{{name}}"
  to = google_firebaserules_release.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Release can be imported using one of the formats above. For example:

```
$ terraform import google_firebaserules_release.default projects/{{project}}/releases/{{name}}
```



