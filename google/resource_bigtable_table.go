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
		Update: resourceBigtableTableUpdate,
		Delete: resourceBigtableTableDestroy,

		Importer: &schema.ResourceImporter{
			State: resourceBigtableTableImport,
		},

		// ----------------------------------------------------------------------
		// IMPORTANT: Do not add any additional ForceNew fields to this resource.
		// Destroying/recreating tables can lead to data loss for users.
		// ----------------------------------------------------------------------
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"column_family": {
				Type:     schema.TypeSet,
				Optional: true,
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
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareResourceNames,
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

	instanceName := GetResourceNameFromSelfLink(d.Get("instance_name").(string))
	c, err := config.bigtableClientFactory.NewAdminClient(project, instanceName)
	if err != nil {
		return fmt.Errorf("Error starting admin client. %s", err)
	}
	d.Set("instance_name", instanceName)

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

	id, err := replaceVars(d, config, "projects/{{project}}/instances/{{instance_name}}/tables/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return resourceBigtableTableRead(d, meta)
}

func resourceBigtableTableRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	instanceName := GetResourceNameFromSelfLink(d.Get("instance_name").(string))
	c, err := config.bigtableClientFactory.NewAdminClient(project, instanceName)
	if err != nil {
		return fmt.Errorf("Error starting admin client. %s", err)
	}

	defer c.Close()

	name := d.Get("name").(string)
	table, err := c.TableInfo(ctx, name)
	if err != nil {
		log.Printf("[WARN] Removing %s because it's gone", name)
		d.SetId("")
		return nil
	}

	d.Set("project", project)
	d.Set("column_family", flattenColumnFamily(table.Families))

	return nil
}

func resourceBigtableTableUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	instanceName := GetResourceNameFromSelfLink(d.Get("instance_name").(string))
	c, err := config.bigtableClientFactory.NewAdminClient(project, instanceName)
	if err != nil {
		return fmt.Errorf("Error starting admin client. %s", err)
	}
	defer c.Close()

	o, n := d.GetChange("column_family")
	oSet := o.(*schema.Set)
	nSet := n.(*schema.Set)
	name := d.Get("name").(string)

	// Add column families that are in new but not in old
	for _, new := range nSet.Difference(oSet).List() {
		column := new.(map[string]interface{})

		if v, ok := column["family"]; ok {
			log.Printf("[DEBUG] adding column family %q", v)
			if err := c.CreateColumnFamily(ctx, name, v.(string)); err != nil {
				return fmt.Errorf("Error creating column family %q: %s", v, err)
			}
		}
	}

	// Remove column families that are in old but not in new
	for _, old := range oSet.Difference(nSet).List() {
		column := old.(map[string]interface{})

		if v, ok := column["family"]; ok {
			log.Printf("[DEBUG] removing column family %q", v)
			if err := c.DeleteColumnFamily(ctx, name, v.(string)); err != nil {
				return fmt.Errorf("Error deleting column family %q: %s", v, err)
			}
		}
	}

	return resourceBigtableTableRead(d, meta)
}

func resourceBigtableTableDestroy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	instanceName := GetResourceNameFromSelfLink(d.Get("instance_name").(string))
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

//TODO(rileykarson): Fix the stored import format after rebasing 3.0.0
func resourceBigtableTableImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/instances/(?P<instance_name>[^/]+)/tables/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<instance_name>[^/]+)/(?P<name>[^/]+)",
		"(?P<instance_name>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/instances/{{instance_name}}/tables/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
