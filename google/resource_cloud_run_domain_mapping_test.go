package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestDomainMappingLabelDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		K, Old, New        string
		ExpectDiffSuppress bool
	}{
		"missing run.googleapis.com/overrideAt": {
			K:                  "metadata.0.labels.run.googleapis.com/overrideAt",
			Old:                "2021-04-20T22:38:23.584Z",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"explicit run.googleapis.com/overrideAt": {
			K:                  "metadata.0.labels.run.googleapis.com/overrideAt",
			Old:                "2021-04-20T22:38:23.584Z",
			New:                "2022-04-20T22:38:23.584Z",
			ExpectDiffSuppress: false,
		},
		"missing cloud.googleapis.com/location": {
			K:                  "metadata.0.labels.cloud.googleapis.com/location",
			Old:                "us-central1",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"explicit cloud.googleapis.com/location": {
			K:                  "metadata.0.labels.cloud.googleapis.com/location",
			Old:                "us-central1",
			New:                "us-central2",
			ExpectDiffSuppress: false,
		},
		"labels.%": {
			K:                  "metadata.0.labels.%",
			Old:                "3",
			New:                "1",
			ExpectDiffSuppress: true,
		},
		"deleted custom key": {
			K:                  "metadata.0.labels.my-label",
			Old:                "my-value",
			New:                "",
			ExpectDiffSuppress: false,
		},
		"added custom key": {
			K:                  "metadata.0.labels.my-label",
			Old:                "",
			New:                "my-value",
			ExpectDiffSuppress: false,
		},
	}
	for tn, tc := range cases {
		if domainMappingLabelDiffSuppress(tc.K, tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Errorf("bad: %s, %q: %q => %q expect DiffSuppress to return %t", tn, tc.K, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

// Destroy and recreate the mapping, testing that Terraform doesn't return a 409
func TestAccCloudRunDomainMapping_foregroundDeletion(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"namespace":     getTestProjectFromEnv(),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
