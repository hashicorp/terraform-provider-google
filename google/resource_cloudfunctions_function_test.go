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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
const testSecretEnvVarFunctionPath = "./test-fixtures/cloudfunctions/secret_environment_variables.js"
const testSecretVolumesMountFunctionPath = "./test-fixtures/cloudfunctions/secret_volumes_mount.js"
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
		"StartsUpperCase",
		"endsUpperCasE",
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
		"aCloudFunctionsFunctionNameThatIsSeventyFiveCharactersLongWhichIsMoreThan63",
	}
	for _, tc := range invalidNames {
		_, errs := validateResourceCloudFunctionsFunctionName(tc, "function.name")
		if len(errs) == 0 {
			t.Errorf("Expected errors for invalid test name %q, got none", tc)
		}
	}
}

func TestValidLabelKeys(t *testing.T) {
	testCases := []struct {
		labelKey string
		valid    bool
	}{
		{
			"test-label", true,
		},
		{
			"test_label", true,
		},
		{
			"MixedCase", false,
		},
		{
			"number-09-dash", true,
		},
		{
			"", false,
		},
		{
			"test-label", true,
		},
		{
			"mixed*symbol", false,
		},
		{
			"intérnätional", true,
		},
	}

	for _, tc := range testCases {
		labels := make(map[string]interface{})
		labels[tc.labelKey] = "test value"

		_, errs := labelKeyValidator(labels, "")
		if tc.valid && len(errs) > 0 {
			t.Errorf("Validation failure, key: '%s' should be valid but actual errors were %q", tc.labelKey, errs)
		}
		if !tc.valid && len(errs) < 1 {
			t.Errorf("Validation failure, key: '%s' should fail but actual errors were %q", tc.labelKey, errs)
		}
	}
}

func TestAccCloudFunctionsFunction_basic(t *testing.T) {
	t.Parallel()

	var function cloudfunctions.CloudFunction

	funcResourceName := "google_cloudfunctions_function.function"
	functionName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", randInt(t))
	zipFilePath := createZIPArchiveForCloudFunctionSource(t, testHTTPTriggerPath)
	defer os.Remove(zipFilePath) // clean up

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_basic(functionName, bucketName, zipFilePath),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudFunctionsFunctionExists(
						t, funcResourceName, &function),
					resource.TestCheckResourceAttr(funcResourceName,
						"name", functionName),
					resource.TestCheckResourceAttr(funcResourceName,
						"description", "test function"),
					resource.TestCheckResourceAttr(funcResourceName,
						"docker_registry", "CONTAINER_REGISTRY"),
					resource.TestCheckResourceAttr(funcResourceName,
						"available_memory_mb", "128"),
					resource.TestCheckResourceAttr(funcResourceName,
						"max_instances", "10"),
					resource.TestCheckResourceAttr(funcResourceName,
						"min_instances", "3"),
					resource.TestCheckResourceAttr(funcResourceName,
						"ingress_settings", "ALLOW_INTERNAL_ONLY"),
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
				ResourceName:            funcResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"build_environment_variables"},
			},
		},
	})
}

func TestAccCloudFunctionsFunction_update(t *testing.T) {
	t.Parallel()

	var function cloudfunctions.CloudFunction

	funcResourceName := "google_cloudfunctions_function.function"
	functionName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", randInt(t))
	zipFilePath := createZIPArchiveForCloudFunctionSource(t, testHTTPTriggerPath)
	zipFileUpdatePath := createZIPArchiveForCloudFunctionSource(t, testHTTPTriggerUpdatePath)
	random_suffix := randString(t, 10)
	defer os.Remove(zipFilePath) // clean up

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_basic(functionName, bucketName, zipFilePath),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudFunctionsFunctionExists(
						t, funcResourceName, &function),
					resource.TestCheckResourceAttr(funcResourceName,
						"available_memory_mb", "128"),
					testAccCloudFunctionsFunctionHasLabel("my-label", "my-label-value", &function),
				),
			},
			{
				ResourceName:            funcResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"build_environment_variables"},
			},
			{
				Config: testAccCloudFunctionsFunction_updated(functionName, bucketName, zipFileUpdatePath, random_suffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudFunctionsFunctionExists(
						t, funcResourceName, &function),
					resource.TestCheckResourceAttr(funcResourceName,
						"available_memory_mb", "256"),
					resource.TestCheckResourceAttr(funcResourceName,
						"description", "test function updated"),
					resource.TestCheckResourceAttr(funcResourceName,
						"docker_registry", "ARTIFACT_REGISTRY"),
					resource.TestCheckResourceAttr(funcResourceName,
						"timeout", "91"),
					resource.TestCheckResourceAttr(funcResourceName,
						"max_instances", "15"),
					resource.TestCheckResourceAttr(funcResourceName,
						"min_instances", "5"),
					resource.TestCheckResourceAttr(funcResourceName,
						"ingress_settings", "ALLOW_ALL"),
					testAccCloudFunctionsFunctionHasLabel("my-label", "my-updated-label-value", &function),
					testAccCloudFunctionsFunctionHasLabel("a-new-label", "a-new-label-value", &function),
					testAccCloudFunctionsFunctionHasEnvironmentVariable("TEST_ENV_VARIABLE",
						"test-env-variable-value", &function),
					testAccCloudFunctionsFunctionHasEnvironmentVariable("NEW_ENV_VARIABLE",
						"new-env-variable-value", &function),
				),
			},
			{
				ResourceName:            funcResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"build_environment_variables"},
			},
		},
	})
}

func TestAccCloudFunctionsFunction_buildworkerpool(t *testing.T) {
	t.Parallel()

	var function cloudfunctions.CloudFunction

	funcResourceName := "google_cloudfunctions_function.function"
	functionName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", randInt(t))
	location := "us-central1"
	zipFilePath := createZIPArchiveForCloudFunctionSource(t, testHTTPTriggerPath)
	proj := getTestProjectFromEnv()

	defer os.Remove(zipFilePath) // clean up

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_buildworkerpool(functionName, bucketName, zipFilePath, location),
				Check: resource.ComposeTestCheckFunc(
					testAccCloudFunctionsFunctionExists(
						t, funcResourceName, &function),
					resource.TestCheckResourceAttr(funcResourceName,
						"name", functionName),
					resource.TestCheckResourceAttr(funcResourceName,
						"build_worker_pool", fmt.Sprintf("projects/%s/locations/%s/workerPools/pool-%s", proj, location, functionName)),
				),
			},
			{
				ResourceName:            funcResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"build_environment_variables"},
			},
		},
	})
}

func TestAccCloudFunctionsFunction_pubsub(t *testing.T) {
	t.Parallel()

	funcResourceName := "google_cloudfunctions_function.function"
	functionName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", randInt(t))
	topicName := fmt.Sprintf("tf-test-sub-%s", randString(t, 10))
	zipFilePath := createZIPArchiveForCloudFunctionSource(t, testPubSubTriggerPath)
	defer os.Remove(zipFilePath) // clean up

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_pubsub(functionName, bucketName,
					topicName, zipFilePath),
			},
			{
				ResourceName:            funcResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"build_environment_variables"},
			},
		},
	})
}

func TestAccCloudFunctionsFunction_bucket(t *testing.T) {
	t.Parallel()
	funcResourceName := "google_cloudfunctions_function.function"
	functionName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", randInt(t))
	zipFilePath := createZIPArchiveForCloudFunctionSource(t, testBucketTriggerPath)
	defer os.Remove(zipFilePath) // clean up

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_bucket(functionName, bucketName, zipFilePath),
			},
			{
				ResourceName:            funcResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"build_environment_variables"},
			},
			{
				Config: testAccCloudFunctionsFunction_bucketNoRetry(functionName, bucketName, zipFilePath),
			},
			{
				ResourceName:            funcResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"build_environment_variables"},
			},
		},
	})
}

func TestAccCloudFunctionsFunction_firestore(t *testing.T) {
	t.Parallel()
	funcResourceName := "google_cloudfunctions_function.function"
	functionName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", randInt(t))
	zipFilePath := createZIPArchiveForCloudFunctionSource(t, testFirestoreTriggerPath)
	defer os.Remove(zipFilePath) // clean up

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_firestore(functionName, bucketName, zipFilePath),
			},
			{
				ResourceName:            funcResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"build_environment_variables"},
			},
		},
	})
}

func TestAccCloudFunctionsFunction_sourceRepo(t *testing.T) {
	t.Parallel()

	funcResourceName := "google_cloudfunctions_function.function"
	functionName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	proj := getTestProjectFromEnv()

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_sourceRepo(functionName, proj),
			},
			{
				ResourceName:            funcResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"build_environment_variables"},
			},
		},
	})
}

func TestAccCloudFunctionsFunction_serviceAccountEmail(t *testing.T) {
	t.Parallel()

	funcResourceName := "google_cloudfunctions_function.function"
	functionName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", randInt(t))
	zipFilePath := createZIPArchiveForCloudFunctionSource(t, testHTTPTriggerPath)
	defer os.Remove(zipFilePath) // clean up

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_serviceAccountEmail(functionName, bucketName, zipFilePath),
			},
			{
				ResourceName:            funcResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"build_environment_variables"},
			},
		},
	})
}

func TestAccCloudFunctionsFunction_vpcConnector(t *testing.T) {
	t.Parallel()

	funcResourceName := "google_cloudfunctions_function.function"
	functionName := fmt.Sprintf("tf-test-%s", randString(t, 10))
	bucketName := fmt.Sprintf("tf-test-bucket-%d", randInt(t))
	networkName := fmt.Sprintf("tf-test-net-%d", randInt(t))
	vpcConnectorName := fmt.Sprintf("tf-test-conn-%s", randString(t, 5))
	zipFilePath := createZIPArchiveForCloudFunctionSource(t, testHTTPTriggerPath)
	projectNumber := os.Getenv("GOOGLE_PROJECT_NUMBER")
	defer os.Remove(zipFilePath) // clean up

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_vpcConnector(projectNumber, networkName, functionName, bucketName, zipFilePath, "10.10.0.0/28", vpcConnectorName),
			},
			{
				ResourceName:            funcResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"build_environment_variables"},
			},
			{
				Config: testAccCloudFunctionsFunction_vpcConnector(projectNumber, networkName, functionName, bucketName, zipFilePath, "10.20.0.0/28", vpcConnectorName+"-update"),
			},
			{
				ResourceName:            funcResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"build_environment_variables"},
			},
		},
	})
}

func TestAccCloudFunctionsFunction_secretEnvVar(t *testing.T) {
	t.Parallel()

	randomSecretSuffix := randString(t, 10)
	accountId := fmt.Sprintf("tf-test-account-%s", randomSecretSuffix)
	secretName := fmt.Sprintf("tf-test-secret-%s", randomSecretSuffix)
	versionName1 := fmt.Sprintf("tf-test-version1-%s", randomSecretSuffix)
	versionName2 := fmt.Sprintf("tf-test-version2-%s", randomSecretSuffix)
	bucketName := fmt.Sprintf("tf-test-bucket-%d", randInt(t))
	functionName := fmt.Sprintf("tf-test-%s", randomSecretSuffix)
	zipFilePath := createZIPArchiveForCloudFunctionSource(t, testSecretEnvVarFunctionPath)
	funcResourceName := "google_cloudfunctions_function.function"
	defer os.Remove(zipFilePath) // clean up

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_secretEnvVar(secretName, versionName1, bucketName, functionName, "1", zipFilePath, accountId),
			},
			{
				ResourceName:            funcResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"build_environment_variables"},
			},
			{
				Config: testAccCloudFunctionsFunction_secretEnvVar(secretName, versionName2, bucketName+"-update", functionName, "2", zipFilePath, accountId),
			},
			{
				ResourceName:            funcResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"build_environment_variables"},
			},
		},
	})
}

func TestAccCloudFunctionsFunction_secretMount(t *testing.T) {
	t.Parallel()

	projectNumber := os.Getenv("GOOGLE_PROJECT_NUMBER")
	randomSecretSuffix := randString(t, 10)
	accountId := fmt.Sprintf("tf-test-account-%s", randomSecretSuffix)
	secretName := fmt.Sprintf("tf-test-secret-%s", randomSecretSuffix)
	versionName1 := fmt.Sprintf("tf-test-version1-%s", randomSecretSuffix)
	versionName2 := fmt.Sprintf("tf-test-version2-%s", randomSecretSuffix)
	bucketName := fmt.Sprintf("tf-test-bucket-%d", randInt(t))
	functionName := fmt.Sprintf("tf-test-%s", randomSecretSuffix)
	zipFilePath := createZIPArchiveForCloudFunctionSource(t, testSecretVolumesMountFunctionPath)
	funcResourceName := "google_cloudfunctions_function.function"
	defer os.Remove(zipFilePath) // clean up

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudFunctionsFunctionDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccCloudFunctionsFunction_secretMount(projectNumber, secretName, versionName1, bucketName, functionName, "1", zipFilePath, accountId),
			},
			{
				ResourceName:            funcResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"build_environment_variables"},
			},
			{
				Config: testAccCloudFunctionsFunction_secretMount(projectNumber, secretName, versionName2, bucketName, functionName, "2", zipFilePath, accountId),
			},
			{
				ResourceName:            funcResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"build_environment_variables"},
			},
		},
	})
}

func testAccCheckCloudFunctionsFunctionDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		config := googleProviderConfig(t)

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
			_, err := config.NewCloudFunctionsClient(config.userAgent).Projects.Locations.Functions.Get(cloudFuncId.cloudFunctionId()).Do()
			if err == nil {
				return fmt.Errorf("Function still exists")
			}

		}

		return nil
	}
}

func testAccCloudFunctionsFunctionExists(t *testing.T, n string, function *cloudfunctions.CloudFunction) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		config := googleProviderConfig(t)
		name := rs.Primary.Attributes["name"]
		project := rs.Primary.Attributes["project"]
		region := rs.Primary.Attributes["region"]
		cloudFuncId := &cloudFunctionId{
			Project: project,
			Region:  region,
			Name:    name,
		}
		found, err := config.NewCloudFunctionsClient(config.userAgent).Projects.Locations.Functions.Get(cloudFuncId.cloudFunctionId()).Do()
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
		log.Printf("Error reading files: %s", err)
		return nil
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if strings.HasPrefix(f.Name(), testFunctionsSourceArchivePrefix) {
			filepath := fmt.Sprintf("%s/%s", os.TempDir(), f.Name())
			if err := os.Remove(filepath); err != nil {
				log.Printf("Error removing files: %s", err)
				return nil
			}
			log.Printf("[INFO] cloud functions sweeper removed old file %s", filepath)
		}
	}
	return nil
}

func testAccCloudFunctionsFunction_basic(functionName string, bucketName string, zipFilePath string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name                  = "%s"
  runtime               = "nodejs10"
  description           = "test function"
  docker_registry       = "CONTAINER_REGISTRY"
  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name
  trigger_http          = true
  timeout               = 61
  entry_point           = "helloGET"
  ingress_settings      = "ALLOW_INTERNAL_ONLY"
  labels = {
    my-label = "my-label-value"
  }
  environment_variables = {
    TEST_ENV_VARIABLE = "test-env-variable-value"
  }
  build_environment_variables = {
    TEST_ENV_VARIABLE = "test-build-env-variable-value"
  }
  max_instances = 10
  min_instances = 3
}
`, bucketName, zipFilePath, functionName)
}

func testAccCloudFunctionsFunction_updated(functionName string, bucketName string, zipFilePath string, randomSuffix string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index_update.zip"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name                         = "%s"
  description                  = "test function updated"
  docker_registry              = "ARTIFACT_REGISTRY"
  docker_repository = google_artifact_registry_repository.my-repo.id
  available_memory_mb          = 256
  source_archive_bucket        = google_storage_bucket.bucket.name
  source_archive_object        = google_storage_bucket_object.archive.name
  trigger_http                 = true
  https_trigger_security_level = "SECURE_ALWAYS"
  runtime                      = "nodejs10"
  timeout                      = 91
  entry_point                  = "helloGET"
  ingress_settings             = "ALLOW_ALL"
  labels = {
    my-label    = "my-updated-label-value"
    a-new-label = "a-new-label-value"
  }
  environment_variables = {
    TEST_ENV_VARIABLE = "test-env-variable-value"
    NEW_ENV_VARIABLE  = "new-env-variable-value"
  }
  build_environment_variables = {
    TEST_ENV_VARIABLE = "test-build-env-variable-value"
    NEW_ENV_VARIABLE  = "new-build-env-variable-value"
  }
  max_instances = 15
  min_instances = 5
  region = "us-central1"
}

resource "google_artifact_registry_repository" "my-repo" {
	location      = "us-central1"
	repository_id = "tf-test-my-repository%s"
	description   = "example docker repository with cmek"
	format        = "DOCKER"
}
`, bucketName, zipFilePath, functionName, randomSuffix)
}

func testAccCloudFunctionsFunction_buildworkerpool(functionName string, bucketName string, zipFilePath string, location string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_cloudbuild_worker_pool" "pool" {
  name     = "pool-%[3]s"
  location = "%s"
  worker_config {
    disk_size_gb   = 100
    machine_type   = "e2-standard-4"
    no_external_ip = false
  }
}

resource "google_cloudfunctions_function" "function" {
  name                  = "%[3]s"
  runtime               = "nodejs10"
  description           = "test function"
  docker_registry       = "CONTAINER_REGISTRY"
  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name
  trigger_http          = true
  timeout               = 61
  entry_point           = "helloGET"
  build_worker_pool		= google_cloudbuild_worker_pool.pool.id
}`, bucketName, zipFilePath, functionName, location)
}

func testAccCloudFunctionsFunction_pubsub(functionName string, bucketName string,
	topic string, zipFilePath string) string {
	return fmt.Sprintf(`
resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
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
  runtime               = "nodejs10"
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
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name                  = "%s"
  runtime               = "nodejs10"
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
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name                  = "%s"
  runtime               = "nodejs10"
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
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = google_storage_bucket.bucket.name
  source = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name                  = "%s"
  runtime               = "nodejs10"
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
  runtime = "nodejs10"

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
  name     = "%s"
  location = "US"
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
  runtime = "nodejs10"

  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name

  service_account_email = data.google_compute_default_service_account.default.email

  trigger_http = true
  entry_point  = "helloGET"
}
`, bucketName, zipFilePath, functionName)
}

func testAccCloudFunctionsFunction_vpcConnector(projectNumber, networkName, functionName, bucketName, zipFilePath, vpcIp, vpcConnectorName string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_project_iam_member" "gcfadmin" {
  project = data.google_project.project.project_id
  role     = "roles/editor"
  member   = "serviceAccount:service-%s@gcf-admin-robot.iam.gserviceaccount.com"
}

resource "google_compute_network" "vpc" {
	name = "%s"
	auto_create_subnetworks = false
}

resource "google_vpc_access_connector" "%s" {
  name          = "%s"
  region        = "us-central1"
  ip_cidr_range = "%s"
  network       = google_compute_network.vpc.name
}

resource "google_storage_bucket" "bucket" {
  name     = "%s"
  location = "US"
}

resource "google_storage_bucket_object" "archive" {
  name     = "index.zip"
  bucket   = google_storage_bucket.bucket.name
  source   = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name     = "%s"
  runtime  = "nodejs10"

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
  min_instances = 3
  vpc_connector = google_vpc_access_connector.%s.self_link
  vpc_connector_egress_settings = "PRIVATE_RANGES_ONLY"

  depends_on = [google_project_iam_member.gcfadmin]
}
`, projectNumber, networkName, vpcConnectorName, vpcConnectorName, vpcIp, bucketName, zipFilePath, functionName, vpcConnectorName)
}

func testAccCloudFunctionsFunction_secretEnvVar(secretName, versionName, bucketName, functionName, versionNumber, zipFilePath, accountId string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_service_account" "cloud_function_runner" {
  account_id   = "%s"
  display_name = "Testing Cloud Function Secrets integration"
}

resource "google_secret_manager_secret" "test_secret" {
  secret_id = "%s"

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
      replicas {
        location = "us-east1"
      }
    }
  }
}

resource "google_secret_manager_secret_version" "%s" {
  secret      = google_secret_manager_secret.test_secret.id
  secret_data = "This is my secret data."
}

resource "google_secret_manager_secret_iam_member" "cloud_function_iam_member" {
  secret_id = google_secret_manager_secret.test_secret.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.cloud_function_runner.email}"
}

resource "google_storage_bucket" "cloud_functions" {
  name                        = "%s"
  location                    = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "cloud_function_zip_object" {
  name   = "cloud-function.zip"
  bucket = google_storage_bucket.cloud_functions.name
  source = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name                  = "%s"
  runtime               = "nodejs14"
  service_account_email = google_service_account.cloud_function_runner.email
  entry_point           = "echoSecret"
  source_archive_bucket = google_storage_bucket.cloud_functions.id
  source_archive_object = google_storage_bucket_object.cloud_function_zip_object.name
  trigger_http          = true
  secret_environment_variables {
    key     = "MY_SECRET"
    secret  = google_secret_manager_secret.test_secret.secret_id
    version = "%s"
  }

}
`, accountId, secretName, versionName, bucketName, zipFilePath, functionName, versionNumber)
}

func testAccCloudFunctionsFunction_secretMount(projectNumber, secretName, versionName, bucketName, functionName, versionNumber, zipFilePath, accountId string) string {
	return fmt.Sprintf(`
data "google_project" "project" {}

resource "google_service_account" "cloud_function_runner" {
  account_id   = "%s"
  display_name = "Testing Cloud Function Secrets integration"
}

resource "google_secret_manager_secret" "test_secret" {
  secret_id = "%s"

  replication {
    user_managed {
      replicas {
        location = "us-central1"
      }
      replicas {
        location = "us-east1"
      }
    }
  }
}

resource "google_secret_manager_secret_version" "%s" {
  secret      = google_secret_manager_secret.test_secret.id
  secret_data = "This is my secret data."
}

resource "google_secret_manager_secret_iam_member" "cloud_function_iam_member" {
  secret_id = google_secret_manager_secret.test_secret.id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.cloud_function_runner.email}"
}

resource "google_storage_bucket" "cloud_functions" {
  name                        = "%s"
  location                    = "US"
  uniform_bucket_level_access = true
}

resource "google_storage_bucket_object" "cloud_function_zip_object" {
  name   = "cloud-function.zip"
  bucket = google_storage_bucket.cloud_functions.name
  source = "%s"
}

resource "google_cloudfunctions_function" "function" {
  name                  = "%s"
  runtime               = "nodejs14"
  service_account_email = google_service_account.cloud_function_runner.email
  entry_point           = "echoSecret"
  source_archive_bucket = google_storage_bucket.cloud_functions.id
  source_archive_object = google_storage_bucket_object.cloud_function_zip_object.name
  trigger_http          = true
  secret_volumes {
    secret     = google_secret_manager_secret.test_secret.secret_id
    mount_path = "/etc/secrets"
    project_id = "%s"
    versions {
      version = "%s"
      path    = "/test-secret"
    }
  }

}
`, accountId, secretName, versionName, bucketName, zipFilePath, functionName, projectNumber, versionNumber)
}
