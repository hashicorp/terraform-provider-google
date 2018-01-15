package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"google.golang.org/api/cloudiot/v1"
	"strings"
)

const (
	mqttStateUnspecified  = "MQTT_STATE_UNSPECIFIED"
	mqttEnabled           = "MQTT_ENABLED"
	mqttDisabled          = "MQTT_DISABLED"
	httpStateUnspecified  = "HTTP_STATE_UNSPECIFIED"
	httpEnabled           = "HTTP_ENABLED"
	httpDisabled          = "HTTP_DISABLED"
	unspecifiedCertFormat = "UNSPECIFIED_PUBLIC_KEY_CERTIFICATE_FORMAT"
	x509CertificatePEM    = "X509_CERTIFICATE_PEM"
)

func resourceCloudiotRegistry() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloudiotRegistryCreate,
		Update: resourceCloudiotRegistryUpdate,
		Read:   resourceCloudiotRegistryRead,
		Delete: resourceCloudiotRegistryDelete,

		Importer: &schema.ResourceImporter{
			State: resourceCloudiotRegistryStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: linkDiffSuppress,
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
			"event_notification_configs": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				MaxItems: 1,
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
				ForceNew: false,
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
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mqtt_enabled_state": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								[]string{mqttStateUnspecified, mqttEnabled, mqttDisabled}, false),
						},
					},
				},
			},
			"http_config": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http_enabled_state": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								[]string{httpStateUnspecified, httpEnabled, httpDisabled}, false),
						},
					},
				},
			},
			"credentials": &schema.Schema{
				// Removed original "public_key_certificate" wrapper object. Additional nesting caused
				// problems with schema parsing.
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"format": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								[]string{unspecifiedCertFormat, x509CertificatePEM}, false),
						},
						"certificate": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"x509_details": &schema.Schema{
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"issuer": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"subject": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"start_time": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"expiry_time": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"signature_algorithm": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"public_key_type": &schema.Schema{
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

func expandEventNotificationConfigs(configs []interface{}) []*cloudiot.EventNotificationConfig {
	eventConfigs := make([]*cloudiot.EventNotificationConfig, len(configs))
	for i, raw := range configs {
		data := raw.(map[string]interface{})

		eventConfigs[i] = &cloudiot.EventNotificationConfig{
			PubsubTopicName: data["pubsub_topic_name"].(string),
		}
	}
	return eventConfigs
}

func expandStateNotificationConfig(config map[string]interface{}) *cloudiot.StateNotificationConfig {
	if v, ok := config["pubsub_topic_name"]; ok {
		return &cloudiot.StateNotificationConfig{
			PubsubTopicName: v.(string),
		}
	}
	return nil
}

func expandMqttConfig(config map[string]interface{}) *cloudiot.MqttConfig {
	if v, ok := config["mqtt_enabled_state"]; ok {
		return &cloudiot.MqttConfig{
			MqttEnabledState: v.(string),
		}
	}
	return nil
}

func expandHttpConfig(config map[string]interface{}) *cloudiot.HttpConfig {
	if v, ok := config["http_enabled_state"]; ok {
		return &cloudiot.HttpConfig{
			HttpEnabledState: v.(string),
		}
	}
	return nil
}

func expandX509Details(x509Details map[string]interface{}) *cloudiot.X509CertificateDetails {
	return &cloudiot.X509CertificateDetails{
		Issuer:             x509Details["issuer"].(string),
		Subject:            x509Details["subject"].(string),
		StartTime:          x509Details["start_time"].(string),
		ExpiryTime:         x509Details["expiry_time"].(string),
		SignatureAlgorithm: x509Details["signature_algorithm"].(string),
		PublicKeyType:      x509Details["public_key_type"].(string),
	}
}

func expandPublicKeyCertificate(certificate map[string]interface{}) *cloudiot.PublicKeyCertificate {
	return &cloudiot.PublicKeyCertificate{
		Format:      certificate["format"].(string),
		Certificate: certificate["certificate"].(string),
		X509Details: expandX509Details(certificate["x509_details"].(map[string]interface{})),
	}
}

func expandCredentials(credentials []interface{}) []*cloudiot.RegistryCredential {
	certificates := make([]*cloudiot.RegistryCredential, len(credentials))
	for i, raw := range credentials {
		cred := raw.(map[string]interface{})

		certificates[i] = &cloudiot.RegistryCredential{
			PublicKeyCertificate: expandPublicKeyCertificate(cred), // ["public_key_certificate"]
		}
	}
	return certificates
}

func expandDeviceRegistry(d *schema.ResourceData) *cloudiot.DeviceRegistry {
	deviceRegistry := &cloudiot.DeviceRegistry{}
	if v, ok := d.GetOk("event_notification_configs"); ok {
		deviceRegistry.EventNotificationConfigs = expandEventNotificationConfigs(v.([]interface{}))
	}
	if v, ok := d.GetOk("state_notification_config"); ok {
		deviceRegistry.StateNotificationConfig = expandStateNotificationConfig(v.(map[string]interface{}))

	}
	if v, ok := d.GetOk("mqtt_config"); ok {
		deviceRegistry.MqttConfig = expandMqttConfig(v.(map[string]interface{}))

	}
	if v, ok := d.GetOk("http_config"); ok {
		deviceRegistry.HttpConfig = expandHttpConfig(v.(map[string]interface{}))

	}
	if v, ok := d.GetOk("credentials"); ok {
		deviceRegistry.Credentials = expandCredentials(v.([]interface{}))

	}
	return deviceRegistry
}

func resourceCloudiotRegistryCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return err
	}

	parent := fmt.Sprintf("projects/%s/locations/%s", project, region)

	deviceRegistry := expandDeviceRegistry(d)
	deviceRegistry.Id = d.Get("name").(string)

	call := config.clientCloudiot.Projects.Locations.Registries.Create(parent, deviceRegistry)
	res, err := call.Do()
	if err != nil {
		return err
	}

	d.SetId(res.Name)

	return resourceCloudiotRegistryRead(d, meta)
}

func resourceCloudiotRegistryUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	updateMask := make([]string, 0, 5)
	hasChanged := false
	deviceRegistry := &cloudiot.DeviceRegistry{}

	d.Partial(true)

	if d.HasChange("event_notification_configs") {
		hasChanged = true
		updateMask = append(updateMask, "event_notification_configs")
		if v, ok := d.GetOk("event_notification_configs"); ok {
			deviceRegistry.EventNotificationConfigs = expandEventNotificationConfigs(v.([]interface{}))
		}
	}
	if d.HasChange("state_notification_config") {
		hasChanged = true
		if v, ok := d.GetOk("state_notification_config"); ok {
			deviceRegistry.StateNotificationConfig = expandStateNotificationConfig(v.(map[string]interface{}))

		}
	}
	if d.HasChange("mqtt_config") {
		hasChanged = true
		if v, ok := d.GetOk("mqtt_config"); ok {
			deviceRegistry.MqttConfig = expandMqttConfig(v.(map[string]interface{}))

		}
	}
	if d.HasChange("http_config") {
		hasChanged = true
		if v, ok := d.GetOk("http_config"); ok {
			deviceRegistry.HttpConfig = expandHttpConfig(v.(map[string]interface{}))

		}
	}
	if d.HasChange("credentials") {
		hasChanged = true
		if v, ok := d.GetOk("credentials"); ok {
			deviceRegistry.Credentials = expandCredentials(v.([]interface{}))

		}
	}

	if hasChanged {
		_, err := config.clientCloudiot.Projects.Locations.Registries.Patch(d.Id(),
			deviceRegistry).UpdateMask(strings.Join(updateMask, ",")).Do()
		if err != nil {
			return fmt.Errorf("Error updating registry %s: %s", d.Get("name").(string), err)
		}
		for _, updateMaskItem := range updateMask {
			d.SetPartial(updateMaskItem)
		}
	}

	d.Partial(false)

	return resourceCloudiotRegistryRead(d, meta)
}

func resourceCloudiotRegistryRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	name := d.Id()
	res, err := config.clientCloudiot.Projects.Locations.Registries.Get(name).Do()

	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Registry %q", name))
	}

	d.Set("name", res.Id)

	if res.EventNotificationConfigs != nil {
		eventConfigs := make([]map[string]string, len(res.EventNotificationConfigs))
		for i, config := range res.EventNotificationConfigs {
			eventConfigs[i] = map[string]string{"pubsub_topic_name": config.PubsubTopicName}
		}
		d.Set("event_notification_configs", eventConfigs)
	}
	// to keep the data model lean only changes are pushed if not the default is
	// returned or a config exists for state notification, mqtt and http config.
	if res.StateNotificationConfig != nil {
		pubsubTopicName := res.StateNotificationConfig.PubsubTopicName
		_, hasStateConfig := d.GetOk("state_notification_config")
		if pubsubTopicName != "" || hasStateConfig {
			d.Set("state_notification_config",
				map[string]string{"pubsub_topic_name": pubsubTopicName})
		}
	}
	if res.MqttConfig != nil {
		mqttState := res.MqttConfig.MqttEnabledState
		_, hasMqttConfig := d.GetOk("mqtt_config")
		if mqttState != mqttEnabled || hasMqttConfig {
			d.Set("mqtt_config",
				map[string]string{"mqtt_enabled_state": mqttState})
		}
	}
	if res.HttpConfig != nil {
		httpState := res.HttpConfig.HttpEnabledState
		_, hasHttpConfig := d.GetOk("http_config")
		if httpState != httpEnabled || hasHttpConfig {
			d.Set("http_config",
				map[string]string{"http_enabled_state": httpState})
		}
	}

	if res.Credentials != nil {
		credentials := make([]map[string]interface{}, len(res.Credentials))
		for i, item := range res.Credentials {
			pubcert := make(map[string]interface{})

			pubcert["format"] = item.PublicKeyCertificate.Format
			pubcert["certificate"] = item.PublicKeyCertificate.Certificate

			x509details := item.PublicKeyCertificate.X509Details
			if x509details != nil {
				pubcert["x509_details"] = map[string]interface{}{
					"issuer":              x509details.Issuer,
					"subject":             x509details.Subject,
					"start_time":          x509details.StartTime,
					"expiry_time":         x509details.ExpiryTime,
					"signature_algorithm": x509details.SignatureAlgorithm,
					"public_key_type":     x509details.PublicKeyType,
				}
			}
			credentials[i] = pubcert
		}
		d.Set("credentials", credentials)
	}

	return nil
}

func resourceCloudiotRegistryDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := d.Id()
	call := config.clientCloudiot.Projects.Locations.Registries.Delete(name)
	_, err := call.Do()
	if err != nil {
		return err
	}

	return nil
}

func resourceCloudiotRegistryStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	region, err := getRegion(d, config)
	if err != nil {
		return nil, err
	}

	id := fmt.Sprintf("projects/%s/locations/%s/registries/%s", project, region, d.Id())

	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
