# Example for using Cloud Armor https://cloud.google.com/armor/
#

resource "random_id" "instance_id" {
  byte_length = 4
}

# Configure the Google Cloud provider
provider "google" {
  credentials = file(var.credentials_file_path)
  project     = var.project_name
  region      = var.region
  zone        = var.region_zone
}

# Set up a backend to be proxied to:
# A single instance in a pool running nginx with port 80 open will allow end to end network testing
resource "google_compute_instance" "cluster1" {
  name         = "armor-gce-${random_id.instance_id.hex}"
  machine_type = "f1-micro"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  network_interface {
    network = "default"
    access_config {
      # Ephemeral IP
    }
  }

  metadata_startup_script = "sudo apt-get update; sudo apt-get install -yq nginx; sudo service nginx restart"
}

resource "google_compute_firewall" "cluster1" {
  name    = "armor-firewall"
  network = "default"

  allow {
    protocol = "tcp"
    ports    = ["80", "43"]
  }
}

resource "google_compute_instance_group" "webservers" {
  name        = "instance-group-all"
  description = "An instance group for the single GCE instance"

  instances = [
    google_compute_instance.cluster1.self_link,
  ]

  named_port {
    name = "http"
    port = "80"
  }
}

resource "google_compute_target_pool" "example" {
  name = "armor-pool"

  instances = [
    google_compute_instance.cluster1.self_link,
  ]

  health_checks = [
    google_compute_http_health_check.health.name,
  ]
}

resource "google_compute_http_health_check" "health" {
  name               = "armor-healthcheck"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_compute_backend_service" "website" {
  name        = "armor-backend"
  description = "Our company website"
  port_name   = "http"
  protocol    = "HTTP"
  timeout_sec = 10
  enable_cdn  = false

  backend {
    group = google_compute_instance_group.webservers.self_link
  }

  security_policy = google_compute_security_policy.security-policy-1.self_link

  health_checks = [google_compute_http_health_check.health.self_link]
}

# Cloud Armor Security policies
resource "google_compute_security_policy" "security-policy-1" {
  name        = "armor-security-policy"
  description = "example security policy"

  # Reject all traffic that hasn't been whitelisted.
  rule {
    action   = "deny(403)"
    priority = "2147483647"

    match {
      versioned_expr = "SRC_IPS_V1"

      config {
        src_ip_ranges = ["*"]
      }
    }

    description = "Default rule, higher priority overrides it"
  }

  # Whitelist traffic from certain ip address
  rule {
    action   = "allow"
    priority = "1000"

    match {
      versioned_expr = "SRC_IPS_V1"

      config {
        src_ip_ranges = var.ip_white_list
      }
    }

    description = "allow traffic from 192.0.2.0/24"
  }
}

# Front end of the load balancer
resource "google_compute_global_forwarding_rule" "default" {
  name       = "armor-rule"
  target     = google_compute_target_http_proxy.default.self_link
  port_range = "80"
}

resource "google_compute_target_http_proxy" "default" {
  name    = "armor-proxy"
  url_map = google_compute_url_map.default.self_link
}

resource "google_compute_url_map" "default" {
  name            = "armor-url-map"
  default_service = google_compute_backend_service.website.self_link

  host_rule {
    hosts        = ["mysite.com"]
    path_matcher = "allpaths"
  }

  path_matcher {
    name            = "allpaths"
    default_service = google_compute_backend_service.website.self_link

    path_rule {
      paths   = ["/*"]
      service = google_compute_backend_service.website.self_link
    }
  }
}

output "ip" {
  value = google_compute_global_forwarding_rule.default.ip_address
}
