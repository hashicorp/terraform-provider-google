// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package alloydb

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceAlloydbSupportedDatabaseFlags() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceAlloydbSupportedDatabaseFlagsRead,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Project ID of the project.`,
			},
			"location": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The canonical id for the location. For example: "us-east1".`,
			},
			"supported_database_flags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: `The name of the flag resource, following Google Cloud conventions, e.g.: * projects/{project}/locations/{location}/flags/{flag} This field currently has no semantic meaning.`,
						},
						"flag_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: `The name of the database flag, e.g. "max_allowed_packets". The is a possibly key for the Instance.database_flags map field.`,
						},
						"value_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: `ValueType describes the semantic type of the value that the flag accepts. The supported values are:- 'VALUE_TYPE_UNSPECIFIED', 'STRING', 'INTEGER', 'FLOAT', 'NONE'.`,
						},
						"accepts_multiple_values": {
							Type:        schema.TypeBool,
							Computed:    true,
							Optional:    true,
							Description: `Whether the database flag accepts multiple values. If true, a comma-separated list of stringified values may be specified.`,
						},
						"supported_db_versions": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: `Major database engine versions for which this flag is supported. Supported values are:- 'DATABASE_VERSION_UNSPECIFIED', and 'POSTGRES_14'.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"requires_db_restart": {
							Type:        schema.TypeBool,
							Computed:    true,
							Optional:    true,
							Description: `Whether setting or updating this flag on an Instance requires a database restart. If a flag that requires database restart is set, the backend will automatically restart the database (making sure to satisfy any availability SLO's).`,
						},
						"string_restrictions": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: `Restriction on STRING type value.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allowed_values": {
										Type:        schema.TypeList,
										Computed:    true,
										Optional:    true,
										Description: `The list of allowed values, if bounded. This field will be empty if there is a unbounded number of allowed values.`,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"integer_restrictions": {
							Type:        schema.TypeList,
							Computed:    true,
							Optional:    true,
							Description: `Restriction on INTEGER type value.`,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Optional:    true,
										Description: `The minimum value that can be specified, if applicable.`,
									},
									"max_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Optional:    true,
										Description: `The maximum value that can be specified, if applicable.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlloydbSupportedDatabaseFlagsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	location := d.Get("location").(string)

	billingProject := ""

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for Cluster: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{AlloydbBasePath}}projects/{{project}}/locations/{{location}}/supportedDatabaseFlags")
	if err != nil {
		return fmt.Errorf("Error setting api endpoint")
	}
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("SupportedDatabaseFlags %q", d.Id()))
	}
	var supportedDatabaseFlags []map[string]interface{}
	for {
		result := res["supportedDatabaseFlags"].([]interface{})
		for _, dbFlag := range result {
			supportedDatabaseFlag := make(map[string]interface{})
			flag := dbFlag.(map[string]interface{})
			if flag["name"] != nil {
				supportedDatabaseFlag["name"] = flag["name"].(string)
			}
			if flag["flagName"] != nil {
				supportedDatabaseFlag["flag_name"] = flag["flagName"].(string)
			}
			if flag["valueType"] != nil {
				supportedDatabaseFlag["value_type"] = flag["valueType"].(string)
			}
			if flag["acceptsMultipleValues"] != nil {
				supportedDatabaseFlag["accepts_multiple_values"] = flag["acceptsMultipleValues"].(bool)
			}
			if flag["requiresDbRestart"] != nil {
				supportedDatabaseFlag["requires_db_restart"] = flag["requiresDbRestart"].(bool)
			}
			if flag["supportedDbVersions"] != nil {
				dbVersions := make([]string, 0, len(flag["supportedDbVersions"].([]interface{})))
				for _, supDbVer := range flag["supportedDbVersions"].([]interface{}) {
					dbVersions = append(dbVersions, supDbVer.(string))
				}
				supportedDatabaseFlag["supported_db_versions"] = dbVersions
			}

			if flag["stringRestrictions"] != nil {
				restrictions := make([]map[string][]string, 0, 1)
				fetchedAllowedValues := flag["stringRestrictions"].(map[string]interface{})["allowedValues"]
				if fetchedAllowedValues != nil {
					allowedValues := make([]string, 0, len(fetchedAllowedValues.([]interface{})))
					for _, val := range fetchedAllowedValues.([]interface{}) {
						allowedValues = append(allowedValues, val.(string))
					}
					stringRestrictions := map[string][]string{
						"allowed_values": allowedValues,
					}
					restrictions = append(restrictions, stringRestrictions)
					supportedDatabaseFlag["string_restrictions"] = restrictions
				}
			}
			if flag["integerRestrictions"] != nil {
				restrictions := make([]map[string]string, 0, 1)
				minValue := flag["integerRestrictions"].(map[string]interface{})["minValue"].(string)
				maxValue := flag["integerRestrictions"].(map[string]interface{})["maxValue"].(string)
				integerRestrictions := map[string]string{
					"min_value": minValue,
					"max_value": maxValue,
				}
				restrictions = append(restrictions, integerRestrictions)
				supportedDatabaseFlag["integer_restrictions"] = restrictions
			}
			supportedDatabaseFlags = append(supportedDatabaseFlags, supportedDatabaseFlag)
		}
		if res["pageToken"] == nil || res["pageToken"].(string) == "" {
			break
		}
		url, err = tpgresource.ReplaceVars(d, config, "{{AlloydbBasePath}}projects/{{project}}/locations/{{location}}/supportedDatabaseFlags?pageToken="+res["nextPageToken"].(string))
		if err != nil {
			return fmt.Errorf("Error setting api endpoint")
		}
		res, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
			Config:    config,
			Method:    "GET",
			Project:   billingProject,
			RawURL:    url,
			UserAgent: userAgent,
		})
		if err != nil {
			return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("SupportedDatabaseFlags %q", d.Id()))
		}
	}
	if err := d.Set("supported_database_flags", supportedDatabaseFlags); err != nil {
		return fmt.Errorf("Error setting supported_database_flags: %s", err)
	}
	d.SetId(fmt.Sprintf("projects/%s/locations/%s/supportedDbFlags", project, location))
	return nil
}
