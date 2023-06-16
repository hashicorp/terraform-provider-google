// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package composer

import (
	"testing"
)

func TestComposerImageVersionDiffSuppress(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		old      string
		new      string
		expected bool
	}{
		{"matches", "composer-1.4.0-airflow-1.10.0", "composer-1.4.0-airflow-1.10.0", true},
		{"preview matches", "composer-1.17.0-preview.0-airflow-2.0.1", "composer-1.17.0-preview.0-airflow-2.0.1", true},
		{"old latest", "composer-latest-airflow-1.10.0", "composer-1.4.1-airflow-1.10.0", true},
		{"new latest", "composer-1.4.1-airflow-1.10.0", "composer-latest-airflow-1.10.0", true},
		{"composer major alias equivalent", "composer-1.4.0-airflow-1.10.0", "composer-1-airflow-1.10", true},
		{"composer major alias different", "composer-1.4.0-airflow-2.1.4", "composer-2-airflow-2.2", false},
		{"composer different", "composer-1.4.0-airflow-1.10.0", "composer-1.4.1-airflow-1.10.0", false},
		{"airflow major alias equivalent", "composer-1.4.0-airflow-1.10.0", "composer-1.4.0-airflow-1", true},
		{"airflow major alias different", "composer-1.4.0-airflow-1.10.0", "composer-1.4.0-airflow-2", false},
		{"airflow major.minor alias equivalent", "composer-1.4.0-airflow-1.10.0", "composer-1.4.0-airflow-1.10", true},
		{"airflow major.minor alias different", "composer-1.4.0-airflow-2.1.4", "composer-1.4.0-airflow-2.2", false},
		{"airflow different", "composer-1.4.0-airflow-1.10.0", "composer-1.4.0-airflow-1.9.0", false},
	}

	for _, tc := range cases {
		if actual := composerImageVersionDiffSuppress("", tc.old, tc.new, nil); actual != tc.expected {
			t.Errorf("'%s' failed, expected %v but got %v", tc.name, tc.expected, actual)
		}
	}
}
