// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package storagetransfer

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"google.golang.org/api/storagetransfer/v1"
)

var (
	objectConditionsKeys = []string{
		"transfer_spec.0.object_conditions.0.min_time_elapsed_since_last_modification",
		"transfer_spec.0.object_conditions.0.max_time_elapsed_since_last_modification",
		"transfer_spec.0.object_conditions.0.include_prefixes",
		"transfer_spec.0.object_conditions.0.exclude_prefixes",
		"transfer_spec.0.object_conditions.0.last_modified_since",
		"transfer_spec.0.object_conditions.0.last_modified_before",
	}

	transferOptionsKeys = []string{
		"transfer_spec.0.transfer_options.0.overwrite_objects_already_existing_in_sink",
		"transfer_spec.0.transfer_options.0.delete_objects_unique_in_sink",
		"transfer_spec.0.transfer_options.0.delete_objects_from_source_after_transfer",
		"transfer_spec.0.transfer_options.0.overwrite_when",
	}

	transferSpecDataSourceKeys = []string{
		"transfer_spec.0.gcs_data_source",
		"transfer_spec.0.aws_s3_data_source",
		"transfer_spec.0.http_data_source",
		"transfer_spec.0.azure_blob_storage_data_source",
		"transfer_spec.0.posix_data_source",
	}
	transferSpecDataSinkKeys = []string{
		"transfer_spec.0.gcs_data_sink",
		"transfer_spec.0.posix_data_sink",
	}
	awsS3AuthKeys = []string{
		"transfer_spec.0.aws_s3_data_source.0.aws_access_key",
		"transfer_spec.0.aws_s3_data_source.0.role_arn",
	}
)

func ResourceStorageTransferJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageTransferJobCreate,
		Read:   resourceStorageTransferJobRead,
		Update: resourceStorageTransferJobUpdate,
		Delete: resourceStorageTransferJobDelete,
		Importer: &schema.ResourceImporter{
			State: resourceStorageTransferJobStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the Transfer Job.`,
			},
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 1024),
				Description:  `Unique description to identify the Transfer Job.`,
			},
			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The project in which the resource belongs. If it is not provided, the provider project is used.`,
			},
			"transfer_spec": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_conditions": objectConditionsSchema(),
						"transfer_options":  transferOptionsSchema(),
						"source_agent_pool_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: `Specifies the agent pool name associated with the posix data source. When unspecified, the default name is used.`,
						},
						"sink_agent_pool_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: `Specifies the agent pool name associated with the posix data source. When unspecified, the default name is used.`,
						},
						"gcs_data_sink": {
							Type:         schema.TypeList,
							Optional:     true,
							MaxItems:     1,
							Elem:         gcsDataSchema(),
							ExactlyOneOf: transferSpecDataSinkKeys,
							Description:  `A Google Cloud Storage data sink.`,
						},
						"posix_data_sink": {
							Type:         schema.TypeList,
							Optional:     true,
							MaxItems:     1,
							Elem:         posixDataSchema(),
							ExactlyOneOf: transferSpecDataSinkKeys,
							Description:  `A POSIX filesystem data sink.`,
						},
						"gcs_data_source": {
							Type:         schema.TypeList,
							Optional:     true,
							MaxItems:     1,
							Elem:         gcsDataSchema(),
							ExactlyOneOf: transferSpecDataSourceKeys,
							Description:  `A Google Cloud Storage data source.`,
						},
						"aws_s3_data_source": {
							Type:         schema.TypeList,
							Optional:     true,
							MaxItems:     1,
							Elem:         awsS3DataSchema(),
							ExactlyOneOf: transferSpecDataSourceKeys,
							Description:  `An AWS S3 data source.`,
						},
						"http_data_source": {
							Type:         schema.TypeList,
							Optional:     true,
							MaxItems:     1,
							Elem:         httpDataSchema(),
							ExactlyOneOf: transferSpecDataSourceKeys,
							Description:  `A HTTP URL data source.`,
						},
						"posix_data_source": {
							Type:         schema.TypeList,
							Optional:     true,
							MaxItems:     1,
							Elem:         posixDataSchema(),
							ExactlyOneOf: transferSpecDataSourceKeys,
							Description:  `A POSIX filesystem data source.`,
						},
						"azure_blob_storage_data_source": {
							Type:         schema.TypeList,
							Optional:     true,
							MaxItems:     1,
							Elem:         azureBlobStorageDataSchema(),
							ExactlyOneOf: transferSpecDataSourceKeys,
							Description:  `An Azure Blob Storage data source.`,
						},
					},
				},
				Description: `Transfer specification.`,
			},
			"notification_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pubsub_topic": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The Topic.name of the Pub/Sub topic to which to publish notifications.`,
						},
						"event_types": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validation.StringInSlice([]string{"TRANSFER_OPERATION_SUCCESS", "TRANSFER_OPERATION_FAILED", "TRANSFER_OPERATION_ABORTED"}, false),
							},
							Description: `Event types for which a notification is desired. If empty, send notifications for all event types. The valid types are "TRANSFER_OPERATION_SUCCESS", "TRANSFER_OPERATION_FAILED", "TRANSFER_OPERATION_ABORTED".`,
						},
						"payload_format": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"NONE", "JSON"}, false),
							Description:  `The desired format of the notification message payloads. One of "NONE" or "JSON".`,
						},
					},
				},
				Description: `Notification configuration.`,
			},
			"schedule": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule_start_date": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Elem:        dateObjectSchema(),
							Description: `The first day the recurring transfer is scheduled to run. If schedule_start_date is in the past, the transfer will run for the first time on the following day.`,
						},
						"schedule_end_date": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Elem:        dateObjectSchema(),
							Description: `The last day the recurring transfer will be run. If schedule_end_date is the same as schedule_start_date, the transfer will be executed only once.`,
						},
						"start_time_of_day": {
							Type:             schema.TypeList,
							Optional:         true,
							MaxItems:         1,
							Elem:             timeObjectSchema(),
							DiffSuppressFunc: diffSuppressEmptyStartTimeOfDay,
							Description:      `The time in UTC at which the transfer will be scheduled to start in a day. Transfers may start later than this time. If not specified, recurring and one-time transfers that are scheduled to run today will run immediately; recurring transfers that are scheduled to run on a future date will start at approximately midnight UTC on that date. Note that when configuring a transfer with the Cloud Platform Console, the transfer's start time in a day is specified in your local timezone.`,
						},
						"repeat_interval": {
							Type:         schema.TypeString,
							ValidateFunc: verify.ValidateDuration(),
							Optional:     true,
							Description:  `Interval between the start of each scheduled transfer. If unspecified, the default value is 24 hours. This value may not be less than 1 hour. A duration in seconds with up to nine fractional digits, terminated by 's'. Example: "3.5s".`,
							Default:      "86400s",
						},
					},
				},
				Description: `Schedule specification defining when the Transfer Job should be scheduled to start, end and what time to run.`,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ENABLED",
				ValidateFunc: validation.StringInSlice([]string{"ENABLED", "DISABLED", "DELETED"}, false),
				Description:  `Status of the job. Default: ENABLED. NOTE: The effect of the new job status takes place during a subsequent job run. For example, if you change the job status from ENABLED to DISABLED, and an operation spawned by the transfer is running, the status change would not affect the current operation.`,
			},
			"creation_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `When the Transfer Job was created.`,
			},
			"last_modification_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `When the Transfer Job was last modified.`,
			},
			"deletion_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `When the Transfer Job was deleted.`,
			},
		},
		UseJSONNumber: true,
	}
}

func objectConditionsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"min_time_elapsed_since_last_modification": {
					Type:         schema.TypeString,
					ValidateFunc: verify.ValidateDuration(),
					Optional:     true,
					AtLeastOneOf: objectConditionsKeys,
					Description:  `A duration in seconds with up to nine fractional digits, terminated by 's'. Example: "3.5s".`,
				},
				"max_time_elapsed_since_last_modification": {
					Type:         schema.TypeString,
					ValidateFunc: verify.ValidateDuration(),
					Optional:     true,
					AtLeastOneOf: objectConditionsKeys,
					Description:  `A duration in seconds with up to nine fractional digits, terminated by 's'. Example: "3.5s".`,
				},
				"include_prefixes": {
					Type:         schema.TypeList,
					Optional:     true,
					AtLeastOneOf: objectConditionsKeys,
					Elem: &schema.Schema{
						MaxItems: 1000,
						Type:     schema.TypeString,
					},
					Description: `If include_refixes is specified, objects that satisfy the object conditions must have names that start with one of the include_prefixes and that do not start with any of the exclude_prefixes. If include_prefixes is not specified, all objects except those that have names starting with one of the exclude_prefixes must satisfy the object conditions.`,
				},
				"exclude_prefixes": {
					Type:         schema.TypeList,
					Optional:     true,
					AtLeastOneOf: objectConditionsKeys,
					Elem: &schema.Schema{
						MaxItems: 1000,
						Type:     schema.TypeString,
					},
					Description: `exclude_prefixes must follow the requirements described for include_prefixes.`,
				},
				"last_modified_since": {
					Type:         schema.TypeString,
					ValidateFunc: verify.ValidateRFC3339Date,
					Optional:     true,
					AtLeastOneOf: objectConditionsKeys,
					Description:  `If specified, only objects with a "last modification time" on or after this timestamp and objects that don't have a "last modification time" are transferred. A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine fractional digits. Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".`,
				},
				"last_modified_before": {
					Type:         schema.TypeString,
					ValidateFunc: verify.ValidateRFC3339Date,
					Optional:     true,
					AtLeastOneOf: objectConditionsKeys,
					Description:  `If specified, only objects with a "last modification time" before this timestamp and objects that don't have a "last modification time" are transferred. A timestamp in RFC3339 UTC "Zulu" format, with nanosecond resolution and up to nine fractional digits. Examples: "2014-10-02T15:01:23Z" and "2014-10-02T15:01:23.045123456Z".`,
				},
			},
		},
		Description: `Only objects that satisfy these object conditions are included in the set of data source and data sink objects. Object conditions based on objects' last_modification_time do not exclude objects in a data sink.`,
	}
}

func transferOptionsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"overwrite_objects_already_existing_in_sink": {
					Type:         schema.TypeBool,
					Optional:     true,
					AtLeastOneOf: transferOptionsKeys,
					Description:  `Whether overwriting objects that already exist in the sink is allowed.`,
				},
				"delete_objects_unique_in_sink": {
					Type:          schema.TypeBool,
					Optional:      true,
					AtLeastOneOf:  transferOptionsKeys,
					ConflictsWith: []string{"transfer_spec.transfer_options.delete_objects_from_source_after_transfer"},
					Description:   `Whether objects that exist only in the sink should be deleted. Note that this option and delete_objects_from_source_after_transfer are mutually exclusive.`,
				},
				"delete_objects_from_source_after_transfer": {
					Type:          schema.TypeBool,
					Optional:      true,
					AtLeastOneOf:  transferOptionsKeys,
					ConflictsWith: []string{"transfer_spec.transfer_options.delete_objects_unique_in_sink"},
					Description:   `Whether objects should be deleted from the source after they are transferred to the sink. Note that this option and delete_objects_unique_in_sink are mutually exclusive.`,
				},
				"overwrite_when": {
					Type:         schema.TypeString,
					Optional:     true,
					AtLeastOneOf: transferOptionsKeys,
					ValidateFunc: validation.StringInSlice([]string{"DIFFERENT", "NEVER", "ALWAYS"}, false),
					Description:  `When to overwrite objects that already exist in the sink. If not set, overwrite behavior is determined by overwriteObjectsAlreadyExistingInSink.`,
				},
			},
		},
		Description: `Characteristics of how to treat files from datasource and sink during job. If the option delete_objects_unique_in_sink is true, object conditions based on objects' last_modification_time are ignored and do not exclude objects in a data source or a data sink.`,
	}
}

func timeObjectSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"hours": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 24),
				Description:  `Hours of day in 24 hour format. Should be from 0 to 23.`,
			},
			"minutes": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 59),
				Description:  `Minutes of hour of day. Must be from 0 to 59.`,
			},
			"seconds": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 60),
				Description:  `Seconds of minutes of the time. Must normally be from 0 to 59.`,
			},
			"nanos": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 999999999),
				Description:  `Fractions of seconds in nanoseconds. Must be from 0 to 999,999,999.`,
			},
		},
	}
}

func dateObjectSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"year": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 9999),
				Description:  `Year of date. Must be from 1 to 9999.`,
			},

			"month": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 12),
				Description:  `Month of year. Must be from 1 to 12.`,
			},

			"day": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 31),
				Description:  `Day of month. Must be from 1 to 31 and valid for the year and month.`,
			},
		},
	}
}

func gcsDataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: `Google Cloud Storage bucket name.`,
			},
			"path": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: `Google Cloud Storage path in bucket to transfer`,
			},
		},
	}
}

func awsS3DataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: `S3 Bucket name.`,
			},
			"path": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: `S3 Bucket path in bucket to transfer.`,
			},
			"aws_access_key": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key_id": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: `AWS Key ID.`,
						},
						"secret_access_key": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: `AWS Secret Access Key.`,
						},
					},
				},
				ExactlyOneOf: awsS3AuthKeys,
				Description:  `AWS credentials block.`,
			},
			"role_arn": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: awsS3AuthKeys,
				Description:  `The Amazon Resource Name (ARN) of the role to support temporary credentials via 'AssumeRoleWithWebIdentity'. For more information about ARNs, see [IAM ARNs](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html#identifiers-arns). When a role ARN is provided, Transfer Service fetches temporary credentials for the session using a 'AssumeRoleWithWebIdentity' call for the provided role using the [GoogleServiceAccount][] for this project.`,
			},
		},
	}
}

func httpDataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"list_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The URL that points to the file that stores the object list entries. This file must allow public access. Currently, only URLs with HTTP and HTTPS schemes are supported.`,
			},
		},
	}
}

func posixDataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"root_directory": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Root directory path to the filesystem.`,
			},
		},
	}
}

func azureBlobStorageDataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"storage_account": {
				Required:    true,
				Type:        schema.TypeString,
				Description: `The name of the Azure Storage account.`,
			},
			"container": {
				Required:    true,
				Type:        schema.TypeString,
				Description: `The container to transfer from the Azure Storage account.`,
			},
			"path": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: `Root path to transfer objects. Must be an empty string or full path name that ends with a '/'. This field is treated as an object prefix. As such, it should generally not begin with a '/'.`,
			},
			"azure_credentials": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sas_token": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							Description: `Azure shared access signature.`,
						},
					},
				},
				Description: ` Credentials used to authenticate API requests to Azure.`,
			},
		},
	}
}

func diffSuppressEmptyStartTimeOfDay(k, old, new string, d *schema.ResourceData) bool {
	return k == "schedule.0.start_time_of_day.#" && old == "1" && new == "0"
}

func resourceStorageTransferJobCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	transferJob := &storagetransfer.TransferJob{
		Description:        d.Get("description").(string),
		ProjectId:          project,
		Status:             d.Get("status").(string),
		Schedule:           expandTransferSchedules(d.Get("schedule").([]interface{})),
		TransferSpec:       expandTransferSpecs(d.Get("transfer_spec").([]interface{})),
		NotificationConfig: expandTransferJobNotificationConfig(d.Get("notification_config").([]interface{})),
	}

	var res *storagetransfer.TransferJob

	err = transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			res, err = config.NewStorageTransferClient(userAgent).TransferJobs.Create(transferJob).Do()
			return err
		},
	})

	if err != nil {
		fmt.Printf("Error creating transfer job %v: %v", transferJob, err)
		return err
	}

	if err := d.Set("name", res.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}

	name := tpgresource.GetResourceNameFromSelfLink(res.Name)
	d.SetId(fmt.Sprintf("%s/%s", project, name))

	return resourceStorageTransferJobRead(d, meta)
}

func resourceStorageTransferJobRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	res, err := config.NewStorageTransferClient(userAgent).TransferJobs.Get(name, project).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Transfer Job %q", name))
	}

	if res.Status == "DELETED" {
		d.SetId("")
		return nil
	}

	if err := d.Set("project", res.ProjectId); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}
	if err := d.Set("description", res.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("status", res.Status); err != nil {
		return fmt.Errorf("Error setting status: %s", err)
	}
	if err := d.Set("last_modification_time", res.LastModificationTime); err != nil {
		return fmt.Errorf("Error setting last_modification_time: %s", err)
	}
	if err := d.Set("creation_time", res.CreationTime); err != nil {
		return fmt.Errorf("Error setting creation_time: %s", err)
	}
	if err := d.Set("deletion_time", res.DeletionTime); err != nil {
		return fmt.Errorf("Error setting deletion_time: %s", err)
	}

	err = d.Set("schedule", flattenTransferSchedule(res.Schedule))
	if err != nil {
		return err
	}

	err = d.Set("transfer_spec", flattenTransferSpec(res.TransferSpec, d))
	if err != nil {
		return err
	}

	err = d.Set("notification_config", flattenTransferJobNotificationConfig(res.NotificationConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceStorageTransferJobUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	transferJob := &storagetransfer.TransferJob{}
	fieldMask := []string{}

	if d.HasChange("description") {
		fieldMask = append(fieldMask, "description")
		if v, ok := d.GetOk("description"); ok {
			transferJob.Description = v.(string)
		}
	}

	if d.HasChange("status") {
		fieldMask = append(fieldMask, "status")
		if v, ok := d.GetOk("status"); ok {
			transferJob.Status = v.(string)
		}
	}

	if d.HasChange("schedule") {
		fieldMask = append(fieldMask, "schedule")
		if v, ok := d.GetOk("schedule"); ok {
			transferJob.Schedule = expandTransferSchedules(v.([]interface{}))
		}
	}

	if d.HasChange("transfer_spec") {
		fieldMask = append(fieldMask, "transfer_spec")
		if v, ok := d.GetOk("transfer_spec"); ok {
			transferJob.TransferSpec = expandTransferSpecs(v.([]interface{}))
		}
	}

	if d.HasChange("notification_config") {
		fieldMask = append(fieldMask, "notification_config")
		if v, ok := d.GetOk("notification_config"); ok {
			transferJob.NotificationConfig = expandTransferJobNotificationConfig(v.([]interface{}))
		} else {
			transferJob.NotificationConfig = nil
		}
	}

	if len(fieldMask) == 0 {
		return nil
	}

	updateRequest := &storagetransfer.UpdateTransferJobRequest{
		ProjectId:   project,
		TransferJob: transferJob,
	}

	updateRequest.UpdateTransferJobFieldMask = strings.Join(fieldMask, ",")

	res, err := config.NewStorageTransferClient(userAgent).TransferJobs.Patch(d.Get("name").(string), updateRequest).Do()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Patched transfer job: %v\n\n", res.Name)
	return nil
}

func resourceStorageTransferJobDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	transferJobName := d.Get("name").(string)

	transferJob := &storagetransfer.TransferJob{
		Status: "DELETED",
	}

	fieldMask := "status"

	updateRequest := &storagetransfer.UpdateTransferJobRequest{
		ProjectId:   project,
		TransferJob: transferJob,
	}

	updateRequest.UpdateTransferJobFieldMask = fieldMask

	// Update transfer job with status set to DELETE
	log.Printf("[DEBUG] Setting status to DELETE for: %v\n\n", transferJobName)
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err := config.NewStorageTransferClient(userAgent).TransferJobs.Patch(transferJobName, updateRequest).Do()
		if err != nil {
			return resource.RetryableError(err)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error deleting transfer job %v: %v\n\n", transferJob, err)
		return err
	}

	log.Printf("[DEBUG] Deleted transfer job %v\n\n", transferJob)

	return nil
}

func resourceStorageTransferJobStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	switch len(parts) {
	case 2:
		if err := d.Set("project", parts[0]); err != nil {
			return nil, fmt.Errorf("Error setting project: %s", err)
		}
		if err := d.Set("name", fmt.Sprintf("transferJobs/%s", parts[1])); err != nil {
			return nil, fmt.Errorf("Error setting name: %s", err)
		}
	default:
		return nil, fmt.Errorf("Invalid transfer job specifier. Expecting {projectId}/{transferJobName}")
	}
	return []*schema.ResourceData{d}, nil
}

func expandDates(dates []interface{}) *storagetransfer.Date {
	if len(dates) == 0 || dates[0] == nil {
		return nil
	}

	dateMap := dates[0].(map[string]interface{})
	date := &storagetransfer.Date{}
	if v, ok := dateMap["day"]; ok {
		date.Day = int64(v.(int))
	}

	if v, ok := dateMap["month"]; ok {
		date.Month = int64(v.(int))
	}

	if v, ok := dateMap["year"]; ok {
		date.Year = int64(v.(int))
	}

	log.Printf("[DEBUG] not nil date: %#v", dates)
	return date
}

func flattenDate(date *storagetransfer.Date) []map[string]interface{} {
	data := map[string]interface{}{
		"year":  date.Year,
		"month": date.Month,
		"day":   date.Day,
	}

	return []map[string]interface{}{data}
}

func expandTimeOfDays(times []interface{}) *storagetransfer.TimeOfDay {
	if len(times) == 0 || times[0] == nil {
		return nil
	}

	timeMap := times[0].(map[string]interface{})
	time := &storagetransfer.TimeOfDay{}
	if v, ok := timeMap["hours"]; ok {
		time.Hours = int64(v.(int))
	}

	if v, ok := timeMap["minutes"]; ok {
		time.Minutes = int64(v.(int))
	}

	if v, ok := timeMap["seconds"]; ok {
		time.Seconds = int64(v.(int))
	}

	if v, ok := timeMap["nanos"]; ok {
		time.Nanos = int64(v.(int))
	}

	return time
}

func flattenTimeOfDay(timeOfDay *storagetransfer.TimeOfDay) []map[string]interface{} {
	data := map[string]interface{}{
		"hours":   timeOfDay.Hours,
		"minutes": timeOfDay.Minutes,
		"seconds": timeOfDay.Seconds,
		"nanos":   timeOfDay.Nanos,
	}

	return []map[string]interface{}{data}
}

func expandTransferSchedules(transferSchedules []interface{}) *storagetransfer.Schedule {
	if len(transferSchedules) == 0 || transferSchedules[0] == nil {
		return nil
	}

	schedule := transferSchedules[0].(map[string]interface{})
	return &storagetransfer.Schedule{
		ScheduleStartDate: expandDates(schedule["schedule_start_date"].([]interface{})),
		ScheduleEndDate:   expandDates(schedule["schedule_end_date"].([]interface{})),
		StartTimeOfDay:    expandTimeOfDays(schedule["start_time_of_day"].([]interface{})),
		RepeatInterval:    schedule["repeat_interval"].(string),
	}
}

func flattenTransferSchedule(transferSchedule *storagetransfer.Schedule) []map[string]interface{} {
	if transferSchedule == nil || reflect.DeepEqual(transferSchedule, &storagetransfer.Schedule{}) {
		return nil
	}

	data := map[string]interface{}{
		"schedule_start_date": flattenDate(transferSchedule.ScheduleStartDate),
	}

	if transferSchedule.ScheduleEndDate != nil {
		data["schedule_end_date"] = flattenDate(transferSchedule.ScheduleEndDate)
	}

	if transferSchedule.StartTimeOfDay != nil {
		data["start_time_of_day"] = flattenTimeOfDay(transferSchedule.StartTimeOfDay)
	}

	if transferSchedule.RepeatInterval != "" {
		data["repeat_interval"] = transferSchedule.RepeatInterval
	}

	return []map[string]interface{}{data}
}

func expandGcsData(gcsDatas []interface{}) *storagetransfer.GcsData {
	if len(gcsDatas) == 0 || gcsDatas[0] == nil {
		return nil
	}

	gcsData := gcsDatas[0].(map[string]interface{})
	var apiData = &storagetransfer.GcsData{
		BucketName: gcsData["bucket_name"].(string),
	}
	var path = gcsData["path"].(string)
	apiData.Path = path

	return apiData
}

func flattenGcsData(gcsData *storagetransfer.GcsData) []map[string]interface{} {
	data := map[string]interface{}{
		"bucket_name": gcsData.BucketName,
		"path":        gcsData.Path,
	}
	return []map[string]interface{}{data}
}

func expandAwsAccessKeys(awsAccessKeys []interface{}) *storagetransfer.AwsAccessKey {
	if len(awsAccessKeys) == 0 || awsAccessKeys[0] == nil {
		return nil
	}

	awsAccessKey := awsAccessKeys[0].(map[string]interface{})
	return &storagetransfer.AwsAccessKey{
		AccessKeyId:     awsAccessKey["access_key_id"].(string),
		SecretAccessKey: awsAccessKey["secret_access_key"].(string),
	}
}

func flattenAwsAccessKeys(d *schema.ResourceData) []map[string]interface{} {
	data := map[string]interface{}{
		"access_key_id":     d.Get("transfer_spec.0.aws_s3_data_source.0.aws_access_key.0.access_key_id"),
		"secret_access_key": d.Get("transfer_spec.0.aws_s3_data_source.0.aws_access_key.0.secret_access_key"),
	}

	return []map[string]interface{}{data}
}

func expandAwsS3Data(awsS3Datas []interface{}) *storagetransfer.AwsS3Data {
	if len(awsS3Datas) == 0 || awsS3Datas[0] == nil {
		return nil
	}

	awsS3Data := awsS3Datas[0].(map[string]interface{})
	return &storagetransfer.AwsS3Data{
		BucketName:   awsS3Data["bucket_name"].(string),
		AwsAccessKey: expandAwsAccessKeys(awsS3Data["aws_access_key"].([]interface{})),
		RoleArn:      awsS3Data["role_arn"].(string),
		Path:         awsS3Data["path"].(string),
	}
}

func flattenAwsS3Data(awsS3Data *storagetransfer.AwsS3Data, d *schema.ResourceData) []map[string]interface{} {
	data := map[string]interface{}{
		"bucket_name": awsS3Data.BucketName,
		"path":        awsS3Data.Path,
		"role_arn":    awsS3Data.RoleArn,
	}
	if awsS3Data.AwsAccessKey != nil {
		data["aws_access_key"] = flattenAwsAccessKeys(d)
	}

	return []map[string]interface{}{data}
}

func expandHttpData(httpDatas []interface{}) *storagetransfer.HttpData {
	if len(httpDatas) == 0 || httpDatas[0] == nil {
		return nil
	}

	httpData := httpDatas[0].(map[string]interface{})
	return &storagetransfer.HttpData{
		ListUrl: httpData["list_url"].(string),
	}
}

func flattenHttpData(httpData *storagetransfer.HttpData) []map[string]interface{} {
	data := map[string]interface{}{
		"list_url": httpData.ListUrl,
	}

	return []map[string]interface{}{data}
}

func expandPosixData(posixDatas []interface{}) *storagetransfer.PosixFilesystem {
	if len(posixDatas) == 0 || posixDatas[0] == nil {
		return nil
	}

	posixData := posixDatas[0].(map[string]interface{})
	return &storagetransfer.PosixFilesystem{
		RootDirectory: posixData["root_directory"].(string),
	}
}

func flattenPosixData(posixData *storagetransfer.PosixFilesystem) []map[string]interface{} {
	data := map[string]interface{}{
		"root_directory": posixData.RootDirectory,
	}

	return []map[string]interface{}{data}
}

func expandAzureCredentials(azureCredentials []interface{}) *storagetransfer.AzureCredentials {
	if len(azureCredentials) == 0 || azureCredentials[0] == nil {
		return nil
	}

	azureCredential := azureCredentials[0].(map[string]interface{})
	return &storagetransfer.AzureCredentials{
		SasToken: azureCredential["sas_token"].(string),
	}
}

func flattenAzureCredentials(d *schema.ResourceData) []map[string]interface{} {
	data := map[string]interface{}{
		"sas_token": d.Get("transfer_spec.0.azure_blob_storage_data_source.0.azure_credentials.0.sas_token"),
	}

	return []map[string]interface{}{data}
}

func expandAzureBlobStorageData(azureBlobStorageDatas []interface{}) *storagetransfer.AzureBlobStorageData {
	if len(azureBlobStorageDatas) == 0 || azureBlobStorageDatas[0] == nil {
		return nil
	}

	azureBlobStorageData := azureBlobStorageDatas[0].(map[string]interface{})

	return &storagetransfer.AzureBlobStorageData{
		Container:        azureBlobStorageData["container"].(string),
		Path:             azureBlobStorageData["path"].(string),
		StorageAccount:   azureBlobStorageData["storage_account"].(string),
		AzureCredentials: expandAzureCredentials(azureBlobStorageData["azure_credentials"].([]interface{})),
	}
}

func flattenAzureBlobStorageData(azureBlobStorageData *storagetransfer.AzureBlobStorageData, d *schema.ResourceData) []map[string]interface{} {
	data := map[string]interface{}{
		"container":         azureBlobStorageData.Container,
		"path":              azureBlobStorageData.Path,
		"storage_account":   azureBlobStorageData.StorageAccount,
		"azure_credentials": flattenAzureCredentials(d),
	}

	return []map[string]interface{}{data}
}

func expandObjectConditions(conditions []interface{}) *storagetransfer.ObjectConditions {
	if len(conditions) == 0 || conditions[0] == nil {
		return nil
	}

	condition := conditions[0].(map[string]interface{})
	return &storagetransfer.ObjectConditions{
		ExcludePrefixes:                     tpgresource.ConvertStringArr(condition["exclude_prefixes"].([]interface{})),
		IncludePrefixes:                     tpgresource.ConvertStringArr(condition["include_prefixes"].([]interface{})),
		MaxTimeElapsedSinceLastModification: condition["max_time_elapsed_since_last_modification"].(string),
		MinTimeElapsedSinceLastModification: condition["min_time_elapsed_since_last_modification"].(string),
		LastModifiedSince:                   condition["last_modified_since"].(string),
		LastModifiedBefore:                  condition["last_modified_before"].(string),
	}
}

func flattenObjectCondition(condition *storagetransfer.ObjectConditions) []map[string]interface{} {
	data := map[string]interface{}{
		"exclude_prefixes":                         condition.ExcludePrefixes,
		"include_prefixes":                         condition.IncludePrefixes,
		"max_time_elapsed_since_last_modification": condition.MaxTimeElapsedSinceLastModification,
		"min_time_elapsed_since_last_modification": condition.MinTimeElapsedSinceLastModification,
		"last_modified_since":                      condition.LastModifiedSince,
		"last_modified_before":                     condition.LastModifiedBefore,
	}
	return []map[string]interface{}{data}
}

func expandTransferOptions(options []interface{}) *storagetransfer.TransferOptions {
	if len(options) == 0 || options[0] == nil {
		return nil
	}

	option := options[0].(map[string]interface{})
	return &storagetransfer.TransferOptions{
		DeleteObjectsFromSourceAfterTransfer:  option["delete_objects_from_source_after_transfer"].(bool),
		DeleteObjectsUniqueInSink:             option["delete_objects_unique_in_sink"].(bool),
		OverwriteObjectsAlreadyExistingInSink: option["overwrite_objects_already_existing_in_sink"].(bool),
		OverwriteWhen:                         option["overwrite_when"].(string),
	}
}

func flattenTransferOption(option *storagetransfer.TransferOptions) []map[string]interface{} {
	data := map[string]interface{}{
		"delete_objects_from_source_after_transfer":  option.DeleteObjectsFromSourceAfterTransfer,
		"delete_objects_unique_in_sink":              option.DeleteObjectsUniqueInSink,
		"overwrite_objects_already_existing_in_sink": option.OverwriteObjectsAlreadyExistingInSink,
		"overwrite_when":                             option.OverwriteWhen,
	}

	return []map[string]interface{}{data}
}

func expandTransferSpecs(transferSpecs []interface{}) *storagetransfer.TransferSpec {
	if len(transferSpecs) == 0 || transferSpecs[0] == nil {
		return nil
	}

	transferSpec := transferSpecs[0].(map[string]interface{})
	return &storagetransfer.TransferSpec{
		SourceAgentPoolName:        transferSpec["source_agent_pool_name"].(string),
		SinkAgentPoolName:          transferSpec["sink_agent_pool_name"].(string),
		GcsDataSink:                expandGcsData(transferSpec["gcs_data_sink"].([]interface{})),
		PosixDataSink:              expandPosixData(transferSpec["posix_data_sink"].([]interface{})),
		ObjectConditions:           expandObjectConditions(transferSpec["object_conditions"].([]interface{})),
		TransferOptions:            expandTransferOptions(transferSpec["transfer_options"].([]interface{})),
		GcsDataSource:              expandGcsData(transferSpec["gcs_data_source"].([]interface{})),
		AwsS3DataSource:            expandAwsS3Data(transferSpec["aws_s3_data_source"].([]interface{})),
		HttpDataSource:             expandHttpData(transferSpec["http_data_source"].([]interface{})),
		AzureBlobStorageDataSource: expandAzureBlobStorageData(transferSpec["azure_blob_storage_data_source"].([]interface{})),
		PosixDataSource:            expandPosixData(transferSpec["posix_data_source"].([]interface{})),
	}
}

func flattenTransferSpec(transferSpec *storagetransfer.TransferSpec, d *schema.ResourceData) []map[string]interface{} {

	data := map[string]interface{}{}

	data["sink_agent_pool_name"] = transferSpec.SinkAgentPoolName
	data["source_agent_pool_name"] = transferSpec.SourceAgentPoolName

	if transferSpec.GcsDataSink != nil {
		data["gcs_data_sink"] = flattenGcsData(transferSpec.GcsDataSink)
	}
	if transferSpec.PosixDataSink != nil {
		data["posix_data_sink"] = flattenPosixData(transferSpec.PosixDataSink)
	}

	if transferSpec.ObjectConditions != nil {
		data["object_conditions"] = flattenObjectCondition(transferSpec.ObjectConditions)
	}
	if transferSpec.TransferOptions != nil &&
		(usingPosix(transferSpec) == false ||
			(usingPosix(transferSpec) == true && reflect.DeepEqual(transferSpec.TransferOptions, &storagetransfer.TransferOptions{}) == false)) {
		data["transfer_options"] = flattenTransferOption(transferSpec.TransferOptions)
	}
	if transferSpec.GcsDataSource != nil {
		data["gcs_data_source"] = flattenGcsData(transferSpec.GcsDataSource)
	} else if transferSpec.AwsS3DataSource != nil {
		data["aws_s3_data_source"] = flattenAwsS3Data(transferSpec.AwsS3DataSource, d)
	} else if transferSpec.HttpDataSource != nil {
		data["http_data_source"] = flattenHttpData(transferSpec.HttpDataSource)
	} else if transferSpec.AzureBlobStorageDataSource != nil {
		data["azure_blob_storage_data_source"] = flattenAzureBlobStorageData(transferSpec.AzureBlobStorageDataSource, d)
	} else if transferSpec.PosixDataSource != nil {
		data["posix_data_source"] = flattenPosixData(transferSpec.PosixDataSource)
	}

	return []map[string]interface{}{data}
}

func usingPosix(transferSpec *storagetransfer.TransferSpec) bool {
	return transferSpec.PosixDataSource != nil || transferSpec.PosixDataSink != nil
}

func expandTransferJobNotificationConfig(notificationConfigs []interface{}) *storagetransfer.NotificationConfig {
	if len(notificationConfigs) == 0 || notificationConfigs[0] == nil {
		return nil
	}

	notificationConfig := notificationConfigs[0].(map[string]interface{})
	var apiData = &storagetransfer.NotificationConfig{
		PayloadFormat: notificationConfig["payload_format"].(string),
		PubsubTopic:   notificationConfig["pubsub_topic"].(string),
	}

	if notificationConfig["event_types"] != nil {
		apiData.EventTypes = tpgresource.ConvertStringArr(notificationConfig["event_types"].(*schema.Set).List())
	}

	log.Printf("[DEBUG] apiData: %v\n\n", apiData)
	return apiData
}

func flattenTransferJobNotificationConfig(notificationConfig *storagetransfer.NotificationConfig) []map[string]interface{} {
	if notificationConfig == nil {
		return nil
	}

	data := map[string]interface{}{
		"payload_format": notificationConfig.PayloadFormat,
		"pubsub_topic":   notificationConfig.PubsubTopic,
	}

	if notificationConfig.EventTypes != nil {
		data["event_types"] = tpgresource.ConvertStringArrToInterface(notificationConfig.EventTypes)
	}

	return []map[string]interface{}{data}
}
