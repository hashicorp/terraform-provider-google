// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
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
	return acctest.BootstrapAllPSARoles(t, prefix, agentNames, roles)
}

// BootstrapAllPSARole is a version of BootstrapAllPSARoles for granting a
// single role to multiple service agents.
func BootstrapAllPSARole(t *testing.T, prefix string, agentNames []string, role string) bool {
	return acctest.BootstrapAllPSARole(t, prefix, agentNames, role)
}

// BootstrapPSARoles is a version of BootstrapAllPSARoles for granting roles to
// a single service agent.
func BootstrapPSARoles(t *testing.T, prefix, agentName string, roles []string) bool {
	return acctest.BootstrapPSARoles(t, prefix, agentName, roles)
}

// BootstrapPSARole is a simplified version of BootstrapPSARoles for granting a
// single role to a single service agent.
func BootstrapPSARole(t *testing.T, prefix, agentName, role string) bool {
	return acctest.BootstrapPSARole(t, prefix, agentName, role)
}

// Returns the bindings that are in the first set of bindings but not the second.
//
// Deprecated: For backward compatibility missingBindings is still working,
// but all new code should use MissingBindings in the tpgiamresource package instead.
func missingBindings(a, b []*cloudresourcemanager.Binding) []*cloudresourcemanager.Binding {
	return tpgiamresource.MissingBindings(a, b)
}
