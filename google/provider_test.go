package google

import (
	"bytes"
	"context"
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
	"sync"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/serviceusage/v1"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

var credsEnvVars = []string{
	"GOOGLE_CREDENTIALS",
	"GOOGLE_CLOUD_KEYFILE_JSON",
	"GCLOUD_KEYFILE_JSON",
	"GOOGLE_APPLICATION_CREDENTIALS",
	"GOOGLE_USE_DEFAULT_CREDENTIALS",
}

var projectNumberEnvVars = []string{
	"GOOGLE_PROJECT_NUMBER",
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

// This value is the Customer ID of the GOOGLE_ORG_DOMAIN workspace.
// See https://admin.google.com/ac/accountsettings when logged into an org admin for the value.
var custIdEnvVars = []string{
	"GOOGLE_CUST_ID",
}

// This value is the username of an identity account within the GOOGLE_ORG_DOMAIN workspace.
// For example in the org example.com with a user "foo@example.com", this would be set to "foo".
// See https://admin.google.com/ac/users when logged into an org admin for a list.
var identityUserEnvVars = []string{
	"GOOGLE_IDENTITY_USER",
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

var masterBillingAccountEnvVars = []string{
	"GOOGLE_MASTER_BILLING_ACCOUNT",
}

var configsLock = sync.RWMutex{}
var sourcesLock = sync.RWMutex{}

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
func getCachedConfig(ctx context.Context, d *schema.ResourceData, configureFunc schema.ConfigureContextFunc, testName string) (*Config, diag.Diagnostics) {
	configsLock.RLock()
	v, ok := configs[testName]
	configsLock.RUnlock()
	if ok {
		return v, nil
	}
	c, diags := configureFunc(ctx, d)
	if diags.HasError() {
		return nil, diags
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
		return nil, diag.FromErr(err)
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
	configsLock.Lock()
	configs[testName] = config
	configsLock.Unlock()
	return config, nil
}

// We need to explicitly close the VCR recorder to save the cassette
func closeRecorder(t *testing.T) {
	configsLock.RLock()
	config, ok := configs[t.Name()]
	configsLock.RUnlock()
	if ok {
		// We did not cache the config if it does not use VCR
		if !t.Failed() && isVcrEnabled() {
			// If a test succeeds, write new seed/yaml to files
			err := config.client.Transport.(*recorder.Recorder).Stop()
			if err != nil {
				t.Error(err)
			}
			envPath := os.Getenv("VCR_PATH")

			sourcesLock.RLock()
			vcrSource, ok := sources[t.Name()]
			sourcesLock.RUnlock()
			if ok {
				err = writeSeedToFile(vcrSource.seed, vcrSeedFile(envPath, t.Name()))
				if err != nil {
					t.Error(err)
				}
			}
		}
		// Clean up test config
		configsLock.Lock()
		delete(configs, t.Name())
		configsLock.Unlock()

		sourcesLock.Lock()
		delete(sources, t.Name())
		sourcesLock.Unlock()
	}
}

func googleProviderConfig(t *testing.T) *Config {
	configsLock.RLock()
	config, ok := configs[t.Name()]
	configsLock.RUnlock()
	if ok {
		return config
	}
	return testAccProvider.Meta().(*Config)
}

func getTestAccProviders(testName string, c resource.TestCase) map[string]*schema.Provider {
	prov := Provider()
	if isVcrEnabled() {
		old := prov.ConfigureContextFunc
		prov.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			return getCachedConfig(ctx, d, old, testName)
		}
	} else {
		log.Print("[DEBUG] VCR_PATH or VCR_MODE not set, skipping VCR")
	}
	var testProvider string
	providerMapKeys := reflect.ValueOf(c.Providers).MapKeys()
	if strings.Contains(providerMapKeys[0].String(), "google-beta") {
		testProvider = "google-beta"
	} else {
		testProvider = "google"
	}
	return map[string]*schema.Provider{
		testProvider: prov,
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
		providers := getTestAccProviders(t.Name(), c)
		c.Providers = providers
		defer closeRecorder(t)
	} else if isReleaseDiffEnabled() {
		c = initializeReleaseDiffTest(c)
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
	sourcesLock.RLock()
	s, ok := sources[t.Name()]
	sourcesLock.RUnlock()
	if ok {
		return &s, nil
	}
	switch mode {
	case "RECORDING":
		seed := rand.Int63()
		s := rand.NewSource(seed)
		vcrSource := VcrSource{seed: seed, source: s}
		sourcesLock.Lock()
		sources[t.Name()] = vcrSource
		sourcesLock.Unlock()
		return &vcrSource, nil
	case "REPLAYING":
		seed, err := readSeedFromFile(vcrSeedFile(path, t.Name()))
		if err != nil {
			return nil, fmt.Errorf("no cassette found on disk for %s, please replay this testcase in recording mode - %w", t.Name(), err)
		}
		s := rand.NewSource(seed)
		vcrSource := VcrSource{seed: seed, source: s}
		sourcesLock.Lock()
		sources[t.Name()] = vcrSource
		sourcesLock.Unlock()
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
	return stringToFixed64(seed)
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

	if v := multiEnvSearch(regionEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(regionEnvVars, ", "))
	}

	if v := multiEnvSearch(zoneEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(zoneEnvVars, ", "))
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

func TestAccProviderMeta_setModuleName(t *testing.T) {
	t.Parallel()

	moduleName := "my-module"
	vcrTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeAddressDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccProviderMeta_setModuleName(moduleName, randString(t, 10)),
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
	skipIfVcr(t)
	t.Parallel()

	org := getTestOrgFromEnv(t)
	billing := getTestBillingAccountFromEnv(t)
	pid := "tf-test-" + randString(t, 10)
	topicName := "tf-test-topic-" + randString(t, 10)

	config := BootstrapConfig(t)
	accessToken, err := setupProjectsAndGetAccessToken(org, billing, pid, "pubsub", config)
	if err != nil {
		t.Error(err)
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
	skipIfVcr(t)
	t.Parallel()

	org := getTestOrgFromEnv(t)
	billing := getTestBillingAccountFromEnv(t)
	pid := "tf-test-" + randString(t, 10)

	config := BootstrapConfig(t)
	accessToken, err := setupProjectsAndGetAccessToken(org, billing, pid, "cloudkms", config)
	if err != nil {
		t.Error(err)
	}

	vcrTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
  name     = "address-test-%s"
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
	name = "address-test-%s"
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
func getTestProjectNumberFromEnv() string {
	return multiEnvSearch(projectNumberEnvVars)
}

// testAccPreCheck ensures at least one of the project env variables is set.
func getTestProjectFromEnv() string {
	return multiEnvSearch(projectEnvVars)
}

// testAccPreCheck ensures at least one of the credentials env variables is set.
func getTestCredsFromEnv() string {
	// Return empty string if GOOGLE_USE_DEFAULT_CREDENTIALS is set to true.
	if multiEnvSearch(credsEnvVars) == "true" {
		return ""
	}
	return multiEnvSearch(credsEnvVars)
}

// testAccPreCheck ensures at least one of the region env variables is set.
func getTestRegionFromEnv() string {
	return multiEnvSearch(regionEnvVars)
}

func getTestZoneFromEnv() string {
	return multiEnvSearch(zoneEnvVars)
}

func getTestCustIdFromEnv(t *testing.T) string {
	skipIfEnvNotSet(t, custIdEnvVars...)
	return multiEnvSearch(custIdEnvVars)
}

func getTestIdentityUserFromEnv(t *testing.T) string {
	skipIfEnvNotSet(t, identityUserEnvVars...)
	return multiEnvSearch(identityUserEnvVars)
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

func getTestMasterBillingAccountFromEnv(t *testing.T) string {
	skipIfEnvNotSet(t, masterBillingAccountEnvVars...)
	return multiEnvSearch(masterBillingAccountEnvVars)
}

func getTestServiceAccountFromEnv(t *testing.T) string {
	skipIfEnvNotSet(t, serviceAccountEnvVars...)
	return multiEnvSearch(serviceAccountEnvVars)
}

// Some tests fail during VCR. One common case is race conditions when creating resources.
// If a test config adds two fine-grained resources with the same parent it is undefined
// which will be created first, causing VCR to fail ~50% of the time
func skipIfVcr(t *testing.T) {
	if isVcrEnabled() {
		t.Skipf("VCR enabled, skipping test: %s", t.Name())
	}
}

func sleepInSecondsForTest(t int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		time.Sleep(time.Duration(t) * time.Second)
		return nil
	}
}

func setupProjectsAndGetAccessToken(org, billing, pid, service string, config *Config) (string, error) {
	// Create project-1 and project-2
	rmService := config.NewResourceManagerClient(config.userAgent)

	project := &cloudresourcemanager.Project{
		ProjectId: pid,
		Name:      pname,
		Parent: &cloudresourcemanager.ResourceId{
			Id:   org,
			Type: "organization",
		},
	}

	var op *cloudresourcemanager.Operation
	err := retryTimeDuration(func() (reqErr error) {
		op, reqErr = rmService.Projects.Create(project).Do()
		return reqErr
	}, 5*time.Minute)
	if err != nil {
		return "", err
	}

	// Wait for the operation to complete
	opAsMap, err := ConvertToMap(op)
	if err != nil {
		return "", err
	}

	waitErr := resourceManagerOperationWaitTime(config, opAsMap, "creating project", config.userAgent, 5*time.Minute)
	if waitErr != nil {
		return "", waitErr
	}

	ba := &cloudbilling.ProjectBillingInfo{
		BillingAccountName: fmt.Sprintf("billingAccounts/%s", billing),
	}
	_, err = config.NewBillingClient(config.userAgent).Projects.UpdateBillingInfo(prefixedProject(pid), ba).Do()
	if err != nil {
		return "", err
	}

	p2 := fmt.Sprintf("%s-2", pid)
	project.ProjectId = p2
	project.Name = fmt.Sprintf("%s-2", pname)

	err = retryTimeDuration(func() (reqErr error) {
		op, reqErr = rmService.Projects.Create(project).Do()
		return reqErr
	}, 5*time.Minute)
	if err != nil {
		return "", err
	}

	// Wait for the operation to complete
	opAsMap, err = ConvertToMap(op)
	if err != nil {
		return "", err
	}

	waitErr = resourceManagerOperationWaitTime(config, opAsMap, "creating project", config.userAgent, 5*time.Minute)
	if waitErr != nil {
		return "", waitErr
	}

	_, err = config.NewBillingClient(config.userAgent).Projects.UpdateBillingInfo(prefixedProject(p2), ba).Do()
	if err != nil {
		return "", err
	}

	// Enable the appropriate service in project-2 only
	suService := config.NewServiceUsageClient(config.userAgent)

	serviceReq := &serviceusage.BatchEnableServicesRequest{
		ServiceIds: []string{fmt.Sprintf("%s.googleapis.com", service)},
	}

	_, err = suService.Services.BatchEnable(fmt.Sprintf("projects/%s", p2), serviceReq).Do()
	if err != nil {
		return "", err
	}

	// Enable the test runner to create service accounts and get an access token on behalf of
	// the project 1 service account
	curEmail, err := GetCurrentUserEmail(config, config.userAgent)
	if err != nil {
		return "", err
	}

	proj1SATokenCreator := &cloudresourcemanager.Binding{
		Members: []string{fmt.Sprintf("serviceAccount:%s", curEmail)},
		Role:    "roles/iam.serviceAccountTokenCreator",
	}

	proj1SACreator := &cloudresourcemanager.Binding{
		Members: []string{fmt.Sprintf("serviceAccount:%s", curEmail)},
		Role:    "roles/iam.serviceAccountCreator",
	}

	bindings := mergeBindings([]*cloudresourcemanager.Binding{proj1SATokenCreator, proj1SACreator})

	p, err := rmService.Projects.GetIamPolicy(pid,
		&cloudresourcemanager.GetIamPolicyRequest{
			Options: &cloudresourcemanager.GetPolicyOptions{
				RequestedPolicyVersion: iamPolicyVersion,
			},
		}).Do()
	if err != nil {
		return "", err
	}

	p.Bindings = mergeBindings(append(p.Bindings, bindings...))
	_, err = config.NewResourceManagerClient(config.userAgent).Projects.SetIamPolicy(pid,
		&cloudresourcemanager.SetIamPolicyRequest{
			Policy:     p,
			UpdateMask: "bindings,etag,auditConfigs",
		}).Do()
	if err != nil {
		return "", err
	}

	// Create a service account for project-1
	sa1, err := getOrCreateServiceAccount(config, pid)
	if err != nil {
		return "", err
	}

	// Add permissions to service accounts

	// Permission needed for user_project_override
	proj2ServiceUsageBinding := &cloudresourcemanager.Binding{
		Members: []string{fmt.Sprintf("serviceAccount:%s", sa1.Email)},
		Role:    "roles/serviceusage.serviceUsageConsumer",
	}

	// Admin permission for service
	proj2ServiceAdminBinding := &cloudresourcemanager.Binding{
		Members: []string{fmt.Sprintf("serviceAccount:%s", sa1.Email)},
		Role:    fmt.Sprintf("roles/%s.admin", service),
	}

	bindings = mergeBindings([]*cloudresourcemanager.Binding{proj2ServiceUsageBinding, proj2ServiceAdminBinding})

	// For KMS test only
	if service == "cloudkms" {
		proj2CryptoKeyBinding := &cloudresourcemanager.Binding{
			Members: []string{fmt.Sprintf("serviceAccount:%s", sa1.Email)},
			Role:    "roles/cloudkms.cryptoKeyEncrypter",
		}

		bindings = mergeBindings(append(bindings, proj2CryptoKeyBinding))
	}

	p, err = rmService.Projects.GetIamPolicy(p2,
		&cloudresourcemanager.GetIamPolicyRequest{
			Options: &cloudresourcemanager.GetPolicyOptions{
				RequestedPolicyVersion: iamPolicyVersion,
			},
		}).Do()
	if err != nil {
		return "", err
	}

	p.Bindings = mergeBindings(append(p.Bindings, bindings...))
	_, err = config.NewResourceManagerClient(config.userAgent).Projects.SetIamPolicy(p2,
		&cloudresourcemanager.SetIamPolicyRequest{
			Policy:     p,
			UpdateMask: "bindings,etag,auditConfigs",
		}).Do()
	if err != nil {
		return "", err
	}

	// The token creator IAM API call returns success long before the policy is
	// actually usable. Wait a solid 2 minutes to ensure we can use it.
	time.Sleep(2 * time.Minute)

	iamCredsService := config.NewIamCredentialsClient(config.userAgent)
	tokenRequest := &iamcredentials.GenerateAccessTokenRequest{
		Lifetime: "300s",
		Scope:    []string{"https://www.googleapis.com/auth/cloud-platform"},
	}
	atResp, err := iamCredsService.Projects.ServiceAccounts.GenerateAccessToken(fmt.Sprintf("projects/-/serviceAccounts/%s", sa1.Email), tokenRequest).Do()
	if err != nil {
		return "", err
	}

	accessToken := atResp.AccessToken

	return accessToken, nil
}

func isReleaseDiffEnabled() bool {
	releaseDiff := os.Getenv("RELEASE_DIFF")
	return releaseDiff != ""
}

func initializeReleaseDiffTest(c resource.TestCase) resource.TestCase {
	var releaseProvider string
	packagePath := fmt.Sprint(reflect.TypeOf(Config{}).PkgPath())
	if strings.Contains(packagePath, "google-beta") {
		releaseProvider = "google-beta"
	} else {
		releaseProvider = "google"
	}

	if c.ExternalProviders != nil {
		c.ExternalProviders[releaseProvider] = resource.ExternalProvider{}
	} else {
		c.ExternalProviders = map[string]resource.ExternalProvider{
			releaseProvider: {},
		}
	}

	localProviderName := "google-local"
	localProvider := map[string]*schema.Provider{
		localProviderName: testAccProvider,
	}
	c.Providers = localProvider

	var replacementSteps []resource.TestStep
	for _, testStep := range c.Steps {
		if testStep.Config != "" {
			ogConfig := testStep.Config
			testStep.Config = reformConfigWithProvider(ogConfig, localProviderName)
			if testStep.ExpectError == nil && testStep.PlanOnly == false {
				newStep := resource.TestStep{
					Config: reformConfigWithProvider(ogConfig, releaseProvider),
				}
				testStep.PlanOnly = true
				testStep.ExpectNonEmptyPlan = false
				replacementSteps = append(replacementSteps, newStep)
			}
			replacementSteps = append(replacementSteps, testStep)
		} else {
			replacementSteps = append(replacementSteps, testStep)
		}
	}

	c.Steps = replacementSteps

	return c
}

func reformConfigWithProvider(config, provider string) string {
	configBytes := []byte(config)
	providerReplacement := fmt.Sprintf("provider = %s", provider)
	providerReplacementBytes := []byte(providerReplacement)
	providerBlock := regexp.MustCompile(`provider *=.*google-beta.*`)

	if providerBlock.Match(configBytes) {
		return string(providerBlock.ReplaceAll(configBytes, providerReplacementBytes))
	}

	providerReplacement = fmt.Sprintf("${1}\n\t%s", providerReplacement)
	providerReplacementBytes = []byte(providerReplacement)
	resourceHeader := regexp.MustCompile(`(resource .*google_.* .*\w+.*\{.*)`)
	return string(resourceHeader.ReplaceAll(configBytes, providerReplacementBytes))
}
