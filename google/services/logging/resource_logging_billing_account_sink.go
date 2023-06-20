// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package logging

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func ResourceLoggingBillingAccountSink() *schema.Resource {
	schm := &schema.Resource{
		Create: resourceLoggingBillingAccountSinkCreate,
		Read:   resourceLoggingBillingAccountSinkRead,
		Delete: resourceLoggingBillingAccountSinkDelete,
		Update: resourceLoggingBillingAccountSinkUpdate,
		Schema: resourceLoggingSinkSchema(),
		Importer: &schema.ResourceImporter{
			State: resourceLoggingSinkImportState("billing_account"),
		},
		UseJSONNumber: true,
	}
	schm.Schema["billing_account"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: `The billing account exported to the sink.`,
	}
	return schm
}

func resourceLoggingBillingAccountSinkCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	id, sink := expandResourceLoggingSink(d, "billingAccounts", d.Get("billing_account").(string))

	// The API will reject any requests that don't explicitly set 'uniqueWriterIdentity' to true.
	_, err = config.NewLoggingClient(userAgent).BillingAccounts.Sinks.Create(id.parent(), sink).UniqueWriterIdentity(true).Do()
	if err != nil {
		return err
	}

	d.SetId(id.canonicalId())
	return resourceLoggingBillingAccountSinkRead(d, meta)
}

func resourceLoggingBillingAccountSinkRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	sink, err := config.NewLoggingClient(userAgent).BillingAccounts.Sinks.Get(d.Id()).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Billing Logging Sink %s", d.Get("name").(string)))
	}

	if err := flattenResourceLoggingSink(d, sink); err != nil {
		return err
	}

	return nil
}

func resourceLoggingBillingAccountSinkUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	sink, updateMask := expandResourceLoggingSinkForUpdate(d)

	// The API will reject any requests that don't explicitly set 'uniqueWriterIdentity' to true.
	_, err = config.NewLoggingClient(userAgent).BillingAccounts.Sinks.Patch(d.Id(), sink).
		UpdateMask(updateMask).UniqueWriterIdentity(true).Do()
	if err != nil {
		return err
	}

	return resourceLoggingBillingAccountSinkRead(d, meta)
}

func resourceLoggingBillingAccountSinkDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	_, err = config.NewLoggingClient(userAgent).Projects.Sinks.Delete(d.Id()).Do()
	if err != nil {
		return err
	}

	return nil
}
