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
	return &ProjectFromIdFunction{
		name: "project_from_id",
	}
}

type ProjectFromIdFunction struct {
	name string // Makes function name available in Run logic for logging purposes
}

func (f ProjectFromIdFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = f.name
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
	resp.Error = function.ConcatFuncErrors(req.Arguments.GetArgument(ctx, 0, &arg0))
	if resp.Error != nil {
		return
	}

	// Prepare how we'll identify project id from input string
	regex := regexp.MustCompile("projects/(?P<ProjectId>[^/]+)/") // Should match the pattern below
	template := "$ProjectId"                                      // Should match the submatch identifier in the regex
	pattern := "projects/{project}/"                              // Human-readable pseudo-regex pattern used in errors and warnings

	// Validate input
	resp.Error = function.ConcatFuncErrors(ValidateElementFromIdArguments(ctx, arg0, regex, pattern, f.name))
	if resp.Error != nil {
		return
	}

	// Get and return element from input string
	projectId := GetElementFromId(arg0, regex, template)
	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, projectId))
}
