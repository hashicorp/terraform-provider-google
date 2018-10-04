package google

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"google.golang.org/api/googleapi"
	"google.golang.org/api/storage/v1"
)

func resourceStorageBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageBucketCreate,
		Read:   resourceStorageBucketRead,
		Update: resourceStorageBucketUpdate,
		Delete: resourceStorageBucketDelete,
		Importer: &schema.ResourceImporter{
			State: resourceStorageBucketStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"encryption": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_kms_key_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"force_destroy": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"location": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "US",
				Optional: true,
				ForceNew: true,
				StateFunc: func(s interface{}) string {
					return strings.ToUpper(s.(string))
				},
			},

			"predefined_acl": &schema.Schema{
				Type:     schema.TypeString,
				Removed:  "Please use resource \"storage_bucket_acl.predefined_acl\" instead.",
				Optional: true,
				ForceNew: true,
			},

			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"url": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"storage_class": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "STANDARD",
				ForceNew: true,
			},

			"lifecycle_rule": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 100,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							MaxItems: 1,
							Set:      resourceGCSBucketLifecycleRuleActionHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Required: true,
									},
									"storage_class": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"condition": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							MaxItems: 1,
							Set:      resourceGCSBucketLifecycleRuleConditionHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"age": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"created_before": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"is_live": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"matches_storage_class": {
										Type:     schema.TypeList,
										Optional: true,
										MinItems: 1,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"num_newer_versions": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"versioning": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"website": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"main_page_suffix": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"not_found_page": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"cors": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"method": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"response_header": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"max_age_seconds": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"logging": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_bucket": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"log_object_prefix": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceStorageBucketCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Get the bucket and location
	bucket := d.Get("name").(string)
	location := d.Get("location").(string)

	// Create a bucket, setting the labels, location and name.
	sb := &storage.Bucket{
		Name:     bucket,
		Labels:   expandLabels(d),
		Location: location,
	}

	if v, ok := d.GetOk("storage_class"); ok {
		sb.StorageClass = v.(string)
	}

	if err := resourceGCSBucketLifecycleCreateOrUpdate(d, sb); err != nil {
		return err
	}

	if v, ok := d.GetOk("versioning"); ok {
		sb.Versioning = expandBucketVersioning(v)
	}

	if v, ok := d.GetOk("website"); ok {
		websites := v.([]interface{})

		if len(websites) > 1 {
			return fmt.Errorf("At most one website block is allowed")
		}

		sb.Website = &storage.BucketWebsite{}

		website := websites[0].(map[string]interface{})

		if v, ok := website["not_found_page"]; ok {
			sb.Website.NotFoundPage = v.(string)
		}

		if v, ok := website["main_page_suffix"]; ok {
			sb.Website.MainPageSuffix = v.(string)
		}
	}

	if v, ok := d.GetOk("cors"); ok {
		sb.Cors = expandCors(v.([]interface{}))
	}

	if v, ok := d.GetOk("logging"); ok {
		sb.Logging = expandBucketLogging(v.([]interface{}))
	}

	if v, ok := d.GetOk("encryption"); ok {
		sb.Encryption = expandBucketEncryption(v.([]interface{}))
	}

	var res *storage.Bucket

	err = retry(func() error {
		res, err = config.clientStorage.Buckets.Insert(project, sb).Do()
		return err
	})

	if err != nil {
		fmt.Printf("Error creating bucket %s: %v", bucket, err)
		return err
	}

	log.Printf("[DEBUG] Created bucket %v at location %v\n\n", res.Name, res.SelfLink)

	d.SetId(res.Id)
	return resourceStorageBucketRead(d, meta)
}

func resourceStorageBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	sb := &storage.Bucket{}

	if d.HasChange("lifecycle_rule") {
		if err := resourceGCSBucketLifecycleCreateOrUpdate(d, sb); err != nil {
			return err
		}
	}

	if d.HasChange("versioning") {
		if v, ok := d.GetOk("versioning"); ok {
			sb.Versioning = expandBucketVersioning(v)
		}
	}

	if d.HasChange("website") {
		if v, ok := d.GetOk("website"); ok {
			websites := v.([]interface{})

			if len(websites) > 1 {
				return fmt.Errorf("At most one website block is allowed")
			}

			// Setting fields to "" to be explicit that the PATCH call will
			// delete this field.
			if len(websites) == 0 {
				sb.Website.NotFoundPage = ""
				sb.Website.MainPageSuffix = ""
			} else {
				website := websites[0].(map[string]interface{})
				sb.Website = &storage.BucketWebsite{}
				if v, ok := website["not_found_page"]; ok {
					sb.Website.NotFoundPage = v.(string)
				} else {
					sb.Website.NotFoundPage = ""
				}

				if v, ok := website["main_page_suffix"]; ok {
					sb.Website.MainPageSuffix = v.(string)
				} else {
					sb.Website.MainPageSuffix = ""
				}
			}
		}
	}

	if v, ok := d.GetOk("cors"); ok {
		sb.Cors = expandCors(v.([]interface{}))
	}

	if d.HasChange("logging") {
		if v, ok := d.GetOk("logging"); ok {
			sb.Logging = expandBucketLogging(v.([]interface{}))
		} else {
			sb.NullFields = append(sb.NullFields, "Logging")
		}
	}

	if d.HasChange("encryption") {
		if v, ok := d.GetOk("encryption"); ok {
			sb.Encryption = expandBucketEncryption(v.([]interface{}))
		} else {
			sb.NullFields = append(sb.NullFields, "Encryption")
		}
	}

	if d.HasChange("labels") {
		sb.Labels = expandLabels(d)
		if len(sb.Labels) == 0 {
			sb.NullFields = append(sb.NullFields, "Labels")
		}

		// To delete a label using PATCH, we have to explicitly set its value
		// to null.
		old, _ := d.GetChange("labels")
		for k := range old.(map[string]interface{}) {
			if _, ok := sb.Labels[k]; !ok {
				sb.NullFields = append(sb.NullFields, fmt.Sprintf("Labels.%s", k))
			}
		}
	}

	res, err := config.clientStorage.Buckets.Patch(d.Get("name").(string), sb).Do()

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Patched bucket %v at location %v\n\n", res.Name, res.SelfLink)

	// Assign the bucket ID as the resource ID
	d.Set("self_link", res.SelfLink)
	d.SetId(res.Id)

	return nil
}

func resourceStorageBucketRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Get the bucket and acl
	bucket := d.Get("name").(string)
	res, err := config.clientStorage.Buckets.Get(bucket).Do()

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Storage Bucket %q", d.Get("name").(string)))
	}
	log.Printf("[DEBUG] Read bucket %v at location %v\n\n", res.Name, res.SelfLink)

	// We need to get the project associated with this bucket because otherwise import
	// won't work properly.  That means we need to call the projects.get API with the
	// project number, to get the project ID - there's no project ID field in the
	// resource response.  However, this requires a call to the Compute API, which
	// would otherwise not be required for this resource.  So, we're going to
	// intentionally check whether the project is set *on the resource*.  If it is,
	// we will not try to fetch the project name.  If it is not, either because
	// the user intends to use the default provider project, or because the resource
	// is currently being imported, we will read it from the API.
	if _, ok := d.GetOk("project"); !ok {
		proj, err := config.clientCompute.Projects.Get(strconv.FormatUint(res.ProjectNumber, 10)).Do()
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] Bucket %v is in project number %v, which is project ID %s.\n", res.Name, res.ProjectNumber, proj.Name)
		d.Set("project", proj.Name)
	}

	// Update the bucket ID according to the resource ID
	d.Set("self_link", res.SelfLink)
	d.Set("url", fmt.Sprintf("gs://%s", bucket))
	d.Set("storage_class", res.StorageClass)
	d.Set("encryption", flattenBucketEncryption(res.Encryption))
	d.Set("location", res.Location)
	d.Set("cors", flattenCors(res.Cors))
	d.Set("logging", flattenBucketLogging(res.Logging))
	d.Set("versioning", flattenBucketVersioning(res.Versioning))
	d.Set("lifecycle_rule", flattenBucketLifecycle(res.Lifecycle))
	d.Set("labels", res.Labels)
	d.SetId(res.Id)
	return nil
}

func resourceStorageBucketDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Get the bucket
	bucket := d.Get("name").(string)

	for {
		res, err := config.clientStorage.Objects.List(bucket).Versions(true).Do()
		if err != nil {
			fmt.Printf("Error Objects.List failed: %v", err)
			return err
		}

		if len(res.Items) != 0 {
			if d.Get("force_destroy").(bool) {
				// purge the bucket...
				log.Printf("[DEBUG] GCS Bucket attempting to forceDestroy\n\n")

				for _, object := range res.Items {
					log.Printf("[DEBUG] Found %s", object.Name)
					if err := config.clientStorage.Objects.Delete(bucket, object.Name).Generation(object.Generation).Do(); err != nil {
						log.Fatalf("Error trying to delete object: %s %s\n\n", object.Name, err)
					} else {
						log.Printf("Object deleted: %s \n\n", object.Name)
					}
				}

			} else {
				delete_err := errors.New("Error trying to delete a bucket containing objects without `force_destroy` set to true")
				log.Printf("Error! %s : %s\n\n", bucket, delete_err)
				return delete_err
			}
		} else {
			break // 0 items, bucket empty
		}
	}

	// remove empty bucket
	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		err := config.clientStorage.Buckets.Delete(bucket).Do()
		if err == nil {
			return nil
		}
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 429 {
			return resource.RetryableError(gerr)
		}
		return resource.NonRetryableError(err)
	})
	if err != nil {
		fmt.Printf("Error deleting bucket %s: %v\n\n", bucket, err)
		return err
	}
	log.Printf("[DEBUG] Deleted bucket %v\n\n", bucket)

	return nil
}

func resourceStorageBucketStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("name", d.Id())
	d.Set("force_destroy", false)
	return []*schema.ResourceData{d}, nil
}

func expandCors(configured []interface{}) []*storage.BucketCors {
	corsRules := make([]*storage.BucketCors, 0, len(configured))
	for _, raw := range configured {
		data := raw.(map[string]interface{})
		corsRule := storage.BucketCors{
			Origin:         convertStringArr(data["origin"].([]interface{})),
			Method:         convertStringArr(data["method"].([]interface{})),
			ResponseHeader: convertStringArr(data["response_header"].([]interface{})),
			MaxAgeSeconds:  int64(data["max_age_seconds"].(int)),
		}

		corsRules = append(corsRules, &corsRule)
	}
	return corsRules
}

func flattenCors(corsRules []*storage.BucketCors) []map[string]interface{} {
	corsRulesSchema := make([]map[string]interface{}, 0, len(corsRules))
	for _, corsRule := range corsRules {
		data := map[string]interface{}{
			"origin":          corsRule.Origin,
			"method":          corsRule.Method,
			"response_header": corsRule.ResponseHeader,
			"max_age_seconds": corsRule.MaxAgeSeconds,
		}

		corsRulesSchema = append(corsRulesSchema, data)
	}
	return corsRulesSchema
}

func expandBucketEncryption(configured interface{}) *storage.BucketEncryption {
	encs := configured.([]interface{})
	enc := encs[0].(map[string]interface{})
	bucketenc := &storage.BucketEncryption{
		DefaultKmsKeyName: enc["default_kms_key_name"].(string),
	}
	return bucketenc
}

func flattenBucketEncryption(enc *storage.BucketEncryption) []map[string]interface{} {
	encryption := make([]map[string]interface{}, 0, 1)

	if enc == nil {
		return encryption
	}

	encryption = append(encryption, map[string]interface{}{
		"default_kms_key_name": enc.DefaultKmsKeyName,
	})

	return encryption
}

func expandBucketLogging(configured interface{}) *storage.BucketLogging {
	loggings := configured.([]interface{})
	logging := loggings[0].(map[string]interface{})

	bucketLogging := &storage.BucketLogging{
		LogBucket:       logging["log_bucket"].(string),
		LogObjectPrefix: logging["log_object_prefix"].(string),
	}

	return bucketLogging
}

func flattenBucketLogging(bucketLogging *storage.BucketLogging) []map[string]interface{} {
	loggings := make([]map[string]interface{}, 0, 1)

	if bucketLogging == nil {
		return loggings
	}

	logging := map[string]interface{}{
		"log_bucket":        bucketLogging.LogBucket,
		"log_object_prefix": bucketLogging.LogObjectPrefix,
	}

	loggings = append(loggings, logging)
	return loggings
}

func expandBucketVersioning(configured interface{}) *storage.BucketVersioning {
	versionings := configured.([]interface{})
	versioning := versionings[0].(map[string]interface{})

	bucketVersioning := &storage.BucketVersioning{}

	bucketVersioning.Enabled = versioning["enabled"].(bool)
	bucketVersioning.ForceSendFields = append(bucketVersioning.ForceSendFields, "Enabled")

	return bucketVersioning
}

func flattenBucketVersioning(bucketVersioning *storage.BucketVersioning) []map[string]interface{} {
	versionings := make([]map[string]interface{}, 0, 1)

	if bucketVersioning == nil {
		return versionings
	}

	versioning := map[string]interface{}{
		"enabled": bucketVersioning.Enabled,
	}
	versionings = append(versionings, versioning)
	return versionings
}

func flattenBucketLifecycle(lifecycle *storage.BucketLifecycle) []map[string]interface{} {
	if lifecycle == nil || lifecycle.Rule == nil {
		return []map[string]interface{}{}
	}

	rules := make([]map[string]interface{}, 0, len(lifecycle.Rule))

	for _, rule := range lifecycle.Rule {
		rules = append(rules, map[string]interface{}{
			"action":    schema.NewSet(resourceGCSBucketLifecycleRuleActionHash, []interface{}{flattenBucketLifecycleRuleAction(rule.Action)}),
			"condition": schema.NewSet(resourceGCSBucketLifecycleRuleConditionHash, []interface{}{flattenBucketLifecycleRuleCondition(rule.Condition)}),
		})
	}

	return rules
}

func flattenBucketLifecycleRuleAction(action *storage.BucketLifecycleRuleAction) map[string]interface{} {
	return map[string]interface{}{
		"type":          action.Type,
		"storage_class": action.StorageClass,
	}
}

func flattenBucketLifecycleRuleCondition(condition *storage.BucketLifecycleRuleCondition) map[string]interface{} {
	ruleCondition := map[string]interface{}{
		"age":                   int(condition.Age),
		"created_before":        condition.CreatedBefore,
		"matches_storage_class": convertStringArrToInterface(condition.MatchesStorageClass),
		"num_newer_versions":    int(condition.NumNewerVersions),
	}
	if condition.IsLive != nil {
		ruleCondition["is_live"] = *condition.IsLive
	}
	return ruleCondition
}

func resourceGCSBucketLifecycleCreateOrUpdate(d *schema.ResourceData, sb *storage.Bucket) error {
	if v, ok := d.GetOk("lifecycle_rule"); ok {
		lifecycle_rules := v.([]interface{})

		sb.Lifecycle = &storage.BucketLifecycle{}
		sb.Lifecycle.Rule = make([]*storage.BucketLifecycleRule, 0, len(lifecycle_rules))

		for _, raw_lifecycle_rule := range lifecycle_rules {
			lifecycle_rule := raw_lifecycle_rule.(map[string]interface{})

			target_lifecycle_rule := &storage.BucketLifecycleRule{}

			if v, ok := lifecycle_rule["action"]; ok {
				if actions := v.(*schema.Set).List(); len(actions) == 1 {
					action := actions[0].(map[string]interface{})

					target_lifecycle_rule.Action = &storage.BucketLifecycleRuleAction{}

					if v, ok := action["type"]; ok {
						target_lifecycle_rule.Action.Type = v.(string)
					}

					if v, ok := action["storage_class"]; ok {
						target_lifecycle_rule.Action.StorageClass = v.(string)
					}
				} else {
					return fmt.Errorf("Exactly one action is required")
				}
			}

			if v, ok := lifecycle_rule["condition"]; ok {
				if conditions := v.(*schema.Set).List(); len(conditions) == 1 {
					condition := conditions[0].(map[string]interface{})

					target_lifecycle_rule.Condition = &storage.BucketLifecycleRuleCondition{}

					if v, ok := condition["age"]; ok {
						target_lifecycle_rule.Condition.Age = int64(v.(int))
					}

					if v, ok := condition["created_before"]; ok {
						target_lifecycle_rule.Condition.CreatedBefore = v.(string)
					}

					if v, ok := condition["is_live"]; ok {
						target_lifecycle_rule.Condition.IsLive = googleapi.Bool(v.(bool))
					}

					if v, ok := condition["matches_storage_class"]; ok {
						matches_storage_classes := v.([]interface{})

						target_matches_storage_classes := make([]string, 0, len(matches_storage_classes))

						for _, v := range matches_storage_classes {
							target_matches_storage_classes = append(target_matches_storage_classes, v.(string))
						}

						target_lifecycle_rule.Condition.MatchesStorageClass = target_matches_storage_classes
					}

					if v, ok := condition["num_newer_versions"]; ok {
						target_lifecycle_rule.Condition.NumNewerVersions = int64(v.(int))
					}
				} else {
					return fmt.Errorf("Exactly one condition is required")
				}
			}

			sb.Lifecycle.Rule = append(sb.Lifecycle.Rule, target_lifecycle_rule)
		}
	} else {
		sb.Lifecycle = &storage.BucketLifecycle{
			ForceSendFields: []string{"Rule"},
		}
	}

	return nil
}

func resourceGCSBucketLifecycleRuleActionHash(v interface{}) int {
	if v == nil {
		return 0
	}

	var buf bytes.Buffer
	m := v.(map[string]interface{})

	buf.WriteString(fmt.Sprintf("%s-", m["type"].(string)))

	if v, ok := m["storage_class"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	return hashcode.String(buf.String())
}

func resourceGCSBucketLifecycleRuleConditionHash(v interface{}) int {
	if v == nil {
		return 0
	}

	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if v, ok := m["age"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}

	if v, ok := m["created_before"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}

	if v, ok := m["is_live"]; ok {
		buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
	}

	if v, ok := m["matches_storage_class"]; ok {
		matches_storage_classes := v.([]interface{})
		for _, matches_storage_class := range matches_storage_classes {
			buf.WriteString(fmt.Sprintf("%s-", matches_storage_class))
		}
	}

	if v, ok := m["num_newer_versions"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}

	return hashcode.String(buf.String())
}
