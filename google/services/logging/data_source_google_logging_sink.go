// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleLoggingSink() *schema.Resource {
	dsSchema := tpgresource.DatasourceSchemaFromResourceSchema(resourceLoggingSinkSchema())
	dsSchema["id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: `Required. An identifier for the resource in format: "projects/[PROJECT_ID]/sinks/[SINK_NAME]", "organizations/[ORGANIZATION_ID]/sinks/[SINK_NAME]", "billingAccounts/[BILLING_ACCOUNT_ID]/sinks/[SINK_NAME]", "folders/[FOLDER_ID]/sinks/[SINK_NAME]"`,
	}

	return &schema.Resource{
		Read:   dataSourceGoogleLoggingSinkRead,
		Schema: dsSchema,
	}
}

func dataSourceGoogleLoggingSinkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	sinkId := d.Get("id").(string)

	sink, err := config.NewLoggingClient(userAgent).Sinks.Get(sinkId).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Logging Sink %s", d.Id()))
	}

	if err := flattenResourceLoggingSink(d, sink); err != nil {
		return err
	}

	d.SetId(sinkId)

	return nil
}
