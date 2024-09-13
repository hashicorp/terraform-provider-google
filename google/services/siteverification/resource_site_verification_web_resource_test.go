// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package siteverification_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccSiteVerificationWebResource_siteVerificationDomain(t *testing.T) {
	// This test requires manual project configuration.
	acctest.SkipIfVcr(t)

	// This test needs to be able to create DNS records that are publicly
	// resolvable. To run, you'll need a registered domain with a GCP managed zone
	// that has been delegated a subdomain of the registered domain:
	//   1. Create a new GCP Cloud DNS managed zone in your test project for the subdomain
	//   2. Open the NS record that was created and note the nameservers (e.g., "ns-cloud-d1.googledomains.com")
	//   3. At your regular DNS host (i.e., not the new managed zone) add an NS record for the subdomain using the servers from step 2.
	//   4. Update the two variables below and run the test manually.
	subDomain := "subdomain.example.com"
	managedZone := "terraform-test"

	domain := "siteverification-" + acctest.RandString(t, 10) + "." + subDomain
	context := map[string]interface{}{
		"managed_zone": managedZone,
		"domain":       domain,
		"dns_name":     domain + ".", // DNS records require an FQDN
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSiteVerificationWebResourceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSiteVerificationWebResource_siteVerificationDomain(context),
			},
			{
				ResourceName:            "google_site_verification_web_resource.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"verification_method"},
			},
			{
				Config: testAccSiteVerificationWebResource_siteVerificationRemoveDomain(context),
			},
		},
	})
}

func testAccSiteVerificationWebResource_siteVerificationDomain(context map[string]interface{}) string {
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

data "google_site_verification_token" "token" {
  provider            = google.scoped
  type                = "INET_DOMAIN"
  identifier          = "%{domain}"
  verification_method = "DNS_TXT"
}

resource "google_dns_record_set" "example" {
  provider     = google.scoped
  managed_zone = "%{managed_zone}"
  name         = "%{dns_name}"
  type         = "TXT"
  rrdatas      = [data.google_site_verification_token.token.token]
  ttl          = 86400
}

resource "google_site_verification_web_resource" "example" {
  provider = google.scoped
  site {
    type       = data.google_site_verification_token.token.type
    identifier = data.google_site_verification_token.token.identifier
  }
  verification_method = data.google_site_verification_token.token.verification_method

  depends_on = [google_dns_record_set.example]
}
`, context)
}

func testAccSiteVerificationWebResource_siteVerificationRemoveDomain(context map[string]interface{}) string {
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

data "google_site_verification_token" "token" {
  provider            = google.scoped
  type                = "INET_DOMAIN"
  identifier          = "%{domain}"
  verification_method = "DNS_TXT"
}

resource "google_site_verification_web_resource" "example" {
  provider = google.scoped
  site {
    type       = data.google_site_verification_token.token.type
    identifier = data.google_site_verification_token.token.identifier
  }
  verification_method = data.google_site_verification_token.token.verification_method
}
`, context)
}

func testAccCheckSiteVerificationWebResourceDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_site_verification_web_resource" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{SiteVerificationBasePath}}webResource/{{id}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:               config,
				Method:               "GET",
				Project:              billingProject,
				RawURL:               url,
				UserAgent:            config.UserAgent,
				ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsSiteVerificationRetryableError},
			})
			if err == nil {
				return fmt.Errorf("SiteVerificationWebResource still exists at %s", url)
			}
		}

		return nil
	}
}

func TestAccSiteVerificationWebResource_siteVerificationBucket(t *testing.T) {
	t.Parallel()

	bucket := "tf-sitverification-test-" + acctest.RandString(t, 10)
	context := map[string]interface{}{
		"bucket": bucket,
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSiteVerificationWebResourceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSiteVerificationWebResource_siteVerificationBucket(context),
			},
			{
				ResourceName:            "google_site_verification_web_resource.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"verification_method"},
			},
			{
				Config: testAccSiteVerificationWebResource_siteVerificationRemoveBucket(context),
			},
		},
	})
}

func testAccSiteVerificationWebResource_siteVerificationBucket(context map[string]interface{}) string {
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

resource "google_storage_bucket" "bucket" {
  provider = google.scoped
  name     = "%{bucket}"
  location = "US"
}

data "google_site_verification_token" "token" {
  provider            = google.scoped
  type                = "SITE"
  identifier          = "https://${google_storage_bucket.bucket.name}.storage.googleapis.com/"
  verification_method = "FILE"
}

resource "google_storage_bucket_object" "object" {
  provider = google.scoped
  name     = "${data.google_site_verification_token.token.token}"
  content  = "google-site-verification: ${data.google_site_verification_token.token.token}"
  bucket   = google_storage_bucket.bucket.name
}

resource "google_storage_object_access_control" "public_rule" {
  provider = google.scoped
  bucket   = google_storage_bucket.bucket.name
  object   = google_storage_bucket_object.object.name
  role     = "READER"
  entity   = "allUsers"
}

resource "google_site_verification_web_resource" "example" {
  provider = google.scoped
  site {
    type       = data.google_site_verification_token.token.type
    identifier = data.google_site_verification_token.token.identifier
  }
  verification_method = data.google_site_verification_token.token.verification_method
}
`, context)
}

func testAccSiteVerificationWebResource_siteVerificationRemoveBucket(context map[string]interface{}) string {
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

data "google_site_verification_token" "token" {
  provider            = google.scoped
  type                = "SITE"
  identifier          = "https://%{bucket}.storage.googleapis.com/"
  verification_method = "FILE"
}

resource "google_site_verification_web_resource" "example" {
  provider = google.scoped
  site {
    type       = data.google_site_verification_token.token.type
    identifier = data.google_site_verification_token.token.identifier
  }
  verification_method = data.google_site_verification_token.token.verification_method
}
`, context)
}
