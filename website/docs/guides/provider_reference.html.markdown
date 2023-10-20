---
page_title: "Google Provider Configuration Reference"
description: |-
  Configuration reference for the Google provider for Terraform.
---

# Google Provider Configuration Reference

The `google` and `google-beta` provider blocks are used to configure the
credentials you use to authenticate with GCP, as well as a default project and
location (`zone` and/or `region`) for your resources. The same values are
available between the provider versions, but must be configured in separate
provider blocks.

### Example Usage - Basic provider blocks

```hcl
provider "google" {
  project     = "my-project-id"
  region      = "us-central1"
  zone        = "us-central1-c"
}
```

```hcl
provider "google-beta" {
  project     = "my-project-id"
  region      = "us-central1"
  zone        = "us-central1-c"
}
```

### Example Usage - Using beta features with `google-beta`

To use Google Cloud Platform features that are in beta, you need to both:

* Explicitly define a `google-beta` provider block

* explicitly set the provider for your resource to `google-beta`.

See [Provider Versions](https://terraform.io/docs/providers/google/guides/provider_versions.html)
for a full reference on how to use features from different GCP API versions in
the Google provider.

```hcl
resource "google_compute_instance" "ga-instance" {
  provider = google

  # ...
}

resource "google_compute_instance" "beta-instance" {
  provider = google-beta

  # ...
}

provider "google-beta" {}
```

## Authentication

### Running Terraform on your workstation.

If you are using Terraform on your workstation we recommend that you install
`gcloud` and authenticate using [User Application Default Credentials ("ADCs")](https://cloud.google.com/sdk/gcloud/reference/auth/application-default)
as a primary authentication method. You can enable ADCs by running the command
`gcloud auth application-default login`.

Google Cloud reads the quota project for requests will be read automatically
from the `core/project` value. You can override this project by specifying the
`--project` flag when running `gcloud auth application-default login`. `gcloud`
should return this message if you have set the correct billing project:
`Quota project "your-project" was added to ADC which can be used by Google client libraries for billing and quota.`

### Running Terraform on Google Cloud

If you are running Terraform in a machine on Google Cloud, you can configure
that instance or cluster to use a [Google Service Account](https://cloud.google.com/compute/docs/authentication).
This allows Terraform to authenticate to Google Cloud without a separate
credential/authentication file. Ensure that the scope of the VM/Cluster is set
to or includes `https://www.googleapis.com/auth/cloud-platform`.

### Running Terraform Outside of Google Cloud

If you are running Terraform outside of Google Cloud, generate an external
credential configuration file ([example for OIDC based federation](https://cloud.google.com/iam/docs/access-resources-oidc#generate-automatic))
or a service account key file and set the `GOOGLE_APPLICATION_CREDENTIALS`
environment variable to the path of the JSON file. Terraform will use that file
for authentication. Terraform supports the full range of
authentication options [documented for Google Cloud](https://cloud.google.com/docs/authentication).

### Using Terraform Cloud

Place your credentials in a Terraform Cloud [environment variable](https://www.terraform.io/docs/cloud/workspaces/variables.html):
1. Create an environment variable called `GOOGLE_CREDENTIALS` in your Terraform Cloud workspace.
2. Remove the newline characters from your JSON key file and then paste the credentials into the environment variable value field. You can use the tr command to strip newline characters. `cat key.json  | tr -s '\n' ' '`
3. Mark the variable as **Sensitive** and click **Save variable**.

All runs within the workspace will use the `GOOGLE_CREDENTIALS` variable to authenticate with Google Cloud Platform.

### Impersonating Service Accounts

Terraform can [impersonate a Google service account](https://cloud.google.com/iam/docs/creating-short-lived-service-account-credentials),
acting as a service account without managing its key locally.

To impersonate a service account, you must use another authentication method
to act as a primary identity, and the primary identity must have the
`roles/iam.serviceAccountTokenCreator` role on the service account Terraform is
impersonating. Google Cloud Platform checks permissions and quotas against the
impersonated service account regardless of the primary identity in use.

## Authentication Configuration

* `credentials` - (Optional) Either the path to or the contents of a
[service account key file] in JSON format. You can
[manage key files using the Cloud Console]. Your service account key file is
used to complete a two-legged OAuth 2.0 flow to obtain access tokens to
authenticate with the GCP API as needed; Terraform will use it to reauthenticate
automatically when tokens expire. You can alternatively use the
`GOOGLE_CREDENTIALS` environment variable, or any of the following ordered
by precedence.

    * GOOGLE_CREDENTIALS
    * GOOGLE_CLOUD_KEYFILE_JSON
    * GCLOUD_KEYFILE_JSON

    Using Terraform-specific [service accounts] to authenticate with GCP is the
    recommended practice when using Terraform. If no Terraform-specific
    credentials are specified, the provider will fall back to using
    [Google Application Default Credentials][adc]. To use them, you can enter
    the path of your service account key file in the
    `GOOGLE_APPLICATION_CREDENTIALS` environment variable, or configure
    authentication through one of the following;

* If you're running Terraform from a GCE instance, default credentials
are automatically available. See
[Creating and Enabling Service Accounts for Instances][gce-service-account]
for more details.

* On your workstation, you can make your Google identity available by
running [`gcloud auth application-default login`][gcloud adc].

---

* `scopes` - (Optional) The list of OAuth 2.0 [scopes] requested when generating
an access token using the service account key specified in `credentials`.

By default, the following scopes are configured:

    * https://www.googleapis.com/auth/cloud-platform
    * https://www.googleapis.com/auth/userinfo.email

---

* `access_token` - (Optional) A temporary [OAuth 2.0 access token] obtained from
the Google Authorization server, i.e. the `Authorization: Bearer` token used to
authenticate HTTP requests to GCP APIs. This is an alternative to `credentials`,
and ignores the `scopes` field. You can alternatively use the
`GOOGLE_OAUTH_ACCESS_TOKEN` environment variable. If you specify both with
environment variables, Terraform uses the `access_token` instead of the
`credentials` field.

    -> Terraform cannot renew these access tokens, and they will eventually
expire (default `1 hour`). If Terraform needs access for longer than a token's
lifetime, use a service account key with `credentials` instead.

---

* `impersonate_service_account` - (Optional) The service account to impersonate for all Google API Calls.
You must have `roles/iam.serviceAccountTokenCreator` role on that account for the impersonation to succeed.
If you are using a delegation chain, you can specify that using the `impersonate_service_account_delegates` field.
Alternatively, this can be specified using the `GOOGLE_IMPERSONATE_SERVICE_ACCOUNT` environment
variable.

* `impersonate_service_account_delegates` - (Optional) The delegation chain for an impersonating a service account as described [here](https://cloud.google.com/iam/docs/creating-short-lived-service-account-credentials#sa-credentials-delegated).

## Quota Management Configuration

* `user_project_override` - (Optional) Defaults to `false`. Controls the quota
project used in requests to GCP APIs for the purpose of preconditions, quota,
and billing. If `false`, the quota project is determined by the API and may be
the project associated with your credentials, or the resource project. If `true`,
most resources in the provider will explicitly supply their resource project, as
described in their documentation. Otherwise, a `billing_project` value must be
supplied. Alternatively, this can be specified using the `USER_PROJECT_OVERRIDE`
environment variable.

Service account credentials are associated with the project the service account
was created in. Credentials that come from the gcloud tool are associated with a
project owned by Google. In order to properly use credentials that come from
gcloud with Terraform, it is recommended to set this property to true.

`user_project_override` uses the `X-Goog-User-Project`
[system parameter](https://cloud.google.com/apis/docs/system-parameters). When
set to true, the caller must have `serviceusage.services.use` permission on the
quota project.

---

* `billing_project` - (Optional) A quota project to send in `user_project_override`,
used for all requests sent from the provider. If set on a resource that supports
sending the resource project, this value will supersede the resource project.
This field is ignored if `user_project_override` is set to false or unset.
Alternatively, this can be specified using the `GOOGLE_BILLING_PROJECT`
environment variable.

## Provider Default Values Configuration

* `project` - (Optional) The default project to manage resources in. If another
project is specified on a resource, it will take precedence. This can also be
specified using the `GOOGLE_PROJECT` environment variable, or any of the
following ordered by precedence.

    * GOOGLE_PROJECT
    * GOOGLE_CLOUD_PROJECT
    * GCLOUD_PROJECT
    * CLOUDSDK_CORE_PROJECT

---

* `region` - (Optional) The default region to manage resources in. If another
region is specified on a regional resource, it will take precedence.
Alternatively, this can be specified using the `GOOGLE_REGION` environment
variable or any of the following ordered by precedence.

    * GOOGLE_REGION
    * GCLOUD_REGION
    * CLOUDSDK_COMPUTE_REGION

---

* `zone` - (Optional) The default zone to manage resources in. Generally, this
zone should be within the default region you specified. If another zone is
specified on a zonal resource, it will take precedence. Alternatively, this can
be specified using the `GOOGLE_ZONE` environment variable or any of the
following ordered by precedence.

    * GOOGLE_ZONE
    * GCLOUD_ZONE
    * CLOUDSDK_COMPUTE_ZONE

---

* `default_labels` (Optional) Labels that will be applied to all resources
with a top level `labels` field or a `labels` field nested inside a top level
`metadata` field. Setting the same key as a default label at the resource level
will override the default value for that label. These values will be recorded in 
individual resource plans through the `terraform_labels` and `effective_labels`
fields.

```
provider "google" {
  default_labels = {
    my_global_key = "one"
    my_default_key = "two"
  }
}

resource "google_compute_address" "my_address" {
  name     = "my-address"

  labels = {
    my_key = "three"
    # overrides provider-wide setting
    my_default_key = "four"
  }
}
```

## Advanced Settings Configuration

* `request_timeout` - (Optional) A duration string controlling the amount of time
the provider should wait for individual HTTP requests. This will not adjust the
amount of time the provider will wait for a logical operation - use the resource
timeout blocks for that. This will adjust only the amount of time that a single
synchronous request will wait for a response. The default is 120 seconds, and
that should be a suitable value in most cases. Many GCP APIs will cancel a
request if no response is forthcoming within 30 seconds in any event. In
limited cases, such as DNS record set creation, there is a synchronous request
to create the resource. This may help in those cases.


---

* `request_reason` - (Optional) Send a Request Reason [System Parameter](https://cloud.google.com/apis/docs/system-parameters)
for each API call made by the provider.  The `X-Goog-Request-Reason` header
value is used to provide a user-supplied justification into GCP AuditLogs.
Alternatively, this can be specified using the `CLOUDSDK_CORE_REQUEST_REASON`
environment variable.

---

* `{{service}}_custom_endpoint` - (Optional) The endpoint for a service's APIs,
such as `compute_custom_endpoint`. Defaults to the production GCP endpoint for
the service. This can be used to configure the Google provider to communicate
with GCP-like APIs such as [the Cloud Functions emulator](https://github.com/googlearchive/cloud-functions-emulator).
Values are expected to include the version of the service, such as
`https://www.googleapis.com/compute/v1/`:

```
provider "google" {
  alias                   = "compute_beta_endpoint"
  compute_custom_endpoint = "https://www.googleapis.com/compute/beta/"
}
```

Custom endpoints are an advanced feature. To determine the possible values you
can set, consult the implementation in [provider.go](https://github.com/hashicorp/terraform-provider-google-beta/blob/main/google-beta/provider.go)
and [config.go](https://github.com/hashicorp/terraform-provider-google-beta/blob/main/google-beta/config.go).

Support for custom endpoints is on a best-effort basis. The underlying
endpoint and default values for a resource can be changed at any time without
being considered a breaking change.

---

* `universe_domain` - (Optional) Specify the GCP universe to deploy in.

---

* `batching` - (Optional) Controls batching for specific GCP request types
where users have encountered quota or speed issues using many resources of
the same type, typically `google_project_service`.

Batching is not used for every resource/request type and can only group parallel
similar calls for nodes at a similar traversal time in the graph during
`terraform apply` (e.g. resources created using `count` that affect a single
`project`). Thus, it is also bounded by the `terraform`
[`-parallelism`](https://www.terraform.io/docs/commands/apply.html#parallelism-n)
flag, as reducing the number of parallel calls will reduce the number of
simultaneous requests being added to a batcher.

  ~> **NOTE** Most resources/GCP request do not have batching implemented (see
  below for requests which use batching) Batching is really only needed for
  resources where several requests are made at the same time to an underlying
  GCP resource protected by a fairly low default quota, or with very slow
  operations with slower eventual propagation. If you're not completely sure
  what you are doing, avoid setting custom batching configuration.

**So far, batching is implemented for below resources**:

* `google_project_service`
* All `google_*_iam_*` resources

The `batching` block supports the following fields.

* `send_after` - (Optional) A duration string representing the amount of time
after which a request should be sent. Defaults to 10s. Should be a non-negative
integer or float string with a unit suffix, such as "300ms", "1.5h" or "2h45m".
Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".

* `enable_batching` - (Optional) Defaults to true. If false, disables global
batching and each request is sent normally.

---

You can extend the user agent header for each request made by the provider by setting the `GOOGLE_TERRAFORM_USERAGENT_EXTENSION` environment variable. This can be helpful for tracking (e.g. compliance through [audit logs](https://cloud.google.com/logging/docs/audit)) or debugging purposes.

Example:

```sh
export GOOGLE_TERRAFORM_USERAGENT_EXTENSION="my-extension/1.0"
```

See [RFC 9110](https://www.rfc-editor.org/rfc/rfc9110#field.user-agent) for format compliance of user agent header fields. 

[OAuth 2.0 access token]: https://developers.google.com/identity/protocols/OAuth2
[service account key file]: https://cloud.google.com/iam/docs/creating-managing-service-account-keys
[manage key files using the Cloud Console]: https://console.cloud.google.com/apis/credentials/serviceaccountkey
[adc]: https://cloud.google.com/docs/authentication/production
[gce-service-account]: https://cloud.google.com/compute/docs/authentication
[gcloud adc]: https://cloud.google.com/sdk/gcloud/reference/auth/application-default/login
[service accounts]: https://cloud.google.com/docs/authentication/getting-started
[scopes]: https://developers.google.com/identity/protocols/googlescopes
