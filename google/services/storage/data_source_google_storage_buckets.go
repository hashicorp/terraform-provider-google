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

func DataSourceGoogleStorageBuckets() *schema.Resource {
	return &schema.Resource{
		Read: datasourceGoogleStorageBucketsRead,
		Schema: map[string]*schema.Schema{
			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"buckets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"location": {
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

func datasourceGoogleStorageBucketsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	params := make(map[string]string)
	buckets := make([]map[string]interface{}, 0)

	for {
		url := "https://storage.googleapis.com/storage/v1/b"

		params["project"], err = tpgresource.GetProject(d, config)
		if err != nil {
			return fmt.Errorf("Error fetching project for bucket: %s", err)
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
			return fmt.Errorf("Error retrieving buckets: %s", err)
		}

		pageBuckets := flattenDatasourceGoogleBucketsList(res["items"])
		buckets = append(buckets, pageBuckets...)

		pToken, ok := res["nextPageToken"]
		if ok && pToken != nil && pToken.(string) != "" {
			params["pageToken"] = pToken.(string)
		} else {
			break
		}
	}

	if err := d.Set("buckets", buckets); err != nil {
		return fmt.Errorf("Error retrieving buckets: %s", err)
	}

	d.SetId(params["project"])

	return nil
}

func flattenDatasourceGoogleBucketsList(v interface{}) []map[string]interface{} {
	if v == nil {
		return make([]map[string]interface{}, 0)
	}

	ls := v.([]interface{})
	buckets := make([]map[string]interface{}, 0, len(ls))
	for _, raw := range ls {
		o := raw.(map[string]interface{})

		var mLabels, mLocation, mName, mSelfLink, mStorageClass interface{}
		if oLabels, ok := o["labels"]; ok {
			mLabels = oLabels
		}
		if oLocation, ok := o["location"]; ok {
			mLocation = oLocation
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
		buckets = append(buckets, map[string]interface{}{
			"labels":        mLabels,
			"location":      mLocation,
			"name":          mName,
			"self_link":     mSelfLink,
			"storage_class": mStorageClass,
		})
	}

	return buckets
}
