package random

import (
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceUuid() *schema.Resource {
	return &schema.Resource{
		Create: CreateUuid,
		Read:   schema.Noop,
		Delete: schema.RemoveFromState,
		Importer: &schema.ResourceImporter{
			State: ImportUuid,
		},

		Schema: map[string]*schema.Schema{
			"keepers": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"result": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func CreateUuid(d *schema.ResourceData, meta interface{}) error {
	result, err := uuid.GenerateUUID()
	if err != nil {
		return errwrap.Wrapf("error generating uuid: {{err}}", err)
	}
	d.Set("result", result)
	d.SetId(result)
	return nil
}

func ImportUuid(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()
	bytes, err := uuid.ParseUUID(id)
	if err != nil {
		return nil, errwrap.Wrapf("error parsing uuid bytes: {{err}}", err)
	}
	result, err2 := uuid.FormatUUID(bytes)
	if err2 != nil {
		return nil, errwrap.Wrapf("error formatting uuid bytes: {{err2}}", err2)
	}

	d.Set("result", result)
	d.SetId(result)

	return []*schema.ResourceData{d}, nil
}
