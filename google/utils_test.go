package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var fictionalSchema = map[string]*schema.Schema{
	"location": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"region": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"zone": {
		Type:     schema.TypeString,
		Optional: true,
	},
	"project": {
		Type:     schema.TypeString,
		Optional: true,
	},
}

func TestGetProject(t *testing.T) {
	cases := map[string]struct {
		ResourceConfig  map[string]interface{}
		ProviderConfig  map[string]string
		ExpectedProject string
		ExpectedError   bool
	}{
		"project is pulled from resource config instead of provider config": {
			ResourceConfig: map[string]interface{}{
				"project": "resource-project",
			},
			ProviderConfig: map[string]string{
				"project": "provider-project",
			},
			ExpectedProject: "resource-project",
		},
		"project is pulled from provider config when not set on resource": {
			ProviderConfig: map[string]string{
				"project": "provider-project",
			},
			ExpectedProject: "provider-project",
		},
		"error returned when project not set on either provider or resource": {
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange

			// Create provider config
			var config transport_tpg.Config
			if v, ok := tc.ProviderConfig["project"]; ok {
				config.Project = v
			}

			// Create resource config
			// Here use a fictional schema that includes a project field
			d := setupTestResourceDataFromConfigMap(t, fictionalSchema, tc.ResourceConfig)

			// Act
			project, err := tpgresource.GetProject(d, &config)

			// Assert
			if err != nil {
				if tc.ExpectedError {
					return
				}
				t.Fatalf("Unexpected error using test: %s", err)
			}

			if project != tc.ExpectedProject {
				t.Fatalf("Incorrect project: got %s, want %s", project, tc.ExpectedProject)
			}
		})
	}
}

func TestGetLocation(t *testing.T) {
	cases := map[string]struct {
		ResourceConfig   map[string]interface{}
		ProviderConfig   map[string]string
		ExpectedLocation string
		ExpectError      bool
	}{
		"returns the value of the location field in resource config": {
			ResourceConfig: map[string]interface{}{
				"location": "resource-location",
				"region":   "resource-region", // unused
				"zone":     "resource-zone-a", // unused
			},
			ExpectedLocation: "resource-location",
		},
		"does not shorten the location value when it is set as a self link in the resource config": {
			ResourceConfig: map[string]interface{}{
				"location": "https://www.googleapis.com/compute/v1/projects/my-project/locations/resource-location",
			},
			ExpectedLocation: "https://www.googleapis.com/compute/v1/projects/my-project/locations/resource-location", // No shortening takes place
		},
		"returns the region value set in the resource config when location is not in the schema": {
			ResourceConfig: map[string]interface{}{
				"region": "resource-region",
				"zone":   "resource-zone-a", // unused
			},
			ExpectedLocation: "resource-region",
		},
		"does not shorten the region value when it is set as a self link in the resource config": {
			ResourceConfig: map[string]interface{}{
				"region": "https://www.googleapis.com/compute/v1/projects/my-project/region/resource-region",
			},
			ExpectedLocation: "https://www.googleapis.com/compute/v1/projects/my-project/region/resource-region", // No shortening takes place
		},
		"returns the zone value set in the resource config when neither location nor region in the schema": {
			ResourceConfig: map[string]interface{}{
				"zone": "resource-zone-a",
			},
			ExpectedLocation: "resource-zone-a",
		},
		"shortens zone values set as self links in the resource config": {
			// Results from GetLocation using GetZone internally
			// This behaviour makes sense because APIs may return a self link as the zone value
			ResourceConfig: map[string]interface{}{
				"zone": "https://www.googleapis.com/compute/v1/projects/my-project/zones/resource-zone-a",
			},
			ExpectedLocation: "resource-zone-a",
		},
		"returns the zone value from the provider config when none of location/region/zone are set in the resource config": {
			ProviderConfig: map[string]string{
				"zone": "provider-zone-a",
			},
			ExpectedLocation: "provider-zone-a",
		},
		"does not shorten the zone value when it is set as a self link in the provider config": {
			// This behaviour makes sense because provider config values don't originate from APIs
			// Users should always configure the provider with the short names of regions/zones
			ProviderConfig: map[string]string{
				"zone": "https://www.googleapis.com/compute/v1/projects/my-project/zones/provider-zone-a",
			},
			ExpectedLocation: "https://www.googleapis.com/compute/v1/projects/my-project/zones/provider-zone-a",
		},
		// Handling of empty strings
		"returns the region value set in the resource config when location is an empty string": {
			ResourceConfig: map[string]interface{}{
				"location": "",
				"region":   "resource-region",
			},
			ExpectedLocation: "resource-region",
		},
		"returns the zone value set in the resource config when both location or region are empty strings": {
			ResourceConfig: map[string]interface{}{
				"location": "",
				"region":   "",
				"zone":     "resource-zone-a",
			},
			ExpectedLocation: "resource-zone-a",
		},
		"returns the zone value from the provider config when all of location/region/zone are set as empty strings in the resource config": {
			ResourceConfig: map[string]interface{}{
				"location": "",
				"region":   "",
				"zone":     "",
			},
			ProviderConfig: map[string]string{
				"zone": "provider-zone-a",
			},
			ExpectedLocation: "provider-zone-a",
		},
		// Error states
		"returns an error when only a region value is set in the the provider config and none of location/region/zone are set in the resource config": {
			ProviderConfig: map[string]string{
				"region": "provider-region",
			},
			ExpectError: true,
		},
		"returns an error when none of location/region/zone are set on the resource, and neither region or zone is set on the provider": {
			ExpectError: true,
		},
		"returns an error if location/region/zone are set as empty strings in both resource and provider configs": {
			ResourceConfig: map[string]interface{}{
				"location": "",
				"region":   "",
				"zone":     "",
			},
			ProviderConfig: map[string]string{
				"zone": "",
			},
			ExpectError: true,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange

			// Create provider config
			var config transport_tpg.Config
			if v, ok := tc.ProviderConfig["region"]; ok {
				config.Region = v
			}
			if v, ok := tc.ProviderConfig["zone"]; ok {
				config.Zone = v
			}

			// Create resource config
			// Here use a fictional schema as example because we need to have all of
			// location, region, and zone fields present in the schema for the test,
			// and no real resources would contain all of these
			d := setupTestResourceDataFromConfigMap(t, fictionalSchema, tc.ResourceConfig)

			// Act
			location, err := tpgresource.GetLocation(d, &config)

			// Assert
			if err != nil {
				if tc.ExpectError {
					return
				}
				t.Fatalf("unexpected error using test: %s", err)
			}

			if location != tc.ExpectedLocation {
				t.Fatalf("incorrect location: got %s, want %s", location, tc.ExpectedLocation)
			}
		})
	}
}

func TestGetZone(t *testing.T) {
	cases := map[string]struct {
		ResourceConfig map[string]interface{}
		ProviderConfig map[string]string
		ExpectedZone   string
		ExpectedError  bool
	}{
		"returns the value of the zone field in resource config": {
			ResourceConfig: map[string]interface{}{
				"zone": "resource-zone-a",
			},
			ProviderConfig: map[string]string{
				"zone": "provider-zone-a",
			},
			ExpectedZone: "resource-zone-a",
		},
		"shortens zone values set as self links in the resource config": {
			ResourceConfig: map[string]interface{}{
				"zone": "https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-a",
			},
			ExpectedZone: "us-central1-a",
		},
		"returns the value of the zone field in provider config when zone is unset in resource config": {
			ProviderConfig: map[string]string{
				"zone": "provider-zone-a",
			},
			ExpectedZone: "provider-zone-a",
		},
		// Handling of empty strings
		"returns the value of the zone field in provider config when zone is set to an empty string in resource config": {
			ResourceConfig: map[string]interface{}{
				"zone": "",
			},
			ProviderConfig: map[string]string{
				"zone": "provider-zone-a",
			},
			ExpectedZone: "provider-zone-a",
		},
		// Error states
		"returns an error when a zone value can't be found": {
			ResourceConfig: map[string]interface{}{
				"location": "resource-location", // unused
				"region":   "resource-region",   // unused
			},
			ProviderConfig: map[string]string{
				"region": "provider-region", //unused
			},
			ExpectedError: true,
		},
		"returns an error if zone is set as an empty string in both resource and provider configs": {
			ResourceConfig: map[string]interface{}{
				"zone": "",
			},
			ProviderConfig: map[string]string{
				"zone": "",
			},
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange

			// Create provider config
			var config transport_tpg.Config
			if v, ok := tc.ProviderConfig["zone"]; ok {
				config.Zone = v
			}

			// Create resource config
			// Here use a fictional schema as example because we need to have all of
			// location, region, and zone fields present in the schema for the test,
			// and no real resources would contain all of these
			d := setupTestResourceDataFromConfigMap(t, fictionalSchema, tc.ResourceConfig)

			// Act
			zone, err := tpgresource.GetZone(d, &config)

			// Assert
			if err != nil {
				if tc.ExpectedError {
					return
				}
				t.Fatalf("Unexpected error using test: %s", err)
			}

			if zone != tc.ExpectedZone {
				t.Fatalf("Incorrect zone: got %s, want %s", zone, tc.ExpectedZone)
			}
		})
	}
}

func TestGetRegion(t *testing.T) {
	cases := map[string]struct {
		ResourceConfig map[string]interface{}
		ProviderConfig map[string]string
		ExpectedRegion string
		ExpectedError  bool
	}{
		"returns the value of the region field in resource config": {
			ResourceConfig: map[string]interface{}{
				"region":   "resource-region",
				"zone":     "resource-zone-a",
				"location": "resource-location", // unused
			},
			ProviderConfig: map[string]string{
				"region": "provider-region",
				"zone":   "provider-zone-a",
			},
			ExpectedRegion: "resource-region",
		},
		"shortens region values set as self links in the resource config": {
			ResourceConfig: map[string]interface{}{
				"region": "https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1",
			},
			ExpectedRegion: "us-central1",
		},
		"returns a region derived from the zone field in resource config when region is unset": {
			ResourceConfig: map[string]interface{}{
				"zone":     "resource-zone-a",
				"location": "resource-location", // unused
			},
			ExpectedRegion: "resource-zone", // is truncated
		},
		"does not shorten region values when derived from a zone self link set in the resource config": {
			ResourceConfig: map[string]interface{}{
				"zone": "https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-a",
			},
			ExpectedRegion: "https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1", // Value is not shortenedfrom URI to name
		},
		"returns the value of the region field in provider config when region/zone is unset in resource config": {
			ProviderConfig: map[string]string{
				"region": "provider-region",
				"zone":   "provider-zone-a", // unused
			},
			ExpectedRegion: "provider-region",
		},
		"returns a region derived from the zone field in provider config when region unset in both resource and provider config": {
			ProviderConfig: map[string]string{
				"zone": "provider-zone-a",
			},
			ExpectedRegion: "provider-zone", // is truncated
		},
		// Handling of empty strings
		"returns a region derived from the zone field in resource config when region is set as an empty string": {
			ResourceConfig: map[string]interface{}{
				"region": "",
				"zone":   "resource-zone-a",
			},
			ExpectedRegion: "resource-zone", // is truncated
		},
		"returns the value of the region field in provider config when region/zone set as an empty string in resource config": {
			ResourceConfig: map[string]interface{}{
				"region": "",
				"zone":   "",
			},
			ProviderConfig: map[string]string{
				"region": "provider-region",
			},
			ExpectedRegion: "provider-region",
		},
		// Error states
		"returns an error when region values can't be found": {
			ResourceConfig: map[string]interface{}{
				"location": "resource-location",
			},
			ExpectedError: true,
		},
		"returns an error if region and zone set as empty strings in both resource and provider configs": {
			ResourceConfig: map[string]interface{}{
				"region": "",
				"zone":   "",
			},
			ProviderConfig: map[string]string{
				"region": "",
				"zone":   "",
			},
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange

			// Create provider config
			var config transport_tpg.Config
			if v, ok := tc.ProviderConfig["region"]; ok {
				config.Region = v
			}
			if v, ok := tc.ProviderConfig["zone"]; ok {
				config.Zone = v
			}

			// Create resource config
			// Here use a fictional schema as example because we need to have all of
			// location, region, and zone fields present in the schema for the test,
			// and no real resources would contain all of these
			d := setupTestResourceDataFromConfigMap(t, fictionalSchema, tc.ResourceConfig)

			// Act
			region, err := tpgresource.GetRegion(d, &config)

			// Assert
			if err != nil {
				if tc.ExpectedError {
					return
				}
				t.Fatalf("Unexpected error using test: %s", err)
			}

			if region != tc.ExpectedRegion {
				t.Fatalf("Incorrect region: got %s, want %s", region, tc.ExpectedRegion)
			}
		})
	}
}

func TestCheckGCSName(t *testing.T) {
	valid63 := RandString(t, 63)
	cases := map[string]bool{
		// Valid
		"foobar":       true,
		"foobar1":      true,
		"12345":        true,
		"foo_bar_baz":  true,
		"foo-bar-baz":  true,
		"foo-bar_baz1": true,
		"foo--bar":     true,
		"foo__bar":     true,
		"foo-goog":     true,
		"foo.goog":     true,
		valid63:        true,
		fmt.Sprintf("%s.%s.%s", valid63, valid63, valid63): true,

		// Invalid
		"goog-foobar":     false,
		"foobar-google":   false,
		"-foobar":         false,
		"foobar-":         false,
		"_foobar":         false,
		"foobar_":         false,
		"fo":              false,
		"foo$bar":         false,
		"foo..bar":        false,
		RandString(t, 64): false,
		fmt.Sprintf("%s.%s.%s.%s", valid63, valid63, valid63, valid63): false,
	}

	for bucketName, valid := range cases {
		err := tpgresource.CheckGCSName(bucketName)
		if valid && err != nil {
			t.Errorf("The bucket name %s was expected to pass validation and did not pass.", bucketName)
		} else if !valid && err == nil {
			t.Errorf("The bucket name %s was NOT expected to pass validation and passed.", bucketName)
		}
	}
}
