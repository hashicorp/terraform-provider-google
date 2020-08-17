package google

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/logging/v2"
)

func resourceLoggingSinkSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: `The name of the logging sink.`,
		},

		"destination": {
			Type:        schema.TypeString,
			Required:    true,
			Description: `The destination of the sink (or, in other words, where logs are written to). Can be a Cloud Storage bucket, a PubSub topic, or a BigQuery dataset. Examples: "storage.googleapis.com/[GCS_BUCKET]" "bigquery.googleapis.com/projects/[PROJECT_ID]/datasets/[DATASET]" "pubsub.googleapis.com/projects/[PROJECT_ID]/topics/[TOPIC_ID]" The writer associated with the sink must have access to write to the above resource.`,
		},

		"filter": {
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: optionalSurroundingSpacesSuppress,
			Description:      `The filter to apply when exporting logs. Only log entries that match the filter are exported.`,
		},

		"writer_identity": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: `The identity associated with this sink. This identity must be granted write access to the configured destination.`,
		},

		"bigquery_options": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    1,
			Description: `Options that affect sinks exporting data to BigQuery.`,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"use_partitioned_tables": {
						Type:        schema.TypeBool,
						Required:    true,
						Description: `Whether to use BigQuery's partition tables. By default, Logging creates dated tables based on the log entries' timestamps, e.g. syslog_20170523. With partitioned tables the date suffix is no longer present and special query syntax has to be used instead. In both cases, tables are sharded based on UTC timezone.`,
					},
				},
			},
		},
	}
}

func expandResourceLoggingSink(d *schema.ResourceData, resourceType, resourceId string) (LoggingSinkId, *logging.LogSink) {
	id := LoggingSinkId{
		resourceType: resourceType,
		resourceId:   resourceId,
		name:         d.Get("name").(string),
	}

	sink := logging.LogSink{
		Name:            d.Get("name").(string),
		Destination:     d.Get("destination").(string),
		Filter:          d.Get("filter").(string),
		BigqueryOptions: expandLoggingSinkBigqueryOptions(d.Get("bigquery_options")),
	}
	return id, &sink
}

func flattenResourceLoggingSink(d *schema.ResourceData, sink *logging.LogSink) {
	d.Set("name", sink.Name)
	d.Set("destination", sink.Destination)
	d.Set("filter", sink.Filter)
	d.Set("writer_identity", sink.WriterIdentity)
	d.Set("bigquery_options", flattenLoggingSinkBigqueryOptions(sink.BigqueryOptions))
}

func expandResourceLoggingSinkForUpdate(d *schema.ResourceData) (sink *logging.LogSink, updateMask string) {
	// Can only update destination/filter right now. Despite the method below using 'Patch', the API requires both
	// destination and filter (even if unchanged).
	sink = &logging.LogSink{
		Destination:     d.Get("destination").(string),
		Filter:          d.Get("filter").(string),
		ForceSendFields: []string{"Destination", "Filter"},
	}

	updateFields := []string{}
	if d.HasChange("destination") {
		updateFields = append(updateFields, "destination")
	}
	if d.HasChange("filter") {
		updateFields = append(updateFields, "filter")
	}
	if d.HasChange("bigquery_options") {
		sink.BigqueryOptions = expandLoggingSinkBigqueryOptions(d.Get("bigquery_options"))
		updateFields = append(updateFields, "bigqueryOptions")
	}
	updateMask = strings.Join(updateFields, ",")
	return
}

func expandLoggingSinkBigqueryOptions(v interface{}) *logging.BigQueryOptions {
	if v == nil {
		return nil
	}
	optionsSlice := v.([]interface{})
	if len(optionsSlice) == 0 || optionsSlice[0] == nil {
		return nil
	}
	options := optionsSlice[0].(map[string]interface{})
	bo := &logging.BigQueryOptions{}
	if usePartitionedTables, ok := options["use_partitioned_tables"]; ok {
		bo.UsePartitionedTables = usePartitionedTables.(bool)
	}
	return bo
}

func flattenLoggingSinkBigqueryOptions(o *logging.BigQueryOptions) []map[string]interface{} {
	if o == nil {
		return nil
	}
	oMap := map[string]interface{}{
		"use_partitioned_tables": o.UsePartitionedTables,
	}
	return []map[string]interface{}{oMap}
}

func resourceLoggingSinkImportState(sinkType string) schema.StateFunc {
	return func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
		loggingSinkId, err := parseLoggingSinkId(d.Id())
		if err != nil {
			return nil, err
		}

		d.Set(sinkType, loggingSinkId.resourceId)

		return []*schema.ResourceData{d}, nil
	}
}
