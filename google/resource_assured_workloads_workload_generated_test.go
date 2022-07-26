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
	assuredworkloads "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/assuredworkloads"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"
)

func TestAccAssuredWorkloadsWorkload_BasicHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_acct":  getTestBillingAccountFromEnv(t),
		"org_id":        getTestOrgFromEnv(t),
		"region":        getTestRegionFromEnv(),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAssuredWorkloadsWorkloadDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAssuredWorkloadsWorkload_BasicHandWritten(context),
			},
			{
				ResourceName:            "google_assured_workloads_workload.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"billing_account", "kms_settings", "resource_settings", "provisioned_resources_parent"},
			},
			{
				Config: testAccAssuredWorkloadsWorkload_BasicHandWrittenUpdate0(context),
			},
			{
				ResourceName:            "google_assured_workloads_workload.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"billing_account", "kms_settings", "resource_settings", "provisioned_resources_parent"},
			},
		},
	})
}
func TestAccAssuredWorkloadsWorkload_FullHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"billing_acct":  getTestBillingAccountFromEnv(t),
		"org_id":        getTestOrgFromEnv(t),
		"region":        getTestRegionFromEnv(),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAssuredWorkloadsWorkloadDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAssuredWorkloadsWorkload_FullHandWritten(context),
			},
			{
				ResourceName:            "google_assured_workloads_workload.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"billing_account", "kms_settings", "resource_settings", "provisioned_resources_parent"},
			},
		},
	})
}

func testAccAssuredWorkloadsWorkload_BasicHandWritten(context map[string]interface{}) string {
	return Nprintf(`
resource "google_assured_workloads_workload" "primary" {
  display_name = "tf-test-name%{random_suffix}"
  labels = {
    a = "a"
  }
  billing_account = "billingAccounts/%{billing_acct}"
  compliance_regime = "FEDRAMP_MODERATE"
  provisioned_resources_parent = google_folder.folder1.name
  organization = "%{org_id}"
  location = "us-central1"
}

resource "google_folder" "folder1" {
  display_name = "tf-test-name%{random_suffix}"
  parent       = "organizations/%{org_id}"
}
`, context)
}

func testAccAssuredWorkloadsWorkload_BasicHandWrittenUpdate0(context map[string]interface{}) string {
	return Nprintf(`
resource "google_assured_workloads_workload" "primary" {
  display_name = "tf-test-name%{random_suffix}"
  labels = {
    a = "b"
  }
  billing_account = "billingAccounts/%{billing_acct}"
  compliance_regime = "FEDRAMP_MODERATE"
  provisioned_resources_parent = google_folder.folder1.name
  organization = "%{org_id}"
  location = "us-central1"
}

resource "google_folder" "folder1" {
  display_name = "tf-test-name%{random_suffix}"
  parent       = "organizations/%{org_id}"
}
`, context)
}

func testAccAssuredWorkloadsWorkload_FullHandWritten(context map[string]interface{}) string {
	return Nprintf(`
resource "google_assured_workloads_workload" "primary" {
  display_name = "tf-test-name%{random_suffix}"
  billing_account = "billingAccounts/%{billing_acct}"
  compliance_regime = "FEDRAMP_MODERATE"
  organization = "%{org_id}"
  location = "us-central1"
  kms_settings {
    next_rotation_time = "2022-10-02T15:01:23Z"
    rotation_period = "864000s"
  }
  provisioned_resources_parent = google_folder.folder1.name
}

resource "google_folder" "folder1" {
  display_name = "tf-test-name%{random_suffix}"
  parent       = "organizations/%{org_id}"
}

`, context)
}

func testAccCheckAssuredWorkloadsWorkloadDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_assured_workloads_workload" {
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

			obj := &assuredworkloads.Workload{
				BillingAccount:             dcl.String(rs.Primary.Attributes["billing_account"]),
				ComplianceRegime:           assuredworkloads.WorkloadComplianceRegimeEnumRef(rs.Primary.Attributes["compliance_regime"]),
				DisplayName:                dcl.String(rs.Primary.Attributes["display_name"]),
				Location:                   dcl.String(rs.Primary.Attributes["location"]),
				Organization:               dcl.String(rs.Primary.Attributes["organization"]),
				ProvisionedResourcesParent: dcl.String(rs.Primary.Attributes["provisioned_resources_parent"]),
				CreateTime:                 dcl.StringOrNil(rs.Primary.Attributes["create_time"]),
				Name:                       dcl.StringOrNil(rs.Primary.Attributes["name"]),
			}

			client := NewDCLAssuredWorkloadsClient(config, config.userAgent, billingProject, 0)
			_, err := client.GetWorkload(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_assured_workloads_workload still exists %v", obj)
			}
		}
		return nil
	}
}
