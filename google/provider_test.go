package google

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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

func TestAccProviderBasePath_setBasePath(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderBasePath_setBasePath("https://www.googleapis.com/compute/beta/", acctest.RandString(10)),
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

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccProviderBasePath_setBasePath("https://www.example.com/compute/beta/", acctest.RandString(10)),
				ExpectError: regexp.MustCompile("got HTTP response code 404 with body"),
			},
		},
	})
}

func TestAccProviderUserProjectOverride(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	billing := getTestBillingAccountFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	sa := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
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
				Config:      testAccProviderUserProjectOverride_step2(pid, pname, org, billing, sa, false),
				ExpectError: regexp.MustCompile("Binary Authorization API has not been used"),
			},
			{
				Config: testAccProviderUserProjectOverride_step2(pid, pname, org, billing, sa, true),
			},
			{
				ResourceName:      "google_binary_authorization_policy.project-2-policy",
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
	t.Parallel()

	org := getTestOrgFromEnv(t)
	billing := getTestBillingAccountFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	sa := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
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
// binauthz policy in project 2. The binauthz API is only enabled in project 2,
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

resource "google_project_service" "project-2-binauthz" {
	project = google_project.project-2.project_id
	service = "binaryauthorization.googleapis.com"
}

// Permission needed for user_project_override
resource "google_project_iam_member" "project-2-serviceusage" {
	project = google_project.project-2.project_id
	role    = "roles/serviceusage.serviceUsageConsumer"
	member  = "serviceAccount:${google_service_account.project-1.email}"
}

resource "google_project_iam_member" "project-2-binauthz" {
	project = google_project.project-2.project_id
	role    = "roles/binaryauthorization.policyEditor"
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

func testAccProviderUserProjectOverride_step2(pid, name, org, billing, sa string, override bool) string {
	return fmt.Sprintf(`
// See step 3 below, which is really step 2 minus the binauthz policy.
// Step 3 exists because provider configurations can't be removed while objects
// created by that provider still exist in state. Step 3 will remove the
// binauthz policy so the whole config can be deleted.
%s

resource "google_binary_authorization_policy" "project-2-policy" {
	provider = google.project-1-token
	project  = google_project.project-2.project_id

	admission_whitelist_patterns {
		name_pattern= "gcr.io/google_containers/*"
	}

	default_admission_rule {
		evaluation_mode = "ALWAYS_DENY"
		enforcement_mode = "ENFORCED_BLOCK_AND_AUDIT_LOG"
	}
}
`, testAccProviderUserProjectOverride_step3(pid, name, org, billing, sa, override))
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
