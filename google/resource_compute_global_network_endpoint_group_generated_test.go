// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccComputeGlobalNetworkEndpointGroup_globalNetworkEndpointGroupExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckComputeGlobalNetworkEndpointGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeGlobalNetworkEndpointGroup_globalNetworkEndpointGroupExample(context),
			},
			{
				ResourceName:      "google_compute_global_network_endpoint_group.neg",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeGlobalNetworkEndpointGroup_globalNetworkEndpointGroupExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_global_network_endpoint_group" "neg" {
  name                  = "tf-test-my-lb-neg%{random_suffix}"
  default_port          = "90"
  network_endpoint_type = "INTERNET_FQDN_PORT"
}
`, context)
}

func TestAccComputeGlobalNetworkEndpointGroup_globalNetworkEndpointGroupIpAddressExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {},
		},
		CheckDestroy: testAccCheckComputeGlobalNetworkEndpointGroupDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeGlobalNetworkEndpointGroup_globalNetworkEndpointGroupIpAddressExample(context),
			},
			{
				ResourceName:      "google_compute_global_network_endpoint_group.neg",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccComputeGlobalNetworkEndpointGroup_globalNetworkEndpointGroupIpAddressExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_compute_global_network_endpoint_group" "neg" {
  name                  = "tf-test-my-lb-neg%{random_suffix}"
  network_endpoint_type = "INTERNET_IP_PORT"
  default_port          = 90
}
`, context)
}

func testAccCheckComputeGlobalNetworkEndpointGroupDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_compute_global_network_endpoint_group" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := googleProviderConfig(t)

			url, err := replaceVarsForTest(config, rs, "{{ComputeBasePath}}projects/{{project}}/global/networkEndpointGroups/{{name}}")
			if err != nil {
				return err
			}

			_, err = sendRequest(config, "GET", "", url, nil)
			if err == nil {
				return fmt.Errorf("ComputeGlobalNetworkEndpointGroup still exists at %s", url)
			}
		}

		return nil
	}
}
