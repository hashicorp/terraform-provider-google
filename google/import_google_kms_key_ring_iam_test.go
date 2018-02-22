package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccKmsKeyRingIamMember_importBasic(t *testing.T) {
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

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKmsKeyRingIamMember_basic(projectId, orgId, billingAccount, account, keyRingName, roleId),
			},

			resource.TestStep{
				ResourceName:  "google_kms_key_ring_iam_member.foo",
				ImportStateId: fmt.Sprintf("%s %s serviceAccount:%s@%s.iam.gserviceaccount.com", keyRingId.terraformId(), roleId, account, projectId),
				ImportState:   true,
			},
		},
	})
}

func TestAccKmsKeyRingIamPolicy_importBasic(t *testing.T) {
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

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKmsKeyRingIamPolicy_basic(projectId, orgId, billingAccount, account, keyRingName, roleId),
			},

			resource.TestStep{
				ResourceName:  "google_kms_key_ring_iam_policy.foo",
				ImportStateId: keyRingId.terraformId(),
				ImportState:   true,
			},
		},
	})
}

func TestAccKmsKeyRingIamBinding_importBasic(t *testing.T) {
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
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKmsKeyRingIamBinding_basic(projectId, orgId, billingAccount, account, keyRingName, roleId),
			},

			resource.TestStep{
				ResourceName:  "google_kms_key_ring_iam_binding.foo",
				ImportStateId: fmt.Sprintf("%s %s", keyRingId.terraformId(), roleId),
				ImportState:   true,
			},
		},
	})
}
