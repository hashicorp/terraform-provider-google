package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	resourceManagerV3 "google.golang.org/api/cloudresourcemanager/v3"
)

func TestAccFolder_rename(t *testing.T) {
	t.Parallel()

	folderDisplayName := "tf-test-" + randString(t, 10)
	newFolderDisplayName := "tf-test-renamed-" + randString(t, 10)
	org := getTestOrgFromEnv(t)
	parent := "organizations/" + org
	folder := resourceManagerV3.Folder{}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleFolderDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFolder_basic(folderDisplayName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderExists(t, "google_folder.folder1", &folder),
					testAccCheckGoogleFolderParent(&folder, parent),
					testAccCheckGoogleFolderDisplayName(&folder, folderDisplayName),
				),
			},
			{
				Config: testAccFolder_basic(newFolderDisplayName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderExists(t, "google_folder.folder1", &folder),
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

	folder1DisplayName := "tf-test-" + randString(t, 10)
	folder2DisplayName := "tf-test-" + randString(t, 10)
	org := getTestOrgFromEnv(t)
	parent := "organizations/" + org
	folder1 := resourceManagerV3.Folder{}
	folder2 := resourceManagerV3.Folder{}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGoogleFolderDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFolder_basic(folder1DisplayName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderExists(t, "google_folder.folder1", &folder1),
					testAccCheckGoogleFolderParent(&folder1, parent),
					testAccCheckGoogleFolderDisplayName(&folder1, folder1DisplayName),
				),
			},
			{
				Config: testAccFolder_move(folder1DisplayName, folder2DisplayName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleFolderExists(t, "google_folder.folder1", &folder1),
					testAccCheckGoogleFolderDisplayName(&folder1, folder1DisplayName),
					testAccCheckGoogleFolderExists(t, "google_folder.folder2", &folder2),
					testAccCheckGoogleFolderParent(&folder2, parent),
					testAccCheckGoogleFolderDisplayName(&folder2, folder2DisplayName),
				),
			},
		},
	})
}

func testAccCheckGoogleFolderDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_folder" {
				continue
			}

			folder, err := config.NewResourceManagerV3Client(config.userAgent).Folders.Get(rs.Primary.ID).Do()
			if err != nil || folder.State != "DELETE_REQUESTED" {
				return fmt.Errorf("Folder '%s' hasn't been marked for deletion", rs.Primary.Attributes["display_name"])
			}
		}

		return nil
	}
}

func testAccCheckGoogleFolderExists(t *testing.T, n string, folder *resourceManagerV3.Folder) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := googleProviderConfig(t)

		found, err := config.NewResourceManagerV3Client(config.userAgent).Folders.Get(rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		*folder = *found

		return nil
	}
}

func testAccCheckGoogleFolderDisplayName(folder *resourceManagerV3.Folder, displayName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if folder.DisplayName != displayName {
			return fmt.Errorf("Incorrect display name . Expected '%s', got '%s'", displayName, folder.DisplayName)
		}
		return nil
	}
}

func testAccCheckGoogleFolderParent(folder *resourceManagerV3.Folder, parent string) resource.TestCheckFunc {
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
