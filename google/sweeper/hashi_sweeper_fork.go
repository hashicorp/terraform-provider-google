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
	"time"
)

// flagSweep is a flag available when running tests on the command line. It
// contains a comma separated list of regions to for the sweeper functions to
// run in.  This flag bypasses the normal Test path and instead runs functions designed to
// clean up any leaked resources a testing environment could have created. It is
// a best effort attempt, and relies on Provider authors to implement "Sweeper"
// methods for resources.

// Adding Sweeper methods with AddTestSweepers will
// construct a list of sweeper funcs to be called here. We iterate through
// regions provided by the sweep flag, and for each region we iterate through the
// tests, and exit on any errors. At time of writing, sweepers are ran
// sequentially, however they can list dependencies to be ran first. We track
// the sweepers that have been ran, so as to not run a sweeper twice for a given
// region.
//
// WARNING:
// Sweepers are designed to be destructive. You should not use the -sweep flag
// in any environment that is not strictly a test environment. Resources will be
// destroyed.

var (
	flagSweep              *string
	flagSweepAllowFailures *bool
	flagSweepRun           *string
	sweeperFuncs           map[string]*Sweeper
)

// SweeperFunc is a signature for a function that acts as a sweeper. It
// accepts a string for the region that the sweeper is to be ran in. This
// function must be able to construct a valid client for that region.
type SweeperFunc func(r string) error

type Sweeper struct {
	// Name for sweeper. Must be unique to be ran by the Sweeper Runner
	Name string

	// Dependencies list the const names of other Sweeper functions that must be ran
	// prior to running this Sweeper. This is an ordered list that will be invoked
	// recursively at the helper/resource level
	Dependencies []string

	// Sweeper function that when invoked sweeps the Provider of specific
	// resources
	F SweeperFunc
}

func init() {
	sweeperFuncs = make(map[string]*Sweeper)
}

// registerFlags checks for and gets existing flag definitions before trying to redefine them.
// This is needed because this package and terraform-plugin-testing both define the same sweep flags.
// By checking first, we ensure we reuse any existing flags rather than causing a panic from flag redefinition.
// This allows this module to be used alongside terraform-plugin-testing without conflicts.
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

// AddTestSweepers function adds a given name and Sweeper configuration
// pair to the internal sweeperFuncs map. Invoke this function to register a
// resource sweeper to be available for running when the -sweep flag is used
// with `go test`. Sweeper names must be unique to help ensure a given sweeper
// is only ran once per run.
func addTestSweepers(name string, s *Sweeper) {
	if _, ok := sweeperFuncs[name]; ok {
		log.Fatalf("[ERR] Error adding (%s) to sweeperFuncs: function already exists in map", name)
	}

	sweeperFuncs[name] = s
}

// ExecuteSweepers
//
// Sweepers enable infrastructure cleanup functions to be included with
// resource definitions, typically so developers can remove all resources of
// that resource type from testing infrastructure in case of failures that
// prevented the normal resource destruction behavior of acceptance tests.
// Use the AddTestSweepers() function to configure available sweepers.
//
// Sweeper flags added to the "go test" command:
//
//	-sweep: Comma-separated list of locations/regions to run available sweepers.
//	-sweep-allow-failues: Enable to allow other sweepers to run after failures.
//	-sweep-run: Comma-separated list of resource type sweepers to run. Defaults
//	        to all sweepers.
//
// Refer to the Env prefixed constants for environment variables that further
// control testing functionality.
func ExecuteSweepers(t *testing.T) {
	registerFlags()
	flag.Parse()
	if *flagSweep != "" {
		// parse flagSweep contents for regions to run
		regions := strings.Split(*flagSweep, ",")

		// get filtered list of sweepers to run based on sweep-run flag
		sweepers := filterSweepers(*flagSweepRun, sweeperFuncs)

		if err := runSweepers(t, regions, sweepers, *flagSweepAllowFailures); err != nil {
			os.Exit(1)
		}
	} else {
		t.Skip("skipping sweeper run. No region supplied")
	}
}

func runSweepers(t *testing.T, regions []string, sweepers map[string]*Sweeper, allowFailures bool) error {
	// Sort sweepers by dependency order
	sorted, err := validateAndOrderSweepers(sweepers)
	if err != nil {
		return fmt.Errorf("failed to sort sweepers: %v", err)
	}

	// Run each sweeper in dependency order
	for _, sweeper := range sorted {
		sweeper := sweeper // capture for closure
		t.Run(sweeper.Name, func(t *testing.T) {
			for _, region := range regions {
				region := strings.TrimSpace(region)
				log.Printf("[DEBUG] Running Sweeper (%s) in region (%s)", sweeper.Name, region)

				start := time.Now()
				err := sweeper.F(region)
				elapsed := time.Since(start)

				log.Printf("[DEBUG] Completed Sweeper (%s) in region (%s) in %s", sweeper.Name, region, elapsed)

				if err != nil {
					log.Printf("[ERROR] Error running Sweeper (%s) in region (%s): %s", sweeper.Name, region, err)
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

// filterSweepers takes a comma separated string listing the names of sweepers
// to be ran, and returns a filtered set from the list of all of sweepers to
// run based on the names given.
func filterSweepers(f string, source map[string]*Sweeper) map[string]*Sweeper {
	filterSlice := strings.Split(strings.ToLower(f), ",")
	if len(filterSlice) == 1 && filterSlice[0] == "" {
		// if the filter slice is a single element of "" then no sweeper list was
		// given, so just return the full list
		return source
	}

	sweepers := make(map[string]*Sweeper)
	for name := range source {
		for _, s := range filterSlice {
			if strings.Contains(strings.ToLower(name), s) {
				for foundName, foundSweeper := range filterSweeperWithDependencies(name, source) {
					sweepers[foundName] = foundSweeper
				}
			}
		}
	}
	return sweepers
}

// filterSweeperWithDependencies recursively returns sweeper and all dependencies.
// Since filterSweepers performs fuzzy matching, this function is used
// to perform exact sweeper and dependency lookup.
func filterSweeperWithDependencies(name string, source map[string]*Sweeper) map[string]*Sweeper {
	result := make(map[string]*Sweeper)

	currentSweeper, ok := source[name]
	if !ok {
		log.Printf("[WARN] Sweeper has dependency (%s), but that sweeper was not found", name)
		return result
	}

	result[name] = currentSweeper

	for _, dependency := range currentSweeper.Dependencies {
		for foundName, foundSweeper := range filterSweeperWithDependencies(dependency, source) {
			result[foundName] = foundSweeper
		}
	}

	return result
}

// validateAndOrderSweepers performs topological sort on sweepers based on their dependencies.
// It ensures there are no cycles in the dependency graph and all referenced dependencies exist.
// Returns an ordered list of sweepers where each sweeper appears after its dependencies.
// Returns error if there are any cycles or missing dependencies.
func validateAndOrderSweepers(sweepers map[string]*Sweeper) ([]*Sweeper, error) {
	// Detect cycles and get sorted list
	visited := make(map[string]bool)
	inPath := make(map[string]bool)
	sorted := make([]*Sweeper, 0, len(sweepers))

	var visit func(name string) error
	visit = func(name string) error {
		if inPath[name] {
			return fmt.Errorf("dependency cycle detected: %s", name)
		}
		if visited[name] {
			return nil
		}

		inPath[name] = true
		sweeper := sweepers[name]
		for _, dep := range sweeper.Dependencies {
			if _, exists := sweepers[dep]; !exists {
				return fmt.Errorf("sweeper %s depends on %s, but %s not found", name, dep, dep)
			}
			if err := visit(dep); err != nil {
				return err
			}
		}
		inPath[name] = false
		visited[name] = true
		sorted = append(sorted, sweeper)
		return nil
	}

	// Visit all sweepers
	for name := range sweepers {
		if !visited[name] {
			if err := visit(name); err != nil {
				return nil, err
			}
		}
	}

	return sorted, nil
}
