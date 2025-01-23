// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sweeper

import (
	"reflect"
	"strings"
	"testing"
)

func TestValidateAndOrderSweepers_Simple(t *testing.T) {
	sweepers := map[string]*Sweeper{
		"B": {
			Name:         "B",
			Dependencies: []string{"A"},
		},
		"C": {
			Name:         "C",
			Dependencies: []string{"B"},
		},
		"A": {
			Name:         "A",
			Dependencies: []string{},
		},
	}

	sorted, err := validateAndOrderSweepers(sweepers)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify order: A should come before B, B before C
	names := make([]string, len(sorted))
	for i, s := range sorted {
		names[i] = s.Name
	}

	expected := []string{"A", "B", "C"}
	if !reflect.DeepEqual(names, expected) {
		t.Errorf("Expected order %v, got %v", expected, names)
	}
}

func TestValidateAndOrderSweepers_Cycle(t *testing.T) {
	testCases := []struct {
		name     string
		sweepers map[string]*Sweeper
		wantErr  bool
		errMsg   string
	}{
		{
			name: "direct_cycle",
			sweepers: map[string]*Sweeper{
				"A": {
					Name:         "A",
					Dependencies: []string{"B"},
				},
				"B": {
					Name:         "B",
					Dependencies: []string{"A"},
				},
			},
			wantErr: true,
			errMsg:  "dependency cycle detected",
		},
		{
			name: "indirect_cycle",
			sweepers: map[string]*Sweeper{
				"A": {
					Name:         "A",
					Dependencies: []string{"B"},
				},
				"B": {
					Name:         "B",
					Dependencies: []string{"C"},
				},
				"C": {
					Name:         "C",
					Dependencies: []string{"A"},
				},
			},
			wantErr: true,
			errMsg:  "dependency cycle detected",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := validateAndOrderSweepers(tc.sweepers)
			if tc.wantErr {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if !strings.Contains(err.Error(), tc.errMsg) {
					t.Errorf("Expected error containing %q, got %v", tc.errMsg, err)
				}
			} else if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
		})
	}
}

func TestValidateAndOrderSweepers_MissingDependency(t *testing.T) {
	sweepers := map[string]*Sweeper{
		"A": {
			Name:         "A",
			Dependencies: []string{"NonExistent"},
		},
	}

	_, err := validateAndOrderSweepers(sweepers)
	if err == nil {
		t.Fatal("Expected error for missing dependency, got nil")
	}
	expected := "sweeper A depends on NonExistent, but NonExistent not found"
	if err.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, err.Error())
	}
}

func TestValidateAndOrderSweepers_Complex(t *testing.T) {
	sweepers := map[string]*Sweeper{
		"A": {
			Name:         "A",
			Dependencies: []string{},
		},
		"B": {
			Name:         "B",
			Dependencies: []string{"A"},
		},
		"C": {
			Name:         "C",
			Dependencies: []string{"A"},
		},
		"D": {
			Name:         "D",
			Dependencies: []string{"B", "C"},
		},
		"E": {
			Name:         "E",
			Dependencies: []string{"C"},
		},
	}

	sorted, err := validateAndOrderSweepers(sweepers)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Create a map to check relative positions
	positions := make(map[string]int)
	for i, s := range sorted {
		positions[s.Name] = i
	}

	// Verify dependencies come before dependents
	checks := []struct {
		dependency string
		dependent  string
	}{
		{"A", "B"},
		{"A", "C"},
		{"B", "D"},
		{"C", "D"},
		{"C", "E"},
	}

	for _, check := range checks {
		if positions[check.dependency] >= positions[check.dependent] {
			t.Errorf("Expected %s to come before %s, but got positions %d and %d",
				check.dependency, check.dependent,
				positions[check.dependency], positions[check.dependent])
		}
	}
}

func TestValidateAndOrderSweepers_Empty(t *testing.T) {
	sweepers := map[string]*Sweeper{}

	sorted, err := validateAndOrderSweepers(sweepers)
	if err != nil {
		t.Fatalf("Expected no error for empty sweepers, got: %v", err)
	}
	if len(sorted) != 0 {
		t.Errorf("Expected empty result for empty input, got %d items", len(sorted))
	}
}

func TestValidateAndOrderSweepers_SelfDependency(t *testing.T) {
	sweepers := map[string]*Sweeper{
		"A": {
			Name:         "A",
			Dependencies: []string{"A"},
		},
	}

	_, err := validateAndOrderSweepers(sweepers)
	if err == nil {
		t.Fatal("Expected error for self-dependency, got nil")
	}
	if !strings.Contains(err.Error(), "dependency cycle detected") {
		t.Errorf("Expected cycle detection error, got: %v", err)
	}
}

func TestFilterSweepers(t *testing.T) {
	testCases := []struct {
		name           string
		filter         string
		sourceSweepers map[string]*Sweeper
		expected       map[string]*Sweeper
	}{
		{
			name:   "empty_filter",
			filter: "",
			sourceSweepers: map[string]*Sweeper{
				"test": {Name: "test"},
				"prod": {Name: "prod"},
			},
			expected: map[string]*Sweeper{
				"test": {Name: "test"},
				"prod": {Name: "prod"},
			},
		},
		{
			name:   "single_match",
			filter: "test",
			sourceSweepers: map[string]*Sweeper{
				"test":    {Name: "test"},
				"testing": {Name: "testing"},
				"prod":    {Name: "prod"},
			},
			expected: map[string]*Sweeper{
				"test":    {Name: "test"},
				"testing": {Name: "testing"},
			},
		},
		{
			name:   "multiple_filters",
			filter: "test,prod",
			sourceSweepers: map[string]*Sweeper{
				"test":    {Name: "test"},
				"testing": {Name: "testing"},
				"prod":    {Name: "prod"},
				"stage":   {Name: "stage"},
			},
			expected: map[string]*Sweeper{
				"test":    {Name: "test"},
				"testing": {Name: "testing"},
				"prod":    {Name: "prod"},
			},
		},
		{
			name:   "case_insensitive",
			filter: "TEST",
			sourceSweepers: map[string]*Sweeper{
				"test":    {Name: "test"},
				"testing": {Name: "testing"},
			},
			expected: map[string]*Sweeper{
				"test":    {Name: "test"},
				"testing": {Name: "testing"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filterSweepers(tc.filter, tc.sourceSweepers)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestFilterSweeperWithDependencies(t *testing.T) {
	testCases := []struct {
		name           string
		sweeper        string
		sourceSweepers map[string]*Sweeper
		expected       map[string]*Sweeper
	}{
		{
			name:    "no_dependencies",
			sweeper: "test",
			sourceSweepers: map[string]*Sweeper{
				"test": {
					Name:         "test",
					Dependencies: []string{},
				},
			},
			expected: map[string]*Sweeper{
				"test": {
					Name:         "test",
					Dependencies: []string{},
				},
			},
		},
		{
			name:    "with_dependencies",
			sweeper: "test",
			sourceSweepers: map[string]*Sweeper{
				"test": {
					Name:         "test",
					Dependencies: []string{"dep1", "dep2"},
				},
				"dep1": {
					Name:         "dep1",
					Dependencies: []string{},
				},
				"dep2": {
					Name:         "dep2",
					Dependencies: []string{},
				},
			},
			expected: map[string]*Sweeper{
				"test": {
					Name:         "test",
					Dependencies: []string{"dep1", "dep2"},
				},
				"dep1": {
					Name:         "dep1",
					Dependencies: []string{},
				},
				"dep2": {
					Name:         "dep2",
					Dependencies: []string{},
				},
			},
		},
		{
			name:    "missing_dependency",
			sweeper: "test",
			sourceSweepers: map[string]*Sweeper{
				"test": {
					Name:         "test",
					Dependencies: []string{"missing"},
				},
			},
			expected: map[string]*Sweeper{
				"test": {
					Name:         "test",
					Dependencies: []string{"missing"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := filterSweeperWithDependencies(tc.sweeper, tc.sourceSweepers)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}
