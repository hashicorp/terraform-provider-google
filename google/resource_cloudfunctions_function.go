package google

import (
	"regexp"

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

var allowedIngressSettings = []string{
	"ALLOW_ALL",
	"ALLOW_INTERNAL_ONLY",
}

var allowedVpcConnectorEgressSettings = []string{
	"ALL_TRAFFIC",
	"PRIVATE_RANGES_ONLY",
}

const functionDefaultIngressSettings = "ALLOW_ALL"

type cloudFunctionId struct {
	Project string
	Region  string
	Name    string
}

func (s *cloudFunctionId) cloudFunctionId() string {
	return fmt.Sprintf("projects/%s/locations/%s/functions/%s", s.Project, s.Region, s.Name)
}

// matches all international lower case letters, number, underscores and dashes.
var labelKeyRegex = regexp.MustCompile(`^[\p{Ll}0-9_-]+$`)

func labelKeyValidator(val interface{}, key string) (warns []string, errs []error) {
	if val == nil {
		return
	}

	m := val.(map[string]interface{})
	for k := range m {
		if !labelKeyRegex.MatchString(k) {
			errs = append(errs, fmt.Errorf("%q is an invalid label key. See https://cloud.google.com/resource-manager/docs/creating-managing-labels#requirements", k))
		}
	}
	return
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
				Description:  `A user-defined name of the function. Function names must be unique globally.`,
				ValidateFunc: validateResourceCloudFunctionsFunctionName,
			},

			"source_archive_bucket": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The GCS bucket containing the zip archive which contains the function.`,
			},

			"source_archive_object": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The source archive object (file) in archive bucket.`,
			},

			"source_repository": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				Description:   `Represents parameters related to source repository where a function is hosted. Cannot be set alongside source_archive_bucket or source_archive_object.`,
				ConflictsWith: []string{"source_archive_bucket", "source_archive_object"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The URL pointing to the hosted repository where the function is defined.`,
						},
						"deployed_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The URL pointing to the hosted repository where the function was defined at the time of deployment.`,
						},
					},
				},
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Description of the function.`,
			},

			"available_memory_mb": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     functionDefaultAllowedMemoryMb,
				Description: `Memory (in MB), available to the function. Default value is 256MB. Allowed values are: 128MB, 256MB, 512MB, 1024MB, and 2048MB.`,
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
				Description:  `Timeout (in seconds) for the function. Default value is 60 seconds. Cannot be more than 540 seconds.`,
			},

			"entry_point": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Name of the function that will be executed when the Google Cloud Function is triggered.`,
			},

			"ingress_settings": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      functionDefaultIngressSettings,
				ValidateFunc: validation.StringInSlice(allowedIngressSettings, true),
				Description:  `String value that controls what traffic can reach the function. Allowed values are ALLOW_ALL and ALLOW_INTERNAL_ONLY. Changes to this field will recreate the cloud function.`,
			},

			"vpc_connector_egress_settings": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice(allowedVpcConnectorEgressSettings, true),
				Description:  `The egress settings for the connector, controlling what traffic is diverted through it. Allowed values are ALL_TRAFFIC and PRIVATE_RANGES_ONLY. Defaults to PRIVATE_RANGES_ONLY. If unset, this field preserves the previously set value.`,
			},

			"labels": {
				Type:         schema.TypeMap,
				ValidateFunc: labelKeyValidator,
				Optional:     true,
				Description:  `A set of key/value label pairs to assign to the function. Label keys must follow the requirements at https://cloud.google.com/resource-manager/docs/creating-managing-labels#requirements.`,
			},

			"runtime": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The runtime in which the function is going to run. Eg. "nodejs8", "nodejs10", "python37", "go111".`,
			},

			"service_account_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: ` If provided, the self-provided service account to run the function with.`,
			},

			"vpc_connector": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      `The VPC Network Connector that this cloud function can connect to. It can be either the fully-qualified URI, or the short name of the network connector resource. The format of this field is projects/*/locations/*/connectors/*.`,
			},

			"environment_variables": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: `A set of key/value environment variable pairs to assign to the function.`,
			},

			"trigger_http": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Boolean variable. Any HTTP request (of a supported type) to the endpoint will trigger function execution. Supported HTTP request types are: POST, PUT, GET, DELETE, and OPTIONS. Endpoint is returned as https_trigger_url. Cannot be used with trigger_bucket and trigger_topic.`,
			},

			"event_trigger": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"trigger_http"},
				MaxItems:      1,
				Description:   `A source that fires events in response to a condition in another service. Cannot be used with trigger_http.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_type": {
							Type:        schema.TypeString,
							ForceNew:    true,
							Required:    true,
							Description: `The type of event to observe. For example: "google.storage.object.finalize". See the documentation on calling Cloud Functions for a full reference of accepted triggers.`,
						},
						"resource": {
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: compareSelfLinkOrResourceNameWithMultipleParts,
							Description:      `The name or partial URI of the resource from which to observe events. For example, "myBucket" or "projects/my-project/topics/my-topic"`,
						},
						"failure_policy": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: `Specifies policy for failed executions`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"retry": {
										Type: schema.TypeBool,
										// not strictly required, but this way an empty block can't be specified
										Required:    true,
										Description: `Whether the function should be retried on failure. Defaults to false.`,
									},
								}},
						},
					},
				},
			},

			"https_trigger_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `URL which triggers function execution. Returned only if trigger_http is used.`,
			},

			"max_instances": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntAtLeast(0),
				Description:  `The limit on the maximum number of function instances that may coexist at a given time.`,
			},

			"project": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Project of the function. If it is not provided, the provider project is used.`,
			},

			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Region of function. Currently can be only "us-central1". If it is not provided, the provider region is used.`,
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

	if v, ok := d.GetOk("ingress_settings"); ok {
		function.IngressSettings = v.(string)
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

	if v, ok := d.GetOk("vpc_connector_egress_settings"); ok {
		function.VpcConnectorEgressSettings = v.(string)
	}

	if v, ok := d.GetOk("max_instances"); ok {
		function.MaxInstances = int64(v.(int))
	}

	log.Printf("[DEBUG] Creating cloud function: %s", function.Name)

	// We retry the whole create-and-wait because Cloud Functions
	// will sometimes fail a creation operation entirely if it fails to pull
	// source code and we need to try the whole creation again.
	rerr := retryTimeDuration(func() error {
		op, err := config.clientCloudFunctions.Projects.Locations.Functions.Create(
			cloudFuncId.locationId(), function).Do()
		if err != nil {
			return err
		}

		// Name of function should be unique
		d.SetId(cloudFuncId.cloudFunctionId())

		return cloudFunctionsOperationWait(config, op, "Creating CloudFunctions Function",
			d.Timeout(schema.TimeoutCreate))
	}, d.Timeout(schema.TimeoutCreate), isCloudFunctionsSourceCodeError)
	if rerr != nil {
		return rerr
	}
	log.Printf("[DEBUG] Finished creating cloud function: %s", function.Name)
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
	d.Set("ingress_settings", function.IngressSettings)
	d.Set("labels", function.Labels)
	d.Set("runtime", function.Runtime)
	d.Set("service_account_email", function.ServiceAccountEmail)
	d.Set("environment_variables", function.EnvironmentVariables)
	d.Set("vpc_connector", function.VpcConnector)
	d.Set("vpc_connector_egress_settings", function.VpcConnectorEgressSettings)
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

	if d.HasChange("ingress_settings") {
		function.IngressSettings = d.Get("ingress_settings").(string)
		updateMaskArr = append(updateMaskArr, "ingressSettings")
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

	if d.HasChange("vpc_connector_egress_settings") {
		function.VpcConnectorEgressSettings = d.Get("vpc_connector_egress_settings").(string)
		updateMaskArr = append(updateMaskArr, "vpcConnectorEgressSettings")
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
			d.Timeout(schema.TimeoutUpdate))
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
		d.Timeout(schema.TimeoutDelete))
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
