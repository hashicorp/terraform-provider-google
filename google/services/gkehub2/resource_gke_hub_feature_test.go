// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package gkehub2_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccGKEHubFeature_gkehubFeatureFleetObservability(t *testing.T) {
	// VCR fails to handle batched project services
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGKEHubFeatureDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccGKEHubFeature_gkehubFeatureFleetObservability(context),
			},
			{
				ResourceName:      "google_gke_hub_feature.feature",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGKEHubFeature_gkehubFeatureFleetObservabilityUpdate1(context),
			},
			{
				ResourceName:      "google_gke_hub_feature.feature",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGKEHubFeature_gkehubFeatureFleetObservabilityUpdate2(context),
			},
			{
				ResourceName:      "google_gke_hub_feature.feature",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGKEHubFeature_gkehubFeatureFleetObservability(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`
resource "time_sleep" "wait_for_gkehub_enablement" {
  create_duration = "150s"
  depends_on = [google_project_service.gkehub]
}

resource "google_gke_hub_feature" "feature" {
  name = "fleetobservability"
  location = "global"
  project = google_project.project.project_id
  spec {
    fleetobservability {
      logging_config {
        default_config {
    mode = "MOVE"
        }
        fleet_scope_logs_config {
          mode = "COPY"
        }
      }
    }
  }
  depends_on = [time_sleep.wait_for_gkehub_enablement]
}
`, context)
}

func testAccGKEHubFeature_gkehubFeatureFleetObservabilityUpdate1(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`
resource "time_sleep" "wait_for_gkehub_enablement" {
  create_duration = "150s"
  depends_on = [google_project_service.gkehub]
}

resource "google_gke_hub_feature" "feature" {
  name = "fleetobservability"
  location = "global"
  project = google_project.project.project_id
  spec {
    fleetobservability {
      logging_config {
        default_config {
    mode = "MOVE"
        }
      }
    }
  }
  depends_on = [time_sleep.wait_for_gkehub_enablement]
}
`, context)
}

func testAccGKEHubFeature_gkehubFeatureFleetObservabilityUpdate2(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`
resource "time_sleep" "wait_for_gkehub_enablement" {
  create_duration = "150s"
  depends_on = [google_project_service.gkehub]
}

resource "google_gke_hub_feature" "feature" {
  name = "fleetobservability"
  location = "global"
  project = google_project.project.project_id
  spec {
    fleetobservability {
      logging_config {
        fleet_scope_logs_config {
          mode = "COPY"
        }
      }
    }
  }
  depends_on = [time_sleep.wait_for_gkehub_enablement]
}
`, context)
}

func TestAccGKEHubFeature_gkehubFeatureMciUpdate(t *testing.T) {
	// VCR fails to handle batched project services
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGKEHubFeatureDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEHubFeature_gkehubFeatureMciUpdateStart(context),
			},
			{
				ResourceName:            "google_gke_hub_feature.feature",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"update_time"},
			},
			{
				Config: testAccGKEHubFeature_gkehubFeatureMciChangeMembership(context),
			},
			{
				ResourceName:            "google_gke_hub_feature.feature",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"update_time"},
			},
		},
	})
}

func testAccGKEHubFeature_gkehubFeatureMciUpdateStart(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`

resource "google_container_cluster" "primary" {
  name               = "tf-test%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  project = google_project.project.project_id
  depends_on = [google_project_service.mci, google_project_service.container, google_project_service.container, google_project_service.gkehub]
}

resource "google_container_cluster" "secondary" {
  name               = "tf-test2%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  project = google_project.project.project_id
  depends_on = [google_project_service.mci, google_project_service.container, google_project_service.container, google_project_service.gkehub]
}

resource "google_gke_hub_membership" "membership" {
  membership_id = "tf-test%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.primary.id}"
    }
  }
  project = google_project.project.project_id
}

resource "google_gke_hub_membership" "membership_second" {
  membership_id = "tf-test2%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.secondary.id}"
    }
  }
  project = google_project.project.project_id
}

resource "google_gke_hub_feature" "feature" {
  name = "multiclusteringress"
  location = "global"
  spec {
    multiclusteringress {
      config_membership = google_gke_hub_membership.membership.id
    }
  }
  project = google_project.project.project_id
}
`, context)
}

func testAccGKEHubFeature_gkehubFeatureMciChangeMembership(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`
resource "google_container_cluster" "primary" {
  name               = "tf-test%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  project = google_project.project.project_id
  depends_on = [google_project_service.mci, google_project_service.container, google_project_service.container, google_project_service.gkehub]
}

resource "google_container_cluster" "secondary" {
  name               = "tf-test2%{random_suffix}"
  location           = "us-central1-a"
  initial_node_count = 1
  project = google_project.project.project_id
  depends_on = [google_project_service.mci, google_project_service.container, google_project_service.container, google_project_service.gkehub]
}

resource "google_gke_hub_membership" "membership" {
  membership_id = "tf-test%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.primary.id}"
    }
  }
  project = google_project.project.project_id
}

resource "google_gke_hub_membership" "membership_second" {
  membership_id = "tf-test2%{random_suffix}"
  endpoint {
    gke_cluster {
      resource_link = "//container.googleapis.com/${google_container_cluster.secondary.id}"
    }
  }
  project = google_project.project.project_id
}

resource "google_gke_hub_feature" "feature" {
  name = "multiclusteringress"
  location = "global"
  spec {
    multiclusteringress {
      config_membership = google_gke_hub_membership.membership_second.id
    }
  }
  labels = {
    foo = "bar"
  }
  project = google_project.project.project_id
}
`, context)
}

func TestAccGKEHubFeature_gkehubFeatureMcsd(t *testing.T) {
	// VCR fails to handle batched project services
	acctest.SkipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGKEHubFeatureDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccGKEHubFeature_gkehubFeatureMcsd(context),
			},
			{
				ResourceName:            "google_gke_hub_feature.feature",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
			{
				Config: testAccGKEHubFeature_gkehubFeatureMcsdUpdate(context),
			},
			{
				ResourceName:      "google_gke_hub_feature.feature",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGKEHubFeature_gkehubFeatureMcsd(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  name = "multiclusterservicediscovery"
  location = "global"
  project = "projects/${google_project.project.project_id}"
  labels = {
    foo = "bar"
  }
  depends_on = [google_project_service.mcsd]
}
`, context)
}

func testAccGKEHubFeature_gkehubFeatureMcsdUpdate(context map[string]interface{}) string {
	return gkeHubFeatureProjectSetupForGA(context) + acctest.Nprintf(`
resource "google_gke_hub_feature" "feature" {
  name = "multiclusterservicediscovery"
  location = "global"
  project = google_project.project.project_id
  labels = {
    foo = "quux"
    baz = "qux"
  }
  depends_on = [google_project_service.mcsd]
}
`, context)
}

func gkeHubFeatureProjectSetupForGA(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  name            = "tf-test-gkehub%{random_suffix}"
  project_id      = "tf-test-gkehub%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_project_service" "mesh" {
  project = google_project.project.project_id
  service = "meshconfig.googleapis.com"
}

resource "google_project_service" "mci" {
  project = google_project.project.project_id
  service = "multiclusteringress.googleapis.com"
}

resource "google_project_service" "acm" {
  project = google_project.project.project_id
  service = "anthosconfigmanagement.googleapis.com"
}

resource "google_project_service" "mcsd" {
  project = google_project.project.project_id
  service = "multiclusterservicediscovery.googleapis.com"
}

resource "google_project_service" "compute" {
  project = google_project.project.project_id
  service = "compute.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "container" {
  project = google_project.project.project_id
  service = "container.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "gkehub" {
  project = google_project.project.project_id
  service = "gkehub.googleapis.com"
  disable_on_destroy = false
}
`, context)
}

func testAccCheckGKEHubFeatureDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_gke_hub_feature" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{GKEHub2BasePath}}projects/{{project}}/locations/{{location}}/features/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("GKEHubFeature still exists at %s", url)
			}
		}

		return nil
	}
}
