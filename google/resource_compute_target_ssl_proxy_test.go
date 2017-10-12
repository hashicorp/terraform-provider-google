package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"reflect"
)

func TestAccComputeTargetSslProxy_basic(t *testing.T) {
	target := fmt.Sprintf("tssl-test-%s", acctest.RandString(10))
	cert := fmt.Sprintf("tssl-test-%s", acctest.RandString(10))
	backend := fmt.Sprintf("tssl-test-%s", acctest.RandString(10))
	hc := fmt.Sprintf("tssl-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeTargetSslProxyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeTargetSslProxy_basic1(target, cert, backend, hc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetSslProxy(
						"google_compute_target_ssl_proxy.foobar", "NONE", []string{cert}),
				),
			},
		},
	})
}

func TestAccComputeTargetSslProxy_update(t *testing.T) {
	target := fmt.Sprintf("tssl-test-%s", acctest.RandString(10))
	cert1 := fmt.Sprintf("tssl-test-%s", acctest.RandString(10))
	cert2 := fmt.Sprintf("tssl-test-%s", acctest.RandString(10))
	backend1 := fmt.Sprintf("tssl-test-%s", acctest.RandString(10))
	backend2 := fmt.Sprintf("tssl-test-%s", acctest.RandString(10))
	hc := fmt.Sprintf("tssl-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeTargetSslProxyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeTargetSslProxy_basic1(target, cert1, backend1, hc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetSslProxy(
						"google_compute_target_ssl_proxy.foobar", "NONE", []string{cert1}),
				),
			},
			resource.TestStep{
				Config: testAccComputeTargetSslProxy_basic2(target, cert1, cert2, backend1, backend2, hc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetSslProxy(
						"google_compute_target_ssl_proxy.foobar", "PROXY_V1", []string{cert1, cert2}),
				),
			},
		},
	})
}

func testAccCheckComputeTargetSslProxyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_target_ssl_proxy" {
			continue
		}

		_, err := config.clientCompute.TargetSslProxies.Get(
			config.Project, rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("TargetSslProxy still exists")
		}
	}

	return nil
}

func testAccCheckComputeTargetSslProxy(n, proxyHeader string, sslCerts []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.TargetSslProxies.Get(
			config.Project, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("TargetSslProxy not found")
		}

		if found.ProxyHeader != proxyHeader {
			return fmt.Errorf("Wrong proxy header. Expected '%s', got '%s'", proxyHeader, found.ProxyHeader)
		}

		foundCertsName := make([]string, 0, len(found.SslCertificates))
		for _, foundCert := range found.SslCertificates {
			foundCertsName = append(foundCertsName, GetResourceNameFromSelfLink(foundCert))
		}

		if !reflect.DeepEqual(foundCertsName, sslCerts) {
			return fmt.Errorf("Wrong ssl certificates. Expected '%s', got '%s'", sslCerts, found.SslCertificates)
		}

		return nil
	}
}

func testAccComputeTargetSslProxy_basic1(target, sslCert, backend, hc string) string {
	return fmt.Sprintf(`
resource "google_compute_target_ssl_proxy" "foobar" {
	description = "Resource created for Terraform acceptance testing"
	name = "%s"
	backend_service = "${google_compute_backend_service.foo.self_link}"
	ssl_certificates = ["${google_compute_ssl_certificate.foo.self_link}"]
	proxy_header = "NONE"
}

resource "google_compute_ssl_certificate" "foo" {
	name = "%s"
	private_key = "${file("test-fixtures/ssl_cert/test.key")}"
	certificate = "${file("test-fixtures/ssl_cert/test.crt")}"
}

resource "google_compute_backend_service" "foo" {
	name = "%s"
	protocol    = "SSL"
	health_checks = ["${google_compute_health_check.zero.self_link}"]
}

resource "google_compute_health_check" "zero" {
	name = "%s"
	check_interval_sec = 1
	timeout_sec = 1
	tcp_health_check {
		port = "443"
	}
}
`, target, sslCert, backend, hc)
}

func testAccComputeTargetSslProxy_basic2(target, sslCert1, sslCert2, backend1, backend2, hc string) string {
	return fmt.Sprintf(`
resource "google_compute_target_ssl_proxy" "foobar" {
	description = "Resource created for Terraform acceptance testing"
	name = "%s"
	backend_service = "${google_compute_backend_service.bar.self_link}"
	ssl_certificates = [
		"${google_compute_ssl_certificate.foo.self_link}",
		"${google_compute_ssl_certificate.bar.name}",
	]
	proxy_header = "PROXY_V1"
}

resource "google_compute_ssl_certificate" "foo" {
	name = "%s"
	private_key = "${file("test-fixtures/ssl_cert/test.key")}"
	certificate = "${file("test-fixtures/ssl_cert/test.crt")}"
}

resource "google_compute_ssl_certificate" "bar" {
	name = "%s"
	private_key = "${file("test-fixtures/ssl_cert/test.key")}"
	certificate = "${file("test-fixtures/ssl_cert/test.crt")}"
}

resource "google_compute_backend_service" "foo" {
	name = "%s"
	protocol    = "SSL"
	health_checks = ["${google_compute_health_check.zero.self_link}"]
}

resource "google_compute_backend_service" "bar" {
	name = "%s"
	protocol    = "SSL"
	health_checks = ["${google_compute_health_check.zero.self_link}"]
}

resource "google_compute_health_check" "zero" {
	name = "%s"
	check_interval_sec = 1
	timeout_sec = 1
	tcp_health_check {
		port = "443"
	}
}
`, target, sslCert1, sslCert2, backend1, backend2, hc)
}
