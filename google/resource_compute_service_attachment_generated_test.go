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
	compute "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/compute"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"
)

func TestAccComputeServiceAttachment_BasicHandWritten(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        getTestOrgFromEnv(t),
		"project_name":  getTestProjectFromEnv(),
		"region":        getTestRegionFromEnv(),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeServiceAttachmentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServiceAttachment_BasicHandWritten(context),
			},
			{
				ResourceName:            "google_compute_service_attachment.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"nat_subnets.0", "consumer_accept_lists.0.project_id_or_num", "consumer_reject_lists.0"},
			},
			{
				Config: testAccComputeServiceAttachment_BasicHandWrittenUpdate0(context),
			},
			{
				ResourceName:            "google_compute_service_attachment.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"nat_subnets.0", "consumer_accept_lists.0.project_id_or_num", "consumer_reject_lists.0"},
			},
			{
				Config: testAccComputeServiceAttachment_BasicHandWrittenUpdate1(context),
			},
			{
				ResourceName:            "google_compute_service_attachment.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"nat_subnets.0", "consumer_accept_lists.0.project_id_or_num", "consumer_reject_lists.0"},
			},
		},
	})
}

func testAccComputeServiceAttachment_BasicHandWritten(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_service_attachment" "primary" {
  connection_preference = "ACCEPT_MANUAL"
  name                  = "tf-test-test-attachment%{random_suffix}"
  nat_subnets           = [google_compute_subnetwork.first_private_service_connect.self_link]
  target_service        = google_compute_forwarding_rule.first_internal.self_link

  consumer_accept_lists {
    project_id_or_num = google_project.second.name
    connection_limit  = 2
  }

  consumer_reject_lists = [google_project.first.name]
  description           = "A sample service attachment"
  enable_proxy_protocol = false
  project               = "%{project_name}"
  region                = "%{region}"
}

resource "google_compute_forwarding_rule" "first_internal" {
  name                  = "tf-test-test-rule1%{random_suffix}"
  all_ports             = true
  backend_service       = google_compute_region_backend_service.internal.self_link
  description           = "A test forwarding rule with internal load balancing scheme"
  load_balancing_scheme = "INTERNAL"
  network               = google_compute_network.basic.self_link
  network_tier          = "PREMIUM"
  project               = "%{project_name}"
  region                = "%{region}"
  subnetwork            = google_compute_subnetwork.private.self_link
}

resource "google_compute_region_backend_service" "internal" {
  name                  = "tf-test-test-service%{random_suffix}"
  load_balancing_scheme = "INTERNAL"
  region                = "%{region}"
  network               = google_compute_network.basic.self_link
  project               = "%{project_name}"
}

resource "google_compute_subnetwork" "first_private_service_connect" {
  ip_cidr_range = "10.2.0.0/16"
  name          = "tf-test-compute-psc-subnetwork1%{random_suffix}"
  network       = google_compute_network.basic.self_link
  project       = "%{project_name}"
  purpose       = "PRIVATE_SERVICE_CONNECT"
  region        = "%{region}"
}

resource "google_compute_subnetwork" "second_private_service_connect" {
  ip_cidr_range = "10.3.0.0/16"
  name          = "tf-test-compute-psc-subnetwork2%{random_suffix}"
  network       = google_compute_network.basic.self_link
  project       = "%{project_name}"
  purpose       = "PRIVATE_SERVICE_CONNECT"
  region        = "%{region}"
}

resource "google_compute_subnetwork" "private" {
  ip_cidr_range = "10.4.0.0/16"
  name          = "tf-test-compute-private-subnetwork%{random_suffix}"
  network       = google_compute_network.basic.self_link
  project       = "%{project_name}"
  purpose       = "PRIVATE"
  region        = "%{region}"
}

resource "google_compute_network" "basic" {
  name                    = "tf-test-compute-network%{random_suffix}"
  auto_create_subnetworks = false
  project                 = "%{project_name}"
}

resource "google_project" "first" {
  project_id = "tf-test-test-id1%{random_suffix}"
  name       = "tf-test-test-id1%{random_suffix}"
  org_id     = "%{org_id}"
}

resource "google_project" "second" {
  project_id = "tf-test-test-id2%{random_suffix}"
  name       = "tf-test-test-id2%{random_suffix}"
  org_id     = "%{org_id}"
}

`, context)
}

func testAccComputeServiceAttachment_BasicHandWrittenUpdate0(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_service_attachment" "primary" {
  connection_preference = "ACCEPT_MANUAL"
  name                  = "tf-test-test-attachment%{random_suffix}"
  nat_subnets           = [google_compute_subnetwork.first_private_service_connect.self_link]
  target_service        = google_compute_forwarding_rule.first_internal.self_link

  consumer_accept_lists {
    project_id_or_num = google_project.first.name
    connection_limit  = 3
  }

  consumer_reject_lists = [google_project.second.name]
  description           = "A sample service attachment"
  enable_proxy_protocol = false
  project               = "%{project_name}"
  region                = "%{region}"
}

resource "google_compute_forwarding_rule" "first_internal" {
  name                  = "tf-test-test-rule1%{random_suffix}"
  all_ports             = true
  backend_service       = google_compute_region_backend_service.internal.self_link
  description           = "A test forwarding rule with internal load balancing scheme"
  load_balancing_scheme = "INTERNAL"
  network               = google_compute_network.basic.self_link
  network_tier          = "PREMIUM"
  project               = "%{project_name}"
  region                = "%{region}"
  subnetwork            = google_compute_subnetwork.private.self_link
}

resource "google_compute_region_backend_service" "internal" {
  name                  = "tf-test-test-service%{random_suffix}"
  load_balancing_scheme = "INTERNAL"
  region                = "%{region}"
  network               = google_compute_network.basic.self_link
  project               = "%{project_name}"
}

resource "google_compute_subnetwork" "first_private_service_connect" {
  ip_cidr_range = "10.2.0.0/16"
  name          = "tf-test-compute-psc-subnetwork1%{random_suffix}"
  network       = google_compute_network.basic.self_link
  project       = "%{project_name}"
  purpose       = "PRIVATE_SERVICE_CONNECT"
  region        = "%{region}"
}

resource "google_compute_subnetwork" "second_private_service_connect" {
  ip_cidr_range = "10.3.0.0/16"
  name          = "tf-test-compute-psc-subnetwork2%{random_suffix}"
  network       = google_compute_network.basic.self_link
  project       = "%{project_name}"
  purpose       = "PRIVATE_SERVICE_CONNECT"
  region        = "%{region}"
}

resource "google_compute_subnetwork" "private" {
  ip_cidr_range = "10.4.0.0/16"
  name          = "tf-test-compute-private-subnetwork%{random_suffix}"
  network       = google_compute_network.basic.self_link
  project       = "%{project_name}"
  purpose       = "PRIVATE"
  region        = "%{region}"
}

resource "google_compute_network" "basic" {
  name                    = "tf-test-compute-network%{random_suffix}"
  auto_create_subnetworks = false
  project                 = "%{project_name}"
}

resource "google_project" "first" {
  project_id = "tf-test-test-id1%{random_suffix}"
  name       = "tf-test-test-id1%{random_suffix}"
  org_id     = "%{org_id}"
}

resource "google_project" "second" {
  project_id = "tf-test-test-id2%{random_suffix}"
  name       = "tf-test-test-id2%{random_suffix}"
  org_id     = "%{org_id}"
}

`, context)
}

func testAccComputeServiceAttachment_BasicHandWrittenUpdate1(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_service_attachment" "primary" {
  connection_preference = "ACCEPT_AUTOMATIC"
  name                  = "tf-test-test-attachment%{random_suffix}"
  nat_subnets           = [google_compute_subnetwork.first_private_service_connect.self_link]
  target_service        = google_compute_forwarding_rule.first_internal.self_link

  consumer_reject_lists = []
  description           = "A sample service attachment"
  enable_proxy_protocol = false
  project               = "%{project_name}"
  region                = "%{region}"
}

resource "google_compute_forwarding_rule" "first_internal" {
  name                  = "tf-test-test-rule1%{random_suffix}"
  all_ports             = true
  backend_service       = google_compute_region_backend_service.internal.self_link
  description           = "A test forwarding rule with internal load balancing scheme"
  load_balancing_scheme = "INTERNAL"
  network               = google_compute_network.basic.self_link
  network_tier          = "PREMIUM"
  project               = "%{project_name}"
  region                = "%{region}"
  subnetwork            = google_compute_subnetwork.private.self_link
}

resource "google_compute_region_backend_service" "internal" {
  name                  = "tf-test-test-service%{random_suffix}"
  load_balancing_scheme = "INTERNAL"
  region                = "%{region}"
  network               = google_compute_network.basic.self_link
  project               = "%{project_name}"
}

resource "google_compute_subnetwork" "first_private_service_connect" {
  ip_cidr_range = "10.2.0.0/16"
  name          = "tf-test-compute-psc-subnetwork1%{random_suffix}"
  network       = google_compute_network.basic.self_link
  project       = "%{project_name}"
  purpose       = "PRIVATE_SERVICE_CONNECT"
  region        = "%{region}"
}

resource "google_compute_subnetwork" "second_private_service_connect" {
  ip_cidr_range = "10.3.0.0/16"
  name          = "tf-test-compute-psc-subnetwork2%{random_suffix}"
  network       = google_compute_network.basic.self_link
  project       = "%{project_name}"
  purpose       = "PRIVATE_SERVICE_CONNECT"
  region        = "%{region}"
}

resource "google_compute_subnetwork" "private" {
  ip_cidr_range = "10.4.0.0/16"
  name          = "tf-test-compute-private-subnetwork%{random_suffix}"
  network       = google_compute_network.basic.self_link
  project       = "%{project_name}"
  purpose       = "PRIVATE"
  region        = "%{region}"
}

resource "google_compute_network" "basic" {
  name                    = "tf-test-compute-network%{random_suffix}"
  auto_create_subnetworks = false
  project                 = "%{project_name}"
}

resource "google_project" "first" {
  project_id = "tf-test-test-id1%{random_suffix}"
  name       = "tf-test-test-id1%{random_suffix}"
  org_id     = "%{org_id}"
}

resource "google_project" "second" {
  project_id = "tf-test-test-id2%{random_suffix}"
  name       = "tf-test-test-id2%{random_suffix}"
  org_id     = "%{org_id}"
}

`, context)
}

func testAccCheckComputeServiceAttachmentDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_compute_service_attachment" {
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

			obj := &compute.ServiceAttachment{
				ConnectionPreference: compute.ServiceAttachmentConnectionPreferenceEnumRef(rs.Primary.Attributes["connection_preference"]),
				Name:                 dcl.String(rs.Primary.Attributes["name"]),
				TargetService:        dcl.String(rs.Primary.Attributes["target_service"]),
				Description:          dcl.String(rs.Primary.Attributes["description"]),
				EnableProxyProtocol:  dcl.Bool(rs.Primary.Attributes["enable_proxy_protocol"] == "true"),
				Project:              dcl.StringOrNil(rs.Primary.Attributes["project"]),
				Location:             dcl.StringOrNil(rs.Primary.Attributes["region"]),
				Fingerprint:          dcl.StringOrNil(rs.Primary.Attributes["fingerprint"]),
				SelfLink:             dcl.StringOrNil(rs.Primary.Attributes["self_link"]),
			}

			client := NewDCLComputeClient(config, config.userAgent, billingProject, 0)
			_, err := client.GetServiceAttachment(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_compute_service_attachment still exists %v", obj)
			}
		}
		return nil
	}
}
