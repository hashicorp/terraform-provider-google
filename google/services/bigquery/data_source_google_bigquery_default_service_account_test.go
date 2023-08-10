// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigquery_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleBigqueryDefaultServiceAccount_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_bigquery_default_service_account.bq_account"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleBigqueryDefaultServiceAccount_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "email"),
					resource.TestCheckResourceAttrSet(resourceName, "member"),
				),
			},
		},
	})
}

const testAccCheckGoogleBigqueryDefaultServiceAccount_basic = `
data "google_bigquery_default_service_account" "bq_account" {
}
`
