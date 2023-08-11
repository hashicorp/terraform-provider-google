// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Test that an IAM audit config can be applied to a folder
func TestAccFolderIamAuditConfig_basic(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	fname := "tf-test-" + acctest.RandString(t, 10)
	service := "cloudkms.googleapis.com"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(t, org, fname),
				),
			},
			// Apply an IAM audit config
			{
				Config: testAccFolderAssociateAuditConfigBasic(org, fname, service),
			},
		},
	})
}

// Test that multiple IAM audit configs can be applied to a folder, one at a time
func TestAccFolderIamAuditConfig_multiple(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	fname := "tf-test-" + acctest.RandString(t, 10)
	service := "cloudkms.googleapis.com"
	service2 := "cloudsql.googleapis.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(t, org, fname),
				),
			},
			// Apply an IAM audit config
			{
				Config: testAccFolderAssociateAuditConfigBasic(org, fname, service),
			},
			// Apply another IAM audit config
			{
				Config: testAccFolderAssociateAuditConfigMultiple(org, fname, service, service2),
			},
		},
	})
}

// Test that multiple IAM audit configs can be applied to a folder all at once
func TestAccFolderIamAuditConfig_multipleAtOnce(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	fname := "tf-test-" + acctest.RandString(t, 10)
	service := "cloudkms.googleapis.com"
	service2 := "cloudsql.googleapis.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(t, org, fname),
				),
			},
			// Apply an IAM audit config
			{
				Config: testAccFolderAssociateAuditConfigMultiple(org, fname, service, service2),
			},
		},
	})
}

// Test that an IAM audit config can be updated once applied to a folder
func TestAccFolderIamAuditConfig_update(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	fname := "tf-test-" + acctest.RandString(t, 10)
	service := "cloudkms.googleapis.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(t, org, fname),
				),
			},
			// Apply an IAM audit config
			{
				Config: testAccFolderAssociateAuditConfigBasic(org, fname, service),
			},
			// Apply an updated IAM audit config
			{
				Config: testAccFolderAssociateAuditConfigUpdated(org, fname, service),
			},
			// Drop the original member
			{
				Config: testAccFolderAssociateAuditConfigDropMemberFromBasic(org, fname, service),
			},
		},
	})
}

// Test that an IAM audit config can be removed from a folder
func TestAccFolderIamAuditConfig_remove(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	fname := "tf-test-" + acctest.RandString(t, 10)
	service := "cloudkms.googleapis.com"
	service2 := "cloudsql.googleapis.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(t, org, fname),
				),
			},
			// Apply multiple IAM audit configs
			{
				Config: testAccFolderAssociateAuditConfigMultiple(org, fname, service, service2),
			},
			// Remove the audit configs
			{
				Config: testAccFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(t, org, fname),
				),
			},
		},
	})
}

// Test adding exempt first exempt member
func TestAccFolderIamAuditConfig_addFirstExemptMember(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	fname := "tf-test-" + acctest.RandString(t, 10)
	service := "cloudkms.googleapis.com"
	members := []string{}
	members2 := []string{"user:gterraformtest1@gmail.com"}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(t, org, fname),
				),
			},
			// Apply IAM audit config with no members
			{
				Config: testAccFolderAssociateAuditConfigMembers(org, fname, service, members),
			},
			// Apply IAM audit config with one member
			{
				Config: testAccFolderAssociateAuditConfigMembers(org, fname, service, members2),
			},
		},
	})
}

// test removing last exempt member
func TestAccFolderIamAuditConfig_removeLastExemptMember(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	fname := "tf-test-" + acctest.RandString(t, 10)
	service := "cloudkms.googleapis.com"
	members2 := []string{}
	members := []string{"user:gterraformtest1@gmail.com"}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(t, org, fname),
				),
			},
			// Apply IAM audit config with member
			{
				Config: testAccFolderAssociateAuditConfigMembers(org, fname, service, members),
			},
			// Apply IAM audit config with no members
			{
				Config: testAccFolderAssociateAuditConfigMembers(org, fname, service, members2),
			},
		},
	})
}

// test changing log type with no exempt members
func TestAccFolderIamAuditConfig_updateNoExemptMembers(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	fname := "tf-test-" + acctest.RandString(t, 10)
	logType := "DATA_READ"
	logType2 := "DATA_WRITE"
	service := "cloudkms.googleapis.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolderIamBasic(org, fname),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(t, org, fname),
				),
			},
			// Apply IAM audit config with DATA_READ
			{
				Config: testAccFolderAssociateAuditConfigLogType(org, fname, service, logType),
			},
			// Apply IAM audit config with DATA_WRITE
			{
				Config: testAccFolderAssociateAuditConfigLogType(org, fname, service, logType2),
			},
		},
	})
}

func testAccFolderAssociateAuditConfigBasic(org, fname, service string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  parent       = "organizations/%s"
  display_name = "%s"
}

resource "google_folder_iam_audit_config" "acceptance" {
  folder = google_folder.acceptance.name
  service = "%s"
  audit_log_config {
    log_type = "DATA_READ"
    exempted_members = [
      "user:gterraformtest1@gmail.com",
      "user:gterraformtest2@gmail.com",
    ]
  }
}
`, org, fname, service)
}

func testAccFolderAssociateAuditConfigMultiple(org, fname, service, service2 string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  parent       = "organizations/%s"
  display_name = "%s"
}

resource "google_folder_iam_audit_config" "acceptance" {
  folder = google_folder.acceptance.name
  service = "%s"
  audit_log_config {
    log_type = "DATA_READ"
    exempted_members = [
      "user:gterraformtest1@gmail.com",
      "user:gterraformtest2@gmail.com",
    ]
  }
}

resource "google_folder_iam_audit_config" "multiple" {
  folder = google_folder.acceptance.name
  service = "%s"
  audit_log_config {
    log_type = "DATA_WRITE"
  }
}
`, org, fname, service, service2)
}

func testAccFolderAssociateAuditConfigUpdated(org, fname, service string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  parent       = "organizations/%s"
  display_name = "%s"
}

resource "google_folder_iam_audit_config" "acceptance" {
  folder = google_folder.acceptance.name
  service = "%s"
  audit_log_config {
    log_type = "DATA_WRITE"
    exempted_members = [
      "user:admin@hashicorptest.com",
      "user:gterraformtest2@gmail.com",
    ]
  }
}
`, org, fname, service)
}

func testAccFolderAssociateAuditConfigDropMemberFromBasic(org, fname, service string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  parent       = "organizations/%s"
  display_name = "%s"
}

resource "google_folder_iam_audit_config" "acceptance" {
  folder = google_folder.acceptance.name
  service = "%s"
  audit_log_config {
    log_type = "DATA_READ"
    exempted_members = [
      "user:gterraformtest1@gmail.com",
    ]
  }
}
`, org, fname, service)
}

func testAccFolderAssociateAuditConfigMembers(org, fname, service string, members []string) string {
	var memberStr string
	if len(members) > 0 {
		for pos, member := range members {
			members[pos] = "\"" + member + "\","
		}
		memberStr = "\n    exempted_members = [" + strings.Join(members, "\n") + "\n    ]"
	}
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  parent       = "organizations/%s"
  display_name = "%s"
}

resource "google_folder_iam_audit_config" "acceptance" {
  folder = google_folder.acceptance.name
  service = "%s"
  audit_log_config {
    log_type = "DATA_READ"%s
  }
}
`, org, fname, service, memberStr)
}

func testAccFolderAssociateAuditConfigLogType(org, fname, service, logType string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  parent       = "organizations/%s"
  display_name = "%s"
}

resource "google_folder_iam_audit_config" "acceptance" {
  folder = google_folder.acceptance.name
  service = "%s"
  audit_log_config {
    log_type = "%s"
  }
}
`, org, fname, service, logType)
}
