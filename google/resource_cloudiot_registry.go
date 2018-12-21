package google

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/cloudiot/v1"
)

const (
	mqttEnabled        = "MQTT_ENABLED"
	mqttDisabled       = "MQTT_DISABLED"
	httpEnabled        = "HTTP_ENABLED"
	httpDisabled       = "HTTP_DISABLED"
	x509CertificatePEM = "X509_CERTIFICATE_PEM"
)

func resourceCloudIoTRegistry() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudIoTRegistryCreate,
		Update: resourceCloudIoTRegistryUpdate,
		Read:   resourceCloudIoTRegistryRead,
		Delete: resourceCloudIoTRegistryDelete,

		Importer: &schema.ResourceImporter{
			State: resourceCloudIoTRegistryStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateCloudIoTID,
			},
			"project": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"event_notification_config": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pubsub_topic_name": &schema.Schema{
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: compareSelfLinkOrResourceName,
						},
					},
				},
			},
			"state_notification_config": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pubsub_topic_name": &schema.Schema{
							Type:             schema.TypeString,
							Required:         true,
							DiffSuppressFunc: compareSelfLinkOrResourceName,
						},
					},
				},
			},
			"mqtt_config": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mqtt_enabled_state": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								[]string{mqttEnabled, mqttDisabled}, false),
						},
					},
				},
			},
			"http_config": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http_enabled_state": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								[]string{httpEnabled, httpDisabled}, false),
						},
					},
				},
			},
			"credentials": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"public_key_certificate": &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"format": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice(
											[]string{x509CertificatePEM}, false),
									},
									"certificate": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildEventNotificationConfig(config map[string]interface{}) *cloudiot.EventNotificationConfig {
	if v, ok := config["pubsub_topic_name"]; ok {
		return &cloudiot.EventNotificationConfig{
			PubsubTopicName: v.(string),
		}
	}
	return nil
}

func buildStateNotificationConfig(config map[string]interface{}) *cloudiot.StateNotificationConfig {
	if v, ok := config["pubsub_topic_name"]; ok {
		return &cloudiot.StateNotificationConfig{
			PubsubTopicName: v.(string),
		}
	}
	return nil
}

func buildMqttConfig(config map[string]interface{}) *cloudiot.MqttConfig {
	if v, ok := config["mqtt_enabled_state"]; ok {
		return &cloudiot.MqttConfig{
			MqttEnabledState: v.(string),
		}
	}
	return nil
}

func buildHttpConfig(config map[string]interface{}) *cloudiot.HttpConfig {
	if v, ok := config["http_enabled_state"]; ok {
		return &cloudiot.HttpConfig{
			HttpEnabledState: v.(string),
		}
	}
	return nil
}

func buildPublicKeyCertificate(certificate map[string]interface{}) *cloudiot.PublicKeyCertificate {
	cert := &cloudiot.PublicKeyCertificate{
		Format:      certificate["format"].(string),
		Certificate: certificate["certificate"].(string),
	}
	return cert
}

func expandCredentials(credentials []interface{}) []*cloudiot.RegistryCredential {
	certificates := make([]*cloudiot.RegistryCredential, len(credentials))
	for i, raw := range credentials {
		cred := raw.(map[string]interface{})
		certificates[i] = &cloudiot.RegistryCredential{
			PublicKeyCertificate: buildPublicKeyCertificate(cred["public_key_certificate"].(map[string]interface{})),
		}
	}
	return certificates
}

func createDeviceRegistry(d *schema.ResourceData) *cloudiot.DeviceRegistry {
	deviceRegistry := &cloudiot.DeviceRegistry{}
	if v, ok := d.GetOk("event_notification_config"); ok {
		deviceRegistry.EventNotificationConfigs = make([]*cloudiot.EventNotificationConfig, 1, 1)
		deviceRegistry.EventNotificationConfigs[0] = buildEventNotificationConfig(v.(map[string]interface{}))
	}
	if v, ok := d.GetOk("state_notification_config"); ok {
		deviceRegistry.StateNotificationConfig = buildStateNotificationConfig(v.(map[string]interface{}))
	}
	if v, ok := d.GetOk("mqtt_config"); ok {
		deviceRegistry.MqttConfig = buildMqttConfig(v.(map[string]interface{}))
	}
	if v, ok := d.GetOk("http_config"); ok {
		deviceRegistry.HttpConfig = buildHttpConfig(v.(map[string]interface{}))
	}
	if v, ok := d.GetOk("credentials"); ok {
		deviceRegistry.Credentials = expandCredentials(v.([]interface{}))
	}
	return deviceRegistry
}

func resourceCloudIoTRegistryCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	region, err := getRegion(d, config)
	if err != nil {
		return err
	}
	deviceRegistry := createDeviceRegistry(d)
	deviceRegistry.Id = d.Get("name").(string)
	parent := fmt.Sprintf("projects/%s/locations/%s", project, region)
	registryId := fmt.Sprintf("%s/registries/%s", parent, deviceRegistry.Id)
	d.SetId(registryId)

	err = retryTime(func() error {
		_, err := config.clientCloudIoT.Projects.Locations.Registries.Create(parent, deviceRegistry).Do()
		return err
	}, 5)
	if err != nil {
		d.SetId("")
		return err
	}

	// If we infer project and region, they are never actually set so we set them here
	d.Set("project", project)
	d.Set("region", region)

	return resourceCloudIoTRegistryRead(d, meta)
}

func resourceCloudIoTRegistryUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	updateMask := make([]string, 0, 5)
	hasChanged := false
	deviceRegistry := &cloudiot.DeviceRegistry{}

	d.Partial(true)

	if d.HasChange("event_notification_config") {
		hasChanged = true
		updateMask = append(updateMask, "event_notification_configs")
		if v, ok := d.GetOk("event_notification_config"); ok {
			deviceRegistry.EventNotificationConfigs = make([]*cloudiot.EventNotificationConfig, 1, 1)
			deviceRegistry.EventNotificationConfigs[0] = buildEventNotificationConfig(v.(map[string]interface{}))
		}
	}
	if d.HasChange("state_notification_config") {
		hasChanged = true
		updateMask = append(updateMask, "state_notification_config.pubsub_topic_name")
		if v, ok := d.GetOk("state_notification_config"); ok {
			deviceRegistry.StateNotificationConfig = buildStateNotificationConfig(v.(map[string]interface{}))
		}
	}
	if d.HasChange("mqtt_config") {
		hasChanged = true
		updateMask = append(updateMask, "mqtt_config.mqtt_enabled_state")
		if v, ok := d.GetOk("mqtt_config"); ok {
			deviceRegistry.MqttConfig = buildMqttConfig(v.(map[string]interface{}))
		}
	}
	if d.HasChange("http_config") {
		hasChanged = true
		updateMask = append(updateMask, "http_config.http_enabled_state")
		if v, ok := d.GetOk("http_config"); ok {
			deviceRegistry.HttpConfig = buildHttpConfig(v.(map[string]interface{}))
		}
	}
	if d.HasChange("credentials") {
		hasChanged = true
		updateMask = append(updateMask, "credentials")
		if v, ok := d.GetOk("credentials"); ok {
			deviceRegistry.Credentials = expandCredentials(v.([]interface{}))
		}
	}
	if hasChanged {
		_, err := config.clientCloudIoT.Projects.Locations.Registries.Patch(d.Id(),
			deviceRegistry).UpdateMask(strings.Join(updateMask, ",")).Do()
		if err != nil {
			return fmt.Errorf("Error updating registry %s: %s", d.Get("name").(string), err)
		}
		for _, updateMaskItem := range updateMask {
			d.SetPartial(updateMaskItem)
		}
	}
	d.Partial(false)
	return resourceCloudIoTRegistryRead(d, meta)
}

func resourceCloudIoTRegistryRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	name := d.Id()
	res, err := config.clientCloudIoT.Projects.Locations.Registries.Get(name).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Registry %q", name))
	}

	d.Set("name", res.Id)

	if len(res.EventNotificationConfigs) > 0 {
		eventConfig := map[string]string{"pubsub_topic_name": res.EventNotificationConfigs[0].PubsubTopicName}
		d.Set("event_notification_config", eventConfig)
	} else {
		d.Set("event_notification_config", nil)
	}
	pubsubTopicName := res.StateNotificationConfig.PubsubTopicName
	if pubsubTopicName != "" {
		d.Set("state_notification_config",
			map[string]string{"pubsub_topic_name": pubsubTopicName})
	} else {
		d.Set("state_notification_config", nil)
	}

	d.Set("mqtt_config", map[string]string{"mqtt_enabled_state": res.MqttConfig.MqttEnabledState})
	d.Set("http_config", map[string]string{"http_enabled_state": res.HttpConfig.HttpEnabledState})

	credentials := make([]map[string]interface{}, len(res.Credentials))
	for i, item := range res.Credentials {
		pubcert := make(map[string]interface{})
		pubcert["format"] = item.PublicKeyCertificate.Format
		pubcert["certificate"] = item.PublicKeyCertificate.Certificate
		credentials[i] = make(map[string]interface{})
		credentials[i]["public_key_certificate"] = pubcert
	}
	d.Set("credentials", credentials)
	return nil
}

func resourceCloudIoTRegistryDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	name := d.Id()
	call := config.clientCloudIoT.Projects.Locations.Registries.Delete(name)
	_, err := call.Do()
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceCloudIoTRegistryStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	r, _ := regexp.Compile("^projects/(.+)/locations/(.+)/registries/(.+)$")
	if !r.MatchString(d.Id()) {
		return nil, fmt.Errorf("Invalid registry specifier. " +
			"Expecting: projects/{project}/locations/{region}/registries/{name}")
	}
	parms := r.FindAllStringSubmatch(d.Id(), -1)[0]
	project := parms[1]
	region := parms[2]
	name := parms[3]

	id := fmt.Sprintf("projects/%s/locations/%s/registries/%s", project, region, name)
	d.Set("project", project)
	d.Set("region", region)
	d.SetId(id)
	return []*schema.ResourceData{d}, nil
}
