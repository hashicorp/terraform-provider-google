package google

import (
	"fmt"
	"testing"

	"archive/zip"
	"bytes"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/cloudfunctions/v1"
	"io/ioutil"
	"os"
	"strings"
)

const (
	FUNCTION_TRIGGER_HTTP = iota
	FUNCTION_TRIGGER_TOPIC
	FUNCTION_TRIGGER_BUCKET
)

const testHTTPTriggerPath = "./test-fixtures/cloudfunctions/http_trigger.js"
const testPubSubTriggerPath = "./test-fixtures/cloudfunctions/pubsub_trigger.js"
const testBucketTriggerPath = "./test-fixtures/cloudfunctions/bucket_trigger.js"

func TestAccCloudFunctionsFunction_basic(t *testing.T) {
	t.Parallel()

	var function cloudfunctions.CloudFunction

	functionName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	zipFilePath, err := createZIParchiveForIndexJs(testHTTPTriggerPath)
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
	defer os.Remove(zipFilePath) // clean up

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_basic(functionName, bucketName, zipFilePath),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudFunctionsFunctionExists(
						"google_cloudfunctions_function.function", &function),
					testAccCloudFunctionsFunctionName(functionName, &function),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function",
						"description", "test function"),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function",
						"memory", "128"),
					testAccCloudFunctionsFunctionSource(fmt.Sprintf("gs://%s/index.zip", bucketName), &function),
					testAccCloudFunctionsFunctionTrigger(FUNCTION_TRIGGER_HTTP, &function),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function",
						"timeout", "61"),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function",
						"entry_point", "helloGET"),
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

	var function cloudfunctions.CloudFunction

	functionName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	zipFilePath, err := createZIParchiveForIndexJs(testHTTPTriggerPath)
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
	defer os.Remove(zipFilePath) // clean up

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_basic(functionName, bucketName, zipFilePath),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudFunctionsFunctionExists(
						"google_cloudfunctions_function.function", &function),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function",
						"memory", "128"),
					testAccCloudFunctionsFunctionHasLabel("my-label", "my-label-value", &function),
				),
			},
			{
				Config: testAccCloudFunctionsFunction_updated(functionName, bucketName, zipFilePath),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudFunctionsFunctionExists(
						"google_cloudfunctions_function.function", &function),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function",
						"memory", "256"),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function",
						"description", "test function updated"),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function",
						"timeout", "91"),
					testAccCloudFunctionsFunctionHasLabel("my-label", "my-updated-label-value", &function),
					testAccCloudFunctionsFunctionHasLabel("a-new-label", "a-new-label-value", &function),
				),
			},
		},
	})
}

func TestAccCloudFunctionsFunction_pubsub(t *testing.T) {
	t.Parallel()

	var function cloudfunctions.CloudFunction

	functionName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	subscription := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(10))
	zipFilePath, err := createZIParchiveForIndexJs(testPubSubTriggerPath)
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
	defer os.Remove(zipFilePath) // clean up

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_pubsub(functionName, bucketName, subscription, zipFilePath),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudFunctionsFunctionExists(
						"google_cloudfunctions_function.function", &function),
					testAccCloudFunctionsFunctionName(functionName, &function),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function",
						"memory", "128"),
					testAccCloudFunctionsFunctionSource(fmt.Sprintf("gs://%s/index.zip", bucketName), &function),
					testAccCloudFunctionsFunctionTrigger(FUNCTION_TRIGGER_TOPIC, &function),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function",
						"timeout", "61"),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function",
						"entry_point", "helloPubSub"),
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
func TestAccCloudFunctionsFunction_bucket(t *testing.T) {
	t.Parallel()

	var function cloudfunctions.CloudFunction

	functionName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	zipFilePath, err := createZIParchiveForIndexJs(testBucketTriggerPath)
	if err != nil {
		t.Errorf(err.Error())
		t.FailNow()
	}
	defer os.Remove(zipFilePath) // clean up

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_bucket(functionName, bucketName, zipFilePath),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudFunctionsFunctionExists(
						"google_cloudfunctions_function.function", &function),
					testAccCloudFunctionsFunctionName(functionName, &function),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function",
						"memory", "128"),
					testAccCloudFunctionsFunctionSource(fmt.Sprintf("gs://%s/index.zip", bucketName), &function),
					testAccCloudFunctionsFunctionTrigger(FUNCTION_TRIGGER_BUCKET, &function),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function",
						"timeout", "61"),
					resource.TestCheckResourceAttr("google_cloudfunctions_function.function",
						"entry_point", "helloGCS"),
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

func testAccCloudFunctionsFunctionHasLabel(key, value string,
	function *cloudfunctions.CloudFunction) resource.TestCheckFunc {
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

func createZIParchiveForIndexJs(sourcePath string) (string, error) {
	source, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return "", err
	}
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	f, err := w.Create("index.js")
	if err != nil {
		return "", err
	}
	_, err = f.Write(source)
	if err != nil {
		return "", err
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		return "", err
	}
	//Create temp file to write zip to
	tmpfile, err := ioutil.TempFile("", "zip")
	if err != nil {
		return "", err
	}

	if _, err := tmpfile.Write(buf.Bytes()); err != nil {
		return "", err
	}
	if err := tmpfile.Close(); err != nil {
		return "", err
	}
	return tmpfile.Name(), nil
}

func testAccCloudFunctionsFunction_basic(functionName string, bucketName string, zipFilePath string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = "${google_storage_bucket.bucket.name}"
  source = "%s"
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
`, bucketName, zipFilePath, functionName)
}

func testAccCloudFunctionsFunction_updated(functionName string, bucketName string, zipFilePath string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = "${google_storage_bucket.bucket.name}"
  source = "%s"
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
}`, bucketName, zipFilePath, functionName)
}

func testAccCloudFunctionsFunction_pubsub(functionName string, bucketName string,
	subscription string, zipFilePath string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = "${google_storage_bucket.bucket.name}"
  source = "%s"
}

resource "google_pubsub_topic" "sub" {
	name = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name           = "%s"
  memory		 = 128
  storage_bucket = "${google_storage_bucket.bucket.name}"
  storage_object = "${google_storage_bucket_object.archive.name}"
  trigger_topic  = "${google_pubsub_topic.sub.name}"
  timeout		 = 61
  entry_point    = "helloPubSub"
}`, bucketName, zipFilePath, subscription, functionName)
}

func testAccCloudFunctionsFunction_bucket(functionName string, bucketName string,
	zipFilePath string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
	name = "%s"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = "${google_storage_bucket.bucket.name}"
  source = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name           = "%s"
  memory		 = 128
  storage_bucket = "${google_storage_bucket.bucket.name}"
  storage_object = "${google_storage_bucket_object.archive.name}"
  trigger_bucket  = "${google_storage_bucket.bucket.name}"
  timeout		 = 61
  entry_point    = "helloGCS"
}`, bucketName, zipFilePath, functionName)
}