// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package oracledatabase_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"
)

func TestAccOracleDatabaseCloudExadataInfrastructure_basic(t *testing.T) {
	t.Parallel()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOracleDatabaseCloudExadataInfrastructure_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_exadata_infrastructure.my-exadata", "display_name"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_exadata_infrastructure.my-exadata", "name"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_exadata_infrastructure.my-exadata", "gcp_oracle_zone"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_exadata_infrastructure.my-exadata", "properties.#"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_exadata_infrastructure.my-exadata", "properties.0.compute_count"),
					resource.TestCheckResourceAttr("data.google_oracle_database_cloud_exadata_infrastructure.my-exadata", "display_name", "ofake-do-not-delete-tf-exadata"),
					resource.TestCheckResourceAttr("data.google_oracle_database_cloud_exadata_infrastructure.my-exadata", "gcp_oracle_zone", "us-east4-b-r1"),
					resource.TestCheckResourceAttr("data.google_oracle_database_cloud_exadata_infrastructure.my-exadata", "properties.0.state", "AVAILABLE"),
					resource.TestCheckResourceAttr("data.google_oracle_database_cloud_exadata_infrastructure.my-exadata", "properties.0.shape", "Exadata.X9M"),
				),
			},
		},
	})
}

func testAccOracleDatabaseCloudExadataInfrastructure_basic() string {
	return fmt.Sprintf(`
data "google_oracle_database_cloud_exadata_infrastructure" "my-exadata"{
  cloud_exadata_infrastructure_id = "ofake-do-not-delete-tf-exadata"
  project = "oci-terraform-testing-prod"
  location = "us-east4"
}
`)
}
