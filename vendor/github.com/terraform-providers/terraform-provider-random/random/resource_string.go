package random

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceString() *schema.Resource {
	return &schema.Resource{
		Create:        createStringFunc(false),
		Read:          readNil,
		Delete:        schema.RemoveFromState,
		MigrateState:  resourceRandomStringMigrateState,
		SchemaVersion: 1,
		Schema:        stringSchemaV1(false),
	}
}
