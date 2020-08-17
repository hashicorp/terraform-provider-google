package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var loggingBillingAccountBucketConfigSchema = map[string]*schema.Schema{
	"billing_account": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: `The parent resource that contains the logging bucket.`,
	},
}

func billingAccountBucketConfigID(d *schema.ResourceData, config *Config) (string, error) {
	billingAccount := d.Get("billing_account").(string)
	location := d.Get("location").(string)
	bucketID := d.Get("bucket_id").(string)

	if !strings.HasPrefix(billingAccount, "billingAccounts") {
		billingAccount = "billingAccounts/" + billingAccount
	}

	id := fmt.Sprintf("%s/locations/%s/buckets/%s", billingAccount, location, bucketID)
	return id, nil
}

// Create Logging Bucket config
func ResourceLoggingBillingAccountBucketConfig() *schema.Resource {
	return ResourceLoggingBucketConfig("billing_account", loggingBillingAccountBucketConfigSchema, billingAccountBucketConfigID)
}
