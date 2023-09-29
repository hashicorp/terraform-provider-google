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

func TestAccGKEHub2Fleet_gkehubFleetBasicExample_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   acctest.RandString(t, 10),
		"org_id":          envvar.GetTestOrgFromEnv(t),
		"billing_account": envvar.GetTestBillingAccountFromEnv(t),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckGKEHub2FleetDestroyProducer(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccGKEHub2Fleet_basic(context),
			},
			{
				ResourceName:      "google_gke_hub_fleet.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccGKEHub2Fleet_update(context),
			},
			{
				ResourceName:      "google_gke_hub_fleet.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGKEHub2Fleet_basic(context map[string]interface{}) string {
	return gkeHubFleetProjectSetupForGA(context) + acctest.Nprintf(`
resource "google_gke_hub_fleet" "default" {
  project = google_project.project.project_id
  display_name = "my production fleet"

  depends_on = [time_sleep.wait_for_gkehub_enablement]
}
`, context)
}

func testAccGKEHub2Fleet_update(context map[string]interface{}) string {
	return gkeHubFleetProjectSetupForGA(context) + acctest.Nprintf(`
resource "google_gke_hub_fleet" "default" {
  project = google_project.project.project_id
  display_name = "my staging fleet"

  depends_on = [time_sleep.wait_for_gkehub_enablement]
}
`, context)
}

func gkeHubFleetProjectSetupForGA(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_project" "project" {
  name            = "tf-test-gkehub%{random_suffix}"
  project_id      = "tf-test-gkehub%{random_suffix}"
  org_id          = "%{org_id}"
  billing_account = "%{billing_account}"
}

resource "google_project_service" "gkehub" {
  project = google_project.project.project_id
  service = "gkehub.googleapis.com"
  disable_on_destroy = false
}

resource "time_sleep" "wait_for_gkehub_enablement" {
  create_duration = "150s"
  depends_on = [google_project_service.gkehub]
}
`, context)
}

func testAccCheckGKEHub2FleetDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_gke_hub_fleet" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{GKEHub2BasePath}}projects/{{project}}/locations/global/fleets/default")
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
				return fmt.Errorf("GKEHub2Fleet still exists at %s", url)
			}
		}

		return nil
	}
}
