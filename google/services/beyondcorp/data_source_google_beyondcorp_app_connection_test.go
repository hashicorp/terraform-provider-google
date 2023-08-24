// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package beyondcorp_test

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceGoogleBeyondcorpAppConnection_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBeyondcorpAppConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBeyondcorpAppConnection_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_beyondcorp_app_connection.foo", "google_beyondcorp_app_connection.foo"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleBeyondcorpAppConnection_full(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBeyondcorpAppConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBeyondcorpAppConnection_full(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState("data.google_beyondcorp_app_connection.foo", "google_beyondcorp_app_connection.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleBeyondcorpAppConnection_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "service_account" {
	account_id   = "tf-test-my-account%{random_suffix}"
	display_name = "Test Service Account"
}

resource "google_beyondcorp_app_connector" "app_connector" {
	name = "tf-test-appconnector-%{random_suffix}"
	principal_info {
		service_account {
			email = google_service_account.service_account.email
		}
	}
}

resource "google_beyondcorp_app_connection" "foo" {
	name = "tf-test-my-app-connection-%{random_suffix}"
	type = "TCP_PROXY"
	application_endpoint {
		host = "foo-host"
		port = 8080
	}
	connectors = [google_beyondcorp_app_connector.app_connector.id]
}

data "google_beyondcorp_app_connection" "foo" {
	name = google_beyondcorp_app_connection.foo.name
}
`, context)
}

func testAccDataSourceGoogleBeyondcorpAppConnection_full(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "service_account" {
	account_id   = "tf-test-my-account%{random_suffix}"
	display_name = "Test Service Account"
}

resource "google_beyondcorp_app_connector" "app_connector" {
	name = "tf-test-appconnector-%{random_suffix}"
	principal_info {
		service_account {
			email = google_service_account.service_account.email
		}
	}
}

resource "google_beyondcorp_app_connection" "foo" {
	name = "tf-test-my-app-connection-%{random_suffix}"
	type = "TCP_PROXY"
	application_endpoint {
		host = "foo-host"
		port = 8080
	}
	connectors = [google_beyondcorp_app_connector.app_connector.id]
}

data "google_beyondcorp_app_connection" "foo" {
	name    = google_beyondcorp_app_connection.foo.name
	project = google_beyondcorp_app_connection.foo.project
	region  = google_beyondcorp_app_connection.foo.region
}
`, context)
}
