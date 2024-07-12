// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package networkconnectivity_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccNetworkConnectivityHub_BasicHubLongForm(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkConnectivityHubDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConnectivityHub_BasicHubLongForm(context),
			},
			{
				ResourceName:            "google_network_connectivity_hub.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "project"},
			},
			{
				Config: testAccNetworkConnectivityHub_BasicHubLongFormUpdate0(context),
			},
			{
				ResourceName:            "google_network_connectivity_hub.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels", "project"},
			},
		},
	})
}
func TestAccNetworkConnectivityHub_BasicHub(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkConnectivityHubDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkConnectivityHub_BasicHub(context),
			},
			{
				ResourceName:            "google_network_connectivity_hub.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkConnectivityHub_BasicHubUpdate0(context),
			},
			{
				ResourceName:            "google_network_connectivity_hub.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccNetworkConnectivityHub_BasicHubLongForm(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_connectivity_hub" "primary" {
  name        = "tf-test-hub%{random_suffix}"
  description = "A sample hub"
  project     = "projects/%{project_name}"

  labels = {
    label-one = "value-one"
  }
}


`, context)
}

func testAccNetworkConnectivityHub_BasicHubLongFormUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_connectivity_hub" "primary" {
  name        = "tf-test-hub%{random_suffix}"
  description = "An updated sample hub"
  project     = "projects/%{project_name}"

  labels = {
    label-two = "value-one"
  }
}


`, context)
}

func testAccNetworkConnectivityHub_BasicHub(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_connectivity_hub" "primary" {
  name        = "tf-test-hub%{random_suffix}"
  description = "A sample hub"
  project     = "%{project_name}"

  labels = {
    label-one = "value-one"
  }
}


`, context)
}

func testAccNetworkConnectivityHub_BasicHubUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_network_connectivity_hub" "primary" {
  name        = "tf-test-hub%{random_suffix}"
  description = "An updated sample hub"
  project     = "%{project_name}"

  labels = {
    label-two = "value-one"
  }
}


`, context)
}
