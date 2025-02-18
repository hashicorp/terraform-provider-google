// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package spanner_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccSpannerInstancePartition_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstancePartitionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstancePartition_basic(context),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSpannerInstancePartition_update(context),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerInstancePartition_processingUnits(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstancePartitionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstancePartition_processingUnits(context),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSpannerInstancePartition_processingUnitsUpdate(context),
			},
			{
				ResourceName:      "google_spanner_instance_partition.partition",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSpannerInstancePartition_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name         = "tf-test-spanner-main-%{random_suffix}"
  config       = "nam6"
  display_name = "main-instance"
  num_nodes    = 1
}

resource "google_spanner_instance_partition" "partition" {
  name         = "tf-test-partition-%{random_suffix}"
  instance     = google_spanner_instance.main.name
  config       = "regional-us-central1"
  display_name = "test-spanner-partition"
  node_count   = 1
}
`, context)
}

func testAccSpannerInstancePartition_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name         = "tf-test-spanner-main-%{random_suffix}"
  config       = "nam6"
  display_name = "main-instance"
  num_nodes    = 1
}

resource "google_spanner_instance_partition" "partition" {
  name         = "tf-test-partition-%{random_suffix}"
  instance     = google_spanner_instance.main.name
  config       = "regional-us-central1"
  display_name = "updated-spanner-partition"
  node_count   = 2
}
`, context)
}

func testAccSpannerInstancePartition_processingUnits(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name             = "tf-test-spanner-main-%{random_suffix}"
  config           = "nam6"
  display_name     = "main-instance"
  processing_units = 1000
}

resource "google_spanner_instance_partition" "partition" {
  name             = "tf-test-partition-%{random_suffix}"
  instance         = google_spanner_instance.main.name
  config           = "regional-us-central1"
  display_name     = "test-spanner-partition"
  processing_units = 1000
}
`, context)
}

func testAccSpannerInstancePartition_processingUnitsUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  name             = "tf-test-spanner-main-%{random_suffix}"
  config           = "nam6"
  display_name     = "main-instance"
  processing_units = 1000
}

resource "google_spanner_instance_partition" "partition" {
  name             = "tf-test-partition-%{random_suffix}"
  instance         = google_spanner_instance.main.name
  config           = "regional-us-central1"
  display_name     = "updated-spanner-partition"
  processing_units = 2000
}
`, context)
}
