// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudiot_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestAccCloudIoTRegistry_update(t *testing.T) {
	t.Parallel()

	registryName := fmt.Sprintf("psregistry-test-%s", acctest.RandString(t, 10))
	resourceName := fmt.Sprintf("google_cloudiot_registry.%s", registryName)
	deviceStatus := fmt.Sprintf("psregistry-test-devicestatus-%s", acctest.RandString(t, 10))
	defaultTelemetry := fmt.Sprintf("psregistry-test-telemetry-%s", acctest.RandString(t, 10))
	additionalTelemetry := fmt.Sprintf("psregistry-additional-test-telemetry-%s", acctest.RandString(t, 10))

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckCloudIotDeviceRegistryDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIoTRegistryBasic(registryName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudIoTRegistryExtended(registryName, deviceStatus, defaultTelemetry, additionalTelemetry),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudIoTRegistryBasic(registryName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCloudIoTRegistryBasic(registryName string) string {
	return fmt.Sprintf(`

resource "google_cloudiot_registry" "%s" {
  name = "%s"
}
`, registryName, registryName)
}

func testAccCloudIoTRegistryExtended(registryName string, deviceStatus string, defaultTelemetry string, additionalTelemetry string) string {
	return fmt.Sprintf(`

resource "google_pubsub_topic" "default-devicestatus" {
  name = "psregistry-test-devicestatus-%s"
}

resource "google_pubsub_topic" "default-telemetry" {
  name = "psregistry-test-telemetry-%s"
}

resource "google_pubsub_topic" "additional-telemetry" {
  name = "psregistry-additional-test-telemetry-%s"
}

resource "google_cloudiot_registry" "%s" {
  name     = "%s"

  event_notification_configs {
    pubsub_topic_name = google_pubsub_topic.additional-telemetry.id
    subfolder_matches = "test/directory"
  }

  event_notification_configs {
    pubsub_topic_name = google_pubsub_topic.default-telemetry.id
    subfolder_matches = ""
  }

  state_notification_config = {
    pubsub_topic_name = google_pubsub_topic.default-devicestatus.id
  }

  mqtt_config = {
    mqtt_enabled_state = "MQTT_DISABLED"
  }

  http_config = {
    http_enabled_state = "HTTP_DISABLED"
  }

  credentials {
    public_key_certificate = {
      format      = "X509_CERTIFICATE_PEM"
      certificate = file("test-fixtures/rsa_cert.pem")
    }
  }
}
`, deviceStatus, defaultTelemetry, additionalTelemetry, registryName, registryName)
}
