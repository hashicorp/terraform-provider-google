package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccComputeBackendBucket_import(t *testing.T) {
	t.Parallel()

	backendName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	storageName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeBackendBucketDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeBackendBucket_basic(backendName, storageName),
			},
			resource.TestStep{
				ResourceName:      "google_compute_backend_bucket.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
