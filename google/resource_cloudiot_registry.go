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
	mqttStateUnspecified  = "MQTT_STATE_UNSPECIFIED"
	mqttEnabled           = "MQTT_ENABLED"
	mqttDisabled          = "MQTT_DISABLED"
	httpStateUnspecified  = "HTTP_STATE_UNSPECIFIED"
	httpEnabled           = "HTTP_ENABLED"
	httpDisabled          = "HTTP_DISABLED"
	unspecifiedCertFormat = "UNSPECIFIED_PUBLIC_KEY_CERTIFICATE_FORMAT"
	x509CertificatePEM    = "X509_CERTIFICATE_PEM"
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
				Optional: true,
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
											[]string{unspecifiedCertFormat, x509CertificatePEM}, false),
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
	parent := fmt.Sprintf("projects/%s/locations/%s", project, region)

	deviceRegistry := createDeviceRegistry(d)
	deviceRegistry.Id = d.Get("name").(string)

	call := config.clientCloudIoT.Projects.Locations.Registries.Create(parent, deviceRegistry)
	res, err := call.Do()
	if err != nil {
		return err
	}
	d.SetId(res.Name)
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
		updateMask = append(updateMask, "event_notification_config")
		if v, ok := d.GetOk("event_notification_config"); ok {
			deviceRegistry.EventNotificationConfigs = make([]*cloudiot.EventNotificationConfig, 1, 1)
			deviceRegistry.EventNotificationConfigs[0] = buildEventNotificationConfig(v.(map[string]interface{}))
		}
	}
	if d.HasChange("state_notification_config") {
		hasChanged = true
		updateMask = append(updateMask, "state_notification_config")
		if v, ok := d.GetOk("state_notification_config"); ok {
			deviceRegistry.StateNotificationConfig = buildStateNotificationConfig(v.(map[string]interface{}))
		}
	}
	if d.HasChange("mqtt_config") {
		hasChanged = true
		updateMask = append(updateMask, "mqtt_config")
		if v, ok := d.GetOk("mqtt_config"); ok {
			deviceRegistry.MqttConfig = buildMqttConfig(v.(map[string]interface{}))
		}
	}
	if d.HasChange("http_config") {
		hasChanged = true
		updateMask = append(updateMask, "http_config")
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
	d.Set("event_notification_config", nil)
	d.Set("state_notification_config", nil)
	d.Set("mqtt_config", nil)
	d.Set("http_config", nil)
	d.Set("credentials", nil)

	if res.EventNotificationConfigs != nil && len(res.EventNotificationConfigs) > 0 {
		eventConfig := map[string]string{"pubsub_topic_name": res.EventNotificationConfigs[0].PubsubTopicName}
		d.Set("event_notification_config", eventConfig)
	}
	// If no config exist for state notification, mqtt or http config default values are omitted.
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
			credentials[i] = pubcert
		}
		d.Set("credentials", credentials)
	}
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
	return nil
}

func resourceCloudIoTRegistryStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	r, _ := regexp.Compile("projects/(.*)/locations/(.*)/registries/(.*)")
	if r.MatchString(d.Id()) == false {
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
