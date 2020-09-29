package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccActiveDirectoryDomain_update(t *testing.T) {
	t.Parallel()

	domain := fmt.Sprintf("tf-test%s.org1.com", randString(t, 5))
	context := map[string]interface{}{
		"domain":        domain,
		"resource_name": "ad-domain",
	}

	resourceName := Nprintf("google_active_directory_domain.%{resource_name}", context)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
