// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package acctest

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	tpgcompute "github.com/hashicorp/terraform-provider-google/google/services/compute"
	"github.com/hashicorp/terraform-provider-google/google/services/privateca"
	"github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	tpgservicenetworking "github.com/hashicorp/terraform-provider-google/google/services/servicenetworking"
	"github.com/hashicorp/terraform-provider-google/google/services/sql"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/cloudbilling/v1"
	cloudkms "google.golang.org/api/cloudkms/v1"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	iam "google.golang.org/api/iam/v1"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/servicenetworking/v1"
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

	projectID := envvar.GetTestProjectFromEnv()
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

var serviceAccountPrefix = "tf-bootstrap-sa-"
var serviceAccountDisplay = "Bootstrapped Service Account for Terraform tests"

// Some tests need a second service account, other than the test runner, to assert functionality on.
// This provides a well-known service account that can be used when dynamically creating a service
// account isn't an option.
func getOrCreateServiceAccount(config *transport_tpg.Config, project, serviceAccountEmail string) (*iam.ServiceAccount, error) {
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

// A separate testId should be used for each test, to create separate service accounts for each,
// and avoid race conditions where the policy of the same service account is being modified by 2
// tests at once. This is needed as long as the function overwrites the policy on every run.
func BootstrapServiceAccount(t *testing.T, testId, testRunner string) string {
	project := envvar.GetTestProjectFromEnv()
	serviceAccountEmail := serviceAccountPrefix + testId

	config := BootstrapConfig(t)
	if config == nil {
		return ""
	}

	sa, err := getOrCreateServiceAccount(config, project, serviceAccountEmail)
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
	project := envvar.GetTestProjectFromEnv()
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
// Usage 1
// Resources like service_networking_connection use a consumer network and
// create a complementing tenant network which we don't control. These tenant
// networks never get cleaned up and they can accumulate to the point where a
// limit is reached for the organization. By reusing a consumer network across
// test runs, we can reduce the number of tenant networks that are needed.
// See b/146351146 for more context.
//
// Usage 2
// Bootstrap networks used in tests (gke clusters, dataproc clusters...)
// to avoid traffic to the default network
//
// testId specifies the test for which a shared network is used/initialized.
// Note that if the network is being used for a service_networking_connection,
// the same testId should generally not be used across tests, to avoid race
// conditions where multiple tests attempt to modify the connection at once.
//
// Returns the name of a network, creating it if it hasn't been created in the
// test project.
func BootstrapSharedTestNetwork(t *testing.T, testId string) string {
	project := envvar.GetTestProjectFromEnv()
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
		err = tpgcompute.ComputeOperationWaitTime(config, res, project, "Error bootstrapping shared test network", config.UserAgent, 4*time.Minute)
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

type AddressSettings struct {
	PrefixLength int
}

func AddressWithPrefixLength(prefixLength int) func(*AddressSettings) {
	return func(settings *AddressSettings) {
		settings.PrefixLength = prefixLength
	}
}

func NewAddressSettings(options ...func(*AddressSettings)) *AddressSettings {
	settings := &AddressSettings{
		PrefixLength: 16, // default prefix length
	}

	for _, o := range options {
		o(settings)
	}
	return settings
}

const SharedTestGlobalAddressPrefix = "tf-bootstrap-addr-"

// params are the functions to set compute global address
func BootstrapSharedTestGlobalAddress(t *testing.T, testId string, params ...func(*AddressSettings)) string {
	project := envvar.GetTestProjectFromEnv()
	addressName := SharedTestGlobalAddressPrefix + testId
	networkName := BootstrapSharedTestNetwork(t, testId)
	networkId := fmt.Sprintf("projects/%v/global/networks/%v", project, networkName)

	config := BootstrapConfig(t)
	if config == nil {
		return ""
	}

	log.Printf("[DEBUG] Getting shared test global address %q", addressName)
	_, err := config.NewComputeClient(config.UserAgent).GlobalAddresses.Get(project, addressName).Do()
	if err != nil && transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
		log.Printf("[DEBUG] Global address %q not found, bootstrapping", addressName)
		url := fmt.Sprintf("%sprojects/%s/global/addresses", config.ComputeBasePath, project)

		settings := NewAddressSettings(params...)

		netObj := map[string]interface{}{
			"name":          addressName,
			"address_type":  "INTERNAL",
			"purpose":       "VPC_PEERING",
			"prefix_length": settings.PrefixLength,
			"network":       networkId,
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
			t.Fatalf("Error bootstrapping shared test global address %q: %s", addressName, err)
		}

		log.Printf("[DEBUG] Waiting for global address creation to finish")
		err = tpgcompute.ComputeOperationWaitTime(config, res, project, "Error bootstrapping shared test global address", config.UserAgent, 4*time.Minute)
		if err != nil {
			t.Fatalf("Error bootstrapping shared test global address %q: %s", addressName, err)
		}
	}

	address, err := config.NewComputeClient(config.UserAgent).GlobalAddresses.Get(project, addressName).Do()
	if err != nil {
		t.Errorf("Error getting shared test global address %q: %s", addressName, err)
	}
	if address == nil {
		t.Fatalf("Error getting shared test global address %q: is nil", addressName)
	}
	return address.Name
}

type ServiceNetworkSettings struct {
	PrefixLength  int
	ParentService string
}

func ServiceNetworkWithPrefixLength(prefixLength int) func(*ServiceNetworkSettings) {
	return func(settings *ServiceNetworkSettings) {
		settings.PrefixLength = prefixLength
	}
}

func ServiceNetworkWithParentService(parentService string) func(*ServiceNetworkSettings) {
	return func(settings *ServiceNetworkSettings) {
		settings.ParentService = parentService
	}
}

func NewServiceNetworkSettings(options ...func(*ServiceNetworkSettings)) *ServiceNetworkSettings {
	settings := &ServiceNetworkSettings{
		PrefixLength:  16,                                 // default prefix length
		ParentService: "servicenetworking.googleapis.com", // default parent service
	}

	for _, o := range options {
		o(settings)
	}
	return settings
}

// BootstrapSharedServiceNetworkingConnection will create a shared network
// if it hasn't been created in the test project, a global address
// if it hasn't been created in the test project, and a service networking connection
// if it hasn't been created in the test project.
//
// params are the functions to set compute global address
//
// BootstrapSharedServiceNetworkingConnection returns a persistent compute network name
// for a test or set of tests.
//
// To delete a service networking conneciton, all of the service instances that use that connection
// must be deleted first. After the service instances are deleted, some service producers delay the deletion
// utnil a waiting period has passed. For example, after four days that you delete a SQL instance,
// the service networking connection can be deleted.
// That is the reason to use the shared service networking connection for thest resources.
// https://cloud.google.com/vpc/docs/configure-private-services-access#removing-connection
//
// testId specifies the test for which a shared network and a gobal address are used/initialized.
func BootstrapSharedServiceNetworkingConnection(t *testing.T, testId string, params ...func(*ServiceNetworkSettings)) string {
	settings := NewServiceNetworkSettings(params...)
	parentService := "services/" + settings.ParentService
	projectId := envvar.GetTestProjectFromEnv()

	config := BootstrapConfig(t)
	if config == nil {
		return ""
	}

	// Get project number by calling the API
	crmClient := config.NewResourceManagerClient(config.UserAgent)
	project, err := crmClient.Projects.Get(projectId).Do()
	if err != nil {
		t.Fatalf("Error getting project: %s", err)
	}

	networkName := SharedTestNetworkPrefix + testId
	networkId := fmt.Sprintf("projects/%v/global/networks/%v", project.ProjectNumber, networkName)
	globalAddressName := BootstrapSharedTestGlobalAddress(t, testId, AddressWithPrefixLength(settings.PrefixLength))

	readCall := config.NewServiceNetworkingClient(config.UserAgent).Services.Connections.List(parentService).Network(networkId)
	if config.UserProjectOverride {
		readCall.Header().Add("X-Goog-User-Project", projectId)
	}
	response, err := readCall.Do()
	if err != nil {
		t.Errorf("Error getting shared test service networking connection: %s", err)
	}

	var connection *servicenetworking.Connection
	for _, c := range response.Connections {
		if c.Network == networkId {
			connection = c
			break
		}
	}

	if connection == nil {
		log.Printf("[DEBUG] Service networking connection not found, bootstrapping")

		connection := &servicenetworking.Connection{
			Network:               networkId,
			ReservedPeeringRanges: []string{globalAddressName},
		}

		createCall := config.NewServiceNetworkingClient(config.UserAgent).Services.Connections.Create(parentService, connection)
		if config.UserProjectOverride {
			createCall.Header().Add("X-Goog-User-Project", projectId)
		}
		op, err := createCall.Do()
		if err != nil {
			t.Fatalf("Error bootstrapping shared test service networking connection: %s", err)
		}

		log.Printf("[DEBUG] Waiting for service networking connection creation to finish")
		if err := tpgservicenetworking.ServiceNetworkingOperationWaitTimeHW(config, op, "Create Service Networking Connection", config.UserAgent, projectId, 4*time.Minute); err != nil {
			t.Fatalf("Error bootstrapping shared test service networking connection: %s", err)
		}
	}

	log.Printf("[DEBUG] Getting shared test service networking connection")

	return networkName
}

var SharedServicePerimeterProjectPrefix = "tf-bootstrap-sp-"

func BootstrapServicePerimeterProjects(t *testing.T, desiredProjects int) []*cloudresourcemanager.Project {
	config := BootstrapConfig(t)
	if config == nil {
		return nil
	}

	org := envvar.GetTestOrgFromEnv(t)

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

		err = resourcemanager.ResourceManagerOperationWaitTime(config, opAsMap, "creating project", config.UserAgent, 4)
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

	projectIDSuffix := strings.Replace(envvar.GetTestProjectFromEnv(), "ci-test-project-", "", 1)
	projectID := projectIDPrefix + projectIDSuffix

	crmClient := config.NewResourceManagerClient(config.UserAgent)

	project, err := crmClient.Projects.Get(projectID).Do()
	if err != nil {
		if !transport_tpg.IsGoogleApiErrorWithCode(err, 403) {
			t.Fatalf("Error getting bootstrapped project: %s", err)
		}
		org := envvar.GetTestOrgFromEnv(t)

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

		err = resourcemanager.ResourceManagerOperationWaitTime(config, opAsMap, "creating project", config.UserAgent, 4*time.Minute)
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
		err = transport_tpg.Retry(transport_tpg.RetryOptions{
			RetryFunc: func() error {
				var reqErr error
				pbi, reqErr = billingClient.Projects.GetBillingInfo(resourcemanager.PrefixedProject(projectID)).Do()
				return reqErr
			},
			Timeout: 30 * time.Second,
		})
		if err != nil {
			t.Fatalf("Error getting billing info for project %q: %v", projectID, err)
		}
		if strings.TrimPrefix(pbi.BillingAccountName, "billingAccounts/") != billingAccount {
			pbi.BillingAccountName = "billingAccounts/" + billingAccount
			err := transport_tpg.Retry(transport_tpg.RetryOptions{
				RetryFunc: func() error {
					_, err := config.NewBillingClient(config.UserAgent).Projects.UpdateBillingInfo(resourcemanager.PrefixedProject(projectID), pbi).Do()
					return err
				},
				Timeout: 2 * time.Minute,
			})
			if err != nil {
				t.Fatalf("Error setting billing account for project %q to %q: %s", projectID, billingAccount, err)
			}
		}
	}

	if len(services) > 0 {

		enabledServices, err := resourcemanager.ListCurrentlyEnabledServices(projectID, "", config.UserAgent, config, 1*time.Minute)
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
			if err := resourcemanager.EnableServiceUsageProjectServices(servicesToEnable, projectID, "", config.UserAgent, config, 10*time.Minute); err != nil {
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
		Credentials: envvar.GetTestCredsFromEnv(),
		Project:     envvar.GetTestProjectFromEnv(),
		Region:      envvar.GetTestRegionFromEnv(),
		Zone:        envvar.GetTestZoneFromEnv(),
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
	project := envvar.GetTestProjectFromEnv()

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
		err = transport_tpg.Retry(transport_tpg.RetryOptions{
			RetryFunc: func() (operr error) {
				op, operr = config.NewSqlAdminClient(config.UserAgent).Instances.Insert(project, bootstrapInstance).Do()
				return operr
			},
			Timeout:              20 * time.Minute,
			ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsSqlOperationInProgressError},
		})
		if err != nil {
			t.Fatalf("Error, failed to create instance %s: %s", bootstrapInstance.Name, err)
		}
		err = sql.SqlAdminOperationWaitTime(config, op, project, "Create Instance", config.UserAgent, 40*time.Minute)
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
		err = transport_tpg.Retry(transport_tpg.RetryOptions{
			RetryFunc: func() (operr error) {
				op, operr = config.NewSqlAdminClient(config.UserAgent).BackupRuns.Insert(project, bootstrapInstance.Name, backupRun).Do()
				return operr
			},
			Timeout:              20 * time.Minute,
			ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{transport_tpg.IsSqlOperationInProgressError},
		})
		if err != nil {
			t.Fatalf("Error, failed to create instance backup: %s", err)
		}
		err = sql.SqlAdminOperationWaitTime(config, op, project, "Backup Instance", config.UserAgent, 20*time.Minute)
		if err != nil {
			t.Fatalf("Error, failed to create instance backup: %s", err)
		}
	}

	return bootstrapInstance.Name
}

func BootstrapSharedCaPoolInLocation(t *testing.T, location string) string {
	project := envvar.GetTestProjectFromEnv()
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
		err = privateca.PrivatecaOperationWaitTimeWithResponse(
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

func BootstrapSubnet(t *testing.T, subnetName string, networkName string) string {
	projectID := envvar.GetTestProjectFromEnv()
	region := envvar.GetTestRegionFromEnv()

	config := BootstrapConfig(t)
	if config == nil {
		t.Fatal("Could not bootstrap config.")
	}

	computeService := config.NewComputeClient(config.UserAgent)
	if computeService == nil {
		t.Fatal("Could not create compute client.")
	}

	// In order to create a networkAttachment we need to bootstrap a subnet.
	_, err := computeService.Subnetworks.Get(projectID, region, subnetName).Do()
	if err != nil && transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
		log.Printf("[DEBUG] Subnet %q not found, bootstrapping", subnetName)

		networkUrl := fmt.Sprintf("%sprojects/%s/global/networks/%s", config.ComputeBasePath, projectID, networkName)
		url := fmt.Sprintf("%sprojects/%s/regions/%s/subnetworks", config.ComputeBasePath, projectID, region)

		subnetObj := map[string]interface{}{
			"name":        subnetName,
			"region ":     region,
			"network":     networkUrl,
			"ipCidrRange": "10.77.0.0/20",
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   projectID,
			RawURL:    url,
			UserAgent: config.UserAgent,
			Body:      subnetObj,
			Timeout:   4 * time.Minute,
		})

		log.Printf("Response is, %s", res)
		if err != nil {
			t.Fatalf("Error bootstrapping test subnet %s: %s", subnetName, err)
		}

		log.Printf("[DEBUG] Waiting for network creation to finish")
		err = tpgcompute.ComputeOperationWaitTime(config, res, projectID, "Error bootstrapping test subnet", config.UserAgent, 4*time.Minute)
		if err != nil {
			t.Fatalf("Error bootstrapping test subnet %s: %s", subnetName, err)
		}
	}

	subnet, err := computeService.Subnetworks.Get(projectID, region, subnetName).Do()

	if subnet == nil {
		t.Fatalf("Error getting test subnet %s: is nil", subnetName)
	}

	if err != nil {
		t.Fatalf("Error getting test subnet %s: %s", subnetName, err)
	}
	return subnet.Name
}

func BootstrapNetworkAttachment(t *testing.T, networkAttachmentName string, subnetName string) string {
	projectID := envvar.GetTestProjectFromEnv()
	region := envvar.GetTestRegionFromEnv()

	config := BootstrapConfig(t)
	if config == nil {
		return ""
	}

	computeService := config.NewComputeClient(config.UserAgent)
	if computeService == nil {
		return ""
	}

	networkAttachment, err := computeService.NetworkAttachments.Get(projectID, region, networkAttachmentName).Do()
	if err != nil && transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
		// Create Network Attachment Here.
		log.Printf("[DEBUG] Network Attachment %s not found, bootstrapping", networkAttachmentName)
		url := fmt.Sprintf("%sprojects/%s/regions/%s/networkAttachments", config.ComputeBasePath, projectID, region)

		subnetURL := fmt.Sprintf("%sprojects/%s/regions/%s/subnetworks/%s", config.ComputeBasePath, projectID, region, subnetName)
		networkAttachmentObj := map[string]interface{}{
			"name":                 networkAttachmentName,
			"region":               region,
			"subnetworks":          []string{subnetURL},
			"connectionPreference": "ACCEPT_AUTOMATIC",
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   projectID,
			RawURL:    url,
			UserAgent: config.UserAgent,
			Body:      networkAttachmentObj,
			Timeout:   4 * time.Minute,
		})
		if err != nil {
			t.Fatalf("Error bootstrapping test Network Attachment %s: %s", networkAttachmentName, err)
		}

		log.Printf("[DEBUG] Waiting for network creation to finish")
		err = tpgcompute.ComputeOperationWaitTime(config, res, projectID, "Error bootstrapping shared test subnet", config.UserAgent, 4*time.Minute)
		if err != nil {
			t.Fatalf("Error bootstrapping test Network Attachment %s: %s", networkAttachmentName, err)
		}
	}

	networkAttachment, err = computeService.NetworkAttachments.Get(projectID, region, networkAttachmentName).Do()

	if networkAttachment == nil {
		t.Fatalf("Error getting test network attachment %s: is nil", networkAttachmentName)
	}

	if err != nil {
		t.Fatalf("Error getting test Network Attachment %s: %s", networkAttachmentName, err)
	}

	return networkAttachment.Name
}

// The default network within GCP already comes pre configured with
// certain firewall rules open to allow internal communication. As we
// are boostrapping a network for dataproc tests, we need to additionally
// open up similar rules to allow the nodes to talk to each other
// internally as part of their configuration or this will just hang.
const SharedTestFirewallPrefix = "tf-bootstrap-firewall-"

func BootstrapFirewallForDataprocSharedNetwork(t *testing.T, firewallName string, networkName string) string {
	project := envvar.GetTestProjectFromEnv()
	firewallName = SharedTestFirewallPrefix + firewallName

	config := BootstrapConfig(t)
	if config == nil {
		return ""
	}

	log.Printf("[DEBUG] Getting Firewall %q for Network %q", firewallName, networkName)
	_, err := config.NewComputeClient(config.UserAgent).Firewalls.Get(project, firewallName).Do()
	if err != nil && transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
		log.Printf("[DEBUG] firewallName %q not found, bootstrapping", firewallName)
		url := fmt.Sprintf("%sprojects/%s/global/firewalls", config.ComputeBasePath, project)

		networkId := fmt.Sprintf("projects/%s/global/networks/%s", project, networkName)
		allowObj := []interface{}{
			map[string]interface{}{
				"IPProtocol": "icmp",
			},
			map[string]interface{}{
				"IPProtocol": "tcp",
				"ports":      []string{"0-65535"},
			},
			map[string]interface{}{
				"IPProtocol": "udp",
				"ports":      []string{"0-65535"},
			},
		}

		firewallObj := map[string]interface{}{
			"name":    firewallName,
			"network": networkId,
			"allowed": allowObj,
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   project,
			RawURL:    url,
			UserAgent: config.UserAgent,
			Body:      firewallObj,
			Timeout:   4 * time.Minute,
		})
		if err != nil {
			t.Fatalf("Error bootstrapping Firewall %q for Network %q: %s", firewallName, networkName, err)
		}

		log.Printf("[DEBUG] Waiting for Firewall creation to finish")
		err = tpgcompute.ComputeOperationWaitTime(config, res, project, "Error bootstrapping Firewall", config.UserAgent, 4*time.Minute)
		if err != nil {
			t.Fatalf("Error bootstrapping Firewall %q: %s", firewallName, err)
		}
	}

	firewall, err := config.NewComputeClient(config.UserAgent).Firewalls.Get(project, firewallName).Do()
	if err != nil {
		t.Errorf("Error getting Firewall %q: %s", firewallName, err)
	}
	if firewall == nil {
		t.Fatalf("Error getting Firewall %q: is nil", firewallName)
	}
	return firewall.Name
}

const SharedStoragePoolPrefix = "tf-bootstrap-storage-pool-"

func BootstrapComputeStoragePool(t *testing.T, storagePoolName, storagePoolType string) string {
	projectID := envvar.GetTestProjectFromEnv()
	zone := envvar.GetTestZoneFromEnv()

	storagePoolName = SharedStoragePoolPrefix + storagePoolType + "-" + storagePoolName

	config := BootstrapConfig(t)
	if config == nil {
		t.Fatal("Could not bootstrap config.")
	}

	computeService := config.NewComputeClient(config.UserAgent)
	if computeService == nil {
		t.Fatal("Could not create compute client.")
	}

	_, err := computeService.StoragePools.Get(projectID, zone, storagePoolName).Do()
	if err != nil && transport_tpg.IsGoogleApiErrorWithCode(err, 404) {
		log.Printf("[DEBUG] Storage pool %q not found, bootstrapping", storagePoolName)

		url := fmt.Sprintf("%sprojects/%s/zones/%s/storagePools", config.ComputeBasePath, projectID, zone)
		storagePoolTypeUrl := fmt.Sprintf("/projects/%s/zones/%s/storagePoolTypes/%s", projectID, zone, storagePoolType)

		storagePoolObj := map[string]interface{}{
			"name":                      storagePoolName,
			"poolProvisionedCapacityGb": 10240,
			"poolProvisionedThroughput": 180,
			"storagePoolType":           storagePoolTypeUrl,
			"capacityProvisioningType":  "ADVANCED",
		}

		if storagePoolType == "hyperdisk-balanced" {
			storagePoolObj["poolProvisionedIops"] = 10000
			storagePoolObj["poolProvisionedThroughput"] = 1024
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   projectID,
			RawURL:    url,
			UserAgent: config.UserAgent,
			Body:      storagePoolObj,
			Timeout:   20 * time.Minute,
		})

		log.Printf("Response is, %s", res)
		if err != nil {
			t.Fatalf("Error bootstrapping storage pool %s: %s", storagePoolName, err)
		}

		log.Printf("[DEBUG] Waiting for storage pool creation to finish")
		err = tpgcompute.ComputeOperationWaitTime(config, res, projectID, "Error bootstrapping storage pool", config.UserAgent, 4*time.Minute)
		if err != nil {
			t.Fatalf("Error bootstrapping test storage pool %s: %s", storagePoolName, err)
		}
	}

	storagePool, err := computeService.StoragePools.Get(projectID, zone, storagePoolName).Do()

	if storagePool == nil {
		t.Fatalf("Error getting storage pool %s: is nil", storagePoolName)
	}

	if err != nil {
		t.Fatalf("Error getting storage pool %s: %s", storagePoolName, err)
	}
	return storagePool.SelfLink
}

func SetupProjectsAndGetAccessToken(org, billing, pid, service string, config *transport_tpg.Config) (string, error) {
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
	err := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() (reqErr error) {
			op, reqErr = rmService.Projects.Create(project).Do()
			return reqErr
		},
		Timeout: 5 * time.Minute,
	})
	if err != nil {
		return "", err
	}

	// Wait for the operation to complete
	opAsMap, err := tpgresource.ConvertToMap(op)
	if err != nil {
		return "", err
	}

	waitErr := resourcemanager.ResourceManagerOperationWaitTime(config, opAsMap, "creating project", config.UserAgent, 5*time.Minute)
	if waitErr != nil {
		return "", waitErr
	}

	ba := &cloudbilling.ProjectBillingInfo{
		BillingAccountName: fmt.Sprintf("billingAccounts/%s", billing),
	}
	_, err = config.NewBillingClient(config.UserAgent).Projects.UpdateBillingInfo(resourcemanager.PrefixedProject(pid), ba).Do()
	if err != nil {
		return "", err
	}

	p2 := fmt.Sprintf("%s-2", pid)
	project.ProjectId = p2
	project.Name = fmt.Sprintf("%s-2", pid)

	err = transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() (reqErr error) {
			op, reqErr = rmService.Projects.Create(project).Do()
			return reqErr
		},
		Timeout: 5 * time.Minute,
	})
	if err != nil {
		return "", err
	}

	// Wait for the operation to complete
	opAsMap, err = tpgresource.ConvertToMap(op)
	if err != nil {
		return "", err
	}

	waitErr = resourcemanager.ResourceManagerOperationWaitTime(config, opAsMap, "creating project", config.UserAgent, 5*time.Minute)
	if waitErr != nil {
		return "", waitErr
	}

	_, err = config.NewBillingClient(config.UserAgent).Projects.UpdateBillingInfo(resourcemanager.PrefixedProject(p2), ba).Do()
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
	serviceAccountEmail := serviceAccountPrefix + service
	sa1, err := getOrCreateServiceAccount(config, pid, serviceAccountEmail)
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

const sharedTagKeyPrefix = "tf-bootstrap-tagkey"

func BootstrapSharedTestTagKey(t *testing.T, testId string) string {
	org := envvar.GetTestOrgFromEnv(t)
	sharedTagKey := fmt.Sprintf("%s-%s", sharedTagKeyPrefix, testId)
	tagKeyName := fmt.Sprintf("%s/%s", org, sharedTagKey)

	config := BootstrapConfig(t)
	if config == nil {
		return ""
	}

	log.Printf("[DEBUG] Getting shared test tag key %q", sharedTagKey)
	getURL := fmt.Sprintf("%stagKeys/namespaced?name=%s", config.TagsBasePath, tagKeyName)
	_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   config.Project,
		RawURL:    getURL,
		UserAgent: config.UserAgent,
		Timeout:   2 * time.Minute,
	})
	if err != nil && transport_tpg.IsGoogleApiErrorWithCode(err, 403) {
		log.Printf("[DEBUG] TagKey %q not found, bootstrapping", sharedTagKey)
		tagKeyObj := map[string]interface{}{
			"parent":      "organizations/" + org,
			"shortName":   sharedTagKey,
			"description": "Bootstrapped tag key for Terraform Acceptance testing",
		}

		_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   config.Project,
			RawURL:    config.TagsBasePath + "tagKeys/",
			UserAgent: config.UserAgent,
			Body:      tagKeyObj,
			Timeout:   10 * time.Minute,
		})
		if err != nil {
			t.Fatalf("Error bootstrapping shared tag key %q: %s", sharedTagKey, err)
		}

		log.Printf("[DEBUG] Waiting for shared tag key creation to finish")
	}

	_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   config.Project,
		RawURL:    getURL,
		UserAgent: config.UserAgent,
		Timeout:   2 * time.Minute,
	})

	if err != nil {
		t.Fatalf("Error getting shared tag key %q: %s", sharedTagKey, err)
	}

	return sharedTagKey
}

const sharedTagValuePrefix = "tf-bootstrap-tagvalue"

func BootstrapSharedTestTagValue(t *testing.T, testId string, tagKey string) string {
	org := envvar.GetTestOrgFromEnv(t)
	sharedTagValue := fmt.Sprintf("%s-%s", sharedTagValuePrefix, testId)
	tagKeyName := fmt.Sprintf("%s/%s", org, tagKey)
	tagValueName := fmt.Sprintf("%s/%s", tagKeyName, sharedTagValue)

	config := BootstrapConfig(t)
	if config == nil {
		return ""
	}

	log.Printf("[DEBUG] Getting shared test tag value %q", sharedTagValue)
	getURL := fmt.Sprintf("%stagValues/namespaced?name=%s", config.TagsBasePath, tagValueName)
	_, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   config.Project,
		RawURL:    getURL,
		UserAgent: config.UserAgent,
		Timeout:   2 * time.Minute,
	})
	if err != nil && transport_tpg.IsGoogleApiErrorWithCode(err, 403) {
		log.Printf("[DEBUG] TagValue %q not found, bootstrapping", sharedTagValue)
		log.Printf("[DEBUG] Fetching permanent id for tagkey %s", tagKeyName)
		tagKeyGetURL := fmt.Sprintf("%stagKeys/namespaced?name=%s", config.TagsBasePath, tagKeyName)
		tagKeyResponse, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   config.Project,
			RawURL:    tagKeyGetURL,
			UserAgent: config.UserAgent,
			Timeout:   2 * time.Minute,
		})
		if err != nil {
			t.Fatalf("Error getting tag key id for %s : %s", tagKeyName, err)
		}
		tagKeyObj := map[string]interface{}{
			"parent":      tagKeyResponse["name"].(string),
			"shortName":   sharedTagValue,
			"description": "Bootstrapped tag value for Terraform Acceptance testing",
		}

		_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "POST",
			Project:   config.Project,
			RawURL:    config.TagsBasePath + "tagValues/",
			UserAgent: config.UserAgent,
			Body:      tagKeyObj,
			Timeout:   10 * time.Minute,
		})
		if err != nil {
			t.Fatalf("Error bootstrapping shared tag value %q: %s", sharedTagValue, err)
		}

		log.Printf("[DEBUG] Waiting for shared tag value creation to finish")
	}

	_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   config.Project,
		RawURL:    getURL,
		UserAgent: config.UserAgent,
		Timeout:   2 * time.Minute,
	})

	if err != nil {
		t.Fatalf("Error getting shared tag value %q: %s", sharedTagValue, err)
	}

	return sharedTagValue
}
