// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sweeper

import (
	"reflect"
	"strings"
	"testing"
)

// TestUnifyRelationships verifies that parent relationships are correctly
// converted to reverse dependencies
func TestUnifyRelationships(t *testing.T) {
	testCases := []struct {
		name         string
		sweepers     map[string]*Sweeper
		expectedDeps map[string][]string
	}{
		{
			name: "simple_parent_relationship",
			sweepers: map[string]*Sweeper{
				"A": {
					Name:           "A",
					Dependencies:   []string{},
					DeleteFunction: func(region string) error { return nil },
				},
				"B": {
					Name:           "B",
					Parents:        []string{"A"},
					DeleteFunction: func(region string) error { return nil },
				},
			},
			expectedDeps: map[string][]string{
				"A": {"B"}, // A depends on B (parent depends on child)
				"B": {},
			},
		},
		{
			name: "chain_of_parents",
			sweepers: map[string]*Sweeper{
				"A": {
					Name:           "A",
					DeleteFunction: func(region string) error { return nil },
				},
				"B": {
					Name:           "B",
					Parents:        []string{"A"},
					DeleteFunction: func(region string) error { return nil },
				},
				"C": {
					Name:           "C",
					Parents:        []string{"B"},
					DeleteFunction: func(region string) error { return nil },
				},
			},
			expectedDeps: map[string][]string{
				"A": {"B"}, // A depends on B
				"B": {"C"}, // B depends on C
				"C": {},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			unified := unifyRelationships(tc.sweepers)

			// Verify unified dependencies
			for name, expectedDeps := range tc.expectedDeps {
				sweeper, ok := unified[name]
				if !ok {
					t.Fatalf("Sweeper %s missing from unified result", name)
				}

				// Use equal content comparison rather than exact slice equality
				if !hasSameDependencies(sweeper.Dependencies, expectedDeps) {
					t.Errorf("For sweeper %s, expected dependencies %v, got %v",
						name, expectedDeps, sweeper.Dependencies)
				}
			}
		})
	}
}

// hasSameDependencies checks if two slices have the same elements regardless of order
func hasSameDependencies(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	aMap := make(map[string]bool)
	for _, item := range a {
		aMap[item] = true
	}

	for _, item := range b {
		if !aMap[item] {
			return false
		}
	}

	return true
}

// TestCycleDetectionBasic verifies that basic cycles are correctly detected
func TestCycleDetectionBasic(t *testing.T) {
	testCases := []struct {
		name     string
		sweepers map[string]*Sweeper
		wantErr  bool
		errMsg   string
	}{
		{
			name: "direct_dependency_cycle",
			sweepers: map[string]*Sweeper{
				"A": {
					Name:           "A",
					Dependencies:   []string{"B"},
					DeleteFunction: func(region string) error { return nil },
				},
				"B": {
					Name:           "B",
					Dependencies:   []string{"A"},
					DeleteFunction: func(region string) error { return nil },
				},
			},
			wantErr: true,
			errMsg:  "dependency cycle detected",
		},
		{
			name: "self_dependency",
			sweepers: map[string]*Sweeper{
				"A": {
					Name:           "A",
					Dependencies:   []string{"A"},
					DeleteFunction: func(region string) error { return nil },
				},
			},
			wantErr: true,
			errMsg:  "dependency cycle detected",
		},
		{
			name: "no_cycle",
			sweepers: map[string]*Sweeper{
				"A": {
					Name:           "A",
					DeleteFunction: func(region string) error { return nil },
				},
				"B": {
					Name:           "B",
					Dependencies:   []string{"A"},
					DeleteFunction: func(region string) error { return nil },
				},
				"C": {
					Name:           "C",
					Dependencies:   []string{"B"},
					DeleteFunction: func(region string) error { return nil },
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// First unify relationships
			unified := unifyRelationships(tc.sweepers)

			// Then check for cycles
			err := detectCycles(unified)

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

// TestTopologicalSortWithUnifiedModel verifies that the topological sort
// produces the correct order with the unified model
func TestTopologicalSortWithUnifiedModel(t *testing.T) {
	tests := []struct {
		name     string
		sweepers map[string]*Sweeper
		expected []string
	}{
		{
			name: "simple_dependencies",
			sweepers: map[string]*Sweeper{
				"A": {
					Name:           "A",
					DeleteFunction: func(region string) error { return nil },
				},
				"B": {
					Name:           "B",
					Dependencies:   []string{"A"},
					DeleteFunction: func(region string) error { return nil },
				},
				"C": {
					Name:           "C",
					Dependencies:   []string{"B"},
					DeleteFunction: func(region string) error { return nil },
				},
			},
			expected: []string{"A", "B", "C"},
		},
		{
			name: "simple_parent_relationships",
			sweepers: map[string]*Sweeper{
				"A": {
					Name:           "A",
					DeleteFunction: func(region string) error { return nil },
				},
				"B": {
					Name:           "B",
					Parents:        []string{"A"},
					DeleteFunction: func(region string) error { return nil },
				},
				"C": {
					Name:           "C",
					Parents:        []string{"B"},
					DeleteFunction: func(region string) error { return nil },
				},
			},
			expected: []string{"C", "B", "A"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unify relationships
			unified := unifyRelationships(tt.sweepers)

			// Get sorted result
			sorted := topologicalSort(unified)

			// Extract names for comparison
			actual := make([]string, len(sorted))
			for i, s := range sorted {
				actual[i] = s.Name
			}

			// Check if the order is correct
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("Expected order %v, got %v", tt.expected, actual)
			}
		})
	}
}

// TestValidateAndOrderSweepers verifies the full pipeline
func TestValidateAndOrderSweepers(t *testing.T) {
	// Test valid parent-child relationships
	t.Run("parent_child_relationships", func(t *testing.T) {
		sweepers := map[string]*Sweeper{
			"A": {
				Name:           "A",
				DeleteFunction: func(region string) error { return nil },
			},
			"B": {
				Name:           "B",
				Parents:        []string{"A"},
				DeleteFunction: func(region string) error { return nil },
			},
			"C": {
				Name:           "C",
				Parents:        []string{"B"},
				DeleteFunction: func(region string) error { return nil },
			},
		}

		sorted, err := validateAndOrderSweepersWithDependencies(sweepers)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Extract names for comparison
		names := make([]string, len(sorted))
		for i, s := range sorted {
			names[i] = s.Name
		}

		// Expected order: children before parents
		expected := []string{"C", "B", "A"}
		if !reflect.DeepEqual(names, expected) {
			t.Errorf("Expected order %v, got %v", expected, names)
		}
	})

	// Test direct cycle
	t.Run("direct_cycle", func(t *testing.T) {
		sweepers := map[string]*Sweeper{
			"A": {
				Name:           "A",
				Dependencies:   []string{"B"},
				DeleteFunction: func(region string) error { return nil },
			},
			"B": {
				Name:           "B",
				Dependencies:   []string{"A"},
				DeleteFunction: func(region string) error { return nil },
			},
		}

		_, err := validateAndOrderSweepersWithDependencies(sweepers)
		if err == nil {
			t.Error("Expected error for cycle, got nil")
		} else if !strings.Contains(err.Error(), "dependency cycle detected") {
			t.Errorf("Expected error containing 'dependency cycle detected', got: %v", err)
		}
	})
}

// TestFilterSweepers verifies that filtering sweepers includes dependencies
func TestFilterSweepers(t *testing.T) {
	sweepers := map[string]*Sweeper{
		"resource_a": {
			Name:           "resource_a",
			DeleteFunction: func(region string) error { return nil },
		},
		"resource_b": {
			Name:           "resource_b",
			Dependencies:   []string{"resource_a"},
			DeleteFunction: func(region string) error { return nil },
		},
		"resource_c": {
			Name:           "resource_c",
			DeleteFunction: func(region string) error { return nil },
		},
	}

	// Filter on resource_b should include resource_a
	filtered := filterSweepers("resource_b", sweepers)
	if _, ok := filtered["resource_a"]; !ok {
		t.Error("Filtering for resource_b should include its dependency resource_a")
	}
	if _, ok := filtered["resource_c"]; ok {
		t.Error("Filtering for resource_b should not include unrelated resource_c")
	}
}
