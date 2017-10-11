package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/logging/v2"
	"strconv"
	"testing"
)

func TestAccLoggingProjectSink_basic(t *testing.T) {
	sinkName := "tf-test-sink-" + acctest.RandString(10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(10)

	var sink logging.LogSink

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingProjectSinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectSink_basic(sinkName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingProjectSinkExists("google_logging_project_sink.basic", &sink),
					testAccCheckLoggingProjectSink(&sink, "google_logging_project_sink.basic"),
				),
			},
		},
	})
}

func TestAccLoggingProjectSink_uniqueWriter(t *testing.T) {
	sinkName := "tf-test-sink-" + acctest.RandString(10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(10)

	var sink logging.LogSink

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingProjectSinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectSink_uniqueWriter(sinkName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingProjectSinkExists("google_logging_project_sink.unique_writer", &sink),
					testAccCheckLoggingProjectSink(&sink, "google_logging_project_sink.unique_writer"),
				),
			},
		},
	})
}

func TestAccLoggingProjectSink_updatePreservesUniqueWriter(t *testing.T) {
	sinkName := "tf-test-sink-" + acctest.RandString(10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(10)
	updatedBucketName := "tf-test-sink-bucket-" + acctest.RandString(10)

	var sinkBefore, sinkAfter logging.LogSink

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingProjectSinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingProjectSink_uniqueWriter(sinkName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingProjectSinkExists("google_logging_project_sink.unique_writer", &sinkBefore),
					testAccCheckLoggingProjectSink(&sinkBefore, "google_logging_project_sink.unique_writer"),
				),
			}, {
				Config: testAccLoggingProjectSink_uniqueWriterUpdated(sinkName, updatedBucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingProjectSinkExists("google_logging_project_sink.unique_writer", &sinkAfter),
					testAccCheckLoggingProjectSink(&sinkAfter, "google_logging_project_sink.unique_writer"),
				),
			},
		},
	})

	// Destination and Filter should have changed, but WriterIdentity should be the same
	if sinkBefore.Destination == sinkAfter.Destination {
		t.Errorf("Expected Destination to change, but it didn't: Destination = %#v", sinkBefore.Destination)
	}
	if sinkBefore.Filter == sinkAfter.Filter {
		t.Errorf("Expected Filter to change, but it didn't: Filter = %#v", sinkBefore.Filter)
	}
	if sinkBefore.WriterIdentity != sinkAfter.WriterIdentity {
		t.Errorf("Expected WriterIdentity to be the same, but it differs: before = %#v, after = %#v",
			sinkBefore.WriterIdentity, sinkAfter.WriterIdentity)
	}
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

func testAccCheckLoggingProjectSinkExists(n string, sink *logging.LogSink) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := getResourceAttributes(n, s)
		if err != nil {
			return err
		}
		config := testAccProvider.Meta().(*Config)

		si, err := config.clientLogging.Projects.Sinks.Get(attributes["id"]).Do()
		if err != nil {
			return err
		}
		*sink = *si

		return nil
	}
}

func testAccCheckLoggingProjectSink(sink *logging.LogSink, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := getResourceAttributes(n, s)
		if err != nil {
			return err
		}

		if sink.Destination != attributes["destination"] {
			return fmt.Errorf("mismatch on destination: api has %s but client has %s", sink.Destination, attributes["destination"])
		}

		if sink.Filter != attributes["filter"] {
			return fmt.Errorf("mismatch on filter: api has %s but client has %s", sink.Filter, attributes["filter"])
		}

		apiLooksUnique := strconv.FormatBool(nonUniqueWriterAccount != attributes["writer_identity"])
		if apiLooksUnique != attributes["unique_writer_identity"] {
			return fmt.Errorf("mismatch on unique_writer_identity: api looks like %s but client has %s", apiLooksUnique, attributes["unique_writer_identity"])
		}

		if sink.WriterIdentity != attributes["writer_identity"] {
			return fmt.Errorf("mismatch on writer_identity: api has %s but client has %s", sink.WriterIdentity, attributes["writer_identity"])
		}

		return nil
	}
}

func testAccLoggingProjectSink_basic(name, bucketName string) string {
	return fmt.Sprintf(`
resource "google_logging_project_sink" "basic" {
	name = "%s"
	destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
	filter = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
	unique_writer_identity = false
}

resource "google_storage_bucket" "log-bucket" {
	name     = "%s"
}`, name, getTestProjectFromEnv(), bucketName)
}

func testAccLoggingProjectSink_uniqueWriter(name, bucketName string) string {
	return fmt.Sprintf(`
resource "google_logging_project_sink" "unique_writer" {
	name = "%s"
	destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
	filter = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
	unique_writer_identity = true
}

resource "google_storage_bucket" "log-bucket" {
	name     = "%s"
}`, name, getTestProjectFromEnv(), bucketName)
}

func testAccLoggingProjectSink_uniqueWriterUpdated(name, bucketName string) string {
	return fmt.Sprintf(`
resource "google_logging_project_sink" "unique_writer" {
	name = "%s"
	destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
	filter = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=WARNING"
	unique_writer_identity = true
}

resource "google_storage_bucket" "log-bucket" {
	name     = "%s"
}`, name, getTestProjectFromEnv(), bucketName)
}
