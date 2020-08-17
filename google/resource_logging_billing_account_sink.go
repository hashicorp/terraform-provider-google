package google

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLoggingBillingAccountSink() *schema.Resource {
	schm := &schema.Resource{
		Create: resourceLoggingBillingAccountSinkCreate,
		Read:   resourceLoggingBillingAccountSinkRead,
		Delete: resourceLoggingBillingAccountSinkDelete,
		Update: resourceLoggingBillingAccountSinkUpdate,
		Schema: resourceLoggingSinkSchema(),
		Importer: &schema.ResourceImporter{
			State: resourceLoggingSinkImportState("billing_account"),
		},
	}
	schm.Schema["billing_account"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: `The billing account exported to the sink.`,
	}
	return schm
}

func resourceLoggingBillingAccountSinkCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	id, sink := expandResourceLoggingSink(d, "billingAccounts", d.Get("billing_account").(string))

	// The API will reject any requests that don't explicitly set 'uniqueWriterIdentity' to true.
	_, err := config.clientLogging.BillingAccounts.Sinks.Create(id.parent(), sink).UniqueWriterIdentity(true).Do()
	if err != nil {
		return err
	}

	d.SetId(id.canonicalId())
	return resourceLoggingBillingAccountSinkRead(d, meta)
}

func resourceLoggingBillingAccountSinkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	sink, err := config.clientLogging.BillingAccounts.Sinks.Get(d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Billing Logging Sink %s", d.Get("name").(string)))
	}

	flattenResourceLoggingSink(d, sink)
	return nil

}

func resourceLoggingBillingAccountSinkUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	sink, updateMask := expandResourceLoggingSinkForUpdate(d)

	// The API will reject any requests that don't explicitly set 'uniqueWriterIdentity' to true.
	_, err := config.clientLogging.BillingAccounts.Sinks.Patch(d.Id(), sink).
		UpdateMask(updateMask).UniqueWriterIdentity(true).Do()
	if err != nil {
		return err
	}

	return resourceLoggingBillingAccountSinkRead(d, meta)
}

func resourceLoggingBillingAccountSinkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	_, err := config.clientLogging.Projects.Sinks.Delete(d.Id()).Do()
	if err != nil {
		return err
	}

	return nil
}
