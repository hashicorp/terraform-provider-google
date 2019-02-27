package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCloudIoTRegistry_basic(t *testing.T) {
	t.Parallel()

	registryName := fmt.Sprintf("psregistry-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudIoTRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIoTRegistry_basic(registryName),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudIoTRegistryExists(
						"google_cloudiot_registry.foobar"),
				),
			},
			{
				ResourceName:      "google_cloudiot_registry.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudIoTRegistry_extended(t *testing.T) {
	t.Parallel()

	registryName := fmt.Sprintf("psregistry-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudIoTRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIoTRegistry_extended(registryName),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudIoTRegistryExists(
						"google_cloudiot_registry.foobar"),
				),
			},
			{
				ResourceName:      "google_cloudiot_registry.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudIoTRegistry_update(t *testing.T) {
	t.Parallel()

	registryName := fmt.Sprintf("psregistry-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudIoTRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIoTRegistry_basic(registryName),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudIoTRegistryExists(
						"google_cloudiot_registry.foobar"),
				),
			},
			{
				Config: testAccCloudIoTRegistry_extended(registryName),
			},
			{
				Config: testAccCloudIoTRegistry_basic(registryName),
			},
			{
				ResourceName:      "google_cloudiot_registry.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCloudIoTRegistryDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_cloudiot_registry" {
			continue
		}
		config := testAccProvider.Meta().(*Config)
		registry, _ := config.clientCloudIoT.Projects.Locations.Registries.Get(rs.Primary.ID).Do()
		if registry != nil {
			return fmt.Errorf("Registry still present")
		}
	}
	return nil
}

func testAccCloudIoTRegistryExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*Config)
		_, err := config.clientCloudIoT.Projects.Locations.Registries.Get(rs.Primary.ID).Do()
		if err != nil {
			return fmt.Errorf("Registry does not exist")
		}
		return nil
	}
}

func testAccCloudIoTRegistry_basic(registryName string) string {
	return fmt.Sprintf(`
resource "google_cloudiot_registry" "foobar" {
	name = "%s"
}`, registryName)
}

func testAccCloudIoTRegistry_extended(registryName string) string {
	return fmt.Sprintf(`
resource "google_project_iam_binding" "cloud-iot-iam-binding" {
  members = ["serviceAccount:cloud-iot@system.gserviceaccount.com"]
  role    = "roles/pubsub.publisher"
}

resource "google_pubsub_topic" "default-devicestatus" {
  name = "psregistry-test-devicestatus-%s"
}

resource "google_pubsub_topic" "default-telemetry" {
  name = "psregistry-test-telemetry-%s"
}

resource "google_cloudiot_registry" "foobar" {
  depends_on = ["google_project_iam_binding.cloud-iot-iam-binding"]

  name = "%s"

  event_notification_config = {
    pubsub_topic_name = "${google_pubsub_topic.default-devicestatus.id}"
  }

  state_notification_config = {
    pubsub_topic_name = "${google_pubsub_topic.default-telemetry.id}"
  }

  http_config = {
    http_enabled_state = "HTTP_DISABLED"
  }

  mqtt_config = {
    mqtt_enabled_state = "MQTT_DISABLED"
  }

  credentials {
    public_key_certificate = {
      format      = "X509_CERTIFICATE_PEM"
      certificate = "${file("test-fixtures/rsa_cert.pem")}"
    }
  }
}
`, acctest.RandString(10), acctest.RandString(10), registryName)
}
