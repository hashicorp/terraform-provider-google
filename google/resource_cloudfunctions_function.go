package google

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"google.golang.org/api/cloudfunctions/v1"

	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Min is 1 second, max is 9 minutes 540 sec
const functionTimeOutMax = 540
const functionTimeOutMin = 1
const functionDefaultTimeout = 60

var functionAllowedMemory = map[int]bool{
	128:  true,
	256:  true,
	512:  true,
	1024: true,
	2048: true,
}

const functionDefaultAllowedMemoryMb = 256

type cloudFunctionId struct {
	Project string
	Region  string
	Name    string
}

func (s *cloudFunctionId) cloudFunctionId() string {
	return fmt.Sprintf("projects/%s/locations/%s/functions/%s", s.Project, s.Region, s.Name)
}

func (s *cloudFunctionId) locationId() string {
	return fmt.Sprintf("projects/%s/locations/%s", s.Project, s.Region)
}

func parseCloudFunctionId(d *schema.ResourceData, config *Config) (*cloudFunctionId, error) {
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<region>[^/]+)/functions/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}
	return &cloudFunctionId{
		Project: d.Get("project").(string),
		Region:  d.Get("region").(string),
		Name:    d.Get("name").(string),
	}, nil
}

func joinMapKeys(mapToJoin *map[int]bool) string {
	var keys []string
	for key := range *mapToJoin {
		keys = append(keys, strconv.Itoa(key))
	}
	return strings.Join(keys, ",")
}

// Differs from validateGcpName because Cloud Functions allow capital letters
// at start/end
func validateResourceCloudFunctionsFunctionName(v interface{}, k string) (ws []string, errors []error) {
	re := `^(?:[a-zA-Z](?:[-_a-zA-Z0-9]{0,61}[a-zA-Z0-9])?)$`
	return validateRegexp(re)(v, k)
}

// based on compareSelfLinkOrResourceName, but less reusable and allows multi-/
// strings in the new state (config) part
func compareSelfLinkOrResourceNameWithMultipleParts(_, old, new string, _ *schema.ResourceData) bool {
	return strings.HasSuffix(old, new)
}

func resourceCloudFunctionsFunction() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudFunctionsCreate,
		Read:   resourceCloudFunctionsRead,
		Update: resourceCloudFunctionsUpdate,
		Delete: resourceCloudFunctionsDestroy,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateResourceCloudFunctionsFunctionName,
			},

			"source_archive_bucket": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"source_archive_object": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"source_repository": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"source_archive_bucket", "source_archive_object"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
						"deployed_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"available_memory_mb": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  functionDefaultAllowedMemoryMb,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					availableMemoryMB := v.(int)

					if !functionAllowedMemory[availableMemoryMB] {
						errors = append(errors, fmt.Errorf("Allowed values for memory (in MB) are: %s . Got %d",
							joinMapKeys(&functionAllowedMemory), availableMemoryMB))
					}
					return
				},
			},

			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      functionDefaultTimeout,
				ValidateFunc: validation.IntBetween(functionTimeOutMin, functionTimeOutMax),
			},

			"entry_point": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"runtime": {
				Type:     schema.TypeString,
				Required: true,
			},

			"service_account_email": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"vpc_connector": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},

			"environment_variables": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"trigger_http": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"event_trigger": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"trigger_http"},
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_type": {
							Type:     schema.TypeString,
							ForceNew: true,
							Required: true,
						},
						"resource": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: compareSelfLinkOrResourceNameWithMultipleParts,
						},
						"failure_policy": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"retry": {
										Type: schema.TypeBool,
										// not strictly required, but this way an empty block can't be specified
										Required: true,
									},
								}},
						},
					},
				},
			},

			"https_trigger_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"max_instances": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
			},

			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCloudFunctionsCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	cloudFuncId := &cloudFunctionId{
		Project: project,
		Region:  region,
		Name:    d.Get("name").(string),
	}

	function := &cloudfunctions.CloudFunction{
		Name:                cloudFuncId.cloudFunctionId(),
		Runtime:             d.Get("runtime").(string),
		ServiceAccountEmail: d.Get("service_account_email").(string),
		ForceSendFields:     []string{},
	}

	sourceRepos := d.Get("source_repository").([]interface{})
	if len(sourceRepos) > 0 {
		function.SourceRepository = expandSourceRepository(sourceRepos)
	} else {
		sourceArchiveBucket := d.Get("source_archive_bucket").(string)
		sourceArchiveObj := d.Get("source_archive_object").(string)
		if sourceArchiveBucket == "" || sourceArchiveObj == "" {
			return fmt.Errorf("either source_repository or both of source_archive_bucket+source_archive_object must be set")
		}
		function.SourceArchiveUrl = fmt.Sprintf("gs://%v/%v", sourceArchiveBucket, sourceArchiveObj)
	}

	if v, ok := d.GetOk("available_memory_mb"); ok {
		availableMemoryMb := v.(int)
		function.AvailableMemoryMb = int64(availableMemoryMb)
	}

	if v, ok := d.GetOk("description"); ok {
		function.Description = v.(string)
	}

	if v, ok := d.GetOk("entry_point"); ok {
		function.EntryPoint = v.(string)
	}

	if v, ok := d.GetOk("timeout"); ok {
		function.Timeout = fmt.Sprintf("%vs", v.(int))
	}

	if v, ok := d.GetOk("event_trigger"); ok {
		function.EventTrigger = expandEventTrigger(v.([]interface{}), project)
	} else if v, ok := d.GetOk("trigger_http"); ok && v.(bool) {
		function.HttpsTrigger = &cloudfunctions.HttpsTrigger{}
	} else {
		return fmt.Errorf("One of `event_trigger` or `trigger_http` is required: " +
			"You must specify a trigger when deploying a new function.")
	}

	if _, ok := d.GetOk("labels"); ok {
		function.Labels = expandLabels(d)
	}

	if _, ok := d.GetOk("environment_variables"); ok {
		function.EnvironmentVariables = expandEnvironmentVariables(d)
	}

	if v, ok := d.GetOk("vpc_connector"); ok {
		function.VpcConnector = v.(string)
	}

	if v, ok := d.GetOk("max_instances"); ok {
		function.MaxInstances = int64(v.(int))
	}

	log.Printf("[DEBUG] Creating cloud function: %s", function.Name)
	op, err := config.clientCloudFunctions.Projects.Locations.Functions.Create(
		cloudFuncId.locationId(), function).Do()
	if err != nil {
		return err
	}

	// Name of function should be unique
	d.SetId(cloudFuncId.cloudFunctionId())

	err = cloudFunctionsOperationWait(config, op, "Creating CloudFunctions Function",
		int(d.Timeout(schema.TimeoutCreate).Minutes()))
	if err != nil {
		return err
	}

	return resourceCloudFunctionsRead(d, meta)
}

func resourceCloudFunctionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	cloudFuncId, err := parseCloudFunctionId(d, config)
	if err != nil {
		return err
	}

	function, err := config.clientCloudFunctions.Projects.Locations.Functions.Get(cloudFuncId.cloudFunctionId()).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Target CloudFunctions Function %q", cloudFuncId.Name))
	}

	d.Set("name", cloudFuncId.Name)
	d.Set("description", function.Description)
	d.Set("entry_point", function.EntryPoint)
	d.Set("available_memory_mb", function.AvailableMemoryMb)
	sRemoved := strings.Replace(function.Timeout, "s", "", -1)
	timeout, err := strconv.Atoi(sRemoved)
	if err != nil {
		return err
	}
	d.Set("timeout", timeout)
	d.Set("labels", function.Labels)
	d.Set("runtime", function.Runtime)
	d.Set("service_account_email", function.ServiceAccountEmail)
	d.Set("environment_variables", function.EnvironmentVariables)
	d.Set("vpc_connector", function.VpcConnector)
	if function.SourceArchiveUrl != "" {
		// sourceArchiveUrl should always be a Google Cloud Storage URL (e.g. gs://bucket/object)
		// https://cloud.google.com/functions/docs/reference/rest/v1/projects.locations.functions
		sourceURL, err := url.Parse(function.SourceArchiveUrl)
		if err != nil {
			return err
		}
		bucket := sourceURL.Host
		object := strings.TrimLeft(sourceURL.Path, "/")
		d.Set("source_archive_bucket", bucket)
		d.Set("source_archive_object", object)
	}
	d.Set("source_repository", flattenSourceRepository(function.SourceRepository))

	if function.HttpsTrigger != nil {
		d.Set("trigger_http", true)
		d.Set("https_trigger_url", function.HttpsTrigger.Url)
	}

	d.Set("event_trigger", flattenEventTrigger(function.EventTrigger))
	d.Set("max_instances", function.MaxInstances)
	d.Set("region", cloudFuncId.Region)
	d.Set("project", cloudFuncId.Project)

	return nil
}

func resourceCloudFunctionsUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]: Updating google_cloudfunctions_function")
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	cloudFuncId, err := parseCloudFunctionId(d, config)
	if err != nil {
		return err
	}

	d.Partial(true)

	function := cloudfunctions.CloudFunction{
		Name: cloudFuncId.cloudFunctionId(),
	}

	var updateMaskArr []string
	if d.HasChange("available_memory_mb") {
		availableMemoryMb := d.Get("available_memory_mb").(int)
		function.AvailableMemoryMb = int64(availableMemoryMb)
		updateMaskArr = append(updateMaskArr, "availableMemoryMb")
	}

	if d.HasChange("source_archive_bucket") || d.HasChange("source_archive_object") {
		sourceArchiveBucket := d.Get("source_archive_bucket").(string)
		sourceArchiveObj := d.Get("source_archive_object").(string)
		function.SourceArchiveUrl = fmt.Sprintf("gs://%v/%v", sourceArchiveBucket, sourceArchiveObj)
		updateMaskArr = append(updateMaskArr, "sourceArchiveUrl")
	}

	if d.HasChange("source_repository") {
		function.SourceRepository = expandSourceRepository(d.Get("source_repository").([]interface{}))
		updateMaskArr = append(updateMaskArr, "sourceRepository")
	}

	if d.HasChange("description") {
		function.Description = d.Get("description").(string)
		updateMaskArr = append(updateMaskArr, "description")
	}

	if d.HasChange("timeout") {
		function.Timeout = fmt.Sprintf("%vs", d.Get("timeout").(int))
		updateMaskArr = append(updateMaskArr, "timeout")
	}

	if d.HasChange("labels") {
		function.Labels = expandLabels(d)
		updateMaskArr = append(updateMaskArr, "labels")
	}

	if d.HasChange("runtime") {
		function.Runtime = d.Get("runtime").(string)
		updateMaskArr = append(updateMaskArr, "runtime")
	}

	if d.HasChange("environment_variables") {
		function.EnvironmentVariables = expandEnvironmentVariables(d)
		updateMaskArr = append(updateMaskArr, "environmentVariables")
	}

	if d.HasChange("vpc_connector") {
		function.VpcConnector = d.Get("vpc_connector").(string)
		updateMaskArr = append(updateMaskArr, "vpcConnector")
	}

	if d.HasChange("event_trigger") {
		function.EventTrigger = expandEventTrigger(d.Get("event_trigger").([]interface{}), project)
		updateMaskArr = append(updateMaskArr, "eventTrigger", "eventTrigger.failurePolicy.retry")
	}

	if d.HasChange("max_instances") {
		function.MaxInstances = int64(d.Get("max_instances").(int))
		updateMaskArr = append(updateMaskArr, "maxInstances")
	}

	if len(updateMaskArr) > 0 {
		log.Printf("[DEBUG] Send Patch CloudFunction Configuration request: %#v", function)
		updateMask := strings.Join(updateMaskArr, ",")
		op, err := config.clientCloudFunctions.Projects.Locations.Functions.Patch(function.Name, &function).
			UpdateMask(updateMask).Do()

		if err != nil {
			return fmt.Errorf("Error while updating cloudfunction configuration: %s", err)
		}

		err = cloudFunctionsOperationWait(config, op, "Updating CloudFunctions Function",
			int(d.Timeout(schema.TimeoutUpdate).Minutes()))
		if err != nil {
			return err
		}
	}
	d.Partial(false)

	return resourceCloudFunctionsRead(d, meta)
}

func resourceCloudFunctionsDestroy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	cloudFuncId, err := parseCloudFunctionId(d, config)
	if err != nil {
		return err
	}

	op, err := config.clientCloudFunctions.Projects.Locations.Functions.Delete(cloudFuncId.cloudFunctionId()).Do()
	if err != nil {
		return err
	}
	err = cloudFunctionsOperationWait(config, op, "Deleting CloudFunctions Function",
		int(d.Timeout(schema.TimeoutDelete).Minutes()))
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func expandEventTrigger(configured []interface{}, project string) *cloudfunctions.EventTrigger {
	if len(configured) == 0 || configured[0] == nil {
		return nil
	}

	data := configured[0].(map[string]interface{})
	eventType := data["event_type"].(string)
	resource := data["resource"].(string)

	// if resource starts with "projects/", we can reasonably assume it's a
	// partial URI. Otherwise, it's a shortname. Construct a partial URI based
	// on the event type if so.
	if !strings.HasPrefix(resource, "projects/") {
		shape := ""
		switch {
		case strings.HasPrefix(eventType, "google.storage.object."):
			shape = "projects/%s/buckets/%s"
		case strings.HasPrefix(eventType, "google.pubsub.topic."):
			shape = "projects/%s/topics/%s"
		// Legacy style triggers
		case strings.HasPrefix(eventType, "providers/cloud.storage/eventTypes/"):
			// Note that this is an uncommon way to refer to buckets; normally,
			// you'd use to the global URL of the bucket and not the project
			// scoped one.
			shape = "projects/%s/buckets/%s"
		case strings.HasPrefix(eventType, "providers/cloud.pubsub/eventTypes/"):
			shape = "projects/%s/topics/%s"
		case strings.HasPrefix(eventType, "providers/cloud.firestore/eventTypes/"):
			// Firestore doesn't not yet support multiple databases, so "(default)" is assumed.
			// https://cloud.google.com/functions/docs/calling/cloud-firestore#deploying_your_function
			shape = "projects/%s/databases/(default)/documents/%s"
		}

		resource = fmt.Sprintf(shape, project, resource)
	}

	return &cloudfunctions.EventTrigger{
		EventType:     eventType,
		Resource:      resource,
		FailurePolicy: expandFailurePolicy(data["failure_policy"].([]interface{})),
	}
}

func flattenEventTrigger(eventTrigger *cloudfunctions.EventTrigger) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if eventTrigger == nil {
		return result
	}

	result = append(result, map[string]interface{}{
		"event_type":     eventTrigger.EventType,
		"resource":       eventTrigger.Resource,
		"failure_policy": flattenFailurePolicy(eventTrigger.FailurePolicy),
	})

	return result
}

func expandFailurePolicy(configured []interface{}) *cloudfunctions.FailurePolicy {
	if len(configured) == 0 || configured[0] == nil {
		return &cloudfunctions.FailurePolicy{}
	}

	if data := configured[0].(map[string]interface{}); data["retry"].(bool) {
		return &cloudfunctions.FailurePolicy{
			Retry: &cloudfunctions.Retry{},
		}
	}

	return nil
}

func flattenFailurePolicy(failurePolicy *cloudfunctions.FailurePolicy) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if failurePolicy == nil {
		return nil
	}

	result = append(result, map[string]interface{}{
		"retry": failurePolicy.Retry != nil,
	})

	return result
}

func expandSourceRepository(configured []interface{}) *cloudfunctions.SourceRepository {
	if len(configured) == 0 || configured[0] == nil {
		return &cloudfunctions.SourceRepository{}
	}

	data := configured[0].(map[string]interface{})
	return &cloudfunctions.SourceRepository{
		Url: data["url"].(string),
	}
}

func flattenSourceRepository(sourceRepo *cloudfunctions.SourceRepository) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if sourceRepo == nil {
		return nil
	}

	result = append(result, map[string]interface{}{
		"url":          sourceRepo.Url,
		"deployed_url": sourceRepo.DeployedUrl,
	})

	return result
}
