// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package siteverification_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccSiteVerificationToken_siteverificationTokenSite(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"site": "https://www.example.com",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSiteVerificationToken_siteverificationTokenSite(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_site_verification_token.site_meta", "token"),
					resource.TestCheckResourceAttr("data.google_site_verification_token.site_meta", "type", "SITE"),
					resource.TestCheckResourceAttr("data.google_site_verification_token.site_meta", "identifier", context["site"].(string)),
					resource.TestCheckResourceAttr("data.google_site_verification_token.site_meta", "verification_method", "META"),
				),
			},
		},
	})
}

func testAccSiteVerificationToken_siteverificationTokenSite(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  alias                 = "scoped"
  user_project_override = true
  scopes = [
    "https://www.googleapis.com/auth/siteverification",
    "https://www.googleapis.com/auth/cloud-platform",
    "https://www.googleapis.com/auth/userinfo.email",
  ]
}

data "google_site_verification_token" "site_meta" {
  provider            = google.scoped
  type                = "SITE"
  identifier          = "%{site}"
  verification_method = "META"
}
`, context)
}

func TestAccSiteVerificationToken_siteverificationTokenDomain(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"domain": "www.example.com",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccSiteVerificationToken_siteverificationTokenDomain(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_site_verification_token.dns_text", "token"),
					resource.TestCheckResourceAttr("data.google_site_verification_token.dns_text", "type", "INET_DOMAIN"),
					resource.TestCheckResourceAttr("data.google_site_verification_token.dns_text", "identifier", context["domain"].(string)),
					resource.TestCheckResourceAttr("data.google_site_verification_token.dns_text", "verification_method", "DNS_TXT"),
				),
			},
		},
	})
}

func testAccSiteVerificationToken_siteverificationTokenDomain(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  alias                 = "scoped"
  user_project_override = true
  scopes = [
    "https://www.googleapis.com/auth/siteverification",
    "https://www.googleapis.com/auth/cloud-platform",
    "https://www.googleapis.com/auth/userinfo.email",
  ]
}

data "google_site_verification_token" "dns_text" {
  provider            = google.scoped
  type                = "INET_DOMAIN"
  identifier          = "%{domain}"
  verification_method = "DNS_TXT"
}
`, context)
}
