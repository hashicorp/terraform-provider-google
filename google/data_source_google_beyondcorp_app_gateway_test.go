package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGoogleBeyondcorpAppGateway_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBeyondcorpAppGatewayDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBeyondcorpAppGateway_basic(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_beyondcorp_app_gateway.foo", "google_beyondcorp_app_gateway.foo"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleBeyondcorpAppGateway_optionalProject(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBeyondcorpAppGatewayDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBeyondcorpAppGateway_optionalProject(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_beyondcorp_app_gateway.foo", "google_beyondcorp_app_gateway.foo"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleBeyondcorpAppGateway_optionalRegion(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBeyondcorpAppGatewayDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBeyondcorpAppGateway_optionalRegion(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_beyondcorp_app_gateway.foo", "google_beyondcorp_app_gateway.foo"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleBeyondcorpAppGateway_optionalProjectRegion(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBeyondcorpAppGatewayDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleBeyondcorpAppGateway_optionalProjectRegion(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_beyondcorp_app_gateway.foo", "google_beyondcorp_app_gateway.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleBeyondcorpAppGateway_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_beyondcorp_app_gateway" "foo" {
	name      = "tf-test-appgateway-%{random_suffix}"
	type      = "TCP_PROXY"
	host_type = "GCP_REGIONAL_MIG"
}

data "google_beyondcorp_app_gateway" "foo" {
	name    = google_beyondcorp_app_gateway.foo.name
	project = google_beyondcorp_app_gateway.foo.project
	region  = google_beyondcorp_app_gateway.foo.region
}
`, context)
}

func testAccDataSourceGoogleBeyondcorpAppGateway_optionalProject(context map[string]interface{}) string {
	return Nprintf(`
resource "google_beyondcorp_app_gateway" "foo" {
	name      = "tf-test-appgateway-%{random_suffix}"
	type      = "TCP_PROXY"
	host_type = "GCP_REGIONAL_MIG"
}

data "google_beyondcorp_app_gateway" "foo" {
	name   = google_beyondcorp_app_gateway.foo.name
	region = google_beyondcorp_app_gateway.foo.region
}
`, context)
}

func testAccDataSourceGoogleBeyondcorpAppGateway_optionalRegion(context map[string]interface{}) string {
	return Nprintf(`
resource "google_beyondcorp_app_gateway" "foo" {
	name      = "tf-test-appgateway-%{random_suffix}"
	type      = "TCP_PROXY"
	host_type = "GCP_REGIONAL_MIG"
}

data "google_beyondcorp_app_gateway" "foo" {
	name    = google_beyondcorp_app_gateway.foo.name
	project = google_beyondcorp_app_gateway.foo.project
}
`, context)
}

func testAccDataSourceGoogleBeyondcorpAppGateway_optionalProjectRegion(context map[string]interface{}) string {
	return Nprintf(`
resource "google_beyondcorp_app_gateway" "foo" {
	name      = "tf-test-appgateway-%{random_suffix}"
	type      = "TCP_PROXY"
	host_type = "GCP_REGIONAL_MIG"
}

data "google_beyondcorp_app_gateway" "foo" {
	name = google_beyondcorp_app_gateway.foo.name
}
`, context)
}
