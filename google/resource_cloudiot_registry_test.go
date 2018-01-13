package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCloudiotRegistryCreate_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudiotRegistryDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCloudiotRegistry_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudiotRegistryExists(
						"google_cloudiot_registry.foobar"),
				),
			},
		},
	})
}

func TestAccCloudiotRegistryCreate_extended(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudiotRegistryDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCloudiotRegistry_extended(),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudiotRegistryExists(
						"google_cloudiot_registry.foobar-extended"),
				),
			},
		},
	})
}

func testAccCheckCloudiotRegistryDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_cloudiot_registry" {
			continue
		}

		config := testAccProvider.Meta().(*Config)
		registry, _ := config.clientCloudiot.Projects.Locations.Registries.Get(rs.Primary.ID).Do()
		if registry != nil {
			return fmt.Errorf("Registry still present")
		}
	}

	return nil
}

func testAccCloudiotRegistryExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*Config)
		_, err := config.clientCloudiot.Projects.Locations.Registries.Get(rs.Primary.ID).Do()
		if err != nil {
			return fmt.Errorf("Topic does not exist")
		}

		return nil
	}
}

func testAccCloudiotRegistry_basic() string {
	return fmt.Sprintf(`
resource "google_cloudiot_registry" "foobar" {
	name = "psregistry-test-%s"
}`, acctest.RandString(10))
}

func testAccCloudiotRegistry_extended() string {
	return fmt.Sprintf(`
resource "google_pubsub_topic" "default-devicestatus" {
  name = "psregistry-test-devicestatus-%s"
}

resource "google_pubsub_topic" "default-telemetry" {
  name = "psregistry-test-telemetry-%s"
}

resource "google_cloudiot_registry" "foobar-extended" {
  name = "psregistry-test-%s"

  event_notification_configs = [{
    pubsub_topic_name = "${google_pubsub_topic.default-devicestatus.id}"
  }]

  state_notification_config = {
    pubsub_topic_name = "${google_pubsub_topic.default-telemetry.id}"
  }

  http_config = {
    http_enabled_state = "HTTP_DISABLED"
  }

  mqtt_config = {
    mqtt_enabled_state = "MQTT_ENABLED"
  }

  credentials = [
    {
      format      = "X509_CERTIFICATE_PEM"
      certificate = "-----BEGIN CERTIFICATE-----\nMIICnjCCAYYCCQC/5gx7LgJFqTANBgkqhkiG9w0BAQsFADARMQ8wDQYDVQQDDAZ1\nbnVzZWQwHhcNMTgwMTEyMjAxMzQzWhcNMjMwMTExMjAxMzQzWjARMQ8wDQYDVQQD\nDAZ1bnVzZWQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDStfQvJzmN\nCYLSWpwvTmyCKn8t19cfWZ69wFaB3OSglxXgYe3w9An0QHybDpKITt61PpfsKov3\nEcnzH5IA+Ox+4jUppBL1mSkO/BWig+sd1dG7pQMbGi4nMxW704A0PRUaNIOarOlR\nrNUJZQrsghkMjLayCTJ2HISBBiPnKKB3f3KCc9sDhj2Z7zy7HfeW0apZ1m6xAQCC\neSZNW0IyGIYKTd9F7HEJFzOWg9JHvabbciBEcFWKGVzM8nQr1q8KU8Xi3iN2mpNK\nJkbRLNnqKhvjPyIZ4s4cDSEZN1/OaGQ4XP2mvU03+4UAoMPoJ8IczBKTB0mFxfX8\nlDZZa5IWU9sNAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAHnkTIghRj/cerR9ctji\nkancnjlsdNEuPiVpMj+SOtOH8cvlgl0oWG6segYTVzk4VEHlq3POB67Yjoz829XM\nCEgUxSqGvDrQ7IaPLPryYy8o5azMLnEZDr+Yd6CUKr/pUZzJoZxHj7z3iqeQZnMW\nS6kb6HYvG5PKlJ7+JUIKLou0RQmaM9BQ0Nln/YDRRIerD0MY9k7No2ZEDbywZqQK\nGRIqT+BlN84oHOR44h2RqWhn9O50tkbcmAIKgmeg/mxwmeAm/6hQ8VrOhDHqsFdT\nzh2l6IeCl8EF8MjNrFRcQx21TTqeU6vGIPgM3E0k8PQUc+s+lir8UFsIzKaOFsIh\nuKU=\n-----END CERTIFICATE-----\n"

      x509_details = {
        issuer              = "CN=unused"
        subject             = "CN=unused"
        start_time          = "2018-01-12T20:13:43Z"
        expiry_time         = "2023-01-11T20:13:43Z"
        signature_algorithm = "sha256WithRSAEncryption"
        public_key_type     = "PK_RSA"
      }
    },
  ]
}

`, acctest.RandString(10), acctest.RandString(10), acctest.RandString(10))
}
