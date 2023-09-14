// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package secretmanager_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSecretManagerSecret_import(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerSecret_basic(context),
			},
			{
				ResourceName:            "google_secret_manager_secret.secret-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl"},
			},
		},
	})
}

func TestAccSecretManagerSecret_cmek(t *testing.T) {
	t.Parallel()

	kmscentral := acctest.BootstrapKMSKeyInLocation(t, "us-central1")
	kmseast := acctest.BootstrapKMSKeyInLocation(t, "us-east1")
	context1 := map[string]interface{}{
		"pid":                  envvar.GetTestProjectFromEnv(),
		"random_suffix":        acctest.RandString(t, 10),
		"kms_key_name_central": kmscentral.CryptoKey.Name,
		"kms_key_name_east":    kmseast.CryptoKey.Name,
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretMangerSecret_cmek(context1),
			},
			{
				ResourceName:            "google_secret_manager_secret.secret-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl"},
			},
		},
	})
}

func TestAccSecretManagerSecret_annotationsUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerSecret_annotationsBasic(context),
			},
			{
				ResourceName:            "google_secret_manager_secret.secret-with-annotations",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl"},
			},
			{
				Config: testAccSecretManagerSecret_annotationsUpdate(context),
			},
			{
				ResourceName:            "google_secret_manager_secret.secret-with-annotations",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl"},
			},
			{
				Config: testAccSecretManagerSecret_annotationsBasic(context),
			},
			{
				ResourceName:            "google_secret_manager_secret.secret-with-annotations",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl"},
			},
		},
	})
}

func TestAccSecretManagerSecret_versionAliasesUpdate(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretManagerSecret_basicWithSecretVersions(context),
			},
			{
				ResourceName:            "google_secret_manager_secret.secret-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl"},
			},
			{
				Config: testAccSecretManagerSecret_versionAliasesBasic(context),
			},
			{
				ResourceName:            "google_secret_manager_secret.secret-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl"},
			},
			{
				Config: testAccSecretManagerSecret_versionAliasesUpdate(context),
			},
			{
				ResourceName:            "google_secret_manager_secret.secret-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl"},
			},
			{
				Config: testAccSecretManagerSecret_basicWithSecretVersions(context),
			},
			{
				ResourceName:            "google_secret_manager_secret.secret-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl"},
			},
		},
	})
}

func TestAccSecretManagerSecret_userManagedCmekUpdate(t *testing.T) {
	t.Parallel()

	kmscentral := acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-secret-manager-managed-central-key1")
	kmseast := acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-east1", "tf-secret-manager-managed-east-key1")
	kmscentralother := acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "us-central1", "tf-secret-manager-managed-central-key2")
	context := map[string]interface{}{
		"pid":                        envvar.GetTestProjectFromEnv(),
		"random_suffix":              acctest.RandString(t, 10),
		"kms_key_name_central":       kmscentral.CryptoKey.Name,
		"kms_key_name_east":          kmseast.CryptoKey.Name,
		"kms_key_name_central_other": kmscentralother.CryptoKey.Name,
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretMangerSecret_userManagedCmekBasic(context),
			},
			{
				ResourceName:            "google_secret_manager_secret.secret-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl"},
			},
			{
				Config: testAccSecretMangerSecret_userManagedCmekUpdate(context),
			},
			{
				ResourceName:            "google_secret_manager_secret.secret-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl"},
			},
			{
				Config: testAccSecretMangerSecret_userManagedCmekUpdate2(context),
			},
			{
				ResourceName:            "google_secret_manager_secret.secret-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl"},
			},
			{
				Config: testAccSecretMangerSecret_userManagedCmekBasic(context),
			},
			{
				ResourceName:            "google_secret_manager_secret.secret-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl"},
			},
		},
	})
}

func TestAccSecretManagerSecret_automaticCmekUpdate(t *testing.T) {
	t.Parallel()

	suffix := acctest.RandString(t, 10)
	key1 := acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "global", "tf-secret-manager-automatic-key1")
	key2 := acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, "ENCRYPT_DECRYPT", "global", "tf-secret-manager-automatic-key2")
	context := map[string]interface{}{
		"pid":            envvar.GetTestProjectFromEnv(),
		"random_suffix":  suffix,
		"kms_key_name_1": key1.CryptoKey.Name,
		"kms_key_name_2": key2.CryptoKey.Name,
	}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSecretManagerSecretDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretMangerSecret_automaticBasic(context),
			},
			{
				ResourceName:            "google_secret_manager_secret.secret-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl", "replication.0.automatic", "replication.0.auto"},
			},
			{
				Config: testAccSecretMangerSecret_automaticCmekBasic(context),
			},
			{
				ResourceName:            "google_secret_manager_secret.secret-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl"},
			},
			{
				Config: testAccSecretMangerSecret_automaticCmekUpdate(context),
			},
			{
				ResourceName:            "google_secret_manager_secret.secret-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl"},
			},
			{
				Config: testAccSecretMangerSecret_automaticCmekUpdate2(context),
			},
			{
				ResourceName:            "google_secret_manager_secret.secret-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl"},
			},
			{
				Config: testAccSecretMangerSecret_automaticCmekBasic(context),
			},
			{
				ResourceName:            "google_secret_manager_secret.secret-basic",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ttl"},
			},
		},
	})
}

func testAccSecretManagerSecret_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-test-secret-%{random_suffix}"
  
  labels = {
    label = "my-label"
  }

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
      replicas {
        location = "us-east1"
      }
    }
  }

  ttl = "3600s"

}
`, context)
}

func testAccSecretMangerSecret_cmek(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
  project_id = "%{pid}"
}
resource "google_project_iam_member" "kms-secret-binding" {
  project = data.google_project.project.project_id
  role    = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-test-secret-%{random_suffix}"
  
  labels = {
    label = "my-label"
  }
  replication {
    user_managed {
      replicas {
		location = "us-central1"
		customer_managed_encryption {
			kms_key_name = "%{kms_key_name_central}"
		}
	  }
	replicas {
		location = "us-east1"
		customer_managed_encryption {
			kms_key_name = "%{kms_key_name_east}"
		}
      }
	  
    }
  }
  project   = google_project_iam_member.kms-secret-binding.project
}
`, context)
}

func testAccSecretManagerSecret_annotationsBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "secret-with-annotations" {
  secret_id = "tf-test-secret-%{random_suffix}"

  labels = {
    label = "my-label"
  }

  annotations = {
    key1 = "someval"
    key2 = "someval2"
    key3 = "someval3"
    key4 = "someval4"
    key5 = "someval5"
  }

  replication {
    auto {}
  }
}
`, context)
}

func testAccSecretManagerSecret_annotationsUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "secret-with-annotations" {
  secret_id = "tf-test-secret-%{random_suffix}"

  labels = {
    label = "my-label"
  }

  annotations = {
    key1 = "someval"
    key2update = "someval2"
    key3 = "someval3update"
    key4update = "someval4update"
  }

  replication {
    auto {}
  }
}
`, context)
}

func testAccSecretManagerSecret_basicWithSecretVersions(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-test-secret-%{random_suffix}"
  
  labels = {
    label = "my-label"
  }

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
      replicas {
        location = "us-east1"
      }
    }
  }
}

resource "google_secret_manager_secret_version" "secret-version-1" {
  secret = google_secret_manager_secret.secret-basic.id

  secret_data = "some-secret-data-%{random_suffix}-1"
}

resource "google_secret_manager_secret_version" "secret-version-2" {
  secret = google_secret_manager_secret.secret-basic.id

  secret_data = "some-secret-data-%{random_suffix}-2"
}

resource "google_secret_manager_secret_version" "secret-version-3" {
  secret = google_secret_manager_secret.secret-basic.id

  secret_data = "some-secret-data-%{random_suffix}-3"
}

resource "google_secret_manager_secret_version" "secret-version-4" {
  secret = google_secret_manager_secret.secret-basic.id

  secret_data = "some-secret-data-%{random_suffix}-4"
}
`, context)
}

func testAccSecretManagerSecret_versionAliasesBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-test-secret-%{random_suffix}"
  
  labels = {
    label = "my-label"
  }

  version_aliases = {
    firstalias = "1",
    secondalias = "2",
    thirdalias = "3",
    otheralias = "2",
    somealias = "3"
  }

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
      replicas {
        location = "us-east1"
      }
    }
  }
}

resource "google_secret_manager_secret_version" "secret-version-1" {
  secret = google_secret_manager_secret.secret-basic.id

  secret_data = "some-secret-data-%{random_suffix}-1"
}

resource "google_secret_manager_secret_version" "secret-version-2" {
  secret = google_secret_manager_secret.secret-basic.id

  secret_data = "some-secret-data-%{random_suffix}-2"
}

resource "google_secret_manager_secret_version" "secret-version-3" {
  secret = google_secret_manager_secret.secret-basic.id

  secret_data = "some-secret-data-%{random_suffix}-3"
}

resource "google_secret_manager_secret_version" "secret-version-4" {
  secret = google_secret_manager_secret.secret-basic.id

  secret_data = "some-secret-data-%{random_suffix}-4"
}
`, context)
}

func testAccSecretManagerSecret_versionAliasesUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-test-secret-%{random_suffix}"
  
  labels = {
    label = "my-label"
  }

  version_aliases = {
    firstalias = "1",
    secondaliasupdated = "2",
    otheralias = "1",
    somealias = "3",
    fourthalias = "4"
  }

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
      replicas {
        location = "us-east1"
      }
    }
  }
}

resource "google_secret_manager_secret_version" "secret-version-1" {
  secret = google_secret_manager_secret.secret-basic.id

  secret_data = "some-secret-data-%{random_suffix}-1"
}

resource "google_secret_manager_secret_version" "secret-version-2" {
  secret = google_secret_manager_secret.secret-basic.id

  secret_data = "some-secret-data-%{random_suffix}-2"
}

resource "google_secret_manager_secret_version" "secret-version-3" {
  secret = google_secret_manager_secret.secret-basic.id

  secret_data = "some-secret-data-%{random_suffix}-3"
}

resource "google_secret_manager_secret_version" "secret-version-4" {
  secret = google_secret_manager_secret.secret-basic.id

  secret_data = "some-secret-data-%{random_suffix}-4"
}
`, context)
}

func testAccSecretMangerSecret_userManagedCmekBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
  project_id = "%{pid}"
}
resource "google_kms_crypto_key_iam_member" "kms-central-binding-1" {
  crypto_key_id = "%{kms_key_name_central}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}
resource "google_kms_crypto_key_iam_member" "kms-central-binding-2" {
  crypto_key_id = "%{kms_key_name_central_other}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}
resource "google_kms_crypto_key_iam_member" "kms-east-binding" {
  crypto_key_id = "%{kms_key_name_east}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-test-secret-%{random_suffix}"
  
  labels = {
    label = "my-label"
  }
  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
      replicas {
        location = "us-east1"
      }
    }
  }
  depends_on = [
    google_kms_crypto_key_iam_member.kms-central-binding-1,
    google_kms_crypto_key_iam_member.kms-central-binding-2,
    google_kms_crypto_key_iam_member.kms-east-binding,
  ]
}
`, context)
}

func testAccSecretMangerSecret_userManagedCmekUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
  project_id = "%{pid}"
}
resource "google_kms_crypto_key_iam_member" "kms-central-binding-1" {
  crypto_key_id = "%{kms_key_name_central}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}
resource "google_kms_crypto_key_iam_member" "kms-central-binding-2" {
  crypto_key_id = "%{kms_key_name_central_other}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}
resource "google_kms_crypto_key_iam_member" "kms-east-binding" {
  crypto_key_id = "%{kms_key_name_east}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-test-secret-%{random_suffix}"
  
  labels = {
    label = "my-label"
  }
  replication {
    user_managed {
      replicas {
        location = "us-central1"
        customer_managed_encryption {
          kms_key_name = "%{kms_key_name_central}"
        }
      }
      replicas {
        location = "us-east1"
        customer_managed_encryption {
          kms_key_name = "%{kms_key_name_east}"
        }
      }
    }
  }
  depends_on = [
    google_kms_crypto_key_iam_member.kms-central-binding-1,
    google_kms_crypto_key_iam_member.kms-central-binding-2,
    google_kms_crypto_key_iam_member.kms-east-binding,
  ]
}
`, context)
}

func testAccSecretMangerSecret_userManagedCmekUpdate2(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
  project_id = "%{pid}"
}
resource "google_kms_crypto_key_iam_member" "kms-central-binding-1" {
  crypto_key_id = "%{kms_key_name_central}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}
resource "google_kms_crypto_key_iam_member" "kms-central-binding-2" {
  crypto_key_id = "%{kms_key_name_central_other}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}
resource "google_kms_crypto_key_iam_member" "kms-east-binding" {
  crypto_key_id = "%{kms_key_name_east}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-test-secret-%{random_suffix}"
  
  labels = {
    label = "my-label"
  }
  replication {
    user_managed {
      replicas {
        location = "us-central1"
        customer_managed_encryption {
          kms_key_name = "%{kms_key_name_central_other}"
        }
      }
      replicas {
        location = "us-east1"
        customer_managed_encryption {
          kms_key_name = "%{kms_key_name_east}"
        }
      }
    }
  }
  depends_on = [
    google_kms_crypto_key_iam_member.kms-central-binding-1,
    google_kms_crypto_key_iam_member.kms-central-binding-2,
    google_kms_crypto_key_iam_member.kms-east-binding,
  ]
}
`, context)
}

func testAccSecretMangerSecret_automaticBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
  project_id = "%{pid}"
}
resource "google_kms_crypto_key_iam_member" "kms-secret-binding-1" {
  crypto_key_id = "%{kms_key_name_1}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}
resource "google_kms_crypto_key_iam_member" "kms-secret-binding-2" {
  crypto_key_id = "%{kms_key_name_2}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-test-secret-%{random_suffix}"
  
  labels = {
    label = "my-label"
  }
  replication {
    automatic = true
  }
  depends_on = [
    google_kms_crypto_key_iam_member.kms-secret-binding-1,
    google_kms_crypto_key_iam_member.kms-secret-binding-2,
  ]
}
`, context)
}

func testAccSecretMangerSecret_automaticCmekBasic(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
  project_id = "%{pid}"
}
resource "google_kms_crypto_key_iam_member" "kms-secret-binding-1" {
  crypto_key_id = "%{kms_key_name_1}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}
resource "google_kms_crypto_key_iam_member" "kms-secret-binding-2" {
  crypto_key_id = "%{kms_key_name_2}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-test-secret-%{random_suffix}"
  
  labels = {
    label = "my-label"
  }
  replication {
    auto {}
  }
  depends_on = [
    google_kms_crypto_key_iam_member.kms-secret-binding-1,
    google_kms_crypto_key_iam_member.kms-secret-binding-2,
  ]
}
`, context)
}

func testAccSecretMangerSecret_automaticCmekUpdate(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
  project_id = "%{pid}"
}
resource "google_kms_crypto_key_iam_member" "kms-secret-binding-1" {
  crypto_key_id = "%{kms_key_name_1}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}
resource "google_kms_crypto_key_iam_member" "kms-secret-binding-2" {
  crypto_key_id = "%{kms_key_name_2}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-test-secret-%{random_suffix}"
  
  labels = {
    label = "my-label"
  }
  replication {
    auto {
      customer_managed_encryption {
        kms_key_name = "%{kms_key_name_1}"
      }
    }
  }
  depends_on = [
    google_kms_crypto_key_iam_member.kms-secret-binding-1,
    google_kms_crypto_key_iam_member.kms-secret-binding-2,
  ]
}
`, context)
}

func testAccSecretMangerSecret_automaticCmekUpdate2(context map[string]interface{}) string {
	return acctest.Nprintf(`
data "google_project" "project" {
  project_id = "%{pid}"
}
resource "google_kms_crypto_key_iam_member" "kms-secret-binding-1" {
  crypto_key_id = "%{kms_key_name_1}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}
resource "google_kms_crypto_key_iam_member" "kms-secret-binding-2" {
  crypto_key_id = "%{kms_key_name_2}"
  role          = "roles/cloudkms.cryptoKeyEncrypterDecrypter"
  member        = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-secretmanager.iam.gserviceaccount.com"
}
resource "google_secret_manager_secret" "secret-basic" {
  secret_id = "tf-test-secret-%{random_suffix}"
  
  labels = {
    label = "my-label"
  }
  replication {
    auto {
      customer_managed_encryption {
        kms_key_name = "%{kms_key_name_2}"
      }
    }
  }
  depends_on = [
    google_kms_crypto_key_iam_member.kms-secret-binding-1,
    google_kms_crypto_key_iam_member.kms-secret-binding-2,
  ]
}
`, context)
}
