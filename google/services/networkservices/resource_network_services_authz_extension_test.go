// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networkservices_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccNetworkServicesAuthzExtension_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project":       envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesAuthzExtensionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesAuthzExtension_start(context),
			},
			{
				ResourceName:            "google_network_services_authz_extension.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "service", "terraform_labels"},
			},
			{
				Config: testAccNetworkServicesAuthzExtension_update(context),
			},
			{
				ResourceName:            "google_network_services_authz_extension.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "location", "service", "terraform_labels"},
			},
		},
	})
}

func testAccNetworkServicesAuthzExtension_start(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "default" {
  name                    = "tf-test-lb-network%{random_suffix}"
  project                 = "%{project}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "tf-test-backend-subnet%{random_suffix}"
  project       = "%{project}"
  region        = "us-west1"
  ip_cidr_range = "10.1.2.0/24"
  network       = google_compute_network.default.id
}

resource "google_compute_subnetwork" "proxy_only" {
  name          = "tf-test-proxy-only-subnet%{random_suffix}"
  project       = "%{project}"
  region        = "us-west1"
  ip_cidr_range = "10.129.0.0/23"
  purpose       = "REGIONAL_MANAGED_PROXY"
  role          = "ACTIVE"
  network       = google_compute_network.default.id
}

resource "google_compute_address" "default" {
  name         = "tf-test-l7-ilb-ip-address%{random_suffix}"
  project      = "%{project}"
  region       = "us-west1"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  purpose      = "GCE_ENDPOINT"
}


resource "google_compute_region_health_check" "default" {
  name    = "tf-test-l7-ilb-basic-check%{random_suffix}"
  project = "%{project}"
  region  = "us-west1"

  http_health_check {
    port_specification = "USE_SERVING_PORT"
  }
}

resource "google_compute_region_backend_service" "url_map" {
  name                  = "tf-test-l7-ilb-backend-service%{random_suffix}"
  project               = "%{project}"
  region                = "us-west1"
  load_balancing_scheme = "INTERNAL_MANAGED"

  health_checks = [google_compute_region_health_check.default.id]
}

resource "google_compute_forwarding_rule" "default" {
  name                  = "tf-test-l7-ilb-forwarding-rule%{random_suffix}"
  project               = "%{project}"
  region                = "us-west1"
  load_balancing_scheme = "INTERNAL_MANAGED"
  network               = google_compute_network.default.id
  subnetwork            = google_compute_subnetwork.default.id
  ip_protocol           = "TCP"
  port_range            = "80"
  target                = google_compute_region_target_http_proxy.default.id
  ip_address            = google_compute_address.default.id

  depends_on = [google_compute_subnetwork.proxy_only]
}

resource "google_compute_region_url_map" "default" {
  name            = "tf-test-l7-ilb-map%{random_suffix}"
  project         = "%{project}"
  region          = "us-west1"
  default_service = google_compute_region_backend_service.url_map.id
}

resource "google_compute_region_target_http_proxy" "default" {
  name    = "tf-test-l7-ilb-proxy%{random_suffix}"
  project = "%{project}"
  region  = "us-west1"
  url_map = google_compute_region_url_map.default.id
}

resource "google_compute_region_backend_service" "default" {
  name                  = "tf-test-authz-service%{random_suffix}"
  project               = "%{project}"
  region                = "us-west1"

  protocol              = "HTTP2"
  load_balancing_scheme = "INTERNAL_MANAGED"
  port_name             = "grpc"
}

resource "google_compute_region_backend_service" "updated" {
  name                  = "tf-test-authz-service-updated%{random_suffix}"
  project               = "%{project}"
  region                = "us-west1"
  
  protocol              = "HTTP2"
  load_balancing_scheme = "INTERNAL_MANAGED"
  port_name             = "grpc"
}

resource "google_network_services_authz_extension" "default" {
  name     = "tf-test-my-authz-ext%{random_suffix}"
  project  = "%{project}"
  location = "us-west1"

  description           = "my description"
  load_balancing_scheme = "INTERNAL_MANAGED"
  authority             = "ext11.com"
  service               = google_compute_region_backend_service.default.self_link
  timeout               = "0.1s"
  fail_open             = false
  forward_headers       = ["Authorization"]
}
`, context)
}

func testAccNetworkServicesAuthzExtension_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_compute_network" "default" {
  name                    = "tf-test-lb-network%{random_suffix}"
  project                 = "%{project}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "default" {
  name          = "tf-test-backend-subnet%{random_suffix}"
  project       = "%{project}"
  region        = "us-west1"
  ip_cidr_range = "10.1.2.0/24"
  network       = google_compute_network.default.id
}

resource "google_compute_subnetwork" "proxy_only" {
  name          = "tf-test-proxy-only-subnet%{random_suffix}"
  project       = "%{project}"
  region        = "us-west1"
  ip_cidr_range = "10.129.0.0/23"
  purpose       = "REGIONAL_MANAGED_PROXY"
  role          = "ACTIVE"
  network       = google_compute_network.default.id
}

resource "google_compute_address" "default" {
  name         = "tf-test-l7-ilb-ip-address%{random_suffix}"
  project      = "%{project}"
  region       = "us-west1"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  purpose      = "GCE_ENDPOINT"
}

resource "google_compute_region_health_check" "default" {
  name    = "tf-test-l7-ilb-basic-check%{random_suffix}"
  project = "%{project}"
  region  = "us-west1"

  http_health_check {
    port_specification = "USE_SERVING_PORT"
  }
}

resource "google_compute_region_backend_service" "url_map" {
  name                  = "tf-test-l7-ilb-backend-service%{random_suffix}"
  project               = "%{project}"
  region                = "us-west1"
  load_balancing_scheme = "INTERNAL_MANAGED"

  health_checks = [google_compute_region_health_check.default.id]
}

resource "google_compute_forwarding_rule" "default" {
  name                  = "tf-test-l7-ilb-forwarding-rule%{random_suffix}"
  project               = "%{project}"
  region                = "us-west1"
  load_balancing_scheme = "INTERNAL_MANAGED"
  network               = google_compute_network.default.id
  subnetwork            = google_compute_subnetwork.default.id
  ip_protocol           = "TCP"
  port_range            = "80"
  target                = google_compute_region_target_http_proxy.default.id
  ip_address            = google_compute_address.default.id

  depends_on = [google_compute_subnetwork.proxy_only]
}

resource "google_compute_region_url_map" "default" {
  name            = "tf-test-l7-ilb-map%{random_suffix}"
  project         = "%{project}"
  region          = "us-west1"
  default_service = google_compute_region_backend_service.url_map.id
}

resource "google_compute_region_target_http_proxy" "default" {
  name    = "tf-test-l7-ilb-proxy%{random_suffix}"
  project = "%{project}"
  region  = "us-west1"
  url_map = google_compute_region_url_map.default.id
}

resource "google_compute_region_backend_service" "default" {
  name                  = "tf-test-authz-service%{random_suffix}"
  project               = "%{project}"
  region                = "us-west1"

  protocol              = "HTTP2"
  load_balancing_scheme = "INTERNAL_MANAGED"
  port_name             = "grpc"
}

resource "google_compute_region_backend_service" "updated" {
  name                  = "tf-test-authz-service-updated%{random_suffix}"
  project               = "%{project}"
  region                = "us-west1"
  
  protocol              = "HTTP2"
  load_balancing_scheme = "INTERNAL_MANAGED"
  port_name             = "grpc"
}

resource "google_network_services_authz_extension" "default" {
  name     = "tf-test-my-authz-ext%{random_suffix}"
  project  = "%{project}"
  location = "us-west1"

  description           = "updated description"
  load_balancing_scheme = "INTERNAL_MANAGED"
  authority             = "ext11.com"
  service               = google_compute_region_backend_service.updated.self_link
  timeout               = "0.1s"
  fail_open             = false
  forward_headers       = ["Authorization"]

  metadata = {
    forwarding_rule_id = google_compute_forwarding_rule.default.id
  }

  labels = {
    foo = "bar"
  }
}
`, context)
}
