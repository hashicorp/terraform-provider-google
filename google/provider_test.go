package google

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func TestProvider_noDuplicatesInResourceMap(t *testing.T) {
	_, err := ResourceMapWithErrors()
	if err != nil {
		t.Error(err)
	}
}

func TestProvider_validateCredentials(t *testing.T) {
	cases := map[string]struct {
		ConfigValue      func(t *testing.T) interface{}
		ValueNotProvided bool
		ExpectedWarnings []string
		ExpectedErrors   []error
	}{
		"configuring credentials as a path to a credentials JSON file is valid": {
			ConfigValue: func(t *testing.T) interface{} {
				return testFakeCredentialsPath // Path to a test fixture
			},
		},
		"configuring credentials as a path to a non-existant file is NOT valid": {
			ConfigValue: func(t *testing.T) interface{} {
				return "./this/path/doesnt/exist.json" // Doesn't exist
			},
			ExpectedErrors: []error{
				// As the file doesn't exist, so the function attempts to parse it as a JSON
				errors.New("JSON credentials are not valid: invalid character '.' looking for beginning of value"),
			},
		},
		"configuring credentials as a credentials JSON string is valid": {
			ConfigValue: func(t *testing.T) interface{} {
				contents, err := ioutil.ReadFile(testFakeCredentialsPath)
				if err != nil {
					t.Fatalf("Unexpected error: %s", err)
				}
				return string(contents)
			},
		},
		"configuring credentials as an empty string is valid": {
			ConfigValue: func(t *testing.T) interface{} {
				return ""
			},
		},
		"leaving credentials unconfigured is valid": {
			ValueNotProvided: true,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange
			var configValue interface{}
			if !tc.ValueNotProvided {
				configValue = tc.ConfigValue(t)
			}

			// Act
			// Note: second argument is currently unused by the function but is necessary to fulfill the SchemaValidateFunc type's function signature
			ws, es := validateCredentials(configValue, "")

			// Assert
			if len(ws) != len(tc.ExpectedWarnings) {
				t.Errorf("Expected %d warnings, got %d: %v", len(tc.ExpectedWarnings), len(ws), ws)
			}
			if len(es) != len(tc.ExpectedErrors) {
				t.Errorf("Expected %d errors, got %d: %v", len(tc.ExpectedErrors), len(es), es)
			}

			if len(tc.ExpectedErrors) > 0 {
				if es[0].Error() != tc.ExpectedErrors[0].Error() {
					t.Errorf("Expected first error to be \"%s\", got \"%s\"", tc.ExpectedErrors[0], es[0])
				}
			}
		})
	}
}

func TestAccProviderBasePath_setBasePath(t *testing.T) {
	t.Parallel()

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProviderBasePath_setBasePath("https://www.googleapis.com/compute/beta/", RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_address.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccProviderBasePath_setInvalidBasePath(t *testing.T) {
	t.Parallel()

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccProviderBasePath_setBasePath("https://www.example.com/compute/beta/", RandString(t, 10)),
				ExpectError: regexp.MustCompile("got HTTP response code 404 with body"),
			},
		},
	})
}

func TestAccProviderMeta_setModuleName(t *testing.T) {
	t.Parallel()

	moduleName := "my-module"
	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProviderMeta_setModuleName(moduleName, RandString(t, 10)),
			},
			{
				ResourceName:      "google_compute_address.default",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccProviderUserProjectOverride(t *testing.T) {
	// Parallel fine-grained resource creation
	SkipIfVcr(t)
	t.Parallel()

	org := GetTestOrgFromEnv(t)
	billing := GetTestBillingAccountFromEnv(t)
	pid := "tf-test-" + RandString(t, 10)
	topicName := "tf-test-topic-" + RandString(t, 10)

	config := BootstrapConfig(t)
	accessToken, err := setupProjectsAndGetAccessToken(org, billing, pid, "pubsub", config)
	if err != nil {
		t.Error(err)
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		// No TestDestroy since that's not really the point of this test
		Steps: []resource.TestStep{
			{
				Config:      testAccProviderUserProjectOverride_step2(accessToken, pid, false, topicName),
				ExpectError: regexp.MustCompile("Cloud Pub/Sub API has not been used"),
			},
			{
				Config: testAccProviderUserProjectOverride_step2(accessToken, pid, true, topicName),
			},
			{
				ResourceName:      "google_pubsub_topic.project-2-topic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccProviderUserProjectOverride_step3(accessToken, true),
			},
		},
	})
}

// Do the same thing as TestAccProviderUserProjectOverride, but using a resource that gets its project via
// a reference to a different resource instead of a project field.
func TestAccProviderIndirectUserProjectOverride(t *testing.T) {
	// Parallel fine-grained resource creation
	SkipIfVcr(t)
	t.Parallel()

	org := GetTestOrgFromEnv(t)
	billing := GetTestBillingAccountFromEnv(t)
	pid := "tf-test-" + RandString(t, 10)

	config := BootstrapConfig(t)
	accessToken, err := setupProjectsAndGetAccessToken(org, billing, pid, "cloudkms", config)
	if err != nil {
		t.Error(err)
	}

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { AccTestPreCheck(t) },
		ProtoV5ProviderFactories: ProtoV5ProviderFactories(t),
		// No TestDestroy since that's not really the point of this test
		Steps: []resource.TestStep{
			{
				Config:      testAccProviderIndirectUserProjectOverride_step2(pid, accessToken, false),
				ExpectError: regexp.MustCompile(`Cloud Key Management Service \(KMS\) API has not been used`),
			},
			{
				Config: testAccProviderIndirectUserProjectOverride_step2(pid, accessToken, true),
			},
			{
				ResourceName:      "google_kms_crypto_key.project-2-key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccProviderIndirectUserProjectOverride_step3(accessToken, true),
			},
		},
	})
}

func testAccProviderBasePath_setBasePath(endpoint, name string) string {
	return fmt.Sprintf(`
provider "google" {
  alias                   = "compute_custom_endpoint"
  compute_custom_endpoint = "%s"
}

resource "google_compute_address" "default" {
  provider = google.compute_custom_endpoint
  name     = "tf-test-address-%s"
}`, endpoint, name)
}

func testAccProviderMeta_setModuleName(key, name string) string {
	return fmt.Sprintf(`
terraform {
  provider_meta "google" {
    module_name = "%s"
  }
}

resource "google_compute_address" "default" {
	name = "tf-test-address-%s"
}`, key, name)
}

// Set up two projects. Project 1 has a service account that is used to create a
// pubsub topic in project 2. The pubsub API is only enabled in project 2,
// which causes the create to fail unless user_project_override is set to true.

func testAccProviderUserProjectOverride_step2(accessToken, pid string, override bool, topicName string) string {
	return fmt.Sprintf(`
// See step 3 below, which is really step 2 minus the pubsub topic.
// Step 3 exists because provider configurations can't be removed while objects
// created by that provider still exist in state. Step 3 will remove the
// pubsub topic so the whole config can be deleted.
%s

resource "google_pubsub_topic" "project-2-topic" {
	provider = google.project-1-token
	project  = "%s-2"

	name = "%s"
	labels = {
	  foo = "bar"
	}
}
`, testAccProviderUserProjectOverride_step3(accessToken, override), pid, topicName)
}

func testAccProviderUserProjectOverride_step3(accessToken string, override bool) string {
	return fmt.Sprintf(`
provider "google" {
	alias  = "project-1-token"
	access_token = "%s"
	user_project_override = %v
}
`, accessToken, override)
}

func testAccProviderIndirectUserProjectOverride_step2(pid, accessToken string, override bool) string {
	return fmt.Sprintf(`
// See step 3 below, which is really step 2 minus the kms resources.
// Step 3 exists because provider configurations can't be removed while objects
// created by that provider still exist in state. Step 3 will remove the
// kms resources so the whole config can be deleted.
%s

resource "google_kms_key_ring" "project-2-keyring" {
	provider = google.project-1-token
	project  = "%s-2"

	name     = "%s"
	location = "us-central1"
}

resource "google_kms_crypto_key" "project-2-key" {
	provider = google.project-1-token
	name     = "%s"
	key_ring = google_kms_key_ring.project-2-keyring.id
}

data "google_kms_secret_ciphertext" "project-2-ciphertext" {
	provider   = google.project-1-token
	crypto_key = google_kms_crypto_key.project-2-key.id
	plaintext  = "my-secret"
}
`, testAccProviderIndirectUserProjectOverride_step3(accessToken, override), pid, pid, pid)
}

func testAccProviderIndirectUserProjectOverride_step3(accessToken string, override bool) string {
	return fmt.Sprintf(`
provider "google" {
	alias = "project-1-token"

	access_token          = "%s"
	user_project_override = %v
}
`, accessToken, override)
}
