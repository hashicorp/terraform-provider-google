package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSqlUser_mysql(t *testing.T) {
	// Multiple fine-grained resources
	skipIfVcr(t)
	t.Parallel()

	instance := fmt.Sprintf("tf-test-%d", randInt(t))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_mysql(instance, "password"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user1"),
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user2"),
				),
			},
			{
				// Update password
				Config: testGoogleSqlUser_mysql(instance, "new_password"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user1"),
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user2"),
				),
			},
			{
				ResourceName:            "google_sql_user.user2",
				ImportStateId:           fmt.Sprintf("%s/%s/gmail.com/admin", getTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccSqlUser_iamUser(t *testing.T) {
	// Multiple fine-grained resources
	skipIfVcr(t)
	t.Parallel()

	instance := fmt.Sprintf("tf-test-%d", randInt(t))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_iamUser(instance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user1"),
				),
			},
			{
				ResourceName:      "google_sql_user.user1",
				ImportStateId:     fmt.Sprintf("%s/%s/%%/%s@%s.iam.gserviceaccount.com", getTestProjectFromEnv(), instance, instance, getTestProjectFromEnv()),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSqlUser_postgres(t *testing.T) {
	t.Parallel()

	instance := fmt.Sprintf("tf-test-%d", randInt(t))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_postgres(instance, "password"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user"),
				),
			},
			{
				// Update password
				Config: testGoogleSqlUser_postgres(instance, "new_password"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user"),
				),
			},
			{
				ResourceName:            "google_sql_user.user",
				ImportStateId:           fmt.Sprintf("%s/%s/admin", getTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccSqlUser_postgresIAM(t *testing.T) {
	t.Parallel()

	instance := fmt.Sprintf("tf-test-%d", randInt(t))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_postgresIAM(instance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user"),
				),
			},
			{
				ResourceName:            "google_sql_user.user",
				ImportStateId:           fmt.Sprintf("%s/%s/admin", getTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccSqlUser_postgresAbandon(t *testing.T) {
	t.Parallel()

	instance := fmt.Sprintf("tf-test-%d", randInt(t))
	userName := "admin"
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_postgresAbandon(instance, userName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user"),
				),
			},
			{
				ResourceName:            "google_sql_user.user",
				ImportStateId:           fmt.Sprintf("%s/%s/admin", getTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "deletion_policy"},
			},
			{
				// Abandon user
				Config: testGoogleSqlUser_postgresNoUser(instance),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExistsWithName(t, instance, userName),
				),
			},
		},
	})
}

func testAccCheckGoogleSqlUserExists(t *testing.T, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		name := rs.Primary.Attributes["name"]
		instance := rs.Primary.Attributes["instance"]
		host := rs.Primary.Attributes["host"]
		users, err := config.NewSqlAdminClient(config.userAgent).Users.List(config.Project,
			instance).Do()

		if err != nil {
			return err
		}

		for _, user := range users.Items {
			if user.Name == name && user.Host == host {
				return nil
			}
		}

		return fmt.Errorf("Not found: %s: %s", n, err)
	}
}

func testAccCheckGoogleSqlUserExistsWithName(t *testing.T, instance, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

		users, err := config.NewSqlAdminClient(config.userAgent).Users.List(config.Project,
			instance).Do()

		if err != nil {
			return err
		}

		for _, user := range users.Items {
			if user.Name == name {
				return nil
			}
		}

		return fmt.Errorf("Not found: User: %s in instance: %s: %s", name, instance, err)
	}
}

func testAccSqlUserDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			config := googleProviderConfig(t)
			if rs.Type != "google_sql_database" {
				continue
			}

			name := rs.Primary.Attributes["name"]
			instance := rs.Primary.Attributes["instance"]
			host := rs.Primary.Attributes["host"]
			users, err := config.NewSqlAdminClient(config.userAgent).Users.List(config.Project,
				instance).Do()

			if users == nil {
				return nil
			}

			for _, user := range users.Items {
				if user.Name == name && user.Host == host {
					return fmt.Errorf("User still %s exists %s", name, err)
				}
			}

			return nil
		}

		return nil
	}
}

func TestAccSqlUser_mysqlPasswordPolicy(t *testing.T) {
	// Multiple fine-grained resources
	skipIfVcr(t)
	t.Parallel()

	instance := fmt.Sprintf("tf-test-i%d", randInt(t))
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlUserDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testGoogleSqlUser_mysqlPasswordPolicy(instance, "password"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user1"),
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user2"),
				),
			},
			{
				// Update password
				Config: testGoogleSqlUser_mysqlPasswordPolicy(instance, "new_password"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user1"),
					testAccCheckGoogleSqlUserExists(t, "google_sql_user.user2"),
				),
			},
			{
				ResourceName:            "google_sql_user.user2",
				ImportStateId:           fmt.Sprintf("%s/%s/gmail.com/admin", getTestProjectFromEnv(), instance),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testGoogleSqlUser_mysql(instance, password string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_5_7"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_user" "user1" {
  name     = "admin"
  instance = google_sql_database_instance.instance.name
  host     = "google.com"
  password = "%s"
}

resource "google_sql_user" "user2" {
  name     = "admin"
  instance = google_sql_database_instance.instance.name
  host     = "gmail.com"
  password = "hunter2"
}
`, instance, password)
}

func testGoogleSqlUser_mysqlPasswordPolicy(instance, password string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name                = "%s"
  region              = "us-central1"
  database_version    = "MYSQL_8_0"
  deletion_protection = false
  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_user" "user1" {
  name     = "admin"
  instance = google_sql_database_instance.instance.name
  host     = "google.com"
  password = "%s"

  password_policy {
    allowed_failed_attempts  = 6
    password_expiration_duration  =  "2592000s"
    enable_failed_attempts_check = true
    enable_password_verification = true
  }
}

resource "google_sql_user" "user2" {
  name     = "admin"
  instance = google_sql_database_instance.instance.name
  host     = "gmail.com"
  password = "hunter2"
  password_policy {
    allowed_failed_attempts  = 6
    enable_failed_attempts_check = true
  }
}
`, instance, password)
}

func testGoogleSqlUser_postgres(instance, password string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-central1"
  database_version = "POSTGRES_9_6"
  deletion_protection = false

  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_user" "user" {
  name     = "admin"
  instance = google_sql_database_instance.instance.name
  password = "%s"
}
`, instance, password)
}

func testGoogleSqlUser_postgresIAM(instance string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-central1"
  database_version = "POSTGRES_9_6"
  deletion_protection = false

  settings {
    tier = "db-f1-micro"
    database_flags {
      name  = "cloudsql.iam_authentication"
      value = "on"
    }
  }
}

resource "google_sql_user" "user" {
  name     = "admin"
  instance = google_sql_database_instance.instance.name
  type     = "CLOUD_IAM_USER"
}
`, instance)
}

func testGoogleSqlUser_postgresAbandon(instance, name string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-central1"
  database_version = "POSTGRES_9_6"
  deletion_protection = false

  settings {
    tier = "db-f1-micro"
  }
}

resource "google_sql_user" "user" {
  name     = "%s"
  instance = google_sql_database_instance.instance.name
  password = "password"
  deletion_policy = "ABANDON"
}
`, instance, name)
}

func testGoogleSqlUser_postgresNoUser(instance string) string {
	return fmt.Sprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "%s"
  region           = "us-central1"
  database_version = "POSTGRES_9_6"
  deletion_protection = false

  settings {
    tier = "db-f1-micro"
  }
}
`, instance)
}

func testGoogleSqlUser_iamUser(instance string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_sql_database_instance" "instance" {
  database_version = "MYSQL_8_0"
  name             = "%s"
  region           = "us-central1"

  settings {
    tier              = "db-f1-micro"
    availability_type = "REGIONAL"

    backup_configuration {
      enabled            = true
      binary_log_enabled = true
    }

    database_flags {
      name  = "cloudsql_iam_authentication"
      value = "on"
    }
  }

  deletion_protection = false
}

resource "google_sql_database" "db" {
  name     = "%s"
  instance = google_sql_database_instance.instance.name
}

resource "google_service_account" "sa" {
  account_id   = "%s"
  display_name = "%s"
}

resource "google_service_account_key" "sa_key" {
  service_account_id = google_service_account.sa.email
}

resource "google_sql_user" "user1" {
  name     = google_service_account.sa.email
  instance = google_sql_database_instance.instance.name
  type     = "CLOUD_IAM_SERVICE_ACCOUNT"
}

resource "google_project_iam_member" "instance_user" {
  project = data.google_project.project.project_id
  role    = "roles/cloudsql.instanceUser"
  member  = "serviceAccount:${google_service_account.sa.email}"
}

resource "google_project_iam_member" "sa_user" {
  project = data.google_project.project.project_id
  role    = "roles/iam.serviceAccountUser"
  member  = "serviceAccount:${google_service_account.sa.email}"
}
`, instance, instance, instance, instance)
}
