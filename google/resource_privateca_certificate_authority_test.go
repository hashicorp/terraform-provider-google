package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPrivatecaCertificateAuthority_privatecaCertificateAuthorityUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"pool_name":           BootstrapSharedCaPoolInLocation(t, "us-central1"),
		"pool_location":       "us-central1",
		"deletion_protection": false,
		"random_suffix":       randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPrivatecaCertificateAuthorityDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivatecaCertificateAuthority_privatecaCertificateAuthorityBasicRoot(context),
				// we added a `desired_state` field in https://github.com/GoogleCloudPlatform/magic-modules/pull/5934, this ensures
				// we don't regress and that CAs are enabled by default
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("google_privateca_certificate_authority.default", "state", "ENABLED"),
				),
			},
			{
				ResourceName:            "google_privateca_certificate_authority.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_active_certificates_on_deletion", "location", "certificate_authority_id", "pool", "deletion_protection", "skip_grace_period"},
			},
			{
				Config: testAccPrivatecaCertificateAuthority_privatecaCertificateAuthorityEnd(context),
			},
			{
				ResourceName:            "google_privateca_certificate_authority.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_active_certificates_on_deletion", "location", "certificate_authority_id", "pool", "deletion_protection", "skip_grace_period"},
			},
			{
				Config: testAccPrivatecaCertificateAuthority_privatecaCertificateAuthorityBasicRoot(context),
			},
			{
				ResourceName:            "google_privateca_certificate_authority.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ignore_active_certificates_on_deletion", "location", "certificate_authority_id", "pool", "deletion_protection", "skip_grace_period"},
			},
		},
	})
}

func TestAccPrivatecaCertificateAuthority_rootCaManageDesiredState(t *testing.T) {
	t.Parallel()

	random_suffix := randString(t, 10)
	context_staged := map[string]interface{}{
		"pool_name":           BootstrapSharedCaPoolInLocation(t, "us-central1"),
		"pool_location":       "us-central1",
		"deletion_protection": false,
		"random_suffix":       random_suffix,
		"desired_state":       "STAGED",
	}

	context_enabled := map[string]interface{}{
		"pool_name":           BootstrapSharedCaPoolInLocation(t, "us-central1"),
		"pool_location":       "us-central1",
		"deletion_protection": false,
		"random_suffix":       random_suffix,
		"desired_state":       "ENABLED",
	}

	context_disabled := map[string]interface{}{
		"pool_name":           BootstrapSharedCaPoolInLocation(t, "us-central1"),
		"pool_location":       "us-central1",
		"deletion_protection": false,
		"random_suffix":       random_suffix,
		"desired_state":       "DISABLED",
	}

	resourceName := "google_privateca_certificate_authority.default"
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPrivatecaCertificateAuthorityDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivatecaCertificateAuthority_privatecaCertificateAuthorityWithDesiredState(context_staged),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "state", "STAGED"),
				),
			},
			{
				Config: testAccPrivatecaCertificateAuthority_privatecaCertificateAuthorityWithDesiredState(context_enabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "state", "ENABLED"),
				),
			},
			{
				Config: testAccPrivatecaCertificateAuthority_privatecaCertificateAuthorityWithDesiredState(context_disabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "state", "DISABLED"),
				),
			},
			{
				Config: testAccPrivatecaCertificateAuthority_privatecaCertificateAuthorityWithDesiredState(context_enabled),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "state", "ENABLED"),
				),
			},
		},
	})
}

func testAccPrivatecaCertificateAuthority_privatecaCertificateAuthorityBasicRoot(context map[string]interface{}) string {
	return Nprintf(`
resource "google_privateca_certificate_authority" "default" {
	// This example assumes this pool already exists.
	// Pools cannot be deleted in normal test circumstances, so we depend on static pools
	pool = "%{pool_name}"
	certificate_authority_id = "tf-test-my-certificate-authority-%{random_suffix}"
	location = "%{pool_location}"
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
`, context)
}

func testAccPrivatecaCertificateAuthority_privatecaCertificateAuthorityEnd(context map[string]interface{}) string {
	return Nprintf(`
resource "google_privateca_certificate_authority" "default" {
	// This example assumes this pool already exists.
	// Pools cannot be deleted in normal test circumstances, so we depend on static pools
	pool = "%{pool_name}"
	certificate_authority_id = "tf-test-my-certificate-authority-%{random_suffix}"
	location = "%{pool_location}"
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
	labels = {
		foo = "bar"
	}
}
`, context)
}

func testAccPrivatecaCertificateAuthority_privatecaCertificateAuthorityWithDesiredState(context map[string]interface{}) string {
	return Nprintf(`
resource "google_privateca_certificate_authority" "default" {
	// This example assumes this pool already exists.
	// Pools cannot be deleted in normal test circumstances, so we depend on static pools
	pool = "%{pool_name}"
	certificate_authority_id = "tf-test-my-certificate-authority-%{random_suffix}"
	location = "%{pool_location}"
	desired_state = "%{desired_state}"
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
`, context)
}
