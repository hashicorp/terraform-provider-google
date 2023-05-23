package google

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/cloudbilling/v1"
	cloudkms "google.golang.org/api/cloudkms/v1"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	iam "google.golang.org/api/iam/v1"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/serviceusage/v1"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

var SharedKeyRing = "tftest-shared-keyring-1"
var SharedCryptoKey = map[string]string{
	"ENCRYPT_DECRYPT":    "tftest-shared-key-1",
	"ASYMMETRIC_SIGN":    "tftest-shared-sign-key-1",
	"ASYMMETRIC_DECRYPT": "tftest-shared-decrypt-key-1",
}

type BootstrappedKMS struct {
	*cloudkms.KeyRing
	*cloudkms.CryptoKey
}

func BootstrapKMSKey(t *testing.T) BootstrappedKMS {
	return BootstrapKMSKeyInLocation(t, "global")
}

func BootstrapKMSKeyInLocation(t *testing.T, locationID string) BootstrappedKMS {
	return BootstrapKMSKeyWithPurposeInLocation(t, "ENCRYPT_DECRYPT", locationID)
}

// BootstrapKMSKeyWithPurpose returns a KMS key in the "global" location.
// See BootstrapKMSKeyWithPurposeInLocation.
func BootstrapKMSKeyWithPurpose(t *testing.T, purpose string) BootstrappedKMS {
	return BootstrapKMSKeyWithPurposeInLocation(t, purpose, "global")
}

/**
* BootstrapKMSKeyWithPurposeInLocation will return a KMS key in a
* particular location with the given purpose that can be used
* in tests that are testing KMS integration with other resources.
*
* This will either return an existing key or create one if it hasn't been created
* in the project yet. The motivation is because keyrings don't get deleted and we
* don't want a linear growth of disabled keyrings in a project. We also don't want
* to incur the overhead of creating a new project for each test that needs to use
* a KMS key.
**/
func BootstrapKMSKeyWithPurposeInLocation(t *testing.T, purpose, locationID string) BootstrappedKMS {
	return BootstrapKMSKeyWithPurposeInLocationAndName(t, purpose, locationID, SharedCryptoKey[purpose])
}

func BootstrapKMSKeyWithPurposeInLocationAndName(t *testing.T, purpose, locationID, keyShortName string) BootstrappedKMS {
	config := BootstrapConfig(t)
	if config == nil {
		return BootstrappedKMS{
			&cloudkms.KeyRing{},
			&cloudkms.CryptoKey{},
		}
	}

	projectID := acctest.GetTestProjectFromEnv()
	keyRingParent := fmt.Sprintf("projects/%s/locations/%s", projectID, locationID)
	keyRingName := fmt.Sprintf("%s/keyRings/%s", keyRingParent, SharedKeyRing)
	keyParent := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s", projectID, locationID, SharedKeyRing)
	keyName := fmt.Sprintf("%s/cryptoKeys/%s", keyParent, keyShortName)

	// Get or Create the hard coded shared keyring for testing
	kmsClient := config.NewKmsClient(config.UserAgent)
	keyRing, err := kmsClient.Projects.Locations.KeyRings.Get(keyRingName).Do()
	if err != nil {
		if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			keyRing, err = kmsClient.Projects.Locations.KeyRings.Create(keyRingParent, &cloudkms.KeyRing{}).
				KeyRingId(SharedKeyRing).Do()
			if err != nil {
				t.Errorf("Unable to bootstrap KMS key. Cannot create keyRing: %s", err)
			}
		} else {
			t.Errorf("Unable to bootstrap KMS key. Cannot retrieve keyRing: %s", err)
		}
	}

	if keyRing == nil {
		t.Fatalf("Unable to bootstrap KMS key. keyRing is nil!")
	}

	// Get or Create the hard coded, shared crypto key for testing
	cryptoKey, err := kmsClient.Projects.Locations.KeyRings.CryptoKeys.Get(keyName).Do()
	if err != nil {
		if transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
			algos := map[string]string{
				"ENCRYPT_DECRYPT":    "GOOGLE_SYMMETRIC_ENCRYPTION",
				"ASYMMETRIC_SIGN":    "RSA_SIGN_PKCS1_4096_SHA512",
				"ASYMMETRIC_DECRYPT": "RSA_DECRYPT_OAEP_4096_SHA512",
			}
			template := cloudkms.CryptoKeyVersionTemplate{
				Algorithm: algos[purpose],
			}

			newKey := cloudkms.CryptoKey{
				Purpose:         purpose,
				VersionTemplate: &template,
			}

			cryptoKey, err = kmsClient.Projects.Locations.KeyRings.CryptoKeys.Create(keyParent, &newKey).
				CryptoKeyId(keyShortName).Do()
			if err != nil {
				t.Errorf("Unable to bootstrap KMS key. Cannot create new CryptoKey: %s", err)
			}

		} else {
			t.Errorf("Unable to bootstrap KMS key. Cannot call CryptoKey service: %s", err)
		}
	}

	if cryptoKey == nil {
		t.Fatalf("Unable to bootstrap KMS key. CryptoKey is nil!")
	}

	return BootstrappedKMS{
		keyRing,
		cryptoKey,
	}
}

var serviceAccountEmail = "tf-bootstrap-service-account"
var serviceAccountDisplay = "Bootstrapped Service Account for Terraform tests"

// Some tests need a second service account, other than the test runner, to assert functionality on.
// This provides a well-known service account that can be used when dynamically creating a service
// account isn't an option.
func getOrCreateServiceAccount(config *transport_tpg.Config, project string) (*iam.ServiceAccount, error) {
	name := fmt.Sprintf("projects/%s/serviceAccounts/%s@%s.iam.gserviceaccount.com", project, serviceAccountEmail, project)
	log.Printf("[DEBUG] Verifying %s as bootstrapped service account.\n", name)

	sa, err := config.NewIamClient(config.UserAgent).Projects.ServiceAccounts.Get(name).Do()
	if err != nil && !transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
		return nil, err
	}

	if sa == nil {
		log.Printf("[DEBUG] Account missing. Creating %s as bootstrapped service account.\n", name)
		sa = &iam.ServiceAccount{
			DisplayName: serviceAccountDisplay,
		}

		r := &iam.CreateServiceAccountRequest{
			AccountId:      serviceAccountEmail,
			ServiceAccount: sa,
		}
		sa, err = config.NewIamClient(config.UserAgent).Projects.ServiceAccounts.Create("projects/"+project, r).Do()
		if err != nil {
			return nil, err
		}
	}

	return sa, nil
}

// In order to test impersonation we need to grant the testRunner's account the ability to grant tokens
// on a different service account. Granting permissions takes time and there is no operation to wait on
// so instead this creates a single service account once per test-suite with the correct permissions.
// The first time this test is run it will fail, but subsequent runs will succeed.
func impersonationServiceAccountPermissions(config *transport_tpg.Config, sa *iam.ServiceAccount, testRunner string) error {
	log.Printf("[DEBUG] Setting service account permissions.\n")
	policy := iam.Policy{
		Bindings: []*iam.Binding{},
	}

	binding := &iam.Binding{
		Role:    "roles/iam.serviceAccountTokenCreator",
		Members: []string{"serviceAccount:" + sa.Email, "serviceAccount:" + testRunner},
	}
	policy.Bindings = append(policy.Bindings, binding)

	// Overwrite the roles each time on this service account. This is because this account is
	// only created for the test suite and will stop snowflaking of permissions to get tests
	// to run. Overwriting permissions on 1 service account shouldn't affect others.
	_, err := config.NewIamClient(config.UserAgent).Projects.ServiceAccounts.SetIamPolicy(sa.Name, &iam.SetIamPolicyRequest{
		Policy: &policy,
	}).Do()
	if err != nil {
		return err
	}

	return nil
}

func BootstrapServiceAccount(t *testing.T, project, testRunner string) string {
	config := BootstrapConfig(t)
	if config == nil {
		return ""
	}

	sa, err := getOrCreateServiceAccount(config, project)
	if err != nil {
		t.Fatalf("Bootstrapping failed. Cannot retrieve service account, %s", err)
	}

	err = impersonationServiceAccountPermissions(config, sa, testRunner)
	if err != nil {
		t.Fatalf("Bootstrapping failed. Cannot set service account permissions, %s", err)
	}

	return sa.Email
}

const SharedTestADDomainPrefix = "tf-bootstrap-ad"

func BootstrapSharedTestADDomain(t *testing.T, testId string, networkName string) string {
	project := acctest.GetTestProjectFromEnv()
	sharedADDomain := fmt.Sprintf("%s.%s.com", SharedTestADDomainPrefix, testId)
	adDomainName := fmt.Sprintf("projects/%s/locations/global/domains/%s", project, sharedADDomain)

	config := BootstrapConfig(t)
	if config == nil {
		return ""
	}

	log.Printf("[DEBUG] Getting shared test active directory domain %q", adDomainName)
	getURL := fmt.Sprintf("%s%s", config.ActiveDirectoryBasePath, adDomainName)
	_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    getURL,
		UserAgent: config.UserAgent,
		Timeout:   4 * time.Minute,
	})
	if err != nil && transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
		log.Printf("[DEBUG] AD domain %q not found, bootstrapping", sharedADDomain)
		postURL := fmt.Sprintf("%sprojects/%s/locations/global/domains?domainName=%s", config.ActiveDirectoryBasePath, project, sharedADDomain)
		domainObj := map[string]interface{}{
			"locations":          []string{"us-central1"},
			"reservedIpRange":    "10.0.1.0/24",
			"authorizedNetworks": []string{fmt.Sprintf("projects/%s/global/networks/%s", project, networkName)},
		}

		_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   project,
			RawURL:    postURL,
			UserAgent: config.UserAgent,
			Body:      domainObj,
			Timeout:   60 * time.Minute,
		})
		if err != nil {
			t.Fatalf("Error bootstrapping shared active directory domain %q: %s", adDomainName, err)
		}

		log.Printf("[DEBUG] Waiting for active directory domain creation to finish")
	}

	_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    getURL,
		UserAgent: config.UserAgent,
		Timeout:   4 * time.Minute,
	})

	if err != nil {
		t.Fatalf("Error getting shared active directory domain %q: %s", adDomainName, err)
	}

	return sharedADDomain
}

const SharedTestNetworkPrefix = "tf-bootstrap-net-"

// BootstrapSharedTestNetwork will return a persistent compute network for a
// test or set of tests.
//
// Resources like service_networking_connection use a consumer network and
// create a complementing tenant network which we don't control. These tenant
// networks never get cleaned up and they can accumulate to the point where a
// limit is reached for the organization. By reusing a consumer network across
// test runs, we can reduce the number of tenant networks that are needed.
// See b/146351146 for more context.
//
// testId specifies the test for which a shared network is used/initialized.
// Note that if the network is being used for a service_networking_connection,
// the same testId should generally not be used across tests, to avoid race
// conditions where multiple tests attempt to modify the connection at once.
//
// Returns the name of a network, creating it if it hasn't been created in the
// test project.
func BootstrapSharedTestNetwork(t *testing.T, testId string) string {
	project := acctest.GetTestProjectFromEnv()
	networkName := SharedTestNetworkPrefix + testId

	config := BootstrapConfig(t)
	if config == nil {
		return ""
	}

	log.Printf("[DEBUG] Getting shared test network %q", networkName)
	_, err := config.NewComputeClient(config.UserAgent).Networks.Get(project, networkName).Do()
	if err != nil && transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
		log.Printf("[DEBUG] Network %q not found, bootstrapping", networkName)
		url := fmt.Sprintf("%sprojects/%s/global/networks", config.ComputeBasePath, project)
		netObj := map[string]interface{}{
			"name":                  networkName,
			"autoCreateSubnetworks": false,
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   project,
			RawURL:    url,
			UserAgent: config.UserAgent,
			Body:      netObj,
			Timeout:   4 * time.Minute,
		})
		if err != nil {
			t.Fatalf("Error bootstrapping shared test network %q: %s", networkName, err)
		}

		log.Printf("[DEBUG] Waiting for network creation to finish")
		err = ComputeOperationWaitTime(config, res, project, "Error bootstrapping shared test network", config.UserAgent, 4*time.Minute)
		if err != nil {
			t.Fatalf("Error bootstrapping shared test network %q: %s", networkName, err)
		}
	}

	network, err := config.NewComputeClient(config.UserAgent).Networks.Get(project, networkName).Do()
	if err != nil {
		t.Errorf("Error getting shared test network %q: %s", networkName, err)
	}
	if network == nil {
		t.Fatalf("Error getting shared test network %q: is nil", networkName)
	}
	return network.Name
}

var SharedServicePerimeterProjectPrefix = "tf-bootstrap-sp-"

func BootstrapServicePerimeterProjects(t *testing.T, desiredProjects int) []*cloudresourcemanager.Project {
	config := BootstrapConfig(t)
	if config == nil {
		return nil
	}

	org := acctest.GetTestOrgFromEnv(t)

	// The filter endpoint works differently if you provide both the parent id and parent type, and
	// doesn't seem to allow for prefix matching. Don't change this to include the parent type unless
	// that API behavior changes.
	prefixFilter := fmt.Sprintf("id:%s* parent.id:%s", SharedServicePerimeterProjectPrefix, org)
	res, err := config.NewResourceManagerClient(config.UserAgent).Projects.List().Filter(prefixFilter).Do()
	if err != nil {
		t.Fatalf("Error getting shared test projects: %s", err)
	}

	projects := res.Projects
	for len(projects) < desiredProjects {
		pid := SharedServicePerimeterProjectPrefix + RandString(t, 10)
		project := &cloudresourcemanager.Project{
			ProjectId: pid,
			Name:      "TF Service Perimeter Test",
			Parent: &cloudresourcemanager.ResourceId{
				Type: "organization",
				Id:   org,
			},
		}
		op, err := config.NewResourceManagerClient(config.UserAgent).Projects.Create(project).Do()
		if err != nil {
			t.Fatalf("Error bootstrapping shared test project: %s", err)
		}

		opAsMap, err := tpgresource.ConvertToMap(op)
		if err != nil {
			t.Fatalf("Error bootstrapping shared test project: %s", err)
		}

		err = ResourceManagerOperationWaitTime(config, opAsMap, "creating project", config.UserAgent, 4)
		if err != nil {
			t.Fatalf("Error bootstrapping shared test project: %s", err)
		}

		p, err := config.NewResourceManagerClient(config.UserAgent).Projects.Get(pid).Do()
		if err != nil {
			t.Fatalf("Error getting shared test project: %s", err)
		}
		projects = append(projects, p)
	}

	return projects
}

func RemoveContainerServiceAgentRoleFromContainerEngineRobot(t *testing.T, project *cloudresourcemanager.Project) {
	config := BootstrapConfig(t)
	if config == nil {
		return
	}

	client := config.NewResourceManagerClient(config.UserAgent)
	containerEngineRobot := fmt.Sprintf("serviceAccount:service-%d@container-engine-robot.iam.gserviceaccount.com", project.ProjectNumber)
	getPolicyRequest := &cloudresourcemanager.GetIamPolicyRequest{}
	policy, err := client.Projects.GetIamPolicy(project.ProjectId, getPolicyRequest).Do()
	if err != nil {
		t.Fatalf("error getting project iam policy: %v", err)
	}
	roleFound := false
	changed := false
	for _, binding := range policy.Bindings {
		if binding.Role == "roles/container.serviceAgent" {
			memberFound := false
			for i, member := range binding.Members {
				if member == containerEngineRobot {
					binding.Members[i] = binding.Members[len(binding.Members)-1]
					memberFound = true
				}
			}
			if memberFound {
				binding.Members = binding.Members[:len(binding.Members)-1]
				changed = true
			}
		} else if binding.Role == "roles/editor" {
			memberFound := false
			for _, member := range binding.Members {
				if member == containerEngineRobot {
					memberFound = true
					break
				}
			}
			if !memberFound {
				binding.Members = append(binding.Members, containerEngineRobot)
				changed = true
			}
			roleFound = true
		}
	}
	if !roleFound {
		policy.Bindings = append(policy.Bindings, &cloudresourcemanager.Binding{
			Members: []string{containerEngineRobot},
			Role:    "roles/editor",
		})
		changed = true
	}
	if changed {
		setPolicyRequest := &cloudresourcemanager.SetIamPolicyRequest{Policy: policy}
		policy, err = client.Projects.SetIamPolicy(project.ProjectId, setPolicyRequest).Do()
		if err != nil {
			t.Fatalf("error setting project iam policy: %v", err)
		}
	}
}

// BootstrapProject will create or get a project named
// "<projectIDPrefix><projectIDSuffix>" that will persist across test runs,
// where projectIDSuffix is based off of getTestProjectFromEnv(). The reason
// for the naming is to isolate bootstrapped projects by test environment.
// Given the existing projects being used by our team, the prefix provided to
// this function can be no longer than 18 characters.
func BootstrapProject(t *testing.T, projectIDPrefix, billingAccount string, services []string) *cloudresourcemanager.Project {
	config := BootstrapConfig(t)
	if config == nil {
		return nil
	}

	projectIDSuffix := strings.Replace(acctest.GetTestProjectFromEnv(), "ci-test-project-", "", 1)
	projectID := projectIDPrefix + projectIDSuffix

	crmClient := config.NewResourceManagerClient(config.UserAgent)

	project, err := crmClient.Projects.Get(projectID).Do()
	if err != nil {
		if !transport_tpg.IsGoogleApiErrorWithCode(err, 403) {
			t.Fatalf("Error getting bootstrapped project: %s", err)
		}
		org := acctest.GetTestOrgFromEnv(t)

		op, err := crmClient.Projects.Create(&cloudresourcemanager.Project{
			ProjectId: projectID,
			Name:      "Bootstrapped Test Project",
			Parent: &cloudresourcemanager.ResourceId{
				Type: "organization",
				Id:   org,
			},
		}).Do()
		if err != nil {
			t.Fatalf("Error creating bootstrapped test project: %s", err)
		}

		opAsMap, err := tpgresource.ConvertToMap(op)
		if err != nil {
			t.Fatalf("Error converting create project operation to map: %s", err)
		}

		err = ResourceManagerOperationWaitTime(config, opAsMap, "creating project", config.UserAgent, 4*time.Minute)
		if err != nil {
			t.Fatalf("Error waiting for create project operation: %s", err)
		}

		project, err = crmClient.Projects.Get(projectID).Do()
		if err != nil {
			t.Fatalf("Error getting bootstrapped project: %s", err)
		}

	}

	if project.LifecycleState == "DELETE_REQUESTED" {
		_, err := crmClient.Projects.Undelete(projectID, &cloudresourcemanager.UndeleteProjectRequest{}).Do()
		if err != nil {
			t.Fatalf("Error undeleting bootstrapped project: %s", err)
		}
	}

	if billingAccount != "" {
		billingClient := config.NewBillingClient(config.UserAgent)
		var pbi *cloudbilling.ProjectBillingInfo
		err = transport_tpg.RetryTimeDuration(func() error {
			var reqErr error
			pbi, reqErr = billingClient.Projects.GetBillingInfo(PrefixedProject(projectID)).Do()
			return reqErr
		}, 30*time.Second)
		if err != nil {
			t.Fatalf("Error getting billing info for project %q: %v", projectID, err)
		}
		if strings.TrimPrefix(pbi.BillingAccountName, "billingAccounts/") != billingAccount {
			pbi.BillingAccountName = "billingAccounts/" + billingAccount
			err := transport_tpg.RetryTimeDuration(func() error {
				_, err := config.NewBillingClient(config.UserAgent).Projects.UpdateBillingInfo(PrefixedProject(projectID), pbi).Do()
				return err
			}, 2*time.Minute)
			if err != nil {
				t.Fatalf("Error setting billing account for project %q to %q: %s", projectID, billingAccount, err)
			}
		}
	}

	if len(services) > 0 {

		enabledServices, err := ListCurrentlyEnabledServices(projectID, "", config.UserAgent, config, 1*time.Minute)
		if err != nil {
			t.Fatalf("Error listing services for project %q: %s", projectID, err)
		}

		servicesToEnable := make([]string, 0, len(services))
		for _, service := range services {
			if _, ok := enabledServices[service]; !ok {
				servicesToEnable = append(servicesToEnable, service)
			}
		}

		if len(servicesToEnable) > 0 {
			if err := EnableServiceUsageProjectServices(servicesToEnable, projectID, "", config.UserAgent, config, 10*time.Minute); err != nil {
				t.Fatalf("Error enabling services for project %q: %s", projectID, err)
			}
		}
	}

	return project
}

// BootstrapConfig returns a Config pulled from the environment.
func BootstrapConfig(t *testing.T) *transport_tpg.Config {
	if v := os.Getenv("TF_ACC"); v == "" {
		t.Skip("Acceptance tests and bootstrapping skipped unless env 'TF_ACC' set")
		return nil
	}

	config := &transport_tpg.Config{
		Credentials: acctest.GetTestCredsFromEnv(),
		Project:     acctest.GetTestProjectFromEnv(),
		Region:      acctest.GetTestRegionFromEnv(),
		Zone:        acctest.GetTestZoneFromEnv(),
	}

	transport_tpg.ConfigureBasePaths(config)

	if err := config.LoadAndValidate(context.Background()); err != nil {
		t.Fatalf("Bootstrapping failed. Unable to load test config: %s", err)
	}
	return config
}

// SQL Instance names are not reusable for a week after deletion
const SharedTestSQLInstanceNamePrefix = "tf-bootstrap-"

// BootstrapSharedSQLInstanceBackupRun will return a shared SQL db instance that
// has a backup created for it.
func BootstrapSharedSQLInstanceBackupRun(t *testing.T) string {
	project := acctest.GetTestProjectFromEnv()

	config := BootstrapConfig(t)
	if config == nil {
		return ""
	}

	log.Printf("[DEBUG] Getting list of existing sql instances")

	instances, err := config.NewSqlAdminClient(config.UserAgent).Instances.List(project).Do()
	if err != nil {
		t.Fatalf("Unable to bootstrap SQL Instance. Cannot retrieve instance list: %s", err)
	}

	var bootstrapInstance *sqladmin.DatabaseInstance

	// Look for any existing bootstrap instances
	for _, i := range instances.Items {
		if strings.HasPrefix(i.Name, SharedTestSQLInstanceNamePrefix) {
			bootstrapInstance = i
			break
		}
	}

	if bootstrapInstance == nil {
		bootstrapInstanceName := SharedTestSQLInstanceNamePrefix + RandString(t, 10)
		log.Printf("[DEBUG] Bootstrap SQL Instance not found, bootstrapping new instance %s", bootstrapInstanceName)

		backupConfig := &sqladmin.BackupConfiguration{
			Enabled:                    true,
			PointInTimeRecoveryEnabled: true,
		}
		settings := &sqladmin.Settings{
			Tier:                "db-f1-micro",
			BackupConfiguration: backupConfig,
		}
		bootstrapInstance = &sqladmin.DatabaseInstance{
			Name:            bootstrapInstanceName,
			Region:          "us-central1",
			Settings:        settings,
			DatabaseVersion: "POSTGRES_11",
		}

		var op *sqladmin.Operation
		err = transport_tpg.RetryTimeDuration(func() (operr error) {
			op, operr = config.NewSqlAdminClient(config.UserAgent).Instances.Insert(project, bootstrapInstance).Do()
			return operr
		}, time.Duration(20)*time.Minute, transport_tpg.IsSqlOperationInProgressError)
		if err != nil {
			t.Fatalf("Error, failed to create instance %s: %s", bootstrapInstance.Name, err)
		}
		err = SqlAdminOperationWaitTime(config, op, project, "Create Instance", config.UserAgent, time.Duration(40)*time.Minute)
		if err != nil {
			t.Fatalf("Error, failed to create instance %s: %s", bootstrapInstance.Name, err)
		}
	}

	// Look for backups in bootstrap instance
	res, err := config.NewSqlAdminClient(config.UserAgent).BackupRuns.List(project, bootstrapInstance.Name).Do()
	if err != nil {
		t.Fatalf("Unable to bootstrap SQL Instance. Cannot retrieve backup list: %s", err)
	}
	backupsList := res.Items
	if len(backupsList) == 0 {
		log.Printf("[DEBUG] No backups found for %s, creating backup", bootstrapInstance.Name)
		backupRun := &sqladmin.BackupRun{
			Instance: bootstrapInstance.Name,
		}

		var op *sqladmin.Operation
		err = transport_tpg.RetryTimeDuration(func() (operr error) {
			op, operr = config.NewSqlAdminClient(config.UserAgent).BackupRuns.Insert(project, bootstrapInstance.Name, backupRun).Do()
			return operr
		}, time.Duration(20)*time.Minute, transport_tpg.IsSqlOperationInProgressError)
		if err != nil {
			t.Fatalf("Error, failed to create instance backup: %s", err)
		}
		err = SqlAdminOperationWaitTime(config, op, project, "Backup Instance", config.UserAgent, time.Duration(20)*time.Minute)
		if err != nil {
			t.Fatalf("Error, failed to create instance backup: %s", err)
		}
	}

	return bootstrapInstance.Name
}

func BootstrapSharedCaPoolInLocation(t *testing.T, location string) string {
	project := acctest.GetTestProjectFromEnv()
	poolName := "static-ca-pool"

	config := BootstrapConfig(t)
	if config == nil {
		return ""
	}

	log.Printf("[DEBUG] Getting shared CA pool %q", poolName)
	url := fmt.Sprintf("%sprojects/%s/locations/%s/caPools/%s", config.PrivatecaBasePath, project, location, poolName)
	_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: config.UserAgent,
	})
	if err != nil {
		log.Printf("[DEBUG] CA pool %q not found, bootstrapping", poolName)
		poolObj := map[string]interface{}{
			"tier": "ENTERPRISE",
		}
		createUrl := fmt.Sprintf("%sprojects/%s/locations/%s/caPools?caPoolId=%s", config.PrivatecaBasePath, project, location, poolName)
		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   project,
			RawURL:    createUrl,
			UserAgent: config.UserAgent,
			Body:      poolObj,
			Timeout:   4 * time.Minute,
		})
		if err != nil {
			t.Fatalf("Error bootstrapping shared CA pool %q: %s", poolName, err)
		}

		log.Printf("[DEBUG] Waiting for CA pool creation to finish")
		var opRes map[string]interface{}
		err = PrivatecaOperationWaitTimeWithResponse(
			config, res, &opRes, project, "Creating CA pool", config.UserAgent,
			4*time.Minute)
		if err != nil {
			t.Errorf("Error getting shared CA pool %q: %s", poolName, err)
		}
		_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   project,
			RawURL:    url,
			UserAgent: config.UserAgent,
		})
		if err != nil {
			t.Errorf("Error getting shared CA pool %q: %s", poolName, err)
		}
	}
	return poolName
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
	opAsMap, err := tpgresource.ConvertToMap(op)
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
	opAsMap, err = tpgresource.ConvertToMap(op)
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

	bindings := tpgiamresource.MergeBindings([]*cloudresourcemanager.Binding{proj1SATokenCreator, proj1SACreator})

	p, err := rmService.Projects.GetIamPolicy(pid,
		&cloudresourcemanager.GetIamPolicyRequest{
			Options: &cloudresourcemanager.GetPolicyOptions{
				RequestedPolicyVersion: tpgiamresource.IamPolicyVersion,
			},
		}).Do()
	if err != nil {
		return "", err
	}

	p.Bindings = tpgiamresource.MergeBindings(append(p.Bindings, bindings...))
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

	bindings = tpgiamresource.MergeBindings([]*cloudresourcemanager.Binding{proj2ServiceUsageBinding, proj2ServiceAdminBinding})

	// For KMS test only
	if service == "cloudkms" {
		proj2CryptoKeyBinding := &cloudresourcemanager.Binding{
			Members: []string{fmt.Sprintf("serviceAccount:%s", sa1.Email)},
			Role:    "roles/cloudkms.cryptoKeyEncrypter",
		}

		bindings = tpgiamresource.MergeBindings(append(bindings, proj2CryptoKeyBinding))
	}

	p, err = rmService.Projects.GetIamPolicy(p2,
		&cloudresourcemanager.GetIamPolicyRequest{
			Options: &cloudresourcemanager.GetPolicyOptions{
				RequestedPolicyVersion: tpgiamresource.IamPolicyVersion,
			},
		}).Do()
	if err != nil {
		return "", err
	}

	p.Bindings = tpgiamresource.MergeBindings(append(p.Bindings, bindings...))
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
