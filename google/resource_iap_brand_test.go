package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIapBrand_iapBrandExample(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        GetTestOrgFromEnv(t),
		"org_domain":    GetTestOrgDomainFromEnv(t),
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIapBrand_iapBrandExample(context),
			},
			{
				ResourceName:            "google_iap_brand.project_brand",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project"},
			},
		},
	})
}

func testAccIapBrand_iapBrandExample(context map[string]interface{}) string {
	return Nprintf(`
resource "google_project" "project" {
  project_id = "tf-test%{random_suffix}"
  name       = "tf-test%{random_suffix}"
  org_id     = "%{org_id}"
}

resource "google_project_service" "project_service" {
  project = google_project.project.project_id
  service = "iap.googleapis.com"
}

resource "google_iap_brand" "project_brand" {
  support_email     = "support@%{org_domain}"
  application_title = "Cloud IAP protected Application"
  project           = google_project_service.project_service.project
}
`, context)
}
