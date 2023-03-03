package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccActiveDirectoryDomainTrust_activeDirectoryDomainTrustBasicExample(t *testing.T) {
	// skip the test until Active Directory setup issue got resolved
	t.Skip()

	// This test continues to fail due to AD setup required
	// Skipping in VCR to allow for fully successful test runs
	SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckActiveDirectoryDomainTrustDestroyProducer(t),
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
	return Nprintf(`
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
	return Nprintf(`
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

			config := GoogleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{ActiveDirectoryBasePath}}projects/{{project}}/locations/global/domains/{{domain}}")
			if err != nil {
				return err
			}

			res, _ := SendRequest(config, "GET", "", url, config.UserAgent, nil)

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
