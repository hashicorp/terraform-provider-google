package google

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"strings"

	"google.golang.org/api/googleapi"
)

// Unit Tests

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

func TestImportSpannerInstanceId_projectId(t *testing.T) {
	shouldPass := []string{
		"project-id/instance",
		"123123/instance",
		"hashicorptest.net:project-123/instance",
		"123/456",
	}

	shouldFail := []string{
		"project-id#/instance",
		"project-id/instance#",
		"hashicorptest.net:project-123:invalid:project/instance",
		"hashicorptest.net:/instance",
	}

	for _, element := range shouldPass {
		_, e := importSpannerInstanceId(element)
		if e != nil {
			t.Error("importSpannerInstanceId should pass on '" + element + "' but doesn't")
		}
	}

	for _, element := range shouldFail {
		_, e := importSpannerInstanceId(element)
		if e == nil {
			t.Error("importSpannerInstanceId should fail on '" + element + "' but doesn't")
		}
	}
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

	idName := fmt.Sprintf("spanner-test-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpannerInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_basic(idName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "state"),
				),
			},
			{
				ResourceName:      "google_spanner_instance.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerInstance_basicWithAutogenName(t *testing.T) {
	t.Parallel()

	displayName := fmt.Sprintf("spanner-test-%s-dname", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpannerInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_basicWithAutogenName(displayName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_instance.basic", "name"),
				),
			},
			{
				ResourceName:      "google_spanner_instance.basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerInstance_update(t *testing.T) {
	t.Parallel()

	dName1 := fmt.Sprintf("spanner-dname1-%s", acctest.RandString(10))
	dName2 := fmt.Sprintf("spanner-dname2-%s", acctest.RandString(10))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpannerInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstance_update(dName1, 1, false),
			},
			{
				ResourceName:      "google_spanner_instance.updater",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSpannerInstance_update(dName2, 2, true),
			},
			{
				ResourceName:      "google_spanner_instance.updater",
				ImportState:       true,
				ImportStateVerify: true,
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

		if err == nil {
			return fmt.Errorf("Spanner instance still exists")
		}

		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == http.StatusNotFound {
			return nil
		}
		return errwrap.Wrapf("Error verifying spanner instance deleted: {{err}}", err)
	}

	return nil
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

func testAccSpannerInstance_basicWithAutogenName(name string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  config        = "regional-us-central1"
  display_name  = "%s"
  num_nodes     = 1
}
`, name)
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
