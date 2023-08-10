// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sql_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceSqlDatabaseInstances_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSqlDatabaseInstances_basic(context),
				Check: resource.ComposeTestCheckFunc(
					checkListDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_sql_database_instances.qa",
						"google_sql_database_instance.main",
						"google_sql_database_instance.main2",
						map[string]struct{}{
							"deletion_protection": {},
							"id":                  {},
						},
					),
				),
			},
		},
	})
}

func TestAccDataSourceSqlDatabaseInstances_databaseVersionFilter(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSqlDatabaseInstances_databaseVersionFilter(context),
				Check: resource.ComposeTestCheckFunc(
					checkListDataSourceStateMatchesResourceStateWithIgnoresForAppliedFilter(
						"data.google_sql_database_instances.qa",
						"google_sql_database_instance.main",
						"google_sql_database_instance.main2",
						map[string]struct{}{
							"deletion_protection": {},
							"id":                  {},
						},
					),
				),
			},
		},
	})
}

func TestAccDataSourceSqlDatabaseInstances_regionFilter(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSqlDatabaseInstances_regionFilter(context),
				Check: resource.ComposeTestCheckFunc(
					checkListDataSourceStateMatchesResourceStateWithIgnoresForAppliedFilter(
						"data.google_sql_database_instances.qa",
						"google_sql_database_instance.main",
						"google_sql_database_instance.main2",
						map[string]struct{}{
							"deletion_protection": {},
							"id":                  {},
						},
					),
				),
			},
		},
	})
}

func TestAccDataSourceSqlDatabaseInstances_tierFilter(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSqlDatabaseInstances_tierFilter(context),
				Check: resource.ComposeTestCheckFunc(
					checkListDataSourceStateMatchesResourceStateWithIgnoresForAppliedFilter(
						"data.google_sql_database_instances.qa",
						"google_sql_database_instance.main",
						"google_sql_database_instance.main2",
						map[string]struct{}{
							"deletion_protection": {},
							"id":                  {},
						},
					),
				),
			},
		},
	})
}

func testAccDataSourceSqlDatabaseInstances_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sql_database_instance" "main" {
  name             = "tf-test-instance-%{random_suffix}"
  database_version = "POSTGRES_14"
  region           = "us-central1"

  settings {
    # Second-generation instance tiers are based on the machine
    # type. See argument reference below.
    tier = "db-f1-micro"
  }

  deletion_protection = false
}

resource "google_sql_database_instance" "main2" {
	name             = "tf-test-instance-2-%{random_suffix}"
	database_version = "MYSQL_8_0"
	region           = "us-central1"
  
	settings {
	  # Second-generation instance tiers are based on the machine
	  # type. See argument reference below.
	  tier = "db-f1-micro"
	}
  
	deletion_protection = false
  }


data "google_sql_database_instances" "qa" {
	depends_on = [
		google_sql_database_instance.main2,
		google_sql_database_instance.main
	]
}
`, context)
}

func testAccDataSourceSqlDatabaseInstances_databaseVersionFilter(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sql_database_instance" "main" {
  name             = "tf-test-instance-%{random_suffix}"
  database_version = "POSTGRES_14"
  region           = "us-central1"

  settings {
    # Second-generation instance tiers are based on the machine
    # type. See argument reference below.
    tier = "db-f1-micro"
  }

  deletion_protection = false
}

resource "google_sql_database_instance" "main2" {
	name             = "tf-test-instance-2-%{random_suffix}"
	database_version = "MYSQL_8_0"
	region           = "us-central1"
  
	settings {
	  # Second-generation instance tiers are based on the machine
	  # type. See argument reference below.
	  tier = "db-f1-micro"
	}
  
	deletion_protection = false
  }


data "google_sql_database_instances" "qa" {
	database_version = "MYSQL_8_0"
	depends_on = [
		google_sql_database_instance.main2,
		google_sql_database_instance.main
	]
}
`, context)
}

func testAccDataSourceSqlDatabaseInstances_regionFilter(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sql_database_instance" "main" {
  name             = "tf-test-instance-%{random_suffix}"
  database_version = "POSTGRES_14"
  region           = "us-central1"

  settings {
    # Second-generation instance tiers are based on the machine
    # type. See argument reference below.
    tier = "db-f1-micro"
  }

  deletion_protection = false
}

resource "google_sql_database_instance" "main2" {
	name             = "tf-test-instance-2-%{random_suffix}"
	database_version = "MYSQL_8_0"
	region           = "us-east1"
  
	settings {
	  # Second-generation instance tiers are based on the machine
	  # type. See argument reference below.
	  tier = "db-f1-micro"
	}
  
	deletion_protection = false
  }


data "google_sql_database_instances" "qa" {
	region = "us-east1"
	depends_on = [
		google_sql_database_instance.main2,
		google_sql_database_instance.main
	]
}
`, context)
}

func testAccDataSourceSqlDatabaseInstances_tierFilter(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_sql_database_instance" "main" {
  name             = "tf-test-instance-%{random_suffix}"
  database_version = "POSTGRES_14"
  region           = "us-central1"

  settings {
    # Second-generation instance tiers are based on the machine
    # type. See argument reference below.
    tier = "db-f1-micro"
  }

  deletion_protection = false
}

resource "google_sql_database_instance" "main2" {
	name             = "tf-test-instance-2-%{random_suffix}"
	database_version = "MYSQL_8_0"
	region           = "us-central1"
  
	settings {
	  # Second-generation instance tiers are based on the machine
	  # type. See argument reference below.
	  tier = "db-custom-2-13312"
	}
  
	deletion_protection = false
  }


data "google_sql_database_instances" "qa" {
	region = "us-central1"
	tier = "db-custom-2-13312"
	depends_on = [
		google_sql_database_instance.main2,
		google_sql_database_instance.main
	]
}
`, context)
}

// This function checks data source state matches for resorceName database instance state
func checkListDataSourceStateMatchesResourceStateWithIgnores(dataSourceName, resourceName, resourceName2 string, ignoreFields map[string]struct{}) func(*terraform.State) error {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", dataSourceName)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}

		rs2, ok := s.RootModule().Resources[resourceName2]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName2)
		}

		dsAttr := ds.Primary.Attributes
		rsAttr := rs.Primary.Attributes
		rsAttr2 := rs2.Primary.Attributes

		err := checkFieldsMatchForDataSourceStateAndResourceState(dsAttr, rsAttr, ignoreFields)
		if err != nil {
			return err
		}
		err = checkFieldsMatchForDataSourceStateAndResourceState(dsAttr, rsAttr2, ignoreFields)
		return err

	}
}

// This function checks state match for resorceName2 and asserts the absense of resorceName in data source
func checkListDataSourceStateMatchesResourceStateWithIgnoresForAppliedFilter(dataSourceName, resourceName, resourceName2 string, ignoreFields map[string]struct{}) func(*terraform.State) error {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", dataSourceName)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}

		rs2, ok := s.RootModule().Resources[resourceName2]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName2)
		}

		dsAttr := ds.Primary.Attributes
		rsAttr := rs.Primary.Attributes
		rsAttr2 := rs2.Primary.Attributes

		err := checkResourceAbsentInDataSourceAfterFilterApllied(dsAttr, rsAttr)
		if err != nil {
			return err
		}
		err = checkFieldsMatchForDataSourceStateAndResourceState(dsAttr, rsAttr2, ignoreFields)
		return err

	}
}

// This function asserts the absence of the database instance resource which would not be included in the data source list due to the filter applied.
func checkResourceAbsentInDataSourceAfterFilterApllied(dsAttr, rsAttr map[string]string) error {
	totalInstances, err := strconv.Atoi(dsAttr["instances.#"])
	if err != nil {
		return errors.New("Couldn't convert length of instances list to integer")
	}
	for i := 0; i < totalInstances; i++ {
		if dsAttr["instances."+strconv.Itoa(i)+".name"] == rsAttr["name"] {
			return errors.New("The resource is present in data source event after filter applied")
		}
	}
	return nil
}

// This function checks whether all the attributes of the database instance resource and the attributes of the datbase instance inside the data source list are the same
func checkFieldsMatchForDataSourceStateAndResourceState(dsAttr, rsAttr map[string]string, ignoreFields map[string]struct{}) error {
	totalInstances, err := strconv.Atoi(dsAttr["instances.#"])
	if err != nil {
		return errors.New("Couldn't convert length of instances list to integer")
	}
	index := "-1"
	for i := 0; i < totalInstances; i++ {
		if dsAttr["instances."+strconv.Itoa(i)+".name"] == rsAttr["name"] {
			index = strconv.Itoa(i)
		}
	}

	if index == "-1" {
		return errors.New("The newly created intance is not found in the data source")
	}

	errMsg := ""
	// Data sources are often derived from resources, so iterate over the resource fields to
	// make sure all fields are accounted for in the data source.
	// If a field exists in the data source but not in the resource, its expected value should
	// be checked separately.
	for k := range rsAttr {
		if _, ok := ignoreFields[k]; ok {
			continue
		}
		if k == "%" {
			continue
		}
		if dsAttr["instances."+index+"."+k] != rsAttr[k] {
			// ignore data sources where an empty list is being compared against a null list.
			if k[len(k)-1:] == "#" && (dsAttr["instances."+index+"."+k] == "" || dsAttr["instances."+index+"."+k] == "0") && (rsAttr[k] == "" || rsAttr[k] == "0") {
				continue
			}
			errMsg += fmt.Sprintf("%s is %s; want %s\n", k, dsAttr["instances."+index+"."+k], rsAttr[k])
		}
	}

	if errMsg != "" {
		return errors.New(errMsg)
	}

	return nil
}
