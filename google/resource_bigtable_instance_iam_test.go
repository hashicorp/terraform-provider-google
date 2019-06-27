package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccBigtableInstanceIamBinding(t *testing.T) {
	t.Parallel()

	instance := "tf-bigtable-iam-" + acctest.RandString(10)
	account := "tf-bigtable-iam-" + acctest.RandString(10)
	role := "roles/bigtable.user"

	importId := fmt.Sprintf("projects/%s/instances/%s %s",
		getTestProjectFromEnv(), instance, role)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccBigtableInstanceIamBinding_basic(instance, account, role),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_bigtable_instance_iam_binding.binding", "role", role),
				),
			},
			{
				ResourceName:      "google_bigtable_instance_iam_binding.binding",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test IAM Binding update
				Config: testAccBigtableInstanceIamBinding_update(instance, account, role),
			},
			{
				ResourceName:      "google_bigtable_instance_iam_binding.binding",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableInstanceIamMember(t *testing.T) {
	t.Parallel()

	instance := "tf-bigtable-iam-" + acctest.RandString(10)
	account := "tf-bigtable-iam-" + acctest.RandString(10)
	role := "roles/bigtable.user"

	importId := fmt.Sprintf("projects/%s/instances/%s %s serviceAccount:%s",
		getTestProjectFromEnv(),
		instance,
		role,
		serviceAccountCanonicalEmail(account))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccBigtableInstanceIamMember(instance, account, role),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_bigtable_instance_iam_member.member", "role", role),
					resource.TestCheckResourceAttr(
						"google_bigtable_instance_iam_member.member", "member", "serviceAccount:"+serviceAccountCanonicalEmail(account)),
				),
			},
			{
				ResourceName:      "google_bigtable_instance_iam_member.member",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableInstanceIamPolicy(t *testing.T) {
	t.Parallel()

	instance := "tf-bigtable-iam-" + acctest.RandString(10)
	account := "tf-bigtable-iam-" + acctest.RandString(10)
	role := "roles/bigtable.user"

	importId := fmt.Sprintf("projects/%s/instances/%s",
		getTestProjectFromEnv(), instance)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccBigtableInstanceIamPolicy(instance, account, role),
			},
			{
				ResourceName:      "google_bigtable_instance_iam_policy.policy",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBigtableInstanceIamBinding_basic(instance, account, role string) string {
	return fmt.Sprintf(testBigtableInstanceIam+`

resource "google_service_account" "test-account1" {
  account_id   = "%s-1"
  display_name = "Dataproc IAM Testing Account"
}

resource "google_service_account" "test-account2" {
  account_id   = "%s-2"
  display_name = "Iam Testing Account"
}

resource "google_bigtable_instance_iam_binding" "binding" {
  instance      = "${google_bigtable_instance.instance.name}"
  role         = "%s"
  members      = [
    "serviceAccount:${google_service_account.test-account1.email}",
  ]
}
`, instance, acctest.RandString(10), account, account, role)
}

func testAccBigtableInstanceIamBinding_update(instance, account, role string) string {
	return fmt.Sprintf(testBigtableInstanceIam+`
resource "google_service_account" "test-account1" {
  account_id   = "%s-1"
  display_name = "Dataproc IAM Testing Account"
}

resource "google_service_account" "test-account2" {
  account_id   = "%s-2"
  display_name = "Iam Testing Account"
}

resource "google_bigtable_instance_iam_binding" "binding" {
  instance      = "${google_bigtable_instance.instance.name}"
  role         = "%s"
  members      = [
    "serviceAccount:${google_service_account.test-account1.email}",
    "serviceAccount:${google_service_account.test-account2.email}",
  ]
}
`, instance, acctest.RandString(10), account, account, role)
}

func testAccBigtableInstanceIamMember(instance, account, role string) string {
	return fmt.Sprintf(testBigtableInstanceIam+`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Dataproc IAM Testing Account"
}

resource "google_bigtable_instance_iam_member" "member" {
  instance      = "${google_bigtable_instance.instance.name}"
  role         = "%s"
  member       = "serviceAccount:${google_service_account.test-account.email}"
}
`, instance, acctest.RandString(10), account, role)
}

func testAccBigtableInstanceIamPolicy(instance, account, role string) string {
	return fmt.Sprintf(testBigtableInstanceIam+`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Dataproc IAM Testing Account"
}

data "google_iam_policy" "policy" {
	binding {
		role    = "%s"
		members = ["serviceAccount:${google_service_account.test-account.email}"]
	}
}

resource "google_bigtable_instance_iam_policy" "policy" {
  instance      = "${google_bigtable_instance.instance.name}"
  policy_data  = "${data.google_iam_policy.policy.policy_data}"
}
`, instance, acctest.RandString(10), account, role)
}

// Smallest instance possible for testing
var testBigtableInstanceIam = `
resource "google_bigtable_instance" "instance" {
	name                  = "%s"
    instance_type = "DEVELOPMENT"

    cluster {
      cluster_id   = "c-%s"
      zone         = "us-central1-b"
      storage_type = "HDD"
    }
}`
