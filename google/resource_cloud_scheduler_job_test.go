package google

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestCloudScheduler_FlattenHttpHeaders(t *testing.T) {

	cases := []struct {
		Input  map[string]interface{}
		Output map[string]interface{}
	}{
		// simple, no headers included
		{
			Input: map[string]interface{}{
				"My-Header": "my-header-value",
			},
			Output: map[string]interface{}{
				"My-Header": "my-header-value",
			},
		},

		// include the User-Agent header value Google-Cloud-Scheduler
		// Tests Removing User-Agent header
		{
			Input: map[string]interface{}{
				"User-Agent": "Google-Cloud-Scheduler",
				"My-Header":  "my-header-value",
			},
			Output: map[string]interface{}{
				"My-Header": "my-header-value",
			},
		},

		// include the User-Agent header
		// Tests removing value AppEngine-Google; (+http://code.google.com/appengine)
		{
			Input: map[string]interface{}{
				"User-Agent": "My-User-Agent AppEngine-Google; (+http://code.google.com/appengine)",
				"My-Header":  "my-header-value",
			},
			Output: map[string]interface{}{
				"User-Agent": "My-User-Agent",
				"My-Header":  "my-header-value",
			},
		},

		// include the Content-Type header value application/octet-stream.
		// Tests Removing Content-Type header
		{
			Input: map[string]interface{}{
				"Content-Type": "application/octet-stream",
				"My-Header":    "my-header-value",
			},
			Output: map[string]interface{}{
				"My-Header": "my-header-value",
			},
		},

		// include the Content-Length header
		// Tests Removing Content-Length header
		{
			Input: map[string]interface{}{
				"Content-Length": 7,
				"My-Header":      "my-header-value",
			},
			Output: map[string]interface{}{
				"My-Header": "my-header-value",
			},
		},

		// include the X-Google- header
		// Tests Removing X-Google- header
		{
			Input: map[string]interface{}{
				"X-Google-My-Header": "x-google-my-header-value",
				"My-Header":          "my-header-value",
			},
			Output: map[string]interface{}{
				"My-Header": "my-header-value",
			},
		},
	}

	for _, c := range cases {
		d := &schema.ResourceData{}
		output := flattenCloudSchedulerJobAppEngineHttpTargetHeaders(c.Input, d, &Config{})
		if !reflect.DeepEqual(output, c.Output) {
			t.Fatalf("Error matching output and expected: %#v vs %#v", output, c.Output)
		}
	}
}
