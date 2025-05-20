// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccComputeInstanceSettings_update(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeInstanceSettingsDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeInstanceSettings_basic(context),
			},
			{
				ResourceName:            "google_compute_instance_settings.gce_instance_settings",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "fingerprint"},
			},
			{
				Config: testAccComputeInstanceSettings_update(context),
			},
			{
				ResourceName:            "google_compute_instance_settings.gce_instance_settings",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "fingerprint"},
			},
			{
				Config: testAccComputeInstanceSettings_delete(context),
			},
			{
				ResourceName:            "google_compute_instance_settings.gce_instance_settings",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"zone", "fingerprint"},
			},
		},
	})
}

func testAccComputeInstanceSettings_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_compute_instance_settings" "gce_instance_settings" {
  zone = "us-east5-c"
  metadata {
    items = {
      foo = "baz"
    }
  }
}

`, context)
}

func testAccComputeInstanceSettings_update(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_compute_instance_settings" "gce_instance_settings" {
  zone = "us-east5-c"
  metadata {
    items = {
      foo = "bar"
      baz = "qux"
    }
  }
}

`, context)
}

func testAccComputeInstanceSettings_delete(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_compute_instance_settings" "gce_instance_settings" {
  zone = "us-east5-c"
  metadata {
    items = {
      baz = "qux"
    }
  }
}

`, context)
}
