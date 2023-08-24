// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sql_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleSqlTiers_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_sql_tiers.default"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleSqlTiers_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project"),
					resource.TestCheckResourceAttrSet(resourceName, "tiers.0.%"),
					resource.TestCheckResourceAttrSet(resourceName, "tiers.0.tier"),
					resource.TestCheckResourceAttrSet(resourceName, "tiers.0.ram"),
					resource.TestCheckResourceAttrSet(resourceName, "tiers.0.disk_quota"),
					resource.TestCheckResourceAttrSet(resourceName, "tiers.0.region.0"),
				),
			},
		},
	})
}

const testAccCheckGoogleSqlTiers_basic = `
data "google_sql_tiers" "default" {
}
`
