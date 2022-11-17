package google

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/bigtable"
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
						t, "google_bigtable_gc_policy.policy", false),
				),
			},
		},
	})
}

func TestAccBigtableGCPolicy_abandoned(t *testing.T) {
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
				Config: testAccBigtableGCPolicyToBeAbandoned(instanceName, tableName, familyName),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableGCPolicyExists(
						t, "google_bigtable_gc_policy.policy", false),
				),
			},
			// Verify that the remote infrastructure GC policy still exists after it is removed in the config.
			{
				Config: testAccBigtableGCPolicyNoPolicy(instanceName, tableName, familyName),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableRemoteGCPolicyExists(
						t, "google_bigtable_table.table"),
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
						t, "google_bigtable_gc_policy.policy", false),
					// Verify can write some data.
					testAccBigtableCanWriteData(
						t, "google_bigtable_gc_policy.policy", 10),
				),
			},
			{
				Config: testAccBigtableGCPolicy(instanceName, tableName, familyName),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableGCPolicyExists(
						t, "google_bigtable_gc_policy.policy", false),
					// Verify no data loss after the GC policy update.
					testAccBigtableCanReadData(
						t, "google_bigtable_gc_policy.policy", 10),
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
						t, "google_bigtable_gc_policy.policy", false),
				),
			},
		},
	})
}

// Testing multiple GC policies; one per column family.
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
						t, "google_bigtable_gc_policy.policyA", false),
					testAccBigtableGCPolicyExists(
						t, "google_bigtable_gc_policy.policyB", false),
					testAccBigtableGCPolicyExists(
						t, "google_bigtable_gc_policy.policyC", false),
				),
			},
		},
	})
}

func TestAccBigtableGCPolicy_gcRulesPolicy(t *testing.T) {
	skipIfVcr(t)
	t.Parallel()

	instanceName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	tableName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	familyName := fmt.Sprintf("tf-test-%s", randString(t, 10))

	gcRulesOriginal := "{\"mode\":\"intersection\",\"rules\":[{\"max_age\":\"10h\"},{\"max_version\":2}]}"
	gcRulesUpdate := "{\"mode\":\"intersection\",\"rules\":[{\"max_age\":\"16h\"},{\"max_version\":1}]}"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBigtableGCPolicyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBigtableGCPolicy_gcRulesCreate(instanceName, tableName, familyName),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableGCPolicyExists(t, "google_bigtable_gc_policy.policy", true),
					resource.TestCheckResourceAttr("google_bigtable_gc_policy.policy", "gc_rules", gcRulesOriginal),
				),
			},
			// Testing gc_rules update
			// TODO: Add test to verify no data loss
			{
				Config: testAccBigtableGCPolicy_gcRulesUpdate(instanceName, tableName, familyName),
				Check: resource.ComposeTestCheckFunc(
					testAccBigtableGCPolicyExists(t, "google_bigtable_gc_policy.policy", true),
					resource.TestCheckResourceAttr("google_bigtable_gc_policy.policy", "gc_rules", gcRulesUpdate),
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

type testUnitBigtableGCPolicyJSONRules struct {
	name          string
	gcJSONString  string
	want          string
	errorExpected bool
}

var testUnitBigtableGCPolicyRulesTestCases = []testUnitBigtableGCPolicyJSONRules{
	{
		name:          "Simple policy",
		gcJSONString:  `{"rules":[{"max_age":"10h"}]}`,
		want:          "age() > 10h",
		errorExpected: false,
	},
	{
		name:          "Simple multiple policies",
		gcJSONString:  `{"mode":"union", "rules":[{"max_age":"10h"},{"max_version":2}]}`,
		want:          "(age() > 10h || versions() > 2)",
		errorExpected: false,
	},
	{
		name:          "Nested policy",
		gcJSONString:  `{"mode":"union", "rules":[{"max_age":"10h"},{"mode": "intersection", "rules":[{"max_age":"2h"}, {"max_version":2}]}]}`,
		want:          "(age() > 10h || (age() > 2h && versions() > 2))",
		errorExpected: false,
	},
	{
		name:          "JSON with no `rules`",
		gcJSONString:  `{"mode": "union"}`,
		errorExpected: true,
	},
	{
		name:          "Empty JSON",
		gcJSONString:  "{}",
		errorExpected: true,
	},
	{
		name:          "Invalid duration string",
		errorExpected: true,
		gcJSONString:  `{"mode":"union","rules":[{"max_age":"12o"},{"max_version":2}]}`,
	},
	{
		name:          "Empty mode policy with more than 1 rules",
		gcJSONString:  `{"rules":[{"max_age":"10h"}, {"max_version":2}]}`,
		errorExpected: true,
	},
	{
		name:          "Less than 2 rules with mode specified",
		gcJSONString:  `{"mode":"union", "rules":[{"max_version":2}]}`,
		errorExpected: true,
	},
	{
		name:          "Invalid GC rule object",
		gcJSONString:  `{"mode": "union", "rules": [{"mode": "intersection"}]}`,
		errorExpected: true,
	},
	{
		name:          "Invalid GC rule field: not max_version or max_age",
		gcJSONString:  `{"mode": "union", "rules": [{"max_versions": 2}]}`,
		errorExpected: true,
	},
	{
		name:          "Invalid GC rule field: additional fields",
		gcJSONString:  `{"mode": "union", "rules": [{"max_age": "10h", "something_else": 100}]}`,
		errorExpected: true,
	},
	{
		name:          "Invalid GC rule field: more than 2 fields in a gc rule object",
		gcJSONString:  `{"mode": "union", "rules": [{"max_age": "10h", "max_version": 10, "something": 100}]}`,
		errorExpected: true,
	},
	{
		name:          "Invalid GC rule field: max_version or max_age is in the wrong type",
		gcJSONString:  `{"mode": "union", "rules": [{"max_age": "10d", "max_version": 2}]}`,
		errorExpected: true,
	},
	{
		name:          "Invalid GC rule: wrong data type for child gc_rule",
		gcJSONString:  `{"rules": {"max_version": "456"}}`,
		errorExpected: true,
	},
}

func TestUnitBigtableGCPolicy_getGCPolicyFromJSON(t *testing.T) {
	for _, tc := range testUnitBigtableGCPolicyRulesTestCases {
		t.Run(tc.name, func(t *testing.T) {
			var topLevelPolicy map[string]interface{}
			err := json.Unmarshal([]byte(tc.gcJSONString), &topLevelPolicy)
			if err != nil {
				t.Fatalf("error unmarshalling JSON string: %v", err)
			}
			got, err := getGCPolicyFromJSON(topLevelPolicy /*isTopLevel=*/, true)
			if tc.errorExpected && err == nil {
				t.Fatal("expect error, got nil")
			} else if !tc.errorExpected && err != nil {
				t.Fatalf("unexpected error: %v", err)
			} else {
				if got != nil && got.String() != tc.want {
					t.Errorf("error getting policy from JSON, got: %v, want: %v", got, tc.want)
				}
			}
		})
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

type testUnitGcPolicyToGCRuleString struct {
	name          string
	policy        bigtable.GCPolicy
	topLevel      bool
	want          string
	errorExpected bool
}

var testUnitGcPolicyToGCRuleStringTestCases = []testUnitGcPolicyToGCRuleString{
	{
		name:          "NoGcPolicy",
		policy:        bigtable.NoGcPolicy(),
		topLevel:      true,
		want:          `{"rules":[{"max_version":1}]}`,
		errorExpected: true,
	},
	{
		name:          "MaxVersionPolicy",
		policy:        bigtable.MaxVersionsPolicy(1),
		topLevel:      true,
		want:          `{"rules":[{"max_version":1}]}`,
		errorExpected: false,
	},
	{
		name:          "MaxAgePolicy",
		policy:        bigtable.MaxAgePolicy(time.Hour),
		topLevel:      true,
		want:          `{"rules":[{"max_age":"1h"}]}`,
		errorExpected: false,
	},
	{
		name:          "UnionPolicy",
		policy:        bigtable.UnionPolicy(bigtable.MaxVersionsPolicy(1), bigtable.MaxAgePolicy(time.Hour)),
		topLevel:      true,
		want:          `{"mode":"union","rules":[{"max_version":1},{"max_age":"1h"}]}`,
		errorExpected: false,
	},
	{
		name:          "IntersectionPolicy",
		policy:        bigtable.IntersectionPolicy(bigtable.MaxVersionsPolicy(1), bigtable.MaxAgePolicy(time.Hour)),
		topLevel:      true,
		want:          `{"mode":"intersection","rules":[{"max_version":1},{"max_age":"1h"}]}`,
		errorExpected: false,
	},
	{
		name:          "NestedPolicy",
		policy:        bigtable.UnionPolicy(bigtable.IntersectionPolicy(bigtable.MaxVersionsPolicy(1), bigtable.MaxAgePolicy(3*time.Hour)), bigtable.MaxAgePolicy(time.Hour)),
		topLevel:      true,
		want:          `{"mode":"union","rules":[{"mode":"intersection","rules":[{"max_version":1},{"max_age":"3h"}]},{"max_age":"1h"}]}`,
		errorExpected: false,
	},
	{
		name:          "MaxVersionPolicyNotTopeLevel",
		policy:        bigtable.MaxVersionsPolicy(1),
		topLevel:      false,
		want:          `{"max_version":1}`,
		errorExpected: false,
	},
	{
		name:          "MaxAgePolicyNotTopeLevel",
		policy:        bigtable.MaxAgePolicy(time.Hour),
		topLevel:      false,
		want:          `{"max_age":"1h"}`,
		errorExpected: false,
	},
}

func TestUnitBigtableGCPolicy_gcPolicyToGCRuleString(t *testing.T) {
	for _, tc := range testUnitGcPolicyToGCRuleStringTestCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := gcPolicyToGCRuleString(tc.policy, tc.topLevel)
			if tc.errorExpected && err == nil {
				t.Fatal("expect error, got nil")
			} else if !tc.errorExpected && err != nil {
				t.Fatalf("unexpected error: %v", err)
			} else {
				if got != nil {
					gcRuleJsonString, err := json.Marshal(got)
					if err != nil {
						t.Fatalf("Error marshaling GC policy to json: %s", err)
					}
					if string(gcRuleJsonString) != tc.want {
						t.Errorf("Unexpected GC policy, got: %v, want: %v", string(gcRuleJsonString), tc.want)
					}
				}
			}
		})
	}
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

func testAccBigtableGCPolicyExists(t *testing.T, n string, compareGcRules bool) resource.TestCheckFunc {
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

		for _, familyInfo := range table.FamilyInfos {
			if familyInfo.Name == rs.Primary.Attributes["column_family"] && familyInfo.GCPolicy == rs.Primary.ID {
				// Ensure the remote GC policy matches the local copy if `compareGcRules` is set to true.
				if !compareGcRules {
					return nil
				}
				gcRuleString, err := gcPolicyToGCRuleString(familyInfo.FullGCPolicy /*isTopLevel=*/, true)
				if err != nil {
					return fmt.Errorf("Error converting GC policy to JSON string: %s", err)
				}
				gcRuleJsonString, err := json.Marshal(gcRuleString)
				if err != nil {
					return fmt.Errorf("Error marshaling GC Policy to JSON: %s", err)
				}
				if string(gcRuleJsonString) == rs.Primary.Attributes["gc_rules"] {
					return nil
				}
				return fmt.Errorf("Found differences in the local and the remote GC policies: %s vs %s", rs.Primary.Attributes["gc_rules"], string(gcRuleJsonString))
			}
		}

		return fmt.Errorf("Error retrieving gc policy. Could not find policy in family %s", rs.Primary.Attributes["column_family"])
	}
}

func testAccBigtableRemoteGCPolicyExists(t *testing.T, table_name_space string) resource.TestCheckFunc {
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

		// We expect a single local column family in the table.
		family, ok := rs.Primary.Attributes["column_family.0.family"]
		if !ok {
			return fmt.Errorf("Error retrieving the local family")
		}

		for _, familyInfo := range table.FamilyInfos {
			if familyInfo.Name == family {
				if familyInfo.GCPolicy == "" {
					return fmt.Errorf("The remote GC policy is missing in family %s", family)
				}
				return nil
			}
		}

		return fmt.Errorf("Error retrieving GC policy. Could not find the column family")
	}
}

func testAccBigtableCanWriteData(t *testing.T, n string, numberOfRows int) resource.TestCheckFunc {
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
		c, err := config.BigTableClientFactory(config.userAgent).NewClient(config.Project, rs.Primary.Attributes["instance_name"])
		if err != nil {
			return fmt.Errorf("Error starting client. %s", err)
		}

		defer c.Close()

		table := c.Open(rs.Primary.Attributes["table"])
		rowKeys := make([]string, numberOfRows)
		mutations := make([]*bigtable.Mutation, numberOfRows)
		columnFamily := rs.Primary.Attributes["column_family"]
		for i := 0; i < 10; i++ {
			rowKeys[i] = fmt.Sprintf("row%d", i)
			mutations[i] = bigtable.NewMutation()
			mutations[i].Set(columnFamily, "column", bigtable.Now(), []byte(fmt.Sprintf("value%d", i)))
		}

		rowErrs, err := table.ApplyBulk(ctx, rowKeys, mutations)
		if err != nil {
			return fmt.Errorf("could not write elements to bigtable: %v", err)
		}
		for _, rowErr := range rowErrs {
			return fmt.Errorf("could not write element to bigtable: %v", rowErr)
		}
		return nil
	}
}

func testAccBigtableCanReadData(t *testing.T, n string, numberOfRows int) resource.TestCheckFunc {
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
		c, err := config.BigTableClientFactory(config.userAgent).NewClient(config.Project, rs.Primary.Attributes["instance_name"])
		if err != nil {
			return fmt.Errorf("Error starting client. %s", err)
		}

		defer c.Close()

		table := c.Open(rs.Primary.Attributes["table"])
		var rows []bigtable.Row
		if err := table.ReadRows(ctx, bigtable.InfiniteRange(""), func(row bigtable.Row) bool {
			rows = append(rows, row)
			return true
		}); err != nil {
			return fmt.Errorf("Could not read elements from bigtable: %v", err)
		}

		if len(rows) != numberOfRows {
			return fmt.Errorf("Expecting %d rows but got: %d", numberOfRows, len(rows))
		}

		return nil
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

func testAccBigtableGCPolicyToBeAbandoned(instanceName, tableName, family string) string {
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

  deletion_policy = "ABANDON"
}
`, instanceName, instanceName, tableName, family, family)
}

func testAccBigtableGCPolicyNoPolicy(instanceName, tableName, family string) string {
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
`, instanceName, instanceName, tableName, family)
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
    family = "%sA"
  }
  column_family {
    family = "%sB"
  }
  column_family {
    family = "%sC"
  }
}

resource "google_bigtable_gc_policy" "policyA" {
  instance_name = google_bigtable_instance.instance.id
  table         = google_bigtable_table.table.name
  column_family = "%sA"

  max_age {
    days = 30
  }
}

resource "google_bigtable_gc_policy" "policyB" {
  instance_name = google_bigtable_instance.instance.id
  table         = google_bigtable_table.table.name
  column_family = "%sB"

  max_version {
    number = 8
  }
}

resource "google_bigtable_gc_policy" "policyC" {
	instance_name = google_bigtable_instance.instance.id
  table         = google_bigtable_table.table.name
  column_family = "%sC"

  max_age {
    days = 7
  }

  max_version {
    number = 10
  }

  mode        = "UNION"
}
`, instanceName, instanceName, tableName, family, family, family, family, family, family)
}

func testAccBigtableGCPolicy_gcRulesCreate(instanceName, tableName, family string) string {
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

		gc_rules = "{\"mode\":\"intersection\", \"rules\":[{\"max_age\":\"10h\"},{\"max_version\":2}]}"
	}
`, instanceName, instanceName, tableName, family, family)
}

func testAccBigtableGCPolicy_gcRulesUpdate(instanceName, tableName, family string) string {
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

		gc_rules = "{\"mode\":\"intersection\", \"rules\":[{\"max_age\":\"16h\"},{\"max_version\":1}]}"
	}
`, instanceName, instanceName, tableName, family, family)
}
