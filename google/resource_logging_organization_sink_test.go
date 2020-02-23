package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"google.golang.org/api/logging/v2"
	"strconv"
)

func TestAccLoggingOrganizationSink_basic(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	sinkName := "tf-test-sink-" + acctest.RandString(10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(10)

	var sink logging.LogSink

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingOrganizationSinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingOrganizationSink_basic(sinkName, bucketName, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingOrganizationSinkExists("google_logging_organization_sink.basic", &sink),
					testAccCheckLoggingOrganizationSink(&sink, "google_logging_organization_sink.basic"),
				),
			}, {
				ResourceName:      "google_logging_organization_sink.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingOrganizationSink_update(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	sinkName := "tf-test-sink-" + acctest.RandString(10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(10)
	updatedBucketName := "tf-test-sink-bucket-" + acctest.RandString(10)

	var sinkBefore, sinkAfter logging.LogSink

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingOrganizationSinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingOrganizationSink_update(sinkName, bucketName, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingOrganizationSinkExists("google_logging_organization_sink.update", &sinkBefore),
					testAccCheckLoggingOrganizationSink(&sinkBefore, "google_logging_organization_sink.update"),
				),
			}, {
				Config: testAccLoggingOrganizationSink_update(sinkName, updatedBucketName, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingOrganizationSinkExists("google_logging_organization_sink.update", &sinkAfter),
					testAccCheckLoggingOrganizationSink(&sinkAfter, "google_logging_organization_sink.update"),
				),
			}, {
				ResourceName:      "google_logging_organization_sink.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})

	// Destination should have changed, but WriterIdentity should be the same
	if sinkBefore.Destination == sinkAfter.Destination {
		t.Errorf("Expected Destination to change, but it didn't: Destination = %#v", sinkBefore.Destination)
	}
	if sinkBefore.WriterIdentity != sinkAfter.WriterIdentity {
		t.Errorf("Expected WriterIdentity to be the same, but it differs: before = %#v, after = %#v",
			sinkBefore.WriterIdentity, sinkAfter.WriterIdentity)
	}
}

func TestAccLoggingOrganizationSink_updateBigquerySink(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	sinkName := "tf-test-sink-" + acctest.RandString(10)
	bqDatasetID := "tf_test_sink_" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingOrganizationSinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingOrganizationSink_bigquery_before(sinkName, bqDatasetID, org),
			},
			{
				ResourceName:      "google_logging_organization_sink.bigquery",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingOrganizationSink_bigquery_after(sinkName, bqDatasetID, org),
			},
			{
				ResourceName:      "google_logging_organization_sink.bigquery",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingOrganizationSink_heredoc(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	sinkName := "tf-test-sink-" + acctest.RandString(10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(10)

	var sink logging.LogSink

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingOrganizationSinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingOrganizationSink_heredoc(sinkName, bucketName, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingOrganizationSinkExists("google_logging_organization_sink.heredoc", &sink),
					testAccCheckLoggingOrganizationSink(&sink, "google_logging_organization_sink.heredoc"),
				),
			}, {
				ResourceName:      "google_logging_organization_sink.heredoc",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLoggingOrganizationSinkDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_logging_organization_sink" {
			continue
		}

		attributes := rs.Primary.Attributes

		_, err := config.clientLogging.Organizations.Sinks.Get(attributes["id"]).Do()
		if err == nil {
			return fmt.Errorf("organization sink still exists")
		}
	}

	return nil
}

func testAccCheckLoggingOrganizationSinkExists(n string, sink *logging.LogSink) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := getResourceAttributes(n, s)
		if err != nil {
			return err
		}
		config := testAccProvider.Meta().(*Config)

		si, err := config.clientLogging.Organizations.Sinks.Get(attributes["id"]).Do()
		if err != nil {
			return err
		}
		*sink = *si

		return nil
	}
}

func testAccCheckLoggingOrganizationSink(sink *logging.LogSink, n string) resource.TestCheckFunc {
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

		if sink.WriterIdentity != attributes["writer_identity"] {
			return fmt.Errorf("mismatch on writer_identity: api has %s but client has %s", sink.WriterIdentity, attributes["writer_identity"])
		}

		includeChildren := false
		if attributes["include_children"] != "" {
			includeChildren, err = strconv.ParseBool(attributes["include_children"])
			if err != nil {
				return err
			}
		}
		if sink.IncludeChildren != includeChildren {
			return fmt.Errorf("mismatch on include_children: api has %v but client has %v", sink.IncludeChildren, includeChildren)
		}

		return nil
	}
}

func testAccLoggingOrganizationSink_basic(sinkName, bucketName, orgId string) string {
	return fmt.Sprintf(`
resource "google_logging_organization_sink" "basic" {
  name             = "%s"
  org_id           = "%s"
  destination      = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter           = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
  include_children = true
}

resource "google_storage_bucket" "log-bucket" {
  name = "%s"
}
`, sinkName, orgId, getTestProjectFromEnv(), bucketName)
}

func testAccLoggingOrganizationSink_update(sinkName, bucketName, orgId string) string {
	return fmt.Sprintf(`
resource "google_logging_organization_sink" "update" {
  name             = "%s"
  org_id           = "%s"
  destination      = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter           = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
  include_children = false
}

resource "google_storage_bucket" "log-bucket" {
  name = "%s"
}
`, sinkName, orgId, getTestProjectFromEnv(), bucketName)
}

func testAccLoggingOrganizationSink_heredoc(sinkName, bucketName, orgId string) string {
	return fmt.Sprintf(`
resource "google_logging_organization_sink" "heredoc" {
  name        = "%s"
  org_id      = "%s"
  destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter      = <<EOS

	logName="projects/%s/logs/compute.googleapis.com%%2Factivity_log"
AND severity>=ERROR



EOS

  include_children = true
}

resource "google_storage_bucket" "log-bucket" {
  name = "%s"
}
`, sinkName, orgId, getTestProjectFromEnv(), bucketName)
}

func testAccLoggingOrganizationSink_bigquery_before(sinkName, bqDatasetID, orgId string) string {
	return fmt.Sprintf(`
resource "google_logging_organization_sink" "bigquery" {
  name             = "%s"
  org_id           = "%s"
  destination      = "bigquery.googleapis.com/projects/%s/datasets/${google_bigquery_dataset.logging_sink.dataset_id}"
  filter           = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
  include_children = true

  bigquery_options {
    use_partitioned_tables = true
  }
}

resource "google_bigquery_dataset" "logging_sink" {
  dataset_id  = "%s"
  description = "Log sink (generated during acc test of terraform-provider-google(-beta))."
}`, sinkName, orgId, getTestProjectFromEnv(), getTestProjectFromEnv(), bqDatasetID)
}

func testAccLoggingOrganizationSink_bigquery_after(sinkName, bqDatasetID, orgId string) string {
	return fmt.Sprintf(`
resource "google_logging_organization_sink" "bigquery" {
  name             = "%s"
  org_id           = "%s"
  destination      = "bigquery.googleapis.com/projects/%s/datasets/${google_bigquery_dataset.logging_sink.dataset_id}"
  filter           = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=WARNING"
  include_children = true
}

resource "google_bigquery_dataset" "logging_sink" {
  dataset_id  = "%s"
  description = "Log sink (generated during acc test of terraform-provider-google(-beta))."
}`, sinkName, orgId, getTestProjectFromEnv(), getTestProjectFromEnv(), bqDatasetID)
}
