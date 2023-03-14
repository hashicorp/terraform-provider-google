package google

import (
	"fmt"
	"log"
	"testing"

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
	project, err := client.Projects.Get(GetTestProjectFromEnv()).Do()
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

	mergedBindings := MergeBindings(append(policy.Bindings, newBindings...))

	if !compareBindings(policy.Bindings, mergedBindings) {
		addedBindings := missingBindings(policy.Bindings, mergedBindings)
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

// Returns a map representing iam bindings that are in the first map but not the second.
func missingBindingsMap(aMap, bMap map[iamBindingKey]map[string]struct{}) map[iamBindingKey]map[string]struct{} {
	results := make(map[iamBindingKey]map[string]struct{})
	for key, aMembers := range aMap {
		if bMembers, ok := bMap[key]; ok {
			// The key is in both maps.
			resultMembers := make(map[string]struct{})

			for aMember := range aMembers {
				if _, ok := bMembers[aMember]; !ok {
					// The member is in a but not in b.
					resultMembers[aMember] = struct{}{}
				}
			}
			for bMember := range bMembers {
				if _, ok := aMembers[bMember]; !ok {
					// The member is in b but not in a.
					resultMembers[bMember] = struct{}{}
				}
			}

			if len(resultMembers) > 0 {
				results[key] = resultMembers
			}
		} else {
			// The key is in map a but not map b.
			results[key] = aMembers
		}
	}

	for key, bMembers := range bMap {
		if _, ok := aMap[key]; !ok {
			// The key is in map b but not map a.
			results[key] = bMembers
		}
	}

	return results
}

// Returns the bindings that are in the first set of bindings but not the second.
func missingBindings(a, b []*cloudresourcemanager.Binding) []*cloudresourcemanager.Binding {
	aMap := createIamBindingsMap(a)
	bMap := createIamBindingsMap(b)

	var results []*cloudresourcemanager.Binding
	for key, membersSet := range missingBindingsMap(aMap, bMap) {
		members := make([]string, 0, len(membersSet))
		for member := range membersSet {
			members = append(members, member)
		}
		results = append(results, &cloudresourcemanager.Binding{
			Role:    key.Role,
			Members: members,
		})
	}
	return results
}
