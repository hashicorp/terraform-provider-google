package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/dataproc/v1"
)

func resourceDataprocJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataprocJobCreate,
		Read:   resourceDataprocJobRead,
		//Update: resourceDataprocJobUpdate,
		Delete: resourceDataprocJobDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "global",
				ForceNew: true,
			},

			"force_delete": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
				ForceNew: true,
			},

			"cluster": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     schema.TypeString,
			},

			"pyspark_config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"main_python_file": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"additional_python_files": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"jar_files": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"args": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},

						"properties": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
							Elem:     schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func resourceDataprocJobCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	clusterName := d.Get("cluster").(string)
	region := d.Get("region").(string)

	submitReq := &dataproc.SubmitJobRequest{
		Job: &dataproc.Job{
			Placement: &dataproc.JobPlacement{
				ClusterName: clusterName,
			},
		},
	}

	if v, ok := d.GetOk("labels"); ok {
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		submitReq.Job.Labels = m
	}

	if v, ok := d.GetOk("pyspark_config"); ok {
		configs := v.([]interface{})
		config := configs[0].(map[string]interface{})

		submitReq.Job.PysparkJob = &dataproc.PySparkJob{}
		if v, ok = config["main_python_file"]; ok {
			submitReq.Job.PysparkJob.MainPythonFileUri = v.(string)
		}
		if v, ok = config["args"]; ok {
			arrList := v.([]interface{})
			arr := []string{}
			for _, v := range arrList {
				arr = append(arr, v.(string))
			}
			submitReq.Job.PysparkJob.Args = arr
		}
		if v, ok = config["properties"]; ok {
			m := make(map[string]string)
			for k, val := range v.(map[string]interface{}) {
				m[k] = val.(string)
			}
			submitReq.Job.PysparkJob.Properties = m
		}
		if v, ok = config["jar_files"]; ok {
			arrList := v.([]interface{})
			arr := []string{}
			for _, v := range arrList {
				arr = append(arr, v.(string))
			}
			submitReq.Job.PysparkJob.JarFileUris = arr
		}
		if v, ok = config["additional_python_files"]; ok {
			arrList := v.([]interface{})
			arr := []string{}
			for _, v := range arrList {
				arr = append(arr, v.(string))
			}
			submitReq.Job.PysparkJob.PythonFileUris = arr
		}
	}

	// Submit the job
	job, err := config.clientDataproc.Projects.Regions.Jobs.Submit(
		project, region, submitReq).Do()
	if err != nil {
		return err
	}
	d.SetId(job.Reference.JobId)

	// We don't bother to wait and check the status as the non error code
	// from the above in this case is good enough ...

	log.Printf("[INFO] Dataproc job %s has been created", job.Reference.JobId)

	e := resourceDataprocJobRead(d, meta)
	if e != nil {
		log.Print("[INFO] Got an error reading back dataproc job after creating, \n", e)
	}
	return e
}

/*
func resourceDataprocJobUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	return nil

	// The only items which are currently able to be updated, without a
	// forceNew in place are the labels and/or the number of worker nodes in a cluster
	if !d.HasChange("labels") {
		return errors.New("*** programmer issue - update resource called however item is not allowed to be changed - investigate ***")
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)
	jobId := d.Id()
	timeoutInMinutes := int(d.Timeout(schema.TimeoutUpdate).Minutes())

	job := &dataproc.Job{}
	patch := config.clientDataproc.Projects.Regions.Jobs.Patch(
		project, region, jobId, job)

	updMask := ""

	if d.HasChange("labels") {

		v := d.Get("labels")
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		job.Labels = m

		updMask = "labels"
	}

	patch.UpdateMask(updMask)

	_, err = patch.Do()
	if err != nil {
		return err
	}

	// Wait until it's updated
	waitErr := dataprocJobOperationWait(config, region, project, jobId, "updating Dataproc job ", timeoutInMinutes, 2)
	if waitErr != nil {
		return waitErr
	}

	log.Printf("[INFO] Dataproc job %s has been updated ", jobId)
	return resourceDataprocJobRead(d, meta)
}
*/

func resourceDataprocJobRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := d.Get("region").(string)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	job, err := config.clientDataproc.Projects.Regions.Jobs.Get(
		project, region, d.Id()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Dataproc Job %q", d.Id()))
	}

	d.Set("labels", job.Labels)
	d.Set("cluster", job.Placement.ClusterName)

	if job.PysparkJob != nil {

		pySparkConfig := []map[string]interface{}{
			{
				"args":                    job.PysparkJob.Args,
				"properties":              job.PysparkJob.Properties,
				"jar_files":               job.PysparkJob.JarFileUris,
				"main_python_file":        job.PysparkJob.MainPythonFileUri,
				"additional_python_files": job.PysparkJob.PythonFileUris,
			},
		}

		d.Set("pyspark_config", pySparkConfig)

	}
	return nil
}

func resourceDataprocJobDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region := d.Get("region").(string)
	forceDelete := d.Get("force_delete").(bool)

	if forceDelete {
		log.Printf("[DEBUG] Attempting to first cancel Dataproc job %s if its still running ...", d.Id())

		_, _ = config.clientDataproc.Projects.Regions.Jobs.Cancel(
			project, region, d.Id(), &dataproc.CancelJobRequest{}).Do()

		// ignore error if we get one (Job may not need to be cancelled, simply proceed to delete)
		waitErr := dataprocJobOperationWait(config, region, project, d.Id(),
			"Cancelling Dataproc job", 2, 1)
		if waitErr != nil {
			return waitErr
		}

	}

	log.Printf("[DEBUG] Deleting Dataproc job %s", d.Id())
	_, err = config.clientDataproc.Projects.Regions.Jobs.Delete(
		project, region, d.Id()).Do()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Dataproc job %s has been deleted", d.Id())
	d.SetId("")

	return nil
}
