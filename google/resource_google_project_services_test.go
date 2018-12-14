package google

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// Test that services can be enabled and disabled on a project
func TestAccProjectServices_basic(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	services1 := []string{"logging.googleapis.com", "cloudresourcemanager.googleapis.com"}
	services2 := []string{"cloudresourcemanager.googleapis.com"}
	oobService := "logging.googleapis.com"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project with some services
			{
				Config: testAccProjectAssociateServicesBasic(services1, pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testProjectServicesMatch(services1, pid),
				),
			},
			// Update services to remove one
			{
				Config: testAccProjectAssociateServicesBasic(services2, pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testProjectServicesMatch(services2, pid),
				),
			},
			// Add a service out-of-band and ensure it is removed
			{
				PreConfig: func() {
					config := testAccProvider.Meta().(*Config)
					enableService(oobService, pid, config)
				},
				Config: testAccProjectAssociateServicesBasic(services2, pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testProjectServicesMatch(services2, pid),
				),
			},
			{
				ResourceName:            "google_project_services.acceptance",
				ImportState:             true,
				ImportStateId:           pid,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"disable_on_destroy"},
			},
		},
	})
}

// Test that services are authoritative when a project has existing
// sevices not represented in config
func TestAccProjectServices_authoritative(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	services := []string{"cloudresourcemanager.googleapis.com"}
	oobService := "logging.googleapis.com"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project with no services
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectExists("google_project.acceptance", pid),
				),
			},
			// Add a service out-of-band, then apply a config that creates a service.
			// It should remove the out-of-band service.
			{
				PreConfig: func() {
					config := testAccProvider.Meta().(*Config)
					enableService(oobService, pid, config)
				},
				Config: testAccProjectAssociateServicesBasic(services, pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testProjectServicesMatch(services, pid),
				),
			},
		},
	})
}

// Test that services are authoritative when a project has existing
// sevices, some which are represented in the config and others
// that are not
func TestAccProjectServices_authoritative2(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	oobServices := []string{"logging.googleapis.com", "cloudresourcemanager.googleapis.com"}
	services := []string{"logging.googleapis.com"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project with no services
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleProjectExists("google_project.acceptance", pid),
				),
			},
			// Add a service out-of-band, then apply a config that creates a service.
			// It should remove the out-of-band service.
			{
				PreConfig: func() {
					config := testAccProvider.Meta().(*Config)
					for _, s := range oobServices {
						enableService(s, pid, config)
					}
				},
				Config: testAccProjectAssociateServicesBasic(services, pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testProjectServicesMatch(services, pid),
				),
			},
		},
	})
}

// Test that services that can't be enabled on their own (such as dataproc-control.googleapis.com)
// don't end up causing diffs when they are enabled as a side-effect of a different service's
// enablement.
func TestAccProjectServices_ignoreUnenablableServices(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	billingId := getTestBillingAccountFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	services := []string{
		"dataproc.googleapis.com",
		// The following services are enabled as a side-effect of dataproc's enablement
		"storage-component.googleapis.com",
		"deploymentmanager.googleapis.com",
		"replicapool.googleapis.com",
		"replicapoolupdater.googleapis.com",
		"resourceviews.googleapis.com",
		"compute.googleapis.com",
		"container.googleapis.com",
		"containerregistry.googleapis.com",
		"storage-api.googleapis.com",
		"pubsub.googleapis.com",
		"oslogin.googleapis.com",
		"bigquery-json.googleapis.com",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectAssociateServicesBasic_withBilling(services, pid, pname, org, billingId),
				Check: resource.ComposeTestCheckFunc(
					testProjectServicesMatch(services, pid),
				),
			},
		},
	})
}

func TestAccProjectServices_pagination(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	billingId := getTestBillingAccountFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)

	// we need at least 50 services (doesn't matter what they are) to exercise the
	// pagination handling code.
	services := []string{
		"actions.googleapis.com",
		"appengine.googleapis.com",
		"appengineflex.googleapis.com",
		"bigquery-json.googleapis.com",
		"bigquerydatatransfer.googleapis.com",
		"bigtableadmin.googleapis.com",
		"bigtabletableadmin.googleapis.com",
		"cloudbuild.googleapis.com",
		"clouderrorreporting.googleapis.com",
		"cloudfunctions.googleapis.com",
		"cloudiot.googleapis.com",
		"cloudkms.googleapis.com",
		"cloudmonitoring.googleapis.com",
		"cloudresourcemanager.googleapis.com",
		"cloudtrace.googleapis.com",
		"compute.googleapis.com",
		"container.googleapis.com",
		"containerregistry.googleapis.com",
		"dataflow.googleapis.com",
		"dataproc.googleapis.com",
		"datastore.googleapis.com",
		"deploymentmanager.googleapis.com",
		"dialogflow.googleapis.com",
		"dns.googleapis.com",
		"endpoints.googleapis.com",
		"firebaserules.googleapis.com",
		"firestore.googleapis.com",
		"genomics.googleapis.com",
		"iam.googleapis.com",
		"iamcredentials.googleapis.com",
		"language.googleapis.com",
		"logging.googleapis.com",
		"ml.googleapis.com",
		"monitoring.googleapis.com",
		"oslogin.googleapis.com",
		"pubsub.googleapis.com",
		"replicapool.googleapis.com",
		"replicapoolupdater.googleapis.com",
		"resourceviews.googleapis.com",
		"runtimeconfig.googleapis.com",
		"servicecontrol.googleapis.com",
		"servicemanagement.googleapis.com",
		"sourcerepo.googleapis.com",
		"spanner.googleapis.com",
		"speech.googleapis.com",
		"sql-component.googleapis.com",
		"storage-api.googleapis.com",
		"storage-component.googleapis.com",
		"storagetransfer.googleapis.com",
		"testing.googleapis.com",
		"toolresults.googleapis.com",
		"translate.googleapis.com",
		"videointelligence.googleapis.com",
		"vision.googleapis.com",
		"zync.googleapis.com",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectAssociateServicesBasic_withBilling(services, pid, pname, org, billingId),
				Check: resource.ComposeTestCheckFunc(
					testProjectServicesMatch(services, pid),
				),
			},
		},
	})
}

func testAccProjectAssociateServicesBasic(services []string, pid, name, org string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}
resource "google_project_services" "acceptance" {
  project            = "${google_project.acceptance.project_id}"
  services           = [%s]
  disable_on_destroy = true
}
`, pid, name, org, testStringsToString(services))
}

func testAccProjectAssociateServicesBasic_withBilling(services []string, pid, name, org, billing string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id      = "%s"
  name            = "%s"
  org_id          = "%s"
  billing_account = "%s"
}
resource "google_project_services" "acceptance" {
  project            = "${google_project.acceptance.project_id}"
  services           = [%s]
  disable_on_destroy = false
}
`, pid, name, org, billing, testStringsToString(services))
}

func testProjectServicesMatch(services []string, pid string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		apiServices, err := getApiServices(pid, config, ignoreProjectServices)
		if err != nil {
			return fmt.Errorf("Error listing services for project %q: %v", pid, err)
		}

		sort.Strings(services)
		sort.Strings(apiServices)
		if !reflect.DeepEqual(services, apiServices) {
			return fmt.Errorf("Services in config (%v) do not exactly match services returned by API (%v)", services, apiServices)
		}

		return nil
	}
}

func testStringsToString(s []string) string {
	var b bytes.Buffer
	for i, v := range s {
		b.WriteString(fmt.Sprintf("\"%s\"", v))
		if i < len(s)-1 {
			b.WriteString(",")
		}
	}
	r := b.String()
	log.Printf("[DEBUG]: Converted list of strings to %s", r)
	return b.String()
}
