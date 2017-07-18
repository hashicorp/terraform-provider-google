package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/compute/v1"
	"regexp"
)

func TestAccComputeTargetHttpsProxy_basic(t *testing.T) {
	var proxy compute.TargetHttpsProxy
	resourceSuffix := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeTargetHttpsProxyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeTargetHttpsProxy_basic1(resourceSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetHttpsProxyExists(
						"google_compute_target_https_proxy.foobar", &proxy),
					testAccComputeTargetHttpsProxyDescription("Resource created for Terraform acceptance testing", &proxy),
					testAccComputeTargetHttpsProxyHasSslCertificate("httpsproxy-test-cert1-"+resourceSuffix, &proxy),
				),
			},
		},
	})
}

func TestAccComputeTargetHttpsProxy_update(t *testing.T) {
	var proxy compute.TargetHttpsProxy
	resourceSuffix := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeTargetHttpsProxyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeTargetHttpsProxy_basic1(resourceSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetHttpsProxyExists(
						"google_compute_target_https_proxy.foobar", &proxy),
					testAccComputeTargetHttpsProxyDescription("Resource created for Terraform acceptance testing", &proxy),
					testAccComputeTargetHttpsProxyHasSslCertificate("httpsproxy-test-cert1-"+resourceSuffix, &proxy),
				),
			},

			resource.TestStep{
				Config: testAccComputeTargetHttpsProxy_basic2(resourceSuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeTargetHttpsProxyExists(
						"google_compute_target_https_proxy.foobar", &proxy),
					testAccComputeTargetHttpsProxyDescription("Resource created for Terraform acceptance testing (updated)", &proxy),
					testAccComputeTargetHttpsProxyHasSslCertificate("httpsproxy-test-cert1-"+resourceSuffix, &proxy),
					testAccComputeTargetHttpsProxyHasSslCertificate("httpsproxy-test-cert2-"+resourceSuffix, &proxy),
				),
			},
		},
	})
}

func TestAccComputeTargetHttpsProxy_invalidCertificate(t *testing.T) {
	resourceSuffix := acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeTargetHttpsProxyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testAccComputeTargetHttpsProxy_invalidCertificate(resourceSuffix),
				ExpectError: regexp.MustCompile("ssl_certificate"),
			},
		},
	})
}

func testAccCheckComputeTargetHttpsProxyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_compute_target_https_proxy" {
			continue
		}

		_, err := config.clientCompute.TargetHttpsProxies.Get(
			config.Project, rs.Primary.ID).Do()
		if err == nil {
			return fmt.Errorf("TargetHttpsProxy still exists")
		}
	}

	return nil
}

func testAccCheckComputeTargetHttpsProxyExists(n string, proxy *compute.TargetHttpsProxy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.TargetHttpsProxies.Get(
			config.Project, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
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

func testAccComputeTargetHttpsProxyHasSslCertificate(cert string, proxy *compute.TargetHttpsProxy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		certUrl := fmt.Sprintf(canonicalSslCertificateTemplate, config.Project, cert)

		for _, sslCertificate := range proxy.SslCertificates {
			if sslCertificate == certUrl {
				return nil
			}
		}

		return fmt.Errorf("Ssl certificate not found: expected'%s'", certUrl)
	}
}

func testAccComputeTargetHttpsProxy_basic1(id string) string {
	return fmt.Sprintf(`
resource "google_compute_target_https_proxy" "foobar" {
	description = "Resource created for Terraform acceptance testing"
	name = "httpsproxy-test-%s"
	url_map = "${google_compute_url_map.foobar.self_link}"
	ssl_certificates = ["${google_compute_ssl_certificate.foobar1.self_link}"]
}

resource "google_compute_backend_service" "foobar" {
	name = "httpsproxy-test-backend-%s"
	health_checks = ["${google_compute_http_health_check.zero.self_link}"]
}

resource "google_compute_http_health_check" "zero" {
	name = "httpsproxy-test-health-check-%s"
	request_path = "/"
	check_interval_sec = 1
	timeout_sec = 1
}

resource "google_compute_url_map" "foobar" {
	name = "httpsproxy-test-url-map-%s"
	default_service = "${google_compute_backend_service.foobar.self_link}"
	host_rule {
		hosts = ["mysite.com", "myothersite.com"]
		path_matcher = "boop"
	}
	path_matcher {
		default_service = "${google_compute_backend_service.foobar.self_link}"
		name = "boop"
		path_rule {
			paths = ["/*"]
			service = "${google_compute_backend_service.foobar.self_link}"
		}
	}
	test {
		host = "mysite.com"
		path = "/*"
		service = "${google_compute_backend_service.foobar.self_link}"
	}
}

resource "google_compute_ssl_certificate" "foobar1" {
	name = "httpsproxy-test-cert1-%s"
	description = "very descriptive"
	private_key = "${file("test-fixtures/ssl_cert/test.key")}"
	certificate = "${file("test-fixtures/ssl_cert/test.crt")}"
}

resource "google_compute_ssl_certificate" "foobar2" {
	name = "httpsproxy-test-cert2-%s"
	description = "very descriptive"
	private_key = "${file("test-fixtures/ssl_cert/test.key")}"
	certificate = "${file("test-fixtures/ssl_cert/test.crt")}"
}
`, id, id, id, id, id, id)
}

func testAccComputeTargetHttpsProxy_basic2(id string) string {
	return fmt.Sprintf(`
resource "google_compute_target_https_proxy" "foobar" {
	description = "Resource created for Terraform acceptance testing (updated)"
	name = "httpsproxy-test-%s"
	url_map = "${google_compute_url_map.foobar.self_link}"
	ssl_certificates = [
		"${google_compute_ssl_certificate.foobar1.self_link}",
		"${google_compute_ssl_certificate.foobar2.self_link}",
	]
}

resource "google_compute_backend_service" "foobar" {
	name = "httpsproxy-test-backend-%s"
	health_checks = ["${google_compute_http_health_check.zero.self_link}"]
}

resource "google_compute_http_health_check" "zero" {
	name = "httpsproxy-test-health-check-%s"
	request_path = "/"
	check_interval_sec = 1
	timeout_sec = 1
}

resource "google_compute_url_map" "foobar" {
	name = "httpsproxy-test-url-map-%s"
	default_service = "${google_compute_backend_service.foobar.self_link}"
	host_rule {
		hosts = ["mysite.com", "myothersite.com"]
		path_matcher = "boop"
	}
	path_matcher {
		default_service = "${google_compute_backend_service.foobar.self_link}"
		name = "boop"
		path_rule {
			paths = ["/*"]
			service = "${google_compute_backend_service.foobar.self_link}"
		}
	}
	test {
		host = "mysite.com"
		path = "/*"
		service = "${google_compute_backend_service.foobar.self_link}"
	}
}

resource "google_compute_ssl_certificate" "foobar1" {
	name = "httpsproxy-test-cert1-%s"
	description = "very descriptive"
	private_key = "${file("test-fixtures/ssl_cert/test.key")}"
	certificate = "${file("test-fixtures/ssl_cert/test.crt")}"
}

resource "google_compute_ssl_certificate" "foobar2" {
	name = "httpsproxy-test-cert2-%s"
	description = "very descriptive"
	private_key = "${file("test-fixtures/ssl_cert/test.key")}"
	certificate = "${file("test-fixtures/ssl_cert/test.crt")}"
}
`, id, id, id, id, id, id)
}

func testAccComputeTargetHttpsProxy_invalidCertificate(id string) string {
	return fmt.Sprintf(`
resource "google_compute_target_https_proxy" "foobar" {
name = "httpsproxy-test-%s"
url_map = "some-url-map"
ssl_certificates = ["invalid-certificate-reference"]
}
`, id)
}
