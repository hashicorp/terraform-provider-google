package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccSpannerDatabase_importInstanceDatabase(t *testing.T) {
	resourceName := "google_spanner_database.basic"
	instanceName := fmt.Sprintf("span-iname-%s", acctest.RandString(10))
	dbName := fmt.Sprintf("span-dbname-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpannerDatabaseDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSpannerDatabase_basicImport(instanceName, dbName),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportStateId:     instanceName + "/" + dbName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerDatabase_importProjectInstanceDatabase(t *testing.T) {
	resourceName := "google_spanner_database.basic"
	instanceName := fmt.Sprintf("span-iname-%s", acctest.RandString(10))
	dbName := fmt.Sprintf("span-dbname-%s", acctest.RandString(10))
	projectId := getTestProjectFromEnv()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpannerDatabaseDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSpannerDatabase_basicImportWithProject(projectId, instanceName, dbName),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
