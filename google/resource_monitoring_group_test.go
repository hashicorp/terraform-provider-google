package google

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func init() {
	resource.AddTestSweepers("gcp_monitoring_group", &resource.Sweeper{
		Name: "gcp_monitoring_group",
		F:    testSweepMonitoringGroups,
	})
}

func testSweepMonitoringGroups(region string) error {
	project := getTestProjectFromEnv()
	config, err := sharedConfigForRegion(region)
	if err != nil {
		log.Fatalf("error getting shared config for region: %s", err)
	}

	err = config.LoadAndValidate()
	if err != nil {
		log.Fatalf("error loading: %s", err)
	}

	url := fmt.Sprintf("%sprojects/%s/groups", config.MonitoringBasePath, project)
	res, err := sendRequest(config, "GET", project, url, nil)
	if err != nil {
		log.Fatalf("Unable to list Monitoring Groups: %s", err)
	}

	groups, ok := res["group"]
	if !ok {
		log.Fatalf("Fatal - no groups found in Monitoring Groups response")
	}
	gs := groups.([]interface{})

	for _, gi := range gs {
		g := gi.(map[string]interface{})

		// Only sweep monitoring groups with the test prefix
		if g["name"] != nil && strings.HasPrefix(g["name"].(string), "tf-test") {
			url := fmt.Sprintf("%s%s", config.MonitoringBasePath, g["name"].(string))
			log.Printf("Sweeping Monitoring Group: %s", g["name"].(string))

			_, err = sendRequest(config, "DELETE", project, url, nil)
			if err != nil {
				log.Printf("Error deleting monitoring group: %s", err)
			}
		}
	}

	return nil
}

func TestAccMonitoringGroup_update(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringGroup_update("europe-west1"),
			},
			{
				ResourceName:      "google_monitoring_group.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMonitoringGroup_update("europe-west2"),
			},
			{
				ResourceName:      "google_monitoring_group.update",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMonitoringGroup_update(zone string) string {
	return fmt.Sprintf(`
resource "google_monitoring_group" "update" {
  display_name = "tf-test Integration Test Group"

  filter = "resource.metadata.region=\"%s\""
}
`, zone,
	)
}
