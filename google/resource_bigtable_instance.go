package google

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"cloud.google.com/go/bigtable"
	"golang.org/x/net/context"
)

func resourceBigtableInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigtableInstanceCreate,
		Read:   resourceBigtableInstanceRead,
		Delete: resourceBigtableInstanceDestroy,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"num_nodes": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},

			"instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "PRODUCTION",
				ValidateFunc: validation.StringInSlice([]string{"DEVELOPMENT", "PRODUCTION"}, false),
			},

			"storage_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "SSD",
				ValidateFunc: validation.StringInSlice([]string{"SSD", "HDD"}, false),
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceBigtableInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	displayName, ok := d.GetOk("display_name")
	if !ok {
		displayName = name
	}

	var storageType bigtable.StorageType
	switch value := d.Get("storage_type"); value {
	case "HDD":
		storageType = bigtable.HDD
	case "SSD":
		storageType = bigtable.SSD
	}

	numNodes := int32(d.Get("num_nodes").(int))
	var instanceType bigtable.InstanceType
	switch value := d.Get("instance_type"); value {
	case "DEVELOPMENT":
		instanceType = bigtable.DEVELOPMENT

		if numNodes > 0 {
			return fmt.Errorf("Can't specify a non-zero number of nodes: %d for DEVELOPMENT Bigtable instance: %s", numNodes, name)
		}
	case "PRODUCTION":
		instanceType = bigtable.PRODUCTION
	}

	instanceConf := &bigtable.InstanceConf{
		InstanceId:   name,
		DisplayName:  displayName.(string),
		ClusterId:    d.Get("cluster_id").(string),
		NumNodes:     numNodes,
		InstanceType: instanceType,
		StorageType:  storageType,
		Zone:         d.Get("zone").(string),
	}

	c, err := config.bigtableClientFactory.NewInstanceAdminClient(project)
	if err != nil {
		return fmt.Errorf("Error starting instance admin client. %s", err)
	}

	defer c.Close()

	err = c.CreateInstance(ctx, instanceConf)
	if err != nil {
		return fmt.Errorf("Error creating instance. %s", err)
	}

	d.SetId(name)

	return resourceBigtableInstanceRead(d, meta)
}

func resourceBigtableInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	c, err := config.bigtableClientFactory.NewInstanceAdminClient(project)
	if err != nil {
		return fmt.Errorf("Error starting instance admin client. %s", err)
	}

	defer c.Close()

	instance, err := c.InstanceInfo(ctx, d.Id())
	if err != nil {
		log.Printf("[WARN] Removing %s because it's gone", d.Id())
		d.SetId("")
		return fmt.Errorf("Error retrieving instance. Could not find %s. %s", d.Id(), err)
	}

	d.Set("name", instance.Name)
	d.Set("display_name", instance.DisplayName)

	return nil
}

func resourceBigtableInstanceDestroy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	c, err := config.bigtableClientFactory.NewInstanceAdminClient(project)
	if err != nil {
		return fmt.Errorf("Error starting instance admin client. %s", err)
	}

	defer c.Close()

	name := d.Id()
	err = c.DeleteInstance(ctx, name)
	if err != nil {
		return fmt.Errorf("Error deleting instance. %s", err)
	}

	d.SetId("")

	return nil
}
