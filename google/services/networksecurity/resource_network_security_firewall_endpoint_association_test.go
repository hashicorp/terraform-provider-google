// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networksecurity_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccNetworkSecurityFirewallEndpointAssociations_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"orgId":            envvar.GetTestOrgFromEnv(t),
		"randomSuffix":     acctest.RandString(t, 10),
		"billingProjectId": envvar.GetTestProjectFromEnv(),
		"disabled":         strconv.FormatBool(false),
	}

	testResourceName := "google_network_security_firewall_endpoint_association.foobar"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkSecurityFirewallEndpointDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSecurityFirewallEndpointAssociation_basic(context),
			},
			{
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkSecurityFirewallEndpointAssociation_update(context),
			},
			{
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func TestAccNetworkSecurityFirewallEndpointAssociations_disabled(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"orgId":            envvar.GetTestOrgFromEnv(t),
		"randomSuffix":     acctest.RandString(t, 10),
		"billingProjectId": envvar.GetTestProjectFromEnv(),
	}

	testResourceName := "google_network_security_firewall_endpoint_association.foobar"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkSecurityFirewallEndpointDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkSecurityFirewallEndpointAssociation_basic(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "disabled", "false"),
				),
			},
			{
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkSecurityFirewallEndpointAssociation_update(testContextMapDisabledField(context, true)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "disabled", "true"),
				),
			},
			{
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkSecurityFirewallEndpointAssociation_update(testContextMapDisabledField(context, false)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(testResourceName, "disabled", "false"),
				),
			},
			{
				ResourceName:            testResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testContextMapDisabledField(context map[string]interface{}, disabled bool) map[string]interface{} {
	context["disabled"] = strconv.FormatBool(disabled)
	return context
}

func testAccNetworkSecurityFirewallEndpointAssociation_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "foobar" {
  name                    = "tf-test-my-vpc%{randomSuffix}"
  auto_create_subnetworks = false
}

resource "google_network_security_firewall_endpoint" "foobar" {
  name               = "tf-test-my-firewall-endpoint%{randomSuffix}"
  parent             = "organizations/%{orgId}"
  location           = "us-central1-a"
  billing_project_id = "%{billingProjectId}"
}

# TODO: add tlsInspectionPolicy once resource is ready
resource "google_network_security_firewall_endpoint_association" "foobar" {
  name              = "tf-test-my-firewall-endpoint-association%{randomSuffix}"
  parent            = "projects/%{billingProjectId}"
  location          = "us-central1-a"
  firewall_endpoint = google_network_security_firewall_endpoint.foobar.id
  network           = google_compute_network.foobar.id

  labels = {
    foo = "bar"
  }
}
`, context)
}

func testAccNetworkSecurityFirewallEndpointAssociation_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "foobar" {
  name                    = "tf-test-my-vpc%{randomSuffix}"
  auto_create_subnetworks = false
}

resource "google_network_security_firewall_endpoint" "foobar" {
  name               = "tf-test-my-firewall-endpoint%{randomSuffix}"
  parent             = "organizations/%{orgId}"
  location           = "us-central1-a"
  billing_project_id = "%{billingProjectId}"
}

# TODO: add tlsInspectionPolicy once resource is ready
resource "google_network_security_firewall_endpoint_association" "foobar" {
  name              = "tf-test-my-firewall-endpoint-association%{randomSuffix}"
  parent            = "projects/%{billingProjectId}"
  location          = "us-central1-a"
  firewall_endpoint = google_network_security_firewall_endpoint.foobar.id
  network           = google_compute_network.foobar.id
  disabled          = "%{disabled}"

  labels = {
    foo = "bar-updated"
  }
}
`, context)
}

func testAccCheckNetworkSecurityFirewallEndpointAssociationDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_network_security_firewall_endpoint_association" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{NetworkSecurityBasePath}}{{parent}}/locations/{{location}}/firewallEndpointAssociations/{{name}}")
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
				return fmt.Errorf("NetworkSecurityFirewallEndpointAssociation still exists at %s", url)
			}
		}

		return nil
	}
}
