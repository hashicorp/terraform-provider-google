package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccBigtableAppProfile_update(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableAppProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableAppProfile_multiClusterRouting(instanceName),
			},
			{
				ResourceName:            "google_bigtable_app_profile.ap",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_warnings"},
			},
			{
				Config: testAccBigtableAppProfile_update(instanceName),
			},
			{
				ResourceName:            "google_bigtable_app_profile.ap",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_warnings"},
			},
		},
	})
}

func testAccBigtableAppProfile_multiClusterRouting(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
  cluster {
    cluster_id   = "%s"
    zone         = "us-central1-b"
    num_nodes    = 3
    storage_type = "HDD"
  }
}

resource "google_bigtable_app_profile" "ap" {
  instance       = google_bigtable_instance.instance.id
  app_profile_id = "test"

  multi_cluster_routing_use_any = true
  ignore_warnings               = true
}
`, instanceName, instanceName)
}

func testAccBigtableAppProfile_update(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"
  cluster {
    cluster_id   = "%s"
    zone         = "us-central1-b"
    num_nodes    = 3
    storage_type = "HDD"
  }
}

resource "google_bigtable_app_profile" "ap" {
  instance       = google_bigtable_instance.instance.id
  app_profile_id = "test"
  description    = "add a description"

  multi_cluster_routing_use_any = true
  ignore_warnings               = true
}
`, instanceName, instanceName)
}
