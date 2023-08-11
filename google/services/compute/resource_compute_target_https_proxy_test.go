// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/api/compute/v1"
)

const (
	canonicalSslCertificateTemplate = "https://www.googleapis.com/compute/v1/projects/%s/global/sslCertificates/%s"
	canonicalCertificateMapTemplate = "//certificatemanager.googleapis.com/projects/%s/locations/global/certificateMaps/%s"
)

func TestAccComputeTargetHttpsProxy_update(t *testing.T) {
	t.Parallel()

	var proxy compute.TargetHttpsProxy
	resourceSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeTargetHttpsProxyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeTargetHttpsProxy_basic1(resourceSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetHttpsProxyExists(
						t, "google_compute_target_https_proxy.foobar", &proxy),
					testAccComputeTargetHttpsProxyDescription("Resource created for Terraform acceptance testing", &proxy),
					testAccComputeTargetHttpsProxyHasSslCertificate(t, "tf-test-httpsproxy-cert1-"+resourceSuffix, &proxy),
				),
			},

			{
				Config: testAccComputeTargetHttpsProxy_basic2(resourceSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetHttpsProxyExists(
						t, "google_compute_target_https_proxy.foobar", &proxy),
					testAccComputeTargetHttpsProxyDescription("Resource created for Terraform acceptance testing", &proxy),
					testAccComputeTargetHttpsProxyHasSslCertificate(t, "tf-test-httpsproxy-cert1-"+resourceSuffix, &proxy),
					testAccComputeTargetHttpsProxyHasSslCertificate(t, "tf-test-httpsproxy-cert2-"+resourceSuffix, &proxy),
				),
			},
		},
	})
}

func TestAccComputeTargetHttpsProxy_certificateMap(t *testing.T) {
	t.Parallel()

	var proxy compute.TargetHttpsProxy
	resourceSuffix := acctest.RandString(t, 10)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeTargetHttpsProxyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeTargetHttpsProxy_certificateMap(resourceSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetHttpsProxyExists(
						t, "google_compute_target_https_proxy.foobar", &proxy),
					testAccComputeTargetHttpsProxyDescription("Resource created for Terraform acceptance testing", &proxy),
					testAccComputeTargetHttpsProxyHasCertificateMap(t, "tf-test-certmap-"+resourceSuffix, &proxy),
				),
			},
		},
	})
}

func testAccCheckComputeTargetHttpsProxyExists(t *testing.T, n string, proxy *compute.TargetHttpsProxy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.GoogleProviderConfig(t)
		name := rs.Primary.Attributes["name"]

		found, err := config.NewComputeClient(config.UserAgent).TargetHttpsProxies.Get(
			config.Project, name).Do()
		if err != nil {
			return err
		}

		if found.Name != name {
			return fmt.Errorf("TargetHttpsProxy not found")
		}

		*proxy = *found

		return nil
	}
}

func testAccComputeTargetHttpsProxyDescription(description string, proxy *compute.TargetHttpsProxy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if proxy.Description != description {
			return fmt.Errorf("Wrong description: expected '%s' got '%s'", description, proxy.Description)
		}
		return nil
	}
}

func testAccComputeTargetHttpsProxyHasSslCertificate(t *testing.T, cert string, proxy *compute.TargetHttpsProxy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		certUrl := fmt.Sprintf(canonicalSslCertificateTemplate, config.Project, cert)

		for _, sslCertificate := range proxy.SslCertificates {
			if tpgresource.ConvertSelfLinkToV1(sslCertificate) == certUrl {
				return nil
			}
		}

		return fmt.Errorf("Ssl certificate not found: expected'%s'", certUrl)
	}
}

func testAccComputeTargetHttpsProxyHasCertificateMap(t *testing.T, certificateMap string, proxy *compute.TargetHttpsProxy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acctest.GoogleProviderConfig(t)
		certificateMapUrl := fmt.Sprintf(canonicalCertificateMapTemplate, config.Project, certificateMap)

		if tpgresource.ConvertSelfLinkToV1(proxy.CertificateMap) == certificateMapUrl {
			return nil
		}

		return fmt.Errorf("certificate map not found: expected'%s'", certificateMapUrl)
	}
}

func testAccComputeTargetHttpsProxy_basic1(id string) string {
	return fmt.Sprintf(`
resource "google_compute_target_https_proxy" "foobar" {
  description      = "Resource created for Terraform acceptance testing"
  name             = "tf-test-httpsproxy-%s"
  url_map          = google_compute_url_map.foobar.self_link
  ssl_certificates = [google_compute_ssl_certificate.foobar1.self_link]
  ssl_policy       = google_compute_ssl_policy.foobar.self_link
}

resource "google_compute_backend_service" "foobar" {
  name          = "tf-test-httpsproxy-backend-%s"
  health_checks = [google_compute_http_health_check.zero.self_link]
}

resource "google_compute_http_health_check" "zero" {
  name               = "tf-test-httpsproxy-check-%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_compute_url_map" "foobar" {
  name            = "tf-test-httpsproxy-urlmap-%s"
  default_service = google_compute_backend_service.foobar.self_link
  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_backend_service.foobar.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_backend_service.foobar.self_link
    }
  }
  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_backend_service.foobar.self_link
  }
}

resource "google_compute_ssl_policy" "foobar" {
  name            = "tf-test-sslproxy-%s"
  description     = "my-description"
  min_tls_version = "TLS_1_2"
  profile         = "MODERN"
}

resource "google_compute_ssl_certificate" "foobar1" {
  name        = "tf-test-httpsproxy-cert1-%s"
  description = "very descriptive"
  private_key = file("test-fixtures/test.key")
  certificate = file("test-fixtures/test.crt")
}

resource "google_compute_ssl_certificate" "foobar2" {
  name        = "tf-test-httpsproxy-cert2-%s"
  description = "very descriptive"
  private_key = file("test-fixtures/test.key")
  certificate = file("test-fixtures/test.crt")
}
`, id, id, id, id, id, id, id)
}

func testAccComputeTargetHttpsProxy_basic2(id string) string {
	return fmt.Sprintf(`
resource "google_compute_target_https_proxy" "foobar" {
  description = "Resource created for Terraform acceptance testing"
  name        = "tf-test-httpsproxy-%s"
  url_map     = google_compute_url_map.foobar.self_link
  ssl_certificates = [
    google_compute_ssl_certificate.foobar1.self_link,
    google_compute_ssl_certificate.foobar2.self_link,
  ]
  quic_override = "ENABLE"
}

resource "google_compute_backend_service" "foobar" {
  name          = "tf-test-httpsproxy-backend-%s"
  health_checks = [google_compute_http_health_check.zero.self_link]
}

resource "google_compute_http_health_check" "zero" {
  name               = "tf-test-httpsproxy-check-%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_compute_url_map" "foobar" {
  name            = "tf-test-httpsproxy-urlmap-%s"
  default_service = google_compute_backend_service.foobar.self_link
  host_rule {
    hosts        = ["mysite.com", "myothersite.com"]
    path_matcher = "boop"
  }
  path_matcher {
    default_service = google_compute_backend_service.foobar.self_link
    name            = "boop"
    path_rule {
      paths   = ["/*"]
      service = google_compute_backend_service.foobar.self_link
    }
  }
  test {
    host    = "mysite.com"
    path    = "/*"
    service = google_compute_backend_service.foobar.self_link
  }
}

resource "google_compute_ssl_policy" "foobar" {
  name            = "tf-test-sslproxy-%s"
  description     = "my-description"
  min_tls_version = "TLS_1_2"
  profile         = "MODERN"
}

resource "google_compute_ssl_certificate" "foobar1" {
  name        = "tf-test-httpsproxy-cert1-%s"
  description = "very descriptive"
  private_key = file("test-fixtures/test.key")
  certificate = file("test-fixtures/test.crt")
}

resource "google_compute_ssl_certificate" "foobar2" {
  name        = "tf-test-httpsproxy-cert2-%s"
  description = "very descriptive"
  private_key = file("test-fixtures/test.key")
  certificate = file("test-fixtures/test.crt")
}
`, id, id, id, id, id, id, id)
}

func testAccComputeTargetHttpsProxy_certificateMap(id string) string {
	return fmt.Sprintf(`
resource "google_compute_target_https_proxy" "foobar" {
  description      = "Resource created for Terraform acceptance testing"
  name             = "tf-test-httpsproxy-%s"
  url_map          = google_compute_url_map.foobar.self_link
  certificate_map = "//certificatemanager.googleapis.com/${google_certificate_manager_certificate_map.map.id}"
}

resource "google_compute_backend_service" "foobar" {
  name          = "tf-test-httpsproxy-backend-%s"
  health_checks = [google_compute_http_health_check.zero.self_link]
}

resource "google_compute_http_health_check" "zero" {
  name               = "tf-test-httpsproxy-check-%s"
  request_path       = "/"
  check_interval_sec = 1
  timeout_sec        = 1
}

resource "google_compute_url_map" "foobar" {
  name            = "tf-test-httpsproxy-urlmap-%s"
  default_service = google_compute_backend_service.foobar.self_link
}

resource "google_certificate_manager_certificate_map" "map" {
  name = "tf-test-certmap-%s"
}

resource "google_certificate_manager_certificate_map_entry" "map_entry" {
  name         = "tf-test-certmapentry-%s"
  map          = google_certificate_manager_certificate_map.map.name
  certificates = [google_certificate_manager_certificate.certificate.id]
  matcher      = "PRIMARY"
}

resource "google_certificate_manager_certificate" "certificate" {
  name        = "tf-test-cert-%s"
  scope       = "DEFAULT"
  managed {
    domains = [
      google_certificate_manager_dns_authorization.instance.domain,
    ]
    dns_authorizations = [
      google_certificate_manager_dns_authorization.instance.id,
    ]
  }
}

resource "google_certificate_manager_dns_authorization" "instance" {
  name   = "tf-test-dnsauthz-%s"
  domain = "mysite.com"
}

`, id, id, id, id, id, id, id, id)
}
