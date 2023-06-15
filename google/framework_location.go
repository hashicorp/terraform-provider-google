// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package google

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LocationDescriber interface {
	GetLocationDescription(providerConfig *frameworkProvider) LocationDescription
}

type LocationDescription struct {
	// Location - not configurable on provider
	LocationSchemaField types.String
	ResourceLocation    types.String

	// Region
	RegionSchemaField types.String
	ResourceRegion    types.String
	ProviderRegion    types.String

	// Zone
	ZoneSchemaField types.String
	ResourceZone    types.String
	ProviderZone    types.String
}

func (ld *LocationDescription) GetLocation() (types.String, error) {
	// Location from resource config
	if !ld.ResourceLocation.IsNull() && !ld.ResourceLocation.IsUnknown() && !ld.ResourceLocation.Equal(types.StringValue("")) {
		return ld.ResourceLocation, nil
	}

	// Location from region in resource config
	if !ld.ResourceRegion.IsNull() && !ld.ResourceRegion.IsUnknown() && !ld.ResourceRegion.Equal(types.StringValue("")) {
		return ld.ResourceRegion, nil
	}

	// Location from zone in resource config
	if !ld.ResourceZone.IsNull() && !ld.ResourceZone.IsUnknown() && !ld.ResourceZone.Equal(types.StringValue("")) {
		location := GetResourceNameFromSelfLink(ld.ResourceZone.ValueString()) // Zone could be a self link
		return types.StringValue(location), nil
	}

	// Location from zone in provider config
	if !ld.ProviderZone.IsNull() && !ld.ProviderZone.IsUnknown() && !ld.ProviderZone.Equal(types.StringValue("")) {
		return ld.ProviderZone, nil
	}

	var err error
	if !ld.LocationSchemaField.IsNull() {
		err = fmt.Errorf("location could not be identified, please add `%s` in your resource or set `region` in your provider configuration block", ld.LocationSchemaField.ValueString())
	} else {
		err = errors.New("location could not be identified, please add `location` in your resource or `region` in your provider configuration block")
	}
	return types.StringNull(), err
}

func (ld *LocationDescription) GetRegion() (types.String, error) {
	// TODO(SarahFrench): Make empty strings not ignored, see https://github.com/hashicorp/terraform-provider-google/issues/14447
	// For all checks in this function body

	// Region from resource config
	if !ld.ResourceRegion.IsNull() && !ld.ResourceRegion.IsUnknown() && !ld.ResourceRegion.Equal(types.StringValue("")) {
		region := GetResourceNameFromSelfLink(ld.ResourceRegion.ValueString()) // Region could be a self link
		return types.StringValue(region), nil
	}
	// Region from zone in resource config
	if !ld.ResourceZone.IsNull() && !ld.ResourceZone.IsUnknown() && !ld.ResourceZone.Equal(types.StringValue("")) {
		region := getRegionFromZone(ld.ResourceZone.ValueString())
		return types.StringValue(region), nil
	}
	// Region from provider config
	if !ld.ProviderRegion.IsNull() && !ld.ProviderRegion.IsUnknown() && !ld.ProviderRegion.Equal(types.StringValue("")) {
		return ld.ProviderRegion, nil
	}
	// Region from zone in provider config
	if !ld.ProviderZone.IsNull() && !ld.ProviderZone.IsUnknown() && !ld.ProviderZone.Equal(types.StringValue("")) {
		region := getRegionFromZone(ld.ProviderZone.ValueString())
		return types.StringValue(region), nil
	}

	var err error
	if !ld.RegionSchemaField.IsNull() {
		err = fmt.Errorf("region could not be identified, please add `%s` in your resource or set `region` in your provider configuration block", ld.RegionSchemaField.ValueString())
	} else {
		err = errors.New("region could not be identified, please add `region` in your resource or provider configuration block")
	}
	return types.StringNull(), err
}

func (ld *LocationDescription) GetZone() (types.String, error) {
	// TODO(SarahFrench): Make empty strings not ignored, see https://github.com/hashicorp/terraform-provider-google/issues/14447
	// For all checks in this function body

	if !ld.ResourceZone.IsNull() && !ld.ResourceZone.IsUnknown() && !ld.ResourceZone.Equal(types.StringValue("")) {
		// Zone could be a self link
		zone := GetResourceNameFromSelfLink(ld.ResourceZone.ValueString())
		return types.StringValue(zone), nil
	}
	if !ld.ProviderZone.IsNull() && !ld.ProviderZone.IsUnknown() && !ld.ProviderZone.Equal(types.StringValue("")) {
		return ld.ProviderZone, nil
	}

	var err error
	if !ld.ZoneSchemaField.IsNull() {
		err = fmt.Errorf("zone could not be identified, please add `%s` in your resource or `zone` in your provider configuration block", ld.ZoneSchemaField.ValueString())
	} else {
		err = errors.New("zone could not be identified, please add `zone` in your resource or provider configuration block")
	}
	return types.StringNull(), err
}
