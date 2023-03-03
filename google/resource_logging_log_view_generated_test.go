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
	logging "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"
)

func TestAccLoggingLogView_BasicHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  GetTestProjectFromEnv(),
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckLoggingLogViewDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLoggingLogView_BasicHandWritten(context),
			},
			{
				ResourceName:            "google_logging_log_view.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bucket"},
			},
		},
	})
}

func testAccLoggingLogView_BasicHandWritten(context map[string]interface{}) string {
	return Nprintf(`
resource "google_logging_log_view" "primary" {
  name        = "tf-test-view%{random_suffix}"
  bucket      = google_logging_project_bucket_config.basic.id
  description = "A logging view configured with Terraform"
  filter      = "SOURCE(\"projects/myproject\") AND resource.type = \"gce_instance\" AND LOG_ID(\"stdout\")"
}

resource "google_logging_project_bucket_config" "basic" {
    project        = "%{project_name}"
    location       = "global"
    retention_days = 30
    bucket_id      = "_Default"
}

`, context)
}

func testAccCheckLoggingLogViewDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_logging_log_view" {
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

			obj := &logging.LogView{
				Bucket:      dcl.String(rs.Primary.Attributes["bucket"]),
				Name:        dcl.String(rs.Primary.Attributes["name"]),
				Description: dcl.String(rs.Primary.Attributes["description"]),
				Filter:      dcl.String(rs.Primary.Attributes["filter"]),
				Location:    dcl.StringOrNil(rs.Primary.Attributes["location"]),
				Parent:      dcl.StringOrNil(rs.Primary.Attributes["parent"]),
				CreateTime:  dcl.StringOrNil(rs.Primary.Attributes["create_time"]),
				UpdateTime:  dcl.StringOrNil(rs.Primary.Attributes["update_time"]),
			}

			client := NewDCLLoggingClient(config, config.UserAgent, billingProject, 0)
			_, err := client.GetLogView(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_logging_log_view still exists %v", obj)
			}
		}
		return nil
	}
}
