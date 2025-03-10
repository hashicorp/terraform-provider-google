// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networkservices_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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
				ResourceName:            "google_network_services_gateway.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkServicesGateway_update(gatewayName),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("google_network_services_gateway.foobar", plancheck.ResourceActionUpdate),
					},
				},
			},
			{
				ResourceName:            "google_network_services_gateway.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
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
  ports       = [1000]
  description = "update description"
  labels      = {
    foo = "bar"
  }
}
`, gatewayName)
}

func TestAccNetworkServicesGateway_networkServicesGatewaySecureWebProxyWithoutAddresses(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesGatewayDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesGateway_networkServicesGatewaySecureWebProxy(context, false),
			},
			{
				ResourceName:            "google_network_services_gateway.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "delete_swg_autogen_router_on_destroy", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccNetworkServicesGateway_networkServicesGatewaySecureWebProxy(context map[string]interface{}, withAddresses bool) string {
	config := ""
	config += acctest.Nprintf(`
resource "google_certificate_manager_certificate" "default" {
  name        = "tf-test-my-certificate-%{random_suffix}"
  location    = "us-central1"
  self_managed {
    pem_certificate = file("test-fixtures/cert.pem")
    pem_private_key = file("test-fixtures/private-key.pem")
  }
}

resource "google_compute_network" "default" {
  name                    = "tf-test-my-network-%{random_suffix}"
  routing_mode            = "REGIONAL"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "tf-test-my-subnetwork-name-%{random_suffix}"
  purpose       = "PRIVATE"
  ip_cidr_range = "10.128.0.0/20"
  region        = "us-central1"
  network       = google_compute_network.default.id
  role          = "ACTIVE"
}

resource "google_compute_subnetwork" "proxyonlysubnet" {
  name          = "tf-test-my-proxy-only-subnetwork-%{random_suffix}"
  purpose       = "REGIONAL_MANAGED_PROXY"
  ip_cidr_range = "192.168.0.0/23"
  region        = "us-central1"
  network       = google_compute_network.default.id
  role          = "ACTIVE"
}

resource "google_network_security_gateway_security_policy" "default" {
  name        = "tf-test-my-policy-name-%{random_suffix}"
  location    = "us-central1"
}

resource "google_network_security_gateway_security_policy_rule" "default" {
  name                    = "tf-test-my-policyrule-name-%{random_suffix}"
  location                = "us-central1"
  gateway_security_policy = google_network_security_gateway_security_policy.default.name
  enabled                 = true
  priority                = 1
  session_matcher         = "host() == 'example.com'"
  basic_profile           = "ALLOW"
}

resource "google_network_services_gateway" "default" {
  name                                 = "tf-test-my-gateway-%{random_suffix}"
  location                             = "us-central1"`, context)

	if withAddresses {
		config += `
  addresses                            = ["10.128.0.99"]`
	}

	config += acctest.Nprintf(`
  type                                 = "SECURE_WEB_GATEWAY"
  ports                                = [443]
  scope                                = "tf-test-my-default-scope-%{random_suffix}"
  certificate_urls                     = [google_certificate_manager_certificate.default.id]
  gateway_security_policy              = google_network_security_gateway_security_policy.default.id
  network                              = google_compute_network.default.id
  subnetwork                           = google_compute_subnetwork.default.id
  delete_swg_autogen_router_on_destroy = true
  depends_on                           = [google_compute_subnetwork.proxyonlysubnet]
}
`, context)

	return config
}

func TestAccNetworkServicesGateway_swpUpdate(t *testing.T) {
	cmName := fmt.Sprintf("tf-test-gateway-swp-cm-%s", acctest.RandString(t, 10))
	netName := fmt.Sprintf("tf-test-gateway-swp-net-%s", acctest.RandString(t, 10))
	subnetName := fmt.Sprintf("tf-test-gateway-swp-subnet-%s", acctest.RandString(t, 10))
	pSubnetName := fmt.Sprintf("tf-test-gateway-swp-proxyonly-%s", acctest.RandString(t, 10))
	policyName := fmt.Sprintf("tf-test-gateway-swp-policy-%s", acctest.RandString(t, 10))
	ruleName := fmt.Sprintf("tf-test-gateway-swp-rule-%s", acctest.RandString(t, 10))
	gatewayScope := fmt.Sprintf("tf-test-gateway-swp-scope-%s", acctest.RandString(t, 10))
	gatewayName := fmt.Sprintf("tf-test-gateway-swp-%s", acctest.RandString(t, 10))
	serverTlsName := fmt.Sprintf("tf-test-gateway-swp-servertls-%s", acctest.RandString(t, 10))
	// updates
	newCmName := fmt.Sprintf("tf-test-gateway-swp-newcm-%s", acctest.RandString(t, 10))
	newPolicyName := fmt.Sprintf("tf-test-gateway-swp-newpolicy-%s", acctest.RandString(t, 10))
	newRuleName := fmt.Sprintf("tf-test-gateway-swp-newrule-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesGatewayDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesGateway_basicSwp(cmName, netName, subnetName, pSubnetName, policyName, ruleName, serverTlsName, gatewayName, gatewayScope),
			},
			{
				ResourceName:            "google_network_services_gateway.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "delete_swg_autogen_router_on_destroy"},
			},
			{
				Config: testAccNetworkServicesGateway_updateSwp(cmName, newCmName, netName, subnetName, pSubnetName, policyName, newPolicyName, ruleName, newRuleName, serverTlsName, gatewayName, gatewayScope),
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

func testAccNetworkServicesGateway_basicSwp(cmName, netName, subnetName, pSubnetName, policyName, ruleName, serverTlsName, gatewayName, gatewayScope string) string {
	return fmt.Sprintf(`
resource "google_certificate_manager_certificate" "default" {
  name        = "%s"
  location    = "us-east1"
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
  region        = "us-east1"
  network       = google_compute_network.default.id
  role          = "ACTIVE"
}

resource "google_compute_subnetwork" "default" {
  name          = "%s"
  purpose       = "PRIVATE"
  ip_cidr_range = "10.128.0.0/20"
  region        = "us-east1"
  network       = google_compute_network.default.id
  role          = "ACTIVE"
}

resource "google_network_security_gateway_security_policy" "default" {
  name     = "%s"
  location = "us-east1"
}

resource "google_network_security_gateway_security_policy_rule" "default" {
  name                    = "%s"
  location                = "us-east1"
  gateway_security_policy = google_network_security_gateway_security_policy.default.name
  enabled                 = true
  priority                = 1
  session_matcher         = "host() == 'example.com'"
  basic_profile           = "ALLOW"
}

resource "google_network_security_server_tls_policy" "servertls" {
  name                   = "%s"
  labels                 = {
    foo = "bar"
  }
  description            = "my description"
  location               = "us-east1"
  allow_open             = "false"
  mtls_policy {
    client_validation_mode = "ALLOW_INVALID_OR_MISSING_CLIENT_CERT"
  }
}

resource "google_network_services_gateway" "foobar" {
  name                                 = "%s"
  location                             = "us-east1"
  addresses                            = ["10.128.0.99"]
  type                                 = "SECURE_WEB_GATEWAY"
  ports                                = [443]
  description                          = "my description"
  scope                                = "%s"
  certificate_urls                     = [google_certificate_manager_certificate.default.id]
  gateway_security_policy              = google_network_security_gateway_security_policy.default.id
  network                              = google_compute_network.default.id
  subnetwork                           = google_compute_subnetwork.default.id
  delete_swg_autogen_router_on_destroy = true
  envoy_headers                        = "NONE"
  ip_version                           = "IPV4"
  server_tls_policy                    = google_network_security_server_tls_policy.servertls.id
  depends_on                           = [google_compute_subnetwork.proxyonlysubnet]
}

`, cmName, netName, subnetName, pSubnetName, policyName, ruleName, serverTlsName, gatewayName, gatewayScope)
}

func testAccNetworkServicesGateway_updateSwp(cmName, newCmName, netName, subnetName, pSubnetName, policyName, newPolicyName, ruleName, newRuleName, serverTlsName, gatewayName, gatewayScope string) string {
	return fmt.Sprintf(`
resource "google_certificate_manager_certificate" "default" {
  name        = "%s"
  location    = "us-east1"
  self_managed {
    pem_certificate = file("test-fixtures/cert.pem")
	  pem_private_key = file("test-fixtures/private-key.pem")
  }
}

resource "google_certificate_manager_certificate" "newcm" {
  name        = "%s"
  location    = "us-east1"
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
  region        = "us-east1"
  network       = google_compute_network.default.id
  role          = "ACTIVE"
}

resource "google_compute_subnetwork" "default" {
  name          = "%s"
  purpose       = "PRIVATE"
  ip_cidr_range = "10.128.0.0/20"
  region        = "us-east1"
  network       = google_compute_network.default.id
  role          = "ACTIVE"
}

resource "google_network_security_gateway_security_policy" "default" {
  name     = "%s"
  location = "us-east1"
}

resource "google_network_security_gateway_security_policy_rule" "default" {
  name                    = "%s"
  location                = "us-east1"
  gateway_security_policy = google_network_security_gateway_security_policy.default.name
  enabled                 = true
  priority                = 1
  session_matcher         = "host() == 'example.com'"
  basic_profile           = "ALLOW"
}

resource "google_network_security_gateway_security_policy" "newpolicy" {
  name     = "%s"
  location = "us-east1"
}

resource "google_network_security_gateway_security_policy_rule" "newrule" {
  name                    = "%s"
  location                = "us-east1"
  gateway_security_policy = google_network_security_gateway_security_policy.newpolicy.name
  enabled                 = true
  priority                = 1
  session_matcher         = "host() == 'example.com'"
  basic_profile           = "ALLOW"
}

resource "google_network_security_server_tls_policy" "servertls" {
  name                   = "%s"
  labels                 = {
    foo = "bar"
  }
  description            = "my description"
  location               = "us-east1"
  allow_open             = "false"
  mtls_policy {
    client_validation_mode = "ALLOW_INVALID_OR_MISSING_CLIENT_CERT"
  }
}

resource "google_network_services_gateway" "foobar" {
  name                                 = "%s"
  location                             = "us-east1"
  addresses                            = ["10.128.0.99"]
  type                                 = "SECURE_WEB_GATEWAY"
  ports                                = [443]
  description                          = "updated description"
  scope                                = "%s"
  certificate_urls                     = [google_certificate_manager_certificate.newcm.id]
  gateway_security_policy              = google_network_security_gateway_security_policy.newpolicy.id
  network                              = google_compute_network.default.id
  subnetwork                           = google_compute_subnetwork.default.id
  delete_swg_autogen_router_on_destroy = true
  envoy_headers                        = "NONE"
  ip_version                           = "IPV4"
  server_tls_policy                    = google_network_security_server_tls_policy.servertls.id
  depends_on                           = [google_compute_subnetwork.proxyonlysubnet]
}

`, cmName, newCmName, netName, subnetName, pSubnetName, policyName, newPolicyName, ruleName, newRuleName, serverTlsName, gatewayName, gatewayScope)
}

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
  name     = "%s"
  location = "us-west1"
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
  name     = "%s"
  location = "us-west2"
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
  name     = "%s"
  location = "us-west2"
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
  name     = "%s"
  location = "us-central1"
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

func TestAccNetworkServicesGateway_swpAsNextHop(t *testing.T) {
	context := map[string]interface{}{
		"region":        "us-east1",
		"random_suffix": fmt.Sprintf("-%s", acctest.RandString(t, 10)),
		"name_prefix":   "tf-test-gateway-",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesGatewayDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesGateway_swpAsNextHop(context),
			},
			{
				ResourceName:            "google_network_services_gateway.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "delete_swg_autogen_router_on_destroy"},
			},
		},
	})
}

func testAccNetworkServicesGateway_swpAsNextHop(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "default" {
  name                    = "%{name_prefix}network%{random_suffix}"
  routing_mode            = "REGIONAL"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "proxyonlysubnet" {
  name          = "%{name_prefix}proxysubnet%{random_suffix}"
  purpose       = "REGIONAL_MANAGED_PROXY"
  ip_cidr_range = "192.168.0.0/23"
  region        = "%{region}"
  network       = google_compute_network.default.id
  role          = "ACTIVE"
}

resource "google_compute_subnetwork" "default" {
  name          = "%{name_prefix}subnet%{random_suffix}"
  purpose       = "PRIVATE"
  ip_cidr_range = "10.128.0.0/20"
  region        = "%{region}"
  network       = google_compute_network.default.id
  role          = "ACTIVE"
}

resource "google_privateca_ca_pool" "default" {
  name     = "%{name_prefix}ca-pool%{random_suffix}"
  location = "%{region}"
  tier     = "DEVOPS"

  publishing_options {
    publish_ca_cert = false
    publish_crl     = false
  }

  issuance_policy {
    maximum_lifetime = "1209600s"
    baseline_values {
      ca_options {
        is_ca = false
      }
      key_usage {
        base_key_usage {}
        extended_key_usage {
          server_auth = true
        }
      }
    }
  }
}
  
resource "google_privateca_certificate_authority" "default" {
  pool                                   = google_privateca_ca_pool.default.name
  certificate_authority_id               = "%{name_prefix}certificate-authority%{random_suffix}"
  location                               = "%{region}"
  lifetime                               = "86400s"
  type                                   = "SELF_SIGNED"
  deletion_protection                    = false
  skip_grace_period                      = true
  ignore_active_certificates_on_deletion = true

  config {
    subject_config {
      subject {
        organization = "Test LLC"
        common_name  = "private-certificate-authority"
      }
    }
    x509_config {
      ca_options {
        is_ca = true
      }
      key_usage {
        base_key_usage {
          cert_sign = true
          crl_sign  = true
        }
        extended_key_usage {
          server_auth = false
        }
      }
    }
  }

  key_spec {
    algorithm = "RSA_PKCS1_4096_SHA256"
  }
}

resource "google_certificate_manager_certificate" "default" {
  name     = "%{name_prefix}certificate%{random_suffix}"
  location = "%{region}"

  self_managed {
    pem_certificate = file("test-fixtures/cert.pem")
    pem_private_key = file("test-fixtures/private-key.pem")
  }
}

resource "google_network_security_tls_inspection_policy" "default" {
  name     = "%{name_prefix}tls-insp-policy%{random_suffix}"
  location = "%{region}"
  ca_pool  = google_privateca_ca_pool.default.id

  depends_on = [
    google_privateca_ca_pool.default,
    google_privateca_certificate_authority.default
  ]
}

resource "google_network_security_gateway_security_policy" "default" {
  name                  = "%{name_prefix}sec-policy%{random_suffix}"
  location              = "%{region}"
  description           = "my description"
  tls_inspection_policy = google_network_security_tls_inspection_policy.default.id

  depends_on = [
    google_network_security_tls_inspection_policy.default
  ]
}

resource "google_network_security_gateway_security_policy_rule" "default" {
  name                    = "%{name_prefix}sec-policy-rule%{random_suffix}"
  location                = "%{region}"
  gateway_security_policy = google_network_security_gateway_security_policy.default.name
  enabled                 = true
  description             = "my description"
  priority                = 0
  session_matcher         = "host() == 'example.com'"
  application_matcher     = "request.method == 'POST'"
  tls_inspection_enabled  = true
  basic_profile           = "ALLOW"
}

resource "google_network_services_gateway" "default" {
  name                                 = "%{name_prefix}swp%{random_suffix}"
  location                             = "%{region}"
  addresses                            = ["10.128.0.99"]
  type                                 = "SECURE_WEB_GATEWAY"
  routing_mode                         = "NEXT_HOP_ROUTING_MODE"
  ports                                = [443]
  description                          = "my description"
  scope                                = "%s"
  certificate_urls                     = [google_certificate_manager_certificate.default.id]
  gateway_security_policy              = google_network_security_gateway_security_policy.default.id
  network                              = google_compute_network.default.id
  subnetwork                           = google_compute_subnetwork.default.id
  delete_swg_autogen_router_on_destroy = true
  depends_on                           = [google_compute_subnetwork.proxyonlysubnet]
}

resource "google_compute_route" "default" {
  name        = "%{name_prefix}route%{random_suffix}"
  dest_range  = "15.0.0.0/24"
  network     = google_compute_network.default.name
  next_hop_ip = google_network_services_gateway.default.addresses[0]
  priority    = 100
}

resource "google_network_connectivity_policy_based_route" "swproute" {
  name            = "%{name_prefix}policy-based-swp-route%{random_suffix}"
  description     = "My routing policy"
  network         = google_compute_network.default.id
  next_hop_ilb_ip = google_network_services_gateway.default.addresses[0]
  priority        = 2

  filter {
    protocol_version = "IPV4"
    src_range        = "10.0.0.0/24"
    dest_range       = "15.0.0.0/24"
  }
}

resource "google_network_connectivity_policy_based_route" "default" {
  name                  = "%{name_prefix}policy-based-route%{random_suffix}"
  description           = "My routing policy"
  network               = google_compute_network.default.id
  next_hop_other_routes = "DEFAULT_ROUTING"
  priority              = 1

  filter {
    protocol_version = "IPV4"
    src_range        = "10.0.0.0/24"
    dest_range       = "15.0.0.0/24"
  }
}
	`, context)
}
