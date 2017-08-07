package main

import (
	"reflect"
	"testing"

	"google.golang.org/api/discovery/v1"
)

func TestGenerateFields_primitive(t *testing.T) {
	schema := map[string]discovery.JsonSchema{
		"Resource": {
			Type: "object",
			Properties: map[string]discovery.JsonSchema{
				"stringField": {
					Type:        "string",
					Description: "string field",
				},
				"numberField": {
					Type:        "number",
					Description: "Required. number field",
				},
				"intField": {
					Type:        "integer",
					Description: "integer field",
				},
				"boolField": {
					Type:        "boolean",
					Description: "Output-only. boolean field",
				},
				"mapField": {
					Type:        "object",
					Description: "object field",
				},
			},
		},
	}

	reqFields, optFields, comFields := generateFields(schema, "Resource")

	expectedReq := map[string]string{
		"number_field": "{\nType: schema.TypeFloat,\nDescription: \"Required. number field\",\nRequired: true,\nForceNew: true,\n}",
	}

	expectedOpt := map[string]string{
		"string_field": "{\nType: schema.TypeString,\nDescription: \"string field\",\nOptional: true,\nForceNew: true,\n}",
		"int_field":    "{\nType: schema.TypeInt,\nDescription: \"integer field\",\nOptional: true,\nForceNew: true,\n}",
		"map_field":    "{\nType: schema.TypeMap,\nDescription: \"object field\",\nOptional: true,\nForceNew: true,\n}",
	}

	expectedCom := map[string]string{
		"bool_field": "{\nType: schema.TypeBool,\nDescription: \"Output-only. boolean field\",\nForceNew: true,\nComputed: true,\n}",
	}

	if !reflect.DeepEqual(reqFields, expectedReq) {
		t.Fatalf("Expected: %+v\n\nGiven: %+v\n", expectedReq, reqFields)
	}
	if !reflect.DeepEqual(optFields, expectedOpt) {
		t.Fatalf("Expected: %+v\n\nGiven: %+v\n", expectedOpt, optFields)
	}
	if !reflect.DeepEqual(comFields, expectedCom) {
		t.Fatalf("Expected: %+v\n\nGiven: %+v\n", expectedCom, comFields)
	}
}

func TestGenerateFields_listOfPrimitives(t *testing.T) {
	schema := map[string]discovery.JsonSchema{
		"Resource": {
			Type: "object",
			Properties: map[string]discovery.JsonSchema{
				"stringsField": {
					Type: "array",
					Items: &discovery.JsonSchema{
						Type: "string",
					},
				},
				"numbersField": {
					Type: "array",
					Items: &discovery.JsonSchema{
						Type: "number",
					},
				},
				"intsField": {
					Type: "array",
					Items: &discovery.JsonSchema{
						Type: "integer",
					},
				},
				"boolsField": {
					Type: "array",
					Items: &discovery.JsonSchema{
						Type: "boolean",
					},
				},
			},
		},
	}

	_, optFields, _ := generateFields(schema, "Resource")

	expected := map[string]string{
		"strings_field": "{\nType: schema.TypeList,\nOptional: true,\nForceNew: true,\nElem: &schema.Schema{Type: schema.TypeString,},\n}",
		"numbers_field": "{\nType: schema.TypeList,\nOptional: true,\nForceNew: true,\nElem: &schema.Schema{Type: schema.TypeFloat,},\n}",
		"ints_field":    "{\nType: schema.TypeList,\nOptional: true,\nForceNew: true,\nElem: &schema.Schema{Type: schema.TypeInt,},\n}",
		"bools_field":   "{\nType: schema.TypeList,\nOptional: true,\nForceNew: true,\nElem: &schema.Schema{Type: schema.TypeBool,},\n}",
	}

	if !reflect.DeepEqual(optFields, expected) {
		t.Fatalf("Expected: %+v\n\nGiven: %+v\n", expected, optFields)
	}
}

func TestGenerateFields_nested(t *testing.T) {
	schema := map[string]discovery.JsonSchema{
		"Resource": {
			Type: "object",
			Properties: map[string]discovery.JsonSchema{
				"nestedField": {
					Ref: "OtherThing",
				},
			},
		},
		"OtherThing": {
			Type: "object",
			Properties: map[string]discovery.JsonSchema{
				"intField": {
					Type: "integer",
				},
				"stringField": {
					Type: "string",
				},
			},
		},
	}

	_, optFields, _ := generateFields(schema, "Resource")

	expected := map[string]string{
		"nested_field": "{\nType: schema.TypeList,\nOptional: true,\nForceNew: true,\nMaxItems: 1,\nElem: &schema.Resource{\nSchema: map[string]*schema.Schema{\n\"int_field\": {\nType: schema.TypeInt,\nOptional: true,\nForceNew: true,\n},\n\"string_field\": {\nType: schema.TypeString,\nOptional: true,\nForceNew: true,\n},\n},\n},\n}",
	}

	if !reflect.DeepEqual(optFields, expected) {
		t.Fatalf("Expected: %+v\n\nGiven: %+v\n", expected, optFields)
	}
}

func TestGenerateFields_nestedList(t *testing.T) {
	schema := map[string]discovery.JsonSchema{
		"Resource": {
			Type: "object",
			Properties: map[string]discovery.JsonSchema{
				"nestedField": {
					Type: "array",
					Items: &discovery.JsonSchema{
						Ref: "OtherThing",
					},
				},
			},
		},
		"OtherThing": {
			Type: "object",
			Properties: map[string]discovery.JsonSchema{
				"intField": {
					Type: "integer",
				},
				"stringField": {
					Type: "string",
				},
			},
		},
	}

	_, optFields, _ := generateFields(schema, "Resource")

	expected := map[string]string{
		"nested_field": "{\nType: schema.TypeList,\nOptional: true,\nForceNew: true,\nElem: &schema.Resource{\nSchema: map[string]*schema.Schema{\n\"int_field\": {\nType: schema.TypeInt,\nOptional: true,\nForceNew: true,\n},\n\"string_field\": {\nType: schema.TypeString,\nOptional: true,\nForceNew: true,\n},\n},\n},\n}",
	}

	if !reflect.DeepEqual(optFields, expected) {
		t.Fatalf("Expected: %+v\n\nGiven: %+v\n", expected, optFields)
	}
}

func TestUnderscore(t *testing.T) {
	testCases := map[string]string{
		"camelCase":           "camel_case",
		"CamelCase":           "camel_case",
		"HTTPResponseCode":    "http_response_code",
		"HTTPResponseCodeXYZ": "http_response_code_xyz",
		"getHTTPResponseCode": "get_http_response_code",
		"ISCSI":               "iscsi",
		"externalIPs":         "external_ips",
	}

	for from, to := range testCases {
		converted := underscore(from)
		if converted != to {
			t.Fatalf("Expected %q after conversion, given: %q", to, converted)
		}
	}
}
