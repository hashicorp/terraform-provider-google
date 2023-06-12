// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigquery

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestBigQueryTableSchemaDiffSuppress(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"empty schema": {
			Old:                "null",
			New:                "[]",
			ExpectDiffSuppress: true,
		},
		"empty schema -> non-empty": {
			Old: "null",
			New: `[
				{
					"name": "PageNo",
					"type": "INTEGER"
				}
			]`,
			ExpectDiffSuppress: false,
		},
		"no change": {
			Old:                "[{\"name\": \"someValue\", \"type\": \"INT64\", \"anotherKey\" : \"anotherValue\", \"finalKey\" : {} }]",
			New:                "[{\"name\": \"someValue\", \"type\": \"INT64\", \"anotherKey\" : \"anotherValue\", \"finalKey\" : {} }]",
			ExpectDiffSuppress: true,
		},
		"remove key": {
			Old:                "[{\"name\": \"someValue\", \"type\": \"INT64\", \"anotherKey\" : \"anotherValue\", \"finalKey\" : {} }]",
			New:                "[{\"name\": \"someValue\", \"type\": \"INT64\", \"finalKey\" : {} }]",
			ExpectDiffSuppress: false,
		},
		"empty description -> default description (empty)": {
			Old:                "[{\"name\": \"someValue\", \"type\": \"INT64\", \"anotherKey\" : \"anotherValue\", \"description\": \"\"  }]",
			New:                "[{\"name\": \"someValue\", \"type\": \"INT64\", \"anotherKey\" : \"anotherValue\" }]",
			ExpectDiffSuppress: true,
		},
		"empty description -> other description": {
			Old:                "[{\"name\": \"someValue\", \"type\": \"INT64\", \"anotherKey\" : \"anotherValue\", \"description\": \"\"  }]",
			New:                "[{\"name\": \"someValue\", \"type\": \"INT64\", \"anotherKey\" : \"anotherValue\", \"description\": \"somethingRandom\"  }]",
			ExpectDiffSuppress: false,
		},
		"mode NULLABLE -> other mode": {
			Old:                "[{\"name\": \"someValue\", \"type\": \"INT64\", \"anotherKey\" : \"anotherValue\", \"mode\": \"NULLABLE\"  }]",
			New:                "[{\"name\": \"someValue\", \"type\": \"INT64\", \"anotherKey\" : \"anotherValue\", \"mode\": \"somethingRandom\"  }]",
			ExpectDiffSuppress: false,
		},
		"mode NULLABLE -> default mode (also NULLABLE)": {
			Old: `[
				{
					"mode": "NULLABLE",
					"name": "PageNo",
					"type": "INTEGER"
				}
			]`,
			New: `[
				{
					"name": "PageNo",
					"type": "INTEGER"
				}
			]`,
			ExpectDiffSuppress: true,
		},
		"mode & type uppercase -> lowercase": {
			Old: `[
				{
					"mode": "NULLABLE",
					"name": "PageNo",
					"type": "INTEGER"
				}
			]`,
			New: `[
				{
					"mode": "nullable",
					"name": "PageNo",
					"type": "integer"
				}
			]`,
			ExpectDiffSuppress: true,
		},
		"type INTEGER -> INT64": {
			Old:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"INTEGER\"  }]",
			New:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"INT64\"  }]",
			ExpectDiffSuppress: true,
		},
		"type INTEGER -> other": {
			Old:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"INTEGER\"  }]",
			New:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"somethingRandom\"  }]",
			ExpectDiffSuppress: false,
		},
		"type FLOAT -> FLOAT64": {
			Old:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"FLOAT\"  }]",
			New:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"FLOAT64\"  }]",
			ExpectDiffSuppress: true,
		},
		"type FLOAT -> other": {
			Old:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"FLOAT\"  }]",
			New:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"somethingRandom\" }]",
			ExpectDiffSuppress: false,
		},
		"type BOOLEAN -> BOOL": {
			Old:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"BOOLEAN\"  }]",
			New:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"BOOL\"  }]",
			ExpectDiffSuppress: true,
		},
		"type BOOLEAN -> other": {
			Old:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"BOOLEAN\"  }]",
			New:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"somethingRandom\" }]",
			ExpectDiffSuppress: false,
		},
		// this is invalid but we need to make sure we don't cause a panic
		// if users provide an invalid schema
		"invalid - missing type for old": {
			Old:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\" }]",
			New:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"BOOLEAN\" }]",
			ExpectDiffSuppress: false,
		},
		// this is invalid but we need to make sure we don't cause a panic
		// if users provide an invalid schema
		"invalid - missing type for new": {
			Old:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"BOOLEAN\" }]",
			New:                "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\" }]",
			ExpectDiffSuppress: false,
		},
		"reordering fields": {
			Old: `[
				{
					"name": "PageNo",
					"type": "INTEGER"
				},
				{
					"name": "IngestTime",
					"type": "TIMESTAMP"
				}
			]`,
			New: `[
				{
					"name": "IngestTime",
					"type": "TIMESTAMP"
				},
				{
					"name": "PageNo",
					"type": "INTEGER"
				}
			]`,
			ExpectDiffSuppress: true,
		},
		"reordering fields with value change": {
			Old: `[
				{
					"name": "PageNo",
					"type": "INTEGER",
					"description": "someVal"
				},
				{
					"name": "IngestTime",
					"type": "TIMESTAMP"
				}
			]`,
			New: `[
				{
					"name": "IngestTime",
					"type": "TIMESTAMP"
				},
				{
					"name": "PageNo",
					"type": "INTEGER",
					"description": "otherVal"
				}
			]`,
			ExpectDiffSuppress: false,
		},
		"nested field ordering changes": {
			Old: `[
				{
					"name": "someValue",
					"type": "INTEGER",
					"fields": [
						{
							"name": "value1",
							"type": "INTEGER",
							"mode": "NULLABLE",
							"description": "someVal"
						},
						{
							"name": "value2",
							"type": "BOOLEAN",
							"mode": "NULLABLE",
							"description": "someVal"
						}
					]
				}
			]`,
			New: `[
				{
					"name": "someValue",
					"type": "INTEGER",
					"fields": [
						{
							"name": "value2",
							"type": "BOOLEAN",
							"mode": "NULLABLE",
							"description": "someVal"
						},
						{
							"name": "value1",
							"type": "INTEGER",
							"mode": "NULLABLE",
							"description": "someVal"
						}
					]
				}
			]`,
			ExpectDiffSuppress: true,
		},
		"policyTags": {
			Old: `[
				{
					"mode": "NULLABLE",
					"name": "providerphone",
					"policyTags": {
						"names": [
							"projects/my-project/locations/us/taxonomies/12345678/policyTags/12345678"
						]
					},
					"type":"STRING"
				}
			]`,
			New: `[
			  {
			    "name": "providerphone",
			    "type": "STRING",
			    "policyTags": {
			          "names": ["projects/my-project/locations/us/taxonomies/12345678/policyTags/12345678"]
			        }
			  }
			]`,
			ExpectDiffSuppress: true,
		},
		"multiple levels of reordering with policyTags set": {
			Old: `[
				{
					"mode": "NULLABLE",
					"name": "providerphone",
					"type":"STRING",
					"policyTags": {
						"names": [
							"projects/my-project/locations/us/taxonomies/12345678/policyTags/12345678"
						]
					},
					"fields": [
						{
							"name": "value1",
							"type": "INTEGER",
							"mode": "NULLABLE",
							"description": "someVal",
							"policyTags": {
								"names": [
									"projects/my-project/locations/us/taxonomies/12345678/policyTags/12345678"
								]
							}
						},
						{
							"name": "value2",
							"type": "BOOLEAN",
							"mode": "NULLABLE",
							"description": "someVal"
						}
					]
				},
				{
					"name": "PageNo",
					"type": "INTEGER"
				},
				{
					"name": "IngestTime",
					"type": "TIMESTAMP",
					"fields": [
						{
							"name": "value3",
							"type": "INTEGER",
							"mode": "NULLABLE",
							"description": "someVal",
							"policyTags": {
								"names": [
									"projects/my-project/locations/us/taxonomies/12345678/policyTags/12345678"
								]
							}
						},
						{
							"name": "value4",
							"type": "BOOLEAN",
							"mode": "NULLABLE",
							"description": "someVal"
						}
					]
				}
			]`,
			New: `[
				{
					"name": "IngestTime",
					"type": "TIMESTAMP",
					"fields": [
						{
							"name": "value4",
							"type": "BOOLEAN",
							"mode": "NULLABLE",
							"description": "someVal"
						},
						{
							"name": "value3",
							"type": "INTEGER",
							"mode": "NULLABLE",
							"description": "someVal",
							"policyTags": {
								"names": [
									"projects/my-project/locations/us/taxonomies/12345678/policyTags/12345678"
								]
							}
						}
					]
				},
				{
					"mode": "NULLABLE",
					"name": "providerphone",
					"type":"STRING",
					"policyTags": {
						"names": [
							"projects/my-project/locations/us/taxonomies/12345678/policyTags/12345678"
						]
					},
					"fields": [
						{
							"name": "value1",
							"type": "INTEGER",
							"mode": "NULLABLE",
							"description": "someVal",
							"policyTags": {
								"names": [
									"projects/my-project/locations/us/taxonomies/12345678/policyTags/12345678"
								]
							}
						},
						{
							"name": "value2",
							"type": "BOOLEAN",
							"mode": "NULLABLE",
							"description": "someVal"
						}
					]
				},
				{
					"name": "PageNo",
					"type": "INTEGER"
				}
			]`,
			ExpectDiffSuppress: true,
		},
	}

	for tn, tc := range cases {
		tn := tn
		tc := tc
		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			var a, b interface{}
			if err := json.Unmarshal([]byte(tc.Old), &a); err != nil {
				t.Fatalf(fmt.Sprintf("unable to unmarshal old json - %v", err))
			}
			if err := json.Unmarshal([]byte(tc.New), &b); err != nil {
				t.Fatalf(fmt.Sprintf("unable to unmarshal new json - %v", err))
			}
			if bigQueryTableSchemaDiffSuppress("schema", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
				t.Fatalf("bad: %s, %q => %q expect DiffSuppress to return %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
			}
		})
	}
}

type testUnitBigQueryDataTableJSONChangeableTestCase struct {
	name       string
	jsonOld    string
	jsonNew    string
	changeable bool
}

func (testcase *testUnitBigQueryDataTableJSONChangeableTestCase) check(t *testing.T) {
	var old, new interface{}
	if err := json.Unmarshal([]byte(testcase.jsonOld), &old); err != nil {
		t.Fatalf("unable to unmarshal json - %v", err)
	}
	if err := json.Unmarshal([]byte(testcase.jsonNew), &new); err != nil {
		t.Fatalf("unable to unmarshal json - %v", err)
	}
	changeable, err := resourceBigQueryTableSchemaIsChangeable(old, new)
	if err != nil {
		t.Errorf("%s failed unexpectedly: %s", testcase.name, err)
	}
	if changeable != testcase.changeable {
		t.Errorf("expected changeable result of %v but got %v for testcase %s", testcase.changeable, changeable, testcase.name)
	}

	d := &tpgresource.ResourceDiffMock{
		Before: map[string]interface{}{},
		After:  map[string]interface{}{},
	}

	d.Before["schema"] = testcase.jsonOld
	d.After["schema"] = testcase.jsonNew

	err = resourceBigQueryTableSchemaCustomizeDiffFunc(d)
	if err != nil {
		t.Errorf("error on testcase %s - %v", testcase.name, err)
	}
	if !testcase.changeable != d.IsForceNew {
		t.Errorf("%s: expected d.IsForceNew to be %v, but was %v", testcase.name, !testcase.changeable, d.IsForceNew)
	}
}

var testUnitBigQueryDataTableIsChangableTestCases = []testUnitBigQueryDataTableJSONChangeableTestCase{
	{
		name:       "defaultEquality",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		changeable: true,
	},
	{
		name:       "arraySizeIncreases",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }, {\"name\": \"asomeValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		changeable: true,
	},
	{
		name:       "arraySizeDecreases",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }, {\"name\": \"asomeValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		changeable: false,
	},
	{
		name:       "descriptionChanges",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"some new value\" }]",
		changeable: true,
	},
	{
		name:       "typeInteger",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"INT64\", \"mode\" : \"NULLABLE\", \"description\" : \"some new value\" }]",
		changeable: true,
	},
	{
		name:       "typeFloat",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"FLOAT\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"FLOAT64\", \"mode\" : \"NULLABLE\", \"description\" : \"some new value\" }]",
		changeable: true,
	},
	{
		name:       "typeBool",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"BOOLEAN\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"BOOL\", \"mode\" : \"NULLABLE\", \"description\" : \"some new value\" }]",
		changeable: true,
	},
	{
		name:       "typeChangeIncompatible",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"BOOLEAN\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"DATETIME\", \"mode\" : \"NULLABLE\", \"description\" : \"some new value\" }]",
		changeable: false,
	},
	// this is invalid but we need to make sure we don't cause a panic
	// if users provide an invalid schema
	{
		name:       "typeChangeIgnoreNewMissingType",
		jsonOld:    "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"BOOLEAN\" }]",
		changeable: true,
	},
	// this is invalid but we need to make sure we don't cause a panic
	// if users provide an invalid schema
	{
		name:       "typeChangeIgnoreOldMissingType",
		jsonOld:    "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"anotherKey\" : \"anotherValue\", \"type\": \"BOOLEAN\" }]",
		changeable: true,
	},
	{
		name:       "typeModeReqToNull",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"BOOLEAN\", \"mode\" : \"REQUIRED\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"BOOLEAN\", \"mode\" : \"NULLABLE\", \"description\" : \"some new value\" }]",
		changeable: true,
	},
	{
		name:       "typeModeIncompatible",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"BOOLEAN\", \"mode\" : \"REQUIRED\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"BOOLEAN\", \"mode\" : \"REPEATED\", \"description\" : \"some new value\" }]",
		changeable: false,
	},
	{
		name:       "modeToDefaultNullable",
		jsonOld:    "[{\"name\": \"someValue\", \"type\" : \"BOOLEAN\", \"mode\" : \"REQUIRED\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"someValue\", \"type\" : \"BOOLEAN\", \"description\" : \"some new value\" }]",
		changeable: true,
	},
	{
		name:       "orderOfArrayChangesAndDescriptionChanges",
		jsonOld:    "[{\"name\": \"value1\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }, {\"name\": \"value2\", \"type\" : \"BOOLEAN\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"value2\", \"type\" : \"BOOLEAN\", \"mode\" : \"NULLABLE\", \"description\" : \"newVal\" },  {\"name\": \"value1\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		changeable: true,
	},
	{
		name:       "orderOfArrayChangesAndNameChanges",
		jsonOld:    "[{\"name\": \"value1\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }, {\"name\": \"value2\", \"type\" : \"BOOLEAN\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		jsonNew:    "[{\"name\": \"value3\", \"type\" : \"BOOLEAN\", \"mode\" : \"NULLABLE\", \"description\" : \"newVal\" },  {\"name\": \"value1\", \"type\" : \"INTEGER\", \"mode\" : \"NULLABLE\", \"description\" : \"someVal\" }]",
		changeable: false,
	},
	{
		name: "policyTags",
		jsonOld: `[
			{
				"mode": "NULLABLE",
				"name": "providerphone",
				"policyTags": {
					"names": ["projects/my-project/locations/us/taxonomies/12345678/policyTags/12345678"]
				},
				"type":"STRING"
			}
		]`,
		jsonNew: `[
			{
				"name": "providerphone",
				"type": "STRING",
				"policyTags": {
					"names": ["projects/my-project/locations/us/taxonomies/12345678/policyTags/12345678"]
				}
			}
		]`,
		changeable: true,
	},
}

func TestUnitBigQueryDataTable_schemaIsChangable(t *testing.T) {
	t.Parallel()
	for _, testcase := range testUnitBigQueryDataTableIsChangableTestCases {
		testcase.check(t)
		testcaseNested := &testUnitBigQueryDataTableJSONChangeableTestCase{
			testcase.name + "Nested",
			fmt.Sprintf("[{\"name\": \"someValue\", \"type\" : \"INTEGER\", \"fields\" : %s }]", testcase.jsonOld),
			fmt.Sprintf("[{\"name\": \"someValue\", \"type\" : \"INT64\", \"fields\" : %s }]", testcase.jsonNew),
			testcase.changeable,
		}
		testcaseNested.check(t)
	}
}
