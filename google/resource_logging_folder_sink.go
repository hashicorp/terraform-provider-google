package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLoggingFolderSink() *schema.Resource {
	schm := &schema.Resource{
		Create: resourceLoggingFolderSinkCreate,
		Read:   resourceLoggingFolderSinkRead,
		Delete: resourceLoggingFolderSinkDelete,
		Update: resourceLoggingFolderSinkUpdate,
		Schema: resourceLoggingSinkSchema(),
		Importer: &schema.ResourceImporter{
			State: resourceLoggingSinkImportState("folder"),
		},
	}
	schm.Schema["folder"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
		StateFunc: func(v interface{}) string {
			return strings.Replace(v.(string), "folders/", "", 1)
		},
	}
	schm.Schema["include_children"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		ForceNew: true,
		Default:  false,
	}

	return schm
}

func resourceLoggingFolderSinkCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	folder := parseFolderId(d.Get("folder"))
	id, sink := expandResourceLoggingSink(d, "folders", folder)
	sink.IncludeChildren = d.Get("include_children").(bool)

	// The API will reject any requests that don't explicitly set 'uniqueWriterIdentity' to true.
	_, err := config.clientLogging.Folders.Sinks.Create(id.parent(), sink).UniqueWriterIdentity(true).Do()
	if err != nil {
		return err
	}

	d.SetId(id.canonicalId())
	return resourceLoggingFolderSinkRead(d, meta)
}

func resourceLoggingFolderSinkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	sink, err := config.clientLogging.Folders.Sinks.Get(d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Folder Logging Sink %s", d.Get("name").(string)))
	}

	flattenResourceLoggingSink(d, sink)
	d.Set("include_children", sink.IncludeChildren)

	return nil
}

func resourceLoggingFolderSinkUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	sink := expandResourceLoggingSinkForUpdate(d)
	// It seems the API might actually accept an update for include_children; this is not in the list of updatable
	// properties though and might break in the future. Always include the value to prevent it changing.
	sink.IncludeChildren = d.Get("include_children").(bool)
	sink.ForceSendFields = append(sink.ForceSendFields, "IncludeChildren")

	// The API will reject any requests that don't explicitly set 'uniqueWriterIdentity' to true.
	_, err := config.clientLogging.Folders.Sinks.Patch(d.Id(), sink).
		UpdateMask(defaultLogSinkUpdateMask).UniqueWriterIdentity(true).Do()
	if err != nil {
		return err
	}

	return resourceLoggingFolderSinkRead(d, meta)
}

func resourceLoggingFolderSinkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	_, err := config.clientLogging.Projects.Sinks.Delete(d.Id()).Do()
	if err != nil {
		return err
	}

	return nil
}
