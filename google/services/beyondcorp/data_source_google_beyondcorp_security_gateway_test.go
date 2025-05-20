// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package beyondcorp_test

import (
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleBeyondcorpSecurityGateway_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBeyondcorpSecurityGatewayDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBeyondcorpSecurityGateway_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_beyondcorp_security_gateway.foo", "google_beyondcorp_security_gateway.foo"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleBeyondcorpSecurityGateway_full(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBeyondcorpSecurityGatewayDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBeyondcorpSecurityGateway_full(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_beyondcorp_security_gateway.foo", "google_beyondcorp_security_gateway.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleBeyondcorpSecurityGateway_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_beyondcorp_security_gateway" "foo" {
  security_gateway_id = "default-foo-sg-basic-%{random_suffix}"
  display_name = "My Security Gateway resource"
  hubs { region = "us-central1" }
}

data "google_beyondcorp_security_gateway" "foo" {
	security_gateway_id = google_beyondcorp_security_gateway.foo.security_gateway_id
}
`, context)
}

func testAccDataSourceGoogleBeyondcorpSecurityGateway_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_beyondcorp_security_gateway" "foo" {
  security_gateway_id = "default-foo-sg-full-%{random_suffix}"
  display_name = "My Security Gateway resource"
  hubs { region = "us-central1" }
}

data "google_beyondcorp_security_gateway" "foo" {
	security_gateway_id = google_beyondcorp_security_gateway.foo.security_gateway_id
	project = google_beyondcorp_security_gateway.foo.project
}
`, context)
}
