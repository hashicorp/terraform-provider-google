// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var runOrgIamAuditConfigTestEnvVar = "TF_RUN_ORG_IAM_AUDIT_CONFIG"

func organizationIamAuditConfigImportStep(resourceName, org, service string) resource.TestStep {
	return resource.TestStep{
		ResourceName:      resourceName,
		ImportStateId:     fmt.Sprintf("%s %s", org, service),
		ImportState:       true,
		ImportStateVerify: true,
	}
}

// Test that an IAM audit config can be applied to an organization
func TestAccOrganizationIamAuditConfig_basic(t *testing.T) {
	if os.Getenv(runOrgIamAuditConfigTestEnvVar) != "true" {
		t.Skipf("Environment variable %s is not set, skipping.", runOrgIamAuditConfigTestEnvVar)
	}
	org := envvar.GetTestOrgFromEnv(t)
	service := "cloudkms.googleapis.com"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Apply an IAM audit config
			{
				Config: testAccOrganizationAssociateAuditConfigBasic(org, service),
			},
			organizationIamAuditConfigImportStep("google_organization_iam_audit_config.acceptance", org, service),
		},
	})
}

// Test that multiple IAM audit configs can be applied to an organization, one at a time
func TestAccOrganizationIamAuditConfig_multiple(t *testing.T) {
	if os.Getenv(runOrgIamAuditConfigTestEnvVar) != "true" {
		t.Skipf("Environment variable %s is not set, skipping.", runOrgIamAuditConfigTestEnvVar)
	}
	org := envvar.GetTestOrgFromEnv(t)
	service := "cloudkms.googleapis.com"
	service2 := "cloudsql.googleapis.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Apply an IAM audit config
			{
				Config: testAccOrganizationAssociateAuditConfigBasic(org, service),
			},
			// Apply another IAM audit config
			{
				Config: testAccOrganizationAssociateAuditConfigMultiple(org, service, service2),
			},
			organizationIamAuditConfigImportStep("google_organization_iam_audit_config.acceptance", org, service),
			organizationIamAuditConfigImportStep("google_organization_iam_audit_config.multiple", org, service2),
		},
	})
}

// Test that multiple IAM audit configs can be applied to an organization all at once
func TestAccOrganizationIamAuditConfig_multipleAtOnce(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	if os.Getenv(runOrgIamAuditConfigTestEnvVar) != "true" {
		t.Skipf("Environment variable %s is not set, skipping.", runOrgIamAuditConfigTestEnvVar)
	}
	org := envvar.GetTestOrgFromEnv(t)
	service := "cloudkms.googleapis.com"
	service2 := "cloudsql.googleapis.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Apply an IAM audit config
			{
				Config: testAccOrganizationAssociateAuditConfigMultiple(org, service, service2),
			},
			organizationIamAuditConfigImportStep("google_organization_iam_audit_config.acceptance", org, service),
			organizationIamAuditConfigImportStep("google_organization_iam_audit_config.multiple", org, service2),
		},
	})
}

// Test that an IAM audit config can be updated once applied to an organization
func TestAccOrganizationIamAuditConfig_update(t *testing.T) {
	if os.Getenv(runOrgIamAuditConfigTestEnvVar) != "true" {
		t.Skipf("Environment variable %s is not set, skipping.", runOrgIamAuditConfigTestEnvVar)
	}
	org := envvar.GetTestOrgFromEnv(t)
	service := "cloudkms.googleapis.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Apply an IAM audit config
			{
				Config: testAccOrganizationAssociateAuditConfigBasic(org, service),
			},
			organizationIamAuditConfigImportStep("google_organization_iam_audit_config.acceptance", org, service),

			// Apply an updated IAM audit config
			{
				Config: testAccOrganizationAssociateAuditConfigUpdated(org, service),
			},
			organizationIamAuditConfigImportStep("google_organization_iam_audit_config.acceptance", org, service),

			// Drop the original member
			{
				Config: testAccOrganizationAssociateAuditConfigDropMemberFromBasic(org, service),
			},
			organizationIamAuditConfigImportStep("google_organization_iam_audit_config.acceptance", org, service),
		},
	})
}

// Test that an IAM audit config can be removed from an organization
func TestAccOrganizationIamAuditConfig_remove(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	if os.Getenv(runOrgIamAuditConfigTestEnvVar) != "true" {
		t.Skipf("Environment variable %s is not set, skipping.", runOrgIamAuditConfigTestEnvVar)
	}
	org := envvar.GetTestOrgFromEnv(t)
	service := "cloudkms.googleapis.com"
	service2 := "cloudsql.googleapis.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Apply multiple IAM audit configs
			{
				Config: testAccOrganizationAssociateAuditConfigMultiple(org, service, service2),
			},
			organizationIamAuditConfigImportStep("google_organization_iam_audit_config.acceptance", org, service),
			organizationIamAuditConfigImportStep("google_organization_iam_audit_config.multiple", org, service2),

			// Remove one IAM audit config
			{
				Config: testAccOrganizationAssociateAuditConfigBasic(org, service),
			},
			organizationIamAuditConfigImportStep("google_organization_iam_audit_config.acceptance", org, service),
		},
	})
}

// Test adding exempt first exempt member
func TestAccOrganizationIamAuditConfig_addFirstExemptMember(t *testing.T) {
	if os.Getenv(runOrgIamAuditConfigTestEnvVar) != "true" {
		t.Skipf("Environment variable %s is not set, skipping.", runOrgIamAuditConfigTestEnvVar)
	}
	org := envvar.GetTestOrgFromEnv(t)
	service := "cloudkms.googleapis.com"
	members := []string{}
	members2 := []string{"user:gterraformtest1@gmail.com"}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Apply IAM audit config with no members
			{
				Config: testAccOrganizationAssociateAuditConfigMembers(org, service, members),
			},
			organizationIamAuditConfigImportStep("google_organization_iam_audit_config.acceptance", org, service),

			// Apply IAM audit config with one member
			{
				Config: testAccOrganizationAssociateAuditConfigMembers(org, service, members2),
			},
			organizationIamAuditConfigImportStep("google_organization_iam_audit_config.acceptance", org, service),
		},
	})
}

// test removing last exempt member
func TestAccOrganizationIamAuditConfig_removeLastExemptMember(t *testing.T) {
	if os.Getenv(runOrgIamAuditConfigTestEnvVar) != "true" {
		t.Skipf("Environment variable %s is not set, skipping.", runOrgIamAuditConfigTestEnvVar)
	}
	org := envvar.GetTestOrgFromEnv(t)
	service := "cloudkms.googleapis.com"
	members := []string{"user:gterraformtest1@gmail.com"}
	members2 := []string{}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Apply IAM audit config with member
			{
				Config: testAccOrganizationAssociateAuditConfigMembers(org, service, members),
			},
			organizationIamAuditConfigImportStep("google_organization_iam_audit_config.acceptance", org, service),

			// Apply IAM audit config with no members
			{
				Config: testAccOrganizationAssociateAuditConfigMembers(org, service, members2),
			},
			organizationIamAuditConfigImportStep("google_organization_iam_audit_config.acceptance", org, service),
		},
	})
}

// test changing service with no exempt members
func TestAccOrganizationIamAuditConfig_updateNoExemptMembers(t *testing.T) {
	if os.Getenv(runOrgIamAuditConfigTestEnvVar) != "true" {
		t.Skipf("Environment variable %s is not set, skipping.", runOrgIamAuditConfigTestEnvVar)
	}
	org := envvar.GetTestOrgFromEnv(t)
	logType := "DATA_READ"
	logType2 := "DATA_WRITE"
	service := "cloudkms.googleapis.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Apply IAM audit config with DATA_READ
			{
				Config: testAccOrganizationAssociateAuditConfigLogType(org, service, logType),
			},
			organizationIamAuditConfigImportStep("google_organization_iam_audit_config.acceptance", org, service),

			// Apply IAM audit config with DATA_WRITe
			{
				Config: testAccOrganizationAssociateAuditConfigLogType(org, service, logType2),
			},
			organizationIamAuditConfigImportStep("google_organization_iam_audit_config.acceptance", org, service),
		},
	})
}

func testAccOrganizationAssociateAuditConfigBasic(org, service string) string {
	return fmt.Sprintf(`
resource "google_organization_iam_audit_config" "acceptance" {
  org_id = "%s"
  service = "%s"
  audit_log_config {
    log_type = "DATA_READ"
    exempted_members = [
      "user:gterraformtest1@gmail.com",
      "user:gterraformtest2@gmail.com",
    ]
  }
}
`, org, service)
}

func testAccOrganizationAssociateAuditConfigMultiple(org, service, service2 string) string {
	return fmt.Sprintf(`
resource "google_organization_iam_audit_config" "acceptance" {
  org_id = "%s"
  service = "%s"
  audit_log_config {
    log_type = "DATA_READ"
    exempted_members = [
      "user:gterraformtest1@gmail.com",
      "user:gterraformtest2@gmail.com",
    ]
  }
}

resource "google_organization_iam_audit_config" "multiple" {
  org_id = "%s"
  service = "%s"
  audit_log_config {
    log_type = "DATA_WRITE"
  }
}
`, org, service, org, service2)
}

func testAccOrganizationAssociateAuditConfigUpdated(org, service string) string {
	return fmt.Sprintf(`
resource "google_organization_iam_audit_config" "acceptance" {
  org_id = "%s"
  service = "%s"
  audit_log_config {
    log_type = "DATA_WRITE"
    exempted_members = [
      "user:admin@hashicorptest.com",
      "user:gterraformtest2@gmail.com",
    ]
  }
}
`, org, service)
}

func testAccOrganizationAssociateAuditConfigDropMemberFromBasic(org, service string) string {
	return fmt.Sprintf(`
resource "google_organization_iam_audit_config" "acceptance" {
  org_id = "%s"
  service = "%s"
  audit_log_config {
    log_type = "DATA_READ"
    exempted_members = [
      "user:gterraformtest1@gmail.com",
    ]
  }
}
`, org, service)
}

func testAccOrganizationAssociateAuditConfigMembers(org, service string, members []string) string {
	var memberStr string
	if len(members) > 0 {
		for pos, member := range members {
			members[pos] = "\"" + member + "\","
		}
		memberStr = "\n    exempted_members = [" + strings.Join(members, "\n") + "\n    ]"
	}
	return fmt.Sprintf(`
resource "google_organization_iam_audit_config" "acceptance" {
  org_id = "%s"
  service = "%s"
  audit_log_config {
    log_type = "DATA_READ"%s
  }
}
`, org, service, memberStr)
}

func testAccOrganizationAssociateAuditConfigLogType(org, service, logType string) string {
	return fmt.Sprintf(`
resource "google_organization_iam_audit_config" "acceptance" {
  org_id = "%s"
  service = "%s"
  audit_log_config {
    log_type = "%s"
  }
}
`, org, service, logType)
}
