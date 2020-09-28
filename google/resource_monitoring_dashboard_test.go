package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccMonitoringDashboard_basic(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringDashboardDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringDashboard_basic(),
			},
			{
				ResourceName:      "google_monitoring_dashboard.dashboard",
				ImportState:       true,
				ImportStateVerify: true,
				// Default import format uses the ID, which contains the project #
				// Testing import formats with the project name don't work because we set
				// the ID on import to what the user specified, which won't match the ID
				// from the apply
				ImportStateVerifyIgnore: []string{"project"},
			},
		},
	})
}

func TestAccMonitoringDashboard_gridLayout(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringDashboardDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringDashboard_gridLayout(),
			},
			{
				ResourceName:            "google_monitoring_dashboard.dashboard",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
		},
	})
}

func TestAccMonitoringDashboard_rowLayout(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringDashboardDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringDashboard_rowLayout(),
			},
			{
				ResourceName:            "google_monitoring_dashboard.dashboard",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
		},
	})
}

func TestAccMonitoringDashboard_update(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMonitoringDashboardDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccMonitoringDashboard_rowLayout(),
			},
			{
				ResourceName:            "google_monitoring_dashboard.dashboard",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
			{
				Config: testAccMonitoringDashboard_basic(),
			},
			{
				ResourceName:            "google_monitoring_dashboard.dashboard",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
			{
				Config: testAccMonitoringDashboard_gridLayout(),
			},
			{
				ResourceName:            "google_monitoring_dashboard.dashboard",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
		},
	})
}

func testAccCheckMonitoringDashboardDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_monitoring_dashboard" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := googleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{MonitoringBasePath}}v1/{{name}}")
			if err != nil {
				return err
			}

			_, err = sendRequest(config, "GET", "", url, config.userAgent, nil, isMonitoringConcurrentEditError)
			if err == nil {
				return fmt.Errorf("MonitoringDashboard still exists at %s", url)
			}
		}

		return nil
	}
}

func testAccMonitoringDashboard_basic() string {
	return fmt.Sprintf(`
resource "google_monitoring_dashboard" "dashboard" {
  dashboard_json = <<EOF
{
  "displayName": "Demo Dashboard",
  "gridLayout": {
    "widgets": [
      {
        "blank": {}
      }
    ]
  }
}

EOF
}
`)
}

func testAccMonitoringDashboard_gridLayout() string {
	return fmt.Sprintf(`
resource "google_monitoring_dashboard" "dashboard" {
  dashboard_json = <<EOF
{
  "displayName": "Grid Layout Example",
  "gridLayout": {
    "columns": "2",
    "widgets": [
      {
        "title": "Widget 1",
        "xyChart": {
          "dataSets": [{
            "timeSeriesQuery": {
              "timeSeriesFilter": {
                "filter": "metric.type=\"agent.googleapis.com/nginx/connections/accepted_count\"",
                "aggregation": {
                  "perSeriesAligner": "ALIGN_RATE"
                }
              },
              "unitOverride": "1"
            },
            "plotType": "LINE"
          }],
          "timeshiftDuration": "0s",
          "yAxis": {
            "label": "y1Axis",
            "scale": "LINEAR"
          }
        }
      },
      {
        "text": {
          "content": "Widget 2",
          "format": "MARKDOWN"
        }
      },
      {
        "title": "Widget 3",
        "xyChart": {
          "dataSets": [{
            "timeSeriesQuery": {
              "timeSeriesFilter": {
                "filter": "metric.type=\"agent.googleapis.com/nginx/connections/accepted_count\"",
                "aggregation": {
                  "perSeriesAligner": "ALIGN_RATE"
                }
              },
              "unitOverride": "1"
            },
            "plotType": "STACKED_BAR"
          }],
          "timeshiftDuration": "0s",
          "yAxis": {
            "label": "y1Axis",
            "scale": "LINEAR"
          }
        }
      }
    ]
  }
}

EOF
}
`)
}

func testAccMonitoringDashboard_rowLayout() string {
	return fmt.Sprintf(`
resource "google_monitoring_dashboard" "dashboard" {
  dashboard_json = <<EOF
{
  "displayName": "Row Layout Example",
  "rowLayout": {
    "rows": [
      {
        "weight": "1",
        "widgets": [
          {
            "text": {
              "content": "Widget 1",
              "format": "MARKDOWN"
            }
          },
          {
            "text": {
              "content": "Widget 3",
              "format": "MARKDOWN"
            }
          },
          {
            "text": {
              "content": "Widget 2",
              "format": "MARKDOWN"
            }
          }
        ]
      }
    ]
  }
}

EOF
}
`)
}
