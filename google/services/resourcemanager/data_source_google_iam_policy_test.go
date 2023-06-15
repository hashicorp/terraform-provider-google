// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDataSourceGoogleIamPolicyRead(t *testing.T) {
	cases := map[string]struct {
		Bindings                  []interface{}
		OriginalBindingCount      int
		ExpectedFinalBindingCount int
		ExpectedPolicyDataString  string
	}{
		"members are sorted alphabetically within a single binding": {
			Bindings: []interface{}{
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:c",
						"user:a",
						"user:b",
					},
				},
			},
			OriginalBindingCount:      1,
			ExpectedFinalBindingCount: 1,
			ExpectedPolicyDataString:  "{\"bindings\":[{\"members\":[\"user:a\",\"user:b\",\"user:c\"],\"role\":\"role/A\"}]}",
		},
		"bindings are sorted by role (regardless of conditions being present)": {
			Bindings: []interface{}{
				map[string]interface{}{
					"role": "role/B",
					"members": []interface{}{
						"user:a",
					},
				},
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"description": "descriptionA",
							"expression":  "expressionA",
							"title":       "titleA",
						},
					},
				},
				map[string]interface{}{
					"role": "role/C",
					"members": []interface{}{
						"user:a",
					},
				},
			},
			OriginalBindingCount:      3,
			ExpectedFinalBindingCount: 3,
			ExpectedPolicyDataString:  "{\"bindings\":[{\"condition\":{\"description\":\"descriptionA\",\"expression\":\"expressionA\",\"title\":\"titleA\"},\"members\":[\"user:a\"],\"role\":\"role/A\"},{\"members\":[\"user:a\"],\"role\":\"role/B\"},{\"members\":[\"user:a\"],\"role\":\"role/C\"}]}",
		},
		"equivalent bindings (with no conditions) are combined into one binding with a larger member list": {
			Bindings: []interface{}{
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
				},
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:b",
					},
				},
			},
			OriginalBindingCount:      2,
			ExpectedFinalBindingCount: 1, // This test combines bindings
			ExpectedPolicyDataString:  "{\"bindings\":[{\"members\":[\"user:a\",\"user:b\"],\"role\":\"role/A\"}]}",
		},
		"exact duplicate bindings are removed before `policy_data` is set": {
			Bindings: []interface{}{
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
				},
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
				},
			},
			OriginalBindingCount:      1, // Duplicates are identified and removed before producing the policy_data string
			ExpectedFinalBindingCount: 1,
			ExpectedPolicyDataString:  "{\"bindings\":[{\"members\":[\"user:a\"],\"role\":\"role/A\"}]}",
		},
		"equivalent bindings (with conditions) are combined into one binding with a larger member list": {
			Bindings: []interface{}{
				// Should not be consolidated into the other bindings as there's no condition
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:c",
					},
				},
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:b",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"description": "descriptionA",
							"expression":  "expressionA",
							"title":       "titleA",
						},
					},
				},
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"description": "descriptionA",
							"expression":  "expressionA",
							"title":       "titleA",
						},
					},
				},
			},
			OriginalBindingCount:      3,
			ExpectedFinalBindingCount: 2, // This test combines bindings
			ExpectedPolicyDataString:  "{\"bindings\":[{\"members\":[\"user:c\"],\"role\":\"role/A\"},{\"condition\":{\"description\":\"descriptionA\",\"expression\":\"expressionA\",\"title\":\"titleA\"},\"members\":[\"user:a\",\"user:b\"],\"role\":\"role/A\"}]}",
		},
		"bindings on the same role are sorted to place bindings without conditions first": {
			Bindings: []interface{}{
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"description": "descriptionA",
							"expression":  "expressionA",
							"title":       "titleA",
						},
					},
				},
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
				},
			},
			OriginalBindingCount:      2,
			ExpectedFinalBindingCount: 2,
			ExpectedPolicyDataString:  "{\"bindings\":[{\"members\":[\"user:a\"],\"role\":\"role/A\"},{\"condition\":{\"description\":\"descriptionA\",\"expression\":\"expressionA\",\"title\":\"titleA\"},\"members\":[\"user:a\"],\"role\":\"role/A\"}]}",
		},
		"bindings (with conditions) on the same role are first sorted by condition expressions": {
			Bindings: []interface{}{
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"expression":  "A",
							"title":       "B",
							"description": "C",
						},
					},
				},
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"expression":  "B",
							"title":       "C",
							"description": "A",
						},
					},
				},
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"expression":  "C",
							"title":       "A",
							"description": "B",
						},
					},
				},
			},
			OriginalBindingCount:      3,
			ExpectedFinalBindingCount: 3,
			ExpectedPolicyDataString:  "{\"bindings\":[{\"condition\":{\"description\":\"C\",\"expression\":\"A\",\"title\":\"B\"},\"members\":[\"user:a\"],\"role\":\"role/A\"},{\"condition\":{\"description\":\"A\",\"expression\":\"B\",\"title\":\"C\"},\"members\":[\"user:a\"],\"role\":\"role/A\"},{\"condition\":{\"description\":\"B\",\"expression\":\"C\",\"title\":\"A\"},\"members\":[\"user:a\"],\"role\":\"role/A\"}]}",
		},
		"bindings (with conditions) on the same role, with matching condition expressions, are next sorted by condition title": {
			Bindings: []interface{}{
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"expression":  "same expression",
							"title":       "B",
							"description": "A",
						},
					},
				},
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"expression":  "same expression",
							"title":       "A",
							"description": "B",
						},
					},
				},
			},
			OriginalBindingCount:      2,
			ExpectedFinalBindingCount: 2,
			ExpectedPolicyDataString:  "{\"bindings\":[{\"condition\":{\"description\":\"B\",\"expression\":\"same expression\",\"title\":\"A\"},\"members\":[\"user:a\"],\"role\":\"role/A\"},{\"condition\":{\"description\":\"A\",\"expression\":\"same expression\",\"title\":\"B\"},\"members\":[\"user:a\"],\"role\":\"role/A\"}]}",
		},
		"bindings (with conditions) on the same role, with matching condition expressions and titles, are next sorted by condition description": {
			Bindings: []interface{}{
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"expression":  "same expression",
							"title":       "same title",
							"description": "B",
						},
					},
				},
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"expression":  "same expression",
							"title":       "same title",
							"description": "A",
						},
					},
				},
			},
			OriginalBindingCount:      2,
			ExpectedFinalBindingCount: 2,
			ExpectedPolicyDataString:  "{\"bindings\":[{\"condition\":{\"description\":\"A\",\"expression\":\"same expression\",\"title\":\"same title\"},\"members\":[\"user:a\"],\"role\":\"role/A\"},{\"condition\":{\"description\":\"B\",\"expression\":\"same expression\",\"title\":\"same title\"},\"members\":[\"user:a\"],\"role\":\"role/A\"}]}",
		},
		"bindings for different roles and conditions are sorted firstly by role, and within a role block sorting is based on conditions": {
			Bindings: []interface{}{
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"expression":  "A",
							"title":       "A",
							"description": "A",
						},
					},
				},
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
				},
				map[string]interface{}{
					"role": "role/A",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"expression":  "B",
							"title":       "B",
							"description": "B",
						},
					},
				},
				map[string]interface{}{
					"role": "role/B",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"expression":  "A",
							"title":       "A",
							"description": "A",
						},
					},
				},
				map[string]interface{}{
					"role": "role/B",
					"members": []interface{}{
						"user:a",
					},
				},
				map[string]interface{}{
					"role": "role/B",
					"members": []interface{}{
						"user:a",
					},
					"condition": []interface{}{
						map[string]interface{}{
							"expression":  "B",
							"title":       "B",
							"description": "B",
						},
					},
				},
			},
			OriginalBindingCount:      6,
			ExpectedFinalBindingCount: 6,
			ExpectedPolicyDataString:  "{\"bindings\":[{\"members\":[\"user:a\"],\"role\":\"role/A\"},{\"condition\":{\"description\":\"A\",\"expression\":\"A\",\"title\":\"A\"},\"members\":[\"user:a\"],\"role\":\"role/A\"},{\"condition\":{\"description\":\"B\",\"expression\":\"B\",\"title\":\"B\"},\"members\":[\"user:a\"],\"role\":\"role/A\"},{\"members\":[\"user:a\"],\"role\":\"role/B\"},{\"condition\":{\"description\":\"A\",\"expression\":\"A\",\"title\":\"A\"},\"members\":[\"user:a\"],\"role\":\"role/B\"},{\"condition\":{\"description\":\"B\",\"expression\":\"B\",\"title\":\"B\"},\"members\":[\"user:a\"],\"role\":\"role/B\"}]}",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// ARRANGE - Create schema.ResourceData variable as test input
			rawData := map[string]interface{}{
				"binding":      tc.Bindings,
				"policy_data":  "",              // Not set
				"audit_config": []interface{}{}, // Not set
			}
			// Note: for TestResourceDataRaw to process rawData ok, test inputs' data types have to be
			// either primitive types, []interface{} or map[string]interface{}
			d := schema.TestResourceDataRaw(t, DataSourceGoogleIamPolicy().Schema, rawData)

			// ACT - Update resource data using `dataSourceGoogleIamPolicyRead`
			var meta interface{}
			err := dataSourceGoogleIamPolicyRead(d, meta)
			if err != nil {
				t.Error(err)
			}

			// ASSERT
			policyData := d.Get("policy_data").(string)
			var jsonObjs interface{}
			json.Unmarshal([]byte(policyData), &jsonObjs)
			objSlice, ok := jsonObjs.(map[string]interface{})
			if !ok {
				t.Errorf("cannot convert the JSON string")
			}
			policyDataBindings := objSlice["bindings"].([]interface{})
			if len(policyDataBindings) != tc.ExpectedFinalBindingCount {
				t.Errorf("expected there to be %d bindings in the policy_data string, got: %d", tc.ExpectedFinalBindingCount, len(policyDataBindings))
			}
			if policyData != tc.ExpectedPolicyDataString {
				t.Errorf("expected `policy_data` to be %s, got: %s", tc.ExpectedPolicyDataString, policyData)
			}

			bset := d.Get("binding").(*schema.Set)
			if bset.Len() != tc.OriginalBindingCount {
				t.Errorf("expected there to be %d bindings in the data source internals, got: %d", tc.OriginalBindingCount, bset.Len())
			}
		})
	}
}
