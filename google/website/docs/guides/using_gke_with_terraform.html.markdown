---
page_title: "Using GKE with Terraform"
description: |-
  Recommendations and best practices for using GKE with Terraform.
---

# Using GKE with Terraform

-> Visit the [Provision a GKE Cluster (Google Cloud)](https://learn.hashicorp.com/tutorials/terraform/gke?in=terraform/kubernetes&utm_source=WEBSITE&utm_medium=WEB_IO&utm_offer=ARTICLE_PAGE&utm_content=DOCS) tutorial to learn how to provision and interact
with a GKE cluster.

This page is a brief overview of GKE usage with Terraform, based on the content
available in the [How-to guides for GKE](https://cloud.google.com/kubernetes-engine/docs/how-to).
It's intended as a supplement for intermediate users, covering cases that are
unintuitive or confusing when using Terraform instead of `gcloud`/the Cloud
Console.

Additionally, you may consider using Google's [`kubernetes-engine`](https://registry.terraform.io/modules/terraform-google-modules/kubernetes-engine/google)
module, which implements many of these practices for you.

If the information on this page conflicts with recommendations available on
`cloud.google.com`, `cloud.google.com` should be considered the correct source.

## Interacting with Kubernetes

After creating a `google_container_cluster` with Terraform, you can use `gcloud` to
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
[statically defined credentials](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs#credentials-config)
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
  host  = "https://${data.google_container_cluster.my_cluster.endpoint}"
  token = data.google_client_config.provider.access_token
  cluster_ca_certificate = base64decode(
    data.google_container_cluster.my_cluster.master_auth[0].cluster_ca_certificate,
  )
}
```
Although the above can result in authentication errors, over time, as the token recorded in the google_client_cofig data resource is short lived (thus it expires) and it's stored in state.  Fortunately, the [kubernetes provider can accept valid credentials from an exec-based plugin](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs#exec-plugins) to fetch a new token before each Terraform operation (so long as you have the [gke-cloud-auth-plugin for kubectl installed](https://cloud.google.com/blog/products/containers-kubernetes/kubectl-auth-changes-in-gke)), like so:
 
```hcl
# Retrieve an access token as the Terraform runner
data "google_client_config" "provider" {}

data "google_container_cluster" "my_cluster" {
  name     = "my-cluster"
  location = "us-central1"
}

provider "kubernetes" {
  host  = "https://${data.google_container_cluster.my_cluster.endpoint}"
  token = data.google_client_config.provider.access_token
  cluster_ca_certificate = base64decode(
    data.google_container_cluster.my_cluster.master_auth[0].cluster_ca_certificate,
  )
  exec {
    api_version = "client.authentication.k8s.io/v1beta1"
    command     = "gke-gcloud-auth-plugin"
  }
}
```

Alternatively, you can authenticate as another service account on which your
Terraform user has been granted the `roles/iam.serviceAccountTokenCreator`
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


In both `gcloud` and the Cloud Console, VPC-native is the default for new
clusters and many managed products such as CloudSQL, Memorystore and others
require VPC Native Clusters to work properly. In Terraform however, the default
behaviour is to create a routes-based cluster for backwards compatibility.

It's recommended that you create a VPC-native cluster, done by specifying the
`ip_allocation_policy` block or using secondary ranges on existing subnet. Configuration will look like the following:

```hcl
resource "google_compute_subnetwork" "custom" {
  name          = "test-subnetwork"
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.custom.id
  secondary_ip_range {
    range_name    = "services-range"
    ip_cidr_range = "192.168.1.0/24"
  }

  secondary_ip_range {
    range_name    = "pod-ranges"
    ip_cidr_range = "192.168.64.0/22"
  }
}

resource "google_compute_network" "custom" {
  name                    = "test-network"
  auto_create_subnetworks = false
}

resource "google_container_cluster" "my_vpc_native_cluster" {
  name               = "my-vpc-native-cluster"
  location           = "us-central1"
  initial_node_count = 1

  network    = google_compute_network.custom.id
  subnetwork = google_compute_subnetwork.custom.id

  ip_allocation_policy {
    cluster_secondary_range_name  = "pod-ranges"
    services_secondary_range_name = google_compute_subnetwork.custom.secondary_ip_range.0.range_name
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

### Windows Node Pools

You can add
[Windows Server node pools](https://cloud.google.com/kubernetes-engine/docs/how-to/creating-a-cluster-windows)
to your GKE cluster by adding `google_container_node_pool` to your Terraform
configuration with `image_type=WINDOWS_LTSC` or `WINDOWS_SAC`.

```hcl
resource "google_container_cluster" "demo_cluster" {
  project  = "" # Replace with your Project ID, https://cloud.google.com/resource-manager/docs/creating-managing-projects#identifying_projects
  name     = "demo-cluster"
  location = "us-west1-a"

  min_master_version = "1.27"

  # Enable Alias IPs to allow Windows Server networking.
  ip_allocation_policy {
    cluster_ipv4_cidr_block  = "/14"
    services_ipv4_cidr_block = "/20"
  }

  # Removes the implicit default node pool, recommended when using
  # google_container_node_pool.
  remove_default_node_pool = true
  initial_node_count = 1
}

# Small Linux node pool to run some Linux-only Kubernetes Pods.
resource "google_container_node_pool" "linux_pool" {
  name               = "linux-pool"
  project            = google_container_cluster.demo_cluster.project
  cluster            = google_container_cluster.demo_cluster.name
  location           = google_container_cluster.demo_cluster.location

  node_config {
    image_type   = "COS_CONTAINERD"
  }
}

# Node pool of Windows Server machines.
resource "google_container_node_pool" "windows_pool" {
  name               = "windows-pool"
  project            = google_container_cluster.demo_cluster.project
  cluster            = google_container_cluster.demo_cluster.name
  location           = google_container_cluster.demo_cluster.location

  node_config {
    machine_type = "e2-standard-4"
    image_type   = "WINDOWS_LTSC" # Or WINDOWS_SAC for new features.
  }

  # The Linux node pool must be created before the Windows Server node pool.
  depends_on = [google_container_node_pool.linux_pool]
}
```

The example above creates a cluster with a small Linux node pool and a Windows
Server node pool. The Linux node pool is necessary since some critical pods are
not yet supported on Windows. Please see
[Limitations](https://cloud.google.com/kubernetes-engine/docs/how-to/creating-a-cluster-windows#limitations)
for details on features that are not supported by Windows Server node pools.
