// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"maps"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	tpgresourcemanager "github.com/hashicorp/terraform-provider-google/google/services/resourcemanager"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestAccServiceAccount_ResourceIdentity(t *testing.T) {
	t.Parallel()

	accountId := "a" + acctest.RandString(t, 10)
	project := envvar.GetTestProjectFromEnv()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_12_0),
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceAccountBasic(accountId, "identity_test", "identity_test_desc"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectIdentity(
						"google_service_account.identity_test",
						map[string]knownvalue.Check{
							"account_id": knownvalue.StringExact(accountId),
							"project":    knownvalue.StringExact(project),
						},
					),
				},
			},
			{
				ResourceName:    "google_service_account.identity_test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithResourceIdentity,
			},
		},
	})

}

// Test that a service account resource can be created, updated, and destroyed
func TestAccServiceAccount_basic(t *testing.T) {
	t.Parallel()

	accountId := "a" + acctest.RandString(t, 10)
	uniqueId := ""
	displayName := "Terraform Test"
	displayName2 := "Terraform Test Update"
	desc := "test description"
	desc2 := ""
	project := envvar.GetTestProjectFromEnv()
	expectedEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", accountId, project)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// The first step creates a basic service account
			{
				Config: testAccServiceAccountBasic(accountId, displayName, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "member", "serviceAccount:"+expectedEmail),
				),
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportStateId:     fmt.Sprintf("projects/%s/serviceAccounts/%s", project, expectedEmail),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportStateId:     fmt.Sprintf("%s/%s", project, expectedEmail),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportStateId:     expectedEmail,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// The second step updates the service account
			{
				Config: testAccServiceAccountBasic(accountId, displayName2, desc2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
					testAccStoreServiceAccountUniqueId(&uniqueId),
				),
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// The third step explicitly adds the same default project to the service account configuration
			// and ensure the service account is not recreated by comparing the value of its unique_id with the one from the previous step
			{
				Config: testAccServiceAccountWithProject(project, accountId, displayName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
					resource.TestCheckResourceAttrPtr(
						"google_service_account.acceptance", "unique_id", &uniqueId),
				),
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test the option to ignore ALREADY_EXISTS error from creating a service account.
func TestAccServiceAccount_createIgnoreAlreadyExists(t *testing.T) {
	t.Parallel()

	accountId := "a" + acctest.RandString(t, 10)
	displayName := "Terraform Test"
	desc := "test description"
	project := envvar.GetTestProjectFromEnv()
	expectedEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", accountId, project)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// The first step creates a basic service account
			{
				Config: testAccServiceAccountBasic(accountId, displayName, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "member", "serviceAccount:"+expectedEmail),
				),
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportStateId:     fmt.Sprintf("projects/%s/serviceAccounts/%s", project, expectedEmail),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// The second step creates a new resource that duplicates with the existing service account.
			{
				Config: testAccServiceAccountDuplicateIgnoreAlreadyExists(accountId, displayName, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.duplicate", "member", "serviceAccount:"+expectedEmail),
				),
			},
		},
	})
}

// Test setting create_ignore_already_exists on an existing resource
func TestAccServiceAccount_existingResourceCreateIgnoreAlreadyExists(t *testing.T) {
	t.Parallel()

	project := envvar.GetTestProjectFromEnv()
	accountId := "a" + acctest.RandString(t, 10)
	displayName := "Terraform Test"
	desc := "test description"

	expectedEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", accountId, project)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// The first step creates a new resource with create_ignore_already_exists=false
			{
				Config: testAccServiceAccountCreateIgnoreAlreadyExists(accountId, displayName, desc, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "member", "serviceAccount:"+expectedEmail),
				),
			},
			{
				ResourceName:            "google_service_account.acceptance",
				ImportStateId:           fmt.Sprintf("projects/%s/serviceAccounts/%s", project, expectedEmail),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"create_ignore_already_exists"}, // Import leaves this field out when false
			},
			// The second step updates the resource to have create_ignore_already_exists=true
			{
				Config: testAccServiceAccountCreateIgnoreAlreadyExists(accountId, displayName, desc, true),
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(
					"google_service_account.acceptance", "project", project),
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "member", "serviceAccount:"+expectedEmail),
				),
			},
		},
	})
}

func TestAccServiceAccount_Disabled(t *testing.T) {
	t.Parallel()

	accountId := "a" + acctest.RandString(t, 10)
	uniqueId := ""
	displayName := "Terraform Test"
	desc := "test description"
	project := envvar.GetTestProjectFromEnv()
	expectedEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", accountId, project)
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// The first step creates a basic service account
			{
				Config: testAccServiceAccountBasic(accountId, displayName, desc),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "member", "serviceAccount:"+expectedEmail),
				),
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportStateId:     fmt.Sprintf("projects/%s/serviceAccounts/%s", project, expectedEmail),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// The second step disables the service account
			{
				Config: testAccServiceAccountDisabled(accountId, displayName, desc, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
					testAccStoreServiceAccountUniqueId(&uniqueId),
				),
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// The third step enables the disabled service account
			{
				Config: testAccServiceAccountDisabled(accountId, displayName, desc, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"google_service_account.acceptance", "project", project),
					testAccStoreServiceAccountUniqueId(&uniqueId),
				),
			},
			{
				ResourceName:      "google_service_account.acceptance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccStoreServiceAccountUniqueId(uniqueId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		*uniqueId = s.RootModule().Resources["google_service_account.acceptance"].Primary.Attributes["unique_id"]
		return nil
	}
}

func testAccServiceAccountBasic(account, name, desc string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
  account_id   = "%v"
  display_name = "%v"
  description  = "%v"
}
`, account, name, desc)
}

func testAccServiceAccountCreateIgnoreAlreadyExists(account, name, desc string, ignore_already_exists bool) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
  account_id   = "%v"
  display_name = "%v"
  description  = "%v"
  create_ignore_already_exists = %t
}
`, account, name, desc, ignore_already_exists)
}

func testAccServiceAccountDuplicateIgnoreAlreadyExists(account, name, desc string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
  account_id   = "%v"
  display_name = "%v"
  description  = "%v"
}
resource "google_service_account" "duplicate" {
  account_id   = "%v"
  display_name = "%v"
  description  = "%v"
  create_ignore_already_exists = true
}
`, account, name, desc, account, name, desc)
}

func testAccServiceAccountWithProject(project, account, name string) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
  project      = "%v"
  account_id   = "%v"
  display_name = "%v"
  description  = "foo"
}
`, project, account, name)
}

func testAccServiceAccountDisabled(account, name, desc string, disabled bool) string {
	return fmt.Sprintf(`
resource "google_service_account" "acceptance" {
  account_id   = "%v"
  display_name = "%v"
  description  = "%v"
  disabled      = "%t"
}
`, account, name, desc, disabled)
}

func TestResourceServiceAccountCustomDiff(t *testing.T) {
	t.Parallel()

	accountId := "a" + acctest.RandString(t, 10)
	project := envvar.GetTestProjectFromEnv()
	if project == "" {
		project = "test-project"
	}
	expectedEmail := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", accountId, project)
	expectedMember := "serviceAccount:" + expectedEmail

	cases := []struct {
		name       string
		before     map[string]interface{}
		after      map[string]interface{}
		wantEmail  string
		wantMember string
	}{
		{
			name:   "normal (new)",
			before: map[string]interface{}{},
			after: map[string]interface{}{
				"account_id": accountId,
				"name":       "", // Empty name indicates a new resource
				"project":    project,
			},
			wantEmail:  expectedEmail,
			wantMember: expectedMember,
		},
		{
			name: "no change",
			before: map[string]interface{}{
				"account_id": accountId,
				"email":      "dontchange",
				"member":     "dontchange",
				"project":    project,
			},
			after: map[string]interface{}{
				"account_id": accountId,
				"name":       "unimportant",
				"project":    project,
			},
			wantEmail:  "",
			wantMember: "",
		},
		{
			name: "recreate (new)",
			before: map[string]interface{}{
				"account_id": "recreate-account",
				"email":      "recreate-email",
				"member":     "recreate-member",
				"project":    project,
			},
			after: map[string]interface{}{
				"account_id": accountId,
				"name":       "",
				"project":    project,
			},
			wantEmail:  expectedEmail,
			wantMember: expectedMember,
		},
		{
			name:   "missing account_id (new)",
			before: map[string]interface{}{},
			after: map[string]interface{}{
				"account_id": "",
				"name":       "",
				"project":    project,
			},
			wantEmail:  "",
			wantMember: "",
		},
		{
			name:   "missing project (new)",
			before: map[string]interface{}{},
			after: map[string]interface{}{
				"account_id": accountId,
				"name":       "",
				"project":    "",
			},
			wantEmail:  "",
			wantMember: "",
		},
	}
	for _, tc := range cases {
		result := maps.Clone(tc.after)
		if tc.wantEmail != "" || tc.wantMember != "" {
			result["email"] = tc.wantEmail
			result["member"] = tc.wantMember
		}
		t.Run(tc.name, func(t *testing.T) {
			diff := &tpgresource.ResourceDiffMock{
				Before: tc.before,
				After:  tc.after,
				Schema: tpgresourcemanager.ResourceGoogleServiceAccount().Schema,
			}
			tpgresourcemanager.ResourceServiceAccountCustomDiffFunc(diff)
			if d := cmp.Diff(result, diff.After); d != "" {
				t.Fatalf("got unexpected change: %v expected: %v", diff.After, result)
			}
		})
	}
}
