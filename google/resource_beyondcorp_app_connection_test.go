package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBeyondcorpAppConnection_beyondcorpAppConnectionUpdateExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
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
	return Nprintf(`
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
