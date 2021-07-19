package google

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	assuredworkloads "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/assuredworkloads"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				ImportStateVerifyIgnore: []string{"billing_account", "kms_settings", "resource_settings", "provisioned_resources_parent"},
			},
			{
				Config: testAccAssuredWorkloadsWorkload_basicUpdate(context),
				Check:  deleteAssuredWorkloadProvisionedResources(t),
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
				Check:  resource.ComposeTestCheckFunc(deleteAssuredWorkloadProvisionedResources(t)),
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

// deleteAssuredWorkloadProvisionedResources deletes the resources provisioned by
// assured workloads.. this is needed in order to delete the parent resource
func deleteAssuredWorkloadProvisionedResources(t *testing.T) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)
		timeout := *schema.DefaultTimeout(4 * time.Minute)
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "google_assured_workloads_workload" {
				continue
			}
			resourceAttributes := rs.Primary.Attributes
			n, err := strconv.Atoi(resourceAttributes["resources.#"])
			log.Printf("[DEBUG]found %v resources\n", n)
			log.Println(resourceAttributes)
			if err != nil {
				return err
			}

			// first delete the projects
			for i := 0; i < n; i++ {
				typee := resourceAttributes[fmt.Sprintf("resources.%d.resource_type", i)]
				if !strings.Contains(typee, "PROJECT") {
					continue
				}
				resource_id := resourceAttributes[fmt.Sprintf("resources.%d.resource_id", i)]
				log.Printf("[DEBUG] searching for project %s\n", resource_id)
				err := retryTimeDuration(func() (reqErr error) {
					_, reqErr = config.NewResourceManagerClient(config.userAgent).Projects.Get(resource_id).Do()
					return reqErr
				}, timeout)
				if err != nil {
					log.Printf("[DEBUG] did not find project %sn", resource_id)
					continue
				}
				log.Printf("[DEBUG] found project %s\n", resource_id)

				err = retryTimeDuration(func() error {
					_, delErr := config.NewResourceManagerClient(config.userAgent).Projects.Delete(resource_id).Do()
					return delErr
				}, timeout)
				if err != nil {
					log.Printf("Error deleting project '%s': %s\n ", resource_id, err)
					continue
				}
				log.Printf("[DEBUG] deleted project %s\n", resource_id)
			}

			// Then delete the folders
			for i := 0; i < n; i++ {
				typee := resourceAttributes[fmt.Sprintf("resources.%d.resource_type", i)]
				if typee != "CONSUMER_FOLDER" {
					continue
				}
				resource_id := "folders/" + resourceAttributes[fmt.Sprintf("resources.%d.resource_id", i)]
				err := retryTimeDuration(func() error {
					var reqErr error
					_, reqErr = config.NewResourceManagerV2Client(config.userAgent).Folders.Get(resource_id).Do()
					return reqErr
				}, timeout)
				log.Printf("[DEBUG] searching for folder %s\n", resource_id)
				if err != nil {
					log.Printf("[DEBUG] did not find folder %sn", resource_id)
					continue
				}
				log.Printf("[DEBUG] found folder %s\n", resource_id)
				err = retryTimeDuration(func() error {
					_, reqErr := config.NewResourceManagerV2Client(config.userAgent).Folders.Delete(resource_id).Do()
					return reqErr
				}, timeout)
				if err != nil {
					return fmt.Errorf("Error deleting folder '%s': %s\n ", resource_id, err)
				}
				log.Printf("[DEBUG] deleted folder %s\n", resource_id)
			}
		}
		return nil
	}
}

func testAccAssuredWorkloadsWorkload_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_assured_workloads_workload" "meep" {
  display_name = "workloadExample-%{random_suffix}"
  labels = {
    a = "a"
  }
  billing_account = "billingAccounts/%{billing_account}"
  compliance_regime = "FEDRAMP_MODERATE"
  provisioned_resources_parent = google_folder.folder1.name
  organization = "%{org_id}"
  location = "us-central1"
}

resource "google_folder" "folder1" {
  display_name = "tf-test-%{random_suffix}"
  parent       = "organizations/%{org_id}"
}
`, context)
}

func testAccAssuredWorkloadsWorkload_basicUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_assured_workloads_workload" "meep" {
  display_name = "updatedExample-%{random_suffix}"
  labels = {
    a = "b"
  }
  billing_account = "billingAccounts/%{billing_account}"
  compliance_regime = "FEDRAMP_MODERATE"
  provisioned_resources_parent = google_folder.folder1.name
  organization = "%{org_id}"
  location = "us-central1"
}

resource "google_folder" "folder1" {
  display_name = "tf-test-%{random_suffix}"
  parent       = "organizations/%{org_id}"
}
`, context)
}

func testAccAssuredWorkloadsWorkload_full(context map[string]interface{}) string {
	return Nprintf(`
resource "google_assured_workloads_workload" "meep" {
  display_name = "workloadExample-%{random_suffix}"
  billing_account = "billingAccounts/%{billing_account}"
  compliance_regime = "FEDRAMP_MODERATE"
  organization = "%{org_id}"
  location = "us-central1"
  kms_settings {
    next_rotation_time = "2021-10-02T15:01:23Z"
    rotation_period = "864000s"
  }
  provisioned_resources_parent = google_folder.folder1.name
}

resource "google_folder" "folder1" {
  display_name = "tf-test-%{random_suffix}"
  parent       = "organizations/%{org_id}"
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
