package google

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGoogleClientConfig_basic(t *testing.T) {
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
				),
			},
		},
	})
}

const testAccCheckGoogleClientConfig_basic = `
data "google_client_config" "current" { }
`
