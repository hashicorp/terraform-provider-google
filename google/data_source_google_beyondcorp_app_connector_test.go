package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleBeyondcorpAppConnector_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBeyondcorpAppConnectorDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBeyondcorpAppConnector_basic(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_beyondcorp_app_connector.foo", "google_beyondcorp_app_connector.foo"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleBeyondcorpAppConnector_optionalProject(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBeyondcorpAppConnectorDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBeyondcorpAppConnector_optionalProject(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_beyondcorp_app_connector.foo", "google_beyondcorp_app_connector.foo"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleBeyondcorpAppConnector_optionalRegion(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBeyondcorpAppConnectorDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBeyondcorpAppConnector_optionalRegion(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_beyondcorp_app_connector.foo", "google_beyondcorp_app_connector.foo"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleBeyondcorpAppConnector_optionalProjectRegion(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBeyondcorpAppConnectorDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBeyondcorpAppConnector_optionalProjectRegion(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_beyondcorp_app_connector.foo", "google_beyondcorp_app_connector.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleBeyondcorpAppConnector_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_service_account" "service_account" {
	account_id   = "tf-test-my-account%{random_suffix}"
	display_name = "Test Service Account"
}

resource "google_beyondcorp_app_connector" "foo" {
 	name = "tf-test-appconnector-%{random_suffix}"
	principal_info {
		service_account {
			email = google_service_account.service_account.email
		}
	}
}

data "google_beyondcorp_app_connector" "foo" {
	name    = google_beyondcorp_app_connector.foo.name
	project = google_beyondcorp_app_connector.foo.project
	region  = google_beyondcorp_app_connector.foo.region
}
`, context)
}

func testAccDataSourceGoogleBeyondcorpAppConnector_optionalProject(context map[string]interface{}) string {
	return Nprintf(`
resource "google_service_account" "service_account" {
	account_id   = "tf-test-my-account%{random_suffix}"
	display_name = "Test Service Account"
}

resource "google_beyondcorp_app_connector" "foo" {
	name = "tf-test-appconnector-%{random_suffix}"
	principal_info {
		service_account {
			email = google_service_account.service_account.email
		}
	}
}

data "google_beyondcorp_app_connector" "foo" {
	name   = google_beyondcorp_app_connector.foo.name
	region = google_beyondcorp_app_connector.foo.region
}
`, context)
}

func testAccDataSourceGoogleBeyondcorpAppConnector_optionalRegion(context map[string]interface{}) string {
	return Nprintf(`
resource "google_service_account" "service_account" {
	account_id   = "tf-test-my-account%{random_suffix}"
	display_name = "Test Service Account"
}

resource "google_beyondcorp_app_connector" "foo" {
	name = "tf-test-appconnector-%{random_suffix}"
	principal_info {
		service_account {
			email = google_service_account.service_account.email
		}
	}
}

data "google_beyondcorp_app_connector" "foo" {
	name    = google_beyondcorp_app_connector.foo.name
	project = google_beyondcorp_app_connector.foo.project
}
`, context)
}

func testAccDataSourceGoogleBeyondcorpAppConnector_optionalProjectRegion(context map[string]interface{}) string {
	return Nprintf(`
resource "google_service_account" "service_account" {
	account_id   = "tf-test-my-account%{random_suffix}"
	display_name = "Test Service Account"
}

resource "google_beyondcorp_app_connector" "foo" {
	name = "tf-test-appconnector-%{random_suffix}"
	principal_info {
		service_account {
			email = google_service_account.service_account.email
		}
	}
}

data "google_beyondcorp_app_connector" "foo" {
	name = google_beyondcorp_app_connector.foo.name
}
`, context)
}
