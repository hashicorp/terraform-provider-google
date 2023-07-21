// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging_test

import (
	"fmt"
	"testing"

	"strconv"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/api/logging/v2"
)

func TestAccLoggingFolderSink_basic(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)
	folderName := "tf-test-folder-" + acctest.RandString(t, 10)

	var sink logging.LogSink

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingFolderSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingFolderSink_basic(sinkName, bucketName, folderName, "organizations/"+org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingFolderSinkExists(t, "google_logging_folder_sink.basic", &sink),
					testAccCheckLoggingFolderSink(&sink, "google_logging_folder_sink.basic"),
				),
			}, {
				ResourceName:      "google_logging_folder_sink.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingFolderSink_described(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)
	folderName := "tf-test-folder-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingFolderSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingFolderSink_described(sinkName, bucketName, folderName, "organizations/"+org),
			}, {
				ResourceName:      "google_logging_folder_sink.described",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingFolderSink_disabled(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)
	folderName := "tf-test-folder-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingFolderSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingFolderSink_disabled(sinkName, bucketName, folderName, "organizations/"+org),
			}, {
				ResourceName:      "google_logging_folder_sink.disabled",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingFolderSink_removeOptionals(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)
	folderName := "tf-test-folder-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingFolderSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingFolderSink_basic(sinkName, bucketName, folderName, "organizations/"+org),
			},
			{
				ResourceName:      "google_logging_folder_sink.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingFolderSink_removeOptionals(sinkName, bucketName, folderName, "organizations/"+org),
			},
			{
				ResourceName:      "google_logging_folder_sink.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingFolderSink_folderAcceptsFullFolderPath(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)
	folderName := "tf-test-folder-" + acctest.RandString(t, 10)

	var sink logging.LogSink

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingFolderSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingFolderSink_withFullFolderPath(sinkName, bucketName, folderName, "organizations/"+org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingFolderSinkExists(t, "google_logging_folder_sink.basic", &sink),
					testAccCheckLoggingFolderSink(&sink, "google_logging_folder_sink.basic"),
				),
			}, {
				ResourceName:      "google_logging_folder_sink.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingFolderSink_update(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)
	updatedBucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)
	folderName := "tf-test-folder-" + acctest.RandString(t, 10)
	parent := "organizations/" + org

	var sinkBefore, sinkAfter logging.LogSink

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingFolderSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingFolderSink_basic(sinkName, bucketName, folderName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingFolderSinkExists(t, "google_logging_folder_sink.basic", &sinkBefore),
					testAccCheckLoggingFolderSink(&sinkBefore, "google_logging_folder_sink.basic"),
				),
			}, {
				Config: testAccLoggingFolderSink_basic(sinkName, updatedBucketName, folderName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingFolderSinkExists(t, "google_logging_folder_sink.basic", &sinkAfter),
					testAccCheckLoggingFolderSink(&sinkAfter, "google_logging_folder_sink.basic"),
				),
			}, {
				ResourceName:      "google_logging_folder_sink.basic",
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

func TestAccLoggingFolderSink_updateBigquerySink(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bqDatasetID := "tf_test_sink_" + acctest.RandString(t, 10)
	folderName := "tf-test-folder-" + acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingFolderSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingFolderSink_bigquery_before(sinkName, bqDatasetID, folderName, "organizations/"+org),
			},
			{
				ResourceName:      "google_logging_folder_sink.bigquery",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccLoggingFolderSink_bigquery_after(sinkName, bqDatasetID, folderName, "organizations/"+org),
			},
			{
				ResourceName:      "google_logging_folder_sink.bigquery",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccLoggingFolderSink_heredoc(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	sinkName := "tf-test-sink-" + acctest.RandString(t, 10)
	bucketName := "tf-test-sink-bucket-" + acctest.RandString(t, 10)
	folderName := "tf-test-folder-" + acctest.RandString(t, 10)

	var sink logging.LogSink

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLoggingFolderSinkDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingFolderSink_heredoc(sinkName, bucketName, folderName, "organizations/"+org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoggingFolderSinkExists(t, "google_logging_folder_sink.heredoc", &sink),
					testAccCheckLoggingFolderSink(&sink, "google_logging_folder_sink.heredoc"),
				),
			}, {
				ResourceName:      "google_logging_folder_sink.heredoc",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckLoggingFolderSinkDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_logging_folder_sink" {
				continue
			}

			attributes := rs.Primary.Attributes

			_, err := config.NewLoggingClient(config.UserAgent).Folders.Sinks.Get(attributes["id"]).Do()
			if err == nil {
				return fmt.Errorf("folder sink still exists")
			}
		}

		return nil
	}
}

func testAccCheckLoggingFolderSinkExists(t *testing.T, n string, sink *logging.LogSink) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attributes, err := tpgresource.GetResourceAttributes(n, s)
		if err != nil {
			return err
		}
		config := acctest.GoogleProviderConfig(t)

		si, err := config.NewLoggingClient(config.UserAgent).Folders.Sinks.Get(attributes["id"]).Do()
		if err != nil {
			return err
		}
		*sink = *si

		return nil
	}
}

func testAccCheckLoggingFolderSink(sink *logging.LogSink, n string) resource.TestCheckFunc {
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

func testAccLoggingFolderSink_basic(sinkName, bucketName, folderName, folderParent string) string {
	return fmt.Sprintf(`
resource "google_logging_folder_sink" "basic" {
  name             = "%s"
  folder           = element(split("/", google_folder.my-folder.name), 1)
  destination      = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter           = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
  include_children = true
}

resource "google_storage_bucket" "log-bucket" {
  name     = "%s"
  location = "US"
}

resource "google_folder" "my-folder" {
  display_name = "%s"
  parent       = "%s"
}
`, sinkName, envvar.GetTestProjectFromEnv(), bucketName, folderName, folderParent)
}

func testAccLoggingFolderSink_described(sinkName, bucketName, folderName, folderParent string) string {
	return fmt.Sprintf(`
resource "google_logging_folder_sink" "described" {
  name             = "%s"
  folder           = element(split("/", google_folder.my-folder.name), 1)
  destination      = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  description      = "this is a description for a folder level logging sink"
  filter           = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
  include_children = true
}

resource "google_storage_bucket" "log-bucket" {
  name     = "%s"
  location = "US"
}

resource "google_folder" "my-folder" {
  display_name = "%s"
  parent       = "%s"
}
`, sinkName, envvar.GetTestProjectFromEnv(), bucketName, folderName, folderParent)
}

func testAccLoggingFolderSink_disabled(sinkName, bucketName, folderName, folderParent string) string {
	return fmt.Sprintf(`
resource "google_logging_folder_sink" "disabled" {
  name             = "%s"
  folder           = element(split("/", google_folder.my-folder.name), 1)
  destination      = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  disabled         = true
  filter           = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
  include_children = true
}

resource "google_storage_bucket" "log-bucket" {
  name     = "%s"
  location = "US"
}

resource "google_folder" "my-folder" {
  display_name = "%s"
  parent       = "%s"
}
`, sinkName, envvar.GetTestProjectFromEnv(), bucketName, folderName, folderParent)
}

func testAccLoggingFolderSink_removeOptionals(sinkName, bucketName, folderName, folderParent string) string {
	return fmt.Sprintf(`
resource "google_logging_folder_sink" "basic" {
  name             = "%s"
  folder           = "${element(split("/", google_folder.my-folder.name), 1)}"
  destination      = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter           = ""
  include_children = true
}

resource "google_storage_bucket" "log-bucket" {
  name     = "%s"
  location = "US"
}

resource "google_folder" "my-folder" {
  display_name = "%s"
    parent       = "%s"
}`, sinkName, bucketName, folderName, folderParent)
}

func testAccLoggingFolderSink_withFullFolderPath(sinkName, bucketName, folderName, folderParent string) string {
	return fmt.Sprintf(`
resource "google_logging_folder_sink" "basic" {
  name             = "%s"
  folder           = google_folder.my-folder.name
  destination      = "storage.googleapis.com/${google_storage_bucket.log-bucket.name}"
  filter           = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=ERROR"
  include_children = false
}

resource "google_storage_bucket" "log-bucket" {
  name     = "%s"
  location = "US"
}

resource "google_folder" "my-folder" {
  display_name = "%s"
  parent       = "%s"
}
`, sinkName, envvar.GetTestProjectFromEnv(), bucketName, folderName, folderParent)
}

func testAccLoggingFolderSink_heredoc(sinkName, bucketName, folderName, folderParent string) string {
	return fmt.Sprintf(`
resource "google_logging_folder_sink" "heredoc" {
  name        = "%s"
  folder      = element(split("/", google_folder.my-folder.name), 1)
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

resource "google_folder" "my-folder" {
  display_name = "%s"
  parent       = "%s"
}
`, sinkName, envvar.GetTestProjectFromEnv(), bucketName, folderName, folderParent)
}

func testAccLoggingFolderSink_bigquery_before(sinkName, bqDatasetID, folderName, folderParent string) string {
	return fmt.Sprintf(`
resource "google_logging_folder_sink" "bigquery" {
  name             = "%s"
  folder           = "${element(split("/", google_folder.my-folder.name), 1)}"
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
}

resource "google_folder" "my-folder" {
  display_name = "%s"
  parent       = "%s"
}`, sinkName, envvar.GetTestProjectFromEnv(), envvar.GetTestProjectFromEnv(), bqDatasetID, folderName, folderParent)
}

func testAccLoggingFolderSink_bigquery_after(sinkName, bqDatasetID, folderName, folderParent string) string {
	return fmt.Sprintf(`
resource "google_logging_folder_sink" "bigquery" {
  name             = "%s"
  folder           = "${element(split("/", google_folder.my-folder.name), 1)}"
  destination      = "bigquery.googleapis.com/projects/%s/datasets/${google_bigquery_dataset.logging_sink.dataset_id}"
  filter           = "logName=\"projects/%s/logs/compute.googleapis.com%%2Factivity_log\" AND severity>=WARNING"
  include_children = true
}

resource "google_bigquery_dataset" "logging_sink" {
  dataset_id  = "%s"
  description = "Log sink (generated during acc test of terraform-provider-google(-beta))."
}

resource "google_folder" "my-folder" {
  display_name = "%s"
  parent       = "%s"
}`, sinkName, envvar.GetTestProjectFromEnv(), envvar.GetTestProjectFromEnv(), bqDatasetID, folderName, folderParent)
}
