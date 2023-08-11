// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networkservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkServicesGateway_update(t *testing.T) {
	t.Parallel()

	gatewayName := fmt.Sprintf("tf-test-gateway-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesGatewayDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesGateway_basic(gatewayName),
			},
			{
				ResourceName:      "google_network_services_gateway.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccNetworkServicesGateway_update(gatewayName),
			},
			{
				ResourceName:      "google_network_services_gateway.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetworkServicesGateway_basic(gatewayName string) string {
	return fmt.Sprintf(`
resource "google_network_services_gateway" "foobar" {
  name        = "%s"
  scope       = "default-scope-update"
  type        = "OPEN_MESH"
  ports       = [443]
  description = "my description"
}
`, gatewayName)
}

func testAccNetworkServicesGateway_update(gatewayName string) string {
	return fmt.Sprintf(`
resource "google_network_services_gateway" "foobar" {
  name        = "%s"
  scope       = "default-scope-update"
  type        = "OPEN_MESH"
  ports       = [443]
  description = "update description"
  labels      = {
    foo = "bar"
  } 
}
`, gatewayName)
}

// TODO(#14600): Enable the test once the api allows to update the fields for secure web gateway type.
//func TestAccNetworkServicesGateway_updateSwp(t *testing.T) {
//cmName := fmt.Sprintf("tf-test-gateway-swp-cm-%s", acctest.RandString(t, 10))
//	netName := fmt.Sprintf("tf-test-gateway-swp-net-%s", acctest.RandString(t, 10))
//	subnetName := fmt.Sprintf("tf-test-gateway-swp-subnet-%s", acctest.RandString(t, 10))
//	pSubnetName := fmt.Sprintf("tf-test-gateway-swp-proxyonly-%s", acctest.RandString(t, 10))
//	policyName := fmt.Sprintf("tf-test-gateway-swp-policy-%s", acctest.RandString(t, 10))
//	ruleName := fmt.Sprintf("tf-test-gateway-swp-rule-%s", acctest.RandString(t, 10))
//  gatewayScope := fmt.Sprintf("tf-test-gateway-swp-scope-%s", acctest.RandString(t, 10))
//	gatewayName := fmt.Sprintf("tf-test-gateway-swp-%s", acctest.RandString(t, 10))
//	// updates
//	newCmName := fmt.Sprintf("tf-test-gateway-swp-newcm-%s", acctest.RandString(t, 10))
//	newPolicyName := fmt.Sprintf("tf-test-gateway-swp-newpolicy-%s", acctest.RandString(t, 10))
//	newRuleName := fmt.Sprintf("tf-test-gateway-swp-newrule-%s", acctest.RandString(t, 10))
//
//	acctest.VcrTest(t, resource.TestCase{
//		PreCheck:     func() { acctest.AccTestPreCheck(t) },
//		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
//		CheckDestroy: testAccCheckNetworkServicesGatewayDestroyProducer(t),
//		Steps: []resource.TestStep{
//			{
//				Config: testAccNetworkServicesGateway_basicSwp(cmName, netName, subnetName, pSubnetName, policyName, ruleName, gatewayName, gatewayScope),
//			},
//			{
//				ResourceName:      "google_network_services_gateway.foobar",
//				ImportState:       true,
//				ImportStateVerify: true,
//        ImportStateVerifyIgnore: []string{"name", "location", "delete_swg_autogen_router_on_destroy"},
//			},
//			{
//				Config: testAccNetworkServicesGateway_updateSwp(cmName, newCmName, netName, subnetName, pSubnetName, policyName, newPolicyName, ruleName, newRuleName, gatewayName, gatewayScope),
//			},
//			{
//				ResourceName:      "google_network_services_gateway.foobar",
//				ImportState:       true,
//				ImportStateVerify: true,
//        ImportStateVerifyIgnore: []string{"name", "location", "delete_swg_autogen_router_on_destroy"},
//			},
//		},
//	})
//}

//func testAccNetworkServicesGateway_basicSwp(cmName, netName, subnetName, pSubnetName, policyName, ruleName, gatewayName, gatewayScope string) string {
//	return fmt.Sprintf(`
//resource "google_certificate_manager_certificate" "default" {
//  name        = "%s"
//  location    = "us-east1"
//  self_managed {
//    pem_certificate = file("test-fixtures/cert.pem")
//	pem_private_key = file("test-fixtures/private-key.pem")
//  }
//}
//
//resource "google_compute_network" "default" {
//  name                    = "%s"
//  routing_mode            = "REGIONAL"
//  auto_create_subnetworks = false
//}
//
//resource "google_compute_subnetwork" "proxyonlysubnet" {
//  name          = "%s"
//  purpose       = "REGIONAL_MANAGED_PROXY"
//  ip_cidr_range = "192.168.0.0/23"
//  region        = "us-east1"
//  network       = google_compute_network.default.id
//  role          = "ACTIVE"
//}
//
//resource "google_compute_subnetwork" "default" {
//  name          = "%s"
//  purpose       = "PRIVATE"
//  ip_cidr_range = "10.128.0.0/20"
//  region        = "us-east1"
//  network       = google_compute_network.default.id
//  role          = "ACTIVE"
//}
//
//resource "google_network_security_gateway_security_policy" "default" {
//  name        = "%s"
//  location    = "us-east1"
//}
//
//resource "google_network_security_gateway_security_policy_rule" "default" {
//  name                    = "%s"
//  location                = "us-east1"
//  gateway_security_policy = google_network_security_gateway_security_policy.default.name
//  enabled                 = true
//  priority                = 1
//  session_matcher         = "host() == 'example.com'"
//  basic_profile           = "ALLOW"
//}
//
//resource "google_network_services_gateway" "foobar" {
//  name                                 = "%s"
//  location                             = "us-east1"
//  addresses                            = ["10.128.0.99"]
//  type                                 = "SECURE_WEB_GATEWAY"
//  ports                                = [443]
//  description                          = "my description"
//  scope                                = "%s"
//  certificate_urls                     = [google_certificate_manager_certificate.default.id]
//  gateway_security_policy              = google_network_security_gateway_security_policy.default.id
//  network                              = google_compute_network.default.id
//  subnetwork                           = google_compute_subnetwork.default.id
//  delete_swg_autogen_router_on_destroy = true
//  depends_on                           = [google_compute_subnetwork.proxyonlysubnet]
//
//}
//`, cmName, netName, subnetName, pSubnetName, policyName, ruleName, gatewayName, gatewayScope)
//}

//func testAccNetworkServicesGateway_updateSwp(cmName, newCmName, netName, subnetName, pSubnetName, policyName, newPolicyName, ruleName, newRuleName, gatewayName, gatewayScope string) string {
//	return fmt.Sprintf(`
//resource "google_certificate_manager_certificate" "default" {
//  name        = "%s"
//  location    = "us-east1"
//  self_managed {
//    pem_certificate = file("test-fixtures/cert.pem")
//	pem_private_key = file("test-fixtures/private-key.pem")
//  }
//}
//
//resource "google_certificate_manager_certificate" "newcm" {
//  name        = "%s"
//  location    = "us-east1"
//  self_managed {
//    pem_certificate = file("test-fixtures/cert.pem")
//	pem_private_key = file("test-fixtures/private-key.pem")
//  }
//}
//
//resource "google_compute_network" "default" {
//  name                    = "%s"
//  routing_mode            = "REGIONAL"
//  auto_create_subnetworks = false
//}
//
//resource "google_compute_subnetwork" "proxyonlysubnet" {
//  name          = "%s"
//  purpose       = "REGIONAL_MANAGED_PROXY"
//  ip_cidr_range = "192.168.0.0/23"
//  region        = "us-east1"
//  network       = google_compute_network.default.id
//  role          = "ACTIVE"
//}
//
//resource "google_compute_subnetwork" "default" {
//  name          = "%s"
//  purpose       = "PRIVATE"
//  ip_cidr_range = "10.128.0.0/20"
//  region        = "us-east1"
//  network       = google_compute_network.default.id
//  role          = "ACTIVE"
//}
//
//resource "google_network_security_gateway_security_policy" "default" {
//  name        = "%s"
//  location    = "us-east1"
//}
//
//resource "google_network_security_gateway_security_policy_rule" "default" {
//  name                    = "%s"
//  location                = "us-east1"
//  gateway_security_policy = google_network_security_gateway_security_policy.default.name
//  enabled                 = true
//  priority                = 1
//  session_matcher         = "host() == 'example.com'"
//  basic_profile           = "ALLOW"
//}
//
//# TODO(#14600): this field will be updatable soon so this test should also cover it.
//# resource "google_network_security_gateway_security_policy" "newpolicy" {
//#   name        = "%s"
//#   location    = "us-east1"
//# }
//
//# resource "google_network_security_gateway_security_policy_rule" "newrule" {
//#   name                    = "%s"
//#   location                = "us-east1"
//#   gateway_security_policy = google_network_security_gateway_security_policy.newpolicy.name
//#   enabled                 = true
//#   priority                = 1
//#   session_matcher         = "host() == 'example.com'"
//#   basic_profile           = "ALLOW"
//# }
//
//resource "google_network_services_gateway" "foobar" {
//  name                                 = "%s"
//  location                             = "us-east1"
//  addresses                            = ["10.128.0.99"]
//  type                                 = "SECURE_WEB_GATEWAY"
//  ports                                = [443]
//  description                          = "updated description"
//  scope                                = "%s"
//  certificate_urls                     = [google_certificate_manager_certificate.default.id, google_certificate_manager_certificate.newcm.id]
//  gateway_security_policy              = google_network_security_gateway_security_policy.default.id
//  # TODO(#14600): this field will be updatable soon so this test should also cover it.
//  # gateway_security_policy              = google_network_security_gateway_security_policy.newpolicy.id
//  network                              = google_compute_network.default.id
//  subnetwork                           = google_compute_subnetwork.default.id
//  delete_swg_autogen_router_on_destroy = true
//  depends_on                           = [google_compute_subnetwork.proxyonlysubnet]
//
//}
//`, cmName, newCmName, netName, subnetName, pSubnetName, policyName, newPolicyName, ruleName, newRuleName, gatewayName, gatewayScope)
//}

func TestAccNetworkServicesGateway_multipleSwpGatewaysDifferentSubnetwork(t *testing.T) {
	cmName := fmt.Sprintf("tf-test-gateway-multiswp-cm-%s", acctest.RandString(t, 10))
	netName := fmt.Sprintf("tf-test-gateway-multiswp-net-%s", acctest.RandString(t, 10))
	subnetName := fmt.Sprintf("tf-test-gateway-multiswp-subnet-%s", acctest.RandString(t, 10))
	pSubnetName := fmt.Sprintf("tf-test-gateway-multiswp-proxyonly-%s", acctest.RandString(t, 10))
	policyName := fmt.Sprintf("tf-test-gateway-multiswp-policy-%s", acctest.RandString(t, 10))
	ruleName := fmt.Sprintf("tf-test-gateway-multiswp-rule-%s", acctest.RandString(t, 10))
	gatewayScope := fmt.Sprintf("tf-test-gateway-multiswp-scope-%s", acctest.RandString(t, 10))
	gatewayName := fmt.Sprintf("tf-test-gateway-multiswp-%s", acctest.RandString(t, 10))
	subnet2Name := fmt.Sprintf("tf-test-gateway-multiswp-subnet2-%s", acctest.RandString(t, 10))
	gateway2Name := fmt.Sprintf("tf-test-gateway-multiswp2-%s", acctest.RandString(t, 10))
	gateway2Scope := fmt.Sprintf("tf-test-gateway-multiswp-scope2-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesGatewayDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesGateway_multipleSwpGatewaysDifferentSubnetwork(cmName, netName, subnetName, pSubnetName, policyName, ruleName, gatewayName, gatewayScope, subnet2Name, gateway2Name, gateway2Scope),
			},
			{
				ResourceName:            "google_network_services_gateway.gateway1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "delete_swg_autogen_router_on_destroy"},
			},
			{
				Config: testAccNetworkServicesGateway_multipleSwpGatewaysDifferentSubnetworkRemoveGateway2(cmName, netName, subnetName, pSubnetName, policyName, ruleName, gatewayName, gatewayScope, subnet2Name),
			},
			{
				ResourceName:            "google_network_services_gateway.gateway1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "delete_swg_autogen_router_on_destroy"},
			},
		},
	})
}

func testAccNetworkServicesGateway_multipleSwpGatewaysDifferentSubnetwork(cmName, netName, subnetName, pSubnetName, policyName, ruleName, gatewayName, gatewayScope, subnet2Name, gateway2Name, gateway2Scope string) string {
	return fmt.Sprintf(`
resource "google_certificate_manager_certificate" "default" {
  name        = "%s"
  location    = "us-west1"
  self_managed {
    pem_certificate = file("test-fixtures/cert.pem")
	pem_private_key = file("test-fixtures/private-key.pem")
  }
}

resource "google_compute_network" "default" {
  name                    = "%s"
  routing_mode            = "REGIONAL"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "proxyonlysubnet" {
  name          = "%s"
  purpose       = "REGIONAL_MANAGED_PROXY"
  ip_cidr_range = "192.168.0.0/23"
  region        = "us-west1"
  network       = google_compute_network.default.id
  role          = "ACTIVE"
}

resource "google_compute_subnetwork" "subnet1" {
  name          = "%s"
  purpose       = "PRIVATE"
  ip_cidr_range = "10.128.0.0/20"
  region        = "us-west1"
  network       = google_compute_network.default.id
  role          = "ACTIVE"
}

resource "google_network_security_gateway_security_policy" "default" {
  name        = "%s"
  location    = "us-west1"
}

resource "google_network_security_gateway_security_policy_rule" "default" {
  name                    = "%s"
  location                = "us-west1"
  gateway_security_policy = google_network_security_gateway_security_policy.default.name
  enabled                 = true  
  priority                = 1
  session_matcher         = "host() == 'example.com'"
  basic_profile           = "ALLOW"
}
	  
resource "google_network_services_gateway" "gateway1" {
  name                                 = "%s"
  location                             = "us-west1"
  addresses                            = ["10.128.0.99"]
  type                                 = "SECURE_WEB_GATEWAY"
  ports                                = [443]
  description                          = "gateway1_subnet1"
  scope                                = "%s"
  certificate_urls                     = [google_certificate_manager_certificate.default.id]
  gateway_security_policy              = google_network_security_gateway_security_policy.default.id
  network                              = google_compute_network.default.id
  subnetwork                           = google_compute_subnetwork.subnet1.id
  delete_swg_autogen_router_on_destroy = true
  depends_on                           = [google_compute_subnetwork.proxyonlysubnet]
}

resource "google_compute_subnetwork" "subnet2" {
  name          = "%s"
  purpose       = "PRIVATE"
  ip_cidr_range = "10.142.0.0/20"
  region        = "us-west1"
  network       = google_compute_network.default.id
  role          = "ACTIVE"
}

resource "google_network_services_gateway" "gateway2" {
  name                                 = "%s"
  location                             = "us-west1"
  addresses                            = ["10.142.0.99"]
  type                                 = "SECURE_WEB_GATEWAY"
  ports                                = [443]
  description                          = "gateway2_subnet2"
  scope                                = "%s"
  certificate_urls                     = [google_certificate_manager_certificate.default.id]
  gateway_security_policy              = google_network_security_gateway_security_policy.default.id
  network                              = google_compute_network.default.id
  subnetwork                           = google_compute_subnetwork.subnet2.id
  delete_swg_autogen_router_on_destroy = true
  depends_on                           = [google_compute_subnetwork.proxyonlysubnet]
}

`, cmName, netName, subnetName, pSubnetName, policyName, ruleName, gatewayName, gatewayScope, subnet2Name, gateway2Name, gateway2Scope)
}

func testAccNetworkServicesGateway_multipleSwpGatewaysDifferentSubnetworkRemoveGateway2(cmName, netName, subnetName, pSubnetName, policyName, ruleName, gatewayName, gatewayScope, subnet2Name string) string {
	return fmt.Sprintf(`
resource "google_certificate_manager_certificate" "default" {
  name        = "%s"
  location    = "us-west1"
  self_managed {
    pem_certificate = file("test-fixtures/cert.pem")
	pem_private_key = file("test-fixtures/private-key.pem")
  }
}

resource "google_compute_network" "default" {
  name                    = "%s"
  routing_mode            = "REGIONAL"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "proxyonlysubnet" {
  name          = "%s"
  purpose       = "REGIONAL_MANAGED_PROXY"
  ip_cidr_range = "192.168.0.0/23"
  region        = "us-west1"
  network       = google_compute_network.default.id
  role          = "ACTIVE"
}

resource "google_compute_subnetwork" "subnet1" {
  name          = "%s"
  purpose       = "PRIVATE"
  ip_cidr_range = "10.128.0.0/20"
  region        = "us-west1"
  network       = google_compute_network.default.id
  role          = "ACTIVE"
}

resource "google_network_security_gateway_security_policy" "default" {
  name        = "%s"
  location    = "us-west1"
}

resource "google_network_security_gateway_security_policy_rule" "default" {
  name                    = "%s"
  location                = "us-west1"
  gateway_security_policy = google_network_security_gateway_security_policy.default.name
  enabled                 = true  
  priority                = 1
  session_matcher         = "host() == 'example.com'"
  basic_profile           = "ALLOW"
}
	  
resource "google_network_services_gateway" "gateway1" {
  name                                 = "%s"
  location                             = "us-west1"
  addresses                            = ["10.128.0.99"]
  type                                 = "SECURE_WEB_GATEWAY"
  ports                                = [443]
  description                          = "gateway1_subnet1"
  scope                                = "%s"
  certificate_urls                     = [google_certificate_manager_certificate.default.id]
  gateway_security_policy              = google_network_security_gateway_security_policy.default.id
  network                              = google_compute_network.default.id
  subnetwork                           = google_compute_subnetwork.subnet1.id
  delete_swg_autogen_router_on_destroy = true
  depends_on                           = [google_compute_subnetwork.proxyonlysubnet]
}

resource "google_compute_subnetwork" "subnet2" {
  name          = "%s"
  purpose       = "PRIVATE"
  ip_cidr_range = "10.142.0.0/20"
  region        = "us-west1"
  network       = google_compute_network.default.id
  role          = "ACTIVE"
}

# Destroying gateway2 so it allows to test if there is still a gateway remaining under the same network so the swg_autogen_router is kept.

`, cmName, netName, subnetName, pSubnetName, policyName, ruleName, gatewayName, gatewayScope, subnet2Name)
}

func TestAccNetworkServicesGateway_multipleSwpGatewaysDifferentNetwork(t *testing.T) {
	cmName := fmt.Sprintf("tf-test-gateway-diffswp-cm-%s", acctest.RandString(t, 10))
	netName := fmt.Sprintf("tf-test-gateway-diffswp-net-%s", acctest.RandString(t, 10))
	subnetName := fmt.Sprintf("tf-test-gateway-diffswp-subnet-%s", acctest.RandString(t, 10))
	pSubnetName := fmt.Sprintf("tf-test-gateway-diffswp-proxyonly-%s", acctest.RandString(t, 10))
	policyName := fmt.Sprintf("tf-test-gateway-diffswp-policy-%s", acctest.RandString(t, 10))
	ruleName := fmt.Sprintf("tf-test-gateway-diffswp-rule-%s", acctest.RandString(t, 10))
	gatewayName := fmt.Sprintf("tf-test-gateway-diffswp-%s", acctest.RandString(t, 10))
	gatewayScope := fmt.Sprintf("tf-test-gateway-diffswp-scope-%s", acctest.RandString(t, 10))
	net2Name := fmt.Sprintf("tf-test-gateway-diffswp-net2-%s", acctest.RandString(t, 10))
	subnet2Name := fmt.Sprintf("tf-test-gateway-diffswp-subnet2-%s", acctest.RandString(t, 10))
	pSubnet2Name := fmt.Sprintf("tf-test-gateway-diffswp-proxyonly2-%s", acctest.RandString(t, 10))
	gateway2Name := fmt.Sprintf("tf-test-gateway-diffswp2-%s", acctest.RandString(t, 10))
	gateway2Scope := fmt.Sprintf("tf-test-gateway-diffswp-scope2-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesGatewayDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesGateway_multipleSwpGatewaysDifferentNetwork(cmName, netName, subnetName, pSubnetName, policyName, ruleName, gatewayName, gatewayScope, net2Name, subnet2Name, pSubnet2Name, gateway2Name, gateway2Scope),
			},
			{
				ResourceName:            "google_network_services_gateway.gateway1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "delete_swg_autogen_router_on_destroy"},
			},
			{
				Config: testAccNetworkServicesGateway_multipleSwpGatewaysDifferentNetworkRemoveGateway2(cmName, netName, subnetName, pSubnetName, policyName, ruleName, gatewayName, gatewayScope, net2Name, subnet2Name, pSubnet2Name),
			},
			{
				ResourceName:            "google_network_services_gateway.gateway1",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "delete_swg_autogen_router_on_destroy"},
			},
		},
	})
}

func testAccNetworkServicesGateway_multipleSwpGatewaysDifferentNetwork(cmName, netName, subnetName, pSubnetName, policyName, ruleName, gatewayName, gatewayScope, net2Name, subnet2Name, pSubnet2Name, gateway2Name, gateway2Scope string) string {
	return fmt.Sprintf(`
resource "google_certificate_manager_certificate" "default" {
  name        = "%s"
  location    = "us-west2"
  self_managed {
    pem_certificate = file("test-fixtures/cert.pem")
	  pem_private_key = file("test-fixtures/private-key.pem")
  }
}

resource "google_compute_network" "default" {
  name                    = "%s"
  routing_mode            = "REGIONAL"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "proxyonlysubnet" {
  name          = "%s"
  purpose       = "REGIONAL_MANAGED_PROXY"
  ip_cidr_range = "192.168.0.0/23"
  region        = "us-west2"
  network       = google_compute_network.default.id
  role          = "ACTIVE"
}

resource "google_compute_subnetwork" "subnet1" {
  name          = "%s"
  purpose       = "PRIVATE"
  ip_cidr_range = "10.128.0.0/20"
  region        = "us-west2"
  network       = google_compute_network.default.id
  role          = "ACTIVE"
}

resource "google_network_security_gateway_security_policy" "default" {
  name        = "%s"
  location    = "us-west2"
}

resource "google_network_security_gateway_security_policy_rule" "default" {
  name                    = "%s"
  location                = "us-west2"
  gateway_security_policy = google_network_security_gateway_security_policy.default.name
  enabled                 = true  
  priority                = 1
  session_matcher         = "host() == 'example.com'"
  basic_profile           = "ALLOW"
}
	  
resource "google_network_services_gateway" "gateway1" {
  name                                 = "%s"
  location                             = "us-west2"
  addresses                            = ["10.128.0.99"]
  type                                 = "SECURE_WEB_GATEWAY"
  ports                                = [443]
  description                          = "gateway1_subnet1"
  scope                                = "%s"
  certificate_urls                     = [google_certificate_manager_certificate.default.id]
  gateway_security_policy              = google_network_security_gateway_security_policy.default.id
  network                              = google_compute_network.default.id
  subnetwork                           = google_compute_subnetwork.subnet1.id
  delete_swg_autogen_router_on_destroy = true
  depends_on                           = [google_compute_subnetwork.proxyonlysubnet]
}

resource "google_compute_network" "network2" {
  name                    = "%s"
  routing_mode            = "REGIONAL"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnet2" {
  name          = "%s"
  purpose       = "PRIVATE"
  ip_cidr_range = "10.142.0.0/20"
  region        = "us-west2"
  network       = google_compute_network.network2.id
  role          = "ACTIVE"
}

resource "google_compute_subnetwork" "proxyonlysubnet2" {
  region        = "us-west2"
  name          = "%s"
  purpose       = "REGIONAL_MANAGED_PROXY"
  ip_cidr_range = "192.170.0.0/23"
  network       = google_compute_network.network2.id
  role          = "ACTIVE"
}

resource "google_network_services_gateway" "gateway2" {
  name                                 = "%s"
  location                             = "us-west2"
  addresses                            = ["10.142.0.99"]
  type                                 = "SECURE_WEB_GATEWAY"
  ports                                = [443]
  description                          = "gateway2_subnet2"
  scope                                = "%s"
  certificate_urls                     = [google_certificate_manager_certificate.default.id]
  gateway_security_policy              = google_network_security_gateway_security_policy.default.id
  network                              = google_compute_network.network2.id
  subnetwork                           = google_compute_subnetwork.subnet2.id
  delete_swg_autogen_router_on_destroy = true
  depends_on                           = [google_compute_subnetwork.proxyonlysubnet2]
}

`, cmName, netName, subnetName, pSubnetName, policyName, ruleName, gatewayName, gatewayScope, net2Name, subnet2Name, pSubnet2Name, gateway2Name, gateway2Scope)
}

func testAccNetworkServicesGateway_multipleSwpGatewaysDifferentNetworkRemoveGateway2(cmName, netName, subnetName, pSubnetName, policyName, ruleName, gatewayName, gatewayScope, net2Name, subnet2Name, pSubnet2Name string) string {
	return fmt.Sprintf(`
resource "google_certificate_manager_certificate" "default" {
  name        = "%s"
  location    = "us-west2"
  self_managed {
    pem_certificate = file("test-fixtures/cert.pem")
	  pem_private_key = file("test-fixtures/private-key.pem")
  }
}

resource "google_compute_network" "default" {
  name                    = "%s"
  routing_mode            = "REGIONAL"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "proxyonlysubnet" {
  name          = "%s"
  purpose       = "REGIONAL_MANAGED_PROXY"
  ip_cidr_range = "192.168.0.0/23"
  region        = "us-west2"
  network       = google_compute_network.default.id
  role          = "ACTIVE"
}

resource "google_compute_subnetwork" "subnet1" {
  name          = "%s"
  purpose       = "PRIVATE"
  ip_cidr_range = "10.128.0.0/20"
  region        = "us-west2"
  network       = google_compute_network.default.id
  role          = "ACTIVE"
}

resource "google_network_security_gateway_security_policy" "default" {
  name        = "%s"
  location    = "us-west2"
}

resource "google_network_security_gateway_security_policy_rule" "default" {
  name                    = "%s"
  location                = "us-west2"
  gateway_security_policy = google_network_security_gateway_security_policy.default.name
  enabled                 = true  
  priority                = 1
  session_matcher         = "host() == 'example.com'"
  basic_profile           = "ALLOW"
}
	  
resource "google_network_services_gateway" "gateway1" {
  name                                 = "%s"
  location                             = "us-west2"
  addresses                            = ["10.128.0.99"]
  type                                 = "SECURE_WEB_GATEWAY"
  ports                                = [443]
  description                          = "gateway1_subnet1"
  scope                                = "%s"
  certificate_urls                     = [google_certificate_manager_certificate.default.id]
  gateway_security_policy              = google_network_security_gateway_security_policy.default.id
  network                              = google_compute_network.default.id
  subnetwork                           = google_compute_subnetwork.subnet1.id
  delete_swg_autogen_router_on_destroy = true
  depends_on                           = [google_compute_subnetwork.proxyonlysubnet]
}

resource "google_compute_network" "network2" {
  name                    = "%s"
  routing_mode            = "REGIONAL"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnet2" {
  name          = "%s"
  purpose       = "PRIVATE"
  ip_cidr_range = "10.142.0.0/20"
  region        = "us-west2"
  network       = google_compute_network.network2.id
  role          = "ACTIVE"
}

resource "google_compute_subnetwork" "proxyonlysubnet2" {
  region        = "us-west2"
  name          = "%s"
  purpose       = "REGIONAL_MANAGED_PROXY"
  ip_cidr_range = "192.170.0.0/23"
  network       = google_compute_network.network2.id
  role          = "ACTIVE"
}

# Destroying gateway2 so it allows to test that there is no gateway remaining under the same network so the swg_autogen_router is deleted.

`, cmName, netName, subnetName, pSubnetName, policyName, ruleName, gatewayName, gatewayScope, net2Name, subnet2Name, pSubnet2Name)
}

func TestAccNetworkServicesGateway_minimalSwp(t *testing.T) {
	netName := fmt.Sprintf("tf-test-gateway-swp-net-%s", acctest.RandString(t, 10))
	subnetName := fmt.Sprintf("tf-test-gateway-swp-subnet-%s", acctest.RandString(t, 10))
	pSubnetName := fmt.Sprintf("tf-test-gateway-swp-proxyonly-%s", acctest.RandString(t, 10))
	policyName := fmt.Sprintf("tf-test-gateway-swp-policy-%s", acctest.RandString(t, 10))
	ruleName := fmt.Sprintf("tf-test-gateway-swp-rule-%s", acctest.RandString(t, 10))
	gatewayName := fmt.Sprintf("tf-test-gateway-swp-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesGatewayDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesGateway_minimalSwp(netName, subnetName, pSubnetName, policyName, ruleName, gatewayName),
			},
			{
				ResourceName:            "google_network_services_gateway.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "delete_swg_autogen_router_on_destroy"},
			},
		},
	})
}

func testAccNetworkServicesGateway_minimalSwp(netName, subnetName, pSubnetName, policyName, ruleName, gatewayName string) string {
	return fmt.Sprintf(`
resource "google_compute_network" "default" {
  name                    = "%s"
  routing_mode            = "REGIONAL"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "proxyonlysubnet" {
  name          = "%s"
  purpose       = "REGIONAL_MANAGED_PROXY"
  ip_cidr_range = "192.168.0.0/23"
  region        = "us-central1"
  network       = google_compute_network.default.id
  role          = "ACTIVE"
}

resource "google_compute_subnetwork" "default" {
  name          = "%s"
  purpose       = "PRIVATE"
  ip_cidr_range = "10.128.0.0/20"
  region        = "us-central1"
  network       = google_compute_network.default.id
  role          = "ACTIVE"
}

resource "google_network_security_gateway_security_policy" "default" {
  name        = "%s"
  location    = "us-central1"
}

resource "google_network_security_gateway_security_policy_rule" "default" {
  name                    = "%s"
  location                = "us-central1"
  gateway_security_policy = google_network_security_gateway_security_policy.default.name
  enabled                 = true
  priority                = 1
  session_matcher         = "host() == 'example.com'"
  basic_profile           = "ALLOW"
}

resource "google_network_services_gateway" "foobar" {
  name                                 = "%s"
  location                             = "us-central1"
  addresses                            = ["10.128.0.99"]
  type                                 = "SECURE_WEB_GATEWAY"
  ports                                = [443]
  description                          = "my description"
  gateway_security_policy              = google_network_security_gateway_security_policy.default.id
  network                              = google_compute_network.default.id
  subnetwork                           = google_compute_subnetwork.default.id
  delete_swg_autogen_router_on_destroy = true
  depends_on                           = [google_compute_subnetwork.proxyonlysubnet]
}
`, netName, subnetName, pSubnetName, policyName, ruleName, gatewayName)
}
