// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeSslCertificate_no_name(t *testing.T) {
	// Randomness
	acctest.SkipIfVcr(t)
	t.Parallel()

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeSslCertificateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeSslCertificate_no_name(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeSslCertificateExists(
						t, "google_compute_ssl_certificate.foobar"),
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

func testAccCheckComputeSslCertificateExists(t *testing.T, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acctest.GoogleProviderConfig(t)
		// We don't specify a name, but it is saved during create
		name := rs.Primary.Attributes["name"]

		found, err := config.NewComputeClient(config.UserAgent).SslCertificates.Get(
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
  private_key = file("test-fixtures/test.key")
  certificate = file("test-fixtures/test.crt")
}
`)
}
