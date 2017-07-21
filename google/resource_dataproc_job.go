package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
	"google.golang.org/api/dataproc/v1"
)

func resourceDataprocJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceDataprocJobCreate,
		Read:   resourceDataprocJobRead,
		//Update: not supported,
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

			"pyspark_config": pySparkTFSchema(),
			"spark_config":   sparkTFSchema(),
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

	jobConfCount := 0

	if v, ok := d.GetOk("pyspark_config"); ok {
		jobConfCount++
		configs := v.([]interface{})
		config := configs[0].(map[string]interface{})
		submitReq.Job.PysparkJob = getPySparkJob(config)
	}

	if v, ok := d.GetOk("spark_config"); ok {
		jobConfCount++
		configs := v.([]interface{})
		config := configs[0].(map[string]interface{})
		submitReq.Job.SparkJob = getSparkJob(config)
	}

	if jobConfCount != 1 {
		return errors.New("You must define and configure exactly one xxx_config block")
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
		pySparkConfig := getTfPySparkConfig(job.PysparkJob)
		d.Set("pyspark_config", pySparkConfig)
	}
	if job.SparkJob != nil {
		sparkConfig := getTfSparkConfig(job.SparkJob)
		d.Set("spark_config", sparkConfig)
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
	timeoutInMinutes := int(d.Timeout(schema.TimeoutDelete).Minutes())

	if forceDelete {
		log.Printf("[DEBUG] Attempting to first cancel Dataproc job %s if its still running ...", d.Id())

		_, _ = config.clientDataproc.Projects.Regions.Jobs.Cancel(
			project, region, d.Id(), &dataproc.CancelJobRequest{}).Do()
		// ignore error if we get one - job may be finished already and not need to
		// be cancelled. We do however wait for the state to be one that is
		// at least not active
		waitErr := dataprocJobOperationWait(config, region, project, d.Id(),
			"Cancelling Dataproc job", timeoutInMinutes, 1)
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

	waitErr := dataprocDeleteOperationWait(config, region, project, d.Id(),
		"Deleting Dataproc job", timeoutInMinutes, 1)
	if waitErr != nil {
		return waitErr
	}

	log.Printf("[INFO] Dataproc job %s has been deleted", d.Id())
	d.SetId("")

	return nil
}

// ---- PySpark Job ----

func pySparkTFSchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		ForceNew:      true,
		MaxItems:      1,
		ConflictsWith: []string{"spark_config"},
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
	}
}

func getTfPySparkConfig(job *dataproc.PySparkJob) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"args":                    job.Args,
			"properties":              job.Properties,
			"jar_files":               job.JarFileUris,
			"main_python_file":        job.MainPythonFileUri,
			"additional_python_files": job.PythonFileUris,
		},
	}
}

func getPySparkJob(config map[string]interface{}) *dataproc.PySparkJob {

	job := &dataproc.PySparkJob{}
	if v, ok := config["main_python_file"]; ok {
		job.MainPythonFileUri = v.(string)
	}
	if v, ok := config["args"]; ok {
		arrList := v.([]interface{})
		arr := []string{}
		for _, v := range arrList {
			arr = append(arr, v.(string))
		}
		job.Args = arr
	}
	if v, ok := config["properties"]; ok {
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		job.Properties = m
	}
	if v, ok := config["jar_files"]; ok {
		arrList := v.([]interface{})
		arr := []string{}
		for _, v := range arrList {
			arr = append(arr, v.(string))
		}
		job.JarFileUris = arr
	}
	if v, ok := config["additional_python_files"]; ok {
		arrList := v.([]interface{})
		arr := []string{}
		for _, v := range arrList {
			arr = append(arr, v.(string))
		}
		job.PythonFileUris = arr
	}

	return job

}

// ---- Spark Job ----

func sparkTFSchema() *schema.Schema {
	return &schema.Schema{
		Type:          schema.TypeList,
		Optional:      true,
		ForceNew:      true,
		MaxItems:      1,
		ConflictsWith: []string{"pyspark_config"},
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{

				"main_class": {
					Type:          schema.TypeString,
					Optional:      true,
					ForceNew:      true,
					ConflictsWith: []string{"spark_config.main_jar"},
				},

				"main_jar": {
					Type:          schema.TypeString,
					Optional:      true,
					ForceNew:      true,
					ConflictsWith: []string{"spark_config.main_class"},
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
	}
}

func getTfSparkConfig(job *dataproc.SparkJob) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"main_class": job.MainClass,
			"main_jar":   job.MainJarFileUri,
			"jar_files":  job.JarFileUris,
			"args":       job.Args,
			"properties": job.Properties,
		},
	}
}

func getSparkJob(config map[string]interface{}) *dataproc.SparkJob {

	job := &dataproc.SparkJob{}
	if v, ok := config["main_class"]; ok {
		job.MainClass = v.(string)
	}
	if v, ok := config["main_jar"]; ok {
		job.MainJarFileUri = v.(string)
	}
	if v, ok := config["jar_files"]; ok {
		arrList := v.([]interface{})
		arr := []string{}
		for _, v := range arrList {
			arr = append(arr, v.(string))
		}
		job.JarFileUris = arr
	}
	if v, ok := config["args"]; ok {
		arrList := v.([]interface{})
		arr := []string{}
		for _, v := range arrList {
			arr = append(arr, v.(string))
		}
		job.Args = arr
	}
	if v, ok := config["properties"]; ok {
		m := make(map[string]string)
		for k, val := range v.(map[string]interface{}) {
			m[k] = val.(string)
		}
		job.Properties = m
	}

	return job

}
