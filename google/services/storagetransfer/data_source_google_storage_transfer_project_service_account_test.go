// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storagetransfer_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleStorageTransferProjectServiceAccount_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_storage_transfer_project_service_account.default"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleStorageTransferProjectServiceAccount_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "email"),
					resource.TestCheckResourceAttrSet(resourceName, "subject_id"),
					resource.TestCheckResourceAttrSet(resourceName, "member"),
				),
			},
		},
	})
}

const testAccCheckGoogleStorageTransferProjectServiceAccount_basic = `
data "google_storage_transfer_project_service_account" "default" { }
`
