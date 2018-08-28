package google

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
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
			// log.Println(" - true!")
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

// GetZoneFromSelfLink will attempt to parse the zone if it's in the referenced self link
//
// If there is no zone present or the link is malformed it will return an empty string
func GetZoneFromSelfLink(link string) (string, error) {
	paths := strings.Split(link, "/")
	for i, path := range paths {
		if path == "zones" && i+1 < len(paths) {
			return paths[i+1], nil
		}
	}

	return "", fmt.Errorf("unable to determine zone from self link")
}
