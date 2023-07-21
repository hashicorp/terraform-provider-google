// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package containeranalysis_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestDataSourceGoogleContainerRegistryRepository(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_container_registry_repository.test"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleContainerRegistryRepo_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttr(resourceName, "repository_url", "bar.gcr.io/foo"),
					resource.TestCheckResourceAttrSet(resourceName+"Scoped", "project"),
					resource.TestCheckResourceAttr(resourceName+"Scoped", "repository_url", "bar.gcr.io/example.com/foo"),
				),
			},
		},
	})
}

const testAccCheckGoogleContainerRegistryRepo_basic = `
data "google_container_registry_repository" "test" {
	project = "foo"
	region = "bar"
}
data "google_container_registry_repository" "testScoped" {
	project = "example.com:foo"
	region = "bar"
}
`

func TestDataSourceGoogleContainerRegistryImage(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_container_registry_image.test"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleContainerRegistryImage_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "project"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttr(resourceName, "image_url", "bar.gcr.io/foo/baz"),
					resource.TestCheckResourceAttr(resourceName+"2", "image_url", "bar.gcr.io/foo/baz:qux"),
					resource.TestCheckResourceAttr(resourceName+"3", "image_url", "bar.gcr.io/foo/baz@1234"),
					resource.TestCheckResourceAttrSet(resourceName+"Scoped", "project"),
					resource.TestCheckResourceAttr(resourceName+"Scoped", "image_url", "bar.gcr.io/example.com/foo/baz:qux"),
				),
			},
		},
	})
}

const testAccCheckGoogleContainerRegistryImage_basic = `
data "google_container_registry_image" "test" {
  project = "foo"
  region  = "bar"
  name    = "baz"
}

data "google_container_registry_image" "test2" {
  project = "foo"
  region  = "bar"
  name    = "baz"
  tag     = "qux"
}

data "google_container_registry_image" "test3" {
  project = "foo"
  region  = "bar"
  name    = "baz"
  digest  = "1234"
}

data "google_container_registry_image" "testScoped" {
  project = "example.com:foo"
  region  = "bar"
  name    = "baz"
  tag     = "qux"
}
`
