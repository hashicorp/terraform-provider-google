// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager_test

import (
	"fmt"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func projectIamBindingImportStep(resourceName, pid, role string) resource.TestStep {
	return resource.TestStep{
		ResourceName:      resourceName,
		ImportStateId:     fmt.Sprintf("%s %s", pid, role),
		ImportState:       true,
		ImportStateVerify: true,
	}
}

// Test that an IAM binding can be applied to a project
func TestAccProjectIamBinding_basic(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	pid := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	role := "roles/compute.instanceAdmin"
	member := "user:admin@hashicorptest.com"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(t, pid),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccProjectAssociateBindingBasic(pid, org, role, member),
			},
			projectIamBindingImportStep("google_project_iam_binding.acceptance", pid, role),
		},
	})
}

// Test that multiple IAM bindings can be applied to a project, one at a time
func TestAccProjectIamBinding_multiple(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	pid := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	role := "roles/compute.instanceAdmin"
	role2 := "roles/viewer"
	member := "user:admin@hashicorptest.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(t, pid),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccProjectAssociateBindingBasic(pid, org, role, member),
			},
			// Apply another IAM binding
			{
				Config: testAccProjectAssociateBindingMultiple(pid, org, role, role2),
			},
			projectIamBindingImportStep("google_project_iam_binding.acceptance", pid, role),
			projectIamBindingImportStep("google_project_iam_binding.multiple", pid, role2),
		},
	})
}

// Test that multiple IAM bindings can be applied to a project all at once
func TestAccProjectIamBinding_multipleAtOnce(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	pid := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	role := "roles/compute.instanceAdmin"
	role2 := "roles/viewer"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(t, pid),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccProjectAssociateBindingMultiple(pid, org, role, role2),
			},
			projectIamBindingImportStep("google_project_iam_binding.acceptance", pid, role),
			projectIamBindingImportStep("google_project_iam_binding.multiple", pid, role2),
		},
	})
}

// Test that an IAM binding can be updated once applied to a project
func TestAccProjectIamBinding_update(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	pid := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	role := "roles/compute.instanceAdmin"
	member := "user:admin@hashicorptest.com"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(t, pid),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccProjectAssociateBindingBasic(pid, org, role, member),
			},
			projectIamBindingImportStep("google_project_iam_binding.acceptance", pid, role),

			// Apply an updated IAM binding
			{
				Config: testAccProjectAssociateBindingUpdated(pid, org, role),
			},
			projectIamBindingImportStep("google_project_iam_binding.acceptance", pid, role),

			// Drop the original member
			{
				Config: testAccProjectAssociateBindingDropMemberFromBasic(pid, org, role),
			},
			projectIamBindingImportStep("google_project_iam_binding.acceptance", pid, role),
		},
	})
}

// Test that an IAM binding can be removed from a project
func TestAccProjectIamBinding_remove(t *testing.T) {
	// Multiple fine-grained resources
	acctest.SkipIfVcr(t)
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	pid := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	role := "roles/compute.instanceAdmin"
	role2 := "roles/viewer"

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(t, pid),
				),
			},
			// Apply multiple IAM bindings
			{
				Config: testAccProjectAssociateBindingMultiple(pid, org, role, role2),
			},
			projectIamBindingImportStep("google_project_iam_binding.acceptance", pid, role),
			projectIamBindingImportStep("google_project_iam_binding.multiple", pid, role2),

			// Remove the bindings
			{
				Config: testAccProject_create(pid, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(t, pid),
				),
			},
		},
	})
}

// Test that an IAM binding with no members can be applied to a project
func TestAccProjectIamBinding_noMembers(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	pid := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	role := "roles/compute.instanceAdmin"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(t, pid),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccProjectAssociateBindingNoMembers(pid, org, role),
			},
			projectIamBindingImportStep("google_project_iam_binding.acceptance", pid, role),
		},
	})
}

func TestAccProjectIamBinding_withCondition(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	pid := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	role := "roles/compute.instanceAdmin"
	conditionTitle := "expires_after_2019_12_31"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			// Create a new project
			{
				Config: testAccProject_create(pid, org),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectExistingPolicy(t, pid),
				),
			},
			// Apply an IAM binding
			{
				Config: testAccProjectAssociateBinding_withCondition(pid, org, role, conditionTitle),
			},
			{
				ResourceName:      "google_project_iam_binding.acceptance",
				ImportStateId:     fmt.Sprintf("%s %s %s", pid, role, conditionTitle),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Test that an IAM binding with invalid members returns an error.
func TestAccProjectIamBinding_invalidMembers(t *testing.T) {
	t.Parallel()

	org := envvar.GetTestOrgFromEnv(t)
	pid := fmt.Sprintf("tf-test-%d", acctest.RandInt(t))
	role := "roles/compute.instanceAdmin"
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config:      testAccProjectAssociateBindingBasic(pid, org, role, "admin@hashicorptest.com"),
				ExpectError: regexp.MustCompile("invalid value for members\\.0 \\(IAM members must have one of the values outlined here: https://cloud.google.com/billing/docs/reference/rest/v1/Policy#Binding\\)"),
			},
			{
				Config: testAccProjectAssociateBindingBasic(pid, org, role, "user:admin@hashicorptest.com"),
			},
		},
	})
}

func testAccProjectAssociateBindingBasic(pid, org, role, member string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_binding" "acceptance" {
  project = google_project.acceptance.project_id
  members = ["%s"]
  role    = "%s"
}
`, pid, pid, org, member, role)
}

func testAccProjectAssociateBindingMultiple(pid, org, role, role2 string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_binding" "acceptance" {
  project = google_project.acceptance.project_id
  members = ["user:admin@hashicorptest.com"]
  role    = "%s"
}

resource "google_project_iam_binding" "multiple" {
  project = google_project.acceptance.project_id
  members = ["user:gterraformtest1@gmail.com"]
  role    = "%s"
}
`, pid, pid, org, role, role2)
}

func testAccProjectAssociateBindingUpdated(pid, org, role string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_binding" "acceptance" {
  project = google_project.acceptance.project_id
  members = ["user:admin@hashicorptest.com", "user:gterraformtest1@gmail.com"]
  role    = "%s"
}
`, pid, pid, org, role)
}

func testAccProjectAssociateBindingDropMemberFromBasic(pid, org, role string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_binding" "acceptance" {
  project = google_project.acceptance.project_id
  members = ["user:gterraformtest1@gmail.com"]
  role    = "%s"
}
`, pid, pid, org, role)
}

func testAccProjectAssociateBindingNoMembers(pid, org, role string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_binding" "acceptance" {
  project = google_project.acceptance.project_id
  members = []
  role    = "%s"
}
`, pid, pid, org, role)
}

func testAccProjectAssociateBinding_withCondition(pid, org, role, conditionTitle string) string {
	return fmt.Sprintf(`
resource "google_project" "acceptance" {
  project_id = "%s"
  name       = "%s"
  org_id     = "%s"
}

resource "google_project_iam_binding" "acceptance" {
  project = google_project.acceptance.project_id
  members = ["user:admin@hashicorptest.com"]
  role    = "%s"
  condition {
    title       = "%s"
    description = "Expiring at midnight of 2019-12-31"
    expression  = "request.time < timestamp(\"2020-01-01T00:00:00Z\")"
  }
}
`, pid, pid, org, role, conditionTitle)
}
