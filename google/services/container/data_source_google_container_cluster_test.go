// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package container_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccContainerClusterDatasource_zonal(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerClusterDatasource_zonal(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_container_cluster.kubes",
						"google_container_cluster.kubes",
						// Remove once https://github.com/hashicorp/terraform/issues/21347 is fixed.
						map[string]struct{}{
							"enable_autopilot":             {},
							"enable_tpu":                   {},
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

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccContainerClusterDatasource_regional(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_container_cluster.kubes",
						"google_container_cluster.kubes",
						// Remove once https://github.com/hashicorp/terraform/issues/21347 is fixed.
						map[string]struct{}{
							"enable_autopilot":             {},
							"enable_tpu":                   {},
							"pod_security_policy_config.#": {},
						},
					),
				),
			},
		},
	})
}

func testAccContainerClusterDatasource_zonal(suffix string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "kubes" {
  name               = "tf-test-cluster-%s"
  location           = "us-central1-a"
  initial_node_count = 1
}

data "google_container_cluster" "kubes" {
  name     = google_container_cluster.kubes.name
  location = google_container_cluster.kubes.location
}
`, suffix)
}

func testAccContainerClusterDatasource_regional(suffix string) string {
	return fmt.Sprintf(`
resource "google_container_cluster" "kubes" {
  name               = "tf-test-cluster-%s"
  location           = "us-central1"
  initial_node_count = 1
}

data "google_container_cluster" "kubes" {
  name     = google_container_cluster.kubes.name
  location = google_container_cluster.kubes.location
}
`, suffix)
}
