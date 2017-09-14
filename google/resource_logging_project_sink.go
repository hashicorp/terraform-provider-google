package google

import (
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/logging/v2"
)

const nonUniqueWriterAccount = "serviceAccount:cloud-logs@system.gserviceaccount.com"

func resourceLoggingProjectSink() *schema.Resource {
	return &schema.Resource{
		Create: resourceLoggingProjectSinkCreate,
		Read:   resourceLoggingProjectSinkRead,
		Delete: resourceLoggingProjectSinkDelete,
		Update: resourceLoggingProjectSinkUpdate,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destination": {
				Type:     schema.TypeString,
				Required: true,
			},

			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"unique_writer_identity": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"writer_identity": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceLoggingProjectSinkCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	id := LoggingSinkId{
		resourceType: "projects",
		resourceId:   project,
		name:         name,
	}

	sink := logging.LogSink{
		Name:        d.Get("name").(string),
		Destination: d.Get("destination").(string),
		Filter:      d.Get("filter").(string),
	}

	uniqueWriterIdentity := d.Get("unique_writer_identity").(bool)

	_, err = config.clientLogging.Projects.Sinks.Create(id.parent(), &sink).UniqueWriterIdentity(uniqueWriterIdentity).Do()
	if err != nil {
		return err
	}

	d.SetId(id.canonicalId())

	return resourceLoggingProjectSinkRead(d, meta)
}

func resourceLoggingProjectSinkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	sink, err := config.clientLogging.Projects.Sinks.Get(d.Id()).Do()
	if err != nil {
		return err
	}

	d.Set("name", sink.Name)
	d.Set("destination", sink.Destination)
	d.Set("filter", sink.Filter)
	d.Set("writer_identity", sink.WriterIdentity)
	if sink.WriterIdentity != nonUniqueWriterAccount {
		d.Set("unique_writer_identity", true)
	} else {
		d.Set("unique_writer_identity", false)
	}
	return nil
}

func resourceLoggingProjectSinkUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Can only update destination/filter right now. Despite the method below using 'Patch', the API requires both
	// destination and filter (even if unchanged).
	sink := logging.LogSink{
		Destination: d.Get("destination").(string),
		Filter:      d.Get("filter").(string),
	}

	if d.HasChange("destination") {
		sink.ForceSendFields = append(sink.ForceSendFields, "Destination")
	}
	if d.HasChange("filter") {
		sink.ForceSendFields = append(sink.ForceSendFields, "Filter")
	}

	uniqueWriterIdentity := d.Get("unique_writer_identity").(bool)

	_, err := config.clientLogging.Projects.Sinks.Patch(d.Id(), &sink).UniqueWriterIdentity(uniqueWriterIdentity).Do()
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
