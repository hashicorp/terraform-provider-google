// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package filestore_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccFilestoreInstance_restore(t *testing.T) {
	t.Parallel()

	srcInstancetName := fmt.Sprintf("tf-fs-inst-source-%d", acctest.RandInt(t))
	restoreInstanceName := fmt.Sprintf("tf-fs-inst-restored-%d", acctest.RandInt(t))
	backupName := fmt.Sprintf("tf-fs-bkup-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckFilestoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreInstanceRestore_restore(srcInstancetName, restoreInstanceName, backupName),
			},
			{
				ResourceName:            "google_filestore_instance.instance_source",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccFilestoreInstanceRestore_restore(srcInstancetName, restoreInstanceName, backupName string) string {
	return fmt.Sprintf(`
	resource "google_filestore_instance" "instance_source" {
		name        = "%s"
		location    = "us-central1-b"
		tier        = "BASIC_HDD"
		description = "An instance created during testing."
	  
		file_shares {
		  capacity_gb = 1024
		  name        = "volume1"
		}
	  
		networks {
		  network      = "default"
		  modes        = ["MODE_IPV4"]
		  connect_mode = "DIRECT_PEERING"
		}
	}

	resource "google_filestore_instance" "instance_restored" {
		name        = "%s"
		location    = "us-central1-b"
		tier        = "BASIC_HDD"
		description = "An instance created during testing."
	  
		file_shares {
		  capacity_gb = 1024
		  name        = "volume1"
		  source_backup = google_filestore_backup.backup.id
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
		source_instance   = google_filestore_instance.instance_source.id
		source_file_share = "volume1"
	  
		description = "This is a filestore backup for the test instance"
	}
	  
	`, srcInstancetName, restoreInstanceName, backupName)
}
