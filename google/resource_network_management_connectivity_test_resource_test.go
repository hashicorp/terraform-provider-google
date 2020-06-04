package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNetworkManagementConnectivityTest_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkManagementConnectivityTestDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkManagementConnectivityTest_instanceToInstance(context),
			},
			{
				ResourceName:      "google_network_management_connectivity_test.conn-test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkManagementConnectivityTest_instanceToAddr(context),
			},
			{
				ResourceName:      "google_network_management_connectivity_test.conn-test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkManagementConnectivityTest_instanceToInstance(context map[string]interface{}) string {
	connTestCfg := Nprintf(`
resource "google_network_management_connectivity_test" "conn-test" {
  name = "tf-test-conntest%{random_suffix}"
  source {
    instance = google_compute_instance.vm1.id
  }

  destination {
    instance = google_compute_instance.vm2.id
  }

  protocol = "TCP"
}
`, context)
	return fmt.Sprintf("%s\n\n%s\n\n", connTestCfg, testAccNetworkManagementConnectivityTest_baseResources(context))
}

func testAccNetworkManagementConnectivityTest_instanceToAddr(context map[string]interface{}) string {
	connTestCfg := Nprintf(`
resource "google_network_management_connectivity_test" "conn-test" {
  name = "tf-test-conntest%{random_suffix}"
  source {
	instance = google_compute_instance.vm1.id
	network = google_compute_network.vpc.id
	port = 50
  }

  destination {
	ip_address = google_compute_address.addr.address
	project_id =  google_compute_address.addr.address
	network = google_compute_network.vpc.id
	port = 80
  }

  protocol = "TCP"
}
`, context)
	return fmt.Sprintf("%s\n\n%s\n\n", connTestCfg, testAccNetworkManagementConnectivityTest_baseResources(context))
}

func testAccNetworkManagementConnectivityTest_baseResources(context map[string]interface{}) string {
	return Nprintf(`

resource "google_compute_address" "addr" {
	name         = "tf-test-addr%{random_suffix}"
	subnetwork   = google_compute_subnetwork.subnet.id
	address_type = "INTERNAL"
	address      = "10.0.43.43"
	region       = "us-central1"
}

resource "google_compute_instance" "vm1" {
  	name = "tf-test-src-vm%{random_suffix}"
	machine_type = "n1-standard-1"
	boot_disk {
	  initialize_params {
	    image = data.google_compute_image.debian_9.self_link
	  }
	}	
	network_interface {
	  network = google_compute_network.vpc.id
	}
}

resource "google_compute_instance" "vm2" {
	name = "tf-test-vm-dest%{random_suffix}"
	machine_type = "n1-standard-1"
  
	boot_disk {
	  initialize_params {
		image = data.google_compute_image.debian_9.self_link
	  }
	}
  
	network_interface {
	  network = google_compute_network.vpc.id

	}
}

resource "google_compute_network" "vpc" {
	name = "tf-test-connnet%{random_suffix}"
}

resource "google_compute_subnetwork" "subnet" {
	name          = "tf-test-connet%{random_suffix}"
	ip_cidr_range = "10.0.0.0/16"
	region        = "us-central1"
	network       = google_compute_network.vpc.id
}	

data "google_compute_image" "debian_9" {
	family  = "debian-9"
	project = "debian-cloud"
}
`, context)
}
