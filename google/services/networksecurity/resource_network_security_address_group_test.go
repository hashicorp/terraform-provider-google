// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networksecurity_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccNetworkSecurityAddressGroups_update(t *testing.T) {
	t.Parallel()

	addressGroupsName := fmt.Sprintf("tf-test-address-group-%s", acctest.RandString(t, 10))
	projectName := envvar.GetTestProjectFromEnv()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkSecurityAddressGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSecurityAddressGroups_basic(addressGroupsName, projectName),
			},
			{
				ResourceName:      "google_network_security_address_group.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkSecurityAddressGroups_update(addressGroupsName, projectName),
			},
			{
				ResourceName:      "google_network_security_address_group.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkSecurityAddressGroups_basic(addressGroupsName, projectName string) string {
	return fmt.Sprintf(`
resource "google_network_security_address_group" "foobar" {
    name        = "%s"
    parent 		= "projects/%s"
    location    = "us-central1"
    description = "my address groups"
    type        = "IPV4"
    capacity    = "100"
    labels      = {
		foo = "bar"
    }
    items 		= ["208.80.154.224/32"]
}
`, addressGroupsName, projectName)
}

func testAccNetworkSecurityAddressGroups_update(addressGroupsName, projectName string) string {
	return fmt.Sprintf(`
resource "google_network_security_address_group" "foobar" {
    name        = "%s"
	parent 		= "projects/%s"
    location    = "us-central1"
    description = "my address groups. Update"
    type        = "IPV4"
    capacity    = "100"
    labels      = {
		foo = "foo"
    }
    items 		= ["208.80.155.224/32", "208.80.154.224/32"]
}
`, addressGroupsName, projectName)
}
