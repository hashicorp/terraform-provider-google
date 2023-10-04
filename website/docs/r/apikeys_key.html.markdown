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
subcategory: "Apikeys"
description: |-
  The Apikeys Key resource
---

# google_apikeys_key

The Apikeys Key resource

## Example Usage - android_key
A basic example of a android api keys key
```hcl
resource "google_apikeys_key" "primary" {
  name         = "key"
  display_name = "sample-key"
  project      = google_project.basic.name

  restrictions {
    android_key_restrictions {
      allowed_applications {
        package_name     = "com.example.app123"
        sha1_fingerprint = "1699466a142d4682a5f91b50fdf400f2358e2b0b"
      }
    }

    api_targets {
      service = "translate.googleapis.com"
      methods = ["GET*"]
    }
  }
}

resource "google_project" "basic" {
  project_id = "app"
  name       = "app"
  org_id     = "123456789"
}


```
## Example Usage - basic_key
A basic example of a api keys key
```hcl
resource "google_apikeys_key" "primary" {
  name         = "key"
  display_name = "sample-key"
  project      = google_project.basic.name

  restrictions {
    api_targets {
      service = "translate.googleapis.com"
      methods = ["GET*"]
    }

    browser_key_restrictions {
      allowed_referrers = [".*"]
    }
  }
}

resource "google_project" "basic" {
  project_id = "app"
  name       = "app"
  org_id     = "123456789"
}


```
## Example Usage - ios_key
A basic example of a ios api keys key
```hcl
resource "google_apikeys_key" "primary" {
  name         = "key"
  display_name = "sample-key"
  project      = google_project.basic.name

  restrictions {
    api_targets {
      service = "translate.googleapis.com"
      methods = ["GET*"]
    }

    ios_key_restrictions {
      allowed_bundle_ids = ["com.google.app.macos"]
    }
  }
}

resource "google_project" "basic" {
  project_id = "app"
  name       = "app"
  org_id     = "123456789"
}


```
## Example Usage - minimal_key
A minimal example of a api keys key
```hcl
resource "google_apikeys_key" "primary" {
  name         = "key"
  display_name = "sample-key"
  project      = google_project.basic.name
}

resource "google_project" "basic" {
  project_id = "app"
  name       = "app"
  org_id     = "123456789"
}


```
## Example Usage - server_key
A basic example of a server api keys key
```hcl
resource "google_apikeys_key" "primary" {
  name         = "key"
  display_name = "sample-key"
  project      = google_project.basic.name

  restrictions {
    api_targets {
      service = "translate.googleapis.com"
      methods = ["GET*"]
    }

    server_key_restrictions {
      allowed_ips = ["127.0.0.1"]
    }
  }
}

resource "google_project" "basic" {
  project_id = "app"
  name       = "app"
  org_id     = "123456789"
}


```

## Argument Reference

The following arguments are supported:

* `name` -
  (Required)
  The resource name of the key. The name must be unique within the project, must conform with RFC-1034, is restricted to lower-cased letters, and has a maximum length of 63 characters. In another word, the name must match the regular expression: `[a-z]([a-z0-9-]{0,61}[a-z0-9])?`.
  


The `allowed_applications` block supports:
    
* `package_name` -
  (Required)
  The package name of the application.
    
* `sha1_fingerprint` -
  (Required)
  The SHA1 fingerprint of the application. For example, both sha1 formats are acceptable : DA:39:A3:EE:5E:6B:4B:0D:32:55:BF:EF:95:60:18:90:AF:D8:07:09 or DA39A3EE5E6B4B0D3255BFEF95601890AFD80709. Output format is the latter.
    
- - -

* `display_name` -
  (Optional)
  Human-readable display name of this API key. Modifiable by user.
  
* `project` -
  (Optional)
  The project for the resource
  
* `restrictions` -
  (Optional)
  Key restrictions.
  


The `restrictions` block supports:
    
* `android_key_restrictions` -
  (Optional)
  The Android apps that are allowed to use the key.
    
* `api_targets` -
  (Optional)
  A restriction for a specific service and optionally one or more specific methods. Requests are allowed if they match any of these restrictions. If no restrictions are specified, all targets are allowed.
    
* `browser_key_restrictions` -
  (Optional)
  The HTTP referrers (websites) that are allowed to use the key.
    
* `ios_key_restrictions` -
  (Optional)
  The iOS apps that are allowed to use the key.
    
* `server_key_restrictions` -
  (Optional)
  The IP addresses of callers that are allowed to use the key.
    
The `android_key_restrictions` block supports:
    
* `allowed_applications` -
  (Required)
  A list of Android applications that are allowed to make API calls with this key.
    
The `api_targets` block supports:
    
* `methods` -
  (Optional)
  Optional. List of one or more methods that can be called. If empty, all methods for the service are allowed. A wildcard (*) can be used as the last symbol. Valid examples: `google.cloud.translate.v2.TranslateService.GetSupportedLanguage` `TranslateText` `Get*` `translate.googleapis.com.Get*`
    
* `service` -
  (Required)
  The service for this restriction. It should be the canonical service name, for example: `translate.googleapis.com`. You can use `gcloud services list` to get a list of services that are enabled in the project.
    
The `browser_key_restrictions` block supports:
    
* `allowed_referrers` -
  (Required)
  A list of regular expressions for the referrer URLs that are allowed to make API calls with this key.
    
The `ios_key_restrictions` block supports:
    
* `allowed_bundle_ids` -
  (Required)
  A list of bundle IDs that are allowed when making API calls with this key.
    
The `server_key_restrictions` block supports:
    
* `allowed_ips` -
  (Required)
  A list of the caller IP addresses that are allowed to make API calls with this key.
    
## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/global/keys/{{name}}`

* `key_string` -
  Output only. An encrypted and signed value held by this key. This field can be accessed only through the `GetKeyString` method.
  
* `uid` -
  Output only. Unique id in UUID4 format.
  
## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import

Key can be imported using any of these accepted formats:
* `projects/{{project}}/locations/global/keys/{{name}}`
* `{{project}}/{{name}}`
* `{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Key using one of the formats above. For example:


```tf
import {
  id = "projects/{{project}}/locations/global/keys/{{name}}"
  to = google_apikeys_key.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Key can be imported using one of the formats above. For example:

```
$ terraform import google_apikeys_key.default projects/{{project}}/locations/global/keys/{{name}}
$ terraform import google_apikeys_key.default {{project}}/{{name}}
$ terraform import google_apikeys_key.default {{name}}
```



