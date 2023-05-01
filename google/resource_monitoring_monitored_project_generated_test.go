// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package google

import (
	"context"
	"fmt"
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	monitoring "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/monitoring"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccMonitoringMonitoredProject_BasicMonitoredProject(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        acctest.GetTestOrgFromEnv(t),
		"project_name":  acctest.GetTestProjectFromEnv(),
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckMonitoringMonitoredProjectDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringMonitoredProject_BasicMonitoredProject(context),
			},
			{
				ResourceName:      "google_monitoring_monitored_project.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMonitoringMonitoredProject_BasicMonitoredProject(context map[string]interface{}) string {
	return Nprintf(`
resource "google_monitoring_monitored_project" "primary" {
  metrics_scope = "%{project_name}"
  name          = google_project.basic.name
}

resource "google_project" "basic" {
  project_id = "tf-test-id%{random_suffix}"
  name       = "tf-test-id%{random_suffix}"
  org_id     = "%{org_id}"
}


`, context)
}

func testAccCheckMonitoringMonitoredProjectDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_monitoring_monitored_project" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := GoogleProviderConfig(t)

			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			obj := &monitoring.MonitoredProject{
				MetricsScope: dcl.String(rs.Primary.Attributes["metrics_scope"]),
				Name:         dcl.String(rs.Primary.Attributes["name"]),
				CreateTime:   dcl.StringOrNil(rs.Primary.Attributes["create_time"]),
			}

			client := transport_tpg.NewDCLMonitoringClient(config, config.UserAgent, billingProject, 0)
			_, err := client.GetMonitoredProject(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_monitoring_monitored_project still exists %v", obj)
			}
		}
		return nil
	}
}
