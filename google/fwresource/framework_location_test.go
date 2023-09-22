// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package fwresource

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestLocationDescription_GetZone(t *testing.T) {
	cases := map[string]struct {
		ld            LocationDescription
		ExpectedZone  types.String
		ExpectedError bool
	}{
		"returns the value of the zone field in resource config": {
			ld: LocationDescription{
				// A resource would not have all 3 fields set, but if they were all present zone is used
				ResourceZone:     types.StringValue("resource-zone-a"),
				ResourceRegion:   types.StringValue("resource-region"),
				ResourceLocation: types.StringValue("resource-location"),
				// Provider config doesn't override resource config
				ProviderRegion: types.StringValue("provider-region"),
				ProviderZone:   types.StringValue("provider-zone-a"),
			},
			ExpectedZone: types.StringValue("resource-zone-a"),
		},
		"shortens zone values set as self links in the resource config": {
			ld: LocationDescription{
				ResourceZone: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/zones/resource-zone-a"),
			},
			ExpectedZone: types.StringValue("resource-zone-a"),
		},
		"returns the value of the zone field in provider config when zone is unset in resource config": {
			ld: LocationDescription{
				ResourceLocation: types.StringValue("resource-location"), // unused
				ResourceRegion:   types.StringValue("resource-region"),   // unused
				ProviderZone:     types.StringValue("provider-zone-a"),
			},
			ExpectedZone: types.StringValue("provider-zone-a"),
		},
		// Handling of empty strings
		"returns the value of the zone field in provider config when zone is set to an empty string in resource config": {
			ld: LocationDescription{
				ResourceZone: types.StringValue(""),
				ProviderZone: types.StringValue("provider-zone-a"),
			},
			ExpectedZone: types.StringValue("provider-zone-a"),
		},
		// Error states
		"returns an error when a zone value can't be found": {
			ld: LocationDescription{
				ResourceLocation: types.StringValue("resource-location"), // unused
				ResourceRegion:   types.StringValue("resource-region"),   // unused
			},
			ExpectedError: true,
		},
		"returns an error if zone is set as an empty string in both resource and provider configs": {
			ld: LocationDescription{
				ResourceZone: types.StringValue(""),
				ProviderZone: types.StringValue(""),
			},
			ExpectedError: true,
		},
		"returns an error that mention non-standard schema field names when a zone value can't be found": {
			ld: LocationDescription{
				ZoneSchemaField: types.StringValue("foobar"),
			},
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			zone, err := tc.ld.GetZone()

			if err != nil {
				if tc.ExpectedError {
					if !tc.ld.ZoneSchemaField.IsNull() {
						if !strings.Contains(err.Error(), tc.ld.ZoneSchemaField.ValueString()) {
							t.Fatalf("expected error to use provider schema field value %s, instead got: %s", tc.ld.ZoneSchemaField.ValueString(), err)
						}
					}
					return
				}
				t.Fatalf("unexpected error using test: %s", err)
			}
			if err == nil && tc.ExpectedError {
				t.Fatal("expected error but got none")
			}
			if zone != tc.ExpectedZone {
				t.Fatalf("Incorrect zone: got %s, want %s", zone, tc.ExpectedZone)
			}
		})
	}
}

func TestLocationDescription_GetRegion(t *testing.T) {
	cases := map[string]struct {
		ld             LocationDescription
		ExpectedRegion types.String
		ExpectedError  bool
	}{
		"returns the value of the region field in resource config": {
			ld: LocationDescription{
				// A resource would not have all 3 fields set, but if they were all present region is used first
				ResourceRegion:   types.StringValue("resource-region"),
				ResourceLocation: types.StringValue("resource-location"),
				ResourceZone:     types.StringValue("resource-zone-a"),
				// Provider config doesn't override resource config
				ProviderRegion: types.StringValue("provider-region"),
				ProviderZone:   types.StringValue("provider-zone-a"),
			},
			ExpectedRegion: types.StringValue("resource-region"),
		},
		"shortens region values set as self links in the resource config": {
			ld: LocationDescription{
				ResourceRegion: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1"),
			},
			ExpectedRegion: types.StringValue("us-central1"),
		},
		"returns a region derived from the zone field in resource config when region is unset": {
			ld: LocationDescription{
				ResourceZone:     types.StringValue("provider-zone-a"),
				ResourceLocation: types.StringValue("resource-location"), // unused
			},
			ExpectedRegion: types.StringValue("provider-zone"), // is truncated
		},
		"shortens region values when derived from a zone self link set in the resource config": {
			ld: LocationDescription{
				ResourceZone: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-a"),
			},
			ExpectedRegion: types.StringValue("us-central1"),
		},
		"shortens region values set as self links in the provider config": {
			ld: LocationDescription{
				ProviderRegion: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1"),
			},
			ExpectedRegion: types.StringValue("us-central1"),
		},
		"shortens region values when derived from a zone self link set in the provider config": {
			ld: LocationDescription{
				ProviderZone: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-a"),
			},
			ExpectedRegion: types.StringValue("us-central1"),
		},
		"returns the value of the region field in provider config when region/zone is unset in resource config": {
			ld: LocationDescription{
				ProviderRegion: types.StringValue("provider-region"),
				ProviderZone:   types.StringValue("provider-zone-a"), // unused
			},
			ExpectedRegion: types.StringValue("provider-region"),
		},
		"returns a region derived from the zone field in provider config when region unset in both resource and provider config": {
			ld: LocationDescription{
				ProviderZone: types.StringValue("provider-zone-a"),
			},
			ExpectedRegion: types.StringValue("provider-zone"), // is truncated
		},
		// Handling of empty strings
		"returns a region derived from the zone field in resource config when region is set as an empty string": {
			ld: LocationDescription{
				ResourceRegion: types.StringValue(""),
				ResourceZone:   types.StringValue("provider-zone-a"),
			},
			ExpectedRegion: types.StringValue("provider-zone"), // is truncated
		},
		"returns the value of the region field in provider config when region/zone set as an empty string in resource config": {
			ld: LocationDescription{
				ResourceRegion: types.StringValue(""),
				ResourceZone:   types.StringValue(""),
				ProviderRegion: types.StringValue("provider-region"),
			},
			ExpectedRegion: types.StringValue("provider-region"),
		},
		// Error states
		"returns an error when region/zone values can't be found (location is ignored)": {
			ld: LocationDescription{
				ResourceLocation: types.StringValue("resource-location"),
			},
			ExpectedError: true,
		},
		"returns an error if region and zone set as empty strings in both resource and provider configs": {
			ld: LocationDescription{
				ResourceRegion: types.StringValue(""),
				ResourceZone:   types.StringValue(""),
				ProviderRegion: types.StringValue(""),
				ProviderZone:   types.StringValue(""),
			},
			ExpectedError: true,
		},
		"returns an error that mention non-standard schema field names when region value can't be found": {
			ld: LocationDescription{
				RegionSchemaField: types.StringValue("foobar"),
			},
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			region, err := tc.ld.GetRegion()

			if err != nil {
				if tc.ExpectedError {
					if !tc.ld.RegionSchemaField.IsNull() {
						if !strings.Contains(err.Error(), tc.ld.RegionSchemaField.ValueString()) {
							t.Fatalf("expected error to use provider schema field value %s, instead got: %s", tc.ld.RegionSchemaField.ValueString(), err)
						}
					}
					return
				}
				t.Fatalf("unexpected error using test: %s", err)
			}
			if err == nil && tc.ExpectedError {
				t.Fatal("expected error but got none")
			}
			if region != tc.ExpectedRegion {
				t.Fatalf("Incorrect region: got %s, want %s", region, tc.ExpectedRegion)
			}
		})
	}
}

func TestLocationDescription_GetLocation(t *testing.T) {
	cases := map[string]struct {
		ld               LocationDescription
		ExpectedLocation types.String
		ExpectedError    bool
	}{
		"returns the value of the location field in resource config": {
			ld: LocationDescription{
				// A resource would not have all 3 fields set, but if they were all present location is used first
				ResourceLocation: types.StringValue("resource-location"),
				ResourceRegion:   types.StringValue("resource-region"),
				ResourceZone:     types.StringValue("resource-zone-a"),
				// Provider config doesn't override resource config
				ProviderRegion: types.StringValue("provider-region"),
				ProviderZone:   types.StringValue("provider-zone-a"),
			},
			ExpectedLocation: types.StringValue("resource-location"),
		},
		"shortens the location value when it is set as a self link in the resource config": {
			ld: LocationDescription{
				ResourceLocation: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/locations/resource-location"),
			},
			ExpectedLocation: types.StringValue("resource-location"),
		},
		"returns the region value set in the resource config when location is not in the schema": {
			ld: LocationDescription{
				ResourceRegion: types.StringValue("resource-region"),
				ResourceZone:   types.StringValue("resource-zone-a"), // unused
			},
			ExpectedLocation: types.StringValue("resource-region"),
		},
		"shortens the region value when it is set as a self link in the resource config": {
			ld: LocationDescription{
				ResourceRegion: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/regions/resource-region"),
			},
			ExpectedLocation: types.StringValue("resource-region"),
		},
		"returns the zone value set in the resource config when neither location nor region in the schema": {
			ld: LocationDescription{
				ResourceZone: types.StringValue("resource-zone-a"),
			},
			ExpectedLocation: types.StringValue("resource-zone-a"),
		},
		"shortens zone values set as self links in the resource config": {
			ld: LocationDescription{
				ResourceZone: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/zones/resource-zone-a"),
			},
			ExpectedLocation: types.StringValue("resource-zone-a"),
		},
		"returns the region value from the provider config when none of location/region/zone are set in the resource config": {
			ld: LocationDescription{
				ProviderRegion: types.StringValue("provider-region"), // Preferred to use region value over zone value if both are set
				ProviderZone:   types.StringValue("provider-zone-a"),
			},
			ExpectedLocation: types.StringValue("provider-region"),
		},
		"returns the zone value from the provider config when none of location/region/zone are set in the resource config and region is not set in the provider config": {
			ld: LocationDescription{
				ProviderZone: types.StringValue("provider-zone-a"),
			},
			ExpectedLocation: types.StringValue("provider-zone-a"),
		},
		"shortens the zone value when it is set as a self link in the provider config": {
			ld: LocationDescription{
				ProviderZone: types.StringValue("https://www.googleapis.com/compute/v1/projects/my-project/zones/provider-zone-a"),
			},
			ExpectedLocation: types.StringValue("provider-zone-a"),
		},
		// Handling of empty strings
		"returns the region value set in the resource config when location is an empty string": {
			ld: LocationDescription{
				ResourceLocation: types.StringValue(""),
				ResourceRegion:   types.StringValue("resource-region"),
			},
			ExpectedLocation: types.StringValue("resource-region"),
		},
		"returns the zone value set in the resource config when both location or region are empty strings": {
			ld: LocationDescription{
				ResourceLocation: types.StringValue(""),
				ResourceRegion:   types.StringValue(""),
				ResourceZone:     types.StringValue("resource-zone-a"),
			},
			ExpectedLocation: types.StringValue("resource-zone-a"),
		},
		"returns the zone value from the provider config when all of location/region/zone are set as empty strings in the resource config": {
			ld: LocationDescription{
				ResourceLocation: types.StringValue(""),
				ResourceRegion:   types.StringValue(""),
				ResourceZone:     types.StringValue(""),
				ProviderZone:     types.StringValue("provider-zone-a"),
			},
			ExpectedLocation: types.StringValue("provider-zone-a"),
		},
		"returns an error when none of location/region/zone are set on the resource, and neither region or zone is set on the provider": {
			ExpectedError: true,
		},
		"returns an error if location/region/zone are set as empty strings in both resource and provider configs": {
			ld: LocationDescription{
				ResourceLocation: types.StringValue(""),
				ResourceRegion:   types.StringValue(""),
				ResourceZone:     types.StringValue(""),
				ProviderRegion:   types.StringValue(""),
				ProviderZone:     types.StringValue(""),
			},
			ExpectedError: true,
		},
		"returns an error that mention non-standard schema field names when location value can't be found": {
			ld: LocationDescription{
				LocationSchemaField: types.StringValue("foobar"),
			},
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			region, err := tc.ld.GetLocation()

			if err != nil {
				if tc.ExpectedError {
					if !tc.ld.LocationSchemaField.IsNull() {
						if !strings.Contains(err.Error(), tc.ld.LocationSchemaField.ValueString()) {
							t.Fatalf("expected error to use provider schema field value %s, instead got: %s", tc.ld.LocationSchemaField.ValueString(), err)
						}
					}
					return
				}
				t.Fatalf("unexpected error using test: %s", err)
			}
			if err == nil && tc.ExpectedError {
				t.Fatal("expected error but got none")
			}
			if region != tc.ExpectedLocation {
				t.Fatalf("Incorrect location: got %s, want %s", region, tc.ExpectedLocation)
			}
		})
	}
}
