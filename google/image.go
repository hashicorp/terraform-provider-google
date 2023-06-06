// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"github.com/hashicorp/terraform-provider-google/google/services/compute"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// built-in projects to look for images/families containing the string
// on the left in
var imageMap = compute.ImageMap

// If the given name is a URL, return it.
// If it's in the form projects/{project}/global/images/{image}, return it
// If it's in the form projects/{project}/global/images/family/{family}, return it
// If it's in the form global/images/{image}, return it
// If it's in the form global/images/family/{family}, return it
// If it's in the form family/{family}, check if it's a family in the current project. If it is, return it as global/images/family/{family}.
//
//	If not, check if it could be a GCP-provided family, and if it exists. If it does, return it as projects/{project}/global/images/family/{family}.
//
// If it's in the form {project}/{family-or-image}, check if it's an image in the named project. If it is, return it as projects/{project}/global/images/{image}.
//
//	If not, check if it's a family in the named project. If it is, return it as projects/{project}/global/images/family/{family}.
//
// If it's in the form {family-or-image}, check if it's an image in the current project. If it is, return it as global/images/{image}.
//
//	If not, check if it could be a GCP-provided image, and if it exists. If it does, return it as projects/{project}/global/images/{image}.
//	If not, check if it's a family in the current project. If it is, return it as global/images/family/{family}.
//	If not, check if it could be a GCP-provided family, and if it exists. If it does, return it as projects/{project}/global/images/family/{family}
//
// Deprecated: For backward compatibility resolveImage is still working,
// but all new code should use ResolveImage in the compute package instead.
func resolveImage(c *transport_tpg.Config, project, name, userAgent string) (string, error) {
	return compute.ResolveImage(c, project, name, userAgent)
}
