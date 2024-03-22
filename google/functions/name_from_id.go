// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package functions

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = NameFromIdFunction{}

func NewNameFromIdFunction() function.Function {
	return &NameFromIdFunction{
		name: "name_from_id",
	}
}

type NameFromIdFunction struct {
	name string
}

func (f NameFromIdFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = f.name
}

func (f NameFromIdFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Returns the short-form name of a resource within a provided resource's id, resource URI, self link, or full resource name.",
		Description: "Takes a single string argument, which should be a resource's id, resource URI, self link, or full resource name. This function will return the short-form name of a resource from the input string, or raise an error due to a problem with the input string. The function returns the final element in the input string as the resource's name, e.g. when the function is passed the id \"projects/my-project/zones/us-central1-c/instances/my-instance\" as an argument it will return \"my-instance\".",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "id",
				Description: "A string of a resource's id, resource URI, self link, or full resource name. For example, \"projects/my-project/zones/us-central1-c/instances/my-instance\", \"https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-c/instances/my-instance\" and \"//gkehub.googleapis.com/projects/my-project/locations/us-central1/memberships/my-membership\" are valid values",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f NameFromIdFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	// Load arguments from function call
	var arg0 string
	resp.Error = function.ConcatFuncErrors(req.Arguments.GetArgument(ctx, 0, &arg0))
	if resp.Error != nil {
		return
	}

	// Prepare how we'll identify resource name from input string
	regex := regexp.MustCompile("/(?P<ResourceName>[^/]+)$") // Should match the pattern below
	template := "$ResourceName"                              // Should match the submatch identifier in the regex
	pattern := "resourceType/{name}$"                        // Human-readable pseudo-regex pattern used in errors and warnings

	// Validate input
	resp.Error = function.ConcatFuncErrors(ValidateElementFromIdArguments(ctx, arg0, regex, pattern, f.name))
	if resp.Error != nil {
		return
	}

	// Get and return element from input string
	name := GetElementFromId(arg0, regex, template)
	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, name))
}
