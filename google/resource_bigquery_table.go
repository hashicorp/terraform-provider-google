package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"google.golang.org/api/bigquery/v2"
)

func resourceBigQueryTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigQueryTableCreate,
		Read:   resourceBigQueryTableRead,
		Delete: resourceBigQueryTableDelete,
		Update: resourceBigQueryTableUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceBigQueryTableImport,
		},
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
							Description: `The data format. Supported values are: "CSV", "GOOGLE_SHEETS", "NEWLINE_DELIMITED_JSON", "AVRO", "PARQUET", and "DATSTORE_BACKUP". To use "GOOGLE_SHEETS" the scopes must include "googleapis.com/auth/drive.readonly".`,
							ValidateFunc: validation.StringInSlice([]string{
								"CSV", "GOOGLE_SHEETS", "NEWLINE_DELIMITED_JSON", "AVRO", "DATSTORE_BACKUP", "PARQUET",
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
							ValidateFunc: validation.ValidateJsonString,
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
				ValidateFunc: validation.ValidateJsonString,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
				Description: `A JSON schema for the table.`,
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
							Description: `Number of milliseconds for which to keep the storage for a partition.`,
						},

						// Type: [Required] The supported types are DAY and HOUR, which will generate
						// one partition per day or hour based on data loading time.
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  `The supported types are DAY and HOUR, which will generate one partition per day or hour based on data loading time.`,
							ValidateFunc: validation.StringInSlice([]string{"DAY", "HOUR"}, false),
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
				ForceNew:    true,
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
		},
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

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	table, err := resourceTable(d, meta)
	if err != nil {
		return err
	}

	datasetID := d.Get("dataset_id").(string)

	log.Printf("[INFO] Creating BigQuery table: %s", table.TableReference.TableId)

	res, err := config.clientBigQuery.Tables.Insert(project, datasetID, table).Do()
	if err != nil {
		return err
	}

	log.Printf("[INFO] BigQuery table %s has been created", res.Id)
	d.SetId(fmt.Sprintf("projects/%s/datasets/%s/tables/%s", res.TableReference.ProjectId, res.TableReference.DatasetId, res.TableReference.TableId))

	return resourceBigQueryTableRead(d, meta)
}

func resourceBigQueryTableRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	log.Printf("[INFO] Reading BigQuery table: %s", d.Id())

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	datasetID := d.Get("dataset_id").(string)
	tableID := d.Get("table_id").(string)

	res, err := config.clientBigQuery.Tables.Get(project, datasetID, tableID).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("BigQuery table %q", tableID))
	}

	d.Set("project", project)
	d.Set("description", res.Description)
	d.Set("expiration_time", res.ExpirationTime)
	d.Set("friendly_name", res.FriendlyName)
	d.Set("labels", res.Labels)
	d.Set("creation_time", res.CreationTime)
	d.Set("etag", res.Etag)
	d.Set("last_modified_time", res.LastModifiedTime)
	d.Set("location", res.Location)
	d.Set("num_bytes", res.NumBytes)
	d.Set("table_id", res.TableReference.TableId)
	d.Set("dataset_id", res.TableReference.DatasetId)
	d.Set("num_long_term_bytes", res.NumLongTermBytes)
	d.Set("num_rows", res.NumRows)
	d.Set("self_link", res.SelfLink)
	d.Set("type", res.Type)

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

		d.Set("external_data_configuration", externalDataConfiguration)
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
		d.Set("clustering", res.Clustering.Fields)
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

		d.Set("schema", schema)
	}

	if res.View != nil {
		view := flattenView(res.View)
		d.Set("view", view)
	}

	return nil
}

func resourceBigQueryTableUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

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

	if _, err = config.clientBigQuery.Tables.Update(project, datasetID, tableID, table).Do(); err != nil {
		return err
	}

	return resourceBigQueryTableRead(d, meta)
}

func resourceBigQueryTableDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	log.Printf("[INFO] Deleting BigQuery table: %s", d.Id())

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	datasetID := d.Get("dataset_id").(string)
	tableID := d.Get("table_id").(string)

	if err := config.clientBigQuery.Tables.Delete(project, datasetID, tableID).Do(); err != nil {
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

	if edc.IgnoreUnknownValues {
		result["ignore_unknown_values"] = edc.IgnoreUnknownValues
	}
	if edc.MaxBadRecords != 0 {
		result["max_bad_records"] = edc.MaxBadRecords
	}

	if edc.SourceFormat != "" {
		result["source_format"] = edc.SourceFormat
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
	}

	if v, ok := raw["allow_quoted_newlines"]; ok {
		opts.AllowQuotedNewlines = v.(bool)
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

	if opts.SourceUriPrefix != "" {
		result["source_uri_prefix"] = opts.SourceUriPrefix
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
	return []map[string]interface{}{{"kms_key_name": ec.KmsKeyName}}
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

func resourceBigQueryTableImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/datasets/(?P<dataset_id>[^/]+)/tables/(?P<table_id>[^/]+)",
		"(?P<project>[^/]+)/(?P<dataset_id>[^/]+)/(?P<table_id>[^/]+)",
		"(?P<dataset_id>[^/]+)/(?P<table_id>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/datasets/{{dataset_id}}/tables/{{table_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
