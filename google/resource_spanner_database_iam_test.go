package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccSpannerDatabaseIamBinding(t *testing.T) {
	t.Parallel()

	account := acctest.RandomWithPrefix("tf-test")
	role := "roles/spanner.databaseAdmin"
	project := getTestProjectFromEnv()
	database := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	instance := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerDatabaseIamBinding_basic(account, instance, database, role),
			},
			{
				ResourceName: "google_spanner_database_iam_binding.foo",
				ImportStateId: fmt.Sprintf("%s %s", spannerDatabaseId{
					Project:  project,
					Instance: instance,
					Database: database,
				}.terraformId(), role),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccSpannerDatabaseIamBinding_update(account, instance, database, role),
			},
			{
				ResourceName: "google_spanner_database_iam_binding.foo",
				ImportStateId: fmt.Sprintf("%s %s", spannerDatabaseId{
					Project:  project,
					Instance: instance,
					Database: database,
				}.terraformId(), role),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerDatabaseIamMember(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	account := acctest.RandomWithPrefix("tf-test")
	role := "roles/spanner.databaseAdmin"
	database := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	instance := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccSpannerDatabaseIamMember_basic(account, instance, database, role),
			},
			{
				ResourceName: "google_spanner_database_iam_member.foo",
				ImportStateId: fmt.Sprintf("%s %s serviceAccount:%s@%s.iam.gserviceaccount.com", spannerDatabaseId{
					Instance: instance,
					Database: database,
					Project:  project,
				}.terraformId(), role, account, project),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerDatabaseIamPolicy(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	account := acctest.RandomWithPrefix("tf-test")
	role := "roles/spanner.databaseAdmin"
	database := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	instance := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerDatabaseIamPolicy_basic(account, instance, database, role),
			},
			// Test a few import formats
			{
				ResourceName: "google_spanner_database_iam_policy.foo",
				ImportStateId: spannerDatabaseId{
					Instance: instance,
					Database: database,
					Project:  project,
				}.terraformId(),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSpannerDatabaseIamBinding_basic(account, instance, database, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Spanner Database Iam Testing Account"
}

resource "google_spanner_instance" "instance" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s"
  num_nodes    = 1
}

resource "google_spanner_database" "database" {
  instance = google_spanner_instance.instance.name
  name     = "%s"
}

resource "google_spanner_database_iam_binding" "foo" {
  project  = google_spanner_database.database.project
  database = google_spanner_database.database.name
  instance = google_spanner_database.database.instance
  role     = "%s"
  members  = ["serviceAccount:${google_service_account.test_account.email}"]
}
`, account, instance, instance, database, roleId)
}

func testAccSpannerDatabaseIamBinding_update(account, instance, database, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Spanner Database Iam Testing Account"
}

resource "google_service_account" "test_account_2" {
  account_id   = "%s-2"
  display_name = "Spanner Database Iam Testing Account"
}

resource "google_spanner_instance" "instance" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s"
  num_nodes    = 1
}

resource "google_spanner_database" "database" {
  instance = google_spanner_instance.instance.name
  name     = "%s"
}

resource "google_spanner_database_iam_binding" "foo" {
  project  = google_spanner_database.database.project
  database = google_spanner_database.database.name
  instance = google_spanner_database.database.instance
  role     = "%s"
  members = [
    "serviceAccount:${google_service_account.test_account.email}",
    "serviceAccount:${google_service_account.test_account_2.email}",
  ]
}
`, account, account, instance, instance, database, roleId)
}

func testAccSpannerDatabaseIamMember_basic(account, instance, database, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Spanner Database Iam Testing Account"
}

resource "google_spanner_instance" "instance" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s"
  num_nodes    = 1
}

resource "google_spanner_database" "database" {
  instance = google_spanner_instance.instance.name
  name     = "%s"
}

resource "google_spanner_database_iam_member" "foo" {
  project  = google_spanner_database.database.project
  database = google_spanner_database.database.name
  instance = google_spanner_database.database.instance
  role     = "%s"
  member   = "serviceAccount:${google_service_account.test_account.email}"
}
`, account, instance, instance, database, roleId)
}

func testAccSpannerDatabaseIamPolicy_basic(account, instance, database, roleId string) string {
	return fmt.Sprintf(`
resource "google_service_account" "test_account" {
  account_id   = "%s"
  display_name = "Spanner Database Iam Testing Account"
}

resource "google_spanner_instance" "instance" {
  name         = "%s"
  config       = "regional-us-central1"
  display_name = "%s"
  num_nodes    = 1
}

resource "google_spanner_database" "database" {
  instance = google_spanner_instance.instance.name
  name     = "%s"
}

data "google_iam_policy" "foo" {
  binding {
    role = "%s"

    members = ["serviceAccount:${google_service_account.test_account.email}"]
  }
}

resource "google_spanner_database_iam_policy" "foo" {
  project     = google_spanner_database.database.project
  database    = google_spanner_database.database.name
  instance    = google_spanner_database.database.instance
  policy_data = data.google_iam_policy.foo.policy_data
}
`, account, instance, instance, database, roleId)
}
