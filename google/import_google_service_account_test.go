package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"os"
)

func TestAccGoogleServiceAccount_importBasic(t *testing.T) {
	t.Parallel()

	resourceName := "google_service_account.acceptance"
	project := os.Getenv("GOOGLE_PROJECT")
	sa_name := "terraform-" + acctest.RandString(10)
	conf := testAccGoogleServiceAccount_import(project, sa_name)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: conf,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGoogleServiceAccount_import(project, sa_name string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
    project = "%s"
    account_id = "%s"
    display_name = "%s"
}`, project, sa_name, sa_name)
}
