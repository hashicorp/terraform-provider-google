---
layout: "google"
page_title: "Google: google_container_engine_versions"
sidebar_current: "docs-google-datasource-container-versions"
description: |-
  Provides lists of available Google Container Engine versions for masters and nodes.
---

# google\_container\_engine\_versions

Provides access to available Google Container Engine versions in a zone or region for a given project.

```hcl
data "google_container_engine_versions" "central1b" {
  zone = "us-central1-b"
}

resource "google_container_cluster" "foo" {
  name               = "terraform-test-cluster"
  zone               = "us-central1-b"
  node_version       = "${data.google_container_engine_versions.central1b.latest_node_version}"
  initial_node_count = 1

  master_auth {
    username = "mr.yoda"
    password = "adoy.rm"
  }
}
```

## Argument Reference

The following arguments are supported:

* `zone` (optional) - Zone to list available cluster versions for. Should match the zone the cluster will be deployed in.
    One of zone or region must be given or inferred from provider

* `region` (optional) - Region to list available cluster versions for. Should match the region the cluster will be deployed in.
    One of zone or region must be given or inferred from provider

* `project` (optional) - ID of the project to list available cluster versions for. Should match the project the cluster will be deployed to.
  Defaults to the project that the provider is authenticated with.

Terraform will attempt to get or infer the location (zone or region) in the following order:
1. `zone` from data source config
2. `region` from data source config
3. Provider-level zone
4. Provider-level region if no zone is given at provider-level. If provider-level zone is given but you wish to use a
   regional location, please specify it in the data source.


## Attributes Reference

The following attributes are exported:

* `valid_master_versions` - A list of versions available in the given zone for use with master instances.
* `valid_node_versions` - A list of versions available in the given zone for use with node instances.
* `latest_master_version` - The latest version available in the given zone for use with master instances.
* `latest_node_version` - The latest version available in the given zone for use with node instances.
* `default_cluster_version` - Version of Kubernetes the service deploys by default.
