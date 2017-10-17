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

func TestExtractInstanceConfigFromUri_withFullPath(t *testing.T) {
	actual := extractInstanceConfigFromUri("projects/project123/instanceConfigs/conf987")
	expected := "conf987"
	expectEquals(t, expected, actual)
}

func TestExtractInstanceConfigFromUri_withNoPath(t *testing.T) {
	actual := extractInstanceConfigFromUri("conf987")
	expected := "conf987"
	expectEquals(t, expected, actual)
}

func TestExtractInstanceNameFromUri_withFullPath(t *testing.T) {
	actual := extractInstanceNameFromUri("projects/project123/instances/instance456")
	expected := "instance456"
	expectEquals(t, expected, actual)
}

func TestExtractInstanceNameFromUri_withNoPath(t *testing.T) {
	actual := extractInstanceConfigFromUri("instance456")
	expected := "instance456"
	expectEquals(t, expected, actual)
}

func TestSpannerInstanceId_instanceUri(t *testing.T) {
	id := spannerInstanceId{
		Project:  "project123",
		Instance: "instance456",
	}
	actual := id.instanceUri()
	expected := "projects/project123/instances/instance456"
	expectEquals(t, expected, actual)
}

func TestSpannerInstanceId_instanceConfigUri(t *testing.T) {
	id := spannerInstanceId{
		Project:  "project123",
		Instance: "instance456",
	}
	actual := id.instanceConfigUri("conf987")
	expected := "projects/project123/instanceConfigs/conf987"
	expectEquals(t, expected, actual)
}

func TestSpannerInstanceId_parentProjectUri(t *testing.T) {
	id := spannerInstanceId{
		Project:  "project123",
		Instance: "instance456",
	}
	actual := id.parentProjectUri()
	expected := "projects/project123"
	expectEquals(t, expected, actual)
}

func TestGenSpannerInstanceName(t *testing.T) {
	s := genSpannerInstanceName()
	if len(s) != 30 {
		t.Fatalf("Expected a 30 char ID to be generated, instead found %d chars", len(s))
	}
}

func TestImportSpannerInstanceId(t *testing.T) {
	sid, e := importSpannerInstanceId("instance456")
	if e != nil {
		t.Errorf("Error should have been nil")
	}
	expectEquals(t, "", sid.Project)
	expectEquals(t, "instance456", sid.Instance)
}

func TestImportSpannerInstanceId_projectAndInstance(t *testing.T) {
	sid, e := importSpannerInstanceId("project123/instance456")
	if e != nil {
		t.Errorf("Error should have been nil")
	}
	expectEquals(t, "project123", sid.Project)
	expectEquals(t, "instance456", sid.Instance)
}

func TestImportSpannerInstanceId_invalidLeadingSlash(t *testing.T) {
	sid, e := importSpannerInstanceId("/instance456")
	expectInvalidSpannerInstanceImport(t, sid, e)
}

func TestImportSpannerInstanceId_invalidTrailingSlash(t *testing.T) {
	sid, e := importSpannerInstanceId("project123/")
	expectInvalidSpannerInstanceImport(t, sid, e)
}

func TestImportSpannerInstanceId_invalidSingleSlash(t *testing.T) {
	sid, e := importSpannerInstanceId("/")
	expectInvalidSpannerInstanceImport(t, sid, e)
}

func TestImportSpannerInstanceId_invalidMultiSlash(t *testing.T) {
	sid, e := importSpannerInstanceId("project123/instance456/db789")
	expectInvalidSpannerInstanceImport(t, sid, e)
}

func expectInvalidSpannerInstanceImport(t *testing.T, sid *spannerInstanceId, e error) {
	if sid != nil {
		t.Errorf("Expected spannerInstanceId to be nil")
		return
	}
	if e == nil {
		t.Errorf("Expected an Error but did not get one")
		return
	}
	if !strings.HasPrefix(e.Error(), "Invalid spanner instance specifier") {
		t.Errorf("Expecting Error starting with 'Invalid spanner instance specifier'")
	}
}

func expectEquals(t *testing.T, expected, actual string) {
	if actual != expected {
		t.Fatalf("Expected %s, but got %s", expected, actual)
	}
}

// Acceptance Tests

func TestAccSpannerInstance_basic(t *testing.T) {
	t.Parallel()

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
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
		},
	})
}

func TestAccSpannerInstance_basicWithAutogenName(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
					fmt.Sprintf("Error, the name %s is not unique within project", idName)),
			},
		},
	})
}

func TestAccSpannerInstance_update(t *testing.T) {
	t.Parallel()

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
				Config: testAccSpannerInstance_update(dName1, 1, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpannerInstanceExists("google_spanner_instance.updater", &instance),
					resource.TestCheckResourceAttr("google_spanner_instance.updater", "display_name", dName1),
					resource.TestCheckResourceAttr("google_spanner_instance.updater", "num_nodes", "1"),
					resource.TestCheckResourceAttr("google_spanner_instance.updater", "labels.%", "1"),
				),
			},
			{
				Config: testAccSpannerInstance_update(dName2, 2, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSpannerInstanceExists("google_spanner_instance.updater", &instance),
					resource.TestCheckResourceAttr("google_spanner_instance.updater", "display_name", dName2),
					resource.TestCheckResourceAttr("google_spanner_instance.updater", "num_nodes", "2"),
					resource.TestCheckResourceAttr("google_spanner_instance.updater", "labels.%", "2"),
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
		project, err := getTestProject(rs.Primary, config)
		if err != nil {
			return err
		}

		id := spannerInstanceId{
			Project:  project,
			Instance: instanceName,
		}
		_, err = config.clientSpanner.Projects.Instances.Get(
			id.instanceUri()).Do()

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

		id, err := extractSpannerInstanceId(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := config.clientSpanner.Projects.Instances.Get(
			id.instanceUri()).Do()
		if err != nil {
			return err
		}

		fName := extractInstanceNameFromUri(found.Name)
		if fName != extractInstanceNameFromUri(rs.Primary.ID) {
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

func testAccSpannerInstance_update(name string, nodes int, addLabel bool) string {

	extraLabel := ""
	if addLabel {
		extraLabel = "\"key2\" = \"value2\""
	}
	return fmt.Sprintf(`
resource "google_spanner_instance" "updater" {
  config        = "regional-us-central1"
  display_name  = "%s"
  num_nodes     = %d

  labels {
     "key1" = "value1"
     %s
  }
}
`, name, nodes, extraLabel)
}
