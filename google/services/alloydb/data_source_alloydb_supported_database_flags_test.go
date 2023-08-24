// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package alloydb_test

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceAlloydbSupportedDatabaseFlags_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAlloydbSupportedDatabaseFlags_basic(context),
				Check: resource.ComposeTestCheckFunc(
					validateAlloydbSupportedDatabaseFlagsResult(
						"data.google_alloydb_supported_database_flags.qa",
					),
				),
			},
		},
	})
}

func testAccDataSourceAlloydbSupportedDatabaseFlags_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_alloydb_supported_database_flags" "qa" {
	location = "us-central1"
}
`, context)
}

func validateAlloydbSupportedDatabaseFlagsResult(dataSourceName string) func(*terraform.State) error {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", dataSourceName)
		}

		var dsAttr map[string]string
		dsAttr = ds.Primary.Attributes

		totalFlags, err := strconv.Atoi(dsAttr["supported_database_flags.#"])
		if err != nil {
			return errors.New("Couldn't convert length of flags list to integer")
		}
		if totalFlags == 0 {
			return errors.New("No supported database flags are fetched from location 'us-central1'")
		}
		for i := 0; i < totalFlags; i++ {
			if dsAttr["supported_database_flags."+strconv.Itoa(i)+".name"] == "" {
				return errors.New("name parameter is not set for the flag")
			}
			if dsAttr["supported_database_flags."+strconv.Itoa(i)+".flag_name"] == "" {
				return errors.New("flag_name parameter is not set for the flag")
			}
			if len(dsAttr["supported_database_flags."+strconv.Itoa(i)+".string_restrictions"]) > 0 && len(dsAttr["supported_database_flags."+strconv.Itoa(i)+".integer_restrictions"]) > 0 {
				return errors.New("Both string restriction and integer restriction cannot be set for a union restriction field")
			}
		}
		return nil
	}
}
