package tpgresource

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Compare only the resource name of two self links/paths.
func CompareResourceNames(_, old, new string, _ *schema.ResourceData) bool {
	return GetResourceNameFromSelfLink(old) == GetResourceNameFromSelfLink(new)
}

// Compare only the relative path of two self links.
func CompareSelfLinkRelativePaths(_, old, new string, _ *schema.ResourceData) bool {
	oldStripped, err := GetRelativePath(old)
	if err != nil {
		return false
	}

	newStripped, err := GetRelativePath(new)
	if err != nil {
		return false
	}

	if oldStripped == newStripped {
		return true
	}

	return false
}

// CompareSelfLinkOrResourceName checks if two resources are the same resource
//
// Use this method when the field accepts either a name or a self_link referencing a resource.
// The value we store (i.e. `old` in this method), must be a self_link.
func CompareSelfLinkOrResourceName(_, old, new string, _ *schema.ResourceData) bool {
	newParts := strings.Split(new, "/")

	if len(newParts) == 1 {
		// `new` is a name
		// `old` is always a self_link
		if GetResourceNameFromSelfLink(old) == newParts[0] {
			return true
		}
	}

	// The `new` string is a self_link
	return CompareSelfLinkRelativePaths("", old, new, nil)
}

// Hash the relative path of a self link.
func SelfLinkRelativePathHash(selfLink interface{}) int {
	path, _ := GetRelativePath(selfLink.(string))
	return Hashcode(path)
}

func GetRelativePath(selfLink string) (string, error) {
	stringParts := strings.SplitAfterN(selfLink, "projects/", 2)
	if len(stringParts) != 2 {
		return "", fmt.Errorf("String was not a self link: %s", selfLink)
	}

	return "projects/" + stringParts[1], nil
}

// Hash the name path of a self link.
func SelfLinkNameHash(selfLink interface{}) int {
	name := GetResourceNameFromSelfLink(selfLink.(string))
	return Hashcode(name)
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

func GetZonalResourcePropertiesFromSelfLinkOrSchema(d *schema.ResourceData, config *transport_tpg.Config) (string, string, string, error) {
	return getResourcePropertiesFromSelfLinkOrSchema(d, config, Zonal)
}

func GetRegionalResourcePropertiesFromSelfLinkOrSchema(d *schema.ResourceData, config *transport_tpg.Config) (string, string, string, error) {
	return getResourcePropertiesFromSelfLinkOrSchema(d, config, Regional)
}

func getResourcePropertiesFromSelfLinkOrSchema(d *schema.ResourceData, config *transport_tpg.Config, locationType LocationType) (string, string, string, error) {
	if selfLink, ok := d.GetOk("self_link"); ok {
		return GetLocationalResourcePropertiesFromSelfLinkString(selfLink.(string))
	} else {
		project, err := GetProject(d, config)
		if err != nil {
			return "", "", "", err
		}

		location := ""
		if locationType == Regional {
			location, err = GetRegion(d, config)
			if err != nil {
				return "", "", "", err
			}
		} else if locationType == Zonal {
			location, err = GetZone(d, config)
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

	// This is a pretty bad way to tell if this is a self link, but stops us
	// from accessing an index out of bounds and causing a panic. generally, we
	// expect bad values to be partial URIs and names, so this will catch them
	if len(s) < 9 {
		return "", "", "", fmt.Errorf("value %s was not a self link", selfLink)
	}

	return s[4], s[6], s[8], nil
}

// This function supports selflinks that have regions and locations in their paths
func GetRegionFromRegionalSelfLink(selfLink string) string {
	re := regexp.MustCompile("projects/[a-zA-Z0-9-]*/(?:locations|regions)/([a-zA-Z0-9-]*)")
	switch {
	case re.MatchString(selfLink):
		if res := re.FindStringSubmatch(selfLink); len(res) == 2 && res[1] != "" {
			return res[1]
		}
	}
	return selfLink
}
