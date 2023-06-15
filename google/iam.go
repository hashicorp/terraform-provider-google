// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Utils for modifying IAM policies for resources across GCP
package google

import (
	"github.com/hashicorp/terraform-provider-google/google/tpgiamresource"
	"google.golang.org/api/cloudresourcemanager/v1"
)

// Deprecated: For backward compatibility IamPolicyVersion is still working,
// but all new code should use IamPolicyVersion in the tpgiamresource package instead.
const IamPolicyVersion = tpgiamresource.IamPolicyVersion

type (
	// The ResourceIamUpdater interface is implemented for each GCP resource supporting IAM policy.
	// Implementations should be created per resource and should keep track of the resource identifier.
	//
	// Deprecated: For backward compatibility ResourceIamUpdater is still working,
	// but all new code should use ResourceIamUpdater in the tpgiamresource package instead.
	ResourceIamUpdater = tpgiamresource.ResourceIamUpdater

	// Parser for Terraform resource identifier (d.Id) for resource whose IAM policy is being changed
	//
	// Deprecated: For backward compatibility ResourceIdParserFunc is still working,
	// but all new code should use ResourceIdParserFunc in the tpgiamresource package instead.
	ResourceIdParserFunc = tpgiamresource.ResourceIdParserFunc
)

// Flattens a list of Bindings so each role+condition has a single Binding with combined members
//
// Deprecated: For backward compatibility MergeBindings is still working,
// but all new code should use MergeBindings in the tpgiamresource package instead.
func MergeBindings(bindings []*cloudresourcemanager.Binding) []*cloudresourcemanager.Binding {
	return tpgiamresource.MergeBindings(bindings)
}

// Deprecated: For backward compatibility compareBindings is still working,
// but all new code should use CompareBindings in the tpgiamresource package instead.
func compareBindings(a, b []*cloudresourcemanager.Binding) bool {
	return tpgiamresource.CompareBindings(a, b)
}

// Deprecated: For backward compatibility compareAuditConfigs is still working,
// but all new code should use CompareAuditConfigs in the tpgiamresource package instead.
func compareAuditConfigs(a, b []*cloudresourcemanager.AuditConfig) bool {
	return tpgiamresource.CompareAuditConfigs(a, b)
}

// Util to deref and print auditConfigs
//
// Deprecated: For backward compatibility debugPrintAuditConfigs is still working,
// but all new code should use DebugPrintAuditConfigs in the tpgiamresource package instead.
func debugPrintAuditConfigs(bs []*cloudresourcemanager.AuditConfig) string {
	return tpgiamresource.DebugPrintAuditConfigs(bs)
}

// Util to deref and print bindings
//
// Deprecated: For backward compatibility debugPrintBindings is still working,
// but all new code should use DebugPrintBindings in the tpgiamresource package instead.
func debugPrintBindings(bs []*cloudresourcemanager.Binding) string {
	return tpgiamresource.DebugPrintBindings(bs)
}
