// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package dataproc

import (
	"testing"
)

func TestDataprocExtractInitTimeout(t *testing.T) {
	t.Parallel()

	actual, err := extractInitTimeout("500s")
	expected := 500
	if err != nil {
		t.Fatalf("Expected %d, but got error %v", expected, err)
	}
	if actual != expected {
		t.Fatalf("Expected %d, but got %d", expected, actual)
	}
}

func TestDataprocExtractInitTimeout_nonSeconds(t *testing.T) {
	t.Parallel()

	actual, err := extractInitTimeout("5m")
	expected := 300
	if err != nil {
		t.Fatalf("Expected %d, but got error %v", expected, err)
	}
	if actual != expected {
		t.Fatalf("Expected %d, but got %d", expected, actual)
	}
}

func TestDataprocExtractInitTimeout_empty(t *testing.T) {
	t.Parallel()

	_, err := extractInitTimeout("")
	expected := "time: invalid duration"
	if err != nil && err.Error() != expected {
		return
	}
	t.Fatalf("Expected an error with message '%s', but got %v", expected, err.Error())
}

func TestDataprocParseImageVersion(t *testing.T) {
	t.Parallel()

	testCases := map[string]dataprocImageVersion{
		"1.2":             {"1", "2", "", ""},
		"1.2.3":           {"1", "2", "3", ""},
		"1.2.3rc":         {"1", "2", "3rc", ""},
		"1.2-debian9":     {"1", "2", "", "debian9"},
		"1.2.3-debian9":   {"1", "2", "3", "debian9"},
		"1.2.3rc-debian9": {"1", "2", "3rc", "debian9"},
	}

	for v, expected := range testCases {
		actual, err := parseDataprocImageVersion(v)
		if actual.major != expected.major {
			t.Errorf("parsing version %q returned error: %v", v, err)
		}
		if err != nil {
			t.Errorf("parsing version %q returned error: %v", v, err)
		}
		if actual.minor != expected.minor {
			t.Errorf("parsing version %q returned error: %v", v, err)
		}
		if actual.subminor != expected.subminor {
			t.Errorf("parsing version %q returned error: %v", v, err)
		}
		if actual.osName != expected.osName {
			t.Errorf("parsing version %q returned error: %v", v, err)
		}
	}

	errorTestCases := []string{
		"",
		"1",
		"notaversion",
		"1-debian",
	}
	for _, v := range errorTestCases {
		if _, err := parseDataprocImageVersion(v); err == nil {
			t.Errorf("expected parsing invalid version %q to return error", v)
		}
	}
}

func TestDataprocDiffSuppress(t *testing.T) {
	t.Parallel()

	doSuppress := [][]string{
		{"1.3.10-debian9", "1.3"},
		{"1.3.10-debian9", "1.3-debian9"},
		{"1.3.10", "1.3"},
		{"1.3-debian9", "1.3"},
	}

	noSuppress := [][]string{
		{"1.3.10-debian9", "1.3.10-ubuntu"},
		{"1.3.10-debian9", "1.3.9-debian9"},
		{"1.3.10-debian9", "1.3-ubuntu"},
		{"1.3.10-debian9", "1.3.9"},
		{"1.3.10-debian9", "1.4"},
		{"1.3.10-debian9", "2.3"},
		{"1.3.10", "1.3.10-debian9"},
		{"1.3", "1.3.10"},
		{"1.3", "1.3.10-debian9"},
		{"1.3", "1.3-debian9"},
	}

	for _, tup := range doSuppress {
		if !dataprocImageVersionDiffSuppress("", tup[0], tup[1], nil) {
			t.Errorf("expected (old: %q, new: %q) to be suppressed", tup[0], tup[1])
		}
	}
	for _, tup := range noSuppress {
		if dataprocImageVersionDiffSuppress("", tup[0], tup[1], nil) {
			t.Errorf("expected (old: %q, new: %q) to not be suppressed", tup[0], tup[1])
		}
	}
}
