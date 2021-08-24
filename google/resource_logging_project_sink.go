package google

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const nonUniqueWriterAccount = "serviceAccount:cloud-logs@system.gserviceaccount.com"

var projectLoggingSinkSchema = map[string]*schema.Schema{
	"project": {
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
		Description: `The ID of the project to create the sink in. If omitted, the project associated with the provider is used.`,
	},
	"unique_writer_identity": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		ForceNew:    true,
		Description: `Whether or not to create a unique identity associated with this sink. If false (the default), then the writer_identity used is serviceAccount:cloud-logs@system.gserviceaccount.com. If true, then a unique service account is created and used for this sink. If you wish to publish logs across projects, you must set unique_writer_identity to true.`,
	},
}

func projectLoggingSinkID(d *schema.ResourceData, config *Config) (string, error) {
	project, _ := getProject(d, config)
	sinkName := d.Get("name").(string)
	id := fmt.Sprintf("projects/%s/sinks/%s", project, sinkName)
	return id, nil
}

func projectLoggingSinksPath(d *schema.ResourceData, config *Config) (string, error) {
	project, _ := getProject(d, config)
	id := fmt.Sprintf("projects/%s/sinks", project)
	return id, nil
}

func resourceLoggingProjectSink() *schema.Resource {
	projectLoggingSinkRead := resourceLoggingSinkRead(flattenProjectLoggingSink)
	projectLoggingSinkCreate := resourceLoggingSinkCreate(projectLoggingSinksPath, expandProjectLoggingSinkForCreate, projectLoggingSinkRead)
	projectLoggingSinkUpdate := resourceLoggingSinkUpdate(expandProjectLoggingSinkForUpdate, projectLoggingSinkRead)

	return &schema.Resource{
		Create: resourceLoggingSinkAcquireOrCreate(projectLoggingSinkID, projectLoggingSinkCreate, projectLoggingSinkUpdate),
		Read:   projectLoggingSinkRead,
		Update: projectLoggingSinkUpdate,
		Delete: resourceLoggingSinkDelete,
		Importer: &schema.ResourceImporter{
			State: resourceLoggingSinkImportState("project"),
		},
		CustomizeDiff: resourceLoggingProjectSinkCustomizeDiff,
		Schema:        resourceLoggingSinkSchema(projectLoggingSinkSchema),
		UseJSONNumber: true,
	}
}

func expandProjectLoggingSinkForCreate(d *schema.ResourceData, config *Config) (obj map[string]interface{}, uniqueWriterIdentity bool) {
	obj = expandLoggingSink(d)
	uniqueWriterIdentity = d.Get("unique_writer_identity").(bool)
	return
}

func flattenProjectLoggingSink(d *schema.ResourceData, res map[string]interface{}, config *Config) error {
	if err := flattenLoggingSinkBase(d, res); err != nil {
		return err
	}

	project, _ := getProject(d, config)
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	if res["writerIdentity"] != nonUniqueWriterAccount {
		if err := d.Set("unique_writer_identity", true); err != nil {
			return fmt.Errorf("Error setting unique_writer_identity: %s", err)
		}
	} else {
		if err := d.Set("unique_writer_identity", false); err != nil {
			return fmt.Errorf("Error setting unique_writer_identity: %s", err)
		}
	}
	return nil
}

func expandProjectLoggingSinkForUpdate(d *schema.ResourceData, config *Config) (obj map[string]interface{}, updateMask string, uniqueWriterIdentity bool) {
	obj, updateFields := expandResourceLoggingSinkForUpdateBase(d)
	uniqueWriterIdentity = d.Get("unique_writer_identity").(bool)
	updateMask = strings.Join(updateFields, ",")
	return
}

// if bigquery_options is set unique_writer_identity must be true
func resourceLoggingProjectSinkCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// separate func to allow unit testing
	return resourceLoggingProjectSinkCustomizeDiffFunc(d)
}

func resourceLoggingProjectSinkCustomizeDiffFunc(diff TerraformResourceDiff) error {
	if !diff.HasChange("bigquery_options.#") {
		return nil
	}

	bigqueryOptions := diff.Get("bigquery_options.#").(int)
	if bigqueryOptions > 0 {
		uwi := diff.Get("unique_writer_identity")
		if !uwi.(bool) {
			return errors.New("unique_writer_identity must be true when bigquery_options is supplied")
		}
	}
	return nil
}
