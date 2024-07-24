// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networksecurity_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccNetworkSecuritySecurityProfileGroups_update(t *testing.T) {
	t.Parallel()

	orgId := envvar.GetTestOrgFromEnv(t)
	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkSecuritySecurityProfileGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSecuritySecurityProfileGroups_basic(orgId, randomSuffix),
			},
			{
				ResourceName:            "google_network_security_security_profile_group.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkSecuritySecurityProfileGroups_update(orgId, randomSuffix),
			},
			{
				ResourceName:            "google_network_security_security_profile_group.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccNetworkSecuritySecurityProfileGroups_basic(orgId string, randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_network_security_security_profile" "foobar" {
    name        = "tf-test-my-security-profile%s"
    type        = "THREAT_PREVENTION"
    parent      = "organizations/%s"
    location    = "global"
}

resource "google_network_security_security_profile_group" "foobar" {
    name                      = "tf-test-my-security-profile-group%s"
    parent                    = "organizations/%s"
    location                  = "global"
    description               = "My security profile group."
    threat_prevention_profile = google_network_security_security_profile.foobar.id

    labels = {
        foo = "bar"
    }
}
`, randomSuffix, orgId, randomSuffix, orgId)
}

func testAccNetworkSecuritySecurityProfileGroups_update(orgId string, randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_network_security_security_profile" "foobar" {
    name        = "tf-test-my-security-profile%s"
    type        = "THREAT_PREVENTION"
    parent      = "organizations/%s"
    location    = "global"
}

resource "google_network_security_security_profile" "foobar_updated" {
    name        = "tf-test-my-security-profile-updated%s"
    type        = "THREAT_PREVENTION"
    parent      = "organizations/%s"
    location    = "global"
}

resource "google_network_security_security_profile_group" "foobar" {
    name                      = "tf-test-my-security-profile-group%s"
    parent                    = "organizations/%s"
    location                  = "global"
    description               = "My security profile group. Update"
    threat_prevention_profile = google_network_security_security_profile.foobar_updated.id

    labels = {
        foo = "foo"
    }
}
`, randomSuffix, orgId, randomSuffix, orgId, randomSuffix, orgId)
}
