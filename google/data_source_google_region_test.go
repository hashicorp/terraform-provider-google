package google

import (
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"fmt"
)

func TestAccRegion_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckProviderRegionConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleRegionCheck("data.google_region.current"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleRegionCheck(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		return nil
	}
}

var testAccCheckProviderRegionConfig = `
data "google_region" "current" {}
`
