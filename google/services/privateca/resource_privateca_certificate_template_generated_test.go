// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: DCL     ***
//
// ----------------------------------------------------------------------------
//
//     This file is managed by Magic Modules (https://github.com/GoogleCloudPlatform/magic-modules)
//     and is based on the DCL (https://github.com/GoogleCloudPlatform/declarative-resource-client-library).
//     Changes will need to be made to the DCL or Magic Modules instead of here.
//
//     We are not currently able to accept contributions to this file. If changes
//     are required, please file an issue at https://github.com/hashicorp/terraform-provider-google/issues/new/choose
//
// ----------------------------------------------------------------------------

package privateca_test

import (
	"context"
	"fmt"
	dcl "github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
	privateca "github.com/GoogleCloudPlatform/declarative-resource-client-library/services/google/privateca"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func TestAccPrivatecaCertificateTemplate_BasicCertificateTemplate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"project_name":  envvar.GetTestProjectFromEnv(),
		"region":        envvar.GetTestRegionFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckPrivatecaCertificateTemplateDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccPrivatecaCertificateTemplate_BasicCertificateTemplate(context),
			},
			{
				ResourceName:            "google_privateca_certificate_template.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"predefined_values.0.key_usage.0.extended_key_usage"},
			},
			{
				Config: testAccPrivatecaCertificateTemplate_BasicCertificateTemplateUpdate0(context),
			},
			{
				ResourceName:            "google_privateca_certificate_template.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"predefined_values.0.key_usage.0.extended_key_usage"},
			},
		},
	})
}

func testAccPrivatecaCertificateTemplate_BasicCertificateTemplate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_certificate_template" "primary" {
  location    = "%{region}"
  name        = "tf-test-template%{random_suffix}"
  description = "An updated sample certificate template"

  identity_constraints {
    allow_subject_alt_names_passthrough = true
    allow_subject_passthrough           = true

    cel_expression {
      description = "Always true"
      expression  = "true"
      location    = "any.file.anywhere"
      title       = "Sample expression"
    }
  }

  labels = {
    label-two = "value-two"
  }

  passthrough_extensions {
    additional_extensions {
      object_id_path = [1, 6]
    }

    known_extensions = ["EXTENDED_KEY_USAGE"]
  }

  predefined_values {
    additional_extensions {
      object_id {
        object_id_path = [1, 6]
      }

      value    = "c3RyaW5nCg=="
      critical = true
    }

    aia_ocsp_servers = ["string"]

    ca_options {
      is_ca                  = false
      max_issuer_path_length = 6
    }

    key_usage {
      base_key_usage {
        cert_sign          = false
        content_commitment = true
        crl_sign           = false
        data_encipherment  = true
        decipher_only      = true
        digital_signature  = true
        encipher_only      = true
        key_agreement      = true
        key_encipherment   = true
      }

      extended_key_usage {
        client_auth      = true
        code_signing     = true
        email_protection = true
        ocsp_signing     = true
        server_auth      = true
        time_stamping    = true
      }

      unknown_extended_key_usages {
        object_id_path = [1, 6]
      }
    }

    policy_ids {
      object_id_path = [1, 6]
    }
  }

  project = "%{project_name}"
}


`, context)
}

func testAccPrivatecaCertificateTemplate_BasicCertificateTemplateUpdate0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_privateca_certificate_template" "primary" {
  location    = "%{region}"
  name        = "tf-test-template%{random_suffix}"
  description = "A sample certificate template"

  identity_constraints {
    allow_subject_alt_names_passthrough = false
    allow_subject_passthrough           = false

    cel_expression {
      description = "Always false"
      expression  = "false"
      location    = "update.certificate_template.json"
      title       = "New sample expression"
    }
  }

  labels = {
    label-one = "value-one"
  }

  passthrough_extensions {
    additional_extensions {
      object_id_path = [1, 7]
    }

    known_extensions = ["BASE_KEY_USAGE"]
  }

  predefined_values {
    additional_extensions {
      object_id {
        object_id_path = [1, 7]
      }

      value    = "bmV3LXN0cmluZw=="
      critical = false
    }

    aia_ocsp_servers = ["new-string"]

    ca_options {
      is_ca                  = true
      max_issuer_path_length = 7
    }

    key_usage {
      base_key_usage {
        cert_sign          = true
        content_commitment = false
        crl_sign           = true
        data_encipherment  = false
        decipher_only      = false
        digital_signature  = false
        encipher_only      = false
        key_agreement      = false
        key_encipherment   = false
      }

      extended_key_usage {
        client_auth      = false
        code_signing     = false
        email_protection = false
        ocsp_signing     = false
        server_auth      = false
        time_stamping    = false
      }

      unknown_extended_key_usages {
        object_id_path = [1, 7]
      }
    }

    policy_ids {
      object_id_path = [1, 7]
    }
  }

  project = "%{project_name}"
}


`, context)
}

func testAccCheckPrivatecaCertificateTemplateDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "rs.google_privateca_certificate_template" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			billingProject := ""
			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			obj := &privateca.CertificateTemplate{
				Location:    dcl.String(rs.Primary.Attributes["location"]),
				Name:        dcl.String(rs.Primary.Attributes["name"]),
				Description: dcl.String(rs.Primary.Attributes["description"]),
				Project:     dcl.StringOrNil(rs.Primary.Attributes["project"]),
				CreateTime:  dcl.StringOrNil(rs.Primary.Attributes["create_time"]),
				UpdateTime:  dcl.StringOrNil(rs.Primary.Attributes["update_time"]),
			}

			client := transport_tpg.NewDCLPrivatecaClient(config, config.UserAgent, billingProject, 0)
			_, err := client.GetCertificateTemplate(context.Background(), obj)
			if err == nil {
				return fmt.Errorf("google_privateca_certificate_template still exists %v", obj)
			}
		}
		return nil
	}
}
