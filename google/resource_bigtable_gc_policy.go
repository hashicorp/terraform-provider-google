package google

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	GCPolicyModeIntersection = "INTERSECTION"
	GCPolicyModeUnion        = "UNION"
)

func resourceBigtableGCPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigtableGCPolicyCreate,
		Read:   resourceBigtableGCPolicyRead,
		Delete: resourceBigtableGCPolicyDestroy,

		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareResourceNames,
				Description:      `The name of the Bigtable instance.`,
			},

			"table": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the table.`,
			},

			"column_family": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the column family.`,
			},

			"mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  `If multiple policies are set, you should choose between UNION OR INTERSECTION.`,
				ValidateFunc: validation.StringInSlice([]string{GCPolicyModeIntersection, GCPolicyModeUnion}, false),
			},

			"max_age": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `GC policy that applies to all cells older than the given age.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"days": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `Number of days before applying GC policy.`,
						},
					},
				},
			},

			"max_version": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: `GC policy that applies to all versions of a cell except for the most recent.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"number": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `Number of version before applying the GC policy.`,
						},
					},
				},
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},
		},
	}
}

func resourceBigtableGCPolicyCreate(d *schema.ResourceData, meta interface{}) error {
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

	gcPolicy, err := generateBigtableGCPolicy(d)
	if err != nil {
		return err
	}

	tableName := d.Get("table").(string)
	columnFamily := d.Get("column_family").(string)

	if err := c.SetGCPolicy(ctx, tableName, columnFamily, gcPolicy); err != nil {
		return err
	}

	table, err := c.TableInfo(ctx, tableName)
	if err != nil {
		return fmt.Errorf("Error retrieving table. Could not find %s in %s. %s", tableName, instanceName, err)
	}

	for _, i := range table.FamilyInfos {
		if i.Name == columnFamily {
			d.SetId(i.GCPolicy)
		}
	}

	return resourceBigtableGCPolicyRead(d, meta)
}

func resourceBigtableGCPolicyRead(d *schema.ResourceData, meta interface{}) error {
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

	name := d.Get("table").(string)
	ti, err := c.TableInfo(ctx, name)
	if err != nil {
		log.Printf("[WARN] Removing %s because it's gone", name)
		d.SetId("")
		return nil
	}

	for _, fi := range ti.FamilyInfos {
		if fi.Name == name {
			d.SetId(fi.GCPolicy)
			break
		}
	}

	d.Set("project", project)

	return nil
}

func resourceBigtableGCPolicyDestroy(d *schema.ResourceData, meta interface{}) error {
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

	if err := c.SetGCPolicy(ctx, d.Get("table").(string), d.Get("column_family").(string), bigtable.NoGcPolicy()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func generateBigtableGCPolicy(d *schema.ResourceData) (bigtable.GCPolicy, error) {
	var policies []bigtable.GCPolicy
	mode := d.Get("mode").(string)
	ma, aok := d.GetOk("max_age")
	mv, vok := d.GetOk("max_version")

	if !aok && !vok {
		return bigtable.NoGcPolicy(), nil
	}

	if mode == "" && aok && vok {
		return nil, fmt.Errorf("If multiple policies are set, mode can't be empty")
	}

	if aok {
		l, _ := ma.([]interface{})
		d, _ := l[0].(map[string]interface{})["days"].(int)

		policies = append(policies, bigtable.MaxAgePolicy(time.Duration(d)*time.Hour*24))
	}

	if vok {
		l, _ := mv.([]interface{})
		n, _ := l[0].(map[string]interface{})["number"].(int)

		policies = append(policies, bigtable.MaxVersionsPolicy(n))
	}

	switch mode {
	case GCPolicyModeUnion:
		return bigtable.UnionPolicy(policies...), nil
	case GCPolicyModeIntersection:
		return bigtable.IntersectionPolicy(policies...), nil
	}

	return policies[0], nil
}
