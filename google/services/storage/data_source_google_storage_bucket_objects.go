// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleStorageBucketObjects() *schema.Resource {
	return &schema.Resource{
		Read: datasourceGoogleStorageBucketObjectsRead,
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"match_glob": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bucket_objects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"media_link": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"self_link": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceGoogleStorageBucketObjectsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	params := make(map[string]string)
	bucketObjects := make([]map[string]interface{}, 0)

	for {
		bucket := d.Get("bucket").(string)
		url, err := tpgresource.ReplaceVars(d, config, fmt.Sprintf("{{StorageBasePath}}b/%s/o", bucket))
		if err != nil {
			return err
		}

		if v, ok := d.GetOk("match_glob"); ok {
			params["matchGlob"] = v.(string)
		}

		if v, ok := d.GetOk("prefix"); ok {
			params["prefix"] = v.(string)
		}

		url, err = transport_tpg.AddQueryParams(url, params)
		if err != nil {
			return err
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			RawURL:    url,
			UserAgent: userAgent,
		})
		if err != nil {
			return fmt.Errorf("Error retrieving bucket objects: %s", err)
		}

		pageBucketObjects := flattenDatasourceGoogleBucketObjectsList(res["items"])
		bucketObjects = append(bucketObjects, pageBucketObjects...)

		pToken, ok := res["nextPageToken"]
		if ok && pToken != nil && pToken.(string) != "" {
			params["pageToken"] = pToken.(string)
		} else {
			break
		}
	}

	if err := d.Set("bucket_objects", bucketObjects); err != nil {
		return fmt.Errorf("Error retrieving bucket_objects: %s", err)
	}

	d.SetId(d.Get("bucket").(string))

	return nil
}

func flattenDatasourceGoogleBucketObjectsList(v interface{}) []map[string]interface{} {
	if v == nil {
		return make([]map[string]interface{}, 0)
	}

	ls := v.([]interface{})
	bucketObjects := make([]map[string]interface{}, 0, len(ls))
	for _, raw := range ls {
		o := raw.(map[string]interface{})

		var mContentType, mMediaLink, mName, mSelfLink, mStorageClass interface{}
		if oContentType, ok := o["contentType"]; ok {
			mContentType = oContentType
		}
		if oMediaLink, ok := o["mediaLink"]; ok {
			mMediaLink = oMediaLink
		}
		if oName, ok := o["name"]; ok {
			mName = oName
		}
		if oSelfLink, ok := o["selfLink"]; ok {
			mSelfLink = oSelfLink
		}
		if oStorageClass, ok := o["storageClass"]; ok {
			mStorageClass = oStorageClass
		}
		bucketObjects = append(bucketObjects, map[string]interface{}{
			"content_type":  mContentType,
			"media_link":    mMediaLink,
			"name":          mName,
			"self_link":     mSelfLink,
			"storage_class": mStorageClass,
		})
	}

	return bucketObjects
}
