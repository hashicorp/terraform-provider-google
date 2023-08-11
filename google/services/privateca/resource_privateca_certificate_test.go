// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package privateca_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccPrivatecaCertificate_privatecaCertificateUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"deletion_protection": false,
		"random_suffix":       acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPrivatecaCertificateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivatecaCertificate_privatecaCertificateStart(context),
			},
			{
				ResourceName:            "google_privateca_certificate.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pool", "name", "location", "certificate_authority"},
			},
			{
				Config: testAccPrivatecaCertificate_privatecaCertificateEnd(context),
			},
			{
				ResourceName:            "google_privateca_certificate.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pool", "name", "location", "certificate_authority"},
			},
			{
				Config: testAccPrivatecaCertificate_privatecaCertificateStart(context),
			},
			{
				ResourceName:            "google_privateca_certificate.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"pool", "name", "location", "certificate_authority"},
			},
		},
	})
}

func testAccPrivatecaCertificate_privatecaCertificateStart(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_ca_pool" "default" {
  location = "us-central1"
  name = "my-pool-%{random_suffix}"
  tier = "ENTERPRISE"
}

resource "google_privateca_certificate_authority" "default" {
	// This example assumes this pool already exists.
	// Pools cannot be deleted in normal test circumstances, so we depend on static pools
  	location = "us-central1"
	pool = google_privateca_ca_pool.default.name
	certificate_authority_id = "tf-test-my-certificate-authority-%{random_suffix}"
	deletion_protection = false
	skip_grace_period = true
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
					cert_sign = true
					crl_sign = true
				}
				extended_key_usage {}
			}
		}
	}
	lifetime = "86400s"
	key_spec {
		algorithm = "RSA_PKCS1_4096_SHA256"
	}
}

resource "google_privateca_certificate" "default" {
	pool = google_privateca_ca_pool.default.name
  	location = "us-central1"
	certificate_authority = google_privateca_certificate_authority.default.certificate_authority_id
	lifetime = "860s"
	name = "my-certificate-%{random_suffix}"
	config {
	  subject_config  {
		subject {
			common_name = "san1.example.com"
			organization = "HashiCorp"
		} 
		subject_alt_name {
		  email_addresses = ["email@example.com"]
		}
	  }
	  x509_config {
		ca_options {
		  is_ca = false
		}
		key_usage {
		  base_key_usage {
			crl_sign = false
			decipher_only = false
		  }
		  extended_key_usage {
			server_auth = false
		  }
		}
	  }
	  public_key {
		format = "PEM"
		key = filebase64("test-fixtures/rsa_public.pem")
	  }
	}
}
`, context)
}

func testAccPrivatecaCertificate_privatecaCertificateEnd(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_ca_pool" "default" {
  location = "us-central1"
  name = "my-pool-%{random_suffix}"
  tier = "ENTERPRISE"
}

resource "google_privateca_certificate_authority" "default" {
	// This example assumes this pool already exists.
	// Pools cannot be deleted in normal test circumstances, so we depend on static pools
	location = "us-central1"
	pool = google_privateca_ca_pool.default.name
	certificate_authority_id = "tf-test-my-certificate-authority-%{random_suffix}"
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
					cert_sign = true
					crl_sign = true
				}
				extended_key_usage {}
			}
		}
	}
	lifetime = "86400s"
	key_spec {
		algorithm = "RSA_PKCS1_4096_SHA256"
	}

	// Disable CA deletion related safe checks for easier cleanup.
	deletion_protection                    = false
	skip_grace_period                      = true
	ignore_active_certificates_on_deletion = true
}

resource "google_privateca_certificate" "default" {
	location = "us-central1"
	pool = google_privateca_ca_pool.default.name
	certificate_authority = google_privateca_certificate_authority.default.certificate_authority_id
	lifetime = "860s"
	name = "my-certificate-%{random_suffix}"
	config {
	  subject_config  {
		subject {
			common_name = "san1.example.com"
			organization = "HashiCorp"
		} 
		subject_alt_name {
		  email_addresses = ["email@example.com"]
		}
	  }
	  x509_config {
		ca_options {
		  is_ca = false
		}
		key_usage {
		  base_key_usage {
			crl_sign = false
			decipher_only = false
		  }
		  extended_key_usage {
			server_auth = false
		  }
		}
	  }
	  public_key {
		format = "PEM"
		key = filebase64("test-fixtures/rsa_public.pem")
	  }
	}
	labels = {
		foo = "bar"
	}
}
`, context)
}
