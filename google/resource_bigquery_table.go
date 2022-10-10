package google

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"google.golang.org/api/bigquery/v2"
)

func bigQueryTableSortArrayByName(array []interface{}) {
	sort.Slice(array, func(i, k int) bool {
		return array[i].(map[string]interface{})["name"].(string) < array[k].(map[string]interface{})["name"].(string)
	})
}

func bigQueryArrayToMapIndexedByName(array []interface{}) map[string]interface{} {
	out := map[string]interface{}{}
	for _, v := range array {
		name := v.(map[string]interface{})["name"].(string)
		out[name] = v
	}
	return out
}

func bigQueryTablecheckNameExists(jsonList []interface{}) error {
	for _, m := range jsonList {
		if _, ok := m.(map[string]interface{})["name"]; !ok {
			return fmt.Errorf("No name in schema %+v", m)
		}
	}

	return nil
}

// Compares two json's while optionally taking in a compareMapKeyVal function.
// This function will override any comparison of a given map[string]interface{}
// on a specific key value allowing for a separate equality in specific scenarios
func jsonCompareWithMapKeyOverride(key string, a, b interface{}, compareMapKeyVal func(key string, val1, val2 map[string]interface{}) bool) (bool, error) {
	switch a.(type) {
	case []interface{}:
		arrayA := a.([]interface{})
		arrayB, ok := b.([]interface{})
		if !ok {
			return false, nil
		} else if len(arrayA) != len(arrayB) {
			return false, nil
		}

		// Sort fields by name so reordering them doesn't cause a diff.
		if key == "schema" || key == "fields" {
			if err := bigQueryTablecheckNameExists(arrayA); err != nil {
				return false, err
			}
			bigQueryTableSortArrayByName(arrayA)
			if err := bigQueryTablecheckNameExists(arrayB); err != nil {
				return false, err
			}
			bigQueryTableSortArrayByName(arrayB)
		}
		for i := range arrayA {
			eq, err := jsonCompareWithMapKeyOverride(strconv.Itoa(i), arrayA[i], arrayB[i], compareMapKeyVal)
			if err != nil {
				return false, err
			} else if !eq {
				return false, nil
			}
		}
		return true, nil
	case map[string]interface{}:
		objectA := a.(map[string]interface{})
		objectB, ok := b.(map[string]interface{})
		if !ok {
			return false, nil
		}

		var unionOfKeys map[string]bool = make(map[string]bool)
		for subKey := range objectA {
			unionOfKeys[subKey] = true
		}
		for subKey := range objectB {
			unionOfKeys[subKey] = true
		}

		for subKey := range unionOfKeys {
			eq := compareMapKeyVal(subKey, objectA, objectB)
			if !eq {
				valA, ok1 := objectA[subKey]
				valB, ok2 := objectB[subKey]
				if !ok1 || !ok2 {
					return false, nil
				}
				eq, err := jsonCompareWithMapKeyOverride(subKey, valA, valB, compareMapKeyVal)
				if err != nil || !eq {
					return false, err
				}
			}
		}
		return true, nil
	case string, float64, bool, nil:
		return a == b, nil
	default:
		log.Printf("[DEBUG] tried to iterate through json but encountered a non native type to json deserialization... please ensure you are passing a json object from json.Unmarshall")
		return false, errors.New("unable to compare values")
	}
}

// checks if the value is within the array, only works for generics
// because objects and arrays will take the reference comparison
func valueIsInArray(value interface{}, array []interface{}) bool {
	for _, item := range array {
		if item == value {
			return true
		}
	}
	return false
}

func bigQueryTableMapKeyOverride(key string, objectA, objectB map[string]interface{}) bool {
	// we rely on the fallback to nil if the object does not have the key
	valA := objectA[key]
	valB := objectB[key]
	switch key {
	case "mode":
		eq := bigQueryTableNormalizeMode(valA) == bigQueryTableNormalizeMode(valB)
		return eq
	case "description":
		equivalentSet := []interface{}{nil, ""}
		eq := valueIsInArray(valA, equivalentSet) && valueIsInArray(valB, equivalentSet)
		return eq
	case "type":
		if valA == nil || valB == nil {
			return false
		}
		return bigQueryTableTypeEq(valA.(string), valB.(string))
	}

	// otherwise rely on default behavior
	return false
}

// Compare the JSON strings are equal
func bigQueryTableSchemaDiffSuppress(name, old, new string, _ *schema.ResourceData) bool {
	// The API can return an empty schema which gets encoded to "null" during read.
	if old == "null" {
		old = "[]"
	}
	var a, b interface{}
	if err := json.Unmarshal([]byte(old), &a); err != nil {
		log.Printf("[DEBUG] unable to unmarshal old json - %v", err)
	}
	if err := json.Unmarshal([]byte(new), &b); err != nil {
		log.Printf("[DEBUG] unable to unmarshal new json - %v", err)
	}

	eq, err := jsonCompareWithMapKeyOverride(name, a, b, bigQueryTableMapKeyOverride)
	if err != nil {
		log.Printf("[DEBUG] %v", err)
		log.Printf("[DEBUG] Error comparing JSON: %v, %v", old, new)
	}

	return eq
}

func bigQueryTableTypeEq(old, new string) bool {
	// Do case-insensitive comparison. https://github.com/hashicorp/terraform-provider-google/issues/9472
	oldUpper := strings.ToUpper(old)
	newUpper := strings.ToUpper(new)

	equivalentSet1 := []interface{}{"INTEGER", "INT64"}
	equivalentSet2 := []interface{}{"FLOAT", "FLOAT64"}
	equivalentSet3 := []interface{}{"BOOLEAN", "BOOL"}
	eq0 := oldUpper == newUpper
	eq1 := valueIsInArray(oldUpper, equivalentSet1) && valueIsInArray(newUpper, equivalentSet1)
	eq2 := valueIsInArray(oldUpper, equivalentSet2) && valueIsInArray(newUpper, equivalentSet2)
	eq3 := valueIsInArray(oldUpper, equivalentSet3) && valueIsInArray(newUpper, equivalentSet3)
	eq := eq0 || eq1 || eq2 || eq3
	return eq
}

func bigQueryTableNormalizeMode(mode interface{}) string {
	if mode == nil {
		return "NULLABLE"
	}
	// Upper-case to get case-insensitive comparisons. https://github.com/hashicorp/terraform-provider-google/issues/9472
	return strings.ToUpper(mode.(string))
}

func bigQueryTableModeIsForceNew(old, new string) bool {
	eq := old == new
	reqToNull := old == "REQUIRED" && new == "NULLABLE"
	return !eq && !reqToNull
}

// Compares two existing schema implementations and decides if
// it is changeable.. pairs with a force new on not changeable
func resourceBigQueryTableSchemaIsChangeable(old, new interface{}) (bool, error) {
	switch old.(type) {
	case []interface{}:
		arrayOld := old.([]interface{})
		arrayNew, ok := new.([]interface{})
		if !ok {
			// if not both arrays not changeable
			return false, nil
		}
		if len(arrayOld) > len(arrayNew) {
			// if not growing not changeable
			return false, nil
		}
		if err := bigQueryTablecheckNameExists(arrayOld); err != nil {
			return false, err
		}
		mapOld := bigQueryArrayToMapIndexedByName(arrayOld)
		if err := bigQueryTablecheckNameExists(arrayNew); err != nil {
			return false, err
		}
		mapNew := bigQueryArrayToMapIndexedByName(arrayNew)
		for key := range mapNew {
			// making unchangeable if an newly added column is with REQUIRED mode
			if _, ok := mapOld[key]; !ok {
				items := mapNew[key].(map[string]interface{})
				for k := range items {
					if k == "mode" && fmt.Sprintf("%v", items[k]) == "REQUIRED" {
						return false, nil
					}
				}
			}
		}
		for key := range mapOld {
			// all old keys should be represented in the new config
			if _, ok := mapNew[key]; !ok {
				return false, nil
			}
			if isChangable, err :=
				resourceBigQueryTableSchemaIsChangeable(mapOld[key], mapNew[key]); err != nil || !isChangable {
				return false, err
			}
		}
		return true, nil
	case map[string]interface{}:
		objectOld := old.(map[string]interface{})
		objectNew, ok := new.(map[string]interface{})
		if !ok {
			// if both aren't objects
			return false, nil
		}
		var unionOfKeys map[string]bool = make(map[string]bool)
		for key := range objectOld {
			unionOfKeys[key] = true
		}
		for key := range objectNew {
			unionOfKeys[key] = true
		}
		for key := range unionOfKeys {
			valOld := objectOld[key]
			valNew := objectNew[key]
			switch key {
			case "name":
				if valOld != valNew {
					return false, nil
				}
			case "type":
				if valOld == nil || valNew == nil {
					// This is invalid, so it shouldn't require a ForceNew
					return true, nil
				}
				if !bigQueryTableTypeEq(valOld.(string), valNew.(string)) {
					return false, nil
				}
			case "mode":
				if bigQueryTableModeIsForceNew(
					bigQueryTableNormalizeMode(valOld),
					bigQueryTableNormalizeMode(valNew),
				) {
					return false, nil
				}
			case "fields":
				return resourceBigQueryTableSchemaIsChangeable(valOld, valNew)

				// other parameters: description, policyTags and
				// policyTags.names[] are changeable
			}
		}
		return true, nil
	case string, float64, bool, nil:
		// realistically this shouldn't hit
		log.Printf("[DEBUG] comparison of generics hit... not expected")
		return old == new, nil
	default:
		log.Printf("[DEBUG] tried to iterate through json but encountered a non native type to json deserialization... please ensure you are passing a json object from json.Unmarshall")
		return false, errors.New("unable to compare values")
	}
}

func resourceBigQueryTableSchemaCustomizeDiffFunc(d TerraformResourceDiff) error {
	if _, hasSchema := d.GetOk("schema"); hasSchema {
		oldSchema, newSchema := d.GetChange("schema")
		oldSchemaText := oldSchema.(string)
		newSchemaText := newSchema.(string)
		if oldSchemaText == "null" {
			// The API can return an empty schema which gets encoded to "null" during read.
			oldSchemaText = "[]"
		}
		if newSchemaText == "null" {
			newSchemaText = "[]"
		}
		var old, new interface{}
		if err := json.Unmarshal([]byte(oldSchemaText), &old); err != nil {
			// don't return error, its possible we are going from no schema to schema
			// this case will be cover on the conparision regardless.
			log.Printf("[DEBUG] unable to unmarshal json customized diff - %v", err)
		}
		if err := json.Unmarshal([]byte(newSchemaText), &new); err != nil {
			// same as above
			log.Printf("[DEBUG] unable to unmarshal json customized diff - %v", err)
		}
		isChangeable, err := resourceBigQueryTableSchemaIsChangeable(old, new)
		if err != nil {
			return err
		}
		if !isChangeable {
			if err := d.ForceNew("schema"); err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

func resourceBigQueryTableSchemaCustomizeDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	return resourceBigQueryTableSchemaCustomizeDiffFunc(d)
}

func resourceBigQueryTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigQueryTableCreate,
		Read:   resourceBigQueryTableRead,
		Delete: resourceBigQueryTableDelete,
		Update: resourceBigQueryTableUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceBigQueryTableImport,
		},
		CustomizeDiff: customdiff.All(
			resourceBigQueryTableSchemaCustomizeDiff,
		),
		Schema: map[string]*schema.Schema{
			// TableId: [Required] The ID of the table. The ID must contain only
			// letters (a-z, A-Z), numbers (0-9), or underscores (_). The maximum
			// length is 1,024 characters.
			"table_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `A unique ID for the resource. Changing this forces a new resource to be created.`,
			},

			// DatasetId: [Required] The ID of the dataset containing this table.
			"dataset_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The dataset ID to create the table in. Changing this forces a new resource to be created.`,
			},

			// ProjectId: [Required] The ID of the project containing this table.
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs.`,
			},

			// Description: [Optional] A user-friendly description of this table.
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The field description.`,
			},

			// ExpirationTime: [Optional] The time when this table expires, in
			// milliseconds since the epoch. If not present, the table will persist
			// indefinitely. Expired tables will be deleted and their storage
			// reclaimed.
			"expiration_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The time when this table expires, in milliseconds since the epoch. If not present, the table will persist indefinitely. Expired tables will be deleted and their storage reclaimed.`,
			},

			// ExternalDataConfiguration [Optional] Describes the data format,
			// location, and other properties of a table stored outside of BigQuery.
			// By defining these properties, the data source can then be queried as
			// if it were a standard BigQuery table.
			"external_data_configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `Describes the data format, location, and other properties of a table stored outside of BigQuery. By defining these properties, the data source can then be queried as if it were a standard BigQuery table.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Autodetect : [Required] If true, let BigQuery try to autodetect the
						// schema and format of the table.
						"autodetect": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Let BigQuery try to autodetect the schema and format of the table.`,
						},
						// SourceFormat [Required] The data format.
						"source_format": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The data format. Supported values are: "CSV", "GOOGLE_SHEETS", "NEWLINE_DELIMITED_JSON", "AVRO", "PARQUET", "ORC" and "DATASTORE_BACKUP". To use "GOOGLE_SHEETS" the scopes must include "googleapis.com/auth/drive.readonly".`,
							ValidateFunc: validation.StringInSlice([]string{
								"CSV", "GOOGLE_SHEETS", "NEWLINE_DELIMITED_JSON", "AVRO", "DATASTORE_BACKUP", "PARQUET", "ORC", "BIGTABLE",
							}, false),
						},
						// SourceURIs [Required] The fully-qualified URIs that point to your data in Google Cloud.
						"source_uris": {
							Type:        schema.TypeList,
							Required:    true,
							Description: `A list of the fully-qualified URIs that point to your data in Google Cloud.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						// Compression: [Optional] The compression type of the data source.
						"compression": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"NONE", "GZIP"}, false),
							Default:      "NONE",
							Description:  `The compression type of the data source. Valid values are "NONE" or "GZIP".`,
						},
						// Schema: Optional] The schema for the  data.
						// Schema is required for CSV and JSON formats if autodetect is not on.
						// Schema is disallowed for Google Cloud Bigtable, Cloud Datastore backups, Avro, ORC and Parquet formats.
						"schema": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsJSON,
							StateFunc: func(v interface{}) string {
								json, _ := structure.NormalizeJsonString(v)
								return json
							},
							Description: `A JSON schema for the external table. Schema is required for CSV and JSON formats and is disallowed for Google Cloud Bigtable, Cloud Datastore backups, and Avro formats when using external tables.`,
						},
						// CsvOptions: [Optional] Additional properties to set if
						// sourceFormat is set to CSV.
						"csv_options": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: `Additional properties to set if source_format is set to "CSV".`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// Quote: [Required] The value that is used to quote data
									// sections in a CSV file.
									"quote": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The value that is used to quote data sections in a CSV file. If your data does not contain quoted sections, set the property value to an empty string. If your data contains quoted newline characters, you must also set the allow_quoted_newlines property to true. The API-side default is ", specified in Terraform escaped as \". Due to limitations with Terraform default values, this value is required to be explicitly set.`,
									},
									// AllowJaggedRows: [Optional] Indicates if BigQuery should
									// accept rows that are missing trailing optional columns.
									"allow_jagged_rows": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: `Indicates if BigQuery should accept rows that are missing trailing optional columns.`,
									},
									// AllowQuotedNewlines: [Optional] Indicates if BigQuery
									// should allow quoted data sections that contain newline
									// characters in a CSV file. The default value is false.
									"allow_quoted_newlines": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
										Description: `Indicates if BigQuery should allow quoted data sections that contain newline characters in a CSV file. The default value is false.`,
									},
									// Encoding: [Optional] The character encoding of the data.
									// The supported values are UTF-8 or ISO-8859-1.
									"encoding": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"ISO-8859-1", "UTF-8"}, false),
										Default:      "UTF-8",
										Description:  `The character encoding of the data. The supported values are UTF-8 or ISO-8859-1.`,
									},
									// FieldDelimiter: [Optional] The separator for fields in a CSV file.
									"field_delimiter": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     ",",
										Description: `The separator for fields in a CSV file.`,
									},
									// SkipLeadingRows: [Optional] The number of rows at the top
									// of a CSV file that BigQuery will skip when reading the data.
									"skip_leading_rows": {
										Type:        schema.TypeInt,
										Optional:    true,
										Default:     0,
										Description: `The number of rows at the top of a CSV file that BigQuery will skip when reading the data.`,
									},
								},
							},
						},
						// GoogleSheetsOptions: [Optional] Additional options if sourceFormat is set to GOOGLE_SHEETS.
						"google_sheets_options": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: `Additional options if source_format is set to "GOOGLE_SHEETS".`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// Range: [Optional] Range of a sheet to query from. Only used when non-empty.
									// Typical format: !:
									"range": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Range of a sheet to query from. Only used when non-empty. At least one of range or skip_leading_rows must be set. Typical format: "sheet_name!top_left_cell_id:bottom_right_cell_id" For example: "sheet1!A1:B20"`,
										AtLeastOneOf: []string{
											"external_data_configuration.0.google_sheets_options.0.skip_leading_rows",
											"external_data_configuration.0.google_sheets_options.0.range",
										},
									},
									// SkipLeadingRows: [Optional] The number of rows at the top
									// of the sheet that BigQuery will skip when reading the data.
									"skip_leading_rows": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `The number of rows at the top of the sheet that BigQuery will skip when reading the data. At least one of range or skip_leading_rows must be set.`,
										AtLeastOneOf: []string{
											"external_data_configuration.0.google_sheets_options.0.skip_leading_rows",
											"external_data_configuration.0.google_sheets_options.0.range",
										},
									},
								},
							},
						},

						// HivePartitioningOptions:: [Optional] Options for configuring hive partitioning detect.
						"hive_partitioning_options": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: `When set, configures hive partitioning support. Not all storage formats support hive partitioning -- requesting hive partitioning on an unsupported format will lead to an error, as will providing an invalid specification.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// Mode: [Optional] [Experimental] When set, what mode of hive partitioning to use when reading data.
									// Two modes are supported.
									//* AUTO: automatically infer partition key name(s) and type(s).
									//* STRINGS: automatically infer partition key name(s).
									"mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `When set, what mode of hive partitioning to use when reading data.`,
									},
									// RequirePartitionFilter: [Optional] If set to true, queries over this table
									// require a partition filter that can be used for partition elimination to be
									// specified.
									"require_partition_filter": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `If set to true, queries over this table require a partition filter that can be used for partition elimination to be specified.`,
									},
									// SourceUriPrefix: [Optional] [Experimental] When hive partition detection is requested, a common for all source uris must be required.
									// The prefix must end immediately before the partition key encoding begins.
									"source_uri_prefix": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `When hive partition detection is requested, a common for all source uris must be required. The prefix must end immediately before the partition key encoding begins.`,
									},
								},
							},
						},
						// AvroOptions: [Optional] Additional options if sourceFormat is set to AVRO.
						"avro_options": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: `Additional options if source_format is set to "AVRO"`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"use_avro_logical_types": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: `If sourceFormat is set to "AVRO", indicates whether to interpret logical types as the corresponding BigQuery data type (for example, TIMESTAMP), instead of using the raw type (for example, INTEGER).`,
									},
								},
							},
						},

						// IgnoreUnknownValues: [Optional] Indicates if BigQuery should
						// allow extra values that are not represented in the table schema.
						// If true, the extra values are ignored. If false, records with
						// extra columns are treated as bad records, and if there are too
						// many bad records, an invalid error is returned in the job result.
						// The default value is false.
						"ignore_unknown_values": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Indicates if BigQuery should allow extra values that are not represented in the table schema. If true, the extra values are ignored. If false, records with extra columns are treated as bad records, and if there are too many bad records, an invalid error is returned in the job result. The default value is false.`,
						},
						// MaxBadRecords: [Optional] The maximum number of bad records that
						// BigQuery can ignore when reading data.
						"max_bad_records": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `The maximum number of bad records that BigQuery can ignore when reading data.`,
						},
						// ConnectionId: [Optional] The connection specifying the credentials
						// to be used to read external storage, such as Azure Blob,
						// Cloud Storage, or S3. The connectionId can have the form
						// "{{project}}.{{location}}.{{connection_id}}" or
						// "projects/{{project}}/locations/{{location}}/connections/{{connection_id}}".
						"connection_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The connection specifying the credentials to be used to read external storage, such as Azure Blob, Cloud Storage, or S3. The connectionId can have the form "{{project}}.{{location}}.{{connection_id}}" or "projects/{{project}}/locations/{{location}}/connections/{{connection_id}}".`,
						},
					},
				},
			},

			// FriendlyName: [Optional] A descriptive name for this table.
			"friendly_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `A descriptive name for the table.`,
			},

			// Labels: [Experimental] The labels associated with this table. You can
			// use these to organize and group your tables. Label keys and values
			// can be no longer than 63 characters, can only contain lowercase
			// letters, numeric characters, underscores and dashes. International
			// characters are allowed. Label values are optional. Label keys must
			// start with a letter and each label in the list must have a different
			// key.
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `A mapping of labels to assign to the resource.`,
			},

			// Schema: [Optional] Describes the schema of this table.
			"schema": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsJSON,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
				DiffSuppressFunc: bigQueryTableSchemaDiffSuppress,
				Description:      `A JSON schema for the table.`,
			},
			// View: [Optional] If specified, configures this table as a view.
			"view": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `If specified, configures this table as a view.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Query: [Required] A query that BigQuery executes when the view is
						// referenced.
						"query": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `A query that BigQuery executes when the view is referenced.`,
						},

						// UseLegacySQL: [Optional] Specifies whether to use BigQuery's
						// legacy SQL for this view. The default value is true. If set to
						// false, the view will use BigQuery's standard SQL:
						"use_legacy_sql": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: `Specifies whether to use BigQuery's legacy SQL for this view. The default value is true. If set to false, the view will use BigQuery's standard SQL`,
						},
					},
				},
			},

			// Materialized View: [Optional] If specified, configures this table as a materialized view.
			"materialized_view": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `If specified, configures this table as a materialized view.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// EnableRefresh: [Optional] Enable automatic refresh of
						// the materialized view when the base table is updated. The default
						// value is "true".
						"enable_refresh": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: `Specifies if BigQuery should automatically refresh materialized view when the base table is updated. The default is true.`,
						},

						// RefreshIntervalMs: [Optional] The maximum frequency
						// at which this materialized view will be refreshed. The default value
						// is 1800000 (30 minutes).
						"refresh_interval_ms": {
							Type:        schema.TypeInt,
							Default:     1800000,
							Optional:    true,
							Description: `Specifies maximum frequency at which this materialized view will be refreshed. The default is 1800000`,
						},

						// Query: [Required] A query whose result is persisted
						"query": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: `A query whose result is persisted.`,
						},
					},
				},
			},

			// TimePartitioning: [Experimental] If specified, configures time-based
			// partitioning for this table.
			"time_partitioning": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `If specified, configures time-based partitioning for this table.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// ExpirationMs: [Optional] Number of milliseconds for which to keep the
						// storage for a partition.
						"expiration_ms": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: `Number of milliseconds for which to keep the storage for a partition.`,
						},

						// Type: [Required] The supported types are DAY, HOUR, MONTH, and YEAR, which will generate
						// one partition per day, hour, month, and year, respectively.
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  `The supported types are DAY, HOUR, MONTH, and YEAR, which will generate one partition per day, hour, month, and year, respectively.`,
							ValidateFunc: validation.StringInSlice([]string{"DAY", "HOUR", "MONTH", "YEAR"}, false),
						},

						// Field: [Optional] The field used to determine how to create a time-based
						// partition. If time-based partitioning is enabled without this value, the
						// table is partitioned based on the load time.
						"field": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `The field used to determine how to create a time-based partition. If time-based partitioning is enabled without this value, the table is partitioned based on the load time.`,
						},

						// RequirePartitionFilter: [Optional] If set to true, queries over this table
						// require a partition filter that can be used for partition elimination to be
						// specified.
						"require_partition_filter": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `If set to true, queries over this table require a partition filter that can be used for partition elimination to be specified.`,
						},
					},
				},
			},

			// RangePartitioning: [Optional] If specified, configures range-based
			// partitioning for this table.
			"range_partitioning": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `If specified, configures range-based partitioning for this table.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Field: [Required] The field used to determine how to create a range-based
						// partition.
						"field": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: `The field used to determine how to create a range-based partition.`,
						},

						// Range: [Required] Information required to partition based on ranges.
						"range": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: `Information required to partition based on ranges. Structure is documented below.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// Start: [Required] Start of the range partitioning, inclusive.
									"start": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: `Start of the range partitioning, inclusive.`,
									},

									// End: [Required] End of the range partitioning, exclusive.
									"end": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: `End of the range partitioning, exclusive.`,
									},

									// Interval: [Required] The width of each range within the partition.
									"interval": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: `The width of each range within the partition.`,
									},
								},
							},
						},
					},
				},
			},

			// Clustering: [Optional] Specifies column names to use for data clustering.  Up to four
			// top-level columns are allowed, and should be specified in descending priority order.
			"clustering": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    4,
				Description: `Specifies column names to use for data clustering. Up to four top-level columns are allowed, and should be specified in descending priority order.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"encryption_configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: `Specifies how the table should be encrypted. If left blank, the table will be encrypted with a Google-managed key; that process is transparent to the user.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kms_key_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The self link or full name of a key which should be used to encrypt this table. Note that the default bigquery service account will need to have encrypt/decrypt permissions on this key - you may want to see the google_bigquery_default_service_account datasource and the google_kms_crypto_key_iam_binding resource.`,
						},
						"kms_key_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The self link or full name of the kms key version used to encrypt this table.`,
						},
					},
				},
			},

			// CreationTime: [Output-only] The time when this table was created, in
			// milliseconds since the epoch.
			"creation_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The time when this table was created, in milliseconds since the epoch.`,
			},

			// Etag: [Output-only] A hash of this resource.
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `A hash of the resource.`,
			},

			// LastModifiedTime: [Output-only] The time when this table was last
			// modified, in milliseconds since the epoch.
			"last_modified_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The time when this table was last modified, in milliseconds since the epoch.`,
			},

			// Location: [Output-only] The geographic location where the table
			// resides. This value is inherited from the dataset.
			"location": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The geographic location where the table resides. This value is inherited from the dataset.`,
			},

			// NumBytes: [Output-only] The size of this table in bytes, excluding
			// any data in the streaming buffer.
			"num_bytes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The geographic location where the table resides. This value is inherited from the dataset.`,
			},

			// NumLongTermBytes: [Output-only] The number of bytes in the table that
			// are considered "long-term storage".
			"num_long_term_bytes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of bytes in the table that are considered "long-term storage".`,
			},

			// NumRows: [Output-only] The number of rows of data in this table,
			// excluding any data in the streaming buffer.
			"num_rows": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of rows of data in this table, excluding any data in the streaming buffer.`,
			},

			// SelfLink: [Output-only] A URL that can be used to access this
			// resource again.
			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URI of the created resource.`,
			},

			// Type: [Output-only] Describes the table type. The following values
			// are supported: TABLE: A normal BigQuery table. VIEW: A virtual table
			// defined by a SQL query. EXTERNAL: A table that references data stored
			// in an external storage system, such as Google Cloud Storage. The
			// default value is TABLE.
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Describes the table type.`,
			},

			"deletion_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Whether or not to allow Terraform to destroy the instance. Unless this field is set to false in Terraform state, a terraform destroy or terraform apply that would delete the instance will fail.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceTable(d *schema.ResourceData, meta interface{}) (*bigquery.Table, error) {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	table := &bigquery.Table{
		TableReference: &bigquery.TableReference{
			DatasetId: d.Get("dataset_id").(string),
			TableId:   d.Get("table_id").(string),
			ProjectId: project,
		},
	}

	if v, ok := d.GetOk("view"); ok {
		table.View = expandView(v)
	}

	if v, ok := d.GetOk("materialized_view"); ok {
		table.MaterializedView = expandMaterializedView(v)
	}

	if v, ok := d.GetOk("description"); ok {
		table.Description = v.(string)
	}

	if v, ok := d.GetOk("expiration_time"); ok {
		table.ExpirationTime = int64(v.(int))
	}

	if v, ok := d.GetOk("external_data_configuration"); ok {
		externalDataConfiguration, err := expandExternalDataConfiguration(v)
		if err != nil {
			return nil, err
		}

		table.ExternalDataConfiguration = externalDataConfiguration
	}

	if v, ok := d.GetOk("friendly_name"); ok {
		table.FriendlyName = v.(string)
	}

	if v, ok := d.GetOk("encryption_configuration.0.kms_key_name"); ok {
		table.EncryptionConfiguration = &bigquery.EncryptionConfiguration{
			KmsKeyName: v.(string),
		}
	}

	if v, ok := d.GetOk("labels"); ok {
		labels := map[string]string{}

		for k, v := range v.(map[string]interface{}) {
			labels[k] = v.(string)
		}

		table.Labels = labels
	}

	if v, ok := d.GetOk("schema"); ok {
		schema, err := expandSchema(v)
		if err != nil {
			return nil, err
		}
		table.Schema = schema
	}

	if v, ok := d.GetOk("time_partitioning"); ok {
		table.TimePartitioning = expandTimePartitioning(v)
	}

	if v, ok := d.GetOk("range_partitioning"); ok {
		rangePartitioning, err := expandRangePartitioning(v)
		if err != nil {
			return nil, err
		}

		table.RangePartitioning = rangePartitioning
	}

	if v, ok := d.GetOk("clustering"); ok {
		table.Clustering = &bigquery.Clustering{
			Fields:          convertStringArr(v.([]interface{})),
			ForceSendFields: []string{"Fields"},
		}
	}

	return table, nil
}

func resourceBigQueryTableCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	table, err := resourceTable(d, meta)
	if err != nil {
		return err
	}

	datasetID := d.Get("dataset_id").(string)

	if table.View != nil && table.Schema != nil {

		log.Printf("[INFO] Removing schema from table definition because big query does not support setting schema on view creation")
		schemaBack := table.Schema
		table.Schema = nil

		log.Printf("[INFO] Creating BigQuery table: %s without schema", table.TableReference.TableId)

		res, err := config.NewBigQueryClient(userAgent).Tables.Insert(project, datasetID, table).Do()
		if err != nil {
			return err
		}

		log.Printf("[INFO] BigQuery table %s has been created", res.Id)
		d.SetId(fmt.Sprintf("projects/%s/datasets/%s/tables/%s", res.TableReference.ProjectId, res.TableReference.DatasetId, res.TableReference.TableId))

		table.Schema = schemaBack
		log.Printf("[INFO] Updating BigQuery table: %s with schema", table.TableReference.TableId)
		if _, err = config.NewBigQueryClient(userAgent).Tables.Update(project, datasetID, res.TableReference.TableId, table).Do(); err != nil {
			return err
		}

		log.Printf("[INFO] BigQuery table %s has been update with schema", res.Id)
	} else {
		log.Printf("[INFO] Creating BigQuery table: %s", table.TableReference.TableId)

		res, err := config.NewBigQueryClient(userAgent).Tables.Insert(project, datasetID, table).Do()
		if err != nil {
			return err
		}

		log.Printf("[INFO] BigQuery table %s has been created", res.Id)
		d.SetId(fmt.Sprintf("projects/%s/datasets/%s/tables/%s", res.TableReference.ProjectId, res.TableReference.DatasetId, res.TableReference.TableId))
	}

	return resourceBigQueryTableRead(d, meta)
}

func resourceBigQueryTableRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading BigQuery table: %s", d.Id())

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	datasetID := d.Get("dataset_id").(string)
	tableID := d.Get("table_id").(string)

	res, err := config.NewBigQueryClient(userAgent).Tables.Get(project, datasetID, tableID).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("BigQuery table %q", tableID))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("description", res.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("expiration_time", res.ExpirationTime); err != nil {
		return fmt.Errorf("Error setting expiration_time: %s", err)
	}
	if err := d.Set("friendly_name", res.FriendlyName); err != nil {
		return fmt.Errorf("Error setting friendly_name: %s", err)
	}
	if err := d.Set("labels", res.Labels); err != nil {
		return fmt.Errorf("Error setting labels: %s", err)
	}
	if err := d.Set("creation_time", res.CreationTime); err != nil {
		return fmt.Errorf("Error setting creation_time: %s", err)
	}
	if err := d.Set("etag", res.Etag); err != nil {
		return fmt.Errorf("Error setting etag: %s", err)
	}
	if err := d.Set("last_modified_time", res.LastModifiedTime); err != nil {
		return fmt.Errorf("Error setting last_modified_time: %s", err)
	}
	if err := d.Set("location", res.Location); err != nil {
		return fmt.Errorf("Error setting location: %s", err)
	}
	if err := d.Set("num_bytes", res.NumBytes); err != nil {
		return fmt.Errorf("Error setting num_bytes: %s", err)
	}
	if err := d.Set("table_id", res.TableReference.TableId); err != nil {
		return fmt.Errorf("Error setting table_id: %s", err)
	}
	if err := d.Set("dataset_id", res.TableReference.DatasetId); err != nil {
		return fmt.Errorf("Error setting dataset_id: %s", err)
	}
	if err := d.Set("num_long_term_bytes", res.NumLongTermBytes); err != nil {
		return fmt.Errorf("Error setting num_long_term_bytes: %s", err)
	}
	if err := d.Set("num_rows", res.NumRows); err != nil {
		return fmt.Errorf("Error setting num_rows: %s", err)
	}
	if err := d.Set("self_link", res.SelfLink); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err := d.Set("type", res.Type); err != nil {
		return fmt.Errorf("Error setting type: %s", err)
	}

	if res.ExternalDataConfiguration != nil {
		externalDataConfiguration, err := flattenExternalDataConfiguration(res.ExternalDataConfiguration)
		if err != nil {
			return err
		}

		if v, ok := d.GetOk("external_data_configuration"); ok {
			// The API response doesn't return the `external_data_configuration.schema`
			// used when creating the table and it cannot be queried.
			// After creation, a computed schema is stored in the toplevel `schema`,
			// which combines `external_data_configuration.schema`
			// with any hive partioning fields found in the `source_uri_prefix`.
			// So just assume the configured schema has been applied after successful
			// creation, by copying the configured value back into the resource schema.
			// This avoids that reading back this field will be identified as a change.
			// The `ForceNew=true` on `external_data_configuration.schema` will ensure
			// the users' expectation that changing the configured  input schema will
			// recreate the resource.
			edc := v.([]interface{})[0].(map[string]interface{})
			if edc["schema"] != nil {
				externalDataConfiguration[0]["schema"] = edc["schema"]
			}
		}

		if err := d.Set("external_data_configuration", externalDataConfiguration); err != nil {
			return fmt.Errorf("Error setting external_data_configuration: %s", err)
		}
	}

	if res.TimePartitioning != nil {
		if err := d.Set("time_partitioning", flattenTimePartitioning(res.TimePartitioning)); err != nil {
			return err
		}
	}

	if res.RangePartitioning != nil {
		if err := d.Set("range_partitioning", flattenRangePartitioning(res.RangePartitioning)); err != nil {
			return err
		}
	}

	if res.Clustering != nil {
		if err := d.Set("clustering", res.Clustering.Fields); err != nil {
			return fmt.Errorf("Error setting clustering: %s", err)
		}
	}
	if res.EncryptionConfiguration != nil {
		if err := d.Set("encryption_configuration", flattenEncryptionConfiguration(res.EncryptionConfiguration)); err != nil {
			return err
		}
	}

	if res.Schema != nil {
		schema, err := flattenSchema(res.Schema)
		if err != nil {
			return err
		}
		if err := d.Set("schema", schema); err != nil {
			return fmt.Errorf("Error setting schema: %s", err)
		}
	}

	if res.View != nil {
		view := flattenView(res.View)
		if err := d.Set("view", view); err != nil {
			return fmt.Errorf("Error setting view: %s", err)
		}
	}

	if res.MaterializedView != nil {
		materialized_view := flattenMaterializedView(res.MaterializedView)

		if err := d.Set("materialized_view", materialized_view); err != nil {
			return fmt.Errorf("Error setting materialized view: %s", err)
		}
	}

	return nil
}

func resourceBigQueryTableUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	table, err := resourceTable(d, meta)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating BigQuery table: %s", d.Id())

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	datasetID := d.Get("dataset_id").(string)
	tableID := d.Get("table_id").(string)

	if _, err = config.NewBigQueryClient(userAgent).Tables.Update(project, datasetID, tableID, table).Do(); err != nil {
		return err
	}

	return resourceBigQueryTableRead(d, meta)
}

func resourceBigQueryTableDelete(d *schema.ResourceData, meta interface{}) error {
	if d.Get("deletion_protection").(bool) {
		return fmt.Errorf("cannot destroy instance without setting deletion_protection=false and running `terraform apply`")
	}
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting BigQuery table: %s", d.Id())

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	datasetID := d.Get("dataset_id").(string)
	tableID := d.Get("table_id").(string)

	if err := config.NewBigQueryClient(userAgent).Tables.Delete(project, datasetID, tableID).Do(); err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func expandExternalDataConfiguration(cfg interface{}) (*bigquery.ExternalDataConfiguration, error) {
	raw := cfg.([]interface{})[0].(map[string]interface{})

	edc := &bigquery.ExternalDataConfiguration{
		Autodetect: raw["autodetect"].(bool),
	}

	sourceUris := []string{}
	for _, rawSourceUri := range raw["source_uris"].([]interface{}) {
		sourceUris = append(sourceUris, rawSourceUri.(string))
	}
	if len(sourceUris) > 0 {
		edc.SourceUris = sourceUris
	}

	if v, ok := raw["compression"]; ok {
		edc.Compression = v.(string)
	}
	if v, ok := raw["csv_options"]; ok {
		edc.CsvOptions = expandCsvOptions(v)
	}
	if v, ok := raw["google_sheets_options"]; ok {
		edc.GoogleSheetsOptions = expandGoogleSheetsOptions(v)
	}
	if v, ok := raw["hive_partitioning_options"]; ok {
		edc.HivePartitioningOptions = expandHivePartitioningOptions(v)
	}
	if v, ok := raw["avro_options"]; ok {
		edc.AvroOptions = expandAvroOptions(v)
	}
	if v, ok := raw["ignore_unknown_values"]; ok {
		edc.IgnoreUnknownValues = v.(bool)
	}
	if v, ok := raw["max_bad_records"]; ok {
		edc.MaxBadRecords = int64(v.(int))
	}
	if v, ok := raw["schema"]; ok {
		schema, err := expandSchema(v)
		if err != nil {
			return nil, err
		}
		edc.Schema = schema
	}
	if v, ok := raw["source_format"]; ok {
		edc.SourceFormat = v.(string)
	}
	if v, ok := raw["connection_id"]; ok {
		edc.ConnectionId = v.(string)
	}

	return edc, nil

}

func flattenExternalDataConfiguration(edc *bigquery.ExternalDataConfiguration) ([]map[string]interface{}, error) {
	result := map[string]interface{}{}

	result["autodetect"] = edc.Autodetect
	result["source_uris"] = edc.SourceUris

	if edc.Compression != "" {
		result["compression"] = edc.Compression
	}

	if edc.CsvOptions != nil {
		result["csv_options"] = flattenCsvOptions(edc.CsvOptions)
	}

	if edc.GoogleSheetsOptions != nil {
		result["google_sheets_options"] = flattenGoogleSheetsOptions(edc.GoogleSheetsOptions)
	}

	if edc.HivePartitioningOptions != nil {
		result["hive_partitioning_options"] = flattenHivePartitioningOptions(edc.HivePartitioningOptions)
	}

	if edc.AvroOptions != nil {
		result["avro_options"] = flattenAvroOptions(edc.AvroOptions)
	}

	if edc.IgnoreUnknownValues {
		result["ignore_unknown_values"] = edc.IgnoreUnknownValues
	}
	if edc.MaxBadRecords != 0 {
		result["max_bad_records"] = edc.MaxBadRecords
	}

	if edc.SourceFormat != "" {
		result["source_format"] = edc.SourceFormat
	}

	if edc.ConnectionId != "" {
		result["connection_id"] = edc.ConnectionId
	}

	return []map[string]interface{}{result}, nil
}

func expandCsvOptions(configured interface{}) *bigquery.CsvOptions {
	if len(configured.([]interface{})) == 0 {
		return nil
	}

	raw := configured.([]interface{})[0].(map[string]interface{})
	opts := &bigquery.CsvOptions{}

	if v, ok := raw["allow_jagged_rows"]; ok {
		opts.AllowJaggedRows = v.(bool)
		opts.ForceSendFields = append(opts.ForceSendFields, "allow_jagged_rows")
	}

	if v, ok := raw["allow_quoted_newlines"]; ok {
		opts.AllowQuotedNewlines = v.(bool)
		opts.ForceSendFields = append(opts.ForceSendFields, "allow_quoted_newlines")
	}

	if v, ok := raw["encoding"]; ok {
		opts.Encoding = v.(string)
	}

	if v, ok := raw["field_delimiter"]; ok {
		opts.FieldDelimiter = v.(string)
	}

	if v, ok := raw["skip_leading_rows"]; ok {
		opts.SkipLeadingRows = int64(v.(int))
	}

	if v, ok := raw["quote"]; ok {
		quote := v.(string)
		opts.Quote = &quote
	}

	opts.ForceSendFields = []string{"Quote"}

	return opts
}

func flattenCsvOptions(opts *bigquery.CsvOptions) []map[string]interface{} {
	result := map[string]interface{}{}

	if opts.AllowJaggedRows {
		result["allow_jagged_rows"] = opts.AllowJaggedRows
	}

	if opts.AllowQuotedNewlines {
		result["allow_quoted_newlines"] = opts.AllowQuotedNewlines
	}

	if opts.Encoding != "" {
		result["encoding"] = opts.Encoding
	}

	if opts.FieldDelimiter != "" {
		result["field_delimiter"] = opts.FieldDelimiter
	}

	if opts.SkipLeadingRows != 0 {
		result["skip_leading_rows"] = opts.SkipLeadingRows
	}

	if opts.Quote != nil {
		result["quote"] = *opts.Quote
	}

	return []map[string]interface{}{result}
}

func expandGoogleSheetsOptions(configured interface{}) *bigquery.GoogleSheetsOptions {
	if len(configured.([]interface{})) == 0 {
		return nil
	}

	raw := configured.([]interface{})[0].(map[string]interface{})
	opts := &bigquery.GoogleSheetsOptions{}

	if v, ok := raw["range"]; ok {
		opts.Range = v.(string)
	}

	if v, ok := raw["skip_leading_rows"]; ok {
		opts.SkipLeadingRows = int64(v.(int))
	}
	return opts
}

func flattenGoogleSheetsOptions(opts *bigquery.GoogleSheetsOptions) []map[string]interface{} {
	result := map[string]interface{}{}

	if opts.Range != "" {
		result["range"] = opts.Range
	}

	if opts.SkipLeadingRows != 0 {
		result["skip_leading_rows"] = opts.SkipLeadingRows
	}

	return []map[string]interface{}{result}
}

func expandHivePartitioningOptions(configured interface{}) *bigquery.HivePartitioningOptions {
	if len(configured.([]interface{})) == 0 {
		return nil
	}

	raw := configured.([]interface{})[0].(map[string]interface{})
	opts := &bigquery.HivePartitioningOptions{}

	if v, ok := raw["mode"]; ok {
		opts.Mode = v.(string)
	}

	if v, ok := raw["require_partition_filter"]; ok {
		opts.RequirePartitionFilter = v.(bool)
	}

	if v, ok := raw["source_uri_prefix"]; ok {
		opts.SourceUriPrefix = v.(string)
	}

	return opts
}

func flattenHivePartitioningOptions(opts *bigquery.HivePartitioningOptions) []map[string]interface{} {
	result := map[string]interface{}{}

	if opts.Mode != "" {
		result["mode"] = opts.Mode
	}

	if opts.RequirePartitionFilter {
		result["require_partition_filter"] = opts.RequirePartitionFilter
	}

	if opts.SourceUriPrefix != "" {
		result["source_uri_prefix"] = opts.SourceUriPrefix
	}

	return []map[string]interface{}{result}
}

func expandAvroOptions(configured interface{}) *bigquery.AvroOptions {
	if len(configured.([]interface{})) == 0 {
		return nil
	}

	raw := configured.([]interface{})[0].(map[string]interface{})
	opts := &bigquery.AvroOptions{}

	if v, ok := raw["use_avro_logical_types"]; ok {
		opts.UseAvroLogicalTypes = v.(bool)
	}

	return opts
}

func flattenAvroOptions(opts *bigquery.AvroOptions) []map[string]interface{} {
	result := map[string]interface{}{}

	if opts.UseAvroLogicalTypes {
		result["use_avro_logical_types"] = opts.UseAvroLogicalTypes
	}

	return []map[string]interface{}{result}
}

func expandSchema(raw interface{}) (*bigquery.TableSchema, error) {
	var fields []*bigquery.TableFieldSchema

	if len(raw.(string)) == 0 {
		return nil, nil
	}

	if err := json.Unmarshal([]byte(raw.(string)), &fields); err != nil {
		return nil, err
	}

	return &bigquery.TableSchema{Fields: fields}, nil
}

func flattenSchema(tableSchema *bigquery.TableSchema) (string, error) {
	schema, err := json.Marshal(tableSchema.Fields)
	if err != nil {
		return "", err
	}

	return string(schema), nil
}

func expandTimePartitioning(configured interface{}) *bigquery.TimePartitioning {
	raw := configured.([]interface{})[0].(map[string]interface{})
	tp := &bigquery.TimePartitioning{Type: raw["type"].(string)}

	if v, ok := raw["field"]; ok {
		tp.Field = v.(string)
	}

	if v, ok := raw["expiration_ms"]; ok {
		tp.ExpirationMs = int64(v.(int))
	}

	if v, ok := raw["require_partition_filter"]; ok {
		tp.RequirePartitionFilter = v.(bool)
	}

	return tp
}

func expandRangePartitioning(configured interface{}) (*bigquery.RangePartitioning, error) {
	if configured == nil {
		return nil, nil
	}

	rpList := configured.([]interface{})
	if len(rpList) == 0 || rpList[0] == nil {
		return nil, errors.New("Error casting range partitioning interface to expected structure")
	}

	rangePartJson := rpList[0].(map[string]interface{})
	rp := &bigquery.RangePartitioning{
		Field: rangePartJson["field"].(string),
	}

	if v, ok := rangePartJson["range"]; ok && v != nil {
		rangeLs := v.([]interface{})
		if len(rangeLs) != 1 || rangeLs[0] == nil {
			return nil, errors.New("Non-empty range must be given for range partitioning")
		}

		rangeJson := rangeLs[0].(map[string]interface{})
		rp.Range = &bigquery.RangePartitioningRange{
			Start:           int64(rangeJson["start"].(int)),
			End:             int64(rangeJson["end"].(int)),
			Interval:        int64(rangeJson["interval"].(int)),
			ForceSendFields: []string{"Start"},
		}
	}

	return rp, nil
}

func flattenEncryptionConfiguration(ec *bigquery.EncryptionConfiguration) []map[string]interface{} {
	re := regexp.MustCompile(`(projects/.*/locations/.*/keyRings/.*/cryptoKeys/.*)/cryptoKeyVersions/.*`)
	paths := re.FindStringSubmatch(ec.KmsKeyName)

	if len(paths) > 0 {
		return []map[string]interface{}{
			{
				"kms_key_name":    paths[1],
				"kms_key_version": ec.KmsKeyName,
			},
		}
	}

	//	The key name was returned, no need to set the version
	return []map[string]interface{}{{"kms_key_name": ec.KmsKeyName, "kms_key_version": ""}}
}

func flattenTimePartitioning(tp *bigquery.TimePartitioning) []map[string]interface{} {
	result := map[string]interface{}{"type": tp.Type}

	if tp.Field != "" {
		result["field"] = tp.Field
	}

	if tp.ExpirationMs != 0 {
		result["expiration_ms"] = tp.ExpirationMs
	}

	if tp.RequirePartitionFilter {
		result["require_partition_filter"] = tp.RequirePartitionFilter
	}

	return []map[string]interface{}{result}
}

func flattenRangePartitioning(rp *bigquery.RangePartitioning) []map[string]interface{} {
	result := map[string]interface{}{
		"field": rp.Field,
		"range": []map[string]interface{}{
			{
				"start":    rp.Range.Start,
				"end":      rp.Range.End,
				"interval": rp.Range.Interval,
			},
		},
	}

	return []map[string]interface{}{result}
}

func expandView(configured interface{}) *bigquery.ViewDefinition {
	raw := configured.([]interface{})[0].(map[string]interface{})
	vd := &bigquery.ViewDefinition{Query: raw["query"].(string)}

	if v, ok := raw["use_legacy_sql"]; ok {
		vd.UseLegacySql = v.(bool)
		vd.ForceSendFields = append(vd.ForceSendFields, "UseLegacySql")
	}

	return vd
}

func flattenView(vd *bigquery.ViewDefinition) []map[string]interface{} {
	result := map[string]interface{}{"query": vd.Query}
	result["use_legacy_sql"] = vd.UseLegacySql

	return []map[string]interface{}{result}
}

func expandMaterializedView(configured interface{}) *bigquery.MaterializedViewDefinition {
	raw := configured.([]interface{})[0].(map[string]interface{})
	mvd := &bigquery.MaterializedViewDefinition{Query: raw["query"].(string)}

	if v, ok := raw["enable_refresh"]; ok {
		mvd.EnableRefresh = v.(bool)
		mvd.ForceSendFields = append(mvd.ForceSendFields, "EnableRefresh")
	}

	if v, ok := raw["refresh_interval_ms"]; ok {
		mvd.RefreshIntervalMs = int64(v.(int))
		mvd.ForceSendFields = append(mvd.ForceSendFields, "RefreshIntervalMs")
	}

	return mvd
}

func flattenMaterializedView(mvd *bigquery.MaterializedViewDefinition) []map[string]interface{} {
	result := map[string]interface{}{"query": mvd.Query}
	result["enable_refresh"] = mvd.EnableRefresh
	result["refresh_interval_ms"] = mvd.RefreshIntervalMs

	return []map[string]interface{}{result}
}

func resourceBigQueryTableImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/datasets/(?P<dataset_id>[^/]+)/tables/(?P<table_id>[^/]+)",
		"(?P<project>[^/]+)/(?P<dataset_id>[^/]+)/(?P<table_id>[^/]+)",
		"(?P<dataset_id>[^/]+)/(?P<table_id>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Explicitly set virtual fields to default values on import
	if err := d.Set("deletion_protection", true); err != nil {
		return nil, fmt.Errorf("Error setting deletion_protection: %s", err)
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/datasets/{{dataset_id}}/tables/{{table_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
