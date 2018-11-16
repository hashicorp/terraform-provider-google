package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDnsManagedZone_update(t *testing.T) {
	t.Parallel()

	zoneSuffix := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsManagedZoneDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDnsManagedZone_basic(zoneSuffix, "description1"),
			},
			resource.TestStep{
				ResourceName:      "google_dns_managed_zone.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				Config: testAccDnsManagedZone_basic(zoneSuffix, "description2"),
			},
			resource.TestStep{
				ResourceName:      "google_dns_managed_zone.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDnsManagedZone_basic(suffix, description string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "foobar" {
	name = "mzone-test-%s"
	dns_name = "tf-acctest-%s.hashicorptest.com."
	description = "%s"
	labels = {
		foo = "bar"
	}
}`, suffix, suffix, description)
}

func TestDnsManagedZoneImport_parseImportId(t *testing.T) {
	zoneRegexes := []string{
		"projects/(?P<project>[^/]+)/managedZones/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/managedZones/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}

	cases := map[string]struct {
		ImportId             string
		IdRegexes            []string
		Config               *Config
		ExpectedSchemaValues map[string]interface{}
		ExpectError          bool
	}{
		"full self_link": {
			IdRegexes: zoneRegexes,
			ImportId:  "https://www.googleapis.com/dns/v1/projects/my-project/managedZones/my-zone",
			ExpectedSchemaValues: map[string]interface{}{
				"project": "my-project",
				"name":    "my-zone",
			},
		},
		"relative self_link": {
			IdRegexes: zoneRegexes,
			ImportId:  "projects/my-project/managedZones/my-zone",
			ExpectedSchemaValues: map[string]interface{}{
				"project": "my-project",
				"name":    "my-zone",
			},
		},
		"short id": {
			IdRegexes: zoneRegexes,
			ImportId:  "my-project/managedZones/my-zone",
			ExpectedSchemaValues: map[string]interface{}{
				"project": "my-project",
				"name":    "my-zone",
			},
		},
		"short id with default project and region": {
			IdRegexes: zoneRegexes,
			ImportId:  "my-zone",
			Config: &Config{
				Project: "default-project",
			},
			ExpectedSchemaValues: map[string]interface{}{
				"project": "default-project",
				"name":    "my-zone",
			},
		},
	}

	for tn, tc := range cases {
		d := &ResourceDataMock{
			FieldsInSchema: make(map[string]interface{}),
			id:             tc.ImportId,
		}
		config := tc.Config
		if config == nil {
			config = &Config{}
		}
		//
		if err := parseImportId(tc.IdRegexes, d, config); err == nil {
			for k, expectedValue := range tc.ExpectedSchemaValues {
				if v, ok := d.GetOk(k); ok {
					if v != expectedValue {
						t.Errorf("%s failed; Expected value %q for field %q, got %q", tn, expectedValue, k, v)
					}
				} else {
					t.Errorf("%s failed; Expected a value for field %q", tn, k)
				}
			}
		} else if !tc.ExpectError {
			t.Errorf("%s failed; unexpected error: %s", tn, err)
		}
	}
}

func TestAccDnsManagedZone_importWithProject(t *testing.T) {
	t.Parallel()

	zoneSuffix := acctest.RandString(10)
	project := getTestProjectFromEnv()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsManagedZoneDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDnsManagedZone_basicWithProject(zoneSuffix, "description1", project),
			},
			resource.TestStep{
				ResourceName:      "google_dns_managed_zone.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDnsManagedZone_basicWithProject(suffix, description, project string) string {
	return fmt.Sprintf(`
resource "google_dns_managed_zone" "foobar" {
	name = "mzone-test-%s"
	dns_name = "tf-acctest-%s.hashicorptest.com."
	description = "%s"
	project = "%s"
}`, suffix, suffix, description, project)
}
