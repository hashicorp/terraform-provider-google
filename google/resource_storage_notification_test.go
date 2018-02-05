package google

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/storage/v1"
)

var (
	gcsServiceAccount = fmt.Sprintf("serviceAccount:%s@gs-project-accounts.iam.gserviceaccount.com", os.Getenv("GOOGLE_PROJECT"))
	payload           = "JSON_API_V1"
)

func TestAccGoogleStorageNotification_basic(t *testing.T) {
	t.Parallel()

	skipIfEnvNotSet(t, "GOOGLE_PROJECT")

	var notification storage.Notification
	bucketName := testBucketName()
	topicName := fmt.Sprintf("tf-pstopic-test-%d", acctest.RandInt())
	topic := fmt.Sprintf("//pubsub.googleapis.com/projects/%s/topics/%s", os.Getenv("GOOGLE_PROJECT"), topicName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleStorageNotificationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleStorageNotificationBasic(bucketName, topicName, topic),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageNotificationExists(
						"google_storage_notification.notification", &notification),
					resource.TestCheckResourceAttr(
						"google_storage_notification.notification", "bucket", bucketName),
					resource.TestCheckResourceAttr(
						"google_storage_notification.notification", "topic", topic),
					resource.TestCheckResourceAttr(
						"google_storage_notification.notification", "payload_format", payload),
					resource.TestCheckResourceAttr(
						"google_storage_notification.notification_with_prefix", "object_name_prefix", "foobar"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_storage_notification.notification",
				ImportState:       true,
				ImportStateVerify: true,
			},
			resource.TestStep{
				ResourceName:      "google_storage_notification.notification_with_prefix",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGoogleStorageNotification_withEventsAndAttributes(t *testing.T) {
	t.Parallel()

	skipIfEnvNotSet(t, "GOOGLE_PROJECT")

	var notification storage.Notification
	bucketName := testBucketName()
	topicName := fmt.Sprintf("tf-pstopic-test-%d", acctest.RandInt())
	topic := fmt.Sprintf("//pubsub.googleapis.com/projects/%s/topics/%s", os.Getenv("GOOGLE_PROJECT"), topicName)
	eventType1 := "OBJECT_FINALIZE"
	eventType2 := "OBJECT_ARCHIVE"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccGoogleStorageNotificationDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testGoogleStorageNotificationOptionalEventsAttributes(bucketName, topicName, topic, eventType1, eventType2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStorageNotificationExists(
						"google_storage_notification.notification", &notification),
					resource.TestCheckResourceAttr(
						"google_storage_notification.notification", "bucket", bucketName),
					resource.TestCheckResourceAttr(
						"google_storage_notification.notification", "topic", topic),
					resource.TestCheckResourceAttr(
						"google_storage_notification.notification", "payload_format", payload),
					testAccCheckStorageNotificationCheckEventType(
						&notification, []string{eventType1, eventType2}),
					testAccCheckStorageNotificationCheckAttributes(
						&notification, "new-attribute", "new-attribute-value"),
				),
			},
			resource.TestStep{
				ResourceName:      "google_storage_notification.notification",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGoogleStorageNotificationDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_storage_notification" {
			continue
		}

		bucket, notificationID := resourceStorageNotificationParseID(rs.Primary.ID)

		_, err := config.clientStorage.Notifications.Get(bucket, notificationID).Do()
		if err == nil {
			return fmt.Errorf("Notification configuration still exists")
		}
	}

	return nil
}

func testAccCheckStorageNotificationExists(resource string, notification *storage.Notification) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		bucket, notificationID := resourceStorageNotificationParseID(rs.Primary.ID)

		found, err := config.clientStorage.Notifications.Get(bucket, notificationID).Do()
		if err != nil {
			return err
		}

		if found.Id != notificationID {
			return fmt.Errorf("Storage notification configuration not found")
		}

		*notification = *found

		return nil
	}
}

func testAccCheckStorageNotificationCheckEventType(notification *storage.Notification, eventTypes []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !reflect.DeepEqual(notification.EventTypes, eventTypes) {
			return fmt.Errorf("Target event types are incorrect. Expected %s, got %s", eventTypes, notification.EventTypes)
		}
		return nil
	}
}

func testAccCheckStorageNotificationCheckAttributes(notification *storage.Notification, key, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		val, ok := notification.CustomAttributes[key]
		if !ok {
			return fmt.Errorf("Custom attribute with key %s not found", key)
		}

		if val != value {
			return fmt.Errorf("Custom attribute value did not match for key %s: expected %s but found %s", key, value, val)
		}
		return nil
	}
}

func testGoogleStorageNotificationBasic(bucketName, topicName, topic string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
}
		
resource "google_pubsub_topic" "topic" {
	name = "%s"
}
// We have to provide GCS default storage account with the permission
// to publish to a Cloud Pub/Sub topic from this project
// Otherwise notification configuration won't work
resource "google_pubsub_topic_iam_binding" "binding" {
	topic   = "${google_pubsub_topic.topic.name}"
	role    = "roles/pubsub.publisher"
		  
	members = ["%s"]
}

resource "google_storage_notification" "notification" {
	bucket         = "${google_storage_bucket.bucket.name}"
	payload_format = "JSON_API_V1"
	topic          = "${google_pubsub_topic.topic.id}"
	depends_on     = ["google_pubsub_topic_iam_binding.binding"]
}

resource "google_storage_notification" "notification_with_prefix" {
	bucket             = "${google_storage_bucket.bucket.name}"
	payload_format     = "JSON_API_V1"
	topic              = "${google_pubsub_topic.topic.id}"
	object_name_prefix = "foobar"
	depends_on         = ["google_pubsub_topic_iam_binding.binding"]
}

`, bucketName, topicName, gcsServiceAccount)
}

func testGoogleStorageNotificationOptionalEventsAttributes(bucketName, topicName, topic, eventType1, eventType2 string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
}
		
resource "google_pubsub_topic" "topic" {
	name = "%s"
}
// We have to provide GCS default storage account with the permission
// to publish to a Cloud Pub/Sub topic from this project
// Otherwise notification configuration won't work
resource "google_pubsub_topic_iam_binding" "binding" {
	topic       = "${google_pubsub_topic.topic.name}"
	role        = "roles/pubsub.publisher"
		  
	members     = ["%s"]
}

resource "google_storage_notification" "notification" {
	bucket            = "${google_storage_bucket.bucket.name}"
	payload_format    = "JSON_API_V1"
	topic             = "${google_pubsub_topic.topic.id}"
	event_types       = ["%s","%s"]
	custom_attributes {
		new-attribute = "new-attribute-value"
	}
	depends_on        = ["google_pubsub_topic_iam_binding.binding"]
}

`, bucketName, topicName, gcsServiceAccount, eventType1, eventType2)
}
