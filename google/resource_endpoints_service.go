package google

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/servicemanagement/v1"
)

func resourceEndpointsService() *schema.Resource {
	return &schema.Resource{
		Create: resourceEndpointsServiceCreate,
		Read:   resourceEndpointsServiceRead,
		Delete: resourceEndpointsServiceDelete,
		Update: resourceEndpointsServiceUpdate,
		Schema: map[string]*schema.Schema{
			"service_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"openapi_config": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"grpc_config": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"protoc_output": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"config_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"apis": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"syntax": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"version": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"methods": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"syntax": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"request_type": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"response_type": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"dns_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoints": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"address": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func getOpenApiConfigSource(config_text string) servicemanagement.ConfigSource {
	// We need to provide a ConfigSource object to the API whenever submitting a
	// new config.  A ConfigSource contains a ConfigFile which contains the b64
	// encoded contents of the file.  OpenAPI requires only one file.
	configfile := servicemanagement.ConfigFile{
		FileContents: base64.StdEncoding.EncodeToString([]byte(config_text)),
		FileType:     "OPEN_API_YAML",
		FilePath:     "heredoc.yaml",
	}
	return servicemanagement.ConfigSource{
		Files: []*servicemanagement.ConfigFile{&configfile},
	}
}

func getGrpcConfigSource(service_config, proto_config string) servicemanagement.ConfigSource {
	// gRPC requires both the file specifying the service and the compiled protobuf,
	// but they can be in any order.
	yml_configfile := servicemanagement.ConfigFile{
		FileContents: base64.StdEncoding.EncodeToString([]byte(service_config)),
		FileType:     "SERVICE_CONFIG_YAML",
		FilePath:     "heredoc.yaml",
	}
	proto_configfile := servicemanagement.ConfigFile{
		FileContents: base64.StdEncoding.EncodeToString([]byte(proto_config)),
		FileType:     "FILE_DESCRIPTOR_SET_PROTO",
		FilePath:     "api_def.pb",
	}
	return servicemanagement.ConfigSource{
		Files: []*servicemanagement.ConfigFile{&yml_configfile, &proto_configfile},
	}
}

func resourceEndpointsServiceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	// If the service doesn't exist, we'll need to create it, but if it does, it
	// will be reused.
	serviceName := d.Get("service_name").(string)
	servicesService := servicemanagement.NewServicesService(config.clientServiceMan)
	_, err = servicesService.Get(serviceName).Do()
	if err != nil {
		_, err = servicesService.Create(&servicemanagement.ManagedService{ProducerProjectId: project, ServiceName: serviceName}).Do()
		if err != nil {
			return err
		}
	}
	// Do a rollout using the update mechanism.
	err = resourceEndpointsServiceUpdate(d, meta)
	if err != nil {
		return err
	}

	d.SetId(serviceName)
	return resourceEndpointsServiceRead(d, meta)
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
	openapi_config, ok := d.GetOk("openapi_config")
	var source servicemanagement.ConfigSource
	if ok {
		source = getOpenApiConfigSource(openapi_config.(string))
	} else {
		grpc_config, g_ok := d.GetOk("grpc_config")
		protoc_output, p_ok := d.GetOk("protoc_output")
		if g_ok && p_ok {
			source = getGrpcConfigSource(grpc_config.(string), protoc_output.(string))
		} else {
			return errors.New("Could not decypher config - please either set openapi_config or set both grpc_config and protoc_output.")
		}
	}

	configService := servicemanagement.NewServicesConfigsService(config.clientServiceMan)
	// The difference between "submit" and "create" is that submit parses the config
	// you provide, where "create" requires the config in a pre-parsed format.
	// "submit" will be a lot more flexible for users and will always be up-to-date
	// with any new features that arise - this is why you provide a YAML config
	// instead of providing the config in HCL.
	op, err := configService.Submit(serviceName, &servicemanagement.SubmitConfigSourceRequest{ConfigSource: &source}).Do()
	if err != nil {
		return err
	}
	s, err := serviceManagementOperationWait(config, op, "Submitting service config.")
	if err != nil {
		return err
	}
	var serviceConfig servicemanagement.SubmitConfigSourceResponse
	json.Unmarshal(s, &serviceConfig)

	// Next, we create a new rollout with the new config value, and wait for it to complete.
	rolloutService := servicemanagement.NewServicesRolloutsService(config.clientServiceMan)
	rollout := servicemanagement.Rollout{
		ServiceName: serviceName,
		TrafficPercentStrategy: &servicemanagement.TrafficPercentStrategy{
			Percentages: map[string]float64{serviceConfig.ServiceConfig.Id: 100.0},
		},
	}
	op, err = rolloutService.Create(serviceName, &rollout).Do()
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
	servicesService := servicemanagement.NewServicesService(config.clientServiceMan)
	op, err := servicesService.Delete(d.Get("service_name").(string)).Do()
	if err != nil {
		return err
	}
	_, err = serviceManagementOperationWait(config, op, "Deleting service.")
	d.SetId("")
	return err
}

func resourceEndpointsServiceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	servicesService := servicemanagement.NewServicesService(config.clientServiceMan)
	service, err := servicesService.GetConfig(d.Get("service_name").(string)).Do()
	if err != nil {
		return err
	}
	d.Set("config_id", service.Id)
	d.Set("dns_address", service.Name)
	d.Set("apis", flattenServiceManagementApi(service.Apis))
	d.Set("endpoints", flattenServiceManagementEndpoints(service.Endpoints))

	return nil
}

func flattenServiceManagementApi(apis []*servicemanagement.Api) []map[string]interface{} {
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
