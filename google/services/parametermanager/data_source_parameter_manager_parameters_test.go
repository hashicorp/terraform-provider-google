// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package parametermanager_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceParameterManagerParameters_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckParameterManagerParameterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceParameterManagerParameters_basic(context),
				Check: resource.ComposeTestCheckFunc(
					checkListDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_parameter_manager_parameters.parameters-datasource",
						"google_parameter_manager_parameter.parameters",
						map[string]struct{}{
							"id":      {},
							"project": {},
						},
					),
				),
			},
		},
	})
}

func testAccDataSourceParameterManagerParameters_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  add_terraform_attribution_label = false
}

resource "google_parameter_manager_parameter" "parameters" {
  parameter_id = "tf_test_parameter%{random_suffix}"
  format = "YAML"

  labels = {
    key1 = "val1"
    key2 = "val2"
    key3 = "val3"
    key4 = "val4"
    key5 = "val5"
  }
}

data "google_parameter_manager_parameters" "parameters-datasource" {
  depends_on = [
    google_parameter_manager_parameter.parameters
  ]
}
`, context)
}

func TestAccDataSourceParameterManagerParameters_filter(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckParameterManagerParameterDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceParameterManagerParameters_filter(context),
				Check: resource.ComposeTestCheckFunc(
					checkListDataSourceStateMatchesResourceStateWithIgnoresForAppliedFilter(
						"data.google_parameter_manager_parameters.parameters-datasource-filter",
						"google_parameter_manager_parameter.parameters-1",
						"google_parameter_manager_parameter.parameters-2",
						map[string]struct{}{
							"id":      {},
							"project": {},
						},
					),
				),
			},
		},
	})
}

func testAccDataSourceParameterManagerParameters_filter(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  add_terraform_attribution_label = false
}

resource "google_parameter_manager_parameter" "parameters-1" {
  parameter_id = "tf_test_parameter%{random_suffix}"
  format = "JSON"

  labels = {
    key1 = "val1"
  }
}

resource "google_parameter_manager_parameter" "parameters-2" {
  parameter_id = "tf_test_parameter_2_%{random_suffix}"
  format = "YAML"

  labels = {
    keyoth1 = "valoth1"
  }
}

data "google_parameter_manager_parameters" "parameters-datasource-filter" {
  filter = "format:JSON"
  depends_on = [
    google_parameter_manager_parameter.parameters-1,
	google_parameter_manager_parameter.parameters-2
  ]
}
`, context)
}

// This function checks data source state matches for resourceName parameter manager parameter state
func checkListDataSourceStateMatchesResourceStateWithIgnores(dataSourceName, resourceName string, ignoreFields map[string]struct{}) func(*terraform.State) error {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", dataSourceName)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", resourceName)
		}

		dsAttr := ds.Primary.Attributes
		rsAttr := rs.Primary.Attributes

		err := checkFieldsMatchForDataSourceStateAndResourceState(dsAttr, rsAttr, ignoreFields)
		if err != nil {
			return err
		}
		return nil
	}
}

// This function checks whether all the attributes of the parameter manager parameter resource and the attributes of the parameter manager parameter inside the data source list are the same
func checkFieldsMatchForDataSourceStateAndResourceState(dsAttr, rsAttr map[string]string, ignoreFields map[string]struct{}) error {
	totalParameters, err := strconv.Atoi(dsAttr["parameters.#"])
	if err != nil {
		return errors.New("couldn't convert length of parameters list to integer")
	}
	index := "-1"
	for i := 0; i < totalParameters; i++ {
		if dsAttr["parameters."+strconv.Itoa(i)+".name"] == rsAttr["name"] {
			index = strconv.Itoa(i)
		}
	}

	if index == "-1" {
		return errors.New("the newly created parameter is not found in the data source")
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
		if dsAttr["parameters."+index+"."+k] != rsAttr[k] {
			// ignore data sources where an empty list is being compared against a null list.
			if k[len(k)-1:] == "#" && (dsAttr["parameters."+index+"."+k] == "" || dsAttr["parameters."+index+"."+k] == "0") && (rsAttr[k] == "" || rsAttr[k] == "0") {
				continue
			}
			errMsg += fmt.Sprintf("%s is %s; want %s\n", k, dsAttr["parameters."+index+"."+k], rsAttr[k])
		}
	}

	if errMsg != "" {
		return errors.New(errMsg)
	}

	return nil
}

// This function checks state match for resourceName and asserts the absense of resourceName2 in data source
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

		err := checkFieldsMatchForDataSourceStateAndResourceState(dsAttr, rsAttr, ignoreFields)
		if err != nil {
			return err
		}
		err = checkResourceAbsentInDataSourceAfterFilterApplied(dsAttr, rsAttr2)
		return err
	}
}

// This function asserts the absence of the parameter manager parameter resource which would not be included in the data source list due to the filter applied.
func checkResourceAbsentInDataSourceAfterFilterApplied(dsAttr, rsAttr map[string]string) error {
	totalParameters, err := strconv.Atoi(dsAttr["parameters.#"])
	if err != nil {
		return errors.New("couldn't convert length of parameters list to integer")
	}
	for i := 0; i < totalParameters; i++ {
		if dsAttr["parameters."+strconv.Itoa(i)+".name"] == rsAttr["name"] {
			return errors.New("the resource is present in the data source even after the filter is applied")
		}
	}
	return nil
}
