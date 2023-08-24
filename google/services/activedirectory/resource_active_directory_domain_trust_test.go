// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package activedirectory_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccActiveDirectoryDomainTrust_activeDirectoryDomainTrustBasicExample(t *testing.T) {
	// skip the test until Active Directory setup issue got resolved
	t.Skip()

	// This test continues to fail due to AD setup required
	// Skipping in VCR to allow for fully successful test runs
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckActiveDirectoryDomainTrustDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccActiveDirectoryDomainTrust_activeDirectoryDomainTrustBasicExample(context),
			},
			{
				ResourceName:            "google_active_directory_domain_trust.ad-domain-trust",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"trust_handshake_secret", "domain"},
			},
			{
				Config: testAccActiveDirectoryDomainTrust_activeDirectoryDomainTrustUpdate(context),
			},
			{
				ResourceName:            "google_active_directory_domain_trust.ad-domain-trust",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"trust_handshake_secret", "domain"},
			},
		},
	})
}

func testAccActiveDirectoryDomainTrust_activeDirectoryDomainTrustBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_active_directory_domain_trust" "ad-domain-trust" {
    domain     = "ci-managed-ad.com"
    target_domain_name = "example-gcp.com"
    target_dns_ip_addresses = ["10.1.0.100"]
    trust_direction         = "OUTBOUND"
    trust_type              = "FOREST"
    trust_handshake_secret  = "Testing1!"
}
`, context)
}

func testAccActiveDirectoryDomainTrust_activeDirectoryDomainTrustUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_active_directory_domain_trust" "ad-domain-trust" {
    domain     = "ci-managed-ad.com"
    target_domain_name = "example-gcp.com"
    target_dns_ip_addresses = ["10.2.0.100"]
    trust_direction         = "OUTBOUND"
    trust_type              = "FOREST"
    trust_handshake_secret  = "Testing1!"
}
`, context)
}

func testAccCheckActiveDirectoryDomainTrustDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_active_directory_domain_trust" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ActiveDirectoryBasePath}}projects/{{project}}/locations/global/domains/{{domain}}")
			if err != nil {
				return err
			}

			res, _ := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})

			var v interface{}
			var ok bool

			v, ok = res["trusts"]
			if ok || v != nil {
				return fmt.Errorf("ActiveDirectoryDomainTrust still exists at %s", url)
			}
		}
		return nil
	}
}
