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
	dataplex "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/dataplex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"
)

func TestAccDataplexLake_BasicLake(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  GetTestProjectFromEnv(),
		"region":        GetTestRegionFromEnv(),
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckDataplexLakeDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataplexLake_BasicLake(context),
			},
			{
				ResourceName:      "google_dataplex_lake.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccDataplexLake_BasicLakeUpdate0(context),
			},
			{
				ResourceName:      "google_dataplex_lake.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDataplexLake_BasicLake(context map[string]interface{}) string {
	return Nprintf(`
resource "google_dataplex_lake" "primary" {
  location     = "%{region}"
  name         = "tf-test-lake%{random_suffix}"
  description  = "Lake for DCL"
  display_name = "Lake for DCL"

  labels = {
    my-lake = "exists"
  }

  project = "%{project_name}"
}


`, context)
}

func testAccDataplexLake_BasicLakeUpdate0(context map[string]interface{}) string {
	return Nprintf(`
resource "google_dataplex_lake" "primary" {
  location     = "%{region}"
  name         = "tf-test-lake%{random_suffix}"
  description  = "Updated description for lake"
  display_name = "Lake for DCL"

  labels = {
    my-lake = "exists"
  }

  project = "%{project_name}"
}


`, context)
}

func testAccCheckDataplexLakeDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_dataplex_lake" {
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

			obj := &dataplex.Lake{
				Location:       dcl.String(rs.Primary.Attributes["location"]),
				Name:           dcl.String(rs.Primary.Attributes["name"]),
				Description:    dcl.String(rs.Primary.Attributes["description"]),
				DisplayName:    dcl.String(rs.Primary.Attributes["display_name"]),
				Project:        dcl.StringOrNil(rs.Primary.Attributes["project"]),
				CreateTime:     dcl.StringOrNil(rs.Primary.Attributes["create_time"]),
				ServiceAccount: dcl.StringOrNil(rs.Primary.Attributes["service_account"]),
				State:          dataplex.LakeStateEnumRef(rs.Primary.Attributes["state"]),
				Uid:            dcl.StringOrNil(rs.Primary.Attributes["uid"]),
				UpdateTime:     dcl.StringOrNil(rs.Primary.Attributes["update_time"]),
			}

			client := NewDCLDataplexClient(config, config.UserAgent, billingProject, 0)
			_, err := client.GetLake(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_dataplex_lake still exists %v", obj)
			}
		}
		return nil
	}
}
