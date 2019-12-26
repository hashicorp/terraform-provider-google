package google

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/servicemanagement/v1"
)

func resourceEndpointsService() *schema.Resource {
	return &schema.Resource{
		Create: resourceEndpointsServiceCreate,
		Read:   resourceEndpointsServiceRead,
		Delete: resourceEndpointsServiceDelete,
		Update: resourceEndpointsServiceUpdate,

		// Migrates protoc_output -> protoc_output_base64.
		SchemaVersion: 1,
		MigrateState:  migrateEndpointsService,

		Schema: map[string]*schema.Schema{
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"openapi_config": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"grpc_config", "protoc_output_base64"},
			},
			"grpc_config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protoc_output_base64": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"config_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"apis": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"syntax": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"methods": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"syntax": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"request_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"response_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"dns_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func getEndpointServiceOpenAPIConfigSource(configText string) *servicemanagement.ConfigSource {
	// We need to provide a ConfigSource object to the API whenever submitting a
	// new config.  A ConfigSource contains a ConfigFile which contains the b64
	// encoded contents of the file.  OpenAPI requires only one file.
	configfile := servicemanagement.ConfigFile{
		FileContents: base64.StdEncoding.EncodeToString([]byte(configText)),
		FileType:     "OPEN_API_YAML",
		FilePath:     "heredoc.yaml",
	}
	return &servicemanagement.ConfigSource{
		Files: []*servicemanagement.ConfigFile{&configfile},
	}
}

func getEndpointServiceGRPCConfigSource(serviceConfig, protoConfig string) *servicemanagement.ConfigSource {
	// gRPC requires both the file specifying the service and the compiled protobuf,
	// but they can be in any order.
	ymlConfigfile := servicemanagement.ConfigFile{
		FileContents: base64.StdEncoding.EncodeToString([]byte(serviceConfig)),
		FileType:     "SERVICE_CONFIG_YAML",
		FilePath:     "heredoc.yaml",
	}
	protoConfigfile := servicemanagement.ConfigFile{
		FileContents: protoConfig,
		FileType:     "FILE_DESCRIPTOR_SET_PROTO",
		FilePath:     "api_def.pb",
	}
	return &servicemanagement.ConfigSource{
		Files: []*servicemanagement.ConfigFile{&ymlConfigfile, &protoConfigfile},
	}
}

func resourceEndpointsServiceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	// If the service doesn't exist, we'll need to create it, but if it does, it
	// will be reused.  This is unusual for Terraform, but it causes the behavior
	// that users will want and accept.  Users of Endpoints are not thinking in
	// terms of services, configs, and rollouts - they just want the setup declared
	// in their config to happen.  The fact that a service may need to be created
	// is not interesting to them.  Consequently, we create this service if necessary
	// so that we can perform the rollout without further disruption, which is the
	// action that a user running `terraform apply` is going to want.
	serviceName := d.Get("service_name").(string)
	log.Printf("[DEBUG] Create Endpoint Service %q", serviceName)

	log.Printf("[DEBUG] Checking for existing ManagedService %q", serviceName)
	_, err = config.clientServiceMan.Services.Get(serviceName).Do()
	if err != nil {
		log.Printf("[DEBUG] Creating new ServiceManagement ManagedService %q", serviceName)
		op, err := config.clientServiceMan.Services.Create(
			&servicemanagement.ManagedService{
				ProducerProjectId: project,
				ServiceName:       serviceName,
			}).Do()
		if err != nil {
			return err
		}

		_, err = serviceManagementOperationWait(config, op, "Creating new ManagedService.")
		if err != nil {
			return err
		}
	}

	// Use update to set other fields like config.
	err = resourceEndpointsServiceUpdate(d, meta)
	if err != nil {
		return err
	}

	d.SetId(serviceName)
	return resourceEndpointsServiceRead(d, meta)
}

func expandEndpointServiceConfigSource(d *schema.ResourceData, meta interface{}) (*servicemanagement.ConfigSource, error) {
	if openapiConfig, ok := d.GetOk("openapi_config"); ok {
		return getEndpointServiceOpenAPIConfigSource(openapiConfig.(string)), nil
	}

	grpcConfig, gok := d.GetOk("grpc_config")
	protocOutput, pok := d.GetOk("protoc_output_base64")
	if gok && pok {
		return getEndpointServiceGRPCConfigSource(grpcConfig.(string), protocOutput.(string)), nil
	}

	return nil, errors.New("Could not parse config - either openapi_config or both grpc_config and protoc_output_base64 must be set.")
}

func resourceEndpointsServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	// This update is not quite standard for a terraform resource.  Instead of
	// using the go client library to send an HTTP request to update something
	// serverside, we have to push a new configuration, wait for it to be
	// parsed and loaded, then create and push a rollout and wait for that
	// rollout to be completed.
	// There's a lot of moving parts there, and all of them have knobs that can
	// be tweaked if the user is using gcloud.  In the interest of simplicity,
	// we currently only support full rollouts - anyone trying to do incremental
	// rollouts or A/B testing is going to need a more precise tool than this resource.
	config := meta.(*Config)
	serviceName := d.Get("service_name").(string)

	log.Printf("[DEBUG] Updating ManagedService %q", serviceName)

	cfgSource, err := expandEndpointServiceConfigSource(d, meta)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating ManagedService %q", serviceName)
	// The difference between "submit" and "create" is that submit parses the config
	// you provide, where "create" requires the config in a pre-parsed format.
	// "submit" will be a lot more flexible for users and will always be up-to-date
	// with any new features that arise - this is why you provide a YAML config
	// instead of providing the config in HCL.
	log.Printf("[DEBUG] Submitting config for ManagedService %q", serviceName)
	op, err := config.clientServiceMan.Services.Configs.Submit(
		serviceName,
		&servicemanagement.SubmitConfigSourceRequest{
			ConfigSource: cfgSource,
		}).Do()
	if err != nil {
		return err
	}
	s, err := serviceManagementOperationWait(config, op, "Submitting service config.")
	if err != nil {
		return err
	}
	var serviceConfig servicemanagement.SubmitConfigSourceResponse
	if err := json.Unmarshal(s, &serviceConfig); err != nil {
		return err
	}

	// Next, we create a new rollout with the new config value, and wait for it to complete.
	rollout := servicemanagement.Rollout{
		ServiceName: serviceName,
		TrafficPercentStrategy: &servicemanagement.TrafficPercentStrategy{
			Percentages: map[string]float64{serviceConfig.ServiceConfig.Id: 100.0},
		},
	}

	log.Printf("[DEBUG] Creating new rollout for ManagedService %q", serviceName)
	op, err = config.clientServiceMan.Services.Rollouts.Create(serviceName, &rollout).Do()
	if err != nil {
		return err
	}
	_, err = serviceManagementOperationWait(config, op, "Performing service rollout.")
	if err != nil {
		return err
	}

	return resourceEndpointsServiceRead(d, meta)
}

func resourceEndpointsServiceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	log.Printf("[DEBUG] Deleting ManagedService %q", d.Id())

	op, err := config.clientServiceMan.Services.Delete(d.Get("service_name").(string)).Do()
	if err != nil {
		return err
	}
	_, err = serviceManagementOperationWait(config, op, "Deleting service.")
	d.SetId("")
	return err
}

func resourceEndpointsServiceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	log.Printf("[DEBUG] Reading ManagedService %q", d.Id())

	service, err := config.clientServiceMan.Services.GetConfig(d.Get("service_name").(string)).Do()
	if err != nil {
		return err
	}

	d.Set("config_id", service.Id)
	d.Set("dns_address", service.Name)
	d.Set("apis", flattenServiceManagementAPIs(service.Apis))
	d.Set("endpoints", flattenServiceManagementEndpoints(service.Endpoints))

	return nil
}

func flattenServiceManagementAPIs(apis []*servicemanagement.Api) []map[string]interface{} {
	flattened := make([]map[string]interface{}, len(apis))
	for i, a := range apis {
		flattened[i] = map[string]interface{}{
			"name":    a.Name,
			"version": a.Version,
			"syntax":  a.Syntax,
			"methods": flattenServiceManagementMethods(a.Methods),
		}
	}
	return flattened
}

func flattenServiceManagementMethods(methods []*servicemanagement.Method) []map[string]interface{} {
	flattened := make([]map[string]interface{}, len(methods))
	for i, m := range methods {
		flattened[i] = map[string]interface{}{
			"name":          m.Name,
			"syntax":        m.Syntax,
			"request_type":  m.RequestTypeUrl,
			"response_type": m.ResponseTypeUrl,
		}
	}
	return flattened
}

func flattenServiceManagementEndpoints(endpoints []*servicemanagement.Endpoint) []map[string]interface{} {
	flattened := make([]map[string]interface{}, len(endpoints))
	for i, e := range endpoints {
		flattened[i] = map[string]interface{}{
			"name":    e.Name,
			"address": e.Target,
		}
	}
	return flattened
}
