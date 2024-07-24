// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/services/compute"
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

func TestUnitComputeManagedSslCertificate_AbsoluteDomainSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"new trailing dot": {
			Old:                "sslcert.tf-test.club",
			New:                "sslcert.tf-test.club.",
			ExpectDiffSuppress: true,
		},
		"old trailing dot": {
			Old:                "sslcert.tf-test.club.",
			New:                "sslcert.tf-test.club",
			ExpectDiffSuppress: true,
		},
		"same trailing dot": {
			Old:                "sslcert.tf-test.club.",
			New:                "sslcert.tf-test.club.",
			ExpectDiffSuppress: false,
		},
		"different trailing dot": {
			Old:                "sslcert.tf-test.club.",
			New:                "sslcert.tf-test.clubs.",
			ExpectDiffSuppress: false,
		},
		"different no trailing dot": {
			Old:                "sslcert.tf-test.club",
			New:                "sslcert.tf-test.clubs",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if compute.AbsoluteDomainSuppress("managed.0.domains.", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
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
