// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package secretmanagerregional_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSecretManagerRegionalRegionalSecret_import(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerRegionalSecret_basic(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
		},
	})
}

func TestAccSecretManagerRegionalRegionalSecret_labelsUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerRegionalSecret_withoutLabels(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-labels",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_labelsUpdate(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-labels",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_labelsUpdateOther(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-labels",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_withoutLabels(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-labels",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
		},
	})
}

func TestAccSecretManagerRegionalRegionalSecret_annotationsUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerRegionalSecret_withoutAnnotations(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-annotations",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_annotationsUpdate(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-annotations",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_annotationsUpdateOther(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-annotations",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_withoutAnnotations(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-annotations",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
		},
	})
}

func TestAccSecretManagerRegionalRegionalSecret_cmekUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"kms_key_name":       acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-secret-manager-managed-central-key3").CryptoKey.Name,
		"kms_key_name_other": acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-secret-manager-managed-central-key4").CryptoKey.Name,
		"random_suffix":      acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerRegionalSecret_withoutCmek(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-cmek-update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_cmekUpdate(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-cmek-update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_cmekUpdateOther(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-cmek-update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_withoutCmek(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-cmek-update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
		},
	})
}

func TestAccSecretManagerRegionalRegionalSecret_topicsUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerRegionalSecret_withoutTopics(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-topics",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_topicsUpdate(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-topics",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_topicsUpdateOther(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-topics",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_withoutTopics(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-topics",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
		},
	})
}

func TestAccSecretManagerRegionalRegionalSecret_rotationInfoUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"timestamp_1":   "2114-11-30T00:00:00Z",
		"timestamp_2":   "2116-11-30T00:00:00Z",
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerRegionalSecret_rotationBasic(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-rotation-update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_rotationTimeUpdate(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-rotation-update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_rotationPeriodUpdate(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-rotation-update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_rotationBasic(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-rotation-update",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
		},
	})
}

func TestAccSecretManagerRegionalRegionalSecret_expireTimeUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"timestamp_1":   "2114-11-30T00:00:00Z",
		"timestamp_2":   "2116-11-30T00:00:00Z",
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerRegionalSecret_withoutTtlAndExpireTime(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-expiration",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_expireTimeBasic(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-expiration",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_expireTimeUpdate(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-expiration",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_withoutTtlAndExpireTime(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-expiration",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
		},
	})
}

func TestAccSecretManagerRegionalRegionalSecret_ttlUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerRegionalSecret_withoutTtlAndExpireTime(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-expiration",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl", "annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_ttlBasic(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-expiration",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl", "annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_ttlUpdate(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-expiration",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl", "annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_withoutTtlAndExpireTime(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-expiration",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl", "annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
		},
	})
}

func TestAccSecretManagerRegionalRegionalSecret_updateBetweenTtlAndExpireTime(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"timestamp_1":   "2114-11-30T00:00:00Z",
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerRegionalSecret_ttlBasic(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-expiration",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl", "annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_expireTimeBasic(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-expiration",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl", "annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_ttlBasic(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-expiration",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl", "annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
		},
	})
}

func TestAccSecretManagerRegionalRegionalSecret_versionDestroyTtlUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerRegionalSecret_withoutVersionDestroyTtl(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-version-destroy-ttl",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_versionDestroyTtlBasic(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-version-destroy-ttl",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_versionDestroyTtlUpdate(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-version-destroy-ttl",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
			{
				Config: testAccSecretManagerRegionalSecret_withoutVersionDestroyTtl(context),
			},
			{
				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-version-destroy-ttl",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"annotations", "labels", "location", "secret_id", "terraform_labels"},
			},
		},
	})
}

// TODO: Uncomment once google_secret_manager_regional_secret_version is added
// func TestAccSecretManagerRegionalRegionalSecret_versionAliasesUpdate(t *testing.T) {
// 	t.Parallel()
//
// 	context := map[string]interface{}{
// 		"random_suffix": acctest.RandString(t, 10),
// 	}
//
// 	acctest.VcrTest(t, resource.TestCase{
// 		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
// 		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
// 		CheckDestroy:             testAccCheckSecretManagerRegionalRegionalSecretDestroyProducer(t),
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccSecretManagerRegionalSecret_basicRegionalSecretWithVersions(context),
// 			},
// 			{
// 				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-version-aliases",
// 				ImportState:             true,
// 				ImportStateVerify:       true,
// 				ImportStateVerifyIgnore: []string{"ttl", "annotations", "labels", "location", "secret_id", "terraform_labels"},
// 			},
// 			{
// 				Config: testAccSecretManagerRegionalSecret_versionAliasesBasic(context),
// 			},
// 			{
// 				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-version-aliases",
// 				ImportState:             true,
// 				ImportStateVerify:       true,
// 				ImportStateVerifyIgnore: []string{"ttl", "annotations", "labels", "location", "secret_id", "terraform_labels"},
// 			},
// 			{
// 				Config: testAccSecretManagerRegionalSecret_versionAliasesUpdate(context),
// 			},
// 			{
// 				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-version-aliases",
// 				ImportState:             true,
// 				ImportStateVerify:       true,
// 				ImportStateVerifyIgnore: []string{"ttl", "annotations", "labels", "location", "secret_id", "terraform_labels"},
// 			},
// 			{
// 				Config: testAccSecretManagerRegionalSecret_basicRegionalSecretWithVersions(context),
// 			},
// 			{
// 				ResourceName:            "google_secret_manager_regional_secret.regional-secret-with-version-aliases",
// 				ImportState:             true,
// 				ImportStateVerify:       true,
// 				ImportStateVerifyIgnore: []string{"ttl", "annotations", "labels", "location", "secret_id", "terraform_labels"},
// 			},
// 		},
// 	})
// }

func testAccSecretManagerRegionalSecret_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "regional-secret-basic" {
  secret_id = "tf-test-reg-secret-%{random_suffix}"
  location = "us-central1"
}
`, context)
}

func testAccSecretManagerRegionalSecret_withoutLabels(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "regional-secret-with-labels" {
  secret_id = "tf-test-reg-secret-%{random_suffix}"
  location = "us-central1"

  annotations = {
    annotationkey = "annotation-value"
  }
}
`, context)
}

func testAccSecretManagerRegionalSecret_labelsUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "regional-secret-with-labels" {
  secret_id = "tf-test-reg-secret-%{random_suffix}"
  location = "us-central1"

  labels = {
    key1 = "value1"
    key2 = "value2"
    key3 = "value3"
    key4 = "value4"
    key5 = "value5"
  }

  annotations = {
    annotationkey = "annotation-value"
  }
}
`, context)
}

func testAccSecretManagerRegionalSecret_labelsUpdateOther(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "regional-secret-with-labels" {
  secret_id = "tf-test-reg-secret-%{random_suffix}"
  location = "us-central1"

  labels = {
    key1 = "value1"
    key2 = "updatevalue2"
    updatekey3 = "value3"
    updatekey4 = "updatevalue4"
    key6 = "value6"
  }

  annotations = {
    annotationkey = "annotation-value"
  }
}
`, context)
}

func testAccSecretManagerRegionalSecret_withoutAnnotations(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "regional-secret-with-annotations" {
  secret_id = "tf-test-reg-secret-%{random_suffix}"
  location = "us-central1"

  labels = {
    mykey = "myvalue"
  }
}
`, context)
}

func testAccSecretManagerRegionalSecret_annotationsUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "regional-secret-with-annotations" {
  secret_id = "tf-test-reg-secret-%{random_suffix}"
  location = "us-central1"

  labels = {
    mykey = "myvalue"
  }

  annotations = {
    key1 = "value1"
    key2 = "value2"
    key3 = "value3"
    key4 = "value4"
    key5 = "value5"
  }
}
`, context)
}

func testAccSecretManagerRegionalSecret_annotationsUpdateOther(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "regional-secret-with-annotations" {
  secret_id = "tf-test-reg-secret-%{random_suffix}"
  location = "us-central1"

  labels = {
    mykey = "myvalue"
  }

  annotations = {
    key1 = "value1"
    key2 = "updatevalue2"
    updatekey3 = "value3"
    updatekey4 = "updatevalue4"
    key6 = "value6"
  }
}
`, context)
}

func testAccSecretManagerRegionalSecret_withoutCmek(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_kms_crypto_key_iam_member" "kms-regional-secret-binding-1" {
  crypto_key_id = "%{kms_key_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "kms-regional-secret-binding-2" {
  crypto_key_id = "%{kms_key_name_other}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}

resource "google_secret_manager_regional_secret" "regional-secret-cmek-update" {
  secret_id = "tf-test-secret%{random_suffix}"
  location = "us-central1"

  depends_on = [
    google_kms_crypto_key_iam_member.kms-regional-secret-binding-1,
    google_kms_crypto_key_iam_member.kms-regional-secret-binding-2,
  ]
}
`, context)
}

func testAccSecretManagerRegionalSecret_cmekUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_kms_crypto_key_iam_member" "kms-regional-secret-binding-1" {
  crypto_key_id = "%{kms_key_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "kms-regional-secret-binding-2" {
  crypto_key_id = "%{kms_key_name_other}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}

resource "google_secret_manager_regional_secret" "regional-secret-cmek-update" {
  secret_id = "tf-test-secret%{random_suffix}"
  location = "us-central1"

  customer_managed_encryption {
    kms_key_name = "%{kms_key_name}"
  }

  depends_on = [
    google_kms_crypto_key_iam_member.kms-regional-secret-binding-1,
    google_kms_crypto_key_iam_member.kms-regional-secret-binding-2,
  ]
}
`, context)
}

func testAccSecretManagerRegionalSecret_cmekUpdateOther(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_kms_crypto_key_iam_member" "kms-regional-secret-binding-1" {
  crypto_key_id = "%{kms_key_name}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}

resource "google_kms_crypto_key_iam_member" "kms-regional-secret-binding-2" {
  crypto_key_id = "%{kms_key_name_other}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}

resource "google_secret_manager_regional_secret" "regional-secret-cmek-update" {
  secret_id = "tf-test-secret%{random_suffix}"
  location = "us-central1"

  customer_managed_encryption {
    kms_key_name = "%{kms_key_name_other}"
  }

  depends_on = [
    google_kms_crypto_key_iam_member.kms-regional-secret-binding-1,
    google_kms_crypto_key_iam_member.kms-regional-secret-binding-2,
  ]
}
`, context)
}

func testAccSecretManagerRegionalSecret_withoutTopics(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_pubsub_topic" "topic-1" {
  name = "tf-test-topic-1-%{random_suffix}"
}

resource "google_pubsub_topic" "topic-2" {
  name = "tf-test-topic-2-%{random_suffix}"
}

resource "google_pubsub_topic" "topic-3" {
  name = "tf-test-topic-3-%{random_suffix}"
}

resource "google_pubsub_topic_iam_member" "secrets_manager_access_1" {
  topic  = google_pubsub_topic.topic-1.name
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
  role   = "roles/pubsub.publisher"
}

resource "google_pubsub_topic_iam_member" "secrets_manager_access_2" {
  topic  = google_pubsub_topic.topic-2.name
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
  role   = "roles/pubsub.publisher"
}

resource "google_pubsub_topic_iam_member" "secrets_manager_access_3" {
  topic  = google_pubsub_topic.topic-3.name
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
  role   = "roles/pubsub.publisher"
}

resource "google_secret_manager_regional_secret" "regional-secret-with-topics" {
  secret_id = "tf-test-reg-secret-%{random_suffix}"
  location = "us-central1"

  labels = {
    mykey = "myvalue"
  }

  depends_on = [
    google_pubsub_topic_iam_member.secrets_manager_access_1,
    google_pubsub_topic_iam_member.secrets_manager_access_2,
    google_pubsub_topic_iam_member.secrets_manager_access_3,
  ]
}
`, context)
}

func testAccSecretManagerRegionalSecret_topicsUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_pubsub_topic" "topic-1" {
  name = "tf-test-topic-1-%{random_suffix}"
}

resource "google_pubsub_topic" "topic-2" {
  name = "tf-test-topic-2-%{random_suffix}"
}

resource "google_pubsub_topic" "topic-3" {
  name = "tf-test-topic-3-%{random_suffix}"
}

resource "google_pubsub_topic_iam_member" "secrets_manager_access_1" {
  topic  = google_pubsub_topic.topic-1.name
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
  role   = "roles/pubsub.publisher"
}

resource "google_pubsub_topic_iam_member" "secrets_manager_access_2" {
  topic  = google_pubsub_topic.topic-2.name
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
  role   = "roles/pubsub.publisher"
}

resource "google_pubsub_topic_iam_member" "secrets_manager_access_3" {
  topic  = google_pubsub_topic.topic-3.name
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
  role   = "roles/pubsub.publisher"
}

resource "google_secret_manager_regional_secret" "regional-secret-with-topics" {
  secret_id = "tf-test-reg-secret-%{random_suffix}"
  location = "us-central1"

  labels = {
    mykey = "myvalue"
  }

  topics {
    name = google_pubsub_topic.topic-1.id
  }

  topics {
    name = google_pubsub_topic.topic-2.id
  }

  depends_on = [
    google_pubsub_topic_iam_member.secrets_manager_access_1,
    google_pubsub_topic_iam_member.secrets_manager_access_2,
    google_pubsub_topic_iam_member.secrets_manager_access_3,
  ]
}
`, context)
}

func testAccSecretManagerRegionalSecret_topicsUpdateOther(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_pubsub_topic" "topic-1" {
  name = "tf-test-topic-1-%{random_suffix}"
}

resource "google_pubsub_topic" "topic-2" {
  name = "tf-test-topic-2-%{random_suffix}"
}

resource "google_pubsub_topic" "topic-3" {
  name = "tf-test-topic-3-%{random_suffix}"
}

resource "google_pubsub_topic_iam_member" "secrets_manager_access_1" {
  topic  = google_pubsub_topic.topic-1.name
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
  role   = "roles/pubsub.publisher"
}

resource "google_pubsub_topic_iam_member" "secrets_manager_access_2" {
  topic  = google_pubsub_topic.topic-2.name
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
  role   = "roles/pubsub.publisher"
}

resource "google_pubsub_topic_iam_member" "secrets_manager_access_3" {
  topic  = google_pubsub_topic.topic-3.name
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
  role   = "roles/pubsub.publisher"
}

resource "google_secret_manager_regional_secret" "regional-secret-with-topics" {
  secret_id = "tf-test-reg-secret-%{random_suffix}"
  location = "us-central1"

  labels = {
    mykey = "myvalue"
  }

  topics {
    name = google_pubsub_topic.topic-1.id
  }

  topics {
    name = google_pubsub_topic.topic-3.id
  }

  depends_on = [
    google_pubsub_topic_iam_member.secrets_manager_access_1,
    google_pubsub_topic_iam_member.secrets_manager_access_2,
    google_pubsub_topic_iam_member.secrets_manager_access_3,
  ]
}
`, context)
}

func testAccSecretManagerRegionalSecret_rotationBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_pubsub_topic" "topic-update" {
  name = "tf-test-topic%{random_suffix}"
}

resource "google_pubsub_topic_iam_member" "secrets_manager_topic_access" {
  topic  = google_pubsub_topic.topic-update.name
  role   = "roles/pubsub.publisher"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}

resource "google_secret_manager_regional_secret" "regional-secret-with-rotation-update" {
  secret_id = "tf-test-reg-secret%{random_suffix}"
  location = "us-central1"

  topics {
    name = google_pubsub_topic.topic-update.id
  }

  rotation {
    rotation_period = "7200s"
    next_rotation_time = "%{timestamp_1}"
  }

  depends_on = [
    google_pubsub_topic_iam_member.secrets_manager_topic_access,
  ]
}
`, context)
}

func testAccSecretManagerRegionalSecret_rotationTimeUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_pubsub_topic" "topic-update" {
  name = "tf-test-topic%{random_suffix}"
}

resource "google_pubsub_topic_iam_member" "secrets_manager_topic_access" {
  topic  = google_pubsub_topic.topic-update.name
  role   = "roles/pubsub.publisher"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}

resource "google_secret_manager_regional_secret" "regional-secret-with-rotation-update" {
  secret_id = "tf-test-reg-secret%{random_suffix}"
  location = "us-central1"

  topics {
    name = google_pubsub_topic.topic-update.id
  }

  rotation {
    rotation_period = "7200s"
    next_rotation_time = "%{timestamp_2}"
  }

  depends_on = [
    google_pubsub_topic_iam_member.secrets_manager_topic_access,
  ]
}
`, context)
}

func testAccSecretManagerRegionalSecret_rotationPeriodUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {}

resource "google_pubsub_topic" "topic-update" {
  name = "tf-test-topic%{random_suffix}"
}

resource "google_pubsub_topic_iam_member" "secrets_manager_topic_access" {
  topic  = google_pubsub_topic.topic-update.name
  role   = "roles/pubsub.publisher"
  member = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}

resource "google_secret_manager_regional_secret" "regional-secret-with-rotation-update" {
  secret_id = "tf-test-reg-secret%{random_suffix}"
  location = "us-central1"

  topics {
    name = google_pubsub_topic.topic-update.id
  }

  rotation {
    rotation_period = "10800s"
    next_rotation_time = "%{timestamp_2}"
  }

  depends_on = [
    google_pubsub_topic_iam_member.secrets_manager_topic_access,
  ]
}
`, context)
}

func testAccSecretManagerRegionalSecret_withoutTtlAndExpireTime(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "regional-secret-with-expiration" {
  secret_id = "tf-test-reg-secret%{random_suffix}"
  location = "us-central1"

  labels = {
    mylabel = "mykey"
  }

  annotations = {
    myannot = "myannotkey"
  }
}
`, context)
}

func testAccSecretManagerRegionalSecret_expireTimeBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "regional-secret-with-expiration" {
  secret_id = "tf-test-reg-secret%{random_suffix}"
  location = "us-central1"

  labels = {
    mylabel = "mykey"
  }

  annotations = {
    myannot = "myannotkey"
  }

  expire_time = "%{timestamp_1}"
}
`, context)
}

func testAccSecretManagerRegionalSecret_expireTimeUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "regional-secret-with-expiration" {
  secret_id = "tf-test-reg-secret%{random_suffix}"
  location = "us-central1"

  labels = {
    mylabel = "mykey"
  }

  annotations = {
    myannot = "myannotkey"
  }

  expire_time = "%{timestamp_2}"
}
`, context)
}

func testAccSecretManagerRegionalSecret_ttlBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "regional-secret-with-expiration" {
  secret_id = "tf-test-reg-secret%{random_suffix}"
  location = "us-central1"

  labels = {
    mylabel = "mykey"
  }

  annotations = {
    myannot = "myannotkey"
  }

  ttl = "360000s"
}
`, context)
}

func testAccSecretManagerRegionalSecret_ttlUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "regional-secret-with-expiration" {
  secret_id = "tf-test-reg-secret%{random_suffix}"
  location = "us-central1"

  labels = {
    mylabel = "mykey"
  }

  annotations = {
    myannot = "myannotkey"
  }

  ttl = "720000s"
}
`, context)
}

func testAccSecretManagerRegionalSecret_withoutVersionDestroyTtl(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "regional-secret-with-version-destroy-ttl" {
  secret_id = "tf-test-reg-secret%{random_suffix}"
  location = "us-central1"

  labels = {
    mylabel = "mykey"
  }

  annotations = {
    myannot = "myannotkey"
  }
}
`, context)
}

func testAccSecretManagerRegionalSecret_versionDestroyTtlBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "regional-secret-with-version-destroy-ttl" {
  secret_id = "tf-test-reg-secret%{random_suffix}"
  location = "us-central1"

  labels = {
    mylabel = "mykey"
  }

  annotations = {
    myannot = "myannotkey"
  }

  version_destroy_ttl = "90000s"
}
`, context)
}

func testAccSecretManagerRegionalSecret_versionDestroyTtlUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_regional_secret" "regional-secret-with-version-destroy-ttl" {
  secret_id = "tf-test-reg-secret%{random_suffix}"
  location = "us-central1"

  labels = {
    mylabel = "mykey"
  }

  annotations = {
    myannot = "myannotkey"
  }

  version_destroy_ttl = "360000s"
}
`, context)
}

// TODO: Uncomment once google_secret_manager_regional_secret_version is added
// func testAccSecretManagerRegionalSecret_basicRegionalSecretWithVersions(context map[string]interface{}) string {
// 	return acctest.Nprintf(`
// resource "google_secret_manager_regional_secret" "regional-secret-with-version-aliases" {
//   secret_id = "tf-test-reg-secret%{random_suffix}"
//   location = "us-central1"
//
//   labels = {
//     mylabel = "mykey"
//   }
// }
//
// resource "google_secret_manager_regional_secret_version" "reg-secret-version-1" {
//   secret = google_secret_manager_regional_secret.regional-secret-with-version-aliases.id
//
//   secret_data = "very secret data keep it down %{random_suffix}-1"
// }
//
// resource "google_secret_manager_regional_secret_version" "reg-secret-version-2" {
//   secret = google_secret_manager_regional_secret.regional-secret-with-version-aliases.id
//
//   secret_data = "very secret data keep it down %{random_suffix}-2"
// }
//
// resource "google_secret_manager_regional_secret_version" "reg-secret-version-3" {
//   secret = google_secret_manager_regional_secret.regional-secret-with-version-aliases.id
//
//   secret_data = "very secret data keep it down %{random_suffix}-3"
// }
//
// resource "google_secret_manager_regional_secret_version" "reg-secret-version-4" {
//   secret = google_secret_manager_regional_secret.regional-secret-with-version-aliases.id
//
//   secret_data = "very secret data keep it down %{random_suffix}-4"
// }
// `, context)
// }
//
// func testAccSecretManagerRegionalSecret_versionAliasesBasic(context map[string]interface{}) string {
// 	return acctest.Nprintf(`
// resource "google_secret_manager_regional_secret" "regional-secret-with-version-aliases" {
//   secret_id = "tf-test-reg-secret%{random_suffix}"
//   location = "us-central1"
//
//   version_aliases = {
//     firstalias = "1",
//     secondalias = "2",
//     thirdalias = "3",
//     otheralias = "2",
//     somealias = "3"
//   }
//
//   labels = {
//     mylabel = "mykey"
//   }
// }
//
// resource "google_secret_manager_regional_secret_version" "reg-secret-version-1" {
//   secret = google_secret_manager_regional_secret.regional-secret-with-version-aliases.id
//
//   secret_data = "very secret data keep it down %{random_suffix}-1"
// }
//
// resource "google_secret_manager_regional_secret_version" "reg-secret-version-2" {
//   secret = google_secret_manager_regional_secret.regional-secret-with-version-aliases.id
//
//   secret_data = "very secret data keep it down %{random_suffix}-2"
// }
//
// resource "google_secret_manager_regional_secret_version" "reg-secret-version-3" {
//   secret = google_secret_manager_regional_secret.regional-secret-with-version-aliases.id
//
//   secret_data = "very secret data keep it down %{random_suffix}-3"
// }
//
// resource "google_secret_manager_regional_secret_version" "reg-secret-version-4" {
//   secret = google_secret_manager_regional_secret.regional-secret-with-version-aliases.id
//
//   secret_data = "very secret data keep it down %{random_suffix}-4"
// }
// `, context)
// }
//
// func testAccSecretManagerRegionalSecret_versionAliasesUpdate(context map[string]interface{}) string {
// 	return acctest.Nprintf(`
// resource "google_secret_manager_regional_secret" "regional-secret-with-version-aliases" {
//   secret_id = "tf-test-reg-secret%{random_suffix}"
//   location = "us-central1"
//
//   version_aliases = {
//     firstalias = "1",
//     secondaliasupdated = "2",
//     otheralias = "1",
//     somealias = "3",
//     fourthalias = "4"
//   }
//
//   labels = {
//     mylabel = "mykey"
//   }
// }
//
// resource "google_secret_manager_regional_secret_version" "reg-secret-version-1" {
//   secret = google_secret_manager_regional_secret.regional-secret-with-version-aliases.id
//
//   secret_data = "very secret data keep it down %{random_suffix}-1"
// }
//
// resource "google_secret_manager_regional_secret_version" "reg-secret-version-2" {
//   secret = google_secret_manager_regional_secret.regional-secret-with-version-aliases.id
//
//   secret_data = "very secret data keep it down %{random_suffix}-2"
// }
//
// resource "google_secret_manager_regional_secret_version" "reg-secret-version-3" {
//   secret = google_secret_manager_regional_secret.regional-secret-with-version-aliases.id
//
//   secret_data = "very secret data keep it down %{random_suffix}-3"
// }
//
// resource "google_secret_manager_regional_secret_version" "reg-secret-version-4" {
//   secret = google_secret_manager_regional_secret.regional-secret-with-version-aliases.id
//
//   secret_data = "very secret data keep it down %{random_suffix}-4"
// }
// `, context)
// }
