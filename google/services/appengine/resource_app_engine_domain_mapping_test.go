// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package appengine_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccAppEngineDomainMapping_update(t *testing.T) {
	t.Parallel()

	domainName := fmt.Sprintf("tf-test-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAppEngineDomainMappingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAppEngineDomainMapping_basic(domainName),
			},
			{
				ResourceName:            "google_app_engine_domain_mapping.domain_mapping",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"override_strategy"},
			},
			{
				Config: testAccAppEngineDomainMapping_update(domainName),
			},
			{
				ResourceName:            "google_app_engine_domain_mapping.domain_mapping",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"override_strategy"},
			},
		},
	})
}

func testAccAppEngineDomainMapping_basic(domainName string) string {
	return fmt.Sprintf(`
resource "google_app_engine_domain_mapping" "domain_mapping" {
  domain_name = "%s.gcp.tfacc.hashicorptest.com"

  ssl_settings {
    ssl_management_type = "AUTOMATIC"
  }
}
`, domainName)
}

func testAccAppEngineDomainMapping_update(domainName string) string {
	return fmt.Sprintf(`
resource "google_app_engine_domain_mapping" "domain_mapping" {
  domain_name = "%s.gcp.tfacc.hashicorptest.com"

  ssl_settings {
    certificate_id      = ""
    ssl_management_type = "MANUAL"
  }
}
`, domainName)
}
