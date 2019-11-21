package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccServiceAccountIamBinding(t *testing.T) {
	t.Parallel()

	account := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountIamBinding_basic(account),
				Check:  testAccCheckGoogleServiceAccountIam(account, 1),
			},
			{
				ResourceName:      "google_service_account_iam_binding.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s %s", serviceAccountCanonicalId(account), "roles/iam.serviceAccountUser"),
			},
		},
	})
}

func TestAccServiceAccountIamMember(t *testing.T) {
	t.Parallel()

	account := acctest.RandomWithPrefix("tf-test")
	identity := fmt.Sprintf("serviceAccount:%s", serviceAccountCanonicalEmail(account))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountIamMember_basic(account),
				Check:  testAccCheckGoogleServiceAccountIam(account, 1),
			},
			{
				ResourceName:      "google_service_account_iam_member.foo",
				ImportStateId:     fmt.Sprintf("%s %s %s", serviceAccountCanonicalId(account), "roles/iam.serviceAccountUser", identity),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServiceAccountIamPolicy(t *testing.T) {
	t.Parallel()

	account := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountIamPolicy_basic(account),
			},
			{
				ResourceName:      "google_service_account_iam_policy.foo",
				ImportStateId:     serviceAccountCanonicalId(account),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Ensure that our tests only create the expected number of bindings.
// The content of the binding is tested in the import tests.
func testAccCheckGoogleServiceAccountIam(account string, numBindings int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		p, err := config.clientIAM.Projects.ServiceAccounts.GetIamPolicy(serviceAccountCanonicalId(account)).Do()
		if err != nil {
			return err
		}

		if len(p.Bindings) != numBindings {
			return fmt.Errorf("Expected exactly %d binding(s) for account %q, was %d", numBindings, account, len(p.Bindings))
		}

		return nil
	}
}

func serviceAccountCanonicalId(account string) string {
	return fmt.Sprintf("projects/%s/serviceAccounts/%s@%s.iam.gserviceaccount.com", getTestProjectFromEnv(), account, getTestProjectFromEnv())
}

func serviceAccountCanonicalEmail(account string) string {
	return fmt.Sprintf("%s@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv())
}

func testAccServiceAccountIamBinding_basic(account string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Service Account Iam Testing Account"
}

resource "google_service_account_iam_binding" "foo" {
  service_account_id = google_service_account.test_account.name
  role               = "roles/iam.serviceAccountUser"
  members            = ["user:admin@hashicorptest.com"]
}
`, account)
}

func testAccServiceAccountIamMember_basic(account string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Service Account Iam Testing Account"
}

resource "google_service_account_iam_member" "foo" {
  service_account_id = google_service_account.test_account.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.test_account.email}"
}
`, account)
}

func testAccServiceAccountIamPolicy_basic(account string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Service Account Iam Testing Account"
}

data "google_iam_policy" "foo" {
  binding {
    role = "roles/iam.serviceAccountUser"

    members = ["serviceAccount:${google_service_account.test_account.email}"]
  }
}

resource "google_service_account_iam_policy" "foo" {
  service_account_id = google_service_account.test_account.name
  policy_data        = data.google_iam_policy.foo.policy_data
}
`, account)
}
