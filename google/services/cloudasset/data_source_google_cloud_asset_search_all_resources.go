// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudasset

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleCloudAssetSearchAllResources() *schema.Resource {
	return &schema.Resource{
		Read: datasourceGoogleCloudAssetSearchAllResourcesRead,
		Schema: map[string]*schema.Schema{
			"scope": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"asset_types": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"asset_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"folders": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"organization": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"location": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"network_tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"kms_keys": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_full_resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_asset_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceGoogleCloudAssetSearchAllResourcesRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	params := make(map[string]string)
	results := make([]map[string]interface{}, 0)

	scope := d.Get("scope").(string)
	query := d.Get("query").(string)
	assetTypes := d.Get("asset_types").([]interface{})

	url := fmt.Sprintf("https://cloudasset.googleapis.com/v1/%s:searchAllResources", scope)
	params["query"] = query

	url, err = transport_tpg.AddArrayQueryParams(url, "asset_types", assetTypes)
	if err != nil {
		return fmt.Errorf("Error setting asset_types: %s", err)
	}

	for {
		url, err := transport_tpg.AddQueryParams(url, params)
		if err != nil {
			return err
		}

		var project string
		if config.UserProjectOverride && config.BillingProject != "" {
			project = config.BillingProject
		}

		res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Project:   project,
			Method:    "GET",
			RawURL:    url,
			UserAgent: userAgent,
		})
		if err != nil {
			return fmt.Errorf("Error searching resources: %s", err)
		}

		pageResults := flattenDatasourceGoogleCloudAssetSearchAllResources(res["results"])
		results = append(results, pageResults...)

		pToken, ok := res["nextPageToken"]
		if ok && pToken != nil && pToken.(string) != "" {
			params["pageToken"] = pToken.(string)
		} else {
			break
		}
	}

	if err := d.Set("results", results); err != nil {
		return fmt.Errorf("Error searching resources: %s", err)
	}

	if err := d.Set("query", query); err != nil {
		return fmt.Errorf("Error setting query: %s", err)
	}

	if err := d.Set("asset_types", assetTypes); err != nil {
		return fmt.Errorf("Error setting asset_types: %s", err)
	}

	d.SetId(scope)

	return nil
}

func flattenDatasourceGoogleCloudAssetSearchAllResources(v interface{}) []map[string]interface{} {
	if v == nil {
		return make([]map[string]interface{}, 0)
	}

	ls := v.([]interface{})
	results := make([]map[string]interface{}, 0, len(ls))
	for _, raw := range ls {
		p := raw.(map[string]interface{})

		var mName, mAssetType, mProject, mFolders, mOrganization, mDisplayName, mDescription, mLocation, mLabels, mNetworkTags, mKmsKeys, mCreateTime, mUpdateTime, mState, mParentFullResourceName, mParentAssetType interface{}
		if pName, ok := p["name"]; ok {
			mName = pName
		}
		if pAssetType, ok := p["assetType"]; ok {
			mAssetType = pAssetType
		}
		if pProject, ok := p["project"]; ok {
			mProject = pProject
		}
		if pFolders, ok := p["folders"]; ok {
			mFolders = pFolders
		}
		if pOrganization, ok := p["organization"]; ok {
			mOrganization = pOrganization
		}
		if pDisplayName, ok := p["displayName"]; ok {
			mDisplayName = pDisplayName
		}
		if pDescription, ok := p["description"]; ok {
			mDescription = pDescription
		}
		if pLocation, ok := p["location"]; ok {
			mLocation = pLocation
		}
		if pLabels, ok := p["labels"]; ok {
			mLabels = pLabels
		}
		if pNetworkTags, ok := p["networkTags"]; ok {
			mNetworkTags = pNetworkTags
		}
		if pKmsKeys, ok := p["kmsKeys"]; ok {
			mKmsKeys = pKmsKeys
		}
		if pCreateTime, ok := p["createTime"]; ok {
			mCreateTime = pCreateTime
		}
		if pUpdateTime, ok := p["updateTime"]; ok {
			mUpdateTime = pUpdateTime
		}
		if pState, ok := p["state"]; ok {
			mState = pState
		}
		if pParentFullResourceName, ok := p["parentFullResourceName"]; ok {
			mParentFullResourceName = pParentFullResourceName
		}
		if pParentAssetType, ok := p["parentAssetType"]; ok {
			mParentAssetType = pParentAssetType
		}
		results = append(results, map[string]interface{}{
			"name":                      mName,
			"asset_type":                mAssetType,
			"project":                   mProject,
			"folders":                   mFolders,
			"organization":              mOrganization,
			"display_name":              mDisplayName,
			"description":               mDescription,
			"location":                  mLocation,
			"labels":                    mLabels,
			"network_tags":              mNetworkTags,
			"kms_keys":                  mKmsKeys,
			"create_time":               mCreateTime,
			"update_time":               mUpdateTime,
			"state":                     mState,
			"parent_full_resource_name": mParentFullResourceName,
			"parent_asset_type":         mParentAssetType,
		})
	}

	return results
}
