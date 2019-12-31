package google

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func folderIamAuditConfigImportStep(resourceName, org, folderDisplayName, service string) resource.TestStep {
	return resource.TestStep{
		ResourceName: resourceName,
		ImportStateIdFunc: func(_ *terraform.State) (string, error) {
			c := testAccProvider.Meta().(*Config)
			var err error
			name, err := getFolderNameByParentAndDisplayName("organizations/"+org, folderDisplayName, c)
			if err != nil {
				return "", fmt.Errorf("Failed to retrieve IAM Policy for folder %q: %s", folderDisplayName, err)
			}
			return fmt.Sprintf("%s %s", name, service), nil
		},
		ImportState:       true,
		ImportStateVerify: true,
	}
}

// Test that an IAM audit config can be applied to a folder
func TestAccFolderIamAuditConfig_basic(t *testing.T) {
	t.Parallel()

	folderDisplayName := "tf-test-" + acctest.RandString(10)
	org := getTestOrgFromEnv(t)
	parent := "organizations/" + org
	service := "cloudkms.googleapis.com"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolder_create(folderDisplayName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(org, folderDisplayName),
				),
			},
			// Apply an IAM audit config
			{
				Config: testAccFolderAssociateAuditConfigBasic(folderDisplayName, parent, service),
			},
			folderIamAuditConfigImportStep("google_folder_iam_audit_config.acceptance", org, folderDisplayName, service),
		},
	})
}

// Test that multiple IAM audit configs can be applied to a folder, one at a time
func TestAccFolderIamAuditConfig_multiple(t *testing.T) {
	t.Parallel()

	folderDisplayName := "tf-test-" + acctest.RandString(10)
	org := getTestOrgFromEnv(t)
	parent := "organizations/" + org
	service := "cloudkms.googleapis.com"
	service2 := "cloudsql.googleapis.com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolder_create(folderDisplayName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(org, folderDisplayName),
				),
			},
			// Apply an IAM audit config
			{
				Config: testAccFolderAssociateAuditConfigBasic(folderDisplayName, parent, service),
			},
			// Apply another IAM audit config
			{
				Config: testAccFolderAssociateAuditConfigMultiple(folderDisplayName, parent, service, service2),
			},
			folderIamAuditConfigImportStep("google_folder_iam_audit_config.acceptance", org, folderDisplayName, service),
			folderIamAuditConfigImportStep("google_folder_iam_audit_config.multiple", org, folderDisplayName, service2),
		},
	})
}

// Test that multiple IAM audit configs can be applied to a folder all at once
func TestAccFolderIamAuditConfig_multipleAtOnce(t *testing.T) {
	t.Parallel()

	folderDisplayName := "tf-test-" + acctest.RandString(10)
	org := getTestOrgFromEnv(t)
	parent := "organizations/" + org
	service := "cloudkms.googleapis.com"
	service2 := "cloudsql.googleapis.com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolder_create(folderDisplayName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(org, folderDisplayName),
				),
			},
			// Apply an IAM audit config
			{
				Config: testAccFolderAssociateAuditConfigMultiple(folderDisplayName, parent, service, service2),
			},
			folderIamAuditConfigImportStep("google_folder_iam_audit_config.acceptance", org, folderDisplayName, service),
			folderIamAuditConfigImportStep("google_folder_iam_audit_config.multiple", org, folderDisplayName, service2),
		},
	})
}

// Test that an IAM audit config can be updated once applied to a folder
func TestAccFolderIamAuditConfig_update(t *testing.T) {
	t.Parallel()

	folderDisplayName := "tf-test-" + acctest.RandString(10)
	org := getTestOrgFromEnv(t)
	parent := "organizations/" + org
	service := "cloudkms.googleapis.com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolder_create(folderDisplayName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(org, folderDisplayName),
				),
			},
			// Apply an IAM audit config
			{
				Config: testAccFolderAssociateAuditConfigBasic(folderDisplayName, parent, service),
			},
			folderIamAuditConfigImportStep("google_folder_iam_audit_config.acceptance", org, folderDisplayName, service),

			// Apply an updated IAM audit config
			{
				Config: testAccFolderAssociateAuditConfigUpdated(folderDisplayName, parent, service),
			},
			folderIamAuditConfigImportStep("google_folder_iam_audit_config.acceptance", org, folderDisplayName, service),

			// Drop the original member
			{
				Config: testAccFolderAssociateAuditConfigDropMemberFromBasic(folderDisplayName, parent, service),
			},
			folderIamAuditConfigImportStep("google_folder_iam_audit_config.acceptance", org, folderDisplayName, service),
		},
	})
}

// Test that an IAM audit config can be removed from a folder
func TestAccFolderIamAuditConfig_remove(t *testing.T) {
	t.Parallel()

	folderDisplayName := "tf-test-" + acctest.RandString(10)
	org := getTestOrgFromEnv(t)
	parent := "organizations/" + org
	service := "cloudkms.googleapis.com"
	service2 := "cloudsql.googleapis.com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolder_create(folderDisplayName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(org, folderDisplayName),
				),
			},
			// Apply multiple IAM audit configs
			{
				Config: testAccFolderAssociateAuditConfigMultiple(folderDisplayName, parent, service, service2),
			},
			folderIamAuditConfigImportStep("google_folder_iam_audit_config.acceptance", org, folderDisplayName, service),
			folderIamAuditConfigImportStep("google_folder_iam_audit_config.multiple", org, folderDisplayName, service2),

			// Remove the audit configs
			{
				Config: testAccFolder_create(folderDisplayName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(org, folderDisplayName),
				),
			},
		},
	})
}

// Test adding exempt first exempt member
func TestAccFolderIamAuditConfig_addFirstExemptMember(t *testing.T) {
	t.Parallel()

	folderDisplayName := "tf-test-" + acctest.RandString(10)
	org := getTestOrgFromEnv(t)
	parent := "organizations/" + org
	service := "cloudkms.googleapis.com"
	members := []string{}
	members2 := []string{"user:paddy@hashicorp.com"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolder_create(folderDisplayName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(org, folderDisplayName),
				),
			},
			// Apply IAM audit config with no members
			{
				Config: testAccFolderAssociateAuditConfigMembers(folderDisplayName, parent, service, members),
			},
			folderIamAuditConfigImportStep("google_folder_iam_audit_config.acceptance", org, folderDisplayName, service),

			// Apply IAM audit config with one member
			{
				Config: testAccFolderAssociateAuditConfigMembers(folderDisplayName, parent, service, members2),
			},
			folderIamAuditConfigImportStep("google_folder_iam_audit_config.acceptance", org, folderDisplayName, service),
		},
	})
}

// Test removing last exempt member
func TestAccFolderIamAuditConfig_removeLastExemptMember(t *testing.T) {
	t.Parallel()

	folderDisplayName := "tf-test-" + acctest.RandString(10)
	org := getTestOrgFromEnv(t)
	parent := "organizations/" + org
	service := "cloudkms.googleapis.com"
	members2 := []string{}
	members := []string{"user:paddy@hashicorp.com"}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolder_create(folderDisplayName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(org, folderDisplayName),
				),
			},
			// Apply IAM audit config with member
			{
				Config: testAccFolderAssociateAuditConfigMembers(folderDisplayName, parent, service, members),
			},
			folderIamAuditConfigImportStep("google_folder_iam_audit_config.acceptance", org, folderDisplayName, service),

			// Apply IAM audit config with no members
			{
				Config: testAccFolderAssociateAuditConfigMembers(folderDisplayName, parent, service, members2),
			},
			folderIamAuditConfigImportStep("google_folder_iam_audit_config.acceptance", org, folderDisplayName, service),
		},
	})
}

// Test changing service with no exempt members
func TestAccFolderIamAuditConfig_updateNoExemptMembers(t *testing.T) {
	t.Parallel()

	folderDisplayName := "tf-test-" + acctest.RandString(10)
	org := getTestOrgFromEnv(t)
	parent := "organizations/" + org
	logType := "DATA_READ"
	logType2 := "DATA_WRITE"
	service := "cloudkms.googleapis.com"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create a new folder
			{
				Config: testAccFolder_create(folderDisplayName, parent),
				Check: resource.ComposeTestCheckFunc(
					testAccFolderExistingPolicy(org, folderDisplayName),
				),
			},
			// Apply IAM audit config with DATA_READ
			{
				Config: testAccFolderAssociateAuditConfigLogType(folderDisplayName, parent, service, logType),
			},
			folderIamAuditConfigImportStep("google_folder_iam_audit_config.acceptance", org, folderDisplayName, service),

			// Apply IAM audit config with DATA_WRITE
			{
				Config: testAccFolderAssociateAuditConfigLogType(folderDisplayName, parent, service, logType2),
			},
			folderIamAuditConfigImportStep("google_folder_iam_audit_config.acceptance", org, folderDisplayName, service),
		},
	})
}

func testAccFolderAssociateAuditConfigBasic(folder, parent, service string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  display_name = "%s"
  parent = "%s"
}
resource "google_folder_iam_audit_config" "acceptance" {
  folder = "${google_folder.acceptance.name}"
  service = "%s"
  audit_log_config {
    log_type = "DATA_READ"
    exempted_members = [
      "user:paddy@hashicorp.com",
      "user:paddy@carvers.co"
    ]
  }
}
`, folder, parent, service)
}

func testAccFolderAssociateAuditConfigMultiple(folder, parent, service, service2 string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  display_name = "%s"
  parent = "%s"
}
resource "google_folder_iam_audit_config" "acceptance" {
  folder = "${google_folder.acceptance.name}"
  service = "%s"
  audit_log_config {
    log_type = "DATA_READ" 
    exempted_members = [
      "user:paddy@hashicorp.com",
      "user:paddy@carvers.co"
    ]
  }
}
resource "google_folder_iam_audit_config" "multiple" {
  folder = "${google_folder.acceptance.name}"
  service = "%s"
  audit_log_config {
    log_type = "DATA_WRITE"
  }
}
`, folder, parent, service, service2)
}

func testAccFolderAssociateAuditConfigUpdated(folder, parent, service string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  display_name = "%s"
  parent = "%s"
}
resource "google_folder_iam_audit_config" "acceptance" {
  folder = "${google_folder.acceptance.name}"
  service = "%s"
  audit_log_config {
    log_type = "DATA_WRITE"
    exempted_members = [
      "user:admin@hashicorptest.com",
      "user:paddy@carvers.co"
    ]
  }
}
`, folder, parent, service)
}

func testAccFolderAssociateAuditConfigDropMemberFromBasic(folder, parent, service string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  display_name = "%s"
  parent = "%s"
}
resource "google_folder_iam_audit_config" "acceptance" {
  folder = "${google_folder.acceptance.name}"
  service = "%s"
  audit_log_config {
    log_type = "DATA_READ"
    exempted_members = [
      "user:paddy@hashicorp.com",
    ]
  }
}
`, folder, parent, service)
}

func testAccFolderAssociateAuditConfigMembers(folder, parent, service string, members []string) string {
	var memberStr string
	if len(members) > 0 {
		for pos, member := range members {
			members[pos] = "\"" + member + "\","
		}
		memberStr = "\n    exempted_members = [" + strings.Join(members, "\n") + "\n    ]"
	}
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  display_name = "%s"
  parent = "%s"
}
resource "google_folder_iam_audit_config" "acceptance" {
  folder = "${google_folder.acceptance.name}"
  service = "%s"
  audit_log_config {
    log_type = "DATA_READ"%s
  }
}
`, folder, parent, service, memberStr)
}

func testAccFolderAssociateAuditConfigLogType(folder, parent, service, logType string) string {
	return fmt.Sprintf(`
resource "google_folder" "acceptance" {
  display_name = "%s"
  parent = "%s"
}
resource "google_folder_iam_audit_config" "acceptance" {
  folder = "${google_folder.acceptance.name}"
  service = "%s"
  audit_log_config {
    log_type = "%s"
  }
}
`, folder, parent, service, logType)
}
