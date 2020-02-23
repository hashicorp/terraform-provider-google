package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceGoogleClientConfig_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_client_config.current"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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

const testAccCheckGoogleClientConfig_basic = `
data "google_client_config" "current" { }
`
