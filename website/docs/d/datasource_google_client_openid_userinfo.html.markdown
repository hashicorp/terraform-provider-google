---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_client_openid_userinfo"
sidebar_current: "docs-google-datasource-client-openid-userinfo"
description: |-
  Get OpenID userinfo about the credentials used with the Google provider, specifically the email.
---

# google\_client\_openid\_userinfo

Get OpenID userinfo about the credentials used with the Google provider,
specifically the email.

This datasource enables you to export the email of the account you've
authenticated the provider with; this can be used alongside
`data.google_client_config`'s `access_token` to perform OpenID Connect
authentication with GKE and configure an RBAC role for the email used.

~> This resource will only work as expected if the provider is configured to
use the `https://www.googleapis.com/auth/userinfo.email` scope! You will
receive an error otherwise.

## Example Usage - exporting an email

```hcl
data "google_client_openid_userinfo" "me" {
}

output "my-email" {
  value = data.google_client_openid_userinfo.me.email
}
```

## Example Usage - OpenID Connect w/ Kubernetes provider + RBAC IAM role

```hcl
data "google_client_openid_userinfo" "provider_identity" {
}

data "google_client_config" "provider" {
}

data "google_container_cluster" "my_cluster" {
  name = "my-cluster"
  zone = "us-east1-a"
}

provider "kubernetes" {
  load_config_file = false

  host  = "https://${data.google_container_cluster.my_cluster.endpoint}"
  token = data.google_client_config.provider.access_token
  cluster_ca_certificate = base64decode(
    data.google_container_cluster.my_cluster.master_auth[0].cluster_ca_certificate,
  )
}

resource "kubernetes_cluster_role_binding" "user" {
  metadata {
    name = "provider-user-admin"
  }

  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = "cluster-admin"
  }

  subject {
    kind = "User"
    name = data.google_client_openid_userinfo.provider_identity.email
  }
}
```

## Argument Reference

There are no arguments available for this data source.

## Attributes Reference

The following attributes are exported:

* `email` - The email of the account used by the provider to authenticate with GCP.
