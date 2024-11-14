// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package composer_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComposerUserWorkloadsConfigMap_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"env_name":        fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t)),
		"config_map_name": fmt.Sprintf("tf-test-composer-config-map-%d", acctest.RandInt(t)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComposerUserWorkloadsConfigMap_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_composer_user_workloads_config_map.test",
						"google_composer_user_workloads_config_map.test"),
				),
			},
		},
	})
}

func testAccDataSourceComposerUserWorkloadsConfigMap_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_composer_environment" "test" {
  name   = "%{env_name}"
  config {
    software_config {
      image_version = "composer-3-airflow-2"
    }
  }
}
resource "google_composer_user_workloads_config_map" "test" {
  environment = google_composer_environment.test.name
  name = "%{config_map_name}"
  data = {
    db_host: "dbhost:5432",
    api_host: "apihost:443",
  }
}
data "google_composer_user_workloads_config_map" "test" {
  name        = google_composer_user_workloads_config_map.test.name
  environment = google_composer_environment.test.name
}
`, context)
}
