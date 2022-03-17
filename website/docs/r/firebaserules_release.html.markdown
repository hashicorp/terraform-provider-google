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
layout: "google"
page_title: "Google: google_firebaserules_release"
description: |-
  The Firebaserules Release resource
---

# google_firebaserules_release

The Firebaserules Release resource

## Example Usage - basic_release
Creates a basic Firebase Rules Release
```hcl
resource "google_firebaserules_release" "primary" {
  name         = "release"
  ruleset_name = "projects/my-project-name/rulesets/${google_firebaserules_ruleset.basic.name}"
  project      = "my-project-name"
}

resource "google_firebaserules_ruleset" "basic" {
  source {
    files {
      content     = "service cloud.firestore {match /databases/{database}/documents { match /{document=**} { allow read, write: if false; } } }"
      name        = "firestore.rules"
      fingerprint = ""
    }

    language = ""
  }

  project = "my-project-name"
}

resource "google_firebaserules_ruleset" "minimal" {
  source {
    files {
      content = "service cloud.firestore {match /databases/{database}/documents { match /{document=**} { allow read, write: if false; } } }"
      name    = "firestore.rules"
    }
  }

  project = "my-project-name"
}


```
## Example Usage - minimal_release
Creates a minimal Firebase Rules Release
```hcl
resource "google_firebaserules_release" "primary" {
  name         = "prod/release"
  ruleset_name = "projects/my-project-name/rulesets/${google_firebaserules_ruleset.minimal.name}"
  project      = "my-project-name"
}

resource "google_firebaserules_ruleset" "minimal" {
  source {
    files {
      content = "service cloud.firestore {match /databases/{database}/documents { match /{document=**} { allow read, write: if false; } } }"
      name    = "firestore.rules"
    }
  }

  project = "my-project-name"
}


```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  Format: `projects/{project_id}/releases/{release_id}`
  
* `ruleset_name` -
  (Required)
  Name of the `Ruleset` referred to by this `Release`. The `Ruleset` must exist the `Release` to be created.
  


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
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Release can be imported using any of these accepted formats:

```
$ terraform import google_firebaserules_release.default projects/{{project}}/releases/{{name}}
```



