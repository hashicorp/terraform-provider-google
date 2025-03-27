// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package acctest

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

type IamMember struct {
	Member, Role string
}

// BootstrapIamMembers ensures that a given set of member/role pairs exist in the default
// test project. This should be used to avoid race conditions that can happen on the
// default project due to parallel tests managing the same member/role pairings. Members
// will have `{project_number}` replaced with the default test project's project number.
func BootstrapIamMembers(t *testing.T, members []IamMember) {
	config := BootstrapConfig(t)
	if config == nil {
		t.Fatal("Could not bootstrap a config for BootstrapIamMembers.")
	}
	client := config.NewResourceManagerClient(config.UserAgent)

	// Separate the given members into two groups: project-level vs. org-level
	var projectMembers []IamMember
	var orgMembers []IamMember
	for _, member := range members {
		// If the member has an {organization_id} token, we'll handle it as an org binding
		if strings.Contains(member.Member, "{organization_id}") {
			orgMembers = append(orgMembers, member)
		} else {
			// Otherwise, treat as project-level (this also covers {project_number} or none)
			projectMembers = append(projectMembers, member)
		}
	}

	if len(projectMembers) > 0 {
		// Get the project since we need its number, id, and policy.
		project, err := client.Projects.Get(envvar.GetTestProjectFromEnv()).Do()
		if err != nil {
			t.Fatalf("Error getting project with id %q: %s", project.ProjectId, err)
		}

		var projectBindings []*cloudresourcemanager.Binding
		for _, pm := range projectMembers {
			replacedMember := strings.ReplaceAll(pm.Member, "{project_number}", strconv.FormatInt(project.ProjectNumber, 10))
			projectBindings = append(projectBindings, &cloudresourcemanager.Binding{
				Role:    pm.Role,
				Members: []string{replacedMember},
			})
		}
		applyProjectIamBindings(t, client, project.ProjectId, projectBindings)
	}

	if len(orgMembers) > 0 {
		// Get the organization ID from environment if any
		orgId := envvar.GetTestOrgTargetFromEnv(t)
		if orgId == "" {
			t.Fatal("Error: Org-level IAM was requested, but no target organization ID was set in the environment.")
		}

		var orgBindings []*cloudresourcemanager.Binding
		for _, om := range orgMembers {
			replacedMember := strings.ReplaceAll(om.Member, "{organization_id}", orgId)
			orgBindings = append(orgBindings, &cloudresourcemanager.Binding{
				Role:    om.Role,
				Members: []string{replacedMember},
			})
		}
		orgName := "organizations/" + orgId
		applyOrgIamBindings(t, client, orgName, orgBindings)
	}
}

func applyProjectIamBindings(t *testing.T,
	client *cloudresourcemanager.Service,
	projectId string,
	newBindings []*cloudresourcemanager.Binding) {

	// Retry bootstrapping with exponential backoff for concurrent writes
	backoff := time.Second
	for {
		getPolicyRequest := &cloudresourcemanager.GetIamPolicyRequest{}
		policy, err := client.Projects.GetIamPolicy(projectId, getPolicyRequest).Do()
		if transport_tpg.IsGoogleApiErrorWithCode(err, 429) {
			t.Logf("[DEBUG] 429 while attempting to read policy for project %s, waiting %v before attempting again", projectId, backoff)
			time.Sleep(backoff)
			continue
		} else if err != nil {
			t.Fatalf("Error getting iam policy for project %s: %v\n", projectId, err)
		}

		mergedBindings := tpgiamresource.MergeBindings(append(policy.Bindings, newBindings...))

		if tpgiamresource.CompareBindings(policy.Bindings, mergedBindings) {
			t.Logf("[DEBUG] All bindings already present for project %s", projectId)
			break
		}
		// The policy must change.
		policy.Bindings = mergedBindings
		setPolicyRequest := &cloudresourcemanager.SetIamPolicyRequest{Policy: policy}
		policy, err = client.Projects.SetIamPolicy(projectId, setPolicyRequest).Do()
		if err == nil {
			t.Logf("[DEBUG] Waiting for IAM bootstrapping to propagate for project %s.", projectId)
			time.Sleep(3 * time.Minute)
			break
		}
		if tpgresource.IsConflictError(err) {
			t.Logf("[DEBUG]: Concurrent policy changes, restarting read-modify-write after %s", backoff)
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > 30*time.Second {
				t.Fatalf("Error applying IAM policy to %s: Too many conflicts.  Latest error: %s", projectId, err)
			}
			continue
		}
		t.Fatalf("Error setting project iam policy: %v", err)
	}
}

func applyOrgIamBindings(
	t *testing.T,
	client *cloudresourcemanager.Service,
	orgName string,
	newBindings []*cloudresourcemanager.Binding) {

	// Retry bootstrapping with exponential backoff for concurrent writes
	backoff := time.Second
	for {
		getPolicyRequest := &cloudresourcemanager.GetIamPolicyRequest{}
		policy, err := client.Organizations.GetIamPolicy(orgName, getPolicyRequest).Do()
		if transport_tpg.IsGoogleApiErrorWithCode(err, 429) {
			t.Logf("[DEBUG] 429 while attempting to read policy for org %s, waiting %v before attempting again", orgName, backoff)
			time.Sleep(backoff)
			continue
		} else if err != nil {
			t.Fatalf("Error getting iam policy for org %s: %v\n", orgName, err)
		}

		mergedBindings := tpgiamresource.MergeBindings(append(policy.Bindings, newBindings...))

		if tpgiamresource.CompareBindings(policy.Bindings, mergedBindings) {
			t.Logf("[DEBUG] All bindings already present for org %s", orgName)
			break
		}
		// The policy must change.
		policy.Bindings = mergedBindings
		setPolicyRequest := &cloudresourcemanager.SetIamPolicyRequest{Policy: policy}
		policy, err = client.Organizations.SetIamPolicy(orgName, setPolicyRequest).Do()
		if err == nil {
			t.Logf("[DEBUG] Waiting for IAM bootstrapping to propagate for org %s.", orgName)
			time.Sleep(3 * time.Minute)
			break
		}
		if tpgresource.IsConflictError(err) {
			t.Logf("[DEBUG]: Concurrent policy changes, restarting read-modify-write after %s", backoff)
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > 30*time.Second {
				t.Fatalf("Error applying IAM policy to %s: Too many conflicts.  Latest error: %s", orgName, err)
			}
			continue
		}
		t.Fatalf("Error setting org iam policy: %v", err)
	}
}
