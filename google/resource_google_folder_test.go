package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	resourceManagerV2Beta1 "google.golang.org/api/cloudresourcemanager/v2beta1"
)

func TestAccFolder_rename(t *testing.T) {
	t.Parallel()

	folderDisplayName := "tf-test-" + acctest.RandString(10)
	newFolderDisplayName := "tf-test-renamed-" + acctest.RandString(10)
	org := getTestOrgFromEnv(t)
	parent := "organizations/" + org
	folder := resourceManagerV2Beta1.Folder{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFolder_basic(folderDisplayName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderExists("google_folder.folder1", &folder),
					testAccCheckGoogleFolderParent(&folder, parent),
					testAccCheckGoogleFolderDisplayName(&folder, folderDisplayName),
				),
			},
			{
				Config: testAccFolder_basic(newFolderDisplayName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderExists("google_folder.folder1", &folder),
					testAccCheckGoogleFolderParent(&folder, parent),
					testAccCheckGoogleFolderDisplayName(&folder, newFolderDisplayName),
				)},
			{
				ResourceName:      "google_folder.folder1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccFolder_moveParent(t *testing.T) {
	t.Parallel()

	folder1DisplayName := "tf-test-" + acctest.RandString(10)
	folder2DisplayName := "tf-test-" + acctest.RandString(10)
	org := getTestOrgFromEnv(t)
	parent := "organizations/" + org
	folder1 := resourceManagerV2Beta1.Folder{}
	folder2 := resourceManagerV2Beta1.Folder{}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleFolderDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFolder_basic(folder1DisplayName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderExists("google_folder.folder1", &folder1),
					testAccCheckGoogleFolderParent(&folder1, parent),
					testAccCheckGoogleFolderDisplayName(&folder1, folder1DisplayName),
				),
			},
			{
				Config: testAccFolder_move(folder1DisplayName, folder2DisplayName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderExists("google_folder.folder1", &folder1),
					testAccCheckGoogleFolderDisplayName(&folder1, folder1DisplayName),
					testAccCheckGoogleFolderExists("google_folder.folder2", &folder2),
					testAccCheckGoogleFolderParent(&folder2, parent),
					testAccCheckGoogleFolderDisplayName(&folder2, folder2DisplayName),
				),
			},
		},
	})
}

func testAccCheckGoogleFolderDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_folder" {
			continue
		}

		folder, err := config.clientResourceManagerV2Beta1.Folders.Get(rs.Primary.ID).Do()
		if err != nil || folder.LifecycleState != "DELETE_REQUESTED" {
			return fmt.Errorf("Folder '%s' hasn't been marked for deletion", rs.Primary.Attributes["display_name"])
		}
	}

	return nil
}

func testAccCheckGoogleFolderExists(n string, folder *resourceManagerV2Beta1.Folder) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientResourceManagerV2Beta1.Folders.Get(rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		*folder = *found

		return nil
	}
}

func testAccCheckGoogleFolderDisplayName(folder *resourceManagerV2Beta1.Folder, displayName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if folder.DisplayName != displayName {
			return fmt.Errorf("Incorrect display name . Expected '%s', got '%s'", displayName, folder.DisplayName)
		}
		return nil
	}
}

func testAccCheckGoogleFolderParent(folder *resourceManagerV2Beta1.Folder, parent string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if folder.Parent != parent {
			return fmt.Errorf("Incorrect parent. Expected '%s', got '%s'", parent, folder.Parent)
		}
		return nil
	}
}

func testAccFolder_basic(folder, parent string) string {
	return fmt.Sprintf(`
resource "google_folder" "folder1" {
  display_name = "%s"
  parent       = "%s"
}
`, folder, parent)
}

func testAccFolder_move(folder1, folder2, parent string) string {
	return fmt.Sprintf(`
resource "google_folder" "folder1" {
  display_name = "%s"
  parent       = google_folder.folder2.name
}

resource "google_folder" "folder2" {
  display_name = "%s"
  parent       = "%s"
}
`, folder1, folder2, parent)
}
