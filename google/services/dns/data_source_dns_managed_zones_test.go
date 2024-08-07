// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dns_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceDnsManagedZones_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"name-1": fmt.Sprintf("tf-test-zone-%s", acctest.RandString(t, 10)),
		"name-2": fmt.Sprintf("tf-test-zone-%s", acctest.RandString(t, 10)),
	}

	project := envvar.GetTestProjectFromEnv()
	expectedId := fmt.Sprintf("projects/%s/managedZones", project)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: testAccCheckDNSManagedZoneDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
				Config:                   testAccDataSourceDnsManagedZones_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_dns_managed_zones.qa", "id", expectedId),
					resource.TestMatchResourceAttr("data.google_dns_managed_zones.qa", "managed_zones.#", regexp.MustCompile("^[1-9]")), // Non-zero number length

					// Checks below ensure that fields in the first element are set. We can't always make assertions about exact values.
					resource.TestCheckResourceAttr("data.google_dns_managed_zones.qa", "managed_zones.0.project", project),
					resource.TestCheckResourceAttrSet("data.google_dns_managed_zones.qa", "managed_zones.0.name"),
					resource.TestCheckResourceAttrSet("data.google_dns_managed_zones.qa", "managed_zones.0.dns_name"),
					resource.TestCheckResourceAttrSet("data.google_dns_managed_zones.qa", "managed_zones.0.managed_zone_id"),
					resource.TestCheckResourceAttrSet("data.google_dns_managed_zones.qa", "managed_zones.0.visibility"),
					resource.TestCheckResourceAttrSet("data.google_dns_managed_zones.qa", "managed_zones.0.id"),
				),
			},
		},
	})
}

func testAccDataSourceDnsManagedZones_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_dns_managed_zone" "one" {
  name        = "%{name-1}"
  dns_name    = "%{name-1}.hashicorptest.com."
  description = "tf test DNS zone"
}

resource "google_dns_managed_zone" "two" {
  name        = "%{name-2}"
  dns_name    = "%{name-2}.hashicorptest.com."
  description = "tf test DNS zone"
}

data "google_dns_managed_zones" "qa" {
}
`, context)
}
