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
		t.Fatal("Could not bootstrap a config for BootstrapAllPSARoles.")
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

// BootstrapAllPSARoles ensures that the given project's IAM
// policy grants the given service agents the given roles.
// prefix is usually "service-" and indicates the service agent should have the
// given prefix before the project number.
// This is important to bootstrap because using iam policy resources means that
// deleting them removes permissions for concurrent tests.
// Return whether the bindings changed.
func BootstrapAllPSARoles(t *testing.T, prefix string, agentNames, roles []string) bool {
	var members []IamMember
	for _, agentName := range agentNames {
		member := fmt.Sprintf("serviceAccount:%s{project_number}@%s.iam.gserviceaccount.com", prefix, agentName)
		for _, role := range roles {
			members = append(members, IamMember{
				Member: member,
				Role:   role,
			})
		}
	}
	BootstrapIamMembers(t, members)
	// Always return false because we now wait for IAM propagation.
	return false
}

// BootstrapAllPSARole is a version of BootstrapAllPSARoles for granting a
// single role to multiple service agents.
func BootstrapAllPSARole(t *testing.T, prefix string, agentNames []string, role string) bool {
	return BootstrapAllPSARoles(t, prefix, agentNames, []string{role})
}

// BootstrapPSARoles is a version of BootstrapAllPSARoles for granting roles to
// a single service agent.
func BootstrapPSARoles(t *testing.T, prefix, agentName string, roles []string) bool {
	return BootstrapAllPSARoles(t, prefix, []string{agentName}, roles)
}

// BootstrapPSARole is a simplified version of BootstrapPSARoles for granting a
// single role to a single service agent.
func BootstrapPSARole(t *testing.T, prefix, agentName, role string) bool {
	return BootstrapPSARoles(t, prefix, agentName, []string{role})
}
