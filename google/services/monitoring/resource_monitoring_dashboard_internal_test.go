package monitoring

import (
	"testing"
)

func TestRemoveEmptyConfigEntries(t *testing.T) {
	cases := []struct {
		name     string
		old      map[string]interface{}
		new      map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "removes empty arrays not in API response",
			old: map[string]interface{}{
				"displayName": "test",
			},
			new: map[string]interface{}{
				"displayName":      "test",
				"dashboardFilters": []interface{}{},
			},
			expected: map[string]interface{}{
				"displayName": "test",
			},
		},
		{
			name: "removes empty maps not in API response",
			old: map[string]interface{}{
				"displayName": "test",
			},
			new: map[string]interface{}{
				"displayName": "test",
				"labels":      map[string]interface{}{},
			},
			expected: map[string]interface{}{
				"displayName": "test",
			},
		},
		{
			name: "removes empty strings not in API response",
			old: map[string]interface{}{
				"targetAxis": "Y1",
			},
			new: map[string]interface{}{
				"targetAxis": "Y1",
				"label":      "",
			},
			expected: map[string]interface{}{
				"targetAxis": "Y1",
			},
		},
		{
			name: "preserves non-empty arrays",
			old: map[string]interface{}{
				"items": []interface{}{"a"},
			},
			new: map[string]interface{}{
				"items": []interface{}{"a"},
			},
			expected: map[string]interface{}{
				"items": []interface{}{"a"},
			},
		},
		{
			name: "handles nested maps recursively",
			old: map[string]interface{}{
				"widget": map[string]interface{}{
					"title": "test",
				},
			},
			new: map[string]interface{}{
				"widget": map[string]interface{}{
					"title":      "test",
					"breakdowns": []interface{}{},
				},
			},
			expected: map[string]interface{}{
				"widget": map[string]interface{}{
					"title": "test",
				},
			},
		},
		{
			name: "handles nested slices of maps recursively",
			old: map[string]interface{}{
				"dataSets": []interface{}{
					map[string]interface{}{
						"plotType": "LINE",
					},
				},
			},
			new: map[string]interface{}{
				"dataSets": []interface{}{
					map[string]interface{}{
						"plotType":   "LINE",
						"breakdowns": []interface{}{},
						"dimensions": []interface{}{},
					},
				},
			},
			expected: map[string]interface{}{
				"dataSets": []interface{}{
					map[string]interface{}{
						"plotType": "LINE",
					},
				},
			},
		},
		{
			name: "does not remove empty values that exist in API response",
			old: map[string]interface{}{
				"items": []interface{}{},
			},
			new: map[string]interface{}{
				"items": []interface{}{},
			},
			expected: map[string]interface{}{
				"items": []interface{}{},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := removeEmptyConfigEntries(tc.old, tc.new)
			if len(result) != len(tc.expected) {
				t.Errorf("expected %d keys, got %d: %v", len(tc.expected), len(result), result)
				return
			}
			for k, v := range tc.expected {
				if result[k] == nil && v != nil {
					t.Errorf("key %q missing from result", k)
				}
			}
			for k := range result {
				if tc.expected[k] == nil {
					t.Errorf("unexpected key %q in result", k)
				}
			}
		})
	}
}

func TestMonitoringDashboardDiffSuppress_emptyArrays(t *testing.T) {
	// Simulates the perma-diff scenario from issue #16173
	// API response (old/state) omits empty arrays, config (new) includes them
	old := `{"displayName":"test","mosaicLayout":{"columns":12,"tiles":[{"height":4,"widget":{"title":"SLO","xyChart":{"chartOptions":{"mode":"COLOR"},"dataSets":[{"plotType":"LINE","targetAxis":"Y1","timeSeriesQuery":{"timeSeriesFilter":{"aggregation":{"perSeriesAligner":"ALIGN_NEXT_OLDER"},"filter":"test"},"unitOverride":"10^2.%"}}],"thresholds":[{"targetAxis":"Y1","value":0.999}]}},"width":12}]}}`
	new := `{"dashboardFilters":[],"displayName":"test","labels":{},"mosaicLayout":{"columns":12,"tiles":[{"height":4,"widget":{"title":"SLO","xyChart":{"chartOptions":{"mode":"COLOR"},"dataSets":[{"breakdowns":[],"dimensions":[],"measures":[],"plotType":"LINE","targetAxis":"Y1","timeSeriesQuery":{"timeSeriesFilter":{"aggregation":{"perSeriesAligner":"ALIGN_NEXT_OLDER"},"filter":"test"},"unitOverride":"10^2.%"}}],"thresholds":[{"label":"","targetAxis":"Y1","value":0.999}]}},"width":12}]}}`

	if !monitoringDashboardDiffSuppress("dashboard_json", old, new, nil) {
		t.Error("expected diff to be suppressed for empty arrays/maps/strings in config")
	}
}
