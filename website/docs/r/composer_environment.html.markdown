---
subcategory: "Cloud Composer"
layout: "google"
page_title: "Google: google_composer_environment"
sidebar_current: "docs-google-composer-environment"
description: |-
  An environment for running orchestration tasks.
---

# google\_composer\_environment

An environment for running orchestration tasks.

Environments run Apache Airflow software on Google infrastructure.

To get more information about Environments, see:

* [API documentation](https://cloud.google.com/composer/docs/reference/rest/)
* How-to Guides
    * [Official Documentation](https://cloud.google.com/composer/docs)
    * [Configuring Shared VPC for Composer Environments](https://cloud.google.com/composer/docs/how-to/managing/configuring-shared-vpc)
* [Apache Airflow Documentation](http://airflow.apache.org/)

~> **Warning:** We **STRONGLY** recommend  you read the [GCP guides](https://cloud.google.com/composer/docs/how-to)
  as the Environment resource requires a long deployment process and involves several layers of GCP infrastructure, 
  including a Kubernetes Engine cluster, Cloud Storage, and Compute networking resources. Due to limitations of the API,
  Terraform will not be able to automatically find or manage many of these underlying resources. In particular:
  * It can take up to one hour to create or update an environment resource. In addition, GCP may only detect some 
    errors in configuration when they are used (e.g. ~40-50 minutes into the creation process), and is prone to limited
    error reporting. If you encounter confusing or uninformative errors, please verify your configuration is valid 
    against GCP Cloud Composer before filing bugs against the Terraform provider. 
  * **Environments create Google Cloud Storage buckets that do not get cleaned up automatically** on environment 
    deletion. [More about Composer's use of Cloud Storage](https://cloud.google.com/composer/docs/concepts/cloud-storage).

## Example Usage

### Basic Usage
```hcl
resource "google_composer_environment" "test" {
  name   = "my-composer-env"
  region = "us-central1"
}
```

### With GKE and Compute Resource Dependencies

**NOTE** To use service accounts, you need to give `role/composer.worker` to the service account on any resources that may be created for the environment
(i.e. at a project level). This will probably require an explicit dependency
on the IAM policy binding (see `google_project_iam_member` below).

```hcl
resource "google_composer_environment" "test" {
  name   = "mycomposer"
  region = "us-central1"
  config {
    node_count = 4

    node_config {
      zone         = "us-central1-a"
      machine_type = "n1-standard-1"

      network    = google_compute_network.test.id
      subnetwork = google_compute_subnetwork.test.id

      service_account = google_service_account.test.name
    }
  }

  depends_on = [google_project_iam_member.composer-worker]
}

resource "google_compute_network" "test" {
  name                    = "composer-test-network"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "test" {
  name          = "composer-test-subnetwork"
  ip_cidr_range = "10.2.0.0/16"
  region        = "us-central1"
  network       = google_compute_network.test.id
}

resource "google_service_account" "test" {
  account_id   = "composer-env-account"
  display_name = "Test Service Account for Composer Environment"
}

resource "google_project_iam_member" "composer-worker" {
  role   = "roles/composer.worker"
  member = "serviceAccount:${google_service_account.test.email}"
}
```

### With Software (Airflow) Config
```hcl
resource "google_composer_environment" "test" {
  name   = "mycomposer"
  region = "us-central1"

  config {
    software_config {
      airflow_config_overrides = {
        core-load_example = "True"
      }

      pypi_packages = {
        numpy = ""
        scipy = "==1.1.0"
      }

      env_variables = {
        FOO = "bar"
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:


* `name` -
  (Required)
  Name of the environment


- - -
* `config` -
  (Optional)
  Configuration parameters for this environment  Structure is documented below.

* `labels` -
  (Optional)
  User-defined labels for this environment. The labels map can contain
  no more than 64 entries. Entries of the labels map are UTF8 strings
  that comply with the following restrictions:
  Label keys must be between 1 and 63 characters long and must conform
  to the following regular expression: `[a-z]([-a-z0-9]*[a-z0-9])?`.
  Label values must be between 0 and 63 characters long and must
  conform to the regular expression `([a-z]([-a-z0-9]*[a-z0-9])?)?`.
  No more than 64 labels can be associated with a given environment.
  Both keys and values must be <= 128 bytes in size.

* `region` -
  (Optional)
  The location or Compute Engine region for the environment.
* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.

The `config` block supports:

* `node_count` -
  (Optional)
  The number of nodes in the Kubernetes Engine cluster that
  will be used to run this environment.

* `node_config` -
  (Optional)
  The configuration used for the Kubernetes Engine cluster.  Structure is documented below.

* `software_config` -
  (Optional)
  The configuration settings for software inside the environment.  Structure is documented below.

* `private_environment_config` -
  (Optional)
  The configuration used for the Private IP Cloud Composer environment. Structure is documented below.

* `web_server_network_access_control` -
  (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html))
  The network-level access control policy for the Airflow web server. If unspecified, no network-level access restrictions will be applied.


The `node_config` block supports:

* `zone` -
  (Required)
  The Compute Engine zone in which to deploy the VMs running the
  Apache Airflow software, specified as the zone name or
  relative resource name (e.g. "projects/{project}/zones/{zone}"). Must belong to the enclosing environment's project 
  and region.

* `machine_type` -
  (Optional)
  The Compute Engine machine type used for cluster instances,
  specified as a name or relative resource name. For example:
  "projects/{project}/zones/{zone}/machineTypes/{machineType}". Must belong to the enclosing environment's project and 
  region/zone.

* `network` -
  (Optional)
  The Compute Engine network to be used for machine
  communications, specified as a self-link, relative resource name 
  (e.g. "projects/{project}/global/networks/{network}"), by name.

  The network must belong to the environment's project. If unspecified, the "default" network ID in the environment's 
  project is used. If a Custom Subnet Network is provided, subnetwork must also be provided.

* `subnetwork` -
  (Optional)
  The Compute Engine subnetwork to be used for machine
  communications, , specified as a self-link, relative resource name (e.g.
  "projects/{project}/regions/{region}/subnetworks/{subnetwork}"), or by name. If subnetwork is provided, 
  network must also be provided and the subnetwork must belong to the enclosing environment's project and region.

* `disk_size_gb` -
  (Optional)
  The disk size in GB used for node VMs. Minimum size is 20GB.
  If unspecified, defaults to 100GB. Cannot be updated.

* `oauth_scopes` -
  (Optional)
  The set of Google API scopes to be made available on all node
  VMs. Cannot be updated. If empty, defaults to
  `["https://www.googleapis.com/auth/cloud-platform"]`

* `service_account` -
  (Optional)
  The Google Cloud Platform Service Account to be used by the
  node VMs. If a service account is not specified, the "default"
  Compute Engine service account is used. Cannot be updated. If given,
  note that the service account must have `roles/composer.worker` 
  for any GCP resources created under the Cloud Composer Environment.

* `tags` -
  (Optional)
  The list of instance tags applied to all node VMs. Tags are
  used to identify valid sources or targets for network
  firewalls. Each tag within the list must comply with RFC1035.
  Cannot be updated.

* `ip_allocation_policy` -
  (Optional)
  Configuration for controlling how IPs are allocated in the GKE cluster.
  Structure is documented below.
  Cannot be updated.

The `software_config` block supports:

* `airflow_config_overrides` -
  (Optional) Apache Airflow configuration properties to override. Property keys contain the section and property names, 
  separated by a hyphen, for example "core-dags_are_paused_at_creation".

  Section names must not contain hyphens ("-"), opening square brackets ("["), or closing square brackets ("]"). 
  The property name must not be empty and cannot contain "=" or ";". Section and property names cannot contain 
  characters: "." Apache Airflow configuration property names must be written in snake_case. Property values can 
  contain any character, and can be written in any lower/upper case format. Certain Apache Airflow configuration 
  property values are [blacklisted](https://cloud.google.com/composer/docs/concepts/airflow-configurations#airflow_configuration_blacklists), 
  and cannot be overridden.

* `pypi_packages` -
  (Optional)
  Custom Python Package Index (PyPI) packages to be installed
  in the environment. Keys refer to the lowercase package name (e.g. "numpy"). Values are the lowercase extras and 
  version specifier (e.g. "==1.12.0", "[devel,gcp_api]", "[devel]>=1.8.2, <1.9.2"). To specify a package without 
  pinning it to a version specifier, use the empty string as the value.

* `env_variables` -
  (Optional)
  Additional environment variables to provide to the Apache Airflow scheduler, worker, and webserver processes. 
  Environment variable names must match the regular expression `[a-zA-Z_][a-zA-Z0-9_]*`. 
  They cannot specify Apache Airflow software configuration overrides (they cannot match the regular expression 
  `AIRFLOW__[A-Z0-9_]+__[A-Z0-9_]+`), and they cannot match any of the following reserved names:
  ```
  AIRFLOW_HOME
  C_FORCE_ROOT
  CONTAINER_NAME
  DAGS_FOLDER
  GCP_PROJECT
  GCS_BUCKET
  GKE_CLUSTER_NAME
  SQL_DATABASE
  SQL_INSTANCE
  SQL_PASSWORD
  SQL_PROJECT
  SQL_REGION
  SQL_USER
  ```

* `image_version` (Optional) -
  The version of the software running in the environment. This encapsulates both the version of Cloud Composer
  functionality and the version of Apache Airflow. It must match the regular expression 
  `composer-[0-9]+\.[0-9]+(\.[0-9]+)?-airflow-[0-9]+\.[0-9]+(\.[0-9]+.*)?`.
  The Cloud Composer portion of the version is a semantic version. 
  The portion of the image version following 'airflow-' is an official Apache Airflow repository release name.
  See [documentation](https://cloud.google.com/composer/docs/reference/rest/v1beta1/projects.locations.environments#softwareconfig)
  for allowed release names.

* `python_version` (Optional) -
  The major version of Python used to run the Apache Airflow scheduler, worker, and webserver processes.
  Can be set to '2' or '3'. If not specified, the default is '2'. Cannot be updated.

See [documentation](https://cloud.google.com/composer/docs/how-to/managing/configuring-private-ip) for setting up private environments. The `private_environment_config` block supports:

* `enable_private_endpoint` -
  If true, access to the public endpoint of the GKE cluster is denied.

* `master_ipv4_cidr_block` -
  (Optional)
  The IP range in CIDR notation to use for the hosted master network. This range is used
  for assigning internal IP addresses to the cluster master or set of masters and to the
  internal load balancer virtual IP. This range must not overlap with any other ranges
  in use within the cluster's network.
  If left blank, the default value of '172.16.0.0/28' is used.

* `cloud_sql_ipv4_cidr_block` -
  (Optional)
  The CIDR block from which IP range in tenant project will be reserved for Cloud SQL. Needs to be disjoint from `web_server_ipv4_cidr_block`

* `web_server_ipv4_cidr_block` -
  (Optional)
  The CIDR block from which IP range for web server will be reserved. Needs to be disjoint from `master_ipv4_cidr_block` and `cloud_sql_ipv4_cidr_block`.

The `web_server_network_access_control` supports:

* `allowed_ip_range` -
  A collection of allowed IP ranges with descriptions. Structure is documented below.

The `allowed_ip_range` supports:

* `value` -
  (Required)
  IP address or range, defined using CIDR notation, of requests that this rule applies to.
  Examples: `192.168.1.1` or `192.168.0.0/16` or `2001:db8::/32` or `2001:0db8:0000:0042:0000:8a2e:0370:7334`.
  IP range prefixes should be properly truncated. For example,
  `1.2.3.4/24` should be truncated to `1.2.3.0/24`. Similarly, for IPv6, `2001:db8::1/32` should be truncated to `2001:db8::/32`.

* `description` -
  (Optional)
  A description of this ip range.

The `ip_allocation_policy` block supports:

* `use_ip_aliases` -
  (Required)
  Whether or not to enable Alias IPs in the GKE cluster. If true, a VPC-native cluster is created.
  Defaults to true if the `ip_allocation_block` is present in config.

* `cluster_secondary_range_name` -
  (Optional)
  The name of the cluster's secondary range used to allocate IP addresses to pods.
  Specify either `cluster_secondary_range_name` or `cluster_ipv4_cidr_block` but not both.
  This field is applicable only when `use_ip_aliases` is true.

* `services_secondary_range_name` -
  (Optional)
  The name of the services' secondary range used to allocate IP addresses to the cluster.
  Specify either `services_secondary_range_name` or `services_ipv4_cidr_block` but not both.
  This field is applicable only when `use_ip_aliases` is true.

* `cluster_ipv4_cidr_block` -
  (Optional)
  The IP address range used to allocate IP addresses to pods in the cluster.
  Set to blank to have GKE choose a range with the default size.
  Set to /netmask (e.g. /14) to have GKE choose a range with a specific netmask.
  Set to a CIDR notation (e.g. 10.96.0.0/14) from the RFC-1918 private networks
  (e.g. 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) to pick a specific range to use.
  Specify either `cluster_secondary_range_name` or `cluster_ipv4_cidr_block` but not both.

* `services_ipv4_cidr_block` -
  (Optional)
  The IP address range used to allocate IP addresses in this cluster.
  Set to blank to have GKE choose a range with the default size.
  Set to /netmask (e.g. /14) to have GKE choose a range with a specific netmask.
  Set to a CIDR notation (e.g. 10.96.0.0/14) from the RFC-1918 private networks
  (e.g. 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) to pick a specific range to use.
  Specify either `services_secondary_range_name` or `services_ipv4_cidr_block` but not both.



## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{region}}/environments/{{name}}`

* `config.0.gke_cluster` -
  The Kubernetes Engine cluster used to run this environment.

* `config.0.dag_gcs_prefix` -
  The Cloud Storage prefix of the DAGs for this environment.
  Although Cloud Storage objects reside in a flat namespace, a
  hierarchical file tree can be simulated using '/'-delimited
  object name prefixes. DAG objects for this environment
  reside in a simulated directory with this prefix.

* `config.0.airflow_uri` -
  The URI of the Apache Airflow Web UI hosted within this
  environment.

## Timeouts

This resource provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - Default is 60 minutes.
- `update` - Default is 60 minutes.
- `delete` - Default is 6 minutes.

## Import

Environment can be imported using any of these accepted formats:

```
$ terraform import google_composer_environment.default projects/{{project}}/locations/{{region}}/environments/{{name}}
$ terraform import google_composer_environment.default {{project}}/{{region}}/{{name}}
$ terraform import google_composer_environment.default {{name}}
```
