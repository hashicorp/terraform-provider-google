// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package acctest

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Deprecated: For backward compatibility TestEnvVar is still working,
// but all new code should use TestEnvVar in the envvar package instead.
const TestEnvVar = envvar.TestEnvVar

// Deprecated: For backward compatibility CredsEnvVars is still working,
// but all new code should use CredsEnvVars in the envvar package instead.
var CredsEnvVars = envvar.CredsEnvVars

// Deprecated: For backward compatibility ProjectNumberEnvVars is still working,
// but all new code should use ProjectNumberEnvVars in the envvar package instead.
var ProjectNumberEnvVars = envvar.ProjectNumberEnvVars

// Deprecated: For backward compatibility ProjectEnvVars is still working,
// but all new code should use ProjectEnvVars in the envvar package instead.
var ProjectEnvVars = envvar.ProjectEnvVars

// Deprecated: For backward compatibility FirestoreProjectEnvVars is still working,
// but all new code should use FirestoreProjectEnvVars in the envvar package instead.
var FirestoreProjectEnvVars = envvar.FirestoreProjectEnvVars

// Deprecated: For backward compatibility RegionEnvVars is still working,
// but all new code should use RegionEnvVars in the envvar package instead.
var RegionEnvVars = envvar.RegionEnvVars

// Deprecated: For backward compatibility ZoneEnvVars is still working,
// but all new code should use ZoneEnvVars in the envvar package instead.
var ZoneEnvVars = envvar.ZoneEnvVars

// Deprecated: For backward compatibility OrgEnvVars is still working,
// but all new code should use OrgEnvVars in the envvar package instead.
var OrgEnvVars = envvar.OrgEnvVars

// This value is the Customer ID of the GOOGLE_ORG_DOMAIN workspace.
// See https://admin.google.com/ac/accountsettings when logged into an org admin for the value.
//
// Deprecated: For backward compatibility CustIdEnvVars is still working,
// but all new code should use CustIdEnvVars in the envvar package instead.
var CustIdEnvVars = envvar.CustIdEnvVars

// This value is the username of an identity account within the GOOGLE_ORG_DOMAIN workspace.
// For example in the org example.com with a user "foo@example.com", this would be set to "foo".
// See https://admin.google.com/ac/users when logged into an org admin for a list.
//
// Deprecated: For backward compatibility IdentityUserEnvVars is still working,
// but all new code should use IdentityUserEnvVars in the envvar package instead.
var IdentityUserEnvVars = envvar.IdentityUserEnvVars

// Deprecated: For backward compatibility OrgEnvDomainVars is still working,
// but all new code should use OrgEnvDomainVars in the envvar package instead.
var OrgEnvDomainVars = envvar.OrgEnvDomainVars

// Deprecated: For backward compatibility ServiceAccountEnvVars is still working,
// but all new code should use ServiceAccountEnvVars in the envvar package instead.
var ServiceAccountEnvVars = envvar.ServiceAccountEnvVars

// Deprecated: For backward compatibility OrgTargetEnvVars is still working,
// but all new code should use OrgTargetEnvVars in the envvar package instead.
var OrgTargetEnvVars = envvar.OrgTargetEnvVars

// This is the billing account that will be charged for the infrastructure used during testing. For
// that reason, it is also the billing account used for creating new projects.
//
// Deprecated: For backward compatibility BillingAccountEnvVars is still working,
// but all new code should use BillingAccountEnvVars in the envvar package instead.
var BillingAccountEnvVars = envvar.BillingAccountEnvVars

// This is the billing account that will be modified to test billing-related functionality. It is
// expected to have more permissions granted to the test user and support subaccounts.
//
// Deprecated: For backward compatibility MasterBillingAccountEnvVars is still working,
// but all new code should use MasterBillingAccountEnvVars in the envvar package instead.
var MasterBillingAccountEnvVars = envvar.MasterBillingAccountEnvVars

// This value is the description used for test PublicAdvertisedPrefix setup to avoid required DNS
// setup. This is only used during integration tests and would be invalid to surface to users
//
// Deprecated: For backward compatibility PapDescriptionEnvVars is still working,
// but all new code should use PapDescriptionEnvVars in the envvar package instead.
var PapDescriptionEnvVars = envvar.PapDescriptionEnvVars

func AccTestPreCheck(t *testing.T) {
	if v := os.Getenv("GOOGLE_CREDENTIALS_FILE"); v != "" {
		creds, err := ioutil.ReadFile(v)
		if err != nil {
			t.Fatalf("Error reading GOOGLE_CREDENTIALS_FILE path: %s", err)
		}
		os.Setenv("GOOGLE_CREDENTIALS", string(creds))
	}

	if v := transport_tpg.MultiEnvSearch(CredsEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(CredsEnvVars, ", "))
	}

	if v := transport_tpg.MultiEnvSearch(ProjectEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(ProjectEnvVars, ", "))
	}

	if v := transport_tpg.MultiEnvSearch(RegionEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(RegionEnvVars, ", "))
	}

	if v := transport_tpg.MultiEnvSearch(ZoneEnvVars); v == "" {
		t.Fatalf("One of %s must be set for acceptance tests", strings.Join(ZoneEnvVars, ", "))
	}
}

// GetTestRegion has the same logic as the provider's GetRegion, to be used in tests.
func GetTestRegion(is *terraform.InstanceState, config *transport_tpg.Config) (string, error) {
	if res, ok := is.Attributes["region"]; ok {
		return res, nil
	}
	if config.Region != "" {
		return config.Region, nil
	}
	return "", fmt.Errorf("%q: required field is not set", "region")
}

// GetTestProject has the same logic as the provider's GetProject, to be used in tests.
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
//
// Deprecated: For backward compatibility GetTestProjectNumberFromEnv is still working,
// but all new code should use GetTestProjectNumberFromEnv in the envvar package instead.
func GetTestProjectNumberFromEnv() string {
	return envvar.GetTestProjectNumberFromEnv()
}

// AccTestPreCheck ensures at least one of the project env variables is set.
//
// Deprecated: For backward compatibility GetTestProjectFromEnv is still working,
// but all new code should use GetTestProjectFromEnv in the envvar package instead.
func GetTestProjectFromEnv() string {
	return envvar.GetTestProjectFromEnv()
}

// AccTestPreCheck ensures at least one of the credentials env variables is set.
//
// Deprecated: For backward compatibility GetTestCredsFromEnv is still working,
// but all new code should use GetTestCredsFromEnv in the envvar package instead.
func GetTestCredsFromEnv() string {
	return envvar.GetTestCredsFromEnv()
}

// AccTestPreCheck ensures at least one of the region env variables is set.
//
// Deprecated: For backward compatibility GetTestRegionFromEnv is still working,
// but all new code should use GetTestRegionFromEnv in the envvar package instead.
func GetTestRegionFromEnv() string {
	return envvar.GetTestRegionFromEnv()
}

// Deprecated: For backward compatibility GetTestZoneFromEnv is still working,
// but all new code should use GetTestZoneFromEnv in the envvar package instead.
func GetTestZoneFromEnv() string {
	return envvar.GetTestZoneFromEnv()
}

// Deprecated: For backward compatibility GetTestCustIdFromEnv is still working,
// but all new code should use GetTestCustIdFromEnv in the envvar package instead.
func GetTestCustIdFromEnv(t *testing.T) string {
	return envvar.GetTestCustIdFromEnv(t)
}

// Deprecated: For backward compatibility GetTestIdentityUserFromEnv is still working,
// but all new code should use GetTestIdentityUserFromEnv in the envvar package instead.
func GetTestIdentityUserFromEnv(t *testing.T) string {
	return envvar.GetTestIdentityUserFromEnv(t)
}

// Firestore can't be enabled at the same time as Datastore, so we need a new
// project to manage it until we can enable Firestore programmatically.
//
// Deprecated: For backward compatibility GetTestFirestoreProjectFromEnv is still working,
// but all new code should use GetTestFirestoreProjectFromEnv in the envvar package instead.
func GetTestFirestoreProjectFromEnv(t *testing.T) string {
	return envvar.GetTestFirestoreProjectFromEnv(t)
}

// Returns the raw organization id like 1234567890, skipping the test if one is
// not found.
//
// Deprecated: For backward compatibility GetTestOrgFromEnv is still working,
// but all new code should use GetTestOrgFromEnv in the envvar package instead.
func GetTestOrgFromEnv(t *testing.T) string {
	return envvar.GetTestOrgFromEnv(t)
}

// Alternative to GetTestOrgFromEnv that doesn't need *testing.T
// If using this, you need to process unset values at the call site
//
// Deprecated: For backward compatibility UnsafeGetTestOrgFromEnv is still working,
// but all new code should use UnsafeGetTestOrgFromEnv in the envvar package instead.
func UnsafeGetTestOrgFromEnv() string {
	return envvar.UnsafeGetTestOrgFromEnv()
}

// Deprecated: For backward compatibility GetTestOrgDomainFromEnv is still working,
// but all new code should use GetTestOrgDomainFromEnv in the envvar package instead.
func GetTestOrgDomainFromEnv(t *testing.T) string {
	return envvar.GetTestOrgDomainFromEnv(t)
}

// Deprecated: For backward compatibility GetTestOrgTargetFromEnv is still working,
// but all new code should use GetTestOrgTargetFromEnv in the envvar package instead.
func GetTestOrgTargetFromEnv(t *testing.T) string {
	return envvar.GetTestOrgTargetFromEnv(t)
}

// This is the billing account that will be charged for the infrastructure used during testing. For
// that reason, it is also the billing account used for creating new projects.
//
// Deprecated: For backward compatibility GetTestBillingAccountFromEnv is still working,
// but all new code should use GetTestBillingAccountFromEnv in the envvar package instead.
func GetTestBillingAccountFromEnv(t *testing.T) string {
	return envvar.GetTestBillingAccountFromEnv(t)
}

// This is the billing account that will be modified to test billing-related functionality. It is
// expected to have more permissions granted to the test user and support subaccounts.
//
// Deprecated: For backward compatibility GetTestMasterBillingAccountFromEnv is still working,
// but all new code should use GetTestMasterBillingAccountFromEnv in the envvar package instead.
func GetTestMasterBillingAccountFromEnv(t *testing.T) string {
	return envvar.GetTestMasterBillingAccountFromEnv(t)
}

// Deprecated: For backward compatibility GetTestServiceAccountFromEnv is still working,
// but all new code should use GetTestServiceAccountFromEnv in the envvar package instead.
func GetTestServiceAccountFromEnv(t *testing.T) string {
	return envvar.GetTestServiceAccountFromEnv(t)
}

// Deprecated: For backward compatibility GetTestPublicAdvertisedPrefixDescriptionFromEnv is still working,
// but all new code should use GetTestPublicAdvertisedPrefixDescriptionFromEnv in the envvar package instead.
func GetTestPublicAdvertisedPrefixDescriptionFromEnv(t *testing.T) string {
	return envvar.GetTestPublicAdvertisedPrefixDescriptionFromEnv(t)
}

// Some tests fail during VCR. One common case is race conditions when creating resources.
// If a test config adds two fine-grained resources with the same parent it is undefined
// which will be created first, causing VCR to fail ~50% of the time
func SkipIfVcr(t *testing.T) {
	if IsVcrEnabled() {
		t.Skipf("VCR enabled, skipping test: %s", t.Name())
	}
}

func SleepInSecondsForTest(t int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		time.Sleep(time.Duration(t) * time.Second)
		return nil
	}
}

// Deprecated: For backward compatibility SkipIfEnvNotSet is still working,
// but all new code should use SkipIfEnvNotSet in the envvar package instead.
func SkipIfEnvNotSet(t *testing.T, envs ...string) {
	envvar.SkipIfEnvNotSet(t, envs...)
}
