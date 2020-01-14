package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestValidateCloudIoTID(t *testing.T) {
	x := []StringValidationTestCase{
		// No errors
		{TestName: "basic", Value: "foobar"},
		{TestName: "with numbers", Value: "foobar123"},
		{TestName: "short", Value: "foo"},
		{TestName: "long", Value: "foobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoobarfoo"},
		{TestName: "has a hyphen", Value: "foo-bar"},

		// With errors
		{TestName: "empty", Value: "", ExpectError: true},
		{TestName: "starts with a goog", Value: "googfoobar", ExpectError: true},
		{TestName: "starts with a number", Value: "1foobar", ExpectError: true},
		{TestName: "has an slash", Value: "foo/bar", ExpectError: true},
		{TestName: "has an backslash", Value: "foo\bar", ExpectError: true},
		{TestName: "too long", Value: strings.Repeat("f", 260), ExpectError: true},
	}

	es := testStringValidationCases(x, validateCloudIotID)
	if len(es) > 0 {
		t.Errorf("Failed to validate CloudIoT ID names: %v", es)
	}
}

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
			},
			{
				ResourceName:      "google_cloudiot_registry.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudIoTRegistry_extended(registryName),
			},
			{
				ResourceName:      "google_cloudiot_registry.foobar",
				ImportState:       true,
				ImportStateVerify: true,
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

func TestAccCloudIoTRegistry_eventNotificationConfigsSingle(t *testing.T) {
	t.Parallel()

	registryName := fmt.Sprintf("tf-registry-test-%s", acctest.RandString(10))
	topic := fmt.Sprintf("tf-registry-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudIoTRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIoTRegistry_singleEventNotificationConfigs(topic, registryName),
			},
			{
				ResourceName:      "google_cloudiot_registry.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudIoTRegistry_eventNotificationConfigsMultiple(t *testing.T) {
	t.Parallel()

	registryName := fmt.Sprintf("tf-registry-test-%s", acctest.RandString(10))
	topic := fmt.Sprintf("tf-registry-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudIoTRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudIoTRegistry_multipleEventNotificationConfigs(topic, registryName),
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

func testAccCloudIoTRegistry_basic(registryName string) string {
	return fmt.Sprintf(`
resource "google_cloudiot_registry" "foobar" {
  name = "%s"
}
`, registryName)
}

func testAccCloudIoTRegistry_extended(registryName string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "default-devicestatus" {
  name = "psregistry-test-devicestatus-%s"
}

resource "google_pubsub_topic" "default-telemetry" {
  name = "psregistry-test-telemetry-%s"
}

resource "google_cloudiot_registry" "foobar" {
  name = "%s"

  event_notification_configs {
    pubsub_topic_name = google_pubsub_topic.default-devicestatus.id
  }

  state_notification_config = {
    pubsub_topic_name = google_pubsub_topic.default-telemetry.id
  }

  http_config = {
    http_enabled_state = "HTTP_DISABLED"
  }

  mqtt_config = {
    mqtt_enabled_state = "MQTT_DISABLED"
  }

  log_level = "INFO"

  credentials {
    public_key_certificate = {
      format      = "X509_CERTIFICATE_PEM"
      certificate = file("test-fixtures/rsa_cert.pem")
    }
  }
}
`, acctest.RandString(10), acctest.RandString(10), registryName)
}

func testAccCloudIoTRegistry_singleEventNotificationConfigs(topic, registryName string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "event-topic-1" {
  name = "%s"
}

resource "google_cloudiot_registry" "foobar" {
  name = "%s"

  event_notification_configs {
    pubsub_topic_name = google_pubsub_topic.event-topic-1.id
    subfolder_matches = ""
  }
}
`, topic, registryName)
}

func testAccCloudIoTRegistry_multipleEventNotificationConfigs(topic, registryName string) string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "event-topic-1" {
  name = "%s"
}

resource "google_pubsub_topic" "event-topic-2" {
  name = "%s-alt"
}

resource "google_cloudiot_registry" "foobar" {
  name = "%s"

  event_notification_configs {
    pubsub_topic_name = google_pubsub_topic.event-topic-1.id
    subfolder_matches = "test"
  }

  event_notification_configs {
    pubsub_topic_name = google_pubsub_topic.event-topic-2.id
    subfolder_matches = ""
  }
}
`, topic, topic, registryName)
}
