// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package filestore_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccFilestoreSnapshot_filestoreSnapshotBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFilestoreSnapshotDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreSnapshot_filestoreSnapshotBasicExample(context),
			},
			{
				ResourceName:            "google_filestore_snapshot.snapshot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance", "labels", "location", "name", "terraform_labels"},
			},
		},
	})
}

func testAccFilestoreSnapshot_filestoreSnapshotBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_filestore_snapshot" "snapshot" {
  name     = "tf-test-test-snapshot%{random_suffix}"
  instance = google_filestore_instance.instance.name
  location = "us-east1"
}

resource "google_filestore_instance" "instance" {
  name     = "tf-test-test-instance-for-snapshot%{random_suffix}"
  location = "us-east1"
  tier     = "ENTERPRISE"

  file_shares {
    capacity_gb = 1024
    name        = "share1"
  }

  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
  }
}
`, context)
}

func TestAccFilestoreSnapshot_filestoreSnapshotFullExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFilestoreSnapshotDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreSnapshot_filestoreSnapshotFullExample(context),
			},
			{
				ResourceName:            "google_filestore_snapshot.snapshot",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance", "labels", "location", "name", "terraform_labels"},
			},
		},
	})
}

func testAccFilestoreSnapshot_filestoreSnapshotFullExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_filestore_snapshot" "snapshot" {
  name     = "tf-test-test-snapshot%{random_suffix}"
  instance = google_filestore_instance.instance.name
  location = "us-west1"

  description = "Snapshot of tf-test-test-instance-for-snapshot%{random_suffix}"

  labels = {
    my_label = "value"
  }
}

resource "google_filestore_instance" "instance" {
  name     = "tf-test-test-instance-for-snapshot%{random_suffix}"
  location = "us-west1"
  tier     = "ENTERPRISE"

  file_shares {
    capacity_gb = 1024
    name        = "share1"
  }

  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
  }
}
`, context)
}

func testAccCheckFilestoreSnapshotDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_filestore_snapshot" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{FilestoreBasePath}}projects/{{project}}/locations/{{location}}/instances/{{instance}}/snapshots/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:               config,
				Method:               "GET",
				Project:              billingProject,
				RawURL:               url,
				UserAgent:            config.UserAgent,
				ErrorAbortPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.Is429QuotaError},
			})
			if err == nil {
				return fmt.Errorf("FilestoreSnapshot still exists at %s", url)
			}
		}

		return nil
	}
}
