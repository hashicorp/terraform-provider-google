package google

import (
	"encoding/json"
	"reflect"
	"testing"

	"google.golang.org/api/cloudresourcemanager/v1"
)

func TestIamMergeBindings(t *testing.T) {
	testCases := []struct {
		input  []*cloudresourcemanager.Binding
		expect []*cloudresourcemanager.Binding
	}{
		// Nothing to merge - return same list
		{
			input:  []*cloudresourcemanager.Binding{},
			expect: []*cloudresourcemanager.Binding{},
		},
		// No members returns no binding
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role: "role-1",
				},
				{
					Role:    "role-2",
					Members: []string{"member-2"},
				},
			},
			expect: []*cloudresourcemanager.Binding{
				{
					Role:    "role-2",
					Members: []string{"member-2"},
				},
			},
		},
		// Nothing to merge - return same list
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1"},
				},
			},
			expect: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1"},
				},
			},
		},
		// Nothing to merge - return same list
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1"},
				},
				{
					Role:    "role-2",
					Members: []string{"member-2"},
				},
			},
			expect: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1"},
				},
				{
					Role:    "role-2",
					Members: []string{"member-2"},
				},
			},
		},
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1"},
				},
				{
					Role:    "role-1",
					Members: []string{"member-2"},
				},
			},
			expect: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-2"},
				},
			},
		},
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-2"},
				},
				{
					Role:    "role-1",
					Members: []string{"member-3"},
				},
			},
			expect: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-2", "member-3"},
				},
			},
		},
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-3", "member-4"},
				},
				{
					Role:    "role-1",
					Members: []string{"member-2", "member-1"},
				},
				{
					Role:    "role-2",
					Members: []string{"member-1"},
				},
				{
					Role:    "role-1",
					Members: []string{"member-5"},
				},
				{
					Role:    "role-3",
					Members: []string{"member-1"},
				},
				{
					Role:    "role-2",
					Members: []string{"member-2"},
				},
				{Role: "empty-role", Members: []string{}},
			},
			expect: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-2", "member-3", "member-4", "member-5"},
				},
				{
					Role:    "role-2",
					Members: []string{"member-1", "member-2"},
				},
				{
					Role:    "role-3",
					Members: []string{"member-1"},
				},
			},
		},
	}

	for _, tc := range testCases {
		got := mergeBindings(tc.input)
		if !compareBindings(got, tc.expect) {
			t.Errorf("Unexpected value for mergeBindings(%s).\nActual: %s\nExpected: %s\n",
				debugPrintBindings(tc.input), debugPrintBindings(got), debugPrintBindings(tc.expect))
		}
	}
}

func TestIamFilterBindingsWithRoleAndCondition(t *testing.T) {
	testCases := []struct {
		input          []*cloudresourcemanager.Binding
		role           string
		conditionTitle string
		expect         []*cloudresourcemanager.Binding
	}{
		// No-op
		{
			input:  []*cloudresourcemanager.Binding{},
			role:   "role-1",
			expect: []*cloudresourcemanager.Binding{},
		},
		// Remove one binding
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-2"},
				},
			},
			role:   "role-1",
			expect: []*cloudresourcemanager.Binding{},
		},
		// Remove multiple bindings
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-2"},
				},
				{
					Role:    "role-1",
					Members: []string{"member-3"},
				},
			},
			role:   "role-1",
			expect: []*cloudresourcemanager.Binding{},
		},
		// Remove multiple bindings and leave some.
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-2"},
				},
				{
					Role:    "role-2",
					Members: []string{"member-1"},
				},
				{
					Role:    "role-3",
					Members: []string{"member-1", "member-3"},
				},
				{
					Role:    "role-1",
					Members: []string{"member-2"},
				},
				{
					Role:    "role-2",
					Members: []string{"member-1", "member-2"},
				},
			},
			role: "role-1",
			expect: []*cloudresourcemanager.Binding{
				{
					Role:    "role-2",
					Members: []string{"member-1"},
				},
				{
					Role:    "role-3",
					Members: []string{"member-1", "member-3"},
				},
				{
					Role:    "role-2",
					Members: []string{"member-1", "member-2"},
				},
			},
		},
	}

	for _, tc := range testCases {
		got := filterBindingsWithRoleAndCondition(tc.input, tc.role, &cloudresourcemanager.Expr{Title: tc.conditionTitle})
		if !compareBindings(got, tc.expect) {
			t.Errorf("Got unexpected value for removeAllBindingsWithRole(%s, %s).\nActual: %s\nExpected: %s",
				debugPrintBindings(tc.input), tc.role, debugPrintBindings(got), debugPrintBindings(tc.expect))
		}
	}
}

func TestIamSubtractFromBindings(t *testing.T) {
	testCases := []struct {
		input  []*cloudresourcemanager.Binding
		remove []*cloudresourcemanager.Binding
		expect []*cloudresourcemanager.Binding
	}{
		{
			input:  []*cloudresourcemanager.Binding{},
			remove: []*cloudresourcemanager.Binding{},
			expect: []*cloudresourcemanager.Binding{},
		},
		// Empty input should no-op return empty
		{
			input: []*cloudresourcemanager.Binding{},
			remove: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-2"},
				},
			},
			expect: []*cloudresourcemanager.Binding{},
		},
		// Empty removal should return original expect
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-2"},
				},
			},
			remove: []*cloudresourcemanager.Binding{},
			expect: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-2"},
				},
			},
		},
		// Removal not in input should no-op
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-1+"},
				},
			},
			remove: []*cloudresourcemanager.Binding{
				{
					Role:    "role-2",
					Members: []string{"member-2"},
				},
			},
			expect: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-1+"},
				},
			},
		},
		// Same input/remove should return empty
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-2"},
				},
			},
			remove: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-2"},
				},
			},
			expect: []*cloudresourcemanager.Binding{},
		},
		// Single removal
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-2"},
				},
			},
			remove: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1"},
				},
			},
			expect: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-2"},
				},
			},
		},
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-2", "member-3"},
				},
				{
					Role:    "role-2",
					Members: []string{"member-1"},
				},
				{
					Role:    "role-1",
					Members: []string{"member-1"},
				},
				{
					Role:    "role-3",
					Members: []string{"member-1"},
				},
				{
					Role:    "role-2",
					Members: []string{"member-2"},
				},
			},
			remove: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-2", "member-4"},
				},
				{
					Role:    "role-2",
					Members: []string{"member-2"},
				},
			},
			expect: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"member-1", "member-3"},
				},
				{
					Role:    "role-2",
					Members: []string{"member-1"},
				},
				{
					Role:    "role-3",
					Members: []string{"member-1"},
				},
			},
		},
	}

	for _, tc := range testCases {
		got := subtractFromBindings(tc.input, tc.remove...)
		if !compareBindings(got, tc.expect) {
			t.Errorf("Unexpected value for subtractFromBindings(%s, %s).\nActual: %s\nExpected: %s\n",
				debugPrintBindings(tc.input), debugPrintBindings(tc.remove), debugPrintBindings(got), debugPrintBindings(tc.expect))
		}
	}
}

func TestIamCreateIamBindingsMap(t *testing.T) {
	testCases := []struct {
		input  []*cloudresourcemanager.Binding
		expect map[iamBindingKey]map[string]struct{}
	}{
		{
			input:  []*cloudresourcemanager.Binding{},
			expect: map[iamBindingKey]map[string]struct{}{},
		},
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"user-1", "user-2"},
				},
			},
			expect: map[iamBindingKey]map[string]struct{}{
				{"role-1", conditionKey{}}: {"user-1": {}, "user-2": {}},
			},
		},
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"user-1", "user-2"},
				},
				{
					Role:    "role-1",
					Members: []string{"user-3"},
				},
			},
			expect: map[iamBindingKey]map[string]struct{}{
				{"role-1", conditionKey{}}: {"user-1": {}, "user-2": {}, "user-3": {}},
			},
		},
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"user-1", "user-2"},
				},
				{
					Role:    "role-2",
					Members: []string{"user-1"},
				},
			},
			expect: map[iamBindingKey]map[string]struct{}{
				{"role-1", conditionKey{}}: {"user-1": {}, "user-2": {}},
				{"role-2", conditionKey{}}: {"user-1": {}},
			},
		},
		{
			input: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"user-1", "user-2"},
				},
				{
					Role:    "role-2",
					Members: []string{"user-1"},
				},
				{
					Role:    "role-1",
					Members: []string{"user-3"},
				},
				{
					Role:    "role-2",
					Members: []string{"user-2"},
				},
				{
					Role:    "role-3",
					Members: []string{"user-3"},
				},
			},
			expect: map[iamBindingKey]map[string]struct{}{
				{"role-1", conditionKey{}}: {"user-1": {}, "user-2": {}, "user-3": {}},
				{"role-2", conditionKey{}}: {"user-1": {}, "user-2": {}},
				{"role-3", conditionKey{}}: {"user-3": {}},
			},
		},
	}

	for _, tc := range testCases {
		got := createIamBindingsMap(tc.input)
		if !reflect.DeepEqual(got, tc.expect) {
			t.Errorf("Unexpected value for createIamBindingsMap(%s).\nActual: %#v\nExpected: %#v\n",
				debugPrintBindings(tc.input), got, tc.expect)
		}
	}
}

func TestIamListFromIamBindingMap(t *testing.T) {
	testCases := []struct {
		input  map[iamBindingKey]map[string]struct{}
		expect []*cloudresourcemanager.Binding
	}{
		{
			input:  map[iamBindingKey]map[string]struct{}{},
			expect: []*cloudresourcemanager.Binding{},
		},
		{
			input: map[iamBindingKey]map[string]struct{}{
				{"role-1", conditionKey{}}: {"user-1": {}, "user-2": {}},
			},
			expect: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"user-1", "user-2"},
				},
			},
		},
		{
			input: map[iamBindingKey]map[string]struct{}{
				{"role-1", conditionKey{}}: {"user-1": {}},
				{"role-2", conditionKey{}}: {"user-1": {}, "user-2": {}},
			},
			expect: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"user-1"},
				},
				{
					Role:    "role-2",
					Members: []string{"user-1", "user-2"},
				},
			},
		},
		{
			input: map[iamBindingKey]map[string]struct{}{
				{"role-1", conditionKey{}}: {"user-1": {}, "user-2": {}},
				{"role-2", conditionKey{}}: {},
			},
			expect: []*cloudresourcemanager.Binding{
				{
					Role:    "role-1",
					Members: []string{"user-1", "user-2"},
				},
			},
		},
	}

	for _, tc := range testCases {
		got := listFromIamBindingMap(tc.input)
		if !compareBindings(got, tc.expect) {
			t.Errorf("Unexpected value for subtractFromBindings(%v).\nActual: %#v\nExpected: %#v\n",
				tc.input, debugPrintBindings(got), debugPrintBindings(tc.expect))
		}
	}
}

func TestIamRemoveAllAuditConfigsWithService(t *testing.T) {
	testCases := []struct {
		input   []*cloudresourcemanager.AuditConfig
		service string
		expect  []*cloudresourcemanager.AuditConfig
	}{
		// No-op
		{
			service: "foo.googleapis.com",
			input:   []*cloudresourcemanager.AuditConfig{},
			expect:  []*cloudresourcemanager.AuditConfig{},
		},
		// No-op - service not in audit configs
		{
			service: "bar.googleapis.com",
			input: []*cloudresourcemanager.AuditConfig{
				{
					Service: "foo.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType: "ADMIN_READ",
						},
					},
				},
			},
			expect: []*cloudresourcemanager.AuditConfig{
				{
					Service: "foo.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType: "ADMIN_READ",
						},
					},
				},
			},
		},
		// Single removal
		{
			service: "foo.googleapis.com",
			input: []*cloudresourcemanager.AuditConfig{
				{
					Service: "foo.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType: "ADMIN_READ",
						},
					},
				},
			},
			expect: []*cloudresourcemanager.AuditConfig{},
		},
		// Multiple removal/merge
		{
			service: "kms.googleapis.com",
			input: []*cloudresourcemanager.AuditConfig{
				{
					Service: "kms.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType: "ADMIN_READ",
						},
						{
							LogType:         "DATA_WRITE",
							ExemptedMembers: []string{"user-1"},
						},
					},
				},
				{
					Service: "iam.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType:         "ADMIN_READ",
							ExemptedMembers: []string{"user-1"},
						},
					},
				},
				{
					Service: "kms.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType:         "DATA_WRITE",
							ExemptedMembers: []string{"user-2"},
						},
					},
				},
				{
					Service: "iam.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType:         "ADMIN_READ",
							ExemptedMembers: []string{"user-2"},
						},
					},
				},
				{
					Service: "foo.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType:         "DATA_WRITE",
							ExemptedMembers: []string{"user-1"},
						},
					},
				},
				{
					Service: "kms.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType:         "DATA_WRITE",
							ExemptedMembers: []string{"user-3", "user-4"},
						},
						{
							LogType:         "DATA_READ",
							ExemptedMembers: []string{"user-1", "user-2"},
						},
					},
				},
			},
			expect: []*cloudresourcemanager.AuditConfig{
				{
					Service: "iam.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType:         "ADMIN_READ",
							ExemptedMembers: []string{"user-1", "user-2"},
						},
					},
				},
				{
					Service: "foo.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType:         "DATA_WRITE",
							ExemptedMembers: []string{"user-1"},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		got := removeAllAuditConfigsWithService(tc.input, tc.service)
		if !compareAuditConfigs(got, tc.expect) {
			t.Errorf("Got unexpected value for removeAllAuditConfigsWithService(%s, %s).\nActual: %s\nExpected: %s",
				debugPrintAuditConfigs(tc.input), tc.service, debugPrintAuditConfigs(got), debugPrintAuditConfigs(tc.expect))
		}
	}
}

func TestIamCreateIamAuditConfigsMap(t *testing.T) {
	testCases := []struct {
		input  []*cloudresourcemanager.AuditConfig
		expect map[string]map[string]map[string]struct{}
	}{
		{
			input:  []*cloudresourcemanager.AuditConfig{},
			expect: make(map[string]map[string]map[string]struct{}),
		},
		{
			input: []*cloudresourcemanager.AuditConfig{
				{
					Service: "foo.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType: "ADMIN_READ",
						},
					},
				},
			},
			expect: map[string]map[string]map[string]struct{}{
				"foo.googleapis.com": {
					"ADMIN_READ": map[string]struct{}{},
				},
			},
		},
		{
			input: []*cloudresourcemanager.AuditConfig{
				{
					Service: "foo.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType:         "ADMIN_READ",
							ExemptedMembers: []string{"user-1", "user-2"},
						},
						{
							LogType:         "DATA_WRITE",
							ExemptedMembers: []string{"user-1"},
						},
					},
				},
			},
			expect: map[string]map[string]map[string]struct{}{
				"foo.googleapis.com": {
					"ADMIN_READ": map[string]struct{}{"user-1": {}, "user-2": {}},
					"DATA_WRITE": map[string]struct{}{"user-1": {}},
				},
			},
		},
		{
			input: []*cloudresourcemanager.AuditConfig{
				{
					Service: "foo.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType:         "ADMIN_READ",
							ExemptedMembers: []string{"user-1", "user-2"},
						},
						{
							LogType:         "DATA_WRITE",
							ExemptedMembers: []string{"user-1"},
						},
					},
				},
				{
					Service: "foo.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType:         "DATA_READ",
							ExemptedMembers: []string{"user-2"},
						},
					},
				},
			},
			expect: map[string]map[string]map[string]struct{}{
				"foo.googleapis.com": {
					"ADMIN_READ": map[string]struct{}{"user-1": {}, "user-2": {}},
					"DATA_WRITE": map[string]struct{}{"user-1": {}},
					"DATA_READ":  map[string]struct{}{"user-2": {}},
				},
			},
		},
		{
			input: []*cloudresourcemanager.AuditConfig{
				{
					Service: "kms.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType: "ADMIN_READ",
						},
					},
				},
				{
					Service: "foo.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType:         "ADMIN_READ",
							ExemptedMembers: []string{"user-1", "user-2"},
						},
						{
							LogType:         "DATA_WRITE",
							ExemptedMembers: []string{"user-1"},
						},
					},
				},
				{
					Service: "kms.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType:         "ADMIN_READ",
							ExemptedMembers: []string{"user-1", "user-2"},
						},
					},
				},
				{
					Service: "foo.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType:         "DATA_READ",
							ExemptedMembers: []string{"user-2"},
						},
					},
				},
			},
			expect: map[string]map[string]map[string]struct{}{
				"kms.googleapis.com": {
					"ADMIN_READ": map[string]struct{}{"user-1": {}, "user-2": {}},
				},
				"foo.googleapis.com": {
					"ADMIN_READ": map[string]struct{}{"user-1": {}, "user-2": {}},
					"DATA_WRITE": map[string]struct{}{"user-1": {}},
					"DATA_READ":  map[string]struct{}{"user-2": {}},
				},
			},
		},
	}

	for _, tc := range testCases {
		got := createIamAuditConfigsMap(tc.input)
		if !reflect.DeepEqual(got, tc.expect) {
			t.Errorf("Unexpected value for createIamAuditConfigsMap(%s).\nActual: %#v\nExpected: %#v\n",
				debugPrintAuditConfigs(tc.input), got, tc.expect)
		}
	}
}

func TestIamListFromIamAuditConfigsMap(t *testing.T) {
	testCases := []struct {
		input  map[string]map[string]map[string]struct{}
		expect []*cloudresourcemanager.AuditConfig
	}{
		{
			input:  make(map[string]map[string]map[string]struct{}),
			expect: []*cloudresourcemanager.AuditConfig{},
		},
		{
			input: map[string]map[string]map[string]struct{}{
				"foo.googleapis.com": {"ADMIN_READ": map[string]struct{}{}},
			},
			expect: []*cloudresourcemanager.AuditConfig{
				{
					Service: "foo.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType: "ADMIN_READ",
						},
					},
				},
			},
		},
		{
			input: map[string]map[string]map[string]struct{}{
				"foo.googleapis.com": {
					"ADMIN_READ": map[string]struct{}{"user-1": {}, "user-2": {}},
					"DATA_WRITE": map[string]struct{}{"user-1": {}},
					"DATA_READ":  map[string]struct{}{},
				},
			},
			expect: []*cloudresourcemanager.AuditConfig{
				{
					Service: "foo.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType:         "ADMIN_READ",
							ExemptedMembers: []string{"user-1", "user-2"},
						},
						{
							LogType:         "DATA_WRITE",
							ExemptedMembers: []string{"user-1"},
						},
						{
							LogType: "DATA_READ",
						},
					},
				},
			},
		},
		{
			input: map[string]map[string]map[string]struct{}{
				"kms.googleapis.com": {
					"ADMIN_READ": map[string]struct{}{},
					"DATA_READ":  map[string]struct{}{"user-1": {}, "user-2": {}},
				},
				"foo.googleapis.com": {
					"ADMIN_READ": map[string]struct{}{"user-1": {}, "user-2": {}},
					"DATA_WRITE": map[string]struct{}{"user-1": {}},
					"DATA_READ":  map[string]struct{}{"user-2": {}},
				},
			},
			expect: []*cloudresourcemanager.AuditConfig{
				{
					Service: "kms.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType: "ADMIN_READ",
						},
						{
							LogType:         "DATA_READ",
							ExemptedMembers: []string{"user-1", "user-2"},
						},
					},
				},
				{
					Service: "foo.googleapis.com",
					AuditLogConfigs: []*cloudresourcemanager.AuditLogConfig{
						{
							LogType:         "ADMIN_READ",
							ExemptedMembers: []string{"user-1", "user-2"},
						},
						{
							LogType:         "DATA_WRITE",
							ExemptedMembers: []string{"user-1"},
						},
						{
							LogType:         "DATA_READ",
							ExemptedMembers: []string{"user-2"},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		got := listFromIamAuditConfigMap(tc.input)
		if !compareAuditConfigs(got, tc.expect) {
			t.Errorf("Unexpected value for listFromIamAuditConfigMap(%+v).\nActual: %s\nExpected: %s\n",
				tc.input, debugPrintAuditConfigs(got), debugPrintAuditConfigs(tc.expect))
		}
	}
}

// Util to deref and print auditConfigs
func debugPrintAuditConfigs(bs []*cloudresourcemanager.AuditConfig) string {
	v, _ := json.MarshalIndent(bs, "", "\t")
	return string(v)
}

// Util to deref and print bindings
func debugPrintBindings(bs []*cloudresourcemanager.Binding) string {
	v, _ := json.MarshalIndent(bs, "", "\t")
	return string(v)
}
