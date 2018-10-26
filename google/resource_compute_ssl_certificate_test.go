package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccComputeSslCertificate_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSslCertificateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSslCertificate_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSslCertificateExists(
						"google_compute_ssl_certificate.foobar"),
				),
			},
			resource.TestStep{
				ResourceName:            "google_compute_ssl_certificate.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key"},
			},
		},
	})
}

func TestAccComputeSslCertificate_no_name(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSslCertificateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSslCertificate_no_name(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSslCertificateExists(
						"google_compute_ssl_certificate.foobar"),
				),
			},
			resource.TestStep{
				ResourceName:            "google_compute_ssl_certificate.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key"},
			},
		},
	})
}

func TestAccComputeSslCertificate_name_prefix(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSslCertificateDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeSslCertificate_name_prefix(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSslCertificateExists(
						"google_compute_ssl_certificate.foobar"),
				),
			},
			resource.TestStep{
				ResourceName:            "google_compute_ssl_certificate.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key", "name_prefix"},
			},
		},
	})
}

func testAccCheckComputeSslCertificateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		found, err := config.clientCompute.SslCertificates.Get(
			config.Project, rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Certificate not found")
		}

		return nil
	}
}

func testAccComputeSslCertificate_basic() string {
	return fmt.Sprintf(`
resource "google_compute_ssl_certificate" "foobar" {
	name = "sslcert-test-%s"
	description = "very descriptive"
	private_key = "${file("test-fixtures/ssl_cert/test.key")}"
	certificate = "${file("test-fixtures/ssl_cert/test.crt")}"
}
`, acctest.RandString(10))
}

func testAccComputeSslCertificate_no_name() string {
	return fmt.Sprintf(`
resource "google_compute_ssl_certificate" "foobar" {
	description = "really descriptive"
	private_key = "${file("test-fixtures/ssl_cert/test.key")}"
	certificate = "${file("test-fixtures/ssl_cert/test.crt")}"
}
`)
}

func testAccComputeSslCertificate_name_prefix() string {
	return fmt.Sprintf(`
resource "google_compute_ssl_certificate" "foobar" {
	name_prefix = "sslcert-test-%s-"
	description = "extremely descriptive"
	private_key = "${file("test-fixtures/ssl_cert/test.key")}"
	certificate = "${file("test-fixtures/ssl_cert/test.crt")}"
}
`, acctest.RandString(10))
}
