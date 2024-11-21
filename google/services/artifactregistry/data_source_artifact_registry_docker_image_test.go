// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package artifactregistry_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceArtifactRegistryDockerImage(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	resourceName := "data.google_artifact_registry_docker_image.test"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceArtifactRegistryDockerImageConfig,
				Check: resource.ComposeTestCheckFunc(
					// Data source using a tag
					checkTaggedDataSources(resourceName+"Tag", "latest"),
					resource.TestCheckResourceAttrSet(resourceName+"Tag", "image_size_bytes"),
					validateTimeStamps(resourceName+"Tag"),

					// url safe docker name using a tag
					checkTaggedDataSources(resourceName+"UrlTag", "latest"),

					// Data source using no tag or digest
					resource.TestCheckResourceAttrSet(resourceName+"None", "repository_id"),
					resource.TestCheckResourceAttrSet(resourceName+"None", "image_name"),
					resource.TestCheckResourceAttrSet(resourceName+"None", "name"),
					resource.TestCheckResourceAttrSet(resourceName+"None", "self_link"),
				),
			},
		},
	})
}

// Test the data source against the public AR repos
// https://console.cloud.google.com/artifacts/docker/cloudrun/us/container
// https://console.cloud.google.com/artifacts/docker/go-containerregistry/us/gcr.io
// Currently, gcr.io does not provide a imageSizeBytes or buildTime field in the JSON response
const testAccDataSourceArtifactRegistryDockerImageConfig = `
data "google_artifact_registry_docker_image" "testTag" {
	project       = "cloudrun"
	location      = "us"
	repository_id = "container"
	image_name    = "hello:latest"
}

data "google_artifact_registry_docker_image" "testUrlTag" {
	project       = "go-containerregistry"
	location      = "us"
	repository_id = "gcr.io"
	image_name    = "krane/debug:latest"
}

data "google_artifact_registry_docker_image" "testNone" {
	project       = "go-containerregistry"
	location      = "us"
	repository_id = "gcr.io"
	image_name    = "crane"
}
`

func checkTaggedDataSources(resourceName string, expectedTag string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "repository_id"),
		resource.TestCheckResourceAttrSet(resourceName, "image_name"),
		resource.TestCheckResourceAttrSet(resourceName, "name"),
		resource.TestCheckResourceAttrSet(resourceName, "self_link"),
		resource.TestCheckTypeSetElemAttr(resourceName, "tags.*", expectedTag),
		resource.TestCheckResourceAttrSet(resourceName, "media_type"),
	)
}

func checkDigestDataSources(resourceName string, expectedName string, expectedSelfLink string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "repository_id"),
		resource.TestCheckResourceAttrSet(resourceName, "image_name"),
		resource.TestCheckResourceAttr(resourceName, "name", expectedName),
		resource.TestCheckResourceAttr(resourceName, "self_link", expectedSelfLink),
		resource.TestCheckResourceAttrSet(resourceName, "media_type"),
	)
}

func validateTimeStamps(dataSourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// check that the timestamps are RFC3339
		ds, ok := s.RootModule().Resources[dataSourceName]
		if !ok {
			return fmt.Errorf("can't find %s in state", dataSourceName)
		}

		if !isRFC3339(ds.Primary.Attributes["upload_time"]) {
			return fmt.Errorf("upload_time is not RFC3339: %s", ds.Primary.Attributes["upload_time"])
		}

		if !isRFC3339(ds.Primary.Attributes["build_time"]) {
			return fmt.Errorf("build_time is not RFC3339: %s", ds.Primary.Attributes["build_time"])
		}

		if !isRFC3339(ds.Primary.Attributes["update_time"]) {
			return fmt.Errorf("update_time is not RFC3339: %s", ds.Primary.Attributes["update_time"])
		}

		return nil
	}
}

func isRFC3339(s string) bool {
	_, err := time.Parse(time.RFC3339, s)
	return err == nil
}
