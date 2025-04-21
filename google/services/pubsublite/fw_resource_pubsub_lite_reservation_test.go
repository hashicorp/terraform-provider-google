// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package pubsublite_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccResourceFWPubsubLiteReservation_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceFWPubsubLiteReservation_basic(context),
			},
			{
				Config: testAccResourceFWPubsubLiteReservation_upgrade(context),
			},
		},
	})
}

func testAccResourceFWPubsubLiteReservation_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_fwprovider_pubsub_lite_reservation" "basic" {
  name = "tf-test-example-reservation%{random_suffix}"
  region = "us-central1"
  project = data.google_project.project.number
  throughput_capacity = 2
}

data "google_project" "project" {
}
`, context)
}

func testAccResourceFWPubsubLiteReservation_upgrade(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_fwprovider_pubsub_lite_reservation" "basic" {
  name = "tf-test-example-reservation%{random_suffix}"
  region = "us-central1"
  project = data.google_project.project.number
  throughput_capacity = 3
}

data "google_project" "project" {
}
`, context)
}
