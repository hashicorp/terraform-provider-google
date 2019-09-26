package google

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBigtableTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigtableTableCreate,
		Read:   resourceBigtableTableRead,
		Delete: resourceBigtableTableDestroy,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"column_family": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"family": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"split_keys": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceBigtableTableCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	instanceName := d.Get("instance_name").(string)
	c, err := config.bigtableClientFactory.NewAdminClient(project, instanceName)
	if err != nil {
		return fmt.Errorf("Error starting admin client. %s", err)
	}

	defer c.Close()

	name := d.Get("name").(string)
	if v, ok := d.GetOk("split_keys"); ok {
		splitKeys := convertStringArr(v.([]interface{}))
		// This method may return before the table's creation is complete - we may need to wait until
		// it exists in the future.
		err = c.CreatePresplitTable(ctx, name, splitKeys)
		if err != nil {
			return fmt.Errorf("Error creating presplit table. %s", err)
		}
	} else {
		// This method may return before the table's creation is complete - we may need to wait until
		// it exists in the future.
		err = c.CreateTable(ctx, name)
		if err != nil {
			return fmt.Errorf("Error creating table. %s", err)
		}
	}

	if d.Get("column_family.#").(int) > 0 {
		columns := d.Get("column_family").(*schema.Set).List()

		for _, co := range columns {
			column := co.(map[string]interface{})

			if v, ok := column["family"]; ok {
				if err := c.CreateColumnFamily(ctx, name, v.(string)); err != nil {
					return fmt.Errorf("Error creating column family %s. %s", v, err)
				}
			}
		}
	}

	d.SetId(name)

	return resourceBigtableTableRead(d, meta)
}

func resourceBigtableTableRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	instanceName := d.Get("instance_name").(string)
	c, err := config.bigtableClientFactory.NewAdminClient(project, instanceName)
	if err != nil {
		return fmt.Errorf("Error starting admin client. %s", err)
	}

	defer c.Close()

	name := d.Id()
	table, err := c.TableInfo(ctx, name)
	if err != nil {
		log.Printf("[WARN] Removing %s because it's gone", name)
		d.SetId("")
		return fmt.Errorf("Error retrieving table. Could not find %s in %s. %s", name, instanceName, err)
	}

	d.Set("project", project)
	d.Set("column_family", flattenColumnFamily(table.Families))

	return nil
}

func resourceBigtableTableDestroy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	instanceName := d.Get("instance_name").(string)
	c, err := config.bigtableClientFactory.NewAdminClient(project, instanceName)
	if err != nil {
		return fmt.Errorf("Error starting admin client. %s", err)
	}

	defer c.Close()

	name := d.Get("name").(string)
	err = c.DeleteTable(ctx, name)
	if err != nil {
		return fmt.Errorf("Error deleting table. %s", err)
	}

	d.SetId("")

	return nil
}

func flattenColumnFamily(families []string) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(families))

	for _, f := range families {
		data := make(map[string]interface{})
		data["family"] = f
		result = append(result, data)
	}

	return result
}
