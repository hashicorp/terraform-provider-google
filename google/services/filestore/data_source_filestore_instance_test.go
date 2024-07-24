// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package filestore_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccFilestoreInstanceDatasource_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFilestoreInstanceDatasourceConfig(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_filestore_instance.filestore", "google_filestore_instance.filestore"),
				),
			},
		},
	})
}

func testAccFilestoreInstanceDatasourceConfig(suffix string) string {
	return fmt.Sprintf(`
resource "google_filestore_instance" "filestore" {
  name        = "tf-instance-%s"
  location    = "us-central1-b"
  tier        = "BASIC_HDD"
  description = "A basic filestore instance created during testing."

  file_shares {
    capacity_gb = 1536
    name        = "share"
  }

  networks {
    network = "default"
    modes   = ["MODE_IPV4"]
  }
}

data "google_filestore_instance" "filestore" {
  name = google_filestore_instance.filestore.name
  location = "us-central1-b"
}
`, suffix)
}
