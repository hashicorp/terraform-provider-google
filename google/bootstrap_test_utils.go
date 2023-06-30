// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

type BootstrappedKMS = acctest.BootstrappedKMS

func BootstrapKMSKey(t *testing.T) acctest.BootstrappedKMS {
	return acctest.BootstrapKMSKey(t)
}

func BootstrapKMSKeyInLocation(t *testing.T, locationID string) acctest.BootstrappedKMS {
	return acctest.BootstrapKMSKeyInLocation(t, locationID)
}

// BootstrapKMSKeyWithPurpose returns a KMS key in the "global" location.
// See BootstrapKMSKeyWithPurposeInLocation.
func BootstrapKMSKeyWithPurpose(t *testing.T, purpose string) acctest.BootstrappedKMS {
	return acctest.BootstrapKMSKeyWithPurpose(t, purpose)
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
func BootstrapKMSKeyWithPurposeInLocation(t *testing.T, purpose, locationID string) acctest.BootstrappedKMS {
	return acctest.BootstrapKMSKeyWithPurposeInLocation(t, purpose, locationID)
}

func BootstrapKMSKeyWithPurposeInLocationAndName(t *testing.T, purpose, locationID, keyShortName string) acctest.BootstrappedKMS {
	return acctest.BootstrapKMSKeyWithPurposeInLocationAndName(t, purpose, locationID, keyShortName)
}

func BootstrapServiceAccount(t *testing.T, project, testRunner string) string {
	return acctest.BootstrapServiceAccount(t, project, testRunner)
}

func BootstrapSharedTestADDomain(t *testing.T, testId string, networkName string) string {
	return acctest.BootstrapSharedTestADDomain(t, testId, networkName)
}

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
	return acctest.BootstrapSharedTestNetwork(t, testId)
}

func BootstrapServicePerimeterProjects(t *testing.T, desiredProjects int) []*cloudresourcemanager.Project {
	return acctest.BootstrapServicePerimeterProjects(t, desiredProjects)
}

func RemoveContainerServiceAgentRoleFromContainerEngineRobot(t *testing.T, project *cloudresourcemanager.Project) {
	acctest.RemoveContainerServiceAgentRoleFromContainerEngineRobot(t, project)
}

// BootstrapProject will create or get a project named
// "<projectIDPrefix><projectIDSuffix>" that will persist across test runs,
// where projectIDSuffix is based off of getTestProjectFromEnv(). The reason
// for the naming is to isolate bootstrapped projects by test environment.
// Given the existing projects being used by our team, the prefix provided to
// this function can be no longer than 18 characters.
func BootstrapProject(t *testing.T, projectIDPrefix, billingAccount string, services []string) *cloudresourcemanager.Project {
	return acctest.BootstrapProject(t, projectIDPrefix, billingAccount, services)
}

// BootstrapConfig returns a Config pulled from the environment.
func BootstrapConfig(t *testing.T) *transport_tpg.Config {
	return acctest.BootstrapConfig(t)
}

// BootstrapSharedSQLInstanceBackupRun will return a shared SQL db instance that
// has a backup created for it.
func BootstrapSharedSQLInstanceBackupRun(t *testing.T) string {
	return acctest.BootstrapSharedSQLInstanceBackupRun(t)
}

func BootstrapSharedCaPoolInLocation(t *testing.T, location string) string {
	return acctest.BootstrapSharedCaPoolInLocation(t, location)
}

func BootstrapSubnet(t *testing.T, subnetName string, networkName string) string {
	return acctest.BootstrapSubnet(t, subnetName, networkName)
}

func BootstrapNetworkAttachment(t *testing.T, networkAttachmentName string, subnetName string) string {
	return acctest.BootstrapNetworkAttachment(t, networkAttachmentName, subnetName)
}

func setupProjectsAndGetAccessToken(org, billing, pid, service string, config *transport_tpg.Config) (string, error) {
	return acctest.SetupProjectsAndGetAccessToken(org, billing, pid, service, config)
}
