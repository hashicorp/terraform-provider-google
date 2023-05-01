package google

import (
	"context"
	"fmt"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"io/ioutil"
	"log"
	"os"
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

const TestEnvVar = "TF_ACC"

var TestAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

// providerConfigEnvNames returns a list of all the environment variables that could be set by a user to configure the provider
func providerConfigEnvNames() []string {

	envs := []string{}

	// Use existing collections of ENV names
	envVarsSets := [][]string{
		CredsEnvVars,   // credentials field
		ProjectEnvVars, // project field
		regionEnvVars,  //region field
		zoneEnvVars,    // zone field
	}
	for _, set := range envVarsSets {
		envs = append(envs, set...)
	}

	// Add remaining ENVs
	envs = append(envs, "GOOGLE_OAUTH_ACCESS_TOKEN")          // access_token field
	envs = append(envs, "GOOGLE_BILLING_PROJECT")             // billing_project field
	envs = append(envs, "GOOGLE_IMPERSONATE_SERVICE_ACCOUNT") // impersonate_service_account field
	envs = append(envs, "USER_PROJECT_OVERRIDE")              // user_project_override field
	envs = append(envs, "CLOUDSDK_CORE_REQUEST_REASON")       // request_reason field

	return envs
}

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

// This value is the description used for test PublicAdvertisedPrefix setup to avoid required DNS
// setup. This is only used during integration tests and would be invalid to surface to users
var papDescriptionEnvVars = []string{
	"GOOGLE_PUBLIC_AVERTISED_PREFIX_DESCRIPTION",
}

func init() {
	configs = make(map[string]*transport_tpg.Config)
	fwProviders = make(map[string]*frameworkTestProvider)
	sources = make(map[string]VcrSource)
	testAccProvider = Provider()
	TestAccProviders = map[string]*schema.Provider{
		"google": testAccProvider,
	}
}

func GoogleProviderConfig(t *testing.T) *transport_tpg.Config {
	configsLock.RLock()
	config, ok := configs[t.Name()]
	configsLock.RUnlock()
	if ok {
		return config
	}

	sdkProvider := Provider()
	rc := terraform.ResourceConfig{}
	sdkProvider.Configure(context.Background(), &rc)
	return sdkProvider.Meta().(*transport_tpg.Config)
}

func AccTestPreCheck(t *testing.T) {
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

// GetTestRegion has the same logic as the provider's getRegion, to be used in tests.
func GetTestRegion(is *terraform.InstanceState, config *transport_tpg.Config) (string, error) {
	if res, ok := is.Attributes["region"]; ok {
		return res, nil
	}
	if config.Region != "" {
		return config.Region, nil
	}
	return "", fmt.Errorf("%q: required field is not set", "region")
}

// GetTestProject has the same logic as the provider's getProject, to be used in tests.
func GetTestProject(is *terraform.InstanceState, config *transport_tpg.Config) (string, error) {
	if res, ok := is.Attributes["project"]; ok {
		return res, nil
	}
	if config.Project != "" {
		return config.Project, nil
	}
	return "", fmt.Errorf("%q: required field is not set", "project")
}

// AccTestPreCheck ensures at least one of the project env variables is set.
func GetTestProjectNumberFromEnv() string {
	return MultiEnvSearch(projectNumberEnvVars)
}

// AccTestPreCheck ensures at least one of the project env variables is set.
func GetTestProjectFromEnv() string {
	return MultiEnvSearch(ProjectEnvVars)
}

// AccTestPreCheck ensures at least one of the credentials env variables is set.
func GetTestCredsFromEnv() string {
	// Return empty string if GOOGLE_USE_DEFAULT_CREDENTIALS is set to true.
	if MultiEnvSearch(CredsEnvVars) == "true" {
		return ""
	}
	return MultiEnvSearch(CredsEnvVars)
}

// AccTestPreCheck ensures at least one of the region env variables is set.
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

func GetTestPublicAdvertisedPrefixDescriptionFromEnv(t *testing.T) string {
	SkipIfEnvNotSet(t, papDescriptionEnvVars...)
	return MultiEnvSearch(papDescriptionEnvVars)
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

func setupProjectsAndGetAccessToken(org, billing, pid, service string, config *transport_tpg.Config) (string, error) {
	// Create project-1 and project-2
	rmService := config.NewResourceManagerClient(config.UserAgent)

	project := &cloudresourcemanager.Project{
		ProjectId: pid,
		Name:      pid,
		Parent: &cloudresourcemanager.ResourceId{
			Id:   org,
			Type: "organization",
		},
	}

	var op *cloudresourcemanager.Operation
	err := transport_tpg.RetryTimeDuration(func() (reqErr error) {
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
	project.Name = fmt.Sprintf("%s-2", pid)

	err = transport_tpg.RetryTimeDuration(func() (reqErr error) {
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
	curEmail, err := transport_tpg.GetCurrentUserEmail(config, config.UserAgent)
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
