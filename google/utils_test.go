package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestCheckGCSName(t *testing.T) {
	valid63 := RandString(t, 63)
	cases := map[string]bool{
		// Valid
		"foobar":       true,
		"foobar1":      true,
		"12345":        true,
		"foo_bar_baz":  true,
		"foo-bar-baz":  true,
		"foo-bar_baz1": true,
		"foo--bar":     true,
		"foo__bar":     true,
		"foo-goog":     true,
		"foo.goog":     true,
		valid63:        true,
		fmt.Sprintf("%s.%s.%s", valid63, valid63, valid63): true,

		// Invalid
		"goog-foobar":     false,
		"foobar-google":   false,
		"-foobar":         false,
		"foobar-":         false,
		"_foobar":         false,
		"foobar_":         false,
		"fo":              false,
		"foo$bar":         false,
		"foo..bar":        false,
		RandString(t, 64): false,
		fmt.Sprintf("%s.%s.%s.%s", valid63, valid63, valid63, valid63): false,
	}

	for bucketName, valid := range cases {
		err := tpgresource.CheckGCSName(bucketName)
		if valid && err != nil {
			t.Errorf("The bucket name %s was expected to pass validation and did not pass.", bucketName)
		} else if !valid && err == nil {
			t.Errorf("The bucket name %s was NOT expected to pass validation and passed.", bucketName)
		}
	}
}
