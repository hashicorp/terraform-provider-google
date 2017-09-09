package google

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

// Compare only the relative path of two self links.
func compareSelfLinkRelativePaths(k, old, new string, d *schema.ResourceData) bool {
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

// Use this method when the field accepts either a name or a self_link referencing a resource.
// The value we store (i.e. `old` in this method), must be a self_link.
func compareSelfLinkOrResourceName(k, old, new string, d *schema.ResourceData) bool {
	oldParts := strings.Split(old, "/") // always a self_link
	newParts := strings.Split(new, "/")

	if len(newParts) == 1 {
		// The `new` string is a name
		if oldParts[len(oldParts)-1] == newParts[0] {
			return true
		}
	}

	// The `new` string is a self_link
	return compareSelfLinkRelativePaths(k, old, new, d)
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

func ConvertSelfLinkToV1(link string) string {
	reg := regexp.MustCompile("/compute/[a-zA-Z0-9]*/projects/")
	return reg.ReplaceAllString(link, "/compute/v1/projects/")
}

func GetResourceNameFromSelfLink(link string) string {
	parts := strings.Split(link, "/")
	return parts[len(parts)-1]
}
