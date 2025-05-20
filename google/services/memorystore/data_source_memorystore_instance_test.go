// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package memorystore_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccMemorystoreInstanceDatasourceConfig(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMemorystoreInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMemorystoreInstanceDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_memorystore_instance.default", "google_memorystore_instance.instance-basic"),
				),
			},
		},
	})
}

func testAccMemorystoreInstanceDatasourceConfig(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_memorystore_instance" "instance-basic" {
  instance_id                 = "tf-test-memorystore-instance%{random_suffix}"
  shard_count                 = 1
  location                    = "us-central1"
  deletion_protection_enabled = false
}

data "google_memorystore_instance" "default" {
  instance_id                 = google_memorystore_instance.instance-basic.instance_id
  location                    = "us-central1"
}
`, context)
}
