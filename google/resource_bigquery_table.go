//
package google

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"

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
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			// TableId: [Required] The ID of the table. The ID must contain only
			// letters (a-z, A-Z), numbers (0-9), or underscores (_). The maximum
			// length is 1,024 characters.
			"table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// DatasetId: [Required] The ID of the dataset containing this table.
			"dataset_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// ProjectId: [Required] The ID of the project containing this table.
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			// Description: [Optional] A user-friendly description of this table.
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// ExpirationTime: [Optional] The time when this table expires, in
			// milliseconds since the epoch. If not present, the table will persist
			// indefinitely. Expired tables will be deleted and their storage
			// reclaimed.
			"expiration_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			// ExternalDataConfiguration [Optional] Describes the data format,
			// location, and other properties of a table stored outside of BigQuery.
			// By defining these properties, the data source can then be queried as
			// if it were a standard BigQuery table.
			"external_data_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Autodetect : [Required] If true, let BigQuery try to autodetect the
						// schema and format of the table.
						"autodetect": {
							Type:     schema.TypeBool,
							Required: true,
						},
						// SourceFormat [Required] The data format.
						"source_format": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"CSV", "GOOGLE_SHEETS", "NEWLINE_DELIMITED_JSON", "AVRO", "DATSTORE_BACKUP",
							}, false),
						},
						// SourceURIs [Required] The fully-qualified URIs that point to your data in Google Cloud.
						"source_uris": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						// Compression: [Optional] The compression type of the data source.
						"compression": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"NONE", "GZIP"}, false),
							Default:      "NONE",
						},
						// CsvOptions: [Optional] Additional properties to set if
						// sourceFormat is set to CSV.
						"csv_options": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// Quote: [Required] The value that is used to quote data
									// sections in a CSV file.
									"quote": {
										Type:     schema.TypeString,
										Required: true,
									},
									// AllowJaggedRows: [Optional] Indicates if BigQuery should
									// accept rows that are missing trailing optional columns.
									"allow_jagged_rows": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									// AllowQuotedNewlines: [Optional] Indicates if BigQuery
									// should allow quoted data sections that contain newline
									// characters in a CSV file. The default value is false.
									"allow_quoted_newlines": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									// Encoding: [Optional] The character encoding of the data.
									// The supported values are UTF-8 or ISO-8859-1.
									"encoding": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"ISO-8859-1", "UTF-8"}, false),
										Default:      "UTF-8",
									},
									// FieldDelimiter: [Optional] The separator for fields in a CSV file.
									"field_delimiter": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  ",",
									},
									// SkipLeadingRows: [Optional] The number of rows at the top
									// of a CSV file that BigQuery will skip when reading the data.
									"skip_leading_rows": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
								},
							},
						},
						// GoogleSheetsOptions: [Optional] Additional options if sourceFormat is set to GOOGLE_SHEETS.
						"google_sheets_options": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// Range: [Optional] Range of a sheet to query from. Only used when non-empty.
									// Typical format: !:
									"range": {
										Removed:  "This field is in beta. Use it in the the google-beta provider instead. See https://terraform.io/docs/providers/google/provider_versions.html for more details.",
										Type:     schema.TypeString,
										Optional: true,
									},
									// SkipLeadingRows: [Optional] The number of rows at the top
									// of the sheet that BigQuery will skip when reading the data.
									"skip_leading_rows": {
										Type:     schema.TypeInt,
										Optional: true,
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
							Type:     schema.TypeBool,
							Optional: true,
						},
						// MaxBadRecords: [Optional] The maximum number of bad records that
						// BigQuery can ignore when reading data.
						"max_bad_records": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},

			// FriendlyName: [Optional] A descriptive name for this table.
			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Labels: [Experimental] The labels associated with this table. You can
			// use these to organize and group your tables. Label keys and values
			// can be no longer than 63 characters, can only contain lowercase
			// letters, numeric characters, underscores and dashes. International
			// characters are allowed. Label values are optional. Label keys must
			// start with a letter and each label in the list must have a different
			// key.
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// Schema: [Optional] Describes the schema of this table.
			// Schema is required for external tables in CSV and JSON formats
			// and disallowed for Google Cloud Bigtable, Cloud Datastore backups,
			// and Avro formats.
			"schema": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.ValidateJsonString,
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
			},

			// View: [Optional] If specified, configures this table as a view.
			"view": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Query: [Required] A query that BigQuery executes when the view is
						// referenced.
						"query": {
							Type:     schema.TypeString,
							Required: true,
						},

						// UseLegacySQL: [Optional] Specifies whether to use BigQuery's
						// legacy SQL for this view. The default value is true. If set to
						// false, the view will use BigQuery's standard SQL:
						"use_legacy_sql": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},

			// TimePartitioning: [Experimental] If specified, configures time-based
			// partitioning for this table.
			"time_partitioning": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// ExpirationMs: [Optional] Number of milliseconds for which to keep the
						// storage for a partition.
						"expiration_ms": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						// Type: [Required] The only type supported is DAY, which will generate
						// one partition per day based on data loading time.
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"DAY"}, false),
						},

						// Field: [Optional] The field used to determine how to create a time-based
						// partition. If time-based partitioning is enabled without this value, the
						// table is partitioned based on the load time.
						"field": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						// RequirePartitionFilter: [Optional] If set to true, queries over this table
						// require a partition filter that can be used for partition elimination to be
						// specified.
						"require_partition_filter": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},

			// Clustering: [Optional] Specifies column names to use for data clustering.  Up to four
			// top-level columns are allowed, and should be specified in descending priority order.
			"clustering": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 4,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// CreationTime: [Output-only] The time when this table was created, in
			// milliseconds since the epoch.
			"creation_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// Etag: [Output-only] A hash of this resource.
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// LastModifiedTime: [Output-only] The time when this table was last
			// modified, in milliseconds since the epoch.
			"last_modified_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// Location: [Output-only] The geographic location where the table
			// resides. This value is inherited from the dataset.
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// NumBytes: [Output-only] The size of this table in bytes, excluding
			// any data in the streaming buffer.
			"num_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// NumLongTermBytes: [Output-only] The number of bytes in the table that
			// are considered "long-term storage".
			"num_long_term_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// NumRows: [Output-only] The number of rows of data in this table,
			// excluding any data in the streaming buffer.
			"num_rows": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// SelfLink: [Output-only] A URL that can be used to access this
			// resource again.
			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// Type: [Output-only] Describes the table type. The following values
			// are supported: TABLE: A normal BigQuery table. VIEW: A virtual table
			// defined by a SQL query. EXTERNAL: A table that references data stored
			// in an external storage system, such as Google Cloud Storage. The
			// default value is TABLE.
			"type": {
				Type:     schema.TypeString,
				Computed: true,
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

	d.SetId(fmt.Sprintf("%s:%s.%s", res.TableReference.ProjectId, res.TableReference.DatasetId, res.TableReference.TableId))

	return resourceBigQueryTableRead(d, meta)
}

func resourceBigQueryTableRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	log.Printf("[INFO] Reading BigQuery table: %s", d.Id())

	id, err := parseBigQueryTableId(d.Id())
	if err != nil {
		return err
	}

	res, err := config.clientBigQuery.Tables.Get(id.Project, id.DatasetId, id.TableId).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("BigQuery table %q", id.TableId))
	}

	d.Set("project", id.Project)
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

		d.Set("external_data_configuration", externalDataConfiguration)
	}

	if res.TimePartitioning != nil {
		if err := d.Set("time_partitioning", flattenTimePartitioning(res.TimePartitioning)); err != nil {
			return err
		}
	}

	if res.Clustering != nil {
		d.Set("clustering", res.Clustering.Fields)
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

	id, err := parseBigQueryTableId(d.Id())
	if err != nil {
		return err
	}

	if _, err = config.clientBigQuery.Tables.Update(id.Project, id.DatasetId, id.TableId, table).Do(); err != nil {
		return err
	}

	return resourceBigQueryTableRead(d, meta)
}

func resourceBigQueryTableDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	log.Printf("[INFO] Deleting BigQuery table: %s", d.Id())

	id, err := parseBigQueryTableId(d.Id())
	if err != nil {
		return err
	}

	if err := config.clientBigQuery.Tables.Delete(id.Project, id.DatasetId, id.TableId).Do(); err != nil {
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
	if v, ok := raw["ignore_unknown_values"]; ok {
		edc.IgnoreUnknownValues = v.(bool)
	}
	if v, ok := raw["max_bad_records"]; ok {
		edc.MaxBadRecords = int64(v.(int))
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

	if edc.IgnoreUnknownValues == true {
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

	if opts.AllowJaggedRows == true {
		result["allow_jagged_rows"] = opts.AllowJaggedRows
	}

	if opts.AllowQuotedNewlines == true {
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

func flattenTimePartitioning(tp *bigquery.TimePartitioning) []map[string]interface{} {
	result := map[string]interface{}{"type": tp.Type}

	if tp.Field != "" {
		result["field"] = tp.Field
	}

	if tp.ExpirationMs != 0 {
		result["expiration_ms"] = tp.ExpirationMs
	}

	if tp.RequirePartitionFilter == true {
		result["require_partition_filter"] = tp.RequirePartitionFilter
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

type bigQueryTableId struct {
	Project, DatasetId, TableId string
}

func parseBigQueryTableId(id string) (*bigQueryTableId, error) {
	// Expected format is "PROJECT:DATASET.TABLE", but the project can itself have . and : in it.
	// Those characters are not valid dataset or table components, so just split on the last two.
	matchRegex := regexp.MustCompile("^(.+):([^:.]+)\\.([^:.]+)$")
	subMatches := matchRegex.FindStringSubmatch(id)
	if subMatches == nil {
		return nil, fmt.Errorf("Invalid BigQuery table specifier. Expecting {project}:{dataset-id}.{table-id}, got %s", id)
	}
	return &bigQueryTableId{
		Project:   subMatches[1],
		DatasetId: subMatches[2],
		TableId:   subMatches[3],
	}, nil
}
