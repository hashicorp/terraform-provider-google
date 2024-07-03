// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package networkservices_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetworkServicesLbRouteExtension_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetworkServicesLbRouteExtensionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkServicesLbRouteExtension_basic(context),
			},
			{
				ResourceName:            "google_network_services_lb_route_extension.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetworkServicesLbRouteExtension_update(context),
			},
			{
				ResourceName:            "google_network_services_lb_route_extension.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccNetworkServicesLbRouteExtension_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
# Internal HTTP load balancer with a managed instance group backend
# VPC network
resource "google_compute_network" "ilb_network" {
  name                    = "tf-test-l7-ilb-network%{random_suffix}"
  auto_create_subnetworks = false
}

# proxy-only subnet
resource "google_compute_subnetwork" "proxy_subnet" {
  name          = "tf-test-l7-ilb-proxy-subnet%{random_suffix}"
  ip_cidr_range = "10.0.0.0/24"
  region        = "us-west1"
  purpose       = "REGIONAL_MANAGED_PROXY"
  role          = "ACTIVE"
  network       = google_compute_network.ilb_network.id
}

# backend subnet
resource "google_compute_subnetwork" "ilb_subnet" {
  name          = "tf-test-l7-ilb-subnet%{random_suffix}"
  ip_cidr_range = "10.0.1.0/24"
  region        = "us-west1"
  network       = google_compute_network.ilb_network.id

  depends_on = [
    google_compute_subnetwork.proxy_subnet
  ]
}

# forwarding rule
resource "google_compute_forwarding_rule" "default" {
  name                  = "tf-test-l7-ilb-forwarding-rule%{random_suffix}"
  region                = "us-west1"
  ip_protocol           = "TCP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  port_range            = "80"
  target                = google_compute_region_target_http_proxy.default.id
  network               = google_compute_network.ilb_network.id
  subnetwork            = google_compute_subnetwork.ilb_subnet.id
  network_tier          = "PREMIUM"

  depends_on = [
    google_compute_subnetwork.proxy_subnet
  ]
}

# HTTP target proxy
resource "google_compute_region_target_http_proxy" "default" {
  name     = "tf-test-l7-ilb-target-http-proxy%{random_suffix}"
  region   = "us-west1"
  url_map  = google_compute_region_url_map.default.id
}

# URL map
resource "google_compute_region_url_map" "default" {
  name            = "tf-test-l7-ilb-regional-url-map%{random_suffix}"
  region          = "us-west1"
  default_service = google_compute_region_backend_service.default.id

  host_rule {
    hosts        = ["service-extensions.com"]
    path_matcher = "callouts"
  }

  path_matcher {
    name            = "callouts"
    default_service = google_compute_region_backend_service.callouts_backend.id
  }
}

# backend service
resource "google_compute_region_backend_service" "default" {
  name                  = "tf-test-l7-ilb-backend-subnet%{random_suffix}"
  region                = "us-west1"
  protocol              = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  timeout_sec           = 10
  health_checks         = [google_compute_region_health_check.default.id]

  backend {
    group           = google_compute_region_instance_group_manager.mig.instance_group
    balancing_mode  = "UTILIZATION"
    capacity_scaler = 1.0
  }
}

# instance template
resource "google_compute_instance_template" "instance_template" {
  name         = "tf-test-l7-ilb-mig-template%{random_suffix}"
  machine_type = "e2-small"
  tags         = ["http-server"]

  network_interface {
    network    = google_compute_network.ilb_network.id
    subnetwork = google_compute_subnetwork.ilb_subnet.id

    access_config {
      # add external ip to fetch packages
    }
  }

  disk {
    source_image = "debian-cloud/debian-12"
    auto_delete  = true
    boot         = true
  }

  # install nginx and serve a simple web page
  metadata = {
    startup-script = <<-EOF1
      #! /bin/bash
      set -euo pipefail

      export DEBIAN_FRONTEND=noninteractive
      apt-get update
      apt-get install -y nginx-light jq

      NAME=$(curl -H "Metadata-Flavor: Google" "http://metadata.google.internal/computeMetadata/v1/instance/hostname")
      IP=$(curl -H "Metadata-Flavor: Google" "http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/ip")
      METADATA=$(curl -f -H "Metadata-Flavor: Google" "http://metadata.google.internal/computeMetadata/v1/instance/attributes/?recursive=True" | jq 'del(.["startup-script"])')

      cat <<EOF > /var/www/html/index.html
      <pre>
      Name: $NAME
      IP: $IP
      Metadata: $METADATA
      </pre>
      EOF
    EOF1
  }

  lifecycle {
    create_before_destroy = true
  }
}

# health check
resource "google_compute_region_health_check" "default" {
  name     = "tf-test-l7-ilb-hc%{random_suffix}"
  region   = "us-west1"

  http_health_check {
    port_specification = "USE_SERVING_PORT"
  }
}

# MIG
resource "google_compute_region_instance_group_manager" "mig" {
  name     = "tf-test-l7-ilb-mig1%{random_suffix}"
  region   = "us-west1"

  base_instance_name = "vm"
  target_size        = 2

  version {
    instance_template = google_compute_instance_template.instance_template.id
    name              = "primary"
  }
}

# allow all access from IAP and health check ranges
resource "google_compute_firewall" "fw_iap" {
  name          = "tf-test-l7-ilb-fw-allow-iap-hc%{random_suffix}"
  direction     = "INGRESS"
  network       = google_compute_network.ilb_network.id
  source_ranges = ["130.211.0.0/22", "35.191.0.0/16", "35.235.240.0/20"]

  allow {
    protocol = "tcp"
  }
}

# allow http from proxy subnet to backends
resource "google_compute_firewall" "fw_ilb_to_backends" {
  name          = "tf-test-l7-ilb-fw-allow-ilb-to-backends%{random_suffix}"
  direction     = "INGRESS"
  network       = google_compute_network.ilb_network.id
  source_ranges = ["10.0.0.0/24"]
  target_tags   = ["http-server"]

  allow {
    protocol = "tcp"
    ports    = ["80", "443", "8080"]
  }

  depends_on = [
    google_compute_firewall.fw_iap
  ]
}

resource "google_network_services_lb_route_extension" "default" {
  name                  = "tf-test-l7-ilb-route-ext%{random_suffix}"
  description           = "my route extension"
  location              = "us-west1"
  load_balancing_scheme = "INTERNAL_MANAGED"
  forwarding_rules      = [google_compute_forwarding_rule.default.self_link]

  extension_chains {
    name = "chain1"

    match_condition {
      cel_expression = "request.path.startsWith('/extensions')"
    }

    extensions {
      name      = "ext11"
      authority = "ext11.com"
      service   = google_compute_region_backend_service.callouts_backend.self_link
      timeout   = "0.1s"
      fail_open = false

      forward_headers  = ["custom-header"]
    }
  }

  labels = {
    foo = "bar"
  }
}

# Route Extension Backend Instance
resource "google_compute_instance" "callouts_instance" {
  name         = "tf-test-l7-ilb-callouts-ins%{random_suffix}"
  zone         = "us-west1-a"
  machine_type = "e2-small"

  labels = {
    "container-vm" = "cos-stable-109-17800-147-54"
  }

  tags = ["allow-ssh","load-balanced-backend"]

  network_interface {
    network    = google_compute_network.ilb_network.id
    subnetwork = google_compute_subnetwork.ilb_subnet.id

    access_config {
      # add external ip to fetch packages
    }
  }

  boot_disk {
    auto_delete  = true

    initialize_params {
      type  = "pd-standard"
      size  = 10
      image = "https://www.googleapis.com/compute/v1/projects/cos-cloud/global/images/cos-stable-109-17800-147-54"
    }
  }

  # Initialize an Envoy's Ext Proc gRPC API based on a docker container
  metadata = {
    startup-script = <<-EOF1
      #! /bin/bash
      apt-get update
      apt-get install apache2 -y
      a2ensite default-ssl
      a2enmod ssl
      echo "Page served from second backend service" | tee /var/www/html/index.html
      systemctl restart apache2'
    EOF1
  }

  lifecycle {
    create_before_destroy = true
  }

  deletion_protection = false
}

// callouts instance group
resource "google_compute_instance_group" "callouts_instance_group" {
  name        = "tf-test-l7-ilb-callouts-ins-group%{random_suffix}"
  description = "Terraform test instance group"
  zone        = "us-west1-a"

  instances = [
    google_compute_instance.callouts_instance.id,
  ]

  named_port {
    name = "http"
    port = "80"
  }

  named_port {
    name = "grpc"
    port = "443"
  }
}

# callout health check
resource "google_compute_region_health_check" "callouts_health_check" {
  name     = "tf-test-l7-ilb-callouts-hc%{random_suffix}"
  region   = "us-west1"

  http_health_check {
    port = 80
  }

  depends_on = [
    google_compute_region_health_check.default
  ]
}

# callout backend service
resource "google_compute_region_backend_service" "callouts_backend" {
  name                  = "tf-test-l7-ilb-callouts-backend%{random_suffix}"
  region                = "us-west1"
  protocol              = "HTTP2"
  load_balancing_scheme = "INTERNAL_MANAGED"
  timeout_sec           = 10
  port_name             = "grpc"
  health_checks         = [google_compute_region_health_check.callouts_health_check.id]

  backend {
    group           = google_compute_instance_group.callouts_instance_group.id
    balancing_mode  = "UTILIZATION"
    capacity_scaler = 1.0
  }

  depends_on = [
    google_compute_region_backend_service.default
  ]
}
`, context)
}

func testAccNetworkServicesLbRouteExtension_update(context map[string]interface{}) string {
	return acctest.Nprintf(`
# Internal HTTP load balancer with a managed instance group backend
# VPC network
resource "google_compute_network" "ilb_network" {
  name                    = "tf-test-l7-ilb-network%{random_suffix}"
  auto_create_subnetworks = false
}

# proxy-only subnet
resource "google_compute_subnetwork" "proxy_subnet" {
  name          = "tf-test-l7-ilb-proxy-subnet%{random_suffix}"
  ip_cidr_range = "10.0.0.0/24"
  region        = "us-west1"
  purpose       = "REGIONAL_MANAGED_PROXY"
  role          = "ACTIVE"
  network       = google_compute_network.ilb_network.id
}

# backend subnet
resource "google_compute_subnetwork" "ilb_subnet" {
  name          = "tf-test-l7-ilb-subnet%{random_suffix}"
  ip_cidr_range = "10.0.1.0/24"
  region        = "us-west1"
  network       = google_compute_network.ilb_network.id

  depends_on = [
    google_compute_subnetwork.proxy_subnet
  ]
}

# forwarding rule
resource "google_compute_forwarding_rule" "default" {
  name                  = "tf-test-l7-ilb-forwarding-rule%{random_suffix}"
  region                = "us-west1"
  ip_protocol           = "TCP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  port_range            = "80"
  target                = google_compute_region_target_http_proxy.default.id
  network               = google_compute_network.ilb_network.id
  subnetwork            = google_compute_subnetwork.ilb_subnet.id
  network_tier          = "PREMIUM"

  depends_on = [
    google_compute_subnetwork.proxy_subnet
  ]
}

# Additional forwarding rule
resource "google_compute_forwarding_rule" "additional_forwarding_rule" {
  name                  = "tf-test-l7-ilb-additional-forwarding-rule%{random_suffix}"
  region                = "us-west1"
  ip_protocol           = "TCP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  port_range            = "80"
  target                = google_compute_region_target_http_proxy.default.id
  network               = google_compute_network.ilb_network.id
  subnetwork            = google_compute_subnetwork.ilb_subnet.id
  network_tier          = "PREMIUM"

  depends_on = [
    google_compute_subnetwork.proxy_subnet,
	google_compute_forwarding_rule.default
  ]
}

# HTTP target proxy
resource "google_compute_region_target_http_proxy" "default" {
  name     = "tf-test-l7-ilb-target-http-proxy%{random_suffix}"
  region   = "us-west1"
  url_map  = google_compute_region_url_map.default.id
}

# URL map
resource "google_compute_region_url_map" "default" {
  name            = "tf-test-l7-ilb-regional-url-map%{random_suffix}"
  region          = "us-west1"
  default_service = google_compute_region_backend_service.default.id

  host_rule {
    hosts        = ["service-extensions.com"]
    path_matcher = "callouts"
  }

  path_matcher {
    name            = "callouts"
    default_service = google_compute_region_backend_service.callouts_backend.id
  }

  host_rule {
    hosts        = ["service-extensions-2.com"]
    path_matcher = "callouts2"
  }

  path_matcher {
    name            = "callouts2"
    default_service = google_compute_region_backend_service.callouts_backend_2.id
  }
}

# backend service
resource "google_compute_region_backend_service" "default" {
  name                  = "tf-test-l7-ilb-backend-subnet%{random_suffix}"
  region                = "us-west1"
  protocol              = "HTTP"
  load_balancing_scheme = "INTERNAL_MANAGED"
  timeout_sec           = 10
  health_checks         = [google_compute_region_health_check.default.id]

  backend {
    group           = google_compute_region_instance_group_manager.mig.instance_group
    balancing_mode  = "UTILIZATION"
    capacity_scaler = 1.0
  }
}

# instance template
resource "google_compute_instance_template" "instance_template" {
  name         = "tf-test-l7-ilb-mig-template%{random_suffix}"
  machine_type = "e2-small"
  tags         = ["http-server"]

  network_interface {
    network    = google_compute_network.ilb_network.id
    subnetwork = google_compute_subnetwork.ilb_subnet.id

    access_config {
      # add external ip to fetch packages
    }
  }

  disk {
    source_image = "debian-cloud/debian-12"
    auto_delete  = true
    boot         = true
  }

  # install nginx and serve a simple web page
  metadata = {
    startup-script = <<-EOF1
      #! /bin/bash
      set -euo pipefail

      export DEBIAN_FRONTEND=noninteractive
      apt-get update
      apt-get install -y nginx-light jq

      NAME=$(curl -H "Metadata-Flavor: Google" "http://metadata.google.internal/computeMetadata/v1/instance/hostname")
      IP=$(curl -H "Metadata-Flavor: Google" "http://metadata.google.internal/computeMetadata/v1/instance/network-interfaces/0/ip")
      METADATA=$(curl -f -H "Metadata-Flavor: Google" "http://metadata.google.internal/computeMetadata/v1/instance/attributes/?recursive=True" | jq 'del(.["startup-script"])')

      cat <<EOF > /var/www/html/index.html
      <pre>
      Name: $NAME
      IP: $IP
      Metadata: $METADATA
      </pre>
      EOF
    EOF1
  }

  lifecycle {
    create_before_destroy = true
  }
}

# health check
resource "google_compute_region_health_check" "default" {
  name     = "tf-test-l7-ilb-hc%{random_suffix}"
  region   = "us-west1"

  http_health_check {
   port_specification = "USE_SERVING_PORT"
  }
}

# MIG
resource "google_compute_region_instance_group_manager" "mig" {
  name     = "tf-test-l7-ilb-mig1%{random_suffix}"
  region   = "us-west1"

  base_instance_name = "vm"
  target_size        = 2

  version {
    instance_template = google_compute_instance_template.instance_template.id
    name              = "primary"
  }
}

# allow all access from IAP and health check ranges
resource "google_compute_firewall" "fw_iap" {
  name          = "tf-test-l7-ilb-fw-allow-iap-hc%{random_suffix}"
  direction     = "INGRESS"
  network       = google_compute_network.ilb_network.id
  source_ranges = ["130.211.0.0/22", "35.191.0.0/16", "35.235.240.0/20"]

  allow {
    protocol = "tcp"
  }
}

# allow http from proxy subnet to backends
resource "google_compute_firewall" "fw_ilb_to_backends" {
  name          = "tf-test-l7-ilb-fw-allow-ilb-to-backends%{random_suffix}"
  direction     = "INGRESS"
  network       = google_compute_network.ilb_network.id
  source_ranges = ["10.0.0.0/24"]
  target_tags   = ["http-server"]

  allow {
    protocol = "tcp"
    ports    = ["80", "443", "8080"]
  }

  depends_on = [
    google_compute_firewall.fw_iap
  ]
}

resource "google_network_services_lb_route_extension" "default" {
  name                  = "tf-test-l7-ilb-route-ext%{random_suffix}"
  description           = "my route extension"
  location              = "us-west1"
  load_balancing_scheme = "INTERNAL_MANAGED"
  forwarding_rules      = [
    google_compute_forwarding_rule.default.self_link,
    google_compute_forwarding_rule.additional_forwarding_rule.self_link
  ]

  extension_chains {
    name = "chain1"

    match_condition {
      cel_expression = "request.path.startsWith('/extensions')"
    }

    extensions {
      name      = "ext12"
      authority = "ext12.com"
      service   = google_compute_region_backend_service.callouts_backend_2.self_link
      timeout   = "0.2s"
      fail_open = false

      forward_headers  = ["custom-header"]
    }
  }

  extension_chains {
    name = "chain2"

    match_condition {
      cel_expression = "request.path.startsWith('/extensions2')"
    }

    extensions {
      name      = "ext11"
      authority = "ext11.com"
      service   = google_compute_region_backend_service.callouts_backend.self_link
      timeout   = "0.1s"
      fail_open = false

      forward_headers  = ["custom-header"]
    }
  }

  labels = {
    bar = "foo"
  }
}

# Route Extension Backend Instance
resource "google_compute_instance" "callouts_instance" {
  name         = "tf-test-l7-ilb-callouts-ins%{random_suffix}"
  zone         = "us-west1-a"
  machine_type = "e2-small"

  labels = {
    "container-vm" = "cos-stable-109-17800-147-54"
  }

  tags = ["allow-ssh","load-balanced-backend"]

  network_interface {
    network    = google_compute_network.ilb_network.id
    subnetwork = google_compute_subnetwork.ilb_subnet.id

    access_config {
      # add external ip to fetch packages
    }
  }

  boot_disk {
    auto_delete  = true

    initialize_params {
      type  = "pd-standard"
      size  = 10
      image = "https://www.googleapis.com/compute/v1/projects/cos-cloud/global/images/cos-stable-109-17800-147-54"
    }
  }

  # Initialize an Envoy's Ext Proc gRPC API based on a docker container
  metadata = {
    startup-script = <<-EOF1
      #! /bin/bash
      apt-get update
      apt-get install apache2 -y
      a2ensite default-ssl
      a2enmod ssl
      echo "Page served from second backend service" | tee /var/www/html/index.html
      systemctl restart apache2'
    EOF1
  }

  lifecycle {
    create_before_destroy = true
  }

  deletion_protection = false
}

// callouts instance group
resource "google_compute_instance_group" "callouts_instance_group" {
  name        = "tf-test-l7-ilb-callouts-ins-group%{random_suffix}"
  description = "Terraform test instance group"
  zone        = "us-west1-a"

  instances = [
    google_compute_instance.callouts_instance.id,
  ]

  named_port {
    name = "http"
    port = "80"
  }

  named_port {
    name = "grpc"
    port = "443"
  }
}

# callout health check
resource "google_compute_region_health_check" "callouts_health_check" {
  name     = "tf-test-l7-ilb-callouts-hc%{random_suffix}"
  region   = "us-west1"

  http_health_check {
    port = 80
  }

  depends_on = [
    google_compute_region_health_check.default
  ]
}

# callout backend service
resource "google_compute_region_backend_service" "callouts_backend" {
  name                  = "tf-test-l7-ilb-callouts-backend%{random_suffix}"
  region                = "us-west1"
  protocol              = "HTTP2"
  load_balancing_scheme = "INTERNAL_MANAGED"
  timeout_sec           = 10
  port_name             = "grpc"
  health_checks         = [google_compute_region_health_check.callouts_health_check.id]

  backend {
    group           = google_compute_instance_group.callouts_instance_group.id
    balancing_mode  = "UTILIZATION"
    capacity_scaler = 1.0
  }

  depends_on = [
    google_compute_region_backend_service.default
  ]
}

# route extension backend instance 2
resource "google_compute_instance" "callouts_instance_2" {
  name         = "tf-test-l7-ilb-callouts-ins-2%{random_suffix}"
  zone         = "us-west1-a"
  machine_type = "e2-small"

  labels = {
    "container-vm" = "cos-stable-109-17800-147-54"
  }

  tags = ["allow-ssh","load-balanced-backend"]

  network_interface {
    network    = google_compute_network.ilb_network.id
    subnetwork = google_compute_subnetwork.ilb_subnet.id

    access_config {
      # add external ip to fetch packages
    }
  }

  boot_disk {
    auto_delete  = true

    initialize_params {
      type  = "pd-standard"
      size  = 10
      image = "https://www.googleapis.com/compute/v1/projects/cos-cloud/global/images/cos-stable-109-17800-147-54"
    }
  }

  # Initialize an Envoy's Ext Proc gRPC API based on a docker container
  metadata = {
    startup-script = <<-EOF1
      #! /bin/bash
      apt-get update
      apt-get install apache2 -y
      a2ensite default-ssl
      a2enmod ssl
      echo "Page served from second backend service" | tee /var/www/html/index.html
      systemctl restart apache2'
    EOF1
  }

  lifecycle {
    create_before_destroy = true
  }

  deletion_protection = false

  depends_on = [
    google_compute_instance.callouts_instance
  ]
}

// callouts instance group 2
resource "google_compute_instance_group" "callouts_instance_group_2" {
  name        = "tf-test-l7-ilb-callouts-ins-group-2%{random_suffix}"
  description = "Terraform test instance group"
  zone        = "us-west1-a"

  instances = [
    google_compute_instance.callouts_instance_2.id,
  ]

  named_port {
    name = "http"
    port = "80"
  }

  named_port {
    name = "grpc"
    port = "443"
  }

  depends_on = [
    google_compute_instance_group.callouts_instance_group
  ]
}

# callout health check 2
resource "google_compute_region_health_check" "callouts_health_check_2" {
  name     = "tf-test-l7-ilb-callouts-hc-2%{random_suffix}"
  region   = "us-west1"

  http_health_check {
    port = 80
  }

  depends_on = [
    google_compute_region_health_check.callouts_health_check
  ]
}

# callout backend service
resource "google_compute_region_backend_service" "callouts_backend_2" {
  name                  = "tf-test-l7-ilb-callouts-backend-2%{random_suffix}"
  region                = "us-west1"
  protocol              = "HTTP2"
  load_balancing_scheme = "INTERNAL_MANAGED"
  timeout_sec           = 10
  port_name             = "grpc"
  health_checks         = [google_compute_region_health_check.callouts_health_check_2.id]

  backend {
    group           = google_compute_instance_group.callouts_instance_group_2.id
    balancing_mode  = "UTILIZATION"
    capacity_scaler = 1.0
  }

  depends_on = [
    google_compute_region_backend_service.callouts_backend
  ]
}
`, context)
}
