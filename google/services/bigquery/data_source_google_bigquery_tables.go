// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigquery

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleBigQueryTables() *schema.Resource {

	dsSchema := map[string]*schema.Schema{
		"dataset_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The ID of the dataset containing the tables.",
		},
		"project": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The ID of the project in which the dataset is located. If it is not provided, the provider project is used.",
		},
		"tables": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"labels": {
						Type:     schema.TypeMap,
						Computed: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"table_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}

	return &schema.Resource{
		Read:   DataSourceGoogleBigQueryTablesRead,
		Schema: dsSchema,
	}
}

func DataSourceGoogleBigQueryTablesRead(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	datasetID := d.Get("dataset_id").(string)

	project, err := tpgresource.GetProject(d, config)

	if err != nil {
		return fmt.Errorf("Error fetching project: %s", err)
	}

	params := make(map[string]string)
	tables := make([]map[string]interface{}, 0)

	for {

		url, err := tpgresource.ReplaceVars(d, config, "{{BigQueryBasePath}}projects/{{project}}/datasets/{{dataset_id}}/tables")
		if err != nil {
			return err
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
			return fmt.Errorf("Error retrieving tables: %s", err)
		}

		pageTables := flattenDataSourceGoogleBigQueryTablesList(res["tables"])
		tables = append(tables, pageTables...)

		pToken, ok := res["nextPageToken"]
		if ok && pToken != nil && pToken.(string) != "" {
			params["pageToken"] = pToken.(string)
		} else {
			break
		}
	}

	if err := d.Set("tables", tables); err != nil {
		return fmt.Errorf("Error retrieving tables: %s", err)
	}

	id := fmt.Sprintf("projects/%s/datasets/%s/tables", project, datasetID)
	d.SetId(id)

	return nil
}

func flattenDataSourceGoogleBigQueryTablesList(res interface{}) []map[string]interface{} {

	if res == nil {
		return make([]map[string]interface{}, 0)
	}

	ls := res.([]interface{})

	tables := make([]map[string]interface{}, 0, len(ls))

	for _, raw := range ls {
		output := raw.(map[string]interface{})

		var mLabels map[string]interface{}
		var mTableName string

		if oLabels, ok := output["labels"].(map[string]interface{}); ok {
			mLabels = oLabels
		} else {
			mLabels = make(map[string]interface{}) // Initialize as an empty map if labels are missing
		}

		if oTableReference, ok := output["tableReference"].(map[string]interface{}); ok {
			if tableID, ok := oTableReference["tableId"].(string); ok {
				mTableName = tableID
			}
		}
		tables = append(tables, map[string]interface{}{
			"labels":   mLabels,
			"table_id": mTableName,
		})
	}

	return tables
}
