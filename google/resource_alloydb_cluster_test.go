package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlloydbCluster_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAlloydbClusterDestroyProducer(t),
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
	foo = "bar" 
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
