// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package spanner

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

// Unit Tests for ForceNew when the change in ddl
func TestSpannerDatabase_resourceSpannerDBDdlCustomDiffFuncForceNew(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		before   interface{}
		after    interface{}
		forcenew bool
	}{
		"remove_old_statements": {
			before: []interface{}{
				"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)"},
			after: []interface{}{
				"CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)"},
			forcenew: true,
		},
		"append_new_statements": {
			before: []interface{}{
				"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)"},
			after: []interface{}{
				"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
				"CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)",
			},
			forcenew: false,
		},
		"no_change": {
			before: []interface{}{
				"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)"},
			after: []interface{}{
				"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)"},
			forcenew: false,
		},
		"order_of_statments_change": {
			before: []interface{}{
				"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
				"CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)",
				"CREATE TABLE t3 (t3 INT64 NOT NULL,) PRIMARY KEY(t3)",
			},
			after: []interface{}{
				"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
				"CREATE TABLE t3 (t3 INT64 NOT NULL,) PRIMARY KEY(t3)",
				"CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)",
			},
			forcenew: true,
		},
		"missing_an_old_statement": {
			before: []interface{}{
				"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
				"CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)",
				"CREATE TABLE t3 (t3 INT64 NOT NULL,) PRIMARY KEY(t3)",
			},
			after: []interface{}{
				"CREATE TABLE t1 (t1 INT64 NOT NULL,) PRIMARY KEY(t1)",
				"CREATE TABLE t2 (t2 INT64 NOT NULL,) PRIMARY KEY(t2)",
			},
			forcenew: true,
		},
	}

	for tn, tc := range cases {
		d := &tpgresource.ResourceDiffMock{
			Before: map[string]interface{}{
				"ddl": tc.before,
			},
			After: map[string]interface{}{
				"ddl": tc.after,
			},
		}
		err := resourceSpannerDBDdlCustomDiffFunc(d)
		if err != nil {
			t.Errorf("failed, expected no error but received - %s for the condition %s", err, tn)
		}
		if d.IsForceNew != tc.forcenew {
			t.Errorf("ForceNew not setup correctly for the condition-'%s', expected:%v;actual:%v", tn, tc.forcenew, d.IsForceNew)
		}
	}
}

// Unit Tests for type SpannerDatabaseId
func TestDatabaseNameForApi(t *testing.T) {
	id := SpannerDatabaseId{
		Project:  "project123",
		Instance: "instance456",
		Database: "db789",
	}
	actual := id.databaseUri()
	expected := "projects/project123/instances/instance456/databases/db789"
	expectEquals(t, expected, actual)
}
