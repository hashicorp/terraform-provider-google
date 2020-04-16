package google

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccBigtableGCPolicy_basic(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	familyName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableGCPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableGCPolicy(instanceName, tableName, familyName),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableGCPolicyExists(
						t, "google_bigtable_gc_policy.policy"),
				),
			},
		},
	})
}

func TestAccBigtableGCPolicy_union(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	familyName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableGCPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableGCPolicyUnion(instanceName, tableName, familyName),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableGCPolicyExists(
						t, "google_bigtable_gc_policy.policy"),
				),
			},
		},
	})
}

func testAccCheckBigtableGCPolicyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		var ctx = context.Background()
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_bigtable_gc_policy" {
				continue
			}

			config := googleProviderConfig(t)
			c, err := config.bigtableClientFactory.NewAdminClient(config.Project, rs.Primary.Attributes["instance_name"])
			if err != nil {
				// The instance is already gone
				return nil
			}

			table, err := c.TableInfo(ctx, rs.Primary.Attributes["name"])
			if err != nil {
				// The table is already gone
				return nil
			}

			for _, i := range table.FamilyInfos {
				if i.Name == rs.Primary.Attributes["column_family"] {
					if i.GCPolicy != "<never>" {
						return fmt.Errorf("GC Policy still present. Found %s in %s.", i.GCPolicy, rs.Primary.Attributes["column_family"])
					}
				}
			}

			c.Close()
		}

		return nil
	}
}

func testAccBigtableGCPolicyExists(t *testing.T, n string) resource.TestCheckFunc {
	var ctx = context.Background()
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := googleProviderConfig(t)
		c, err := config.bigtableClientFactory.NewAdminClient(config.Project, rs.Primary.Attributes["instance_name"])
		if err != nil {
			return fmt.Errorf("Error starting admin client. %s", err)
		}

		defer c.Close()

		table, err := c.TableInfo(ctx, rs.Primary.Attributes["table"])
		if err != nil {
			return fmt.Errorf("Error retrieving table. Could not find %s in %s.", rs.Primary.Attributes["table"], rs.Primary.Attributes["instance_name"])
		}

		for _, i := range table.FamilyInfos {
			if i.Name == rs.Primary.Attributes["column_family"] {
				return nil
			}
		}

		return fmt.Errorf("Error retrieving gc policy. Could not find policy in family %s", rs.Primary.Attributes["column_family"])
	}
}

func testAccBigtableGCPolicy(instanceName, tableName, family string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.id

  column_family {
    family = "%s"
  }
}

resource "google_bigtable_gc_policy" "policy" {
  instance_name = google_bigtable_instance.instance.id
  table         = google_bigtable_table.table.name
  column_family = "%s"

  max_age {
    days = 3
  }
}
`, instanceName, instanceName, tableName, family, family)
}

func testAccBigtableGCPolicyUnion(instanceName, tableName, family string) string {
	return fmt.Sprintf(`
resource "google_bigtable_instance" "instance" {
  name = "%s"

  cluster {
    cluster_id = "%s"
    zone       = "us-central1-b"
  }

  instance_type = "DEVELOPMENT"
}

resource "google_bigtable_table" "table" {
  name          = "%s"
  instance_name = google_bigtable_instance.instance.name

  column_family {
    family = "%s"
  }
}

resource "google_bigtable_gc_policy" "policy" {
  instance_name = google_bigtable_instance.instance.name
  table         = google_bigtable_table.table.name
  column_family = "%s"

  mode = "UNION"

  max_age {
    days = 3
  }

  max_version {
    number = 10
  }
}
`, instanceName, instanceName, tableName, family, family)
}
