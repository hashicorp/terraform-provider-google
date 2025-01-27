// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package acctest

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
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

	getPolicyRequest := &cloudresourcemanager.GetIamPolicyRequest{}
	policy, err := client.Projects.GetIamPolicy(project.ProjectId, getPolicyRequest).Do()
	if err != nil {
		t.Fatalf("Error getting project iam policy: %v", err)
	}

	// Create the bindings we need to add to the policy.
	var newBindings []*cloudresourcemanager.Binding
	for _, member := range members {
		newBindings = append(newBindings, &cloudresourcemanager.Binding{
			Role:    member.Role,
			Members: []string{strings.ReplaceAll(member.Member, "{project_number}", strconv.FormatInt(project.ProjectNumber, 10))},
		})
	}

	mergedBindings := tpgiamresource.MergeBindings(append(policy.Bindings, newBindings...))

	if !tpgiamresource.CompareBindings(policy.Bindings, mergedBindings) {
		addedBindings := tpgiamresource.MissingBindings(policy.Bindings, mergedBindings)
		for _, missingBinding := range addedBindings {
			log.Printf("[DEBUG] Adding binding: %+v", missingBinding)
		}
		// The policy must change.
		policy.Bindings = mergedBindings
		setPolicyRequest := &cloudresourcemanager.SetIamPolicyRequest{Policy: policy}
		policy, err = client.Projects.SetIamPolicy(project.ProjectId, setPolicyRequest).Do()
		if err != nil {
			t.Fatalf("Error setting project iam policy: %v", err)
		}
		msg := "Added the following bindings to the test project's IAM policy:\n"
		for _, binding := range addedBindings {
			msg += fmt.Sprintf("Members: %q, Role: %q\n", binding.Members, binding.Role)
		}
		msg += "Waiting for IAM to propagate."
		t.Log(msg)
		time.Sleep(3 * time.Minute)
	}
}
