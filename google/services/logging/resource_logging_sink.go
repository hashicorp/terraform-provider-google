// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
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
			DiffSuppressFunc: tpgresource.OptionalSurroundingSpacesSuppress,
			Description:      `The filter to apply when exporting logs. Only log entries that match the filter are exported.`,
		},

		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: `A description of this sink. The maximum length of the description is 8000 characters.`,
		},

		"disabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: `If set to True, then this sink is disabled and it does not export any log entries.`,
		},

		"exclusions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: `Log entries that match any of the exclusion filters will not be exported. If a log entry is matched by both filter and one of exclusion's filters, it will not be exported.`,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: `A client-assigned identifier, such as "load-balancer-exclusion". Identifiers are limited to 100 characters and can include only letters, digits, underscores, hyphens, and periods. First character has to be alphanumeric.`,
					},
					"description": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: `A description of this exclusion.`,
					},
					"filter": {
						Type:        schema.TypeString,
						Required:    true,
						Description: `An advanced logs filter that matches the log entries to be excluded. By using the sample function, you can exclude less than 100% of the matching log entries`,
					},
					"disabled": {
						Type:        schema.TypeBool,
						Optional:    true,
						Default:     false,
						Description: `If set to True, then this exclusion is disabled and it does not exclude any log entries`,
					},
				},
			},
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
		Description:     d.Get("description").(string),
		Disabled:        d.Get("disabled").(bool),
		Exclusions:      expandLoggingSinkExclusions(d.Get("exclusions")),
		BigqueryOptions: expandLoggingSinkBigqueryOptions(d.Get("bigquery_options")),
	}
	return id, &sink
}

func flattenResourceLoggingSink(d *schema.ResourceData, sink *logging.LogSink) error {
	if err := d.Set("name", sink.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("destination", sink.Destination); err != nil {
		return fmt.Errorf("Error setting destination: %s", err)
	}
	if err := d.Set("filter", sink.Filter); err != nil {
		return fmt.Errorf("Error setting filter: %s", err)
	}
	if err := d.Set("description", sink.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("disabled", sink.Disabled); err != nil {
		return fmt.Errorf("Error setting disabled: %s", err)
	}
	if err := d.Set("writer_identity", sink.WriterIdentity); err != nil {
		return fmt.Errorf("Error setting writer_identity: %s", err)
	}
	if err := d.Set("exclusions", flattenLoggingSinkExclusion(sink.Exclusions)); err != nil {
		return fmt.Errorf("Error setting exclusions: %s", err)
	}
	if err := d.Set("bigquery_options", flattenLoggingSinkBigqueryOptions(sink.BigqueryOptions)); err != nil {
		return fmt.Errorf("Error setting bigquery_options: %s", err)
	}

	return nil
}

func expandResourceLoggingSinkForUpdate(d *schema.ResourceData) (sink *logging.LogSink, updateMask string) {
	// Can only update destination/filter right now. Despite the method below using 'Patch', the API requires both
	// destination and filter (even if unchanged).
	sink = &logging.LogSink{
		Destination:     d.Get("destination").(string),
		Filter:          d.Get("filter").(string),
		Disabled:        d.Get("disabled").(bool),
		Description:     d.Get("description").(string),
		ForceSendFields: []string{"Destination", "Filter", "Disabled"},
	}

	updateFields := []string{}
	if d.HasChange("destination") {
		updateFields = append(updateFields, "destination")
	}
	if d.HasChange("filter") {
		updateFields = append(updateFields, "filter")
	}
	if d.HasChange("description") {
		updateFields = append(updateFields, "description")
	}
	if d.HasChange("disabled") {
		updateFields = append(updateFields, "disabled")
	}
	if d.HasChange("exclusions") {
		sink.Exclusions = expandLoggingSinkExclusions(d.Get("exclusions"))
		updateFields = append(updateFields, "exclusions")
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

func expandLoggingSinkExclusions(v interface{}) []*logging.LogExclusion {
	if v == nil {
		return nil
	}
	exclusions := v.([]interface{})
	if len(exclusions) == 0 {
		return nil
	}
	results := make([]*logging.LogExclusion, 0, len(exclusions))
	for _, e := range exclusions {
		exclusion := e.(map[string]interface{})
		results = append(results, &logging.LogExclusion{
			Name:        exclusion["name"].(string),
			Description: exclusion["description"].(string),
			Filter:      exclusion["filter"].(string),
			Disabled:    exclusion["disabled"].(bool),
		})
	}
	return results
}

func flattenLoggingSinkExclusion(exclusions []*logging.LogExclusion) []map[string]interface{} {
	if exclusions == nil {
		return nil
	}
	flattenedExclusions := make([]map[string]interface{}, 0, len(exclusions))
	for _, e := range exclusions {
		flattenedExclusion := map[string]interface{}{
			"name":        e.Name,
			"description": e.Description,
			"filter":      e.Filter,
			"disabled":    e.Disabled,
		}
		flattenedExclusions = append(flattenedExclusions, flattenedExclusion)

	}

	return flattenedExclusions
}

func resourceLoggingSinkImportState(sinkType string) schema.StateFunc {
	return func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
		loggingSinkId, err := ParseLoggingSinkId(d.Id())
		if err != nil {
			return nil, err
		}

		if err := d.Set(sinkType, loggingSinkId.resourceId); err != nil {
			return nil, fmt.Errorf("Error setting sinkType: %s", err)
		}

		return []*schema.ResourceData{d}, nil
	}
}
