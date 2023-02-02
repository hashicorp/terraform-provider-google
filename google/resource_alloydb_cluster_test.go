package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlloydbCluster_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlloydbClusterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAlloydbCluster_alloydbClusterBasicExample(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_update(context),
			},
			{
				ResourceName:            "google_alloydb_cluster.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_user", "cluster_id", "location"},
			},
			{
				Config: testAccAlloydbCluster_alloydbClusterBasicExample(context),
			},
		},
	})
}

func testAccAlloydbCluster_update(context map[string]interface{}) string {
	return Nprintf(`
resource "google_alloydb_cluster" "default" {
  cluster_id = "tf-test-alloydb-cluster%{random_suffix}"
  location   = "us-central1"
  network    = "projects/${data.google_project.project.number}/global/networks/${google_compute_network.default.name}"

  labels = {
	update = "true" 
  }

  lifecycle {
    prevent_destroy = true
  }  
}

data "google_project" "project" {
}

resource "google_compute_network" "default" {
  name = "tf-test-alloydb-cluster%{random_suffix}"
}
`, context)
}
