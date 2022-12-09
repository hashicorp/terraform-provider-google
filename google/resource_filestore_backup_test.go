package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFilestoreBackup_update(t *testing.T) {
	t.Parallel()

	instName := fmt.Sprintf("tf-fs-inst-%d", randInt(t))
	bkupName := fmt.Sprintf("tf-fs-bkup-%d", randInt(t))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFilestoreBackupDestroyProducer(t),
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
				ImportStateVerifyIgnore: []string{"labels", "description", "location"},
			},
		},
	})
}

func testAccFilestoreBackup_create(instName string, bkupName string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "instance" {
  name = "%s"
  location = "us-central1-b"
  tier = "BASIC_SSD"

  file_shares {
    capacity_gb = 2560
    name        = "share22"
  }

  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
    connect_mode = "DIRECT_PEERING"
  }
  description = "An instance created during testing."
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
  name = "%s"
  location = "us-central1-b"
  tier = "BASIC_SSD"

  file_shares {
    capacity_gb = 2560
    name        = "share22"
  }

  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
    connect_mode = "DIRECT_PEERING"
  }

  labels = {
	"files":"label1",
	"other-label": "update"
  }

  description = "A modified instance during testing."
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
