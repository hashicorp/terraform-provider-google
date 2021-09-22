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
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("MonitoringMetrics_scope", &resource.Sweeper{
		Name: "MonitoringMetrics_scope",
		F:    testSweepMonitoringMetrics_scope,
	})
}
func testSweepMonitoringMetrics_scope(region string) error {
	log.Print("[INFO][SWEEPER_LOG] No-op sweeper called for undeletable MonitoringMetrics_scope")
	return nil
}
