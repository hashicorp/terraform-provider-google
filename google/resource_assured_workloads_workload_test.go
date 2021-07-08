package google

import (
	"context"
	"fmt"
	"strings"
	"testing"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	assuredworkloads "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/assuredworkloads"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAssuredWorkloadsWorkload_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   randString(t, 10),
		"org_id":          getTestOrgFromEnv(t),
		"billing_account": getTestBillingAccountFromEnv(t),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAssuredWorkloadsWorkloadDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAssuredWorkloadsWorkload_basic(context),
			},
			{
				ResourceName:            "google_assured_workloads_workload.meep",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"billing_account"},
			},
			{
				Config: testAccAssuredWorkloadsWorkload_basicUpdate(context),
			},
			{
				ResourceName:            "google_assured_workloads_workload.meep",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"billing_account"},
			},
		},
	})
}

func TestAccAssuredWorkloadsWorkload_full(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix":   randString(t, 10),
		"org_id":          getTestOrgFromEnv(t),
		"billing_account": getTestBillingAccountFromEnv(t),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAssuredWorkloadsWorkloadDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAssuredWorkloadsWorkload_full(context),
			},
			{
				ResourceName:            "google_assured_workloads_workload.meep",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"billing_account", "kms_settings", "resource_settings", "provisioned_resources_parent"},
			},
		},
	})
}

func testAccAssuredWorkloadsWorkload_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_assured_workloads_workload" "meep" {
	display_name = "workloadExample"
	labels = {
		a = "a"
	}
	billing_account = "billingAccounts/%{billing_account}"
	compliance_regime = "FEDRAMP_MODERATE"
	organization = "%{org_id}"
	location = "us-central1"
}
`, context)
}

func testAccAssuredWorkloadsWorkload_basicUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_assured_workloads_workload" "meep" {
	display_name = "updatedExample"
	labels = {
		a = "b"
	}
	billing_account = "billingAccounts/%{billing_account}"
	compliance_regime = "FEDRAMP_MODERATE"
	organization = "%{org_id}"
	location = "us-central1"
}
`, context)
}

func testAccAssuredWorkloadsWorkload_full(context map[string]interface{}) string {
	return Nprintf(`
resource "google_assured_workloads_workload" "meep" {
	display_name = "workloadExample"
	billing_account = "billingAccounts/%{billing_account}"
	compliance_regime = "FEDRAMP_MODERATE"
	organization = "%{org_id}"
	location = "us-central1"
	kms_settings {
		next_rotation_time = "2021-10-02T15:01:23Z"
		rotation_period = "864000s"
	}
	provisioned_resources_parent = "folders/177863664720"
	resource_settings {
		resource_id = "tf-test-prj-%{random_suffix}"
		resource_type = "CONSUMER_PROJECT"
	}
}
`, context)
}

func testAccCheckAssuredWorkloadsWorkloadDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_assured_workloads_workload" {
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

			billingAccount := rs.Primary.Attributes["billing_account"]
			location := rs.Primary.Attributes["location"]
			name := rs.Primary.Attributes["name"]

			obj := &assuredworkloads.Workload{
				BillingAccount: dcl.String(rs.Primary.Attributes["billing_account"]),
				Location:       dcl.String(rs.Primary.Attributes["location"]),
				Name:           dcl.StringOrNil(rs.Primary.Attributes["name"]),
			}

			client := NewDCLAssuredWorkloadsClient(config, config.userAgent, billingProject)
			_, err := client.GetWorkload(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("AssuredWorkloadsWorkloadResource still exists at %s, %s, %s", billingAccount, location, name)
			}
		}

		return nil
	}
}
