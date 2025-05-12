// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigtable_test

import (
	"fmt"

	"testing"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBigtableLogicalView_update(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	acctest.SkipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))
	lvName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableLogicalView_update(instanceName, tableName, lvName, "col1", true),
			},
			{
				ResourceName:      "google_bigtable_logical_view.logical_view",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigtableLogicalView_update(instanceName, tableName, lvName, "col2", false),

				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_bigtable_logical_view.logical_view", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:      "google_bigtable_logical_view.logical_view",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBigtableLogicalView_update(instanceName, tableName, lvName, colName string, dp bool) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  cluster {
    cluster_id = "%s-c"
    zone       = "us-east1-b"
  }

  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.id

  column_family {
	family = "CF"
  }
}

resource "google_bigtable_logical_view" "logical_view" {
  logical_view_id = "%s"
  instance        = google_bigtable_instance.instance.name
  deletion_protection = %v
  query = <<EOT
SELECT _key, CF['%s'] AS column
FROM %s
EOT  

  depends_on = [
    google_bigtable_table.table
  ]
}
`, instanceName, instanceName, tableName, lvName, dp, colName, fmt.Sprintf("`%s`", tableName))
}
