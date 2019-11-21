package google

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func projectIamAuditConfigImportStep(resourceName, pid, service string) resource.TestStep {
	return resource.TestStep{
		ResourceName:      resourceName,
		ImportStateId:     fmt.Sprintf("%s %s", pid, service),
		ImportState:       true,
		ImportStateVerify: true,
	}
}

// Test that an IAM audit config can be applied to a project
func TestAccProjectIamAuditConfig_basic(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	service := "cloudkms.googleapis.com"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(pid),
				),
			},
			// Apply an IAM audit config
			{
				Config: testAccProjectAssociateAuditConfigBasic(pid, pname, org, service),
			},
			projectIamAuditConfigImportStep("google_project_iam_audit_config.acceptance", pid, service),
		},
	})
}

// Test that multiple IAM audit configs can be applied to a project, one at a time
func TestAccProjectIamAuditConfig_multiple(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	service := "cloudkms.googleapis.com"
	service2 := "cloudsql.googleapis.com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(pid),
				),
			},
			// Apply an IAM audit config
			{
				Config: testAccProjectAssociateAuditConfigBasic(pid, pname, org, service),
			},
			// Apply another IAM audit config
			{
				Config: testAccProjectAssociateAuditConfigMultiple(pid, pname, org, service, service2),
			},
			projectIamAuditConfigImportStep("google_project_iam_audit_config.acceptance", pid, service),
			projectIamAuditConfigImportStep("google_project_iam_audit_config.multiple", pid, service2),
		},
	})
}

// Test that multiple IAM audit configs can be applied to a project all at once
func TestAccProjectIamAuditConfig_multipleAtOnce(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	service := "cloudkms.googleapis.com"
	service2 := "cloudsql.googleapis.com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(pid),
				),
			},
			// Apply an IAM audit config
			{
				Config: testAccProjectAssociateAuditConfigMultiple(pid, pname, org, service, service2),
			},
			projectIamAuditConfigImportStep("google_project_iam_audit_config.acceptance", pid, service),
			projectIamAuditConfigImportStep("google_project_iam_audit_config.multiple", pid, service2),
		},
	})
}

// Test that an IAM audit config can be updated once applied to a project
func TestAccProjectIamAuditConfig_update(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	service := "cloudkms.googleapis.com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(pid),
				),
			},
			// Apply an IAM audit config
			{
				Config: testAccProjectAssociateAuditConfigBasic(pid, pname, org, service),
			},
			projectIamAuditConfigImportStep("google_project_iam_audit_config.acceptance", pid, service),

			// Apply an updated IAM audit config
			{
				Config: testAccProjectAssociateAuditConfigUpdated(pid, pname, org, service),
			},
			projectIamAuditConfigImportStep("google_project_iam_audit_config.acceptance", pid, service),

			// Drop the original member
			{
				Config: testAccProjectAssociateAuditConfigDropMemberFromBasic(pid, pname, org, service),
			},
			projectIamAuditConfigImportStep("google_project_iam_audit_config.acceptance", pid, service),
		},
	})
}

// Test that an IAM audit config can be removed from a project
func TestAccProjectIamAuditConfig_remove(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	service := "cloudkms.googleapis.com"
	service2 := "cloudsql.googleapis.com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(pid),
				),
			},
			// Apply multiple IAM audit configs
			{
				Config: testAccProjectAssociateAuditConfigMultiple(pid, pname, org, service, service2),
			},
			projectIamAuditConfigImportStep("google_project_iam_audit_config.acceptance", pid, service),
			projectIamAuditConfigImportStep("google_project_iam_audit_config.multiple", pid, service2),

			// Remove the audit configs
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(pid),
				),
			},
		},
	})
}

// Test adding exempt first exempt member
func TestAccProjectIamAuditConfig_addFirstExemptMember(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	service := "cloudkms.googleapis.com"
	members := []string{}
	members2 := []string{"user:paddy@hashicorp.com"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(pid),
				),
			},
			// Apply IAM audit config with no members
			{
				Config: testAccProjectAssociateAuditConfigMembers(pid, pname, org, service, members),
			},
			projectIamAuditConfigImportStep("google_project_iam_audit_config.acceptance", pid, service),

			// Apply IAM audit config with one member
			{
				Config: testAccProjectAssociateAuditConfigMembers(pid, pname, org, service, members2),
			},
			projectIamAuditConfigImportStep("google_project_iam_audit_config.acceptance", pid, service),
		},
	})
}

// test removing last exempt member
func TestAccProjectIamAuditConfig_removeLastExemptMember(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	service := "cloudkms.googleapis.com"
	members2 := []string{}
	members := []string{"user:paddy@hashicorp.com"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(pid),
				),
			},
			// Apply IAM audit config with member
			{
				Config: testAccProjectAssociateAuditConfigMembers(pid, pname, org, service, members),
			},
			projectIamAuditConfigImportStep("google_project_iam_audit_config.acceptance", pid, service),

			// Apply IAM audit config with no members
			{
				Config: testAccProjectAssociateAuditConfigMembers(pid, pname, org, service, members2),
			},
			projectIamAuditConfigImportStep("google_project_iam_audit_config.acceptance", pid, service),
		},
	})
}

// test changing service with no exempt members
func TestAccProjectIamAuditConfig_updateNoExemptMembers(t *testing.T) {
	t.Parallel()

	org := getTestOrgFromEnv(t)
	pid := "terraform-" + acctest.RandString(10)
	logType := "DATA_READ"
	logType2 := "DATA_WRITE"
	service := "cloudkms.googleapis.com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, pname, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(pid),
				),
			},
			// Apply IAM audit config with DATA_READ
			{
				Config: testAccProjectAssociateAuditConfigLogType(pid, pname, org, service, logType),
			},
			projectIamAuditConfigImportStep("google_project_iam_audit_config.acceptance", pid, service),

			// Apply IAM audit config with DATA_WRITe
			{
				Config: testAccProjectAssociateAuditConfigLogType(pid, pname, org, service, logType2),
			},
			projectIamAuditConfigImportStep("google_project_iam_audit_config.acceptance", pid, service),
		},
	})
}

func testAccProjectAssociateAuditConfigBasic(pid, name, org, service string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_audit_config" "acceptance" {
  project = google_project.acceptance.project_id
  service = "%s"
  audit_log_config {
    log_type = "DATA_READ"
    exempted_members = [
      "user:paddy@hashicorp.com",
      "user:paddy@carvers.co",
    ]
  }
}
`, pid, name, org, service)
}

func testAccProjectAssociateAuditConfigMultiple(pid, name, org, service, service2 string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_audit_config" "acceptance" {
  project = google_project.acceptance.project_id
  service = "%s"
  audit_log_config {
    log_type = "DATA_READ"
    exempted_members = [
      "user:paddy@hashicorp.com",
      "user:paddy@carvers.co",
    ]
  }
}

resource "google_project_iam_audit_config" "multiple" {
  project = google_project.acceptance.project_id
  service = "%s"
  audit_log_config {
    log_type = "DATA_WRITE"
  }
}
`, pid, name, org, service, service2)
}

func testAccProjectAssociateAuditConfigUpdated(pid, name, org, service string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_audit_config" "acceptance" {
  project = google_project.acceptance.project_id
  service = "%s"
  audit_log_config {
    log_type = "DATA_WRITE"
    exempted_members = [
      "user:admin@hashicorptest.com",
      "user:paddy@carvers.co",
    ]
  }
}
`, pid, name, org, service)
}

func testAccProjectAssociateAuditConfigDropMemberFromBasic(pid, name, org, service string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_audit_config" "acceptance" {
  project = google_project.acceptance.project_id
  service = "%s"
  audit_log_config {
    log_type = "DATA_READ"
    exempted_members = [
      "user:paddy@hashicorp.com",
    ]
  }
}
`, pid, name, org, service)
}

func testAccProjectAssociateAuditConfigMembers(pid, name, org, service string, members []string) string {
	var memberStr string
	if len(members) > 0 {
		for pos, member := range members {
			members[pos] = "\"" + member + "\","
		}
		memberStr = "\n    exempted_members = [" + strings.Join(members, "\n") + "\n    ]"
	}
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_audit_config" "acceptance" {
  project = google_project.acceptance.project_id
  service = "%s"
  audit_log_config {
    log_type = "DATA_READ"%s
  }
}
`, pid, name, org, service, memberStr)
}

func testAccProjectAssociateAuditConfigLogType(pid, name, org, service, logType string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_audit_config" "acceptance" {
  project = google_project.acceptance.project_id
  service = "%s"
  audit_log_config {
    log_type = "%s"
  }
}
`, pid, name, org, service, logType)
}
