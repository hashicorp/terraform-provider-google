// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigqueryreservation_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccBigqueryReservationBiReservation_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
			"time":   {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccBigqueryReservationBiReservation_full(context),
			},
			{
				ResourceName:            "google_bigquery_bi_reservation.reservation",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccBigqueryReservationBiReservation_updateProperties(context),
			},
			{
				ResourceName:            "google_bigquery_bi_reservation.reservation",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
			{
				Config: testAccBigqueryReservationBiReservation_updateLocation(context),
			},
			{
				ResourceName:            "google_bigquery_bi_reservation.reservation",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location"},
			},
		},
	})
}

func testAccBigqueryReservationBiReservation_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_table" "foo" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.bar.dataset_id
  table_id   = "table_%{random_suffix}"
}

resource "google_bigquery_table" "foo2" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.bar2.dataset_id
  table_id   = "table2_%{random_suffix}"
}

resource "google_bigquery_dataset" "bar" {
  dataset_id                  = "dataset_%{random_suffix}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "EU"
}

resource "google_bigquery_dataset" "bar2" {
  dataset_id                  = "dataset2_%{random_suffix}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "EU"
}

resource "google_bigquery_bi_reservation" "reservation" {
	location  = "europe-west1"
	size      = "2800000000"
  preferred_tables {
      project_id  = "%{project}"
      dataset_id  = google_bigquery_dataset.bar.dataset_id
      table_id    = google_bigquery_table.foo.table_id
  }
  preferred_tables {
      project_id  = "%{project}"
      dataset_id  = google_bigquery_dataset.bar2.dataset_id
      table_id    = google_bigquery_table.foo2.table_id
  }
}
`, context)
}

func testAccBigqueryReservationBiReservation_updateProperties(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_table" "foo" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.bar.dataset_id
  table_id   = "table_%{random_suffix}"
}

resource "google_bigquery_table" "foo2" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.bar2.dataset_id
  table_id   = "table2_%{random_suffix}"
}

resource "google_bigquery_dataset" "bar" {
  dataset_id                  = "dataset_%{random_suffix}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "EU"
}

resource "google_bigquery_dataset" "bar2" {
  dataset_id                  = "dataset2_%{random_suffix}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "EU"
}

resource "google_bigquery_bi_reservation" "reservation" {
	location  = "europe-west1"
	size      = "3200000000"
  preferred_tables {
      project_id  = "%{project}"
      dataset_id  = google_bigquery_dataset.bar2.dataset_id
      table_id    = google_bigquery_table.foo2.table_id
  }
}
`, context)
}

func testAccBigqueryReservationBiReservation_updateLocation(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_bigquery_table" "foo" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.bar.dataset_id
  table_id   = "table_%{random_suffix}"
}

resource "google_bigquery_table" "foo2" {
  deletion_protection = false
  dataset_id = google_bigquery_dataset.bar2.dataset_id
  table_id   = "table2_%{random_suffix}"
}

resource "google_bigquery_dataset" "bar" {
  dataset_id                  = "dataset_%{random_suffix}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "EU"
}

resource "google_bigquery_dataset" "bar2" {
  dataset_id                  = "dataset2_%{random_suffix}"
  friendly_name               = "test"
  description                 = "This is a test description"
  location                    = "EU"
}

resource "google_bigquery_bi_reservation" "reservation" {
	location  = "asia-southeast1"
	size      = "3200000000"
  preferred_tables {
      project_id  = "%{project}"
      dataset_id  = google_bigquery_dataset.bar2.dataset_id
      table_id    = google_bigquery_table.foo2.table_id
  }
}
`, context)
}
