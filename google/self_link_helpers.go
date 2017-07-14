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
