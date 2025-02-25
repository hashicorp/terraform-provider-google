// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"

	"google.golang.org/api/compute/v1"
)

func TestAccComputeInstantSnapshot_basicFeatures(t *testing.T) {
	var is compute.InstantSnapshot
	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstantSnapshotDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstantSnapshot_basicFeatures(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstantSnapshotExists(t, "google_compute_instant_snapshot.foobar", envvar.GetTestProjectFromEnv(), &is),
				),
			},
		},
	})
}

func TestAccComputeInstantSnapshot_labelsUpdate(t *testing.T) {
	var is compute.InstantSnapshot
	context_1 := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
		"label_key":     "test-1",
		"label_value":   "test-1",
	}
	context_2 := map[string]interface{}{
		"random_suffix": context_1["random_suffix"],
		"label_key":     "test-1",
		"label_value":   "test-2",
	}
	context_3 := map[string]interface{}{
		"random_suffix": context_1["random_suffix"],
		"label_key":     "test-2",
		"label_value":   "test-2",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstantSnapshotDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstantSnapshot_labelsUpdate(context_1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstantSnapshotExists(t, "google_compute_instant_snapshot.foobar", envvar.GetTestProjectFromEnv(), &is),
				),
			},
			{
				Config: testAccComputeInstantSnapshot_labelsUpdate(context_2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstantSnapshotExists(t, "google_compute_instant_snapshot.foobar", envvar.GetTestProjectFromEnv(), &is),
				),
			},
			{
				Config: testAccComputeInstantSnapshot_labelsUpdate(context_3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstantSnapshotExists(t, "google_compute_instant_snapshot.foobar", envvar.GetTestProjectFromEnv(), &is),
				),
			},
		},
	})
}

func testAccCheckComputeInstantSnapshotExists(t *testing.T, n, p string, is *compute.InstantSnapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.GoogleProviderConfig(t)

		zone := tpgresource.GetResourceNameFromSelfLink(rs.Primary.Attributes["zone"])

		found, err := config.NewComputeClient(config.UserAgent).InstantSnapshots.Get(
			p, zone, rs.Primary.Attributes["name"]).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("Instant Snapshot not found")
		}

		*is = *found

		return nil
	}
}

func testAccComputeInstantSnapshot_basicFeatures(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_disk" "disk" {
  name = "tf-test-disk-%{random_suffix}"
  type = "pd-standard"
  zone = "us-central1-a"
  size = 10
}

resource "google_compute_instant_snapshot" "foobar" {
  name = "tf-test-instant-snapshot-%{random_suffix}"
  source_disk = google_compute_disk.disk.self_link
  zone = google_compute_disk.disk.zone

  description = "A test snapshot"
  labels = {
	foo = "bar"
  }
}
`, context)
}

func testAccComputeInstantSnapshot_labelsUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_disk" "disk" {
  name = "tf-test-disk-%{random_suffix}"
  type = "pd-standard"
  zone = "us-central1-a"
  size = 10
}

resource "google_compute_instant_snapshot" "foobar" {
  name = "tf-test-instant-snapshot-%{random_suffix}"
  source_disk = google_compute_disk.disk.self_link
  zone = google_compute_disk.disk.zone

  labels = {
	%{label_key} = "%{label_value}"
  }
}
`, context)
}
