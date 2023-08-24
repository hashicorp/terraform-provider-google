// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package privateca_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccPrivatecaCaPool_privatecaCapoolUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPrivatecaCaPoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivatecaCaPool_privatecaCapoolStart(context),
			},
			{
				ResourceName:            "google_privateca_ca_pool.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
			{
				Config: testAccPrivatecaCaPool_privatecaCapoolEnd(context),
			},
			{
				ResourceName:            "google_privateca_ca_pool.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
			{
				Config: testAccPrivatecaCaPool_privatecaCapoolStart(context),
			},
			{
				ResourceName:            "google_privateca_ca_pool.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
		},
	})
}

func testAccPrivatecaCaPool_privatecaCapoolStart(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_ca_pool" "default" {
  name = "tf-test-my-capool%{random_suffix}"
  location = "us-central1"
  tier = "ENTERPRISE"
  publishing_options {
    publish_ca_cert = false
    publish_crl = true
  }
  labels = {
    foo = "bar"
  }
  issuance_policy {
    allowed_key_types {
      elliptic_curve {
        signature_algorithm = "ECDSA_P256"
      }
    }
    allowed_key_types {
      rsa {
        min_modulus_size = 5
        max_modulus_size = 10
      }
    }
    maximum_lifetime = "50000s"
    allowed_issuance_modes {
      allow_csr_based_issuance = true
      allow_config_based_issuance = false
    }
    identity_constraints {
      allow_subject_passthrough = false
      allow_subject_alt_names_passthrough = true
      cel_expression {
        expression = "subject_alt_names.all(san, san.type == DNS || san.type == EMAIL )"
        title = "My title"
      }
    }
    baseline_values {
      aia_ocsp_servers = ["example.com"]
      additional_extensions {
        critical = true
        value = "asdf"
        object_id {
          object_id_path = [1, 5]
        }
      }
      policy_ids {
        object_id_path = [1, 7]
      }
      policy_ids {
        object_id_path = [1,5,7]
      }
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
          cert_sign = false
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
}
`, context)
}

func testAccPrivatecaCaPool_privatecaCapoolEnd(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_ca_pool" "default" {
  name = "tf-test-my-capool%{random_suffix}"
  location = "us-central1"
  tier = "ENTERPRISE"
  publishing_options {
    publish_ca_cert = true
    publish_crl = true
  }
  labels = {
    foo = "bar"
    baz = "qux"
  }
  issuance_policy {
    allowed_key_types {
      elliptic_curve {
        signature_algorithm = "ECDSA_P256"
      }
    }
    allowed_key_types {
      rsa {
        min_modulus_size = 6
      }
    }
    maximum_lifetime = "3000s"
    allowed_issuance_modes {
      allow_csr_based_issuance = true
      allow_config_based_issuance = true
    }
    identity_constraints {
      allow_subject_passthrough = true
      allow_subject_alt_names_passthrough = true
      cel_expression {
        expression = "subject_alt_names.all(san, san.type == DNS || san.type == EMAIL )"
        title = "My title3"
      }
    }
    baseline_values {
      aia_ocsp_servers = ["example.com", "hashicorp.com"]
      additional_extensions {
        critical = true
        value = "asdf"
        object_id {
          object_id_path = [1, 7]
        }
      }
      policy_ids {
        object_id_path = [1, 5]
      }
      policy_ids {
        object_id_path = [1, 7]
      }
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
          key_agreement = false
          cert_sign = false
          crl_sign = true
          decipher_only = false
        }
        extended_key_usage {
          server_auth = false
          client_auth = true
          email_protection = true
          code_signing = true
          time_stamping = false
        }
      }
    }
  }
}
`, context)
}

func TestAccPrivatecaCaPool_privatecaCapoolEmptyBaseline(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPrivatecaCaPoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivatecaCaPool_privatecaCapoolEmptyBaseline(context),
			},
			{
				ResourceName:            "google_privateca_ca_pool.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
		},
	})
}

func testAccPrivatecaCaPool_privatecaCapoolEmptyBaseline(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_ca_pool" "default" {
  name = "tf-test-my-capool%{random_suffix}"
  location = "us-central1"
  tier = "ENTERPRISE"
  publishing_options {
    publish_ca_cert = false
    publish_crl = true
  }
  labels = {
    foo = "bar"
  }
  issuance_policy {
    baseline_values {
      additional_extensions {
        critical = false
        value = "asdf"
        object_id {
          object_id_path = [1, 6]
        }
      }
      ca_options {
        is_ca = false
      }
      key_usage {
        base_key_usage {
          digital_signature = false
        }
        extended_key_usage {
          server_auth = false
        }
      }
    }
  }
}
`, context)
}

func TestAccPrivatecaCaPool_privatecaCapoolEmptyPublishingOptions(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPrivatecaCaPoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivatecaCaPool_privatecaCapoolEmptyPublishingOptions(context),
			},
			{
				ResourceName:            "google_privateca_ca_pool.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
		},
	})
}

func testAccPrivatecaCaPool_privatecaCapoolEmptyPublishingOptions(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_ca_pool" "default" {
  name = "tf-test-my-capool%{random_suffix}"
  location = "us-central1"
  tier = "ENTERPRISE"
  publishing_options {
    publish_ca_cert = false
    publish_crl = false
  }
  labels = {
    foo = "bar"
  }
}
`, context)
}

func TestAccPrivatecaCaPool_updateCaOption(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPrivatecaCaPoolDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivatecaCaPool_privatecaCapoolCaOptionIsCaIsTrueAndMaxPathIsPositive(context),
			},
			{
				ResourceName:            "google_privateca_ca_pool.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
			{
				Config: testAccPrivatecaCaPool_privatecaCapoolCaOptionIsCaIsFalse(context),
			},
			{
				ResourceName:            "google_privateca_ca_pool.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
			{
				Config: testAccPrivatecaCaPool_privatecaCapoolCaOptionMaxIssuerPathLenghIsZero(context),
			},
			{
				ResourceName:            "google_privateca_ca_pool.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location"},
			},
		},
	})
}

func testAccPrivatecaCaPool_privatecaCapoolCaOptionIsCaIsTrueAndMaxPathIsPositive(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_ca_pool" "default" {
  name = "tf-test-my-capool%{random_suffix}"
  location = "us-central1"
  tier = "ENTERPRISE"

  issuance_policy {
    baseline_values {
      ca_options {
        is_ca = true
        max_issuer_path_length = 10
      }
      key_usage {
        base_key_usage {
          digital_signature = true
        }
        extended_key_usage {
          server_auth = true
        }
      }
    }
  }
}
`, context)
}

func testAccPrivatecaCaPool_privatecaCapoolCaOptionIsCaIsFalse(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_ca_pool" "default" {
  name = "tf-test-my-capool%{random_suffix}"
  location = "us-central1"
  tier = "ENTERPRISE"

  issuance_policy {
    baseline_values {
      ca_options {
        non_ca = true
        is_ca = false
      }
      key_usage {
        base_key_usage {
          digital_signature = true
        }
        extended_key_usage {
          server_auth = true
        }
      }
    }
  }
}
`, context)
}

func testAccPrivatecaCaPool_privatecaCapoolCaOptionMaxIssuerPathLenghIsZero(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_ca_pool" "default" {
  name = "tf-test-my-capool%{random_suffix}"
  location = "us-central1"
  tier = "ENTERPRISE"

  issuance_policy {
    baseline_values {
      ca_options {
        zero_max_issuer_path_length = true
        max_issuer_path_length = 0
      }
      key_usage {
        base_key_usage {
          digital_signature = true
        }
        extended_key_usage {
          server_auth = true
        }
      }
    }
  }
}
`, context)
}
