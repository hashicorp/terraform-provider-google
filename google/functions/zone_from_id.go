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
	return &ZoneFromIdFunction{}
}

type ZoneFromIdFunction struct{}

func (f ZoneFromIdFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "zone_from_id"
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
	resp.Diagnostics.Append(req.Arguments.GetArgument(ctx, 0, &arg0)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare how we'll identify zone name from input string
	regex := regexp.MustCompile("zones/(?P<ZoneName>[^/]+)/") // Should match the pattern below
	template := "$ZoneName"                                   // Should match the submatch identifier in the regex
	pattern := "zones/{zone}/"                                // Human-readable pseudo-regex pattern used in errors and warnings

	// Validate input
	ValidateElementFromIdArguments(arg0, regex, pattern, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get and return element from input string
	zone := GetElementFromId(arg0, regex, template)
	resp.Diagnostics.Append(resp.Result.Set(ctx, zone)...)
}
