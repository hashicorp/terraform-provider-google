package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var loggingOrganizationBucketConfigSchema = map[string]*schema.Schema{
	"organization": {
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	},
}

func organizationBucketConfigID(d *schema.ResourceData, config *Config) (string, error) {
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
