package tpgresource

import (
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/googleapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func TestConvertStringArr(t *testing.T) {
	input := make([]interface{}, 3)
	input[0] = "aaa"
	input[1] = "bbb"
	input[2] = "aaa"

	expected := []string{"aaa", "bbb", "ccc"}
	actual := ConvertStringArr(input)

	if reflect.DeepEqual(expected, actual) {
		t.Fatalf("(%s) did not match expected value: %s", actual, expected)
	}
}

func TestConvertAndMapStringArr(t *testing.T) {
	input := make([]interface{}, 3)
	input[0] = "aaa"
	input[1] = "bbb"
	input[2] = "aaa"

	expected := []string{"AAA", "BBB", "CCC"}
	actual := ConvertAndMapStringArr(input, strings.ToUpper)

	if reflect.DeepEqual(expected, actual) {
		t.Fatalf("(%s) did not match expected value: %s", actual, expected)
	}
}

func TestConvertStringMap(t *testing.T) {
	input := make(map[string]interface{}, 3)
	input["one"] = "1"
	input["two"] = "2"
	input["three"] = "3"

	expected := map[string]string{"one": "1", "two": "2", "three": "3"}
	actual := ConvertStringMap(input)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("%s did not match expected value: %s", actual, expected)
	}
}

func TestIpCidrRangeDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"single ip address": {
			Old:                "10.2.3.4",
			New:                "10.2.3.5",
			ExpectDiffSuppress: false,
		},
		"cidr format string": {
			Old:                "10.1.2.0/24",
			New:                "10.1.3.0/24",
			ExpectDiffSuppress: false,
		},
		"netmask same mask": {
			Old:                "10.1.2.0/24",
			New:                "/24",
			ExpectDiffSuppress: true,
		},
		"netmask different mask": {
			Old:                "10.1.2.0/24",
			New:                "/32",
			ExpectDiffSuppress: false,
		},
		"add netmask": {
			Old:                "",
			New:                "/24",
			ExpectDiffSuppress: false,
		},
		"remove netmask": {
			Old:                "/24",
			New:                "",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if IpCidrRangeDiffSuppress("ip_cidr_range", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestRfc3339TimeDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"same time, format changed to have leading zero": {
			Old:                "2:00",
			New:                "02:00",
			ExpectDiffSuppress: true,
		},
		"same time, format changed not to have leading zero": {
			Old:                "02:00",
			New:                "2:00",
			ExpectDiffSuppress: true,
		},
		"different time, both without leading zero": {
			Old:                "2:00",
			New:                "3:00",
			ExpectDiffSuppress: false,
		},
		"different time, old with leading zero, new without": {
			Old:                "02:00",
			New:                "3:00",
			ExpectDiffSuppress: false,
		},
		"different time, new with leading zero, oldwithout": {
			Old:                "2:00",
			New:                "03:00",
			ExpectDiffSuppress: false,
		},
		"different time, both with leading zero": {
			Old:                "02:00",
			New:                "03:00",
			ExpectDiffSuppress: false,
		},
	}
	for tn, tc := range cases {
		if Rfc3339TimeDiffSuppress("time", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Errorf("bad: %s, '%s' => '%s' expect DiffSuppress to return %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
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
			d := SetupTestResourceDataFromConfigMap(t, fictionalSchema, tc.ResourceConfig)

			// Act
			project, err := GetProject(d, &config)

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
			d := SetupTestResourceDataFromConfigMap(t, fictionalSchema, tc.ResourceConfig)

			// Act
			location, err := GetLocation(d, &config)

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
			d := SetupTestResourceDataFromConfigMap(t, fictionalSchema, tc.ResourceConfig)

			// Act
			zone, err := GetZone(d, &config)

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
			d := SetupTestResourceDataFromConfigMap(t, fictionalSchema, tc.ResourceConfig)

			// Act
			region, err := GetRegion(d, &config)

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

func TestGetRegionFromZone(t *testing.T) {
	expected := "us-central1"
	actual := GetRegionFromZone("us-central1-f")
	if expected != actual {
		t.Fatalf("Region (%s) did not match expected value: %s", actual, expected)
	}
}

func TestDatasourceSchemaFromResourceSchema(t *testing.T) {
	type args struct {
		rs map[string]*schema.Schema
	}
	tests := []struct {
		name string
		args args
		want map[string]*schema.Schema
	}{
		{
			name: "string",
			args: args{
				rs: map[string]*schema.Schema{
					"foo": {
						Type:        schema.TypeString,
						Required:    true,
						ForceNew:    true,
						Description: "foo of schema",
					},
				},
			},
			want: map[string]*schema.Schema{
				"foo": {
					Type:        schema.TypeString,
					Required:    false,
					ForceNew:    false,
					Computed:    true,
					Elem:        nil,
					Description: "foo of schema",
				},
			},
		},
		{
			name: "map",
			args: args{
				rs: map[string]*schema.Schema{
					"foo": {
						Type:        schema.TypeMap,
						Required:    true,
						ForceNew:    true,
						Description: "map of strings",
						Elem:        schema.TypeString,
					},
				},
			},
			want: map[string]*schema.Schema{
				"foo": {
					Type:        schema.TypeMap,
					Required:    false,
					ForceNew:    false,
					Computed:    true,
					Description: "map of strings",
					Elem:        schema.TypeString,
				},
			},
		},
		{
			name: "list_of_strings",
			args: args{
				rs: map[string]*schema.Schema{
					"foo": {
						Type:     schema.TypeList,
						Required: true,
						ForceNew: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				},
			},
			want: map[string]*schema.Schema{
				"foo": {
					Type:     schema.TypeList,
					Required: false,
					ForceNew: false,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
					MaxItems: 0,
					MinItems: 0,
				},
			},
		},
		{
			name: "list_subresource",
			args: args{
				rs: map[string]*schema.Schema{
					"foo": {
						Type:     schema.TypeList,
						Required: true,
						ForceNew: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"subresource": {
									Type:     schema.TypeList,
									Optional: true,
									Computed: true,
									MaxItems: 1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"disabled": {
												Type:     schema.TypeBool,
												Optional: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			want: map[string]*schema.Schema{
				"foo": {
					Type:     schema.TypeList,
					Required: false,
					ForceNew: false,
					Optional: false,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"subresource": {
								Type:     schema.TypeList,
								Optional: false,
								Computed: true,
								MaxItems: 0,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"disabled": {
											Type:     schema.TypeBool,
											Optional: false,
											Computed: true,
										},
									},
								},
							},
						},
					},
					MaxItems: 0,
					MinItems: 0,
				},
			},
		},
		{
			name: "set_of_strings",
			args: args{
				rs: map[string]*schema.Schema{
					"foo": {
						Type:     schema.TypeSet,
						Required: true,
						ForceNew: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
				},
			},
			want: map[string]*schema.Schema{
				"foo": {
					Type:     schema.TypeSet,
					Required: false,
					ForceNew: false,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
					MaxItems: 0,
					MinItems: 0,
				},
			},
		},
		{
			name: "set_subresource",
			args: args{
				rs: map[string]*schema.Schema{
					"foo": {
						Type:     schema.TypeSet,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"subresource": {
									Type:     schema.TypeInt,
									Optional: true,
									Computed: true,
									MaxItems: 1,
									Elem:     &schema.Schema{Type: schema.TypeInt},
								},
							},
						},
					},
				},
			},
			want: map[string]*schema.Schema{
				"foo": {
					Type:     schema.TypeSet,
					Required: false,
					ForceNew: false,
					Optional: false,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"subresource": {
								Type:     schema.TypeInt,
								Optional: false,
								Computed: true,
								MaxItems: 0,
								Elem:     &schema.Schema{Type: schema.TypeInt},
							},
						},
					},
					MaxItems: 0,
					MinItems: 0,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DatasourceSchemaFromResourceSchema(tt.args.rs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DatasourceSchemaFromResourceSchema() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestEmptyOrDefaultStringSuppress(t *testing.T) {
	testFunc := EmptyOrDefaultStringSuppress("default value")

	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"same value, format changed from empty to default": {
			Old:                "",
			New:                "default value",
			ExpectDiffSuppress: true,
		},
		"same value, format changed from default to empty": {
			Old:                "default value",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"different value, format changed from empty to non-default": {
			Old:                "",
			New:                "not default new",
			ExpectDiffSuppress: false,
		},
		"different value, format changed from non-default to empty": {
			Old:                "not default old",
			New:                "",
			ExpectDiffSuppress: false,
		},
		"different value, format changed from non-default to non-default": {
			Old:                "not default 1",
			New:                "not default 2",
			ExpectDiffSuppress: false,
		},
	}
	for tn, tc := range cases {
		if testFunc("", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Errorf("bad: %s, '%s' => '%s' expect DiffSuppress to return %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestServiceAccountFQN(t *testing.T) {
	// Every test case should produce this fully qualified service account name
	serviceAccountExpected := "projects/-/serviceAccounts/test-service-account@test-project.iam.gserviceaccount.com"
	cases := map[string]struct {
		serviceAccount string
		project        string
	}{
		"service account fully qualified name from account id": {
			serviceAccount: "test-service-account",
			project:        "test-project",
		},
		"service account fully qualified name from account email": {
			serviceAccount: "test-service-account@test-project.iam.gserviceaccount.com",
		},
		"service account fully qualified name from account name": {
			serviceAccount: "projects/-/serviceAccounts/test-service-account@test-project.iam.gserviceaccount.com",
		},
	}

	for tn, tc := range cases {
		config := &transport_tpg.Config{Project: tc.project}
		d := &schema.ResourceData{}
		serviceAccountName, err := ServiceAccountFQN(tc.serviceAccount, d, config)
		if err != nil {
			t.Fatalf("unexpected error for service account FQN: %s", err)
		}
		if serviceAccountName != serviceAccountExpected {
			t.Errorf("bad: %s, expected '%s' but returned '%s", tn, serviceAccountExpected, serviceAccountName)
		}
	}
}

func TestConflictError(t *testing.T) {
	confErr := &googleapi.Error{
		Code: 409,
	}
	if !IsConflictError(confErr) {
		t.Error("did not find that a 409 was a conflict error.")
	}
	if !IsConflictError(errwrap.Wrapf("wrap", confErr)) {
		t.Error("did not find that a wrapped 409 was a conflict error.")
	}
	confErr = &googleapi.Error{
		Code: 412,
	}
	if !IsConflictError(confErr) {
		t.Error("did not find that a 412 was a conflict error.")
	}
	if !IsConflictError(errwrap.Wrapf("wrap", confErr)) {
		t.Error("did not find that a wrapped 412 was a conflict error.")
	}
	// skipping negative tests as other cases may be added later.
}

func TestIsNotFoundGrpcErrort(t *testing.T) {
	error_status := status.New(codes.FailedPrecondition, "FailedPrecondition error")
	if IsNotFoundGrpcError(error_status.Err()) {
		t.Error("found FailedPrecondition as a NotFound error")
	}
	error_status = status.New(codes.OK, "OK")
	if IsNotFoundGrpcError(error_status.Err()) {
		t.Error("found OK as a NotFound error")
	}
	error_status = status.New(codes.NotFound, "NotFound error")
	if !IsNotFoundGrpcError(error_status.Err()) {
		t.Error("expect a NotFound error")
	}
}

func TestSnakeToPascalCase(t *testing.T) {
	input := "boot_disk"
	expected := "BootDisk"
	actual := SnakeToPascalCase(input)

	if actual != expected {
		t.Fatalf("(%s) did not match expected value: %s", actual, expected)
	}
}

func TestCheckGoogleIamPolicy(t *testing.T) {
	cases := []struct {
		valid bool
		json  string
	}{
		{
			valid: false,
			json:  `{"bindings":[{"condition":{"description":"","expression":"request.time \u003c timestamp(\"2020-01-01T00:00:00Z\")","title":"expires_after_2019_12_31-no-description"},"members":["user:admin@example.com"],"role":"roles/privateca.certificateManager"},{"condition":{"description":"Expiring at midnight of 2019-12-31","expression":"request.time \u003c timestamp(\"2020-01-01T00:00:00Z\")","title":"expires_after_2019_12_31"},"members":["user:admin@example.com"],"role":"roles/privateca.certificateManager"}]}`,
		},
		{
			valid: true,
			json:  `{"bindings":[{"condition":{"expression":"request.time \u003c timestamp(\"2020-01-01T00:00:00Z\")","title":"expires_after_2019_12_31-no-description"},"members":["user:admin@example.com"],"role":"roles/privateca.certificateManager"},{"condition":{"description":"Expiring at midnight of 2019-12-31","expression":"request.time \u003c timestamp(\"2020-01-01T00:00:00Z\")","title":"expires_after_2019_12_31"},"members":["user:admin@example.com"],"role":"roles/privateca.certificateManager"}]}`,
		},
	}

	for _, tc := range cases {
		err := CheckGoogleIamPolicy(tc.json)
		if tc.valid && err != nil {
			t.Errorf("The JSON is marked as valid but triggered an error: %s", tc.json)
		} else if !tc.valid && err == nil {
			t.Errorf("The JSON is marked as not valid but failed to trigger an error: %s", tc.json)
		}
	}
}

func TestReplaceVars(t *testing.T) {
	cases := map[string]struct {
		Template      string
		SchemaValues  map[string]interface{}
		Config        *transport_tpg.Config
		Expected      string
		ExpectedError bool
	}{
		"unspecified project fails": {
			Template:      "projects/{{project}}/global/images",
			ExpectedError: true,
		},
		"unspecified region fails": {
			Template: "projects/{{project}}/regions/{{region}}/subnetworks",
			Config: &transport_tpg.Config{
				Project: "default-project",
			},
			ExpectedError: true,
		},
		"unspecified zone fails": {
			Template: "projects/{{project}}/zones/{{zone}}/instances",
			Config: &transport_tpg.Config{
				Project: "default-project",
			},
			ExpectedError: true,
		},
		"regional with default values": {
			Template: "projects/{{project}}/regions/{{region}}/subnetworks",
			Config: &transport_tpg.Config{
				Project: "default-project",
				Region:  "default-region",
			},
			Expected: "projects/default-project/regions/default-region/subnetworks",
		},
		"zonal with default values": {
			Template: "projects/{{project}}/zones/{{zone}}/instances",
			Config: &transport_tpg.Config{
				Project: "default-project",
				Zone:    "default-zone",
			},
			Expected: "projects/default-project/zones/default-zone/instances",
		},
		"regional schema values": {
			Template: "projects/{{project}}/regions/{{region}}/subnetworks/{{name}}",
			SchemaValues: map[string]interface{}{
				"project": "project1",
				"region":  "region1",
				"name":    "subnetwork1",
			},
			Expected: "projects/project1/regions/region1/subnetworks/subnetwork1",
		},
		"regional schema self-link region": {
			Template: "projects/{{project}}/regions/{{region}}/subnetworks/{{name}}",
			SchemaValues: map[string]interface{}{
				"project": "project1",
				"region":  "https://www.googleapis.com/compute/v1/projects/project1/regions/region1",
				"name":    "subnetwork1",
			},
			Expected: "projects/project1/regions/region1/subnetworks/subnetwork1",
		},
		"zonal schema values": {
			Template: "projects/{{project}}/zones/{{zone}}/instances/{{name}}",
			SchemaValues: map[string]interface{}{
				"project": "project1",
				"zone":    "zone1",
				"name":    "instance1",
			},
			Expected: "projects/project1/zones/zone1/instances/instance1",
		},
		"zonal schema self-link zone": {
			Template: "projects/{{project}}/zones/{{zone}}/instances/{{name}}",
			SchemaValues: map[string]interface{}{
				"project": "project1",
				"zone":    "https://www.googleapis.com/compute/v1/projects/project1/zones/zone1",
				"name":    "instance1",
			},
			Expected: "projects/project1/zones/zone1/instances/instance1",
		},
		"zonal schema recursive replacement": {
			Template: "projects/{{project}}/zones/{{zone}}/instances/{{name}}",
			SchemaValues: map[string]interface{}{
				"project":   "project1",
				"zone":      "wrapper{{innerzone}}wrapper",
				"name":      "instance1",
				"innerzone": "inner",
			},
			Expected: "projects/project1/zones/wrapperinnerwrapper/instances/instance1",
		},
		"base path recursive replacement": {
			Template: "{{CloudRunBasePath}}namespaces/{{project}}/services",
			Config: &transport_tpg.Config{
				Project:          "default-project",
				Region:           "default-region",
				CloudRunBasePath: "https://{{region}}-run.googleapis.com/",
			},
			Expected: "https://default-region-run.googleapis.com/namespaces/default-project/services",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			d := &ResourceDataMock{
				FieldsInSchema: tc.SchemaValues,
			}

			config := tc.Config
			if config == nil {
				config = &transport_tpg.Config{}
			}

			v, err := ReplaceVars(d, config, tc.Template)

			if err != nil {
				if !tc.ExpectedError {
					t.Errorf("bad: %s; unexpected error %s", tn, err)
				}
				return
			}

			if tc.ExpectedError {
				t.Errorf("bad: %s; expected error", tn)
			}

			if v != tc.Expected {
				t.Errorf("bad: %s; expected %q, got %q", tn, tc.Expected, v)
			}
		})
	}
}
