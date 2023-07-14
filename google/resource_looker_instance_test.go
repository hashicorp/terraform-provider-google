// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccLookerInstance_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": RandString(t, 10),
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckLookerInstanceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccLookerInstance_lookerInstanceBasicExample(context),
			},
			{
				ResourceName:            "google_looker_instance.looker-instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oauth_config", "region"},
			},
			{
				Config: testAccLookerInstance_lookerInstanceFullExample(context),
			},
			{
				ResourceName:            "google_looker_instance.looker-instance",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"oauth_config", "region"},
			},
		},
	})
}
