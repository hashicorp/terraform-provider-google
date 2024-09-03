// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package certificatemanager_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
)

func TestAccDataSourceGoogleCertificateManagerCertificates_basic(t *testing.T) {
	t.Parallel()

	// Resource identifier used for content testing
	name := fmt.Sprintf("tf-test-certificate-%d", acctest.RandInt(t))
	description := "My acceptance data source test certificates"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCertificateManagerCertificates_basic(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.google_certificate_manager_certificates.certificates", "certificates.#", regexp.MustCompile("^[1-9]")),

					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "region"),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "region", "global"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleCertificateManagerCertificates_full(t *testing.T) {
	t.Parallel()

	// Resource identifier used for content testing
	region := "global"
	id := fmt.Sprintf("projects/%s/locations/%s/certificates", envvar.GetTestProjectFromEnv(), region)
	name := fmt.Sprintf("tf-test-certificate-%d", acctest.RandInt(t))
	description := "My acceptance data source test certificates"
	certificateName := fmt.Sprintf("projects/%s/locations/%s/certificates/%s", envvar.GetTestProjectFromEnv(), region, name)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCertificateManagerCertificates_full(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.google_certificate_manager_certificates.certificates", "certificates.#", regexp.MustCompile("^[1-9]")),

					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "id"),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "id", id),

					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "region"),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "region", region),

					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "certificates.0.name"),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "certificates.0.name", certificateName),

					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "certificates.0.description"),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "certificates.0.description", description),

					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "certificates.0.labels.%"),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "certificates.0.labels.%", "3"),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "certificates.0.labels.terraform", "true"),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "certificates.0.labels.acc-test", "true"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleCertificateManagerCertificates_regionBasic(t *testing.T) {
	t.Parallel()

	// Resource identifier used for content testing
	region := envvar.GetTestRegionFromEnv()
	id := fmt.Sprintf("projects/%s/locations/%s/certificates", envvar.GetTestProjectFromEnv(), region)
	name := fmt.Sprintf("tf-test-certificate-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCertificateManagerCertificates_regionBasic(name, region),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.google_certificate_manager_certificates.certificates", "certificates.#", regexp.MustCompile("^[1-9]")),

					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "id"),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "id", id),

					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "region"),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "region", region),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleCertificateManagerCertificates_managedCertificate(t *testing.T) {
	t.Parallel()

	// Resource identifier used for content testing
	region := "global"
	id := fmt.Sprintf("projects/%s/locations/%s/certificates", envvar.GetTestProjectFromEnv(), region)
	name := fmt.Sprintf("tf-test-certificate-%d", acctest.RandInt(t))
	certificateName := fmt.Sprintf("projects/%s/locations/%s/certificates/%s", envvar.GetTestProjectFromEnv(), region, name)

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCertificateManagerCertificates_managedCertificateBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.google_certificate_manager_certificates.certificates", "certificates.#", regexp.MustCompile("^[1-9]")),

					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "id"),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "id", id),

					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "region"),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "region", region),

					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "certificates.0.name"),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "certificates.0.name", certificateName),

					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "certificates.0.scope"),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "certificates.0.scope", "EDGE_CACHE"),

					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "certificates.0.managed.#"),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "certificates.0.managed.0.domains.#", "1"),

					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "certificates.0.managed.0.state"),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "certificates.0.managed.0.state", "PROVISIONING"),

					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "certificates.0.managed.0.authorization_attempt_info.#"),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "certificates.0.managed.0.authorization_attempt_info.0.details", ""),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "certificates.0.managed.0.authorization_attempt_info.0.domain", "terraform.subdomain1.com"),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "certificates.0.managed.0.authorization_attempt_info.0.failure_reason", ""),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "certificates.0.managed.0.authorization_attempt_info.0.state", "AUTHORIZING"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleCertificateManagerCertificates_managedCertificateDNSAuthorization(t *testing.T) {
	t.Parallel()

	// Resource identifier used for content testing
	region := "global"
	id := fmt.Sprintf("projects/%s/locations/%s/certificates", envvar.GetTestProjectFromEnv(), region)
	name := fmt.Sprintf("tf-test-certificate-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCertificateManagerCertificates_managedCertificateDNSAuthorization(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.google_certificate_manager_certificates.certificates", "certificates.#", regexp.MustCompile("^[1-9]")),

					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "id"),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "id", id),

					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "certificates.0.managed.0.dns_authorizations.#"),
				),
			},
		},
	})
}

func TestAccDataSourceGoogleCertificateManagerCertificates_managedCertificateIssuerConfig(t *testing.T) {
	t.Parallel()

	// Resource identifier used for content testing
	region := "global"
	id := fmt.Sprintf("projects/%s/locations/%s/certificates", envvar.GetTestProjectFromEnv(), region)
	name := fmt.Sprintf("tf-test-certificate-%d", acctest.RandInt(t))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleCertificateManagerCertificates_managedCertificateIssuerConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.google_certificate_manager_certificates.certificates", "certificates.#", regexp.MustCompile("^[1-9]")),

					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "id"),
					resource.TestCheckResourceAttr("data.google_certificate_manager_certificates.certificates", "id", id),

					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "certificates.0.managed.0.issuance_config"),

					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "certificates.0.managed.0.provisioning_issue.#"),
					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "certificates.0.managed.0.provisioning_issue.0.details"),
					resource.TestCheckResourceAttrSet("data.google_certificate_manager_certificates.certificates", "certificates.0.managed.0.provisioning_issue.0.reason"),
				),
			},
		},
	})
}

func testAccDataSourceGoogleCertificateManagerCertificates_basic(certificateName, certificateDescription string) string {
	return fmt.Sprintf(`
resource "google_certificate_manager_certificate" "default" {
  name        = "%s"
  description = "%s"
  self_managed {
    pem_certificate = file("test-fixtures/cert.pem")
    pem_private_key = file("test-fixtures/private-key.pem")
  }

  labels = {
    "terraform" : true,
    "acc-test" : true,
  }
}

data "google_certificate_manager_certificates" "certificates" {
  depends_on = [google_certificate_manager_certificate.default]
}
`, certificateName, certificateDescription)
}

func testAccDataSourceGoogleCertificateManagerCertificates_full(certificateName, certificateDescription string) string {
	return fmt.Sprintf(`
resource "google_certificate_manager_certificate" "default" {
  name        = "%s"
  description = "%s"
  self_managed {
    pem_certificate = file("test-fixtures/cert.pem")
    pem_private_key = file("test-fixtures/private-key.pem")
  }

  labels = {
    "terraform" : true,
    "acc-test" : true,
  }
}

data "google_certificate_manager_certificates" "certificates" {
  filter     = "name:${google_certificate_manager_certificate.default.id}"
  depends_on = [google_certificate_manager_certificate.default]
}
`, certificateName, certificateDescription)
}

func testAccDataSourceGoogleCertificateManagerCertificates_regionBasic(certificateName, region string) string {
	return fmt.Sprintf(`
resource "google_certificate_manager_certificate" "default" {
  name        = "%s"
  location    = "%s"
  self_managed {
    pem_certificate = file("test-fixtures/cert.pem")
    pem_private_key = file("test-fixtures/private-key.pem")
  }

  labels = {
    "terraform" : true,
    "acc-test" : true,
  }
}

data "google_certificate_manager_certificates" "certificates" { 
  filter     = "name:${google_certificate_manager_certificate.default.id}"
  region     = "%s"
  depends_on = [google_certificate_manager_certificate.default]
}
`, certificateName, region, region)
}

func testAccDataSourceGoogleCertificateManagerCertificates_managedCertificateBasic(certificateName string) string {
	return fmt.Sprintf(`
resource "google_certificate_manager_certificate" "default" {
  name        = "%s"
  scope       = "EDGE_CACHE"
  managed {
    domains = [
      "terraform.subdomain1.com"
    ]
  }
}

data "google_certificate_manager_certificates" "certificates" {
  filter     = "name:${google_certificate_manager_certificate.default.id}"
  depends_on = [google_certificate_manager_certificate.default]
}
`, certificateName)
}

func testAccDataSourceGoogleCertificateManagerCertificates_managedCertificateDNSAuthorization(certificateName string) string {
	return fmt.Sprintf(`
resource "google_certificate_manager_certificate" "default" {
  name        = "%s"
  scope       = "EDGE_CACHE"
  managed {
    domains = [
      google_certificate_manager_dns_authorization.default.domain,
    ]
    dns_authorizations = [
      google_certificate_manager_dns_authorization.default.id
    ]
  }
}

resource "google_certificate_manager_dns_authorization" "default" {
  name   = "%s"
  domain = "terraform.subdomain1.com"
}

data "google_certificate_manager_certificates" "certificates" {
  filter     = "name:${google_certificate_manager_certificate.default.id}"
  depends_on = [google_certificate_manager_certificate.default]
}
`, certificateName, certificateName)
}

func testAccDataSourceGoogleCertificateManagerCertificates_managedCertificateIssuerConfig(id string) string {
	return fmt.Sprintf(`
resource "google_certificate_manager_certificate" "default" {
  name        = "%s"
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
  name        = "%s"
  description = "sample description for the certificate issuanceConfigs"
  certificate_authority_config {
    certificate_authority_service_config {
      ca_pool = google_privateca_ca_pool.pool.id
    }
  }
  lifetime                   = "1814400s"
  rotation_window_percentage = 34
  key_algorithm              = "ECDSA_P256"
  depends_on                 = [google_privateca_certificate_authority.ca_authority]
}

resource "google_privateca_ca_pool" "pool" {
  name     = "%s"
  location = "us-central1"
  tier     = "ENTERPRISE"
}

resource "google_privateca_certificate_authority" "ca_authority" {
  location                 = "us-central1"
  pool                     = google_privateca_ca_pool.pool.name
  certificate_authority_id = "%s"
  config {
    subject_config {
      subject {
        organization = "HashiCorp"
        common_name  = "my-certificate-authority"
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
          crl_sign  = true
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

data "google_certificate_manager_certificates" "certificates" {
  filter     = "name:${google_certificate_manager_certificate.default.id}"
  depends_on = [google_certificate_manager_certificate.default]
}
`, id, id, id, id)
}
