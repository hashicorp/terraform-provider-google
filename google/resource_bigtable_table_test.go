package google

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBigtableTable_basic(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable(instanceName, tableName),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableTable_splitKeys(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable_splitKeys(instanceName, tableName),
			},
			{
				ResourceName:            "google_bigtable_table.table",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"split_keys"},
			},
		},
	})
}

func TestAccBigtableTable_family(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	family := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable_family(instanceName, tableName, family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableTable_deletion_protection_protected(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	family := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			// creating a table with a column family and deletion protection equals to protected
			{
				Config: testAccBigtableTable_deletion_protection(instanceName, tableName, "PROTECTED", family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// it is not possible to delete column families in the table with deletion protection equals to protected
			{
				Config:      testAccBigtableTable(instanceName, tableName),
				ExpectError: regexp.MustCompile(".*deletion protection field is set to true.*"),
			},
			// it is not possible to delete the table because of deletion protection equals to protected
			{
				Config:      testAccBigtableTable_destroyTable(instanceName),
				ExpectError: regexp.MustCompile(".*deletion protection field is set to true.*"),
			},
			// changing deletion protection field to unprotected without changing the column families
			// checking if the table and the column family exists
			{
				Config: testAccBigtableTable_deletion_protection(instanceName, tableName, "UNPROTECTED", family),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableColumnFamilyExists(t, "google_bigtable_table.table", family),
				),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// destroying the table is possible when deletion protection is equals to unprotected
			{
				Config: testAccBigtableTable_destroyTable(instanceName),
			},
			{
				ResourceName:            "google_bigtable_instance.instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection", "instance_type"},
			},
		},
	})
}

func TestAccBigtableTable_deletion_protection_unprotected(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	family := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			// creating a table with a column family and deletion protection equals to unprotected
			{
				Config: testAccBigtableTable_deletion_protection(instanceName, tableName, "UNPROTECTED", family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// removing the column family is possible because the deletion protection field is unprotected
			{
				Config: testAccBigtableTable(instanceName, tableName),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// changing the deletion protection field to protected
			{
				Config: testAccBigtableTable_deletion_protection(instanceName, tableName, "PROTECTED", family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// it is not possible to delete the table because of deletion protection equals to protected
			{
				Config:      testAccBigtableTable_destroyTable(instanceName),
				ExpectError: regexp.MustCompile(".*deletion protection field is set to true.*"),
			},
			// changing the deletion protection field to unprotected so that the sources can properly be destroyed
			{
				Config: testAccBigtableTable_deletion_protection(instanceName, tableName, "UNPROTECTED", family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableTable_familyMany(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	family := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable_familyMany(instanceName, tableName, family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableTable_familyUpdate(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	family := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableTableDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableTable_familyMany(instanceName, tableName, family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccBigtableTable_familyUpdate(instanceName, tableName, family),
			},
			{
				ResourceName:      "google_bigtable_table.table",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckBigtableTableDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		var ctx = context.Background()
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_bigtable_table" {
				continue
			}

			config := googleProviderConfig(t)
			c, err := config.BigTableClientFactory(config.userAgent).NewAdminClient(config.Project, rs.Primary.Attributes["instance_name"])
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
}

func testAccBigtableColumnFamilyExists(t *testing.T, table_name_space, family string) resource.TestCheckFunc {
	var ctx = context.Background()
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[table_name_space]
		if !ok {
			return fmt.Errorf("Table not found: %s", table_name_space)
		}

		config := googleProviderConfig(t)
		c, err := config.BigTableClientFactory(config.userAgent).NewAdminClient(config.Project, rs.Primary.Attributes["instance_name"])
		if err != nil {
			return fmt.Errorf("Error starting admin client. %s", err)
		}

		defer c.Close()

		table, err := c.TableInfo(ctx, rs.Primary.Attributes["name"])
		if err != nil {
			return fmt.Errorf("Error retrieving table. Could not find %s in %s.", rs.Primary.Attributes["name"], rs.Primary.Attributes["instance_name"])
		}
		for _, data := range flattenColumnFamily(table.Families) {
			if data["family"] != family {
				return fmt.Errorf("Error checking column family. Could not find column family %s in %s.", family, rs.Primary.Attributes["name"])
			}
		}

		return nil
	}
}

func testAccBigtableTable(instanceName, tableName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  instance_type = "DEVELOPMENT"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.id
}
`, instanceName, instanceName, tableName)
}

func testAccBigtableTable_splitKeys(instanceName, tableName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name          = "%s"
  instance_type = "DEVELOPMENT"
  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.id
  split_keys    = ["a", "b", "c"]
}
`, instanceName, instanceName, tableName)
}

func testAccBigtableTable_family(instanceName, tableName, family string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.name

  column_family {
    family = "%s"
  }
}
`, instanceName, instanceName, tableName, family)
}

func testAccBigtableTable_deletion_protection(instanceName, tableName, deletionProtection, family string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.name
  deletion_protection = "%s"

  column_family {
    family = "%s"
  }
}
`, instanceName, instanceName, tableName, deletionProtection, family)
}

func testAccBigtableTable_familyMany(instanceName, tableName, family string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.name

  column_family {
    family = "%s-first"
  }

  column_family {
    family = "%s-second"
  }
}
`, instanceName, instanceName, tableName, family, family)
}

func testAccBigtableTable_familyUpdate(instanceName, tableName, family string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
  deletion_protection = false
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.name

  column_family {
    family = "%s-third"
  }

  column_family {
    family = "%s-fourth"
  }

  column_family {
    family = "%s-second"
  }
}
`, instanceName, instanceName, tableName, family, family, family)
}

func testAccBigtableTable_destroyTable(instanceName string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
  deletion_protection = false
}
`, instanceName, instanceName)
}
