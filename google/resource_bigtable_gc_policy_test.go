package google

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBigtableGCPolicy_basic(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
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

func TestAccBigtableGCPolicy_swapOffDeprecated(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
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
				Config: testAccBigtableGCPolicy_days(instanceName, tableName, familyName),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableGCPolicyExists(
						t, "google_bigtable_gc_policy.policy"),
				),
			},
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
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
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

func TestAccBigtableGCPolicy_multiplePolicies(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
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
				Config: testAccBigtableGCPolicy_multiplePolicies(instanceName, tableName, familyName),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableGCPolicyExists(
						t, "google_bigtable_gc_policy.policyA"),
					testAccBigtableGCPolicyExists(
						t, "google_bigtable_gc_policy.policyB"),
					testAccBigtableGCPolicyExists(
						t, "google_bigtable_gc_policy.policyC"),
				),
			},
		},
	})
}

func TestUnitBigtableGCPolicy_customizeDiff(t *testing.T) {
	for _, tc := range testUnitBigtableGCPolicyCustomizeDiffTestcases {
		tc.check(t)
	}
}

func (testcase *testUnitBigtableGCPolicyCustomizeDiffTestcase) check(t *testing.T) {
	d := &ResourceDiffMock{
		Before: map[string]interface{}{},
		After:  map[string]interface{}{},
	}

	d.Before["max_age.0.days"] = testcase.oldDays
	d.Before["max_age.0.duration"] = testcase.oldDuration

	d.After["max_age.#"] = testcase.arraySize
	d.After["max_age.0.days"] = testcase.newDays
	d.After["max_age.0.duration"] = testcase.newDuration

	err := resourceBigtableGCPolicyCustomizeDiffFunc(d)
	if err != nil {
		t.Errorf("error on testcase %s - %v", testcase.testName, err)
	}

	var cleared bool = d.Cleared != nil && d.Cleared["max_age.0.duration"] == true && d.Cleared["max_age.0.days"] == true
	if cleared != testcase.cleared {
		t.Errorf("%s: expected diff clear to be %v, but was %v", testcase.testName, testcase.cleared, cleared)
	}
}

type testUnitBigtableGCPolicyCustomizeDiffTestcase struct {
	testName    string
	arraySize   int
	oldDays     int
	newDays     int
	oldDuration string
	newDuration string
	cleared     bool
}

var testUnitBigtableGCPolicyCustomizeDiffTestcases = []testUnitBigtableGCPolicyCustomizeDiffTestcase{
	{
		testName:  "ArraySize0",
		arraySize: 0,
		cleared:   false,
	},
	{
		testName:  "DaysChange",
		arraySize: 1,
		oldDays:   3,
		newDays:   2,
		cleared:   false,
	},
	{
		testName:    "DurationChanges",
		arraySize:   1,
		oldDuration: "3h",
		newDuration: "4h",
		cleared:     false,
	},
	{
		testName:    "DaysToDurationEq",
		arraySize:   1,
		oldDays:     3,
		newDuration: "72h",
		cleared:     true,
	},
	{
		testName:    "DaysToDurationNotEq",
		arraySize:   1,
		oldDays:     3,
		newDuration: "70h",
		cleared:     false,
	},
}

func testAccCheckBigtableGCPolicyDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		var ctx = context.Background()
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_bigtable_gc_policy" {
				continue
			}

			config := googleProviderConfig(t)
			c, err := config.BigTableClientFactory(config.userAgent).NewAdminClient(config.Project, rs.Primary.Attributes["instance_name"])
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
		c, err := config.BigTableClientFactory(config.userAgent).NewAdminClient(config.Project, rs.Primary.Attributes["instance_name"])
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

func testAccBigtableGCPolicy_days(instanceName, tableName, family string) string {
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

func testAccBigtableGCPolicy(instanceName, tableName, family string) string {
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
    duration = "72h"
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
  deletion_protection = false
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
    duration = "72h"
  }

  max_version {
    number = 10
  }
}
`, instanceName, instanceName, tableName, family, family)
}

func testAccBigtableGCPolicy_multiplePolicies(instanceName, tableName, family string) string {
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
  instance_name = google_bigtable_instance.instance.id

  column_family {
    family = "%s"
  }
}

resource "google_bigtable_gc_policy" "policyA" {
  instance_name = google_bigtable_instance.instance.id
  table         = google_bigtable_table.table.name
  column_family = "%s"

  max_age {
    days = 30
  }
}

resource "google_bigtable_gc_policy" "policyB" {
  instance_name = google_bigtable_instance.instance.id
  table         = google_bigtable_table.table.name
  column_family = "%s"

  max_version {
    number = 8
  }
}

resource "google_bigtable_gc_policy" "policyC" {
	instance_name = google_bigtable_instance.instance.id
  table         = google_bigtable_table.table.name
  column_family = "%s"

  max_age {
    days = 7
  }

  max_version {
    number = 10
  }

  mode        = "UNION"
}
`, instanceName, instanceName, tableName, family, family, family, family)
}
