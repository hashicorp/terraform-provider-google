// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package eventarc

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func expandToLongForm(pattern string, v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if strings.HasPrefix(v.(string), "projects/") || v.(string) == "" {
		// If empty or the long-form input is provided, send it as-is.
		return v, nil
	}

	// Otherwise, extract the project, and accept long-form inputs.
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(project, "/")
	project = parts[len(parts)-1]

	return fmt.Sprintf(pattern, project, v.(string)), nil
}

func expandToRegionalLongForm(pattern string, v interface{}, d tpgresource.TerraformResourceData, config *transport_tpg.Config) (interface{}, error) {
	if strings.HasPrefix(v.(string), "projects/") || v.(string) == "" {
		// If empty or the long-form input is provided, send it as-is.
		return v, nil
	}

	// Otherwise, extract the project, and accept long-form inputs.
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return nil, err
	}
	parts := strings.Split(project, "/")
	project = parts[len(parts)-1]

	// Extract the location, and accept long-form inputs.
	location, err := tpgresource.GetLocation(d, config)
	if err != nil {
		return nil, err
	}
	parts = strings.Split(location, "/")
	location = parts[len(parts)-1]

	return fmt.Sprintf(pattern, project, location, v.(string)), nil
}
