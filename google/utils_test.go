package google

import (
	"reflect"
	"strings"
	"testing"
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

func TestExtractLastResourceFromUri_withUrl(t *testing.T) {
	actual := extractLastResourceFromUri("http://something.com/one/two/three")
	expected := "three"
	if actual != expected {
		t.Fatalf("Expected %s, but got %s", expected, actual)
	}
}

func TestExtractLastResourceFromUri_WithStaticValue(t *testing.T) {
	actual := extractLastResourceFromUri("three")
	expected := "three"
	if actual != expected {
		t.Fatalf("Expected %s, but got %s", expected, actual)
	}
}

func TestIpCidrRangeDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New          string
		ExpectDiffSupress bool
	}{
		"single ip address": {
			Old:               "10.2.3.4",
			New:               "10.2.3.5",
			ExpectDiffSupress: false,
		},
		"cidr format string": {
			Old:               "10.1.2.0/24",
			New:               "10.1.3.0/24",
			ExpectDiffSupress: false,
		},
		"netmask same mask": {
			Old:               "10.1.2.0/24",
			New:               "/24",
			ExpectDiffSupress: true,
		},
		"netmask different mask": {
			Old:               "10.1.2.0/24",
			New:               "/32",
			ExpectDiffSupress: false,
		},
		"add netmask": {
			Old:               "",
			New:               "/24",
			ExpectDiffSupress: false,
		},
		"remove netmask": {
			Old:               "/24",
			New:               "",
			ExpectDiffSupress: false,
		},
	}

	for tn, tc := range cases {
		if ipCidrRangeDiffSuppress("ip_cidr_range", tc.Old, tc.New, nil) != tc.ExpectDiffSupress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSupress)
		}
	}
}

func TestRfc3339TimeDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New          string
		ExpectDiffSupress bool
	}{
		"same time, format changed to have leading zero": {
			Old:               "2:00",
			New:               "02:00",
			ExpectDiffSupress: true,
		},
		"same time, format changed not to have leading zero": {
			Old:               "02:00",
			New:               "2:00",
			ExpectDiffSupress: true,
		},
		"different time, both without leading zero": {
			Old:               "2:00",
			New:               "3:00",
			ExpectDiffSupress: false,
		},
		"different time, old with leading zero, new without": {
			Old:               "02:00",
			New:               "3:00",
			ExpectDiffSupress: false,
		},
		"different time, new with leading zero, oldwithout": {
			Old:               "2:00",
			New:               "03:00",
			ExpectDiffSupress: false,
		},
		"different time, both with leading zero": {
			Old:               "02:00",
			New:               "03:00",
			ExpectDiffSupress: false,
		},
	}
	for tn, tc := range cases {
		if rfc3339TimeDiffSuppress("time", tc.Old, tc.New, nil) != tc.ExpectDiffSupress {
			t.Errorf("bad: %s, '%s' => '%s' expect DiffSuppress to return %t", tn, tc.Old, tc.New, tc.ExpectDiffSupress)
		}
	}
}
