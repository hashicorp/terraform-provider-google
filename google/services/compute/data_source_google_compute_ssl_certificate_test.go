// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceComputeSslCertificate(t *testing.T) {
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComputeSslCertificateConfig(acctest.RandString(t, 10)),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceStateWithIgnores(
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
  private_key = file("test-fixtures/test.key")
  certificate = file("test-fixtures/test.crt")
}

data "google_compute_ssl_certificate" "cert" {
  name = google_compute_ssl_certificate.foobar.name
}
`, certName)
}
