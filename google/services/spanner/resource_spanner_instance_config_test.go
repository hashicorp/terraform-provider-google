// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package spanner_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Acceptance Tests

func TestAccSpannerInstanceConfig_update(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("custom-tf-test-config-%s", acctest.RandString(t, 10))
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckSpannerInstanceConfigDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccSpannerInstanceConfig_update(name, "display name", false),
			},
			{
				ResourceName:            "google_spanner_instance_config.updater",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
			{
				Config: testAccSpannerInstanceConfig_update(name, "display name updated", true),
			},
			{
				ResourceName:            "google_spanner_instance_config.updater",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"labels", "terraform_labels"},
			},
		},
	})
}

func testAccSpannerInstanceConfig_update(name, displayName string, addLabel bool) string {
	extraLabel := ""
	if addLabel {
		extraLabel = "\"key2\" = \"value2\""
	}
	return fmt.Sprintf(`
resource "google_spanner_instance_config" "updater" {
  name          = "%s"
  display_name  = "%s-dname"
  base_config    = "nam11"
  replicas     {
      location = "us-west1"
      type = "READ_ONLY"
 }
 labels = {
     "key1" = "value1"
     %s
   }
}
`, name, displayName, extraLabel)
}
