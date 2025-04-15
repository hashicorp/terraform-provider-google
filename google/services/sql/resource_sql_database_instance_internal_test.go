// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sql

import (
	"testing"
)

func TestMaintenanceVersionDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New       string
		ShouldSuppress bool
	}{
		"older configuration maintenance version than current version should suppress diff": {
			Old:            "MYSQL_8_0_26.R20220508.01_09",
			New:            "MYSQL_5_7_37.R20210508.01_03",
			ShouldSuppress: true,
		},
		"older configuration maintenance version than current version should suppress diff with lexicographically smaller database version": {
			Old:            "MYSQL_5_8_10.R20220508.01_09",
			New:            "MYSQL_5_8_7.R20210508.01_03",
			ShouldSuppress: true,
		},
		"newer configuration maintenance version than current version should not suppress diff": {
			Old:            "MYSQL_5_7_37.R20210508.01_03",
			New:            "MYSQL_8_0_26.R20220508.01_09",
			ShouldSuppress: false,
		},
	}

	for tn, tc := range cases {
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()
			if maintenanceVersionDiffSuppress("version", tc.Old, tc.New, nil) != tc.ShouldSuppress {
				t.Fatalf("%q => %q expect DiffSuppress to return %t", tc.Old, tc.New, tc.ShouldSuppress)
			}
		})
	}
}

func TestDatabaseVersionDiffSuppress(t *testing.T) {
	testCases := map[string]struct {
		oldVersion, newVersion string
		shouldSuppressDiff     bool
	}{
		"MySQL 5.6 (non-supported for auto-upgrade) to MySQL 5.7 (non-supported for auto-upgrade) change should not suppress diff": {
			oldVersion:         "MYSQL_5_6",
			newVersion:         "MYSQL_5_7",
			shouldSuppressDiff: false,
		},
		"MySQL 5.7 (non-supported for auto-upgrade) to MySQL 8.0.31 (non-supported for auto-upgrade) change should not suppress diff": {
			oldVersion:         "MYSQL_5_7",
			newVersion:         "MYSQL_8_0_31",
			shouldSuppressDiff: false,
		},
		"MySQL 5.7 (non-supported for auto-upgrade) to MySQL 8.0.40 (supported for auto-upgrade) change should not suppress diff": {
			oldVersion:         "MYSQL_5_7",
			newVersion:         "MYSQL_8_0_40",
			shouldSuppressDiff: false,
		},
		"MySQL 5.7 (non-supported for auto-upgrade) to MySQL 8.0 (supported for auto-upgrade) change should not suppress diff": {
			oldVersion:         "MYSQL_5_7",
			newVersion:         "MYSQL_8_0",
			shouldSuppressDiff: false,
		},
		"MySQL 8.0.31 (non-supported for auto-upgrade) to MySQL 8.0.35 (supported for auto-upgrade) change should not suppress diff": {
			oldVersion:         "MYSQL_8_0_31",
			newVersion:         "MYSQL_8_0_35",
			shouldSuppressDiff: false,
		},
		"MySQL 8.0.31 (non-supported for auto-upgrade) to MySQL 8.0.40 (supported for auto-upgrade) change should not suppress diff": {
			oldVersion:         "MYSQL_8_0_31",
			newVersion:         "MYSQL_8_0_40",
			shouldSuppressDiff: false,
		},
		"MySQL 8.0.31 (non-supported for auto-upgrade) to MySQL 8.0 (supported for auto-upgrade) change should not suppress diff": {
			oldVersion:         "MYSQL_8_0_31",
			newVersion:         "MYSQL_8_0",
			shouldSuppressDiff: false,
		},
		"MySQL 8.0.35 (supported for auto-upgrade) to MySQL 8.0.40 (supported for auto-upgrade) change should not suppress diff": {
			oldVersion:         "MYSQL_8_0_35",
			newVersion:         "MYSQL_8_0_40",
			shouldSuppressDiff: false,
		},
		"MySQL 8.0.35 (supported for auto-upgrade) to MySQL 8.0 (supported for auto-upgrade) change should suppress diff": {
			oldVersion:         "MYSQL_8_0_35",
			newVersion:         "MYSQL_8_0",
			shouldSuppressDiff: true,
		},
		"MySQL 8.0.37 (supported for auto-upgrade) to MySQL 8.0 (supported for auto-upgrade) change should suppress diff": {
			oldVersion:         "MYSQL_8_0_37",
			newVersion:         "MYSQL_8_0",
			shouldSuppressDiff: true,
		},
		"MySQL 8.0.40 (supported for auto-upgrade) to MySQL 8.0 (supported for auto-upgrade) change should suppress diff": {
			oldVersion:         "MYSQL_8_0_40",
			newVersion:         "MYSQL_8_0",
			shouldSuppressDiff: true,
		},
		"MySQL 8.0.41 (supported for auto-upgrade) to MySQL 8.0 (supported for auto-upgrade) change should suppress diff": {
			oldVersion:         "MYSQL_8_0_41",
			newVersion:         "MYSQL_8_0",
			shouldSuppressDiff: true,
		},
		"MySQL 8.0.37 (supported for auto-upgrade) to MySQL 8.4 (non-supported for auto-upgrade) change should not suppress diff": {
			oldVersion:         "MYSQL_8_0_37",
			newVersion:         "MYSQL_8_4",
			shouldSuppressDiff: false,
		},
		"Postgres (or any non-MySQL) versions should not suppress diff": {
			oldVersion:         "POSTGRES_14",
			newVersion:         "POSTGRES_15",
			shouldSuppressDiff: false,
		},
	}

	for testNumber, testCase := range testCases {
		t.Run(testNumber, func(t *testing.T) {
			t.Parallel()
			if databaseVersionDiffSuppress("version", testCase.oldVersion, testCase.newVersion, nil) != testCase.shouldSuppressDiff {
				t.Fatalf("%q => %q expect DiffSuppress to return %t", testCase.oldVersion, testCase.newVersion, testCase.shouldSuppressDiff)
			}
		})
	}
}
