---
layout: "google"
page_title: "Using GKE with Terraform"
sidebar_current: "docs-google-provider-guides-using-gke"
description: |-
  Recommendations and best practices for using GKE with Terraform.
---

# Using GKE with Terraform

This page is a brief overview of GKE usage with Terraform, based on the content
available in the [How-to guides for GKE](https://cloud.google.com/kubernetes-engine/docs/how-to).
It's intended as a supplement for intermediate users, covering cases that are
unintuitive or confusing when using Terraform instead of `gcloud`/the Cloud
Console.

Additionally, you may consider using Google's [`kubernetes-engine`](https://registry.terraform.io/modules/terraform-google-modules/kubernetes-engine/google)
module, which implements many of this practices for you.

If the information on this page conflicts with recommendations available on
`cloud.google.com`, `cloud.google.com` should be considered the correct source.

## Interacting with Kubernetes

After creating a `google_container_cluster` with Terraform, authentication to
the cluster are often a challenge. In most cases, you can use `gcloud` to
configure cluster access, [generating a `kubeconfig` entry](https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-access-for-kubectl#generate_kubeconfig_entry):

```bash
gcloud container clusters get-credentials cluster-name
```

Using this command, `gcloud` will generate a `kubeconfig` entry that uses
`gcloud` as an authentication mechanism. However, sometimes performing
authentication inline with Terraform or a static config without `gcloud` is more
desirable.

### Using the Kubernetes and Helm Providers

When using the `kubernetes` and `helm` providers,
[statically defined credentials](https://www.terraform.io/docs/providers/kubernetes/index.html#statically-defined-credentials)
can allow you to connect to clusters defined in the same config or in a remote
state. You can configure either using configuration such as the following:

```hcl
# Retrieve an access token as the Terraform runner
data "google_client_config" "provider" {}

data "google_container_cluster" "my_cluster" {
  name     = "my-cluster"
  location = "us-central1"
}

provider "kubernetes" {
  load_config_file = false

  host  = "https://${data.google_container_cluster.my_cluster.endpoint}"
  token = data.google_client_config.provider.access_token
  cluster_ca_certificate = base64decode(
    data.google_container_cluster.my_cluster.master_auth[0].cluster_ca_certificate,
  )
}
```

Alternatively, you can authenticate as another service account on which your
Terraform runner has been granted the `roles/iam.serviceAccountTokenCreator`
role:

```hcl
data "google_service_account_access_token" "my_kubernetes_sa" {
  target_service_account = "{{service_account}}"
  scopes                 = ["userinfo-email", "cloud-platform"]
  lifetime               = "3600s"
}

data "google_container_cluster" "my_cluster" {
  name     = "my-cluster"
  location = "us-central1"
}

provider "kubernetes" {
  load_config_file = false

  host  = "https://${data.google_container_cluster.my_cluster.endpoint}"
  token = data.google_service_account_access_token.my_kubernetes_sa.access_token
  cluster_ca_certificate = base64decode(
    data.google_container_cluster.my_cluster.master_auth[0].cluster_ca_certificate,
  )
}
```

### Using kubectl / kubeconfig

It's possible to interface with `kubectl` or other `.kubeconfig`-based tools by
providing them a `.kubeconfig` directly. For situations where `gcloud` can't be
used as an authentication mechanism, you can generate a static `.kubeconfig`
file instead.

An authentication submodule, `auth`, is provided as part of Google's
[`kubernetes-engine`](https://registry.terraform.io/modules/terraform-google-modules/kubernetes-engine/google)
module. You can use it through the module registry, or [in the module source](https://github.com/terraform-google-modules/terraform-google-kubernetes-engine/tree/master/modules/auth).

Authenticating using this method will use a Terraform-generated access token
which persists for 1 hour. For longer-lasting sessions, or cases where a single
persistent config is required, using `gcloud` is advised.

## VPC-native Clusters

[VPC-native clusters](https://cloud.google.com/kubernetes-engine/docs/how-to/alias-ips)
are GKE clusters that use [alias IP ranges](https://cloud.google.com/vpc/docs/alias-ip).
VPC-native clusters route traffic between pods using a VPC network, and are able
to route to other VPCs across network peerings along with [several other benefits](https://cloud.google.com/kubernetes-engine/docs/how-to/alias-ips).

This is in contrast to [routes-based clusters](https://cloud.google.com/kubernetes-engine/docs/how-to/routes-based-cluster),
which route pod traffic using GCP routes.

In both `gcloud` and the Cloud Console, VPC-native is the default for new
clusters and increasingly, GKE features such as [Standalone Network Endpoint Groups (NEGs)](https://cloud.google.com/kubernetes-engine/docs/how-to/standalone-neg#pod_readiness)
have relied on clusters being VPC-native. In Terraform however, the default
behaviour is to create a routes-based cluster for backwards compatibility.

It's recommended that you create a VPC-native cluster, done by specifying the
`ip_allocation_policy` block. Configuration will look like the following:

```hcl
resource "google_container_cluster" "my_vpc_native_cluster" {
  name               = "my-vpc-native-cluster"
  location           = "us-central1"
  initial_node_count = 1

  network    = "default"
  subnetwork = "default"

  ip_allocation_policy {
    cluster_ipv4_cidr_block  = "/16"
    services_ipv4_cidr_block = "/22"
  }

  # other settings...
}
```

## Node Pool Management

In Terraform, we recommend managing your node pools using the 
`google_container_node_pool` resource, separate from the
`google_container_cluster` resource. This separates cluster-level configuration
like networking and Kubernetes features from the configuration of your nodes.
Additionally, it helps ensure your cluster isn't inadvertently deleted.
Terraform struggles to handle complex changes to subresources, and may attempt
to delete a cluster based on changes to inline node pools.

However, the GKE API doesn't allow creating a cluster without nodes. It's common
for Terraform users to define a block such as the following:

```hcl
resource "google_container_cluster" "my-gke-cluster" {
  name     = "my-gke-cluster"
  location = "us-central1"

  # We can't create a cluster with no node pool defined, but we want to only use
  # separately managed node pools. So we create the smallest possible default
  # node pool and immediately delete it.
  remove_default_node_pool = true
  initial_node_count       = 1

  # other settings...
}
```

This creates `initial_node_count` nodes per zone the cluster has nodes in,
typically 1 zone if the cluster `location` is a zone, and 3 if it's a `region`.
Your cluster's initial GKE masters will be sized based on the
`initial_node_count` provided. If subsequent node pools add a large number of
nodes to your cluster, GKE may cause a resizing event immediately after adding a
node pool.

The initial node pool will be created using the
[Compute Engine default service account](https://cloud.google.com/compute/docs/access/service-accounts#default_service_account)
as the [`service_account`](https://cloud.google.com/compute/docs/access/service-accounts#default_service_account).
If you've disabled that service account, or want to use a
[least privilege Google service account](https://cloud.google.com/kubernetes-engine/docs/how-to/hardening-your-cluster#use_least_privilege_sa)
for the temporary  node pool, you can add the following configuration to your
`google_container_cluster` block:

```hcl
resource "google_container_cluster" "my-gke-cluster" {
  # other settings...

  node_config {
    service_account = "{{service_account}}"
  }
  
  lifecycle {
    ignore_changes = ["node_config"]
  }

  # other settings...
}
```
