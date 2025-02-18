// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package backupdr

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"strings"
)

func DataSourceGoogleCloudBackupDRService() *schema.Resource {

	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceBackupDRManagementServer().Schema)
	tpgresource.AddRequiredFieldsToSchema(dsSchema, "location")

	return &schema.Resource{
		Read:   dataSourceGoogleCloudBackupDRServiceRead,
		Schema: dsSchema,
	}
}

func flattenBackupDRManagementServerResourceResp(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) map[string]interface{} {
	if v == nil {
		fmt.Printf("Interface is nil: %s", v)
	}
	fmt.Printf("Interface is : %s", v)
	l := v.([]interface{})
	for _, raw := range l {
		// Management server is a singleton resource. It is only present in one location per project. Hence returning only resource present.
		return flattenBackupDRManagementServerResource(raw, d, config)
	}
	return nil
}
func flattenBackupDRManagementServerResource(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) map[string]interface{} {
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["type"] = flattenBackupDRManagementServerType(original["type"], d, config)
	transformed["networks"] = flattenBackupDRManagementServerNetworks(original["networks"], d, config)
	transformed["oauth2ClientId"] = flattenBackupDRManagementServerOauth2ClientId(original["oauth2ClientId"], d, config)
	transformed["managementUri"] = flattenBackupDRManagementServerManagementUri(original["managementUri"], d, config)
	transformed["name"] = flattenBackupDRManagementServerName(original["name"], d, config)
	return transformed
}

func flattenBackupDRManagementServerName(v interface{}, d *schema.ResourceData, config *transport_tpg.Config) interface{} {
	return v
}

func dataSourceGoogleCloudBackupDRServiceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	billingProject := project
	url, err := tpgresource.ReplaceVars(d, config, `{{BackupDRBasePath}}projects/{{project}}/locations/{{location}}/managementServers`)
	if err != nil {
		return err
	}
	if bp, err := tpgresource.GetBillingProject(d, config); err == nil {
		billingProject = bp
	}
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   billingProject,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return fmt.Errorf("Error reading ManagementServer: %s", err)
	}
	managementServersResponse := res["managementServers"]
	resourceResponse := flattenBackupDRManagementServerResourceResp(managementServersResponse, d, config)
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading ManagementServer: %s", err)
	}

	if err := d.Set("type", resourceResponse["type"]); err != nil {
		return fmt.Errorf("Error reading ManagementServer: %s", err)
	}
	if err := d.Set("networks", resourceResponse["networks"]); err != nil {
		return fmt.Errorf("Error reading ManagementServer: %s", err)
	}
	if err := d.Set("oauth2_client_id", resourceResponse["oauth2ClientId"]); err != nil {
		return fmt.Errorf("Error reading ManagementServer: %s", err)
	}
	if err := d.Set("management_uri", resourceResponse["managementUri"]); err != nil {
		return fmt.Errorf("Error reading ManagementServer: %s", err)
	}

	id := fmt.Sprintf("%s", resourceResponse["name"])
	d.SetId(id)
	name := id[strings.LastIndex(id, "/")+1:]
	d.Set("name", name)
	return nil
}
