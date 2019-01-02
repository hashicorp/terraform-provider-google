package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGoogleKmsKeyRing_basic(t *testing.T) {
	projectId := "terraform-" + acctest.RandString(10)
	projectOrg := getTestOrgFromEnv(t)
	projectBillingAccount := getTestBillingAccountFromEnv(t)
	keyRingName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGoogleKmsKeyRingConfig(projectId, projectOrg, projectBillingAccount, keyRingName),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceMatchesResourceCheck("data.google_kms_key_ring.test", "google_kms_key_ring.project", []string{"project", "location", "self_link"}),
				),
			},
		},
	})
}

func testAccDataSourceGoogleKmsKeyRingConfig(projectId, projectOrg, projectBillingAccount, keyRingName string) string {
	return fmt.Sprintf(`
%s

data "google_kms_key_ring" "test" {
  name     = "${google_kms_key_ring.key_ring.name}"
  location = "us-central1"
}`, testGoogleKmsKeyRing_basic(projectId, projectOrg, projectBillingAccount, keyRingName))
}
