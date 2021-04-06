package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccServiceAccountIamBinding(t *testing.T) {
	t.Parallel()

	account := fmt.Sprintf("tf-test-%d", randInt(t))

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountIamBinding_basic(account),
				Check:  testAccCheckGoogleServiceAccountIam(t, account, 1),
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

func TestAccServiceAccountIamBinding_withCondition(t *testing.T) {
	t.Parallel()

	account := fmt.Sprintf("tf-test-%d", randInt(t))
	conditionExpr := `request.time < timestamp(\"2020-01-01T00:00:00Z\")`
	conditionTitle := "expires_after_2019_12_31"

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountIamBinding_withCondition(account, "user:admin@hashicorptest.com", conditionTitle, conditionExpr),
				Check:  testAccCheckGoogleServiceAccountIam(t, account, 1),
			},
			{
				ResourceName:      "google_service_account_iam_binding.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s %s %s", serviceAccountCanonicalId(account), "roles/iam.serviceAccountUser", conditionTitle),
			},
		},
	})
}

func TestAccServiceAccountIamBinding_withAndWithoutCondition(t *testing.T) {
	// Resource creation race condition
	skipIfVcr(t)
	t.Parallel()

	account := fmt.Sprintf("tf-test-%d", randInt(t))
	conditionExpr := `request.time < timestamp(\"2020-01-01T00:00:00Z\")`
	conditionTitle := "expires_after_2019_12_31"

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountIamBinding_withAndWithoutCondition(account, "user:admin@hashicorptest.com", conditionTitle, conditionExpr),
				Check:  testAccCheckGoogleServiceAccountIam(t, account, 2),
			},
			{
				ResourceName:      "google_service_account_iam_binding.foo",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s %s", serviceAccountCanonicalId(account), "roles/iam.serviceAccountUser"),
			},
			{
				ResourceName:      "google_service_account_iam_binding.foo2",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     fmt.Sprintf("%s %s %s", serviceAccountCanonicalId(account), "roles/iam.serviceAccountUser", conditionTitle),
			},
		},
	})
}

func TestAccServiceAccountIamMember(t *testing.T) {
	t.Parallel()

	account := fmt.Sprintf("tf-test-%d", randInt(t))
	email := serviceAccountCanonicalEmail(account)
	identity := fmt.Sprintf("serviceAccount:%s", email)

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountIamMember_basic(account, email),
				Check:  testAccCheckGoogleServiceAccountIam(t, account, 1),
			},
			{
				ResourceName:      "google_service_account_iam_member.foo",
				ImportStateId:     fmt.Sprintf("%s %s %s", serviceAccountCanonicalId(account), "roles/iam.serviceAccountUser", identity),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccServiceAccountIamMember_basic(account, strings.ToUpper(email)),
				PlanOnly: true,
			},
			{
				Config:   testAccServiceAccountIamMember_basic(account, strings.Title(email)),
				PlanOnly: true,
			},
		},
	})
}

func TestAccServiceAccountIamMember_withCondition(t *testing.T) {
	t.Parallel()

	account := fmt.Sprintf("tf-test-%d", randInt(t))
	identity := fmt.Sprintf("serviceAccount:%s", serviceAccountCanonicalEmail(account))
	conditionTitle := "expires_after_2019_12_31"

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountIamMember_withCondition(account, conditionTitle),
				Check:  testAccCheckGoogleServiceAccountIam(t, account, 1),
			},
			{
				ResourceName:      "google_service_account_iam_member.foo",
				ImportStateId:     fmt.Sprintf("%s %s %s %s", serviceAccountCanonicalId(account), "roles/iam.serviceAccountUser", identity, conditionTitle),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServiceAccountIamMember_withAndWithoutCondition(t *testing.T) {
	// Resource creation race condition
	skipIfVcr(t)
	t.Parallel()

	account := fmt.Sprintf("tf-test-%d", randInt(t))
	identity := fmt.Sprintf("serviceAccount:%s", serviceAccountCanonicalEmail(account))
	conditionTitle := "expires_after_2019_12_31"

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountIamMember_withAndWithoutCondition(account, conditionTitle),
				Check:  testAccCheckGoogleServiceAccountIam(t, account, 2),
			},
			{
				ResourceName:      "google_service_account_iam_member.foo",
				ImportStateId:     fmt.Sprintf("%s %s %s", serviceAccountCanonicalId(account), "roles/iam.serviceAccountUser", identity),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_service_account_iam_member.foo2",
				ImportStateId:     fmt.Sprintf("%s %s %s %s", serviceAccountCanonicalId(account), "roles/iam.serviceAccountUser", identity, conditionTitle),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServiceAccountIamPolicy(t *testing.T) {
	t.Parallel()

	account := fmt.Sprintf("tf-test-%d", randInt(t))

	vcrTest(t, resource.TestCase{
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

func TestAccServiceAccountIamPolicy_withCondition(t *testing.T) {
	t.Parallel()

	account := fmt.Sprintf("tf-test-%d", randInt(t))

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountIamPolicy_withCondition(account),
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

func TestAccServiceAccountIamMember_federatedIdentity(t *testing.T) {
	t.Parallel()

	account := fmt.Sprintf("tf-test-%d", randInt(t))
	pool := fmt.Sprintf("tf-test-%d", randInt(t))

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountIamMember_federatedIdentity(account, pool),
			},
			{
				ResourceName:      "google_service_account_iam_member.impersonate",
				ImportStateIdFunc: testAccServiceAccountIamMember_generateFederatedIdentityStateId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccServiceAccountIamBinding_federatedIdentity(t *testing.T) {
	t.Parallel()

	account := fmt.Sprintf("tf-test-%d", randInt(t))
	pool := fmt.Sprintf("tf-test-%d", randInt(t))

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountIamBinding_federatedIdentity(account, pool),
				Check:  testAccCheckGoogleServiceAccountIam(t, account, 1),
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

func TestAccServiceAccountIamPolicy_federatedIdentity(t *testing.T) {
	t.Parallel()

	account := fmt.Sprintf("tf-test-%d", randInt(t))
	pool := fmt.Sprintf("tf-test-%d", randInt(t))

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountIamPolicy_federatedIdentity(account, pool),
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

func testAccServiceAccountIamMember_generateFederatedIdentityStateId(state *terraform.State) (string, error) {
	resourceName := "google_service_account_iam_member.impersonate"
	var rawState map[string]string
	for _, m := range state.Modules {
		if len(m.Resources) > 0 {
			if v, ok := m.Resources[resourceName]; ok {
				rawState = v.Primary.Attributes
			}
		}
	}
	return fmt.Sprintf("%s %s %s", rawState["service_account_id"], rawState["role"], rawState["member"]), nil
}

// Ensure that our tests only create the expected number of bindings.
// The content of the binding is tested in the import tests.
func testAccCheckGoogleServiceAccountIam(t *testing.T, account string, numBindings int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)
		p, err := config.NewIamClient(config.userAgent).Projects.ServiceAccounts.GetIamPolicy(serviceAccountCanonicalId(account)).OptionsRequestedPolicyVersion(iamPolicyVersion).Do()
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

func testAccServiceAccountIamBinding_withCondition(account, member, conditionTitle, conditionExpr string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Service Account Iam Testing Account"
}

resource "google_service_account_iam_binding" "foo" {
  service_account_id = google_service_account.test_account.name
  role               = "roles/iam.serviceAccountUser"
  members            = ["%s"]
  condition {
    title       = "%s"
    description = "Expiring at midnight of 2019-12-31"
    expression  = "%s"
  }
}
`, account, member, conditionTitle, conditionExpr)
}

func testAccServiceAccountIamBinding_withAndWithoutCondition(account, member, conditionTitle, conditionExpr string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Service Account Iam Testing Account"
}

resource "google_service_account_iam_binding" "foo" {
  service_account_id = google_service_account.test_account.name
  role               = "roles/iam.serviceAccountUser"
  members            = ["%s"]
}

resource "google_service_account_iam_binding" "foo2" {
  service_account_id = google_service_account.test_account.name
  role               = "roles/iam.serviceAccountUser"
  members            = ["%s"]
  condition {
    title       = "%s"
    description = "Expiring at midnight of 2019-12-31"
    expression  = "%s"
  }
}
`, account, member, member, conditionTitle, conditionExpr)
}

func testAccServiceAccountIamMember_basic(account, email string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Service Account Iam Testing Account"
}

resource "google_service_account_iam_member" "foo" {
  service_account_id = google_service_account.test_account.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:%s"
	depends_on = [google_service_account.test_account]
}
`, account, email)
}

func testAccServiceAccountIamMember_withCondition(account, conditionTitle string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Service Account Iam Testing Account"
}

resource "google_service_account_iam_member" "foo" {
  service_account_id = google_service_account.test_account.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.test_account.email}"
  condition {
    title       = "%s"
    description = "Expiring at midnight of 2019-12-31"
    expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
  }
}
`, account, conditionTitle)
}

func testAccServiceAccountIamMember_withAndWithoutCondition(account, conditionTitle string) string {
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

resource "google_service_account_iam_member" "foo2" {
  service_account_id = google_service_account.test_account.name
  role               = "roles/iam.serviceAccountUser"
  member             = "serviceAccount:${google_service_account.test_account.email}"
  condition {
    title       = "%s"
    description = "Expiring at midnight of 2019-12-31"
    expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
  }
}
`, account, conditionTitle)
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

func testAccServiceAccountIamPolicy_withCondition(account string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Service Account Iam Testing Account"
}

data "google_iam_policy" "foo" {
  binding {
    role = "roles/iam.serviceAccountUser"

    members = ["serviceAccount:${google_service_account.test_account.email}"]
    condition {
      title       = "expires_after_2019_12_31"
      description = "Expiring at midnight of 2019-12-31"
      expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
    }
  }
}

resource "google_service_account_iam_policy" "foo" {
  service_account_id = google_service_account.test_account.name
  policy_data        = data.google_iam_policy.foo.policy_data
}
`, account)
}

func testAccServiceAccountIamMember_federatedIdentity(account, poolId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
	account_id   = "%s"
	display_name = "Service Account Iam Testing Account"
}

resource "google_service_account_iam_member" "impersonate" {
	service_account_id = google_service_account.test_account.name
	role               = "roles/iam.workloadIdentityUser"
	member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.my_pool.name}/attribute.aws_role/arn:aws:sts::999999999999:assumed-role/stack-eu-central-1-lambdaRole"
}

resource "google_iam_workload_identity_pool" "my_pool" {
	workload_identity_pool_id = "%s"
}
`, account, poolId)
}

func testAccServiceAccountIamBinding_federatedIdentity(account, poolId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Service Account Iam Testing Account"
}

resource "google_service_account_iam_binding" "foo" {
  service_account_id = google_service_account.test_account.name
  role               = "roles/iam.serviceAccountUser"
  members            = ["principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.my_pool.name}/attribute.aws_role/arn:aws:sts::999999999999:assumed-role/stack-eu-central-1-lambdaRole"]
}

resource "google_iam_workload_identity_pool" "my_pool" {
	workload_identity_pool_id = "%s"
}
`, account, poolId)
}

func testAccServiceAccountIamPolicy_federatedIdentity(account, poolId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Service Account Iam Testing Account"
}

data "google_iam_policy" "foo" {
  binding {
    role = "roles/iam.serviceAccountUser"

    members = ["principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.my_pool.name}/attribute.aws_role/arn:aws:sts::999999999999:assumed-role/stack-eu-central-1-lambdaRole"]
  }
}

resource "google_service_account_iam_policy" "foo" {
  service_account_id = google_service_account.test_account.name
  policy_data        = data.google_iam_policy.foo.policy_data
}

resource "google_iam_workload_identity_pool" "my_pool" {
	workload_identity_pool_id = "%s"
}
`, account, poolId)
}
