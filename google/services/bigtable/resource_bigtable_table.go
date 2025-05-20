// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package bigtable

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func familyHash(v interface{}) int {
	m := v.(map[string]interface{})
	cf := m["family"].(string)
	t, err := getType(m["type"])
	if err != nil {
		panic(err)
	}
	if t == nil {
		// no specified type.
		return tpgresource.Hashcode(cf)
	}
	b, err := bigtable.MarshalJSON(t)
	if err != nil {
		panic(err)
	}
	return tpgresource.Hashcode(cf + string(b))
}

func ResourceBigtableTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceBigtableTableCreate,
		Read:   resourceBigtableTableRead,
		Update: resourceBigtableTableUpdate,
		Delete: resourceBigtableTableDestroy,

		Importer: &schema.ResourceImporter{
			State: resourceBigtableTableImport,
		},

		// Set a longer timeout for table creation as adding column families can be slow.
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			tpgresource.DefaultProviderProject,
			abpDiffFunc,
		),
		// ----------------------------------------------------------------------
		// IMPORTANT: Do not add any additional ForceNew fields to this resource.
		// Destroying/recreating tables can lead to data loss for users.
		// ----------------------------------------------------------------------
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the table. Must be 1-50 characters and must only contain hyphens, underscores, periods, letters and numbers.`,
			},

			"column_family": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `A group of columns within a table which share a common configuration. This can be specified multiple times.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"family": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the column family.`,
						},
						"type": {
							Type:             schema.TypeString,
							Optional:         true,
							Description:      `The type of the column family.`,
							DiffSuppressFunc: typeDiffFunc,
						},
					},
				},
				Set: familyHash,
			},

			"instance_name": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: tpgresource.CompareResourceNames,
				Description:      `The name of the Bigtable instance.`,
			},

			"split_keys": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `A list of predefined keys to split the table on. !> Warning: Modifying the split_keys of an existing table will cause Terraform to delete/recreate the entire google_bigtable_table resource.`,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},

			"deletion_protection": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"PROTECTED", "UNPROTECTED"}, false),
				Elem:         &schema.Schema{Type: schema.TypeString},
				Description:  `A field to make the table protected against data loss i.e. when set to PROTECTED, deleting the table, the column families in the table, and the instance containing the table would be prohibited. If not provided, currently deletion protection will be set to UNPROTECTED as it is the API default value.`,
			},

			"change_stream_retention": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: verify.ValidateDuration(),
				Description:  `Duration to retain change stream data for the table. Set to 0 to disable. Must be between 1 and 7 days.`,
			},

			"automated_backup_policy": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"retention_period": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: verify.ValidateDuration(),
							Description:  `How long the automated backups should be retained.`,
						},
						"frequency": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: verify.ValidateDuration(),
							Description:  `How frequently automated backups should occur.`,
						},
					},
				},
				Description: `Defines an automated backup policy for a table, specified by Retention Period and Frequency. To _create_ a table with automated backup disabled, omit this argument. To disable automated backup on an _existing_ table that has automated backup enabled, set both Retention Period and Frequency to "0". If this argument is not provided in the configuration on update, the resource's automated backup policy will _not_ be modified.`,
			},
		},
		UseJSONNumber: true,
	}
}

func typeDiffFunc(k, oldValue, newValue string, d *schema.ResourceData) bool {
	old, err := getType(oldValue)
	if err != nil {
		panic(fmt.Sprintf("old error: %v", err))
	}
	new, err := getType(newValue)
	if err != nil {
		panic(fmt.Sprintf("new error: %v", err))
	}
	return bigtable.Equal(old, new)
}

// The API uses nil to indicate disabled automated backups. Terraform uses nil
// during creation, or a zero-retention/frequency policy during updates, for the
// same purpose. This diff function recognizes these representations as equivalent.
func abpDiffFunc(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	if !diff.HasChange("automated_backup_policy") {
		return nil
	}

	old, new := diff.GetChange("automated_backup_policy")
	oldAbpSet, ok := old.(*schema.Set)
	if !ok {
		fmt.Errorf("error parsing old automated backup policy: %v", old)
	}
	newAbpSet, ok := new.(*schema.Set)
	if !ok {
		fmt.Errorf("error parsing new automated backup policy: %v", new)
	}

	// If the state contains nil automated_backup_policy and configuration contains
	// automated_backup_policy with zeros, the two are equivalent
	if oldAbpSet.Len() == 0 && newAbpSet.Len() == 1 {
		newAbpMap := newAbpSet.List()[0].(map[string]interface{})
		if newAbpMap["retention_period"] == "0" && newAbpMap["frequency"] == "0" {
			log.Printf("[DEBUG] Suppressing diff for zero automated backup policy")
			diff.Clear("automated_backup_policy")
		}
	}

	return nil
}

func resourceBigtableTableCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	ctx := context.Background()

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	instanceName := tpgresource.GetResourceNameFromSelfLink(d.Get("instance_name").(string))
	c, err := config.BigTableClientFactory(userAgent).NewAdminClient(project, instanceName)
	if err != nil {
		return fmt.Errorf("Error starting admin client. %s", err)
	}
	if err := d.Set("instance_name", instanceName); err != nil {
		return fmt.Errorf("Error setting instance_name: %s", err)
	}

	defer c.Close()

	tableId := d.Get("name").(string)
	tblConf := bigtable.TableConf{TableID: tableId}

	// Check if deletion protection is given
	// If not given, currently tblConf.DeletionProtection will be set to false in the API
	deletionProtection := d.Get("deletion_protection")
	if deletionProtection == "PROTECTED" {
		tblConf.DeletionProtection = bigtable.Protected
	} else if deletionProtection == "UNPROTECTED" {
		tblConf.DeletionProtection = bigtable.Unprotected
	}

	if changeStreamRetention, ok := d.GetOk("change_stream_retention"); ok {
		tblConf.ChangeStreamRetention, err = time.ParseDuration(changeStreamRetention.(string))
		if err != nil {
			return fmt.Errorf("Error parsing change stream retention: %s", err)
		}
	}

	if automatedBackupPolicyField, ok := d.GetOk("automated_backup_policy"); ok {
		automatedBackupPolicyElements := automatedBackupPolicyField.(*schema.Set).List()
		if len(automatedBackupPolicyElements) == 0 {
			return fmt.Errorf("Incomplete automated_backup_policy")
		} else {
			automatedBackupPolicy := automatedBackupPolicyElements[0].(map[string]interface{})
			abpRetentionPeriodField, retentionPeriodExists := automatedBackupPolicy["retention_period"]
			if !retentionPeriodExists {
				return fmt.Errorf("Automated backup policy retention period must be specified")
			}
			abpFrequencyField, frequencyExists := automatedBackupPolicy["frequency"]
			if !frequencyExists {
				return fmt.Errorf("Automated backup policy frequency must be specified")
			}
			abpRetentionPeriod, err := ParseDuration(abpRetentionPeriodField.(string))
			if err != nil {
				return fmt.Errorf("Error parsing automated backup policy retention period: %s", err)
			}
			abpFrequency, err := ParseDuration(abpFrequencyField.(string))
			if err != nil {
				return fmt.Errorf("Error parsing automated backup policy frequency: %s", err)
			}
			tblConf.AutomatedBackupConfig = &bigtable.TableAutomatedBackupPolicy{
				RetentionPeriod: abpRetentionPeriod,
				Frequency:       abpFrequency,
			}
		}
	}

	// Set the split keys if given.
	if v, ok := d.GetOk("split_keys"); ok {
		tblConf.SplitKeys = tpgresource.ConvertStringArr(v.([]interface{}))
	}

	// Set the column families if given.
	columnFamilies := make(map[string]bigtable.Family)
	if d.Get("column_family.#").(int) > 0 {
		columns := d.Get("column_family").(*schema.Set).List()

		for _, co := range columns {
			column := co.(map[string]interface{})

			if v, ok := column["family"]; ok {
				valueType, err := getType(column["type"])
				if err != nil {
					return err
				}
				columnFamilies[v.(string)] = bigtable.Family{
					// By default, there is no GC rules.
					GCPolicy:  bigtable.NoGcPolicy(),
					ValueType: valueType,
				}
			}
		}
	}
	tblConf.ColumnFamilies = columnFamilies

	// This method may return before the table's creation is complete - we may need to wait until
	// it exists in the future.
	// Set a longer timeout as creating table and adding column families can be pretty slow.
	ctxWithTimeout, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel() // Always call cancel.
	err = c.CreateTableFromConf(ctxWithTimeout, &tblConf)
	if err != nil {
		return fmt.Errorf("Error creating table. %s", err)
	}

	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/instances/{{instance_name}}/tables/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return resourceBigtableTableRead(d, meta)
}

func resourceBigtableTableRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	ctx := context.Background()

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	instanceName := tpgresource.GetResourceNameFromSelfLink(d.Get("instance_name").(string))
	c, err := config.BigTableClientFactory(userAgent).NewAdminClient(project, instanceName)
	if err != nil {
		return fmt.Errorf("Error starting admin client. %s", err)
	}

	defer c.Close()

	name := d.Get("name").(string)
	table, err := c.TableInfo(ctx, name)
	if err != nil {
		if tpgresource.IsNotFoundGrpcError(err) {
			log.Printf("[WARN] Removing %s because it's gone", name)
			d.SetId("")
			return nil
		}
		return err
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	families, err := FlattenColumnFamily(table.FamilyInfos)
	if err != nil {
		return fmt.Errorf("Error flatenning column families: %v", err)
	}
	if err := d.Set("column_family", families); err != nil {
		return fmt.Errorf("Error setting column_family: %s", err)
	}

	deletionProtection := table.DeletionProtection
	if deletionProtection == bigtable.Protected {
		if err := d.Set("deletion_protection", "PROTECTED"); err != nil {
			return fmt.Errorf("Error setting deletion_protection: %s", err)
		}
	} else if deletionProtection == bigtable.Unprotected {
		if err := d.Set("deletion_protection", "UNPROTECTED"); err != nil {
			return fmt.Errorf("Error setting deletion_protection: %s", err)
		}
	} else {
		return fmt.Errorf("Error setting deletion_protection, it should be either PROTECTED or UNPROTECTED")
	}

	changeStreamRetention := table.ChangeStreamRetention
	if changeStreamRetention != nil {
		if err := d.Set("change_stream_retention", changeStreamRetention.(time.Duration).String()); err != nil {
			return fmt.Errorf("Error setting change_stream_retention: %s", err)
		}
	}

	if table.AutomatedBackupConfig != nil {
		switch automatedBackupConfig := table.AutomatedBackupConfig.(type) {
		case *bigtable.TableAutomatedBackupPolicy:
			var tableAbp bigtable.TableAutomatedBackupPolicy = *automatedBackupConfig
			abpRetentionPeriod := tableAbp.RetentionPeriod.(time.Duration).String()
			abpFrequency := tableAbp.Frequency.(time.Duration).String()
			abp := []interface{}{
				map[string]interface{}{
					"retention_period": abpRetentionPeriod,
					"frequency":        abpFrequency,
				},
			}
			if err := d.Set("automated_backup_policy", abp); err != nil {
				return fmt.Errorf("Error setting automated_backup_policy: %s", err)
			}
		default:
			return fmt.Errorf("error: Unknown type of automated backup configuration")
		}
	}

	return nil
}

func toFamilyMap(set *schema.Set) (map[string]bigtable.Family, error) {
	result := map[string]bigtable.Family{}
	for _, item := range set.List() {
		column := item.(map[string]interface{})

		if v, ok := column["family"]; ok && v != "" {
			valueType, err := getType(column["type"])
			if err != nil {
				return nil, err
			}
			result[v.(string)] = bigtable.Family{
				ValueType: valueType,
			}
		}
	}
	return result, nil
}

// familyMapDiffKeys returns a new map that is the result of a-b, comparing keys
func familyMapDiffKeys(a, b map[string]bigtable.Family) map[string]bigtable.Family {
	result := map[string]bigtable.Family{}
	for k, v := range a {
		if _, ok := b[k]; !ok {
			result[k] = v
		}
	}
	return result
}

// familyMapDiffValueTypes returns a new map that is the result of a-b, where a and b share keys but have different value types
func familyMapDiffValueTypes(a, b map[string]bigtable.Family) map[string]bigtable.Family {
	result := map[string]bigtable.Family{}
	for k, va := range a {
		if vb, ok := b[k]; ok {
			if !bigtable.Equal(va.ValueType, vb.ValueType) {
				result[k] = va
			}
		}
	}
	return result
}

func resourceBigtableTableUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}
	ctx := context.Background()

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	instanceName := tpgresource.GetResourceNameFromSelfLink(d.Get("instance_name").(string))
	c, err := config.BigTableClientFactory(userAgent).NewAdminClient(project, instanceName)
	if err != nil {
		return fmt.Errorf("Error starting admin client. %s", err)
	}
	defer c.Close()

	o, n := d.GetChange("column_family")
	name := d.Get("name").(string)

	oMap, err := toFamilyMap(o.(*schema.Set))
	if err != nil {
		return err
	}
	nMap, err := toFamilyMap(n.(*schema.Set))
	if err != nil {
		return err
	}

	for cfn, cf := range familyMapDiffKeys(nMap, oMap) {
		log.Printf("[DEBUG] adding column family %q", cfn)
		if err := c.CreateColumnFamilyWithConfig(ctx, name, cfn, cf); err != nil {
			return fmt.Errorf("Error creating column family %q: %s", cfn, err)
		}
	}
	for cfn := range familyMapDiffKeys(oMap, nMap) {
		log.Printf("[DEBUG] removing column family %q", cfn)
		if err := c.DeleteColumnFamily(ctx, name, cfn); err != nil {
			return fmt.Errorf("Error deleting column family %q: %s", cfn, err)
		}
	}
	for cfn, cf := range familyMapDiffValueTypes(nMap, oMap) {
		log.Printf("[DEBUG] updating column family: %q", cfn)
		if err := c.UpdateFamily(ctx, name, cfn, cf); err != nil {
			return fmt.Errorf("Error update column family %q: %s", cfn, err)
		}
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()
	if d.HasChange("deletion_protection") {
		deletionProtection := d.Get("deletion_protection")
		if deletionProtection == "PROTECTED" {
			if err := c.UpdateTableWithDeletionProtection(ctxWithTimeout, name, bigtable.Protected); err != nil {
				return fmt.Errorf("Error updating deletion protection in table %v: %s", name, err)
			}
		} else if deletionProtection == "UNPROTECTED" {
			if err := c.UpdateTableWithDeletionProtection(ctxWithTimeout, name, bigtable.Unprotected); err != nil {
				return fmt.Errorf("Error updating deletion protection in table %v: %s", name, err)
			}
		}
	}

	if d.HasChange("change_stream_retention") {
		changeStreamRetention := d.Get("change_stream_retention")
		changeStream, err := time.ParseDuration(changeStreamRetention.(string))
		if err != nil {
			return fmt.Errorf("Error parsing change stream retention: %s", err)
		}
		if changeStream == 0 {
			if err := c.UpdateTableDisableChangeStream(ctxWithTimeout, name); err != nil {
				return fmt.Errorf("Error disabling change stream retention in table %v: %s", name, err)
			}
		} else {
			if err := c.UpdateTableWithChangeStream(ctxWithTimeout, name, changeStream); err != nil {
				return fmt.Errorf("Error updating change stream retention in table %v: %s", name, err)
			}
		}
	}

	if d.HasChange("automated_backup_policy") {
		automatedBackupPolicyField := d.Get("automated_backup_policy").(*schema.Set)
		automatedBackupPolicyElements := automatedBackupPolicyField.List()
		// If the automated_backup_policy field is being removed, do not modify the automated_backup_policy on the resource when applying changes
		if len(automatedBackupPolicyElements) == 0 {
			log.Printf("[DEBUG] automated_backup_policy field removed from configuration, will not modify automated_backup_policy on resource")
		} else {
			automatedBackupPolicy := automatedBackupPolicyElements[0].(map[string]interface{})
			abp := bigtable.TableAutomatedBackupPolicy{}

			abpRetentionPeriodField, retentionPeriodExists := automatedBackupPolicy["retention_period"]
			if retentionPeriodExists && abpRetentionPeriodField != "" {
				abpRetentionPeriod, err := ParseDuration(abpRetentionPeriodField.(string))
				if err != nil {
					return fmt.Errorf("Error parsing automated backup policy retention period: %s", err)
				}
				abp.RetentionPeriod = abpRetentionPeriod
			}

			abpFrequencyField, frequencyExists := automatedBackupPolicy["frequency"]
			if frequencyExists && abpFrequencyField != "" {
				abpFrequency, err := ParseDuration(abpFrequencyField.(string))
				if err != nil {
					return fmt.Errorf("Error parsing automated backup policy frequency: %s", err)
				}
				abp.Frequency = abpFrequency
			}

			if abp.RetentionPeriod != nil && abp.RetentionPeriod.(time.Duration) == 0 && abp.Frequency != nil && abp.Frequency.(time.Duration) == 0 {
				// Disable Automated Backups
				if err := c.UpdateTableDisableAutomatedBackupPolicy(ctxWithTimeout, name); err != nil {
					return fmt.Errorf("Error disabling automated backup configuration on table %v: %s", name, err)
				}
			} else {
				// Update Automated Backups config
				if err := c.UpdateTableWithAutomatedBackupPolicy(ctxWithTimeout, name, abp); err != nil {
					return fmt.Errorf("Error updating automated backup configuration on table %v: %s", name, err)
				}
			}
		}
	}

	return resourceBigtableTableRead(d, meta)
}

func resourceBigtableTableDestroy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	ctx := context.Background()

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	instanceName := tpgresource.GetResourceNameFromSelfLink(d.Get("instance_name").(string))
	c, err := config.BigTableClientFactory(userAgent).NewAdminClient(project, instanceName)
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

func FlattenColumnFamily(families []bigtable.FamilyInfo) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0, len(families))

	for _, f := range families {
		data := make(map[string]interface{})
		data["family"] = f.Name
		if _, ok := f.ValueType.(bigtable.AggregateType); ok {
			marshalled, err := bigtable.MarshalJSON(f.ValueType)
			if err != nil {
				return nil, err
			}
			data["type"] = string(marshalled)
		}
		result = append(result, data)
	}

	return result, nil
}

// TODO(rileykarson): Fix the stored import format after rebasing 3.0.0
func resourceBigtableTableImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*transport_tpg.Config)
	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/instances/(?P<instance_name>[^/]+)/tables/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<instance_name>[^/]+)/(?P<name>[^/]+)",
		"(?P<instance_name>[^/]+)/(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := tpgresource.ReplaceVars(d, config, "projects/{{project}}/instances/{{instance_name}}/tables/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func getType(input interface{}) (bigtable.Type, error) {
	if input == nil || input.(string) == "" {
		return nil, nil
	}
	inputType := strings.TrimSuffix(input.(string), "\n")
	switch inputType {
	case "intsum":
		return bigtable.AggregateType{
			Input:      bigtable.Int64Type{},
			Aggregator: bigtable.SumAggregator{},
		}, nil
	case "intmin":
		return bigtable.AggregateType{
			Input:      bigtable.Int64Type{},
			Aggregator: bigtable.MinAggregator{},
		}, nil
	case "intmax":
		return bigtable.AggregateType{
			Input:      bigtable.Int64Type{},
			Aggregator: bigtable.MaxAggregator{},
		}, nil
	case "inthll":
		return bigtable.AggregateType{
			Input:      bigtable.Int64Type{},
			Aggregator: bigtable.HllppUniqueCountAggregator{},
		}, nil
	}

	output, err := bigtable.UnmarshalJSON([]byte(inputType))
	if err != nil {
		return nil, err
	}
	return output, nil
}
