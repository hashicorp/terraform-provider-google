// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package parametermanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceParameterManagerParameter_basic(t *testing.T) {
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
				Config: testAccDataSourceParameterManagerParameter_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState(
						"data.google_parameter_manager_parameter.parameter-datasource",
						"google_parameter_manager_parameter.parameter",
					),
				),
			},
		},
	})
}

func testAccDataSourceParameterManagerParameter_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_parameter" "parameter" {
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

data "google_parameter_manager_parameter" "parameter-datasource" {
  parameter_id = google_parameter_manager_parameter.parameter.parameter_id
}
`, context)
}
