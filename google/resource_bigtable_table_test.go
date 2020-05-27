package google

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
