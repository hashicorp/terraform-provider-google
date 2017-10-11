package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccSpannerInstance_importInstance(t *testing.T) {
	resourceName := "google_spanner_instance.basic"
	instanceName := fmt.Sprintf("span-itest-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpannerInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSpannerInstance_basic(instanceName),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportStateId:     instanceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSpannerInstance_importProjectInstance(t *testing.T) {
	resourceName := "google_spanner_instance.basic"
	instanceName := fmt.Sprintf("span-itest-%s", acctest.RandString(10))
	projectId := getTestProjectFromEnv()
	if projectId == "" {
		t.Skip("Unable to locate projectId via environment variables ... skipping ")
		return
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSpannerInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSpannerInstance_basicWithProject(projectId, instanceName),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
