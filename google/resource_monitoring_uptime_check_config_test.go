package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMonitoringUptimeCheckConfig_update(t *testing.T) {
	t.Parallel()
	project := GetTestProjectFromEnv()
	host := "192.168.1.1"

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckMonitoringUptimeCheckConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringUptimeCheckConfig_update(RandString(t, 4), "mypath", "password1", project, host),
			},
			{
				ResourceName:            "google_monitoring_uptime_check_config.http",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"http_check.0.auth_info.0.password"},
			},
			{
				Config: testAccMonitoringUptimeCheckConfig_update(RandString(t, 4), "", "password2", project, host),
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
	project := GetTestProjectFromEnv()

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckMonitoringUptimeCheckConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringUptimeCheckConfig_update(RandString(t, 4), "mypath", "password1", project, "192.168.1.1"),
			},
			{
				ResourceName:            "google_monitoring_uptime_check_config.http",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"http_check.0.auth_info.0.password"},
			},
			{
				Config: testAccMonitoringUptimeCheckConfig_update(RandString(t, 4), "mypath", "password1", project, "192.168.1.2"),
			},
			{
				ResourceName:            "google_monitoring_uptime_check_config.http",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"http_check.0.auth_info.0.password"},
			},
			{
				Config: testAccMonitoringUptimeCheckConfig_update(RandString(t, 4), "mypath", "password2", project, "192.168.1.2"),
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

func TestAccMonitoringUptimeCheckConfig_jsonPathUpdate(t *testing.T) {
	t.Parallel()
	project := GetTestProjectFromEnv()
	host := "192.168.1.1"
	suffix := RandString(t, 4)

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckMonitoringUptimeCheckConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringUptimeCheckConfig_jsonPathUpdate(suffix, project, host, "123", "$.path", "EXACT_MATCH"),
			},
			{
				ResourceName:      "google_monitoring_uptime_check_config.http",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMonitoringUptimeCheckConfig_jsonPathUpdate(suffix, project, host, "content", "$.different", "REGEX_MATCH"),
			},
			{
				ResourceName:      "google_monitoring_uptime_check_config.http",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMonitoringUptimeCheckConfig_update(suffix, path, pwd, project, host string) string {
	return fmt.Sprintf(`
resource "google_monitoring_uptime_check_config" "http" {
  display_name = "http-uptime-check-%s"
  timeout      = "60s"

  http_check {
    path = "/%s"
    port = "8010"
    request_method = "GET"
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
    matcher = "CONTAINS_STRING"
  }
}
`, suffix, path, pwd, project, host,
	)
}

func testAccMonitoringUptimeCheckConfig_jsonPathUpdate(suffix, project, host, content, json_path, json_path_matcher string) string {
	return fmt.Sprintf(`
resource "google_monitoring_uptime_check_config" "http" {
  display_name = "http-uptime-check-%s"
  timeout      = "60s"

  http_check {
    path = "a-path"
    port = "80"
    request_method = "GET"
  }

  monitored_resource {
    type = "uptime_url"
    labels = {
      project_id = "%s"
      host       = "%s"
    }
  }

  content_matchers {
    content = "%s"
    matcher = "MATCHES_JSON_PATH"
	json_path_matcher {
		json_path = "%s"
		json_matcher = "%s"
	}
  }
}
`, suffix, project, host, content, json_path, json_path_matcher,
	)
}
