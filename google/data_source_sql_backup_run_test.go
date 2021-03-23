package google

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSqlBackupRun_basic(t *testing.T) {
	// Sqladmin client
	skipIfVcr(t)
	t.Parallel()

	instance := BootstrapSharedSQLInstanceBackupRun(t)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSqlBackupRun_basic(instance),
				Check:  resource.TestMatchResourceAttr("data.google_sql_backup_run.backup", "status", regexp.MustCompile("SUCCESSFUL")),
			},
		},
	})
}

func TestAccDataSourceSqlBackupRun_notFound(t *testing.T) {
	// Sqladmin client
	skipIfVcr(t)
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": randString(t, 10),
	}

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccSqlDatabaseInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceSqlBackupRun_notFound(context),
				ExpectError: regexp.MustCompile("No backups found for SQL Database Instance"),
			},
		},
	})
}

func testAccDataSourceSqlBackupRun_basic(instance string) string {
	return fmt.Sprintf(`
data "google_sql_backup_run" "backup" {
	instance = "%s"
	most_recent = true
}
`, instance)
}

func testAccDataSourceSqlBackupRun_notFound(context map[string]interface{}) string {
	return Nprintf(`
resource "google_sql_database_instance" "instance" {
  name             = "tf-test-instance-%{random_suffix}"
  database_version = "POSTGRES_11"
  region           = "us-central1"

  settings {
	tier = "db-f1-micro"
	backup_configuration {
		enabled            = "false"
	}
  }

  deletion_protection = false
}

data "google_sql_backup_run" "backup" {
	instance = google_sql_database_instance.instance.name
	most_recent = true
	depends_on = [google_sql_database_instance.instance]
}
`, context)
}
