package google

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"google.golang.org/api/storagetransfer/v1"
	"log"
	"strings"
	"time"
)

func resourceStorageTransferJob() *schema.Resource {
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(0, 1024),
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"transfer_spec": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"object_conditions": objectConditionsSchema(),
						"transfer_options":  transferOptionsSchema(),
						"gcs_data_sink": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     gcsDataSchema(),
						},
						"gcs_data_source": {
							Type:          schema.TypeList,
							Optional:      true,
							MaxItems:      1,
							Elem:          gcsDataSchema(),
							ConflictsWith: []string{"transfer_spec.aws_s3_data_source", "transfer_spec.http_data_source"},
						},
						"aws_s3_data_source": {
							Type:          schema.TypeList,
							Optional:      true,
							MaxItems:      1,
							Elem:          awsS3DataSchema(),
							ConflictsWith: []string{"transfer_spec.gcs_data_source", "transfer_spec.http_data_source"},
						},
						"http_data_source": {
							Type:          schema.TypeList,
							Optional:      true,
							MaxItems:      1,
							Elem:          httpDataSchema(),
							ConflictsWith: []string{"transfer_spec.aws_s3_data_source", "transfer_spec.gcs_data_source"},
						},
					},
				},
			},
			"schedule": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule_start_date": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							MaxItems: 1,
							Elem:     dateObjectSchema(),
						},
						"schedule_end_date": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem:     dateObjectSchema(),
						},
						"start_time_of_day": {
							Type:             schema.TypeList,
							Optional:         true,
							ForceNew:         true,
							MaxItems:         1,
							Elem:             timeObjectSchema(),
							DiffSuppressFunc: diffSuppressEmptyStartTimeOfDay,
						},
					},
				},
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ENABLED",
				ValidateFunc: validation.StringInSlice([]string{"ENABLED", "DISABLED", "DELETED"}, false),
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modification_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deletion_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
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
					ValidateFunc: validateDuration(),
					Optional:     true,
				},
				"max_time_elapsed_since_last_modification": {
					Type:         schema.TypeString,
					ValidateFunc: validateDuration(),
					Optional:     true,
				},
				"include_prefixes": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						MaxItems: 1000,
						Type:     schema.TypeString,
					},
				},
				"exclude_prefixes": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						MaxItems: 1000,
						Type:     schema.TypeString,
					},
				},
			},
		},
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
					Type:     schema.TypeBool,
					Optional: true,
				},
				"delete_objects_unique_in_sink": {
					Type:          schema.TypeBool,
					Optional:      true,
					ConflictsWith: []string{"transfer_spec.transfer_options.delete_objects_from_source_after_transfer"},
				},
				"delete_objects_from_source_after_transfer": {
					Type:          schema.TypeBool,
					Optional:      true,
					ConflictsWith: []string{"transfer_spec.transfer_options.delete_objects_unique_in_sink"},
				},
			},
		},
	}
}

func timeObjectSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"hours": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 24),
			},
			"minutes": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 59),
			},
			"seconds": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 60),
			},
			"nanos": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 999999999),
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
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 9999),
			},

			"month": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 12),
			},

			"day": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(0, 31),
			},
		},
	}
}

func gcsDataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Required: true,
				Type:     schema.TypeString,
			},
		},
	}
}

func awsS3DataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket_name": {
				Required: true,
				Type:     schema.TypeString,
			},
			"aws_access_key": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key_id": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"secret_access_key": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
					},
				},
			},
		},
	}
}

func httpDataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"list_url": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func diffSuppressEmptyStartTimeOfDay(k, old, new string, d *schema.ResourceData) bool {
	return k == "schedule.0.start_time_of_day.#" && old == "1" && new == "0"
}

func resourceStorageTransferJobCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	transferJob := &storagetransfer.TransferJob{
		Description:  d.Get("description").(string),
		ProjectId:    project,
		Status:       d.Get("status").(string),
		Schedule:     expandTransferSchedules(d.Get("schedule").([]interface{})),
		TransferSpec: expandTransferSpecs(d.Get("transfer_spec").([]interface{})),
	}

	var res *storagetransfer.TransferJob

	err = retry(func() error {
		res, err = config.clientStorageTransfer.TransferJobs.Create(transferJob).Do()
		return err
	})

	if err != nil {
		fmt.Printf("Error creating transfer job %v: %v", transferJob, err)
		return err
	}

	d.Set("name", res.Name)

	name := GetResourceNameFromSelfLink(res.Name)
	d.SetId(fmt.Sprintf("%s/%s", project, name))

	return resourceStorageTransferJobRead(d, meta)
}

func resourceStorageTransferJobRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	res, err := config.clientStorageTransfer.TransferJobs.Get(name).ProjectId(project).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Transfer Job %q", name))
	}
	log.Printf("[DEBUG] Read transfer job: %v in project: %v \n\n", res.Name, res.ProjectId)

	d.Set("project", res.ProjectId)
	d.Set("description", res.Description)
	d.Set("status", res.Status)
	d.Set("last_modification_time", res.LastModificationTime)
	d.Set("creation_time", res.CreationTime)
	d.Set("deletion_time", res.DeletionTime)

	err = d.Set("schedule", flattenTransferSchedule(res.Schedule))
	if err != nil {
		return err
	}

	err = d.Set("transfer_spec", flattenTransferSpec(res.TransferSpec, d))
	if err != nil {
		return err
	}

	return nil
}

func resourceStorageTransferJobUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	transferJob := &storagetransfer.TransferJob{}
	fieldMask := []string{}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			fieldMask = append(fieldMask, "description")
			transferJob.Description = v.(string)
		}
	}

	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			fieldMask = append(fieldMask, "status")
			transferJob.Status = v.(string)
		}
	}

	if d.HasChange("schedule") {
		if v, ok := d.GetOk("schedule"); ok {
			fieldMask = append(fieldMask, "schedule")
			transferJob.Schedule = expandTransferSchedules(v.([]interface{}))
		}
	}

	if d.HasChange("transfer_spec") {
		if v, ok := d.GetOk("transfer_spec"); ok {
			fieldMask = append(fieldMask, "transfer_spec")
			transferJob.TransferSpec = expandTransferSpecs(v.([]interface{}))
		}
	}

	updateRequest := &storagetransfer.UpdateTransferJobRequest{
		ProjectId:   project,
		TransferJob: transferJob,
	}

	updateRequest.UpdateTransferJobFieldMask = strings.Join(fieldMask, ",")

	res, err := config.clientStorageTransfer.TransferJobs.Patch(d.Get("name").(string), updateRequest).Do()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Patched transfer job: %v\n\n", res.Name)
	return nil
}

func resourceStorageTransferJobDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
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
		_, err := config.clientStorageTransfer.TransferJobs.Patch(transferJobName, updateRequest).Do()
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
		d.Set("project", parts[0])
		d.Set("name", fmt.Sprintf("transferJobs/%s", parts[1]))
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
	}
}

func flattenTransferSchedule(transferSchedule *storagetransfer.Schedule) []map[string][]map[string]interface{} {
	data := map[string][]map[string]interface{}{
		"schedule_start_date": flattenDate(transferSchedule.ScheduleStartDate),
	}

	if transferSchedule.ScheduleEndDate != nil {
		data["schedule_end_date"] = flattenDate(transferSchedule.ScheduleEndDate)
	}

	if transferSchedule.StartTimeOfDay != nil {
		data["start_time_of_day"] = flattenTimeOfDay(transferSchedule.StartTimeOfDay)
	}

	return []map[string][]map[string]interface{}{data}
}

func expandGcsData(gcsDatas []interface{}) *storagetransfer.GcsData {
	if len(gcsDatas) == 0 || gcsDatas[0] == nil {
		return nil
	}

	gcsData := gcsDatas[0].(map[string]interface{})
	return &storagetransfer.GcsData{
		BucketName: gcsData["bucket_name"].(string),
	}
}

func flattenGcsData(gcsData *storagetransfer.GcsData) []map[string]interface{} {
	data := map[string]interface{}{
		"bucket_name": gcsData.BucketName,
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
	}
}

func flattenAwsS3Data(awsS3Data *storagetransfer.AwsS3Data, d *schema.ResourceData) []map[string]interface{} {
	data := map[string]interface{}{
		"bucket_name":    awsS3Data.BucketName,
		"aws_access_key": flattenAwsAccessKeys(d),
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

func expandObjectConditions(conditions []interface{}) *storagetransfer.ObjectConditions {
	if len(conditions) == 0 || conditions[0] == nil {
		return nil
	}

	condition := conditions[0].(map[string]interface{})
	return &storagetransfer.ObjectConditions{
		ExcludePrefixes:                     convertStringArr(condition["exclude_prefixes"].([]interface{})),
		IncludePrefixes:                     convertStringArr(condition["include_prefixes"].([]interface{})),
		MaxTimeElapsedSinceLastModification: condition["max_time_elapsed_since_last_modification"].(string),
		MinTimeElapsedSinceLastModification: condition["min_time_elapsed_since_last_modification"].(string),
	}
}

func flattenObjectCondition(condition *storagetransfer.ObjectConditions) []map[string]interface{} {
	data := map[string]interface{}{
		"exclude_prefixes":                         condition.ExcludePrefixes,
		"include_prefixes":                         condition.IncludePrefixes,
		"max_time_elapsed_since_last_modification": condition.MaxTimeElapsedSinceLastModification,
		"min_time_elapsed_since_last_modification": condition.MinTimeElapsedSinceLastModification,
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
	}
}

func flattenTransferOption(option *storagetransfer.TransferOptions) []map[string]interface{} {
	data := map[string]interface{}{
		"delete_objects_from_source_after_transfer":  option.DeleteObjectsFromSourceAfterTransfer,
		"delete_objects_unique_in_sink":              option.DeleteObjectsUniqueInSink,
		"overwrite_objects_already_existing_in_sink": option.OverwriteObjectsAlreadyExistingInSink,
	}

	return []map[string]interface{}{data}
}

func expandTransferSpecs(transferSpecs []interface{}) *storagetransfer.TransferSpec {
	if len(transferSpecs) == 0 || transferSpecs[0] == nil {
		return nil
	}

	transferSpec := transferSpecs[0].(map[string]interface{})
	return &storagetransfer.TransferSpec{
		GcsDataSink:      expandGcsData(transferSpec["gcs_data_sink"].([]interface{})),
		ObjectConditions: expandObjectConditions(transferSpec["object_conditions"].([]interface{})),
		TransferOptions:  expandTransferOptions(transferSpec["transfer_options"].([]interface{})),
		GcsDataSource:    expandGcsData(transferSpec["gcs_data_source"].([]interface{})),
		AwsS3DataSource:  expandAwsS3Data(transferSpec["aws_s3_data_source"].([]interface{})),
		HttpDataSource:   expandHttpData(transferSpec["http_data_source"].([]interface{})),
	}
}

func flattenTransferSpec(transferSpec *storagetransfer.TransferSpec, d *schema.ResourceData) []map[string][]map[string]interface{} {

	data := map[string][]map[string]interface{}{
		"gcs_data_sink": flattenGcsData(transferSpec.GcsDataSink),
	}

	if transferSpec.ObjectConditions != nil {
		data["object_conditions"] = flattenObjectCondition(transferSpec.ObjectConditions)
	}
	if transferSpec.TransferOptions != nil {
		data["transfer_options"] = flattenTransferOption(transferSpec.TransferOptions)
	}
	if transferSpec.GcsDataSource != nil {
		data["gcs_data_source"] = flattenGcsData(transferSpec.GcsDataSource)
	} else if transferSpec.AwsS3DataSource != nil {
		data["aws_s3_data_source"] = flattenAwsS3Data(transferSpec.AwsS3DataSource, d)
	} else if transferSpec.HttpDataSource != nil {
		data["http_data_source"] = flattenHttpData(transferSpec.HttpDataSource)
	}

	return []map[string][]map[string]interface{}{data}
}
