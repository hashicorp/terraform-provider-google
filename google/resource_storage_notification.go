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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the bucket.`,
			},

			"payload_format": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"JSON_API_V1", "NONE"}, false),
				Description:  `The desired content of the Payload. One of "JSON_API_V1" or "NONE".`,
			},

			"topic": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description:      `The Cloud PubSub topic to which this subscription publishes. Expects either the  topic name, assumed to belong to the default GCP provider project, or the project-level name,  i.e. projects/my-gcp-project/topics/my-topic or my-topic. If the project is not set in the provider, you will need to use the project-level name.`,
			},

			"custom_attributes": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: ` A set of key/value attribute pairs to attach to each Cloud PubSub message published for this notification subscription`,
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
				Description: `List of event type filters for this notification config. If not specified, Cloud Storage will send notifications for all event types. The valid types are: "OBJECT_FINALIZE", "OBJECT_METADATA_UPDATE", "OBJECT_DELETE", "OBJECT_ARCHIVE"`,
			},

			"object_name_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies a prefix path filter for this notification config. Cloud Storage will only send notifications for objects in this bucket whose names begin with the specified prefix.`,
			},

			"notification_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the created notification.`,
			},

			"self_link": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URI of the created resource.`,
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
	d.Set("notification_id", notificationID)
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
