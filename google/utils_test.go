package google

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"

	"google.golang.org/api/googleapi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestConvertStringArr(t *testing.T) {
	input := make([]interface{}, 3)
	input[0] = "aaa"
	input[1] = "bbb"
	input[2] = "aaa"

	expected := []string{"aaa", "bbb", "ccc"}
	actual := convertStringArr(input)

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
	actual := convertAndMapStringArr(input, strings.ToUpper)

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
	actual := convertStringMap(input)

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
		if tpgresource.IpCidrRangeDiffSuppress("ip_cidr_range", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
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
		if tpgresource.Rfc3339TimeDiffSuppress("time", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Errorf("bad: %s, '%s' => '%s' expect DiffSuppress to return %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestGetProject(t *testing.T) {
	cases := map[string]struct {
		ResourceProject string
		ProviderProject string
		ExpectedProject string
		ExpectedError   bool
	}{
		"project is pulled from resource config instead of provider config": {
			ResourceProject: "foo",
			ProviderProject: "bar",
			ExpectedProject: "foo",
		},
		"project is pulled from provider config when not set on resource": {
			ProviderProject: "bar",
			ExpectedProject: "bar",
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
			if tc.ProviderProject != "" {
				config.Project = tc.ProviderProject
			}

			// Create resource config
			// Here use ResourceComputeDisk schema as example
			emptyConfigMap := map[string]interface{}{}
			d := schema.TestResourceDataRaw(t, ResourceComputeDisk().Schema, emptyConfigMap)
			if tc.ResourceProject != "" {
				if err := d.Set("project", tc.ResourceProject); err != nil {
					t.Fatalf("Cannot set project: %s", err)
				}
			}

			// Act
			project, err := getProject(d, &config)

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
		ResourceConfig   map[string]string
		ProviderConfig   map[string]string
		ExpectedLocation string
		ExpectError      bool
	}{
		"returns the location value set in the resource config": {
			ResourceConfig: map[string]string{
				"location": "resource-location",
			},
			ExpectedLocation: "resource-location",
		},
		"returned location values set as self links are not shortened": {
			ResourceConfig: map[string]string{
				"location": "https://www.googleapis.com/compute/v1/projects/my-project/locations/resource-location",
			},
			ExpectedLocation: "https://www.googleapis.com/compute/v1/projects/my-project/locations/resource-location", // No shortening takes place
		},
		"returns the region value set in the resource config when location is not in the schema": {
			ResourceConfig: map[string]string{
				"region": "resource-region",
			},
			ExpectedLocation: "resource-region",
		},
		"returned region values set as self links are not shortened": {
			ResourceConfig: map[string]string{
				"region": "https://www.googleapis.com/compute/v1/projects/my-project/region/resource-region",
			},
			ExpectedLocation: "https://www.googleapis.com/compute/v1/projects/my-project/region/resource-region", // No shortening takes place
		},
		"returns the zone value set in the resource config when neither location nor region in the schema": {
			ResourceConfig: map[string]string{
				"zone": "resource-zone",
			},
			ExpectedLocation: "resource-zone",
		},
		"returned zone values set as self links in the resource config ARE shortened": {
			// Results from getLocation using getZone internally
			// This behaviour makes sense because APIs may return a self link as the zone value
			ResourceConfig: map[string]string{
				"zone": "https://www.googleapis.com/compute/v1/projects/my-project/zones/resource-zone",
			},
			ExpectedLocation: "resource-zone",
		},
		"returns the zone value from the provider config when none of location/region/zone are set in the resource config": {
			ProviderConfig: map[string]string{
				"zone": "provider-zone",
			},
			ExpectedLocation: "provider-zone",
		},
		"returned zone values set as self links in the provider config are NOT shortened": {
			// This behaviour makes sense because provider config values don't originate from APIs
			// Users should always configure the provider with the short names of regions/zones
			ProviderConfig: map[string]string{
				"zone": "https://www.googleapis.com/compute/v1/projects/my-project/zones/provider-zone",
			},
			ExpectedLocation: "https://www.googleapis.com/compute/v1/projects/my-project/zones/provider-zone",
		},
		"returns an error when only a region value is set in the the provider config and none of location/region/zone are set in the resource config": {
			ProviderConfig: map[string]string{
				"region": "provider-region",
			},
			ExpectError: true,
		},
		"an error is returned when none of location/region/zone are set on the resource, and neither region or zone is set on the provider": {
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
			fictionalSchema := map[string]*schema.Schema{
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
			}
			emptyConfigMap := map[string]interface{}{}
			d := schema.TestResourceDataRaw(t, fictionalSchema, emptyConfigMap)

			// Load Terraform resource config data
			if len(tc.ResourceConfig) > 0 {
				for k, v := range tc.ResourceConfig {
					err := d.Set(k, v)
					if err != nil {
						t.Fatalf("error during test setup: %v", err)
					}
				}
			}

			// Act
			location, err := getLocation(d, &config)

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
		ResourceZone  string
		ProviderZone  string
		ExpectedZone  string
		ExpectedError bool
	}{
		"zone is pulled from resource config instead of provider config": {
			ResourceZone: "foo",
			ProviderZone: "bar",
			ExpectedZone: "foo",
		},
		"zone value from resource can be a self link": {
			ResourceZone: "https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-a",
			ExpectedZone: "us-central1-a",
		},
		"zone is pulled from provider config when not set on resource": {
			ProviderZone: "bar",
			ExpectedZone: "bar",
		},
		"error returned when zone not set on either provider or resource": {
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange

			// Create provider config
			var config transport_tpg.Config
			if tc.ProviderZone != "" {
				config.Zone = tc.ProviderZone
			}

			// Create resource config
			// Here use ResourceComputeDisk schema as example - because it has a zone field in schema
			emptyConfigMap := map[string]interface{}{}
			d := schema.TestResourceDataRaw(t, ResourceComputeDisk().Schema, emptyConfigMap)
			if tc.ResourceZone != "" {
				if err := d.Set("zone", tc.ResourceZone); err != nil {
					t.Fatalf("Cannot set zone: %s", err)
				}
			}

			// Act
			zone, err := getZone(d, &config)

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
		ResourceRegion string
		ProviderRegion string
		ProviderZone   string
		ExpectedRegion string
		ExpectedZone   string
		ExpectedError  bool
	}{
		"region is pulled from resource config instead of provider config": {
			ResourceRegion: "foo",
			ProviderRegion: "bar",
			ProviderZone:   "lol-a",
			ExpectedRegion: "foo",
		},
		"region pulled from resource config can be a self link": {
			ResourceRegion: "https://www.googleapis.com/compute/v1/projects/my-project/regions/us-central1",
			ExpectedRegion: "us-central1",
		},
		"region is pulled from region on provider config when region unset in resource config": {
			ProviderRegion: "bar",
			ProviderZone:   "lol-a",
			ExpectedRegion: "bar",
		},
		"region is pulled from zone on provider config when region unset in both resource and provider config": {
			ProviderZone:   "lol-a",
			ExpectedRegion: "lol",
		},
		"error returned when region not set on resource and neither region or zone set on provider": {
			ExpectedError: true,
		},
	}
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Arrange

			// Create provider config
			var config transport_tpg.Config
			if tc.ProviderRegion != "" {
				config.Region = tc.ProviderRegion
			}
			if tc.ProviderZone != "" {
				config.Zone = tc.ProviderZone
			}

			// Create resource config
			// Here use ResourceComputeSubnetwork schema as example - because it has a region field in schema
			emptyConfigMap := map[string]interface{}{}
			d := schema.TestResourceDataRaw(t, ResourceComputeSubnetwork().Schema, emptyConfigMap)
			if tc.ResourceRegion != "" {
				if err := d.Set("region", tc.ResourceRegion); err != nil {
					t.Fatalf("Cannot set region: %s", err)
				}
			}

			// Act
			region, err := getRegion(d, &config)

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
	actual := getRegionFromZone("us-central1-f")
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
			if got := datasourceSchemaFromResourceSchema(tt.args.rs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("datasourceSchemaFromResourceSchema() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestEmptyOrDefaultStringSuppress(t *testing.T) {
	testFunc := tpgresource.EmptyOrDefaultStringSuppress("default value")

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
		serviceAccountName, err := serviceAccountFQN(tc.serviceAccount, d, config)
		if err != nil {
			t.Fatalf("unexpected error for service account FQN: %s", err)
		}
		if serviceAccountName != serviceAccountExpected {
			t.Errorf("bad: %s, expected '%s' but returned '%s", tn, serviceAccountExpected, serviceAccountName)
		}
	}
}

func TestRetryTimeDuration(t *testing.T) {
	i := 0
	f := func() error {
		i++
		return &googleapi.Error{
			Code: 500,
		}
	}
	if err := transport_tpg.RetryTimeDuration(f, time.Duration(1000)*time.Millisecond); err == nil || err.(*googleapi.Error).Code != 500 {
		t.Errorf("unexpected error retrying: %v", err)
	}
	if i < 2 {
		t.Errorf("expected error function to be called at least twice, but was called %d times", i)
	}
}

func TestRetryTimeDuration_wrapped(t *testing.T) {
	i := 0
	f := func() error {
		i++
		err := &googleapi.Error{
			Code: 500,
		}
		return errwrap.Wrapf("nested error: {{err}}", err)
	}
	if err := transport_tpg.RetryTimeDuration(f, time.Duration(1000)*time.Millisecond); err == nil {
		t.Errorf("unexpected nil error, expected an error")
	} else {
		innerErr := errwrap.GetType(err, &googleapi.Error{})
		if innerErr == nil {
			t.Errorf("unexpected error %v does not have a google api error", err)
		}
		gerr := innerErr.(*googleapi.Error)
		if gerr.Code != 500 {
			t.Errorf("unexpected googleapi error expected code 500, error: %v", gerr)
		}
	}
	if i < 2 {
		t.Errorf("expected error function to be called at least twice, but was called %d times", i)
	}
}

func TestRetryTimeDuration_noretry(t *testing.T) {
	i := 0
	f := func() error {
		i++
		return &googleapi.Error{
			Code: 400,
		}
	}
	if err := transport_tpg.RetryTimeDuration(f, time.Duration(1000)*time.Millisecond); err == nil || err.(*googleapi.Error).Code != 400 {
		t.Errorf("unexpected error retrying: %v", err)
	}
	if i != 1 {
		t.Errorf("expected error function to be called exactly once, but was called %d times", i)
	}
}

func TestRetryTimeDuration_URLTimeoutsShouldRetry(t *testing.T) {
	runCount := 0
	retryFunc := func() error {
		runCount++
		if runCount == 1 {
			return &url.Error{
				Err: &TimeoutError{timeout: true},
			}
		}
		return nil
	}
	err := transport_tpg.RetryTimeDuration(retryFunc, 1*time.Minute)
	if err != nil {
		t.Errorf("unexpected error: got '%v' want 'nil'", err)
	}
	expectedRunCount := 2
	if runCount != expectedRunCount {
		t.Errorf("expected the retryFunc to be called %v time(s), instead was called %v time(s)", expectedRunCount, runCount)
	}
}

func TestRetryWithPolling_noRetry(t *testing.T) {
	retryCount := 0
	retryFunc := func() (interface{}, error) {
		retryCount++
		return "", &googleapi.Error{
			Code: 400,
		}
	}
	result, err := transport_tpg.RetryWithPolling(retryFunc, time.Duration(1000)*time.Millisecond, time.Duration(100)*time.Millisecond)
	if err == nil || err.(*googleapi.Error).Code != 400 || result.(string) != "" {
		t.Errorf("unexpected error %v and result %v", err, result)
	}
	if retryCount != 1 {
		t.Errorf("expected error function to be called exactly once, but was called %d times", retryCount)
	}
}

func TestRetryWithPolling_notRetryable(t *testing.T) {
	retryCount := 0
	retryFunc := func() (interface{}, error) {
		retryCount++
		return "", &googleapi.Error{
			Code: 400,
		}
	}
	// Retryable if the error code is not 400.
	isRetryableFunc := func(err error) (bool, string) {
		return err.(*googleapi.Error).Code != 400, ""
	}
	result, err := transport_tpg.RetryWithPolling(retryFunc, time.Duration(1000)*time.Millisecond, time.Duration(100)*time.Millisecond, isRetryableFunc)
	if err == nil || err.(*googleapi.Error).Code != 400 || result.(string) != "" {
		t.Errorf("unexpected error %v and result %v", err, result)
	}
	if retryCount != 1 {
		t.Errorf("expected error function to be called exactly once, but was called %d times", retryCount)
	}
}

func TestRetryWithPolling_retriedAndSucceeded(t *testing.T) {
	retryCount := 0
	// Retry once and succeeds.
	retryFunc := func() (interface{}, error) {
		retryCount++
		// Error code of 200 is retryable.
		if retryCount < 2 {
			return "", &googleapi.Error{
				Code: 200,
			}
		}
		return "Ok", nil
	}
	// Retryable if the error code is not 400.
	isRetryableFunc := func(err error) (bool, string) {
		return err.(*googleapi.Error).Code != 400, ""
	}
	result, err := transport_tpg.RetryWithPolling(retryFunc, time.Duration(1000)*time.Millisecond, time.Duration(100)*time.Millisecond, isRetryableFunc)
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if result.(string) != "Ok" {
		t.Errorf("unexpected result %v", result)
	}
	if retryCount != 2 {
		t.Errorf("expected error function to be called exactly twice, but was called %d times", retryCount)
	}
}

func TestRetryWithPolling_retriedAndFailed(t *testing.T) {
	retryCount := 0
	// Retry once and fails.
	retryFunc := func() (interface{}, error) {
		retryCount++
		// Error code of 200 is retryable.
		if retryCount < 2 {
			return "", &googleapi.Error{
				Code: 200,
			}
		}
		return "", &googleapi.Error{
			Code: 400,
		}
	}
	// Retryable if the error code is not 400.
	isRetryableFunc := func(err error) (bool, string) {
		return err.(*googleapi.Error).Code != 400, ""
	}
	result, err := transport_tpg.RetryWithPolling(retryFunc, time.Duration(1000)*time.Millisecond, time.Duration(100)*time.Millisecond, isRetryableFunc)
	if err == nil || err.(*googleapi.Error).Code != 400 || result.(string) != "" {
		t.Errorf("unexpected error %v and result %v", err, result)
	}
	if retryCount != 2 {
		t.Errorf("expected error function to be called exactly twice, but was called %d times", retryCount)
	}
}

func TestConflictError(t *testing.T) {
	confErr := &googleapi.Error{
		Code: 409,
	}
	if !isConflictError(confErr) {
		t.Error("did not find that a 409 was a conflict error.")
	}
	if !isConflictError(errwrap.Wrapf("wrap", confErr)) {
		t.Error("did not find that a wrapped 409 was a conflict error.")
	}
	confErr = &googleapi.Error{
		Code: 412,
	}
	if !isConflictError(confErr) {
		t.Error("did not find that a 412 was a conflict error.")
	}
	if !isConflictError(errwrap.Wrapf("wrap", confErr)) {
		t.Error("did not find that a wrapped 412 was a conflict error.")
	}
	// skipping negative tests as other cases may be added later.
}

func TestIsNotFoundGrpcErrort(t *testing.T) {
	error_status := status.New(codes.FailedPrecondition, "FailedPrecondition error")
	if isNotFoundGrpcError(error_status.Err()) {
		t.Error("found FailedPrecondition as a NotFound error")
	}
	error_status = status.New(codes.OK, "OK")
	if isNotFoundGrpcError(error_status.Err()) {
		t.Error("found OK as a NotFound error")
	}
	error_status = status.New(codes.NotFound, "NotFound error")
	if !isNotFoundGrpcError(error_status.Err()) {
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
		err := checkGCSName(bucketName)
		if valid && err != nil {
			t.Errorf("The bucket name %s was expected to pass validation and did not pass.", bucketName)
		} else if !valid && err == nil {
			t.Errorf("The bucket name %s was NOT expected to pass validation and passed.", bucketName)
		}
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
		err := checkGoogleIamPolicy(tc.json)
		if tc.valid && err != nil {
			t.Errorf("The JSON is marked as valid but triggered an error: %s", tc.json)
		} else if !tc.valid && err == nil {
			t.Errorf("The JSON is marked as not valid but failed to trigger an error: %s", tc.json)
		}
	}
}
