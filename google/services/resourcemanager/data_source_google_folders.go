// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package resourcemanager

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleFolders() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGoogleFoldersRead,
		Schema: map[string]*schema.Schema{
			"parent_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"folders": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delete_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"etag": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGoogleFoldersRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	params := make(map[string]string)
	folders := make([]map[string]interface{}, 0)

	for {
		params["parent"] = d.Get("parent_id").(string)
		url := "https://cloudresourcemanager.googleapis.com/v3/folders"

		url, err := transport_tpg.AddQueryParams(url, params)
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
			return fmt.Errorf("Error retrieving folders: %s", err)
		}

		pageFolders := flattenDataSourceGoogleFoldersList(res["folders"])
		folders = append(folders, pageFolders...)

		pToken, ok := res["nextPageToken"]
		if ok && pToken != nil && pToken.(string) != "" {
			params["pageToken"] = pToken.(string)
		} else {
			break
		}
	}

	if err := d.Set("folders", folders); err != nil {
		return fmt.Errorf("Error retrieving folders: %s", err)
	}

	d.SetId(d.Get("parent_id").(string))

	return nil
}

func flattenDataSourceGoogleFoldersList(v interface{}) []map[string]interface{} {
	if v == nil {
		return make([]map[string]interface{}, 0)
	}

	ls := v.([]interface{})
	folders := make([]map[string]interface{}, 0, len(ls))
	for _, raw := range ls {
		f := raw.(map[string]interface{})

		var mState, mName, mCreateTime, mUpdateTime, mDeleteTime, mParent, mDisplayName, mEtag interface{}
		if fName, ok := f["name"]; ok {
			mName = fName
		}
		if fState, ok := f["state"]; ok {
			mState = fState
		}
		if fCreateTime, ok := f["createTime"]; ok {
			mCreateTime = fCreateTime
		}
		if fUpdateTime, ok := f["updateTime"]; ok {
			mUpdateTime = fUpdateTime
		}
		if fDeleteTime, ok := f["deleteTime"]; ok {
			mDeleteTime = fDeleteTime
		}
		if fParent, ok := f["parent"]; ok {
			mParent = fParent
		}
		if fDisplayName, ok := f["displayName"]; ok {
			mDisplayName = fDisplayName
		}
		if fEtag, ok := f["etag"]; ok {
			mEtag = fEtag
		}
		folders = append(folders, map[string]interface{}{
			"name":         mName,
			"state":        mState,
			"create_time":  mCreateTime,
			"update_time":  mUpdateTime,
			"delete_time":  mDeleteTime,
			"parent":       mParent,
			"display_name": mDisplayName,
			"etag":         mEtag,
		})
	}

	return folders
}
