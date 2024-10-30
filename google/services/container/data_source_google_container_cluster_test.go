// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package container_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccContainerClusterDatasource_zonal(t *testing.T) {
	t.Parallel()

	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerClusterDatasource_zonal(acctest.RandString(t, 10), networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_container_cluster.kubes",
						"google_container_cluster.kubes",
						map[string]struct{}{"deletion_protection": {}},
					),
				),
			},
		},
	})
}

func TestAccContainerClusterDatasource_regional(t *testing.T) {
	t.Parallel()

	networkName := acctest.BootstrapSharedTestNetwork(t, "gke-cluster")
	subnetworkName := acctest.BootstrapSubnet(t, "gke-cluster", networkName)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerClusterDatasource_regional(acctest.RandString(t, 10), networkName, subnetworkName),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_container_cluster.kubes",
						"google_container_cluster.kubes",
						map[string]struct{}{
							"deletion_protection": {},
							"resource_labels.%":   {},
						},
					),
				),
			},
		},
	})
}

func testAccContainerClusterDatasource_zonal(suffix, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "kubes" {
  name               = "tf-test-cluster-%s"
  location           = "us-central1-a"
  initial_node_count = 1

  network    = "%s"
  subnetwork = "%s"

  deletion_protection = false
}

data "google_container_cluster" "kubes" {
  name     = google_container_cluster.kubes.name
  location = google_container_cluster.kubes.location
}
`, suffix, networkName, subnetworkName)
}

func testAccContainerClusterDatasource_regional(suffix, networkName, subnetworkName string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "kubes" {
  name               = "tf-test-cluster-%s"
  location           = "us-central1"
  initial_node_count = 1
  resource_labels = {
    created-by = "terraform"
  }
  network    = "%s"
  subnetwork = "%s"

  deletion_protection = false
}

data "google_container_cluster" "kubes" {
  name     = google_container_cluster.kubes.name
  location = google_container_cluster.kubes.location
}
`, suffix, networkName, subnetworkName)
}
