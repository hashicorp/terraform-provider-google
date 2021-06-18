package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSecretManagerSecret_import(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecretManagerSecretDestroyProducer(t),
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

	kmscentral := BootstrapKMSKeyInLocation(t, "us-central1")
	kmseast := BootstrapKMSKeyInLocation(t, "us-east1")
	context1 := map[string]interface{}{
		"pid":                  getTestProjectFromEnv(),
		"random_suffix":        randString(t, 10),
		"kms_key_name_central": kmscentral.CryptoKey.Name,
		"kms_key_name_east":    kmseast.CryptoKey.Name,
	}
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecretManagerSecretDestroyProducer(t),
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

func testAccSecretManagerSecret_basic(context map[string]interface{}) string {
	return Nprintf(`
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
	return Nprintf(`
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
