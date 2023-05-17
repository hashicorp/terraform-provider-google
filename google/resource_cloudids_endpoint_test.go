package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccCloudIdsEndpoint_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudIdsEndpointDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testCloudIds_basic(context),
			},
			{
				ResourceName:      "google_cloud_ids_endpoint.endpoint",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testCloudIds_basicUpdate(context),
			},
			{
				ResourceName:      "google_cloud_ids_endpoint.endpoint",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCloudIds_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_network" "default" {
	name = "tf-test-my-network%{random_suffix}"
}
resource "google_compute_global_address" "service_range" {
	name          = "address"
	purpose       = "VPC_PEERING"
	address_type  = "INTERNAL"
	prefix_length = 16
	network       = google_compute_network.default.id
}
resource "google_service_networking_connection" "private_service_connection" {
	network                 = google_compute_network.default.id
	service                 = "servicenetworking.googleapis.com"
	reserved_peering_ranges = [google_compute_global_address.service_range.name]
}
  
resource "google_cloud_ids_endpoint" "endpoint" {
	name              = "cloud-ids-test-%{random_suffix}"
	location          = "us-central1-f"
	network           = google_compute_network.default.id
	severity          = "INFORMATIONAL"
	threat_exceptions = ["12", "67"]
	depends_on        = [google_service_networking_connection.private_service_connection]
}
`, context)
}

func testCloudIds_basicUpdate(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_network" "default" {
	name = "tf-test-my-network%{random_suffix}"
}
resource "google_compute_global_address" "service_range" {
	name          = "address"
	purpose       = "VPC_PEERING"
	address_type  = "INTERNAL"
	prefix_length = 16
	network       = google_compute_network.default.id
}
resource "google_service_networking_connection" "private_service_connection" {
	network                 = google_compute_network.default.id
	service                 = "servicenetworking.googleapis.com"
	reserved_peering_ranges = [google_compute_global_address.service_range.name]
}
  
resource "google_cloud_ids_endpoint" "endpoint" {
	name              = "cloud-ids-test-%{random_suffix}"
	location          = "us-central1-f"
	network           = google_compute_network.default.id
	severity          = "INFORMATIONAL"
	threat_exceptions = ["33"]
	depends_on        = [google_service_networking_connection.private_service_connection]
}
`, context)
}

func testAccCheckCloudIdsEndpointDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_cloud_ids_endpoint" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{CloudIdsBasePath}}projects/{{project}}/locations/{{location}}/endpoints/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("CloudIdsEndpoint still exists at %s", url)
			}
		}

		return nil
	}
}
