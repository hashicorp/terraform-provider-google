// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var loggingOrganizationBucketConfigSchema = map[string]*schema.Schema{
	"organization": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: `The parent resource that contains the logging bucket.`,
	},
}

func organizationBucketConfigID(d *schema.ResourceData, config *transport_tpg.Config) (string, error) {
	organization := d.Get("organization").(string)
	location := d.Get("location").(string)
	bucketID := d.Get("bucket_id").(string)

	if !strings.HasPrefix(organization, "organization") {
		organization = "organizations/" + organization
	}

	id := fmt.Sprintf("%s/locations/%s/buckets/%s", organization, location, bucketID)
	return id, nil
}

// Create Logging Bucket config
func ResourceLoggingOrganizationBucketConfig() *schema.Resource {
	return ResourceLoggingBucketConfig("organization", loggingOrganizationBucketConfigSchema, organizationBucketConfigID)
}
