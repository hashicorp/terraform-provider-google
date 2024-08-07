// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package netapp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccNetappkmsconfig_kmsConfigCreateExample_Update(t *testing.T) {
	// t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckNetappkmsconfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccNetappkmsconfig_kmsConfigCreateExample_Full(context),
			},
			{
				ResourceName:            "google_netapp_kmsconfig.kmsConfig",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			},
			{
				Config: testAccNetappkmsconfig_kmsConfigCreateExample_Update(context),
			},
			{
				ResourceName:            "google_netapp_kmsconfig.kmsConfig",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"location", "name", "labels", "terraform_labels"},
			},
		},
	})
}

func testAccNetappkmsconfig_kmsConfigCreateExample_Full(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_kms_key_ring" "keyring" {
		name     = "tf-test-key-ring%{random_suffix}"
		location = "us-east4"
	}
	  
	resource "google_kms_crypto_key" "crypto_key" {
		name            = "tf-test-crypto-name%{random_suffix}"
		key_ring        = google_kms_key_ring.keyring.id
	}
	  
	resource "google_netapp_kmsconfig" "kmsConfig" {
		name = "tf-test-kms-test%{random_suffix}"
		description="this is a test description"
		crypto_key_name=google_kms_crypto_key.crypto_key.id
		location="us-east4"
	}
`, context)
}

func testAccNetappkmsconfig_kmsConfigCreateExample_Update(context map[string]interface{}) string {
	return acctest.Nprintf(`
	resource "google_kms_key_ring" "keyring" {
		name     = "tf-test-key-ring%{random_suffix}"
		location = "us-east4"
	}
	  
	resource "google_kms_crypto_key" "crypto_key" {
		name            = "tf-test-crypto-name%{random_suffix}"
		key_ring        = google_kms_key_ring.keyring.id
	}
	  
	resource "google_netapp_kmsconfig" "kmsConfig" {
		name = "tf-test-kms-test%{random_suffix}"
		description="kmsconfig update"
		crypto_key_name=google_kms_crypto_key.crypto_key.id
		location="us-east4"
		labels = {
			"foo": "bar",
		}
	}
`, context)
}
