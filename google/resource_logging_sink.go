package google

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strconv"
)

var loggingSinkSchema = map[string]*schema.Schema{
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
		Description: `Log entries that match any of the exclusion filters will not be exported. If a log entry is matched by both filter and one of exclusion_filters it will not be exported.`,
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

func resourceLoggingSinkSchema(parentSpecificSchema map[string]*schema.Schema) map[string]*schema.Schema {
	return mergeSchemas(loggingSinkSchema, parentSpecificSchema)
}

type loggingSinkIDFunc func(d *schema.ResourceData, config *Config) (string, error)
type loggingSinkCreateFunc func(d *schema.ResourceData, meta interface{}) error
type loggingSinkReadFunc func(d *schema.ResourceData, meta interface{}) error
type loggingSinkUpdateFunc func(d *schema.ResourceData, meta interface{}) error

func resourceLoggingSinkAcquireOrCreate(iDFunc loggingSinkIDFunc, sinkCreateFunc loggingSinkCreateFunc, sinkUpdateFunc loggingSinkUpdateFunc) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		userAgent, err := generateUserAgentString(d, config.userAgent)
		if err != nil {
			return err
		}

		id, err := iDFunc(d, config)
		if err != nil {
			return err
		}

		d.SetId(id)
		log.Printf("[DEBUG] Fetching Logging Sink: %#v", d.Id())
		url, err := replaceVars(d, config, fmt.Sprintf("{{LoggingBasePath}}%s", d.Id()))
		if err != nil {
			return err
		}

		res, _ := sendRequest(config, "GET", "", url, userAgent, nil)
		if res == nil {
			log.Printf("[DEGUG] Logging Sink does not exist %s", d.Id())
			return sinkCreateFunc(d, meta)
		}

		return sinkUpdateFunc(d, meta)
	}
}

type loggingSinksPathFunc func(d *schema.ResourceData, config *Config) (string, error)
type expandForCreateFunc func(d *schema.ResourceData, config *Config) (map[string]interface{}, bool)

func resourceLoggingSinkCreate(loggingSinksPath loggingSinksPathFunc, expandForCreate expandForCreateFunc, sinkReadFunc loggingSinkReadFunc) func(d *schema.ResourceData, meta interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		userAgent, err := generateUserAgentString(d, config.userAgent)
		if err != nil {
			return err
		}

		sinksPath, err := loggingSinksPath(d, config)
		if err != nil {
			return err
		}
		url, err := replaceVars(d, config, fmt.Sprintf("{{LoggingBasePath}}%s", sinksPath))
		if err != nil {
			return err
		}

		obj, uniqueWriterIdentity := expandForCreate(d, config)
		url, err = addQueryParams(url, map[string]string{
			"uniqueWriterIdentity": strconv.FormatBool(uniqueWriterIdentity),
		})
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Creating new Logging Sink: %#v", obj)
		billingProject := ""

		project, err := getProject(d, config)
		if err != nil {
			return err
		}
		billingProject = project

		// err == nil indicates that the billing_project value was found
		if bp, err := getBillingProject(d, config); err == nil {
			billingProject = bp
		}

		res, err := sendRequestWithTimeout(config, "POST", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return fmt.Errorf("Error creating Logging Sink: %s", err)
		}

		log.Printf("[DEBUG] Finished creating Logging Sink %q: %#v", d.Id(), res)

		return sinkReadFunc(d, meta)
	}
}

func expandLoggingSink(d *schema.ResourceData) (obj map[string]interface{}) {
	obj = make(map[string]interface{})
	obj["name"] = d.Get("name")
	obj["description"] = d.Get("description")
	obj["disabled"] = d.Get("disabled")
	obj["filter"] = d.Get("filter")
	obj["destination"] = d.Get("destination")
	obj["exclusions"] = expandLoggingSinkExclusions(d.Get("exclusions"))
	obj["bigqueryOptions"] = expandLoggingSinkBigqueryOptions(d.Get("bigquery_options"))
	return
}

func expandLoggingSinkExclusions(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	exclusions := v.([]interface{})
	if len(exclusions) == 0 {
		return nil
	}
	results := make([]map[string]interface{}, 0, len(exclusions))
	for _, e := range exclusions {
		exclusion := e.(map[string]interface{})
		results = append(results, map[string]interface{}{
			"name":        exclusion["name"].(string),
			"description": exclusion["description"].(string),
			"filter":      exclusion["filter"].(string),
			"disabled":    exclusion["disabled"].(bool),
		})
	}
	return results
}

func expandLoggingSinkBigqueryOptions(v interface{}) map[string]interface{} {
	if v == nil {
		return nil
	}
	bigQueryOptions := v.([]interface{})
	if len(bigQueryOptions) == 0 || bigQueryOptions[0] == nil {
		return nil
	}
	options := bigQueryOptions[0].(map[string]interface{})
	bo := make(map[string]interface{})
	if usePartitionedTables, ok := options["use_partitioned_tables"]; ok {
		bo["usePartitionedTables"] = usePartitionedTables.(bool)
	}
	return bo
}

type flattenLoggingSinkFunc func(d *schema.ResourceData, res map[string]interface{}, config *Config) error

func resourceLoggingSinkRead(flattenLoginSink flattenLoggingSinkFunc) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		userAgent, err := generateUserAgentString(d, config.userAgent)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Fetching Logging Sink: %#v", d.Id())

		url, err := replaceVars(d, config, fmt.Sprintf("{{LoggingBasePath}}%s", d.Id()))
		if err != nil {
			return err
		}

		res, err := sendRequest(config, "GET", "", url, userAgent, nil)
		if err != nil {
			log.Printf("[WARN] Unable to acquire Logging Sink at %s", d.Id())

			d.SetId("")
			return err
		}

		return flattenLoginSink(d, res, config)
	}
}

func flattenLoggingSinkBase(d *schema.ResourceData, res map[string]interface{}) error {
	if err := d.Set("name", res["name"]); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("description", res["description"]); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("disabled", res["disabled"]); err != nil {
		return fmt.Errorf("Error setting disabled: %s", err)
	}
	if err := d.Set("filter", res["filter"]); err != nil {
		return fmt.Errorf("Error setting filter: %s", err)
	}
	if err := d.Set("destination", res["destination"]); err != nil {
		return fmt.Errorf("Error setting destination: %s", err)
	}
	if err := d.Set("writer_identity", res["writerIdentity"]); err != nil {
		return fmt.Errorf("Error setting writer_identity: %s", err)
	}
	if err := d.Set("exclusions", flattenLoggingSinkExclusions(res["exclusions"])); err != nil {
		return fmt.Errorf("Error setting exclusions: %s", err)
	}
	if err := d.Set("bigquery_options", flattenLoggingSinkBigqueryOptions(res["bigqueryOptions"])); err != nil {
		return fmt.Errorf("Error setting bigquery_options: %s", err)
	}
	return nil
}

func flattenLoggingSinkExclusions(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	exclusions := v.([]interface{})
	if len(exclusions) == 0 {
		return nil
	}
	flattenedExclusions := make([]map[string]interface{}, 0, len(exclusions))
	for _, e := range exclusions {
		exclusion := e.(map[string]interface{})
		flattenedExclusion := map[string]interface{}{
			"name":        exclusion["name"],
			"description": exclusion["description"],
			"filter":      exclusion["filter"],
			"disabled":    exclusion["disabled"],
		}
		flattenedExclusions = append(flattenedExclusions, flattenedExclusion)
	}
	return flattenedExclusions
}

func flattenLoggingSinkBigqueryOptions(v interface{}) []map[string]interface{} {
	if v == nil {
		return nil
	}
	bigqueryOptions := v.(map[string]interface{})
	oMap := map[string]interface{}{
		"use_partitioned_tables": bigqueryOptions["usePartitionedTables"],
	}
	return []map[string]interface{}{oMap}
}

type expandLoggingSinkForUpdateFunc func(d *schema.ResourceData, config *Config) (map[string]interface{}, string, bool)

func resourceLoggingSinkUpdate(expandLoggingSinkForUpdate expandLoggingSinkForUpdateFunc, sinkReadFunc loggingSinkReadFunc) func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		config := meta.(*Config)
		userAgent, err := generateUserAgentString(d, config.userAgent)
		if err != nil {
			return err
		}

		url, err := replaceVars(d, config, fmt.Sprintf("{{LoggingBasePath}}%s", d.Id()))
		if err != nil {
			return err
		}

		obj, updateMask, uniqueWriterIdentity := expandLoggingSinkForUpdate(d, config)
		url, err = addQueryParams(url, map[string]string{
			"updateMask":           updateMask,
			"uniqueWriterIdentity": strconv.FormatBool(uniqueWriterIdentity),
		})
		if err != nil {
			return err
		}
		_, err = sendRequestWithTimeout(config, "PATCH", "", url, userAgent, obj, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error updating Logging Sink %q: %s", d.Id(), err)
		}

		return sinkReadFunc(d, meta)
	}
}

func expandResourceLoggingSinkForUpdateBase(d *schema.ResourceData) (obj map[string]interface{}, updateFields []string) {
	obj = make(map[string]interface{})
	updateFields = []string{}
	if d.HasChange("destination") {
		obj["destination"] = d.Get("destination")
		updateFields = append(updateFields, "destination")
	}
	if d.HasChange("filter") {
		obj["filter"] = d.Get("filter")
		updateFields = append(updateFields, "filter")
	}
	if d.HasChange("description") {
		obj["description"] = d.Get("description").(string)
		updateFields = append(updateFields, "description")
	}
	if d.HasChange("disabled") {
		obj["disabled"] = d.Get("disabled").(bool)
		updateFields = append(updateFields, "disabled")
	}
	if d.HasChange("exclusions") {
		obj["exclusions"] = expandLoggingSinkExclusions(d.Get("exclusions"))
		updateFields = append(updateFields, "exclusions")
	}
	if d.HasChange("bigquery_options") {
		obj["bigqueryOptions"] = expandLoggingSinkBigqueryOptions(d.Get("bigquery_options"))
		updateFields = append(updateFields, "bigqueryOptions")
	}
	return
}

func isAutomaticallyCreatedSink(d *schema.ResourceData) bool {
	name := d.Get("name").(string)
	return name == "_Default" || name == "_Required"
}

func resourceLoggingSinkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if !isAutomaticallyCreatedSink(d) {
		userAgent, err := generateUserAgentString(d, config.userAgent)
		if err != nil {
			return err
		}

		url, err := replaceVars(d, config, fmt.Sprintf("{{LoggingBasePath}}%s", d.Id()))
		if err != nil {
			return err
		}

		_, err = sendRequestWithTimeout(config, "DELETE", "", url, userAgent, nil, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error deleting Logging Sink %q: %s", d.Id(), err)
		}
	}

	d.SetId("")
	return nil
}

func resourceLoggingSinkImportState(idField string) schema.StateFunc {
	return func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
		loggingSinkParentId, err := parseLoggingSinkParentId(d.Id())
		if err != nil {
			return nil, err
		}

		if err := d.Set(idField, loggingSinkParentId); err != nil {
			return nil, fmt.Errorf("Error setting idField: %s", err)
		}

		return []*schema.ResourceData{d}, nil
	}
}
