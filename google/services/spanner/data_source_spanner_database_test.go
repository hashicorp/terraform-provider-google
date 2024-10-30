// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package spanner_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceSpannerDatabase_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerDatabaseDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSpannerDatabaseBasic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_spanner_database.bar",
						"google_spanner_database.foo",
						map[string]struct{}{
							"ddl.#":               {},
							"ddl.0":               {},
							"deletion_protection": {},
						},
					),
				),
			},
		},
	})
}

func testAccDataSourceSpannerDatabaseBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "instance" {
  name         = "tf-test-instance-%{random_suffix}"
  display_name = "Test spanner instance"

  config           = "regional-us-central1"
  processing_units = 200
}

resource "google_spanner_database" "foo" {
  name     = "tf-test-db-%{random_suffix}"
  instance = google_spanner_instance.instance.name
  ddl = [
    "CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
  ]

  deletion_protection = false
}

data "google_spanner_database" "bar" {
  name     = google_spanner_database.foo.name
  instance = google_spanner_instance.instance.name
}
`, context)
}
