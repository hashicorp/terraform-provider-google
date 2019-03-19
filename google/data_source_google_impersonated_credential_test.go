package google

import (
	"testing"

	"fmt"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGoogleImpersonatedCredential_basic(t *testing.T) {
	t.Parallel()

	resourceName := "data.google_impersonated_credential.default"

	sourceServiceAccountEmail := getTestServiceAccountFromEnv(t)
	targetServiceAccountID := acctest.RandomWithPrefix("tf-test")
	targetServiceAccountEmail := fmt.Sprintf(
		"%s@%s.iam.gserviceaccount.com",
		targetServiceAccountID,
		getTestProjectFromEnv(),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckGoogleImpersonatedCredential_datasource(sourceServiceAccountEmail, targetServiceAccountID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "target_service_account", targetServiceAccountEmail),
					resource.TestCheckOutput("target-email", targetServiceAccountEmail),
				),
			},
		},
	})
}

func testAccCheckGoogleImpersonatedCredential_datasource(sourceServiceAccountEmail string, targetServiceAccountID string) string {

	return fmt.Sprintf(`

	provider "google" {
		scopes = [
			"https://www.googleapis.com/auth/cloud-platform",
		]
	}

	resource "google_service_account" "targetSA" {
		account_id   = "%s"
	}

	resource "google_service_account_iam_binding" "token-creator-iam" {
		service_account_id = "${google_service_account.targetSA.name}"
		role               = "roles/iam.serviceAccountTokenCreator"
		members = [
			"serviceAccount:%s",
		]
	}

	data "google_impersonated_credential" "default" {
		target_service_account = "${google_service_account.targetSA.email}"
		scopes = ["https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/cloud-platform"]
		lifetime = "60s"

		depends_on = ["google_service_account_iam_binding.token-creator-iam"]
	}

	provider "google" {
		alias  = "impersonated"
		scopes = [
			"https://www.googleapis.com/auth/cloud-platform",
			"https://www.googleapis.com/auth/userinfo.email",
		]
		access_token = "${data.google_impersonated_credential.default.access_token}"
	}

	data "google_client_openid_userinfo" "me" {
		provider = "google.impersonated"
	}

	output "target-email" {
		value = "${data.google_client_openid_userinfo.me.email}"
	}
	`, targetServiceAccountID, sourceServiceAccountEmail)
}
