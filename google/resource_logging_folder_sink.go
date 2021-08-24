package google

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

var folderLoggingSinkSchema = map[string]*schema.Schema{
	"folder": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: `The folder to be exported to the sink. Note that either [FOLDER_ID] or "folders/[FOLDER_ID]" is accepted.`,
		StateFunc: func(v interface{}) string {
			return strings.Replace(v.(string), "folders/", "", 1)
		},
	},
	"include_children": {
		Type:        schema.TypeBool,
		Optional:    true,
		ForceNew:    true,
		Default:     false,
		Description: `Whether or not to include children folders in the sink export. If true, logs associated with child projects are also exported; otherwise only logs relating to the provided folder are included.`,
	},
}

func folderLoggingSinkID(d *schema.ResourceData, config *Config) (string, error) {
	folder := d.Get("folder").(string)
	sinkName := d.Get("name").(string)

	if !strings.HasPrefix(folder, "folder") {
		folder = "folders/" + folder
	}

	id := fmt.Sprintf("%s/sinks/%s", folder, sinkName)
	return id, nil
}

func folderLoggingSinksPath(d *schema.ResourceData, config *Config) (string, error) {
	folder := d.Get("folder").(string)

	if !strings.HasPrefix(folder, "folder") {
		folder = "folders/" + folder
	}

	id := fmt.Sprintf("%s/sinks", folder)
	return id, nil
}

func resourceLoggingFolderSink() *schema.Resource {
	folderLoggingSinkRead := resourceLoggingSinkRead(flattenFolderLoggingSink)
	folderLoggingSinkCreate := resourceLoggingSinkCreate(folderLoggingSinksPath, expandFolderLoggingSinkForCreate, folderLoggingSinkRead)
	folderLoggingSinkUpdate := resourceLoggingSinkUpdate(expandFolderLoggingSinkForUpdate, folderLoggingSinkRead)

	return &schema.Resource{
		Create: resourceLoggingSinkAcquireOrCreate(folderLoggingSinkID, folderLoggingSinkCreate, folderLoggingSinkUpdate),
		Read:   folderLoggingSinkRead,
		Update: folderLoggingSinkUpdate,
		Delete: resourceLoggingSinkDelete,
		Importer: &schema.ResourceImporter{
			State: resourceLoggingSinkImportState("folder"),
		},
		Schema:        resourceLoggingSinkSchema(folderLoggingSinkSchema),
		UseJSONNumber: true,
	}
}

func expandFolderLoggingSinkForCreate(d *schema.ResourceData, config *Config) (obj map[string]interface{}, uniqueWriterIdentity bool) {
	obj = expandLoggingSink(d)

	obj["includeChildren"] = d.Get("include_children").(bool)
	uniqueWriterIdentity = true
	return
}

func flattenFolderLoggingSink(d *schema.ResourceData, res map[string]interface{}, config *Config) error {
	if err := flattenLoggingSinkBase(d, res); err != nil {
		return err
	}

	if err := d.Set("include_children", res["includeChildren"]); err != nil {
		return fmt.Errorf("Error setting include_children: %s", err)
	}

	return nil
}

func expandFolderLoggingSinkForUpdate(d *schema.ResourceData, config *Config) (obj map[string]interface{}, updateMask string, uniqueWriterIdentity bool) {
	obj, updateFields := expandResourceLoggingSinkForUpdateBase(d)

	obj["includeChildren"] = d.Get("include_children").(bool)
	updateFields = append(updateFields, "include_children")

	updateMask = strings.Join(updateFields, ",")
	uniqueWriterIdentity = true
	return
}
