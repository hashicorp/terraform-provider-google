package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var billingAccountLoggingSinkSchema = map[string]*schema.Schema{
	"billing_account": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: `The billing account exported to the sink.`,
	},
}

func billingAccountLoggingSinkID(d *schema.ResourceData, config *Config) (string, error) {
	billingAccount := d.Get("billing_account").(string)
	sinkName := d.Get("name").(string)
	id := fmt.Sprintf("billingAccounts/%s/sinks/%s", billingAccount, sinkName)
	return id, nil
}

func billingAccountLoggingSinksPath(d *schema.ResourceData, config *Config) (string, error) {
	billingAccount := d.Get("billing_account").(string)
	id := fmt.Sprintf("billingAccounts/%s/sinks", billingAccount)
	return id, nil
}

func resourceLoggingBillingAccountSink() *schema.Resource {
	billingAccountLoggingSinkRead := resourceLoggingSinkRead(flattenBillingAccountLoggingSink)
	billingAccountLoggingSinkCreate := resourceLoggingSinkCreate(billingAccountLoggingSinksPath, expandBillingAccountLoggingSinkForCreate, billingAccountLoggingSinkRead)
	billingAccountLoggingSinkUpdate := resourceLoggingSinkUpdate(expandBillingAccountLoggingSinkForUpdate, billingAccountLoggingSinkRead)

	return &schema.Resource{
		Create: resourceLoggingSinkAcquireOrCreate(billingAccountLoggingSinkID, billingAccountLoggingSinkCreate, billingAccountLoggingSinkUpdate),
		Read:   billingAccountLoggingSinkRead,
		Update: billingAccountLoggingSinkUpdate,
		Delete: resourceLoggingSinkDelete,
		Importer: &schema.ResourceImporter{
			State: resourceLoggingSinkImportState("billing_account"),
		},
		Schema:        resourceLoggingSinkSchema(billingAccountLoggingSinkSchema),
		UseJSONNumber: true,
	}
}

func expandBillingAccountLoggingSinkForCreate(d *schema.ResourceData, config *Config) (obj map[string]interface{}, uniqueWriterIdentity bool) {
	obj = expandLoggingSink(d)
	uniqueWriterIdentity = true
	return
}

func flattenBillingAccountLoggingSink(d *schema.ResourceData, res map[string]interface{}, config *Config) error {
	if err := flattenLoggingSinkBase(d, res); err != nil {
		return err
	}
	return nil
}

func expandBillingAccountLoggingSinkForUpdate(d *schema.ResourceData, config *Config) (obj map[string]interface{}, updateMask string, uniqueWriterIdentity bool) {
	obj, updateMaskList := expandResourceLoggingSinkForUpdateBase(d)
	uniqueWriterIdentity = true
	updateMask = strings.Join(updateMaskList, ",")
	return
}
