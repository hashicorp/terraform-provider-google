package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccKmsCryptoKeyIamMember_importBasic(t *testing.T) {
	t.Parallel()

	orgId := getTestOrgFromEnv(t)
	projectId := acctest.RandomWithPrefix("tf-test")
	billingAccount := getTestBillingAccountFromEnv(t)
	account := acctest.RandomWithPrefix("tf-test")
	roleId := "roles/cloudkms.cryptoKeyEncrypter"
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	keyRingId := &kmsKeyRingId{
		Project:  projectId,
		Location: DEFAULT_KMS_TEST_LOCATION,
		Name:     keyRingName,
	}
	cryptoKeyName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKmsCryptoKeyIamMember_basic(projectId, orgId, billingAccount, account, keyRingName, cryptoKeyName, roleId),
			},

			resource.TestStep{
				ResourceName:  "google_kms_crypto_key_iam_member.foo",
				ImportStateId: fmt.Sprintf("%s/%s %s serviceAccount:%s@%s.iam.gserviceaccount.com", keyRingId.terraformId(), cryptoKeyName, roleId, account, projectId),
				ImportState:   true,
			},
		},
	})
}

func TestAccKmsCryptoKeyIamBinding_importBasic(t *testing.T) {
	t.Parallel()

	orgId := getTestOrgFromEnv(t)
	projectId := acctest.RandomWithPrefix("tf-test")
	billingAccount := getTestBillingAccountFromEnv(t)
	account := acctest.RandomWithPrefix("tf-test")
	roleId := "roles/cloudkms.cryptoKeyEncrypter"
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	keyRingId := &kmsKeyRingId{
		Project:  projectId,
		Location: DEFAULT_KMS_TEST_LOCATION,
		Name:     keyRingName,
	}
	cryptoKeyName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKmsCryptoKeyIamBinding_basic(projectId, orgId, billingAccount, account, keyRingName, cryptoKeyName, roleId),
			},

			resource.TestStep{
				ResourceName:  "google_kms_crypto_key_iam_binding.foo",
				ImportStateId: fmt.Sprintf("%s/%s %s", keyRingId.terraformId(), cryptoKeyName, roleId),
				ImportState:   true,
			},
		},
	})
}
