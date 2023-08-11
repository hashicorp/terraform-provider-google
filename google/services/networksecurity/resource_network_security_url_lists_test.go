// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networksecurity_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkSecurityUrlLists_update(t *testing.T) {
	t.Parallel()

	urlListsName := fmt.Sprintf("tf-test-url-lists-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkSecurityUrlListsDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSecurityUrlLists_basic(urlListsName),
			},
			{
				ResourceName:      "google_network_security_url_lists.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkSecurityUrlLists_update(urlListsName),
			},
			{
				ResourceName:      "google_network_security_url_lists.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkSecurityUrlLists_basic(urlListsName string) string {
	return fmt.Sprintf(`
resource "google_network_security_url_lists" "foobar" {
    name        = "%s"
    location    = "us-central1"
    values = ["www.example.com"]
}
`, urlListsName)
}

func testAccNetworkSecurityUrlLists_update(urlListsName string) string {
	return fmt.Sprintf(`
resource "google_network_security_url_lists" "foobar" {
    name        = "%s"
    location    = "us-central1"
    description = "my description"
    values = ["www.example.com", "about.example.com", "github.com/example-org/*"]
}
`, urlListsName)
}
