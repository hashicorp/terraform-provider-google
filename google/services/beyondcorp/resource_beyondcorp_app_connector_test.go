// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package beyondcorp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccBeyondcorpAppConnector_beyondcorpAppConnectorUpdateExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckBeyondcorpAppConnectorDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccBeyondcorpAppConnector_beyondcorpAppConnectorBasicExample(context),
			},
			{
				ResourceName:            "google_beyondcorp_app_connector.app_connector",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "region"},
			},
			{
				Config: testAccBeyondcorpAppConnector_beyondcorpAppConnectorUpdateExample(context),
			},
			{
				ResourceName:            "google_beyondcorp_app_connector.app_connector",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "region"},
			},
			{
				Config: testAccBeyondcorpAppConnector_beyondcorpAppConnectorBasicExample(context),
			},
		},
	})
}

func testAccBeyondcorpAppConnector_beyondcorpAppConnectorUpdateExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_service_account" "service_account" {
  account_id   = "tf-test-my-account%{random_suffix}"
  display_name = "Test Service Account"
}

resource "google_beyondcorp_app_connector" "app_connector" {
  name = "tf-test-my-app-connector%{random_suffix}"
  principal_info {
    service_account {
     email = google_service_account.service_account.email
    }
  }
  display_name = "Some display name"
}
`, context)
}
