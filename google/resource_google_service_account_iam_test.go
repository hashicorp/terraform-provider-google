package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"reflect"
	"sort"
	"testing"
)

func TestAccGoogleServiceAccountIamBinding(t *testing.T) {
	t.Parallel()

	account := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleServiceAccountIamBinding_basic(account),
				Check: testAccCheckGoogleServiceAccountIam("google_service_account.test_account", "roles/viewer", []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
		},
	})
}

func TestAccGoogleServiceAccountIamMember(t *testing.T) {
	t.Parallel()

	account := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleServiceAccountIamMember_basic(account),
				Check: testAccCheckGoogleServiceAccountIam("google_service_account.test_account", "roles/editor", []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
		},
	})
}

func TestAccGoogleServiceAccountIamPolicy(t *testing.T) {
	t.Parallel()

	account := acctest.RandomWithPrefix("tf-test")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGoogleServiceAccountIamPolicy_basic(account),
				Check: testAccCheckGoogleServiceAccountIam("google_service_account.test_account", "roles/owner", []string{
					fmt.Sprintf("serviceAccount:%s@%s.iam.gserviceaccount.com", account, getTestProjectFromEnv()),
				}),
			},
		},
	})
}

func testAccCheckGoogleServiceAccountIam(n, role string, members []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		p, err := config.clientIAM.Projects.ServiceAccounts.GetIamPolicy(rs.Primary.ID).Do()
		if err != nil {
			return err
		}

		for _, binding := range p.Bindings {
			if binding.Role == role {
				sort.Strings(members)
				sort.Strings(binding.Members)

				if reflect.DeepEqual(members, binding.Members) {
					return nil
				}

				return fmt.Errorf("Binding found but expected members is %v, got %v", members, binding.Members)
			}
		}

		return fmt.Errorf("No binding for role %q", role)
	}
}

func testAccGoogleServiceAccountIamBinding_basic(account string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Iam Testing Account"
}

resource "google_service_account_iam_binding" "foo" {
  service_account_id = "${google_service_account.test_account.id}"
  role        = "roles/viewer"
  members     = ["serviceAccount:${google_service_account.test_account.email}"]
}
`, account)
}

func testAccGoogleServiceAccountIamMember_basic(account string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Iam Testing Account"
}

resource "google_service_account_iam_member" "foo" {
  service_account_id = "${google_service_account.test_account.id}"
  role   = "roles/editor"
  member = "serviceAccount:${google_service_account.test_account.email}"
}
`, account)
}

func testAccGoogleServiceAccountIamPolicy_basic(account string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Iam Testing Account"
}

data "google_iam_policy" "foo" {
	binding {
		role = "roles/owner"

		members = ["serviceAccount:${google_service_account.test_account.email}"]
	}
}

resource "google_service_account_iam_policy" "foo" {
  service_account_id = "${google_service_account.test_account.id}"
  policy_data = "${data.google_iam_policy.foo.policy_data}"
}
`, account)
}
