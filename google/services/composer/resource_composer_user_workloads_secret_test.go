// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package composer_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/services/composer"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testComposerUserWorkloadsSecretPrefix = "tf-test-composer-secret"

func TestAccComposerUserWorkloadsSecret_basic(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	secretName := fmt.Sprintf("%s-%d", testComposerUserWorkloadsSecretPrefix, acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		// CheckDestroy:               testAccComposerUserWorkloadsSecretDestroy(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerUserWorkloadsSecret_basic(envName, secretName, envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_composer_user_workloads_secret.test", "data.username"),
					resource.TestCheckResourceAttrSet("google_composer_user_workloads_secret.test", "data.password"),
				),
			},
			{
				ResourceName: "google_composer_user_workloads_secret.test",
				ImportState:  true,
			},
		},
	})
}

func TestAccComposerUserWorkloadsSecret_update(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	secretName := fmt.Sprintf("%s-%d", testComposerUserWorkloadsSecretPrefix, acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerUserWorkloadsSecret_basic(envName, secretName, envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv()),
			},
			{
				Config: testAccComposerUserWorkloadsSecret_update(envName, secretName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("google_composer_user_workloads_secret.test", "data.email"),
					resource.TestCheckResourceAttrSet("google_composer_user_workloads_secret.test", "data.password"),
					resource.TestCheckNoResourceAttr("google_composer_user_workloads_secret.test", "data.username"),
				),
			},
		},
	})
}

func TestAccComposerUserWorkloadsSecret_delete(t *testing.T) {
	t.Parallel()

	envName := fmt.Sprintf("%s-%d", testComposerEnvironmentPrefix, acctest.RandInt(t))
	secretName := fmt.Sprintf("%s-%d", testComposerUserWorkloadsSecretPrefix, acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComposerUserWorkloadsSecret_basic(envName, secretName, envvar.GetTestProjectFromEnv(), envvar.GetTestRegionFromEnv()),
			},
			{
				Config: testAccComposerUserWorkloadsSecret_delete(envName),
				Check: resource.ComposeTestCheckFunc(
					testAccComposerUserWorkloadsSecretDestroyed(t),
				),
			},
		},
	})
}

func testAccComposerUserWorkloadsSecret_basic(envName, secretName, project, region string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  config {
    software_config {
      image_version = "composer-3-airflow-2"
    }
  }
}
resource "google_composer_user_workloads_secret" "test" {
  environment = google_composer_environment.test.name
  name = "%s"
  project = "%s"
  region = "%s"
  data = {
    username: base64encode("username"),
    password: base64encode("password"),
  }
}
`, envName, secretName, project, region)
}

func testAccComposerUserWorkloadsSecret_update(envName, secretName string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  config {
    software_config {
      image_version = "composer-3-airflow-2"
    }
  }
}
resource "google_composer_user_workloads_secret" "test" {
  environment = google_composer_environment.test.name
  name = "%s"
  data = {
		email:    base64encode("email"),
    password: base64encode("password"),
  }
}
`, envName, secretName)
}

func testAccComposerUserWorkloadsSecret_delete(envName string) string {
	return fmt.Sprintf(`
resource "google_composer_environment" "test" {
  name   = "%s"
  config {
    software_config {
      image_version = "composer-3-airflow-2"
    }
  }
}
`, envName)
}

func testAccComposerUserWorkloadsSecretDestroyed(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_composer_user_workloads_secret" {
				continue
			}

			idTokens := strings.Split(rs.Primary.ID, "/")
			if len(idTokens) != 8 {
				return fmt.Errorf("Invalid ID %q, expected format projects/{project}/regions/{region}/environments/{environment}/userWorkloadsSecrets/{name}", rs.Primary.ID)
			}
			secretName := &composer.UserWorkloadsSecretsName{
				Project:     idTokens[1],
				Region:      idTokens[3],
				Environment: idTokens[5],
				Secret:      idTokens[7],
			}

			_, err := config.NewComposerClient(config.UserAgent).Projects.Locations.Environments.UserWorkloadsSecrets.Get(secretName.ResourceName()).Do()
			if err == nil {
				return fmt.Errorf("secret %s still exists", secretName.ResourceName())
			}
		}

		return nil
	}
}
