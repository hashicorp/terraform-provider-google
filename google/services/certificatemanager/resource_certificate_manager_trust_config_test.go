// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package certificatemanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccCertificateManagerTrustConfig_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCertificateManagerTrustConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificateManagerTrustConfig_update0(context),
			},
			{
				ResourceName:            "google_certificate_manager_trust_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccCertificateManagerTrustConfig_update1(context),
			},
			{
				ResourceName:            "google_certificate_manager_trust_config.default",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccCertificateManagerTrustConfig_update0(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_certificate_manager_trust_config" "default" {
  name        = "tf-test-trust-config%{random_suffix}"
  description = "sample description for the trust config"
  location    = "us-central1"

  trust_stores {
    trust_anchors { 
      pem_certificate = file("test-fixtures/cert.pem")
    }
    intermediate_cas { 
      pem_certificate = file("test-fixtures/cert.pem")
    }
  }

  labels = {
    "foo" = "bar"
  }
}
`, context)
}

func testAccCertificateManagerTrustConfig_update1(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_certificate_manager_trust_config" "default" {
  name        = "tf-test-trust-config%{random_suffix}"
  description = "sample description for the trust config 2"
  location    = "us-central1"

  trust_stores {
    trust_anchors { 
      pem_certificate = file("test-fixtures/cert2.pem")
    }
    intermediate_cas { 
      pem_certificate = file("test-fixtures/cert2.pem")
    }
  }

  labels = {
    "bar" = "foo"
  }
}
`, context)
}
