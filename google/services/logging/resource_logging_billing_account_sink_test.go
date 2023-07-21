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
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	"google.golang.org/api/logging/v2"
)

func TestAccLoggingBillingAccountSink_basic(t *testing.T) {
	t.Parallel()

	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)
	billingAccount := envvar.GetTestMasterBillingAccountFromEnv(t)

	var sink logging.LogSink

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingBillingAccountSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBillingAccountSink_basic(sinkName, bucketName, billingAccount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingBillingAccountSinkExists(t, "google_logging_billing_account_sink.basic", &sink),
					testAccCheckLoggingBillingAccountSink(&sink, "google_logging_billing_account_sink.basic"),
				),
			}, {
				ResourceName:      "google_logging_billing_account_sink.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingBillingAccountSink_update(t *testing.T) {
	t.Parallel()

	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)
	updatedBucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)
	billingAccount := envvar.GetTestMasterBillingAccountFromEnv(t)

	var sinkBefore, sinkAfter logging.LogSink

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingBillingAccountSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBillingAccountSink_update(sinkName, bucketName, billingAccount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingBillingAccountSinkExists(t, "google_logging_billing_account_sink.update", &sinkBefore),
					testAccCheckLoggingBillingAccountSink(&sinkBefore, "google_logging_billing_account_sink.update"),
				),
			}, {
				Config: testAccLoggingBillingAccountSink_update(sinkName, updatedBucketName, billingAccount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingBillingAccountSinkExists(t, "google_logging_billing_account_sink.update", &sinkAfter),
					testAccCheckLoggingBillingAccountSink(&sinkAfter, "google_logging_billing_account_sink.update"),
				),
			}, {
				ResourceName:      "google_logging_billing_account_sink.update",
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

func TestAccLoggingBillingAccountSink_described(t *testing.T) {
	t.Parallel()

	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)
	billingAccount := envvar.GetTestMasterBillingAccountFromEnv(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingBillingAccountSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBillingAccountSink_described(sinkName, bucketName, billingAccount),
			}, {
				ResourceName:      "google_logging_billing_account_sink.described",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingBillingAccountSink_disabled(t *testing.T) {
	t.Parallel()

	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)
	billingAccount := envvar.GetTestMasterBillingAccountFromEnv(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingBillingAccountSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBillingAccountSink_disabled(sinkName, bucketName, billingAccount),
			}, {
				ResourceName:      "google_logging_billing_account_sink.disabled",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingBillingAccountSink_updateBigquerySink(t *testing.T) {
	t.Parallel()

	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bqDatasetID := "tf_test_sink_" + acctest.RandString(t, 10)
	billingAccount := envvar.GetTestMasterBillingAccountFromEnv(t)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingBillingAccountSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBillingAccountSink_bigquery_before(sinkName, bqDatasetID, billingAccount),
			},
			{
				ResourceName:      "google_logging_billing_account_sink.bigquery",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingBillingAccountSink_bigquery_after(sinkName, bqDatasetID, billingAccount),
			},
			{
				ResourceName:      "google_logging_billing_account_sink.bigquery",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingBillingAccountSink_heredoc(t *testing.T) {
	t.Parallel()

	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)
	billingAccount := envvar.GetTestMasterBillingAccountFromEnv(t)

	var sink logging.LogSink

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingBillingAccountSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBillingAccountSink_heredoc(sinkName, bucketName, billingAccount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingBillingAccountSinkExists(t, "google_logging_billing_account_sink.heredoc", &sink),
					testAccCheckLoggingBillingAccountSink(&sink, "google_logging_billing_account_sink.heredoc"),
				),
			}, {
				ResourceName:      "google_logging_billing_account_sink.heredoc",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLoggingBillingAccountSinkDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_logging_billing_account_sink" {
				continue
			}

			attributes := rs.Primary.Attributes

			_, err := config.NewLoggingClient(config.UserAgent).BillingAccounts.Sinks.Get(attributes["id"]).Do()
			if err == nil {
				return fmt.Errorf("billing sink still exists")
			}
		}

		return nil
	}
}

func testAccCheckLoggingBillingAccountSinkExists(t *testing.T, n string, sink *logging.LogSink) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := tpgresource.GetResourceAttributes(n, s)
		if err != nil {
			return err
		}
		config := acctest.GoogleProviderConfig(t)

		si, err := config.NewLoggingClient(config.UserAgent).BillingAccounts.Sinks.Get(attributes["id"]).Do()
		if err != nil {
			return err
		}
		*sink = *si

		return nil
	}
}

func testAccCheckLoggingBillingAccountSink(sink *logging.LogSink, n string) resource.TestCheckFunc {
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

		return nil
	}
}

func testAccLoggingBillingAccountSink_basic(name, bucketName, billingAccount string) string {
	return fmt.Sprintf(`
resource "google_logging_billing_account_sink" "basic" {
  name            = "%s"
  billing_account = "%s"
  destination     = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter          = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}

resource "google_storage_bucket" "log-bucket" {
  name     = "%s"
  location = "US"
}
`, name, billingAccount, envvar.GetTestProjectFromEnv(), bucketName)
}

func testAccLoggingBillingAccountSink_described(name, bucketName, billingAccount string) string {
	return fmt.Sprintf(`
resource "google_logging_billing_account_sink" "described" {
  name            = "%s"
  description     = "this is a description"
  billing_account = "%s"
  destination     = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter          = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}

resource "google_storage_bucket" "log-bucket" {
  name     = "%s"
  location = "US"
}
`, name, billingAccount, envvar.GetTestProjectFromEnv(), bucketName)
}

func testAccLoggingBillingAccountSink_disabled(name, bucketName, billingAccount string) string {
	return fmt.Sprintf(`
resource "google_logging_billing_account_sink" "disabled" {
  name            = "%s"
  billing_account = "%s"
  destination     = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter          = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}

resource "google_storage_bucket" "log-bucket" {
  name     = "%s"
  location = "US"
}
`, name, billingAccount, envvar.GetTestProjectFromEnv(), bucketName)
}

func testAccLoggingBillingAccountSink_update(name, bucketName, billingAccount string) string {
	return fmt.Sprintf(`
resource "google_logging_billing_account_sink" "update" {
  name            = "%s"
  billing_account = "%s"
  destination     = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  disabled         = true
  filter          = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}

resource "google_storage_bucket" "log-bucket" {
  name     = "%s"
  location = "US"
}
`, name, billingAccount, envvar.GetTestProjectFromEnv(), bucketName)
}

func testAccLoggingBillingAccountSink_heredoc(name, bucketName, billingAccount string) string {
	return fmt.Sprintf(`
resource "google_logging_billing_account_sink" "heredoc" {
  name            = "%s"
  billing_account = "%s"
  destination     = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter          = <<EOS

  logName="projects/%s/logs/compute.googleapis.com%%2Factivity_log"
AND severity>=ERROR


EOS

}

resource "google_storage_bucket" "log-bucket" {
  name     = "%s"
  location = "US"
}
`, name, billingAccount, envvar.GetTestProjectFromEnv(), bucketName)
}

func testAccLoggingBillingAccountSink_bigquery_before(sinkName, bqDatasetID, billingAccount string) string {
	return fmt.Sprintf(`
resource "google_logging_billing_account_sink" "bigquery" {
  name             = "%s"
  billing_account  = "%s"
  destination      = "bigquery.googleapis.com/projects/%s/datasets/${google_bigquery_dataset.logging_sink.dataset_id}"
  filter           = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"

  bigquery_options {
    use_partitioned_tables = true
  }
}

resource "google_bigquery_dataset" "logging_sink" {
  dataset_id  = "%s"
  description = "Log sink (generated during acc test of terraform-provider-google(-beta))."
}`, sinkName, billingAccount, envvar.GetTestProjectFromEnv(), envvar.GetTestProjectFromEnv(), bqDatasetID)
}

func testAccLoggingBillingAccountSink_bigquery_after(sinkName, bqDatasetID, billingAccount string) string {
	return fmt.Sprintf(`
resource "google_logging_billing_account_sink" "bigquery" {
  name             = "%s"
  billing_account  = "%s"
  destination      = "bigquery.googleapis.com/projects/%s/datasets/${google_bigquery_dataset.logging_sink.dataset_id}"
  filter           = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=WARNING"
}

resource "google_bigquery_dataset" "logging_sink" {
  dataset_id  = "%s"
  description = "Log sink (generated during acc test of terraform-provider-google(-beta))."
}`, sinkName, billingAccount, envvar.GetTestProjectFromEnv(), envvar.GetTestProjectFromEnv(), bqDatasetID)
}
