// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package composer_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComposerUserWorkloadsSecret_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"env_name":    fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t)),
		"secret_name": fmt.Sprintf("%s-%d", testComposerUserWorkloadsSecretPrefix, acctest.RandInt(t)),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComposerUserWorkloadsSecret_basic(context),
				Check: resource.ComposeTestCheckFunc(
					checkSecretDataSourceMatchesResource(),
				),
			},
		},
	})
}

func checkSecretDataSourceMatchesResource() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources["data.google_composer_user_workloads_secret.test"]
		if !ok {
			return fmt.Errorf("can't find %s in state", "data.google_composer_user_workloads_secret.test")
		}
		rs, ok := s.RootModule().Resources["google_composer_user_workloads_secret.test"]
		if !ok {
			return fmt.Errorf("can't find %s in state", "google_composer_user_workloads_secret.test")
		}

		dsAttr := ds.Primary.Attributes
		rsAttr := rs.Primary.Attributes
		errMsg := ""

		for k := range rsAttr {
			if k == "%" || k == "data.%" {
				continue
			}
			// ignore diff if it's due to secrets being masked.
			if strings.HasPrefix(k, "data.") {
				if _, ok := dsAttr[k]; !ok {
					errMsg += fmt.Sprintf("%s is defined in resource and not in datasource\n", k)
				}
				if dsAttr[k] == "**********" {
					continue
				}
			}
			if dsAttr[k] != rsAttr[k] {
				errMsg += fmt.Sprintf("%s is %s; want %s\n", k, dsAttr[k], rsAttr[k])
			}
		}

		if errMsg != "" {
			return errors.New(errMsg)
		}

		return nil
	}
}

func testAccDataSourceComposerUserWorkloadsSecret_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_composer_environment" "test" {
  name   = "%{env_name}"
  config {
    software_config {
      image_version = "composer-3-airflow-2"
    }
  }
}
resource "google_composer_user_workloads_secret" "test" {
  environment = google_composer_environment.test.name
  name = "%{secret_name}"
  data = {
    username: base64encode("username"),
    password: base64encode("password"),
  }
}
data "google_composer_user_workloads_secret" "test" {
  name        = google_composer_user_workloads_secret.test.name
  environment = google_composer_environment.test.name
}
`, context)
}
