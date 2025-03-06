// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package eventarc_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEventarcChannel_cryptoKeyUpdate(t *testing.T) {
	t.Parallel()

	region := envvar.GetTestRegionFromEnv()
	context := map[string]interface{}{
		"region":         region,
		"project_name":   envvar.GetTestProjectFromEnv(),
		"project_number": envvar.GetTestProjectNumberFromEnv(),
		"key1":           acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", region, "tf-bootstrap-eventarc-channel-key1").CryptoKey.Name,
		"key2":           acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", region, "tf-bootstrap-eventarc-channel-key2").CryptoKey.Name,
		"random_suffix":  acctest.RandString(t, 10),
	}
	acctest.BootstrapIamMembers(t, []acctest.IamMember{
		{
			Member: "serviceAccount:service-{project_number}@gcp-sa-eventarc.iam.gserviceaccount.com",
			Role:   "roles/cloudkms.cryptoKeyEncrypterDecrypter",
		},
	})

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEventarcChannelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcChannel_setCryptoKey(context),
			},
			{
				ResourceName:      "google_eventarc_channel.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEventarcChannel_cryptoKeyUpdate(context),
			},
			{
				ResourceName:      "google_eventarc_channel.primary",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccEventarcChannel_setCryptoKey(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_channel" "primary" {
  location             = "%{region}"
  name                 = "tf-test-name%{random_suffix}"
  crypto_key_name      = "%{key1}"
  third_party_provider = "projects/%{project_name}/locations/%{region}/providers/datadog"
}
`, context)
}

func testAccEventarcChannel_cryptoKeyUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_channel" "primary" {
  location             = "%{region}"
  name                 = "tf-test-name%{random_suffix}"
  crypto_key_name      = "%{key2}"
  third_party_provider = "projects/%{project_name}/locations/%{region}/providers/datadog"
}
`, context)
}

func TestAccEventarcChannel_longForm(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"region":        envvar.GetTestRegionFromEnv(),
		"project_name":  envvar.GetTestProjectFromEnv(),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckEventarcChannelDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccEventarcChannel_longForm(context),
			},
			{
				ResourceName:            "google_eventarc_channel.primary",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"name", "location", "project"},
			},
		},
	})
}

func testAccEventarcChannel_longForm(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_eventarc_channel" "primary" {
  location             = "long/form/%{region}"
  project              = "projects/%{project_name}"
  name                 = "projects/%{project_name}/locations/%{region}/channels/tf-test-name%{random_suffix}"
  third_party_provider = "projects/%{project_name}/locations/%{region}/providers/datadog"
}
`, context)
}
