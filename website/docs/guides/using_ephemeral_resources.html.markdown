---
page_title: "Use ephemeral resources in the Google Cloud provider"
description: |-
  How to use ephemeral resources in the Google Cloud provider
---

# Ephemeral Resources in the Google Cloud provider

Ephemeral resources are Terraform resources that are essentially temporary. They allow users to access and use data in their configurations without that data being stored in Terraform state.

Ephemeral resources are available in Terraform v1.10 and later. For more information, see the [official HashiCorp documentation for Ephemeral Resources](https://developer.hashicorp.com/terraform/language/resources/ephemeral).

To mark the launch of the ephemeral resources feature, the Google Cloud provider has added four ephemeral resources:
- [`google_service_account_access_token`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/ephemeral-resources/service_account_access_token)
- [`google_service_account_id_token`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/ephemeral-resources/service_account_id_token)
- [`google_service_account_jwt`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/ephemeral-resources/service_account_jwt)
- [`google_service_account_key`](https://registry.terraform.io/providers/hashicorp/google/latest/docs/ephemeral-resources/service_account_key)

These are based on existing data sources already in the provider. In future you may wish to update your configurations to use these ephemeral versions, as they will allow you to avoid storing tokens and credentials values in your Terraform state.

## Use the Google Cloud provider's new ephemeral resources

Ephemeral resources are a source of ephemeral data, and they can be referenced in your configuration just like the attributes of resources and data sources. However, a field that references an ephemeral resource must be capable of handling ephemeral data. Due to this, resources in the Google Cloud provider will need to be updated so they include write-only attributes that are capable of using ephemeral data while not storing those values in the resource's state. 

Until then, ephemeral resources can only be used to pass values into the provider block, which is already capable of receiving ephemeral values.

The following sections show two examples from the new ephemeral resources' documentation pages, which can be used to test out the ephemeral resources in their current form.

### See how ephemeral resources behave during `terraform plan` and `terraform apply`

The [documentation](https://registry.terraform.io/providers/hashicorp/google/latest/docs/ephemeral-resources/service_account_key) for the `google_service_account_key` ephemeral resource has a simple example that you can use to view how ephemeral resources behave during plan and apply operations:

```hcl
resource "google_service_account" "myaccount" {
  account_id = "dev-foo-account"
}

resource "google_service_account_key" "mykey" {
  service_account_id = google_service_account.myaccount.name
}

ephemeral "google_service_account_key" "mykey" {
  name            = google_service_account_key.mykey.name
  public_key_type = "TYPE_X509_PEM_FILE"
}
```

During `terraform plan` you will see that the ephemeral resource is deferred, as it depends on other resources for its arguments:

```
ephemeral.google_service_account_key.mykey: Configuration unknown, deferring...

Terraform used the selected providers to generate the
following execution plan. Resource actions are indicated
with the following symbols:
  + create

Terraform will perform the following actions:

  # google_service_account.myaccount will be created
  + resource "google_service_account" "myaccount" {
    ...
```

During `terrform apply` you will see the ephemeral resource is the final resource to be evaluated, because it depends on the two other resources, and the ephemeral resource is not reflected in the statistics about how many resources were created during the apply action:

```
ephemeral.google_service_account_key.mykey: Opening...
ephemeral.google_service_account_key.mykey: Opening complete after 1s
ephemeral.google_service_account_key.mykey: Closing...
ephemeral.google_service_account_key.mykey: Closing complete after 0s

Apply complete! Resources: 2 added, 0 changed, 0 destroyed.
```

If you run the example using the local backend you can also inspect the state, where you will see that the ephemeral resource is not represented.


### Use ephemeral resources to configure the Google Cloud provider

The [documentation](https://registry.terraform.io/providers/hashicorp/google/latest/docs/ephemeral-resources/service_account_access_token) for the `google_service_account_access_token` ephemeral resource demonstrates how it can be used to configure the provider. Check that ephemeral resource's documentation for details about the IAM permissions required for this example to work:

```hcl
provider "google" {
}


ephemeral "google_service_account_access_token" "default" {
  provider               = google
  target_service_account = "service_B@projectB.iam.gserviceaccount.com"
  scopes                 = ["userinfo-email", "cloud-platform"]
  lifetime               = "300s"
}

provider "google" {
  alias        = "impersonated"
  access_token = ephemeral.google_service_account_access_token.default.access_token
}

data "google_client_openid_userinfo" "me" {
  provider = google.impersonated
}

output "target-email" {
  value = data.google_client_openid_userinfo.me.email
}
```

