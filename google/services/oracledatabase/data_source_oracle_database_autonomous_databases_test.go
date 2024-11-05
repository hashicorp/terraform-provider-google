// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package oracledatabase_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccOracleDatabaseAutonomousDatabases_basic(t *testing.T) {
	t.Parallel()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOracleDatabaseAutonomousDatabases_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_oracle_database_autonomous_databases.my-adbs", "autonomous_databases.#"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_autonomous_databases.my-adbs", "autonomous_databases.0.display_name"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_autonomous_databases.my-adbs", "autonomous_databases.0.cidr"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_autonomous_databases.my-adbs", "autonomous_databases.0.network"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_autonomous_databases.my-adbs", "autonomous_databases.0.entitlement_id"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_autonomous_databases.my-adbs", "autonomous_databases.0.database"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_autonomous_databases.my-adbs", "autonomous_databases.0.properties.#"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_autonomous_databases.my-adbs", "autonomous_databases.0.properties.0.state"),
				),
			},
		},
	})
}

func testAccOracleDatabaseAutonomousDatabases_basic() string {
	return fmt.Sprintf(`
data "google_oracle_database_autonomous_databases" "my-adbs"{
  location = "us-east4"
  project = "oci-terraform-testing"
}
`)
}
