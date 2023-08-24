// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package spanner_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/spanner"
)

func TestAccSpannerDatabase_basic(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	rnd := acctest.RandString(t, 10)
	instanceName := fmt.Sprintf("tf-test-%s", rnd)
	databaseName := fmt.Sprintf("tfgen_%s", rnd)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerDatabaseDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerDatabase_virtualUpdate(instanceName, databaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_database.basic", "state"),
					resource.TestCheckResourceAttr("google_spanner_database.basic", "version_retention_period", "1h"), // default set by API
				),
			},
			{
				Config: testAccSpannerDatabase_basic(instanceName, databaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_database.basic", "state"),
					resource.TestCheckResourceAttr("google_spanner_database.basic", "version_retention_period", "1h"), // default set by API
				),
			},
			{
				// Test import with default Terraform ID
				ResourceName:            "google_spanner_database.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ddl", "deletion_protection"},
			},
			{
				Config: testAccSpannerDatabase_basicUpdate(instanceName, databaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_database.basic", "state"),
					resource.TestCheckResourceAttr("google_spanner_database.basic", "version_retention_period", "2d"),
				),
			},
			{
				// Test import with default Terraform ID
				ResourceName:            "google_spanner_database.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ddl", "deletion_protection"},
			},
			{
				ResourceName:            "google_spanner_database.basic",
				ImportStateId:           fmt.Sprintf("projects/%s/instances/%s/databases/%s", project, instanceName, databaseName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ddl", "deletion_protection"},
			},
			{
				ResourceName:            "google_spanner_database.basic",
				ImportStateId:           fmt.Sprintf("instances/%s/databases/%s", instanceName, databaseName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ddl", "deletion_protection"},
			},
			{
				ResourceName:            "google_spanner_database.basic",
				ImportStateId:           fmt.Sprintf("%s/%s", instanceName, databaseName),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ddl", "deletion_protection"},
			},
		},
	})
}

func testAccSpannerDatabase_basic(instanceName, databaseName string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s-display"
  num_nodes    = 1
}

resource "google_spanner_database" "basic" {
  instance = google_spanner_instance.basic.name
  name     = "%s"
  ddl = [
	"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
	"CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)",
  ]
  deletion_protection = false
}
`, instanceName, instanceName, databaseName)
}

func testAccSpannerDatabase_basicUpdate(instanceName, databaseName string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s-display"
  num_nodes    = 1
}

resource "google_spanner_database" "basic" {
  instance = google_spanner_instance.basic.name
  name     = "%s"
  version_retention_period = "2d" # increase from default 1h
  ddl = [
	"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
	"CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)",
	"CREATE TABLE t3 (t3 INT64 NOT NULL,) PRIMARY KEY(t3)",
	"CREATE TABLE t4 (t4 INT64 NOT NULL,) PRIMARY KEY(t4)",
  ]
  deletion_protection = false
}
`, instanceName, instanceName, databaseName)
}

func testAccSpannerDatabase_virtualUpdate(instanceName, databaseName string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s-display"
  num_nodes    = 1
}

resource "google_spanner_database" "basic" {
  instance = google_spanner_instance.basic.name
  name     = "%s"
  ddl = [
	"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
	"CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)",
  ]
  deletion_protection = true
}
`, instanceName, instanceName, databaseName)
}

func TestAccSpannerDatabase_postgres(t *testing.T) {
	t.Parallel()

	rnd := acctest.RandString(t, 10)
	instanceName := fmt.Sprintf("tf-test-%s", rnd)
	databaseName := fmt.Sprintf("tfgen_%s", rnd)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerDatabaseDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerDatabase_postgres(instanceName, databaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_database.basic_spangres", "state"),
				),
			},
			{
				// Test import with default Terraform ID
				ResourceName:            "google_spanner_database.basic_spangres",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ddl", "deletion_protection"},
			},
			{
				Config: testAccSpannerDatabase_postgresUpdate(instanceName, databaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_database.basic_spangres", "state"),
				),
			},
			{
				// Test import with default Terraform ID
				ResourceName:            "google_spanner_database.basic_spangres",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ddl", "deletion_protection"},
			},
		},
	})
}

func testAccSpannerDatabase_postgres(instanceName, databaseName string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s-display"
  num_nodes    = 1
}

resource "google_spanner_database" "basic_spangres" {
  instance = google_spanner_instance.basic.name
  name     = "%s-spangres"
  database_dialect = "POSTGRESQL"
  // Confirm that DDL can be run at creation time for POSTGRESQL
  version_retention_period = "2h"
  ddl = [
     "CREATE TABLE t1 (t1 bigint NOT NULL PRIMARY KEY)",
  ]
  deletion_protection = false
}
`, instanceName, instanceName, databaseName)
}

func testAccSpannerDatabase_postgresUpdate(instanceName, databaseName string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s-display"
  num_nodes    = 1
}

resource "google_spanner_database" "basic_spangres" {
  instance = google_spanner_instance.basic.name
  name     = "%s-spangres"
  database_dialect = "POSTGRESQL"
  version_retention_period = "4d"
  ddl = [
     "CREATE TABLE t2 (t2 bigint NOT NULL PRIMARY KEY)",
     "CREATE TABLE t3 (t3 bigint NOT NULL PRIMARY KEY)",
     "CREATE TABLE t4 (t4 bigint NOT NULL PRIMARY KEY)",
  ]
  deletion_protection = false
}
`, instanceName, instanceName, databaseName)
}

func TestAccSpannerDatabase_versionRetentionPeriod(t *testing.T) {
	t.Parallel()

	rnd := acctest.RandString(t, 10)
	instanceName := fmt.Sprintf("tf-test-%s", rnd)
	databaseName := fmt.Sprintf("tfgen_%s", rnd)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerDatabaseDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				// Test creating a database with `version_retention_period` set
				Config: testAccSpannerDatabase_versionRetentionPeriod(instanceName, databaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_database.basic", "state"),
					resource.TestCheckResourceAttr("google_spanner_database.basic", "version_retention_period", "2h"),
				),
			},
			{
				// Test removing `version_retention_period` and setting retention period to a new value with a DDL statement in `ddl`
				Config: testAccSpannerDatabase_versionRetentionPeriodUpdate1(instanceName, databaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_database.basic", "state"),
					resource.TestCheckResourceAttr("google_spanner_database.basic", "version_retention_period", "4h"),
				),
			},
			{
				// Test that adding `version_retention_period` controls retention time, regardless of any previous statements in `ddl`
				Config: testAccSpannerDatabase_versionRetentionPeriodUpdate2(instanceName, databaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_database.basic", "state"),
					resource.TestCheckResourceAttr("google_spanner_database.basic", "version_retention_period", "2h"),
				),
			},
			{
				// Test that changing the retention value via DDL when `version_retention_period` is set:
				// - changes the value (from 2h to 8h)
				// - is unstable; non-empty plan afterwards due to conflict
				Config:             testAccSpannerDatabase_versionRetentionPeriodUpdate3(instanceName, databaseName),
				ExpectNonEmptyPlan: true, // is unstable
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_database.basic", "state"),
					resource.TestCheckResourceAttr("google_spanner_database.basic", "version_retention_period", "8h"),
				),
			},
			{
				// Test that when the above config is reapplied:
				// - changes the value (reverts to set value of `version_retention_period`, 2h)
				// - is stable; no further conflict
				Config:             testAccSpannerDatabase_versionRetentionPeriodUpdate3(instanceName, databaseName), //same as previous step
				ExpectNonEmptyPlan: false,                                                                            // is stable
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_spanner_database.basic", "state"),
					resource.TestCheckResourceAttr("google_spanner_database.basic", "version_retention_period", "2h"),
				),
			},
		},
	})
}

func testAccSpannerDatabase_versionRetentionPeriod(instanceName, databaseName string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s-display"
  num_nodes    = 1
}

resource "google_spanner_database" "basic" {
  instance = google_spanner_instance.basic.name
  name     = "%s"
  version_retention_period = "2h"
  ddl = [
     "CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
  ]
  deletion_protection = false
}
`, instanceName, instanceName, databaseName)
}

func testAccSpannerDatabase_versionRetentionPeriodUpdate1(instanceName, databaseName string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s-display"
  num_nodes    = 1
}

resource "google_spanner_database" "basic" {
  instance = google_spanner_instance.basic.name
  name     = "%s"
  // Change 1/2 : deleted version_retention_period argument
  ddl = [
    "CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
    "ALTER DATABASE %s SET OPTIONS (version_retention_period=\"4h\")",  // Change 2/2 : set retention with new DDL
  ]
  deletion_protection = false
}
`, instanceName, instanceName, databaseName, databaseName)
}

func testAccSpannerDatabase_versionRetentionPeriodUpdate2(instanceName, databaseName string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s-display"
  num_nodes    = 1
}

resource "google_spanner_database" "basic" {
  instance = google_spanner_instance.basic.name
  name     = "%s"
  version_retention_period = "2h" // Change : added version_retention_period argument
  ddl = [
    "CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
    "ALTER DATABASE %s SET OPTIONS (version_retention_period=\"4h\")",
  ]
  deletion_protection = false
}
`, instanceName, instanceName, databaseName, databaseName)
}

func testAccSpannerDatabase_versionRetentionPeriodUpdate3(instanceName, databaseName string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s-display"
  num_nodes    = 1
}

resource "google_spanner_database" "basic" {
  instance = google_spanner_instance.basic.name
  name     = "%s"
  version_retention_period = "2h"
  ddl = [
    "CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
    "ALTER DATABASE %s SET OPTIONS (version_retention_period=\"4h\")",
    "ALTER DATABASE %s SET OPTIONS (version_retention_period=\"8h\")",  // Change : set retention with new DDL
  ]
  deletion_protection = false
}
`, instanceName, instanceName, databaseName, databaseName, databaseName)
}

func TestAccSpannerDatabase_enableDropProtection(t *testing.T) {
	t.Parallel()

	rnd := acctest.RandString(t, 10)
	instanceName := fmt.Sprintf("tf-test-%s", rnd)
	databaseName := fmt.Sprintf("tfgen_%s", rnd)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerDatabaseDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerDatabase_enableDropProtection(instanceName, databaseName),
			},
			{
				ResourceName:            "google_spanner_database.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ddl", "deletion_protection"},
			},
			{
				Config: testAccSpannerDatabase_enableDropProtectionUpdate(instanceName, databaseName),
			},
			{
				ResourceName:            "google_spanner_database.basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ddl", "deletion_protection"},
			},
		},
	})
}

func testAccSpannerDatabase_enableDropProtection(instanceName, databaseName string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s-display"
  num_nodes    = 1
}

resource "google_spanner_database" "basic" {
  instance = google_spanner_instance.basic.name
  name     = "%s"
  enable_drop_protection = true
  deletion_protection = false
  ddl = [
     "CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
  ]
}
`, instanceName, instanceName, databaseName)
}

func testAccSpannerDatabase_enableDropProtectionUpdate(instanceName, databaseName string) string {
	return fmt.Sprintf(`
resource "google_spanner_instance" "basic" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s-display"
  num_nodes    = 1
}

resource "google_spanner_database" "basic" {
  instance = google_spanner_instance.basic.name
  name     = "%s"
  enable_drop_protection = false
  deletion_protection = false
  ddl = [
     "CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
  ]
}
`, instanceName, instanceName, databaseName)
}

// Unit Tests for validation of retention period argument
func TestValidateDatabaseRetentionPeriod(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		input       string
		expectError bool
	}{
		// Not valid input
		"empty_string": {
			input:       "",
			expectError: true,
		},
		"number_with_no_unit": {
			input:       "1",
			expectError: true,
		},
		"less_than_1h": {
			input:       "59m",
			expectError: true,
		},
		"more_than_7days": {
			input:       "8d",
			expectError: true,
		},
		// Valid input
		"1_hour_in_secs": {
			input:       "3600s",
			expectError: false,
		},
		"1_hour_in_mins": {
			input:       "60m",
			expectError: false,
		},
		"1_hour_in_hours": {
			input:       "1h",
			expectError: false,
		},
		"7_days_in_secs": {
			input:       fmt.Sprintf("%ds", 7*24*60*60),
			expectError: false,
		},
		"7_days_in_mins": {
			input:       fmt.Sprintf("%dm", 7*24*60),
			expectError: false,
		},
		"7_days_in_hours": {
			input:       fmt.Sprintf("%dh", 7*24),
			expectError: false,
		},
		"7_days_in_days": {
			input:       "7d",
			expectError: false,
		},
	}

	for tn, tc := range testCases {
		t.Run(tn, func(t *testing.T) {
			_, errs := spanner.ValidateDatabaseRetentionPeriod(tc.input, "foobar")
			var wantErrCount string
			if tc.expectError {
				wantErrCount = "1+"
			} else {
				wantErrCount = "0"
			}
			if (len(errs) > 0 && tc.expectError == false) || (len(errs) == 0 && tc.expectError == true) {
				t.Errorf("failed, expected `%s` test case validation to have %s errors", tn, wantErrCount)
			}
		})
	}
}

func TestAccSpannerDatabase_deletionProtection(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerDatabaseDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerDatabase_deletionProtection(context),
			},
			{
				ResourceName:            "google_spanner_database.database",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ddl", "instance", "deletion_protection"},
			},
			{
				Config:      testAccSpannerDatabase_deletionProtection(context),
				Destroy:     true,
				ExpectError: regexp.MustCompile("deletion_protection"),
			},
			{
				Config: testAccSpannerDatabase_spannerDatabaseBasicExample(context),
			},
		},
	})
}

func testAccSpannerDatabase_deletionProtection(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_spanner_instance" "main" {
  config       = "regional-europe-west1"
  display_name = "main-instance"
  num_nodes    = 1
}

resource "google_spanner_database" "database" {
  instance = google_spanner_instance.main.name
  name     = "tf-test-my-database%{random_suffix}"
  ddl = [
    "CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
    "CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)",
  ]
}
`, context)
}
