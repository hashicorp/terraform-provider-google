package google

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"archive/zip"
	"io/ioutil"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
const testFunctionsSourceArchivePrefix = "cloudfunczip"

func init() {
	resource.AddTestSweepers("gcp_cloud_function_source_archive", &resource.Sweeper{
		Name: "gcp_cloud_function_source_archive",
		F:    sweepCloudFunctionSourceZipArchives,
	})
}

func TestCloudFunctionsFunction_nameValidator(t *testing.T) {
	validNames := []string{
		"a",
		"aA",
		"a0",
		"has-hyphen",
		"has_underscore",
		"hasUpperCase",
		"allChars_-A0",
	}
	for _, tc := range validNames {
		wrns, errs := validateResourceCloudFunctionsFunctionName(tc, "function.name")
		if len(wrns) > 0 {
			t.Errorf("Expected no validation warnings for test case %q, got: %+v", tc, wrns)
		}
		if len(errs) > 0 {
			t.Errorf("Expected no validation errors for test name %q, got: %+v", tc, errs)
		}
	}

	invalidNames := []string{
		"0startsWithNumber",
		"endsWith_",
		"endsWith-",
		"bad*Character",
		"aFunctionsNameThatIsLongerThanFortyEightCharacters",
	}
	for _, tc := range invalidNames {
		_, errs := validateResourceCloudFunctionsFunctionName(tc, "function.name")
		if len(errs) == 0 {
			t.Errorf("Expected errors for invalid test name %q, got none", tc)
		}
	}
}

func TestAccCloudFunctionsFunction_basic(t *testing.T) {
	t.Parallel()

	var function cloudfunctions.CloudFunction

	funcResourceName := "google_cloudfunctions_function.function"
	functionName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	zipFilePath := createZIPArchiveForCloudFunctionSource(t, testHTTPTriggerPath)
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
					resource.TestCheckResourceAttr(funcResourceName,
						"max_instances", "10"),
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
	zipFilePath := createZIPArchiveForCloudFunctionSource(t, testHTTPTriggerPath)
	zipFileUpdatePath := createZIPArchiveForCloudFunctionSource(t, testHTTPTriggerUpdatePath)
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
					resource.TestCheckResourceAttr(funcResourceName,
						"max_instances", "15"),
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
	zipFilePath := createZIPArchiveForCloudFunctionSource(t, testPubSubTriggerPath)
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
	zipFilePath := createZIPArchiveForCloudFunctionSource(t, testBucketTriggerPath)
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
	zipFilePath := createZIPArchiveForCloudFunctionSource(t, testFirestoreTriggerPath)
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
	proj := getTestProjectFromEnv()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_sourceRepo(functionName, proj),
			},
			{
				ResourceName:      funcResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCloudFunctionsFunction_serviceAccountEmail(t *testing.T) {
	t.Parallel()

	funcResourceName := "google_cloudfunctions_function.function"
	functionName := fmt.Sprintf("tf-test-%s", acctest.RandString(10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	zipFilePath := createZIPArchiveForCloudFunctionSource(t, testHTTPTriggerPath)
	defer os.Remove(zipFilePath) // clean up

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_serviceAccountEmail(functionName, bucketName, zipFilePath),
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

func createZIPArchiveForCloudFunctionSource(t *testing.T, sourcePath string) string {
	source, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		t.Fatal(err.Error())
	}
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	f, err := w.Create("index.js")
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = f.Write(source)
	if err != nil {
		t.Fatal(err.Error())
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		t.Fatal(err.Error())
	}
	// Create temp file to write zip to
	tmpfile, err := ioutil.TempFile("", "sourceArchivePrefix")
	if err != nil {
		t.Fatal(err.Error())
	}

	if _, err := tmpfile.Write(buf.Bytes()); err != nil {
		t.Fatal(err.Error())
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err.Error())
	}
	return tmpfile.Name()
}

func sweepCloudFunctionSourceZipArchives(_ string) error {
	files, err := ioutil.ReadDir(os.TempDir())
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasPrefix(f.Name(), testFunctionsSourceArchivePrefix) {
			filepath := fmt.Sprintf("%s/%s", os.TempDir(), f.Name())
			if err := os.Remove(filepath); err != nil {
				return err
			}
			log.Printf("[INFO] cloud functions sweeper removed old file %s", filepath)
		}
	}
	return nil
}

func testAccCloudFunctionsFunction_basic(functionName string, bucketName string, zipFilePath string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name                  = "%s"
  runtime               = "nodejs8"
  description           = "test function"
  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name
  trigger_http          = true
  timeout               = 61
  entry_point           = "helloGET"
  labels = {
    my-label = "my-label-value"
  }
  environment_variables = {
    TEST_ENV_VARIABLE = "test-env-variable-value"
  }
  max_instances = 10
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
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name                  = "%s"
  description           = "test function updated"
  available_memory_mb   = 256
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name
  trigger_http          = true
  runtime               = "nodejs8"
  timeout               = 91
  entry_point           = "helloGET"
  labels = {
    my-label    = "my-updated-label-value"
    a-new-label = "a-new-label-value"
  }
  environment_variables = {
    TEST_ENV_VARIABLE = "test-env-variable-value"
    NEW_ENV_VARIABLE  = "new-env-variable-value"
  }
  max_instances = 15
}
`, bucketName, zipFilePath, functionName)
}

func testAccCloudFunctionsFunction_pubsub(functionName string, bucketName string,
	topic string, zipFilePath string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_pubsub_topic" "sub" {
  name = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name                  = "%s"
  runtime               = "nodejs8"
  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name
  timeout               = 61
  entry_point           = "helloPubSub"
  event_trigger {
    event_type = "providers/cloud.pubsub/eventTypes/topic.publish"
    resource   = google_pubsub_topic.sub.name
    failure_policy {
      retry = false
    }
  }
}
`, bucketName, zipFilePath, topic, functionName)
}

func testAccCloudFunctionsFunction_bucket(functionName string, bucketName string,
	zipFilePath string) string {
	return fmt.Sprintf(`
data "google_client_config" "current" {
}

resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name                  = "%s"
  runtime               = "nodejs8"
  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name
  timeout               = 61
  entry_point           = "helloGCS"
  event_trigger {
    event_type = "google.storage.object.finalize"
    resource   = "projects/${data.google_client_config.current.project}/buckets/${google_storage_bucket.bucket.name}"
    failure_policy {
      retry = true
    }
  }
}
`, bucketName, zipFilePath, functionName)
}

func testAccCloudFunctionsFunction_bucketNoRetry(functionName string, bucketName string,
	zipFilePath string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name                  = "%s"
  runtime               = "nodejs8"
  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name
  timeout               = 61
  entry_point           = "helloGCS"
  event_trigger {
    event_type = "google.storage.object.finalize"
    resource   = google_storage_bucket.bucket.name
  }
}
`, bucketName, zipFilePath, functionName)
}

func testAccCloudFunctionsFunction_firestore(functionName string, bucketName string,
	zipFilePath string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name                  = "%s"
  runtime               = "nodejs8"
  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name
  timeout               = 61
  entry_point           = "helloFirestore"
  event_trigger {
    event_type = "providers/cloud.firestore/eventTypes/document.write"
    resource   = "messages/{messageId}"
  }
}
`, bucketName, zipFilePath, functionName)
}

func testAccCloudFunctionsFunction_sourceRepo(functionName, project string) string {
	return fmt.Sprintf(`
resource "google_cloudfunctions_function" "function" {
  name    = "%s"
  runtime = "nodejs8"

  source_repository {
    // There isn't yet an API that'll allow us to create a source repository and
    // put code in it, so we created this repository outside the test to be used
    // here. If this test is run outside of CI, you may need to create your own
    // source repo.
    url = "https://source.developers.google.com/projects/%s/repos/cloudfunctions-test-do-not-delete/moveable-aliases/master/paths/"
  }

  trigger_http = true
  entry_point  = "helloGET"
}
`, functionName, project)
}

func testAccCloudFunctionsFunction_serviceAccountEmail(functionName, bucketName, zipFilePath string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name = "%s"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

data "google_compute_default_service_account" "default" {
}

resource "google_cloudfunctions_function" "function" {
  name    = "%s"
  runtime = "nodejs8"

  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name

  service_account_email = data.google_compute_default_service_account.default.email

  trigger_http = true
  entry_point  = "helloGET"
}
`, bucketName, zipFilePath, functionName)
}
