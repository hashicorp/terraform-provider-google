package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceGoogleStorageProjectServiceAccount_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_storage_project_service_account.gcs_account"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleStorageProjectServiceAccount_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "email_address"),
				),
			},
		},
	})
}

const testAccCheckGoogleStorageProjectServiceAccount_basic = `
data "google_storage_project_service_account" "gcs_account" {
}
`
