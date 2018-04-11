package google

import "strings"

func isZone(location string) bool {
	return len(strings.Split(location, "-")) == 3
}
