// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleClientConfig_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_client_config.current"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleClientConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttrSet(resourceName, "zone"),
					resource.TestCheckResourceAttrSet(resourceName, "access_token"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleClientConfig_omitLocation(t *testing.T) {
	t.Setenv("GOOGLE_REGION", "")
	t.Setenv("GOOGLE_ZONE", "")

	resourceName := "data.google_client_config.current"

	acctest.VcrTest(t, resource.TestCase{
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleClientConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project"),
					resource.TestCheckResourceAttrSet(resourceName, "access_token"),
				),
			},
		},
	})
}

const testAccCheckGoogleClientConfig_basic = `
data "google_client_config" "current" { }
`
