package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccActiveDirectoryDomain_update(t *testing.T) {
	// skip the test until Active Directory setup issue got resolved
	t.Skip()

	t.Parallel()

	domain := fmt.Sprintf("tf-test%s.org1.com", RandString(t, 5))
	context := map[string]interface{}{
		"domain":        domain,
		"resource_name": "ad-domain",
	}

	resourceName := Nprintf("google_active_directory_domain.%{resource_name}", context)

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckActiveDirectoryDomainDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccADDomainBasic(context),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"domain_name"},
			},
			{
				Config: testAccADDomainUpdate(context),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"domain_name"},
			},
			{
				Config: testAccADDomainBasic(context),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"domain_name"},
			},
		},
	})
}

func testAccADDomainBasic(context map[string]interface{}) string {

	return Nprintf(`
	resource "google_active_directory_domain" "%{resource_name}" {
	  domain_name       = "%{domain}"
	  locations         = ["us-central1"]
	  reserved_ip_range = "192.168.255.0/24" 
	}
	`, context)
}

func testAccADDomainUpdate(context map[string]interface{}) string {
	return Nprintf(`
	resource "google_active_directory_domain" "%{resource_name}" {
	  domain_name       = "%{domain}"	
	  locations         = ["us-central1", "us-west1"]
	  reserved_ip_range = "192.168.255.0/24" 
	  labels = {
		  env = "test"
	  }
	}
	`, context)

}

func testAccCheckActiveDirectoryDomainDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_active_directory_domain" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := GoogleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{ActiveDirectoryBasePath}}{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = SendRequest(config, "GET", billingProject, url, config.UserAgent, nil)
			if err == nil {
				return fmt.Errorf("ActiveDirectoryDomain still exists at %s", url)
			}
		}

		return nil
	}
}
