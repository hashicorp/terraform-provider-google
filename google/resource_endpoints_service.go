package google

import (
	"encoding/base64"
	"encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/servicemanagement/v1"
)

func resourceEndpointsService() *schema.Resource {
	return &schema.Resource{
		Create: resourceEndpointsServiceCreate,
		Read:   resourceEndpointsServiceRead,
		Delete: resourceEndpointsServiceDelete,
		Schema: map[string]*schema.Schema{
			"config_text": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"config_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"service_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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

func resourceEndpointsServiceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	configfile := servicemanagement.ConfigFile{
		FileContents: base64.StdEncoding.EncodeToString([]byte(d.Get("config_text").(string))),
		FileType:     "OPEN_API_YAML",
		FilePath:     "heredoc.yaml",
	}
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	source := servicemanagement.ConfigSource{
		Files: []*servicemanagement.ConfigFile{&configfile},
	}
	serviceName := d.Get("service_name").(string)
	servicesService := servicemanagement.NewServicesService(config.clientServiceMan)
	_, err = servicesService.Get(serviceName).Do()
	if err != nil {
		_, err = servicesService.Create(&servicemanagement.ManagedService{ProducerProjectId: project, ServiceName: serviceName}).Do()
		if err != nil {
			return err
		}
	}
	configService := servicemanagement.NewServicesConfigsService(config.clientServiceMan)
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

	rolloutService := servicemanagement.NewServicesRolloutsService(config.clientServiceMan)
	d.Set("config_id", serviceConfig.ServiceConfig.Id)
	d.Set("dns_address", serviceConfig.ServiceConfig.Name)
	d.Set("apis", flattenServiceManagementApi(serviceConfig.ServiceConfig.Apis))
	d.Set("endpoints", flattenServiceManagementEndpoints(serviceConfig.ServiceConfig.Endpoints))
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
	d.SetId(serviceName)
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
