// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudbuild_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

// Tests schema version migration by creating a trigger with an old version of the provider (4.30.0)
// and then updating it with the current version the provider.
func TestAccCloudBuildTrigger_migration(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()
	name := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))

	oldVersion := map[string]resource.ExternalProvider{
		"google": {
			VersionConstraint: "4.30.0", // a version that doesn't support location yet.
			Source:            "registry.terraform.io/hashicorp/google",
		},
	}
	newVersion := map[string]func() (*schema.Provider, error){
		"mynewprovider": func() (*schema.Provider, error) { return acctest.TestAccProviders["google"], nil },
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: testAccCheckCloudBuildTriggerDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:            configWithFilename(name),
				ExternalProviders: oldVersion,
			},
			{
				ResourceName:            "google_cloudbuild_trigger.simple-trigger",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"location"},
				ExternalProviders:       oldVersion,
			},
			{
				Config:            newConfigWithFilename(name),
				ProviderFactories: newVersion,
			},
			{
				ResourceName:      "google_cloudbuild_trigger.simple-trigger",
				ImportState:       true,
				ImportStateVerify: true,
				ProviderFactories: newVersion,
			},
		},
	})
}

func configWithFilename(name string) string {
	return fmt.Sprintf(`
	resource "google_cloudbuild_trigger" "simple-trigger" {
		trigger_template {
		  branch_name = "main"
		  repo_name   = "my-repo"
		}
		substitutions = {
		  _FOO = "bar"
		  _BAZ = "qux"
		}
		name = "%s"
		filename = "oldfile.yaml"
	}
	`, name)
}

func newConfigWithFilename(name string) string {
	return fmt.Sprintf(`
	provider "mynewprovider" {}

	resource "google_cloudbuild_trigger" "simple-trigger" {
		provider = mynewprovider
		trigger_template {
		  branch_name = "main"
		  repo_name   = "my-repo"
		}
		substitutions = {
		  _FOO = "bar"
		  _BAZ = "qux"
		}
		name = "%s"
		filename = "newfile.yaml"
	}
	`, name)
}
