// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package functions

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = ZoneFromIdFunction{}

func NewZoneFromIdFunction() function.Function {
	return &ZoneFromIdFunction{
		name: "zone_from_id",
	}
}

type ZoneFromIdFunction struct {
	name string // Makes function name available in Run logic for logging purposes
}

func (f ZoneFromIdFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = f.name
}

func (f ZoneFromIdFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Returns the zone name within the resource id or self link provided as an argument.",
		Description: "Takes a single string argument, which should be an id or self link of a resource. This function will either return the zone name from the input string or raise an error due to no zone being present in the string. The function uses the presence of \"zones/{{zone}}/\" in the input string to identify the zone name, e.g. when the function is passed the id \"projects/my-project/zones/us-central1-c/instances/my-instance\" as an argument it will return \"us-central1-c\".",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "id",
				Description: "An id of a resouce, or a self link. For example, both \"projects/my-project/zones/us-central1-c/instances/my-instance\" and \"https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-c/instances/my-instance\" are valid inputs",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f ZoneFromIdFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	// Load arguments from function call
	var arg0 string
	resp.Error = function.ConcatFuncErrors(req.Arguments.GetArgument(ctx, 0, &arg0))
	if resp.Error != nil {
		return
	}

	// Prepare how we'll identify zone name from input string
	regex := regexp.MustCompile("zones/(?P<ZoneName>[^/]+)/") // Should match the pattern below
	template := "$ZoneName"                                   // Should match the submatch identifier in the regex
	pattern := "zones/{zone}/"                                // Human-readable pseudo-regex pattern used in errors and warnings

	// Validate input
	resp.Error = function.ConcatFuncErrors(ValidateElementFromIdArguments(ctx, arg0, regex, pattern, f.name))
	if resp.Error != nil {
		return
	}

	// Get and return element from input string
	zone := GetElementFromId(arg0, regex, template)
	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, zone))
}
