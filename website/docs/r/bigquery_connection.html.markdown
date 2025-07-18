---
# ----------------------------------------------------------------------------
#
#     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
#
# ----------------------------------------------------------------------------
#
#     This code is generated by Magic Modules using the following:
#
#     Configuration: https:#github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/products/bigqueryconnection/Connection.yaml
#     Template:      https:#github.com/GoogleCloudPlatform/magic-modules/tree/main/mmv1/templates/terraform/resource.html.markdown.tmpl
#
#     DO NOT EDIT this file directly. Any changes made to this file will be
#     overwritten during the next generation cycle.
#
# ----------------------------------------------------------------------------
subcategory: "BigQuery Connection"
description: |-
  A connection allows BigQuery connections to external data sources.
---

# google_bigquery_connection

A connection allows BigQuery connections to external data sources..


To get more information about Connection, see:

* [API documentation](https://cloud.google.com/bigquery/docs/reference/bigqueryconnection/rest/v1/projects.locations.connections/create)
* How-to Guides
    * [Cloud SQL federated queries](https://cloud.google.com/bigquery/docs/cloud-sql-federated-queries)

~> **Warning:** All arguments including the following potentially sensitive
values will be stored in the raw state as plain text: `cloud_sql.credential.password`.
[Read more about sensitive data in state](https://www.terraform.io/language/state/sensitive-data).

<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=bigquery_connection_cloud_resource&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Bigquery Connection Cloud Resource


```hcl
resource "google_bigquery_connection" "connection" {
   connection_id = "my-connection"
   location      = "US"
   friendly_name = "👋"
   description   = "a riveting description"
   cloud_resource {}
}
```
<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=bigquery_connection_basic&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Bigquery Connection Basic


```hcl
resource "google_sql_database_instance" "instance" {
    name             = "my-database-instance"
    database_version = "POSTGRES_11"
    region           = "us-central1"
    settings {
		tier = "db-f1-micro"
	}

    deletion_protection  = true
}

resource "google_sql_database" "db" {
    instance = google_sql_database_instance.instance.name
    name     = "db"
}

resource "random_password" "pwd" {
    length = 16
    special = false
}

resource "google_sql_user" "user" {
    name = "user"
    instance = google_sql_database_instance.instance.name
    password = random_password.pwd.result
}

resource "google_bigquery_connection" "connection" {
    friendly_name = "👋"
    description   = "a riveting description"
    location      = "US"
    cloud_sql {
        instance_id = google_sql_database_instance.instance.connection_name
        database    = google_sql_database.db.name
        type        = "POSTGRES"
        credential {
          username = google_sql_user.user.name
          password = google_sql_user.user.password
        }
    }
}
```
<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=bigquery_connection_full&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Bigquery Connection Full


```hcl
resource "google_sql_database_instance" "instance" {
    name             = "my-database-instance"
    database_version = "POSTGRES_11"
    region           = "us-central1"
    settings {
		tier = "db-f1-micro"
	}

    deletion_protection  = true
}

resource "google_sql_database" "db" {
    instance = google_sql_database_instance.instance.name
    name     = "db"
}

resource "random_password" "pwd" {
    length = 16
    special = false
}

resource "google_sql_user" "user" {
    name = "user"
    instance = google_sql_database_instance.instance.name
    password = random_password.pwd.result
}

resource "google_bigquery_connection" "connection" {
    connection_id = "my-connection"
    location      = "US"
    friendly_name = "👋"
    description   = "a riveting description"
    cloud_sql {
        instance_id = google_sql_database_instance.instance.connection_name
        database    = google_sql_database.db.name
        type        = "POSTGRES"
        credential {
          username = google_sql_user.user.name
          password = google_sql_user.user.password
        }
    }
}
```
<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=bigquery_connection_aws&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Bigquery Connection Aws


```hcl
resource "google_bigquery_connection" "connection" {
   connection_id = "my-connection"
   location      = "aws-us-east-1"
   friendly_name = "👋"
   description   = "a riveting description"
   aws { 
      access_role {
         iam_role_id =  "arn:aws:iam::999999999999:role/omnirole"
      }
   }
}
```
<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=bigquery_connection_azure&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Bigquery Connection Azure


```hcl
resource "google_bigquery_connection" "connection" {
   connection_id = "my-connection"
   location      = "azure-eastus2"
   friendly_name = "👋"
   description   = "a riveting description"
   azure {
      customer_tenant_id = "customer-tenant-id"
      federated_application_client_id = "b43eeeee-eeee-eeee-eeee-a480155501ce"
   }
}
```
<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=bigquery_connection_cloudspanner&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Bigquery Connection Cloudspanner


```hcl
resource "google_bigquery_connection" "connection" {
   connection_id = "my-connection"
   location      = "US"
   friendly_name = "👋"
   description   = "a riveting description"
   cloud_spanner { 
      database = "projects/project/instances/instance/databases/database"
      database_role = "database_role"
   }
}
```
<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=bigquery_connection_cloudspanner_databoost&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Bigquery Connection Cloudspanner Databoost


```hcl
resource "google_bigquery_connection" "connection" {
   connection_id = "my-connection"
   location      = "US"
   friendly_name = "👋"
   description   = "a riveting description"
   cloud_spanner { 
      database        = "projects/project/instances/instance/databases/database"
      use_parallelism = true
      use_data_boost  = true
      max_parallelism = 100
   }
}
```
<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=bigquery_connection_spark&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Bigquery Connection Spark


```hcl
resource "google_bigquery_connection" "connection" {
   connection_id = "my-connection"
   location      = "US"
   friendly_name = "👋"
   description   = "a riveting description"
   spark {
      spark_history_server_config {
         dataproc_cluster = google_dataproc_cluster.basic.id
      }
   }
}

resource "google_dataproc_cluster" "basic" {
   name   = "my-connection"
   region = "us-central1"

   cluster_config {
     # Keep the costs down with smallest config we can get away with
     software_config {
       override_properties = {
         "dataproc:dataproc.allow.zero.workers" = "true"
       }
     }
 
     master_config {
       num_instances = 1
       machine_type  = "e2-standard-2"
       disk_config {
         boot_disk_size_gb = 35
       }
     }
   }   
 }
```
<div class = "oics-button" style="float: right; margin: 0 0 -15px">
  <a href="https://console.cloud.google.com/cloudshell/open?cloudshell_git_repo=https%3A%2F%2Fgithub.com%2Fterraform-google-modules%2Fdocs-examples.git&cloudshell_image=gcr.io%2Fcloudshell-images%2Fcloudshell%3Alatest&cloudshell_print=.%2Fmotd&cloudshell_tutorial=.%2Ftutorial.md&cloudshell_working_dir=bigquery_connection_sql_with_cmek&open_in_editor=main.tf" target="_blank">
    <img alt="Open in Cloud Shell" src="//gstatic.com/cloudssh/images/open-btn.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>
## Example Usage - Bigquery Connection Sql With Cmek


```hcl
resource "google_sql_database_instance" "instance" {
  name             = "my-database-instance"
  region           = "us-central1"

  database_version = "POSTGRES_11"
  settings {
    tier = "db-f1-micro"
  }

  deletion_protection  = true
}

resource "google_sql_database" "db" {
  instance = google_sql_database_instance.instance.name
  name     = "db"
}

resource "google_sql_user" "user" {
  name = "user"
  instance = google_sql_database_instance.instance.name
  password = "tf-test-my-password%{random_suffix}"
}

resource "google_bigquery_connection" "bq-connection-cmek" {
  friendly_name = "👋"
  description   = "a riveting description"
  location      = "US"
  kms_key_name  = "projects/project/locations/us-central1/keyRings/us-central1/cryptoKeys/bq-key"
  cloud_sql {
    instance_id = google_sql_database_instance.instance.connection_name
    database    = google_sql_database.db.name
    type        = "POSTGRES"
    credential {
      username = google_sql_user.user.name
      password = google_sql_user.user.password
    }
  }
}
```

## Argument Reference

The following arguments are supported:



* `connection_id` -
  (Optional)
  Optional connection id that should be assigned to the created connection.

* `location` -
  (Optional)
  The geographic location where the connection should reside.
  Cloud SQL instance must be in the same location as the connection
  with following exceptions: Cloud SQL us-central1 maps to BigQuery US, Cloud SQL europe-west1 maps to BigQuery EU.
  Examples: US, EU, asia-northeast1, us-central1, europe-west1.
  Spanner Connections same as spanner region
  AWS allowed regions are aws-us-east-1
  Azure allowed regions are azure-eastus2

* `friendly_name` -
  (Optional)
  A descriptive name for the connection

* `description` -
  (Optional)
  A descriptive description for the connection

* `kms_key_name` -
  (Optional)
  Optional. The Cloud KMS key that is used for encryption.
  Example: projects/[kms_project_id]/locations/[region]/keyRings/[key_region]/cryptoKeys/[key]

* `cloud_sql` -
  (Optional)
  Connection properties specific to the Cloud SQL.
  Structure is [documented below](#nested_cloud_sql).

* `aws` -
  (Optional)
  Connection properties specific to Amazon Web Services.
  Structure is [documented below](#nested_aws).

* `azure` -
  (Optional)
  Container for connection properties specific to Azure.
  Structure is [documented below](#nested_azure).

* `cloud_spanner` -
  (Optional)
  Connection properties specific to Cloud Spanner
  Structure is [documented below](#nested_cloud_spanner).

* `cloud_resource` -
  (Optional)
  Container for connection properties for delegation of access to GCP resources.
  Structure is [documented below](#nested_cloud_resource).

* `spark` -
  (Optional)
  Container for connection properties to execute stored procedures for Apache Spark. resources.
  Structure is [documented below](#nested_spark).

* `project` - (Optional) The ID of the project in which the resource belongs.
    If it is not provided, the provider project is used.



<a name="nested_cloud_sql"></a>The `cloud_sql` block supports:

* `instance_id` -
  (Required)
  Cloud SQL instance ID in the form project:location:instance.

* `database` -
  (Required)
  Database name.

* `credential` -
  (Required)
  Cloud SQL properties.
  Structure is [documented below](#nested_cloud_sql_credential).

* `type` -
  (Required)
  Type of the Cloud SQL database.
  Possible values are: `DATABASE_TYPE_UNSPECIFIED`, `POSTGRES`, `MYSQL`.

* `service_account_id` -
  (Output)
  When the connection is used in the context of an operation in BigQuery, this service account will serve as the identity being used for connecting to the CloudSQL instance specified in this connection.


<a name="nested_cloud_sql_credential"></a>The `credential` block supports:

* `username` -
  (Required)
  Username for database.

* `password` -
  (Required)
  Password for database.
  **Note**: This property is sensitive and will not be displayed in the plan.

<a name="nested_aws"></a>The `aws` block supports:

* `access_role` -
  (Required)
  Authentication using Google owned service account to assume into customer's AWS IAM Role.
  Structure is [documented below](#nested_aws_access_role).


<a name="nested_aws_access_role"></a>The `access_role` block supports:

* `iam_role_id` -
  (Required)
  The user’s AWS IAM Role that trusts the Google-owned AWS IAM user Connection.

* `identity` -
  (Output)
  A unique Google-owned and Google-generated identity for the Connection. This identity will be used to access the user's AWS IAM Role.

<a name="nested_azure"></a>The `azure` block supports:

* `application` -
  (Output)
  The name of the Azure Active Directory Application.

* `client_id` -
  (Output)
  The client id of the Azure Active Directory Application.

* `object_id` -
  (Output)
  The object id of the Azure Active Directory Application.

* `customer_tenant_id` -
  (Required)
  The id of customer's directory that host the data.

* `federated_application_client_id` -
  (Optional)
  The Azure Application (client) ID where the federated credentials will be hosted.

* `redirect_uri` -
  (Output)
  The URL user will be redirected to after granting consent during connection setup.

* `identity` -
  (Output)
  A unique Google-owned and Google-generated identity for the Connection. This identity will be used to access the user's Azure Active Directory Application.

<a name="nested_cloud_spanner"></a>The `cloud_spanner` block supports:

* `database` -
  (Required)
  Cloud Spanner database in the form `project/instance/database'.

* `use_parallelism` -
  (Optional)
  If parallelism should be used when reading from Cloud Spanner.

* `max_parallelism` -
  (Optional)
  Allows setting max parallelism per query when executing on Spanner independent compute resources. If unspecified, default values of parallelism are chosen that are dependent on the Cloud Spanner instance configuration. `useParallelism` and `useDataBoost` must be set when setting max parallelism.

* `use_data_boost` -
  (Optional)
  If set, the request will be executed via Spanner independent compute resources. `use_parallelism` must be set when using data boost.

* `database_role` -
  (Optional)
  Cloud Spanner database role for fine-grained access control. The Cloud Spanner admin should have provisioned the database role with appropriate permissions, such as `SELECT` and `INSERT`. Other users should only use roles provided by their Cloud Spanner admins. The database role name must start with a letter, and can only contain letters, numbers, and underscores. For more details, see https://cloud.google.com/spanner/docs/fgac-about.

* `use_serverless_analytics` -
  (Optional, Deprecated)
  If the serverless analytics service should be used to read data from Cloud Spanner. `useParallelism` must be set when using serverless analytics.

  ~> **Warning:** `useServerlessAnalytics` is deprecated and will be removed in a future major release. Use `useDataBoost` instead.

<a name="nested_cloud_resource"></a>The `cloud_resource` block supports:

* `service_account_id` -
  (Output)
  The account ID of the service created for the purpose of this connection.

<a name="nested_spark"></a>The `spark` block supports:

* `service_account_id` -
  (Output)
  The account ID of the service created for the purpose of this connection.

* `metastore_service_config` -
  (Optional)
  Dataproc Metastore Service configuration for the connection.
  Structure is [documented below](#nested_spark_metastore_service_config).

* `spark_history_server_config` -
  (Optional)
  Spark History Server configuration for the connection.
  Structure is [documented below](#nested_spark_spark_history_server_config).


<a name="nested_spark_metastore_service_config"></a>The `metastore_service_config` block supports:

* `metastore_service` -
  (Optional)
  Resource name of an existing Dataproc Metastore service in the form of projects/[projectId]/locations/[region]/services/[serviceId].

<a name="nested_spark_spark_history_server_config"></a>The `spark_history_server_config` block supports:

* `dataproc_cluster` -
  (Optional)
  Resource name of an existing Dataproc Cluster to act as a Spark History Server for the connection if the form of projects/[projectId]/regions/[region]/clusters/[cluster_name].

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `id` - an identifier for the resource with format `projects/{{project}}/locations/{{location}}/connections/{{connection_id}}`

* `name` -
  The resource name of the connection in the form of:
  "projects/{project_id}/locations/{location_id}/connections/{connectionId}"

* `has_credential` -
  True if the connection has credential assigned.


## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options:

- `create` - Default is 20 minutes.
- `update` - Default is 20 minutes.
- `delete` - Default is 20 minutes.

## Import


Connection can be imported using any of these accepted formats:

* `projects/{{project}}/locations/{{location}}/connections/{{connection_id}}`
* `{{project}}/{{location}}/{{connection_id}}`
* `{{location}}/{{connection_id}}`


In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Connection using one of the formats above. For example:

```tf
import {
  id = "projects/{{project}}/locations/{{location}}/connections/{{connection_id}}"
  to = google_bigquery_connection.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Connection can be imported using one of the formats above. For example:

```
$ terraform import google_bigquery_connection.default projects/{{project}}/locations/{{location}}/connections/{{connection_id}}
$ terraform import google_bigquery_connection.default {{project}}/{{location}}/{{connection_id}}
$ terraform import google_bigquery_connection.default {{location}}/{{connection_id}}
```

## User Project Overrides

This resource supports [User Project Overrides](https://registry.terraform.io/providers/hashicorp/google/latest/docs/guides/provider_reference#user_project_override).
