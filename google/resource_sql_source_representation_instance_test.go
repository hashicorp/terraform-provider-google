package google

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("gcp_sql_source_representation_instance", &resource.Sweeper{
		Name: "gcp_sql_source_representation_instance",
		F:    testSweepSourceRepresentationInstanceDatabases,
	})
}

func testSweepSourceRepresentationInstanceDatabases(region string) error {
	config, err := sharedConfigForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting shared config for region: %s", err)
	}

	err = config.LoadAndValidate()
	if err != nil {
		log.Fatalf("error loading: %s", err)
	}

	found, err := config.clientSqlAdmin.Instances.List(config.Project).Do()
	if err != nil {
		log.Fatalf("error listing databases: %s", err)
	}

	if len(found.Items) == 0 {
		log.Printf("No databases found")
		return nil
	}

	for _, d := range found.Items {
		if !strings.HasPrefix(d.Name, "tf-test-source-representation-instance-") {
			continue
		}

		log.Printf("Destroying Source Representation Instance (%s)", d.Name)

		// destroy instances, replicas first
		op, err := config.clientSqlAdmin.Instances.Delete(config.Project, d.Name).Do()
		if err != nil {
			if strings.Contains(err.Error(), "409") {
				// the GCP api can return a 409 error after the delete operation
				// reaches a successful end
				log.Printf("Operation not found, got 409 response")
				continue
			}

			return fmt.Errorf("Error, failed to delete source representation instance %s: %s", d.Name, err)
		}

		err = sqladminOperationWait(config, op, config.Project, "Delete Source Representation Instance")
		if err != nil {
			if strings.Contains(err.Error(), "does not exist") {
				log.Printf("SQL Source Representation Instance not found")
				continue
			}
			return err
		}
	}

	return nil
}

func TestAccSqlSourceRepresentationInstance_minimal(t *testing.T) {
	t.Parallel()

	databaseName := "minimal-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlSourceRepresentationInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlSourceRepresentationInstance_minimal, databaseName),
			},
			{
				ResourceName:      "google_sql_source_representation_instance.master",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSqlSourceRepresentationInstance_full(t *testing.T) {
	t.Parallel()

	databaseName := "full-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlSourceRepresentationInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testGoogleSqlSourceRepresentationInstance_full, databaseName),
			},
			{
				ResourceName:      "google_sql_source_representation_instance.master",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSqlSourceRepresentationInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		config := testAccProvider.Meta().(*Config)
		if rs.Type != "google_sql_source_representation_instance" {
			continue
		}

		_, err := config.clientSqlAdmin.Instances.Get(config.Project,
			rs.Primary.Attributes["name"]).Do()
		if err == nil {
			return fmt.Errorf("Source representation instance still exists")
		}
	}

	return nil
}

var testGoogleSqlSourceRepresentationInstance_minimal = `
resource "google_sql_source_representation_instance" "master" {
  name = "tf-test-source-representation-instance-%s"
  host = "10.20.30.40"
}
`

var testGoogleSqlSourceRepresentationInstance_full = `
resource "google_sql_source_representation_instance" "master" {
  name                 = "tf-test-source-representation-instance-%s"
  database_version     = "MYSQL_5_6"
  region               = "us-east1"
  host                 = "10.20.30.40"
  port                 = 33006
}
`
