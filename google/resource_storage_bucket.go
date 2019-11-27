package google

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

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
		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIfChange("retention_policy.0.is_locked", isPolicyLocked),
		),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"encryption": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_kms_key_name": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"requester_pays": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"location": {
				Type:     schema.TypeString,
				Default:  "US",
				Optional: true,
				ForceNew: true,
				StateFunc: func(s interface{}) string {
					return strings.ToUpper(s.(string))
				},
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"storage_class": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "STANDARD",
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
										Removed:  "Please use `with_state` instead",
									},
									"with_state": {
										Type:         schema.TypeString,
										Computed:     true,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"LIVE", "ARCHIVED", "ANY", ""}, false),
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

			"versioning": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
					},
				},
			},

			"website": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"main_page_suffix": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: []string{"website.0.not_found_page", "website.0.main_page_suffix"},
						},
						"not_found_page": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: []string{"website.0.main_page_suffix", "website.0.not_found_page"},
						},
					},
				},
			},

			"retention_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_locked": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"retention_period": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, math.MaxInt32),
						},
					},
				},
			},

			"cors": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"method": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"response_header": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"max_age_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"logging": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_bucket": {
							Type:     schema.TypeString,
							Required: true,
						},
						"log_object_prefix": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"bucket_policy_only": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

// Is the old bucket retention policy locked?
func isPolicyLocked(old, new, _ interface{}) bool {
	if old == nil || new == nil {
		return false
	}

	// if the old policy is locked, but the new policy is not
	if old.(bool) && !new.(bool) {
		return true
	}

	return false
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
		Name:             bucket,
		Labels:           expandLabels(d),
		Location:         location,
		IamConfiguration: expandIamConfiguration(d),
	}

	if v, ok := d.GetOk("storage_class"); ok {
		sb.StorageClass = v.(string)
	}

	lifecycle, err := expandStorageBucketLifecycle(d.Get("lifecycle_rule"))
	if err != nil {
		return err
	}
	sb.Lifecycle = lifecycle

	if v, ok := d.GetOk("versioning"); ok {
		sb.Versioning = expandBucketVersioning(v)
	}

	if v, ok := d.GetOk("website"); ok {
		sb.Website = expandBucketWebsite(v.([]interface{}))
	}

	if v, ok := d.GetOk("retention_policy"); ok {
		retention_policies := v.([]interface{})

		sb.RetentionPolicy = &storage.BucketRetentionPolicy{}

		retentionPolicy := retention_policies[0].(map[string]interface{})

		if v, ok := retentionPolicy["retention_period"]; ok {
			sb.RetentionPolicy.RetentionPeriod = int64(v.(int))
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

	if v, ok := d.GetOk("requester_pays"); ok {
		sb.Billing = &storage.BucketBilling{
			RequesterPays: v.(bool),
		}
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

	// If the retention policy is not already locked, check if it
	// needs to be locked.
	if v, ok := d.GetOk("retention_policy"); ok && !res.RetentionPolicy.IsLocked {
		retention_policies := v.([]interface{})

		sb.RetentionPolicy = &storage.BucketRetentionPolicy{}

		retentionPolicy := retention_policies[0].(map[string]interface{})

		if locked, ok := retentionPolicy["is_locked"]; ok && locked.(bool) {
			err = lockRetentionPolicy(config.clientStorage.Buckets, bucket, res.Metageneration)
			if err != nil {
				return err
			}

			log.Printf("[DEBUG] Locked bucket %v at location %v\n\n", res.Name, res.SelfLink)
		}
	}

	return resourceStorageBucketRead(d, meta)
}

func resourceStorageBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	sb := &storage.Bucket{}

	if d.HasChange("lifecycle_rule") {
		lifecycle, err := expandStorageBucketLifecycle(d.Get("lifecycle_rule"))
		if err != nil {
			return err
		}
		sb.Lifecycle = lifecycle
	}

	if d.HasChange("requester_pays") {
		v := d.Get("requester_pays")
		sb.Billing = &storage.BucketBilling{
			RequesterPays:   v.(bool),
			ForceSendFields: []string{"RequesterPays"},
		}
	}

	if d.HasChange("versioning") {
		if v, ok := d.GetOk("versioning"); ok {
			sb.Versioning = expandBucketVersioning(v)
		}
	}

	if d.HasChange("website") {
		sb.Website = expandBucketWebsite(d.Get("website"))
	}

	if d.HasChange("retention_policy") {
		if v, ok := d.GetOk("retention_policy"); ok {
			sb.RetentionPolicy = expandBucketRetentionPolicy(v.([]interface{}))
		} else {
			sb.NullFields = append(sb.NullFields, "RetentionPolicy")
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

	if d.HasChange("storage_class") {
		if v, ok := d.GetOk("storage_class"); ok {
			sb.StorageClass = v.(string)
		}
	}

	if d.HasChange("bucket_policy_only") {
		sb.IamConfiguration = expandIamConfiguration(d)
	}

	res, err := config.clientStorage.Buckets.Patch(d.Get("name").(string), sb).Do()

	if err != nil {
		return err
	}

	// Assign the bucket ID as the resource ID
	d.Set("self_link", res.SelfLink)

	if d.HasChange("retention_policy") {
		if v, ok := d.GetOk("retention_policy"); ok {
			retention_policies := v.([]interface{})

			sb.RetentionPolicy = &storage.BucketRetentionPolicy{}

			retentionPolicy := retention_policies[0].(map[string]interface{})

			if locked, ok := retentionPolicy["is_locked"]; ok && locked.(bool) && d.HasChange("retention_policy.0.is_locked") {
				err = lockRetentionPolicy(config.clientStorage.Buckets, d.Get("name").(string), res.Metageneration)
				if err != nil {
					return err
				}
			}
		}
	}

	log.Printf("[DEBUG] Patched bucket %v at location %v\n\n", res.Name, res.SelfLink)

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

	// We are trying to support several different use cases for bucket. Buckets are globally
	// unique but they are associated with projects internally, but some users want to use
	// buckets in a project agnostic way. Thus we will check to see if the project ID has been
	// explicitly set and use that first. However if no project is explicitly set, such as during
	// import, we will look up the ID from the compute API using the project Number from the
	// bucket API response.
	// If you are working in a project-agnostic way and have not set the project ID in the provider
	// block, or the resource or an environment variable, we use the compute API to lookup the projectID
	// from the projectNumber which is included in the bucket API response
	if d.Get("project") == "" {
		project, _ := getProject(d, config)
		d.Set("project", project)
	}
	if d.Get("project") == "" {
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
	d.Set("website", flattenBucketWebsite(res.Website))
	d.Set("retention_policy", flattenBucketRetentionPolicy(res.RetentionPolicy))

	if res.IamConfiguration != nil && res.IamConfiguration.BucketPolicyOnly != nil {
		d.Set("bucket_policy_only", res.IamConfiguration.BucketPolicyOnly.Enabled)
	} else {
		d.Set("bucket_policy_only", false)
	}

	if res.Billing == nil {
		d.Set("requester_pays", nil)
	} else {
		d.Set("requester_pays", res.Billing.RequesterPays)
	}

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
			if d.Get("retention_policy.0.is_locked").(bool) {
				for _, item := range res.Items {
					expiration, err := time.Parse(time.RFC3339, item.RetentionExpirationTime)
					if err != nil {
						return err
					}
					if expiration.After(time.Now()) {
						deleteErr := errors.New("Bucket '" + d.Get("name").(string) + "' contains objects that have not met the retention period yet and cannot be deleted.")
						log.Printf("Error! %s : %s\n\n", bucket, deleteErr)
						return deleteErr
					}
				}
			}

			if d.Get("force_destroy").(bool) {
				// GCS requires that a bucket be empty (have no objects or object
				// versions) before it can be deleted.
				log.Printf("[DEBUG] GCS Bucket attempting to forceDestroy\n\n")

				// Create a workerpool for parallel deletion of resources. In the
				// future, it would be great to expose Terraform's global parallelism
				// flag here, but that's currently reserved for core use. Testing
				// shows that NumCPUs-1 is the most performant on average networks.
				//
				// The challenge with making this user-configurable is that the
				// configuration would reside in the Terraform configuration file,
				// decreasing its portability. Ideally we'd want this to connect to
				// Terraform's top-level -parallelism flag, but that's not plumbed nor
				// is it scheduled to be plumbed to individual providers.
				wp := workerpool.New(runtime.NumCPU() - 1)

				for _, object := range res.Items {
					log.Printf("[DEBUG] Found %s", object.Name)
					object := object

					wp.Submit(func() {
						log.Printf("[TRACE] Attempting to delete %s", object.Name)
						if err := config.clientStorage.Objects.Delete(bucket, object.Name).Generation(object.Generation).Do(); err != nil {
							// We should really return an error here, but it doesn't really
							// matter since the following step (bucket deletion) will fail
							// with an error indicating objects are still present, and this
							// log line will point to that object.
							log.Printf("[ERR] Failed to delete storage object %s: %s", object.Name, err)
						} else {
							log.Printf("[TRACE] Successfully deleted %s", object.Name)
						}
					})
				}

				// Wait for everything to finish.
				wp.StopWait()
			} else {
				deleteErr := errors.New("Error trying to delete a bucket containing objects without `force_destroy` set to true")
				log.Printf("Error! %s : %s\n\n", bucket, deleteErr)
				return deleteErr
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
	// We need to support project/bucket_name and bucket_name formats. This will allow
	// importing a bucket that is in a different project than the provider default.
	// ParseImportID can't be used because having no project will cause an error but it
	// is a valid state as the project_id will be retrieved in READ
	parts := strings.Split(d.Id(), "/")
	if len(parts) == 1 {
		d.Set("name", parts[0])
	} else if len(parts) > 1 {
		d.Set("project", parts[0])
		d.Set("name", parts[1])
	}

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
	if len(encs) == 0 || encs[0] == nil {
		return nil
	}
	enc := encs[0].(map[string]interface{})
	keyname := enc["default_kms_key_name"]
	if keyname == nil || keyname.(string) == "" {
		return nil
	}
	bucketenc := &storage.BucketEncryption{
		DefaultKmsKeyName: keyname.(string),
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
	if len(loggings) == 0 {
		return nil
	}

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

func expandBucketRetentionPolicy(configured interface{}) *storage.BucketRetentionPolicy {
	retentionPolicies := configured.([]interface{})
	retentionPolicy := retentionPolicies[0].(map[string]interface{})

	bucketRetentionPolicy := &storage.BucketRetentionPolicy{
		IsLocked:        retentionPolicy["is_locked"].(bool),
		RetentionPeriod: int64(retentionPolicy["retention_period"].(int)),
	}

	return bucketRetentionPolicy
}

func flattenBucketRetentionPolicy(bucketRetentionPolicy *storage.BucketRetentionPolicy) []map[string]interface{} {
	bucketRetentionPolicies := make([]map[string]interface{}, 0, 1)

	if bucketRetentionPolicy == nil {
		return bucketRetentionPolicies
	}

	retentionPolicy := map[string]interface{}{
		"is_locked":        bucketRetentionPolicy.IsLocked,
		"retention_period": bucketRetentionPolicy.RetentionPeriod,
	}

	bucketRetentionPolicies = append(bucketRetentionPolicies, retentionPolicy)
	return bucketRetentionPolicies
}

func expandBucketVersioning(configured interface{}) *storage.BucketVersioning {
	versionings := configured.([]interface{})
	if len(versionings) == 0 {
		return nil
	}

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
	if condition.IsLive == nil {
		ruleCondition["with_state"] = "ANY"
	} else {
		if *condition.IsLive {
			ruleCondition["with_state"] = "LIVE"
		} else {
			ruleCondition["with_state"] = "ARCHIVED"
		}
	}
	return ruleCondition
}

func flattenBucketWebsite(website *storage.BucketWebsite) []map[string]interface{} {
	if website == nil {
		return nil
	}
	websites := make([]map[string]interface{}, 0, 1)
	websites = append(websites, map[string]interface{}{
		"main_page_suffix": website.MainPageSuffix,
		"not_found_page":   website.NotFoundPage,
	})

	return websites
}

func expandBucketWebsite(v interface{}) *storage.BucketWebsite {
	if v == nil {
		return nil
	}
	vs := v.([]interface{})

	if len(vs) < 1 || vs[0] == nil {
		return nil
	}

	website := vs[0].(map[string]interface{})
	w := &storage.BucketWebsite{}

	if v := website["not_found_page"]; v != "" {
		w.NotFoundPage = v.(string)
	}

	if v := website["main_page_suffix"]; v != "" {
		w.MainPageSuffix = v.(string)
	}

	return w
}

func expandIamConfiguration(d *schema.ResourceData) *storage.BucketIamConfiguration {
	return &storage.BucketIamConfiguration{
		ForceSendFields: []string{"BucketPolicyOnly"},
		BucketPolicyOnly: &storage.BucketIamConfigurationBucketPolicyOnly{
			Enabled:         d.Get("bucket_policy_only").(bool),
			ForceSendFields: []string{"Enabled"},
		},
	}
}

func expandStorageBucketLifecycle(v interface{}) (*storage.BucketLifecycle, error) {
	if v == nil {
		return &storage.BucketLifecycle{
			ForceSendFields: []string{"Rule"},
		}, nil
	}
	lifecycleRules := v.([]interface{})
	transformedRules := make([]*storage.BucketLifecycleRule, 0, len(lifecycleRules))

	for _, v := range lifecycleRules {
		rule, err := expandStorageBucketLifecycleRule(v)
		if err != nil {
			return nil, err
		}
		transformedRules = append(transformedRules, rule)
	}

	if len(transformedRules) == 0 {
		return &storage.BucketLifecycle{
			ForceSendFields: []string{"Rule"},
		}, nil
	}

	return &storage.BucketLifecycle{
		Rule: transformedRules,
	}, nil
}

func expandStorageBucketLifecycleRule(v interface{}) (*storage.BucketLifecycleRule, error) {
	if v == nil {
		return nil, nil
	}

	rule := v.(map[string]interface{})
	transformed := &storage.BucketLifecycleRule{}

	if v, ok := rule["action"]; ok {
		action, err := expandStorageBucketLifecycleRuleAction(v)
		if err != nil {
			return nil, err
		}
		transformed.Action = action
	} else {
		return nil, fmt.Errorf("exactly one action is required for lifecycle_rule")
	}

	if v, ok := rule["condition"]; ok {
		cond, err := expandStorageBucketLifecycleRuleCondition(v)
		if err != nil {
			return nil, err
		}
		transformed.Condition = cond
	}

	return transformed, nil
}

func expandStorageBucketLifecycleRuleAction(v interface{}) (*storage.BucketLifecycleRuleAction, error) {
	if v == nil {
		return nil, fmt.Errorf("exactly one action is required for lifecycle_rule")
	}

	actions := v.(*schema.Set).List()
	if len(actions) != 1 {
		return nil, fmt.Errorf("exactly one action is required for lifecycle_rule")
	}

	action := actions[0].(map[string]interface{})
	transformed := &storage.BucketLifecycleRuleAction{}

	if v, ok := action["type"]; ok {
		transformed.Type = v.(string)
	}

	if v, ok := action["storage_class"]; ok {
		transformed.StorageClass = v.(string)
	}

	return transformed, nil
}

func expandStorageBucketLifecycleRuleCondition(v interface{}) (*storage.BucketLifecycleRuleCondition, error) {
	if v == nil {
		return nil, nil
	}
	conditions := v.(*schema.Set).List()
	if len(conditions) != 1 {
		return nil, fmt.Errorf("One and only one condition can be provided per lifecycle_rule")
	}

	condition := conditions[0].(map[string]interface{})
	transformed := &storage.BucketLifecycleRuleCondition{}

	if v, ok := condition["age"]; ok {
		transformed.Age = int64(v.(int))
	}

	if v, ok := condition["created_before"]; ok {
		transformed.CreatedBefore = v.(string)
	}

	withStateV, withStateOk := condition["with_state"]
	// Because TF schema, withStateOk currently will always be true,
	// do the check just in case.
	if withStateOk {
		switch withStateV.(string) {
		case "LIVE":
			transformed.IsLive = googleapi.Bool(true)
		case "ARCHIVED":
			transformed.IsLive = googleapi.Bool(false)
		case "ANY", "":
			// This is unnecessary, but set explicitly to nil for readability.
			transformed.IsLive = nil
		default:
			return nil, fmt.Errorf("unexpected value %q for condition.with_state", withStateV.(string))
		}
	}

	if v, ok := condition["matches_storage_class"]; ok {
		classes := v.([]interface{})
		transformedClasses := make([]string, 0, len(classes))

		for _, v := range classes {
			transformedClasses = append(transformedClasses, v.(string))
		}
		transformed.MatchesStorageClass = transformedClasses
	}

	if v, ok := condition["num_newer_versions"]; ok {
		transformed.NumNewerVersions = int64(v.(int))
	}

	return transformed, nil
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

	withStateV, withStateOk := m["with_state"]
	if withStateOk {
		switch withStateV.(string) {
		case "LIVE":
			buf.WriteString(fmt.Sprintf("%t-", true))
		case "ARCHIVED":
			buf.WriteString(fmt.Sprintf("%t-", false))
		}
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

func lockRetentionPolicy(bucketsService *storage.BucketsService, bucketName string, metageneration int64) error {
	lockPolicyCall := bucketsService.LockRetentionPolicy(bucketName, metageneration)
	if _, err := lockPolicyCall.Do(); err != nil {
		return err
	}

	return nil
}
