// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging_test

import (
	"fmt"
	"testing"

	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	"google.golang.org/api/logging/v2"
)

func TestAccLoggingOrganizationSink_basic(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)

	var sink logging.LogSink

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingOrganizationSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingOrganizationSink_basic(sinkName, bucketName, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingOrganizationSinkExists(t, "google_logging_organization_sink.basic", &sink),
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

	org := envvar.GetTestOrgFromEnv(t)
	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)
	updatedBucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)

	var sinkBefore, sinkAfter logging.LogSink

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingOrganizationSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingOrganizationSink_update(sinkName, bucketName, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingOrganizationSinkExists(t, "google_logging_organization_sink.update", &sinkBefore),
					testAccCheckLoggingOrganizationSink(&sinkBefore, "google_logging_organization_sink.update"),
				),
			}, {
				Config: testAccLoggingOrganizationSink_update(sinkName, updatedBucketName, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingOrganizationSinkExists(t, "google_logging_organization_sink.update", &sinkAfter),
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

func TestAccLoggingOrganizationSink_described(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingOrganizationSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingOrganizationSink_described(sinkName, bucketName, org),
			}, {
				ResourceName:      "google_logging_organization_sink.described",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingOrganizationSink_disabled(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingOrganizationSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingOrganizationSink_disabled(sinkName, bucketName, org),
			}, {
				ResourceName:      "google_logging_organization_sink.disabled",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingOrganizationSink_updateBigquerySink(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bqDatasetID := "tf_test_sink_" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingOrganizationSinkDestroyProducer(t),
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

	org := envvar.GetTestOrgFromEnv(t)
	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)

	var sink logging.LogSink

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingOrganizationSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingOrganizationSink_heredoc(sinkName, bucketName, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingOrganizationSinkExists(t, "google_logging_organization_sink.heredoc", &sink),
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

func testAccCheckLoggingOrganizationSinkDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_logging_organization_sink" {
				continue
			}

			attributes := rs.Primary.Attributes

			_, err := config.NewLoggingClient(config.UserAgent).Organizations.Sinks.Get(attributes["id"]).Do()
			if err == nil {
				return fmt.Errorf("organization sink still exists")
			}
		}

		return nil
	}
}

func testAccCheckLoggingOrganizationSinkExists(t *testing.T, n string, sink *logging.LogSink) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := tpgresource.GetResourceAttributes(n, s)
		if err != nil {
			return err
		}
		config := acctest.GoogleProviderConfig(t)

		si, err := config.NewLoggingClient(config.UserAgent).Organizations.Sinks.Get(attributes["id"]).Do()
		if err != nil {
			return err
		}
		*sink = *si

		return nil
	}
}

func testAccCheckLoggingOrganizationSink(sink *logging.LogSink, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := tpgresource.GetResourceAttributes(n, s)
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
  name     = "%s"
  location = "US"
}
`, sinkName, orgId, envvar.GetTestProjectFromEnv(), bucketName)
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
  name     = "%s"
  location = "US"
}
`, sinkName, orgId, envvar.GetTestProjectFromEnv(), bucketName)
}

func testAccLoggingOrganizationSink_described(sinkName, bucketName, orgId string) string {
	return fmt.Sprintf(`
resource "google_logging_organization_sink" "described" {
  name        = "%s"
  org_id      = "%s"
  destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
  description = "this is a description for an organization level logging sink"
}

resource "google_storage_bucket" "log-bucket" {
  name     = "%s"
  location = "US"
}
`, sinkName, orgId, envvar.GetTestProjectFromEnv(), bucketName)
}

func testAccLoggingOrganizationSink_disabled(sinkName, bucketName, orgId string) string {
	return fmt.Sprintf(`
resource "google_logging_organization_sink" "disabled" {
  name        = "%s"
  org_id      = "%s"
  destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter      = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
  disabled    = true
}

resource "google_storage_bucket" "log-bucket" {
  name     = "%s"
  location = "US"
}
`, sinkName, orgId, envvar.GetTestProjectFromEnv(), bucketName)
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
  name     = "%s"
  location = "US"
}
`, sinkName, orgId, envvar.GetTestProjectFromEnv(), bucketName)
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
}`, sinkName, orgId, envvar.GetTestProjectFromEnv(), envvar.GetTestProjectFromEnv(), bqDatasetID)
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
}`, sinkName, orgId, envvar.GetTestProjectFromEnv(), envvar.GetTestProjectFromEnv(), bqDatasetID)
}
