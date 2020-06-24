package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const nonUniqueWriterAccount = "serviceAccount:cloud-logs@system.gserviceaccount.com"

func resourceLoggingProjectSink() *schema.Resource {
	schm := &schema.Resource{
		Create: resourceLoggingProjectSinkCreate,
		Read:   resourceLoggingProjectSinkRead,
		Delete: resourceLoggingProjectSinkDelete,
		Update: resourceLoggingProjectSinkUpdate,
		Schema: resourceLoggingSinkSchema(),
		Importer: &schema.ResourceImporter{
			State: resourceLoggingSinkImportState("project"),
		},
	}
	schm.Schema["project"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
		Description: `The ID of the project to create the sink in. If omitted, the project associated with the provider is used.`,
	}
	schm.Schema["unique_writer_identity"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		ForceNew:    true,
		Description: `Whether or not to create a unique identity associated with this sink. If false (the default), then the writer_identity used is serviceAccount:cloud-logs@system.gserviceaccount.com. If true, then a unique service account is created and used for this sink. If you wish to publish logs across projects, you must set unique_writer_identity to true.`,
	}
	return schm
}

func resourceLoggingProjectSinkCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	id, sink := expandResourceLoggingSink(d, "projects", project)
	uniqueWriterIdentity := d.Get("unique_writer_identity").(bool)

	_, err = config.clientLogging.Projects.Sinks.Create(id.parent(), sink).UniqueWriterIdentity(uniqueWriterIdentity).Do()
	if err != nil {
		return err
	}

	d.SetId(id.canonicalId())

	return resourceLoggingProjectSinkRead(d, meta)
}

func resourceLoggingProjectSinkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	sink, err := config.clientLogging.Projects.Sinks.Get(d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Project Logging Sink %s", d.Get("name").(string)))
	}

	d.Set("project", project)
	flattenResourceLoggingSink(d, sink)
	if sink.WriterIdentity != nonUniqueWriterAccount {
		d.Set("unique_writer_identity", true)
	} else {
		d.Set("unique_writer_identity", false)
	}
	return nil
}

func resourceLoggingProjectSinkUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	sink, updateMask := expandResourceLoggingSinkForUpdate(d)
	uniqueWriterIdentity := d.Get("unique_writer_identity").(bool)

	_, err := config.clientLogging.Projects.Sinks.Patch(d.Id(), sink).
		UpdateMask(updateMask).UniqueWriterIdentity(uniqueWriterIdentity).Do()
	if err != nil {
		return err
	}

	return resourceLoggingProjectSinkRead(d, meta)
}

func resourceLoggingProjectSinkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	_, err := config.clientLogging.Projects.Sinks.Delete(d.Id()).Do()
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
