// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package privateca_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourcePrivatecaCertificateAuthority_privatecaCertificateAuthorityBasicExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"pool_name":     acctest.BootstrapSharedCaPoolInLocation(t, "us-central1"),
		"pool_location": "us-central1",
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPrivatecaCertificateAuthorityDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePrivatecaCertificateAuthority_privatecaCertificateAuthorityBasicExample(context),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.google_privateca_certificate_authority.default", "pem_csr"),
				),
			},
		},
	})
}

func testAccDataSourcePrivatecaCertificateAuthority_privatecaCertificateAuthorityBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_certificate_authority" "default" {
  // This example assumes this pool already exists.
  // Pools cannot be deleted in normal test circumstances, so we depend on static pools
  pool = "%{pool_name}"
  certificate_authority_id = "tf-test-my-certificate-authority%{random_suffix}"
  location = "%{pool_location}"
  type = "SUBORDINATE"
  deletion_protection = false
  config {
    subject_config {
      subject {
        organization = "HashiCorp"
        common_name = "my-certificate-authority"
      }
      subject_alt_name {
        dns_names = ["hashicorp.com"]
      }
    }
    x509_config {
      ca_options {
        is_ca = true
        max_issuer_path_length = 10
      }
      key_usage {
        base_key_usage {
          digital_signature = true
          content_commitment = true
          key_encipherment = false
          data_encipherment = true
          key_agreement = true
          cert_sign = true
          crl_sign = true
          decipher_only = true
        }
        extended_key_usage {
          server_auth = true
          client_auth = false
          email_protection = true
          code_signing = true
          time_stamping = true
        }
      }
    }
  }
  lifetime = "86400s"
  key_spec {
    algorithm = "RSA_PKCS1_4096_SHA256"
  }
}

data "google_privateca_certificate_authority" "default" {
	location = google_privateca_certificate_authority.default.location
	pool = google_privateca_certificate_authority.default.pool
	certificate_authority_id = google_privateca_certificate_authority.default.certificate_authority_id
}

output "csr" {
	value = data.google_privateca_certificate_authority.default.pem_csr
}
`, context)
}
