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

func resourceStorageTransferJobCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	description := d.Get("description").(string)
	status := d.Get("status").(string)
	schedules := d.Get("schedule").([]interface{})
	transferSpecs := d.Get("transfer_spec").([]interface{})

	transferJob := &storagetransfer.TransferJob{
		Description:  description,
		ProjectId:    project,
		Status:       status,
		Schedule:     expandTransferSchedules(schedules)[0],
		TransferSpec: expandTransferSpecs(transferSpecs)[0],
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
			transferJob.Schedule = expandTransferSchedules(v.([]interface{}))[0]
		}
	}

	if d.HasChange("transfer_spec") {
		if v, ok := d.GetOk("transfer_spec"); ok {
			fieldMask = append(fieldMask, "transfer_spec")
			transferJob.TransferSpec = expandTransferSpecs(v.([]interface{}))[0]
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

	err = d.Set("schedule", flattenTransferSchedules([]*storagetransfer.Schedule{res.Schedule}))
	if err != nil {
		return err
	}

	d.Set("transfer_spec", flattenTransferSpecs([]*storagetransfer.TransferSpec{res.TransferSpec}))
	if err != nil {
		return err
	}

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
	log.Printf("[DEBUG] Setting status to DELETE for: %v\n\n", transferJobName)
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
		fmt.Printf("Error deleting transfer job %v: %v\n\n", transferJob, err)
		return err
	}

	log.Printf("[DEBUG] Deleted transfer job %v\n\n", transferJob)

	return nil
}

func resourceStorageTransferJobStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	d.Set("name", d.Id())
	return []*schema.ResourceData{d}, nil
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
