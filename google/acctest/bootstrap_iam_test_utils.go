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

	// Get the project since we need its number, id, and policy.
	project, err := client.Projects.Get(envvar.GetTestProjectFromEnv()).Do()
	if err != nil {
		t.Fatalf("Error getting project with id %q: %s", project.ProjectId, err)
	}

	// Create the bindings we need to add to the policy.
	var newBindings []*cloudresourcemanager.Binding
	for _, member := range members {
		newBindings = append(newBindings, &cloudresourcemanager.Binding{
			Role:    member.Role,
			Members: []string{strings.ReplaceAll(member.Member, "{project_number}", strconv.FormatInt(project.ProjectNumber, 10))},
		})
	}

	// Retry bootstrapping with exponential backoff for concurrent writes
	backoff := time.Second
	for {
		getPolicyRequest := &cloudresourcemanager.GetIamPolicyRequest{}
		policy, err := client.Projects.GetIamPolicy(project.ProjectId, getPolicyRequest).Do()
		if transport_tpg.IsGoogleApiErrorWithCode(err, 429) {
			t.Logf("[DEBUG] 429 while attempting to read policy for project %s, waiting %v before attempting again", project.ProjectId, backoff)
			time.Sleep(backoff)
			continue
		} else if err != nil {
			t.Fatalf("Error getting iam policy for project %s: %v\n", project.ProjectId, err)
		}

		mergedBindings := tpgiamresource.MergeBindings(append(policy.Bindings, newBindings...))

		if tpgiamresource.CompareBindings(policy.Bindings, mergedBindings) {
			t.Logf("[DEBUG] All bindings already present for project %s", project.ProjectId)
			break
		}
		// The policy must change.
		policy.Bindings = mergedBindings
		setPolicyRequest := &cloudresourcemanager.SetIamPolicyRequest{Policy: policy}
		policy, err = client.Projects.SetIamPolicy(project.ProjectId, setPolicyRequest).Do()
		if err == nil {
			t.Logf("[DEBUG] Waiting for IAM bootstrapping to propagate for project %s.", project.ProjectId)
			time.Sleep(3 * time.Minute)
			break
		}
		if tpgresource.IsConflictError(err) {
			t.Logf("[DEBUG]: Concurrent policy changes, restarting read-modify-write after %s", backoff)
			time.Sleep(backoff)
			backoff = backoff * 2
			if backoff > 30*time.Second {
				t.Fatalf("Error applying IAM policy to %s: Too many conflicts.  Latest error: %s", project.ProjectId, err)
			}
			continue
		}
		t.Fatalf("Error setting project iam policy: %v", err)
	}
}
