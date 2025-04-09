// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package filestore_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccFilestoreBackup_update(t *testing.T) {
	t.Parallel()

	instName := fmt.Sprintf("tf-test-fs-inst-%d", acctest.RandInt(t))
	bkupName := fmt.Sprintf("tf-test-fs-bkup-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFilestoreBackupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreBackup_create(instName, bkupName),
			},
			{
				ResourceName:            "google_filestore_backup.backup",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccFilestoreBackup_update(instName, bkupName),
			},
			{
				ResourceName:            "google_filestore_backup.backup",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "description", "location"},
			},
		},
	})
}

func testAccFilestoreBackup_create(instName string, bkupName string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name        = "%s"
  location    = "us-central1-b"
  tier        = "BASIC_HDD"
  description = "An instance created during testing."

  file_shares {
    capacity_gb = 1024
    name        = "share22"
  }

  networks {
    network      = "default"
    modes        = ["MODE_IPV4"]
    connect_mode = "DIRECT_PEERING"
  }
}

resource "google_filestore_backup" "backup" {
	name        = "%s"
	location    = "us-central1"
	source_instance   = google_filestore_instance.instance.id
	source_file_share = "share22"

	description = "This is a filestore backup for the test instance"
}

`, instName, bkupName)
}

func testAccFilestoreBackup_update(instName string, bkupName string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name        = "%s"
  location    = "us-central1-b"
  tier        = "BASIC_HDD"
  description = "A modified instance during testing."

  file_shares {
    capacity_gb = 1024
    name        = "share22"
  }

  networks {
    network      = "default"
    modes        = ["MODE_IPV4"]
    connect_mode = "DIRECT_PEERING"
  }

  labels = {
	  "files"      : "label1",
	  "other-label": "update"
  }
}

resource "google_filestore_backup" "backup" {
	name        = "%s"
	location    = "us-central1"
	source_instance   = google_filestore_instance.instance.id
	source_file_share = "share22"

	description = "This is an updated filestore backup for the test instance"
	labels = {
	  "files":"label1",
	  "other-label": "update"
	}
}

`, instName, bkupName)
}

func TestAccFilestoreBackup_tags(t *testing.T) {
	t.Parallel()

	tagKey := acctest.BootstrapSharedTestTagKey(t, "filestore-backups-tagkey")
	context := map[string]interface{}{
		"org":           envvar.GetTestOrgFromEnv(t),
		"tagKey":        tagKey,
		"tagValue":      acctest.BootstrapSharedTestTagValue(t, "filestore-backups-tagvalue", tagKey),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFilestoreBackupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreBackupTags(context),
			},
			{
				ResourceName:            "google_filestore_backup.backup",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "description", "location", "tags"},
			},
		},
	})
}

func testAccFilestoreBackupTags(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_filestore_instance" "instance" {
  name     = "tf-test-fs-inst-%{random_suffix}"
  location = "us-central1-b"
  tier     = "BASIC_HDD"

  file_shares {
    capacity_gb = 1024
    name        = "share1"
  }

  networks {
    network      = "default"
    modes        = ["MODE_IPV4"]
    connect_mode = "DIRECT_PEERING"
  }
}

resource "google_filestore_backup" "backup" {
  name              = "tf-test-fs-bkup-%{random_suffix}"
  location          = "us-central1"
  description       = "This is a filestore backup for the test instance"
  source_instance   = google_filestore_instance.instance.id
  source_file_share = "share1"

  labels = {
    "files":"label1",
    "other-label": "label2"
  }
  tags = {
    "%{org}/%{tagKey}" = "%{tagValue}"
  }
}
`, context)
}
