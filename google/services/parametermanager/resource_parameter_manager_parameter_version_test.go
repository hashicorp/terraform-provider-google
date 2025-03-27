// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package parametermanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccParameterManagerParameterVersion_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckParameterManagerParameterVersionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccParameterManagerParameterVersion_basic(context),
			},
			{
				ResourceName:            "google_parameter_manager_parameter_version.parameter-version-update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parameter", "parameter_version_id"},
			},
			{
				Config: testAccParameterManagerParameterVersion_update(context),
			},
			{
				ResourceName:            "google_parameter_manager_parameter_version.parameter-version-update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parameter", "parameter_version_id"},
			},
			{
				Config: testAccParameterManagerParameterVersion_basic(context),
			},
			{
				ResourceName:            "google_parameter_manager_parameter_version.parameter-version-update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"parameter", "parameter_version_id"},
			},
		},
	})
}

func testAccParameterManagerParameterVersion_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_parameter" "parameter-update" {
  parameter_id = "tf_test_parameter%{random_suffix}"
}

resource "google_parameter_manager_parameter_version" "parameter-version-update" {
  parameter = google_parameter_manager_parameter.parameter-update.id
  parameter_version_id = "tf_test_parameter_version%{random_suffix}"
  parameter_data = "parameter-version-data"
}
`, context)
}

func testAccParameterManagerParameterVersion_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_parameter_manager_parameter" "parameter-update" {
  parameter_id = "tf_test_parameter%{random_suffix}"
}

resource "google_parameter_manager_parameter_version" "parameter-version-update" {
  parameter = google_parameter_manager_parameter.parameter-update.id
  parameter_version_id = "tf_test_parameter_version%{random_suffix}"
  parameter_data = "parameter-version-data"
  disabled = true
}
`, context)
}
