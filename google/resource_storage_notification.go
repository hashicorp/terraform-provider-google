package google

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"google.golang.org/api/storage/v1"
)

func resourceStorageNotification() *schema.Resource {
	return &schema.Resource{
		Create: resourceStorageNotificationCreate,
		Read:   resourceStorageNotificationRead,
		Delete: resourceStorageNotificationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"payload_format": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"JSON_API_V1", "NONE"}, false),
			},

			"topic": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
			},

			"custom_attributes": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"event_types": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"OBJECT_FINALIZE", "OBJECT_METADATA_UPDATE", "OBJECT_DELETE", "OBJECT_ARCHIVE"},
						false),
				},
			},

			"object_name_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"self_link": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceStorageNotificationCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket := d.Get("bucket").(string)

	topicName := d.Get("topic").(string)
	computedTopicName := getComputedTopicName("", topicName)
	if computedTopicName != topicName {
		project, err := getProject(d, config)
		if err != nil {
			return err
		}
		computedTopicName = getComputedTopicName(project, topicName)
	}

	storageNotification := &storage.Notification{
		CustomAttributes: expandStringMap(d, "custom_attributes"),
		EventTypes:       convertStringSet(d.Get("event_types").(*schema.Set)),
		ObjectNamePrefix: d.Get("object_name_prefix").(string),
		PayloadFormat:    d.Get("payload_format").(string),
		Topic:            computedTopicName,
	}

	res, err := config.clientStorage.Notifications.Insert(bucket, storageNotification).Do()
	if err != nil {
		return fmt.Errorf("Error creating notification config for bucket %s: %v", bucket, err)
	}

	d.SetId(fmt.Sprintf("%s/notificationConfigs/%s", bucket, res.Id))

	return resourceStorageNotificationRead(d, meta)
}

func resourceStorageNotificationRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket, notificationID := resourceStorageNotificationParseID(d.Id())

	res, err := config.clientStorage.Notifications.Get(bucket, notificationID).Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Notification configuration %s for bucket %s", notificationID, bucket))
	}

	d.Set("bucket", bucket)
	d.Set("payload_format", res.PayloadFormat)
	d.Set("topic", res.Topic)
	d.Set("object_name_prefix", res.ObjectNamePrefix)
	d.Set("event_types", res.EventTypes)
	d.Set("self_link", res.SelfLink)
	d.Set("custom_attributes", res.CustomAttributes)

	return nil
}

func resourceStorageNotificationDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	bucket, notificationID := resourceStorageNotificationParseID(d.Id())

	err := config.clientStorage.Notifications.Delete(bucket, notificationID).Do()
	if err != nil {
		return fmt.Errorf("Error deleting notification configuration %s for bucket %s: %v", notificationID, bucket, err)
	}

	return nil
}

func resourceStorageNotificationParseID(id string) (string, string) {
	//bucket, NotificationID
	parts := strings.Split(id, "/")

	return parts[0], parts[2]
}
