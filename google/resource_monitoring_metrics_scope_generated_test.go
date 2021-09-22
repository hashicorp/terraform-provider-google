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
)

func TestAccMonitoringMetricsScope_BasicMetricsScope(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  getTestProjectFromEnv(),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringMetricsScopeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringMetricsScope_BasicMetricsScope(context),
			},
			{
				ResourceName:      "google_monitoring_metrics_scope.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMonitoringMetricsScope_BasicMetricsScope(context map[string]interface{}) string {
	return Nprintf(`
resource "google_monitoring_metrics_scope" "primary" {
  name = "%{project_name}"
}


`, context)
}

func testAccCheckMonitoringMetricsScopeDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_monitoring_metrics_scope" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := googleProviderConfig(t)

			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			obj := &monitoring.MetricsScope{
				Name:       dcl.String(rs.Primary.Attributes["name"]),
				CreateTime: dcl.StringOrNil(rs.Primary.Attributes["create_time"]),
				UpdateTime: dcl.StringOrNil(rs.Primary.Attributes["update_time"]),
			}

			client := NewDCLMonitoringClient(config, config.userAgent, billingProject)
			_, err := client.GetMetricsScope(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_monitoring_metrics_scope still exists %v", obj)
			}
		}
		return nil
	}
}
