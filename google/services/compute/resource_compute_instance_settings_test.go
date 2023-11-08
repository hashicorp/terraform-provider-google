// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute_test

import (
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func testAccComputeInstanceSettings_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`

resource "google_compute_instance_settings" "gce_instance_settings" {
  provider = google-beta
  zone = "us-east7-b"
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
  provider = google-beta
  zone = "us-east7-b"
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
  provider = google-beta
  zone = "us-east7-b"
  metadata {
    items = {
      baz = "qux"
    }
  }
}

`, context)
}
