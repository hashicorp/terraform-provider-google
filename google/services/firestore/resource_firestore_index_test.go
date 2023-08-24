// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package firestore_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/services/firestore"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestUnitFirestoreIndex_firestoreIFieldsDiffSuppress(t *testing.T) {
	for _, tc := range firestoreIndexDiffSuppressTestCases {
		tc.Test(t)
	}
}

type FirestoreIndexDiffSuppressTestCase struct {
	Name           string
	KeysToSuppress []string
	Before         map[string]interface{}
	After          map[string]interface{}
}

var firestoreIndexDiffSuppressTestCases = []FirestoreIndexDiffSuppressTestCase{
	{
		Name:           "working_as_intended",
		KeysToSuppress: []string{"fields.#", "fields.2.field_path", "fields.2.banana"},
		Before: map[string]interface{}{
			"fields.#":            3,
			"fields.0.field_path": "a",
			"fields.1.field_path": "b",
			"fields.2.field_path": "__name__",
			"fields.2.banana":     "sdc",
		},
		After: map[string]interface{}{
			"fields.#":            2,
			"fields.0.field_path": "a",
			"fields.1.field_path": "b",
		},
	},
	{
		Name:           "same_size_array",
		KeysToSuppress: []string{},
		Before: map[string]interface{}{
			"fields.#":            3,
			"fields.0.field_path": "a",
			"fields.1.field_path": "b",
			"fields.2.field_path": "__name__",
			"fields.2.banana":     "sdc",
		},
		After: map[string]interface{}{
			"fields.#":            3,
			"fields.0.field_path": "a",
			"fields.1.field_path": "b",
			"fields.2.field_path": "__name__",
			"fields.2.banana":     "sdc",
		},
	},
	{
		Name:           "new_array_larger",
		KeysToSuppress: []string{},
		Before: map[string]interface{}{
			"fields.#":            3,
			"fields.0.field_path": "a",
			"fields.1.field_path": "b",
			"fields.2.field_path": "beep",
			"fields.2.banana":     "sdc",
		},
		After: map[string]interface{}{
			"fields.#":            4,
			"fields.0.field_path": "a",
			"fields.1.field_path": "b",
			"fields.2.field_path": "__name__",
			"fields.2.banana":     "sdc",
		},
	},
	{
		Name:           "does_not_clear_other_fields",
		KeysToSuppress: []string{"fields.#", "fields.2.field_path", "fields.2.banana"},
		Before: map[string]interface{}{
			"fields.#":            3,
			"fields.0.field_path": "a",
			"fields.1.field_path": "b",
			"fields.2.field_path": "__name__",
			"fields.2.banana":     "sdc",
		},
		After: map[string]interface{}{
			"fields.#":            2,
			"fields.0.field_path": "b",
			"fields.1.field_path": "b",
		},
	},
}

func (tc *FirestoreIndexDiffSuppressTestCase) Test(t *testing.T) {
	mockResourceDiff := &tpgresource.ResourceDiffMock{
		Before: tc.Before,
		After:  tc.After,
	}

	keySuppressionMap := map[string]bool{}
	for key := range tc.Before {
		keySuppressionMap[key] = false
	}
	for key := range tc.After {
		keySuppressionMap[key] = false
	}

	for _, key := range tc.KeysToSuppress {
		keySuppressionMap[key] = true
	}

	for key, tcSuppress := range keySuppressionMap {
		oldValue, ok := tc.Before[key]
		if !ok {
			oldValue = ""
		}
		newValue, ok := tc.After[key]
		if !ok {
			newValue = ""
		}
		suppressed := firestore.FirestoreIFieldsDiffSuppressFunc(key, fmt.Sprintf("%v", oldValue), fmt.Sprintf("%v", newValue), mockResourceDiff)
		if suppressed != tcSuppress {
			var expectation string
			if tcSuppress {
				expectation = "be"
			} else {
				expectation = "not be"
			}
			t.Errorf("Test %s: expected key `%s` to %s suppressed", tc.Name, key, expectation)
		}
	}
}
