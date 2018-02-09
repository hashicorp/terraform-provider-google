package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccGoogleSqlUser_firstGen(t *testing.T) {
	t.Parallel()

	instance := acctest.RandomWithPrefix("i")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlUserDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleSqlUser_firstGen(instance, "password"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists("google_sql_user.user1"),
					testAccCheckGoogleSqlUserExists("google_sql_user.user2"),
				),
			},
			resource.TestStep{
				// Update password
				Config: testGoogleSqlUser_firstGen(instance, "new_password"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists("google_sql_user.user1"),
					testAccCheckGoogleSqlUserExists("google_sql_user.user2"),
				),
			},
			resource.TestStep{
				ResourceName:            "google_sql_user.user2",
				ImportStateId:           instance + "/gmail.com/admin",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccGoogleSqlUser_secondGen(t *testing.T) {
	t.Parallel()

	instance := acctest.RandomWithPrefix("i")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleSqlUserDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleSqlUser_secondGen(instance, "password"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists("google_sql_user.user"),
				),
			},
			resource.TestStep{
				// Update password
				Config: testGoogleSqlUser_secondGen(instance, "new_password"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGoogleSqlUserExists("google_sql_user.user"),
				),
			},
			resource.TestStep{
				ResourceName:            "google_sql_user.user",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccCheckGoogleSqlUserExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		name := rs.Primary.Attributes["name"]
		instance := rs.Primary.Attributes["instance"]
		host := rs.Primary.Attributes["host"]
		users, err := config.clientSqlAdmin.Users.List(config.Project,
			instance).Do()

		for _, user := range users.Items {
			if user.Name == name && user.Host == host {
				return nil
			}
		}

		return fmt.Errorf("Not found: %s: %s", n, err)
	}
}

func testAccGoogleSqlUserDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		config := testAccProvider.Meta().(*Config)
		if rs.Type != "google_sql_database" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		instance := rs.Primary.Attributes["instance"]
		host := rs.Primary.Attributes["host"]
		users, err := config.clientSqlAdmin.Users.List(config.Project,
			instance).Do()

		for _, user := range users.Items {
			if user.Name == name && user.Host == host {
				return fmt.Errorf("User still %s exists %s", name, err)
			}
		}

		return nil
	}

	return nil
}

func testGoogleSqlUser_firstGen(instance, password string) string {
	return fmt.Sprintf(`
	resource "google_sql_database_instance" "instance" {
		name = "%s"
		region = "us-central"
		settings {
			tier = "D0"
		}
	}

	resource "google_sql_user" "user1" {
		name = "admin"
		instance = "${google_sql_database_instance.instance.name}"
		host = "google.com"
		password = "%s"
	}

	resource "google_sql_user" "user2" {
		name = "admin"
		instance = "${google_sql_database_instance.instance.name}"
		host = "gmail.com"
		password = "hunter2"
	}
	`, instance, password)
}

func testGoogleSqlUser_secondGen(instance, password string) string {
	return fmt.Sprintf(`
	resource "google_sql_database_instance" "instance" {
		name = "%s"
		region = "us-central1"
		database_version = "POSTGRES_9_6"

		settings {
			tier = "db-f1-micro"
		}
	}

	resource "google_sql_user" "user" {
		name = "admin"
		instance = "${google_sql_database_instance.instance.name}"
		password = "%s"
	}
	`, instance, password)
}
