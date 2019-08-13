package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccComputeNetworkEndpoint_networkEndpointsBasic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(10),
		"default_port":  90,
		"modified_port": 100,
		"add1_port":     101,
		"add2_port":     102,
	}
	negId := fmt.Sprintf("projects/%s/zones/%s/networkEndpointGroups/neg-%s",
		getTestProjectFromEnv(), getTestZoneFromEnv(), context["random_suffix"])

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Create one endpoint
				Config: testAccComputeNetworkEndpoint_networkEndpointsBasic(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkEndpointWithPortExists("google_compute_network_endpoint.default", "90"),
				),
			},
			{
				ResourceName:      "google_compute_network_endpoint.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Force-recreate old endpoint
				Config: testAccComputeNetworkEndpoint_networkEndpointsModified(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkEndpointWithPortExists("google_compute_network_endpoint.default", "100"),
					testAccCheckComputeNetworkEndpointWithPortsDestroyed(negId, "90"),
				),
			},
			{
				// Add two new endpoints
				Config: testAccComputeNetworkEndpoint_networkEndpointsAdditional(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkEndpointWithPortExists("google_compute_network_endpoint.default", "100"),
					testAccCheckComputeNetworkEndpointWithPortExists("google_compute_network_endpoint.add1", "101"),
					testAccCheckComputeNetworkEndpointWithPortExists("google_compute_network_endpoint.add2", "102"),
				),
			},
			{
				// delete all endpoints
				Config: testAccComputeNetworkEndpoint_noNetworkEndpoints(context),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeNetworkEndpointWithPortsDestroyed(negId, "100"),
				),
			},
		},
	})
}

func testAccComputeNetworkEndpoint_networkEndpointsBasic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_network_endpoint" "default" {
 zone                   = "us-central1-a"
 network_endpoint_group = "${google_compute_network_endpoint_group.neg.name}"

 instance    = "${google_compute_instance.default.name}"
 ip_address  = "${google_compute_instance.default.network_interface.0.network_ip}"
 port        = "${google_compute_network_endpoint_group.neg.default_port}"
}
`, context) + testAccComputeNetworkEndpoint_noNetworkEndpoints(context)
}

func testAccComputeNetworkEndpoint_networkEndpointsModified(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_network_endpoint" "default" {
 zone                   = "us-central1-a"
 network_endpoint_group = "${google_compute_network_endpoint_group.neg.name}"

 instance    = "${google_compute_instance.default.name}"
 ip_address  = "${google_compute_instance.default.network_interface.0.network_ip}"
 port        = "%{modified_port}"
}`, context) + testAccComputeNetworkEndpoint_noNetworkEndpoints(context)
}

func testAccComputeNetworkEndpoint_networkEndpointsAdditional(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_network_endpoint" "default" {
 zone                   = "us-central1-a"
 network_endpoint_group = "${google_compute_network_endpoint_group.neg.name}"

 instance    = "${google_compute_instance.default.name}"
 ip_address  = "${google_compute_instance.default.network_interface.0.network_ip}"
 port        = "%{modified_port}"
}

resource "google_compute_network_endpoint" "add1" {
 zone                   = "us-central1-a"
 network_endpoint_group = "${google_compute_network_endpoint_group.neg.name}"

 instance    = "${google_compute_instance.default.name}"
 ip_address  = "${google_compute_instance.default.network_interface.0.network_ip}"
 port        = "%{add1_port}"
}

resource "google_compute_network_endpoint" "add2" {
 zone                   = "us-central1-a"
 network_endpoint_group = "${google_compute_network_endpoint_group.neg.name}"

 instance    = "${google_compute_instance.default.name}"
 ip_address  = "${google_compute_instance.default.network_interface.0.network_ip}"
 port        = "%{add2_port}"
}

`, context) + testAccComputeNetworkEndpoint_noNetworkEndpoints(context)
}

func testAccComputeNetworkEndpoint_noNetworkEndpoints(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_network_endpoint_group" "neg" {
 name         = "neg-%{random_suffix}"
 zone         = "us-central1-a"
 network      = "${google_compute_network.default.self_link}"
 subnetwork   = "${google_compute_subnetwork.default.self_link}"
 default_port = "%{default_port}"
}

resource "google_compute_network" "default" {
 name = "neg-network-%{random_suffix}"
 auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
 name          = "neg-subnetwork-%{random_suffix}"
 ip_cidr_range = "10.0.0.0/16"
 region        = "us-central1"
 network       = "${google_compute_network.default.self_link}"
}

resource "google_compute_instance" "default" {
 name         =  "neg-instance1-%{random_suffix}"
 machine_type = "n1-standard-1"

 boot_disk {
   initialize_params{
     image = "${data.google_compute_image.my_image.self_link}"
   }
 }

 network_interface {
   subnetwork = "${google_compute_subnetwork.default.self_link}"
   access_config { }
 }
}

data "google_compute_image" "my_image" {
 family  = "debian-9"
 project = "debian-cloud"
}
`, context)
}

// testAccCheckComputeNetworkEndpointExists makes sure the resource with given
// (Terraform) name exists, and returns identifying information about the
// existing endpoint
func testAccCheckComputeNetworkEndpointWithPortExists(name, port string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource %q not in path %q", name, s.RootModule().Path)
		}

		if rs.Type != "google_compute_network_endpoint" {
			return fmt.Errorf("resource %q has unexpected type %q", name, rs.Type)
		}

		if rs.Primary.Attributes["port"] != port {
			return fmt.Errorf("unexpected port %s for resource %s, expected %s", rs.Primary.Attributes["port"], name, port)
		}

		config := testAccProvider.Meta().(*Config)

		negResourceId, err := replaceVarsForTest(config, rs, "projects/{{project}}/zones/{{zone}}/networkEndpointGroups/{{network_endpoint_group}}")
		if err != nil {
			return fmt.Errorf("creating URL for getting network endpoint %q failed: %v", name, err)
		}

		foundPorts, err := testAccComputeNetworkEndpointsListEndpointPorts(negResourceId)
		if err != nil {
			return fmt.Errorf("unable to confirm endpoints with port %s exists: %v", port, err)
		}
		if _, ok := foundPorts[port]; !ok {
			return fmt.Errorf("did not find endpoint with port %s", port)
		}
		return nil
	}
}

// testAccCheckComputeNetworkEndpointDestroyed makes sure the endpoint with
// given Terraform resource name and previous information (obtained from Exists)
// was destroyed properly.
func testAccCheckComputeNetworkEndpointWithPortsDestroyed(negId string, ports ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		foundPorts, err := testAccComputeNetworkEndpointsListEndpointPorts(negId)
		if err != nil {
			return fmt.Errorf("unable to confirm endpoints with ports %+v was destroyed: %v", ports, err)
		}
		for _, p := range ports {
			if _, ok := foundPorts[p]; ok {
				return fmt.Errorf("network endpoint with port %s still exists", p)
			}
		}

		return nil
	}
}

func testAccComputeNetworkEndpointsListEndpointPorts(negId string) (map[string]struct{}, error) {
	config := testAccProvider.Meta().(*Config)

	url := fmt.Sprintf("https://www.googleapis.com/compute/beta/%s/listNetworkEndpoints", negId)
	res, err := sendRequest(config, "POST", "", url, nil)
	if err != nil {
		return nil, err
	}

	v, ok := res["items"]
	if !ok || v == nil {
		return nil, nil
	}
	items := v.([]interface{})
	ports := make(map[string]struct{})
	for _, item := range items {
		endptWithHealth := item.(map[string]interface{})
		v, ok := endptWithHealth["networkEndpoint"]
		if !ok || v == nil {
			continue
		}
		endpt := v.(map[string]interface{})
		ports[fmt.Sprintf("%v", endpt["port"])] = struct{}{}
	}
	return ports, nil
}
