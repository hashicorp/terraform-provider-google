// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networksecurity_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccNetworkSecuritySecurityProfiles_update(t *testing.T) {
	t.Parallel()

	orgId := envvar.GetTestOrgFromEnv(t)
	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkSecuritySecurityProfileDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSecuritySecurityProfiles_basic(orgId, randomSuffix),
			},
			{
				ResourceName:            "google_network_security_security_profile.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkSecuritySecurityProfiles_update(orgId, randomSuffix),
			},
			{
				ResourceName:            "google_network_security_security_profile.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccNetworkSecuritySecurityProfiles_antivirusOverrides(t *testing.T) {
	t.Parallel()

	orgId := envvar.GetTestOrgFromEnv(t)
	randomSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkSecuritySecurityProfileDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSecuritySecurityProfiles_basic(orgId, randomSuffix),
			},
			{
				ResourceName:            "google_network_security_security_profile.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkSecuritySecurityProfiles_antivirusOverrides(orgId, randomSuffix),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_network_security_security_profile.foobar", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_network_security_security_profile.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkSecuritySecurityProfiles_basic(orgId, randomSuffix),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_network_security_security_profile.foobar", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_network_security_security_profile.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccNetworkSecuritySecurityProfiles_basic(orgId string, randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_network_security_security_profile" "foobar" {
    name        = "tf-test-my-security-profile%s"
    parent      = "organizations/%s"
    location    = "global"
    description = "My security profile."
    type        = "THREAT_PREVENTION"

    labels = {
        foo = "bar"
    }
}
`, randomSuffix, orgId)
}

func testAccNetworkSecuritySecurityProfiles_update(orgId string, randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_network_security_security_profile" "foobar" {
    name        = "tf-test-my-security-profile%s"
    parent      = "organizations/%s"
    location    = "global"
    description = "My security profile. Update"
    type        = "THREAT_PREVENTION"

    labels = {
        foo = "foo"
    }

    threat_prevention_profile {
        severity_overrides {
            action   = "ALLOW"
            severity = "INFORMATIONAL"
        }

        severity_overrides {
            action   = "DENY"
            severity = "HIGH"
        }
    }
}
`, randomSuffix, orgId)
}

func testAccNetworkSecuritySecurityProfiles_antivirusOverrides(orgId string, randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_network_security_security_profile" "foobar" {
    name        = "tf-test-my-security-profile%s"
    parent      = "organizations/%s"
    location    = "global"
    description = "My security profile. Update"
    type        = "THREAT_PREVENTION"

    labels = {
        foo = "foo"
    }

    threat_prevention_profile {
        antivirus_overrides {
            action   = "ALLOW"
            protocol = "FTP"
        }

        antivirus_overrides {
            action   = "DENY"
            protocol = "HTTP"
        }

        antivirus_overrides {
            action   = "ALERT"
            protocol = "HTTP2"
        }
    }
}
`, randomSuffix, orgId)
}
