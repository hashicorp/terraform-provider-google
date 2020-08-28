package google

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGoogleComputeInstanceSerialPort() *schema.Resource {
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
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading project: %s", err)
	}
	zone, err := getZone(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("zone", zone); err != nil {
		return fmt.Errorf("Error reading zone: %s", err)
	}

	port := int64(d.Get("port").(int))
	output, err := config.clientCompute.Instances.GetSerialPortOutput(project, zone, d.Get("instance").(string)).Port(port).Do()
	if err != nil {
		return err
	}

	if err := d.Set("contents", output.Contents); err != nil {
		return fmt.Errorf("Error reading contents: %s", err)
	}
	d.SetId(output.SelfLink)
	return nil
}
