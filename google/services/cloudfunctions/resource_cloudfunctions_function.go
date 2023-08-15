// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package cloudfunctions

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"github.com/hashicorp/terraform-provider-google/google/verify"

	"google.golang.org/api/cloudfunctions/v1"

	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var allowedIngressSettings = []string{
	"ALLOW_ALL",
	"ALLOW_INTERNAL_AND_GCLB",
	"ALLOW_INTERNAL_ONLY",
}

var allowedVpcConnectorEgressSettings = []string{
	"ALL_TRAFFIC",
	"PRIVATE_RANGES_ONLY",
}

type CloudFunctionId struct {
	Project string
	Region  string
	Name    string
}

func (s *CloudFunctionId) CloudFunctionId() string {
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

func (s *CloudFunctionId) locationId() string {
	return fmt.Sprintf("projects/%s/locations/%s", s.Project, s.Region)
}

func parseCloudFunctionId(d *schema.ResourceData, config *transport_tpg.Config) (*CloudFunctionId, error) {
	if err := tpgresource.ParseImportId([]string{
		"projects/(?P<project>[^/]+)/locations/(?P<region>[^/]+)/functions/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<region>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}
	return &CloudFunctionId{
		Project: d.Get("project").(string),
		Region:  d.Get("region").(string),
		Name:    d.Get("name").(string),
	}, nil
}

// Differs from validateGCEName because Cloud Functions allow capital letters
// at start/end
func validateResourceCloudFunctionsFunctionName(v interface{}, k string) (ws []string, errors []error) {
	re := `^(?:[a-zA-Z](?:[-_a-zA-Z0-9]{0,61}[a-zA-Z0-9])?)$`
	return verify.ValidateRegexp(re)(v, k)
}

func partsCompare(a, b, reg string) bool {

	regex := regexp.MustCompile(reg)
	if regex.MatchString(a) && regex.MatchString(b) {
		aParts := regex.FindStringSubmatch(a)
		bParts := regex.FindStringSubmatch(b)
		for i := 0; i < len(aParts); i++ {
			if aParts[i] != bParts[i] {
				return false
			}
		}
	} else if regex.MatchString(a) {
		aParts := regex.FindStringSubmatch(a)
		if aParts[len(aParts)-1] != b {
			return false
		}
	} else if regex.MatchString(b) {
		bParts := regex.FindStringSubmatch(b)
		if bParts[len(bParts)-1] != a {
			return false
		}
	} else {
		if a != b {
			return false
		}
	}

	return true
}

// based on CompareSelfLinkOrResourceName, but less reusable and allows multi-/
// strings in the new state (config) part
func compareSelfLinkOrResourceNameWithMultipleParts(_, old, new string, _ *schema.ResourceData) bool {
	// two formats based on expandEventTrigger()
	regex1 := "projects/(.+)/databases/\\(default\\)/documents/(.+)"
	regex2 := "projects/(.+)/(.+)/(.+)"
	return partsCompare(old, new, regex1) || partsCompare(old, new, regex2)
}

func ResourceCloudFunctionsFunction() *schema.Resource {
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

			"build_worker_pool": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Name of the Cloud Build Custom Worker Pool that should be used to build the function.`,
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

			"docker_registry": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Docker Registry to use for storing the function's Docker images. Allowed values are CONTAINER_REGISTRY (default) and ARTIFACT_REGISTRY.`,
			},

			"docker_repository": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `User managed repository created in Artifact Registry optionally with a customer managed encryption key. If specified, deployments will use Artifact Registry for storing images built with Cloud Build.`,
			},

			"kms_key_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Resource name of a KMS crypto key (managed by the user) used to encrypt/decrypt function resources.`,
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Description of the function.`,
			},

			"available_memory_mb": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     256,
				Description: `Memory (in MB), available to the function. Default value is 256. Possible values include 128, 256, 512, 1024, etc.`,
			},

			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      60,
				ValidateFunc: validation.IntBetween(1, 540),
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
				Default:      "ALLOW_ALL",
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
				Description: `The runtime in which the function is going to run. Eg. "nodejs12", "nodejs14", "python37", "go111".`,
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
				DiffSuppressFunc: tpgresource.CompareSelfLinkOrResourceName,
				Description:      `The VPC Network Connector that this cloud function can connect to. It can be either the fully-qualified URI, or the short name of the network connector resource. The format of this field is projects/*/locations/*/connectors/*.`,
			},

			"environment_variables": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: `A set of key/value environment variable pairs to assign to the function.`,
			},

			"build_environment_variables": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: ` A set of key/value environment variable pairs available during build time.`,
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

			"https_trigger_security_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The security level for the function. Defaults to SECURE_OPTIONAL. Valid only if trigger_http is used.`,
			},

			"max_instances": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(0),
				Description:  `The limit on the maximum number of function instances that may coexist at a given time.`,
			},

			"min_instances": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
				Description:  `The limit on the minimum number of function instances that may coexist at a given time.`,
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
				Description: `Region of function. If it is not provided, the provider region is used.`,
			},

			"secret_environment_variables": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Secret environment variables configuration`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Name of the environment variable.`,
						},
						"project_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Project identifier (due to a known limitation, only project number is supported by this field) of the project that contains the secret. If not set, it will be populated with the function's project, assuming that the secret exists in the same project as of the function.`,
						},
						"secret": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `ID of the secret in secret manager (not the full resource name).`,
						},
						"version": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Version of the secret (version number or the string "latest"). It is recommended to use a numeric version for secret environment variables as any updates to the secret value is not reflected until new clones start.`,
						},
					},
				},
			},

			"secret_volumes": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Secret volumes configuration.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_path": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The path within the container to mount the secret volume. For example, setting the mount_path as "/etc/secrets" would mount the secret value files under the "/etc/secrets" directory. This directory will also be completely shadowed and unavailable to mount any other secrets. Recommended mount paths: "/etc/secrets" Restricted mount paths: "/cloudsql", "/dev/log", "/pod", "/proc", "/var/log".`,
						},
						"project_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Project identifier (due to a known limitation, only project number is supported by this field) of the project that contains the secret. If not set, it will be populated with the function's project, assuming that the secret exists in the same project as of the function.`,
						},
						"secret": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `ID of the secret in secret manager (not the full resource name).`,
						},
						"versions": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `List of secret versions to mount for this secret. If empty, the "latest" version of the secret will be made available in a file named after the secret under the mount point.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Relative path of the file under the mount path where the secret value for this version will be fetched and made available. For example, setting the mount_path as "/etc/secrets" and path as "/secret_foo" would mount the secret value file at "/etc/secrets/secret_foo".`,
									},
									"version": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Version of the secret (version number or the string "latest"). It is preferable to use "latest" version with secret volumes as secret value changes are reflected immediately.`,
									},
								},
							},
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Describes the current stage of a deployment.`,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceCloudFunctionsCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	region, err := tpgresource.GetRegion(d, config)
	if err != nil {
		return err
	}

	cloudFuncId := &CloudFunctionId{
		Project: project,
		Region:  region,
		Name:    d.Get("name").(string),
	}

	function := &cloudfunctions.CloudFunction{
		Name:                cloudFuncId.CloudFunctionId(),
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

	secretEnv := d.Get("secret_environment_variables").([]interface{})
	if len(secretEnv) > 0 {
		function.SecretEnvironmentVariables = expandSecretEnvironmentVariables(secretEnv)
	}

	secretVolume := d.Get("secret_volumes").([]interface{})
	if len(secretVolume) > 0 {
		function.SecretVolumes = expandSecretVolumes(secretVolume)
	}

	if v, ok := d.GetOk("available_memory_mb"); ok {
		availableMemoryMb := v.(int)
		function.AvailableMemoryMb = int64(availableMemoryMb)
	}

	if v, ok := d.GetOk("description"); ok {
		function.Description = v.(string)
	}

	if v, ok := d.GetOk("build_worker_pool"); ok {
		function.BuildWorkerPool = v.(string)
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
		function.HttpsTrigger.SecurityLevel = d.Get("https_trigger_security_level").(string)
	} else {
		return fmt.Errorf("One of `event_trigger` or `trigger_http` is required: " +
			"You must specify a trigger when deploying a new function.")
	}

	if v, ok := d.GetOk("ingress_settings"); ok {
		function.IngressSettings = v.(string)
	}

	if _, ok := d.GetOk("labels"); ok {
		function.Labels = tpgresource.ExpandLabels(d)
	}

	if _, ok := d.GetOk("environment_variables"); ok {
		function.EnvironmentVariables = tpgresource.ExpandEnvironmentVariables(d)
	}

	if _, ok := d.GetOk("build_environment_variables"); ok {
		function.BuildEnvironmentVariables = tpgresource.ExpandBuildEnvironmentVariables(d)
	}

	if v, ok := d.GetOk("vpc_connector"); ok {
		function.VpcConnector = v.(string)
	}

	if v, ok := d.GetOk("vpc_connector_egress_settings"); ok {
		function.VpcConnectorEgressSettings = v.(string)
	}

	if v, ok := d.GetOk("docker_registry"); ok {
		function.DockerRegistry = v.(string)
	}

	if v, ok := d.GetOk("docker_repository"); ok {
		function.DockerRepository = v.(string)
	}

	if v, ok := d.GetOk("kms_key_name"); ok {
		function.KmsKeyName = v.(string)
	}

	if v, ok := d.GetOk("max_instances"); ok {
		function.MaxInstances = int64(v.(int))
	}

	if v, ok := d.GetOk("min_instances"); ok {
		function.MinInstances = int64(v.(int))
	}

	log.Printf("[DEBUG] Creating cloud function: %s", function.Name)

	// We retry the whole create-and-wait because Cloud Functions
	// will sometimes fail a creation operation entirely if it fails to pull
	// source code and we need to try the whole creation again.
	rerr := transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: func() error {
			op, err := config.NewCloudFunctionsClient(userAgent).Projects.Locations.Functions.Create(
				cloudFuncId.locationId(), function).Do()
			if err != nil {
				return err
			}

			// Name of function should be unique
			d.SetId(cloudFuncId.CloudFunctionId())

			return CloudFunctionsOperationWait(config, op, "Creating CloudFunctions Function", userAgent,
				d.Timeout(schema.TimeoutCreate))
		},
		Timeout:              d.Timeout(schema.TimeoutCreate),
		ErrorRetryPredicates: []transport_tpg.RetryErrorPredicateFunc{IsCloudFunctionsSourceCodeError},
	})
	if rerr != nil {
		return rerr
	}
	log.Printf("[DEBUG] Finished creating cloud function: %s", function.Name)
	return resourceCloudFunctionsRead(d, meta)
}

func resourceCloudFunctionsRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	cloudFuncId, err := parseCloudFunctionId(d, config)
	if err != nil {
		return err
	}

	function, err := config.NewCloudFunctionsClient(userAgent).Projects.Locations.Functions.Get(cloudFuncId.CloudFunctionId()).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Target CloudFunctions Function %q", cloudFuncId.Name))
	}

	if err := d.Set("name", cloudFuncId.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err := d.Set("description", function.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err := d.Set("build_worker_pool", function.BuildWorkerPool); err != nil {
		return fmt.Errorf("Error setting build_worker_pool: %s", err)
	}
	if err := d.Set("entry_point", function.EntryPoint); err != nil {
		return fmt.Errorf("Error setting entry_point: %s", err)
	}
	if err := d.Set("available_memory_mb", function.AvailableMemoryMb); err != nil {
		return fmt.Errorf("Error setting available_memory_mb: %s", err)
	}
	sRemoved := strings.Replace(function.Timeout, "s", "", -1)
	timeout, err := strconv.Atoi(sRemoved)
	if err != nil {
		return err
	}
	if err := d.Set("timeout", timeout); err != nil {
		return fmt.Errorf("Error setting timeout: %s", err)
	}
	if err := d.Set("ingress_settings", function.IngressSettings); err != nil {
		return fmt.Errorf("Error setting ingress_settings: %s", err)
	}
	if err := d.Set("labels", function.Labels); err != nil {
		return fmt.Errorf("Error setting labels: %s", err)
	}
	if err := d.Set("runtime", function.Runtime); err != nil {
		return fmt.Errorf("Error setting runtime: %s", err)
	}
	if err := d.Set("service_account_email", function.ServiceAccountEmail); err != nil {
		return fmt.Errorf("Error setting service_account_email: %s", err)
	}
	if err := d.Set("environment_variables", function.EnvironmentVariables); err != nil {
		return fmt.Errorf("Error setting environment_variables: %s", err)
	}
	if err := d.Set("vpc_connector", function.VpcConnector); err != nil {
		return fmt.Errorf("Error setting vpc_connector: %s", err)
	}
	if err := d.Set("vpc_connector_egress_settings", function.VpcConnectorEgressSettings); err != nil {
		return fmt.Errorf("Error setting vpc_connector_egress_settings: %s", err)
	}
	if function.SourceArchiveUrl != "" {
		// sourceArchiveUrl should always be a Google Cloud Storage URL (e.g. gs://bucket/object)
		// https://cloud.google.com/functions/docs/reference/rest/v1/projects.locations.functions
		sourceURL, err := url.Parse(function.SourceArchiveUrl)
		if err != nil {
			return err
		}
		bucket := sourceURL.Host
		object := strings.TrimLeft(sourceURL.Path, "/")
		if err := d.Set("source_archive_bucket", bucket); err != nil {
			return fmt.Errorf("Error setting source_archive_bucket: %s", err)
		}
		if err := d.Set("source_archive_object", object); err != nil {
			return fmt.Errorf("Error setting source_archive_object: %s", err)
		}
	}
	if err := d.Set("source_repository", flattenSourceRepository(function.SourceRepository)); err != nil {
		return fmt.Errorf("Error setting source_repository: %s", err)
	}

	if err := d.Set("secret_environment_variables", flattenSecretEnvironmentVariables(function.SecretEnvironmentVariables)); err != nil {
		return fmt.Errorf("Error setting secret_environment_variables: %s", err)
	}

	if err := d.Set("secret_volumes", flattenSecretVolumes(function.SecretVolumes)); err != nil {
		return fmt.Errorf("Error setting secret_volumes: %s", err)
	}

	if err := d.Set("status", function.Status); err != nil {
		return fmt.Errorf("Error setting status: %s", err)
	}

	if function.HttpsTrigger != nil {
		if err := d.Set("trigger_http", true); err != nil {
			return fmt.Errorf("Error setting trigger_http: %s", err)
		}
		if err := d.Set("https_trigger_url", function.HttpsTrigger.Url); err != nil {
			return fmt.Errorf("Error setting https_trigger_url: %s", err)
		}
		if err := d.Set("https_trigger_security_level", function.HttpsTrigger.SecurityLevel); err != nil {
			return fmt.Errorf("Error setting https_trigger_security_level: %s", err)
		}
	}

	if err := d.Set("event_trigger", flattenEventTrigger(function.EventTrigger)); err != nil {
		return fmt.Errorf("Error setting event_trigger: %s", err)
	}
	if err := d.Set("docker_registry", function.DockerRegistry); err != nil {
		return fmt.Errorf("Error setting docker_registry: %s", err)
	}
	if err := d.Set("docker_repository", function.DockerRepository); err != nil {
		return fmt.Errorf("Error setting docker_repository: %s", err)
	}
	if err := d.Set("kms_key_name", function.KmsKeyName); err != nil {
		return fmt.Errorf("Error setting kms_key_name: %s", err)
	}
	if err := d.Set("max_instances", function.MaxInstances); err != nil {
		return fmt.Errorf("Error setting max_instances: %s", err)
	}
	if err := d.Set("min_instances", function.MinInstances); err != nil {
		return fmt.Errorf("Error setting min_instances: %s", err)
	}
	if err := d.Set("region", cloudFuncId.Region); err != nil {
		return fmt.Errorf("Error setting region: %s", err)
	}
	if err := d.Set("project", cloudFuncId.Project); err != nil {
		return fmt.Errorf("Error setting project: %s", err)
	}

	return nil
}

func resourceCloudFunctionsUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG]: Updating google_cloudfunctions_function")
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	project, err := tpgresource.GetProject(d, config)
	if err != nil {
		return err
	}

	cloudFuncId, err := parseCloudFunctionId(d, config)
	if err != nil {
		return err
	}

	// The full function needs to supplied in the PATCH call to evaluate some Organization Policies. https://github.com/hashicorp/terraform-provider-google/issues/6603
	function, err := config.NewCloudFunctionsClient(userAgent).Projects.Locations.Functions.Get(cloudFuncId.CloudFunctionId()).Do()
	if err != nil {
		return transport_tpg.HandleNotFoundError(err, d, fmt.Sprintf("Target CloudFunctions Function %q", cloudFuncId.Name))
	}

	// The full function may contain a reference to manually uploaded code if the function was imported from gcloud
	// This does not work with Terraform, so zero it out from the function if it exists. See https://github.com/hashicorp/terraform-provider-google/issues/7921
	function.SourceUploadUrl = ""

	d.Partial(true)

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

	if d.HasChange("secret_environment_variables") {
		function.SecretEnvironmentVariables = expandSecretEnvironmentVariables(d.Get("secret_environment_variables").([]interface{}))
		updateMaskArr = append(updateMaskArr, "secretEnvironmentVariables")
	}

	if d.HasChange("secret_volumes") {
		function.SecretVolumes = expandSecretVolumes(d.Get("secret_volumes").([]interface{}))
		updateMaskArr = append(updateMaskArr, "secretVolumes")
	}

	if d.HasChange("description") {
		function.Description = d.Get("description").(string)
		updateMaskArr = append(updateMaskArr, "description")
	}

	if d.HasChange("build_worker_pool") {
		function.BuildWorkerPool = d.Get("build_worker_pool").(string)
		updateMaskArr = append(updateMaskArr, "build_worker_pool")
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
		function.Labels = tpgresource.ExpandLabels(d)
		updateMaskArr = append(updateMaskArr, "labels")
	}

	if d.HasChange("runtime") {
		function.Runtime = d.Get("runtime").(string)
		updateMaskArr = append(updateMaskArr, "runtime")
	}

	if d.HasChange("environment_variables") {
		function.EnvironmentVariables = tpgresource.ExpandEnvironmentVariables(d)
		updateMaskArr = append(updateMaskArr, "environmentVariables")
	}

	if d.HasChange("build_environment_variables") {
		function.BuildEnvironmentVariables = tpgresource.ExpandBuildEnvironmentVariables(d)
		updateMaskArr = append(updateMaskArr, "buildEnvironmentVariables")
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

	if d.HasChange("https_trigger_security_level") {
		if function.HttpsTrigger == nil {
			function.HttpsTrigger = &cloudfunctions.HttpsTrigger{}
		}
		function.HttpsTrigger.SecurityLevel = d.Get("https_trigger_security_level").(string)
		updateMaskArr = append(updateMaskArr, "httpsTrigger", "httpsTrigger.securityLevel")
	}

	if d.HasChange("docker_registry") {
		function.DockerRegistry = d.Get("docker_registry").(string)
		updateMaskArr = append(updateMaskArr, "dockerRegistry")
	}

	if d.HasChange("docker_repository") {
		function.DockerRepository = d.Get("docker_repository").(string)
		updateMaskArr = append(updateMaskArr, "dockerRepository")
	}

	if d.HasChange("kms_key_name") {
		function.KmsKeyName = d.Get("kms_key_name").(string)
		updateMaskArr = append(updateMaskArr, "kmsKeyName")
	}

	if d.HasChange("max_instances") {
		function.MaxInstances = int64(d.Get("max_instances").(int))
		updateMaskArr = append(updateMaskArr, "maxInstances")
	}

	if d.HasChange("min_instances") {
		function.MinInstances = int64(d.Get("min_instances").(int))
		updateMaskArr = append(updateMaskArr, "minInstances")
	}

	if len(updateMaskArr) > 0 {
		log.Printf("[DEBUG] Send Patch CloudFunction Configuration request: %#v", function)
		updateMask := strings.Join(updateMaskArr, ",")
		rerr := transport_tpg.Retry(transport_tpg.RetryOptions{
			RetryFunc: func() error {
				op, err := config.NewCloudFunctionsClient(userAgent).Projects.Locations.Functions.Patch(function.Name, function).
					UpdateMask(updateMask).Do()
				if err != nil {
					return err
				}

				return CloudFunctionsOperationWait(config, op, "Updating CloudFunctions Function", userAgent,
					d.Timeout(schema.TimeoutUpdate))
			},
			Timeout: d.Timeout(schema.TimeoutUpdate),
		})
		if rerr != nil {
			return fmt.Errorf("Error while updating cloudfunction configuration: %s", rerr)
		}
	}
	d.Partial(false)

	return resourceCloudFunctionsRead(d, meta)
}

func resourceCloudFunctionsDestroy(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*transport_tpg.Config)
	userAgent, err := tpgresource.GenerateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	cloudFuncId, err := parseCloudFunctionId(d, config)
	if err != nil {
		return err
	}

	op, err := config.NewCloudFunctionsClient(userAgent).Projects.Locations.Functions.Delete(cloudFuncId.CloudFunctionId()).Do()
	if err != nil {
		return err
	}
	err = CloudFunctionsOperationWait(config, op, "Deleting CloudFunctions Function", userAgent,
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

func expandSecretEnvironmentVariables(configured []interface{}) []*cloudfunctions.SecretEnvVar {
	if len(configured) == 0 {
		return nil
	}
	result := make([]*cloudfunctions.SecretEnvVar, 0, len(configured))
	for _, e := range configured {
		data := e.(map[string]interface{})
		result = append(result, &cloudfunctions.SecretEnvVar{
			Key:       data["key"].(string),
			ProjectId: data["project_id"].(string),
			Secret:    data["secret"].(string),
			Version:   data["version"].(string),
		})
	}
	return result
}

func flattenSecretEnvironmentVariables(envVars []*cloudfunctions.SecretEnvVar) []map[string]interface{} {
	if envVars == nil {
		return nil
	}
	var result []map[string]interface{}

	for _, envVar := range envVars {
		if envVar != nil {
			data := map[string]interface{}{
				"key":        envVar.Key,
				"project_id": envVar.ProjectId,
				"secret":     envVar.Secret,
				"version":    envVar.Version,
			}
			result = append(result, data)
		}
	}
	return result
}

func expandSecretVolumes(configured []interface{}) []*cloudfunctions.SecretVolume {
	if len(configured) == 0 {
		return nil
	}
	result := make([]*cloudfunctions.SecretVolume, 0, len(configured))
	for _, e := range configured {
		data := e.(map[string]interface{})
		result = append(result, &cloudfunctions.SecretVolume{
			MountPath: data["mount_path"].(string),
			ProjectId: data["project_id"].(string),
			Secret:    data["secret"].(string),
			Versions:  expandSecretVersion(data["versions"].([]interface{})), //TODO
		})
	}
	return result
}

func flattenSecretVolumes(secretVolumes []*cloudfunctions.SecretVolume) []map[string]interface{} {
	if secretVolumes == nil {
		return nil
	}
	var result []map[string]interface{}

	for _, secretVolume := range secretVolumes {
		if secretVolume != nil {
			data := map[string]interface{}{
				"mount_path": secretVolume.MountPath,
				"project_id": secretVolume.ProjectId,
				"secret":     secretVolume.Secret,
				"versions":   flattenSecretVersion(secretVolume.Versions),
			}
			result = append(result, data)
		}
	}
	return result
}

func expandSecretVersion(configured []interface{}) []*cloudfunctions.SecretVersion {
	if len(configured) == 0 {
		return nil
	}
	result := make([]*cloudfunctions.SecretVersion, 0, len(configured))
	for _, e := range configured {
		data := e.(map[string]interface{})
		result = append(result, &cloudfunctions.SecretVersion{
			Path:    data["path"].(string),
			Version: data["version"].(string),
		})
	}
	return result
}

func flattenSecretVersion(secretVersions []*cloudfunctions.SecretVersion) []map[string]interface{} {
	if secretVersions == nil {
		return nil
	}
	var result []map[string]interface{}

	for _, secretVersion := range secretVersions {
		if secretVersion != nil {
			data := map[string]interface{}{
				"path":    secretVersion.Path,
				"version": secretVersion.Version,
			}
			result = append(result, data)
		}
	}
	return result
}
