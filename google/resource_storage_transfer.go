
package google

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"time"
	"google.golang.org/api/storagetransfer/v1"
	"log"
)

// https://cloud.google.com/storage/transfer/create-manage-transfer-program
// https://cloud.google.com/storage/transfer/reference/rest/v1/transferJobs/patch
// https://cloud.google.com/storage/transfer/reference/rest/v1/transferJobs#Status
func resourceStorageTransfer() *schema.Resource {
	return &schema.Resource{

		Create: resourceStorageTransferCreate,
		Read:   resourceStorageTransferRead,
		Update: resourceStorageTransferUpdate,
		Delete: resourceStorageTransferDelete,
		Importer: &schema.ResourceImporter{
			State: resourceStorageTransferStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
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
						},
						// enrich with HTTP data, S3 bucket
					},
				},
			},
			"schedule": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule_start_date": dateObject(),
						"schedule_end_date":   dateObject(),
						"start_time_of_day":    timeObject(),
					},
				},
			},
			"status": {
				Type:         schema.TypeString,
				Required:     true,
				Default:      "ENABLED",
				ValidateFunc: validation.StringInSlice([]string{"STATUS_UNSPECIFIED", "ENABLED", "DISABLED", "DELETED"}, true),
			},
			"creation_time": {
				Type: schema.TypeString,
			},
			"last_modification_time": {
				Type: schema.TypeString,
			},
			"deletion_time": {
				Type: schema.TypeString,
			},
		},
	}
}

func gcsData() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"bucketName": &schema.Schema{
			Required: true,
			Type:     schema.TypeString,
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
				},
				"max_time_elapsed_since_last_modification": &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateDuration(),
				},
				"include_prefixes": &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						MaxItems: 1000,
						Type:     schema.TypeString,
					},
				},
				"exclude_prefixes": &schema.Schema{
					Type: schema.TypeList,
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
					Type: schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"delete_objects_unique_in_sink": &schema.Schema{
					Type: schema.TypeBool,
					Optional: true,
					Computed: true,
				},
				"delete_objects_from_source_after_transfer": &schema.Schema{
					Type: schema.TypeBool,
					Optional: true,
					Computed: true,
				},
			},
		},
	}
}

func timeObject() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: false,
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

func dateObject() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: false,
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

func resourceStorageTransferCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var transferJob *storagetransfer.TransferJob

	var res *storagetransfer.TransferJob


	var err error
	err = retry(func() error {
		res, err = config.clientStorageTransfer.TransferJobs.Create(transferJob).Do()

		return err
	})

	if err != nil {
		fmt.Printf("Error creating transferJob %s: %v", transferJob, err)
		return err
	}

	log.Printf("[DEBUG] Created transferjob %v \n\n", res.Name)

	d.SetId(res.Name)
	return resourceStorageBucketRead(d, meta)

}

func resourceStorageTransferUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceStorageTransferRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Get the bucket and acl
	name := d.Get("name").(string)
	res, err := config.clientStorageTransfer.TransferJobs.Get(name).Do()

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Transfer %q", d.Get("name").(string)))
	}
	log.Printf("[DEBUG] Read transfer %v \n\n", res.Name)

	// marshal fields

	d.SetId(res.Name)
	return nil
}

func resourceStorageTransferDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceStorageTransferStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return nil, nil
}
