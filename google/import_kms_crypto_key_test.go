package google

import (
	"testing"

	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"os"
)

func TestAccGoogleKmsCryptoKey_importBasic(t *testing.T) {
	skipIfEnvNotSet(t,
		[]string{
			"GOOGLE_ORG",
			"GOOGLE_BILLING_ACCOUNT",
		}...,
	)

	resourceName := "google_kms_crypto_key.crypto_key"

	//projectId := "terraform-" + acctest.RandString(10)
	projectId := "god-00100"
	projectOrg := os.Getenv("GOOGLE_ORG")
	projectBillingAccount := os.Getenv("GOOGLE_BILLING_ACCOUNT")
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	cryptoKeyName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleKmsCryptoKey_basic(projectId, projectOrg, projectBillingAccount, keyRingName, cryptoKeyName),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
