// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

var loggingFolderBucketConfigSchema = map[string]*schema.Schema{
	"folder": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: `The parent resource that contains the logging bucket.`,
	},
}

func folderBucketConfigID(d *schema.ResourceData, config *transport_tpg.Config) (string, error) {
	folder := d.Get("folder").(string)
	location := d.Get("location").(string)
	bucketID := d.Get("bucket_id").(string)

	if !strings.HasPrefix(folder, "folder") {
		folder = "folders/" + folder
	}

	id := fmt.Sprintf("%s/locations/%s/buckets/%s", folder, location, bucketID)
	return id, nil
}

// Create Logging Bucket config
func ResourceLoggingFolderBucketConfig() *schema.Resource {
	return ResourceLoggingBucketConfig("folder", loggingFolderBucketConfigSchema, folderBucketConfigID)
}
