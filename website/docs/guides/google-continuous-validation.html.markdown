---
page_title: "Using Terraform Cloud's Continuous Validation feature with the Google Provider"
description: |-
  Continuous validation helps identify issues immediately and continuously instead of waiting until customers encounter problems. This guide shows how continuous validation can be used with the Google provider.
---

# Using Terraform Cloud's Continuous Validation feature with the Google Provider

The Continuous Validation feature in Terraform Cloud (TFC) allows users to make assertions about their infrastructure between applied runs. This helps users to identify issues at the time they first appear and avoid situations where a change is only identified during a future terraform plan/apply or once it causes a user-facing problem.

Users can add checks to their Terraform configuration using an HCL language feature called [check{} blocks](https://developer.hashicorp.com/terraform/language/checks). Check blocks contain assertions that are defined with a custom condition expression and an error message. When the condition expression evaluates to true the check passes, but when the expression evaluates to false Terraform will show a warning message that includes the user-defined error message.

Custom conditions can be created using data from Terraform providers’ resources and data sources. Data can also be combined from multiple sources; for example, you can use checks to monitor expirable resources by comparing a resource’s expiration date attribute to the current time returned by Terraform’s built-in time functions. These include the [plantimestamp function](https://developer.hashicorp.com/terraform/language/functions/plantimestamp), which was added in Terraform 1.5.

For more information about continuous validation visit the [Workspace Health page in the Terraform Cloud documentation](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/health#continuous-validation).

Below, this guide shows examples of how data returned by the Google provider can be used to define checks in your Terraform configuration. In each example it is assumed that the Google provider is configured with a default project, region, and zone.

~> Check blocks and the plantime function are available in Terraform 1.5 and later

## Example - Assert a VM is in a running state (`google_compute_instance`)

VM instances provisioned using Compute Engine can pass through several states as part of the [VM instance lifecycle](https://cloud.google.com/compute/docs/instances/instance-life-cycle). Once a VM is provisioned it could experience an error, or a user could suspend or stop that VM in the Google Cloud console, without that change being detected until the next Terraform plan is generated. Continuous validation can be used to assert the state of a VM and detect if there are any unexpected status changes that occur out-of-band.

The example below shows how a check block can be used to assert that a VM is in the running state.

You can force the check to fail in this example by provisioning the VM, manually stopping it in the Google Cloud console, and then triggering a health check in TFC. The check will fail and report that the VM is not running.

```hcl
data "google_compute_network" "default" {
  name = "default"
}
resource "google_compute_instance" "vm_instance" {
  name         = "my-instance"
  machine_type = "f1-micro"
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }
  network_interface {
    network = data.google_compute_network.default.name
    access_config {
    }
  }
}
check "check_vm_status" {
  data "google_compute_instance" "vm_instance" {
    name = google_compute_instance.vm_instance.name
  }
  assert {
    condition = data.google_compute_instance.vm_instance.current_status == "RUNNING"
    error_message = format("Provisioned VMs should be in a RUNNING status, instead the VM `%s` has status: %s",
      data.google_compute_instance.vm_instance.name,
      data.google_compute_instance.vm_instance.current_status
    )
  }
}
```

## Example - Check if a certificate will expire within a certain timeframe (`google_privateca_certificate`)

Certificates can be provisioned using either the Certificate Manager, Certificate Authority Service (‘Private CA’), and Compute Engine APIs. In this example we provision a certificate via the Certificate Authority Service that has a user-supplied [lifetime argument](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/privateca_certificate#lifetime). After the lifetime duration passes the certificate is automatically deleted in GCP. By creating a check that asserts the certificate’s expiration date is more than 30 days away we can be notified by TFC health checks when the certificate is approaching expiration and needs manual intervention.

In the example below we provision a certificate with a lifetime of 30 days and 2 minutes (see `local.month_and_2min_in_second_duration`) and create a check that asserts certificates should be valid for the next month (see `local.month_in_hour_duration`).

We can see the check begin to fail by waiting 2 minutes after the certificate is provisioned and then triggering a health check in TFC. The check will fail and report that the certificate is due to expire in less than a month.

```hcl
locals {
  month_in_hour_duration            = "${24 * 30}h"
  month_and_2min_in_second_duration = "${(60 * 60 * 24 * 30) + (60 * 2)}s"
}
resource "tls_private_key" "example" {
  algorithm = "RSA"
}
resource "tls_cert_request" "example" {
  private_key_pem = tls_private_key.example.private_key_pem
  subject {
    common_name  = "example.com"
    organization = "ACME Examples, Inc"
  }
}
resource "google_privateca_ca_pool" "default" {
  name     = "my-ca-pool"
  location = "us-central1"
  tier     = "ENTERPRISE"
  publishing_options {
    publish_ca_cert = true
    publish_crl     = true
  }
  labels = {
    terraform = true
  }
  issuance_policy {
    baseline_values {
      ca_options {
        is_ca = false
      }
      key_usage {
        base_key_usage {
          digital_signature = true
          key_encipherment  = true
        }
        extended_key_usage {
          server_auth = true
        }
      }
    }
  }
}
resource "google_privateca_certificate_authority" "test-ca" {
  Deletion_protection      = false
  certificate_authority_id = "my-authority"
  location                 = google_privateca_ca_pool.default.location
  pool                     = google_privateca_ca_pool.default.name
  config {
    subject_config {
      subject {
        country_code        = "us"
        organization        = "google"
        organizational_unit = "enterprise"
        locality            = "mountain view"
        province            = "california"
        street_address      = "1600 amphitheatre parkway"
        postal_code         = "94109"
        common_name         = "my-certificate-authority"
      }
    }
    x509_config {
      ca_options {
        is_ca = true
      }
      key_usage {
        base_key_usage {
          cert_sign = true
          crl_sign  = true
        }
        extended_key_usage {
          server_auth = true
        }
      }
    }
  }
  type = "SELF_SIGNED"
  key_spec {
    algorithm = "RSA_PKCS1_4096_SHA256"
  }
}
resource "google_privateca_certificate" "default" {
  name                  = "my-certificate"
  pool                  = google_privateca_ca_pool.default.name
  certificate_authority = google_privateca_certificate_authority.test-ca.certificate_authority_id
  location              = google_privateca_ca_pool.default.location
  lifetime              = local.month_and_2min_in_second_duration # lifetime is 2mins over the threshold in the check block below
  pem_csr               = tls_cert_request.example.cert_request_pem
}
check "check_certificate_state" {
  assert {
    condition = timecmp(plantimestamp(), timeadd(
      google_privateca_certificate.default.certificate_description[0].subject_description[0].not_after_time,
    "-${local.month_in_hour_duration}")) < 0
    error_message = format("Provisioned certificates should be valid for at least 30 days, but `%s`is due to expire on `%s`.",
      google_privateca_certificate.default.name,

      google_privateca_certificate.default.certificate_description[0].subject_description[0].not_after_time
    )
  }
}
```

## Example - Validate the status of a Cloud Function (`google_cloudfunctions2_function`)

Cloud Functions can have multiple statuses depending on issues that occur during deployment or triggering the function. These are: ACTIVE, FAILED, DEPLOYING, DELETING
 
In the example below we create a 2nd generation cloud function that uses source code stored as a .zip file in a GCS bucket. A .zip file containing the files needed by the function is uploaded by Terraform from the local machine. In the check we use the `google_cloudfunctions2_function` data source’s [state attribute](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/cloudfunctions2_function#state) to access the function’s state and assert that the function is active.

```hcl
resource "google_storage_bucket" "bucket" {
  name                        = "my-bucket"
  location                    = "US"
  uniform_bucket_level_access = true
}
resource "google_storage_bucket_object" "object" {
  name   = "function-source.zip"
  bucket = google_storage_bucket.bucket.name
  source = "./function-source.zip"
}
resource "google_cloudfunctions2_function" "my-function" {
  name        = "my-function"
  location    = "us-central1"
  description = "a new function"
  build_config {
    runtime     = "nodejs12"
    entry_point = "helloHttp"
    source {
      storage_source {
        bucket = google_storage_bucket.bucket.name
        object = google_storage_bucket_object.object.name
      }
    }
  }
  service_config {
    max_instance_count = 1
    available_memory   = "1536Mi"
    timeout_seconds    = 30
  }
}
check "check_cf_state" {
  data "google_cloudfunctions2_function" "my-function" {
    name     = google_cloudfunctions2_function.my-function.name
    location = google_cloudfunctions2_function.my-function.location
  }
  assert {
    condition = data.google_cloudfunctions2_function.my-function.state == "ACTIVE"
    error_message = format("Provisioned Cloud Functions should be in an ACTIVE state, instead the function `%s` has state: %s",
      data.google_cloudfunctions2_function.my-function.name,
      data.google_cloudfunctions2_function.my-function.state
    )
  }
}
```

