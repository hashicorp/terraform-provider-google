// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigquery_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
				Config: testAccBigQueryDataset_withoutLabels(datasetID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_bigquery_dataset.test", "labels.%"),
					resource.TestCheckNoResourceAttr("google_bigquery_dataset.test", "effective_labels.%"),
				),
			},
			{
				ResourceName:      "google_bigquery_dataset.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigQueryDataset(datasetID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.env", "foo"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.default_table_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.%", "2"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.env", "foo"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.default_table_expiration_ms", "3600000"),
				),
			},
			{
				ResourceName:      "google_bigquery_dataset.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The labels field in the state is decided by the configuration.
				// During importing, the configuration is unavailable, so the labels field in the state after importing is empty.
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDatasetUpdated(datasetID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.env", "bar"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.default_table_expiration_ms", "7200000"),

					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.%", "2"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.env", "bar"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.default_table_expiration_ms", "7200000"),
				),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDatasetUpdated2(datasetID),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDataset_withoutLabels(datasetID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_bigquery_dataset.test", "labels.%"),
					resource.TestCheckNoResourceAttr("google_bigquery_dataset.test", "effective_labels.%"),
				),
			},
			{
				ResourceName:      "google_bigquery_dataset.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigQueryDataset_withComputedLabels(t *testing.T) {
	// Skip it in VCR test because of the randomness of uuid in "labels" field
	// which causes the replaying mode after recording mode failing in VCR test
	acctest.SkipIfVcr(t)
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDataset(datasetID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.env", "foo"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.default_table_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.%", "3"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.env", "foo"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.default_table_expiration_ms", "3600000"),
				),
			},
			{
				ResourceName:      "google_bigquery_dataset.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The labels field in the state is decided by the configuration.
				// During importing, the configuration is unavailable, so the labels field in the state after importing is empty.
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDatasetUpdated_withComputedLabels(datasetID),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccBigQueryDataset_withProvider5(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))
	oldVersion := map[string]resource.ExternalProvider{
		"google": {
			VersionConstraint: "4.75.0", // a version that doesn't separate user defined labels and system labels
			Source:            "registry.terraform.io/hashicorp/google",
		},
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:            testAccBigQueryDataset_withoutLabelsV4(datasetID),
				ExternalProviders: oldVersion,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("google_bigquery_dataset.test", "labels.%"),
					resource.TestCheckNoResourceAttr("google_bigquery_dataset.test", "effective_labels.%"),
				),
			},
			{
				Config:                   testAccBigQueryDataset(datasetID),
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.%", "2"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.env", "foo"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "labels.default_table_expiration_ms", "3600000"),

					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.%", "2"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.env", "foo"),
					resource.TestCheckResourceAttr("google_bigquery_dataset.test", "effective_labels.default_table_expiration_ms", "3600000"),
				),
			},
		},
	})
}

func TestAccBigQueryDataset_withOutOfBandLabels(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	datasetID := fmt.Sprintf("tf_test_%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDataset(datasetID),
				Check:  addOutOfBandLabels(t, datasetID),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_contents_on_destroy", "labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDatasetUpdated(datasetID),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_contents_on_destroy", "labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDatasetUpdated_withOutOfBandLabels(datasetID),
			},
			{
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_contents_on_destroy", "labels", "terraform_labels"},
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
				ImportStateVerifyIgnore: []string{"delete_contents_on_destroy", "labels", "terraform_labels"},
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
				ResourceName:            "google_bigquery_dataset.access_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDatasetWithThreeAccess(datasetID),
			},
			{
				ResourceName:            "google_bigquery_dataset.access_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDatasetWithOneAccess(datasetID),
			},
			{
				ResourceName:            "google_bigquery_dataset.access_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDatasetWithViewAccess(datasetID, otherDatasetID, otherTableID),
			},
			{
				ResourceName:            "google_bigquery_dataset.access_test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
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
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
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
				ResourceName:            "google_bigquery_dataset.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccBigQueryDataset_invalidCharacterInID(t *testing.T) {
	t.Parallel()
	// Not an acceptance test.
	acctest.SkipIfVcr(t)

	datasetID := fmt.Sprintf("tf_test_%s-with-hyphens", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccBigQueryDataset(datasetID),
				ExpectError: regexp.MustCompile("must contain only letters.+numbers.+or underscores.+"),
			},
		},
	})
}

func TestAccBigQueryDataset_invalidLongID(t *testing.T) {
	t.Parallel()
	// Not an acceptance test.
	acctest.SkipIfVcr(t)

	datasetSuffix := acctest.RandString(t, 10)
	datasetID := fmt.Sprintf("tf_test_%s", strings.Repeat(datasetSuffix, 200))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccBigQueryDataset(datasetID),
				ExpectError: regexp.MustCompile(".+cannot be greater than 1,024 characters"),
			},
		},
	})
}

func TestAccBigQueryDataset_bigqueryDatasetResourceTags_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBigQueryDatasetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigQueryDataset_bigqueryDatasetResourceTags_basic(context),
			},
			{
				ResourceName:            "google_bigquery_dataset.dataset",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccBigQueryDataset_bigqueryDatasetResourceTags_update(context),
			},
			{
				ResourceName:            "google_bigquery_dataset.dataset",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
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

func addOutOfBandLabels(t *testing.T, datasetID string) resource.TestCheckFunc {
	// Not actually a check, but adds labels independently of terraform
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		dataset, err := config.NewBigQueryClient(config.UserAgent).Datasets.Get(config.Project, datasetID).Do()
		if err != nil {
			return fmt.Errorf("Could not get dataset with ID %s", datasetID)
		}

		dataset.Labels["outband_key"] = "test"
		_, err = config.NewBigQueryClient(config.UserAgent).Datasets.Patch(config.Project, datasetID, dataset).Do()
		if err != nil {
			return fmt.Errorf("Could not update labele for the dataset")
		}
		return nil
	}
}

func testAccBigQueryDataset_withoutLabels(datasetID string) string {
	return fmt.Sprintf(`
provider "google" {
  add_terraform_attribution_label = false
}

resource "google_bigquery_dataset" "test" {
  dataset_id                      = "%s"
  friendly_name                   = "foo"
  description                     = "This is a foo description"
  location                        = "EU"
  default_partition_expiration_ms = 3600000
  default_table_expiration_ms     = 3600000
}
`, datasetID)
}

func testAccBigQueryDataset_withoutLabelsV4(datasetID string) string {
	return fmt.Sprintf(`
resource "google_bigquery_dataset" "test" {
  dataset_id                      = "%s"
  friendly_name                   = "foo"
  description                     = "This is a foo description"
  location                        = "EU"
  default_partition_expiration_ms = 3600000
  default_table_expiration_ms     = 3600000
}
`, datasetID)
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

func testAccBigQueryDatasetUpdated_withOutOfBandLabels(datasetID string) string {
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
	outband_key                 = "test-update"
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

func testAccBigQueryDatasetUpdated_withComputedLabels(datasetID string) string {
	return fmt.Sprintf(`
resource "random_uuid" "test" {
}

resource "google_bigquery_dataset" "test" {
  dataset_id                      = "%s"
  # friendly_name                   = "bar"
  description                     = "This is a bar description"
  location                        = "EU"
  default_partition_expiration_ms = 7200000
  default_table_expiration_ms     = 7200000

  labels = {
    env                         = "${random_uuid.test.result}"
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

func testAccBigQueryDatasetWithThreeAccess(datasetID string) string {
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
  access {
    role       = "READER"
    iam_member = "allUsers"
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

resource "google_kms_crypto_key_iam_member" "kms-member" {
  crypto_key_id = "%s"
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

  depends_on = [google_kms_crypto_key_iam_member.kms-member]
}
`, pid, kmsKey, datasetID, kmsKey)
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

func testAccBigQueryDataset_bigqueryDatasetResourceTags_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_tags_tag_key" "tag_key1" {
  parent     = data.google_project.project.id
  short_name = "tf_test_tag_key1%{random_suffix}"
}

resource "google_tags_tag_value" "tag_value1" {
  parent = google_tags_tag_key.tag_key1.id
  short_name = "tf_test_tag_value1%{random_suffix}"
}

resource "google_tags_tag_key" "tag_key2" {
  parent     = data.google_project.project.id
  short_name = "tf_test_tag_key2%{random_suffix}"
}

resource "google_tags_tag_value" "tag_value2" {
  parent     = google_tags_tag_key.tag_key2.id
  short_name = "tf_test_tag_value2%{random_suffix}"
}

resource "google_bigquery_dataset" "dataset" {
  dataset_id                  = "dataset%{random_suffix}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "EU"

  resource_tags = {
    (google_tags_tag_key.tag_key1.namespaced_name) = google_tags_tag_value.tag_value1.short_name
    (google_tags_tag_key.tag_key2.namespaced_name) = google_tags_tag_value.tag_value2.short_name
  }
}
`, context)
}

func testAccBigQueryDataset_bigqueryDatasetResourceTags_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
}

resource "google_tags_tag_key" "tag_key1" {
  parent     = data.google_project.project.id
  short_name = "tf_test_tag_key1%{random_suffix}"
}

resource "google_tags_tag_value" "tag_value1" {
  parent     = google_tags_tag_key.tag_key1.id
  short_name = "tf_test_tag_value1%{random_suffix}"
}

resource "google_tags_tag_key" "tag_key2" {
  parent     = data.google_project.project.id
  short_name = "tf_test_tag_key2%{random_suffix}"
}

resource "google_tags_tag_value" "tag_value2" {
  parent     = google_tags_tag_key.tag_key2.id
  short_name = "tf_test_tag_value2%{random_suffix}"
}

resource "google_bigquery_dataset" "dataset" {
  dataset_id                  = "dataset%{random_suffix}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "EU"

  resource_tags = {
  }
}
`, context)
}

func testAccBigQueryDataset_externalCatalogDatasetOptions_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_dataset" "dataset" {
  provider = google-beta

  dataset_id    = "dataset%{random_suffix}"
  friendly_name = "test"
  description   = "This is a test description"
  location      = "US"

  external_catalog_dataset_options {
    parameters = {
      "dataset_owner" = "dataset_owner"
    }
    default_storage_location_uri = "gs://test_dataset/tables"
  }
}
`, context)
}

func testAccBigQueryDataset_externalCatalogDatasetOptions_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_dataset" "dataset" {
  provider = google-beta

  dataset_id    = "dataset%{random_suffix}"
  friendly_name = "test"
  description   = "This is a test description"
  location      = "US"

  external_catalog_dataset_options {
    parameters = {
      "new_dataset_owner" = "new_dataset_owner"
    }
    default_storage_location_uri = "gs://new_test_dataset/new_tables"
  }
}
`, context)
}
