package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccMonitoringUptimeCheckConfig_update(t *testing.T) {
	t.Parallel()
	project := getTestProjectFromEnv()
	host := "192.168.1.1"

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringUptimeCheckConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringUptimeCheckConfig_update(randString(t, 4), "mypath", "password1", project, host),
			},
			{
				ResourceName:            "google_monitoring_uptime_check_config.http",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"http_check.0.auth_info.0.password"},
			},
			{
				Config: testAccMonitoringUptimeCheckConfig_update(randString(t, 4), "", "password2", project, host),
			},
			{
				ResourceName:            "google_monitoring_uptime_check_config.http",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"http_check.0.auth_info.0.password"},
			},
		},
	})
}

// The second update should force a recreation of the uptime check because 'monitored_resource' isn't
// updatable in place
func TestAccMonitoringUptimeCheckConfig_changeNonUpdatableFields(t *testing.T) {
	t.Parallel()
	project := getTestProjectFromEnv()

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringUptimeCheckConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringUptimeCheckConfig_update(randString(t, 4), "mypath", "password1", project, "192.168.1.1"),
			},
			{
				ResourceName:            "google_monitoring_uptime_check_config.http",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"http_check.0.auth_info.0.password"},
			},
			{
				Config: testAccMonitoringUptimeCheckConfig_update(randString(t, 4), "mypath", "password1", project, "192.168.1.2"),
			},
			{
				ResourceName:            "google_monitoring_uptime_check_config.http",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"http_check.0.auth_info.0.password"},
			},
			{
				Config: testAccMonitoringUptimeCheckConfig_update(randString(t, 4), "mypath", "password2", project, "192.168.1.2"),
			},
			{
				ResourceName:            "google_monitoring_uptime_check_config.http",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"http_check.0.auth_info.0.password"},
			},
		},
	})
}

func testAccMonitoringUptimeCheckConfig_update(suffix, path, project, pwd, host string) string {
	return fmt.Sprintf(`
resource "google_monitoring_uptime_check_config" "http" {
  display_name = "http-uptime-check-%s"
  timeout      = "60s"

  http_check {
    path = "/%s"
    port = "8010"
    auth_info {
      username = "name"
      password = "%s"
    }
  }

  monitored_resource {
    type = "uptime_url"
    labels = {
      project_id = "%s"
      host       = "%s"
    }
  }

  content_matchers {
    content = "example"
  }
}
`, suffix, path, project, pwd, host,
	)
}
