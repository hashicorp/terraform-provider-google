// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"

	"google.golang.org/api/compute/v1"
)

func TestInstanceGroupManager_parseUniqueId(t *testing.T) {
	expectations := map[string][]string{
		"projects/imre-test/global/instanceTemplates/example-template-custom?uniqueId=123":                                       {"projects/imre-test/global/instanceTemplates/example-template-custom", "123"},
		"https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom?uniqueId=123": {"https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom", "123"},
		"projects/imre-test/global/instanceTemplates/example-template-custom":                                                    {"projects/imre-test/global/instanceTemplates/example-template-custom", ""},
		"https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom":              {"https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom", ""},
		"example-template-custom?uniqueId=123":                                                                                   {"example-template-custom", "123"},

		// this test demonstrates that uniqueIds can't override eachother
		"projects/imre-test/global/instanceTemplates/example?uniqueId=123?uniqueId=456": {"projects/imre-test/global/instanceTemplates/example", "123?uniqueId=456"},
	}

	for k, v := range expectations {
		aName, aUniqueId := parseUniqueId(k)
		if v[0] != aName {
			t.Errorf("parseUniqueId failed; name of %v should be %v, not %v", k, v[0], aName)
		}
		if v[1] != aUniqueId {
			t.Errorf("parseUniqueId failed; uniqueId of %v should be %v, not %v", k, v[1], aUniqueId)
		}
	}
}

func TestInstanceGroupManager_compareInstanceTemplate(t *testing.T) {
	shouldAllMatch := []string{
		// uniqueId not present
		"https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom",
		"projects/imre-test/global/instanceTemplates/example-template-custom",
		// uniqueId present
		"https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom?uniqueId=123",
		"projects/imre-test/global/instanceTemplates/example-template-custom?uniqueId=123",
	}
	shouldNotMatch := map[string]string{
		// mismatching name
		"https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom": "projects/imre-test/global/instanceTemplates/example-template-custom2",
		"projects/imre-test/global/instanceTemplates/example-template-custom":                                       "https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom2",
		// matching name, but mismatching uniqueId
		"https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom?uniqueId=123": "projects/imre-test/global/instanceTemplates/example-template-custom?uniqueId=1234",
		"projects/imre-test/global/instanceTemplates/example-template-custom?uniqueId=123":                                       "https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom?uniqueId=1234",
	}
	for _, v1 := range shouldAllMatch {
		for _, v2 := range shouldAllMatch {
			if !compareSelfLinkRelativePathsIgnoreParams("", v1, v2, nil) {
				t.Fatalf("compareSelfLinkRelativePathsIgnoreParams did not match (and should have) %v and %v", v1, v2)
			}
		}
	}

	for v1, v2 := range shouldNotMatch {
		if compareSelfLinkRelativePathsIgnoreParams("", v1, v2, nil) {
			t.Fatalf("compareSelfLinkRelativePathsIgnoreParams did match (and shouldn't) %v and %v", v1, v2)
		}
	}
}

func TestInstanceGroupManager_convertUniqueId(t *testing.T) {
	matches := map[string]string{
		// uniqueId not present (should return the same)
		"https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom": "https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom",
		"projects/imre-test/global/instanceTemplates/example-template-custom":                                       "projects/imre-test/global/instanceTemplates/example-template-custom",
		// uniqueId present (should return the last component replaced)
		"https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/example-template-custom?uniqueId=123": "https://www.googleapis.com/compute/v1/projects/imre-test/global/instanceTemplates/123",
		"projects/imre-test/global/instanceTemplates/example-template-custom?uniqueId=123":                                       "projects/imre-test/global/instanceTemplates/123",
		"tf-test-igm-8amncgtq22?uniqueId=8361222501423044003":                                                                    "8361222501423044003",
	}
	for input, expected := range matches {
		actual := ConvertToUniqueIdWhenPresent(input)
		if actual != expected {
			t.Fatalf("invalid return value by ConvertToUniqueIdWhenPresent for input %v; expected: %v, actual: %v", input, expected, actual)
		}
	}
}

func TestFlattenStatefulPolicyStatefulIps(t *testing.T) {
	cases := map[string]struct {
		ConfigValues []interface{}
		Ips          map[string]compute.StatefulPolicyPreservedStateNetworkIp
		Expected     []map[string]interface{}
	}{
		"No IPs in config nor API data": {
			ConfigValues: []interface{}{},
			Ips:          map[string]compute.StatefulPolicyPreservedStateNetworkIp{},
			Expected:     []map[string]interface{}{},
		},
		"Single IP (nic0) in config and API data": {
			ConfigValues: []interface{}{
				map[string]interface{}{
					"interface_name": "nic0",
					"delete_rule":    "NEVER",
				},
			},
			Ips: map[string]compute.StatefulPolicyPreservedStateNetworkIp{
				"nic0": {
					AutoDelete: "NEVER",
				},
			},
			Expected: []map[string]interface{}{
				{
					"interface_name": "nic0",
					"delete_rule":    "NEVER",
				},
			},
		},
		"Two IPs (nic0, nic1). Unordered in config and sorted in API data": {
			ConfigValues: []interface{}{
				map[string]interface{}{
					"interface_name": "nic1",
					"delete_rule":    "NEVER",
				},
				map[string]interface{}{
					"interface_name": "nic0",
					"delete_rule":    "NEVER",
				},
			},
			Ips: map[string]compute.StatefulPolicyPreservedStateNetworkIp{
				"nic0": {
					AutoDelete: "NEVER",
				},
				"nic1": {
					AutoDelete: "NEVER",
				},
			},
			Expected: []map[string]interface{}{
				{
					"interface_name": "nic1",
					"delete_rule":    "NEVER",
				},
				{
					"interface_name": "nic0",
					"delete_rule":    "NEVER",
				},
			},
		},
		"Two IPs (nic0, nic1). Only nic0 in config and both stored in API data": {
			ConfigValues: []interface{}{
				map[string]interface{}{
					"interface_name": "nic0",
					"delete_rule":    "NEVER",
				},
			},
			Ips: map[string]compute.StatefulPolicyPreservedStateNetworkIp{
				"nic0": {
					AutoDelete: "NEVER",
				},
				"nic1": {
					AutoDelete: "NEVER",
				},
			},
			Expected: []map[string]interface{}{
				{
					"interface_name": "nic0",
					"delete_rule":    "NEVER",
				},
				{
					"interface_name": "nic1",
					"delete_rule":    "NEVER",
				},
			},
		},
		"Five IPs (nic0 - nic4). None stored in config and all stored in API data": {
			ConfigValues: []interface{}{},
			Ips: map[string]compute.StatefulPolicyPreservedStateNetworkIp{
				// Out of order here to encourage randomness
				"nic3": {
					AutoDelete: "NEVER",
				},
				"nic0": {
					AutoDelete: "NEVER",
				},
				"nic1": {
					AutoDelete: "NEVER",
				},
				"nic4": {
					AutoDelete: "NEVER",
				},
				"nic2": {
					AutoDelete: "NEVER",
				},
			},
			Expected: []map[string]interface{}{
				{
					"interface_name": "nic0",
					"delete_rule":    "NEVER",
				},
				{
					"interface_name": "nic1",
					"delete_rule":    "NEVER",
				},
				{
					"interface_name": "nic2",
					"delete_rule":    "NEVER",
				},
				{
					"interface_name": "nic3",
					"delete_rule":    "NEVER",
				},
				{
					"interface_name": "nic4",
					"delete_rule":    "NEVER",
				},
			},
		},
		"Three IPs (nic0, nic1, nic2). Only nic1, nic2 in config and all 3 stored in API data": {
			ConfigValues: []interface{}{
				map[string]interface{}{
					"interface_name": "nic1",
					"delete_rule":    "NEVER",
				},
				map[string]interface{}{
					"interface_name": "nic2",
					"delete_rule":    "NEVER",
				},
			},
			Ips: map[string]compute.StatefulPolicyPreservedStateNetworkIp{
				"nic0": {
					AutoDelete: "NEVER",
				},
				"nic1": {
					AutoDelete: "NEVER",
				},
				"nic2": {
					AutoDelete: "NEVER",
				},
			},
			Expected: []map[string]interface{}{
				{
					"interface_name": "nic1",
					"delete_rule":    "NEVER",
				},
				{
					"interface_name": "nic2",
					"delete_rule":    "NEVER",
				},
				{
					"interface_name": "nic0",
					"delete_rule":    "NEVER",
				},
			},
		},
		"Three IPs (nic0, nic1, nic2). Only nic0, nic2 in config and only nic1, nic2 stored in API data": {
			ConfigValues: []interface{}{
				map[string]interface{}{
					"interface_name": "nic2",
					"delete_rule":    "NEVER",
				},
				map[string]interface{}{
					"interface_name": "nic0",
					"delete_rule":    "NEVER",
				},
			},
			Ips: map[string]compute.StatefulPolicyPreservedStateNetworkIp{
				"nic1": {
					AutoDelete: "NEVER",
				},
				"nic2": {
					AutoDelete: "NEVER",
				},
			},
			Expected: []map[string]interface{}{
				{
					"interface_name": "nic2",
					"delete_rule":    "NEVER",
				},
				{
					"interface_name": "nic1",
					"delete_rule":    "NEVER",
				},
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {

			// Terraform config
			schema := ResourceComputeRegionInstanceGroupManager().Schema
			config := map[string]interface{}{
				"stateful_external_ip": tc.ConfigValues,
				"stateful_internal_ip": tc.ConfigValues,
			}
			d := tpgresource.SetupTestResourceDataFromConfigMap(t, schema, config)

			// API response
			statefulPolicyPreservedState := compute.StatefulPolicyPreservedState{
				ExternalIPs: tc.Ips,
				InternalIPs: tc.Ips,
			}
			statefulPolicy := compute.StatefulPolicy{
				PreservedState: &statefulPolicyPreservedState,
			}

			outputExternal := flattenStatefulPolicyStatefulExternalIps(d, &statefulPolicy)
			if !reflect.DeepEqual(tc.Expected, outputExternal) {
				t.Fatalf("expected external IPs output to be %#v, but got %#v", tc.Expected, outputExternal)
			}

			outputInternal := flattenStatefulPolicyStatefulInternalIps(d, &statefulPolicy)
			if !reflect.DeepEqual(tc.Expected, outputInternal) {
				t.Fatalf("expected internal IPs output to be %#v, but got %#v", tc.Expected, outputInternal)
			}
		})
	}
}
