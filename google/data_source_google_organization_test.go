package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGoogleOrganization_basic(t *testing.T) {
	orgId := getTestOrgFromEnv(t)
	name := "organizations/" + orgId

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleOrganization_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_organization.org", "id", orgId),
					resource.TestCheckResourceAttr("data.google_organization.org", "name", name),
				),
			},
		},
	})
}

func testAccCheckGoogleOrganization_basic(name string) string {
	return fmt.Sprintf(`
data "google_organization" "org" {
	name = "%s"
}`, name)
}
