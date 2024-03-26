// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package functions

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = RegionFromIdFunction{}

func NewRegionFromIdFunction() function.Function {
	return &RegionFromIdFunction{
		name: "region_from_id",
	}
}

type RegionFromIdFunction struct {
	name string // Makes function name available in Run logic for logging purposes
}

func (f RegionFromIdFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = f.name
}

func (f RegionFromIdFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Returns the region name within a provided resource id, self link, or OP style resource name.",
		Description: "Takes a single string argument, which should be a resource id, self link, or OP style resource name. This function will either return the region name from the input string or raise an error due to no region being present in the string. The function uses the presence of \"regions/{{region}}/\" in the input string to identify the region name, e.g. when the function is passed the id \"projects/my-project/regions/us-central1/subnetworks/my-subnetwork\" as an argument it will return \"us-central1\".",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "id",
				Description: "A string of a resource's id, a resource's self link, or an OP style resource name. For example, \"projects/my-project/regions/us-central1/subnetworks/my-subnetwork\" and \"https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1/subnetworks/my-subnetwork\" are valid values containing regions",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f RegionFromIdFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	// Load arguments from function call
	var arg0 string
	resp.Error = function.ConcatFuncErrors(req.Arguments.GetArgument(ctx, 0, &arg0))
	if resp.Error != nil {
		return
	}

	// Prepare how we'll identify region name from input string
	regex := regexp.MustCompile("regions/(?P<RegionName>[^/]+)/") // Should match the pattern below
	template := "$RegionName"                                     // Should match the submatch identifier in the regex
	pattern := "regions/{region}/"                                // Human-readable pseudo-regex pattern used in errors and warnings

	// Validate input
	resp.Error = function.ConcatFuncErrors(ValidateElementFromIdArguments(ctx, arg0, regex, pattern, f.name))
	if resp.Error != nil {
		return
	}

	// Get and return element from input string
	region := GetElementFromId(arg0, regex, template)
	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, region))
}
