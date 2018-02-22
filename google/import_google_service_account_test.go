package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccServiceAccount_importBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccServiceAccount_import("terraform-" + acctest.RandString(10)),
			},

			resource.TestStep{
				ResourceName:      "google_service_account.acceptance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccServiceAccount_import(saName string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
    account_id = "%s"
    display_name = "%s"
}`, saName, saName)
}

func TestAccServiceAccount_importWithProject(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccServiceAccount_importWithProject(getTestProjectFromEnv(), "terraform-"+acctest.RandString(10)),
			},

			resource.TestStep{
				ResourceName:      "google_service_account.acceptance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccServiceAccount_importWithProject(project, saName string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
    project = "%s"
    account_id = "%s"
    display_name = "%s"
}`, project, saName, saName)
}
