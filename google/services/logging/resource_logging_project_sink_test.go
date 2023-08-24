// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccLoggingProjectSink_basic(t *testing.T) {
	t.Parallel()

	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingProjectSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectSink_basic(sinkName, envvar.GetTestProjectFromEnv(), bucketName),
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

func testAccLoggingProjectSink_basic(name, project, bucketName string) string {
	return fmt.Sprintf(`
resource "google_logging_project_sink" "basic" {
  name        = "%s"
  project     = "%s"
  destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"

  unique_writer_identity = false
}

resource "google_storage_bucket" "log-bucket" {
  name     = "%s"
  location = "US"
}
`, name, project, project, bucketName)
}

func testAccLoggingProjectSink_described(name, project, bucketName string) string {
	return fmt.Sprintf(`
resource "google_logging_project_sink" "described" {
  name        = "%s"
  project     = "%s"
  destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
  description = "this is a description for a project level logging sink"

  unique_writer_identity = false
}

resource "google_storage_bucket" "log-bucket" {
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
  destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
  description = "description updated"

  unique_writer_identity = false
}

resource "google_storage_bucket" "log-bucket" {
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
  destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
  disabled    = true

  unique_writer_identity = false
}

resource "google_storage_bucket" "log-bucket" {
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
  destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
  disabled    = "%s"

  unique_writer_identity = true
}

resource "google_storage_bucket" "log-bucket" {
  name     = "%s"
  location = "US"
}
`, name, project, project, disabled, bucketName)
}

func testAccLoggingProjectSink_uniqueWriter(name, bucketName string) string {
	return fmt.Sprintf(`
resource "google_logging_project_sink" "unique_writer" {
  name        = "%s"
  destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"

  unique_writer_identity = true
}

resource "google_storage_bucket" "log-bucket" {
  name     = "%s"
  location = "US"
}
`, name, envvar.GetTestProjectFromEnv(), bucketName)
}

func testAccLoggingProjectSink_uniqueWriterUpdated(name, bucketName string) string {
	return fmt.Sprintf(`
resource "google_logging_project_sink" "unique_writer" {
  name        = "%s"
  destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=WARNING"

  unique_writer_identity = true
}

resource "google_storage_bucket" "log-bucket" {
  name     = "%s"
  location = "US"
}
`, name, envvar.GetTestProjectFromEnv(), bucketName)
}

func testAccLoggingProjectSink_heredoc(name, project, bucketName string) string {
	return fmt.Sprintf(`
resource "google_logging_project_sink" "heredoc" {
  name        = "%s"
  project     = "%s"
  destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"

  filter = <<EOS

  logName="projects/%s/logs/compute.googleapis.com%%2Factivity_log"
AND severity>=ERROR


EOS

  unique_writer_identity = false
}

resource "google_storage_bucket" "log-bucket" {
  name     = "%s"
  location = "US"
}
`, name, project, project, bucketName)
}

func testAccLoggingProjectSink_bigquery_before(sinkName, bqDatasetID string) string {
	return fmt.Sprintf(`
resource "google_logging_project_sink" "bigquery" {
  name        = "%s"
  destination = "bigquery.googleapis.com/projects/%s/datasets/${google_bigquery_dataset.logging_sink.dataset_id}"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"

  unique_writer_identity = true

  bigquery_options {
    use_partitioned_tables = true
  }
}

resource "google_bigquery_dataset" "logging_sink" {
  dataset_id  = "%s"
  description = "Log sink (generated during acc test of terraform-provider-google(-beta))."
}
`, sinkName, envvar.GetTestProjectFromEnv(), envvar.GetTestProjectFromEnv(), bqDatasetID)
}

func testAccLoggingProjectSink_bigquery_after(sinkName, bqDatasetID string) string {
	return fmt.Sprintf(`
resource "google_logging_project_sink" "bigquery" {
  name        = "%s"
  destination = "bigquery.googleapis.com/projects/%s/datasets/${google_bigquery_dataset.logging_sink.dataset_id}"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=WARNING"

  unique_writer_identity = false
}

resource "google_bigquery_dataset" "logging_sink" {
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

  unique_writer_identity = true
}

`, name, project, project)
}
