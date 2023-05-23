package google

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

// BootstrapAllPSARoles ensures that the given project's IAM
// policy grants the given service agents the given roles.
// prefix is usually "service-" and indicates the service agent should have the
// given prefix before the project number.
// This is important to bootstrap because using iam policy resources means that
// deleting them removes permissions for concurrent tests.
// Return whether the bindings changed.
func BootstrapAllPSARoles(t *testing.T, prefix string, agentNames, roles []string) bool {
	config := BootstrapConfig(t)
	if config == nil {
		t.Fatal("Could not bootstrap a config for BootstrapAllPSARoles.")
	}
	client := config.NewResourceManagerClient(config.UserAgent)

	// Get the project since we need its number, id, and policy.
	project, err := client.Projects.Get(acctest.GetTestProjectFromEnv()).Do()
	if err != nil {
		t.Fatalf("Error getting project with id %q: %s", project.ProjectId, err)
	}

	getPolicyRequest := &cloudresourcemanager.GetIamPolicyRequest{}
	policy, err := client.Projects.GetIamPolicy(project.ProjectId, getPolicyRequest).Do()
	if err != nil {
		t.Fatalf("Error getting project iam policy: %v", err)
	}

	members := make([]string, len(agentNames))
	for i, agentName := range agentNames {
		members[i] = fmt.Sprintf("serviceAccount:%s%d@%s.iam.gserviceaccount.com", prefix, project.ProjectNumber, agentName)
	}

	// Create the bindings we need to add to the policy.
	var newBindings []*cloudresourcemanager.Binding
	for _, role := range roles {
		newBindings = append(newBindings, &cloudresourcemanager.Binding{
			Role:    role,
			Members: members,
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
		msg += "Retry the test in a few minutes."
		t.Error(msg)
		return true
	}
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

// Returns the bindings that are in the first set of bindings but not the second.
//
// Deprecated: For backward compatibility missingBindings is still working,
// but all new code should use MissingBindings in the tpgiamresource package instead.
func missingBindings(a, b []*cloudresourcemanager.Binding) []*cloudresourcemanager.Binding {
	return tpgiamresource.MissingBindings(a, b)
}
