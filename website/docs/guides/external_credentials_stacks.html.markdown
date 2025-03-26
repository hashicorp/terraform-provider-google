---
page_title: "Use external credentials in the Google Cloud provider with Terraform Stacks"
description: |-
  How to use external credentials in the Google Cloud provider with Terraform Stacks
---

# External Credentials in the Google Cloud provider with Terraform Stacks

Apart from using `access_token` and `credential` fields in the provider configuration, you can also use external credentials in the Google Cloud provider that are provided through a Workload Identity Federation (WIF) provider. This can be used to authenticate Terraform Stacks to provision resources in Google Cloud.

## Setting up a Workload Identity Federation (WIF) credentials

## Stacks Setup

A Terraform Stacks Project requires the following:

- A Workload Identity Federation (WIF) provider to authenticate Terraform Stacks
- Components - `components.tfstacks.hcl`
- Deployment - `deployments.tfdeploy.hcl`

## Generating the Workload Identity Federation (WIF) credentials

In the case of Stacks, we need to create the both the workload identity pool and the pool provider in order for Terraform Stacks to authenticate.

Required variables:

- `project_id` - The GCP project ID
- `hcp_tf_organization` - The HCP Terraform organization
- `hcp_tf_stacks_project` - The HCP Terraform Stacks project

```hcl
provider "google" {
  region = "global"
}

variable "project_id" {
  type        = string
  description = "GCP Project ID"
}

variable "hcp_tf_organization" {
  type        = string
  description = "HCP Terraform Organization"
}

variable "hcp_tf_stacks_project" {
  type        = string
  description = "HCP Terraform Stacks Project"
}

# Create a service account for Terraform Stacks
resource "google_service_account" "terraform_stacks_sa" {
  account_id   = "terraform-stacks-sa"
  display_name = "Terraform Stacks Service Account"
  description  = "Service account used by Terraform Stacks for GCP resources"
}

# Enable required APIs for workload identity federation
locals {
  gcp_service_list = [
    "sts.googleapis.com",
    "iam.googleapis.com",
    "iamcredentials.googleapis.com"
  ]
}

resource "google_project_service" "services" {
  for_each                   = toset(local.gcp_service_list)
  project                    = var.project_id
  service                    = each.key
  disable_dependent_services = false
  disable_on_destroy         = false
}

# Create Workload Identity Pool (reference google_project_service to ensure APIs are enabled)
resource "google_iam_workload_identity_pool" "terraform_stacks_pool" {
  depends_on = [google_project_service.services]
  workload_identity_pool_id = "terraform-stacks-pool-3"
  display_name              = "Terraform Stacks Pool-3"
  description               = "Identity pool for Terraform Stacks authentication"
}

# Create Workload Identity Pool Provider
resource "google_iam_workload_identity_pool_provider" "terraform_stacks_provider" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.terraform_stacks_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = "terraform-stacks-provider"
  display_name                       = "Terraform Stacks Provider"
  description                        = "OIDC identity pool provider for Terraform Stacks"
  
  attribute_mapping = {
    "google.subject"                            = "assertion.sub", # WARNING - this value is has to be <=127 bytes, and is "organization:<ORG NAME>:project:<PROJ NAME>:stack:<STACK NAME>:deployment:development:operation:plan
    "attribute.aud"                             = "assertion.aud",
    "attribute.terraform_operation"             = "assertion.terraform_operation",
    "attribute.terraform_project_id"            = "assertion.terraform_project_id",
    "attribute.terraform_project_name"          = "assertion.terraform_project_name",
    "attribute.terraform_stack_id"              = "assertion.terraform_stack_id",
    "attribute.terraform_stack_name"            = "assertion.terraform_stack_name",
    "attribute.terraform_stack_deployment_name" = "assertion.terraform_stack_deployment_name",
    "attribute.terraform_organization_id"       = "assertion.terraform_organization_id",
    "attribute.terraform_organization_name"     = "assertion.terraform_organization_name",
    "attribute.terraform_run_id"                = "assertion.terraform_run_id",
  }
  oidc {
    issuer_uri = "https://app.terraform.io"
    allowed_audiences = ["hcp.workload.identity"]
  }
  attribute_condition = "assertion.sub.startsWith(\"organization:${var.hcp_tf_organization}:project:${var.hcp_tf_stacks_project}:stack\")"
}

# Switch from binding to member for service account IAM
resource "google_service_account_iam_member" "workload_identity_user" {
  service_account_id = google_service_account.terraform_stacks_sa.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.terraform_stacks_pool.name}/*"
}

# Grant additional permissions to the service account
resource "google_project_iam_member" "sa_more_permissions" {
  project = var.project_id
  role    = "roles/iam.serviceAccountTokenCreator"
  member  = "serviceAccount:${google_service_account.terraform_stacks_sa.email}"
}

# Grant editor role to the service account (similar to reference implementation)
resource "google_project_iam_member" "sa_editor" {
  project = var.project_id
  role    = "roles/editor"
  member  = "serviceAccount:${google_service_account.terraform_stacks_sa.email}"
}

# Outputs to be used by Terraform Stacks
output "service_account_email" {
  value       = google_service_account.terraform_stacks_sa.email
  description = "Email of the service account to be used by Terraform Stacks"
}

output "audience" {
  value       = "//iam.googleapis.com/${google_iam_workload_identity_pool_provider.terraform_stacks_provider.name}"
  description = "The audience value to use when generating OIDC tokens"
}
```

The both output values `service_account_email` and `audience` will be used to authenticate Terraform Stacks.

## Terraform Stacks Setup with External Credentials

Initially you'll want to link your [Terraform Project to the desired Stack through VCS](https://developer.hashicorp.com/terraform/cloud-docs/stacks/create#requirements).

Afterwards, you'll want to setup Terraform Stacks with the use of external credentials. A simple configuration is shown below. Check out the [Terraform Stacks Overview](https://developer.hashicorp.com/terraform/language/stacks) for more information.

Upon setup, you'll now be able to provision GCP resources through Terraform Stacks.

`deployments.tfdeploy.hcl`
```hcl
identity_token "jwt" {
  audience = ["hcp.workload.identity"]
}

deployment "staging" {
  inputs = {
    jwt = identity_token.jwt.jwt
  }
}
```

`components.tfstacks.hcl`
```hcl
required_providers {
  google = {
    source = "hashicorp/google"
    version = "6.25.0"
  }
}

provider "google" "this" {
  config {
    external_credentials {
      audience = output.audience // audience from WIF
      service_account_email = output.service_account_email // service account created from WIF
      identity_token = var.jwt
    }
  }
}

variable "jwt" {
  type = string
}

component "storage_buckets" {
    source = "./buckets"

    inputs = {
        jwt = var.jwt
    }

    providers = {
        google    = provider.google.this
    }
}
```