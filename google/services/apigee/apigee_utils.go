// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package apigee

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func resourceApigeeNatAddressActivate(config *transport_tpg.Config, d *schema.ResourceData, billingProject string, userAgent string) error {
	// 1. check prepare for activation
	name := d.Get("name").(string)

	if d.Get("state").(string) != "RESERVED" {
		return fmt.Errorf("Activating NAT address requires the state to become RESERVED")
	}

	// 2. activation
	activateUrl, err := tpgresource.ReplaceVars(d, config, "{{ApigeeBasePath}}{{instance_id}}/natAddresses/{{name}}:activate")
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Activating NAT address: %s", name)
	res, err := transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
		Config:    config,
		Method:    "POST",
		Project:   billingProject,
		RawURL:    activateUrl,
		UserAgent: userAgent,
	})
	if err != nil {
		return fmt.Errorf("Error activating NAT address: %s", err)
	}

	var opRes map[string]interface{}
	err = ApigeeOperationWaitTimeWithResponse(
		config, res, &opRes, "Activating NAT address", userAgent,
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error waiting to actiavte NAT address: %s", err)
	} else {
		log.Printf("[DEBUG] Finished activating NatAddress %q: %#v", d.Id(), res)
	}
	return nil
}
