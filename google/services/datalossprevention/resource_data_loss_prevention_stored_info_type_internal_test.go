// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package datalossprevention

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestAccDataLossPreventionStoredInfoType_dlpStoredInfoTypeCustomDiffFuncForceNew(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		before   map[string]interface{}
		after    map[string]interface{}
		forcenew bool
	}{
		"updating_dictionary": {
			before: map[string]interface{}{
				"dictionary": map[string]interface{}{
					"word_list": map[string]interface{}{
						"word": []string{"word", "word2"},
					},
				},
			},
			after: map[string]interface{}{
				"dictionary": map[string]interface{}{
					"word_list": map[string]interface{}{
						"word": []string{"wordnew", "word2"},
					},
				},
			},
			forcenew: false,
		},
		"updating_large_custom_dictionary": {
			before: map[string]interface{}{
				"large_custom_dictionary": map[string]interface{}{
					"output_path": map[string]interface{}{
						"path": "gs://sample-dlp-bucket/something.json",
					},
				},
			},
			after: map[string]interface{}{
				"large_custom_dictionary": map[string]interface{}{
					"output_path": map[string]interface{}{
						"path": "gs://sample-dlp-bucket/somethingnew.json",
					},
				},
			},
			forcenew: false,
		},
		"updating_regex": {
			before: map[string]interface{}{
				"regex": map[string]interface{}{
					"pattern": "patient",
				},
			},
			after: map[string]interface{}{
				"regex": map[string]interface{}{
					"pattern": "newpatient",
				},
			},
			forcenew: false,
		},
		"changing_from_dictionary_to_large_custom_dictionary": {
			before: map[string]interface{}{
				"dictionary": map[string]interface{}{
					"word_list": map[string]interface{}{
						"word": []string{"word", "word2"},
					},
				},
			},
			after: map[string]interface{}{
				"large_custom_dictionary": map[string]interface{}{
					"output_path": map[string]interface{}{
						"path": "gs://sample-dlp-bucket/something.json",
					},
				},
			},
			forcenew: true,
		},
		"changing_from_dictionary_to_regex": {
			before: map[string]interface{}{
				"dictionary": map[string]interface{}{
					"word_list": map[string]interface{}{
						"word": []string{"word", "word2"},
					},
				},
			},
			after: map[string]interface{}{
				"regex": map[string]interface{}{
					"pattern": "patient",
				},
			},
			forcenew: true,
		},
		"changing_from_large_custom_dictionary_to_regex": {
			before: map[string]interface{}{
				"large_custom_dictionary": map[string]interface{}{
					"output_path": map[string]interface{}{
						"path": "gs://sample-dlp-bucket/something.json",
					},
				},
			},
			after: map[string]interface{}{
				"regex": map[string]interface{}{
					"pattern": "patient",
				},
			},
			forcenew: true,
		},
		"changing_from_large_custom_dictionary_to_dictionary": {
			before: map[string]interface{}{
				"large_custom_dictionary": map[string]interface{}{
					"output_path": map[string]interface{}{
						"path": "gs://sample-dlp-bucket/something.json",
					},
				},
			},
			after: map[string]interface{}{
				"dictionary": map[string]interface{}{
					"word_list": map[string]interface{}{
						"word": []string{"word", "word2"},
					},
				},
			},
			forcenew: true,
		},
		"changing_from_regex_to_dictionary": {
			before: map[string]interface{}{
				"regex": map[string]interface{}{
					"pattern": "patient",
				},
			},
			after: map[string]interface{}{
				"dictionary": map[string]interface{}{
					"word_list": map[string]interface{}{
						"word": []string{"word", "word2"},
					},
				},
			},
			forcenew: true,
		},
		"changing_from_regex_to_large_custom_dictionary": {
			before: map[string]interface{}{
				"regex": map[string]interface{}{
					"pattern": "patient",
				},
			},
			after: map[string]interface{}{
				"large_custom_dictionary": map[string]interface{}{
					"output_path": map[string]interface{}{
						"path": "gs://sample-dlp-bucket/something.json",
					},
				},
			},
			forcenew: true,
		},
	}

	for tn, tc := range cases {

		fieldBefore := ""
		fieldAfter := ""
		switch tn {
		case "updating_dictionary":
			fieldBefore = "dictionary"
			fieldAfter = fieldBefore
		case "updating_large_custom_dictionary":
			fieldBefore = "large_custom_dictionary"
			fieldAfter = fieldBefore
		case "updating_regex":
			fieldBefore = "regex"
			fieldAfter = fieldBefore
		case "changing_from_dictionary_to_large_custom_dictionary":
			fieldBefore = "dictionary"
			fieldAfter = "large_custom_dictionary"
		case "changing_from_dictionary_to_regex":
			fieldBefore = "dictionary"
			fieldAfter = "regex"
		case "changing_from_large_custom_dictionary_to_regex":
			fieldBefore = "large_custom_dictionary"
			fieldAfter = "regex"
		case "changing_from_large_custom_dictionary_to_dictionary":
			fieldBefore = "large_custom_dictionary"
			fieldAfter = "dictionary"
		case "changing_from_regex_to_dictionary":
			fieldBefore = "regex"
			fieldAfter = "dictionary"
		case "changing_from_regex_to_large_custom_dictionary":
			fieldBefore = "regex"
			fieldAfter = "large_custom_dictionary"
		}

		d := &tpgresource.ResourceDiffMock{
			Before: map[string]interface{}{
				fieldBefore: tc.before[fieldBefore],
			},
			After: map[string]interface{}{
				fieldAfter: tc.after[fieldAfter],
			},
		}
		err := storedInfoTypeCustomizeDiffFunc(d)
		if err != nil {
			t.Errorf("failed, expected no error but received - %s for the condition %s", err, tn)
		}
		if d.IsForceNew != tc.forcenew {
			t.Errorf("ForceNew not setup correctly for the condition-'%s', expected:%v; actual:%v", tn, tc.forcenew, d.IsForceNew)
		}
	}
}
