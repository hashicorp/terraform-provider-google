package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"google.golang.org/api/compute/v1"
)

func TestAccComputeTargetSslProxy_update(t *testing.T) {
	target := fmt.Sprintf("tssl-test-%s", RandString(t, 10))
	sslPolicy := fmt.Sprintf("tssl-test-%s", RandString(t, 10))
	cert1 := fmt.Sprintf("tssl-test-%s", RandString(t, 10))
	cert2 := fmt.Sprintf("tssl-test-%s", RandString(t, 10))
	backend1 := fmt.Sprintf("tssl-test-%s", RandString(t, 10))
	backend2 := fmt.Sprintf("tssl-test-%s", RandString(t, 10))
	hc := fmt.Sprintf("tssl-test-%s", RandString(t, 10))

	resourceSuffix := RandString(t, 10)
	var proxy compute.TargetSslProxy

	VcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    TestAccProviders,
		CheckDestroy: testAccCheckComputeTargetSslProxyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeTargetSslProxy_basic1(target, sslPolicy, cert1, backend1, hc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetSslProxyExists(
						t, "google_compute_target_ssl_proxy.foobar", &proxy),
					testAccCheckComputeTargetSslProxyHeader(t, "NONE", &proxy),
					testAccCheckComputeTargetSslProxyHasSslCertificate(t, cert1, &proxy),
				),
			},
			{
				Config: testAccComputeTargetSslProxy_basic2(target, sslPolicy, cert1, cert2, backend1, backend2, hc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetSslProxyExists(
						t, "google_compute_target_ssl_proxy.foobar", &proxy),
					testAccCheckComputeTargetSslProxyHeader(t, "PROXY_V1", &proxy),
					testAccCheckComputeTargetSslProxyHasSslCertificate(t, cert2, &proxy),
				),
			},
			{
				Config: testAccComputeTargetSslProxy_certificateMap1(resourceSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetSslProxyExists(
						t, "google_compute_target_ssl_proxy.with_certificate_map", &proxy),
					testAccCheckComputeTargetSslProxyHasCertificateMap(t, "certificatemap-test-1-"+resourceSuffix, &proxy),
				),
			},
			{
				Config: testAccComputeTargetSslProxy_certificateMap2(resourceSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetSslProxyExists(
						t, "google_compute_target_ssl_proxy.with_certificate_map", &proxy),
					testAccCheckComputeTargetSslProxyHasCertificateMap(t, "certificatemap-test-2-"+resourceSuffix, &proxy),
				),
			},
		},
	})
}

func testAccCheckComputeTargetSslProxyExists(t *testing.T, n string, proxy *compute.TargetSslProxy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := GoogleProviderConfig(t)
		name := rs.Primary.Attributes["name"]

		found, err := config.NewComputeClient(config.UserAgent).TargetSslProxies.Get(
			config.Project, name).Do()
		if err != nil {
			return err
		}

		if found.Name != name {
			return fmt.Errorf("TargetSslProxy not found")
		}

		*proxy = *found

		return nil
	}
}

func testAccCheckComputeTargetSslProxyHeader(t *testing.T, proxyHeader string, proxy *compute.TargetSslProxy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if proxy.ProxyHeader != proxyHeader {
			return fmt.Errorf("Wrong proxy header. Expected '%s', got '%s'", proxyHeader, proxy.ProxyHeader)
		}
		return nil
	}
}

func testAccCheckComputeTargetSslProxyHasSslCertificate(t *testing.T, cert string, proxy *compute.TargetSslProxy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := GoogleProviderConfig(t)
		certURL := fmt.Sprintf(canonicalSslCertificateTemplate, config.Project, cert)

		for _, sslCertificate := range proxy.SslCertificates {
			if ConvertSelfLinkToV1(sslCertificate) == certURL {
				return nil
			}
		}

		return fmt.Errorf("Ssl certificate not found: expected'%s'", certURL)
	}
}

func testAccCheckComputeTargetSslProxyHasCertificateMap(t *testing.T, certificateMap string, proxy *compute.TargetSslProxy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := GoogleProviderConfig(t)
		wantCertMapURL := fmt.Sprintf(canonicalCertificateMapTemplate, config.Project, certificateMap)
		gotCertMapURL := ConvertSelfLinkToV1(proxy.CertificateMap)
		if wantCertMapURL != gotCertMapURL {
			return fmt.Errorf("certificate map not found: got %q, want %q", gotCertMapURL, wantCertMapURL)
		}
		return nil
	}
}

func testAccComputeTargetSslProxy_basic1(target, sslPolicy, sslCert, backend, hc string) string {
	return fmt.Sprintf(`
resource "google_compute_target_ssl_proxy" "foobar" {
  description      = "Resource created for Terraform acceptance testing"
  name             = "%s"
  backend_service  = google_compute_backend_service.foo.self_link
  ssl_certificates = [google_compute_ssl_certificate.foo.self_link]
  proxy_header     = "NONE"
  ssl_policy       = google_compute_ssl_policy.foo.self_link
}

resource "google_compute_ssl_policy" "foo" {
  name            = "%s"
  description     = "Resource created for Terraform acceptance testing"
  min_tls_version = "TLS_1_2"
  profile         = "MODERN"
}

resource "google_compute_ssl_certificate" "foo" {
  name        = "%s"
  private_key = file("test-fixtures/ssl_cert/test.key")
  certificate = file("test-fixtures/ssl_cert/test.crt")
}

resource "google_compute_backend_service" "foo" {
  name          = "%s"
  protocol      = "SSL"
  health_checks = [google_compute_health_check.zero.self_link]
}

resource "google_compute_health_check" "zero" {
  name               = "%s"
  check_interval_sec = 1
  timeout_sec        = 1
  tcp_health_check {
    port = "443"
  }
}
`, target, sslPolicy, sslCert, backend, hc)
}

func testAccComputeTargetSslProxy_basic2(target, sslPolicy, sslCert1, sslCert2, backend1, backend2, hc string) string {
	return fmt.Sprintf(`
resource "google_compute_target_ssl_proxy" "foobar" {
  description      = "Resource created for Terraform acceptance testing"
  name             = "%s"
  backend_service  = google_compute_backend_service.bar.self_link
  ssl_certificates = [google_compute_ssl_certificate.bar.name]
  proxy_header     = "PROXY_V1"
}

resource "google_compute_ssl_policy" "foo" {
  name            = "%s"
  description     = "Resource created for Terraform acceptance testing"
  min_tls_version = "TLS_1_2"
  profile         = "MODERN"
}

resource "google_compute_ssl_certificate" "foo" {
  name        = "%s"
  private_key = file("test-fixtures/ssl_cert/test.key")
  certificate = file("test-fixtures/ssl_cert/test.crt")
}

resource "google_compute_ssl_certificate" "bar" {
  name        = "%s"
  private_key = file("test-fixtures/ssl_cert/test.key")
  certificate = file("test-fixtures/ssl_cert/test.crt")
}

resource "google_compute_backend_service" "foo" {
  name          = "%s"
  protocol      = "SSL"
  health_checks = [google_compute_health_check.zero.self_link]
}

resource "google_compute_backend_service" "bar" {
  name          = "%s"
  protocol      = "SSL"
  health_checks = [google_compute_health_check.zero.self_link]
}

resource "google_compute_health_check" "zero" {
  name               = "%s"
  check_interval_sec = 1
  timeout_sec        = 1
  tcp_health_check {
    port = "443"
  }
}
`, target, sslPolicy, sslCert1, sslCert2, backend1, backend2, hc)
}

func testAccComputeTargetSslProxy_certificateMap1(id string) string {
	return fmt.Sprintf(`
resource "google_compute_target_ssl_proxy" "with_certificate_map" {
  description      = "Resource created for Terraform acceptance testing"
  name             = "ssl-proxy-%s"
  backend_service  = google_compute_backend_service.foo.self_link
	certificate_map = "//certificatemanager.googleapis.com/${google_certificate_manager_certificate_map.map1.id}"
}

resource "google_compute_backend_service" "foo" {
  name          = "backend-service-%s"
  protocol      = "SSL"
  health_checks = [google_compute_health_check.zero.self_link]
}

resource "google_compute_health_check" "zero" {
  name               = "health-check-%s"
  check_interval_sec = 1
  timeout_sec        = 1
  tcp_health_check {
    port = "443"
  }
}

resource "google_certificate_manager_certificate_map" "map1" {
  name = "certificatemap-test-1-%s"
}
resource "google_certificate_manager_certificate_map_entry" "map_entry" {
  name         = "certificatemapentry-test-%s"
  map          = google_certificate_manager_certificate_map.map1.name
  certificates = [google_certificate_manager_certificate.certificate.id]
  matcher      = "PRIMARY"
}

resource "google_certificate_manager_certificate" "certificate" {
  name        = "certificate-test-%s"
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
  name   = "dnsauthorization-test-%s"
  domain = "mysite.com"
}
`, id, id, id, id, id, id, id)
}

func testAccComputeTargetSslProxy_certificateMap2(id string) string {
	return fmt.Sprintf(`
resource "google_compute_target_ssl_proxy" "with_certificate_map" {
  description      = "Resource created for Terraform acceptance testing"
  name             = "ssl-proxy-%s"
  backend_service  = google_compute_backend_service.foo.self_link
	certificate_map = "//certificatemanager.googleapis.com/${google_certificate_manager_certificate_map.map2.id}"
}

resource "google_compute_backend_service" "foo" {
  name          = "backend-service-%s"
  protocol      = "SSL"
  health_checks = [google_compute_health_check.zero.self_link]
}

resource "google_compute_health_check" "zero" {
  name               = "health-check-%s"
  check_interval_sec = 1
  timeout_sec        = 1
  tcp_health_check {
    port = "443"
  }
}

resource "google_certificate_manager_certificate_map" "map1" {
  name = "certificatemap-test-1-%s"
}

resource "google_certificate_manager_certificate_map" "map2" {
  name = "certificatemap-test-2-%s"
}

resource "google_certificate_manager_certificate_map_entry" "map_entry" {
  name         = "certificatemapentry-test-%s"
  map          = google_certificate_manager_certificate_map.map1.name
  certificates = [google_certificate_manager_certificate.certificate.id]
  matcher      = "PRIMARY"
}

resource "google_certificate_manager_certificate" "certificate" {
  name        = "certificate-test-%s"
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
  name   = "dnsauthorization-test-%s"
  domain = "mysite.com"
}
`, id, id, id, id, id, id, id, id)
}
