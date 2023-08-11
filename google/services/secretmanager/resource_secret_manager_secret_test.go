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
    automatic = true
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
    automatic = true
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
