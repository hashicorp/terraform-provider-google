// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storage

import (
	"context"
	"log"
	"math"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-google/google/verify"
)

func resourceStorageBucketV1() *schema.Resource {
	return &schema.Resource{
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceStorageBucketV0().CoreConfigSchema().ImpliedType(),
				Upgrade: ResourceStorageBucketStateUpgradeV0,
				Version: 0,
			},
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  `The name of the bucket.`,
				ValidateFunc: verify.ValidateGCSName,
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
				Type:         schema.TypeMap,
				ValidateFunc: labelKeyValidator,
				Optional:     true,
				Elem:         &schema.Schema{Type: schema.TypeString},
				Description:  `A set of key/value label pairs to assign to the bucket.`,
			},

			"terraform_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `The combination of labels configured directly on the resource and default labels configured on the provider.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"effective_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
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

			"project_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The project number of the project in which the resource belongs.`,
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
									"no_age": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `While set true, age value will be omitted.Required to set true when age is unset in the config file.`,
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
									"send_days_since_noncurrent_time_if_zero": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `While set true, days_since_noncurrent_time value will be sent in the request even for zero value of the field. This field is only useful for setting 0 value to the days_since_noncurrent_time field. It can be used alone or together with days_since_noncurrent_time.`,
									},
									"send_days_since_custom_time_if_zero": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `While set true, days_since_custom_time value will be sent in the request even for zero value of the field. This field is only useful for setting 0 value to the days_since_custom_time field. It can be used alone or together with days_since_custom_time.`,
									},
									"send_num_newer_versions_if_zero": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `While set true, num_newer_versions value will be sent in the request even for zero value of the field. This field is only useful for setting 0 value to the num_newer_versions field. It can be used alone or together with num_newer_versions.`,
									},
								},
							},
							Description: `The Lifecycle Rule's condition configuration.`,
						},
					},
				},
				Description: `The bucket's Lifecycle Rules configuration.`,
			},

			"enable_object_retention": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Enables each object in the bucket to have its own retention policy, which prevents deletion until stored for a specific length of time.`,
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `While set to true, autoclass automatically transitions objects in your bucket to appropriate storage classes based on each object's access pattern.`,
						},
						"terminal_storage_class": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The storage class that objects in the bucket eventually transition to if they are not read for a certain length of time. Supported values include: NEARLINE, ARCHIVE.`,
						},
					},
				},
				Description: `The bucket's autoclass configuration.`,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					_, n := d.GetChange(strings.TrimSuffix(k, ".#"))
					if !strings.HasSuffix(k, ".#") {
						return false
					}
					var l []interface{}
					if new == "1" && old == "0" {
						l = n.([]interface{})
						contents, ok := l[0].(map[string]interface{})
						if !ok {
							return false
						}
						if contents["enabled"] == false {
							return true
						}
					}
					if new == "0" && old == "1" {
						n := d.Get(strings.TrimSuffix(k, ".#"))
						l = n.([]interface{})
						contents := l[0].(map[string]interface{})
						if contents["enabled"] == false {
							return true
						}
					}
					return false
				},
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
								StateFunc: func(s interface{}) string {
									return strings.ToUpper(s.(string))
								},
							},
							Description: `The list of individual regions that comprise a dual-region bucket. See the docs for a list of acceptable regions. Note: If any of the data_locations changes, it will recreate the bucket.`,
						},
					},
				},
				Description: `The bucket's custom location configuration, which specifies the individual regions that comprise a dual-region bucket. If the bucket is designated a single or multi-region, the parameters are empty.`,
			},
			"rpo": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the RPO setting of bucket. If set 'ASYNC_TURBO', The Turbo Replication will be enabled for the dual-region bucket. Value 'DEFAULT' will set RPO setting to default. Turbo Replication is only for buckets in dual-regions.See the docs for more details.`,
			},
			"public_access_prevention": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Prevents public access to a bucket.`,
			},
			"soft_delete_policy": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: `The bucket's soft delete policy, which defines the period of time that soft-deleted objects will be retained, and cannot be permanently deleted. If it is not provided, by default Google Cloud Storage sets this to default soft delete policy`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"retention_duration_seconds": {
							Type:        schema.TypeInt,
							Default:     604800,
							Optional:    true,
							Description: `The duration in seconds that soft-deleted objects in the bucket will be retained and cannot be permanently deleted. Default value is 604800.`,
						},
						"effective_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Server-determined value that indicates the time from which the policy, or one with a greater retention, was effective. This value is in RFC 3339 format.`,
						},
					},
				},
			},
		},
		UseJSONNumber: true,
	}
}

func ResourceStorageBucketStateUpgradeV1(_ context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", rawState)
	if rawState["lifecycle_rule"] != nil {
		rawRules := rawState["lifecycle_rule"].([]interface{})
		for i, r := range rawRules {
			newRule := r.(map[string]interface{})
			if newRule["condition"] != nil {
				newCondition := newRule["condition"].([]interface{})[0].(map[string]interface{})
				newCondition["send_age_if_zero"] = true
				newRule["condition"].([]interface{})[0] = newCondition
			}
			rawState["lifecycle_rule"].([]interface{})[i] = newRule
		}
	}
	log.Printf("[DEBUG] Attributes after migration: %#v", rawState)
	return rawState, nil
}

func resourceStorageBucketV2() *schema.Resource {
	return &schema.Resource{
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceStorageBucketV0().CoreConfigSchema().ImpliedType(),
				Upgrade: ResourceStorageBucketStateUpgradeV0,
				Version: 0,
			},
			{
				Type:    resourceStorageBucketV1().CoreConfigSchema().ImpliedType(),
				Upgrade: ResourceStorageBucketStateUpgradeV1,
				Version: 1,
			},
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  `The name of the bucket.`,
				ValidateFunc: verify.ValidateGCSName,
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
				Type:         schema.TypeMap,
				ValidateFunc: labelKeyValidator,
				Optional:     true,
				Elem:         &schema.Schema{Type: schema.TypeString},
				Description:  `A set of key/value label pairs to assign to the bucket.`,
			},

			"terraform_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `The combination of labels configured directly on the resource and default labels configured on the provider.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"effective_labels": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `All of labels (key/value pairs) present on the resource in GCP, including the labels configured through Terraform, other clients and services.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
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

			"project_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The project number of the project in which the resource belongs.`,
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
									"no_age": {
										Type:        schema.TypeBool,
										Deprecated:  "`no_age` is deprecated and will be removed in a future major release. Use `send_age_if_zero` instead.",
										Optional:    true,
										Description: `While set true, age value will be omitted.Required to set true when age is unset in the config file.`,
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
									"send_age_if_zero": {
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     true,
										Description: `While set true, age value will be sent in the request even for zero value of the field. This field is only useful for setting 0 value to the age field. It can be used alone or together with age.`,
									},
									"send_days_since_noncurrent_time_if_zero": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `While set true, days_since_noncurrent_time value will be sent in the request even for zero value of the field. This field is only useful for setting 0 value to the days_since_noncurrent_time field. It can be used alone or together with days_since_noncurrent_time.`,
									},
									"send_days_since_custom_time_if_zero": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `While set true, days_since_custom_time value will be sent in the request even for zero value of the field. This field is only useful for setting 0 value to the days_since_custom_time field. It can be used alone or together with days_since_custom_time.`,
									},
									"send_num_newer_versions_if_zero": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `While set true, num_newer_versions value will be sent in the request even for zero value of the field. This field is only useful for setting 0 value to the num_newer_versions field. It can be used alone or together with num_newer_versions.`,
									},
								},
							},
							Description: `The Lifecycle Rule's condition configuration.`,
						},
					},
				},
				Description: `The bucket's Lifecycle Rules configuration.`,
			},

			"enable_object_retention": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Enables each object in the bucket to have its own retention policy, which prevents deletion until stored for a specific length of time.`,
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `While set to true, autoclass automatically transitions objects in your bucket to appropriate storage classes based on each object's access pattern.`,
						},
						"terminal_storage_class": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The storage class that objects in the bucket eventually transition to if they are not read for a certain length of time. Supported values include: NEARLINE, ARCHIVE.`,
						},
					},
				},
				Description: `The bucket's autoclass configuration.`,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					_, n := d.GetChange(strings.TrimSuffix(k, ".#"))
					if !strings.HasSuffix(k, ".#") {
						return false
					}
					var l []interface{}
					if new == "1" && old == "0" {
						l = n.([]interface{})
						contents, ok := l[0].(map[string]interface{})
						if !ok {
							return false
						}
						if contents["enabled"] == false {
							return true
						}
					}
					if new == "0" && old == "1" {
						n := d.Get(strings.TrimSuffix(k, ".#"))
						l = n.([]interface{})
						contents := l[0].(map[string]interface{})
						if contents["enabled"] == false {
							return true
						}
					}
					return false
				},
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
								StateFunc: func(s interface{}) string {
									return strings.ToUpper(s.(string))
								},
							},
							Description: `The list of individual regions that comprise a dual-region bucket. See the docs for a list of acceptable regions. Note: If any of the data_locations changes, it will recreate the bucket.`,
						},
					},
				},
				Description: `The bucket's custom location configuration, which specifies the individual regions that comprise a dual-region bucket. If the bucket is designated a single or multi-region, the parameters are empty.`,
			},
			"rpo": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the RPO setting of bucket. If set 'ASYNC_TURBO', The Turbo Replication will be enabled for the dual-region bucket. Value 'DEFAULT' will set RPO setting to default. Turbo Replication is only for buckets in dual-regions.See the docs for more details.`,
			},
			"public_access_prevention": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Prevents public access to a bucket.`,
			},
			"soft_delete_policy": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: `The bucket's soft delete policy, which defines the period of time that soft-deleted objects will be retained, and cannot be permanently deleted. If it is not provided, by default Google Cloud Storage sets this to default soft delete policy`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"retention_duration_seconds": {
							Type:        schema.TypeInt,
							Default:     604800,
							Optional:    true,
							Description: `The duration in seconds that soft-deleted objects in the bucket will be retained and cannot be permanently deleted. Default value is 604800.`,
						},
						"effective_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Server-determined value that indicates the time from which the policy, or one with a greater retention, was effective. This value is in RFC 3339 format.`,
						},
					},
				},
			},
		},
		UseJSONNumber: true,
	}
}

func ResourceStorageBucketStateUpgradeV2(_ context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	log.Printf("[DEBUG] Attributes before migration: %#v", rawState)
	if rawState["lifecycle_rule"] != nil {
		rawRules := rawState["lifecycle_rule"].([]interface{})
		for i, r := range rawRules {
			newRule := r.(map[string]interface{})
			if newRule["condition"] != nil {
				newCondition := newRule["condition"].([]interface{})[0].(map[string]interface{})
				newCondition["send_age_if_zero"] = false
				newRule["condition"].([]interface{})[0] = newCondition
			}
			rawState["lifecycle_rule"].([]interface{})[i] = newRule
		}
	}
	log.Printf("[DEBUG] Attributes after migration: %#v", rawState)
	return rawState, nil
}
