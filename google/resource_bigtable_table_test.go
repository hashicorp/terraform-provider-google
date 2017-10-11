package google

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccBigtableTable_basic(t *testing.T) {
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable(instanceName, tableName),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableTableExists(
						"google_bigtable_table.table"),
				),
			},
		},
	})
}

func TestAccBigtableTable_splitKeys(t *testing.T) {
	instanceName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	tableName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable_splitKeys(instanceName, tableName),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableTableExists(
						"google_bigtable_table.table"),
				),
			},
		},
	})
}

func testAccCheckBigtableTableDestroy(s *terraform.State) error {
	var ctx = context.Background()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_bigtable_table" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		c, err := config.bigtableClientFactory.NewAdminClient(config.Project, rs.Primary.Attributes["instance_name"])
		if err != nil {
			// The instance is already gone
			return nil
		}

		_, err = c.TableInfo(ctx, rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("Table still present. Found %s in %s.", rs.Primary.Attributes["name"], rs.Primary.Attributes["instance_name"])
		}

		c.Close()
	}

	return nil
}

func testAccBigtableTableExists(n string) resource.TestCheckFunc {
	var ctx = context.Background()
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*Config)
		c, err := config.bigtableClientFactory.NewAdminClient(config.Project, rs.Primary.Attributes["instance_name"])
		if err != nil {
			return fmt.Errorf("Error starting admin client. %s", err)
		}

		_, err = c.TableInfo(ctx, rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("Error retrieving table. Could not find %s in %s.", rs.Primary.Attributes["name"], rs.Primary.Attributes["instance_name"])
		}

		c.Close()

		return nil
	}
}

func testAccBigtableTable(instanceName, tableName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  cluster_id    = "%s"
  zone          = "us-central1-b"
  instance_type = "DEVELOPMENT"
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = "${google_bigtable_instance.instance.name}"
}
`, instanceName, instanceName, tableName)
}

func testAccBigtableTable_splitKeys(instanceName, tableName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  cluster_id    = "%s"
  zone          = "us-central1-b"
  instance_type = "DEVELOPMENT"
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = "${google_bigtable_instance.instance.name}"
  split_keys    = ["a", "b", "c"]
}
`, instanceName, instanceName, tableName)
}
