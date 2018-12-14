package google

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"archive/zip"
	"io/ioutil"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"google.golang.org/api/cloudfunctions/v1"
)

const (
	FUNCTION_TRIGGER_HTTP = iota
)

const testHTTPTriggerPath = "./test-fixtures/cloudfunctions/http_trigger.js"
const testHTTPTriggerUpdatePath = "./test-fixtures/cloudfunctions/http_trigger_update.js"
const testPubSubTriggerPath = "./test-fixtures/cloudfunctions/pubsub_trigger.js"
const testBucketTriggerPath = "./test-fixtures/cloudfunctions/bucket_trigger.js"
const testFirestoreTriggerPath = "./test-fixtures/cloudfunctions/firestore_trigger.js"

func TestAccCloudFunctionsFunction_basic(t *testing.T) {
	t.Parallel()

	var function cloudfunctions.CloudFunction

	funcResourceName := "google_cloudfunctions_function.function"
	functionName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	zipFilePath, err := createZIPArchiveForIndexJs(testHTTPTriggerPath)
	if err != nil {
		t.Fatal(err.Error())
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
						funcResourceName, &function),
					resource.TestCheckResourceAttr(funcResourceName,
						"name", functionName),
					resource.TestCheckResourceAttr(funcResourceName,
						"description", "test function"),
					resource.TestCheckResourceAttr(funcResourceName,
						"available_memory_mb", "128"),
					testAccCloudFunctionsFunctionSource(fmt.Sprintf("gs://%s/index.zip", bucketName), &function),
					testAccCloudFunctionsFunctionTrigger(FUNCTION_TRIGGER_HTTP, &function),
					resource.TestCheckResourceAttr(funcResourceName,
						"timeout", "61"),
					resource.TestCheckResourceAttr(funcResourceName,
						"entry_point", "helloGET"),
					resource.TestCheckResourceAttr(funcResourceName,
						"trigger_http", "true"),
					testAccCloudFunctionsFunctionHasLabel("my-label", "my-label-value", &function),
					testAccCloudFunctionsFunctionHasEnvironmentVariable("TEST_ENV_VARIABLE",
						"test-env-variable-value", &function),
				),
			},
			{
				ResourceName:      funcResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudFunctionsFunction_update(t *testing.T) {
	t.Parallel()

	var function cloudfunctions.CloudFunction

	funcResourceName := "google_cloudfunctions_function.function"
	functionName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	zipFilePath, err := createZIPArchiveForIndexJs(testHTTPTriggerPath)
	zipFileUpdatePath, err := createZIPArchiveForIndexJs(testHTTPTriggerUpdatePath)
	if err != nil {
		t.Fatal(err.Error())
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
						funcResourceName, &function),
					resource.TestCheckResourceAttr(funcResourceName,
						"available_memory_mb", "128"),
					testAccCloudFunctionsFunctionHasLabel("my-label", "my-label-value", &function),
				),
			},
			{
				ResourceName:      funcResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudFunctionsFunction_updated(functionName, bucketName, zipFileUpdatePath),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudFunctionsFunctionExists(
						funcResourceName, &function),
					resource.TestCheckResourceAttr(funcResourceName,
						"available_memory_mb", "256"),
					resource.TestCheckResourceAttr(funcResourceName,
						"description", "test function updated"),
					resource.TestCheckResourceAttr(funcResourceName,
						"timeout", "91"),
					testAccCloudFunctionsFunctionHasLabel("my-label", "my-updated-label-value", &function),
					testAccCloudFunctionsFunctionHasLabel("a-new-label", "a-new-label-value", &function),
					testAccCloudFunctionsFunctionHasEnvironmentVariable("TEST_ENV_VARIABLE",
						"test-env-variable-value", &function),
					testAccCloudFunctionsFunctionHasEnvironmentVariable("NEW_ENV_VARIABLE",
						"new-env-variable-value", &function),
				),
			},
			{
				ResourceName:      funcResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudFunctionsFunction_pubsub(t *testing.T) {
	t.Parallel()

	funcResourceName := "google_cloudfunctions_function.function"
	functionName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	topicName := fmt.Sprintf("tf-test-sub-%s", acctest.RandString(10))
	zipFilePath, err := createZIPArchiveForIndexJs(testPubSubTriggerPath)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.Remove(zipFilePath) // clean up

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_pubsub(functionName, bucketName,
					topicName, zipFilePath),
			},
			{
				ResourceName:      funcResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudFunctionsFunction_bucket(t *testing.T) {
	t.Parallel()
	funcResourceName := "google_cloudfunctions_function.function"
	functionName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	zipFilePath, err := createZIPArchiveForIndexJs(testBucketTriggerPath)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.Remove(zipFilePath) // clean up

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_bucket(functionName, bucketName, zipFilePath),
			},
			{
				ResourceName:      funcResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCloudFunctionsFunction_bucketNoRetry(functionName, bucketName, zipFilePath),
			},
			{
				ResourceName:      funcResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudFunctionsFunction_firestore(t *testing.T) {
	t.Parallel()
	funcResourceName := "google_cloudfunctions_function.function"
	functionName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	zipFilePath, err := createZIPArchiveForIndexJs(testFirestoreTriggerPath)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.Remove(zipFilePath) // clean up

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_firestore(functionName, bucketName, zipFilePath),
			},
			{
				ResourceName:      funcResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudFunctionsFunction_sourceRepo(t *testing.T) {
	t.Parallel()

	funcResourceName := "google_cloudfunctions_function.function"
	functionName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_sourceRepo(functionName),
			},
			{
				ResourceName:      funcResourceName,
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
		cloudFuncId := &cloudFunctionId{
			Project: project,
			Region:  region,
			Name:    name,
		}
		_, err := config.clientCloudFunctions.Projects.Locations.Functions.Get(cloudFuncId.cloudFunctionId()).Do()
		if err == nil {
			return fmt.Errorf("Function still exists")
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
		cloudFuncId := &cloudFunctionId{
			Project: project,
			Region:  region,
			Name:    name,
		}
		found, err := config.clientCloudFunctions.Projects.Locations.Functions.Get(cloudFuncId.cloudFunctionId()).Do()
		if err != nil {
			return fmt.Errorf("CloudFunctions Function not present")
		}

		*function = *found

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
				return fmt.Errorf("Expected HttpsTrigger to be set")
			}
		default:
			return fmt.Errorf("testAccCloudFunctionsFunctionTrigger expects only FUNCTION_TRIGGER_HTTP, ")
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

func testAccCloudFunctionsFunctionHasEnvironmentVariable(key, value string,
	function *cloudfunctions.CloudFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if val, ok := function.EnvironmentVariables[key]; ok {
			if val != value {
				return fmt.Errorf("Environment Variable value did not match for key %s: expected %s but found %s",
					key, value, val)
			}
		} else {
			return fmt.Errorf("Environment Variable with key %s not found", key)
		}
		return nil
	}
}

func createZIPArchiveForIndexJs(sourcePath string) (string, error) {
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
	// Create temp file to write zip to
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
  name                  = "%s"
  description           = "test function"
  available_memory_mb   = 128
  source_archive_bucket = "${google_storage_bucket.bucket.name}"
  source_archive_object = "${google_storage_bucket_object.archive.name}"
  trigger_http          = true
  timeout               = 61
  entry_point           = "helloGET"
  labels {
	my-label = "my-label-value"
  }
  environment_variables {
	TEST_ENV_VARIABLE = "test-env-variable-value"
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
  name   = "index_update.zip"
  bucket = "${google_storage_bucket.bucket.name}"
  source = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name                  = "%s"
  description           = "test function updated"
  available_memory_mb   = 256
  source_archive_bucket = "${google_storage_bucket.bucket.name}"
  source_archive_object = "${google_storage_bucket_object.archive.name}"
  trigger_http          = true
  runtime               = "nodejs8"
  timeout               = 91
  entry_point           = "helloGET"
  labels {
	my-label = "my-updated-label-value"
	a-new-label = "a-new-label-value"
  }
  environment_variables {
	TEST_ENV_VARIABLE = "test-env-variable-value"
	NEW_ENV_VARIABLE = "new-env-variable-value"
  }
}`, bucketName, zipFilePath, functionName)
}

func testAccCloudFunctionsFunction_pubsub(functionName string, bucketName string,
	topic string, zipFilePath string) string {
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
  name                  = "%s"
  available_memory_mb   = 128
  source_archive_bucket = "${google_storage_bucket.bucket.name}"
  source_archive_object = "${google_storage_bucket_object.archive.name}"
  timeout               = 61
  entry_point           = "helloPubSub"
  event_trigger {
    event_type = "providers/cloud.pubsub/eventTypes/topic.publish"
    resource   = "${google_pubsub_topic.sub.name}"
    failure_policy {
      retry = false
    }
  }
}`, bucketName, zipFilePath, topic, functionName)
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
  name                  = "%s"
  available_memory_mb   = 128
  source_archive_bucket = "${google_storage_bucket.bucket.name}"
  source_archive_object = "${google_storage_bucket_object.archive.name}"
  timeout               = 61
  entry_point           = "helloGCS"
  event_trigger {
    event_type = "google.storage.object.finalize"
    resource   = "${google_storage_bucket.bucket.name}"
    failure_policy {
      retry = true
    }
  }
}`, bucketName, zipFilePath, functionName)
}

func testAccCloudFunctionsFunction_bucketNoRetry(functionName string, bucketName string,
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
  name                  = "%s"
  available_memory_mb   = 128
  source_archive_bucket = "${google_storage_bucket.bucket.name}"
  source_archive_object = "${google_storage_bucket_object.archive.name}"
  timeout               = 61
  entry_point           = "helloGCS"
  event_trigger {
    event_type = "google.storage.object.finalize"
    resource   = "${google_storage_bucket.bucket.name}"
  }
}`, bucketName, zipFilePath, functionName)
}

func testAccCloudFunctionsFunction_firestore(functionName string, bucketName string,
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
  name                  = "%s"
  available_memory_mb   = 128
  source_archive_bucket = "${google_storage_bucket.bucket.name}"
  source_archive_object = "${google_storage_bucket_object.archive.name}"
  timeout               = 61
  entry_point           = "helloFirestore"
  event_trigger {
    event_type = "providers/cloud.firestore/eventTypes/document.write"
    resource   = "messages/{messageId}"
  }
}`, bucketName, zipFilePath, functionName)
}

func testAccCloudFunctionsFunction_sourceRepo(functionName string) string {
	return fmt.Sprintf(`
resource "google_cloudfunctions_function" "function" {
  name = "%s"

  source_repository {
    // There isn't yet an API that'll allow us to create a source repository and
    // put code in it, so we created this repository outside the test to be used
    // here. If this test is run outside of CI, it may fail because of permissions
    // errors.
    url = "https://source.developers.google.com/projects/hc-terraform-testing/repos/cloudfunctions-test-do-not-delete/moveable-aliases/master/paths/"
  }

  trigger_http = true
  entry_point  = "helloGET"
}
`, functionName)
}
