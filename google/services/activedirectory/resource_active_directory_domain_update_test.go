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

func TestAccActiveDirectoryDomain_update(t *testing.T) {
	// skip the test until Active Directory setup issue got resolved
	t.Skip()

	t.Parallel()

	domain := fmt.Sprintf("tf-test%s.org1.com", acctest.RandString(t, 5))
	context := map[string]interface{}{
		"domain":        domain,
		"resource_name": "ad-domain",
	}

	resourceName := acctest.Nprintf("google_active_directory_domain.%{resource_name}", context)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckActiveDirectoryDomainDestroyProducer(t),
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

	return acctest.Nprintf(`
	resource "google_active_directory_domain" "%{resource_name}" {
	  domain_name       = "%{domain}"
	  locations         = ["us-central1"]
	  reserved_ip_range = "192.168.255.0/24" 
	}
	`, context)
}

func testAccADDomainUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
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

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{ActiveDirectoryBasePath}}{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("ActiveDirectoryDomain still exists at %s", url)
			}
		}

		return nil
	}
}
