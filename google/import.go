package google

import (
	"fmt"
	"regexp"
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
			return fmt.Errorf("Import is not supported. Invalid regex formats.")
		}

		if fieldValues := re.FindStringSubmatch(d.Id()); fieldValues != nil {
			// Starting at index 1, the first match is the full string.
			for i := 1; i < len(fieldValues); i++ {
				fieldName := re.SubexpNames()[i]
				d.Set(fieldName, fieldValues[i])
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
