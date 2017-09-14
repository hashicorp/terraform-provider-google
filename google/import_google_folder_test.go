package google

import (
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"os"
	"testing"
)

func TestAccGoogleFolder_import(t *testing.T) {
	skipIfEnvNotSet(t, "GOOGLE_ORG")

	folderDisplayName := "tf-test-" + acctest.RandString(10)
	org := os.Getenv("GOOGLE_ORG")
	parent := "organizations/" + org

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleFolderDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccGoogleFolder_basic(folderDisplayName, parent),
			},
			resource.TestStep{
				ResourceName:      "google_folder.folder1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
