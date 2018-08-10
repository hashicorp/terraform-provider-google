package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/logging/v2"
)

func TestAccLoggingBillingAccountSink_basic(t *testing.T) {
	t.Parallel()

	sinkName := "tf-test-sink-" + acctest.RandString(10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(10)
	billingAccount := getTestBillingAccountFromEnv(t)

	var sink logging.LogSink

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingBillingAccountSinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBillingAccountSink_basic(sinkName, bucketName, billingAccount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingBillingAccountSinkExists("google_logging_billing_account_sink.basic", &sink),
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

	sinkName := "tf-test-sink-" + acctest.RandString(10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(10)
	updatedBucketName := "tf-test-sink-bucket-" + acctest.RandString(10)
	billingAccount := getTestBillingAccountFromEnv(t)

	var sinkBefore, sinkAfter logging.LogSink

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingBillingAccountSinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBillingAccountSink_update(sinkName, bucketName, billingAccount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingBillingAccountSinkExists("google_logging_billing_account_sink.update", &sinkBefore),
					testAccCheckLoggingBillingAccountSink(&sinkBefore, "google_logging_billing_account_sink.update"),
				),
			}, {
				Config: testAccLoggingBillingAccountSink_update(sinkName, updatedBucketName, billingAccount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingBillingAccountSinkExists("google_logging_billing_account_sink.update", &sinkAfter),
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

func TestAccLoggingBillingAccountSink_heredoc(t *testing.T) {
	t.Parallel()

	sinkName := "tf-test-sink-" + acctest.RandString(10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(10)
	billingAccount := getTestBillingAccountFromEnv(t)

	var sink logging.LogSink

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLoggingBillingAccountSinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingBillingAccountSink_heredoc(sinkName, bucketName, billingAccount),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingBillingAccountSinkExists("google_logging_billing_account_sink.heredoc", &sink),
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

func testAccCheckLoggingBillingAccountSinkDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_logging_billing_account_sink" {
			continue
		}

		attributes := rs.Primary.Attributes

		_, err := config.clientLogging.BillingAccounts.Sinks.Get(attributes["id"]).Do()
		if err == nil {
			return fmt.Errorf("billing sink still exists")
		}
	}

	return nil
}

func testAccCheckLoggingBillingAccountSinkExists(n string, sink *logging.LogSink) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := getResourceAttributes(n, s)
		if err != nil {
			return err
		}
		config := testAccProvider.Meta().(*Config)

		si, err := config.clientLogging.BillingAccounts.Sinks.Get(attributes["id"]).Do()
		if err != nil {
			return err
		}
		*sink = *si

		return nil
	}
}

func testAccCheckLoggingBillingAccountSink(sink *logging.LogSink, n string) resource.TestCheckFunc {
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

		return nil
	}
}

func testAccLoggingBillingAccountSink_basic(name, bucketName, billingAccount string) string {
	return fmt.Sprintf(`
resource "google_logging_billing_account_sink" "basic" {
	name = "%s"
	billing_account = "%s"
	destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
	filter = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}

resource "google_storage_bucket" "log-bucket" {
	name     = "%s"
}`, name, billingAccount, getTestProjectFromEnv(), bucketName)
}

func testAccLoggingBillingAccountSink_update(name, bucketName, billingAccount string) string {
	return fmt.Sprintf(`
resource "google_logging_billing_account_sink" "update" {
	name = "%s"
	billing_account = "%s"
	destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
	filter = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
}

resource "google_storage_bucket" "log-bucket" {
	name     = "%s"
}`, name, billingAccount, getTestProjectFromEnv(), bucketName)
}

func testAccLoggingBillingAccountSink_heredoc(name, bucketName, billingAccount string) string {
	return fmt.Sprintf(`
resource "google_logging_billing_account_sink" "heredoc" {
	name = "%s"
	billing_account = "%s"
	destination = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
	filter = <<EOS

	logName="projects/%s/logs/compute.googleapis.com%%2Factivity_log"
AND severity>=ERROR



  EOS
}

resource "google_storage_bucket" "log-bucket" {
	name     = "%s"
}`, name, billingAccount, getTestProjectFromEnv(), bucketName)
}
