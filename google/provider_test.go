package google

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-random/random"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var testAccRandomProvider *schema.Provider

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

var serviceAccountEnvVars = []string{
	"GOOGLE_SERVICE_ACCOUNT",
}

var orgTargetEnvVars = []string{
	"GOOGLE_ORG_2",
}

var billingAccountEnvVars = []string{
	"GOOGLE_BILLING_ACCOUNT",
}

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccRandomProvider = random.Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"google": testAccProvider,
		"random": testAccRandomProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
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
// project to manage it until we can enable Firestore programatically.
func getTestFirestoreProjectFromEnv(t *testing.T) string {
	skipIfEnvNotSet(t, firestoreProjectEnvVars...)
	return multiEnvSearch(firestoreProjectEnvVars)
}

func getTestOrgFromEnv(t *testing.T) string {
	skipIfEnvNotSet(t, orgEnvVars...)
	return multiEnvSearch(orgEnvVars)
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
