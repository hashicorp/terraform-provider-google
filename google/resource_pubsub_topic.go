package google

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/pubsub/v1"
)

func resourcePubsubTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourcePubsubTopicCreate,
		Read:   resourcePubsubTopicRead,
		Delete: resourcePubsubTopicDelete,

		Importer: &schema.ResourceImporter{
			State: resourcePubsubTopicStateImporter,
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
				ForceNew: true,
			},
		},
	}
}

func resourcePubsubTopicCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	name := fmt.Sprintf("projects/%s/topics/%s", project, d.Get("name").(string))
	topic := &pubsub.Topic{}

	call := config.clientPubsub.Projects.Topics.Create(name, topic)
	res, err := call.Do()
	if err != nil {
		return err
	}

	d.SetId(res.Name)

	return nil
}

func resourcePubsubTopicRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := d.Id()
	call := config.clientPubsub.Projects.Topics.Get(name)
	res, err := call.Do()
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("Pubsub Topic %q", name))
	}

	d.Set("name", res.Name)

	return nil
}

func resourcePubsubTopicDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	name := d.Id()
	call := config.clientPubsub.Projects.Topics.Delete(name)
	_, err := call.Do()
	if err != nil {
		return err
	}

	return nil
}

func resourcePubsubTopicStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	id := fmt.Sprintf("projects/%s/topics/%s", project, d.Id())

	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
