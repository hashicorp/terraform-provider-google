// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sweeper

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Sweeper struct in the sweeper package
type Sweeper struct {
	// Name for sweeper. Must be unique
	Name string

	// Parents list of parent resource names that must be swept after this resource
	Parents []string

	// Dependencies list of resources that must be swept before this resource
	Dependencies []string

	// List function that can list this resource type
	ListAndAction SweeperListFunc

	DeleteFunction func(region string) error
}

// SweeperListFunc defines the signature for resource list functions
type SweeperListFunc func(ResourceAction) error
type ResourceAction func(*transport_tpg.Config, *tpgresource.ResourceDataMock, map[string]interface{}) error

var (
	flagSweep              *string
	flagSweepAllowFailures *bool
	flagSweepRun           *string
	sweeperInventory       map[string]*Sweeper
)

func init() {
	sweeperInventory = make(map[string]*Sweeper)
}

// registerFlags checks for and gets existing flag definitions before trying to redefine them.
func registerFlags() {
	// Check for existing flags in global CommandLine
	if f := flag.Lookup("sweep"); f != nil {
		// Use the Value.Get() interface to get the values
		if getter, ok := f.Value.(flag.Getter); ok {
			vs := getter.Get().(string)
			flagSweep = &vs
		}
		if f := flag.Lookup("sweep-allow-failures"); f != nil {
			if getter, ok := f.Value.(flag.Getter); ok {
				vb := getter.Get().(bool)
				flagSweepAllowFailures = &vb
			}
		}
		if f := flag.Lookup("sweep-run"); f != nil {
			if getter, ok := f.Value.(flag.Getter); ok {
				vs := getter.Get().(string)
				flagSweepRun = &vs
			}
		}
	} else {
		// Define our flags if they don't exist
		fsDefault := ""
		fsafDefault := true
		fsrDefault := ""
		flagSweep = &fsDefault
		flagSweepAllowFailures = &fsafDefault
		flagSweepRun = &fsrDefault
	}
}

// AddTestSweepers function adds a sweeper configuration to the inventory
func AddTestSweepers(s *Sweeper) {
	if s == nil || s.Name == "" {
		log.Fatalf("attempted to add null sweeper to map")
	}

	if _, ok := sweeperInventory[s.Name]; ok {
		log.Fatalf("[ERR] Error adding (%s) to sweeperFuncs: function already exists in map", s.Name)
	}

	sweeperInventory[s.Name] = s
}

// Legacy support for older sweeper format
func AddTestSweepersLegacy(name string, sweeper func(region string) error) {
	AddTestSweepers(&Sweeper{
		Name:           name,
		DeleteFunction: sweeper,
	})
}

// GetSweeper returns a sweeper by name
func GetSweeper(name string) (*Sweeper, bool) {
	s, ok := sweeperInventory[name]
	return s, ok
}

// ExecuteSweepers runs registered sweepers for specified regions
func ExecuteSweepers(t *testing.T) {
	registerFlags()
	flag.Parse()
	if *flagSweep != "" {
		// parse flagSweep contents for regions to run
		regions := strings.Split(*flagSweep, ",")

		// get filtered list of sweepers to run based on sweep-run flag
		sweepers := filterSweepers(*flagSweepRun, sweeperInventory)

		if err := runSweepers(t, regions, sweepers, *flagSweepAllowFailures); err != nil {
			os.Exit(1)
		}
	} else {
		t.Skip("skipping sweeper run. No region supplied")
	}
}

func runSweepers(t *testing.T, regions []string, sweepers map[string]*Sweeper, allowFailures bool) error {
	// First validate that parent sweepers have ListAndAction
	if err := validateParentSweepers(sweepers); err != nil {
		return fmt.Errorf("parent validation failed: %v", err)
	}

	// Sort sweepers by dependency order, considering both dependencies and parents
	sorted, err := validateAndOrderSweepersWithDependencies(sweepers)
	if err != nil {
		return fmt.Errorf("failed to sort sweepers: %v", err)
	}

	// Run each sweeper in dependency order
	for _, sweeper := range sorted {
		sweeper := sweeper // capture for closure
		t.Run(sweeper.Name, func(t *testing.T) {
			for _, region := range regions {
				region := strings.TrimSpace(region)
				err := sweeper.DeleteFunction(region)

				if err != nil {
					if allowFailures {
						t.Errorf("failed in region %s: %s", region, err)
					} else {
						t.Fatalf("failed in region %s: %s", region, err)
					}
				}
			}
		})
	}

	return nil
}

// filterSweepers takes a comma separated string listing the sweepers to run
func filterSweepers(f string, source map[string]*Sweeper) map[string]*Sweeper {
	filterSlice := strings.Split(strings.ToLower(f), ",")
	if len(filterSlice) == 1 && filterSlice[0] == "" {
		// if the filter slice is a single element of "" then no sweeper list was
		// given, so just return the full list
		return source
	}

	// First convert to the unified model to ensure we include all relationships
	unifiedSource := unifyRelationships(source)

	// Then filter based on the unified model
	sweepers := make(map[string]*Sweeper)
	for name := range source {
		for _, s := range filterSlice {
			if strings.Contains(strings.ToLower(name), s) {
				// When we find a match, include it and all of its relationships
				for foundName := range filterSweeperWithDependencies(name, unifiedSource) {
					// Get the original sweeper (not the unified one)
					sweepers[foundName] = source[foundName]
				}
			}
		}
	}
	return sweepers
}

// filterSweeperWithDependencies collects a sweeper and all its dependencies
func filterSweeperWithDependencies(name string, source map[string]*Sweeper) map[string]*Sweeper {
	result := make(map[string]*Sweeper)

	// Get the current sweeper
	currentSweeper, ok := source[name]
	if !ok {
		log.Printf("[WARN] Sweeper (%s) not found", name)
		return result
	}

	// Add the current sweeper
	result[name] = currentSweeper

	// Add all dependencies recursively
	for _, depName := range currentSweeper.Dependencies {
		if depSweeper, ok := source[depName]; ok {
			result[depName] = depSweeper
			// Recursively add dependencies of dependencies
			for foundName, foundSweeper := range filterSweeperWithDependencies(depName, source) {
				result[foundName] = foundSweeper
			}
		}
	}

	return result
}

// validateAndOrderSweepersWithDependencies orders sweepers based on dependencies
// including implicit dependencies from parent-child relationships
func validateAndOrderSweepersWithDependencies(sweepers map[string]*Sweeper) ([]*Sweeper, error) {
	// First check for contradictions that need to be caught before unification
	if err := validateDependenciesAndParents(sweepers); err != nil {
		return nil, err
	}

	// Create a copy of the sweepers map with parent relationships
	// converted to reverse dependencies
	unifiedSweepers := unifyRelationships(sweepers)

	// Check for cycles in the unified dependency graph
	if err := detectCycles(unifiedSweepers); err != nil {
		return nil, err
	}

	// Build dependency graph and perform topological sort
	return topologicalSort(unifiedSweepers), nil
}

// Add this to the unifyRelationships function
func unifyRelationships(sweepers map[string]*Sweeper) map[string]*Sweeper {
	// Create a copy of the original sweepers
	unified := make(map[string]*Sweeper, len(sweepers))
	for name, sweeper := range sweepers {
		// Clone each sweeper with its dependencies
		unified[name] = &Sweeper{
			Name:           sweeper.Name,
			Dependencies:   make([]string, len(sweeper.Dependencies)),
			Parents:        make([]string, len(sweeper.Parents)),
			ListAndAction:  sweeper.ListAndAction,
			DeleteFunction: sweeper.DeleteFunction,
		}
		copy(unified[name].Dependencies, sweeper.Dependencies)
		copy(unified[name].Parents, sweeper.Parents)
	}

	// Convert parent relationships to reverse dependencies
	// If A has parent B, it means B depends on A (B needs A to be deleted first)
	for childName, child := range sweepers {
		for _, parentName := range child.Parents {
			// Add the child as a dependency of the parent
			parent := unified[parentName]
			if parent != nil {
				// Check if the dependency already exists
				exists := false
				for _, dep := range parent.Dependencies {
					if dep == childName {
						exists = true
						break
					}
				}

				// Add if it doesn't exist
				if !exists {
					parent.Dependencies = append(parent.Dependencies, childName)
				}
			}
		}
	}

	return unified
}

func detectCycles(sweepers map[string]*Sweeper) error {
	// Build a directed graph
	graph := make(map[string][]string)

	// Initialize graph
	for name := range sweepers {
		graph[name] = []string{}
	}

	// Add edges for dependencies
	// If A depends on B, then A → B for cycle detection
	// (A needs B to complete first, so there's a dependency from A to B)
	for name, sweeper := range sweepers {
		for _, dep := range sweeper.Dependencies {
			if dep == name {
				log.Printf("Self-dependency detected: %s depends on itself", name)
				return fmt.Errorf("dependency cycle detected: %s depends on itself", name)
			}
			// Add edge: A → B (A depends on B)
			graph[name] = append(graph[name], dep)
		}
	}

	// Check for cycles using DFS
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	var dfs func(node string, path []string) []string
	dfs = func(node string, path []string) []string {
		if recStack[node] {
			// Found a cycle
			log.Printf("Cycle detected! Node %s is already in recursion stack", node)
			for i, n := range path {
				if n == node {
					cycle := append(path[i:], node)
					log.Printf("Cycle is: %v", cycle)
					return cycle
				}
			}
			return []string{node}
		}

		if visited[node] {
			return nil
		}

		visited[node] = true
		recStack[node] = true

		newPath := append(path, node)
		for _, neighbor := range graph[node] {
			if cycle := dfs(neighbor, newPath); cycle != nil {
				return cycle
			}
		}

		recStack[node] = false
		return nil
	}

	// Try each node as a potential start
	for node := range graph {
		visited = make(map[string]bool)
		recStack = make(map[string]bool)

		if cycle := dfs(node, []string{}); cycle != nil {
			log.Printf("Final cycle detected: %s", strings.Join(cycle, " → "))
			return fmt.Errorf("dependency cycle detected: %s", strings.Join(cycle, " → "))
		}
	}

	// Additional check specifically for indirect cycles that might not be caught by the DFS
	// This is helpful for cases like Indirect_cycle_through_different_relationship_types
	// where the relationships are complex
	for name, sweeper := range sweepers {
		// Create a map to track reachability from this node
		reachable := make(map[string]bool)
		reachable[name] = true

		// Start with dependencies
		toCheck := make([]string, len(sweeper.Dependencies))
		copy(toCheck, sweeper.Dependencies)

		// Breadth-first search to find all reachable nodes
		for len(toCheck) > 0 {
			current := toCheck[0]
			toCheck = toCheck[1:]

			if !reachable[current] {
				reachable[current] = true

				// Add all dependencies of this node to check
				if currSweeper, ok := sweepers[current]; ok {
					for _, dep := range currSweeper.Dependencies {
						if dep == name {
							// If we can reach back to our starting node, it's a cycle
							log.Printf("Indirect cycle detected: %s can reach itself through dependencies", name)
							return fmt.Errorf("dependency cycle detected: %s can reach itself through dependencies", name)
						}
						if !reachable[dep] {
							toCheck = append(toCheck, dep)
						}
					}
				}
			}
		}
	}

	return nil
}

// topologicalSort implements Kahn's algorithm to order the sweepers
func topologicalSort(sweepers map[string]*Sweeper) []*Sweeper {
	// Build the graph - if A depends on B, then B → A
	graph := make(map[string][]string)
	inDegree := make(map[string]int)

	// Initialize
	for name := range sweepers {
		graph[name] = []string{}
		inDegree[name] = 0
	}

	// Add edges and count in-degrees
	for name, sweeper := range sweepers {
		for _, dep := range sweeper.Dependencies {
			graph[dep] = append(graph[dep], name)
			inDegree[name]++
		}
	}

	// Find nodes with no incoming edges
	var queue []string
	for node, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, node)
		}
	}

	// Process the queue
	var result []*Sweeper
	for len(queue) > 0 {
		// Dequeue a node
		node := queue[0]
		queue = queue[1:]

		// Add to result
		result = append(result, sweepers[node])

		// Update neighbors
		for _, neighbor := range graph[node] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	return result
}

// validateParentSweepers ensures all referenced parent sweepers have a ListAndAction function
func validateParentSweepers(sweepers map[string]*Sweeper) error {
	var validationErrors []string

	// For each sweeper
	for childName, childSweeper := range sweepers {
		// For each parent referenced by this sweeper
		for _, parentName := range childSweeper.Parents {
			// Check if parent exists
			parentSweeper, exists := sweepers[parentName]
			if !exists {
				validationErrors = append(validationErrors,
					fmt.Sprintf("sweeper %s references parent %s, but parent %s not found",
						childName, parentName, parentName))
				continue
			}

			// Check if parent has ListAndAction function
			if parentSweeper.ListAndAction == nil {
				validationErrors = append(validationErrors,
					fmt.Sprintf("sweeper %s references parent %s, but parent %s is missing ListAndAction function",
						childName, parentName, parentName))
			}
		}
	}

	// If any validation errors were found, return them
	if len(validationErrors) > 0 {
		if len(validationErrors) == 1 {
			return fmt.Errorf(validationErrors[0])
		}
		return fmt.Errorf("multiple parent validation issues: %s", strings.Join(validationErrors, "; "))
	}

	return nil
}

// validateDependenciesAndParents ensures all dependencies and parents exist
// and also checks for contradictions between them
func validateDependenciesAndParents(sweepers map[string]*Sweeper) error {
	// First check that all references exist
	for name, sweeper := range sweepers {
		for _, dep := range sweeper.Dependencies {
			if _, exists := sweepers[dep]; !exists {
				return fmt.Errorf("sweeper %s has dependency %s, but %s not found",
					name, dep, dep)
			}
		}

		for _, parent := range sweeper.Parents {
			if _, exists := sweepers[parent]; !exists {
				return fmt.Errorf("sweeper %s has parent %s, but %s not found",
					name, parent, parent)
			}
		}
	}

	// Now check for contradictions within the same sweeper
	for name, sweeper := range sweepers {
		// Check if the same resource is listed as both a dependency and a parent
		for _, dep := range sweeper.Dependencies {
			for _, parent := range sweeper.Parents {
				if dep == parent {
					return fmt.Errorf("sweeper %s has %s as both a dependency and a parent, which is contradictory",
						name, dep)
				}
			}
		}
	}

	// Check for cross-relationship cycles (the specific pattern tested in TestCrossRelationshipDetection)
	for name, sweeper := range sweepers {
		// For each dependency of this sweeper
		for _, depName := range sweeper.Dependencies {
			depSweeper, exists := sweepers[depName]
			if !exists {
				continue // This should never happen due to the first check
			}

			// Check if the dependency has this sweeper as a parent
			for _, parentOfDep := range depSweeper.Parents {
				if parentOfDep == name {
					return fmt.Errorf("dependency cycle detected: %s depends on %s, but %s has %s as parent",
						name, depName, depName, name)
				}
			}
		}
	}

	return nil
}

// GetFieldOrDefault safely gets a field from ResourceDataMock or returns the default value
func GetFieldOrDefault(d *tpgresource.ResourceDataMock, key, defaultValue string) string {
	if v, ok := d.FieldsInSchema[key]; ok && v != nil {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return defaultValue
}

// GetStringValue tries to get a string representation of a value
func GetStringValue(v interface{}) (string, bool) {
	if v == nil {
		return "", false
	}

	switch val := v.(type) {
	case string:
		return val, true
	case fmt.Stringer:
		return val.String(), true
	default:
		return "", false
	}
}

// ReplaceTemplateVars replaces all {{key}} occurrences in template with values from replacements map
func ReplaceTemplateVars(template string, replacements map[string]string) string {
	result := template
	for key, value := range replacements {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.Replace(result, placeholder, value, -1)
	}
	return result
}

// HasAnyPrefix checks if the input string begins with any prefix from the given slice.
// Returns true if a match is found, false otherwise.
func HasAnyPrefix(input string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(input, p) {
			return true
		}
	}
	return false
}
