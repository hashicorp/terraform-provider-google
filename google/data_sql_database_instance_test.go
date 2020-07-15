package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccDataSourceSqlDatabaseInstance_basic(t *testing.T) {
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
				Config: testAccDataSourceSqlDatabaseInstance_basic(context),
				Check: resource.ComposeTestCheckFunc(
					checkDataSourceStateMatchesResourceState("data.google_sql_database_instance.qa", "google_sql_database_instance.master"),
				),
			},
		},
	})
}

func testAccDataSourceSqlDatabaseInstance_basic(context map[string]interface{}) string {
	return Nprintf(`
resource "google_sql_database_instance" "master" {
  name             = "master-instance-%{random_suffix}"
  database_version = "POSTGRES_11"
  region           = "us-central1"

  settings {
    # Second-generation instance tiers are based on the machine
    # type. See argument reference below.
    tier = "db-f1-micro"
  }
}

data "google_sql_database_instance" "qa" {
    name = google_sql_database_instance.master.name
}
`, context)
}
