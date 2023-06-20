// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var loggingBillingAccountBucketConfigSchema = map[string]*schema.Schema{
	"billing_account": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: `The parent resource that contains the logging bucket.`,
	},
}

func billingAccountBucketConfigID(d *schema.ResourceData, config *transport_tpg.Config) (string, error) {
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
