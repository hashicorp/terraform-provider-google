package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Destroy and recreate the mapping, testing that Terraform doesn't return a 409
func TestAccCloudRunDomainMapping_foregroundDeletion(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"namespace":     getTestProjectFromEnv(),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckCloudRunDomainMappingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudRunDomainMapping_cloudRunDomainMappingUpdated1(context),
			},
			{
				ResourceName:            "google_cloud_run_domain_mapping.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "status", "metadata.0.resource_version"},
			},
			{
				Config: testAccCloudRunDomainMapping_cloudRunDomainMappingUpdated2(context),
			},
			{
				ResourceName:            "google_cloud_run_domain_mapping.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "status", "metadata.0.resource_version"},
			},
		},
	})
}

func testAccCloudRunDomainMapping_cloudRunDomainMappingUpdated1(context map[string]interface{}) string {
	return Nprintf(`
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
	return Nprintf(`
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
