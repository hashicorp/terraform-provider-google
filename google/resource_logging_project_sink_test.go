package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccLoggingProjectSink_basic(t *testing.T) {
	t.Parallel()

	sinkName := "tf-test-sink-" + acctest.RandString(10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingProjectSinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectSink_basic(sinkName, getTestProjectFromEnv(), bucketName),
			},
			{
				ResourceName:      "google_logging_project_sink.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingProjectSink_updatePreservesUniqueWriter(t *testing.T) {
	t.Parallel()

	sinkName := "tf-test-sink-" + acctest.RandString(10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(10)
	updatedBucketName := "tf-test-sink-bucket-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingProjectSinkDestroy,
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

	sinkName := "tf-test-sink-" + acctest.RandString(10)
	bqDatasetID := "tf_test_sink_" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingProjectSinkDestroy,
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

	sinkName := "tf-test-sink-" + acctest.RandString(10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingProjectSinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectSink_heredoc(sinkName, getTestProjectFromEnv(), bucketName),
			},
			{
				ResourceName:      "google_logging_project_sink.heredoc",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLoggingProjectSinkDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_logging_project_sink" {
			continue
		}

		attributes := rs.Primary.Attributes

		_, err := config.clientLogging.Projects.Sinks.Get(attributes["id"]).Do()
		if err == nil {
			return fmt.Errorf("project sink still exists")
		}
	}

	return nil
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
  name = "%s"
}
`, name, project, project, bucketName)
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
  name = "%s"
}
`, name, getTestProjectFromEnv(), bucketName)
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
  name = "%s"
}
`, name, getTestProjectFromEnv(), bucketName)
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
  name = "%s"
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
`, sinkName, getTestProjectFromEnv(), getTestProjectFromEnv(), bqDatasetID)
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
`, sinkName, getTestProjectFromEnv(), getTestProjectFromEnv(), bqDatasetID)
}
