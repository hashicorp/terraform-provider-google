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
	firebaserules "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/firebaserules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"
)

func TestAccFirebaserulesRuleset_BasicRuleset(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  GetTestProjectFromEnv(),
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckFirebaserulesRulesetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFirebaserulesRuleset_BasicRuleset(context),
			},
			{
				ResourceName:      "google_firebaserules_ruleset.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
func TestAccFirebaserulesRuleset_MinimalRuleset(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  GetTestProjectFromEnv(),
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckFirebaserulesRulesetDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccFirebaserulesRuleset_MinimalRuleset(context),
			},
			{
				ResourceName:      "google_firebaserules_ruleset.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccFirebaserulesRuleset_BasicRuleset(context map[string]interface{}) string {
	return Nprintf(`
resource "google_firebaserules_ruleset" "primary" {
  source {
    files {
      content     = "service cloud.firestore {match /databases/{database}/documents { match /{document=**} { allow read, write: if false; } } }"
      name        = "firestore.rules"
      fingerprint = ""
    }

    language = ""
  }

  project = "%{project_name}"
}


`, context)
}

func testAccFirebaserulesRuleset_MinimalRuleset(context map[string]interface{}) string {
	return Nprintf(`
resource "google_firebaserules_ruleset" "primary" {
  source {
    files {
      content = "service cloud.firestore {match /databases/{database}/documents { match /{document=**} { allow read, write: if false; } } }"
      name    = "firestore.rules"
    }
  }

  project = "%{project_name}"
}


`, context)
}

func testAccCheckFirebaserulesRulesetDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_firebaserules_ruleset" {
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

			obj := &firebaserules.Ruleset{
				Project:    dcl.StringOrNil(rs.Primary.Attributes["project"]),
				CreateTime: dcl.StringOrNil(rs.Primary.Attributes["create_time"]),
				Name:       dcl.StringOrNil(rs.Primary.Attributes["name"]),
			}

			client := NewDCLFirebaserulesClient(config, config.UserAgent, billingProject, 0)
			_, err := client.GetRuleset(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_firebaserules_ruleset still exists %v", obj)
			}
		}
		return nil
	}
}
