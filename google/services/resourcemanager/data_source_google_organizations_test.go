// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleOrganizations_basic(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: `data "google_organizations" "test" {}`,
				Check: resource.ComposeTestCheckFunc(
					// We assume that every principal finds at least one organization and we'll only check set-ness
					resource.TestCheckResourceAttrSet("data.google_organizations.test", "organizations.0.directory_customer_id"),
					resource.TestCheckResourceAttrSet("data.google_organizations.test", "organizations.0.display_name"),
					resource.TestCheckResourceAttrSet("data.google_organizations.test", "organizations.0.lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.google_organizations.test", "organizations.0.name"),
					resource.TestCheckResourceAttrSet("data.google_organizations.test", "organizations.0.org_id"),
				),
			},
		},
	})
}
