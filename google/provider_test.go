package google

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

var credsEnvVars = []string{
	"GOOGLE_CREDENTIALS",
	"GOOGLE_CLOUD_KEYFILE_JSON",
	"GCLOUD_KEYFILE_JSON",
	"GOOGLE_USE_DEFAULT_CREDENTIALS",
}

var projectEnvVars = []string{
	"GOOGLE_PROJECT",
	"GCLOUD_PROJECT",
	"CLOUDSDK_CORE_PROJECT",
}

var firestoreProjectEnvVars = []string{
	"GOOGLE_FIRESTORE_PROJECT",
}

var regionEnvVars = []string{
	"GOOGLE_REGION",
	"GCLOUD_REGION",
	"CLOUDSDK_COMPUTE_REGION",
}

var zoneEnvVars = []string{
	"GOOGLE_ZONE",
	"GCLOUD_ZONE",
	"CLOUDSDK_COMPUTE_ZONE",
}

var orgEnvVars = []string{
	"GOOGLE_ORG",
}

var orgEnvDomainVars = []string{
	"GOOGLE_ORG_DOMAIN",
}

var serviceAccountEnvVars = []string{
	"GOOGLE_SERVICE_ACCOUNT",
}

var orgTargetEnvVars = []string{
	"GOOGLE_ORG_2",
}

var billingAccountEnvVars = []string{
	"GOOGLE_BILLING_ACCOUNT",
}

var configs map[string]*Config

// A source for a given VCR test with the value that seeded it
type VcrSource struct {
	seed   int64
	source rand.Source
}

var sources map[string]VcrSource

func init() {
	configs = make(map[string]*Config)
	sources = make(map[string]VcrSource)
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"google": testAccProvider,
	}
}

// Returns a cached config if VCR testing is enabled. This enables us to use a single HTTP transport
// for a given test, allowing for recording of HTTP interactions.
// Why this exists: schema.Provider.ConfigureFunc is called multiple times for a given test
// ConfigureFunc on our provider creates a new HTTP client and sets base paths (config.go LoadAndValidate)
// VCR requires a single HTTP client to handle all interactions so it can record and replay responses so
// this caches HTTP clients per test by replacing ConfigureFunc
func getCachedConfig(d *schema.ResourceData, configureFunc func(d *schema.ResourceData) (interface{}, error), testName string) (*Config, error) {
	if v, ok := configs[testName]; ok {
		return v, nil
	}
	c, err := configureFunc(d)
	if err != nil {
		return nil, err
	}
	config := c.(*Config)
	var vcrMode recorder.Mode
	switch vcrEnv := os.Getenv("VCR_MODE"); vcrEnv {
	case "RECORDING":
		vcrMode = recorder.ModeRecording
	case "REPLAYING":
		vcrMode = recorder.ModeReplaying
		// When replaying, set the poll interval low to speed up tests
		config.PollInterval = 10 * time.Millisecond
	default:
		log.Printf("[DEBUG] No valid environment var set for VCR_MODE, expected RECORDING or REPLAYING, skipping VCR. VCR_MODE: %s", vcrEnv)
		return config, nil
	}

	envPath := os.Getenv("VCR_PATH")
	if envPath == "" {
		log.Print("[DEBUG] No environment var set for VCR_PATH, skipping VCR")
		return config, nil
	}
	path := filepath.Join(envPath, vcrFileName(testName))

	rec, err := recorder.NewAsMode(path, vcrMode, config.client.Transport)
	if err != nil {
		return nil, err
	}
	// Defines how VCR will match requests to responses.
	rec.SetMatcher(func(r *http.Request, i cassette.Request) bool {
		// Default matcher compares method and URL only
		if !cassette.DefaultMatcher(r, i) {
			return false
		}
		if r.Body == nil {
			return true
		}
		contentType := r.Header.Get("Content-Type")
		// If body contains media, don't try to compare
		if strings.Contains(contentType, "multipart/related") {
			return true
		}

		var b bytes.Buffer
		if _, err := b.ReadFrom(r.Body); err != nil {
			log.Printf("[DEBUG] Failed to read request body from cassette: %v", err)
			return false
		}
		r.Body = ioutil.NopCloser(&b)
		reqBody := b.String()
		// If body matches identically, we are done
		if reqBody == i.Body {
			return true
		}

		// JSON might be the same, but reordered. Try parsing json and comparing
		if strings.Contains(contentType, "application/json") {
			var reqJson, cassetteJson interface{}
			if err := json.Unmarshal([]byte(reqBody), &reqJson); err != nil {
				log.Printf("[DEBUG] Failed to unmarshall request json: %v", err)
				return false
			}
			if err := json.Unmarshal([]byte(i.Body), &cassetteJson); err != nil {
				log.Printf("[DEBUG] Failed to unmarshall cassette json: %v", err)
				return false
			}
			return reflect.DeepEqual(reqJson, cassetteJson)
		}
		return false
	})
	config.client.Transport = rec
	config.wrappedPubsubClient.Transport = rec
	config.wrappedBigQueryClient.Transport = rec
	configs[testName] = config
	return config, err
}

// We need to explicitly close the VCR recorder to save the cassette
func closeRecorder(t *testing.T) {
	if config, ok := configs[t.Name()]; ok {
		// We did not cache the config if it does not use VCR
		if !t.Failed() && isVcrEnabled() {
			// If a test succeeds, write new seed/yaml to files
			err := config.client.Transport.(*recorder.Recorder).Stop()
			if err != nil {
				t.Error(err)
			}
			envPath := os.Getenv("VCR_PATH")
			if vcrSource, ok := sources[t.Name()]; ok {
				err = writeSeedToFile(vcrSource.seed, vcrSeedFile(envPath, t.Name()))
				if err != nil {
					t.Error(err)
				}
			}
		}
		// Clean up test config
		delete(configs, t.Name())
		delete(sources, t.Name())
	}
}

func googleProviderConfig(t *testing.T) *Config {
	config, ok := configs[t.Name()]
	if ok {
		return config
	}
	return testAccProvider.Meta().(*Config)
}

func getTestAccProviders(testName string) map[string]*schema.Provider {
	prov := Provider()
	if isVcrEnabled() {
		old := prov.ConfigureFunc
		prov.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
			return getCachedConfig(d, old, testName)
		}
	} else {
		log.Print("[DEBUG] VCR_PATH or VCR_MODE not set, skipping VCR")
	}
	return map[string]*schema.Provider{
		"google":      prov,
		"google-beta": prov,
	}
}

func isVcrEnabled() bool {
	envPath := os.Getenv("VCR_PATH")
	vcrMode := os.Getenv("VCR_MODE")
	return envPath != "" && vcrMode != ""
}

// Wrapper for resource.Test to swap out providers for VCR providers and handle VCR specific things
// Can be called when VCR is not enabled, and it will behave as normal
func vcrTest(t *testing.T, c resource.TestCase) {
	if isVcrEnabled() {
		providers := getTestAccProviders(t.Name())
		c.Providers = providers
		defer closeRecorder(t)
	}
	resource.Test(t, c)
}

// Retrieves a unique test name used for writing files
// replaces all `/` characters that would cause filepath issues
// This matters during tests that dispatch multiple tests, for example TestAccLoggingFolderExclusion
func vcrSeedFile(path, name string) string {
	return filepath.Join(path, fmt.Sprintf("%s.seed", vcrFileName(name)))
}

func vcrFileName(name string) string {
	return strings.ReplaceAll(name, "/", "_")
}

// Produces a rand.Source for VCR testing based on the given mode.
// In RECORDING mode, generates a new seed and saves it to a file, using the seed for the source
// In REPLAYING mode, reads a seed from a file and creates a source from it
func vcrSource(t *testing.T, path, mode string) (*VcrSource, error) {
	if s, ok := sources[t.Name()]; ok {
		return &s, nil
	}
	switch mode {
	case "RECORDING":
		seed := rand.Int63()
		s := rand.NewSource(seed)
		vcrSource := VcrSource{seed: seed, source: s}
		sources[t.Name()] = vcrSource
		return &vcrSource, nil
	case "REPLAYING":
		seed, err := readSeedFromFile(vcrSeedFile(path, t.Name()))
		if err != nil {
			return nil, err
		}
		s := rand.NewSource(seed)
		vcrSource := VcrSource{seed: seed, source: s}
		sources[t.Name()] = vcrSource
		return &vcrSource, nil
	default:
		log.Printf("[DEBUG] No valid environment var set for VCR_MODE, expected RECORDING or REPLAYING, skipping VCR. VCR_MODE: %s", mode)
		return nil, errors.New("No valid VCR_MODE set")
	}
}

func readSeedFromFile(fileName string) (int64, error) {
	// Max number of digits for int64 is 19
	data := make([]byte, 19)
	f, err := os.Open(fileName)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	_, err = f.Read(data)
	if err != nil {
		return 0, err
	}
	// Remove NULL characters from seed
	data = bytes.Trim(data, "\x00")
	seed := string(data)
	return strconv.ParseInt(seed, 10, 64)
}

func writeSeedToFile(seed int64, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(strconv.FormatInt(seed, 10))
	if err != nil {
		return err
	}
	return nil
}

func randString(t *testing.T, length int) string {
	if !isVcrEnabled() {
		return acctest.RandString(length)
	}
	envPath := os.Getenv("VCR_PATH")
	vcrMode := os.Getenv("VCR_MODE")
	s, err := vcrSource(t, envPath, vcrMode)
	if err != nil {
		// At this point we haven't created any resources, so fail fast
		t.Fatal(err)
	}

	r := rand.New(s.source)
	result := make([]byte, length)
	set := "abcdefghijklmnopqrstuvwxyz012346789"
	for i := 0; i < length; i++ {
		result[i] = set[r.Intn(len(set))]
	}
	return string(result)
}

func randInt(t *testing.T) int {
	if !isVcrEnabled() {
		return acctest.RandInt()
	}
	envPath := os.Getenv("VCR_PATH")
	vcrMode := os.Getenv("VCR_MODE")
	s, err := vcrSource(t, envPath, vcrMode)
	if err != nil {
		// At this point we haven't created any resources, so fail fast
		t.Fatal(err)
	}

	return rand.New(s.source).Int()
}

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

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("GOOGLE_CREDENTIALS_FILE"); v != "" {
		creds, err := ioutil.ReadFile(v)
		if err != nil {
			t.Fatalf("Error reading GOOGLE_CREDENTIALS_FILE path: %s", err)
		}
		os.Setenv("GOOGLE_CREDENTIALS", string(creds))
	}

	if v := multiEnvSearch(credsEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(credsEnvVars, ", "))
	}

	if v := multiEnvSearch(projectEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(projectEnvVars, ", "))
	}

	if v := multiEnvSearch(regionEnvVars); v != "us-central1" {
		t.Fatalf("One of %s must be set to us-central1 for acceptance tests", strings.Join(regionEnvVars, ", "))
	}

	if v := multiEnvSearch(zoneEnvVars); v != "us-central1-a" {
		t.Fatalf("One of %s must be set to us-central1-a for acceptance tests", strings.Join(zoneEnvVars, ", "))
	}
}

func TestProvider_getRegionFromZone(t *testing.T) {
	expected := "us-central1"
	actual := getRegionFromZone("us-central1-f")
	if expected != actual {
		t.Fatalf("Region (%s) did not match expected value: %s", actual, expected)
	}
}

func TestProvider_loadCredentialsFromFile(t *testing.T) {
	ws, es := validateCredentials(testFakeCredentialsPath, "")
	if len(ws) != 0 {
		t.Errorf("Expected %d warnings, got %v", len(ws), ws)
	}
	if len(es) != 0 {
		t.Errorf("Expected %d errors, got %v", len(es), es)
	}
}

func TestProvider_loadCredentialsFromJSON(t *testing.T) {
	contents, err := ioutil.ReadFile(testFakeCredentialsPath)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	ws, es := validateCredentials(string(contents), "")
	if len(ws) != 0 {
		t.Errorf("Expected %d warnings, got %v", len(ws), ws)
	}
	if len(es) != 0 {
		t.Errorf("Expected %d errors, got %v", len(es), es)
	}
}

func TestAccProviderBasePath_setBasePath(t *testing.T) {
	t.Parallel()

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProviderBasePath_setBasePath("https://www.googleapis.com/compute/beta/", randString(t, 10)),
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

	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccProviderBasePath_setBasePath("https://www.example.com/compute/beta/", randString(t, 10)),
				ExpectError: regexp.MustCompile("got HTTP response code 404 with body"),
			},
		},
	})
}

func TestAccProviderUserProjectOverride(t *testing.T) {
	// Parallel fine-grained resource creation
	skipIfVcr(t)
	t.Parallel()

	org := getTestOrgFromEnv(t)
	billing := getTestBillingAccountFromEnv(t)
	pid := "tf-test-" + randString(t, 10)
	sa := "tf-test-" + randString(t, 10)
	topicName := "tf-test-topic-" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// No TestDestroy since that's not really the point of this test
		Steps: []resource.TestStep{
			{
				Config: testAccProviderUserProjectOverride(pid, pname, org, billing, sa),
				Check: func(s *terraform.State) error {
					// The token creator IAM API call returns success long before the policy is
					// actually usable. Wait a solid 2 minutes to ensure we can use it.
					time.Sleep(2 * time.Minute)
					return nil
				},
			},
			{
				Config:      testAccProviderUserProjectOverride_step2(pid, pname, org, billing, sa, false, topicName),
				ExpectError: regexp.MustCompile("Cloud Pub/Sub API has not been used"),
			},
			{
				Config: testAccProviderUserProjectOverride_step2(pid, pname, org, billing, sa, true, topicName),
			},
			{
				ResourceName:      "google_pubsub_topic.project-2-topic",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccProviderUserProjectOverride_step3(pid, pname, org, billing, sa, true),
			},
		},
	})
}

// Do the same thing as TestAccProviderUserProjectOverride, but using a resource that gets its project via
// a reference to a different resource instead of a project field.
func TestAccProviderIndirectUserProjectOverride(t *testing.T) {
	// Parallel fine-grained resource creation
	skipIfVcr(t)
	t.Parallel()

	org := getTestOrgFromEnv(t)
	billing := getTestBillingAccountFromEnv(t)
	pid := "tf-test-" + randString(t, 10)
	sa := "tf-test-" + randString(t, 10)

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		// No TestDestroy since that's not really the point of this test
		Steps: []resource.TestStep{
			{
				Config: testAccProviderIndirectUserProjectOverride(pid, pname, org, billing, sa),
				Check: func(s *terraform.State) error {
					// The token creator IAM API call returns success long before the policy is
					// actually usable. Wait a solid 2 minutes to ensure we can use it.
					time.Sleep(2 * time.Minute)
					return nil
				},
			},
			{
				Config:      testAccProviderIndirectUserProjectOverride_step2(pid, pname, org, billing, sa, false),
				ExpectError: regexp.MustCompile(`Cloud Key Management Service \(KMS\) API has not been used`),
			},
			{
				Config: testAccProviderIndirectUserProjectOverride_step2(pid, pname, org, billing, sa, true),
			},
			{
				ResourceName:      "google_kms_crypto_key.project-2-key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccProviderIndirectUserProjectOverride_step3(pid, pname, org, billing, sa, true),
			},
		},
	})
}

func testAccProviderBasePath_setBasePath(endpoint, name string) string {
	return fmt.Sprintf(`
provider "google" {
  compute_custom_endpoint = "%s"
}

resource "google_compute_address" "default" {
	name = "address-test-%s"
}`, endpoint, name)
}

// Set up two projects. Project 1 has a service account that is used to create a
// pubsub topic in project 2. The pubsub API is only enabled in project 2,
// which causes the create to fail unless user_project_override is set to true.
func testAccProviderUserProjectOverride(pid, name, org, billing, sa string) string {
	return fmt.Sprintf(`
resource "google_project" "project-1" {
	project_id      = "%s"
	name            = "%s"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_service_account" "project-1" {
	project    = google_project.project-1.project_id
    account_id = "%s"
}

resource "google_project" "project-2" {
	project_id      = "%s-2"
	name            = "%s-2"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_service" "project-2-pubsub-service" {
	project = google_project.project-2.project_id
	service = "pubsub.googleapis.com"
}

// Permission needed for user_project_override
resource "google_project_iam_member" "project-2-serviceusage" {
	project = google_project.project-2.project_id
	role    = "roles/serviceusage.serviceUsageConsumer"
	member  = "serviceAccount:${google_service_account.project-1.email}"
}

resource "google_project_iam_member" "project-2-pubsub-member" {
	project = google_project.project-2.project_id
	role    = "roles/pubsub.admin"
	member  = "serviceAccount:${google_service_account.project-1.email}"
}

data "google_client_openid_userinfo" "me" {}

// Enable the test runner to get an access token on behalf of
// the project 1 service account
resource "google_service_account_iam_member" "token-creator-iam" {
	service_account_id = google_service_account.project-1.name
	role               = "roles/iam.serviceAccountTokenCreator"
	member             = "serviceAccount:${data.google_client_openid_userinfo.me.email}"
}
`, pid, name, org, billing, sa, pid, name, org, billing)
}

func testAccProviderUserProjectOverride_step2(pid, name, org, billing, sa string, override bool, topicName string) string {
	return fmt.Sprintf(`
// See step 3 below, which is really step 2 minus the pubsub topic.
// Step 3 exists because provider configurations can't be removed while objects
// created by that provider still exist in state. Step 3 will remove the
// pubsub topic so the whole config can be deleted.
%s

resource "google_pubsub_topic" "project-2-topic" {
	provider = google.project-1-token
	project  = google_project.project-2.project_id

	name = "%s"
	labels = {
	  foo = "bar"
	}
}
`, testAccProviderUserProjectOverride_step3(pid, name, org, billing, sa, override), topicName)
}

func testAccProviderUserProjectOverride_step3(pid, name, org, billing, sa string, override bool) string {
	return fmt.Sprintf(`
%s

data "google_service_account_access_token" "project-1-token" {
	// This data source would have a depends_on t
	// google_service_account_iam_binding.token-creator-iam, but depends_on
	// in data sources makes them always have a diff in apply:
	// https://www.terraform.io/docs/configuration/data-sources.html#data-resource-dependencies
	// Instead, rely on the other test step completing before this one.

	target_service_account = google_service_account.project-1.email
	scopes = ["userinfo-email", "https://www.googleapis.com/auth/cloud-platform"]
	lifetime = "300s"
}

provider "google" {
	alias  = "project-1-token"
	access_token = data.google_service_account_access_token.project-1-token.access_token
	user_project_override = %v
}
`, testAccProviderUserProjectOverride(pid, name, org, billing, sa), override)
}

// Set up two projects. Project 1 has a service account that is used to create a
// kms crypto key in project 2. The kms API is only enabled in project 2,
// which causes the create to fail unless user_project_override is set to true.
func testAccProviderIndirectUserProjectOverride(pid, name, org, billing, sa string) string {
	return fmt.Sprintf(`
resource "google_project" "project-1" {
	project_id      = "%s"
	name            = "%s"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_service_account" "project-1" {
	project    = google_project.project-1.project_id
    account_id = "%s"
}

resource "google_project" "project-2" {
	project_id      = "%s-2"
	name            = "%s-2"
	org_id          = "%s"
	billing_account = "%s"
}

resource "google_project_service" "project-2-kms" {
	project = google_project.project-2.project_id
	service = "cloudkms.googleapis.com"
}

// Permission needed for user_project_override
resource "google_project_iam_member" "project-2-serviceusage" {
	project = google_project.project-2.project_id
	role    = "roles/serviceusage.serviceUsageConsumer"
	member  = "serviceAccount:${google_service_account.project-1.email}"
}

resource "google_project_iam_member" "project-2-kms" {
	project = google_project.project-2.project_id
	role    = "roles/cloudkms.admin"
	member  = "serviceAccount:${google_service_account.project-1.email}"
}

resource "google_project_iam_member" "project-2-kms-encrypt" {
	project = google_project.project-2.project_id
	role    = "roles/cloudkms.cryptoKeyEncrypter"
	member  = "serviceAccount:${google_service_account.project-1.email}"
}

data "google_client_openid_userinfo" "me" {}

// Enable the test runner to get an access token on behalf of
// the project 1 service account
resource "google_service_account_iam_member" "token-creator-iam" {
	service_account_id = google_service_account.project-1.name
	role               = "roles/iam.serviceAccountTokenCreator"
	member             = "serviceAccount:${data.google_client_openid_userinfo.me.email}"
}
`, pid, name, org, billing, sa, pid, name, org, billing)
}

func testAccProviderIndirectUserProjectOverride_step2(pid, name, org, billing, sa string, override bool) string {
	return fmt.Sprintf(`
// See step 3 below, which is really step 2 minus the kms resources.
// Step 3 exists because provider configurations can't be removed while objects
// created by that provider still exist in state. Step 3 will remove the
// kms resources so the whole config can be deleted.
%s

resource "google_kms_key_ring" "project-2-keyring" {
	provider = google.project-1-token
	project  = google_project.project-2.project_id

	name     = "%s"
	location = "us-central1"
}

resource "google_kms_crypto_key" "project-2-key" {
	provider = google.project-1-token
	name     = "%s"
	key_ring = google_kms_key_ring.project-2-keyring.self_link
}

data "google_kms_secret_ciphertext" "project-2-ciphertext" {
	provider   = google.project-1-token
	crypto_key = google_kms_crypto_key.project-2-key.self_link
	plaintext  = "my-secret"
}
`, testAccProviderIndirectUserProjectOverride_step3(pid, name, org, billing, sa, override), pid, pid)
}

func testAccProviderIndirectUserProjectOverride_step3(pid, name, org, billing, sa string, override bool) string {
	return fmt.Sprintf(`
%s

data "google_service_account_access_token" "project-1-token" {
	// This data source would have a depends_on to
	// google_service_account_iam_binding.token-creator-iam, but depends_on
	// in data sources makes them always have a diff in apply:
	// https://www.terraform.io/docs/configuration/data-sources.html#data-resource-dependencies
	// Instead, rely on the other test step completing before this one.

	target_service_account = google_service_account.project-1.email
	scopes                 = ["userinfo-email", "https://www.googleapis.com/auth/cloud-platform"]
	lifetime               = "300s"
}

provider "google" {
	alias = "project-1-token"

	access_token          = data.google_service_account_access_token.project-1-token.access_token
	user_project_override = %v
}
`, testAccProviderIndirectUserProjectOverride(pid, name, org, billing, sa), override)
}

// getTestRegion has the same logic as the provider's getRegion, to be used in tests.
func getTestRegion(is *terraform.InstanceState, config *Config) (string, error) {
	if res, ok := is.Attributes["region"]; ok {
		return res, nil
	}
	if config.Region != "" {
		return config.Region, nil
	}
	return "", fmt.Errorf("%q: required field is not set", "region")
}

// getTestProject has the same logic as the provider's getProject, to be used in tests.
func getTestProject(is *terraform.InstanceState, config *Config) (string, error) {
	if res, ok := is.Attributes["project"]; ok {
		return res, nil
	}
	if config.Project != "" {
		return config.Project, nil
	}
	return "", fmt.Errorf("%q: required field is not set", "project")
}

// testAccPreCheck ensures at least one of the project env variables is set.
func getTestProjectFromEnv() string {
	return multiEnvSearch(projectEnvVars)
}

// testAccPreCheck ensures at least one of the credentials env variables is set.
func getTestCredsFromEnv() string {
	return multiEnvSearch(credsEnvVars)
}

// testAccPreCheck ensures at least one of the region env variables is set.
func getTestRegionFromEnv() string {
	return multiEnvSearch(regionEnvVars)
}

func getTestZoneFromEnv() string {
	return multiEnvSearch(zoneEnvVars)
}

// Firestore can't be enabled at the same time as Datastore, so we need a new
// project to manage it until we can enable Firestore programmatically.
func getTestFirestoreProjectFromEnv(t *testing.T) string {
	skipIfEnvNotSet(t, firestoreProjectEnvVars...)
	return multiEnvSearch(firestoreProjectEnvVars)
}

func getTestOrgFromEnv(t *testing.T) string {
	skipIfEnvNotSet(t, orgEnvVars...)
	return multiEnvSearch(orgEnvVars)
}

func getTestOrgDomainFromEnv(t *testing.T) string {
	skipIfEnvNotSet(t, orgEnvDomainVars...)
	return multiEnvSearch(orgEnvDomainVars)
}

func getTestOrgTargetFromEnv(t *testing.T) string {
	skipIfEnvNotSet(t, orgTargetEnvVars...)
	return multiEnvSearch(orgTargetEnvVars)
}

func getTestBillingAccountFromEnv(t *testing.T) string {
	skipIfEnvNotSet(t, billingAccountEnvVars...)
	return multiEnvSearch(billingAccountEnvVars)
}

func getTestServiceAccountFromEnv(t *testing.T) string {
	skipIfEnvNotSet(t, serviceAccountEnvVars...)
	return multiEnvSearch(serviceAccountEnvVars)
}

func multiEnvSearch(ks []string) string {
	for _, k := range ks {
		if v := os.Getenv(k); v != "" {
			return v
		}
	}
	return ""
}

// Some tests fail during VCR. One common case is race conditions when creating resources.
// If a test config adds two fine-grained resources with the same parent it is undefined
// which will be created first, causing VCR to fail ~50% of the time
func skipIfVcr(t *testing.T) {
	if isVcrEnabled() {
		t.Skipf("VCR enabled, skipping test: %s", t.Name())
	}
}
