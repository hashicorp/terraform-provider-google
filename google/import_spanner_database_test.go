package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccGoogleSpannerDatabase_importBasic(t *testing.T) {
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
				ResourceName:        resourceName,
				ImportStateIdPrefix: instanceName + "/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}

func TestAccGoogleSpannerDatabase_importWithProject(t *testing.T) {
	resourceName := "google_spanner_database.basic"
	instanceName := fmt.Sprintf("span-iname-%s", acctest.RandString(10))
	dbName := fmt.Sprintf("span-dbname-%s", acctest.RandString(10))
	var projectId = multiEnvSearch([]string{"GOOGLE_PROJECT", "GCLOUD_PROJECT", "CLOUDSDK_CORE_PROJECT"})

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpannerDatabaseDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSpannerDatabase_basicImportWithProject(projectId, instanceName, dbName),
			},

			resource.TestStep{
				ResourceName:        resourceName,
				ImportStateIdPrefix: projectId + "/" + instanceName + "/",
				ImportState:         true,
				ImportStateVerify:   true,
			},
		},
	})
}
