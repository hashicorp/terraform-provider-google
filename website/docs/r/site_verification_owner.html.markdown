subcategory: "Site Verification"
description: |-
  Manages additional owners on verified web resources.
---

# google_site_verification_owner

An owner is an additional user that may manage a verified web site in the
[Google Search Console](https://www.google.com/webmasters/tools/). There
are two types of web resource owners:

* Verified owners, which are added to a web resource automatically when it
    is created (i.e., when the resource is verified). A verified owner is
    determined by the identity of the user requesting verification.
* Additional owners, which can be added to the resource by verified owners.

`google_site_verification_owner` creates additional owners. If your web site
was verified using the
[`google_site_verification_web_resource`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/google_site_verification_web_resource)
resource then you (or the identity was used to create the resource, such as a
service account) are already an owner.

~> **Note:** The email address of the owner must belong to a Google account,
such as a Gmail account, a Google Workspace account, or a GCP service account.

Working with site verification requires the `https://www.googleapis.com/auth/siteverification`
authentication scope. See the
[Google Provider authentication documentation](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#authentication)
to learn how to configure additional scopes.

To get more information about site owners, see:

* [API documentation](https://developers.google.com/site-verification/v1)
* How-to Guides
    * [Getting Started](https://developers.google.com/site-verification/v1/getting_started)

## Example Usage - Site Verification Storage Bucket

This example uses the `FILE` verification method to verify ownership of web site hosted
in a Google Cloud Storage bucket. Ownership is proved by creating a file with a Google-provided
value in a known location. The user applying this configuration will automatically be
added as a verified owner, and the `google_site_verification_owner` resource will add
`user@example.com` as an additional owner.

```hcl
resource "google_storage_bucket" "bucket" {
  name     = "example-storage-bucket"
  location = "US"
}

data "google_site_verification_token" "token" {
  type                = "SITE"
  identifier          = "https://${google_storage_bucket.bucket.name}.storage.googleapis.com/"
  verification_method = "FILE"
}

resource "google_storage_bucket_object" "object" {
  name     = "${data.google_site_verification_token.token.token}"
  content  = "google-site-verification: ${data.google_site_verification_token.token.token}"
  bucket   = google_storage_bucket.bucket.name
}

resource "google_storage_object_access_control" "public_rule" {
  bucket   = google_storage_bucket.bucket.name
  object   = google_storage_bucket_object.object.name
  role     = "READER"
  entity   = "allUsers"
}

resource "google_site_verification_web_resource" "example" {
  site {
    type       = data.google_site_verification_token.token.type
    identifier = data.google_site_verification_token.token.identifier
  }
  verification_method = data.google_site_verification_token.token.verification_method
}

resource "google_site_verification_owner" "example" {
  web_resource_id = google_site_verification_web_resource.example.id
  email           = "user@example.com"
}
```

## Argument Reference

The following arguments are supported:


* `web_resource_id` -
  (Required)
  The id of of the web resource to which the owner will be added, in the form `webResource/<resource_id>`,
  such as `webResource/https://www.example.com/`

* `email` -
  (Required)
  The email of the user to be added as an owner.

- - -


## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import


Owner can be imported using this format:

* `webResource/{{web_resource_id}}/{{email}}`


In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import a site owner using the format above. For example:

```tf
import {
  id = "webResource/{{web_resource_id}}/{{email}}"
  to = google_site_verification_web_resource.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Site owners can be imported using the format above. For example:

```
$ terraform import google_site_verification_web_resource.default webResource/{{web_resource_id}}/{{email}}
```

~> **Note:** While verified owners can be successfully imported, attempting to later delete the imported resource will fail. The only way to remove
verified owners is to delete the web resource itself.

## User Project Overrides

This resource supports [User Project Overrides](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#user_project_override).