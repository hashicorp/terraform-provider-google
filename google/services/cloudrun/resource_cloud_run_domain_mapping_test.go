// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudrun_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

// Destroy and recreate the mapping, testing that Terraform doesn't return a 409
func TestAccCloudRunDomainMapping_foregroundDeletion(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"namespace":     envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudRunDomainMappingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunDomainMapping_cloudRunDomainMappingUpdated1(context),
			},
			{
				ResourceName:            "google_cloud_run_domain_mapping.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "status", "metadata.0.labels", "metadata.0.terraform_labels", "metadata.0.resource_version"},
			},
			{
				Config: testAccCloudRunDomainMapping_cloudRunDomainMappingUpdated2(context),
			},
			{
				ResourceName:            "google_cloud_run_domain_mapping.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "status", "metadata.0.labels", "metadata.0.terraform_labels", "metadata.0.resource_version"},
			},
		},
	})
}

func testAccCloudRunDomainMapping_cloudRunDomainMappingUpdated1(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_run_service" "default" {
    name     = "tf-test-cloudrun-srv%{random_suffix}"
    location = "us-central1"

    metadata {
      namespace = "%{namespace}"
    }

    template {
      spec {
        containers {
          image = "us-docker.pkg.dev/cloudrun/container/hello"
        }
      }
    }
  }

resource "google_cloud_run_domain_mapping" "default" {
  location = "us-central1"
  name     = "tf-test-domain%{random_suffix}.gcp.tfacc.hashicorptest.com"

  metadata {
    namespace = "%{namespace}"
  }

  spec {
    route_name = google_cloud_run_service.default.name
  }
}
`, context)
}

func testAccCloudRunDomainMapping_cloudRunDomainMappingUpdated2(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_run_service" "default" {
  name     = "tf-test-cloudrun-srv%{random_suffix}"
  location = "us-central1"
  metadata {
    namespace = "%{namespace}"
  }
  template {
    spec {
      containers {
        image = "us-docker.pkg.dev/cloudrun/container/hello"
      }
    }
  }
}
resource "google_cloud_run_domain_mapping" "default" {
  location = "us-central1"
  name     = "tf-test-domain%{random_suffix}.gcp.tfacc.hashicorptest.com"
  metadata {
    namespace = "%{namespace}"
    labels = {
      "my-label" = "my-value"
    }
  }
  spec {
    route_name = google_cloud_run_service.default.name
  }
}
`, context)
}

func TestAccCloudRunDomainMapping_migration(t *testing.T) {
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"namespace":     envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	oldVersion := map[string]resource.ExternalProvider{
		"google": {
			VersionConstraint: "4.84.0", // a version that doesn't separate user defined labels and system labels
			Source:            "registry.terraform.io/hashicorp/google",
		},
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:     func() { acctest.AccTestPreCheck(t) },
		CheckDestroy: testAccCheckCloudRunDomainMappingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:            testAccCloudRunDomainMapping_cloudRunDomainMappingUpdated2(context),
				ExternalProviders: oldVersion,
			},
			{
				Config:                   testAccCloudRunDomainMapping_cloudRunDomainMappingUpdated2(context),
				ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
			},
		},
	})
}
