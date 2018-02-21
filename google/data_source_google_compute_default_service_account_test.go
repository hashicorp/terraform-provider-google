package google

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGoogleComputeDefaultServiceAccount_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_compute_default_service_account.default"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleComputeDefaultServiceAccount_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

const testAccCheckGoogleComputeDefaultServiceAccount_basic = `
data "google_compute_default_service_account" "default" { }
`
