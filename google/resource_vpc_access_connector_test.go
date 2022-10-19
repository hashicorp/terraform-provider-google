package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVPCAccessConnector_vpcAccessConnectorThroughput(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPCAccessConnectorDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCAccessConnector_vpcAccessConnectorThroughput(context),
			},
			{
				ResourceName:      "google_vpc_access_connector.connector",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccVPCAccessConnector_vpcAccessConnectorThroughput(context map[string]interface{}) string {
	return Nprintf(`
resource "google_vpc_access_connector" "connector" {
  name          = "tf-test-vpc-con%{random_suffix}"
  subnet {
    name = google_compute_subnetwork.custom_test.name
  }
  machine_type = "e2-standard-4"
  min_instances = 2
  max_instances = 3
  region        = "us-central1"
}

resource "google_compute_subnetwork" "custom_test" {
  name          = "tf-test-vpc-con%{random_suffix}"
  ip_cidr_range = "10.2.0.0/28"
  region        = "us-central1"
  network       = google_compute_network.custom_test.id
}

resource "google_compute_network" "custom_test" {
  name                    = "tf-test-vpc-con%{random_suffix}"
  auto_create_subnetworks = false
}
`, context)
}
