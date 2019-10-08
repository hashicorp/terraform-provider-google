package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/logging/v2"
)

// Empty update masks will eventually cause updates to fail, currently empty masks default to this string
const defaultLogSinkUpdateMask = "destination,filter,includeChildren"

func resourceLoggingSinkSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: optionalSurroundingSpacesSuppress,
		},

		"writer_identity": {
			Type:     schema.TypeString,
			Computed: true,
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
		Name:        d.Get("name").(string),
		Destination: d.Get("destination").(string),
		Filter:      d.Get("filter").(string),
	}
	return id, &sink
}

func flattenResourceLoggingSink(d *schema.ResourceData, sink *logging.LogSink) {
	d.Set("name", sink.Name)
	d.Set("destination", sink.Destination)
	d.Set("filter", sink.Filter)
	d.Set("writer_identity", sink.WriterIdentity)
}

func expandResourceLoggingSinkForUpdate(d *schema.ResourceData) *logging.LogSink {
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
	return &sink
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
