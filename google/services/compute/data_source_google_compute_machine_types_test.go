// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleComputeMachineTypes_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeMachineTypes,
				Check: resource.ComposeTestCheckFunc(
					// We can't guarantee machine type availability in a given project and zone, so we'll check set-ness rather than correctness
					resource.TestMatchResourceAttr("data.google_compute_machine_types.test", "machine_types.0.name", regexp.MustCompile(`^[a-z0-9-]+$`)),
					resource.TestMatchResourceAttr("data.google_compute_machine_types.test", "machine_types.0.guest_cpus", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr("data.google_compute_machine_types.test", "machine_types.0.memory_mb", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr("data.google_compute_machine_types.test", "machine_types.0.maximum_persistent_disks", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr("data.google_compute_machine_types.test", "machine_types.0.maximum_persistent_disks_size_gb", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr("data.google_compute_machine_types.test", "machine_types.0.description", regexp.MustCompile(`.+`)),
					resource.TestMatchResourceAttr("data.google_compute_machine_types.test", "machine_types.0.is_shared_cpus", regexp.MustCompile(`^true|false$`)),
					resource.TestMatchResourceAttr("data.google_compute_machine_types.test", "machine_types.0.self_link", regexp.MustCompile(`.+`)),
				),
			},
		},
	})
}

const testAccComputeMachineTypes = `
data "google_compute_zones" "available" {}

data "google_compute_machine_types" "test" {
	filter = "guest_cpus > 0"
	zone   = data.google_compute_zones.available.names[0]
}
`
