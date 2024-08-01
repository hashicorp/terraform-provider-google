---
subcategory: "Cloud Composer"
description: |-
  An environment for running orchestration tasks.
---

# google_composer_environment

An environment for running orchestration tasks.

Environments run Apache Airflow software on Google infrastructure.

To get more information about Environments, see:

* [Cloud Composer documentation](https://cloud.google.com/composer/docs)
* [Cloud Composer API documentation](https://cloud.google.com/composer/docs/reference/rest/v1beta1/projects.locations.environments)
* How-to Guides (Cloud Composer 2)
  * [Creating environments](https://cloud.google.com/composer/docs/composer-2/create-environments)
  * [Scaling environments](https://cloud.google.com/composer/docs/composer-2/scale-environments)
  * [Configuring Shared VPC for Composer Environments](https://cloud.google.com/composer/docs/composer-2/configure-shared-vpc)
* How-to Guides (Cloud Composer 3)
  * [Creating environments](https://cloud.google.com/composer/docs/composer-3/create-environments)
  * [Scaling environments](https://cloud.google.com/composer/docs/composer-3/scale-environments)
  * [Change environment networking type (Private or Public IP)](https://cloud.google.com/composer/docs/composer-3/change-networking-type)
  * [Connect an environment to a VPC network](https://cloud.google.com/composer/docs/composer-3/connect-vpc-network)
* [Apache Airflow Documentation](http://airflow.apache.org/)

-> **Note**
  Cloud Composer 1 is in the post-maintenance mode. Google does 
  not release any further updates to Cloud Composer 1, including new versions 
  of Airflow, bugfixes, and security updates. We recommend using
  Cloud Composer 2 or Cloud Composer 3 instead.

Several special considerations apply to managing Cloud Composer environments 
with Terraform:

* The Environment resource is based on several layers of GCP infrastructure. 
    Terraform does not manage these underlying resources. For example, in Cloud 
    Composer 2, this includes a Kubernetes Engine cluster, Cloud Storage, and 
    Compute networking resources.
* Creating or updating an environment usually takes around 25 minutes.
* In some cases errors in the configuration will be detected and reported only 
    during the process of the environment creation. If you encounter such 
    errors, please verify your configuration is valid against GCP Cloud Composer before filing bugs for the Terraform provider.
* **Environments have Google Cloud Storage buckets that are not automatically 
    deleted** with the environment.
    See [Delete environments](https://cloud.google.com/composer/docs/composer-2/delete-environments)
    for more information.
* Please refer to
    [Troubleshooting pages](https://cloud.devsite.corp.google.com/composer/docs/composer-2/troubleshooting-environment-creation) if you encounter
    problems.

## Example Usage

### Basic Usage (Cloud Composer 3)
```hcl
resource "google_composer_environment" "test" {
  name   = "example-composer-env"
  region = "us-central1"
 config {
    software_config {
      image_version = "composer-3-airflow-2"
    }
  }
}
```

### Basic Usage (Cloud Composer 2)
```hcl
resource "google_composer_environment" "test" {
  name   = "example-composer-env"
  region = "us-central1"
 config {
    software_config {
      image_version = "composer-2-airflow-2"
    }
  }
}
```

### Basic Usage (Cloud Composer 1)
```hcl
resource "google_composer_environment" "test" {
  name   = "example-composer-env"
  region = "us-central1"
  config {
    software_config {
      image_version = "composer-1-airflow-2"
    }
  }
}
```

### With GKE and Compute Resource Dependencies

-> **Note**
  To use custom service accounts, you must give at least the
  `role/composer.worker` role to the service account of the Cloud Composer 
  environment. For more information, see the
  [Access Control](https://cloud.google.com/composer/docs/how-to/access-control)
  page in the Cloud Composer documentation.
  You might need to assign additional roles depending on specific workflows 
  that the Airflow DAGs will be running.

#### GKE and Compute Resource Dependencies (Cloud Composer 3)

```hcl
provider "google" {
  project = "example-project"
}

resource "google_composer_environment" "test" {
  name   = "example-composer-env-tf-c3"
  region = "us-central1"
  config {

    software_config {
      image_version = "composer-3-airflow-2"
    }

    workloads_config {
      scheduler {
        cpu        = 0.5
        memory_gb  = 2
        storage_gb = 1
        count      = 1
      }
      triggerer {
        cpu        = 0.5
        memory_gb  = 1
        count      = 1
      }
      dag_processor {
        cpu        = 1
        memory_gb  = 2
        storage_gb = 1
        count      = 1
      }
      web_server {
        cpu        = 0.5
        memory_gb  = 2
        storage_gb = 1
      }
      worker {
        cpu = 0.5
        memory_gb  = 2
        storage_gb = 1
        min_count  = 1
        max_count  = 3
      }

    }
    environment_size = "ENVIRONMENT_SIZE_SMALL"

    node_config {
      service_account = google_service_account.test.name
    }
  }
}

resource "google_service_account" "test" {
  account_id   = "composer-env-account"
  display_name = "Test Service Account for Composer Environment"
}

resource "google_project_iam_member" "composer-worker" {
  project = "your-project-id"
  role    = "roles/composer.worker"
  member  = "serviceAccount:${google_service_account.test.email}"
}
```

#### GKE and Compute Resource Dependencies (Cloud Composer 2)

```hcl
provider "google" {
  project = "example-project"
}

resource "google_composer_environment" "test" {
  name   = "example-composer-env-tf-c2"
  region = "us-central1"
  config {

    software_config {
      image_version = "composer-2-airflow-2"
    }

    workloads_config {
      scheduler {
        cpu        = 0.5
        memory_gb  = 1.875
        storage_gb = 1
        count      = 1
      }
      web_server {
        cpu        = 0.5
        memory_gb  = 1.875
        storage_gb = 1
      }
      worker {
        cpu = 0.5
        memory_gb  = 1.875
        storage_gb = 1
        min_count  = 1
        max_count  = 3
      }


    }
    environment_size = "ENVIRONMENT_SIZE_SMALL"

    node_config {
      network    = google_compute_network.test.id
      subnetwork = google_compute_subnetwork.test.id
      service_account = google_service_account.test.name
    }
  }
}

resource "google_compute_network" "test" {
  name                    = "composer-test-network3"
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
  project = "your-project-id"
  role    = "roles/composer.worker"
  member  = "serviceAccount:${google_service_account.test.email}"
}
```

#### GKE and Compute Resource Dependencies (Cloud Composer 1)

```hcl
resource "google_composer_environment" "test" {
  name   = "example-composer-env"
  region = "us-central1"
  config {
    
    software_config {
      image_version = "composer-1-airflow-2"
    }
    
    node_count = 4

    node_config {
      zone         = "us-central1-a"
      machine_type = "n1-standard-1"

      network    = google_compute_network.test.id
      subnetwork = google_compute_subnetwork.test.id

      service_account = google_service_account.test.name
    }

    database_config {
      machine_type = "db-n1-standard-2"
    }

    web_server_config {
      machine_type = "composer-n1-webserver-2"
    }
  }
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

### Cloud Composer 3 networking configuration

In Cloud Composer 3, networking configuration is simplified compared to
previous versions. You don't need to specify network ranges, and can attach
custom VPC networks to your environment.

-> **Note**
  It's not possible to detach a VPC network using Terraform. Instead, you can
  attach a different VPC network in its place, or detach the network using
  other tools like Google Cloud CLI.

Use Private IP networking:

```hcl
resource "google_composer_environment" "example" {
  name = "example-environment"
  region = "us-central1"

  config {

    enable_private_ip_environment = true

    # ... other configuration parameters
  }
}
```

Attach a custom VPC network (Cloud Composer creates a new network attachment):

```hcl
resource "google_composer_environment" "example" {
  name = "example-environment"
  region = "us-central1"

  config {

    node_config {
      network = "projects/example-project/global/networks/example-network"
      subnetwork = "projects/example-project/regions/us-central1/subnetworks/example-subnetwork"
    }

    # ... other configuration parameters
  }
}
```

Attach a custom VPC network (use existing network attachment):

```hcl
resource "google_composer_environment" "example" {
  name = "example-environment"
  region = "us-central1"

  config {

    node_config {
      composer_network_attachment = projects/example-project/regions/us-central1/networkAttachments/example-network-attachment
    }

    # ... other configuration parameters
  }
}
```

### With Software (Airflow) Config

```hcl
resource "google_composer_environment" "test" {
  name   = "mycomposer"
  region = "us-central1"

  config {
      airflow_config_overrides = {
        core-dags_are_paused_at_creation = "True"
      }

      pypi_packages = {
        numpy = ""
        scipy = "==1.1.0"
      }

      env_variables = {
        EXAMPLE_VARIABLE = "test"
      }
    }
  }
}
```

## Argument Reference - Cloud Composer 1

The following arguments are supported:

* `name` -
  (Required)
  Name of the environment

* `config` -
  (Optional)
  Configuration parameters for this environment  Structure is [documented below](#nested_config_c1).

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

  **Note**: This field is non-authoritative, and will only manage the labels present in your configuration.
  Please refer to the field 'effective_labels' for all of the labels present on the resource.

* `terraform_labels` -
  The combination of labels configured directly on the resource and default labels configured on the provider.

* `effective_labels` -
  All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.

* `region` -
  (Optional)
  The location or Compute Engine region for the environment.

* `project` -
  (Optional) The ID of the project in which the resource belongs.
  If it is not provided, the provider project is used.

<a name="nested_config_c1"></a>The `config` block supports:

* `node_count` -
  (Optional, Cloud Composer 1 only)
  The number of nodes in the Kubernetes Engine cluster of the environment.

* `node_config` -
  (Optional)
  The configuration used for the Kubernetes Engine cluster.  Structure is [documented below](#nested_node_config_c1).

* `software_config` -
  (Optional)
  The configuration settings for software inside the environment.  Structure is [documented below](#nested_software_config_c1).

* `private_environment_config` -
  (Optional)
  The configuration used for the Private IP Cloud Composer environment. Structure is [documented below](#nested_private_environment_config_c1).

* `web_server_network_access_control` -
  The network-level access control policy for the Airflow web server.
  If unspecified, no network-level access restrictions are applied.

* `database_config` -
  (Optional, Cloud Composer 1 only)
  The configuration settings for Cloud SQL instance used internally
  by Apache Airflow software.

* `web_server_config` -
  (Optional, Cloud Composer 1 only)
  The configuration settings for the Airflow web server App Engine instance.

* `encryption_config` -
  (Optional)
  The encryption options for the Cloud Composer environment and its
  dependencies.

* `maintenance_window` -
  (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html))
  The configuration settings for Cloud Composer maintenance windows.

* `master_authorized_networks_config` -
  (Optional)
  Configuration options for the master authorized networks feature. Enabled
  master authorized networks will disallow all external traffic to access
  Kubernetes master through HTTPS except traffic from the given CIDR blocks,
  Google Compute Engine Public IPs and Google Prod IPs. Structure is
  [documented below](#nested_master_authorized_networks_config_c1).

<a name="nested_node_config_c1"></a>The `node_config` block supports:

* `zone` -
  (Optional, Cloud Composer 1 only)
  The Compute Engine zone in which to deploy the VMs running the
  Apache Airflow software, specified as the zone name or
  relative resource name (e.g. "projects/{project}/zones/{zone}"). Must
  belong to the enclosing environment's project and region.

* `machine_type` -
  (Optional, Cloud Composer 1 only)
  The Compute Engine machine type used for cluster instances,
  specified as a name or relative resource name. For example:
  "projects/{project}/zones/{zone}/machineTypes/{machineType}". Must belong
  to the enclosing environment's project and region/zone.

* `network` -
  (Optional)
  The Compute Engine network to be used for machine
  communications, specified as a self-link, relative resource name
  (for example "projects/{project}/global/networks/{network}"), by name.

  The network must belong to the environment's project. If unspecified, the "default" network ID in the environment's
  project is used. If a Custom Subnet Network is provided, subnetwork must also be provided.

* `subnetwork` -
  (Optional)
  The Compute Engine subnetwork to be used for machine
  communications, specified as a self-link, relative resource name (for example,
  "projects/{project}/regions/{region}/subnetworks/{subnetwork}"), or by name. If subnetwork is provided,
  network must also be provided and the subnetwork must belong to the enclosing environment's project and region.

* `disk_size_gb` -
  (Optional, Cloud Composer 1 only)
  The disk size in GB used for node VMs. Minimum size is 20GB.
  If unspecified, defaults to 100GB. Cannot be updated.

* `oauth_scopes` -
  (Optional, Cloud Composer 1 only)
  The set of Google API scopes to be made available on all node
  VMs. Cannot be updated. If empty, defaults to
  `["https://www.googleapis.com/auth/cloud-platform"]`.

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
  Structure is [documented below](#nested_ip_allocation_policy_c1).
  Cannot be updated.

* `max_pods_per_node` -
  (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html),
  Cloud Composer 1 only)
  The maximum pods per node in the GKE cluster allocated during environment
  creation. Lowering this value reduces IP address consumption by the Cloud
  Composer Kubernetes cluster. This value can only be set if the environment is VPC-Native.
  The range of possible values is 8-110, and the default is 32.
  Cannot be updated.

* `enable_ip_masq_agent` -
  (Optional)
  Deploys 'ip-masq-agent' daemon set in the GKE cluster and defines
  nonMasqueradeCIDRs equals to pod IP range so IP masquerading is used for
  all destination addresses, except between pods traffic.
  See the [documentation](https://cloud.google.com/composer/docs/enable-ip-masquerade-agent).

<a name="nested_software_config_c1"></a>The `software_config` block supports:

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
  AIRFLOW_DATABASE_VERSION
  AIRFLOW_HOME
  AIRFLOW_SRC_DIR
  AIRFLOW_WEBSERVER
  AUTO_GKE
  CLOUDSDK_METRICS_ENVIRONMENT
  CLOUD_LOGGING_ONLY
  COMPOSER_ENVIRONMENT
  COMPOSER_GKE_LOCATION
  COMPOSER_GKE_NAME
  COMPOSER_GKE_ZONE
  COMPOSER_LOCATION
  COMPOSER_OPERATION_UUID
  COMPOSER_PYTHON_VERSION
  COMPOSER_VERSION
  CONTAINER_NAME
  C_FORCE_ROOT
  DAGS_FOLDER
  GCP_PROJECT
  GCP_TENANT_PROJECT
  GCSFUSE_EXTRACTED
  GCS_BUCKET
  GKE_CLUSTER_NAME
  GKE_IN_TENANT
  GOOGLE_APPLICATION_CREDENTIALS
  MAJOR_VERSION
  MINOR_VERSION
  PATH
  PIP_DISABLE_PIP_VERSION_CHECK
  PORT
  PROJECT_ID
  PYTHONPYCACHEPREFIX
  SQL_DATABASE
  SQL_HOST
  SQL_INSTANCE
  SQL_PASSWORD
  SQL_PROJECT
  SQL_REGION
  SQL_USER
  ```

* `image_version` -
(Required) In Composer 1, use a specific Composer 1 version in this parameter. If omitted, the default is the latest version of Composer 2.  

  The version of the software running in the environment. This encapsulates both the version of Cloud Composer
  functionality and the version of Apache Airflow. It must match the regular expression
  `composer-([0-9]+(\.[0-9]+\.[0-9]+(-preview\.[0-9]+)?)?|latest)-airflow-([0-9]+(\.[0-9]+(\.[0-9]+)?)?)`.
  The Cloud Composer portion of the image version is a full semantic version, or an alias in the form of major
  version number or 'latest'.
  The Apache Airflow portion of the image version is a full semantic version that points to one of the
  supported Apache Airflow versions, or an alias in the form of only major or major.minor versions specified.
  For more information about Cloud Composer images, see
  [Cloud Composer version list](https://cloud.google.com/composer/docs/concepts/versioning/composer-versions).

* `python_version` -
  (Optional, Cloud Composer 1 only)
  The major version of Python used to run the Apache Airflow scheduler, worker, and webserver processes.
  Can be set to '2' or '3'. If not specified, the default is '3'.

* `scheduler_count` -
  (Optional, Cloud Composer 1 with Airflow 2 only)
  The number of schedulers for Airflow.


See [documentation](https://cloud.google.com/composer/docs/how-to/managing/configuring-private-ip) for setting up private environments. <a name="nested_private_environment_config_c1"></a>The `private_environment_config` block supports:

* `enable_private_endpoint` -
  If true, access to the public endpoint of the GKE cluster is denied.
  If this field is set to true, the `ip_allocation_policy.use_ip_aliases` field must
  also be set to true for Cloud Composer 1 environments.

* `master_ipv4_cidr_block` -
  (Optional)
  The IP range in CIDR notation to use for the hosted master network. This range is used
  for assigning internal IP addresses to the cluster master or set of masters and to the
  internal load balancer virtual IP. This range must not overlap with any other ranges
  in use within the cluster's network.
  If left blank, the default value of is used. See [documentation](https://cloud.google.com/composer/docs/how-to/managing/configuring-private-ip#defaults) for default values per region.

* `cloud_sql_ipv4_cidr_block` -
  (Optional)
  The CIDR block from which IP range in tenant project will be reserved for Cloud SQL. Needs to be disjoint from `web_server_ipv4_cidr_block`

* `web_server_ipv4_cidr_block` -
  (Optional, Cloud Composer 1 only)
  The CIDR block from which IP range for web server will be reserved. Needs to be disjoint from `master_ipv4_cidr_block` and `cloud_sql_ipv4_cidr_block`.

* `enable_privately_used_public_ips` -
  (Optional)
  When enabled, IPs from public (non-RFC1918) ranges can be used for
  `ip_allocation_policy.cluster_ipv4_cidr_block` and `ip_allocation_policy.service_ipv4_cidr_block`.

The `web_server_network_access_control` supports:

* `allowed_ip_range` -
  A collection of allowed IP ranges with descriptions. Structure is [documented below](#nested_allowed_ip_range_c1).

<a name="nested_allowed_ip_range_c1"></a>The `allowed_ip_range` supports:

* `value` -
  (Required)
  IP address or range, defined using CIDR notation, of requests that this rule applies to.
  Examples: `192.168.1.1` or `192.168.0.0/16` or `2001:db8::/32` or `2001:0db8:0000:0042:0000:8a2e:0370:7334`.
  IP range prefixes should be properly truncated. For example,
  `1.2.3.4/24` should be truncated to `1.2.3.0/24`. Similarly, for IPv6, `2001:db8::1/32` should be truncated to `2001:db8::/32`.

* `description` -
  (Optional)
  A description of this ip range.

<a name="nested_ip_allocation_policy_c1"></a>The `ip_allocation_policy` block supports:

* `use_ip_aliases` -
  (Optional, Cloud Composer 1 only)
  Whether or not to enable Alias IPs in the GKE cluster. If true, a VPC-native cluster is created.
  Defaults to true if the `ip_allocation_policy` block is present in config.

* `cluster_secondary_range_name` -
  (Optional)
  The name of the cluster's secondary range used to allocate IP addresses to pods.
  Specify either `cluster_secondary_range_name` or `cluster_ipv4_cidr_block` but not both.
  For Cloud Composer 1 environments, this field is applicable only when `use_ip_aliases` is true.

* `services_secondary_range_name` -
  (Optional)
  The name of the services' secondary range used to allocate IP addresses to the cluster.
  Specify either `services_secondary_range_name` or `services_ipv4_cidr_block` but not both.
  For Cloud Composer 1 environments, this field is applicable only when `use_ip_aliases` is true.

* `cluster_ipv4_cidr_block` -
  (Optional)
  The IP address range used to allocate IP addresses to pods in the cluster.
  For Cloud Composer 1 environments, this field is applicable only when `use_ip_aliases` is true.
  Set to blank to have GKE choose a range with the default size.
  Set to /netmask (e.g. /14) to have GKE choose a range with a specific netmask.
  Set to a CIDR notation (e.g. 10.96.0.0/14) from the RFC-1918 private networks
  (e.g. 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) to pick a specific range to use.
  Specify either `cluster_secondary_range_name` or `cluster_ipv4_cidr_block` but not both.

* `services_ipv4_cidr_block` -
  (Optional)
  The IP address range used to allocate IP addresses in this cluster.
  For Cloud Composer 1 environments, this field is applicable only when `use_ip_aliases` is true.
  Set to blank to have GKE choose a range with the default size.
  Set to /netmask (e.g. /14) to have GKE choose a range with a specific netmask.
  Set to a CIDR notation (e.g. 10.96.0.0/14) from the RFC-1918 private networks
  (e.g. 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) to pick a specific range to use.
  Specify either `services_secondary_range_name` or `services_ipv4_cidr_block` but not both.

<a name="nested_database_config_c1"></a>The `database_config` block supports:

* `machine_type` -
  (Optional)
  Optional. Cloud SQL machine type used by Airflow database. It has to be one of: db-n1-standard-2,
  db-n1-standard-4, db-n1-standard-8 or db-n1-standard-16.

* `Zone` -
  (Optional)
  Preferred Cloud SQL database zone.

<a name="nested_web_server_config_c1"></a>The `web_server_config` block supports:

* `machine_type` -
  (Required)
  Machine type on which Airflow web server is running. It has to be one of: composer-n1-webserver-2,
  composer-n1-webserver-4 or composer-n1-webserver-8.
  Value custom is returned only in response, if Airflow web server parameters were
  manually changed to a non-standard values.

<a name="nested_encryption_config_c1"></a>The `encryption_config` block supports:

* `kms_key_name` -
  (Required)
  Customer-managed Encryption Key available through Google's Key Management Service. It must
  be the fully qualified resource name,
  i.e. projects/project-id/locations/location/keyRings/keyring/cryptoKeys/key. Cannot be updated.

<a name="nested_maintenance_window_c1"></a>The `maintenance_window` block supports:
* `start_time` -
  (Required)
  Start time of the first recurrence of the maintenance window.

* `end_time` -
  (Required)
  Maintenance window end time. It is used only to calculate the duration of the maintenance window.
  The value for end-time must be in the future, relative to 'start_time'.

* `recurrence` -
  (Required)
  Maintenance window recurrence. Format is a subset of RFC-5545 (https://tools.ietf.org/html/rfc5545) 'RRULE'.
  The only allowed values for 'FREQ' field are 'FREQ=DAILY' and 'FREQ=WEEKLY;BYDAY=...'.
  Example values: 'FREQ=WEEKLY;BYDAY=TU,WE', 'FREQ=DAILY'.

<a name="nested_master_authorized_networks_config_c1"></a>The `master_authorized_networks_config` block supports:
* `enabled` -
  (Required)
  Whether or not master authorized networks is enabled.

* `cidr_blocks` -
  `cidr_blocks `define up to 50 external networks that could access Kubernetes master through HTTPS. Structure is [documented below](#nested_cidr_blocks_c1).

<a name="nested_cidr_blocks_c1"></a>The `cidr_blocks` supports:

* `display_name` -
  (Optional)
  `display_name` is a field for users to identify CIDR blocks.

* `cidr_block` -
  (Required)
  `cidr_block` must be specified in CIDR notation.

## Argument Reference - Cloud Composer 2

The following arguments are supported:

* `name` -
  (Required)
  Name of the environment

* `config` -
  (Optional)
  Configuration parameters for this environment. Structure is [documented below](#nested_config_c2).

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

* `project` -
  (Optional) The ID of the project in which the resource belongs.
  If it is not provided, the provider project is used.

* `storage_config` -
  (Optional)
  Configuration options for storage used by Composer environment. Structure is [documented below](#nested_storage_config_c2).


<a name="nested_config_c2"></a>The `config` block supports:

* `node_config` -
  (Optional)
  The configuration used for the Kubernetes Engine cluster. Structure is [documented below](#nested_node_config_c2).

* `recovery_config` -
  (Optional, Cloud Composer 2 only)
  The configuration settings for recovery. Structure is [documented below](#nested_recovery_config_c2).

* `software_config` -
  (Optional)
  The configuration settings for software (Airflow) inside the environment. Structure is
  [documented below](#nested_software_config_c2).

* `private_environment_config` -
  (Optional)
  The configuration used for the Private IP Cloud Composer environment. Structure is [documented below](#nested_private_environment_config_c2).

* `encryption_config` -
  (Optional)
  The encryption options for the Cloud Composer environment and its
  dependencies.

* `maintenance_window` -
  (Optional)
  The configuration settings for Cloud Composer maintenance windows.

* `workloads_config` -
  (Optional)
  The Kubernetes workloads configuration for GKE cluster associated with the
  Cloud Composer environment.

* `environment_size` -
  (Optional)
  The environment size controls the performance parameters of the managed
  Cloud Composer infrastructure that includes the Airflow database. Values for
  environment size are `ENVIRONMENT_SIZE_SMALL`, `ENVIRONMENT_SIZE_MEDIUM`,
  and `ENVIRONMENT_SIZE_LARGE`.

* `resilience_mode` -
  (Optional, Cloud Composer 2.1.15 or newer only)
  The resilience mode states whether high resilience is enabled for 
  the environment or not. Values for resilience mode are `HIGH_RESILIENCE` 
  for high resilience and `STANDARD_RESILIENCE` for standard
  resilience.

* `master_authorized_networks_config` -
  (Optional)
  Configuration options for the master authorized networks feature. Enabled
  master authorized networks will disallow all external traffic to access
  Kubernetes master through HTTPS except traffic from the given CIDR blocks,
  Google Compute Engine Public IPs and Google Prod IPs. Structure is
  [documented below](#nested_master_authorized_networks_config_c2).

<a name="nested_master_authorized_networks_config_c2"></a>The `master_authorized_networks_config` block supports:
* `enabled` -
  (Required)
  Whether or not master authorized networks is enabled.

* `cidr_blocks` -
  `cidr_blocks `define up to 50 external networks that could access Kubernetes master through HTTPS. Structure is [documented below](#nested_cidr_blocks_c2).

<a name="nested_cidr_blocks_c2"></a>The `cidr_blocks` supports:

* `display_name` -
  (Optional)
  `display_name` is a field for users to identify CIDR blocks.

* `cidr_block` -
  (Required)
  `cidr_block` must be specified in CIDR notation.

* `data_retention_config` -
  (Optional, Cloud Composer 2.0.23 or newer only)
  Configuration setting for airflow data rentention mechanism. Structure is
  [documented below](#nested_data_retention_config_c2).

<a name="nested_data_retention_config_c2"></a>The `data_retention_config` block supports:
* `task_logs_retention_config` - 
  (Optional)
  The configuration setting for Task Logs. Structure is
  [documented below](#nested_task_logs_retention_config_c2).

<a name="nested_task_logs_retention_config_c2"></a>The `task_logs_retention_config` block supports:
* `storage_mode` - 
  (Optional)
  The mode of storage for Airflow workers task logs. Values for storage mode are 
  `CLOUD_LOGGING_ONLY` to only store logs in cloud logging and 
  `CLOUD_LOGGING_AND_CLOUD_STORAGE` to store logs in cloud logging and cloud storage.


<a name="nested_storage_config_c2"></a>The `storage_config` block supports:

* `bucket` -
  (Required)
  Name of an existing Cloud Storage bucket to be used by the environment.


<a name="nested_node_config_c2"></a>The `node_config` block supports:

* `network` -
  (Optional)
  The Compute Engine network to be used for machine
  communications, specified as a self-link, relative resource name
  (for example "projects/{project}/global/networks/{network}"), by name.

  The network must belong to the environment's project. If unspecified, the "default" network ID in the environment's
  project is used. If a Custom Subnet Network is provided, subnetwork must also be provided.

* `subnetwork` -
  (Optional)
  The Compute Engine subnetwork to be used for machine
  communications, specified as a self-link, relative resource name (for example,
  "projects/{project}/regions/{region}/subnetworks/{subnetwork}"), or by name. If subnetwork is provided,
  network must also be provided and the subnetwork must belong to the enclosing environment's project and region.

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
  Structure is [documented below](#nested_ip_allocation_policy_c2).
  Cannot be updated.

* `enable_ip_masq_agent` -
  (Optional)
  IP Masq Agent translates Pod IP addresses to node IP addresses, so that 
  destinations and services targeted from Airflow DAGs and tasks only receive 
  packets from node IP addresses instead of Pod IP addresses
  See the [documentation](https://cloud.google.com/composer/docs/enable-ip-masquerade-agent).

<a name="nested_software_config_c2"></a>The `software_config` block supports:

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

* `image_version` -
(Optional) If omitted, the default is the latest version of Composer 2.


  The version of the software running in the environment. This encapsulates both the version of Cloud Composer
  functionality and the version of Apache Airflow. It must match the regular expression
  `composer-([0-9]+(\.[0-9]+\.[0-9]+(-preview\.[0-9]+)?)?|latest)-airflow-([0-9]+(\.[0-9]+(\.[0-9]+)?)?)`.
  The Cloud Composer portion of the image version is a full semantic version, or an alias in the form of major
  version number or 'latest'.
  The Apache Airflow portion of the image version is a full semantic version that points to one of the
  supported Apache Airflow versions, or an alias in the form of only major or major.minor versions specified.
  **Important**: In-place upgrade is only available using `google-beta` provider. It's because updating the
  `image_version` is still in beta. Using `google-beta` provider, you can upgrade in-place between minor or
  patch versions of Cloud Composer or Apache Airflow. For example, you can upgrade your environment from
  `composer-1.16.x` to `composer-1.17.x`, or from `airflow-2.1.x` to `airflow-2.2.x`. You cannot upgrade between
  major Cloud Composer or Apache Airflow versions (from `1.x.x` to `2.x.x`). To do so, create a new environment.

* `cloud_data_lineage_integration` -
  (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html),
  Cloud Composer environments in versions composer-2.1.2-airflow-*.*.* and newer)
  The configuration for Cloud Data Lineage integration. Structure is
  [documented below](#nested_cloud_data_lineage_integration_c2).

<a name="nested_cloud_data_lineage_integration_c2"></a>The `cloud_data_lineage_integration` block supports:
* `enabled` -
  (Required)
  Whether or not Cloud Data Lineage integration is enabled.

<a name="nested_private_environment_config_c2"></a>See [documentation](https://cloud.google.com/composer/docs/how-to/managing/configuring-private-ip) for setting up private environments. The `private_environment_config` block supports:

* `connection_type` -
  (Optional, Cloud Composer 2 only)
  Mode of internal communication within the Composer environment. Must be one
  of `"VPC_PEERING"` or `"PRIVATE_SERVICE_CONNECT"`.

* `enable_private_endpoint` -
  If true, access to the public endpoint of the GKE cluster is denied.

* `master_ipv4_cidr_block` -
  (Optional)
  The IP range in CIDR notation to use for the hosted master network. This range is used
  for assigning internal IP addresses to the cluster master or set of masters and to the
  internal load balancer virtual IP. This range must not overlap with any other ranges
  in use within the cluster's network.
  If left blank, the default value of is used. See [documentation](https://cloud.google.com/composer/docs/how-to/managing/configuring-private-ip#defaults) for default values per region.

* `cloud_sql_ipv4_cidr_block` -
  (Optional)
  The CIDR block from which IP range in tenant project will be reserved for Cloud SQL. Needs to be disjoint from `web_server_ipv4_cidr_block`

* `cloud_composer_network_ipv4_cidr_block"` -
  (Optional, Cloud Composer 2 only)
  The CIDR block from which IP range for Cloud Composer Network in tenant project will be reserved. Needs to be disjoint from private_cluster_config.master_ipv4_cidr_block and cloud_sql_ipv4_cidr_block.

* `enable_privately_used_public_ips` -
  (Optional)
  When enabled, IPs from public (non-RFC1918) ranges can be used for
  `ip_allocation_policy.cluster_ipv4_cidr_block` and `ip_allocation_policy.service_ipv4_cidr_block`.

* `cloud_composer_connection_subnetwork` -
  (Optional)
  When specified, the environment will use Private Service Connect instead of VPC peerings to connect
  to Cloud SQL in the Tenant Project, and the PSC endpoint in the Customer Project will use an IP
  address from this subnetwork. This field is supported for Cloud Composer environments in
  versions `composer-2.*.*-airflow-*.*.*` and newer.


<a name="nested_ip_allocation_policy_c2"></a>The `ip_allocation_policy` block supports:

* `cluster_secondary_range_name` -
  (Optional)
  The name of the cluster's secondary range used to allocate IP addresses to pods.
  Specify either `cluster_secondary_range_name` or `cluster_ipv4_cidr_block` but not both.

* `services_secondary_range_name` -
  (Optional)
  The name of the services' secondary range used to allocate IP addresses to the cluster.
  Specify either `services_secondary_range_name` or `services_ipv4_cidr_block` but not both.

* `cluster_ipv4_cidr_block` -
  (Optional)
  The IP address range used to allocate IP addresses to pods in the cluster.
  For Cloud Composer 1 environments, this field is applicable only when `use_ip_aliases` is true.
  Set to blank to have GKE choose a range with the default size.
  Set to /netmask (e.g. /14) to have GKE choose a range with a specific netmask.
  Set to a CIDR notation (e.g. 10.96.0.0/14) from the RFC-1918 private networks
  (e.g. 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) to pick a specific range to use.
  Specify either `cluster_secondary_range_name` or `cluster_ipv4_cidr_block` but not both.

* `services_ipv4_cidr_block` -
  (Optional)
  The IP address range used to allocate IP addresses in this cluster.
  For Cloud Composer 1 environments, this field is applicable only when `use_ip_aliases` is true.
  Set to blank to have GKE choose a range with the default size.
  Set to /netmask (e.g. /14) to have GKE choose a range with a specific netmask.
  Set to a CIDR notation (e.g. 10.96.0.0/14) from the RFC-1918 private networks
  (e.g. 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) to pick a specific range to use.
  Specify either `services_secondary_range_name` or `services_ipv4_cidr_block` but not both.

<a name="nested_encryption_config_comp_2"></a>The `encryption_config` block supports:

* `kms_key_name` -
  (Required)
  Customer-managed Encryption Key available through Google's Key Management Service. It must
  be the fully qualified resource name,
  i.e. projects/project-id/locations/location/keyRings/keyring/cryptoKeys/key. Cannot be updated.

<a name="nested_maintenance_window_comp_2"></a>The `maintenance_window` block supports:

* `start_time` -
  (Required)
  Start time of the first recurrence of the maintenance window.

* `end_time` -
  (Required)
  Maintenance window end time. It is used only to calculate the duration of the maintenance window.
  The value for end-time must be in the future, relative to 'start_time'.

* `recurrence` -
  (Required)
  Maintenance window recurrence. Format is a subset of RFC-5545 (https://tools.ietf.org/html/rfc5545) 'RRULE'.
  The only allowed values for 'FREQ' field are 'FREQ=DAILY' and 'FREQ=WEEKLY;BYDAY=...'.
  Example values: 'FREQ=WEEKLY;BYDAY=TU,WE', 'FREQ=DAILY'.

<a name="nested_recovery_config_c2"></a>The `recovery_config` block supports:

* `scheduled_snapshots_config` -
  (Optional)
  The recovery configuration settings for the Cloud Composer environment.

The `scheduled_snapshots_config` block supports:

* `enabled` -
  (Optional)
  When enabled, Cloud Composer periodically saves snapshots of your environment to a Cloud Storage bucket.

* `snapshot_location` -
  (Optional)
  The URI of a bucket folder where to save the snapshot.

* `snapshot_creation_schedule` -
  (Optional)
  Snapshot schedule, in the unix-cron format.

* `time_zone` -
  (Optional)
  A time zone for the schedule. This value is a time offset and does not take into account daylight saving time changes. Valid values are from UTC-12 to UTC+12. Examples: UTC, UTC-01, UTC+03.

The `workloads_config` block supports:

* `scheduler` -
  (Optional)
  Configuration for resources used by Airflow schedulers.

* `triggerer` -
  (Optional)
  Configuration for resources used by Airflow triggerer.

* `web_server` -
  (Optional)
  Configuration for resources used by Airflow web server.

* `worker` -
  (Optional)
  Configuration for resources used by Airflow workers.

The `scheduler` block supports:

* `cpu` -
  (Optional)
  The number of CPUs for a single Airflow scheduler.

* `memory_gb` -
  (Optional)
  The amount of memory (GB) for a single Airflow scheduler.

* `storage_gb` -
  (Optional)
  The amount of storage (GB) for a single Airflow scheduler.

* `count` -
  (Optional)
  The number of schedulers.

The `triggerer` block supports:

* `cpu` -
  (Required)
  The number of CPUs for a single Airflow triggerer.

* `memory_gb` -
  (Required)
  The amount of memory (GB) for a single Airflow triggerer.

* `count` -
  (Required)
  The number of Airflow triggerers.

The `web_server` block supports:

* `cpu` -
  (Optional)
  The number of CPUs for the Airflow web server.

* `memory_gb` -
  (Optional)
  The amount of memory (GB) for the Airflow web server.

* `storage_gb` -
  (Optional)
  The amount of storage (GB) for the Airflow web server.

The `worker` block supports:

* `cpu` -
  (Optional)
  The number of CPUs for a single Airflow worker.

* `memory_gb` -
  (Optional)
  The amount of memory (GB) for a single Airflow worker.

* `storage_gb`
  (Optional)
  The amount of storage (GB) for a single Airflow worker.

* `min_count` -
  (Optional)
  The minimum number of Airflow workers that the environment can run. The number of workers in the
  environment does not go above this number, even if a lower number of workers can handle the load.

* `max_count` -
  (Optional)
  The maximum number of Airflow workers that the environment can run. The number of workers in the
  environment does not go above this number, even if a higher number of workers is required to
  handle the load.


## Argument Reference - Cloud Composer 3

The following arguments are supported:

* `name` -
  (Required)
  Name of the environment

* `config` -
  (Optional)
  Configuration parameters for this environment. Structure is [documented below](#nested_config_c3).

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

* `project` -
  (Optional) The ID of the project in which the resource belongs.
  If it is not provided, the provider project is used.

* `storage_config` -
  (Optional)
  Configuration options for storage used by Composer environment. Structure is [documented below](#nested_storage_config_c3).


<a name="nested_config_c3"></a>The `config` block supports:

* `node_config` -
  (Optional)
  The configuration used for the Kubernetes Engine cluster. Structure is [documented below](#nested_node_config_c3).

* `software_config` -
  (Optional)
  The configuration settings for software (Airflow) inside the environment. Structure is [documented below](#nested_software_config_c3).

* `enable_private_environment` -
  (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html), Cloud Composer 3 only)
  If true, a private Composer environment will be created.

* `enable_private_builds_only` -
  (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html), Cloud Composer 3 only)
  If true, builds performed during operations that install Python packages have only private connectivity to Google services.
  If false, the builds also have access to the internet.

* `encryption_config` -
  (Optional)
  The encryption options for the Cloud Composer environment and its
  dependencies.

* `maintenance_window` -
  (Optional)
  The configuration settings for Cloud Composer maintenance windows.

* `workloads_config` -
  (Optional)
  The Kubernetes workloads configuration for GKE cluster associated with the
  Cloud Composer environment.

* `environment_size` -
  (Optional)
  The environment size controls the performance parameters of the managed
  Cloud Composer infrastructure that includes the Airflow database. Values for
  environment size are `ENVIRONMENT_SIZE_SMALL`, `ENVIRONMENT_SIZE_MEDIUM`,
  and `ENVIRONMENT_SIZE_LARGE`.

* `data_retention_config` -
  (Optional, Cloud Composer 2.0.23 or later only)
  Configuration setting for Airflow database retention mechanism. Structure is
  [documented below](#nested_data_retention_config_c3).

<a name="nested_data_retention_config_c3"></a>The `data_retention_config` block supports:
* `task_logs_retention_config` - 
  (Optional)
  The configuration setting for Airflow task logs. Structure is
  [documented below](#nested_task_logs_retention_config_c3).

<a name="nested_task_logs_retention_config_c3"></a>The `task_logs_retention_config` block supports:
* `storage_mode` - 
  (Optional)
  The mode of storage for Airflow task logs. Values for storage mode are 
  `CLOUD_LOGGING_ONLY` to only store logs in cloud logging and 
  `CLOUD_LOGGING_AND_CLOUD_STORAGE` to store logs in cloud logging and cloud storage.


<a name="nested_storage_config_c3"></a>The `storage_config` block supports:

* `bucket` -
  (Required)
  Name of an existing Cloud Storage bucket to be used by the environment.


<a name="nested_node_config_c3"></a>The `node_config` block supports:

* `network` -
  (Optional)
  The Compute Engine network to be used for machine
  communications, specified as a self-link, relative resource name
  (for example "projects/{project}/global/networks/{network}"), by name.

  The network must belong to the environment's project. If unspecified, the "default" network ID in the environment's
  project is used. If a Custom Subnet Network is provided, subnetwork must also be provided.

* `subnetwork` -
  (Optional)
  The Compute Engine subnetwork to be used for machine
  communications, specified as a self-link, relative resource name (for example,
  "projects/{project}/regions/{region}/subnetworks/{subnetwork}"), or by name. If subnetwork is provided,
  network must also be provided and the subnetwork must belong to the enclosing environment's project and region.

* `composer_network_attachment` -
  (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html), Cloud Composer 3 only)
  PSC (Private Service Connect) Network entry point. Customers can pre-create the Network Attachment 
  and point Cloud Composer environment to use. It is possible to share network attachment among many environments, 
  provided enough IP addresses are available.

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

* `composer_internal_ipv4_cidr_block` -
  (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html), Cloud Composer 3 only)
  /20 IPv4 cidr range that will be used by Composer internal components.
  Cannot be updated.

<a name="nested_software_config_c3"></a>The `software_config` block supports:

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

* `image_version` -
  (Required) If omitted, the default is the latest version of Composer 2.

  In Cloud Composer 3, you can only specify 3 in the Cloud Composer portion of the image version. Example: composer-3-airflow-x.y.z-build.t.

  The Apache Airflow portion of the image version is a full semantic version that points to one of the
  supported Apache Airflow versions, or an alias in the form of only major, major.minor or major.minor.patch versions specified.
  Like in Composer 1 and 2, a given Airflow version is released multiple times in Composer, with different patches
  and versions of dependencies. To distinguish between these versions in Composer 3, you can optionally specify a
  build number to pin to a specific Airflow release.
  Example: composer-3-airflow-2.6.3-build.4.

  The image version in Composer 3 must match the regular expression:
  `composer-(([0-9]+)(\.[0-9]+\.[0-9]+(-preview\.[0-9]+)?)?|latest)-airflow-(([0-9]+)((\.[0-9]+)(\.[0-9]+)?)?(-build\.[0-9]+)?)`
  Example: composer-3-airflow-2.6.3-build.4

  **Important**: In-place upgrade for Composer 3 is not yet supported.

* `cloud_data_lineage_integration` -
  (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html),
  Cloud Composer environments in versions composer-2.1.2-airflow-*.*.* and later)
  The configuration for Cloud Data Lineage integration. Structure is
  [documented below](#nested_cloud_data_lineage_integration_c3).

* `web_server_plugins_mode` -
  (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html), Cloud Composer 3 only)
  Web server plugins configuration. Can be either 'ENABLED' or 'DISABLED'. Defaults to 'ENABLED'.

<a name="nested_cloud_data_lineage_integration_c3"></a>The `cloud_data_lineage_integration` block supports:
* `enabled` -
  (Required)
  Whether or not Cloud Data Lineage integration is enabled.

<a name="nested_encryption_config_comp_2"></a>The `encryption_config` block supports:

* `kms_key_name` -
  (Required)
  Customer-managed Encryption Key available through Google's Key Management Service. It must
  be the fully qualified resource name,
  i.e. projects/project-id/locations/location/keyRings/keyring/cryptoKeys/key. Cannot be updated.

<a name="nested_maintenance_window_comp_2"></a>The `maintenance_window` block supports:

* `start_time` -
  (Required)
  Start time of the first recurrence of the maintenance window.

* `end_time` -
  (Required)
  Maintenance window end time. It is used only to calculate the duration of the maintenance window.
  The value for end-time must be in the future, relative to 'start_time'.

* `recurrence` -
  (Required)
  Maintenance window recurrence. Format is a subset of RFC-5545 (https://tools.ietf.org/html/rfc5545) 'RRULE'.
  The only allowed values for 'FREQ' field are 'FREQ=DAILY' and 'FREQ=WEEKLY;BYDAY=...'.
  Example values: 'FREQ=WEEKLY;BYDAY=TU,WE', 'FREQ=DAILY'.

The `workloads_config` block supports:

* `scheduler` -
  (Optional)
  Configuration for resources used by Airflow scheduler.

* `triggerer` -
  (Optional)
  Configuration for resources used by Airflow triggerer.

* `web_server` -
  (Optional)
  Configuration for resources used by Airflow web server.

* `worker` -
  (Optional)
  Configuration for resources used by Airflow workers.

* `dag_processor` -
  (Optional, [Beta](https://terraform.io/docs/providers/google/guides/provider_versions.html), Cloud Composer 3 only)
  Configuration for resources used by DAG processor.

The `scheduler` block supports:

* `cpu` -
  (Optional)
  The number of CPUs for a single Airflow scheduler.

* `memory_gb` -
  (Optional)
  The amount of memory (GB) for a single Airflow scheduler.

* `storage_gb` -
  (Optional)
  The amount of storage (GB) for a single Airflow scheduler.

* `count` -
  (Optional)
  The number of schedulers.

The `triggerer` block supports:

* `cpu` -
  (Required)
  The number of CPUs for a single Airflow triggerer.

* `memory_gb` -
  (Required)
  The amount of memory (GB) for a single Airflow triggerer.

* `count` -
  (Required)
  The number of Airflow triggerers.

The `web_server` block supports:

* `cpu` -
  (Optional)
  The number of CPUs for the Airflow web server.

* `memory_gb` -
  (Optional)
  The amount of memory (GB) for the Airflow web server.

* `storage_gb` -
  (Optional)
  The amount of storage (GB) for the Airflow web server.

The `worker` block supports:

* `cpu` -
  (Optional)
  The number of CPUs for a single Airflow worker.

* `memory_gb` -
  (Optional)
  The amount of memory (GB) for a single Airflow worker.

* `storage_gb`
  (Optional)
  The amount of storage (GB) for a single Airflow worker.

* `min_count` -
  (Optional)
  The minimum number of Airflow workers that the environment can run. The number of workers in the
  environment does not go above this number, even if a lower number of workers can handle the load.

* `max_count` -
  (Optional)
  The maximum number of Airflow workers that the environment can run. The number of workers in the
  environment does not go above this number, even if a higher number of workers is required to
  handle the load.

The `dag_processor` block supports:

* `cpu` -
  (Optional)
  CPU request and limit for DAG processor.

* `memory_gb` -
  (Optional)
  Memory (GB) request and limit for DAG processor.

* `storage_gb`
  (Optional)
  Storage (GB) request and limit for DAG processor.

* `count` -
  (Required)
  The number of Airflow DAG processors.

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
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options: configuration options:

- `create` - Default is 60 minutes.
- `update` - Default is 60 minutes.
- `delete` - Default is 6 minutes.

## Import

Environment can be imported using any of these accepted formats:

* `projects/{{project}}/locations/{{region}}/environments/{{name}}`
* `{{project}}/{{region}}/{{name}}`
* `{{name}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Environment using one of the formats above. For example:

```tf
import {
  id = "projects/{{project}}/locations/{{region}}/environments/{{name}}"
  to = google_composer_environment.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Environment can be imported using one of the formats above. For example:

```
$ terraform import google_composer_environment.default projects/{{project}}/locations/{{region}}/environments/{{name}}
$ terraform import google_composer_environment.default {{project}}/{{region}}/{{name}}
$ terraform import google_composer_environment.default {{name}}
```
