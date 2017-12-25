package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/cloudfunctions/v1"
	"strings"
)

const (
	FUNCTION_TRIGGER_HTTP = iota
	FUNCTION_TRIGGER_TOPIC
	FUNCTION_TRIGGER_BUCKET
)

func TestAccCloudFunctionsFunction_basic(t *testing.T) {
	t.Parallel()

	var function cloudfunctions.CloudFunction

	functionName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_basic(functionName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudFunctionsFunctionExists(
						"google_cloudfunctions_function.function", &function),
					testAccCloudFunctionsFunctionName(functionName, &function),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function", "description", "test function"),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function", "memory", "128"),
					testAccCloudFunctionsFunctionSource(fmt.Sprintf("gs://%s/index.zip", bucketName), &function),
					testAccCloudFunctionsFunctionTrigger(FUNCTION_TRIGGER_HTTP, &function),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function", "timeout", "61"),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function", "entry_point", "helloGET"),
					testAccCloudFunctionsFunctionHasLabel("my-label", "my-label-value", &function),
				),
			},
			{
				ResourceName:      "google_cloudfunctions_function.function",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudFunctionsFunction_update(t *testing.T) {
	t.Parallel()

	functionName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	var function cloudfunctions.CloudFunction

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_basic(functionName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudFunctionsFunctionExists(
						"google_cloudfunctions_function.function", &function),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function", "memory", "128"),
					testAccCloudFunctionsFunctionHasLabel("my-label", "my-label-value", &function),
				),
			},
			{
				Config: testAccCloudFunctionsFunction_updated(functionName, bucketName),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudFunctionsFunctionExists(
						"google_cloudfunctions_function.function", &function),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function", "memory", "256"),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function", "description", "test function updated"),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function", "timeout", "91"),
					testAccCloudFunctionsFunctionHasLabel("my-label", "my-updated-label-value", &function),
					testAccCloudFunctionsFunctionHasLabel("a-new-label", "a-new-label-value", &function),
				),
			},
		},
	})
}

func testAccCheckCloudFunctionsFunctionDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "google_cloudfunctions_function" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		project := rs.Primary.Attributes["project"]
		region := rs.Primary.Attributes["region"]
		_, err := config.clientCloudFunctions.Projects.Locations.Functions.Get(
			createCloudFunctionsPathString(CLOUDFUNCTIONS_FULL_NAME, project, region, name)).Do()
		if err == nil {
			return fmt.Errorf("CloudFunctions still exists")
		}

	}

	return nil
}

func testAccCloudFunctionsFunctionExists(n string, function *cloudfunctions.CloudFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := testAccProvider.Meta().(*Config)
		name := rs.Primary.Attributes["name"]
		project := rs.Primary.Attributes["project"]
		region := rs.Primary.Attributes["region"]
		found, err := config.clientCloudFunctions.Projects.Locations.Functions.Get(
			createCloudFunctionsPathString(CLOUDFUNCTIONS_FULL_NAME, project, region, name)).Do()
		if err != nil {
			return fmt.Errorf("CloudFunctions Function not present")
		}

		*function = *found

		return nil
	}
}

func testAccCloudFunctionsFunctionName(n string, function *cloudfunctions.CloudFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		expected, err := getCloudFunctionName(function.Name)
		if err != nil {
			return err
		}
		if n != expected {
			return fmt.Errorf("Expected function name %s, got %s", n, expected)
		}
		return nil
	}
}

func testAccCloudFunctionsFunctionSource(n string, function *cloudfunctions.CloudFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if n != function.SourceArchiveUrl {
			return fmt.Errorf("Expected source to be %v, got %v", n, function.EntryPoint)
		}
		return nil
	}
}

func testAccCloudFunctionsFunctionTrigger(n int, function *cloudfunctions.CloudFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		switch n {
		case FUNCTION_TRIGGER_HTTP:
			if function.HttpsTrigger == nil {
				return fmt.Errorf("Expected trigger_http to be set")
			}
		case FUNCTION_TRIGGER_BUCKET:
			if function.EventTrigger == nil {
				return fmt.Errorf("Expected trigger_bucket to be set")
			}
			if strings.Index(function.EventTrigger.EventType, "cloud.storage") == -1 {
				return fmt.Errorf("Expected trigger_bucket to be set")
			}
		case FUNCTION_TRIGGER_TOPIC:
			if function.EventTrigger == nil {
				return fmt.Errorf("Expected trigger_bucket to be set")
			}
			if strings.Index(function.EventTrigger.EventType, "cloud.pubsub") == -1 {
				return fmt.Errorf("Expected trigger_topic to be set")
			}
		default:
			return fmt.Errorf("testAccCloudFunctionsFunctionTrigger expects only FUNCTION_TRIGGER_HTTP, " +
				"FUNCTION_TRIGGER_BUCKET or FUNCTION_TRIGGER_TOPIC")
		}
		return nil
	}
}

func testAccCloudFunctionsFunctionHasLabel(key, value string, function *cloudfunctions.CloudFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		val, ok := function.Labels[key]
		if !ok {
			return fmt.Errorf("Label with key %s not found", key)
		}

		if val != value {
			return fmt.Errorf("Label value did not match for key %s: expected %s but found %s", key, value, val)
		}
		return nil
	}
}

func testAccCloudFunctionsFunction_basic(functionName string, bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = "${google_storage_bucket.bucket.name}"
  source = "test-fixtures/index.zip"
}

resource "google_cloudfunctions_function" "function" {
  name           = "%s"
  description    = "test function"
  memory		 = 128
  storage_bucket = "${google_storage_bucket.bucket.name}"
  storage_object = "${google_storage_bucket_object.archive.name}"
  trigger_http   = true
  timeout		 = 61
  entry_point    = "helloGET"
  labels {
	my-label = "my-label-value"
  }
}
`, bucketName, functionName)
}

func testAccCloudFunctionsFunction_updated(functionName string, bucketName string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = "${google_storage_bucket.bucket.name}"
  source = "test-fixtures/index.zip"
}

resource "google_cloudfunctions_function" "function" {
  name          = "%s"
  description   = "test function updated"
  memory		= 256
  storage_bucket = "${google_storage_bucket.bucket.name}"
  storage_object = "${google_storage_bucket_object.archive.name}"
  trigger_http  = true
  timeout		= 91
  entry_point   = "helloGET"
  labels {
	my-label = "my-updated-label-value"
	a-new-label = "a-new-label-value"
  }
}`, bucketName, functionName)
}
