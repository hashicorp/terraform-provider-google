package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccContainerClusterDatasource_zonal(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerClusterDatasource_zonal(),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_container_cluster.kubes",
						"google_container_cluster.kubes",
						// Remove once https://github.com/hashicorp/terraform/issues/21347 is fixed.
						map[string]struct{}{
							"enable_tpu":                   {},
							"enable_binary_authorization":  {},
							"pod_security_policy_config.#": {},
						},
					),
				),
			},
		},
	})
}

func TestAccContainerClusterDatasource_regional(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccContainerClusterDatasource_regional(),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_container_cluster.kubes",
						"google_container_cluster.kubes",
						// Remove once https://github.com/hashicorp/terraform/issues/21347 is fixed.
						map[string]struct{}{
							"enable_tpu":                   {},
							"enable_binary_authorization":  {},
							"pod_security_policy_config.#": {},
						},
					),
				),
			},
		},
	})
}

func testAccContainerClusterDatasource_zonal() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "kubes" {
  name               = "cluster-test-%s"
  location           = "us-central1-a"
  initial_node_count = 1

  master_auth {
    username = "mr.yoda"
    password = "adoy.rm.123456789"
  }
}

data "google_container_cluster" "kubes" {
  name     = google_container_cluster.kubes.name
  location = google_container_cluster.kubes.location
}
`, acctest.RandString(10))
}

func testAccContainerClusterDatasource_regional() string {
	return fmt.Sprintf(`
resource "google_container_cluster" "kubes" {
  name               = "cluster-test-%s"
  location           = "us-central1"
  initial_node_count = 1
}

data "google_container_cluster" "kubes" {
  name     = google_container_cluster.kubes.name
  location = google_container_cluster.kubes.location
}
`, acctest.RandString(10))
}
