// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccLoggingProjectSink_basic(t *testing.T) {
	t.Parallel()

	orgId := envvar.GetTestOrgFromEnv(t)
	billingAccount := envvar.GetTestBillingAccountFromEnv(t)
	projectId := "tf-test" + acctest.RandString(t, 10)
	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingProjectSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectSink_basic(projectId, orgId, billingAccount, sinkName, bucketName, "false"),
			},
			{
				ResourceName:      "google_logging_project_sink.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingProjectSink_default(t *testing.T) {
	t.Parallel()

	orgId := envvar.GetTestOrgFromEnv(t)
	billingAccount := envvar.GetTestBillingAccountFromEnv(t)
	projectId := "tf-test" + acctest.RandString(t, 10)
	sinkName := "_Default"
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Default sink has a permadiff if any value is sent for "disabled" other than "true"
				Config: testAccLoggingProjectSink_basic(projectId, orgId, billingAccount, sinkName, bucketName, "true"),
			},
			{
				ResourceName:      "google_logging_project_sink.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingProjectSink_described(t *testing.T) {
	t.Parallel()

	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingProjectSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectSink_described(sinkName, envvar.GetTestProjectFromEnv(), bucketName),
			},
			{
				ResourceName:      "google_logging_project_sink.described",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingProjectSink_described_update(t *testing.T) {
	t.Parallel()

	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingProjectSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectSink_described(sinkName, envvar.GetTestProjectFromEnv(), bucketName),
			},
			{
				ResourceName:      "google_logging_project_sink.described",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingProjectSink_described_update(sinkName, envvar.GetTestProjectFromEnv(), bucketName),
			},
			{
				ResourceName:      "google_logging_project_sink.described",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingProjectSink_disabled(t *testing.T) {
	t.Parallel()

	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingProjectSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectSink_disabled(sinkName, envvar.GetTestProjectFromEnv(), bucketName),
			},
			{
				ResourceName:      "google_logging_project_sink.disabled",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingProjectSink_updatePreservesUniqueWriter(t *testing.T) {
	t.Parallel()

	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)
	updatedBucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingProjectSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectSink_uniqueWriter(sinkName, bucketName),
			},
			{
				ResourceName:      "google_logging_project_sink.unique_writer",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingProjectSink_uniqueWriterUpdated(sinkName, updatedBucketName),
			},
			{
				ResourceName:      "google_logging_project_sink.unique_writer",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingProjectSink_updatePreservesCustomWriter(t *testing.T) {
	t.Parallel()

	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	account := "tf-test-sink-sa" + acctest.RandString(t, 10)
	accountUpdated := "tf-test-sink-sa" + acctest.RandString(t, 10)
	testProject := envvar.GetTestProjectFromEnv()

	// custom_writer_identity is write-only, and writer_dietity is an output only field
	// verify that the value of writer_identity matches the expected custom_writer_identity.
	expectedWriterIdentity := fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, testProject)
	expectedUpdatedWriterIdentity := fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", accountUpdated, testProject)

	org := envvar.GetTestOrgFromEnv(t)
	billingId := envvar.GetTestBillingAccountFromEnv(t)
	project := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckLoggingProjectSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectSink_customWriter(org, billingId, project, sinkName, account),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_logging_project_sink.custom_writer", "writer_identity", expectedWriterIdentity),
				),
			},
			{
				ResourceName:      "google_logging_project_sink.custom_writer",
				ImportState:       true,
				ImportStateVerify: true,
				// Logging sink create API doesn't return this field in response
				ImportStateVerifyIgnore: []string{"custom_writer_identity"},
			},
			{
				Config: testAccLoggingProjectSink_customWriterUpdated(org, billingId, project, sinkName, accountUpdated),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_logging_project_sink.custom_writer", "writer_identity", expectedUpdatedWriterIdentity),
				),
			},
			{
				ResourceName:            "google_logging_project_sink.custom_writer",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_writer_identity"},
			},
		},
	})
}

func TestAccLoggingProjectSink_updateBigquerySink(t *testing.T) {
	t.Parallel()

	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bqDatasetID := "tf_test_sink_" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingProjectSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectSink_bigquery_before(sinkName, bqDatasetID),
			},
			{
				ResourceName:      "google_logging_project_sink.bigquery",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingProjectSink_bigquery_after(sinkName, bqDatasetID),
			},
			{
				ResourceName:      "google_logging_project_sink.bigquery",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingProjectSink_heredoc(t *testing.T) {
	t.Parallel()

	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingProjectSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectSink_heredoc(sinkName, envvar.GetTestProjectFromEnv(), bucketName),
			},
			{
				ResourceName:      "google_logging_project_sink.heredoc",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingProjectSink_loggingbucket(t *testing.T) {
	t.Parallel()

	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingProjectSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectSink_loggingbucket(sinkName, envvar.GetTestProjectFromEnv()),
			},
			{
				ResourceName:      "google_logging_project_sink.loggingbucket",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingProjectSink_disabled_update(t *testing.T) {
	t.Parallel()

	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingProjectSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectSink_disabled_update(sinkName, envvar.GetTestProjectFromEnv(), bucketName, "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_logging_project_sink.disabled", "disabled", "true"),
				),
			},
			{
				ResourceName:      "google_logging_project_sink.disabled",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingProjectSink_disabled_update(sinkName, envvar.GetTestProjectFromEnv(), bucketName, "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_logging_project_sink.disabled", "disabled", "false"),
				),
			},
			{
				ResourceName:      "google_logging_project_sink.disabled",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingProjectSink_disabled_update(sinkName, envvar.GetTestProjectFromEnv(), bucketName, "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_logging_project_sink.disabled", "disabled", "true"),
				),
			},
			{
				ResourceName:      "google_logging_project_sink.disabled",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLoggingProjectSinkDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_logging_project_sink" {
				continue
			}

			attributes := rs.Primary.Attributes

			_, err := config.NewLoggingClient(config.UserAgent).Projects.Sinks.Get(attributes["id"]).Do()
			if err == nil {
				return fmt.Errorf("project sink still exists")
			}
		}

		return nil
	}
}

func testAccLoggingProjectSink_basic(projectId, orgId, billingAccount, sinkName, bucketName, disabled string) string {
	return fmt.Sprintf(`
resource "google_project" "project" {
	project_id = "%s"
	name       = "%s"
	org_id     = "%s"
  billing_account = "%s"
}

resource "google_project_service" "logging_service" {
	project = google_project.project.project_id
	service = "logging.googleapis.com"
}

resource "google_logging_project_sink" "basic" {
  name        = "%s"
  disabled    = %s
  project     = google_project_service.logging_service.project
  destination = "storage.googleapis.com/${google_storage_bucket.gcs-bucket.name}"
  filter      = "logName=\"projects/${google_project.project.project_id}/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}

resource "google_storage_bucket" "gcs-bucket" {
  name     = "%s"
  project  = google_project.project.project_id
  location = "US"
}
`, projectId, projectId, orgId, billingAccount, sinkName, disabled, bucketName)
}

func testAccLoggingProjectSink_described(name, project, bucketName string) string {
	return fmt.Sprintf(`
resource "google_logging_project_sink" "described" {
  name        = "%s"
  project     = "%s"
  destination = "storage.googleapis.com/${google_storage_bucket.gcs-bucket.name}"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
  description = "this is a description for a project level logging sink"
  
  unique_writer_identity = false
}

resource "google_storage_bucket" "gcs-bucket" {
  name     = "%s"
  location = "US"
}
`, name, project, project, bucketName)
}

func testAccLoggingProjectSink_described_update(name, project, bucketName string) string {
	return fmt.Sprintf(`
resource "google_logging_project_sink" "described" {
  name        = "%s"
  project     = "%s"
  destination = "storage.googleapis.com/${google_storage_bucket.gcs-bucket.name}"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
  description = "description updated"
  
  unique_writer_identity = true
}

resource "google_storage_bucket" "gcs-bucket" {
  name     = "%s"
  location = "US"
}
`, name, project, project, bucketName)
}

func testAccLoggingProjectSink_disabled(name, project, bucketName string) string {
	return fmt.Sprintf(`
resource "google_logging_project_sink" "disabled" {
  name        = "%s"
  project     = "%s"
  destination = "storage.googleapis.com/${google_storage_bucket.gcs-bucket.name}"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
  disabled    = true

  unique_writer_identity = false
}

resource "google_storage_bucket" "gcs-bucket" {
  name     = "%s"
  location = "US"
}
`, name, project, project, bucketName)
}

func testAccLoggingProjectSink_disabled_update(name, project, bucketName, disabled string) string {
	return fmt.Sprintf(`
resource "google_logging_project_sink" "disabled" {
  name        = "%s"
  project     = "%s"
  destination = "storage.googleapis.com/${google_storage_bucket.gcs-bucket.name}"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
  disabled    = "%s"

  unique_writer_identity = true
}

resource "google_storage_bucket" "gcs-bucket" {
  name     = "%s"
  location = "US"
}
`, name, project, project, disabled, bucketName)
}

func testAccLoggingProjectSink_uniqueWriter(name, bucketName string) string {
	return fmt.Sprintf(`
resource "google_logging_project_sink" "unique_writer" {
  name        = "%s"
  destination = "storage.googleapis.com/${google_storage_bucket.gcs-bucket.name}"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"

  unique_writer_identity = true
}

resource "google_storage_bucket" "gcs-bucket" {
  name     = "%s"
  location = "US"
}
`, name, envvar.GetTestProjectFromEnv(), bucketName)
}

func testAccLoggingProjectSink_uniqueWriterUpdated(name, bucketName string) string {
	return fmt.Sprintf(`
resource "google_logging_project_sink" "unique_writer" {
  name        = "%s"
  destination = "storage.googleapis.com/${google_storage_bucket.gcs-bucket.name}"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=WARNING"

  unique_writer_identity = true
}

resource "google_storage_bucket" "gcs-bucket" {
  name     = "%s"
  location = "US"
}
`, name, envvar.GetTestProjectFromEnv(), bucketName)
}

func testAccLoggingProjectSink_customWriter(org, billingId, project, name, serviceAccount string) string {
	return fmt.Sprintf(`
resource "google_project" "destination-project" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}	

resource "google_logging_project_bucket_config" "destination-bucket" {
  project    = google_project.destination-project.project_id
  location  = "us-central1"
  retention_days = 30
  bucket_id = "shared-bucket"
}

resource "google_service_account" "test-account1" {
  account_id   = "%s"
  display_name = "Log Sink Custom WriterIdentity Testing Account"
}

resource "google_project_iam_member" "custom-sa-logbucket-binding" {
  project = google_project.destination-project.project_id
  role   = "roles/logging.bucketWriter"
  member = "serviceAccount:${google_service_account.test-account1.email}"
}

data "google_project" "testing_project" {
  project_id = "%s"
}

locals {
  project_number = data.google_project.testing_project.number
}

resource "google_service_account_iam_member" "loggingsa-customsa-binding" {
  service_account_id = google_service_account.test-account1.name
  role   = "roles/iam.serviceAccountTokenCreator"
  member = "serviceAccount:service-${local.project_number}@gcp-sa-logging.iam.gserviceaccount.com"
}

resource "time_sleep" "wait_60_seconds" {
  depends_on = [google_service_account_iam_member.loggingsa-customsa-binding]
  create_duration = "60s"
}

resource "google_logging_project_sink" "custom_writer" {
  name        = "%s"
  destination = "logging.googleapis.com/projects/${google_project.destination-project.project_id}/locations/us-central1/buckets/shared-bucket"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"

  unique_writer_identity = true
  custom_writer_identity = "serviceAccount:${google_service_account.test-account1.email}"

  depends_on = [
		google_logging_project_bucket_config.destination-bucket,
		time_sleep.wait_60_seconds,
	]
}
`, project, project, org, billingId, serviceAccount, envvar.GetTestProjectFromEnv(), name, envvar.GetTestProjectFromEnv())
}

func testAccLoggingProjectSink_customWriterUpdated(org, billingId, project, name, serviceAccount string) string {
	return fmt.Sprintf(`
resource "google_project" "destination-project" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}	

resource "google_logging_project_bucket_config" "destination-bucket" {
  project    = google_project.destination-project.project_id
  location  = "us-central1"
  retention_days = 30
  bucket_id = "shared-bucket"
}

resource "google_service_account" "test-account2" {
  account_id   = "%s"
  display_name = "Updated Log Sink Custom WriterIdentity Testing Account"
}

resource "google_project_iam_member" "custom-sa-logbucket-binding" {
  project = google_project.destination-project.project_id
  role   = "roles/logging.bucketWriter"
  member = "serviceAccount:${google_service_account.test-account2.email}"
}

data "google_project" "testing_project" {
  project_id = "%s"
}

locals {
  project_number = data.google_project.testing_project.number
}

resource "google_service_account_iam_member" "loggingsa-customsa-binding" {
  service_account_id = google_service_account.test-account2.name
  role   = "roles/iam.serviceAccountTokenCreator"
  member = "serviceAccount:service-${local.project_number}@gcp-sa-logging.iam.gserviceaccount.com"
}

resource "google_logging_project_sink" "custom_writer" {
  name        = "%s"
  destination = "logging.googleapis.com/projects/${google_project.destination-project.project_id}/locations/us-central1/buckets/shared-bucket"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=WARNING"

  unique_writer_identity = true
  custom_writer_identity = "serviceAccount:${google_service_account.test-account2.email}"

  depends_on = [
	google_logging_project_bucket_config.destination-bucket,
	google_service_account_iam_member.loggingsa-customsa-binding,
	]
}
`, project, project, org, billingId, serviceAccount, envvar.GetTestProjectFromEnv(), name, envvar.GetTestProjectFromEnv())
}

func testAccLoggingProjectSink_heredoc(name, project, bucketName string) string {
	return fmt.Sprintf(`
resource "google_logging_project_sink" "heredoc" {
  name        = "%s"
  project     = "%s"
  destination = "storage.googleapis.com/${google_storage_bucket.gcs-bucket.name}"

  filter = <<EOS

  logName="projects/%s/logs/compute.googleapis.com%%2Factivity_log"
AND severity>=ERROR


EOS

  unique_writer_identity = false
}

resource "google_storage_bucket" "gcs-bucket" {
  name     = "%s"
  location = "US"
}
`, name, project, project, bucketName)
}

func testAccLoggingProjectSink_bigquery_before(sinkName, bqDatasetID string) string {
	return fmt.Sprintf(`
resource "google_logging_project_sink" "bigquery" {
  name        = "%s"
  destination = "bigquery.googleapis.com/projects/%s/datasets/${google_bigquery_dataset.bq_dataset.dataset_id}"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"

  unique_writer_identity = true

  bigquery_options {
    use_partitioned_tables = true
  }
}

resource "google_bigquery_dataset" "bq_dataset" {
  dataset_id  = "%s"
  description = "Log sink (generated during acc test of terraform-provider-google(-beta))."
}
`, sinkName, envvar.GetTestProjectFromEnv(), envvar.GetTestProjectFromEnv(), bqDatasetID)
}

func testAccLoggingProjectSink_bigquery_after(sinkName, bqDatasetID string) string {
	return fmt.Sprintf(`
resource "google_logging_project_sink" "bigquery" {
  name        = "%s"
  destination = "bigquery.googleapis.com/projects/%s/datasets/${google_bigquery_dataset.bq_dataset.dataset_id}"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=WARNING"

  unique_writer_identity = true
}

resource "google_bigquery_dataset" "bq_dataset" {
  dataset_id  = "%s"
  description = "Log sink (generated during acc test of terraform-provider-google(-beta))."
}
`, sinkName, envvar.GetTestProjectFromEnv(), envvar.GetTestProjectFromEnv(), bqDatasetID)
}

func testAccLoggingProjectSink_loggingbucket(name, project string) string {
	return fmt.Sprintf(`
resource "google_logging_project_sink" "loggingbucket" {
  name        = "%s"
  project     = "%s"
  destination = "logging.googleapis.com/projects/%s/locations/global/buckets/_Default"
  exclusions {
    name = "ex1"
    description = "test"
    filter = "resource.type = k8s_container"
  }

  exclusions {
    name = "ex2"
    description = "test-2"
    filter = "resource.type = k8s_container"
  }
}

`, name, project, project)
}
