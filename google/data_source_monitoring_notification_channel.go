package google

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceMonitoringNotificationChannel() *schema.Resource {
	dsSchema := datasourceSchemaFromResourceSchema(resourceMonitoringNotificationChannel().Schema)

	// Set 'Optional' schema elements
	addOptionalFieldsToSchema(dsSchema, "display_name")
	addOptionalFieldsToSchema(dsSchema, "project")
	addOptionalFieldsToSchema(dsSchema, "type")

	return &schema.Resource{
		Read:   dataSourceMonitoringNotificationChannelRead,
		Schema: dsSchema,
	}
}

func dataSourceMonitoringNotificationChannelRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "{{MonitoringBasePath}}projects/{{project}}/notificationChannels")
	if err != nil {
		return err
	}

	displayName := d.Get("display_name").(string)
	channelType := d.Get("type").(string)

	if displayName == "" && channelType == "" {
		return errors.New("Must at least provide either `display_name` or `type`")
	}

	filter := ""
	if displayName != "" {
		filter = fmt.Sprintf("display_name=\"%s\"", displayName)
	}

	if channelType != "" {
		channelFilter := fmt.Sprintf("type=\"%s\"", channelType)
		if filter != "" {
			filter += fmt.Sprintf(" AND %s", channelFilter)
		} else {
			filter = channelFilter
		}
	}

	params := make(map[string]string)
	params["filter"] = filter

	url, err = addQueryParams(url, params)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	response, err := sendRequest(config, "GET", project, url, nil)
	if err != nil {
		return fmt.Errorf("Error retrieving NotificationChannels: %s", err)
	}

	var pageMonitoringNotificationChannels []interface{}
	if v, ok := response["notificationChannels"]; ok {
		pageMonitoringNotificationChannels = v.([]interface{})
	}

	if len(pageMonitoringNotificationChannels) == 0 {
		return fmt.Errorf("No NotificationChannel found using filter=%s", filter)
	}

	if len(pageMonitoringNotificationChannels) > 1 {
		return fmt.Errorf("More than one matching NotificationChannel found using filter=%s", filter)
	}

	res := pageMonitoringNotificationChannels[0].(map[string]interface{})

	name := flattenMonitoringNotificationChannelName(res["name"], d).(string)
	d.Set("name", name)
	d.Set("project", project)
	d.Set("labels", flattenMonitoringNotificationChannelLabels(res["labels"], d))
	d.Set("verification_status", flattenMonitoringNotificationChannelVerificationStatus(res["verificationStatus"], d))
	d.Set("type", flattenMonitoringNotificationChannelType(res["type"], d))
	d.Set("user_labels", flattenMonitoringNotificationChannelUserLabels(res["userLabels"], d))
	d.Set("description", flattenMonitoringNotificationChannelDescription(res["descriptionx"], d))
	d.Set("display_name", flattenMonitoringNotificationChannelDisplayName(res["displayName"], d))
	d.Set("enabled", flattenMonitoringNotificationChannelEnabled(res["enabled"], d))
	d.SetId(name)

	return nil
}
