// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package certificatemanager_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccCertificateManagerCertificate_certificateManagerGoogleManagedCertificateDnsExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCertificateManagerCertificateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerCertificate_certificateManagerGoogleManagedCertificateDnsExample(context),
			},
			{
				ResourceName:            "google_certificate_manager_certificate.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"self_managed", "name", "location", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccCertificateManagerCertificate_certificateManagerGoogleManagedCertificateDnsExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_certificate_manager_certificate" "default" {
  name        = "tf-test-dns-cert%{random_suffix}"
  description = "The default cert"
  scope       = "EDGE_CACHE"
  labels = {
    env = "test"
  }
  managed {
    domains = [
      google_certificate_manager_dns_authorization.instance.domain,
      google_certificate_manager_dns_authorization.instance2.domain,
      ]
    dns_authorizations = [
      google_certificate_manager_dns_authorization.instance.id,
      google_certificate_manager_dns_authorization.instance2.id,
      ]
  }
}


resource "google_certificate_manager_dns_authorization" "instance" {
  name        = "tf-test-dns-auth%{random_suffix}"
  description = "The default dnss"
  domain      = "subdomain%{random_suffix}.hashicorptest.com"
}

resource "google_certificate_manager_dns_authorization" "instance2" {
  name        = "tf-test-dns-auth2%{random_suffix}"
  description = "The default dnss"
  domain      = "subdomain2%{random_suffix}.hashicorptest.com"
}
`, context)
}

func TestAccCertificateManagerCertificate_certificateManagerGoogleManagedCertificateIssuanceConfigExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCertificateManagerCertificateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerCertificate_certificateManagerGoogleManagedCertificateIssuanceConfigExample(context),
			},
			{
				ResourceName:            "google_certificate_manager_certificate.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"self_managed", "name", "location", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccCertificateManagerCertificate_certificateManagerGoogleManagedCertificateIssuanceConfigExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_certificate_manager_certificate" "default" {
  name        = "tf-test-issuance-config-cert%{random_suffix}"
  description = "The default cert"
  scope       = "EDGE_CACHE"
  managed {
    domains = [
        "terraform.subdomain1.com"
      ]
    issuance_config = google_certificate_manager_certificate_issuance_config.issuanceconfig.id
  }
}



# creating certificate_issuance_config to use it in the managed certificate
resource "google_certificate_manager_certificate_issuance_config" "issuanceconfig" {
  name    = "tf-test-issuance-config%{random_suffix}"
  description = "sample description for the certificate issuanceConfigs"
  certificate_authority_config {
    certificate_authority_service_config {
        ca_pool = google_privateca_ca_pool.pool.id
    }
  }
  lifetime = "1814400s"
  rotation_window_percentage = 34
  key_algorithm = "ECDSA_P256"
  depends_on=[google_privateca_certificate_authority.ca_authority]
}
  
resource "google_privateca_ca_pool" "pool" {
  name     = "tf-test-ca-pool%{random_suffix}"
  location = "us-central1"
  tier     = "ENTERPRISE"
}

resource "google_privateca_certificate_authority" "ca_authority" {
  location = "us-central1"
  pool = google_privateca_ca_pool.pool.name
  certificate_authority_id = "tf-test-ca-authority%{random_suffix}"
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
      }
      key_usage {
        base_key_usage {
          cert_sign = true
          crl_sign = true
        }
        extended_key_usage {
          server_auth = true
        }
      }
    }
  }
  key_spec {
    algorithm = "RSA_PKCS1_4096_SHA256"
  }

  // Disable CA deletion related safe checks for easier cleanup.
  deletion_protection                    = false
  skip_grace_period                      = true
  ignore_active_certificates_on_deletion = true
}
`, context)
}

func TestAccCertificateManagerCertificate_certificateManagerSelfManagedCertificateExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCertificateManagerCertificateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerCertificate_certificateManagerSelfManagedCertificateExample(context),
			},
			{
				ResourceName:            "google_certificate_manager_certificate.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"self_managed", "name", "location", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccCertificateManagerCertificate_certificateManagerSelfManagedCertificateExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_certificate_manager_certificate" "default" {
  name        = "tf-test-self-managed-cert%{random_suffix}"
  description = "Global cert"
  scope       = "ALL_REGIONS"
  self_managed {
    pem_certificate = file("test-fixtures/cert.pem")
    pem_private_key = file("test-fixtures/private-key.pem")
  }
}
`, context)
}

func TestAccCertificateManagerCertificate_certificateManagerSelfManagedCertificateRegionalExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCertificateManagerCertificateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerCertificate_certificateManagerSelfManagedCertificateRegionalExample(context),
			},
			{
				ResourceName:            "google_certificate_manager_certificate.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"self_managed", "name", "location", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccCertificateManagerCertificate_certificateManagerSelfManagedCertificateRegionalExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_certificate_manager_certificate" "default" {
  name        = "tf-test-self-managed-cert%{random_suffix}"
  description = "Regional cert"
  location    = "us-central1"
  self_managed {
    pem_certificate = file("test-fixtures/cert.pem")
    pem_private_key = file("test-fixtures/private-key.pem")
  }
}
`, context)
}

func testAccCheckCertificateManagerCertificateDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_certificate_manager_certificate" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{CertificateManagerBasePath}}projects/{{project}}/locations/{{location}}/certificates/{{name}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("CertificateManagerCertificate still exists at %s", url)
			}
		}

		return nil
	}
}
