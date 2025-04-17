// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package gkebackup_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGKEBackupBackupChannel_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":             envvar.GetTestProjectFromEnv(),
		"destination_project": "projects/331279474308",
		"random_suffix":       acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGKEBackupBackupChannelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEBackupBackupChannel_basic(context),
			},
			{
				ResourceName:            "google_gke_backup_backup_channel.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "terraform_labels"},
			},
		},
	})
}

func testAccGKEBackupBackupChannel_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_gke_backup_backup_channel" "basic" {
  name = "tf-test-basic-channel%{random_suffix}"
  location = "us-central1"
  description = ""
  destination_project = "%{destination_project}"
  labels = { "key": "some-value" }
}
`, context)
}
