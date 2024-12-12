// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package secretmanager_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceSecretManagerSecrets_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecretManagerSecrets_basic(context),
				Check: resource.ComposeTestCheckFunc(
					checkListDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_secret_manager_secrets.foo",
						"google_secret_manager_secret.foo",
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

func testAccDataSourceSecretManagerSecrets_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  add_terraform_attribution_label = false
}

resource "google_secret_manager_secret" "foo" {
  secret_id = "tf-test-secret-%{random_suffix}"

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
    }
  }

  labels = {
    label = "my-label"
  }

  annotations = {
    key1 = "value1"
  }

  version_destroy_ttl = "360000s"
}

data "google_secret_manager_secrets" "foo" {
  depends_on = [
    google_secret_manager_secret.foo
  ]
}
`, context)
}

func TestAccDataSourceSecretManagerSecrets_filter(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecretManagerSecrets_filter(context),
				Check: resource.ComposeTestCheckFunc(
					checkListDataSourceStateMatchesResourceStateWithIgnoresForAppliedFilter(
						"data.google_secret_manager_secrets.foo",
						"google_secret_manager_secret.foo",
						"google_secret_manager_secret.bar",
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

func testAccDataSourceSecretManagerSecrets_filter(context map[string]interface{}) string {
	return acctest.Nprintf(`
provider "google" {
  add_terraform_attribution_label = false
}

resource "google_secret_manager_secret" "foo" {
  secret_id = "tf-test-secret-%{random_suffix}"

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
    }
  }

  labels = {
    label = "my-label"
  }

  annotations = {
    key1 = "value1"
  }
}

resource "google_secret_manager_secret" "bar" {
  secret_id = "tf-test-secret-2-%{random_suffix}"
  
  replication {
    user_managed {
      replicas {
        location = "us-east5"
      }
    }
  }

  labels = {
    label= "my-label2"
  }

  annotations = {
    key1 = "value1" 
  }
}

data "google_secret_manager_secrets" "foo" {
  filter = "replication.user_managed.replicas.location:us-central1"
  depends_on = [
    google_secret_manager_secret.foo,
    google_secret_manager_secret.bar
  ]
}
`, context)
}

// This function checks data source state matches for resourceName secret manager secret state
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

// This function checks whether all the attributes of the secret manager secret resource and the attributes of the secret manager secret inside the data source list are the same
func checkFieldsMatchForDataSourceStateAndResourceState(dsAttr, rsAttr map[string]string, ignoreFields map[string]struct{}) error {
	totalSecrets, err := strconv.Atoi(dsAttr["secrets.#"])
	if err != nil {
		return errors.New("Couldn't convert length of secrets list to integer")
	}
	index := "-1"
	for i := 0; i < totalSecrets; i++ {
		if dsAttr["secrets."+strconv.Itoa(i)+".name"] == rsAttr["name"] {
			index = strconv.Itoa(i)
		}
	}

	if index == "-1" {
		return errors.New("The newly created secret is not found in the data source")
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
		if dsAttr["secrets."+index+"."+k] != rsAttr[k] {
			// ignore data sources where an empty list is being compared against a null list.
			if k[len(k)-1:] == "#" && (dsAttr["secrets."+index+"."+k] == "" || dsAttr["secrets."+index+"."+k] == "0") && (rsAttr[k] == "" || rsAttr[k] == "0") {
				continue
			}
			errMsg += fmt.Sprintf("%s is %s; want %s\n", k, dsAttr["secrets."+index+"."+k], rsAttr[k])
		}
	}

	if errMsg != "" {
		return errors.New(errMsg)
	}

	return nil
}

// This function checks state match for resourceName and asserts the absence of resourceName2 in data source
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

// This function asserts the absence of the secret manager secret resource which would not be included in the data source list due to the filter applied.
func checkResourceAbsentInDataSourceAfterFilterApplied(dsAttr, rsAttr map[string]string) error {
	totalSecrets, err := strconv.Atoi(dsAttr["secrets.#"])
	if err != nil {
		return errors.New("Couldn't convert length of secrets list to integer")
	}
	for i := 0; i < totalSecrets; i++ {
		if dsAttr["secrets."+strconv.Itoa(i)+".name"] == rsAttr["name"] {
			return errors.New("The resource is present in the data source even after the filter is applied")
		}
	}
	return nil
}
