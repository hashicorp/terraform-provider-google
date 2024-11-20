// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package beyondcorp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccBeyondcorpAppConnection_beyondcorpAppConnectionUpdateExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {},
		},
		CheckDestroy: testAccCheckBeyondcorpAppConnectionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBeyondcorpAppConnection_beyondcorpAppConnectionBasicExample(context),
			},
			{
				ResourceName:            "google_beyondcorp_app_connection.app_connection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "region", "gateway"},
			},
			{
				Config: testAccBeyondcorpAppConnection_beyondcorpAppConnectionUpdateExample(context),
			},
			{
				ResourceName:            "google_beyondcorp_app_connection.app_connection",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "region", "gateway"},
			},
			{
				Config: testAccBeyondcorpAppConnection_beyondcorpAppConnectionBasicExample(context),
			},
		},
	})
}

func testAccBeyondcorpAppConnection_beyondcorpAppConnectionUpdateExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "service_account" {
  account_id   = "tf-test-my-account%{random_suffix}"
  display_name = "Test Service Account"
}

resource "time_sleep" "wait_120_seconds" {
  depends_on = [google_service_account.service_account]

  create_duration = "120s"
}


resource "google_beyondcorp_app_connector" "app_connector" {
  depends_on = [time_sleep.wait_120_seconds]  

  name = "tf-test-my-app-connector%{random_suffix}"
  principal_info {
    service_account {
     email = google_service_account.service_account.email
    }
  }
}

resource "google_beyondcorp_app_connection" "app_connection" {
  name = "tf-test-my-app-connection%{random_suffix}"
  type = "TCP_PROXY"
  region = "us-central1"
  application_endpoint {
    host = "foo-host"
    port = 8080
  }
  connectors = [google_beyondcorp_app_connector.app_connector.id]
  display_name = "Some display name"
}
`, context)
}
