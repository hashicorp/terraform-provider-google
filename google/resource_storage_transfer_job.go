package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/googleapi"
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
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) > 1024 {
						errors = append(errors, fmt.Errorf(
							"%q cannot be greater than 1,024 characters", k))
					}
					return
				},
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
						"object_conditions": objectConditions(),
						"transfer_options":  transferOptions(),
						"gcs_data_sink": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: gcsData(),
							},
						},
						"gcs_data_source": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: gcsData(),
							},
							ConflictsWith: []string{"transfer_spec.aws_s3_data_source", "transfer_spec.http_data_source"},
						},
						"aws_s3_data_source": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: awsS3Data(),
							},
							ConflictsWith: []string{"transfer_spec.gcs_data_source", "transfer_spec.http_data_source"},
						},
						"http_data_source": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: httpData(),
							},
							ConflictsWith: []string{"transfer_spec.aws_s3_data_source", "transfer_spec.gcs_data_source"},
						},
					},
				},
			},
			"schedule": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule_start_date": dateObject(true, false),
						"schedule_end_date":   dateObject(false, true),
						"start_time_of_day":   timeObject(),
					},
				},
			},
			"status": {
				Type:         schema.TypeString,
				Required:     true,
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

func resourceStorageTransferJobCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	description := d.Get("description").(string)
	status := d.Get("status").(string)

	transferJob := &storagetransfer.TransferJob{
		Description:  description,
		ProjectId:    project,
		Status:       status,
		Schedule:     expandTransferSchedule(d),
		TransferSpec: expandTransferSpec(d),
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

	log.Printf("[DEBUG] Created transfer job %v \n\n", res.Name)

	d.SetId(res.Name)
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
			transferJob.Schedule = expandTransferSchedule(v.([]interface{}))
		}
	}

	if d.HasChange("transfer_spec") {
		if v, ok := d.GetOk("transfer_spec"); ok {
			fieldMask = append(fieldMask, "transfer_spec")
			transferJob.TransferSpec = expandTransferSpec(v.([]interface{}))
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

	log.Printf("[DEBUG] Patched transfer job %v\n\n", res.Name)

	d.SetId(res.Name)
	return nil
}

func resourceStorageTransferJobRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	res, err := config.clientStorageTransfer.TransferJobs.Get(d.Get("name").(string)).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Transfer Job %q", d.Get("name").(string)))
	}
	log.Printf("[DEBUG] Read transfer job %v\n\n", res.Name)

	d.Set("project", res.ProjectId)
	d.Set("description", res.Description)
	d.Set("status", res.Status)
	d.Set("schedule", flattenTransferSchedule(res.Schedule))
	d.Set("transfer_spec", flattenTransferSpec(res.TransferSpec))

	d.SetId(res.Name)
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
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		_, err := config.clientStorageTransfer.TransferJobs.Patch(transferJobName, updateRequest).Do()
		if err != nil {
			return resource.RetryableError(err)
		}
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 429 {
			return resource.RetryableError(gerr)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error deleting transfer job %s: %v\n\n", transferJob, err)
		return err
	}

	log.Printf("[DEBUG] Deleted transfer job %v\n\n", transferJob)

	return nil
}

func resourceStorageTransferJobStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("name", d.Id())
	return []*schema.ResourceData{d}, nil
}

func expandTransferSchedule(configured interface{}) *storagetransfer.Schedule {
	schedules := configured.([]interface{})
	schedule := schedules[0].(map[string]map[string]interface{})

	transferSchedule := &storagetransfer.Schedule{
		ScheduleStartDate: &storagetransfer.Date{
			Day:   int64(schedule["schedule_start_date"]["day"].(int)),
			Month: int64(schedule["schedule_start_date"]["month"].(int)),
			Year:  int64(schedule["schedule_start_date"]["year"].(int)),
		},
		ScheduleEndDate: &storagetransfer.Date{
			Day:   int64(schedule["schedule_end_date"]["day"].(int)),
			Month: int64(schedule["schedule_end_date"]["month"].(int)),
			Year:  int64(schedule["schedule_end_date"]["year"].(int)),
		},
		StartTimeOfDay: &storagetransfer.TimeOfDay{
			Hours:   int64(schedule["start_time_of_day"]["hours"].(int)),
			Minutes: int64(schedule["start_time_of_day"]["minutes"].(int)),
			Seconds: int64(schedule["start_time_of_day"]["seconds"].(int)),
			Nanos:   int64(schedule["start_time_of_day"]["nanos"].(int)),
		},
	}

	return transferSchedule
}

func flattenTransferSchedule(transferSchedule *storagetransfer.Schedule) []map[string]map[string]interface{} {
	schedules := make([]map[string]map[string]interface{}, 0, 1)

	if transferSchedule == nil {
		return schedules
	}

	schedule := map[string]map[string]interface{}{
		"schedule_start_date": map[string]interface{}{
			"year":  transferSchedule.ScheduleStartDate.Year,
			"month": transferSchedule.ScheduleStartDate.Month,
			"day":   transferSchedule.ScheduleStartDate.Day,
		},
		"schedule_end_date": map[string]interface{}{
			"year":  transferSchedule.ScheduleEndDate.Year,
			"month": transferSchedule.ScheduleEndDate.Month,
			"day":   transferSchedule.ScheduleEndDate.Day,
		},
		"start_time_of_day": map[string]interface{}{
			"hours":   transferSchedule.StartTimeOfDay.Hours,
			"minutes": transferSchedule.StartTimeOfDay.Minutes,
			"seconds": transferSchedule.StartTimeOfDay.Seconds,
			"nanos":   transferSchedule.StartTimeOfDay.Nanos,
		},
	}

	schedules = append(schedules, schedule)
	return schedules
}

func expandTransferSpec(configured interface{}) *storagetransfer.TransferSpec {
	specs := configured.([]interface{})
	spec := specs[0].(map[string]map[string]interface{})

	transferSpec := &storagetransfer.TransferSpec{}

	// object_conditions
	// transfer_options
	// gcs_data_sink
	// gcs_data_source
	// aws_s3_data_source
	// http_data_source

	return transferSpec
}

func flattenTransferSpec(transferSpec *storagetransfer.TransferSpec) []map[string]map[string]interface{} {
	specs := make([]map[string]map[string]interface{}, 0, 1)

	if transferSpec == nil {
		return specs
	}

	spec := map[string]map[string]interface{}{}

	// object_conditions
	// transfer_options
	// gcs_data_sink
	// gcs_data_source
	// aws_s3_data_source
	// http_data_source

	specs = append(specs, spec)
	return specs
}

func gcsData() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"bucket_name": &schema.Schema{
			Required: true,
			Type:     schema.TypeString,
		},
	}
}

func awsS3Data() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"bucket_name": &schema.Schema{
			Required: true,
			Type:     schema.TypeString,
		},
		"aws_access_key": &schema.Schema{
			Type:     schema.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"access_key_id": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
					"secret_access_key": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
					},
				},
			},
		},
	}
}

func httpData() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"list_url": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func objectConditions() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"min_time_elapsed_since_last_modification": &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateDuration(),
					Optional:     true,
				},
				"max_time_elapsed_since_last_modification": &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateDuration(),
					Optional:     true,
				},
				"include_prefixes": &schema.Schema{
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						MaxItems: 1000,
						Type:     schema.TypeString,
					},
				},
				"exclude_prefixes": &schema.Schema{
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

func validateDuration() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		if _, err := time.ParseDuration(v); err != nil {
			es = append(es, fmt.Errorf("expected %s to be a duration, but parsing gave an error: %s", k, err.Error()))
			return
		}

		return
	}
}

func transferOptions() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"overwrite_objects_already_existing_in_sink": &schema.Schema{
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},
				"delete_objects_unique_in_sink": &schema.Schema{
					Type:          schema.TypeBool,
					Optional:      true,
					Default:       false,
					ConflictsWith: []string{"transfer_spec.transfer_options.delete_objects_from_source_after_transfer"},
				},
				"delete_objects_from_source_after_transfer": &schema.Schema{
					Type:          schema.TypeBool,
					Optional:      true,
					Default:       false,
					ConflictsWith: []string{"transfer_spec.transfer_options.delete_objects_unique_in_sink"},
				},
			},
		},
	}
}

func timeObject() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"hours": &schema.Schema{
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(0, 23),
				},
				"minutes": &schema.Schema{
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(0, 59),
				},
				"seconds": &schema.Schema{
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(0, 59),
				},
				"nanos": &schema.Schema{
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(0, 999999999),
				},
			},
		},
	}

}

func dateObject(required bool, optional bool) *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: required,
		Optional: optional,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"year": &schema.Schema{
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(0, 9999),
				},

				"month": &schema.Schema{
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 12),
				},

				"day": &schema.Schema{
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(0, 31),
				},
			},
		},
	}
}
