package main

import (
	"flag"
	"fmt"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	googleOld "github.com/hashicorp/terraform-provider-clean-google/google/provider"
	google "github.com/hashicorp/terraform-provider-google/google/provider"
)

var verbose bool
var vFlag = flag.Bool("verbose", false, "set to true to produce more verbose diffs")
var resourceFlag = flag.String("resource", "", "the name of the terraform resource to diff")

func main() {
	flag.Parse()
	if resourceFlag == nil || *resourceFlag == "" {
		fmt.Print("resource flag not specified\n")
		panic("the resource to diff must be specified")
	}
	resourceName := *resourceFlag
	verbose = *vFlag
	m := google.ResourceMap()
	res, ok := m[resourceName]
	if !ok {
		panic(fmt.Sprintf("Unable to find resource in TPGB: %s", resourceName))
	}
	m2 := googleOld.ResourceMap()
	res2, ok := m2[resourceName]
	if !ok {
		panic(fmt.Sprintf("Unable to find resource in clean TPGB: %s", resourceName))
	}
	fmt.Printf("------------Diffing resource %s------------\n", resourceName)
	diffSchema(res2.Schema, res.Schema, []string{})
	fmt.Print("------------Done------------\n")
}

// Diffs a Terraform resource schema. Calls itself recursively as some fields
// are implemented using schema.Resource as their element type
func diffSchema(old, new map[string]*schema.Schema, path []string) {
	var sharedKeys []string
	var addedKeys []string
	for k := range new {
		if _, ok := old[k]; ok {
			sharedKeys = append(sharedKeys, k)
		} else {
			// Key not found in old schema
			addedKeys = append(addedKeys, k)
		}
	}
	var missingKeys []string
	for k := range old {
		if _, ok := new[k]; !ok {
			missingKeys = append(missingKeys, k)
		}
	}
	sort.Strings(sharedKeys)
	sort.Strings(addedKeys)
	sort.Strings(missingKeys)
	if len(addedKeys) != 0 {
		var qualifiedKeys []string
		for _, k := range addedKeys {
			qualifiedKeys = append(qualifiedKeys, strings.Join(append(path, k), "."))
		}
		fmt.Printf("Fields added in tpgtools: %v\n", qualifiedKeys)
	}
	if len(missingKeys) != 0 {
		var qualifiedKeys []string
		for _, k := range missingKeys {
			qualifiedKeys = append(qualifiedKeys, strings.Join(append(path, k), "."))
		}
		fmt.Printf("Fields missing in tpgtools: %v\n", qualifiedKeys)
	}
	for _, k := range sharedKeys {
		diffSchemaObject(old[k], new[k], append(path, k))
	}
}

// Diffs a schema.Schema object. Calls itself and diffSchema recursively as
// needed on nested fields.
func diffSchemaObject(old, new *schema.Schema, path []string) {
	if old.Required != new.Required {
		fmt.Printf("Required status different for path %s, was: %t is now %t\n", strings.Join(path, "."), old.Required, new.Required)
	}
	if old.Computed != new.Computed {
		fmt.Printf("Computed status different for path %s, was: %t is now %t\n", strings.Join(path, "."), old.Computed, new.Computed)
	}
	if old.Optional != new.Optional {
		fmt.Printf("Optional status different for path %s, was: %t is now %t\n", strings.Join(path, "."), old.Optional, new.Optional)
	}
	if old.ForceNew != new.ForceNew {
		fmt.Printf("ForceNew status different for path %s, was: %t is now %t\n", strings.Join(path, "."), old.ForceNew, new.ForceNew)
	}
	if old.Type != new.Type {
		fmt.Printf("Type different for path %s, was: %s is now %s\n", strings.Join(path, "."), old.Type, new.Type)
		// Types are different, other diffs won't make sense
		return
	}
	if old.Sensitive != new.Sensitive {
		fmt.Printf("Sensitive status different for path %s, was: %t is now %t\n", strings.Join(path, "."), old.Sensitive, new.Sensitive)
	}
	if old.Deprecated != new.Deprecated {
		fmt.Printf("Deprecated status different for path %s, was: %s is now %s\n", strings.Join(path, "."), old.Deprecated, new.Deprecated)
	}
	if old.MaxItems != new.MaxItems {
		fmt.Printf("MaxItems different for path %s, was: %d is now %d\n", strings.Join(path, "."), old.MaxItems, new.MaxItems)
	}
	if old.MinItems != new.MinItems {
		fmt.Printf("MinItems different for path %s, was: %d is now %d\n", strings.Join(path, "."), old.MinItems, new.MinItems)
	}
	if old.Default != new.Default {
		fmt.Printf("Default value different for path %s, was: %v is now %v\n", strings.Join(path, "."), old.Default, new.Default)
	}
	if old.ConfigMode != new.ConfigMode {
		// This is only set on very few complicated resources (instance, container cluster)
		fmt.Printf("ConfigMode different for path %s, was: %v is now %v\n", strings.Join(path, "."), old.ConfigMode, new.ConfigMode)
	}
	// Verbose diffs. Enabled using --verbose flag
	if verbose && !reflect.DeepEqual(old.ConflictsWith, new.ConflictsWith) {
		fmt.Printf("ConflictsWith different for path %s, was: %v is now %v\n", strings.Join(path, "."), old.ConflictsWith, new.ConflictsWith)
	}
	oldDiffSuppressFunc := findFunctionName(old.DiffSuppressFunc)
	newDiffSuppressFunc := findFunctionName(new.DiffSuppressFunc)
	if verbose && oldDiffSuppressFunc != newDiffSuppressFunc {
		fmt.Printf("DiffSuppressFunc for path %s, was: %s is now %s\n", strings.Join(path, "."), oldDiffSuppressFunc, newDiffSuppressFunc)
	}
	oldStateFunc := findFunctionName(old.StateFunc)
	newStateFunc := findFunctionName(new.StateFunc)
	if verbose && oldStateFunc != newStateFunc {
		fmt.Printf("StateFunc for path %s, was: %s is now %s\n", strings.Join(path, "."), oldStateFunc, newStateFunc)
	}
	oldValidateFunc := findFunctionName(old.ValidateFunc)
	newValidateFunc := findFunctionName(new.ValidateFunc)
	if verbose && oldValidateFunc != newValidateFunc {
		fmt.Printf("ValidateFunc for path %s, was: %s is now %s\n", strings.Join(path, "."), oldValidateFunc, newValidateFunc)
	}
	oldSet := findFunctionName(old.Set)
	newSet := findFunctionName(new.Set)
	if verbose && oldSet != newSet {
		fmt.Printf("Set function for path %s, was: %s is now %s\n", strings.Join(path, "."), oldSet, newSet)
	}
	// Recursive calls for nested objects
	if old.Type == schema.TypeList || old.Type == schema.TypeMap || old.Type == schema.TypeSet {
		oldElem := old.Elem
		newElem := new.Elem
		if reflect.TypeOf(oldElem) != reflect.TypeOf(newElem) {
			fmt.Printf("Elem type different for path %s, was: %T is now %T\n", strings.Join(path, "."), oldElem, newElem)
		}
		switch v := oldElem.(type) {
		case *schema.Resource:
			diffSchema(v.Schema, newElem.(*schema.Resource).Schema, path)
		case *schema.Schema:
			// Primitive unnamed field as only element
			diffSchemaObject(v, newElem.(*schema.Schema), append(path, "elem"))
		}
	}
}
func findFunctionName(f interface{}) string {
	ptr := reflect.ValueOf(f).Pointer()
	fun := runtime.FuncForPC(ptr)
	if fun == nil {
		return ""
	}
	split := strings.Split(fun.Name(), ".")
	return split[len(split)-1]
}
