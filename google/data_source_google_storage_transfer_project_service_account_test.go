package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleStorageTransferProjectServiceAccount_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_storage_transfer_project_service_account.default"

	VcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: TestAccProviders,
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
