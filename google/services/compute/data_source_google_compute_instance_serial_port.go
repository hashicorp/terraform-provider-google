// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package compute

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

func DataSourceGoogleComputeInstanceSerialPort() *schema.Resource {
	return &schema.Resource{
		Read: computeInstanceSerialPortRead,
		Schema: map[string]*schema.Schema{
			"port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"instance": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"contents": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func computeInstanceSerialPortRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	zone, err := tpgresource.GetZone(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("zone", zone); err != nil {
		return fmt.Errorf("Error setting zone: %s", err)
	}

	port := int64(d.Get("port").(int))
	output, err := config.NewComputeClient(userAgent).Instances.GetSerialPortOutput(project, zone, d.Get("instance").(string)).Port(port).Do()
	if err != nil {
		return err
	}

	if err := d.Set("contents", output.Contents); err != nil {
		return fmt.Errorf("Error setting contents: %s", err)
	}
	d.SetId(output.SelfLink)
	return nil
}
