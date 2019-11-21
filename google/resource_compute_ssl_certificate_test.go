package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccComputeSslCertificate_no_name(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeSslCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSslCertificate_no_name(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSslCertificateExists(
						"google_compute_ssl_certificate.foobar"),
				),
			},
			{
				ResourceName:            "google_compute_ssl_certificate.foobar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"private_key"},
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
		// We don't specify a name, but it is saved during create
		name := rs.Primary.Attributes["name"]

		found, err := config.clientCompute.SslCertificates.Get(
			config.Project, name).Do()
		if err != nil {
			return err
		}

		if found.Name != name {
			return fmt.Errorf("Certificate not found")
		}

		return nil
	}
}

func testAccComputeSslCertificate_no_name() string {
	return fmt.Sprintf(`
resource "google_compute_ssl_certificate" "foobar" {
  description = "really descriptive"
  private_key = file("test-fixtures/ssl_cert/test.key")
  certificate = file("test-fixtures/ssl_cert/test.crt")
}
`)
}
