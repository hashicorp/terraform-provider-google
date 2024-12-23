// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package bigquery_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccBigQueryJob_bigqueryJobQueryExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryJob_bigqueryJobQueryExample(context),
			},
			{
				ResourceName:            "google_bigquery_job.job",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "labels", "status.0.state", "terraform_labels"},
			},
		},
	})
}

func testAccBigQueryJob_bigqueryJobQueryExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_table" "foo" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.bar.dataset_id
  table_id   = "tf_test_job_query%{random_suffix}_table"
}

resource "google_bigquery_dataset" "bar" {
  dataset_id                  = "tf_test_job_query%{random_suffix}_dataset"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
}

resource "google_bigquery_job" "job" {
  job_id     = "tf_test_job_query%{random_suffix}"

  labels = {
    "example-label" ="example-value"
  }

  query {
    query = "SELECT state FROM [lookerdata:cdc.project_tycho_reports]"

    destination_table {
      project_id = google_bigquery_table.foo.project
      dataset_id = google_bigquery_table.foo.dataset_id
      table_id   = google_bigquery_table.foo.table_id
    }

    allow_large_results = true
    flatten_results = true

    script_options {
      key_result_statement = "LAST"
    }
  }
}
`, context)
}

func TestAccBigQueryJob_bigqueryJobQueryTableReferenceExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryJob_bigqueryJobQueryTableReferenceExample(context),
			},
			{
				ResourceName:            "google_bigquery_job.job",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "labels", "query.0.default_dataset.0.dataset_id", "query.0.destination_table.0.table_id", "status.0.state", "terraform_labels"},
			},
		},
	})
}

func testAccBigQueryJob_bigqueryJobQueryTableReferenceExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_table" "foo" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.bar.dataset_id
  table_id   = "tf_test_job_query%{random_suffix}_table"
}

resource "google_bigquery_dataset" "bar" {
  dataset_id                  = "tf_test_job_query%{random_suffix}_dataset"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
}

resource "google_bigquery_job" "job" {
  job_id     = "tf_test_job_query%{random_suffix}"

  labels = {
    "example-label" ="example-value"
  }

  query {
    query = "SELECT state FROM [lookerdata:cdc.project_tycho_reports]"

    destination_table {
      table_id = google_bigquery_table.foo.id
    }

    default_dataset {
      dataset_id = google_bigquery_dataset.bar.id
    }

    allow_large_results = true
    flatten_results = true

    script_options {
      key_result_statement = "LAST"
    }
  }
}
`, context)
}

func TestAccBigQueryJob_bigqueryJobLoadExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryJob_bigqueryJobLoadExample(context),
			},
			{
				ResourceName:            "google_bigquery_job.job",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "labels", "status.0.state", "terraform_labels"},
			},
		},
	})
}

func testAccBigQueryJob_bigqueryJobLoadExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_table" "foo" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.bar.dataset_id
  table_id   = "tf_test_job_load%{random_suffix}_table"
}

resource "google_bigquery_dataset" "bar" {
  dataset_id                  = "tf_test_job_load%{random_suffix}_dataset"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
}

resource "google_bigquery_job" "job" {
  job_id     = "tf_test_job_load%{random_suffix}"

  labels = {
    "my_job" ="load"
  }

  load {
    source_uris = [
      "gs://cloud-samples-data/bigquery/us-states/us-states-by-date.csv",
    ]

    destination_table {
      project_id = google_bigquery_table.foo.project
      dataset_id = google_bigquery_table.foo.dataset_id
      table_id   = google_bigquery_table.foo.table_id
    }

    skip_leading_rows = 1
    schema_update_options = ["ALLOW_FIELD_RELAXATION", "ALLOW_FIELD_ADDITION"]

    write_disposition = "WRITE_APPEND"
    autodetect = true
  }
}
`, context)
}

func TestAccBigQueryJob_bigqueryJobLoadGeojsonExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryJob_bigqueryJobLoadGeojsonExample(context),
			},
			{
				ResourceName:            "google_bigquery_job.job",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "labels", "status.0.state", "terraform_labels"},
			},
		},
	})
}

func testAccBigQueryJob_bigqueryJobLoadGeojsonExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
locals {
  project = "%{project}" # Google Cloud Platform Project ID
}

resource "google_storage_bucket" "bucket" {
  name     = "${local.project}-tf-test-bq-geojson%{random_suffix}"  # Every bucket name must be globally unique
  location = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "object" {
  name   = "geojson-data.jsonl"
  bucket = google_storage_bucket.bucket.name
  content = <<EOF
{"type":"Feature","properties":{"continent":"Europe","region":"Scandinavia"},"geometry":{"type":"Polygon","coordinates":[[[-30.94,53.33],[33.05,53.33],[33.05,71.86],[-30.94,71.86],[-30.94,53.33]]]}}
{"type":"Feature","properties":{"continent":"Africa","region":"West Africa"},"geometry":{"type":"Polygon","coordinates":[[[-23.91,0],[11.95,0],[11.95,18.98],[-23.91,18.98],[-23.91,0]]]}}
EOF
}

resource "google_bigquery_table" "foo" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.bar.dataset_id
  table_id   = "tf_test_job_load%{random_suffix}_table"
}

resource "google_bigquery_dataset" "bar" {
  dataset_id                  = "tf_test_job_load%{random_suffix}_dataset"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
}

resource "google_bigquery_job" "job" {
  job_id     = "tf_test_job_load%{random_suffix}"

  labels = {
    "my_job" = "load"
  }

  load {
    source_uris = [
      "gs://${google_storage_bucket_object.object.bucket}/${google_storage_bucket_object.object.name}"
    ]

    destination_table {
      project_id = google_bigquery_table.foo.project
      dataset_id = google_bigquery_table.foo.dataset_id
      table_id   = google_bigquery_table.foo.table_id
    }

    write_disposition = "WRITE_TRUNCATE"
    autodetect = true
    source_format = "NEWLINE_DELIMITED_JSON"
    json_extension = "GEOJSON"
  }

  depends_on = ["google_storage_bucket_object.object"]
}
`, context)
}

func TestAccBigQueryJob_bigqueryJobLoadParquetExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryJob_bigqueryJobLoadParquetExample(context),
			},
			{
				ResourceName:            "google_bigquery_job.job",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "labels", "status.0.state", "terraform_labels"},
			},
		},
	})
}

func testAccBigQueryJob_bigqueryJobLoadParquetExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_storage_bucket" "test" {
  name                        = "tf_test_job_load%{random_suffix}_bucket"
  location                    = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "test" {
  name   =  "tf_test_job_load%{random_suffix}_bucket_object"
  source = "./test-fixtures/test.parquet.gzip"
  bucket = google_storage_bucket.test.name
}

resource "google_bigquery_dataset" "test" {
  dataset_id                  = "tf_test_job_load%{random_suffix}_dataset"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
}

resource "google_bigquery_table" "test" {
  deletion_protection = false
  table_id            = "tf_test_job_load%{random_suffix}_table"
  dataset_id          = google_bigquery_dataset.test.dataset_id
}

resource "google_bigquery_job" "job" {
  job_id = "tf_test_job_load%{random_suffix}"

  labels = {
    "my_job" ="load"
  }

  load {
    source_uris = [
      "gs://${google_storage_bucket_object.test.bucket}/${google_storage_bucket_object.test.name}"
    ]

    destination_table {
      project_id = google_bigquery_table.test.project
      dataset_id = google_bigquery_table.test.dataset_id
      table_id   = google_bigquery_table.test.table_id
    }

    schema_update_options = ["ALLOW_FIELD_RELAXATION", "ALLOW_FIELD_ADDITION"]
    write_disposition     = "WRITE_APPEND"
    source_format         = "PARQUET"
    autodetect            = true

    parquet_options {
      enum_as_string        = true
      enable_list_inference = true
    }
  }
}
`, context)
}

func TestAccBigQueryJob_bigqueryJobLoadTableReferenceExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryJob_bigqueryJobLoadTableReferenceExample(context),
			},
			{
				ResourceName:            "google_bigquery_job.job",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "labels", "load.0.destination_table.0.table_id", "status.0.state", "terraform_labels"},
			},
		},
	})
}

func testAccBigQueryJob_bigqueryJobLoadTableReferenceExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_table" "foo" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.bar.dataset_id
  table_id   = "tf_test_job_load%{random_suffix}_table"
}

resource "google_bigquery_dataset" "bar" {
  dataset_id                  = "tf_test_job_load%{random_suffix}_dataset"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
}

resource "google_bigquery_job" "job" {
  job_id     = "tf_test_job_load%{random_suffix}"

  labels = {
    "my_job" ="load"
  }

  load {
    source_uris = [
      "gs://cloud-samples-data/bigquery/us-states/us-states-by-date.csv",
    ]

    destination_table {
      table_id   = google_bigquery_table.foo.id
    }

    skip_leading_rows = 1
    schema_update_options = ["ALLOW_FIELD_RELAXATION", "ALLOW_FIELD_ADDITION"]

    write_disposition = "WRITE_APPEND"
    autodetect = true
  }
}
`, context)
}

func TestAccBigQueryJob_bigqueryJobCopyExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"kms_key_name":  acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "global", "tf-bootstrap-bigquery-job-key1").CryptoKey.Name,
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryJob_bigqueryJobCopyExample(context),
			},
			{
				ResourceName:            "google_bigquery_job.job",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "labels", "status.0.state", "terraform_labels"},
			},
		},
	})
}

func testAccBigQueryJob_bigqueryJobCopyExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
locals {
  count = 2
}

resource "google_bigquery_table" "source" {
  count = local.count

  dataset_id = google_bigquery_dataset.source[count.index].dataset_id
  table_id   = "tf_test_job_copy%{random_suffix}_${count.index}_table"

  deletion_protection = false

  schema = <<EOF
[
  {
    "name": "name",
    "type": "STRING",
    "mode": "NULLABLE"
  },
  {
    "name": "post_abbr",
    "type": "STRING",
    "mode": "NULLABLE"
  },
  {
    "name": "date",
    "type": "DATE",
    "mode": "NULLABLE"
  }
]
EOF
}

resource "google_bigquery_dataset" "source" {
  count = local.count

  dataset_id                  = "tf_test_job_copy%{random_suffix}_${count.index}_dataset"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
}

resource "google_bigquery_table" "dest" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.dest.dataset_id
  table_id   = "tf_test_job_copy%{random_suffix}_dest_table"

  schema = <<EOF
[
  {
    "name": "name",
    "type": "STRING",
    "mode": "NULLABLE"
  },
  {
    "name": "post_abbr",
    "type": "STRING",
    "mode": "NULLABLE"
  },
  {
    "name": "date",
    "type": "DATE",
    "mode": "NULLABLE"
  }
]
EOF

  encryption_configuration {
    kms_key_name = "%{kms_key_name}"
  }

  depends_on = ["google_kms_crypto_key_iam_member.encrypt_role"]
}

resource "google_bigquery_dataset" "dest" {
  dataset_id    = "tf_test_job_copy%{random_suffix}_dest_dataset"
  friendly_name = "test"
  description   = "This is a test description"
  location      = "US"
}

data "google_project" "project" {
  project_id = "%{project}"
}

resource "google_kms_crypto_key_iam_member" "encrypt_role" {
  crypto_key_id = "%{kms_key_name}"
  role = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_bigquery_job" "job" {
  job_id     = "tf_test_job_copy%{random_suffix}"

  copy {
    source_tables {
      project_id = google_bigquery_table.source.0.project
      dataset_id = google_bigquery_table.source.0.dataset_id
      table_id   = google_bigquery_table.source.0.table_id
    }

    source_tables {
      project_id = google_bigquery_table.source.1.project
      dataset_id = google_bigquery_table.source.1.dataset_id
      table_id   = google_bigquery_table.source.1.table_id
    }

    destination_table {
      project_id = google_bigquery_table.dest.project
      dataset_id = google_bigquery_table.dest.dataset_id
      table_id   = google_bigquery_table.dest.table_id
    }

    destination_encryption_configuration {
      kms_key_name = "%{kms_key_name}"
    }
  }

  depends_on = ["google_kms_crypto_key_iam_member.encrypt_role"]
}
`, context)
}

func TestAccBigQueryJob_bigqueryJobCopyTableReferenceExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"kms_key_name":  acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "global", "tf-bootstrap-bigquery-job-key2").CryptoKey.Name,
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryJob_bigqueryJobCopyTableReferenceExample(context),
			},
			{
				ResourceName:            "google_bigquery_job.job",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"copy.0.destination_table.0.table_id", "copy.0.source_tables.0.table_id", "copy.0.source_tables.1.table_id", "etag", "labels", "status.0.state", "terraform_labels"},
			},
		},
	})
}

func testAccBigQueryJob_bigqueryJobCopyTableReferenceExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
locals {
  count = 2
}

resource "google_bigquery_table" "source" {
  deletion_protection = false
  count = local.count

  dataset_id = google_bigquery_dataset.source[count.index].dataset_id
  table_id   = "tf_test_job_copy%{random_suffix}_${count.index}_table"

  schema = <<EOF
[
  {
    "name": "name",
    "type": "STRING",
    "mode": "NULLABLE"
  },
  {
    "name": "post_abbr",
    "type": "STRING",
    "mode": "NULLABLE"
  },
  {
    "name": "date",
    "type": "DATE",
    "mode": "NULLABLE"
  }
]
EOF

  depends_on = [google_bigquery_dataset.source]
}

resource "google_bigquery_dataset" "source" {
  count = local.count

  dataset_id                  = "tf_test_job_copy%{random_suffix}_${count.index}_dataset"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
}

resource "google_bigquery_table" "dest" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.dest.dataset_id
  table_id   = "tf_test_job_copy%{random_suffix}_dest_table"

  schema = <<EOF
[
  {
    "name": "name",
    "type": "STRING",
    "mode": "NULLABLE"
  },
  {
    "name": "post_abbr",
    "type": "STRING",
    "mode": "NULLABLE"
  },
  {
    "name": "date",
    "type": "DATE",
    "mode": "NULLABLE"
  }
]
EOF

  encryption_configuration {
    kms_key_name = "%{kms_key_name}"
  }

  depends_on = ["google_kms_crypto_key_iam_member.encrypt_role"]
}

resource "google_bigquery_dataset" "dest" {
  dataset_id    = "tf_test_job_copy%{random_suffix}_dest_dataset"
  friendly_name = "test"
  description   = "This is a test description"
  location      = "US"
}

data "google_project" "project" {
  project_id = "%{project}"
}

resource "google_kms_crypto_key_iam_member" "encrypt_role" {
  crypto_key_id = "%{kms_key_name}"
  role = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_bigquery_job" "job" {
  job_id     = "tf_test_job_copy%{random_suffix}"

  copy {
    source_tables {
      table_id   = google_bigquery_table.source.0.id
    }

    source_tables {
      table_id   = google_bigquery_table.source.1.id
    }

    destination_table {
      table_id   = google_bigquery_table.dest.id
    }

    destination_encryption_configuration {
      kms_key_name = "%{kms_key_name}"
    }
  }

  depends_on = ["google_kms_crypto_key_iam_member.encrypt_role"]
}
`, context)
}

func TestAccBigQueryJob_bigqueryJobExtractExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryJob_bigqueryJobExtractExample(context),
			},
			{
				ResourceName:            "google_bigquery_job.job",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "labels", "status.0.state", "terraform_labels"},
			},
		},
	})
}

func testAccBigQueryJob_bigqueryJobExtractExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_table" "source-one" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.source-one.dataset_id
  table_id   = "tf_test_job_extract%{random_suffix}_table"

  schema = <<EOF
[
  {
    "name": "name",
    "type": "STRING",
    "mode": "NULLABLE"
  },
  {
    "name": "post_abbr",
    "type": "STRING",
    "mode": "NULLABLE"
  },
  {
    "name": "date",
    "type": "DATE",
    "mode": "NULLABLE"
  }
]
EOF
}

resource "google_bigquery_dataset" "source-one" {
  dataset_id    = "tf_test_job_extract%{random_suffix}_dataset"
  friendly_name = "test"
  description   = "This is a test description"
  location      = "US"
}

resource "google_storage_bucket" "dest" {
  name          = "tf_test_job_extract%{random_suffix}_bucket"
  location      = "US"
  force_destroy = true
}

resource "google_bigquery_job" "job" {
  job_id     = "tf_test_job_extract%{random_suffix}"

  extract {
    destination_uris = ["${google_storage_bucket.dest.url}/extract"]

    source_table {
      project_id = google_bigquery_table.source-one.project
      dataset_id = google_bigquery_table.source-one.dataset_id
      table_id   = google_bigquery_table.source-one.table_id
    }

    destination_format = "NEWLINE_DELIMITED_JSON"
    compression = "GZIP"
  }
}
`, context)
}

func TestAccBigQueryJob_bigqueryJobExtractTableReferenceExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryJob_bigqueryJobExtractTableReferenceExample(context),
			},
			{
				ResourceName:            "google_bigquery_job.job",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"etag", "extract.0.source_table.0.table_id", "labels", "status.0.state", "terraform_labels"},
			},
		},
	})
}

func testAccBigQueryJob_bigqueryJobExtractTableReferenceExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_table" "source-one" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.source-one.dataset_id
  table_id   = "tf_test_job_extract%{random_suffix}_table"

  schema = <<EOF
[
  {
    "name": "name",
    "type": "STRING",
    "mode": "NULLABLE"
  },
  {
    "name": "post_abbr",
    "type": "STRING",
    "mode": "NULLABLE"
  },
  {
    "name": "date",
    "type": "DATE",
    "mode": "NULLABLE"
  }
]
EOF
}

resource "google_bigquery_dataset" "source-one" {
  dataset_id    = "tf_test_job_extract%{random_suffix}_dataset"
  friendly_name = "test"
  description   = "This is a test description"
  location      = "US"
}

resource "google_storage_bucket" "dest" {
  name          = "tf_test_job_extract%{random_suffix}_bucket"
  location      = "US"
  force_destroy = true
}

resource "google_bigquery_job" "job" {
  job_id     = "tf_test_job_extract%{random_suffix}"

  extract {
    destination_uris = ["${google_storage_bucket.dest.url}/extract"]

    source_table {
      table_id   = google_bigquery_table.source-one.id
    }

    destination_format = "NEWLINE_DELIMITED_JSON"
    compression = "GZIP"
  }
}
`, context)
}
