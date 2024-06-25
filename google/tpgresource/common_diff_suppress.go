// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Contains common diff suppress functions.

package tpgresource

import (
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func EmptyOrDefaultStringSuppress(defaultVal string) schema.SchemaDiffSuppressFunc {
	return func(k, old, new string, d *schema.ResourceData) bool {
		return (old == "" && new == defaultVal) || (new == "" && old == defaultVal)
	}
}

func EmptyOrFalseSuppressBoolean(k, old, new string, d *schema.ResourceData) bool {
	o, n := d.GetChange(k)
	return (o == nil && !n.(bool))
}

func CaseDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	return strings.ToUpper(old) == strings.ToUpper(new)
}

func EmptyOrUnsetBlockDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	o, n := d.GetChange(strings.TrimSuffix(k, ".#"))
	return EmptyOrUnsetBlockDiffSuppressLogic(k, old, new, o, n)
}

// The core logic for EmptyOrUnsetBlockDiffSuppress, in a format that is more conducive
// to unit testing.
func EmptyOrUnsetBlockDiffSuppressLogic(k, old, new string, o, n interface{}) bool {
	if !strings.HasSuffix(k, ".#") {
		return false
	}
	var l []interface{}
	if old == "0" && new == "1" {
		l = n.([]interface{})
	} else if new == "0" && old == "1" {
		l = o.([]interface{})
	} else {
		// we don't have one set and one unset, so don't suppress the diff
		return false
	}

	contents, ok := l[0].(map[string]interface{})
	if !ok {
		return false
	}
	for _, v := range contents {
		if !IsEmptyValue(reflect.ValueOf(v)) {
			return false
		}
	}
	return true
}

func TimestampDiffSuppress(format string) schema.SchemaDiffSuppressFunc {
	return func(_, old, new string, _ *schema.ResourceData) bool {
		oldT, err := time.Parse(format, old)
		if err != nil {
			return false
		}

		newT, err := time.Parse(format, new)
		if err != nil {
			return false
		}

		return oldT == newT
	}
}

// Suppress diffs for duration format. ex "60.0s" and "60s" same
// https://developers.google.com/protocol-buffers/docs/reference/google.protobuf#duration
func DurationDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	oDuration, err := time.ParseDuration(old)
	if err != nil {
		return false
	}
	nDuration, err := time.ParseDuration(new)
	if err != nil {
		return false
	}
	return oDuration == nDuration
}

// Suppress diffs when the value read from api
// has the project number instead of the project name
func ProjectNumberDiffSuppress(_, old, new string, _ *schema.ResourceData) bool {
	var a2, b2 string
	reN := regexp.MustCompile("projects/\\d+")
	re := regexp.MustCompile("projects/[^/]+")
	replacement := []byte("projects/equal")
	a2 = string(reN.ReplaceAll([]byte(old), replacement))
	b2 = string(re.ReplaceAll([]byte(new), replacement))
	return a2 == b2
}

func IsNewResource(diff TerraformResourceDiff) bool {
	name := diff.Get("name")
	return name.(string) == ""
}

func CompareCryptoKeyVersions(_, old, new string, _ *schema.ResourceData) bool {
	// The API can return cryptoKeyVersions even though it wasn't specified.
	// format: projects/<project>/locations/<region>/keyRings/<keyring>/cryptoKeys/<key>/cryptoKeyVersions/1

	kmsKeyWithoutVersions := strings.Split(old, "/cryptoKeyVersions")[0]
	if kmsKeyWithoutVersions == new {
		return true
	}

	return false
}

func CidrOrSizeDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	// If the user specified a size and the API returned a full cidr block, suppress.
	return strings.HasPrefix(new, "/") && strings.HasSuffix(old, new)
}
