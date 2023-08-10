// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package apigee_test

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/services/apigee"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestUnitApigeeInstance_projectListDiffSuppress(t *testing.T) {
	for _, tc := range apigeeInstanceDiffSuppressTestCases {
		tc.Test(t)
	}
}

type ApigeeInstanceDiffSuppressTestCase struct {
	Name           string
	KeysToSuppress []string
	Before         map[string]interface{}
	After          map[string]interface{}
}

var apigeeInstanceDiffSuppressTestCases = []ApigeeInstanceDiffSuppressTestCase{
	{
		Name:           "projects with the same length and one project entry is converted to project id",
		KeysToSuppress: []string{"consumer_accept_list.0"},
		Before: map[string]interface{}{
			"consumer_accept_list.#": 2,
			"consumer_accept_list.0": "45796856818",
			"consumer_accept_list.1": "12345",
		},
		After: map[string]interface{}{
			"consumer_accept_list.#": 2,
			"consumer_accept_list.0": "tf-test8v1bd04pxa",
			"consumer_accept_list.1": "12345",
		},
	},
	{
		Name:           "projects with the same length and no project conversion",
		KeysToSuppress: []string{},
		Before: map[string]interface{}{
			"consumer_accept_list.#": 2,
			"consumer_accept_list.0": "tf-test8v1bd04pxa",
			"consumer_accept_list.1": "12345",
		},
		After: map[string]interface{}{
			"consumer_accept_list.#": 2,
			"consumer_accept_list.0": "tf-test8v1bd04pxa",
			"consumer_accept_list.1": "12345",
		},
	},
	{
		Name:           "projects are empty",
		KeysToSuppress: []string{},
		Before:         map[string]interface{}{},
		After:          map[string]interface{}{},
	},
	{
		Name:           "projects have the different length",
		KeysToSuppress: []string{},
		Before:         map[string]interface{}{},
		After: map[string]interface{}{
			"consumer_accept_list.#": 2,
			"consumer_accept_list.0": "tf-test8v1bd04pxa",
			"consumer_accept_list.1": "12345",
		},
	},
}

func (tc *ApigeeInstanceDiffSuppressTestCase) Test(t *testing.T) {
	mockResourceDiff := &tpgresource.ResourceDiffMock{
		Before: tc.Before,
		After:  tc.After,
	}

	keysHavingDiff := map[string]bool{}

	for key, val1 := range tc.Before {
		val2, ok := tc.After[key]
		if !ok {
			keysHavingDiff[key] = true
		} else if val1 != val2 {
			keysHavingDiff[key] = true
		}
	}

	for key, val1 := range tc.After {
		val2, ok := tc.Before[key]
		if !ok {
			keysHavingDiff[key] = true
		} else if val1 != val2 {
			keysHavingDiff[key] = true
		}
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

	for key := range keysHavingDiff {
		actual := apigee.ProjectListDiffSuppressFunc(mockResourceDiff)
		if actual != keySuppressionMap[key] {
			t.Errorf("Test %s: expected key `%s` to be suppressed", tc.Name, key)
		}
	}
}
