// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package functions

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = RegionFromZoneFunction{}

func NewRegionFromZoneFunction() function.Function {
	return &RegionFromZoneFunction{}
}

type RegionFromZoneFunction struct{}

func (f RegionFromZoneFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "region_from_zone"
}

func (f RegionFromZoneFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Returns the region within a provided resource's zone",
		Description: "Takes a single string argument, which should be a resource's zone.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "zone",
				Description: "A string of a resource's zone.",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f RegionFromZoneFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	// Load arguments from function call
	var zone string
	resp.Error = function.ConcatFuncErrors(req.Arguments.GetArgument(ctx, 0, &zone))
	if resp.Error != nil {
		return
	}

	if zone == "" {
		err := function.NewArgumentFuncError(0, "The input string cannot be empty.")
		resp.Error = function.ConcatFuncErrors(err)
		return
	}

	zoneParts := strings.Split(zone, "-")

	if len(zoneParts) < 3 {
		err := function.NewArgumentFuncError(0, fmt.Sprintf("The input string \"%s\" is not a valid zone name.", zone))
		resp.Error = function.ConcatFuncErrors(err)
	} else {
		resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, strings.Join(zoneParts[:len(zoneParts)-1], "-")))
	}
}
