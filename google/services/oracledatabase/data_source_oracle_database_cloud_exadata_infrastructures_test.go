// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package oracledatabase_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccOracleDatabaseCloudExadataInfrastructures_basic(t *testing.T) {
	t.Parallel()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOracleDatabaseCloudExadataInfrastructures_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_exadata_infrastructures.my_cloud_exadatas", "cloud_exadata_infrastructures.#"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_exadata_infrastructures.my_cloud_exadatas", "cloud_exadata_infrastructures.0.display_name"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_exadata_infrastructures.my_cloud_exadatas", "cloud_exadata_infrastructures.0.gcp_oracle_zone"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_exadata_infrastructures.my_cloud_exadatas", "cloud_exadata_infrastructures.0.name"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_exadata_infrastructures.my_cloud_exadatas", "cloud_exadata_infrastructures.0.properties.#"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_exadata_infrastructures.my_cloud_exadatas", "cloud_exadata_infrastructures.0.properties.0.cpu_count"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_exadata_infrastructures.my_cloud_exadatas", "cloud_exadata_infrastructures.0.properties.0.compute_count"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_exadata_infrastructures.my_cloud_exadatas", "cloud_exadata_infrastructures.0.properties.0.max_cpu_count"),
					resource.TestCheckResourceAttr("data.google_oracle_database_cloud_exadata_infrastructures.my_cloud_exadatas", "cloud_exadata_infrastructures.0.properties.0.max_cpu_count", "252"),
				),
			},
		},
	})
}

func testAccOracleDatabaseCloudExadataInfrastructures_basic() string {
	return fmt.Sprintf(`
data "google_oracle_database_cloud_exadata_infrastructures" "my_cloud_exadatas"{
  location = "us-east4"
  project = "oci-terraform-testing"
}
`)
}
