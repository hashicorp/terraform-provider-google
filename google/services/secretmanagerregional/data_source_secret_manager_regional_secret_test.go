// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package secretmanagerregional_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccDataSourceSecretManagerRegionalRegionalSecret_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"kms_key_name":  acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-secret-manager-managed-central-key6").CryptoKey.Name,
		"random_suffix": acctest.RandString(t, 10),
		"timestamp_1":   "2114-11-30T00:00:00Z",
		"timestamp_2":   "2115-11-30T00:00:00Z",
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecretManagerRegionalRegionalSecret_basic(context),
				Check: resource.ComposeTestCheckFunc(
					acctest.CheckDataSourceStateMatchesResourceState(
						"data.google_secret_manager_regional_secret.reg-secret-datasource",
						"google_secret_manager_regional_secret.reg-secret",
					),
				),
			},
		},
	})
}

func testAccDataSourceSecretManagerRegionalRegionalSecret_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_pubsub_topic" "secret-topic" {
  name = "tf-test-topic%{random_suffix}"
}

resource "google_pubsub_topic_iam_member" "secret_manager_topic_access" {
  topic  = google_pubsub_topic.secret-topic.name
  role   = "roles/pubsub.publisher"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "secret_manager_kms_access" {
  crypto_key_id = "%{kms_key_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}

resource "google_secret_manager_regional_secret" "reg-secret" {
  secret_id = "tf-test-reg-secret-%{random_suffix}"
  location = "us-central1"

  labels = {
    key1 = "val1"
    key2 = "val2"
    key3 = "val3"
  }

  annotations = {
    annotationkey = "annotation-value"
    otherannotation = "othervalue"
  }

  customer_managed_encryption {
    kms_key_name = "%{kms_key_name}"
  }

  topics {
    name = google_pubsub_topic.secret-topic.id
  }

  rotation {
    rotation_period = "7200s"
    next_rotation_time = "%{timestamp_1}"
  }

  expire_time = "%{timestamp_2}"
  version_destroy_ttl = "108000s"

  depends_on = [
    google_pubsub_topic_iam_member.secret_manager_topic_access,
    google_kms_crypto_key_iam_member.secret_manager_kms_access,
  ]
}

data "google_secret_manager_regional_secret" "reg-secret-datasource" {
  secret_id = google_secret_manager_regional_secret.reg-secret.secret_id
  location = google_secret_manager_regional_secret.reg-secret.location
}
`, context)
}
