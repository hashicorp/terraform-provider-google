package google

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	GCPolicyModeIntersection = "INTERSECTION"
	GCPolicyModeUnion        = "UNION"
)

func resourceBigtableGCPolicyCustomizeDiffFunc(diff TerraformResourceDiff) error {
	count := diff.Get("max_age.#").(int)
	if count < 1 {
		return nil
	}

	oldDays, newDays := diff.GetChange("max_age.0.days")
	oldDuration, newDuration := diff.GetChange("max_age.0.duration")
	log.Printf("days: %v %v", oldDays, newDays)
	log.Printf("duration: %v %v", oldDuration, newDuration)

	if oldDuration == "" && newDuration != "" {
		// flatten the old days and the new duration to duration... if they are
		// equal then do nothing.
		do, err := time.ParseDuration(newDuration.(string))
		if err != nil {
			return err
		}
		dn := time.Hour * 24 * time.Duration(oldDays.(int))
		if do == dn {
			err := diff.Clear("max_age.0.days")
			if err != nil {
				return err
			}
			err = diff.Clear("max_age.0.duration")
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func resourceBigtableGCPolicyCustomizeDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	return resourceBigtableGCPolicyCustomizeDiffFunc(d)
}

func resourceBigtableGCPolicy() *schema.Resource {
	return &schema.Resource{
		Create:        resourceBigtableGCPolicyCreate,
		Read:          resourceBigtableGCPolicyRead,
		Delete:        resourceBigtableGCPolicyDestroy,
		CustomizeDiff: resourceBigtableGCPolicyCustomizeDiff,

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
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"days": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							Deprecated:   "Deprecated in favor of duration",
							Description:  `Number of days before applying GC policy.`,
							ExactlyOneOf: []string{"max_age.0.days", "max_age.0.duration"},
						},
						"duration": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							Description:  `Duration before applying GC policy`,
							ValidateFunc: validateDuration(),
							ExactlyOneOf: []string{"max_age.0.days", "max_age.0.duration"},
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
							ForceNew:    true,
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
		UseJSONNumber: true,
	}
}

func resourceBigtableGCPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	ctx := context.Background()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	instanceName := GetResourceNameFromSelfLink(d.Get("instance_name").(string))
	c, err := config.BigTableClientFactory(userAgent).NewAdminClient(project, instanceName)
	if err != nil {
		return fmt.Errorf("Error starting admin client. %s", err)
	}
	if err := d.Set("instance_name", instanceName); err != nil {
		return fmt.Errorf("Error setting instance_name: %s", err)
	}

	defer c.Close()

	gcPolicy, err := generateBigtableGCPolicy(d)
	if err != nil {
		return err
	}

	tableName := d.Get("table").(string)
	columnFamily := d.Get("column_family").(string)

	err = retryTimeDuration(func() error {
		reqErr := c.SetGCPolicy(ctx, tableName, columnFamily, gcPolicy)
		return reqErr
	}, d.Timeout(schema.TimeoutCreate), isBigTableRetryableError)
	if err != nil {
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
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	ctx := context.Background()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	instanceName := GetResourceNameFromSelfLink(d.Get("instance_name").(string))
	c, err := config.BigTableClientFactory(userAgent).NewAdminClient(project, instanceName)
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

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	return nil
}

func resourceBigtableGCPolicyDestroy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}
	ctx := context.Background()

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	instanceName := GetResourceNameFromSelfLink(d.Get("instance_name").(string))
	c, err := config.BigTableClientFactory(userAgent).NewAdminClient(project, instanceName)
	if err != nil {
		return fmt.Errorf("Error starting admin client. %s", err)
	}

	defer c.Close()

	err = retryTimeDuration(func() error {
		reqErr := c.SetGCPolicy(ctx, d.Get("table").(string), d.Get("column_family").(string), bigtable.NoGcPolicy())
		return reqErr
	}, d.Timeout(schema.TimeoutDelete), isBigTableRetryableError)
	if err != nil {
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
		d, err := getMaxAgeDuration(l[0].(map[string]interface{}))
		if err != nil {
			return nil, err
		}

		policies = append(policies, bigtable.MaxAgePolicy(d))
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

func getMaxAgeDuration(values map[string]interface{}) (time.Duration, error) {
	d := values["duration"].(string)
	if d != "" {
		return time.ParseDuration(d)
	}

	days := values["days"].(int)

	return time.Hour * 24 * time.Duration(days), nil
}
