// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package functions

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = LocationFromIdFunction{}

func NewLocationFromIdFunction() function.Function {
	return &LocationFromIdFunction{
		name: "location_from_id",
	}
}

type LocationFromIdFunction struct {
	name string // Makes function name available in Run logic for logging purposes
}

func (f LocationFromIdFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = f.name
}

func (f LocationFromIdFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Returns the location name within a provided resource id, self link, or OP style resource name.",
		Description: "Takes a single string argument, which should be a resource id, self link, or OP style resource name. This function will either return the location name from the input string or raise an error due to no location being present in the string. The function uses the presence of \"locations/{{location}}/\" in the input string to identify the location name, e.g. when the function is passed the id \"projects/my-project/locations/us-central1/services/my-service\" as an argument it will return \"us-central1\".",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "id",
				Description: "A string of a resource's id, a resource's self link, or an OP style resource name. For example, \"projects/my-project/locations/us-central1/services/my-service\" and \"https://run.googleapis.com/v2/projects/my-project/locations/us-central1/services/my-service\" are valid values containing locations",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f LocationFromIdFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	// Load arguments from function call
	var arg0 string
	resp.Error = function.ConcatFuncErrors(req.Arguments.GetArgument(ctx, 0, &arg0))
	if resp.Error != nil {
		return
	}

	// Prepare how we'll identify location name from input string
	regex := regexp.MustCompile("locations/(?P<LocationName>[^/]+)/") // Should match the pattern below
	template := "$LocationName"                                       // Should match the submatch identifier in the regex
	pattern := "locations/{location}/"                                // Human-readable pseudo-regex pattern used in errors and warnings

	// Validate input
	resp.Error = function.ConcatFuncErrors(ValidateElementFromIdArguments(ctx, arg0, regex, pattern, f.name))
	if resp.Error != nil {
		return
	}

	// Get and return element from input string
	location := GetElementFromId(arg0, regex, template)
	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, location))
}
