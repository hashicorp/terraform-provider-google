package google

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/spanner/v1"
	"strings"
)

// Unit Tests

func TestExtractInstanceConfigFromApi_withFullPath(t *testing.T) {
	actual := extractInstanceConfigFromApi("projects/project123/instanceConfigs/conf123")
	expected := "conf123"
	expectEquals(t, expected, actual)
}

func TestExtractInstanceConfigFromApi_withNoPath(t *testing.T) {
	actual := extractInstanceConfigFromApi("conf123")
	expected := "conf123"
	expectEquals(t, expected, actual)
}

func TestExtractInstanceNameFromApi_withFullPath(t *testing.T) {
	actual := extractInstanceNameFromApi("projects/project123/instances/instance123")
	expected := "instance123"
	expectEquals(t, expected, actual)
}

func TestExtractInstanceNameFromApi_withNoPath(t *testing.T) {
	actual := extractInstanceConfigFromApi("instance123")
	expected := "instance123"
	expectEquals(t, expected, actual)
}

func TestInstanceNameForApi(t *testing.T) {
	actual := instanceNameForApi("project123", "instance123")
	expected := "projects/project123/instances/instance123"
	expectEquals(t, expected, actual)
}

func TestInstanceConfigForApi(t *testing.T) {
	actual := instanceConfigForApi("project123", "conf123")
	expected := "projects/project123/instanceConfigs/conf123"
	expectEquals(t, expected, actual)
}

func TestProjectNameForApi(t *testing.T) {
	actual := projectNameForApi("project123")
	expected := "projects/project123"
	expectEquals(t, expected, actual)
}

func TestExtractSpannerInstanceImport(t *testing.T) {
	sid, e := extractSpannerInstanceImportIds("instance456")
	if e != nil {
		t.Errorf("Error should have been nil")
	}
	expectEquals(t, "", sid.Project)
	expectEquals(t, "instance456", sid.Instance)
}

func TestExtractSpannerInstanceImport_projectAndInstance(t *testing.T) {
	sid, e := extractSpannerInstanceImportIds("project123/instance456")
	if e != nil {
		t.Errorf("Error should have been nil")
	}
	expectEquals(t, "project123", sid.Project)
	expectEquals(t, "instance456", sid.Instance)
}

func TestExtractSpannerInstanceImport_invalidLeadingSlash(t *testing.T) {
	sid, e := extractSpannerInstanceImportIds("/instance456")
	expectInvalidSpannerInstanceImport(t, sid, e)
}

func TestExtractSpannerInstanceImport_invalidTrailingSlash(t *testing.T) {
	sid, e := extractSpannerInstanceImportIds("project123/")
	expectInvalidSpannerInstanceImport(t, sid, e)
}

func TestExtractSpannerInstanceImport_invalidSingleSlash(t *testing.T) {
	sid, e := extractSpannerInstanceImportIds("/")
	expectInvalidSpannerInstanceImport(t, sid, e)
}

func TestExtractSpannerInstanceImport_invalidMultiSlash(t *testing.T) {
	sid, e := extractSpannerInstanceImportIds("project123/instance456/db789")
	expectInvalidSpannerInstanceImport(t, sid, e)
}

func expectInvalidSpannerInstanceImport(t *testing.T, sid *spannerInstanceImportId, e error) {
	if sid != nil {
		t.Errorf("Expected spannerInstanceImportId to be nil")
		return
	}
	if e == nil {
		t.Errorf("Expected an Error but did not get one")
		return
	}
	if !strings.HasPrefix(e.Error(), "Invalid spanner database specifier") {
		t.Errorf("Expecting Error starting with 'Invalid spanner database specifier'")
	}
}

func expectEquals(t *testing.T, expected, actual string) {
	if actual != expected {
		t.Fatalf("Expected %s, but got %s", expected, actual)
	}
}

// Acceptance Tests

func TestAccSpannerInstance_basic(t *testing.T) {
	var instance spanner.Instance
	rnd := acctest.RandString(10)
	idName := fmt.Sprintf("spanner-test-%s", rnd)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpannerInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_basic(idName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpannerInstanceExists("google_spanner_instance.basic", &instance),

					resource.TestCheckResourceAttr("google_spanner_instance.basic", "name", idName),
					resource.TestCheckResourceAttr("google_spanner_instance.basic", "display_name", idName+"-dname"),
					resource.TestCheckResourceAttr("google_spanner_instance.basic", "num_nodes", "1"),
				),
			},
		},
	})
}

func TestAccSpannerInstance_basicWithAutogenName(t *testing.T) {
	var instance spanner.Instance
	rnd := acctest.RandString(10)
	displayName := fmt.Sprintf("spanner-test-%s-dname", rnd)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpannerInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_basicWithAutogenName(displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpannerInstanceExists("google_spanner_instance.basic", &instance),

					resource.TestCheckResourceAttr("google_spanner_instance.basic", "display_name", displayName),
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "name"),
				),
			},
		},
	})
}

func TestAccSpannerInstance_duplicateNameError(t *testing.T) {
	var instance spanner.Instance
	rnd := acctest.RandString(10)
	idName := fmt.Sprintf("spanner-test-%s", rnd)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpannerInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_duplicateNameError_part1(idName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpannerInstanceExists("google_spanner_instance.basic1", &instance),
				),
			},
			{
				Config: testAccSpannerInstance_duplicateNameError_part2(idName),
				ExpectError: regexp.MustCompile(
					fmt.Sprintf("Error, the name %s is not unique and already used", idName)),
			},
		},
	})
}

func TestAccSpannerInstance_updateDisplayNameAndNodes(t *testing.T) {
	var instance spanner.Instance
	rnd := acctest.RandString(10)
	dName1 := fmt.Sprintf("spanner-dname1-%s", rnd)
	dName2 := fmt.Sprintf("spanner-dname2-%s", rnd)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpannerInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_updateDisplayNameAndNodes(dName1, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpannerInstanceExists("google_spanner_instance.updater", &instance),
					resource.TestCheckResourceAttr("google_spanner_instance.updater", "display_name", dName1),
					resource.TestCheckResourceAttr("google_spanner_instance.updater", "num_nodes", "1"),
				),
			},
			{
				Config: testAccSpannerInstance_updateDisplayNameAndNodes(dName2, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpannerInstanceExists("google_spanner_instance.updater", &instance),
					resource.TestCheckResourceAttr("google_spanner_instance.updater", "display_name", dName2),
					resource.TestCheckResourceAttr("google_spanner_instance.updater", "num_nodes", "2"),
				),
			},
		},
	})
}

func testAccCheckSpannerInstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_spanner_instance" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Unable to verify delete of spanner instance, ID is empty")
		}

		instanceName := rs.Primary.Attributes["name"]
		project, err := getProjectId(rs, config)
		if err != nil {
			return err
		}

		_, err = config.clientSpanner.Projects.Instances.Get(
			instanceNameForApi(project, instanceName)).Do()

		if err != nil {
			if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusNotFound {
				return nil
			}
			return fmt.Errorf("Error make GCP platform call to verify spanner instance deleted: %s", err.Error())
		}
		return fmt.Errorf("Spanner instance not destroyed - still exists")
	}

	return nil
}

func getProjectId(rs *terraform.ResourceState, config *Config) (string, error) {
	res := rs.Primary.Attributes["project"]
	if res == "" {
		if config.Project != "" {
			return config.Project, nil
		}
		return "", fmt.Errorf("%q: required field is not set (or cannot be determined from provider)", "project")
	}
	return res, nil
}

func testAccCheckSpannerInstanceExists(n string, instance *spanner.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Terraform resource Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set for Spanner instance")
		}

		project, err := getProjectId(rs, config)
		if err != nil {
			return err
		}
		found, err := config.clientSpanner.Projects.Instances.Get(
			instanceNameForApi(project, rs.Primary.ID)).Do()
		if err != nil {
			return err
		}

		fName := extractInstanceNameFromApi(found.Name)
		if fName != rs.Primary.ID {
			return fmt.Errorf("Spanner instance %s not found, found %s instead", rs.Primary.ID, fName)
		}

		*instance = *found

		return nil
	}
}

func testAccSpannerInstance_basic(name string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name          = "%s"
  config        = "regional-us-central1"
  display_name  = "%s-dname"
  num_nodes     = 1
}
`, name, name)
}

func testAccSpannerInstance_basicWithProject(project, name string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  project       = "%s"
  name          = "%s"
  config        = "regional-us-central1"
  display_name  = "%s-dname"
  num_nodes     = 1
}
`, project, name, name)
}

func testAccSpannerInstance_basicWithAutogenName(name string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  config        = "regional-us-central1"
  display_name  = "%s"
  num_nodes     = 1
}
`, name)
}

func testAccSpannerInstance_duplicateNameError_part1(name string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic1" {
  name          = "%s"
  config        = "regional-us-central1"
  display_name  = "%s-dname"
  num_nodes     = 1
}

`, name, name)
}

func testAccSpannerInstance_duplicateNameError_part2(name string) string {
	return fmt.Sprintf(`
%s

resource "google_spanner_instance" "basic2" {
  name          = "%s"
  config        = "regional-us-central1"
  display_name  = "%s-dname"
  num_nodes     = 1
}
`, testAccSpannerInstance_duplicateNameError_part1(name), name, name)
}

func testAccSpannerInstance_updateDisplayNameAndNodes(name string, nodes int) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "updater" {
  config        = "regional-us-central1"
  display_name  = "%s"
  num_nodes     = %d
}
`, name, nodes)
}
