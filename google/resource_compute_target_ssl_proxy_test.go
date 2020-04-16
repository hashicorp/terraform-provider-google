package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccComputeTargetSslProxy_update(t *testing.T) {
	target := fmt.Sprintf("tssl-test-%s", randString(t, 10))
	sslPolicy := fmt.Sprintf("tssl-test-%s", randString(t, 10))
	cert1 := fmt.Sprintf("tssl-test-%s", randString(t, 10))
	cert2 := fmt.Sprintf("tssl-test-%s", randString(t, 10))
	backend1 := fmt.Sprintf("tssl-test-%s", randString(t, 10))
	backend2 := fmt.Sprintf("tssl-test-%s", randString(t, 10))
	hc := fmt.Sprintf("tssl-test-%s", randString(t, 10))

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeTargetSslProxyDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeTargetSslProxy_basic1(target, sslPolicy, cert1, backend1, hc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetSslProxy(
						t, "google_compute_target_ssl_proxy.foobar", "NONE", cert1),
				),
			},
			{
				Config: testAccComputeTargetSslProxy_basic2(target, sslPolicy, cert1, cert2, backend1, backend2, hc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetSslProxy(
						t, "google_compute_target_ssl_proxy.foobar", "PROXY_V1", cert2),
				),
			},
		},
	})
}

func testAccCheckComputeTargetSslProxy(t *testing.T, n, proxyHeader, sslCert string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := googleProviderConfig(t)
		name := rs.Primary.Attributes["name"]

		found, err := config.clientCompute.TargetSslProxies.Get(
			config.Project, name).Do()
		if err != nil {
			return err
		}

		if found.Name != name {
			return fmt.Errorf("TargetSslProxy not found")
		}

		if found.ProxyHeader != proxyHeader {
			return fmt.Errorf("Wrong proxy header. Expected '%s', got '%s'", proxyHeader, found.ProxyHeader)
		}

		foundCertName := GetResourceNameFromSelfLink(found.SslCertificates[0])
		if foundCertName != sslCert {
			return fmt.Errorf("Wrong ssl certificates. Expected '%s', got '%s'", sslCert, foundCertName)
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
