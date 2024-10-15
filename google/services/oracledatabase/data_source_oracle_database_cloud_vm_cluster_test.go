// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package oracledatabase_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccOracleDatabaseCloudVmCluster_basic(t *testing.T) {
	t.Parallel()
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccOracleDatabaseCloudVmCluster_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_vm_cluster.my-vmcluster", "display_name"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_vm_cluster.my-vmcluster", "exadata_infrastructure"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_vm_cluster.my-vmcluster", "cidr"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_vm_cluster.my-vmcluster", "properties.#"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_vm_cluster.my-vmcluster", "properties.0.cpu_core_count"),
					resource.TestCheckResourceAttrSet("data.google_oracle_database_cloud_vm_cluster.my-vmcluster", "properties.0.cluster_name"),
					resource.TestCheckResourceAttr("data.google_oracle_database_cloud_vm_cluster.my-vmcluster", "gcp_oracle_zone", "us-east4-b-r1"),
					resource.TestCheckResourceAttr("data.google_oracle_database_cloud_vm_cluster.my-vmcluster", "properties.0.node_count", "2"),
					resource.TestCheckResourceAttr("data.google_oracle_database_cloud_vm_cluster.my-vmcluster", "properties.0.state", "AVAILABLE"),
					resource.TestCheckResourceAttr("data.google_oracle_database_cloud_vm_cluster.my-vmcluster", "properties.0.shape", "Exadata.X9M"),
				),
			},
		},
	})
}

func testAccOracleDatabaseCloudVmCluster_basic() string {
	return fmt.Sprintf(`
data "google_oracle_database_cloud_vm_cluster" "my-vmcluster"{
  cloud_vm_cluster_id = "ofake-do-not-delete-tf-vmcluster"
  project = "oci-terraform-testing"
  location = "us-east4"
}
`)
}
