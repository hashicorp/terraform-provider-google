package google

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/serviceusage/v1"
)

var TestAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

var CredsEnvVars = []string{
	"GOOGLE_CREDENTIALS",
	"GOOGLE_CLOUD_KEYFILE_JSON",
	"GCLOUD_KEYFILE_JSON",
	"GOOGLE_APPLICATION_CREDENTIALS",
	"GOOGLE_USE_DEFAULT_CREDENTIALS",
}

var projectNumberEnvVars = []string{
	"GOOGLE_PROJECT_NUMBER",
}

var ProjectEnvVars = []string{
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

// This is the billing account that will be charged for the infrastructure used during testing. For
// that reason, it is also the billing account used for creating new projects.
var billingAccountEnvVars = []string{
	"GOOGLE_BILLING_ACCOUNT",
}

// This is the billing account that will be modified to test billing-related functionality. It is
// expected to have more permissions granted to the test user and support subaccounts.
var masterBillingAccountEnvVars = []string{
	"GOOGLE_MASTER_BILLING_ACCOUNT",
}

func init() {
	configs = make(map[string]*Config)
	fwProviders = make(map[string]*frameworkTestProvider)
	sources = make(map[string]VcrSource)
	testAccProvider = Provider()
	TestAccProviders = map[string]*schema.Provider{
		"google": testAccProvider,
	}
}

func GoogleProviderConfig(t *testing.T) *Config {
	configsLock.RLock()
	config, ok := configs[t.Name()]
	configsLock.RUnlock()
	if ok {
		return config
	}

	sdkProvider := Provider()
	rc := terraform.ResourceConfig{}
	sdkProvider.Configure(context.Background(), &rc)
	return sdkProvider.Meta().(*Config)
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

	if v := MultiEnvSearch(CredsEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(CredsEnvVars, ", "))
	}

	if v := MultiEnvSearch(ProjectEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(ProjectEnvVars, ", "))
	}

	if v := MultiEnvSearch(regionEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(regionEnvVars, ", "))
	}

	if v := MultiEnvSearch(zoneEnvVars); v == "" {
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

	VcrTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
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
		PreCheck:                 func() { testAccPreCheck(t) },
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
		PreCheck:                 func() { testAccPreCheck(t) },
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
		PreCheck:                 func() { testAccPreCheck(t) },
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
		PreCheck:                 func() { testAccPreCheck(t) },
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

// GetTestRegion has the same logic as the provider's getRegion, to be used in tests.
func GetTestRegion(is *terraform.InstanceState, config *Config) (string, error) {
	if res, ok := is.Attributes["region"]; ok {
		return res, nil
	}
	if config.Region != "" {
		return config.Region, nil
	}
	return "", fmt.Errorf("%q: required field is not set", "region")
}

// GetTestProject has the same logic as the provider's getProject, to be used in tests.
func GetTestProject(is *terraform.InstanceState, config *Config) (string, error) {
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
	return MultiEnvSearch(projectNumberEnvVars)
}

// testAccPreCheck ensures at least one of the project env variables is set.
func GetTestProjectFromEnv() string {
	return MultiEnvSearch(ProjectEnvVars)
}

// testAccPreCheck ensures at least one of the credentials env variables is set.
func GetTestCredsFromEnv() string {
	// Return empty string if GOOGLE_USE_DEFAULT_CREDENTIALS is set to true.
	if MultiEnvSearch(CredsEnvVars) == "true" {
		return ""
	}
	return MultiEnvSearch(CredsEnvVars)
}

// testAccPreCheck ensures at least one of the region env variables is set.
func GetTestRegionFromEnv() string {
	return MultiEnvSearch(regionEnvVars)
}

func GetTestZoneFromEnv() string {
	return MultiEnvSearch(zoneEnvVars)
}

func GetTestCustIdFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, custIdEnvVars...)
	return MultiEnvSearch(custIdEnvVars)
}

func GetTestIdentityUserFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, identityUserEnvVars...)
	return MultiEnvSearch(identityUserEnvVars)
}

// Firestore can't be enabled at the same time as Datastore, so we need a new
// project to manage it until we can enable Firestore programmatically.
func GetTestFirestoreProjectFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, firestoreProjectEnvVars...)
	return MultiEnvSearch(firestoreProjectEnvVars)
}

// Returns the raw organization id like 1234567890, skipping the test if one is
// not found.
func GetTestOrgFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, orgEnvVars...)
	return MultiEnvSearch(orgEnvVars)
}

// Alternative to GetTestOrgFromEnv that doesn't need *testing.T
// If using this, you need to process unset values at the call site
func UnsafeGetTestOrgFromEnv() string {
	return MultiEnvSearch(orgEnvVars)
}

func GetTestOrgDomainFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, orgEnvDomainVars...)
	return MultiEnvSearch(orgEnvDomainVars)
}

func GetTestOrgTargetFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, orgTargetEnvVars...)
	return MultiEnvSearch(orgTargetEnvVars)
}

// This is the billing account that will be charged for the infrastructure used during testing. For
// that reason, it is also the billing account used for creating new projects.
func GetTestBillingAccountFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, billingAccountEnvVars...)
	return MultiEnvSearch(billingAccountEnvVars)
}

// This is the billing account that will be modified to test billing-related functionality. It is
// expected to have more permissions granted to the test user and support subaccounts.
func GetTestMasterBillingAccountFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, masterBillingAccountEnvVars...)
	return MultiEnvSearch(masterBillingAccountEnvVars)
}

func GetTestServiceAccountFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, serviceAccountEnvVars...)
	return MultiEnvSearch(serviceAccountEnvVars)
}

// Some tests fail during VCR. One common case is race conditions when creating resources.
// If a test config adds two fine-grained resources with the same parent it is undefined
// which will be created first, causing VCR to fail ~50% of the time
func SkipIfVcr(t *testing.T) {
	if isVcrEnabled() {
		t.Skipf("VCR enabled, skipping test: %s", t.Name())
	}
}

func SleepInSecondsForTest(t int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		time.Sleep(time.Duration(t) * time.Second)
		return nil
	}
}

func setupProjectsAndGetAccessToken(org, billing, pid, service string, config *Config) (string, error) {
	// Create project-1 and project-2
	rmService := config.NewResourceManagerClient(config.UserAgent)

	project := &cloudresourcemanager.Project{
		ProjectId: pid,
		Name:      pname,
		Parent: &cloudresourcemanager.ResourceId{
			Id:   org,
			Type: "organization",
		},
	}

	var op *cloudresourcemanager.Operation
	err := RetryTimeDuration(func() (reqErr error) {
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

	waitErr := ResourceManagerOperationWaitTime(config, opAsMap, "creating project", config.UserAgent, 5*time.Minute)
	if waitErr != nil {
		return "", waitErr
	}

	ba := &cloudbilling.ProjectBillingInfo{
		BillingAccountName: fmt.Sprintf("billingAccounts/%s", billing),
	}
	_, err = config.NewBillingClient(config.UserAgent).Projects.UpdateBillingInfo(PrefixedProject(pid), ba).Do()
	if err != nil {
		return "", err
	}

	p2 := fmt.Sprintf("%s-2", pid)
	project.ProjectId = p2
	project.Name = fmt.Sprintf("%s-2", pname)

	err = RetryTimeDuration(func() (reqErr error) {
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

	waitErr = ResourceManagerOperationWaitTime(config, opAsMap, "creating project", config.UserAgent, 5*time.Minute)
	if waitErr != nil {
		return "", waitErr
	}

	_, err = config.NewBillingClient(config.UserAgent).Projects.UpdateBillingInfo(PrefixedProject(p2), ba).Do()
	if err != nil {
		return "", err
	}

	// Enable the appropriate service in project-2 only
	suService := config.NewServiceUsageClient(config.UserAgent)

	serviceReq := &serviceusage.BatchEnableServicesRequest{
		ServiceIds: []string{fmt.Sprintf("%s.googleapis.com", service)},
	}

	_, err = suService.Services.BatchEnable(fmt.Sprintf("projects/%s", p2), serviceReq).Do()
	if err != nil {
		return "", err
	}

	// Enable the test runner to create service accounts and get an access token on behalf of
	// the project 1 service account
	curEmail, err := GetCurrentUserEmail(config, config.UserAgent)
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

	bindings := MergeBindings([]*cloudresourcemanager.Binding{proj1SATokenCreator, proj1SACreator})

	p, err := rmService.Projects.GetIamPolicy(pid,
		&cloudresourcemanager.GetIamPolicyRequest{
			Options: &cloudresourcemanager.GetPolicyOptions{
				RequestedPolicyVersion: IamPolicyVersion,
			},
		}).Do()
	if err != nil {
		return "", err
	}

	p.Bindings = MergeBindings(append(p.Bindings, bindings...))
	_, err = config.NewResourceManagerClient(config.UserAgent).Projects.SetIamPolicy(pid,
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

	bindings = MergeBindings([]*cloudresourcemanager.Binding{proj2ServiceUsageBinding, proj2ServiceAdminBinding})

	// For KMS test only
	if service == "cloudkms" {
		proj2CryptoKeyBinding := &cloudresourcemanager.Binding{
			Members: []string{fmt.Sprintf("serviceAccount:%s", sa1.Email)},
			Role:    "roles/cloudkms.cryptoKeyEncrypter",
		}

		bindings = MergeBindings(append(bindings, proj2CryptoKeyBinding))
	}

	p, err = rmService.Projects.GetIamPolicy(p2,
		&cloudresourcemanager.GetIamPolicyRequest{
			Options: &cloudresourcemanager.GetPolicyOptions{
				RequestedPolicyVersion: IamPolicyVersion,
			},
		}).Do()
	if err != nil {
		return "", err
	}

	p.Bindings = MergeBindings(append(p.Bindings, bindings...))
	_, err = config.NewResourceManagerClient(config.UserAgent).Projects.SetIamPolicy(p2,
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

	iamCredsService := config.NewIamCredentialsClient(config.UserAgent)
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

func SkipIfEnvNotSet(t *testing.T, envs ...string) {
	if t == nil {
		log.Printf("[DEBUG] Not running inside of test - skip skipping")
		return
	}

	for _, k := range envs {
		if os.Getenv(k) == "" {
			log.Printf("[DEBUG] Warning - environment variable %s is not set - skipping test %s", k, t.Name())
			t.Skipf("Environment variable %s is not set", k)
		}
	}
}
