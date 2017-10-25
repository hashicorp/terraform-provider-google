package google

import (
	"fmt"
	"testing"

	"os"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccGoogleServiceAccount_importBasic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccGoogleServiceAccount_import("terraform-" + acctest.RandString(10)),
			},

			resource.TestStep{
				ResourceName:      "google_service_account.acceptance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGoogleServiceAccount_import(saName string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
    account_id = "%s"
    display_name = "%s"
}`, saName, saName)
}

func TestAccGoogleServiceAccount_importWithProject(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccGoogleServiceAccount_importWithProject(os.Getenv("GOOGLE_PROJECT"), "terraform-"+acctest.RandString(10)),
			},

			resource.TestStep{
				ResourceName:      "google_service_account.acceptance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGoogleServiceAccount_importWithProject(project, saName string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
    project = "%s"
    account_id = "%s"
    display_name = "%s"
}`, project, saName, saName)
}
