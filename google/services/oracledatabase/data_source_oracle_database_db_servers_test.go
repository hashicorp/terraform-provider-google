// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package oracledatabase_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccOracleDatabaseDbServers_basic(t *testing.T) {
	t.Parallel()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOracleDatabaseDbServers_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_oracle_database_db_servers.my_db_servers", "db_servers.0.display_name"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_db_servers.my_db_servers", "db_servers.0.properties.#"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_db_servers.my_db_servers", "db_servers.0.properties.0.max_ocpu_count"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_db_servers.my_db_servers", "db_servers.1.display_name"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_db_servers.my_db_servers", "db_servers.1.properties.#"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_db_servers.my_db_servers", "db_servers.1.properties.0.max_ocpu_count"),
					resource.TestCheckResourceAttr("data.google_oracle_database_db_servers.my_db_servers", "db_servers.0.properties.0.max_ocpu_count", "126"),
					resource.TestCheckResourceAttr("data.google_oracle_database_db_servers.my_db_servers", "db_servers.1.properties.0.max_ocpu_count", "126"),
				),
			},
		},
	})
}

const testAccOracleDatabaseDbServers_basic = `
data "google_oracle_database_db_servers" "my_db_servers"{
	location = "us-east4"
	project = "oci-terraform-testing-prod"
	cloud_exadata_infrastructure = "ofake-do-not-delete-tf-exadata"
}
`
