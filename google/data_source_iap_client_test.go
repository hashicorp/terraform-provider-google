package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIapClient_Datasource_basic(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        getTestOrgFromEnv(t),
		"org_domain":    getTestOrgDomainFromEnv(t),
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIapClientDatasourceConfig(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceStateWithIgnores(
						"data.google_iap_client.project_client",
						"google_iap_client.project_client",
						map[string]struct{}{
							"brand": {},
						},
					),
				),
			},
		},
	})
}

func testAccIapClientDatasourceConfig(context map[string]interface{}) string {
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
	  
resource "google_iap_client" "project_client" {
  display_name = "Test Client"
  brand        = google_iap_brand.project_brand.name
}

data "google_iap_client" "project_client" {
  brand = google_iap_client.project_client.brand
  client_id = google_iap_client.project_client.client_id
}
`, context)
}
