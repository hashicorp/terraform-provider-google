---
subcategory: "Cloud Platform"
layout: "google"
page_title: "Google: google_client_config"
sidebar_current: "docs-google-datasource-client-config"
description: |-
  Get information about the configuration of the Google Cloud provider.
---

# google\_client\_config

Use this data source to access the configuration of the Google Cloud provider.

## Example Usage

```tf
data "google_client_config" "current" {
}

output "project" {
  value = data.google_client_config.current.project
}
```

## Example Usage: Configure Kubernetes provider with OAuth2 access token

```tf
data "google_client_config" "default" {
}

data "google_container_cluster" "my_cluster" {
  name = "my-cluster"
  zone = "us-east1-a"
}

provider "kubernetes" {
  load_config_file = false

  host  = "https://${data.google_container_cluster.my_cluster.endpoint}"
  token = data.google_client_config.default.access_token
  cluster_ca_certificate = base64decode(
    data.google_container_cluster.my_cluster.master_auth[0].cluster_ca_certificate,
  )
}
```

## Argument Reference

There are no arguments available for this data source.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `project` - The ID of the project to apply any resources to.

* `region` - The region to operate under.

* `zone` - The zone to operate under.

* `access_token` - The OAuth2 access token used by the client to authenticate against the Google Cloud API.
