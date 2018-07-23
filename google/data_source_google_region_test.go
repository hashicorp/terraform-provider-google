package google

import (
	"testing"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"fmt"
	"os"
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

func TestAccRegion_fromGoogleRegionEnvVar(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckProviderRegionConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleRegionCheck("data.google_region.current"),
					testAccItReturnsTheRegionSetByTheGoogleRegionEnvVar("data.google_region.current"),
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

func testAccItReturnsTheRegionSetByTheGoogleRegionEnvVar(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		currentRegion := os.Getenv("GOOGLE_REGION")

		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		if currentRegion == "" {
			return fmt.Errorf("the environment variable GOOGLE_REGION must be set to something")
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		if rs.Primary.ID != currentRegion {
			return fmt.Errorf("resource ID was meant to be %s but was %s instead", currentRegion, rs.Primary.ID)
		}

		return nil
	}
}

var testAccCheckProviderRegionConfig = `
data "google_region" "current" {}
`
