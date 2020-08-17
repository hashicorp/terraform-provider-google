package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceComputeSslCertificate(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeSslCertificateConfig(randString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_compute_ssl_certificate.cert",
						"google_compute_ssl_certificate.foobar",
						map[string]struct{}{
							"private_key": {},
						},
					),
				),
			},
		},
	})
}

func testAccDataSourceComputeSslCertificateConfig(certName string) string {
	return fmt.Sprintf(`
resource "google_compute_ssl_certificate" "foobar" {
  name        = "cert-test-%s"
  description = "really descriptive"
  private_key = file("test-fixtures/ssl_cert/test.key")
  certificate = file("test-fixtures/ssl_cert/test.crt")
}

data "google_compute_ssl_certificate" "cert" {
  name = google_compute_ssl_certificate.foobar.name
}
`, certName)
}
