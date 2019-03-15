package google

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGoogleImpersonatedCredential_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_impersonated_credential.default"

	targetServiceAccount := getTestServiceAccountFromEnv(t)
	scopes := []string{"storage-ro", "https://www.googleapis.com/auth/cloud-platform"}
	delegates := []string{}
	lifetime := "30s"
	targetProject := getTestProjectFromEnv()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleImpersonatedCredential_datasource(targetServiceAccount, scopes, delegates, lifetime, targetProject),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "target_service_account", targetServiceAccount),
					resource.TestCheckResourceAttr(resourceName, "lifetime", lifetime),
				),
			},
		},
	})
}

func testAccCheckGoogleImpersonatedCredential_datasource(targetServiceAccount string, scopes []string, delegates []string, lifetime string, target_project string) string {
	return fmt.Sprintf(`

	provider "google" {}

	data "google_client_config" "default" {
	  provider = "google"
	}

	data "google_impersonated_credential" "default" {
	 provider = "google"
	 target_service_account = "%s"
	 scopes = ["storage-ro", "https://www.googleapis.com/auth/cloud-platform"]
	 lifetime = "%s"
	}

	provider "google" {
	   alias  = "impersonated"
	   access_token = "${data.google_impersonated_credential.default.access_token}"
	}

	data "google_project" "project" {
	  provider = "google.impersonated"
	  project_id = "%s"
	}

	`, targetServiceAccount, lifetime, target_project)
}
