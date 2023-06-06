// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package monitoring

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceMonitoringNotificationChannel() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(ResourceMonitoringNotificationChannel().Schema)

	// Set 'Optional' schema elements
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "display_name")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "project")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "type")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "labels")
	tpgresource.AddOptionalFieldsToSchema(dsSchema, "user_labels")

	return &schema.Resource{
		Read:   dataSourceMonitoringNotificationChannelRead,
		Schema: dsSchema,
	}
}

func dataSourceMonitoringNotificationChannelRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := tpgresource.ReplaceVars(d, config, "{{MonitoringBasePath}}v3/projects/{{project}}/notificationChannels")
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	channelType := d.Get("type").(string)

	if displayName == "" && channelType == "" {
		return fmt.Errorf("At least one of display_name or type must be provided")
	}

	labels, err := expandMonitoringNotificationChannelLabels(d.Get("labels"), d, config)
	if err != nil {
		return err
	}

	userLabels, err := expandMonitoringNotificationChannelLabels(d.Get("user_labels"), d, config)
	if err != nil {
		return err
	}

	filters := make([]string, 0, len(labels)+2)

	if displayName != "" {
		filters = append(filters, fmt.Sprintf(`display_name="%s"`, displayName))
	}

	if channelType != "" {
		filters = append(filters, fmt.Sprintf(`type="%s"`, channelType))
	}

	for k, v := range labels {
		filters = append(filters, fmt.Sprintf(`labels.%s="%s"`, k, v))
	}

	for k, v := range userLabels {
		filters = append(filters, fmt.Sprintf(`user_labels.%s="%s"`, k, v))
	}

	filter := strings.Join(filters, " AND ")
	params := map[string]string{
		"filter": filter,
	}
	url, err = transport_tpg.AddQueryParams(url, params)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	response, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "GET",
		Project:   project,
		RawURL:    url,
		UserAgent: userAgent,
	})
	if err != nil {
		return fmt.Errorf("Error retrieving NotificationChannels: %s", err)
	}

	var channels []interface{}
	if v, ok := response["notificationChannels"]; ok {
		channels = v.([]interface{})
	}
	if len(channels) == 0 {
		return fmt.Errorf("No NotificationChannel found using filter: %s", filter)
	}
	if len(channels) > 1 {
		return fmt.Errorf("Found more than one 1 NotificationChannel matching specified filter: %s", filter)
	}
	res := channels[0].(map[string]interface{})

	name := flattenMonitoringNotificationChannelName(res["name"], d, config).(string)
	if err := d.Set("name", name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	d.SetId(name)

	return resourceMonitoringNotificationChannelRead(d, meta)
}
