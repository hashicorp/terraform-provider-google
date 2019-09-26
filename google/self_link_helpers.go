package google

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Compare only the relative path of two self links.
func compareSelfLinkRelativePaths(_, old, new string, _ *schema.ResourceData) bool {
	oldStripped, err := getRelativePath(old)
	if err != nil {
		return false
	}

	newStripped, err := getRelativePath(new)
	if err != nil {
		return false
	}

	if oldStripped == newStripped {
		return true
	}

	return false
}

// compareSelfLinkOrResourceName checks if two resources are the same resource
//
// Use this method when the field accepts either a name or a self_link referencing a resource.
// The value we store (i.e. `old` in this method), must be a self_link.
func compareSelfLinkOrResourceName(_, old, new string, _ *schema.ResourceData) bool {
	newParts := strings.Split(new, "/")

	if len(newParts) == 1 {
		// `new` is a name
		// `old` is always a self_link
		if GetResourceNameFromSelfLink(old) == newParts[0] {
			return true
		}
	}

	// The `new` string is a self_link
	return compareSelfLinkRelativePaths("", old, new, nil)
}

// Hash the relative path of a self link.
func selfLinkRelativePathHash(selfLink interface{}) int {
	path, _ := getRelativePath(selfLink.(string))
	return hashcode.String(path)
}

func getRelativePath(selfLink string) (string, error) {
	stringParts := strings.SplitAfterN(selfLink, "projects/", 2)
	if len(stringParts) != 2 {
		return "", fmt.Errorf("String was not a self link: %s", selfLink)
	}

	return "projects/" + stringParts[1], nil
}

// Hash the name path of a self link.
func selfLinkNameHash(selfLink interface{}) int {
	name := GetResourceNameFromSelfLink(selfLink.(string))
	return hashcode.String(name)
}

func ConvertSelfLinkToV1(link string) string {
	reg := regexp.MustCompile("/compute/[a-zA-Z0-9]*/projects/")
	return reg.ReplaceAllString(link, "/compute/v1/projects/")
}

func GetResourceNameFromSelfLink(link string) string {
	parts := strings.Split(link, "/")
	return parts[len(parts)-1]
}

func NameFromSelfLinkStateFunc(v interface{}) string {
	return GetResourceNameFromSelfLink(v.(string))
}

func StoreResourceName(resourceLink interface{}) string {
	return GetResourceNameFromSelfLink(resourceLink.(string))
}

type LocationType int

const (
	Zonal LocationType = iota
	Regional
	Global
)

func GetZonalResourcePropertiesFromSelfLinkOrSchema(d *schema.ResourceData, config *Config) (string, string, string, error) {
	return getResourcePropertiesFromSelfLinkOrSchema(d, config, Zonal)
}

func GetRegionalResourcePropertiesFromSelfLinkOrSchema(d *schema.ResourceData, config *Config) (string, string, string, error) {
	return getResourcePropertiesFromSelfLinkOrSchema(d, config, Regional)
}

func getResourcePropertiesFromSelfLinkOrSchema(d *schema.ResourceData, config *Config, locationType LocationType) (string, string, string, error) {
	if selfLink, ok := d.GetOk("self_link"); ok {
		return GetLocationalResourcePropertiesFromSelfLinkString(selfLink.(string))
	} else {
		project, err := getProject(d, config)
		if err != nil {
			return "", "", "", err
		}

		location := ""
		if locationType == Regional {
			location, err = getRegion(d, config)
			if err != nil {
				return "", "", "", err
			}
		} else if locationType == Zonal {
			location, err = getZone(d, config)
			if err != nil {
				return "", "", "", err
			}
		}

		n, ok := d.GetOk("name")
		name := n.(string)
		if !ok {
			return "", "", "", errors.New("must provide either `self_link` or `name`")
		}
		return project, location, name, nil
	}
}

// given a full locational (non-global) self link, returns the project + region/zone + name or an error
func GetLocationalResourcePropertiesFromSelfLinkString(selfLink string) (string, string, string, error) {
	parsed, err := url.Parse(selfLink)
	if err != nil {
		return "", "", "", err
	}

	s := strings.Split(parsed.Path, "/")
	return s[4], s[6], s[8], nil
}

// return the region a selfLink is referring to
func GetRegionFromRegionSelfLink(selfLink string) string {
	re := regexp.MustCompile("/compute/[a-zA-Z0-9]*/projects/[a-zA-Z0-9-]*/regions/([a-zA-Z0-9-]*)")
	switch {
	case re.MatchString(selfLink):
		if res := re.FindStringSubmatch(selfLink); len(res) == 2 && res[1] != "" {
			return res[1]
		}
	}
	return selfLink
}
