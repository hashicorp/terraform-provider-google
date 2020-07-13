package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCloudIoTDevice_update(t *testing.T) {
	t.Parallel()

	registryName := fmt.Sprintf("psregistry-test-%s", randString(t, 10))
	deviceName := fmt.Sprintf("psdevice-test-%s", randString(t, 10))
	resourceName := fmt.Sprintf("google_cloudiot_device.%s", deviceName)

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudIotDeviceDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIoTDeviceBasic(deviceName, registryName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudIoTDeviceExtended(deviceName, registryName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudIoTDeviceBasic(deviceName, registryName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCloudIoTDeviceBasic(deviceName string, registryName string) string {
	return fmt.Sprintf(`

resource "google_cloudiot_registry" "%s" {
  name = "%s"
}

resource "google_cloudiot_device" "%s" {
  name     = "%s"
  registry = google_cloudiot_registry.%s.id

  gateway_config {
    gateway_auth_method = "DEVICE_AUTH_TOKEN_ONLY"
    gateway_type = "GATEWAY"
  }
}


`, registryName, registryName, deviceName, deviceName, registryName)
}

func testAccCloudIoTDeviceExtended(deviceName string, registryName string) string {
	return fmt.Sprintf(`

resource "google_cloudiot_registry" "%s" {
  name = "%s"
}

resource "google_cloudiot_device" "%s" {
  name     = "%s"
  registry = google_cloudiot_registry.%s.id

  credentials {
    public_key {
      format = "RSA_PEM"
      key = file("test-fixtures/rsa_public.pem")
    }
  }

  blocked = false

  log_level = "INFO"

  metadata = {
    test_key_1 = "test_value_1"
  }

  gateway_config {
    gateway_auth_method = "ASSOCIATION_AND_DEVICE_AUTH_TOKEN"
    gateway_type = "GATEWAY"
  }
}
`, registryName, registryName, deviceName, deviceName, registryName)
}
