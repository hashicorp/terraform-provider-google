package google

import (
	"testing"

	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"os"
)

func TestAccGoogleKmsKeyRing_importBasic(t *testing.T) {
	skipIfEnvNotSet(t,
		[]string{
			"GOOGLE_ORG",
			"GOOGLE_BILLING_ACCOUNT",
		}...,
	)

	resourceName := "google_kms_key_ring.key_ring"

	projectId := "terraform-" + acctest.RandString(10)
	projectOrg := os.Getenv("GOOGLE_ORG")
	projectBillingAccount := os.Getenv("GOOGLE_BILLING_ACCOUNT")
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleKmsKeyRing_basic(projectId, projectOrg, projectBillingAccount, keyRingName),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
