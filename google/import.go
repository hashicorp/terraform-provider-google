package google

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// Parse an import id extracting field values using the given list of regexes.
// They are applied in order. The first in the list is tried first.
//
// e.g:
// - projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/subnetworks/(?P<name>[^/]+) (applied first)
// - (?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+),
// - (?P<name>[^/]+) (applied last)
func parseImportId(idRegexes []string, d TerraformResourceData, config *Config) error {
	for _, idFormat := range idRegexes {
		re, err := regexp.Compile(idFormat)

		if err != nil {
			log.Printf("[DEBUG] Could not compile %s.", idFormat)
			return fmt.Errorf("Import is not supported. Invalid regex formats.")
		}

		if fieldValues := re.FindStringSubmatch(d.Id()); fieldValues != nil {
			log.Printf("[DEBUG] matching ID %s to regex %s.", d.Id(), idFormat)
			// Starting at index 1, the first match is the full string.
			for i := 1; i < len(fieldValues); i++ {
				fieldName := re.SubexpNames()[i]
				fieldValue := fieldValues[i]
				log.Printf("[DEBUG] importing %s = %s", fieldName, fieldValue)
				// Because we do not know at this point whether 'fieldName'
				// corresponds to a TypeString or a TypeInteger in the resource
				// schema, we need to determine the type in an unintuitive way.
				// We call d.Get, because examining the empty value is the easiest
				// way to get that out.  Normally, we would be able to just
				// use a try/catch pattern - try as a string, and if that doesn't
				// work, try as an integer, and if that doesn't work, return the
				// error.  Unfortunately, this is not possible here - during tests,
				// d.Set(...) will panic if there is an error.
				val, _ := d.GetOk(fieldName)
				if _, ok := val.(string); val == nil || ok {
					if err = d.Set(fieldName, fieldValue); err != nil {
						return err
					}
				} else if _, ok := val.(int); ok {
					if intVal, atoiErr := strconv.Atoi(fieldValue); atoiErr == nil {
						// If the value can be parsed as an integer, we try to set the
						// value as an integer.
						if err = d.Set(fieldName, intVal); err != nil {
							return err
						}
					} else {
						return fmt.Errorf("%s appears to be an integer, but %v cannot be parsed as an int", fieldName, fieldValue)
					}
				} else {
					return fmt.Errorf(
						"cannot handle %s, which currently has value %v, and should be set to %#v, during import", fieldName, val, fieldValue)
				}
			}

			// The first id format is applied first and contains all the fields.
			err := setDefaultValues(idRegexes[0], d, config)
			if err != nil {
				return err
			}

			return nil
		}
	}
	return fmt.Errorf("Import id %q doesn't match any of the accepted formats: %v", d.Id(), idRegexes)
}

func setDefaultValues(idRegex string, d TerraformResourceData, config *Config) error {
	if _, ok := d.GetOk("project"); !ok && strings.Contains(idRegex, "?P<project>") {
		project, err := getProject(d, config)
		if err != nil {
			return err
		}
		d.Set("project", project)
	}
	if _, ok := d.GetOk("region"); !ok && strings.Contains(idRegex, "?P<region>") {
		region, err := getRegion(d, config)
		if err != nil {
			return err
		}
		d.Set("region", region)
	}
	if _, ok := d.GetOk("zone"); !ok && strings.Contains(idRegex, "?P<zone>") {
		zone, err := getZone(d, config)
		if err != nil {
			return err
		}
		d.Set("zone", zone)
	}
	return nil
}

// Parse an import id extracting field values using the given list of regexes.
// They are applied in order. The first in the list is tried first.
// This does not mutate any of the parameters, returning a map of matches
// Similar to parseImportId in import.go, but less import specific
//
// e.g:
// - projects/(?P<project>[^/]+)/regions/(?P<region>[^/]+)/subnetworks/(?P<name>[^/]+) (applied first)
// - (?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+),
// - (?P<name>[^/]+) (applied last)
func getImportIdQualifiers(idRegexes []string, d TerraformResourceData, config *Config, id string) (map[string]string, error) {
	for _, idFormat := range idRegexes {
		re, err := regexp.Compile(idFormat)

		if err != nil {
			log.Printf("[DEBUG] Could not compile %s.", idFormat)
			return nil, fmt.Errorf("Import is not supported. Invalid regex formats.")
		}

		if fieldValues := re.FindStringSubmatch(id); fieldValues != nil {
			result := make(map[string]string)
			log.Printf("[DEBUG] matching ID %s to regex %s.", id, idFormat)
			// Starting at index 1, the first match is the full string.
			for i := 1; i < len(fieldValues); i++ {
				fieldName := re.SubexpNames()[i]
				fieldValue := fieldValues[i]
				result[fieldName] = fieldValue
			}

			defaults, err := getDefaultValues(idRegexes[0], d, config)
			if err != nil {
				return nil, err
			}

			for k, v := range defaults {
				if _, ok := result[k]; !ok {
					if v == "" {
						// No default was found and no value was specified in the import ID
						return nil, fmt.Errorf("No value was found for %s during import", k)
					}
					// Set any fields that are defaultable and not specified in import ID
					result[k] = v
				}
			}

			return result, nil
		}
	}
	return nil, fmt.Errorf("Import id %q doesn't match any of the accepted formats: %v", d.Id(), idRegexes)
}

// Returns a set of default values that are contained in a regular expression
// This does not mutate any parameters, instead returning a map of defaults
func getDefaultValues(idRegex string, d TerraformResourceData, config *Config) (map[string]string, error) {
	result := make(map[string]string)
	if _, ok := d.GetOk("project"); !ok && strings.Contains(idRegex, "?P<project>") {
		project, _ := getProject(d, config)
		result["project"] = project
	}
	if _, ok := d.GetOk("region"); !ok && strings.Contains(idRegex, "?P<region>") {
		region, _ := getRegion(d, config)
		result["region"] = region
	}
	if _, ok := d.GetOk("zone"); !ok && strings.Contains(idRegex, "?P<zone>") {
		zone, _ := getZone(d, config)
		result["zone"] = zone
	}
	return result, nil
}
