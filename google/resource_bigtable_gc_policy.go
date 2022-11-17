package google

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
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
		Create:        resourceBigtableGCPolicyUpsert,
		Read:          resourceBigtableGCPolicyRead,
		Delete:        resourceBigtableGCPolicyDestroy,
		Update:        resourceBigtableGCPolicyUpsert,
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

			"gc_rules": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   `Serialized JSON string for garbage collection policy. Conflicts with "mode", "max_age" and "max_version".`,
				ValidateFunc:  validation.StringIsJSON,
				ConflictsWith: []string{"mode", "max_age", "max_version"},
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
			},
			"mode": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   `If multiple policies are set, you should choose between UNION OR INTERSECTION.`,
				ValidateFunc:  validation.StringInSlice([]string{GCPolicyModeIntersection, GCPolicyModeUnion}, false),
				ConflictsWith: []string{"gc_rules"},
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
				ConflictsWith: []string{"gc_rules"},
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
				ConflictsWith: []string{"gc_rules"},
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},

			"deletion_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `The deletion policy for the GC policy. Setting ABANDON allows the resource
				to be abandoned rather than deleted. This is useful for GC policy as it cannot be deleted
				in a replicated instance. Possible values are: "ABANDON".`,
				ValidateFunc: validation.StringInSlice([]string{"ABANDON", ""}, false),
			},
		},
		UseJSONNumber: true,
	}
}

func resourceBigtableGCPolicyUpsert(d *schema.ResourceData, meta interface{}) error {
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

	retryFunc := func() (interface{}, error) {
		reqErr := c.SetGCPolicy(ctx, tableName, columnFamily, gcPolicy)
		return "", reqErr
	}
	// The default create timeout is 20 minutes.
	timeout := d.Timeout(schema.TimeoutCreate)
	pollInterval := time.Duration(30) * time.Second
	// Mutations to gc policies can only happen one-at-a-time and take some amount of time.
	// Use a fixed polling rate of 30s based on the RetryInfo returned by the server rather than
	// the standard up-to-10s exponential backoff for those operations.
	_, err = retryWithPolling(retryFunc, timeout, pollInterval, isBigTableRetryableError)
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
	columnFamily := d.Get("column_family").(string)
	ti, err := c.TableInfo(ctx, name)
	if err != nil {
		if isNotFoundGrpcError(err) {
			log.Printf("[WARN] Removing the GC policy because the parent table %s is gone", name)
			d.SetId("")
			return nil
		}
		return err
	}

	for _, fi := range ti.FamilyInfos {
		if fi.Name != columnFamily {
			continue
		}

		d.SetId(fi.GCPolicy)

		// No GC Policy.
		if fi.FullGCPolicy.String() == "" {
			return nil
		}

		// Only set `gc_rules`` when the legacy fields are not set. We are not planning to support legacy fields.
		maxAge := d.Get("max_age")
		maxVersion := d.Get("max_version")
		if d.Get("mode") == "" && len(maxAge.([]interface{})) == 0 && len(maxVersion.([]interface{})) == 0 {
			gcRuleString, err := gcPolicyToGCRuleString(fi.FullGCPolicy, true)
			if err != nil {
				return err
			}
			gcRuleJsonString, err := json.Marshal(gcRuleString)
			if err != nil {
				return fmt.Errorf("Error marshaling GC policy to json: %s", err)
			}
			d.Set("gc_rules", string(gcRuleJsonString))
		}
		break
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	return nil
}

// Recursively convert Bigtable GC policy to JSON format in a map.
func gcPolicyToGCRuleString(gc bigtable.GCPolicy, topLevel bool) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	switch bigtable.GetPolicyType(gc) {
	case bigtable.PolicyMaxAge:
		age := gc.(bigtable.MaxAgeGCPolicy).GetDurationString()
		if topLevel {
			rule := make(map[string]interface{})
			rule["max_age"] = age
			rules := []interface{}{}
			rules = append(rules, rule)
			result["rules"] = rules
		} else {
			result["max_age"] = age
		}
		break
	case bigtable.PolicyMaxVersion:
		// bigtable.MaxVersionsGCPolicy is an int.
		// Not sure why max_version is a float64.
		// TODO: Maybe change max_version to an int.
		version := float64(int(gc.(bigtable.MaxVersionsGCPolicy)))
		if topLevel {
			rule := make(map[string]interface{})
			rule["max_version"] = version
			rules := []interface{}{}
			rules = append(rules, rule)
			result["rules"] = rules
		} else {
			result["max_version"] = version
		}
		break
	case bigtable.PolicyUnion:
		result["mode"] = "union"
		rules := []interface{}{}
		for _, c := range gc.(bigtable.UnionGCPolicy).Children {
			gcRuleString, err := gcPolicyToGCRuleString(c, false)
			if err != nil {
				return nil, err
			}
			rules = append(rules, gcRuleString)
		}
		result["rules"] = rules
		break
	case bigtable.PolicyIntersection:
		result["mode"] = "intersection"
		rules := []interface{}{}
		for _, c := range gc.(bigtable.IntersectionGCPolicy).Children {
			gcRuleString, err := gcPolicyToGCRuleString(c, false)
			if err != nil {
				return nil, err
			}
			rules = append(rules, gcRuleString)
		}
		result["rules"] = rules
	default:
		break
	}

	if err := validateNestedPolicy(result, topLevel); err != nil {
		return nil, err
	}

	return result, nil
}

func resourceBigtableGCPolicyDestroy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	if deletionPolicy := d.Get("deletion_policy"); deletionPolicy == "ABANDON" {
		// Allows for the GC policy to be abandoned without deletion to avoid possible
		// deletion failure in a replicated instance.
		log.Printf("[WARN] The GC policy is abandoned")
		return nil
	}

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

	retryFunc := func() (interface{}, error) {
		reqErr := c.SetGCPolicy(ctx, d.Get("table").(string), d.Get("column_family").(string), bigtable.NoGcPolicy())
		return "", reqErr
	}
	// The default delete timeout is 20 minutes.
	timeout := d.Timeout(schema.TimeoutDelete)
	pollInterval := time.Duration(30) * time.Second
	_, err = retryWithPolling(retryFunc, timeout, pollInterval, isBigTableRetryableError)
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
	gcRules, gok := d.GetOk("gc_rules")

	if !aok && !vok && !gok {
		return bigtable.NoGcPolicy(), nil
	}

	if mode == "" && aok && vok {
		return nil, fmt.Errorf("if multiple policies are set, mode can't be empty")
	}

	if gok {
		var topLevelPolicy map[string]interface{}
		if err := json.Unmarshal([]byte(gcRules.(string)), &topLevelPolicy); err != nil {
			return nil, err
		}
		return getGCPolicyFromJSON(topLevelPolicy /*isTopLevel=*/, true)
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

func getGCPolicyFromJSON(inputPolicy map[string]interface{}, isTopLevel bool) (bigtable.GCPolicy, error) {
	policy := []bigtable.GCPolicy{}

	if err := validateNestedPolicy(inputPolicy, isTopLevel); err != nil {
		return nil, err
	}

	for _, p := range inputPolicy["rules"].([]interface{}) {
		childPolicy := p.(map[string]interface{})
		if err := validateNestedPolicy(childPolicy /*isTopLevel=*/, false); err != nil {
			return nil, err
		}

		if childPolicy["max_age"] != nil {
			maxAge := childPolicy["max_age"].(string)
			duration, err := time.ParseDuration(maxAge)
			if err != nil {
				return nil, fmt.Errorf("invalid duration string: %v", maxAge)
			}
			policy = append(policy, bigtable.MaxAgePolicy(duration))
		}

		if childPolicy["max_version"] != nil {
			version := childPolicy["max_version"].(float64)
			policy = append(policy, bigtable.MaxVersionsPolicy(int(version)))
		}

		if childPolicy["mode"] != nil {
			n, err := getGCPolicyFromJSON(childPolicy /*isTopLevel=*/, false)
			if err != nil {
				return nil, err
			}
			policy = append(policy, n)
		}
	}

	switch inputPolicy["mode"] {
	case strings.ToLower(GCPolicyModeUnion):
		return bigtable.UnionPolicy(policy...), nil
	case strings.ToLower(GCPolicyModeIntersection):
		return bigtable.IntersectionPolicy(policy...), nil
	default:
		return policy[0], nil
	}
}

func validateNestedPolicy(p map[string]interface{}, isTopLevel bool) error {
	if len(p) > 2 {
		return fmt.Errorf("rules has more than 2 fields")
	}
	maxVersion, maxVersionOk := p["max_version"]
	maxAge, maxAgeOk := p["max_age"]
	rulesObj, rulesOk := p["rules"]

	_, modeOk := p["mode"]
	rules, arrOk := rulesObj.([]interface{})
	_, vCastOk := maxVersion.(float64)
	_, aCastOk := maxAge.(string)

	if rulesOk && !arrOk {
		return fmt.Errorf("`rules` must be array")
	}

	if modeOk && len(rules) < 2 {
		return fmt.Errorf("`rules` need at least 2 GC rule when mode is specified")
	}

	if isTopLevel && !rulesOk {
		return fmt.Errorf("invalid nested policy, need `rules`")
	}

	if isTopLevel && !modeOk && len(rules) != 1 {
		return fmt.Errorf("when `mode` is not specified, `rules` can only have 1 child rule")
	}

	if !isTopLevel && len(p) == 2 && (!modeOk || !rulesOk) {
		return fmt.Errorf("need `mode` and `rules` for child nested policies")
	}

	if !isTopLevel && len(p) == 1 && !maxVersionOk && !maxAgeOk {
		return fmt.Errorf("need `max_version` or `max_age` for the rule")
	}

	if maxVersionOk && !vCastOk {
		return fmt.Errorf("`max_version` must be a number")
	}

	if maxAgeOk && !aCastOk {
		return fmt.Errorf("`max_age must be a string")
	}

	return nil
}

func getMaxAgeDuration(values map[string]interface{}) (time.Duration, error) {
	d := values["duration"].(string)
	if d != "" {
		return time.ParseDuration(d)
	}

	days := values["days"].(int)

	return time.Hour * 24 * time.Duration(days), nil
}
