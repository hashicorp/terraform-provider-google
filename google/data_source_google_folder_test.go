package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceGoogleFolder_byFullName(t *testing.T) {
	folderId := getTestFolderFromEnv(t)
	name := "folders/" + folderId

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleFolder_byName(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_folder.folder", "id", folderId),
					resource.TestCheckResourceAttr("data.google_folder.folder", "name", name),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleFolder_byShortName(t *testing.T) {
	folderId := getTestFolderFromEnv(t)
	name := "folders/" + folderId

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleFolder_byName(folderId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_folder.folder", "id", folderId),
					resource.TestCheckResourceAttr("data.google_folder.folder", "name", name),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleFolder_lookupOrganization(t *testing.T) {
	orgId := getTestOrgFromEnv(t)
	orgName := "organizations/" + orgId
	folderId := getTestFolderFromEnv(t)
	name := "folders/" + folderId

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleFolder_lookupOrganization(folderId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.google_folder.folder", "id", folderId),
					resource.TestCheckResourceAttr("data.google_folder.folder", "name", name),
					resource.TestCheckResourceAttr("data.google_folder.folder", "organization", orgName),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleFolder_byFullNameNotFound(t *testing.T) {
	name := "folders/" + acctest.RandString(16)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckGoogleFolder_byName(name),
				ExpectError: regexp.MustCompile("Folder Not Found : " + name),
			},
		},
	})
}

func TestAccDataSourceGoogleFolder_attributesCheck(t *testing.T) {
	org := getTestOrgFromEnv(t)

	parent := fmt.Sprintf("organizations/%s", org)
	displayName := "terraform-test-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckGoogleFolder_attributesCheckConfig(parent, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGoogleFolderCheck("data.google_folder.folder", "google_folder.foobar"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleFolderCheck(data_source_name string, resource_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[data_source_name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", data_source_name)
		}

		rs, ok := s.RootModule().Resources[resource_name]
		if !ok {
			return fmt.Errorf("can't find %s in state", resource_name)
		}

		ds_attr := ds.Primary.Attributes
		rs_attr := rs.Primary.Attributes
		folder_attrs_to_test := []string{"parent", "display_name", "name"}

		for _, attr_to_check := range folder_attrs_to_test {
			if ds_attr[attr_to_check] != rs_attr[attr_to_check] {
				return fmt.Errorf(
					"%s is %s; want %s",
					attr_to_check,
					ds_attr[attr_to_check],
					rs_attr[attr_to_check],
				)
			}
		}
		return nil
	}
}

func testAccCheckGoogleFolder_byName(name string) string {
	return fmt.Sprintf(`
data "google_folder" "folder" {
  folder = "%s"
}`, name)
}

func testAccCheckGoogleFolder_lookupOrganization(name string) string {
	return fmt.Sprintf(`
data "google_folder" "folder" {
  folder = "%s"
  lookup_organization = true
}`, name)
}

func testAccCheckGoogleFolder_attributesCheckConfig(parent string, displayName string) string {
	return fmt.Sprintf(`
resource "google_folder" "foobar" {
  parent = "%s"
  display_name = "%s"
}

data "google_folder" "folder" {
  folder = "${google_folder.foobar.name}"
}`, parent, displayName)
}
