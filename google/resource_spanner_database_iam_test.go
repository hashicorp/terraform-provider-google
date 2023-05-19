package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccSpannerDatabaseIamBinding(t *testing.T) {
	t.Parallel()

	account := fmt.Sprintf("tf-test-%d", RandInt(t))
	role := "roles/spanner.databaseAdmin"
	project := acctest.GetTestProjectFromEnv()
	database := fmt.Sprintf("tf-test-%s", RandString(t, 10))
	instance := fmt.Sprintf("tf-test-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
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
				}.TerraformId(), role),
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
				}.TerraformId(), role),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerDatabaseIamMember(t *testing.T) {
	t.Parallel()

	project := acctest.GetTestProjectFromEnv()
	account := fmt.Sprintf("tf-test-%d", RandInt(t))
	role := "roles/spanner.databaseAdmin"
	database := fmt.Sprintf("tf-test-%s", RandString(t, 10))
	instance := fmt.Sprintf("tf-test-%s", RandString(t, 10))
	conditionTitle := "Access only database one"

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccSpannerDatabaseIamMember_basic(account, instance, database, role),
			},
			{
				ResourceName: "google_spanner_database_iam_member.foo",
				ImportStateId: fmt.Sprintf("%s %s serviceAccount:%s@%s.iam.gserviceaccount.com %s", spannerDatabaseId{
					Instance: instance,
					Database: database,
					Project:  project,
				}.TerraformId(), role, account, project, conditionTitle),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerDatabaseIamPolicy(t *testing.T) {
	t.Parallel()

	project := acctest.GetTestProjectFromEnv()
	account := fmt.Sprintf("tf-test-%d", RandInt(t))
	role := "roles/spanner.databaseAdmin"
	database := fmt.Sprintf("tf-test-%s", RandString(t, 10))
	instance := fmt.Sprintf("tf-test-%s", RandString(t, 10))

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerDatabaseIamPolicy_basic(account, instance, database, role),
				Check:  resource.TestCheckResourceAttrSet("data.google_spanner_database_iam_policy.foo", "policy_data"),
			},
			// Test a few import formats
			{
				ResourceName: "google_spanner_database_iam_policy.foo",
				ImportStateId: spannerDatabaseId{
					Instance: instance,
					Database: database,
					Project:  project,
				}.TerraformId(),
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
  deletion_protection = false
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
  deletion_protection = false
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
  deletion_protection = false
}

resource "google_spanner_database_iam_member" "foo" {
  project  = google_spanner_database.database.project
  database = google_spanner_database.database.name
  instance = google_spanner_database.database.instance
  role     = "%s"
  member   = "serviceAccount:${google_service_account.test_account.email}"
  condition {
    title      = "Access only database one"
    expression = "resource.type == \"spanner.googleapis.com/DatabaseRole\" && resource.name.endsWith(\"/databaseRoles/parent\")"
  }
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
  deletion_protection = false
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

data "google_spanner_database_iam_policy" "foo" {
  project     = google_spanner_database.database.project
  database    = google_spanner_database.database.name
  instance    = google_spanner_database.database.instance
}
`, account, instance, instance, database, roleId)
}
