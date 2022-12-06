package google

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gammazero/workerpool"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Read:   schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the bucket.`,
			},

			"encryption": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_kms_key_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `A Cloud KMS key that will be used to encrypt objects inserted into this bucket, if no encryption method is specified. You must pay attention to whether the crypto key is available in the location that this bucket is created in. See the docs for more details.`,
						},
					},
				},
				Description: `The bucket's encryption configuration.`,
			},

			"requester_pays": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Enables Requester Pays on a storage bucket.`,
			},

			"force_destroy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `When deleting a bucket, this boolean option will delete all contained objects. If you try to delete a bucket that contains objects, Terraform will fail that run.`,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				// GCP (Dataplex) automatically adds labels
				DiffSuppressFunc: resourceDataplexLabelDiffSuppress,
				Elem:             &schema.Schema{Type: schema.TypeString},
				Description:      `A set of key/value label pairs to assign to the bucket.`,
			},

			"location": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(s interface{}) string {
					return strings.ToUpper(s.(string))
				},
				Description: `The Google Cloud Storage location`,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The ID of the project in which the resource belongs. If it is not provided, the provider project is used.`,
			},

			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URI of the created resource.`,
			},

			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The base URL of the bucket, in the format gs://<bucket-name>.`,
			},

			"storage_class": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "STANDARD",
				Description: `The Storage Class of the new bucket. Supported values include: STANDARD, MULTI_REGIONAL, REGIONAL, NEARLINE, COLDLINE, ARCHIVE.`,
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
										Type:        schema.TypeString,
										Required:    true,
										Description: `The type of the action of this Lifecycle Rule. Supported values include: Delete, SetStorageClass and AbortIncompleteMultipartUpload.`,
									},
									"storage_class": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The target Storage Class of objects affected by this Lifecycle Rule. Supported values include: MULTI_REGIONAL, REGIONAL, NEARLINE, COLDLINE, ARCHIVE.`,
									},
								},
							},
							Description: `The Lifecycle Rule's action configuration. A single block of this type is supported.`,
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
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `Minimum age of an object in days to satisfy this condition.`,
									},
									"created_before": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Creation date of an object in RFC 3339 (e.g. 2017-06-13) to satisfy this condition.`,
									},
									"custom_time_before": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Creation date of an object in RFC 3339 (e.g. 2017-06-13) to satisfy this condition.`,
									},
									"days_since_custom_time": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `Number of days elapsed since the user-specified timestamp set on an object.`,
									},
									"days_since_noncurrent_time": {
										Type:     schema.TypeInt,
										Optional: true,
										Description: `Number of days elapsed since the noncurrent timestamp of an object. This
										condition is relevant only for versioned objects.`,
									},
									"noncurrent_time_before": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Creation date of an object in RFC 3339 (e.g. 2017-06-13) to satisfy this condition.`,
									},
									"with_state": {
										Type:         schema.TypeString,
										Computed:     true,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"LIVE", "ARCHIVED", "ANY", ""}, false),
										Description:  `Match to live and/or archived objects. Unversioned buckets have only live objects. Supported values include: "LIVE", "ARCHIVED", "ANY".`,
									},
									"matches_storage_class": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `Storage Class of objects to satisfy this condition. Supported values include: MULTI_REGIONAL, REGIONAL, NEARLINE, COLDLINE, ARCHIVE, STANDARD, DURABLE_REDUCED_AVAILABILITY.`,
									},
									"num_newer_versions": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `Relevant only for versioned objects. The number of newer versions of an object to satisfy this condition.`,
									},
									"matches_prefix": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `One or more matching name prefixes to satisfy this condition.`,
									},
									"matches_suffix": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `One or more matching name suffixes to satisfy this condition.`,
									},
								},
							},
							Description: `The Lifecycle Rule's condition configuration.`,
						},
					},
				},
				Description: `The bucket's Lifecycle Rules configuration.`,
			},

			"versioning": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `While set to true, versioning is fully enabled for this bucket.`,
						},
					},
				},
				Description: `The bucket's Versioning configuration.`,
			},

			"autoclass": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							ForceNew:    true,
							Description: `While set to true, autoclass automatically transitions objects in your bucket to appropriate storage classes based on each object's access pattern.`,
						},
					},
				},
				Description: `The bucket's autoclass configuration.`,
			},

			"website": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"main_page_suffix": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: []string{"website.0.not_found_page", "website.0.main_page_suffix"},
							Description:  `Behaves as the bucket's directory index where missing objects are treated as potential directories.`,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return old != "" && new == ""
							},
						},
						"not_found_page": {
							Type:         schema.TypeString,
							Optional:     true,
							AtLeastOneOf: []string{"website.0.main_page_suffix", "website.0.not_found_page"},
							Description:  `The custom object to return when a requested resource is not found.`,
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								return old != "" && new == ""
							},
						},
					},
				},
				Description: `Configuration if the bucket acts as a website.`,
			},

			"retention_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_locked": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: `If set to true, the bucket will be locked and permanently restrict edits to the bucket's retention policy.  Caution: Locking a bucket is an irreversible action.`,
						},
						"retention_period": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, math.MaxInt32),
							Description:  `The period of time, in seconds, that objects in the bucket must be retained and cannot be deleted, overwritten, or archived. The value must be less than 3,155,760,000 seconds.`,
						},
					},
				},
				Description: `Configuration of the bucket's data retention policy for how long objects in the bucket should be retained.`,
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
							Description: `The list of Origins eligible to receive CORS response headers. Note: "*" is permitted in the list of origins, and means "any Origin".`,
						},
						"method": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: `The list of HTTP methods on which to include CORS response headers, (GET, OPTIONS, POST, etc) Note: "*" is permitted in the list of methods, and means "any method".`,
						},
						"response_header": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: `The list of HTTP headers other than the simple response headers to give permission for the user-agent to share across domains.`,
						},
						"max_age_seconds": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `The value, in seconds, to return in the Access-Control-Max-Age header used in preflight responses.`,
						},
					},
				},
				Description: `The bucket's Cross-Origin Resource Sharing (CORS) configuration.`,
			},

			"default_event_based_hold": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether or not to automatically apply an eventBasedHold to new objects added to the bucket.`,
			},

			"logging": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_bucket": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The bucket that will receive log objects.`,
						},
						"log_object_prefix": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The object prefix for log objects. If it's not provided, by default Google Cloud Storage sets this to this bucket's name.`,
						},
					},
				},
				Description: `The bucket's Access & Storage Logs configuration.`,
			},
			"uniform_bucket_level_access": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Enables uniform bucket-level access on a bucket.`,
			},
			"custom_placement_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_locations": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: true,
							MaxItems: 2,
							MinItems: 2,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: `The list of individual regions that comprise a dual-region bucket. See the docs for a list of acceptable regions. Note: If any of the data_locations changes, it will recreate the bucket.`,
						},
					},
				},
				Description: `The bucket's custom location configuration, which specifies the individual regions that comprise a dual-region bucket. If the bucket is designated a single or multi-region, the parameters are empty.`,
			},
			"public_access_prevention": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Prevents public access to a bucket.`,
			},
		},
		UseJSONNumber: true,
	}
}

const resourceDataplexGoogleProvidedLabelPrefix = "labels.goog-dataplex"

func resourceDataplexLabelDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if strings.HasPrefix(k, resourceDataplexGoogleProvidedLabelPrefix) && new == "" {
		return true
	}

	// Let diff be determined by labels (above)
	if strings.HasPrefix(k, "labels.%") {
		return true
	}

	// For other keys, don't suppress diff.
	return false
}

// Is the old bucket retention policy locked?
func isPolicyLocked(_ context.Context, old, new, _ interface{}) bool {
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
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// Get the bucket and location
	bucket := d.Get("name").(string)
	if err := checkGCSName(bucket); err != nil {
		return err
	}
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

	if v, ok := d.GetOk("autoclass"); ok {
		sb.Autoclass = expandBucketAutoclass(v)
	}

	if v, ok := d.GetOk("website"); ok {
		sb.Website = expandBucketWebsite(v.([]interface{}))
	}

	if v, ok := d.GetOk("retention_policy"); ok {
		// Not using expandBucketRetentionPolicy() here because `is_locked` cannot be set on creation.
		retention_policies := v.([]interface{})

		if len(retention_policies) > 0 {
			sb.RetentionPolicy = &storage.BucketRetentionPolicy{}

			retentionPolicy := retention_policies[0].(map[string]interface{})

			if v, ok := retentionPolicy["retention_period"]; ok {
				sb.RetentionPolicy.RetentionPeriod = int64(v.(int))
			}
		}
	}

	if v, ok := d.GetOk("default_event_based_hold"); ok {
		sb.DefaultEventBasedHold = v.(bool)
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

	if v, ok := d.GetOk("custom_placement_config"); ok {
		sb.CustomPlacementConfig = expandBucketCustomPlacementConfig(v.([]interface{}))
	}

	var res *storage.Bucket

	err = retry(func() error {
		res, err = config.NewStorageClient(userAgent).Buckets.Insert(project, sb).Do()
		return err
	})

	if err != nil {
		fmt.Printf("Error creating bucket %s: %v", bucket, err)
		return err
	}

	log.Printf("[DEBUG] Created bucket %v at location %v\n\n", res.Name, res.SelfLink)
	d.SetId(res.Id)

	// There seems to be some eventual consistency errors in some cases, so we want to check a few times
	// to make sure it exists before moving on
	err = retryTimeDuration(func() (operr error) {
		_, retryErr := config.NewStorageClient(userAgent).Buckets.Get(res.Name).Do()
		return retryErr
	}, d.Timeout(schema.TimeoutCreate), isNotFoundRetryableError("bucket creation"))

	if err != nil {
		return fmt.Errorf("Error reading bucket after creation: %s", err)
	}

	// If the retention policy is not already locked, check if it
	// needs to be locked.
	if v, ok := d.GetOk("retention_policy"); ok && !res.RetentionPolicy.IsLocked {
		retention_policies := v.([]interface{})

		sb.RetentionPolicy = &storage.BucketRetentionPolicy{}

		retentionPolicy := retention_policies[0].(map[string]interface{})

		if locked, ok := retentionPolicy["is_locked"]; ok && locked.(bool) {
			err = lockRetentionPolicy(config.NewStorageClient(userAgent).Buckets, bucket, res.Metageneration)
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
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	sb := &storage.Bucket{}

	if detectLifecycleChange(d) {
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

	if d.HasChange("autoclass") {
		if v, ok := d.GetOk("autoclass"); ok {
			sb.Autoclass = expandBucketAutoclass(v)
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

	if d.HasChange("cors") {
		if v, ok := d.GetOk("cors"); ok {
			sb.Cors = expandCors(v.([]interface{}))
		} else {
			sb.NullFields = append(sb.NullFields, "Cors")
		}
	}

	if d.HasChange("default_event_based_hold") {
		v := d.Get("default_event_based_hold")
		sb.DefaultEventBasedHold = v.(bool)
		sb.ForceSendFields = append(sb.ForceSendFields, "DefaultEventBasedHold")
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

	if d.HasChange("uniform_bucket_level_access") || d.HasChange("public_access_prevention") {
		sb.IamConfiguration = expandIamConfiguration(d)
	}

	res, err := config.NewStorageClient(userAgent).Buckets.Patch(d.Get("name").(string), sb).Do()
	if err != nil {
		return err
	}

	// Assign the bucket ID as the resource ID
	if err := d.Set("self_link", res.SelfLink); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}

	// There seems to be some eventual consistency errors in some cases, so we want to check a few times
	// to make sure it exists before moving on
	err = retryTimeDuration(func() (operr error) {
		_, retryErr := config.NewStorageClient(userAgent).Buckets.Get(res.Name).Do()
		return retryErr
	}, d.Timeout(schema.TimeoutUpdate), isNotFoundRetryableError("bucket update"))

	if err != nil {
		return fmt.Errorf("Error reading bucket after update: %s", err)
	}

	if d.HasChange("retention_policy") {
		if v, ok := d.GetOk("retention_policy"); ok {
			retention_policies := v.([]interface{})

			sb.RetentionPolicy = &storage.BucketRetentionPolicy{}

			retentionPolicy := retention_policies[0].(map[string]interface{})

			if locked, ok := retentionPolicy["is_locked"]; ok && locked.(bool) && d.HasChange("retention_policy.0.is_locked") {
				err = lockRetentionPolicy(config.NewStorageClient(userAgent).Buckets, d.Get("name").(string), res.Metageneration)
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
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	// Get the bucket and acl
	bucket := d.Get("name").(string)

	var res *storage.Bucket
	// There seems to be some eventual consistency errors in some cases, so we want to check a few times
	// to make sure it exists before moving on
	err = retryTimeDuration(func() (operr error) {
		var retryErr error
		res, retryErr = config.NewStorageClient(userAgent).Buckets.Get(bucket).Do()
		return retryErr
	}, d.Timeout(schema.TimeoutRead), isNotFoundRetryableError("bucket read"))

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Storage Bucket %q", d.Get("name").(string)))
	}
	log.Printf("[DEBUG] Read bucket %v at location %v\n\n", res.Name, res.SelfLink)

	return setStorageBucket(d, config, res, bucket, userAgent)
}

func resourceStorageBucketDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.userAgent)
	if err != nil {
		return err
	}

	// Get the bucket
	bucket := d.Get("name").(string)

	var listError, deleteObjectError error
	for deleteObjectError == nil {
		res, err := config.NewStorageClient(userAgent).Objects.List(bucket).Versions(true).Do()
		if err != nil {
			log.Printf("Error listing contents of bucket %s: %v", bucket, err)
			// If we can't list the contents, try deleting the bucket anyway in case it's empty
			listError = err
			break
		}

		if len(res.Items) == 0 {
			break // 0 items, bucket empty
		}

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

		if !d.Get("force_destroy").(bool) {
			deleteErr := fmt.Errorf("Error trying to delete bucket %s containing objects without `force_destroy` set to true", bucket)
			log.Printf("Error! %s : %s\n\n", bucket, deleteErr)
			return deleteErr
		}
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
				if err := config.NewStorageClient(userAgent).Objects.Delete(bucket, object.Name).Generation(object.Generation).Do(); err != nil {
					deleteObjectError = err
					log.Printf("[ERR] Failed to delete storage object %s: %s", object.Name, err)
				} else {
					log.Printf("[TRACE] Successfully deleted %s", object.Name)
				}
			})
		}

		// Wait for everything to finish.
		wp.StopWait()
	}

	// remove empty bucket
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		err := config.NewStorageClient(userAgent).Buckets.Delete(bucket).Do()
		if err == nil {
			return nil
		}
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 429 {
			return resource.RetryableError(gerr)
		}
		return resource.NonRetryableError(err)
	})
	if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 409 && strings.Contains(gerr.Message, "not empty") && listError != nil {
		return fmt.Errorf("could not delete non-empty bucket due to error when listing contents: %v", listError)
	}
	if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 409 && strings.Contains(gerr.Message, "not empty") && deleteObjectError != nil {
		return fmt.Errorf("could not delete non-empty bucket due to error when deleting contents: %v", deleteObjectError)
	}
	if err != nil {
		log.Printf("Error deleting bucket %s: %v", bucket, err)
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
		if err := d.Set("name", parts[0]); err != nil {
			return nil, fmt.Errorf("Error setting name: %s", err)
		}
	} else if len(parts) > 1 {
		if err := d.Set("project", parts[0]); err != nil {
			return nil, fmt.Errorf("Error setting project: %s", err)
		}
		if err := d.Set("name", parts[1]); err != nil {
			return nil, fmt.Errorf("Error setting name: %s", err)
		}
	}

	if err := d.Set("force_destroy", false); err != nil {
		return nil, fmt.Errorf("Error setting force_destroy: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}

func expandCors(configured []interface{}) []*storage.BucketCors {
	if len(configured) == 0 {
		return nil
	}
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

func expandBucketCustomPlacementConfig(configured interface{}) *storage.BucketCustomPlacementConfig {
	cfcs := configured.([]interface{})
	if len(cfcs) == 0 || cfcs[0] == nil {
		return nil
	}
	cfc := cfcs[0].(map[string]interface{})
	bucketcfc := &storage.BucketCustomPlacementConfig{
		DataLocations: expandBucketDataLocations(cfc["data_locations"]),
	}
	return bucketcfc
}

func flattenBucketCustomPlacementConfig(cfc *storage.BucketCustomPlacementConfig) []map[string]interface{} {
	customPlacementConfig := make([]map[string]interface{}, 0, 1)

	if cfc == nil {
		return customPlacementConfig
	}

	customPlacementConfig = append(customPlacementConfig, map[string]interface{}{
		"data_locations": cfc.DataLocations,
	})

	return customPlacementConfig
}

func expandBucketDataLocations(configured interface{}) []string {
	l := configured.(*schema.Set).List()

	req := make([]string, 0, len(l))
	for _, raw := range l {
		req = append(req, raw.(string))
	}
	return req
}

func expandBucketLogging(configured interface{}) *storage.BucketLogging {
	loggings := configured.([]interface{})
	if len(loggings) == 0 || loggings[0] == nil {
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
	if len(retentionPolicies) == 0 {
		return nil
	}
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

func expandBucketAutoclass(configured interface{}) *storage.BucketAutoclass {
	autoclassList := configured.([]interface{})
	if len(autoclassList) == 0 {
		return nil
	}

	autoclass := autoclassList[0].(map[string]interface{})

	bucketAutoclass := &storage.BucketAutoclass{}

	bucketAutoclass.Enabled = autoclass["enabled"].(bool)
	bucketAutoclass.ForceSendFields = append(bucketAutoclass.ForceSendFields, "Enabled")

	return bucketAutoclass
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

func flattenBucketAutoclass(bucketAutoclass *storage.BucketAutoclass) []map[string]interface{} {
	autoclassList := make([]map[string]interface{}, 0, 1)

	if bucketAutoclass == nil {
		return autoclassList
	}

	autoclass := map[string]interface{}{
		"enabled": bucketAutoclass.Enabled,
	}
	autoclassList = append(autoclassList, autoclass)
	return autoclassList
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
		"created_before":             condition.CreatedBefore,
		"matches_storage_class":      convertStringArrToInterface(condition.MatchesStorageClass),
		"num_newer_versions":         int(condition.NumNewerVersions),
		"custom_time_before":         condition.CustomTimeBefore,
		"days_since_custom_time":     int(condition.DaysSinceCustomTime),
		"days_since_noncurrent_time": int(condition.DaysSinceNoncurrentTime),
		"noncurrent_time_before":     condition.NoncurrentTimeBefore,
		"matches_prefix":             convertStringArrToInterface(condition.MatchesPrefix),
		"matches_suffix":             convertStringArrToInterface(condition.MatchesSuffix),
	}
	if condition.Age != nil {
		ruleCondition["age"] = int(*condition.Age)
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
	cfg := &storage.BucketIamConfiguration{
		ForceSendFields: []string{"UniformBucketLevelAccess"},
		UniformBucketLevelAccess: &storage.BucketIamConfigurationUniformBucketLevelAccess{
			Enabled:         d.Get("uniform_bucket_level_access").(bool),
			ForceSendFields: []string{"Enabled"},
		},
	}

	if v, ok := d.GetOk("public_access_prevention"); ok {
		cfg.PublicAccessPrevention = v.(string)
	}

	return cfg
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
		age := int64(v.(int))
		transformed.Age = &age
		transformed.ForceSendFields = append(transformed.ForceSendFields, "Age")
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

	if v, ok := condition["custom_time_before"]; ok {
		transformed.CustomTimeBefore = v.(string)
	}

	if v, ok := condition["days_since_custom_time"]; ok {
		transformed.DaysSinceCustomTime = int64(v.(int))
	}

	if v, ok := condition["days_since_noncurrent_time"]; ok {
		transformed.DaysSinceNoncurrentTime = int64(v.(int))
	}

	if v, ok := condition["noncurrent_time_before"]; ok {
		transformed.NoncurrentTimeBefore = v.(string)
	}

	if v, ok := condition["matches_prefix"]; ok {
		prefixes := v.([]interface{})
		transformedPrefixes := make([]string, 0, len(prefixes))

		for _, v := range prefixes {
			transformedPrefixes = append(transformedPrefixes, v.(string))
		}
		transformed.MatchesPrefix = transformedPrefixes
	}
	if v, ok := condition["matches_suffix"]; ok {
		suffixes := v.([]interface{})
		transformedSuffixes := make([]string, 0, len(suffixes))

		for _, v := range suffixes {
			transformedSuffixes = append(transformedSuffixes, v.(string))
		}
		transformed.MatchesSuffix = transformedSuffixes
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

	return hashcode(buf.String())
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

	if v, ok := m["days_since_custom_time"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}

	if v, ok := m["days_since_noncurrent_time"]; ok {
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

	if v, ok := m["matches_prefix"]; ok {
		matches_prefixes := v.([]interface{})
		for _, matches_prefix := range matches_prefixes {
			buf.WriteString(fmt.Sprintf("%s-", matches_prefix))
		}
	}
	if v, ok := m["matches_suffix"]; ok {
		matches_suffixes := v.([]interface{})
		for _, matches_suffix := range matches_suffixes {
			buf.WriteString(fmt.Sprintf("%s-", matches_suffix))
		}
	}

	return hashcode(buf.String())
}

func lockRetentionPolicy(bucketsService *storage.BucketsService, bucketName string, metageneration int64) error {
	lockPolicyCall := bucketsService.LockRetentionPolicy(bucketName, metageneration)
	if _, err := lockPolicyCall.Do(); err != nil {
		return err
	}

	return nil
}

// d.HasChange("lifecycle_rule") always returns true, giving false positives. This function detects changes
// to the list size or the actions/conditions of rules directly.
func detectLifecycleChange(d *schema.ResourceData) bool {
	if d.HasChange("lifecycle_rule.#") {
		return true
	}

	if l, ok := d.GetOk("lifecycle_rule"); ok {
		lifecycleRules := l.([]interface{})
		for i := range lifecycleRules {
			if d.HasChange(fmt.Sprintf("lifecycle_rule.%d.action", i)) || d.HasChange(fmt.Sprintf("lifecycle_rule.%d.condition", i)) {
				return true
			}
		}
	}

	return false
}

// Resource Read and DataSource Read both need to set attributes, but Data Sources don't support Timeouts
// so we pulled this portion out separately (https://github.com/hashicorp/terraform-provider-google/issues/11264)
func setStorageBucket(d *schema.ResourceData, config *Config, res *storage.Bucket, bucket, userAgent string) error {
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
		if err := d.Set("project", project); err != nil {
			return fmt.Errorf("Error setting project: %s", err)
		}
	}
	if d.Get("project") == "" {
		proj, err := config.NewComputeClient(userAgent).Projects.Get(strconv.FormatUint(res.ProjectNumber, 10)).Do()
		if err != nil {
			return err
		}
		log.Printf("[DEBUG] Bucket %v is in project number %v, which is project ID %s.\n", res.Name, res.ProjectNumber, proj.Name)
		if err := d.Set("project", proj.Name); err != nil {
			return fmt.Errorf("Error setting project: %s", err)
		}
	}

	// Update the bucket ID according to the resource ID
	if err := d.Set("self_link", res.SelfLink); err != nil {
		return fmt.Errorf("Error setting self_link: %s", err)
	}
	if err := d.Set("url", fmt.Sprintf("gs://%s", bucket)); err != nil {
		return fmt.Errorf("Error setting url: %s", err)
	}
	if err := d.Set("storage_class", res.StorageClass); err != nil {
		return fmt.Errorf("Error setting storage_class: %s", err)
	}
	if err := d.Set("encryption", flattenBucketEncryption(res.Encryption)); err != nil {
		return fmt.Errorf("Error setting encryption: %s", err)
	}
	if err := d.Set("location", res.Location); err != nil {
		return fmt.Errorf("Error setting location: %s", err)
	}
	if err := d.Set("cors", flattenCors(res.Cors)); err != nil {
		return fmt.Errorf("Error setting cors: %s", err)
	}
	if err := d.Set("default_event_based_hold", res.DefaultEventBasedHold); err != nil {
		return fmt.Errorf("Error setting default_event_based_hold: %s", err)
	}
	if err := d.Set("logging", flattenBucketLogging(res.Logging)); err != nil {
		return fmt.Errorf("Error setting logging: %s", err)
	}
	if err := d.Set("versioning", flattenBucketVersioning(res.Versioning)); err != nil {
		return fmt.Errorf("Error setting versioning: %s", err)
	}
	if err := d.Set("autoclass", flattenBucketAutoclass(res.Autoclass)); err != nil {
		return fmt.Errorf("Error setting autoclass: %s", err)
	}
	if err := d.Set("lifecycle_rule", flattenBucketLifecycle(res.Lifecycle)); err != nil {
		return fmt.Errorf("Error setting lifecycle_rule: %s", err)
	}
	if err := d.Set("labels", res.Labels); err != nil {
		return fmt.Errorf("Error setting labels: %s", err)
	}
	if err := d.Set("website", flattenBucketWebsite(res.Website)); err != nil {
		return fmt.Errorf("Error setting website: %s", err)
	}
	if err := d.Set("retention_policy", flattenBucketRetentionPolicy(res.RetentionPolicy)); err != nil {
		return fmt.Errorf("Error setting retention_policy: %s", err)
	}
	if err := d.Set("custom_placement_config", flattenBucketCustomPlacementConfig(res.CustomPlacementConfig)); err != nil {
		return fmt.Errorf("Error setting custom_placement_config: %s", err)
	}

	if res.IamConfiguration != nil && res.IamConfiguration.UniformBucketLevelAccess != nil {
		if err := d.Set("uniform_bucket_level_access", res.IamConfiguration.UniformBucketLevelAccess.Enabled); err != nil {
			return fmt.Errorf("Error setting uniform_bucket_level_access: %s", err)
		}
	} else {
		if err := d.Set("uniform_bucket_level_access", false); err != nil {
			return fmt.Errorf("Error setting uniform_bucket_level_access: %s", err)
		}
	}

	if res.IamConfiguration != nil && res.IamConfiguration.PublicAccessPrevention != "" {
		if err := d.Set("public_access_prevention", res.IamConfiguration.PublicAccessPrevention); err != nil {
			return fmt.Errorf("Error setting public_access_prevention: %s", err)
		}
	}

	if res.Billing == nil {
		if err := d.Set("requester_pays", nil); err != nil {
			return fmt.Errorf("Error setting requester_pays: %s", err)
		}
	} else {
		if err := d.Set("requester_pays", res.Billing.RequesterPays); err != nil {
			return fmt.Errorf("Error setting requester_pays: %s", err)
		}
	}

	d.SetId(res.Id)
	return nil
}
