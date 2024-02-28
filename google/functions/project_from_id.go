// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package functions

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = ProjectFromIdFunction{}

func NewProjectFromIdFunction() function.Function {
	return &ProjectFromIdFunction{}
}

type ProjectFromIdFunction struct{}

func (f ProjectFromIdFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "project_from_id"
}

func (f ProjectFromIdFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Returns the project within a provided resource's id, resource URI, self link, or full resource name.",
		Description: "Takes a single string argument, which should be a resource's id, resource URI, self link, or full resource name. This function will either return the project name from the input string or raise an error due to no project being present in the string. The function uses the presence of \"projects/{{project}}/\" in the input string to identify the project name, e.g. when the function is passed the id \"projects/my-project/zones/us-central1-c/instances/my-instance\" as an argument it will return \"my-project\".",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "id",
				Description: "A string of a resource's id, resource URI, self link, or full resource name. For example, \"projects/my-project/zones/us-central1-c/instances/my-instance\", \"https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-c/instances/my-instance\" and \"//gkehub.googleapis.com/projects/my-project/locations/us-central1/memberships/my-membership\" are valid values",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f ProjectFromIdFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	// Load arguments from function call
	var arg0 string
	resp.Diagnostics.Append(req.Arguments.GetArgument(ctx, 0, &arg0)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare how we'll identify project id from input string
	regex := regexp.MustCompile("projects/(?P<ProjectId>[^/]+)/") // Should match the pattern below
	template := "$ProjectId"                                      // Should match the submatch identifier in the regex
	pattern := "projects/{project}/"                              // Human-readable pseudo-regex pattern used in errors and warnings

	// Validate input
	ValidateElementFromIdArguments(arg0, regex, pattern, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get and return element from input string
	projectId := GetElementFromId(arg0, regex, template)
	resp.Diagnostics.Append(resp.Result.Set(ctx, projectId)...)
}
