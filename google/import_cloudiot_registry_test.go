package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccCloudIoTRegistry_import(t *testing.T) {
	t.Parallel()

	registryName := fmt.Sprintf("tf-test-registry-%d", acctest.RandInt())
	conf := fmt.Sprintf(`
		resource "google_cloudiot_registry" "registry-import-test" {
			project = "%s"
			region = "%s"
			name = "%s"
		}`, getTestProjectFromEnv(), DEFAULT_KMS_TEST_LOCATION, registryName)

	id := fmt.Sprintf("projects/%s/locations/%s/registries/%s",
		getTestProjectFromEnv(), DEFAULT_KMS_TEST_LOCATION, registryName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudIoTRegistryDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: conf,
			},
			resource.TestStep{
				ResourceName:      "google_cloudiot_registry.registry-import-test",
				ImportStateId:     id,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
