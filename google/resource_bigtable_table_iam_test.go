package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBigtableTableIamBinding(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
	t.Parallel()

	instance := "tf-bigtable-iam-" + randString(t, 10)
	cluster := "c-" + randString(t, 10)
	account := "tf-bigtable-iam-" + randString(t, 10)
	role := "roles/bigtable.user"

	importId := fmt.Sprintf("projects/%s/instances/%s/tables/%s %s",
		getTestProjectFromEnv(), instance, cluster, role)

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccBigtableTableIamBinding_basic(instance, cluster, account, role),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_bigtable_table_iam_binding.binding", "role", role),
				),
			},
			{
				ResourceName:      "google_bigtable_table_iam_binding.binding",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test IAM Binding update
				Config: testAccBigtableTableIamBinding_update(instance, cluster, account, role),
			},
			{
				ResourceName:      "google_bigtable_table_iam_binding.binding",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableTableIamMember(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
	t.Parallel()

	instance := "tf-bigtable-iam-" + randString(t, 10)
	cluster := "c-" + randString(t, 10)
	account := "tf-bigtable-iam-" + randString(t, 10)
	role := "roles/bigtable.user"

	importId := fmt.Sprintf("projects/%s/instances/%s/tables/%s %s serviceAccount:%s",
		getTestProjectFromEnv(),
		instance,
		cluster,
		role,
		serviceAccountCanonicalEmail(account))

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccBigtableTableIamMember(instance, cluster, account, role),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_bigtable_table_iam_member.member", "role", role),
					resource.TestCheckResourceAttr(
						"google_bigtable_table_iam_member.member", "member", "serviceAccount:"+serviceAccountCanonicalEmail(account)),
				),
			},
			{
				ResourceName:      "google_bigtable_table_iam_member.member",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBigtableTableIamPolicy(t *testing.T) {
	// bigtable instance does not use the shared HTTP client, this test creates an instance
	skipIfVcr(t)
	t.Parallel()

	instance := "tf-bigtable-iam-" + randString(t, 10)
	cluster := "c-" + randString(t, 10)
	account := "tf-bigtable-iam-" + randString(t, 10)
	role := "roles/bigtable.user"

	importId := fmt.Sprintf("projects/%s/instances/%s/tables/%s",
		getTestProjectFromEnv(), instance, cluster)

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test IAM Binding creation
				Config: testAccBigtableTableIamPolicy(instance, cluster, account, role),
			},
			{
				ResourceName:      "google_bigtable_table_iam_policy.policy",
				ImportStateId:     importId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBigtableTableIamBinding_basic(instance, cluster, account, role string) string {
	return fmt.Sprintf(testBigtableTableIam+`
resource "google_service_account" "test-account1" {
  account_id   = "%s-1"
  display_name = "Bigtable Table IAM Testing Account"
}

resource "google_service_account" "test-account2" {
  account_id   = "%s-2"
  display_name = "Bigtable Table Iam Testing Account"
}

resource "google_bigtable_table_iam_binding" "binding" {
  instance = google_bigtable_instance.instance.name
  table    = google_bigtable_table.table.name
  role     = "%s"
  members  = [
    "serviceAccount:${google_service_account.test-account1.email}",
  ]
}
`, instance, cluster, cluster, account, account, role)
}

func testAccBigtableTableIamBinding_update(instance, cluster, account, role string) string {
	return fmt.Sprintf(testBigtableTableIam+`
resource "google_service_account" "test-account1" {
  account_id   = "%s-1"
  display_name = "Bigtable Table IAM Testing Account"
}

resource "google_service_account" "test-account2" {
  account_id   = "%s-2"
  display_name = "Bigtable Table IAM Testing Account"
}

resource "google_bigtable_table_iam_binding" "binding" {
  instance = google_bigtable_instance.instance.name
  table    = google_bigtable_table.table.name
  role     = "%s"
  members  = [
    "serviceAccount:${google_service_account.test-account1.email}",
    "serviceAccount:${google_service_account.test-account2.email}",
  ]
}
`, instance, cluster, cluster, account, account, role)
}

func testAccBigtableTableIamMember(instance, cluster, account, role string) string {
	return fmt.Sprintf(testBigtableTableIam+`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Bigtable Table IAM Testing Account"
}

resource "google_bigtable_table_iam_member" "member" {
  instance = google_bigtable_instance.instance.name
  table    = google_bigtable_table.table.name
  role     = "%s"
  member   = "serviceAccount:${google_service_account.test-account.email}"
}
`, instance, cluster, cluster, account, role)
}

func testAccBigtableTableIamPolicy(instance, cluster, account, role string) string {
	return fmt.Sprintf(testBigtableTableIam+`
resource "google_service_account" "test-account" {
  account_id   = "%s"
  display_name = "Bigtable Table IAM Testing Account"
}

data "google_iam_policy" "policy" {
  binding {
    role    = "%s"
    members = ["serviceAccount:${google_service_account.test-account.email}"]
  }
}

resource "google_bigtable_table_iam_policy" "policy" {
  instance    = google_bigtable_instance.instance.name
  table       = google_bigtable_table.table.name
  policy_data = data.google_iam_policy.policy.policy_data
}
`, instance, cluster, cluster, account, role)
}

// Smallest instance possible for testing
var testBigtableTableIam = `
resource "google_bigtable_instance" "instance" {
	name                  = "%s"
    instance_type         = "DEVELOPMENT"

    cluster {
      cluster_id   = "%s"
      zone         = "us-central1-b"
      storage_type = "HDD"
    }

    deletion_protection = false
}

resource "google_bigtable_table" "table" {
	name          = "%s"
	instance_name = google_bigtable_instance.instance.name
	split_keys    = ["a", "b", "c"]
}
`
