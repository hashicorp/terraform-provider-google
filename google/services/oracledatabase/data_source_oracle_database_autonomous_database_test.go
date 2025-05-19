// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package oracledatabase_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccOracleDatabaseAutonomousDatabase_basic(t *testing.T) {
	t.Parallel()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOracleDatabaseAutonomousDatabase_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_oracle_database_autonomous_database.my-adb", "display_name"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_autonomous_database.my-adb", "database"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_autonomous_database.my-adb", "cidr"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_autonomous_database.my-adb", "network"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_autonomous_database.my-adb", "properties.#"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_autonomous_database.my-adb", "properties.0.character_set"),
				),
			},
		},
	})
}

func testAccOracleDatabaseAutonomousDatabase_basic() string {
	return fmt.Sprintf(`
data "google_oracle_database_autonomous_database" "my-adb"{
	autonomous_database_id = "do-not-delete-tf-adb"
	location = "us-east4"
	project = "oci-terraform-testing-prod"
}
`)
}
