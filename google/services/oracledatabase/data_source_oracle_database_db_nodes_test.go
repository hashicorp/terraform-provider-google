// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package oracledatabase_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccOracleDatabaseDbNodes_basic(t *testing.T) {
	t.Parallel()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOracleDatabaseDbNodesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_oracle_database_db_nodes.my_db_nodes", "db_nodes.#"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_db_nodes.my_db_nodes", "db_nodes.0.name"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_db_nodes.my_db_nodes", "db_nodes.0.properties.#"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_db_nodes.my_db_nodes", "db_nodes.1.name"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_db_nodes.my_db_nodes", "db_nodes.1.properties.#"),
					resource.TestCheckResourceAttr("data.google_oracle_database_db_nodes.my_db_nodes", "db_nodes.0.properties.0.state", "AVAILABLE"),
					resource.TestCheckResourceAttr("data.google_oracle_database_db_nodes.my_db_nodes", "db_nodes.1.properties.0.state", "AVAILABLE"),
				),
			},
		},
	})
}

func testAccOracleDatabaseDbNodesConfig() string {
	return fmt.Sprintf(`
data "google_oracle_database_db_nodes" "my_db_nodes"{
	location = "us-east4"
	project = "oci-terraform-testing-prod"
	cloud_vm_cluster = "ofake-do-not-delete-tf-vmcluster"
}
`)
}
