package google

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/googleapi"
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
		if ipCidrRangeDiffSuppress("ip_cidr_range", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
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
		if rfc3339TimeDiffSuppress("time", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Errorf("bad: %s, '%s' => '%s' expect DiffSuppress to return %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestGetZone(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceComputeDisk().Schema, map[string]interface{}{
		"zone": "foo",
	})
	var config Config
	if err := d.Set("zone", "foo"); err != nil {
		t.Fatalf("Cannot set zone: %s", err)
	}
	if zone, err := getZone(d, &config); err != nil || zone != "foo" {
		t.Fatalf("Zone '%s' != 'foo', %s", zone, err)
	}
	config.Zone = "bar"
	if zone, err := getZone(d, &config); err != nil || zone != "foo" {
		t.Fatalf("Zone '%s' != 'foo', %s", zone, err)
	}
	d.Set("zone", "")
	if zone, err := getZone(d, &config); err != nil || zone != "bar" {
		t.Fatalf("Zone '%s' != 'bar', %s", zone, err)
	}
	config.Zone = ""
	if zone, err := getZone(d, &config); err == nil || zone != "" {
		t.Fatalf("Zone '%s' != '', err=%s", zone, err)
	}
}

func TestGetRegion(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceComputeDisk().Schema, map[string]interface{}{
		"zone": "foo",
	})
	var config Config
	barRegionName := getRegionFromZone("bar")
	fooRegionName := getRegionFromZone("foo")

	if region, err := getRegion(d, &config); err != nil || region != fooRegionName {
		t.Fatalf("Zone '%s' != '%s', %s", region, fooRegionName, err)
	}

	config.Zone = "bar"
	d.Set("zone", "")
	if region, err := getRegion(d, &config); err != nil || region != barRegionName {
		t.Fatalf("Zone '%s' != '%s', %s", region, barRegionName, err)
	}
	config.Region = "something-else"
	if region, err := getRegion(d, &config); err != nil || region != config.Region {
		t.Fatalf("Zone '%s' != '%s', %s", region, config.Region, err)
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
	testFunc := emptyOrDefaultStringSuppress("default value")

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
		config := &Config{Project: tc.project}
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
	if err := retryTimeDuration(f, time.Duration(1000)*time.Millisecond); err == nil || err.(*googleapi.Error).Code != 500 {
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
	if err := retryTimeDuration(f, time.Duration(1000)*time.Millisecond); err == nil {
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
	if err := retryTimeDuration(f, time.Duration(1000)*time.Millisecond); err == nil || err.(*googleapi.Error).Code != 400 {
		t.Errorf("unexpected error retrying: %v", err)
	}
	if i != 1 {
		t.Errorf("expected error function to be called exactly once, but was called %d times", i)
	}
}

type TimeoutError struct {
	timeout bool
}

func (e *TimeoutError) Timeout() bool {
	return e.timeout
}

func (e *TimeoutError) Error() string {
	return "timeout error"
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
	err := retryTimeDuration(retryFunc, 1*time.Minute)
	if err != nil {
		t.Errorf("unexpected error: got '%v' want 'nil'", err)
	}
	expectedRunCount := 2
	if runCount != expectedRunCount {
		t.Errorf("expected the retryFunc to be called %v time(s), instead was called %v time(s)", expectedRunCount, runCount)
	}
}
