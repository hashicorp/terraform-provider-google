// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigquery_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"google.golang.org/api/bigquery/v2"
)

func TestAccBigQueryDataset_basic(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDataset(datasetID),
			},
			{
				ResourceName:      "google_bigquery_dataset.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigQueryDatasetUpdated(datasetID),
			},
			{
				ResourceName:      "google_bigquery_dataset.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigQueryDatasetUpdated2(datasetID),
			},
			{
				ResourceName:      "google_bigquery_dataset.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigQueryDataset_datasetWithContents(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	tableID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDatasetDeleteContents(datasetID),
				Check:  testAccAddTable(t, datasetID, tableID),
			},
			{
				ResourceName:            "google_bigquery_dataset.contents_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_contents_on_destroy"},
			},
		},
	})
}

func TestAccBigQueryDataset_access(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_access_%s", acctest.RandString(t, 10))
	otherDatasetID := fmt.Sprintf("tf_test_other_%s", acctest.RandString(t, 10))
	otherTableID := fmt.Sprintf("tf_test_other_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDatasetWithOneAccess(datasetID),
			},
			{
				ResourceName:      "google_bigquery_dataset.access_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigQueryDatasetWithTwoAccess(datasetID),
			},
			{
				ResourceName:      "google_bigquery_dataset.access_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigQueryDatasetWithOneAccess(datasetID),
			},
			{
				ResourceName:      "google_bigquery_dataset.access_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigQueryDatasetWithViewAccess(datasetID, otherDatasetID, otherTableID),
			},
			{
				ResourceName:      "google_bigquery_dataset.access_test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigQueryDataset_regionalLocation(t *testing.T) {
	t.Parallel()

	datasetID1 := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryRegionalDataset(datasetID1, "asia-south1"),
			},
			{
				ResourceName:      "google_bigquery_dataset.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigQueryDataset_cmek(t *testing.T) {
	t.Parallel()

	kms := acctest.BootstrapKMSKeyInLocation(t, "us")
	pid := envvar.GetTestProjectFromEnv()
	datasetID1 := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDataset_cmek(pid, datasetID1, kms.CryptoKey.Name),
			},
			{
				ResourceName:      "google_bigquery_dataset.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigQueryDataset_storageBillModel(t *testing.T) {
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDatasetStorageBillingModel(datasetID),
			},
			{
				ResourceName:      "google_bigquery_dataset.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAddTable(t *testing.T, datasetID string, tableID string) resource.TestCheckFunc {
	// Not actually a check, but adds a table independently of terraform
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		table := &bigquery.Table{
			TableReference: &bigquery.TableReference{
				DatasetId: datasetID,
				TableId:   tableID,
				ProjectId: config.Project,
			},
		}
		_, err := config.NewBigQueryClient(config.UserAgent).Tables.Insert(config.Project, datasetID, table).Do()
		if err != nil {
			return fmt.Errorf("Could not create table")
		}
		return nil
	}
}

func testAccBigQueryDataset(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id                      = "%s"
  friendly_name                   = "foo"
  description                     = "This is a foo description"
  location                        = "EU"
  default_partition_expiration_ms = 3600000
  default_table_expiration_ms     = 3600000

  labels = {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}
`, datasetID)
}

func testAccBigQueryDatasetUpdated(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id                      = "%s"
  friendly_name                   = "bar"
  description                     = "This is a bar description"
  location                        = "EU"
  default_partition_expiration_ms = 7200000
  default_table_expiration_ms     = 7200000

  labels = {
    env                         = "bar"
    default_table_expiration_ms = 7200000
  }
}
`, datasetID)
}

func testAccBigQueryDatasetUpdated2(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id                      = "%s"
  # friendly_name                   = "bar"
  description                     = "This is a bar description"
  location                        = "EU"
  default_partition_expiration_ms = 7200000
  default_table_expiration_ms     = 7200000

  labels = {
    env                         = "bar"
    default_table_expiration_ms = 7200000
  }
}
`, datasetID)
}

func testAccBigQueryDatasetDeleteContents(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "contents_test" {
  dataset_id                      = "%s"
  friendly_name                   = "foo"
  description                     = "This is a foo description"
  location                        = "EU"
  default_partition_expiration_ms = 3600000
  default_table_expiration_ms     = 3600000
  delete_contents_on_destroy      = true

  labels = {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}
`, datasetID)
}

func testAccBigQueryRegionalDataset(datasetID string, location string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id                  = "%s"
  friendly_name               = "foo"
  description                 = "This is a foo description"
  location                    = "%s"
  default_table_expiration_ms = 3600000

  labels = {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}
`, datasetID, location)
}

func testAccBigQueryDatasetWithOneAccess(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "access_test" {
  dataset_id = "%s"

  access {
    role          = "OWNER"
    user_by_email = "Joe@example.com"
  }

  labels = {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}
`, datasetID)
}

func testAccBigQueryDatasetWithTwoAccess(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "access_test" {
  dataset_id = "%s"

  access {
    role          = "OWNER"
    user_by_email = "Joe@example.com"
  }
  access {
    role   = "READER"
    domain = "hashicorp.com"
  }

  labels = {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}
`, datasetID)
}

func testAccBigQueryDatasetWithViewAccess(datasetID, otherDatasetID, otherTableID string) string {
	// Note that we have to add a non-view access to prevent BQ from creating 4 default
	// access entries.
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "other_dataset" {
  dataset_id = "%s"
}

resource "google_bigquery_table" "table_with_view" {
  deletion_protection = false
  table_id   = "%s"
  dataset_id = google_bigquery_dataset.other_dataset.dataset_id

  time_partitioning {
    type = "DAY"
  }

  view {
    query          = "SELECT state FROM [lookerdata:cdc.project_tycho_reports]"
    use_legacy_sql = true
  }
}

resource "google_bigquery_dataset" "access_test" {
  dataset_id = "%s"

  access {
    role          = "OWNER"
    user_by_email = "Joe@example.com"
  }
  access {
    view {
      project_id = google_bigquery_dataset.other_dataset.project
      dataset_id = google_bigquery_dataset.other_dataset.dataset_id
      table_id   = google_bigquery_table.table_with_view.table_id
    }
  }

  labels = {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}
`, otherDatasetID, otherTableID, datasetID)
}

func testAccBigQueryDataset_cmek(pid, datasetID, kmsKey string) string {
	return fmt.Sprintf(`
data "google_project" "project" {
  project_id = "%s"
}

resource "google_project_iam_member" "kms-project-binding" {
  project = data.google_project.project.project_id
  role    = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member  = "serviceAccount:bq-${data.google_project.project.number}@bigquery-encryption.iam.gserviceaccount.com"
}

resource "google_bigquery_dataset" "test" {
  dataset_id                  = "%s"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "US"
  default_table_expiration_ms = 3600000

  default_encryption_configuration {
    kms_key_name = "%s"
  }

  project = google_project_iam_member.kms-project-binding.project
}
`, pid, datasetID, kmsKey)
}

func testAccBigQueryDatasetStorageBillingModel(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id                      = "%s"
  friendly_name                   = "foo"
  description                     = "This is a foo description"
  location                        = "EU"
  default_partition_expiration_ms = 3600000
  default_table_expiration_ms     = 3600000
  storage_billing_model           = "PHYSICAL"

  labels = {
    env                         = "foo"
    default_table_expiration_ms = 3600000
  }
}
`, datasetID)
}
