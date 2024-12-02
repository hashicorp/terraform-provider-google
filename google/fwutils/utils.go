// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwutils

import "github.com/hashicorp/terraform-plugin-framework/types/basetypes"

func StringSet(d basetypes.SetValue) []string {

	StringSlice := make([]string, 0)
	for _, v := range d.Elements() {
		StringSlice = append(StringSlice, v.(basetypes.StringValue).ValueString())
	}
	return StringSlice
}
